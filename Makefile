.PHONY: help format verify

help:
	@echo "Available make targets:"
	@echo "  format   Run update-gofmt.sh and update-imports.sh to automatically fix Go code formatting and import order."
	@echo "  verify   Run verify-gofmt.sh and verify-imports.sh to check Go code formatting and import order."
	@echo "  help     Show this help message."

format:
	@echo "Running goimports and gofumpt to format Go code..."
	@./hack/update-imports.sh
	@./hack/update-gofmt.sh

verify:
	@echo "Verifying Go code formatting and import order..."
	@./hack/verify-gofmt.sh
	@./hack/verify-imports.sh

