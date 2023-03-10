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
	"reflect"
	"testing"

	"github.com/k8gb-io/k8gb-light/controllers/depresolver"

	"github.com/stretchr/testify/assert"
)

func TestAnnotations(t *testing.T) {
	var ds = Spec{
		Type:                       "",
		PrimaryGeoTag:              "",
		DNSTtlSeconds:              30,
		SplitBrainThresholdSeconds: 300,
		Weights:                    nil,
	}
	var serr = fmt.Errorf("")
	var tests = []struct {
		name          string
		annotations   map[string]string
		expectedError error
		expectedSpec  Spec
	}{
		{name: "No Annotations", annotations: map[string]string{}, expectedError: nil, expectedSpec: ds},
		{name: "No K8gb Annotations", annotations: map[string]string{"X": "Y"}, expectedError: nil, expectedSpec: ds},
		{name: "Only K8gb No Strategy Annotation", annotations: map[string]string{AnnotationDNSTTLSeconds: "300"}, expectedError: nil,
			expectedSpec: ds}, {name: "Invalid Strategy Annotation", annotations: map[string]string{AnnotationStrategy: "300"},
			expectedError: serr, expectedSpec: ds},
		{name: "Valid Strategy Invalid TTL", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy,
			AnnotationDNSTTLSeconds: "300a"}, expectedError: serr, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy, DNSTtlSeconds: -1,
			Weights: nil, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds, PrimaryGeoTag: ""}},
		{name: "Valid Strategy Invalid SB", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy,
			AnnotationSplitBrainThresholdSeconds: "300a"}, expectedError: serr, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy,
			DNSTtlSeconds: ds.DNSTtlSeconds, Weights: nil, SplitBrainThresholdSeconds: -1, PrimaryGeoTag: ""}},
		{name: "RR", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy}, expectedError: nil,
			expectedSpec: Spec{Type: depresolver.RoundRobinStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: ds.PrimaryGeoTag, Weights: ds.Weights, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "FO", annotations: map[string]string{AnnotationStrategy: depresolver.FailoverStrategy, AnnotationPrimaryGeoTag: "us"}, expectedError: nil,
			expectedSpec: Spec{Type: depresolver.FailoverStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: "us", Weights: ds.Weights, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "FO No PrimaryGeoTag", annotations: map[string]string{AnnotationStrategy: depresolver.FailoverStrategy}, expectedError: serr,
			expectedSpec: Spec{Type: depresolver.FailoverStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: "", Weights: ds.Weights, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "WRR", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy, AnnotationWeightJSON: "eu:10,us:10"},
			expectedError: nil, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: ds.PrimaryGeoTag, Weights: map[string]int{"eu": 10, "us": 10}, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "WRR with spaces", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy,
			AnnotationWeightJSON: " eu: 10, us:10 "},
			expectedError: nil, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: ds.PrimaryGeoTag, Weights: map[string]int{"eu": 10, "us": 10}, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "WRR followed by comma", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy,
			AnnotationWeightJSON: "eu: 10, us:10,"},
			expectedError: nil, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: ds.PrimaryGeoTag, Weights: map[string]int{"eu": 10, "us": 10}, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "WRR repeatable values", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy,
			AnnotationWeightJSON: "eu:10,us:10,eu:3"},
			expectedError: nil, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: ds.PrimaryGeoTag, Weights: map[string]int{"eu": 3, "us": 10}, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "WRR invalid value", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy,
			AnnotationWeightJSON: "eu:10,us:aa"}, expectedError: serr, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy,
			SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds, PrimaryGeoTag: ds.PrimaryGeoTag, Weights: nil, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "WRR invalid format", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy,
			AnnotationWeightJSON: "eu 10,us:aa"}, expectedError: serr, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy,
			SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds, PrimaryGeoTag: ds.PrimaryGeoTag, Weights: nil, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "WRR Empty", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy, AnnotationWeightJSON: ""},
			expectedError: nil, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: ds.PrimaryGeoTag, Weights: map[string]int{}, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "GeoIP", annotations: map[string]string{AnnotationStrategy: depresolver.GeoStrategy}, expectedError: nil,
			expectedSpec: Spec{Type: depresolver.GeoStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				Weights: ds.Weights, DNSTtlSeconds: ds.DNSTtlSeconds}},
		{name: "RR With PrimaryGeoTag", annotations: map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy, AnnotationPrimaryGeoTag: "us"},
			expectedError: nil, expectedSpec: Spec{Type: depresolver.RoundRobinStrategy, SplitBrainThresholdSeconds: ds.SplitBrainThresholdSeconds,
				PrimaryGeoTag: "us", Weights: ds.Weights, DNSTtlSeconds: ds.DNSTtlSeconds}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spec, err := new(LoopState).asSpec(test.annotations)
			assert.Equal(t, test.expectedError != nil, err != nil)
			assert.True(t, reflect.DeepEqual(test.expectedSpec, spec))
		})
	}
}
