package main

import (
	"log"
	"time"

	"github.com/sagnik3788/gitops-controller/git"
	"github.com/sagnik3788/gitops-controller/k8s"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	repoURL := "https://github.com/sagnik3788/Gitops-controller.git"
	branch := "main"
	path := "manifest"

	//setup k8s client
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Go routine to sync the repo and apply manifests to the cluster
	go func() {
		for {
			log.Println("Syncing repo and applying manifests...")
			err := git.CloneOrPull(repoURL, branch, path)
			if err != nil {
				log.Println("Git sync failed:", err)
			} else {
				err = k8s.ApplyManifests(path, mgr.GetClient(), mgr.GetScheme())
				if err != nil {
					log.Println("Apply failed:", err)
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()

	log.Println("GitOps controller running...")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatal(err)
	}
}
