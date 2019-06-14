APPLY := oc apply

.PHONY: check-core-admin core-admin check-core core

check-core-admin:
	make core-admin APPLY="$(APPLY) --dry-run=true"

core-admin:
	make -C core admin-resources APPLY="$(APPLY) --as=system:admin"

check-core:
	make core APPLY="$(APPLY) --dry-run=true"

core:
	make -C core resources APPLY="$(APPLY)"
