/*
Copyright 2023 The KCP Authors.

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

package placement

import (
	"fmt"
	"strings"

	"github.com/kcp-dev/logicalcluster/v3"
	k8sevents "k8s.io/api/events/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachtypes "k8s.io/apimachinery/pkg/types"
	k8ssets "k8s.io/apimachinery/pkg/util/sets"

	edgeapi "github.com/kcp-dev/edge-mc/pkg/apis/edge/v1alpha1"
)

// This file contains declarations of the interfaces between the main parts
// of the placement translator.  Those are as follows.
//
// - a WhereResolver that monitors the "where" resolutions in the
//   SinglePlacementSlice API objects.
// - a WhatResolver that monitors the EdgePlacement objects and resolves
//   their "what" predicates.
// - a SetBinder that keeps track of the bindings between resolved where
//   and resolved what.
// - a WorkloadProjector that maintains the customized workload copies
//   in the mailbox workspaces.
// - a PlacementProjector that maintains the TMC Placement objects that
//   correspond to the EdgePlacement objects.

// WhereResolver is responsible for keeping given consumers eventually
// consistent with the resolution of the "where" predicate for each EdgePlacement
// (identified by cluster and name).
type WhereResolver DynamicMapProducer[edgeapi.ExternalName, ResolvedWhere]

// WhatResolver is responsible for keeping its consumers eventually consistent
// with the resolution of the "what" predicate of each EdgePlacement
// (identified by cluster ane name).
type WhatResolver DynamicMapProducer[edgeapi.ExternalName, WorkloadParts]

// SetBinder is a component that is kept appraised of the "what" and "where"
// resolutions and reorganizing and picking API versions to guide the
// workload projector and the placement projector.
// The implementation may atomize the resolved "where" using a SinglePlacementSliceSetReducer.
// The implementation may use a BindingOrganizer to get from the atomized
// "what" and "where" to the ProjectionMapProducer behavior.
type SetBinder interface {
	AsWhatConsumer() DynamicMapConsumer[edgeapi.ExternalName, WorkloadParts]
	AsWhereConsumer() DynamicMapConsumer[edgeapi.ExternalName, ResolvedWhere]
	ProjectionMapProducer
}

// WorkloadProjector is kept appraised of what goes where
// and is responsible for maintaining the customized workload
// copies in the mailbox workspaces.
type WorkloadProjector Client[ProjectionMapProducer]

// PlacementProjector is responsible for maintaining the TMC Placement
// objects that cause propagation between mailbox workspace and edge cluster.
type PlacementProjector Client[ProjectionMapProducer]

// AssemplePlacementTranslator puts together the top-level pieces.
func AssemplePlacementTranslator(
	whatResolver WhatResolver,
	whereResolver WhereResolver,
	setBinder SetBinder,
	workloadProjector WorkloadProjector,
	placementProjector PlacementProjector,
) {
	whatResolver.AddConsumer(setBinder.AsWhatConsumer().Set)
	whereResolver.AddConsumer(setBinder.AsWhereConsumer().Set)
	workloadProjector.SetProvider(setBinder)
	placementProjector.SetProvider(setBinder)
}

// ResolvedWhere identifies the set of SyncTargets that match a certain
// EdgePlacement's "where" predicate.
// Each `*edgeapi.SinglePlacementSlice` points to an immutable object.
type ResolvedWhere []*edgeapi.SinglePlacementSlice

// WorkloadParts identifies a workload prescription and provides
// ephemeral details of how to access it.
// A workload prescription is the things that match the "what" predicate
// of an EdgePlacement.
//
// Every WorkloadParts that appears in the interfaces here is immutable.
//
// In the case of a Namespace object, this implies that
// all the objects in that namespace are included.
//
// A workload may include objects of kinds that are built into
// the edge cluster.  By built-in we mean that these kinds are both
// already known (regardless of whether it is via being built into
// the apiserver or added to it by either form of aggregation) and
// not managed by edge workload management.
// It is the user's responsibility to make the "what" predicate
// match the corresponding CRD when the workload includes an object
// of a kind that is not built into the edge cluster.
type WorkloadParts map[WorkloadPartID]WorkloadPartDetails

// WorkloadPartID identifies part of a workload.
type WorkloadPartID struct {
	APIGroup string

	// Resource is the lowercase plural way of identifying the kind of object
	Resource string

	Name string
}

// WorkloadPartDetails provides additional details about how the WorkloadPart
// is to be included.
type WorkloadPartDetails struct {
	// APIVersion is version (no group) that the source workspace prefers to serve.
	APIVersion string

	// IncludeNamespaceObject is only interesting for a Namespace part, and
	// indicates whether to include the details of the Namespace object;
	// the objects in the namespace are certainly included.
	// For other parts, this field holds `false`.
	IncludeNamespaceObject bool
}

type WorkloadPart struct {
	WorkloadPartID
	WorkloadPartDetails
}

// ProjectionMapProducer tells the consumers what to project,
// organized into three levels.
type ProjectionMapProducer DynamicMapProducer[ProjectionKey, *ProjectionPerCluster]

// ProjectionKey identifies a source/kind/destination relationship
type ProjectionKey struct {
	metav1.GroupResource
	Destination SinglePlacement
}

// SinglePlacement extends the API struct with the UID of the SyncTarget.
type SinglePlacement struct {
	edgeapi.SinglePlacement
	SyncTargetUID apimachtypes.UID
}

func (sp SinglePlacement) SyncTargetRef() edgeapi.ExternalName {
	return edgeapi.ExternalName{Workspace: sp.Location.Workspace, Name: sp.SyncTargetName}
}

type ProjectionPerCluster struct {
	// APIVersion is the version to read.  Just the version, no group included
	APIVersion string

	PerSourceCluster DynamicMapProducer[logicalcluster.Name, ProjectionDetails]
}

// ProjectionDetails modulates projection
type ProjectionDetails struct {

	// For namespaced resoruces, Namespaces can optionally be non-nil to restrict
	// the namespaces read from.
	Namespaces *k8ssets.String

	// For non-namespaced objects, Names can optionally be non-nil to restrict
	// the objects handled.
	Names *k8ssets.String
}

// BindingOrganizer produces a SingleBinder and a corresponding map producer
// that reflects the result of combining the single bindings and resolving
// the API group version issue.
// A SetBinder implementation will likely use one of these to provide its
// ProjectionMapProducer producer, feeding the SingleBinder atomized changes
// from the incoming ResolvedWhat and ResolvedWhere values.
// The given EventHandler is given events that the organizer produces
// and publishes them somewhere.
type BindingOrganizer func(discovery APIMapProvider, resourceModes ResourceModes, eventHandler EventHandler) (SingleBinder, ProjectionMapProducer)

// SingleBinder is appraised of individual bindings and unbindings,
// but they may come in batches.
// AddBinding calls are ordered by API machinery dependencies.
// RemoveBinding calls are ordered by the reverse of the API machinery dependencies.
type SingleBinder interface {
	// Transact does a collection of adds and removes.
	Transact(func(SingleBindingOps))
}

type SingleBindingOps interface {
	AddBinding(what WorkloadPart, where SinglePlacement)
	RemoveBinding(what WorkloadPart, where SinglePlacement)
}

// SinglePlacementSliceSetReducer keeps track of a ResolvedWhere.
// Typically one of these has a UIDer and
// a SinglePlacementSetChangeConsumer that is kept
// appraised of the set of values in those slices, extended by the
// SyncTarget UIDs.
// A SetBinder will likely use a SinglePlacementSliceSetReducer
// to atomize its incoming ResolvedWhere values.
type SinglePlacementSliceSetReducer interface {
	Set(ResolvedWhere)
}

// SinglePlacementSetChangeConsumer is something that is kept
// incrementally appraised of a set of SinglePlacement values.
type SinglePlacementSetChangeConsumer interface {
	Add(SinglePlacement)
	Remove(SinglePlacement)
}

// APIMapProvider provides API information on a cluster-by-cluster basis,
// as needed by clients.
// This information comes from runtime monitoring of the API resources
// of the clusters.
type APIMapProvider interface {
	// AddClient adds a client for a cluster.
	// All clients for the same cluster get the same producer.
	AddClient(cluster logicalcluster.Name, client Client[ScopedAPIProducer])

	// RemoveClient removes a client for a cluster.
	// Clients must be comparable.
	// Removing the last client for a given cluster causes release of
	// internal computational resources.
	RemoveClient(cluster logicalcluster.Name, client Client[ScopedAPIProducer])
}

// ScopedAPIGroupVersioner is specific to one logical cluster and
// provides a map from API group name to details about it in that cluster.
// A nil pointer means that the group is not defined in the cluster.
type ScopedAPIProducer DynamicMapProducer[string, *APIGroupInfo]

type APIGroupInfo struct {
	Versions  DynamicValueProducer[APIGroupVersions]
	Resources APIResourceDetailsProducer
}

// APIGroupVersions tells about the versions available for the group.
type APIGroupVersions struct {
	// Versions are ordered as semantic versions.
	// This slice is immutable.
	Versions []metav1.GroupVersionForDiscovery

	PreferredVersion metav1.GroupVersionForDiscovery
}

// APIResourceDetailsProducer reveals details about the resoureces in a
// given API group in a given cluster.
// A resource is identified in the usual way, the lowercase plural.
// A nil pointer means the resource is not defined there.
type APIResourceDetailsProducer DynamicMapProducer[string, *ResourceDetails]

// ResourceDetails holds the information needed here about a resource
type ResourceDetails struct {
	Namespaced        bool
	SupportsInformers bool
}

// ResourceModes tells the handling of all the resources that do not
// get default handling, and maybe some that do.
// This information comes from platform configuration and code.
// Immutable.
type ResourceModes map[metav1.GroupResource]ResourceMode

func (modes ResourceModes) Get(gr metav1.GroupResource) ResourceMode {
	if ans, ok := modes[gr]; ok {
		return ans
	}
	return DefaultResourceMode
}

// ResourceMode describes how a given resource is handled regarding
// propagation and denaturing.
type ResourceMode struct {
	PropagationMode PropagationMode
	NatureMode      NatureMode
	BuiltinToEdge   bool
}

// DefaultResourceMode is the handling for every user-defined resource
var DefaultResourceMode = ResourceMode{
	PropagationMode: GoesToEdge,
	NatureMode:      NaturalyDenatured,
	BuiltinToEdge:   false,
}

// PropagationMode describes the relationship between present-in-center and present-in-edge
type PropagationMode string

const (
	ErrorInCenter    PropagationMode = "error"
	TolerateInCenter PropagationMode = "tolerate"
	GoesToEdge       PropagationMode = "propagate"
)

// NatureMode describes the stance regarding whether a resource is denatured in the center.
// All resources that go to the edge are natured (not denatured) at the edge.
type NatureMode string

const (
	// NaturalyDenatured is a resource that is denatured in the center without any special
	// effort in this code.
	NaturalyDenatured NatureMode = "NaturallyDenatured"

	// NaturallyNatured is a resource that is natured in the center and should be that way.
	NaturallyNatured NatureMode = "NaturallyNatured"

	// ForciblyDenatured is a resource that would be given an undesired interpretation in the center
	// if stored normally in the center, so has to be stored differently in the center (but not
	// at the edge).
	ForciblyDenatured NatureMode = "ForciblyDenatured"
)

// UIDer is a source of mapping from object name to UID.
// One of these is specific to one kind of object,
// which is not namespaced.
type UIDer DynamicMapProducer[edgeapi.ExternalName, apimachtypes.UID]

func (sp *SinglePlacement) MailboxWorkspaceName() string {
	return sp.Location.Workspace + WSNameSep + string(sp.SyncTargetUID)
}

const WSNameSep = "-mb-"

// EventHandler can be given Event objects.
type EventHandler interface {
	HandleEvent(*k8sevents.Event)
}

func (where ResolvedWhere) String() string {
	var builder strings.Builder
	builder.WriteRune('[')
	for idx, slice := range where {
		if idx > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", slice))
	}
	builder.WriteRune(']')
	return builder.String()
}
