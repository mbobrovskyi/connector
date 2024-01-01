VERSION := $(shell cat VERSION | sed -Ee 's/^v|-.*//')

.PHONY: version

version:
	@echo v$(VERSION)

SEMVER_TYPES := major minor patch
BUMP_TARGETS := $(addprefix bump-,$(SEMVER_TYPES))

.PHONY: $(BUMP_TARGETS)

$(BUMP_TARGETS):
	$(eval bump_type := $(strip $(word 2,$(subst -, ,$@))))
	$(eval f := $(words $(shell a="$(SEMVER_TYPES)";echo $${a/$(bump_type)*/$(bump_type)} )))
	$(eval new_version := $(shell echo $(VERSION) | awk -F. -v OFS=. -v f=$(f) '{ $$f++ } 1'))
	$(eval $(shell echo $(new_version) > VERSION))
	git add .
	git commit -m "v$(new_version)"
	git push
	git tag -a "v$(new_version)" -m "v$(new_version)"
	git push --tag

