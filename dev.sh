#!/bin/bash

source test/define.sh

while true; do
	make
	RunTest all "-v3 --port 9999" &
	listen_exec "killall blockchain" --once
done
