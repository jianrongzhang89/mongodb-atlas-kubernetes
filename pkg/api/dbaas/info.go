/*
Copyright 2021 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package v1 contains API Schema definitions for the dbaas.redhat.com v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=dbaas.redhat.com
package dbaas

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	DBaaSGroupVersion = schema.GroupVersion{Group: "dbaas.redhat.com", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	DBaaSSchemeBuilder = &scheme.Builder{GroupVersion: DBaaSGroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = DBaaSSchemeBuilder.AddToScheme

	_ = &MongoDBAtlasInventory{}

	_ = &MongoDBAtlasConnection{}
)