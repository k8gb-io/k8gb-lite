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

	"cloud.example.com/annotation-operator/controllers/logging"

	"cloud.example.com/annotation-operator/controllers/providers/metrics"

	"github.com/stretchr/testify/assert"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"cloud.example.com/annotation-operator/controllers/depresolver"
	"cloud.example.com/annotation-operator/controllers/mocks"
	"cloud.example.com/annotation-operator/controllers/reconciliation"
	"cloud.example.com/annotation-operator/controllers/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcile(t *testing.T) {
	// arrange
	const reconcileRequeue = 30 * time.Second

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name    string
		Request reconcile.Request
		Result  controllerruntime.Result
	}{
		{
			Name:    "Empty request",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "", Name: ""}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: reconcileRequeue},
		},
		{
			Name:    "Only namespace filled",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "demo", Name: ""}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: reconcileRequeue},
		},
		{
			Name:    "Only name filled",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "", Name: "ingress"}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: reconcileRequeue},
		},
		{
			Name:    "Non existing resource",
			Request: reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "non-existing", Name: "ingress"}},
			Result:  reconcile.Result{Requeue: false, RequeueAfter: reconcileRequeue},
		},
	}

	r := getMockedReconciler(ctrl)
	r.Config.ReconcileRequeueSeconds = int(reconcileRequeue.Seconds())
	r.Client.(*mocks.MockClient).EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))

	// act
	// assert
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := r.Reconcile(context.TODO(), test.Request)
			assert.Equal(t, test.Result, result)
			require.NoError(t, err)
		})
	}
}

func getMockedReconciler(ctrl *gomock.Controller) *AnnoReconciler {
	c := mocks.NewMockClient(ctrl)
	r := mocks.NewMockGslbResolver(ctrl)
	p := mocks.NewMockProvider(ctrl)
	tr := mocks.NewMockTracer(ctrl)
	trsp := mocks.NewMockSpan(ctrl)

	config := &depresolver.Config{
		ReconcileRequeueSeconds: 30,
	}
	reconciler := &AnnoReconciler{
		Scheme:           runtime.NewScheme(),
		Client:           c,
		DepResolver:      r,
		DNSProvider:      p,
		Config:           config,
		IngressMapper:    reconciliation.NewIngressMapper(c),
		Tracer:           tr,
		ReconcilerResult: utils.NewReconcileResultHandler(config.ReconcileRequeueSeconds),
		Log:              logging.Logger(),
		Metrics:          metrics.Metrics(),
	}
	trsp.EXPECT().End(gomock.Any()).Return().AnyTimes()
	tr.EXPECT().Start(gomock.Any(), gomock.Any()).Return(context.TODO(), trsp).AnyTimes()
	return reconciler
}
