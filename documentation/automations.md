# OpenShift CI Automation
- _Periodically, something does something [to achieve something]_
- Each minute, all imagestreams on build01 import latest images from api.ci: [job](https://deck-ci.apps.ci.l2s4.p1.openshiftapps.com/?job=periodic-ci-image-import-to-build01) | [def](https://github.com/openshift/release/blob/92215e080639899e2f2936d2db1aa35828332b07/ci-operator/jobs/infra-periodics.yaml#L2-L28)
