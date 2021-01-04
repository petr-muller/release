# OpenShift CI Automation
- _Periodically, something does something [to achieve something]_
- Each 30 minutes, [Peribolos](https://github.com/kubernetes/test-infra/tree/master/prow/cmd/peribolos)
  reconciles GitHub organization metadata, membership, teams and repositories of
  [openshift](https://github.com/openshift/) and [openshift-priv](https://github.com/openshift-priv)
  organizations:
  [job](https://prow.ci.openshift.org/?job=periodic-org-sync) |
  [def](https://github.com/openshift/release/blob/6f2025056ed9d620816a4dd31dbaa0865a645f45/ci-operator/jobs/infra-periodics.yaml#L726-L772) |
  [config](https://github.com/openshift/config)
- Every morning at 01:00 AM, `manage-clonerefs` build is started, rebuilding the `clonerefs` image
  consumed by `ci-operator`:
  [def](https://github.com/openshift/release/blob/a08dee5d3fcd7b8735cb834884ac19711f22257a/ci-operator/jobs/infra-periodics.yaml#L2-L23) |
  [BuildConfig](https://console.svc.ci.openshift.org/k8s/ns/ci/buildconfigs/manage-clonerefs)

## GitHub Automation
- Each 12 minutes, the [commenter tool](https://github.com/kubernetes/test-infra/tree/master/robots/commenter)
  queries selected repos for PRs that are reviewed and approved but fail tests,
  selects random five and comments with `/retest`:
  [job](https://prow.ci.openshift.org/?job=periodic-retester) |
  [def](https://github.com/openshift/release/blob/07e7635a82665b4ffa85ab536fe08c886d76abbd/ci-operator/jobs/infra-periodics.yaml#L160-L212)
- Every six hours, the [commenter tool](https://github.com/kubernetes/test-infra/tree/master/robots/commenter) queries
  selected repos for PRs and issues that are rotten and did not see any activity for 30 days, and closes them via `/close`:
  [job](https://prow.ci.openshift.org/?job=periodic-issue-close) |
  [def](https://github.com/openshift/release/blob/ededb5ef15e3386bd82ddb5dcc327972e1059104/ci-operator/jobs/infra-periodics.yaml#L180-L224)
- Every six hours, the [commenter tool](https://github.com/kubernetes/test-infra/tree/master/robots/commenter) queries
  selected repos for PRs and issues that are stale and were not updated for last 30 days and marks them as rotten via `/lifecycle rotten`.:
  [job](https://prow.ci.openshift.org/?job=periodic-issue-rotten) |
  [def](https://github.com/openshift/release/blob/5ee2cd373314273f0be04dec82fa842c2c36c178/ci-operator/jobs/infra-periodics.yaml#L225-L273)
- Every day, the [commenter tool](https://github.com/kubernetes/test-infra/tree/master/robots/commenter)
  queries repos for PRs that would merge if they were not blocked by referring
  to an invalid bug, and posts `/bugzilla refresh` to re-validate the linked
  bug:
   [job](https://prow.ci.openshift.org/?job=periodic-daily-bugzilla-refresh) |
   [def](https://github.com/openshift/release/blob/b4a57433e9181d135c9e22c5eca87e60fbcc2cc8/ci-operator/jobs/infra-periodics.yaml#L62-L105)

## Cluster Maintenance

- Every day, `build01` cluster is upgraded to the most recent version via
  `oc adm upgrade --to-latest`:
   [job](https://prow.ci.openshift.org/?job=periodic-build01-upgrade) |
   [definition](https://github.com/openshift/release/blob/master/ci-operator/jobs/infra-periodics.yaml#L2-L23)
