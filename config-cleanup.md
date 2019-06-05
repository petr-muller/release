# openshift/release config cleanup

## Configuration inventory

### Cluster profiles
1. `cluster/test-deploy/*/{.type,*.yaml}`: cluster profiles

- **Description:** Cluster profiles
- **Owner:** Cloud platform teams, Clayton
- **Changed:** Monthly
- **How it is applied:** `make prow-cluster-jobs`, config-updater/bootstrapper
- **Namespace:** ci/ci-stg

### Image mirroring
1. `cluster/ci/jobs/image-mirror.yaml`: 5 CronJobs that mirror OpenShift images to Quay
2. `cluster/ci/jobs/{toolchain,knative}-image-mirror.yaml`: 4 CronJobs that mirror knative/toolchain images to Quay
3. `cluster/ci/config/mirroring/*`: Configs for the mirroring jobs

- **Description:** Jobs that mirror images (built by ci-operator) from api.ci to Quay
- **Owner:** Clayton, related product teams
- **Changed:** Monthly? New teams need to set up their mirroring from scratch
- **How it is applied:** `make image-pruner-setup` (cronjobs), config-updater/bootstrapper (config files)
- **Namespace:** ci/ci-stg

### Image Pruner

1. `cluster/ci/jobs/image-pruner.yaml`: A single CronJob that prunes images from the cluster

- **Description:** A single CronJob that prunes images from the cluster
- **Owner:** Clayton
- **Changed:** Not often
- **How it is applied:** Makefile (also creates SA and rolebinding!)
- **Namespace:** ci

### Prow monitoring

1. `cluster/ci/monitoring/*_service_monitor.yaml`:  Additional targets for Prometheus to monitor
2. `cluster/ci/monitoring/msg.tmpl`: I guess a template for Slack alert messages
3. `cluster/ci/monitoring/prometheus_rbac.yaml`: A bunch of RBAC for Prometheus
4. `cluster/ci/monitoring/prow.json`: Not exactly sure, looks like a Grafana JSON config?
5. `cluster/ci/monitoring/mixins/grafana_dashboards/*.jsonnet`: Individial Grafana dashboards?
6. `cluster/ci/monitoring/mixins/jsonnetfile.json`: Not sure.
7. `cluster/ci/monitoring/mixins/prometheus/*`:  Alerts over Prometheus?
8. `cluster/ci/monitoring/prometheus_proxy_secret.yaml`: 1 Secret
9. `cluster/ci/monitoring/*.json`: Actual Grafana dashboards, generated from the jsonnet by `mixins/Makefile`
10. `cluster/ci/monitoring/debug/grafana_datasources.yaml`: Datasource configuration for debug grafana
11. `cluster/ci/monitoring/debug/grafana.ini`: Debug grafana config
12. `cluster/ci/monitoring/debug/grafana_deploy.yaml`: Deployment, Service, Route (everything in `prow-monitoring-stage`)
13. `cluster/ci/monitoring/prometheus_operator_deploy.yaml`: Deployment of prometheus-operator (into `prow-monitoring`)
14. `cluster/ci/monitoring/alertmanager.yaml`: AM config file
15. `cluster/ci/monitoring/alert_manager_rbac.yaml`: ServiceAccount, ClusterRole, ClusterRoleBinding
16. `cluster/ci/monitoring/alert_manager_crd.yaml`: AlertManager spec
17. `cluster/ci/monitoring/ci_hook_rbac.yaml`: Role & RoleBinding
18. `cluster/ci/monitoring/datasources.yaml`: Datasource configuration for grafana
19. `cluster/ci/monitoring/grafana.ini`: grafana config
20. `cluster/ci/monitoring/grafana_deploy.yaml`: Deployment, in `prow-monitoring`
21. `cluster/ci/monitoring/grafana_expose.yaml`: Service, Route, in `prow-monitoring`
22. `cluster/ci/monitoring/alert_manager_expose.yaml`: Service, Route, in `prow-monitoring`
23. `cluster/ci/monitoring/alert_manager_proxy_secret.yaml`: 1 Secret
24. `cluster/ci/monitoring/prometheus_expose.yaml`: Service ,Route, in `prow-monitoring`
25. `cluster/ci/monitoring/prometheus_rule_prow.yaml`: PrometheusRule
26. `cluster/ci/monitoring/dashboards.yaml`: Probably some config file for Grafana
27. `cluster/ci/monitoring/grafana_proxy_secret.yaml`: 1 Secret
28. `cluster/ci/monitoring/client_rbac.yaml`: Role + RoleBinding
29. `cluster/ci/monitoring/grafana_rabc.yaml`: ClusterRole, ServiceAccount, ClusterRoleBinding
30. `cluster/ci/monitoring/prometheus_project.yaml`: Project / namepace (probably rename)
31. `cluster/ci/monitoring/prometheus_crd.yaml`: Prometheus

- **Description:** Monitoring stack for Prow
- **Owner:** DPTP
- **Changed:** Often recently, because we are developing it
- **How it is applied:** `make prow-monitoring` calls into `cluster/ci/monitoring/Makefile`
- **Namespace:** prow-monitoring/ci

### Config needed for ci-operator to work

1. `cluster/ci/config/prow/openshift/ci-operator/roles.yaml`: Assorted roles & bindings for ci-operator to work

- **Description:** Assorted roles & bindings for ci-operator to work
- **Owner:** DPTP, owners of namespaces where ci-operator promotes images
- **Changed:** Often
- **How it is applied:** `make prow-jobs`
- **Namespace:** ci,fabric8-services,kiegroup,ocp,ocp-future,openshift,rhcos,ci,ci-stg

### Stage namespace config
1. `cluster/ci/config/prow/openshift/ci-operator/stage.yaml`: Assorted objects for the `ci-stg` namespace.

- **Description:** Assorted objects for the `ci-stg` namespace.
- **Owner:** DPTP
- **Changed:** Rarely (with Prow bumps)
- **How it is applied:** `make prow-ci-stg-ns`
- **Namespace:** ci/ci-stg

### Prow
1. `cluster/ci/config/prow/openshift/adapter_imagestreams.yaml`: ImageStreams for upstream Prow images
2. `cluster/ci/config/prow/openshift/deck.yaml`: Bunch of expected stuff for Deck in a template
3. `cluster/ci/config/prow/openshift/tide.yaml`: Service and Deployment for Tide
4. `cluster/ci/config/prow/openshift/jenkins_operator.yaml`: 3x Deployment + Service
5. `cluster/ci/config/prow/openshift/hook.yaml`: Route, Service, Deployment
6. `cluster/ci/config/prow/openshift/horologium_rbac.yaml`: ServiceAccount, Role, RB
7. `cluster/ci/config/prow/openshift/artifact-uploader_rbac.yaml`: SA, Role, RB, ClusterRole, CRB
8. `cluster/ci/config/prow/openshift/sinker.yaml`: Deployment, Service
9. `cluster/ci/config/prow/openshift/config_updater_rbac.yaml`: SA, Role, RB (in ci, prow-monitoring, openshift, ci-stg)
10. `cluster/ci/config/prow/openshift/cherrypick.yaml`: Deployment, Service
11. `cluster/ci/config/prow/openshift/needs_rebase.yaml`: Deployment, Service
12. `cluster/ci/config/prow/openshift/tot.yaml`: Deployment, Service, PersistentVolumeClaim
13. `cluster/ci/config/prow/openshift/plank_rbac.yaml`: SA, Role, RB
14. `cluster/ci/config/prow/openshift/refresh.yaml`: Deployment, Service
15. `cluster/ci/config/prow/openshift/statusreconciler.yaml`: Deployment (no service?)
16. `cluster/ci/config/prow/openshift/horologium.yaml`: Deployment in a List (no service)
17. `cluster/ci/config/prow/openshift/sinker_rbac.yaml`: SA, Role, RB
18. `cluster/ci/config/prow/openshift/tracer.yaml`: Route, Service, DeploymentConfig, Secret in a Template
19. `cluster/ci/config/prow/openshift/tide_rbac.yaml`: SA, Role, RB
20. `cluster/ci/config/prow/openshift/plank.yaml`: Deployment, Service
21. `cluster/ci/config/prow/openshift/jenkins_operator_rbac.yaml`: SA, Role, RB
22. `cluster/ci/config/prow/openshift/artifact-uploader.yaml`: Deployment in a List
23. `cluster/ci/config/prow/openshift/statusreconciler_rbac.yaml`: SA, Role, RB
24. `cluster/ci/config/prow/openshift/hook_rbac.yaml`: SA, Role, RB
25. `cluster/ci/config/prow/openshift/tracer_rbac.yaml`: SA, Role, RB, ClusterRole, CRB
26. `cluster/ci/config/prow/openshift/deck_rbac.yaml`: SA, Role, RB, ClusterRole, CRB67
27. `cluster/ci/config/prow/openshift/ghproxy.yaml`: Service, Deployment, PersistentVolumeClaim
28. `cluster/ci/config/prow/deck/extensions/*`: Additional Deck config
29. `cluster/ci/config/prow/prow_crd.yaml`: ProwJob CRD
30. `cluster/ci/config/prow/prowjob_access.yaml`: RBAC allowing ci-admins to manipulate ProwJobs

- **Description:** Prow components
- **Owner:** DPTP
- **Changed:** Bumps often, rarely otherwise
- **How it is applied:** `make prow-rbac; make prow-services`
- **Namespace:** ci, config_updater_rbac in stg, openshift and prow-monitoring too

### Prow configuration
1. `cluster/ci/config/prow/config.yaml`: Prow component configuration
2. `cluster/ci/config/prow/plugins.yaml`: Prow plugin configuration
3. `cluster/ci/config/prow/labels.yaml`: Prow labels configuration

- **Description:** Prow config
- **Owner:** DPTP
- **Changed:** Often
- **How it is applied:** `make prow-config`, `make prow-config-update`, config-updater/bootstrapper (to cluster), infra jobs (to GH)
- **Namespace:** ci

### Logging config
1. `cluster/ci/config/logging/*.yaml`: DaemonSet, CM

- **Description:** Configuration for Prow logs to be synced to Stackdriver
- **Owner:** DPTP
- **Changed:** No
- **How it is applied:** `make logging`
- **Namespace:** kube-system

### OAuth error page
1. `cluster/ci/config/page-templates/oauth-error.html`

- **Description:** Customized login error page with docs link
- **Owner:** DPTP
- **Changed:** No
- **How it is applied:** Manually
- **Namespace:** N/A

### Secret mirroring config
1. `cluster/ci/config/secret-mirroring/mapping.yaml`

- **Description:** Config file used by *something?* to mirror secrets between namespaces
- **Owner:** DPTP
- **Changed:** Monthly
- **How it is applied:** `make prow-secrets` + config-bootstrapper/updater
- **Namespace:** ci

### Various roles

1. `cluster/ci/config/roles.yaml`: Assortment of Namespaces, Groups, Roles and Bindings. Kitchen sink, should be broken down?

- **Description:** Assortment of Namespaces, Groups, Roles and Bindings. Kitchen sink, should be broken down?
- **Owner:** DPTP, Clayton
- **Changed:**  Often
- **How it is applied:** `make cluster-roles`
- **Namespace:** Many

### Some general cluster config
1. `cluster/ci/config/web-console-oauthclient.yaml`: OAuthClient
2. `cluster/ci/config/service-ca.yaml`: Not exactly sure (is present on api.ci)
3. `cluster/ci/config/metrics-server.yaml`: Not exactly sure (is NOT present on api.ci)
4. `cluster/ci/config/cluster-autoscaler.yaml`: Autoscaler (is present on api.ci)
5. `cluster/ci/config/origin-web-console.yaml`: Console, but does not seem to be present on the cluster?
6. `cluster/ci/config/gce-pd-storageclass.yaml`: StorageClass (is NOT present on api.ci)

- **Description:** ?
- **Owner:** ?
- **Changed:**  Rarely?
- **How it is applied:** `make deploy` in `cluster/ci/Makefile`. The target does more arcane stuff.
- **Namespace:** Many

### No idea
1. `cluster/ci/config/prow/openshift/rpm-mirrors/docker.repo`: DNF repo with dockertested? Appears unused.

### Azure
1. `projects/azure/rbac.yaml`: Bunch of RBAC in `azure` namespace
2. `projects/azure/cluster-wide.yaml`: Group
3. `projects/azure/base-images/plugin-base.yaml`: ImageStream
4. `projects/azure/base-images/ci-base.yaml`: ImageStream & BuildConfig
5. `projects/azure/base-images/test-base.yaml`: ImageStream & BuildConfig
6. `projects/azure/token-refresh/token-refresh.yaml`: SA, CM and CronJob
7. `projects/azure/azure-team-rbac.yaml`: Self-management RBAC (allow them to manage their group)
8. `projects/azure/image-mirror/image-mirror.yaml`: CM + CronJob, like the "official" ones
9. `projects/azure/azure-purge/cronjob.yaml`: CronJob

- **Description:** Bunch of custom stuff by the Azure team
- **Owner:** Azure team
- **Changed:** Fairly often
- **How it is applied:** `make azure`
- **Namespace:** azure, ci, ci-stg

### Cincinnati
1. `projects/cincinnati/cincinnati.yaml`: ImageStream, BuildConfig

- **Description:** Apparently just takes care of building Cincinnati images.
- **Owner:** steveeJ
- **Changed:** often, new
- **How it is applied:** `make cincinnati`
- **Namespace:** `ci` (but also present in `cincinnati-ci` which we do not track in the repo)

### cluster-operator
1. `projects/cluster-operator/cluster-operator-team-roles.yaml`: RBAC for some OpenShift people to do stuff
2. `projects/cluster-operator/cluster-operator-roles-template.yaml`: RBAC for operating cluster-operator

- **Description:** Not sure
- **Owner:** dgoodwin
- **Changed:** year ago
- **How it is applied:** `make roles -> cluster-operator-roles`
- **Namespace:** N/A

### content-mirror
1. `projects/content-mirror/pipeline.yaml` ImageStream, BuildConfig

- **Description:** Builds https://github.com/openshift/content-mirror?
- **Owner:** Clayton
- **Changed:** last year
- **How it is applied:** `make content-mirror`
- **Namespace:** ci

### gc-daemonset
1. `projects/gc-daemonset/daemonset.yaml`: DaemonSet

- **Description:** Genius DaemonSet that prevents images from being GCed
- **Owner:** DPTP (Steve)
- **Changed:** Rarely
- **How it is applied:** isn't
- **Namespace:** ci

### gcsweb
1. `projects/gcsweb/pipeline.yaml: ImageStream, 2x BuildConfig, DeploymentConfig, Service, Route

- **Description:** Runs https://gcsweb-ci.svc.ci.openshift.org if I understand correctly
- **Owner:** DPTP
- **Changed:** last year
- **How it is applied:** `make gcsweb`
- **Namespace:** ci

### kube-state-metrics
1. `projects/kube-state-metrics/pipeline.yaml` ImageStream, 2x BuildConfig

- **Description:** Builds https://github.com/kubernetes/kube-state-metrics.git
- **Owner:** Clayton
- **Changed:** 2 years ago
- **How it is applied:** `make kube-state-metrics`
- **Namespace:** ci

### kubernetes
1. `projects/kubernetes/node-problem-detector.yaml`: ImageStream, 2x BuildConfig

- **Description:** Builds https://github.com/openshift/node-problem-detector.git
- **Owner:** Clayton
- **Changed:** last year
- **How it is applied:** `make node-problem-detector`
- **Namespace:** ci

### libpod
1. `projects/libpod/libpod.yaml`: ImageStream, 2x BuildConfig

- **Description:** Builds Fedora & CentOS buildroots for libpod
- **Owner:** baude
- **Changed:** infrequently
- **How it is applied:** `make libpod`
- **Namespace:** ci

### metering
1. `projects/metering/project.yaml`: Project
2. `projects/metering/rbac.yaml`: Simple RBAC
3. `projects/metering/group.yaml`: Group

- **Description:** Nicely looking custom project, with a Makefile, OWNERS and a setup script.
- **Owner:**  chancez
- **Changed:** 3m ago
- **How it is applied:** `make metering`
- **Namespace:** metering

### oauth-proxy
1. `projects/oauth-proxy/pipeline.yaml`: ImageStream, 2x BuildConfig

- **Description:** Builds https://github.com/openshift/oauth-proxy.git
- **Owner:** simo5
- **Changed:** last year
- **How it is applied:** `make oauth-proxy` (not in `make projects` though)
- **Namespace:** ci

### origin-release

1. `projects/origin-release/dashboards-validation/`: ImageStream, BuildConfig, Dockerfile
2. `projects/origin-release/golang-1.*`: RPM, Specfile, DNF repo, Dockerfile
3. `projects/origin-release/pipeline.yaml`: ImageStream, BuildConfigs for 5x golang, 2x nodejs
4. `projects/origin-release/nodejs-8*`: 2x Dockerfile
5. `projects/origin-release/python-validation`: ImageStream, BuildConfig, Dockerfile
5. `projects/origin-release/cli.yaml`: BuildConfig. Present on cluster but is not referred from anywhere?

- **Description:** Kitchen-sinky, but looks like *we* caused it with the validation images
- **Owner:** Unclear
- **Changed:** Often
- **How it is applied:** `make origin-release` but also `make python-validation` and `make build-dashboards-validation-image`
- **Namespace:** ci, ocp

### origin-stable
1. `projects/origin-stable/*.yaml`: ImageStreams

- **Description:** Little unclear
- **Owner:** Clayton
- **Changed:** Infrequent
- **How it is applied:** `make origin-stable`
- **Namespace:** openshift

### Prometheus
1. `projects/prometheus/*`

- **Description:** Bunch of monitoring stuff by Michalis, applied, but I wonder if we use it.
- **Owner:** ? (Michalis)
- **Changed:** last year
- **How it is applied:** `make prometheus`
- **Namespace:** ci

### Publishing bot
1. `projects/publishing-bot/storage-class.yaml`

- **Description:**  Single StorageClass. Old, unsure if used.
- **Owner:** sttts
- **Changed:** last year
- **How it is applied:** `make publishing-bot`
- **Namespace:** N/A

### Service idler
1. `projects/service-idler/pipeline.yaml`: ImageStream, 2x BuildConfig.

- **Description:** Old, present in the cluster, unsure if used.
- **Owner:** DirectXMan12
- **Changed:** last year
- **How it is applied:** `make service-idler` (but not in `make projects`)
- **Namespace:** ci

### Projects
1. `projects/binary-pipeline.yaml`

- **Description:** Looks like a template for other `pipeline.yaml` projects
- **Changed:** 2y ago
- **How it is applied:** N/A
- **Namespace:** N/A

### Jobs

1. `ci-operator/jobs/$org/$component/*.yaml`: Component-specific CI jobs
2. `ci-operator/jobs/openshift/release/*periodics.yaml`: Pseudo-periodic jobs backing cluster-bot, release-controller, etc. Config updates. Fast-forward job
3. `ci-operator/jobs/infra-periodics.yaml`: label-sync, org-sync, branchprotector, retester, stale-issue labeller

- **Description:** All Prow jobs
- **Owner:** Component owners. DPTP. Clayton
- **Changed:** often
- **How it is applied:** config-updater/bootstrapper into CMs
- **Namespace:** ci

### ciop config
1. `ci-operator/config/$org/$repo/*.yaml`

- **Description:**: Component-specific CI config
- **Owner:** Component owners
- **Changed:** Often
- **How it is applied:** config-updater/bootstrapper into CMs
- **Namespace:** ci

### ci-chat-bot
1. `ci-operator/infra/openshift/ci-chat-bot/deploy.yaml`: RBAC (Role, Binding, SA) + Deployment

- **Description:** Service installing clusters via Slack bot
- **Owner:** Clayton
- **Changed:** Rarely
- **How it is applied:** `make prow-ci-chat-bot`
- **Namespace:** ci

### ci-search
1. `ci-operator/infra/openshift/ci-search/config.yaml`: Config *for* ci-search
2. `ci-operator/infra/openshift/ci-search/deploy.yaml`: Role, Binding, ImageStream, Service, Route, Deployment

- **Description:** https://search.svc.ci.openshift.org/
- **Owner:** Clayton
- **Changed:** Infrequently
- **How it is applied:** `make prow-ci-search`
- **Namespace:** ci-search (we also have `ci-search-next` which is not present in openshift/release)

### infra/openshift/origin
1. `ci-operator/infra/openshift/origin/*.yaml`

- **Description:**: Ad-hoc servers (Deployments with inline scripts) in the `ci-rpms` namespace
- **Owner:** Clayton?
- **Changed:** Bi-weekly
- **How it is applied:** `make prow-artifacts`
- **Namespace:** ci

### release-controller
1. `ci-operator/infra/openshift/release-controller/deploy-*.yaml`: RBAC, Deployment, Service, Routes
2. `ci-operator/infra/openshift/release-controller/images-*.yaml`: ImageStreams
3. `ci-operator/infra/openshift/release-controller/releases/release-*.json`: configs?
4. `ci-operator/infra/openshift/release-controller/rpms-*.yaml`: Service + Deployment in each, serve some RPMs
5. `ci-operator/infra/openshift/release-controller/repos/*.repo`: DNF repos

- **Description:** release-controllers (like https://openshift-release.svc.ci.openshift.org/)
- **Owner:** Clayton
- **Changed:** Quite often
- **How it is applied:** `make prow-release-controller`
- **Namespace:** ci,ocp,origin,openshift,ci-release

### misc infra
1. `ci-operator/infra/ansible-*-imagestream.yaml`: 2x ImageStream

- **Description:** Not sure
- **Owner:** shawn-hurley/jmontleon
- **Changed:** Infrequently
- **How it is applied:** `make ci-infra-imagestreams` (applies just runner, not operator)
- **Namespace:** openshift

###
1. `ci-operator/infra/src-cache-origin.yaml`: ImageStream, 4x BuildConfig

- **Description:** Not sure
- **Owner:** Clayton
- **Changed:** Infrequently
- **How it is applied:** both `make ci-operator-config` and `make prow-artifacts`
- **Namespace:** ci

### ci-operator templates
1. `ci-operator/templates/openshift/installer/*.yaml`: ci-operator templates for the new installer
2. `ci-operator/templates/openshift/openshift-ansible/*.yaml`: ci-operator templates for the old installer
3. `ci-operator/templates/openshift/openshift-azure/*.yaml`:  ci-operator templates for azure
4. `ci-operator/templates/os.yaml`: SA, 2x Role+Binding. Not sure why this is in `templates`
5. `ci-operator/templates/master-sidecar-{3,4}.yaml`: Additional ci-operator templates

- **Description:** Stuff to support ci-operator template jobs
- **Owner:** Installer & platform teams
- **Changed:** often
- **How it is applied:** `make prow-cluster-jobs` (in ci-stg) except for `os.yaml` which is applied in `make prow-jobs`. Also, CMs are updated by config-bootstrapper/updater
- **Namespace:** ci, ci-stg

### tools/build
1. `tools/build/*`: Bunch of stuff that I don't know and which looks deprecated.

- **Description:** ?
- **Owner:** Steve?
- **Changed:** last year
- **How it is applied:** ?
- **Namespace:** ?

## Apply script inventory

1. `Makefile`: Basic entrypoint, kitchen-sinky, random people add their stuff, little system. Funky NS selection.
2. `cluster/test-deploy/Makefile`: Looks deprecated (uses nonexistent files), does not seem to be interacting with api.ci
3. `cluster/ci/Makefile`: Looks partially deprecated, partially suspicious.

### Monitoring

1. `cluster/ci/monitoring/mixins/Makefile`: Does not interact with the cluster, regenerates files
2. `cluster/ci/monitoring/Makefile`: Creates/deletes monitoring-related object

### Metering

1. `projects/metering/setup-metering-project.sh`: Actually apllies the setup
2. `projects/metering/Makefile`: Simply calls `setup-metering-project.sh`

### Random scripts

1. `hack/config-bootstrapper.sh`: just calls config-bootstrapper, does not seem to be used from anywhere
2. `ci-operator/populate-secrets-from-bitwarden.sh`: creates secrets

### Jobs that apply some part of the config
1. `ci-operator/jobs/openshift/release/openshift-release-periodics.yaml`: make postsubmit-update, direct config-bootstrapper
2. `ci-operator/jobs/openshift/release/openshift-release-master-postsubmits.yaml`: make postsubmit-update, label-sync
3. `ci-operator/jobs/infra-periodics.yaml`: label-sync, branch-protector sync, org sync
