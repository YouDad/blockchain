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
	@-./test.sh 2>&1 |\
		ack --flush --passthru --color --color-match "underline bold red" "\[(ERROR|FAIL)\].*" |\
		ack --flush --passthru --color --color-match "bold cyan" "\[(INFO)\].*" |\
		ack --flush --passthru --color --color-match "bold black" "\[(DEBUG)\].*" |\
		ack --flush --passthru --color --color-match "bold blue" "\[(TEST)\].*" |\
		ack --flush --passthru --color --color-match "bold yellow" "\[(WARN)\].*" |\
		ack --flush --passthru --color --color-match "underline bold red" "(NotImplement).*" |\
		ack --flush --passthru --color --color-match "bold green" "\[(PASS)\].*"

debug: debug_echo clean debug_body
	@echo [debug] finish

debug_echo:
	@echo [debug] start

debug_body:
	@-./test.sh debug
