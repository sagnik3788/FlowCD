package k8s

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// apply manifests on our local cluster
func ApplyManifest(path string, k8sClient client.Client, scheme *runtime.Scheme) error {
	files, err := filepath.Glob(filepath.Join(path, ".yaml"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		decoder := utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader(content), 100)

		for {
			obj := &unstructured.Unstructured{} //to hold decoded yaml
			err = decoder.Decode(obj)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			obj.SetNamespace("default") // for now in manifests dont have ns
			//apply on local cluster
			err = k8sClient.Patch(context.TODO(), obj, client.Apply, client.ForceOwnership, client.FieldOwner("gitops-controller"))
			if err != nil {
				return err
			}
		}

	}
	return nil
}
