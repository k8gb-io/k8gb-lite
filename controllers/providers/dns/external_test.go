package dns

//
///*
//Copyright 2022 The k8gb Contributors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
//Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
//*/
//
//import (
//	"cloud.example.com/annotation-operator/controllers/depresolver"
//	"cloud.example.com/annotation-operator/controllers/providers/assistant"
//	"cloud.example.com/annotation-operator/controllers/rs"
//	"cloud.example.com/annotation-operator/controllers/utils"
//	"fmt"
//	"os"
//	"reflect"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//	corev1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/runtime/schema"
//
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	externaldns "sigs.k8s.io/external-dns/endpoint"
//
//	"k8s.io/apimachinery/pkg/runtime"
//	"sigs.k8s.io/controller-runtime/pkg/client/fake"
//	"sigs.k8s.io/controller-runtime/pkg/scheme"
//)
//
//// test data
//var a = struct {
//	Config              depresolver.Config
//	Gslb                *rs.ReconciliationState
//	TargetIPs           []string
//	TargetNSNamesSorted []string
//}{
//	Config: depresolver.Config{
//		ReconcileRequeueSeconds: 30,
//		ClusterGeoTag:           "us",
//		ExtClustersGeoTags:      []string{"za", "eu"},
//		EdgeDNSServers: []utils.DNSServer{
//			{
//				Host: "dns.cloud.example.com",
//				Port: 53,
//			},
//		},
//		EdgeDNSZone:   "example.com",
//		DNSZone:       "cloud.example.com",
//		K8gbNamespace: "k8gb",
//	},
//	Gslb: func() *rs.ReconciliationState {
//		var crSampleYaml = "../../../deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr.yaml"
//		gslbYaml, _ := os.ReadFile(crSampleYaml)
//		gslb, _ := utils.YamlToGslb(gslbYaml)
//		return gslb
//	}(),
//	TargetIPs: []string{
//		"10.0.1.38",
//		"10.0.1.40",
//		"10.0.1.39",
//	},
//	TargetNSNamesSorted: []string{
//		"gslb-ns-eu-cloud.example.com",
//		"gslb-ns-us-cloud.example.com",
//		"gslb-ns-za-cloud.example.com",
//	},
//}
//
//var expectedDNSEndpoint = &externaldns.DNSEndpoint{
//	ObjectMeta: metav1.ObjectMeta{
//		Name:        fmt.Sprintf("k8gb-ns-%s", externalDNSTypeCommon),
//		Namespace:   a.Config.K8gbNamespace,
//		Annotations: map[string]string{"k8gb.absa.oss/dnstype": string(externalDNSTypeCommon)},
//	},
//	Spec: externaldns.DNSEndpointSpec{
//		Endpoints: []*externaldns.Endpoint{
//			{
//				DNSName:    a.Config.DNSZone,
//				RecordTTL:  30,
//				RecordType: "NS",
//				Targets:    a.TargetNSNamesSorted,
//			},
//			{
//				DNSName:    "gslb-ns-us-cloud.example.com",
//				RecordTTL:  30,
//				RecordType: "A",
//				Targets:    a.TargetIPs,
//			},
//		},
//	},
//}
//
//func TestCreateZoneDelegationOnExternalDNS(t *testing.T) {
//	// arrange
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	m := mocks.NewMockAssistant(ctrl)
//	p := NewExternalDNS(a.Config, m)
//	m.EXPECT().GslbIngressExposedIPs(a.Gslb).Return(a.TargetIPs, nil).Times(1)
//	m.EXPECT().SaveDNSEndpoint(a.Config.K8gbNamespace, gomock.Eq(expectedDNSEndpoint)).Return(nil).Times(1).
//		Do(func(ns string, ep *externaldns.DNSEndpoint) {
//			require.True(t, reflect.DeepEqual(ep, expectedDNSEndpoint))
//			require.Equal(t, ns, a.Config.K8gbNamespace)
//		})
//
//	// act
//	err := p.CreateZoneDelegationForExternalDNS(a.Gslb)
//	// assert
//	assert.NoError(t, err)
//}
//
//func TestSaveNewDNSEndpointOnExternalDNS(t *testing.T) {
//	// arrange
//	var ep = &corev1.Endpoints{
//		TypeMeta: metav1.TypeMeta{
//			APIVersion: "v1",
//			Kind:       "Endpoints",
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "k8gb-ns-extdns",
//			Namespace: "test-gslb",
//		},
//	}
//	endpointToSave := expectedDNSEndpoint
//	endpointToSave.Namespace = a.Gslb.Namespace
//
//	runtimeScheme := runtime.NewScheme()
//	schemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
//	schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
//	require.NoError(t, corev1.AddToScheme(runtimeScheme))
//	require.NoError(t, k8gbv1beta1.AddToScheme(runtimeScheme))
//	require.NoError(t, schemeBuilder.AddToScheme(runtimeScheme))
//
//	var cl = fake.NewClientBuilder().WithScheme(runtimeScheme).WithObjects(ep).Build()
//
//	assistant := assistant.NewGslbAssistant(cl, a.Config.K8gbNamespace, a.Config.EdgeDNSServers)
//	p := NewExternalDNS(a.Config, assistant)
//	// act, assert
//	err := p.SaveDNSEndpoint(a.Gslb, expectedDNSEndpoint)
//	assert.NoError(t, err)
//}
//
//func TestSaveExistingDNSEndpointOnExternalDNS(t *testing.T) {
//	// arrange
//	endpointToSave := expectedDNSEndpoint
//	endpointToSave.Namespace = a.Gslb.Namespace
//
//	runtimeScheme := runtime.NewScheme()
//	schemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
//	schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
//	require.NoError(t, corev1.AddToScheme(runtimeScheme))
//	require.NoError(t, k8gbv1beta1.AddToScheme(runtimeScheme))
//	require.NoError(t, schemeBuilder.AddToScheme(runtimeScheme))
//
//	var cl = fake.NewClientBuilder().WithScheme(runtimeScheme).WithObjects(endpointToSave).Build()
//	assistant := assistant.NewGslbAssistant(cl, a.Config.K8gbNamespace, a.Config.EdgeDNSServers)
//	p := NewExternalDNS(a.Config, assistant)
//	// act, assert
//	err := p.SaveDNSEndpoint(a.Gslb, endpointToSave)
//	assert.NoError(t, err)
//}