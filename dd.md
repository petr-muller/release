# CI cluster configuration in openshift/release

## Background

The `api.ci` cluster and the configuration checked-in to `openshift/release`
repository grew organically. Together with the fact that there are multiple
teams and persons who manage the `api.ci` cluster, the result is that the
configuration and the way it is applied to the cluster is a little chaotic. The
goal of DPTP-XXX is to set up config structure that is easier to manage, modify
and apply. The initial goal is not to clean up the existing config, but to
create a structure into which even the legacy config would be easily migrated.

## Investigation

### Inventory

Existing config files and config apply scripts were [inventorized as of
2019-06-13](link).

### Main use-cases

1. DPTP operate the cluster itself and managing necessary cluster-level
   resources
2. DPTP operating Prow and other resources necessary to run CI jobs
   (ci-operator)
3. OpenShift architects maintain structures that define how the products are
   assembled
4. OpenShift architects maintain tools that build the products
   (release-controller)
5. OpenShift architects and related component owners maintain jobs that publish
   images to Quay
6. Component owners maintain their CI config
7. OpenShift engineers need various images available for different purposes in
   existing namespaces
8. OpenShift engineers operate long-term services themselves for various
   purposes (ci-chat-bot, metering, azure).
9. OpenShhift engineers create short-term projects for experiments

### Main problems

1. Config fragmented in all repository (`ci-operator/infra`,
   `cluster/ci/config`, `projects`)
2. Unclear owners / purpose / lifetime
3. `Makefile` mess (everything goes there, no clear entrypoints except
   `postsubmit-update`)
4. Misplaced items (config or scripts) put somewhere by a wrong guess
5. Chaotic namespace placement. Sometimes the manifests specify `namespace`.
   Sometimes scripts do `oc apply -n`. Sometimes scripts rely on a previous call
   of `oc project`.

## Proposals

### Top-level repo-layout

- ci-operator: ci-operator config, prow jobs, ci-operator templates (personally,
  I would move cluster-profiles here, too)
- "core services": usually dptp or architects-owned services. quality bar, full
  automation. At least a bit widely useful to the org.
- "sandbox": custom services and experiments. lower quality bar, must not harm
  core services
- "scripts": executable stuff operating over all repo (personally, I would break
  this down and move scripts that operate over just some part to that part)
- "tools": buildable/executable stuff that's somehow helpful in more general way

Things could "graduate" from sandbox to core by meeting the quality bar
(adherence to the conventions, automation, documentation).

### Organize core config around high-level services

A "service" is some collection of config that, when deployed, provides some
reasonably separate value to someone. Each service should have its config
co-located within a repo, would have an admin-deployed part,
automatically-deployed part (triggered by a single entry point each), designated
owner and at least a basic doc.

### Get rid of the central Makefile edited by everyone

Let each "unit" have its own Makefile, and just recursively call into
subdirectories. Avoid doing too much in Makefile and call out to shell/python to
perform actions (use Makefiles just for dependencies). Have two targets required
by convention: `make admin-resources` that need admin rights to apply and `make
resources` that applies everything else. Be able to call these separately for
core services and sandbox services.

### Conventions how to create resources

1. Never rely on `oc project` being set to a correct namespace
2. Only use one way how to specify where to put a resource (either in-manifest,
   or `apply -n $NS`)
3. Only create CMs with config-bootstrapper/updater, never with manifest

## Backlog ideas

- Make image mirroring easier to set up
- Implement namespace TTL
