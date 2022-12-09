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
	"context"
	"fmt"
	"reflect"

	"cloud.example.com/annotation-operator/controllers/depresolver"

	"cloud.example.com/annotation-operator/controllers/providers/metrics"
	corev1 "k8s.io/api/core/v1"

	"cloud.example.com/annotation-operator/controllers/utils"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IngressMapper provides API for working with ingress
type IngressMapper struct {
	c      client.Client
	config *depresolver.Config
}

func NewIngressMapper(c client.Client, config *depresolver.Config) *IngressMapper {
	return &IngressMapper{
		c:      c,
		config: config,
	}
}

func (i *IngressMapper) UpdateStatus(state *LoopState) (err error) {
	// check if object has not been deleted
	var r Result
	var s *LoopState
	s, r, err = i.Get(state.NamespacedName)
	switch r {
	case ResultError:
		return err
	case ResultNotFound:
		// object was deleted
		return nil
	}
	// don't do update if nothing has changed
	if s.Ingress.Annotations[AnnotationStatus] == state.Status.String() {
		return nil
	}
	// update the planned object
	s.Ingress.Annotations[AnnotationStatus] = state.Status.String()
	return i.c.Update(context.TODO(), s.Ingress)
}

func (i *IngressMapper) Get(selector types.NamespacedName) (rs *LoopState, result Result, err error) {
	var ing = &netv1.Ingress{}
	err = i.c.Get(context.TODO(), selector, ing)
	result, err = i.getConverterResult(err, ing)
	if result == ResultError {
		return nil, result, err
	}
	rs, err = NewLoopState(ing)
	if err != nil {
		result = ResultError
	}
	return rs, result, err
}

// Equal compares given ingress annotations and Ingres.Spec. If any of ingresses doesn't exist, returns false
func (i *IngressMapper) Equal(rs1 *LoopState, rs2 *LoopState) bool {
	if rs1 == nil || rs2 == nil {
		return false
	}
	if !reflect.DeepEqual(rs1.Spec, rs2.Spec) {
		return false
	}
	if !reflect.DeepEqual(rs1.Ingress.Spec, rs2.Ingress.Spec) {
		return false
	}
	return true
}

func (i *IngressMapper) TryInjectFinalizer(rs *LoopState) (Result, error) {
	if rs == nil || rs.Ingress == nil {
		return ResultError, fmt.Errorf("injecting finalizer from nil values")
	}
	if !utils.Contains(rs.Ingress.GetFinalizers(), Finalizer) {
		rs.Ingress.SetFinalizers(append(rs.Ingress.GetFinalizers(), Finalizer))
		err := i.c.Update(context.TODO(), rs.Ingress)
		if err != nil {
			return ResultError, err
		}
		return ResultFinalizerInstalled, nil
	}
	return ResultContinue, nil
}

func (i *IngressMapper) TryRemoveFinalizer(rs *LoopState, finalize func(*LoopState) error) (Result, error) {
	if rs == nil || rs.Ingress == nil {
		return ResultError, fmt.Errorf("removing finalizer from nil values")
	}
	if utils.Contains(rs.Ingress.GetFinalizers(), Finalizer) {
		isMarkedToBeDeleted := rs.Ingress.GetDeletionTimestamp() != nil
		if !isMarkedToBeDeleted {
			return ResultContinue, nil
		}
		err := finalize(rs)
		if err != nil {
			return ResultError, err
		}
		rs.Ingress.SetFinalizers(utils.Remove(rs.Ingress.GetFinalizers(), Finalizer))
		err = i.c.Update(context.TODO(), rs.Ingress)
		if err != nil {
			return ResultError, err
		}
		return ResultFinalizerRemoved, nil
	}
	return ResultContinue, nil
}

func (i *IngressMapper) GetHealthStatus(rs *LoopState) (map[string]metrics.HealthStatus, error) {
	serviceHealth := make(map[string]metrics.HealthStatus)
	for _, rule := range rs.Ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service == nil || path.Backend.Service.Name == "" {
				serviceHealth[rule.Host] = metrics.NotFound
				continue
			}

			// check if service exists
			selector := types.NamespacedName{Namespace: rs.NamespacedName.Namespace, Name: path.Backend.Service.Name}
			service := &corev1.Service{}
			err := i.c.Get(context.TODO(), selector, service)
			if err != nil {
				if errors.IsNotFound(err) {
					serviceHealth[rule.Host] = metrics.NotFound
					continue
				}
				return serviceHealth, err
			}

			// check if service endpoint exists
			ep := &corev1.Endpoints{}
			err = i.c.Get(context.TODO(), selector, ep)
			if err != nil {
				return serviceHealth, err
			}
			serviceHealth[rule.Host] = metrics.Unhealthy
			for _, subset := range ep.Subsets {
				if len(subset.Addresses) > 0 {
					serviceHealth[rule.Host] = metrics.Healthy
				}
			}
		}
	}
	return serviceHealth, nil
}

func (i *IngressMapper) GetExposedIPs(rs *LoopState) ([]string, error) {
	var exposed []string
	for _, ing := range rs.Ingress.Status.LoadBalancer.Ingress {
		if len(ing.IP) > 0 {
			exposed = append(exposed, ing.IP)
		}
		if len(ing.Hostname) > 0 {
			ips, err := utils.Dig(ing.Hostname, i.config.EdgeDNSServers...)
			if err != nil {
				return nil, err
			}
			exposed = append(exposed, ips...)
		}
	}
	return exposed, nil
}

func (i *IngressMapper) getConverterResult(err error, ing *netv1.Ingress) (Result, error) {
	if err != nil && errors.IsNotFound(err) {
		return ResultNotFound, nil
	} else if err != nil {
		return ResultError, err
	}
	if _, found := ing.GetAnnotations()[AnnotationStrategy]; !found {
		return ResultExistsButNotAnnotationFound, nil
	}
	return ResultExists, nil
}
