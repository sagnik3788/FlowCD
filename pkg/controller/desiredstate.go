package controller

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func GetDesiredState(manifests []unstructured.Unstructured, defaultNamespace string) ([]ResourceState, error) {
	var resources []ResourceState

	for _, manifest := range manifests {
		kind := manifest.GetKind()
		name := manifest.GetName()
		namespace := manifest.GetNamespace()

		if namespace == "" {
			namespace = defaultNamespace
		}

		spec, found := manifest.Object["spec"]
		if !found {
			continue
		}

		resources = append(resources, ResourceState{
			Kind:      kind,
			Name:      name,
			Namespace: namespace,
			Spec:      spec,
		})
	}

	return resources, nil
}
