High level view - User manages binding policies and clusters. Users sends binding policies and workloads to WDS, which deliver wrapped objects to ITS and receive statuses from ITS. The user sends the clusters to ITS which receives the wrapped objects and sends those workloads to clusters. The clusters send back statuses which are sent to WDS.

Kubestellar allows developers to define a binding policy between clusters. It then uses single-cluster tooling to deploy and configure lusters based on these policies.

ITS - Inventory and transport spaces.

WDS - workload defined spaces. 

Spaces - abstraction to represent API service that behaves like a Kubernetes kube-apiserver (including the persistent storage) and the controllers that are concerned with API machinery.

Eg - Kubeflex Controlplane.

A user can then perform some actions - Create WDSes to store definitions of workloads, applicaitons that run on Kubernetes. A workload can be made by one or several objects.
Create ITSes to manage the inventory of clusters and transport of workloads.
Register and label Workload Execution Clusters (WECs) with an ITWS to track available clusters and their characteristics.
Define a BindingPolicy to specify what objects and where should be deployed on the WECs.
Submit objects in the native Kubernetes format to the WDSes, the BindingPolicy will decide which WECs will receive them
Check the status of submitted objects from the WDS


Modules-

KubeFlex - uses kubeflex to track and provide ITS and WDS. Each appears as a ControlPlane object int he Kubeflex hosting cluster.
KubeStellar Controller Manager - Instantitated once per WDS and watches BindingPolicy objects and create a matching Binding object that lists all references to concrete objects and clusters, also updates the statuses in the WDS.
Pluggable Transport Controller (PTC) - Instantiated once per WDS and delivers workload objects from the WDS to the ITS according to the binding objects.
Space manager - manages the lifecycle of spaces.
OCM cluster manager - Instantiated onces per ITS and syncs objects from that ITS to the WECs. In the ITS, each mailbox is associated with one WEC.
OCM Agent - registers the WEC to the OCM, unwraps and syncs objects into the WEC.
OCM Status Add-On Controller - Module is instantiated once per ITS and uses the OCM Add-on Framework o get the OCM Status Add-On Agent installed in each WEC along with supporting RBAC objects.
OCM Status Add-On Agent - Watches AppliedManifestWork objects that are synced by the OCM agent and gets their statuses to send to WorkStatus objects in ITS namespaces associated with a WEC.


Kubestellar controller manager - Manages binding controller and the status controller
The binding ontroller watches BindingPolicy and workload objects on the WDS, maintains a binding object for each binding policy. The binding object contains references to a list of workload objects and list of clusters selected by the bindingpolicy.
The status controller watches for WOrkStatus objects on the ITS and updates the status of the WDS objects when a singleton status is requested in the bindingpolicy for those objects.
The is only one instance of a Kubestellar control manager for each WDS. This controller-manager runs in the KubeFlex hosting cluster and is responsible for installing the required CRDs in the associated WDS. 

Pluggable Transport Controller - 
Watches binding objects on the WDS, maintains a wrapped object per binding object in the ITS.
Pluggable and can be implemented with different options. Currently the only option supported is based on Open Cluster Management Project
Only one instance of PTC for each WDS. This runs in an executable process. 

Space Manager - handles the lifecycle of spaces. Kubestellar uses KubeFlex for space management. In KubeFlex, a space is named a ControlPlane.
Two spaces managed by KubeFlex, ITS and WDS. ITS runs the OCM cluster manager on a vcluster-type control plane, and WDS on a k8s type control plane.
An ITS holds the inventory and mailbox namespaces. This is anchored by the ManagedCluster open-cluster-management objects that describe the WECs. FOr each WEC there may be a ConfigMap object hat carries additional properties of the WEC. The mailbox and its contents are transport implementations. Each mailbox corresponds with a WEC and holds ManifestWork objects managed by the central KubeStellar controllers.
A WDS holds user workload objects. User control objects are BindingPolicies and Binding objects.

API Documentation - 

User controls downsyncs through API objects of BindingPolicyies and Bindings. BindingPolicies are higher level concepts and will be translated into bindings. 

Binding Policy - Two predicates to identifty a subset of WECs in the inventory of the ITS and identifies a subset of workload objects in the WDS.
The WEC selecting predicate is an array of labels in sepc.clusterSelectors. These label selectors test the labels of inventory objects describing the WECs. These are the ones whose inventory objects passes at least one label selector in spec.clusterSelectors.
The workload object selection predicate is in spec.downsync, which holds a list of DownsyncPolicyClauses, each includes a workload object selection predicate and two kinds of info that modulate the downsync.

Example policy
https://mikespreitzer.github.io/kcp-edge-mc/doc-more-propagation-data/direct/binding/