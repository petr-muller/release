# Experimental changes

```
+-- Makefile
+-- core
|   +-- Makefile
|   +-- openshift-monitoring
|   |   +-- Makefile
|   |   +-- deploy.sh
|   |   +-- resource.yaml
|   |   +-- OWNERS
|   |   \-- README.md
|   |
|   \-- _TEMPLATE
|       \-- ...
|
\-- sandbox
```

#### Makefile

Two targets applying core services config: `core` and `core-admin`, and two
associated targets for dry-run checking the config. The targets set `$(APPLY)`
to allow running with `--dry-run=true` or `--as system:admin`. This root
`Makefile` should not be needed to change often. Still need to find good
entrypoints for sandbox project.

#### core/

Directory holding a subdir for each "thing" of core services. We may reuse
`cluster/ci` or `cluster/ci/config` if we don't want to have a new directory,
or we may find a better name.

#### core/Makefile

Just mux the `core` and `core-admin` targets to individual subdirs. We should
only edit this target when we add new "thing", or we can somehow iterate and do
not edit this file at all (or iterate in root Makefile and get rid of this file
entirely)

#### core/openshift-monitoring/Makefile

Thing-specific deployment recipe. Needs to have (even empty) `resources` and
`admin-resources`, need to respect `$(APPLY)` and export it to children
processes. More complex procedures should be done in scripts and simply called
from Makefile. The destination namespace should be passed as `-n` param to
`$(APPLY)` (not defined in resource manifests, or to be relied on current
namespace).

#### core/openshift-monitoring/deploy.sh

(Optional) Script implementing more complicated deploy procedures. Called by
Makefile, should respect `$(APPLY)`.

#### core/openshift-monitoring/resource.yaml

Kubernetes/OpenShift resource manifests. Should not define `namespace` in
metadata.

#### core/openshift-monitoring/{OWNERS,README.md}

Should be mandatory for core services, even having OWNERS for DPTP-owned
services (explicit is better).

#### core/_TEMPLATE

Copy-able template for a core service directory

#### sandbox

TBD. Structure similar to `core` or current `projects`, but with fewer imposed
rules and weaker structure. We should definitely require OWNERS. Automated
updates / validation should not be required, but should be allowed for opt-in,
using a mechanism consistent with core services (provide a Makefile, respect
`$(APPLY)`). These should be separate jobs from those operating on core
services and should not block merging/updates of core service config.

Interesting part of sandbox projects are admin-created resources (namespaces,
RBACs...). I think these should not be defined within the sandbox projects, but
we should have a core service where we would keep these. This way we would
avoid using `--as system:admin` on resources in `sandbox`
