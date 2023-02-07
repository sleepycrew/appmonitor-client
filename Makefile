.PHONY: gen-code list

list:
	@echo Tasks:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

gen-code: transform-spec

transform-spec: gen-spec
	./scripts/remove-marshal.sh


gen-spec:
	go run ./scripts/generate-types.go
