all: cc install
	@echo [all]

cc:
	@echo [cc]
	go build main.go

install:
	@echo [install]
	mv main $(GOPATH)/bin/blockchain

clean:
	@echo [clean]
	rm -f blockchain*.db
	rm -f wallet*.dat
	rm -f *.log

test_main: cc install
	@-./test.sh test/main.sh

test: cc install test_echo clean test_body sedlog
	@echo [test] finish

test_echo:
	@echo [test] start

test_body:
	@-./test.sh test/test*.sh

debug: debug_echo clean debug_body
	@echo [debug] finish

debug_echo:
	@echo [debug] start

debug_body:
	@-./debug.sh test/test*.sh

sedlog:
	@echo [log]
	sed -e "s/[\x1b][[0-9;]*[mK]//g" last.color.log > last.log
