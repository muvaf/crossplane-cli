/*
Copyright 2020 The Crossplane Authors.

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

package crossplane

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// TODO(hasan): Add other resources related to networking and/or iam
	groupKindsClaimCluster = []string{
		"kubernetescluster.compute.crossplane.io",
	}
	groupKindsClaim = []string{
		"mysqlinstance.database.crossplane.io",
		"rediscluster.cache.crossplane.io",
		"postgresqlinstance.database.crossplane.io",
		"bucket.storage.crossplane.io",
	}
	groupKindsManagedCluster = []string{
		"akscluster.compute.azure.crossplane.io",
		"gkecluster.compute.gcp.crossplane.io",
		"ekscluster.compute.aws.crossplane.io",
	}
	groupKindsManaged = []string{
		// Azure
		"redis.cache.azure.crossplane.io",
		"mysqlserver.database.azure.crossplane.io",
		"postgresqlserver.database.azure.crossplane.io",
		"container.storage.azure.crossplane.io",
		"account.storage.azure.crossplane.io",

		// GCP
		"cloudsqlinstance.database.gcp.crossplane.io",
		"cloudmemorystoreinstance.cache.gcp.crossplane.io",
		"bucket.storage.gcp.crossplane.io",

		// AWS
		"replicationgroup.cache.aws.crossplane.io",
		"rdsinstance.database.aws.crossplane.io",
		"s3bucket.storage.aws.crossplane.io",
	}
	groupKindsPortableClass = []string{
		"mysqlinstanceclass.database.crossplane.io",
		"kubernetesclusterclass.compute.crossplane.io",
		"redisclusterclass.cache.crossplane.io",
		"postgresqlinstanceclass.database.crossplane.io",
		"bucketclass.storage.crossplane.io",
	}
	groupKindsNonPortableClass = []string{
		// Azure
		"redisclass.cache.azure.crossplane.io",
		"aksclusterclass.compute.azure.crossplane.io",
		"sqlserverclass.database.azure.crossplane.io",

		// GCP
		"cloudsqlinstanceclass.database.gcp.crossplane.io",
		"gkeclusterclass.compute.gcp.crossplane.io",
		"cloudmemorystoreinstanceclass.cache.gcp.crossplane.io",
		"bucketclass.storage.gcp.crossplane.io",

		// AWS
		"replicationgroupclass.cache.aws.crossplane.io",
		"eksclusterclass.compute.aws.crossplane.io",
		"rdsinstanceclass.database.aws.crossplane.io",
		"s3bucketclass.storage.aws.crossplane.io",
	}
	groupKindsProvider = []string{
		"provider.gcp.crossplane.io",
		"provider.azure.crossplane.io",
		"provider.aws.crossplane.io",
	}
	groupKindsApplication = []string{
		"kubernetesapplication.workload.crossplane.io",
	}
	groupKindsApplicationResource = []string{
		"kubernetesapplicationresource.workload.crossplane.io",
	}
)

// Object is an interface for getting information of Crossplane objects.
type Object interface {
	GetStatus() string
	GetObjectDetails() ObjectDetails
	GetAge() string
	GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error)
	IsReady() bool
}

// ObjectDetails struct contains detailed information of a Crossplane object.
type ObjectDetails struct {
	// Common details for all objects.
	Kind       string
	Name       string
	Namespace  string
	Conditions []map[string]string

	// Details available only for a subset of objects.
	RemoteStatus string

	// Each object provides own AdditionalStatusColumns from status field.
	AdditionalStatusColumns []map[string]string
}

// ObjectFromUnstructured returns an Object is provided unstructured.Unstructured is a Crossplane object.
func ObjectFromUnstructured(u *unstructured.Unstructured) Object {
	gvk := u.GroupVersionKind()
	if isClaim(gvk) {
		return NewClaim(u)
	} else if isManaged(gvk) {
		return NewManaged(u)
	} else if isPortableClass(gvk) {
		return NewPortableClass(u)
	} else if isNonPortableClass(gvk) {
		return NewNonPortableClass(u)
	} else if isProvider(gvk) {
		return NewProvider(u)
	} else if isApplication(gvk) {
		return NewApplication(u)
	} else if isApplicationResource(gvk) {
		return NewApplicationResource(u)
	}
	return nil
}

func isClaim(gvk schema.GroupVersionKind) bool {
	return stringInSlice(normalizedGroupKind(gvk), groupKindsClaim) || stringInSlice(normalizedGroupKind(gvk), groupKindsClaimCluster)
}
func isManaged(gvk schema.GroupVersionKind) bool {
	return stringInSlice(normalizedGroupKind(gvk), groupKindsManaged) || stringInSlice(normalizedGroupKind(gvk), groupKindsManagedCluster)
}
func isNonPortableClass(gvk schema.GroupVersionKind) bool {
	return stringInSlice(normalizedGroupKind(gvk), groupKindsNonPortableClass)
}
func isPortableClass(gvk schema.GroupVersionKind) bool {
	return stringInSlice(normalizedGroupKind(gvk), groupKindsPortableClass)
}
func isProvider(gvk schema.GroupVersionKind) bool {
	return stringInSlice(normalizedGroupKind(gvk), groupKindsProvider)
}
func isApplication(gvk schema.GroupVersionKind) bool {
	return stringInSlice(normalizedGroupKind(gvk), groupKindsApplication)
}
func isApplicationResource(gvk schema.GroupVersionKind) bool {
	return stringInSlice(normalizedGroupKind(gvk), groupKindsApplicationResource)
}
