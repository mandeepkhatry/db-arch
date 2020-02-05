genparser:
	@echo "***Generating parse.go file***"
	@pigeon -o server/internal/engine/parser/parse.go server/internal/engine/parser/parse.peg

dependencies:
	@echo "***Checking go dependencies***"
	@go mod tidy

build: genparser dependencies


