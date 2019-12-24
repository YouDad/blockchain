all: cc install
	@echo [all]

cc:
	@echo [cc]
	go build main.go

install:
	@echo [install]
	mv main $(GOPATH)/bin/blockchain
