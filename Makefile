run:
	@go build -o .build/paybazar_debug cmd/*.go && ./.build/paybazar_debug

build:
	@go build -o .build/paybazar_release cmd/*.go