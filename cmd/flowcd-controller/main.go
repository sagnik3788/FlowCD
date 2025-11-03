package main

import (
	"flag"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	flowcdv1alpha1 "github.com/sagnik3788/FlowCD/pkg/apis/flowcd/v1alpha1"
	"github.com/sagnik3788/FlowCD/pkg/controller"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(flowcdv1alpha1.AddToScheme(scheme))

	metav1.AddToGroupVersion(scheme, schema.GroupVersion{
		Group:   "flowcd.io",
		Version: "v1alpha1",
	})
}

func main() {
	// Setup logger
	ctrl.SetLogger(zap.New())

	// Parse flags (controller-runtime handles kubeconfig automatically)
	flag.Parse()

	// Create manager (automatically uses --kubeconfig flag or in-cluster config)
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		os.Exit(1)
	}

	// Setup FlowCD controller
	if err = (&controller.FlowCDReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("FlowCD"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FlowCD")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
