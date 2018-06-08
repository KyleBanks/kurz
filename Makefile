default: | install example-local

install:
	@go install -v ./cmd/kurz
.PHONY: install

help: | install
	@kurz -h
.PHONY: help

example-local: | install
	@kurz ./README.md
.PHONY: example

example-git: | install
	@kurz github.com/KyleBanks/modoc
.PHONY: example
