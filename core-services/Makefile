ifeq ($(APPLY),)
	APPLY := $(CURDIR)/../hack/apply-config.py
endif

SERVICES = openshift-monitoring

.PHONY: resources admin-resources check

check:
	_hack/validate-core-service.sh _TEMPLATE # keep the template valid
	$(foreach service,$(SERVICES),_hack/validate-core-service.sh $(service);)

resources:
	$(foreach service,$(SERVICES),$(APPLY) $(service);)

admin-resources:
	$(foreach service,$(SERVICES),$(APPLY) $(service);)
