package controllers

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"cloud.example.com/annotation-operator/controllers/status"
	"context"
	"regexp"
	"strings"

	"cloud.example.com/annotation-operator/controllers/rs"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

func (r *AnnoReconciler) updateStatus(rs *rs.ReconciliationState, ep *externaldns.DNSEndpoint) (err error) {
	rs.Status.ServiceHealth, err = r.getServiceHealthStatus(rs)
	if err != nil {
		return err
	}

	m.UpdateIngressHostsPerStatusMetric(rs, rs.Status.ServiceHealth)

	rs.Status.HealthyRecords, err = r.getHealthyRecords(rs)
	if err != nil {
		return err
	}

	rs.Status.GeoTag = r.Config.ClusterGeoTag
	rs.Status.Hosts = r.hostsToCSV(rs)

	m.UpdateHealthyRecordsMetric(rs, rs.Status.HealthyRecords)

	m.UpdateEndpointStatus(ep)

	rs.Ingress.Annotations["k8gb.io/status"] = rs.Status.String()
	err = r.Update(context.TODO(), rs.Ingress)
	return err
}

func (r *AnnoReconciler) getServiceHealthStatus(rs *rs.ReconciliationState) (map[string]status.HealthStatus, error) {
	serviceHealth := make(map[string]status.HealthStatus)
	for _, rule := range rs.Ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service == nil || path.Backend.Service.Name == "" {
				log.Warn().
					Str("gslb", rs.NamespacedName.Name).
					Interface("service", path.Backend.Service).
					Msg("Malformed service definition")
				serviceHealth[rule.Host] = status.NotFound
				continue
			}
			service := &corev1.Service{}
			finder := client.ObjectKey{
				Namespace: rs.NamespacedName.Namespace,
				Name:      path.Backend.Service.Name,
			}
			err := r.Get(context.TODO(), finder, service)
			if err != nil {
				if errors.IsNotFound(err) {
					serviceHealth[rule.Host] = status.NotFound
					continue
				}
				return serviceHealth, err
			}

			endpoints := &corev1.Endpoints{}

			nn := types.NamespacedName{
				Name:      path.Backend.Service.Name,
				Namespace: rs.NamespacedName.Namespace,
			}

			err = r.Get(context.TODO(), nn, endpoints)
			if err != nil {
				return serviceHealth, err
			}

			serviceHealth[rule.Host] = status.Unhealthy
			if len(endpoints.Subsets) > 0 {
				for _, subset := range endpoints.Subsets {
					if len(subset.Addresses) > 0 {
						serviceHealth[rule.Host] = status.Healthy
					}
				}
			}
		}
	}
	return serviceHealth, nil
}

func (r *AnnoReconciler) getHealthyRecords(rs *rs.ReconciliationState) (map[string][]string, error) {

	dnsEndpoint := &externaldns.DNSEndpoint{}

	err := r.Get(context.TODO(), rs.NamespacedName, dnsEndpoint)
	if err != nil {
		return nil, err
	}

	healthyRecords := make(map[string][]string)

	serviceRegex := regexp.MustCompile("^localtargets")
	for _, endpoint := range dnsEndpoint.Spec.Endpoints {
		local := serviceRegex.Match([]byte(endpoint.DNSName))
		if !local && endpoint.RecordType == "A" {
			if len(endpoint.Targets) > 0 {
				healthyRecords[endpoint.DNSName] = endpoint.Targets
			}
		}
	}

	return healthyRecords, nil
}

func (r *AnnoReconciler) hostsToCSV(rs *rs.ReconciliationState) string {
	var hosts []string
	for _, r := range rs.Ingress.Spec.Rules {
		hosts = append(hosts, r.Host)
	}
	return strings.Join(hosts, ", ")
}