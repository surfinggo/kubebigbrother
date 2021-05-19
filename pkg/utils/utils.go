package utils

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func NamespaceKey(s *unstructured.Unstructured) string {
	if len(s.GetNamespace()) > 0 {
		return s.GetNamespace() + "/" + s.GetName()
	}
	return s.GetName()
}
