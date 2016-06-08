PACKAGES = util channel process processManager main


.PHONY: build
build: $(PACKAGES) environment
	cd util; go build
	cd channel; go build
	cd process; go build
	cd processManager; go build
	cd main; go build

.PHONY: install environment
install: $(PACKAGES)
	cd util; go install
	cd channel; go install
	cd process; go install
	cd processManager; go install
	cd main; go install

.PHONY: test environment
test: $(PACKAGES) install
	cd main; go test

.PHONY: clean
clean: 
	rm -rf ../../pkg/*/consensus
	rm -f ../../bin/consensus
	
