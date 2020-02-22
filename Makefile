all: cc install
	@echo [all]

cc:
	@echo [cc]
	go build main.go

install:
	@echo [install]
	mv main $(GOPATH)/bin/blockchain

port := 10000

test: cc install test_echo test_body test_clean
	@echo [test] finish

test_echo:
	@echo [test] start

test_clean:
	@echo [clean]
	rm -f blockchain$(port).db
	rm -f wallet$(port).dat

test_body:
	-./test.sh
