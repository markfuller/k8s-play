package main

import (
	"github.com/lyraproj/crd-mod/informer"
)

func main() {
	//we set the group here to meta.k8s.io because that the requested resource is
	createCRD("meta.k8s.io", "noddy", "noddies")
	informer.Start("noddies")
}
