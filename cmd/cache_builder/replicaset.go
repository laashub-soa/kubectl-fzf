package main

import (
	"fmt"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
)

const ReplicaSetHeader = "Namespace Name Replicas AvailableReplicas ReadyReplicas Selector Age Labels\n"

// ReplicaSet is the summary of a kubernetes replicaSet
type ReplicaSet struct {
	ResourceMeta
	replicas          string
	readyReplicas     string
	availableReplicas string
	selectors         []string
}

// NewReplicaSetFromRuntime builds a k8sresource from informer result
func NewReplicaSetFromRuntime(obj interface{}) K8sResource {
	p := &ReplicaSet{}
	p.FromRuntime(obj)
	return p
}

// FromRuntime builds object from the informer's result
func (r *ReplicaSet) FromRuntime(obj interface{}) {
	replicaSet := obj.(*appsv1.ReplicaSet)
	r.FromObjectMeta(replicaSet.ObjectMeta)
	r.replicas = strconv.Itoa(int(replicaSet.Status.Replicas))
	r.readyReplicas = strconv.Itoa(int(replicaSet.Status.ReadyReplicas))
	r.availableReplicas = strconv.Itoa(int(replicaSet.Status.AvailableReplicas))
	r.selectors = JoinStringMap(replicaSet.Spec.Selector.MatchLabels,
		ExcludedLabels, "=")
}

// HasChanged returns true if the resource'r dump needs to be updated
func (r *ReplicaSet) HasChanged(k K8sResource) bool {
	oldRs := k.(*ReplicaSet)
	return (r.replicas != oldRs.replicas ||
		r.readyReplicas != oldRs.readyReplicas ||
		r.availableReplicas != oldRs.availableReplicas ||
		StringSlicesEqual(r.selectors, oldRs.selectors) ||
		StringMapsEqual(r.labels, oldRs.labels))
}

// ToString serializes the object to strings
func (r *ReplicaSet) ToString() string {
	selectorList := JoinSlicesOrNone(r.selectors, ",")
	line := strings.Join([]string{r.namespace,
		r.name,
		r.replicas,
		r.availableReplicas,
		r.readyReplicas,
		selectorList,
		r.resourceAge(),
		r.labelsString(),
	}, " ")
	return fmt.Sprintf("%s\n", line)
}