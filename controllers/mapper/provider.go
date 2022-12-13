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

	"cloud.example.com/annotation-operator/controllers/depresolver"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ProviderMapper interface {
	Get(types.NamespacedName) (*LoopState, Result, error)
	FromIngress(*netv1.Ingress) (*LoopState, error)
	FromGatewayAPI() (*LoopState, error)
}

type CommonProvider struct {
	c      client.Client
	config *depresolver.Config
}

func NewCommonProvider(c client.Client, config *depresolver.Config) *CommonProvider {
	return &CommonProvider{
		c:      c,
		config: config,
	}
}

func (c *CommonProvider) Get(selector types.NamespacedName) (rs *LoopState, result Result, err error) {
	// TODO: implement gateway part of Get. Only ingress is implemented
	// e.g: You can try read GW first, if not success than Ingress
	var ing = &netv1.Ingress{}
	err = c.c.Get(context.TODO(), selector, ing)
	result, err = c.getConverterResult(err, ing)
	if result == ResultError {
		return nil, result, err
	}
	rs, err = c.FromIngress(ing)
	if err != nil {
		result = ResultError
	}
	return rs, result, err
}

// FromIngress LoopState from Ingress instance
func (c *CommonProvider) FromIngress(ingress *netv1.Ingress) (*LoopState, error) {
	// TODO: check here
	m := NewIngressMapper(c.c, c.config)
	return fromIngress(ingress, m)
}

func (c *CommonProvider) FromGatewayAPI() (*LoopState, error) {
	m := NewGatewayAPIMapper(c.c, c.config)
	return fromGatewayAPI(nil, m)
}

func (c *CommonProvider) getConverterResult(err error, ing *netv1.Ingress) (Result, error) {
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
