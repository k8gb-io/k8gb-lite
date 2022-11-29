/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"os"

	externaldns "sigs.k8s.io/external-dns/endpoint"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	sch "sigs.k8s.io/controller-runtime/pkg/scheme"

	"cloud.example.com/annotation-operator/controllers/depresolver"
	"cloud.example.com/annotation-operator/controllers/logging"
	"cloud.example.com/annotation-operator/controllers/providers/dns"
	"cloud.example.com/annotation-operator/controllers/providers/metrics"
	"cloud.example.com/annotation-operator/controllers/tracing"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"cloud.example.com/annotation-operator/controllers"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
	version  = "development"
	commit   = "none"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	resolver := depresolver.NewDependencyResolver()
	config, err := resolver.ResolveOperatorConfig()
	deprecations := resolver.GetDeprecations()
	// Initialize desired log or default log in case of configuration failed.
	logging.Init(config)
	log := logging.Logger()
	log.Info().
		Str("version", version).
		Str("commit", commit).
		Msg("k8gb info")
	if err != nil {
		log.Err(err).Msg("Can't resolve environment variables")
		return err
	}
	log.Debug().
		Interface("config", config).
		Msg("Resolved config")

	ctrl.SetLogger(logging.NewLogrAdapter(log))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: config.MetricsAddress,
		Port:               9443,
		LeaderElection:     false,
		LeaderElectionID:   "8020e9ff.absa.oss",
	})
	if err != nil {
		log.Err(err).Msg("Unable to create k8gb operator manager")
		return err
	}

	for _, d := range deprecations {
		log.Warn().Msg(d)
	}

	log.Info().Msg("Registering components")

	// Add external-dns DNSEndpoints resource
	// https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#adding-3rd-party-resources-to-your-operator
	schemeBuilder := &sch.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
	schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
	if err = schemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
		log.Err(err).Msg("Unable to register ExternalDNS resource schemas")
		return err
	}

	reconciler := &controllers.AnnoReconciler{
		Config:      config,
		Client:      mgr.GetClient(),
		DepResolver: resolver,
		Scheme:      mgr.GetScheme(),
	}

	log.Info().Msg("Starting metrics")
	metrics.Init(config)
	defer metrics.Metrics().Unregister()
	err = metrics.Metrics().Register()
	if err != nil {
		log.Err(err).Msg("Unable to register metrics")
		return err
	}

	log.Info().Msg("Resolving DNS provider")
	var f *dns.ProviderFactory
	f, err = dns.NewDNSProviderFactory(reconciler.Client, *reconciler.Config)
	if err != nil {
		log.Err(err).Msg("Unable to create DNS provider factory")
		return err
	}
	reconciler.DNSProvider = f.Provider()
	log.Info().
		Str("provider", reconciler.DNSProvider.String()).
		Msg("Started DNS provider")

	if err = reconciler.SetupWithManager(mgr); err != nil {
		log.Err(err).Msg("Unable to create Gslb controller")
		return err
	}
	metrics.Metrics().SetRuntimeInfo(version, commit)

	// tracing
	cfg := tracing.Settings{
		Enabled:       config.TracingEnabled,
		Endpoint:      config.OtelExporterOtlpEndpoint,
		SamplingRatio: config.TracingSamplingRatio,
		Commit:        commit,
		AppVersion:    version,
	}
	cleanup, tracer := tracing.SetupTracing(context.Background(), cfg, log)
	reconciler.Tracer = tracer
	defer cleanup()

	// +kubebuilder:scaffold:builder
	log.Info().Msg("Starting k8gb")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Err(err).Msg("Problem running k8gb")
		return err
	}
	log.Info().Msg("Gracefully finished, bye!")
	return nil
}