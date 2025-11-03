package kubernetes

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewClient(c client.Client, scheme *runtime.Scheme) *Client {
	return &Client{
		client: c,
		scheme: scheme,
	}
}

func (c *Client) Apply(ctx context.Context, manifests []unstructured.Unstructured, namespace string) error {
	for _, manifest := range manifests {
		// Set namespace if not specified
		if manifest.GetNamespace() == "" {
			manifest.SetNamespace(namespace)
		}

		// Create apply options for server-side apply
		applyOpts := []client.PatchOption{
			client.ForceOwnership,
			client.FieldOwner("flowcd-controller"),
		}

		//obj creation
		obj := &unstructured.Unstructured{}
		obj.SetGroupVersionKind(manifest.GroupVersionKind())
		obj.SetName(manifest.GetName())
		obj.SetNamespace(manifest.GetNamespace())
		obj.Object = manifest.Object

		// Apply the manifest using server-side apply
		if err := c.client.Patch(ctx, obj, client.Apply, applyOpts...); err != nil {
			return fmt.Errorf("failed to apply manifest %s/%s: %w",
				manifest.GetNamespace(),
				manifest.GetName(),
				err)
		}
	}

	return nil
}
