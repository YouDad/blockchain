#!/bin/bash
source test/define.sh $0 $1

rm -rf *9999*
VPortA="-v3 --port 9999"
VPortB="-v3 --port 9994"

RunTest create_wallet "${VPortA}" 's#.*: \(.*\)#\1#g'
AddressA="${TestRegMatch}"

RunTest create_wallet "${VPortB}" 's#.*: \(.*\)#\1#g'
AddressB="${TestRegMatch}"

RunTest create_blockchain "${VPortA} --address ${AddressA}"

RunTest all "${VPortA}" &
sleep 1

RunTest sync "${VPortB}"

killall_blockchain() {
	sleep 3
	killall blockchain
}
killall_blockchain &

RunTest mining "${VPortB} --address ${AddressB}"

RunTest get_version "${VPortA}"
