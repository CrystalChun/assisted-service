---
title: cancel-day-2-host-installation
authors:
  - "@CrystalChun"
creation-date: 2025-03-25
last-updated: 2025-03-25
---

# Cancel Day-2 Host Installation

<!-- This is the title of the enhancement. Keep it simple and descriptive. A good
title can help communicate what the enhancement is and should be considered as
part of any review.

The YAML `title` should be lowercased and spaces/punctuation should be
replaced with `-`.


The `Metadata` section above is intended to support the creation of tooling
around the enhancement process.

Make sure that the enhancement covers all scenarios that are relevant:
- REST API and Kubernetes API
- Late-binding and bind-on-discovery
- SNO, multinode, compact, and 3+1 clusters
- Day1 (cluster installation) and day2 (adding hosts to an installed cluster)
- Day2 both for installed clusters and imported clusters
- Observability: feature usage, events, status info, etc.
- Network configuration: static networking vs DHCP, IPv4/v6/dual-stack
- Supported platforms (e.g., baremetal, none, vsphere, nutanix)
- Supported CPU architectures
- Feature support based on OCP version
- OLM operators
- Agent upgrade -->

## Summary

<!-- The `Summary` section is incredibly important for producing high quality
user-focused documentation such as release notes or a development roadmap. It
should be possible to collect this information before implementation begins in
order to avoid requiring implementors to split their attention between writing
release notes and implementing the feature itself.

A good summary is probably at least a paragraph in length. -->


## Motivation

<!-- This section is for explicitly listing the motivation, goals and non-goals of
this proposal. Describe why the change is important and the benefits to users. -->

### Goals

<!-- List the specific goals of the proposal. How will we know that this has succeeded? -->
- Once a day-2 host starts installation, it can be cancelled
- There's a straight-forward way to trigger this cancellation that must be intentionally set (can't accidentally cancel it)
- Host whose installations are cancelled can be reclaimed back to their InfraEnvs and reused
- Host cancellation should not affect other hosts in the cluster and should not affect the cluster itself
- Host cancellation should happen within a reasonable amount of time in order to support dynamic scaling 

### Non-Goals

<!-- What is out of scope for this proposal? Listing non-goals helps to focus discussion
and make progress. -->
- Hosts that are installed in a day-1 cluster should not be allowed to have their installation cancelled at a host level
- This will not focus on cluster installation cancellation 
- This is only for day-2 hosts
- Reclaim/unbind will not be affected by these changes

## Proposal

<!-- This is where we get down to the nitty gritty of what the proposal actually is.
 -->

### User Stories

<!-- Detail the things that people will be able to do if this is implemented.
Include as much detail as possible so that people can understand the "how" of
the system. The goal here is to make this feel real for users without getting
bogged down.

Include a story on how this proposal will be deployed in production: 
lifecycle, monitoring and scale requirements or benefits. -->

#### Story 1

As a cluster admin, when I create a hosted cluster using Hypershift, I'd like to bind and install worker hosts into my cluster. 
However, if I want to remove the host from the cluster before they finish installing, I should be able to do so and not have them stuck.


#### Story 2

As a cluster admin using Hypershift, I often scale up and scale down my cluster in quick succession. When the scaling happens and hosts are bound/unbound from the cluster,
they should not become stuck and should bind and unbind easily.

### Implementation Details/Notes/Constraints [optional]

<!-- What are the caveats to the implementation? What are some important details that
didn't come across above. Go in to as much detail as necessary here. This might
be a good place to talk about core concepts and how they relate.

This is an excellent place to call out changes that need to be made in projects
other than openshift/assisted-service. For example, if a change will need to be
made in the agent (openshift/assisted-installer-agent) and the installer
(openshift/assisted-installer; it should be mentioned here. -->

There is an already existing cancel host installation API
We will need a way to trigger that the host's installation is canceled from the kube API
My suggestion would be when a host becomes unbound and it's in the middle of installation, we know to cancel and reclaim it
Would BMH affect this?

### Risks and Mitigations

<!-- What are the risks of this proposal and how do we mitigate. Think broadly. For
example, consider both security and how this will impact the larger OKD
ecosystem. 

Will choices made here affect adoption of assisted-installer?
Will this work make it harder to integrate with other upstream projects?
How will security be reviewed and by whom? How will UX be reviewed and by whom?

Consider including folks that also work outside your immediate sub-project. -->

#### Risk 1: Host can't be reclaimed 

Host installation cancellation is a completely new flow that could have many potential consequences such as a weird transition state like what
if it becomes stuck? 

##### Mitigation

The way to mitigate this is to fully explore and ensure all paths into and out of this state are covered

#### Risk 2: If a host installation can't be canceled 

The host might be in a stage of installation that prevents assisted from stopping the installation right away 

## Design Details [optional]

If an enhancement is complex enough, design details should be included. When not
included, reviewers reserve the right to ask for this section to be filled in to
enable more thoughtful discussion about the enhancement and it's impact.

### Open Questions

This is where to call out areas of the design that require closure before deciding
to implement the design.  For instance,
 > 1. This requires exposing previously private resources which contain sensitive
  information.  Can we do this?

### UI Impact

No need to go into great detail about the UI changes that will need to be made.
However, this is an excellent time to mention 1) if UI changes are required if
this enhancement were accepted and 2) at a high-level what those UI changes
would be.

### Test Plan

Consider the following in developing a test plan for this enhancement:
- Will there be e2e and integration tests, in addition to unit tests?
- How will it be tested in isolation vs with other components?

No need to outline all of the test cases, just the general strategy. Anything
that would count as tricky in the implementation and anything particularly
challenging to test should be called out.


## Drawbacks

The idea is to find the best form of an argument why this enhancement should _not_ be implemented.

## Alternatives

Similar to the `Drawbacks` section the `Alternatives` section is used to
highlight and record other possible approaches to delivering the value proposed
by an enhancement.
