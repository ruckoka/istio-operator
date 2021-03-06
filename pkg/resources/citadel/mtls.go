/*
Copyright 2019 Banzai Cloud.

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

package citadel

import (
	"github.com/banzaicloud/istio-operator/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// mTLS returns a map to configure the default MeshPolicy
func (r *Reconciler) mTLS() map[string]interface{} {
	if r.Config.Spec.MTLS {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"mode": "PERMISSIVE",
	}
}

// meshPolicy returns an authentication policy to either enable or disable mutual TLS
// for all services (that have sidecar) in the mesh
// https://istio.io/docs/tasks/security/authn-policy/
func (r *Reconciler) meshPolicy() *k8sutil.DynamicObject {
	return &k8sutil.DynamicObject{
		Gvr: schema.GroupVersionResource{
			Group:    "authentication.istio.io",
			Version:  "v1alpha1",
			Resource: "meshpolicies",
		},
		Kind:      "MeshPolicy",
		Name:      "default",
		Namespace: r.Config.Namespace,
		Labels: map[string]string{
			"app": "istio-security",
		},
		Spec: map[string]interface{}{
			"peers": []map[string]interface{}{
				{
					"mtls": r.mTLS(),
				},
			},
		},
		Owner: r.Config,
	}
}

// destinationRuleDefaultMtls returns a destination rule to configure client side to use mutual TLS when talking to
// any service (host) in the mesh
func (r *Reconciler) destinationRuleDefaultMtls() *k8sutil.DynamicObject {
	return &k8sutil.DynamicObject{
		Gvr: schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1alpha3",
			Resource: "destinationrules",
		},
		Kind:      "DestinationRule",
		Name:      "default",
		Namespace: r.Config.Namespace,
		Labels: map[string]string{
			"app": "istio-security",
		},
		Spec: map[string]interface{}{
			"host": "*.local",
			"trafficPolicy": map[string]interface{}{
				"tls": map[string]interface{}{
					"mode": "ISTIO_MUTUAL",
				},
			},
		},
		Owner: r.Config,
	}
}

// destinationRuleApiServerMtls returns a destination rule to disable (m)TLS when talking to API server, as API server doesn't have sidecar
// User should add similar destination rules for other services that don't have sidecar
func (r *Reconciler) destinationRuleApiServerMtls() *k8sutil.DynamicObject {
	return &k8sutil.DynamicObject{
		Gvr: schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1alpha3",
			Resource: "destinationrules",
		},
		Kind:      "DestinationRule",
		Name:      "api-server",
		Namespace: r.Config.Namespace,
		Labels: map[string]string{
			"app": "istio-security",
		},
		Spec: map[string]interface{}{
			"host": "kubernetes.default.svc.cluster.local",
			"trafficPolicy": map[string]interface{}{
				"tls": map[string]interface{}{
					"mode": "DISABLE",
				},
			},
		},
		Owner: r.Config,
	}
}
