package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion define "gdch.mycompany.io/v1"
	GroupVersion = schema.GroupVersion{Group: "gdch.mycompany.io", Version: "v1"}
)
