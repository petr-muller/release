#!/bin/bash

set -o errexit

CACHE="$(mktemp -d)"

while true; do
  ./deck-build-log-pull
  ./deck-build-log-plot \
    "release .* failed: pod .* was already deleted" \
    "waiting for Kubernetes API: context deadline exceeded" \
    "failed to initialize the cluster:.*Cluster operator .* is still updating" \
    "Container setup in pod .* failed" \
    "failed:.*Multi-AZ Clusters should spread the pods of a replication controller across zones" \
    "failed:.*In-tree Volumes.*subPath should support existing single file" \
    "failed:.*Prometheus when installed on the cluster should start and expose a secured proxy and unsecured metrics" \
    "failed:.*Basic StatefulSet functionality .* should perform rolling updates and roll backs of template modifications with PVCs" \
    "Container test in pod .* failed"
  CURRENT="${CACHE}/$(date -Iminutes -u).svg"
  mv deck-build-log.svg "${CURRENT}"
  chromium-browser "${CURRENT}"
  echo "== ${CURRENT}"
  sleep 600
done
