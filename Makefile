all: cc install
	@echo [all]

cc:
	@echo [cc]
	go build main.go

install:
	@echo [install]
	mv main $(GOPATH)/bin/blockchain

port := 10000

clean:
	@echo [clean]
	rm -f blockchain$(port).db
	rm -f wallet$(port).dat

test: cc install test_echo clean test_body
	@echo [test] finish

test_echo:
	@echo [test] start

test_body:
	@-./test.sh 2>&1 | \
		ag --passthrough --color --color-match "4;31" "\[(ERROR|FAIL)\]" | \
		ag --passthrough --color --color-match "4;34" "\[(INFO)\]" | \
		ag --passthrough --color --color-match "4;33" "\[(WARN)\]" | \
		ag --passthrough --color --color-match "4;32" "\[(PASS)\]"

debug: debug_echo clean debug_body
	@echo [debug] finish

debug_echo:
	@echo [debug] start

debug_body:
	@-./debug.sh
