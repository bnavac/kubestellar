# Notes about OCM

Relationship is complicated
Kubestellar relies on OCM, and uses some stuff from it
View Kubestellar as a compliment to OCM
Maybe switch out of OCM?
There is a plug-in architecture (two?)
Documentation needs improvement
Two boundaries, one is the binding objects and work-status objects
The independent part of Kubestellar is the part that produces binding objects and consumes work-status objects
As well as the Kubestellar API
OCM part takes in manifest work objects

OCM is simpler than Kubestellar
They have "hub cluster," a managed cluster which acts as the central cluster
OCM has two ways to manage what's in the managed clusters
Application path centered around manifest work objects
Manfiest work objects are objects that the user puts into the managed cluster
In OCM, there is a 1-1 relationship between a managed cluster and a namespace in the hub cluster
So the user puts the stuff in a manifest work object into the namespace, and then the hub sends it to the managed cluster
A manifest work also has a status section, which contains information about things such as creates and updates
But the info that flows back is only about OCM's progress, particularly in the forward propagation.
It does not really try to propagate the status of the manifest clusters back to the hub
OCM has an add-on architecture, which allows for add-ons
Kubestellar includes an OCM add-on, but is separate
The Kubestellar OCM add-on copies the status of the workstatus object into an ITS

Maybe will take advantage of newer features in OCM?
Would be nice to do it without an add-on
Not enough time to investigate it carefully?

Not all Kubernetes objects have a spec and status, and some may have annotations and labels.
So there may not be a clear state between a desired and preferred state
So the add-on is not exactly what they want, as it only returns status, and they want the general state.

Alternatives to OCM?
Would like to write a plugin to use a CDN to do the forward propagation. But CDN's do not typically do a backwards propagation.
A few options to remedy this: Create a way to do backward propagation while using a CDN for forward propagation.

Combined status - Primary concept of getting a reported state back, which aggregates the reported state back into the network
Started Kubestellar thinking about edge computing, thinking about scalability.
But can't centralized everything because of that, so maybe a possibility of having a hierarchical network.
No plans to merge it right now, but things can change.
In the kubestellar org, the kubestellar ocm add-on is defined there.
Extending the content of what's in the workstatus involves changing the status in both.
Though there is a bit of a design problem, Can potentially add new plugins, or change existing plugins.
There is also work being done in the consumption of the workstatus object, so there is not a desire for any update to the consumption of the workstatus object to require frequent changes to the workstatus object.
So the workstatus object should include a lot of information. 
But the more you put in, the more that goes into binding objects and is transported over.
Adding information in an API object is relatively easy. 
Though there needs to be a change in two repos, one to produce new content, and one to consume new stuff
But it is more difficult to alter or remove content
But right now, in the work status object, there is only the possibility to change the workstatus object, so a change there may allow the returning of more objects. 

Prob do not want to add the manifest work object into the work status object,
Prob do not want to add a lot of information into the work status object, just add the fields in the pr we are working on. (Add those fields to the workstatus object in the OCM add-on).

No need to worry about continuous operation due to no long-standing systems supported, so it is easy to change stuff around and update version.
Continuous operations a future thing.

In the Kubestellar repo, use Prow??? In ocm repo, use the regular GitHub mechanisms. 
In the kubestellar repo, a lot of review is done in PR's using prow. 
OCM does not use Prow, and so review is a lot simpler.
Once an OCM pr is approved, a new release will be pushed out.
But there is no update to go.mod as there are no client stubs
It instead uses dynamic clients, which can deal with any kind of object through a generic API (???)
But that means there is no static type checking as there is no easy way to make sure that the change you make is good
