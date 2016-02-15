.SILENT: check
.PHONY: check
check:
	! gofmt -l `find . -path ./third_party -prune -o -name '*.go' -print` 2>&1 | read
	go vet .
