# Things

A personal index of things.

## OpenShift CI

### openshift/release

### openshift-release-master-config-bootstrapper

### ci/prow/correctly-sharded-config

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

### Prow

### Prow Plugins
