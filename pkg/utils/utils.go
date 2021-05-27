package utils

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// GroupVersionKindName returns group, version, kind, namespace and name string
// for an unstructured resource
//
// examples:
//   /v1, Kind=ConfigMap, demo/demo
//   apps/v1, Kind=Deployment, demo/canary
func GroupVersionKindName(s *unstructured.Unstructured) string {
	return s.GroupVersionKind().String() + ", " + NamespaceKey(s)
}

// NamespaceKey returns namespaced key for an unstructured resource
func NamespaceKey(s *unstructured.Unstructured) string {
	if len(s.GetNamespace()) > 0 {
		return s.GetNamespace() + "/" + s.GetName()
	}
	return s.GetName()
}
