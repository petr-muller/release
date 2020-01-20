# Things

A personal index of things.

## OpenShift CI

### api.ci

### ci/prow/correctly-sharded-config

### DPTP

### openshift/release

### openshift-release-master-config-bootstrapper

## Software

### config-bootstrapper

### config-updater

A [Prow](#prow) [plugin](#prow-plugins) that updates ConfigMaps in clusters when
files change in a Git repository. It is used in OpenShift CI to deploy
configuration changes to production after they are merged to [openshift/release](#openshiftrelease).
It is configured by glob patterns that specify which files are supposed to be
updated in which ConfigMaps. It is triggered by a merge to a Git repository.
After it updates ConfigMaps, it posts a comment to the Pull Request whose merge
triggered it.

![config-updater GitHub comment example](./images/config-updater.png)

An example bot comment informing that ConfigMaps were updated after a merge.

**NOTES**

1. A common pitfall while configuring config-updater is that its configuration
   cannot be changed in a same PR as the content it is supposed to cover. For
   example, when adding new configuration file that is not yet covered by any
   glob pattern, the change to add a glob pattern to config-updater configuration
   must be done first in a separate pull request, otherwise the new files will
   not be deployed to the cluster. The reason for this that each pull request is
   processed by config-updater running with *existing* configuration, not the
   one changed in the pull request.
2. Changes to config-updater configuration in [openshift/release](#openshiftrelease)
   repository are tested by [ci/prow/correctly-sharded-config](#ciprowcorrectly-sharded-config)
   job.

**LINKS**

[README](https://github.com/kubernetes/test-infra/blob/master/prow/plugins/updateconfig/README.md) | [code](https://github.com/kubernetes/test-infra/tree/master/prow/plugins/updateconfig) | [OpenShift CI Configuration](https://github.com/openshift/release/blob/4aa6efe87ae360f63c5f724cb47433da2d979da8/core-services/prow/02_config/_plugins.yaml#L170)

**SEE ALSO**

- [config-bootstrapper](#config-bootstrapper)
- [openshift-release-master-config-bootstrapper](#openshift-release-master-config-bootstrapper)

### Deck

### Prow

Prow is a Kubernetes-based CI/CD system. It is used both in Kubernetes upstream
community and in OpenShift. Its main purpose is to execute containerized testing
workloads and to automate and enforce various tasks and development workflow
steps, primarily in GitHub. Prow consists of several small microservices running
on a Kubernetes cluster. The OpenShift CI instance runs on the
[api.ci](#apici) cluster and is operated by [DPTP](#dptp).

**LINKS**

[upstream](https://github.com/kubernetes/test-infra/tree/master/prow) | [Deck frontend to OpenShift CI instance of Prow](https://prow.svc.ci.openshift.org/)

### Prow Plugins

### Prow

### Prow Plugins
