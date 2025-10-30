package main

import (
    "flag"
    "os"

    "k8s.io/apimachinery/pkg/runtime"
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
    // Add FlowCD types to scheme
    flowcdv1alpha1.AddToScheme(scheme)
}

func main() {
    // Setup logger
    ctrl.SetLogger(zap.New())

    // Parse flags
    var kubeconfig string
    flag.StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file")
    flag.Parse()

    // Create manager
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