#!/bin/bash

${APPLY} -f prometheus-prow-rules_prometheusrule.yaml -n openshift-monitoring
