# Copyright 2025 The ChaosBlade Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: help format verify

help:
	@echo "Available make targets:"
	@echo "  format   Run update-gofmt.sh and update-imports.sh to automatically fix Go code formatting and import order."
	@echo "  verify   Run verify-gofmt.sh and verify-imports.sh to check Go code formatting and import order."
	@echo "  license-check   Check for proper license headers in source files."
	@echo "  help     Show this help message."

format: license-format
	@echo "Running goimports and gofumpt to format Go code..."
	@./hack/update-imports.sh
	@./hack/update-gofmt.sh

verify:
	@echo "Verifying Go code formatting and import order..."
	@./hack/verify-gofmt.sh
	@./hack/verify-imports.sh

.PHONY: license-check
license-check:
	@echo "Checking license headers..."
	docker run -it --rm -v $(shell pwd):/github/workspace ghcr.io/korandoru/hawkeye check

.PHONY: license-format
license-format:
	@echo "Formatting license headers..."
	docker run -it --rm -v $(shell pwd):/github/workspace ghcr.io/korandoru/hawkeye format