# OpenShift CI Automation
- _Periodically, something does something [to achieve something]_
- Each 30m minutes, [Peribolos](https://github.com/kubernetes/test-infra/tree/master/prow/cmd/peribolos) reconciles GitHub organization metadata, membership, teams and repositories of [openshift](https://github.com/openshift/) and [openshift-priv](https://github.com/openshift-priv) organizations: [job](https://prow.ci.openshift.org/?job=periodic-org-sync) | [def](https://github.com/openshift/release/blob/6f2025056ed9d620816a4dd31dbaa0865a645f45/ci-operator/jobs/infra-periodics.yaml#L726-L772) | [config](https://github.com/openshift/config)
- Each 12 minutes, the [commenter tool](https://github.com/kubernetes/test-infra/tree/master/robots/commenter) queries selected repos for PRs that are reviewed and approved but fail tests, selects random five and comments with `/retest`: [job](https://prow.ci.openshift.org/?job=periodic-retester) | [def](https://github.com/openshift/release/blob/07e7635a82665b4ffa85ab536fe08c886d76abbd/ci-operator/jobs/infra-periodics.yaml#L160-L212)

## Cluster Maintenance

- Every day, `build01` cluster is upgraded to the most recent version via `oc adm upgrade --to-latest`: [job](https://prow.ci.openshift.org/?job=periodic-build01-upgrade) | [definition](https://github.com/openshift/release/blob/master/ci-operator/jobs/infra-periodics.yaml#L2-L23)
