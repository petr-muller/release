ifeq ($(APPLY),)
	APPLY := oc apply
endif

.PHONY: resources admin-resources

resources:
	make -C openshift-monitoring  APPLY="$(APPLY)" resources

admin-resources:
	make -C openshift-monitoring  APPLY="$(APPLY)" admin-resources
