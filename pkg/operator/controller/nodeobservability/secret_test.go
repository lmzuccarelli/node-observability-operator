/*
Copyright 2021.

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

package nodeobservabilitycontroller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1alpha1 "github.com/openshift/node-observability-operator/api/v1alpha1"
	"github.com/openshift/node-observability-operator/pkg/operator/controller/test"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestEnsureSecret(t *testing.T) {
	nodeObs := &operatorv1alpha1.NodeObservability{}
	makeSecret := func() *corev1.Secret {
		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:   nodeObs.Namespace,
				Name:        secretName,
				Annotations: annotationsForSecret(),
			},
			Type: corev1.SecretTypeServiceAccountToken,
		}
		return &secret

	}
	testCases := []struct {
		name            string
		existingObjects []runtime.Object
		expectedExist   bool
		expectedSecret  *corev1.Secret
		errExpected     bool
	}{
		{
			name:            "Does not exist",
			existingObjects: []runtime.Object{},
			expectedExist:   true,
			expectedSecret:  makeSecret(),
		},
		{
			name: "Exists",
			existingObjects: []runtime.Object{
				makeSecret(),
			},
			expectedExist:  true,
			expectedSecret: makeSecret(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithRuntimeObjects(tc.existingObjects...).Build()
			r := &NodeObservabilityReconciler{
				Client: cl,
				Scheme: test.Scheme,
				Log:    zap.New(zap.UseDevMode(true)),
			}
			nodeObs := &operatorv1alpha1.NodeObservability{}
			gotExist, _, err := r.ensureSecret(context.TODO(), nodeObs)
			if err != nil {
				if !tc.errExpected {
					t.Fatalf("unexpected error received: %v", err)
				}
				return
			}
			if tc.errExpected {
				t.Fatalf("Error expected but wasn't received")
			}
			if gotExist != tc.expectedExist {
				t.Errorf("expected service account's exist to be %t, got %t", tc.expectedExist, gotExist)
			}
		})
	}
}
