#!/bin/bash
source test/define.sh $0 $1

rm -rf *1111*
VPortA="-v3 -g2 --port 1111"

RunTest create_wallet "${VPortA} --specified 0" 's#.*: \(.*\)#\1#g'
AddressA="${TestRegMatch}"

RunTest create_blockchain "${VPortA} --address ${AddressA}"

function killtest() {
	sleep 2
	killall blockchain
}

killtest &

RunTest send_test "$VPortA --from ${AddressA} --wait 0"
