package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceState struct {
	Kind      string
	Name      string
	Namespace string
	Spec      any
}

func ListResources(ctx context.Context, k8sClient client.Client, namespace string) ([]ResourceState, error) {
	var resources []ResourceState

	deploymentList := &appsv1.DeploymentList{}
	serviceList := &corev1.ServiceList{}

	listOpts := []client.ListOption{}
	if namespace != "" {
		listOpts = append(listOpts, client.InNamespace(namespace))
	}

	err := k8sClient.List(ctx, deploymentList, listOpts...)
	if err != nil {
		return nil, err
	}

	err = k8sClient.List(ctx, serviceList, listOpts...)
	if err != nil {
		return nil, err
	}

	for _, item := range deploymentList.Items {
		resources = append(resources, ResourceState{
			Kind:      "Deployment",
			Name:      item.GetName(),
			Namespace: item.GetNamespace(),
			Spec:      item.Spec,
		})
	}

	for _, item := range serviceList.Items {
		resources = append(resources, ResourceState{
			Kind:      "Service",
			Name:      item.GetName(),
			Namespace: item.GetNamespace(),
			Spec:      item.Spec,
		})
	}

	// TODO: Add more resources

	return resources, nil
}
