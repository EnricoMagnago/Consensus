PACKAGES=util channel process processManager benOr collect_data interactive

.PHONY: environment
environment:
	export GOPATH=/usr/lib/go:$HOME/Projects/GoProjects

.PHONY: build
build: $(PACKAGES) environment
	cd util; go build
	cd channel; go build
	cd process; go build
	cd processManager; go build
	cd benOr ; go build
	cd collect_data; go build
	cd interactive; go build

.PHONY: install environment
install: $(PACKAGES)
	cd util; go install
	cd channel; go install
	cd process; go install
	cd processManager; go install
	cd benOr ; go install
	cd interactive; go install
	cd collect_data; go install

.PHONY: run_interactive environment
run_interactive: $(PACKAGES) install
	../../bin/interactive

.PHONY: run_collect_data environment
run_collect_data: $(PACKAGES) install
	../../bin/collect_data

.PHONY: clean
clean: 
	rm -rf ../../pkg/*/consensus
	rm -f ../../bin/consensus

