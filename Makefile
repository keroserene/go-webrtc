.SILENT: check
.PHONY: check
check:
	! gofmt -l `find . -path ./third_party -prune -o -name '*.go' -print` 2>&1 | read
	go vet .

.PHONY: lint
lint:
	find . -name "*.cc" -o -name "*.hpp" -o -name "*.h" -maxdepth 1 | xargs cpplint.py --extensions=h,hpp,cc
