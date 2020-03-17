#!/bin/bash
make
rm -rf *.log
source test/define.sh $0 $1

RunTest list_address "-v3 --port 9999" 's#.*[0-9:.]\{15,15\} \(.*\)#\1#g'
Address9999="${TestRegMatch}"
RunTest list_address "-v3 --port 10000" 's#.*[0-9:.]\{15,15\} \(.*\)#\1#g'
Address10000="${TestRegMatch}"
RunTest list_address "-v3 --port 10001" 's#.*[0-9:.]\{15,15\} \(.*\)#\1#g'
Address10001="${TestRegMatch}"

RunTest all "-v3 --port 9999" &
sleep 1

RunTest sync "-v3 --port 10000"
RunTest mining "-v3 --port 10000 --address ${Address10000} --speed 1" &

RunTest sync "-v3 --port 10001"
RunTest mining "-v3 --port 10001 --address ${Address10001} --speed 1" &

sleep 5

RunTest send_test "-v3 --port 9999 --from ${Address9999}"
killall blockchain
