.PHONY: build clean test run-list run-match run-generate install help

BINARY_NAME=resume-cli
RESUME_PATH=../resume.yaml

help:
	@echo "Resume CLI - Makefile targets:"
	@echo "  build          - Build the binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  install        - Install binary to GOPATH/bin"
	@echo "  run-list       - List all resume items"
	@echo "  run-match      - Match against sample backend job (requires API key)"
	@echo "  run-generate   - Generate sample resume PDF"
	@echo "  watch          - Start docker compose with auto-rebuild on file changes"

build:
	go build -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)
	rm -f *.pdf *.tex *.aux *.log *.out

test:
	go test ./...

install:
	go install .

run-list: build
	./$(BINARY_NAME) list -r $(RESUME_PATH)

run-match: build
	./$(BINARY_NAME) match -f examples/job-backend-go.txt -r $(RESUME_PATH)

run-generate: build
	./$(BINARY_NAME) generate -o sample-resume.pdf -r $(RESUME_PATH)

run-generate-filtered: build
	./$(BINARY_NAME) generate -o filtered-resume.pdf --tags go,distributed-systems,aws -r $(RESUME_PATH)

watch:
	cd .. && docker compose watch
