PACKAGES = util channel process processManager main

.PHONY: build
build: $(PACKAGES)
	cd util; go build
	cd channel; go build
	cd process; go build
	cd processManager; go build

.PHONY: install
install: $(PACKAGES)
	cd util; go install
	cd channel; go install
	cd process; go install
	cd processManager; go install

.PHONY: test
test: $(PACKAGES) install
	cd main; go test
	
