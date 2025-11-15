package controller

import (
	"context"
	"os"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	flowcdv1alpha1 "github.com/sagnik3788/FlowCD/pkg/apis/flowcd/v1alpha1"
	"github.com/sagnik3788/FlowCD/pkg/git"
	"github.com/sagnik3788/FlowCD/pkg/kubernetes"
)

type FlowCDReconciler struct {
	// client to interact with k8s api server
	client.Client

	// Type scheme
	Scheme *runtime.Scheme

	// Logger
	Log logr.Logger
}

func (r *FlowCDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&flowcdv1alpha1.FlowCD{}).
		Complete(r)
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

	// Initialize Git client with unique directory
	gitClient, err := git.NewGitClient("/tmp")
	if err != nil {
		log.Error(err, "failed to create git client")
		return ctrl.Result{}, err
	}

	// Clone repository
	if err := gitClient.Clone(flowcd.Spec.Source.RepoURL, flowcd.Spec.Source.Branch); err != nil {
		log.Error(err, "failed to clone repository")
		return ctrl.Result{}, err
	}

	// Initialize kubernetes client
	k8sClient := kubernetes.NewClient(r.Client, r.Scheme)

	// Get live state from cluster
	resources, err := ListResources(ctx, r.Client, flowcd.Spec.Destination.Namespace)
	if err != nil {
		log.Error(err, "failed to get live state from cluster")
		return ctrl.Result{}, err
	}

	log.Info("Listed resources", "count", len(resources))

	//temporary
	for _, r := range resources {
		log.Info("Resource found", "kind", r.Kind, "name", r.Name, "namespace", r.Namespace)
	}

	// Get manifest files
	files, err := gitClient.GetManifest(flowcd.Spec.Source.Path)
	if err != nil {
		log.Error(err, "failed to get manifest files")
		return ctrl.Result{}, err
	}

	log.Info("Found manifest files", "files", files)
	if len(files) == 0 {
		log.Info("No manifest files found in path", "path", flowcd.Spec.Source.Path)
		return ctrl.Result{}, nil
	}

	// Collect all manifests from all files
	var allManifests []unstructured.Unstructured
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return ctrl.Result{}, err
		}

		manifests, err := kubernetes.ParseManifests(data)
		if err != nil {
			return ctrl.Result{}, err
		}

		allManifests = append(allManifests, manifests...)
	}

	// pass the manifests to GetDesiredState
	desiredState, err := GetDesiredState(allManifests, flowcd.Spec.Destination.Namespace)
	if err != nil {
		log.Error(err, "failed to get desired state")
		return ctrl.Result{}, err
	}

	log.Info("Desired state resources", "count", len(desiredState))

	// Compare live state vs desired state to detect drift
	drift, err := Compare(resources, desiredState)
	if err != nil {
		log.Error(err, "failed to compare states")
		return ctrl.Result{}, err
	}

	log.Info("Drift detection result",
		"hasChanges", drift.HasChanges,
		"missing", len(drift.Missing),
		"extra", len(drift.Extra),
		"modified", len(drift.Modified),
		"summary", drift.Summary)

	// Only apply if there are changes
	totalApplied := 0
	if drift.HasChanges {
		log.Info("Changes detected, applying manifests", "summary", drift.Summary)

		// Apply all manifests
		if err := k8sClient.Apply(ctx, allManifests, flowcd.Spec.Destination.Namespace); err != nil {
			log.Error(err, "failed to apply manifests")
			return ctrl.Result{}, err
		}

		totalApplied = len(allManifests)
		log.Info("Successfully applied manifests", "count", totalApplied)
	} else {
		log.Info("No changes detected, skipping apply", "summary", drift.Summary)
	}

	log.Info("Deployment completed successfully",
		"totalResources", totalApplied,
		"namespace", flowcd.Spec.Destination.Namespace)

	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}
