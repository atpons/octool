.PHONY: credit
credit:
	gocredits . > CREDITS

.PHONY: build
build: credit
	bazel build
