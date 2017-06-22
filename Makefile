PACKAGES=util channel process processManager benOr collect_data interactive
go_root := $(realpath ${CURDIR}/../..)
SET_PATH = export GOPATH=/usr/lib/go:$(go_root)


.PHONY: build
build: $(PACKAGES)
	$(SET_PATH); cd util; go build
	$(SET_PATH); cd channel; go build
	$(SET_PATH); cd process; go build
	$(SET_PATH); cd processManager; go build
	$(SET_PATH); cd benOr ; go build
	$(SET_PATH); cd collect_data; go build
	$(SET_PATH); cd interactive; go build

.PHONY: install 
install: $(PACKAGES)
	$(SET_PATH); cd util; go install
	$(SET_PATH); cd channel; go install
	$(SET_PATH); cd process; go install
	$(SET_PATH); cd processManager; go install
	$(SET_PATH); cd benOr ; go install
	$(SET_PATH); cd interactive; go install
	$(SET_PATH); cd collect_data; go install

.PHONY: run_interactive 
run_interactive: $(PACKAGES) install
	../../bin/interactive

.PHONY: run_collect_data 
run_collect_data: $(PACKAGES) install
	cd ../../bin/; ./collect_data

.PHONY: clean
clean: 
	rm -rf ../../pkg/*/consensus
	rm -f ../../bin/consensus

