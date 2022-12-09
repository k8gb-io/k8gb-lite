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
	"context"
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"cloud.example.com/annotation-operator/controllers/logging"
	"cloud.example.com/annotation-operator/controllers/mapper"

	"github.com/stretchr/testify/assert"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"cloud.example.com/annotation-operator/controllers/depresolver"
	"cloud.example.com/annotation-operator/controllers/mocks"
	"cloud.example.com/annotation-operator/controllers/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileRequest(t *testing.T) {
	// arrange
	const (
		reconcileRequeue = 30 * time.Second
		ingressName      = "ing"
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name    string
		Request reconcile.Request
		Result  controllerruntime.Result
	}{
		{
			Name:    "Request Empty",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "", Name: ""}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: reconcileRequeue},
		},
		{
			Name:    "Request Namespace",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "demo", Name: ""}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: reconcileRequeue},
		},
		{
			Name:    "Request Name",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "", Name: ingressName}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: reconcileRequeue},
		},
		{
			Name:    "Ingress Error",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "error", Name: ingressName}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: reconcileRequeue},
		},
		{
			Name:    "Ingress NotFound",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "notFound", Name: ingressName}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: 0},
		},
		{
			Name:    "Ingress Found Without Annotation",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "foundNoAnnotation", Name: ingressName}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: 0},
		},
	}

	cl := fakeClient(ctrl, depresolver.Config{})
	cl.Config.ReconcileRequeueSeconds = int(reconcileRequeue.Seconds())
	cl.Client.(*mocks.MockClient).EXPECT().Get(gomock.Any(), types.NamespacedName{Namespace: "notFound", Name: ingressName}, gomock.Any()).
		Return(errors.NewNotFound(schema.GroupResource{}, ingressName)).AnyTimes()
	cl.Client.(*mocks.MockClient).EXPECT().Get(gomock.Any(), types.NamespacedName{Namespace: "error", Name: ingressName}, gomock.Any()).
		Return(fmt.Errorf("random error")).AnyTimes()
	cl.Client.(*mocks.MockClient).EXPECT().Get(gomock.Any(), types.NamespacedName{Namespace: "foundNoAnnotation", Name: ingressName}, gomock.Any()).
		Return(nil).AnyTimes()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// arrange
			result, err := cl.Reconcile(context.TODO(), test.Request)

			// assert
			assert.Equal(t, test.Result, result)
			require.NoError(t, err)
		})
	}
}

func TestFinalizerInReconciliation(t *testing.T) {
	// The test evaluates the state that returns the finalization in the reconciliation loop.
	// The test focuses not only on the state, but also on the correct call of the tracer.
	// test
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const (
		ingressName = "ing"
	)
	var ferr = fmt.Errorf("finalizer err")

	var tests = []struct {
		Name     string
		Result   reconcile.Result
		SetMocks func(*AnnoReconciler)
		HasError bool
	}{
		{
			Name:     "Finalizer Installed",
			Result:   reconcile.Result{Requeue: false, RequeueAfter: 0},
			HasError: false,
			SetMocks: func(r *AnnoReconciler) {

				span := mocks.NewMockSpan(ctrl)
				span.EXPECT().End(gomock.Any()).Times(2)
				r.Tracer = mocks.NewMockTracer(ctrl)
				r.Tracer.(*mocks.MockTracer).EXPECT().Start(gomock.Any(), gomock.Any()).Return(context.TODO(), span).Times(2)

				r.Mapper.(*mocks.MockMapper).EXPECT().Get(gomock.Any()).Return(nil, mapper.ResultExists, nil).Times(1)
				r.DNSProvider.(*mocks.MockProvider).EXPECT().RequireFinalizer().Return(true).Times(1)
				r.Mapper.(*mocks.MockMapper).EXPECT().TryInjectFinalizer(gomock.Any()).Return(mapper.ResultContinue, nil).Times(1)
				r.Mapper.(*mocks.MockMapper).EXPECT().TryRemoveFinalizer(gomock.Any(), gomock.Any()).Return(mapper.ResultFinalizerRemoved, nil).Times(1)
			},
		},
		{
			Name:     "Finalizer Error",
			Result:   reconcile.Result{Requeue: false, RequeueAfter: 0},
			HasError: true,
			SetMocks: func(r *AnnoReconciler) {
				span := mocks.NewMockSpan(ctrl)
				span.EXPECT().End(gomock.Any()).Times(1)

				fspan := mocks.NewMockSpan(ctrl)
				fspan.EXPECT().RecordError(gomock.Any()).Times(1)
				fspan.EXPECT().SetStatus(gomock.Any(), gomock.Any()).Times(1)
				fspan.EXPECT().End(gomock.Any()).Return().Times(1)

				r.Tracer = mocks.NewMockTracer(ctrl)
				r.Tracer.(*mocks.MockTracer).EXPECT().Start(gomock.Any(), gomock.Any()).Return(context.TODO(), span).Times(1)
				r.Tracer.(*mocks.MockTracer).EXPECT().Start(gomock.Any(), gomock.Any()).Return(context.TODO(), fspan).Times(1)
				r.Mapper.(*mocks.MockMapper).EXPECT().Get(gomock.Any()).
					Return(&mapper.LoopState{NamespacedName: types.NamespacedName{Namespace: "ns", Name: ingressName}}, mapper.ResultExists, nil).Times(1)
				r.Mapper.(*mocks.MockMapper).EXPECT().TryInjectFinalizer(gomock.Any()).Return(mapper.ResultError, ferr).Times(1)
				r.DNSProvider.(*mocks.MockProvider).EXPECT().RequireFinalizer().Return(true).Times(1)
			},
		},
		// TODO: cover the rest of states
	}

	// act
	// assert
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			m := fakeMapper(ctrl)
			test.SetMocks(m)
			result, err := m.Reconcile(context.TODO(), reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "exists", Name: "ing"}})
			assert.Equal(t, test.Result, result)
			assert.Equal(t, err != nil, test.HasError)
		})
	}
}

func TestHandleFinalizer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// arrange
	var ferr = fmt.Errorf("finalizer err")
	var tests = []struct {
		Name           string
		ExpectedResult mapper.Result
		SetMocks       func(*mocks.MockMapper)
	}{
		{
			Name:           "Inject Finalizer",
			ExpectedResult: mapper.ResultFinalizerInstalled,
			SetMocks: func(c *mocks.MockMapper) {
				c.EXPECT().TryInjectFinalizer(gomock.Any()).Return(mapper.ResultFinalizerInstalled, nil).Times(1)
			},
		},
		{
			Name:           "Inject Finalizer Error",
			ExpectedResult: mapper.ResultError,
			SetMocks: func(c *mocks.MockMapper) {
				c.EXPECT().TryInjectFinalizer(gomock.Any()).Return(mapper.ResultError, ferr).Times(1)
			},
		},
		{
			Name:           "Remove Finalizer",
			ExpectedResult: mapper.ResultFinalizerRemoved,
			SetMocks: func(c *mocks.MockMapper) {
				c.EXPECT().TryInjectFinalizer(gomock.Any()).Return(mapper.ResultContinue, nil).Times(1)
				c.EXPECT().TryRemoveFinalizer(gomock.Any(), gomock.Any()).Return(mapper.ResultFinalizerRemoved, nil).Times(1)
			},
		},
		{
			Name:           "Remove Finalizer Error",
			ExpectedResult: mapper.ResultError,
			SetMocks: func(c *mocks.MockMapper) {
				c.EXPECT().TryInjectFinalizer(gomock.Any()).Return(mapper.ResultContinue, nil).Times(1)
				c.EXPECT().TryRemoveFinalizer(gomock.Any(), gomock.Any()).Return(mapper.ResultError, ferr).Times(1)
			},
		},
		{
			Name:           "Finalizer Skipped",
			ExpectedResult: mapper.ResultContinue,
			SetMocks: func(c *mocks.MockMapper) {
				c.EXPECT().TryInjectFinalizer(gomock.Any()).Return(mapper.ResultContinue, nil).Times(1)
				c.EXPECT().TryRemoveFinalizer(gomock.Any(), gomock.Any()).Return(mapper.ResultContinue, ferr).Times(1)
			},
		},
	}

	// act
	// assert
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			r := fakeMapper(ctrl)
			test.SetMocks(r.Mapper.(*mocks.MockMapper))
			result, err := r.handleFinalizer(nil)
			assert.Equal(t, test.ExpectedResult, result)
			if result == mapper.ResultError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// if you want to mock client use fakeClient
func fakeClient(ctrl *gomock.Controller, config depresolver.Config) *AnnoReconciler {
	r := fakeMapper(ctrl)
	r.Mapper = mapper.NewIngressMapper(r.Client, &config)
	return r
}

// if you want to mock mapper use fakeMapper
func fakeMapper(ctrl *gomock.Controller) *AnnoReconciler {
	c := mocks.NewMockClient(ctrl)
	r := mocks.NewMockGslbResolver(ctrl)
	p := mocks.NewMockProvider(ctrl)
	defaultTracer := mocks.NewMockTracer(ctrl)
	defaultTracerSpan := mocks.NewMockSpan(ctrl)
	m := mocks.NewMockMapper(ctrl)
	defaultMetrics := mocks.NewMockMetrics(ctrl)

	config := &depresolver.Config{
		ReconcileRequeueSeconds: 30,
	}
	reconciler := &AnnoReconciler{
		Scheme:           runtime.NewScheme(),
		Client:           c,
		DepResolver:      r,
		DNSProvider:      p,
		Config:           config,
		Mapper:           m,
		Tracer:           defaultTracer,
		ReconcilerResult: utils.NewReconcileResultHandler(config.ReconcileRequeueSeconds),
		Log:              logging.Logger(),
		Metrics:          defaultMetrics,
	}
	// providing default tracer and span
	defaultTracerSpan.EXPECT().End(gomock.Any()).Return().AnyTimes()
	defaultTracer.EXPECT().Start(gomock.Any(), gomock.Any()).Return(context.TODO(), defaultTracerSpan).AnyTimes()

	// providing default metrics
	defaultMetrics.EXPECT().IncrementError(gomock.Any()).AnyTimes()
	return reconciler
}
