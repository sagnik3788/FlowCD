package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	flowcdv1alpha1 "github.com/sagnik3788/FlowCD/pkg/apis/flowcd/v1alpha1"
)

type FlowCDReconciler struct {
	// client to interact with k8s api server
	client.Client

	// Type scheme
	Scheme *runtime.Scheme

	// Logger
	Log    logr.Logger
}

func (r *FlowCDReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("flowcd", req.NamespacedName)

	// Get FlowCD resource
	var flowcd flowcdv1alpha1.FlowCD
	if err := r.Get(ctx, req.NamespacedName, &flowcd); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Log what we found
	log.Info("Reconciling FlowCD", 
        "repoURL", flowcd.Spec.Source.RepoURL,
        "branch", flowcd.Spec.Source.Branch,
        "namespace", flowcd.Spec.Destination.Namespace)
    
    // TODO: Add Git operations here
    // TODO: Add K8s operations here  
    // TODO: Update status here
    
    return ctrl.Result{}, nil
}