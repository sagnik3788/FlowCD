package kubernetes

import (
	"bytes"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func ParseManifests(data []byte) ([]unstructured.Unstructured, error) {
	var objects []unstructured.Unstructured

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 4096)
	for {
		obj := &unstructured.Unstructured{}
		if err := decoder.Decode(obj); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if obj.GetKind() != "" {
			objects = append(objects, *obj)
		}
	}

	return objects, nil
}
