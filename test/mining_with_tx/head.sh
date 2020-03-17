#!/bin/bash
make
source test/define.sh $0 $1

rm -rf *.db *.dat

RunTest create_wallet "-v3 --port 9999" 's#.*: \(.*\)#\1#g'
Address9999="${TestRegMatch}"

RunTest create_wallet "-v3 --port 10000" 's#.*: \(.*\)#\1#g'
Address10000="${TestRegMatch}"

RunTest create_wallet "-v3 --port 10001" 's#.*: \(.*\)#\1#g'
Address10001="${TestRegMatch}"

RunTest create_blockchain "-v3 --port 9999 --address ${Address9999}"

RunTest all "-v3 --port 9999" &
sleep 1

RunTest sync "-v3 --port 10000"
RunTest mining "-v3 --port 10000 --address ${Address10000} --speed 100" &

sleep 4
killall blockchain
