package mapper

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
	"fmt"
	"strconv"

	"cloud.example.com/annotation-operator/controllers/utils"

	"cloud.example.com/annotation-operator/controllers/providers/metrics"

	"cloud.example.com/annotation-operator/controllers/depresolver"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
)

const (
	AnnotationPrimaryGeoTag              = "k8gb.io/primary-geotag"
	AnnotationStrategy                   = "k8gb.io/strategy"
	AnnotationDNSTTLSeconds              = "k8gb.io/dns-ttl-seconds"
	AnnotationSplitBrainThresholdSeconds = "k8gb.io/splitbrain-threshold-seconds"
	AnnotationWeightJSON                 = "k8gb.io/weights"
	AnnotationStatus                     = "k8gb.io/status"
	Finalizer                            = "k8gb.io/finalizer"
)

type Spec struct {
	PrimaryGeoTag              string         `json:"primaryGeoTag"`
	Type                       string         `json:"strategy"`
	DNSTtlSeconds              int            `json:"dnsTTLSeconds"`
	SplitBrainThresholdSeconds int            `json:"splitBrainThresholdSeconds"`
	Weights                    map[string]int `json:"weights"`
}

func (s *Spec) String() string {
	return fmt.Sprintf("strategy: %s, geo: %s", s.Type, s.PrimaryGeoTag)
}

// Status defines the observed state of Gslb
type Status struct {
	// Associated Service status
	ServiceHealth map[string]metrics.HealthStatus `json:"serviceHealth"`
	// Current Healthy DNS record structure
	HealthyRecords map[string][]string `json:"healthyRecords"`
	// Cluster Geo Tag
	GeoTag string `json:"geoTag"`
	// Comma-separated list of hosts. Duplicating the value from range .spec.ingress.rules[*].host for printer column
	Hosts string `json:"hosts,omitempty"`
}

func (s Status) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("{%v}", err)
	}
	return string(b)
}

// LoopState wraps information about ingress. Ensures that Ingress entity can't be nil
// TODO: don't allow user to access Ingress directly. Minimize number of operations over ingress or spec, Consider to use Getters and Setters instead
type LoopState struct {
	Mapper
	Ingress        *netv1.Ingress
	Spec           Spec
	NamespacedName types.NamespacedName
	Status         Status
}

func fromIngress(ingress *netv1.Ingress, m Mapper) (rs *LoopState, err error) {
	rs = &LoopState{Mapper: m}
	rs.SetReference(rs)
	if ingress == nil {
		return rs, fmt.Errorf("nil *ingress")
	}
	rs.Status = Status{
		ServiceHealth:  map[string]metrics.HealthStatus{},
		HealthyRecords: map[string][]string{},
		GeoTag:         "",
		Hosts:          "",
	}
	rs.Ingress = ingress
	rs.Spec, err = rs.asSpec(ingress.GetAnnotations())
	rs.NamespacedName = types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name}
	return rs, err
}

func fromGatewayAPI(gw *netv1.Ingress, m Mapper) (rs *LoopState, err error) {
	panic("not implemented")
}

func (rs *LoopState) asSpec(annotations map[string]string) (result Spec, err error) {
	var supportedStrategies = []string{depresolver.GeoStrategy, depresolver.FailoverStrategy, depresolver.RoundRobinStrategy}
	toInt := func(k string, v string) (int, error) {
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return -1, fmt.Errorf("can't parse annotation value %s to int for key %s", v, k)
		}
		return intValue, nil
	}
	result = Spec{
		Type:                       "",
		PrimaryGeoTag:              "",
		DNSTtlSeconds:              30,
		SplitBrainThresholdSeconds: 300,
		Weights:                    nil,
	}

	if value, found := annotations[AnnotationStrategy]; found {
		if !utils.Contains(supportedStrategies, value) {
			return result, fmt.Errorf("unsupported value '%s' for %s", value, AnnotationStrategy)
		}
		result.Type = value
	} else {
		return result, nil
	}

	if value, found := annotations[AnnotationPrimaryGeoTag]; found {
		result.PrimaryGeoTag = value
	}
	if value, found := annotations[AnnotationSplitBrainThresholdSeconds]; found {
		if result.SplitBrainThresholdSeconds, err = toInt(AnnotationSplitBrainThresholdSeconds, value); err != nil {
			return result, err
		}
	}
	if value, found := annotations[AnnotationDNSTTLSeconds]; found {
		if result.DNSTtlSeconds, err = toInt(AnnotationDNSTTLSeconds, value); err != nil {
			return result, err
		}
	}

	if value, found := annotations[AnnotationWeightJSON]; found {
		// e.g: '{"eu":5,"us":10}'
		w := make(map[string]int, 0)
		err = json.Unmarshal([]byte(value), &w)
		if err != nil {
			return result, fmt.Errorf("parsing %s (%v)", AnnotationWeightJSON, err)
		}
		result.Weights = w
	}

	if result.Type == depresolver.FailoverStrategy {
		if len(result.PrimaryGeoTag) == 0 {
			return result, fmt.Errorf("%s strategy requires annotation %s", depresolver.FailoverStrategy, AnnotationPrimaryGeoTag)
		}
	}

	// TODO: expose depresolver validators and use here!
	return result, nil
}
