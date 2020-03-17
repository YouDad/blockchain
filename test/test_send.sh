#!/bin/bash
source test/define.sh $0 $1

rm -rf *9999*
VPortA="-v3 --port 9999"
VPortB="-v3 --port 10003"

RunTest create_wallet "${VPortA}" 's#.*: \(.*\)#\1#g'
AddressA="${TestRegMatch}"

RunTest create_wallet "${VPortB}" 's#.*: \(.*\)#\1#g'
AddressB="${TestRegMatch}"

RunTest create_blockchain "${VPortA} --address ${AddressA}"

RunTest send "${VPortA} --amount 50000000 --from ${AddressA} --to ${AddressB} --mine"

RunTest get_version "${VPortA}"
RunTest get_balance "-v3 ${VPortA} --address ${AddressB}"
RunTest get_balance "-v3 ${VPortA} --address ${AddressA}"

RunTest all "${VPortA}" &
sleep 1

RunTest sync "${VPortB}"
RunTest all "${VPortB}" &
sleep 1

RunTest get_version "${VPortB}"
RunTest send "${VPortB} --amount 1 --from ${AddressB} --to ${AddressA}"
sleep 1
killall blockchain
