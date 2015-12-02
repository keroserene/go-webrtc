.SILENT: check
.PHONY: check
check:
	! gofmt -l . 2>&1 | read
	go vet
