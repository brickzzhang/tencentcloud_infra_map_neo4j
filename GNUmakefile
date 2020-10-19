DIRS=$(shell ls -1 -F | grep "/$$" | grep -v vendor)

all: fmt

fmt:
	@echo "==> Fixing source code with gofmt..."
	@for dir in $(DIRS) ; do `goimports -w $$dir` ; done
	@for dir in $(DIRS) ; do `gofmt -s -w $$dir` ; done
	@echo "==> Fixing source code with gofmt done"

.PHONY: all
