package main

import (
    "context"
    "flag"
    "log"
    "os"

    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "sigs.k8s.io/controller-runtime/pkg/client/config"
    "sigs.k8s.io/controller-runtime/pkg/manager"
    "sigs.k8s.io/controller-runtime/pkg/manager/signals"

    flowcdv1alpha1 "github.com/yourusername/flowcd/pkg/apis/flowcd/v1alpha1"
    "github.com/yourusername/flowcd/pkg/controller"
)

var (
    scheme = runtime.NewScheme()
)

func init() {
    flowcdv1alpha1.AddToScheme(scheme)
}

func main() {
    kubeconfig := flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "absolute path to the kubeconfig file")
    flag.Parse()

    cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        log.Fatalf("Failed to build kubeconfig: %v", err)
    }

    mgr, err := manager.New(cfg, manager.Options{Scheme: scheme})
    if err != nil {
        log.Fatalf("Failed to create manager: %v", err)
    }

    if err := controller.SetupController(mgr); err != nil {
        log.Fatalf("Failed to setup controller: %v", err)
    }

    log.Println("Starting the FlowCD controller...")
    if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
        log.Fatalf("Failed to start manager: %v", err)
    }
}