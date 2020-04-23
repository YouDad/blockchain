#!/bin/bash
source test/define.sh $0 $1

rm -rf *1111*
VPortA="-v3 --port 1111"
VPortB="-v3 --port 1104"

RunTest create_wallet "${VPortA} --specified 0" 's#.*: \(.*\)#\1#g'
AddressA="${TestRegMatch}"

RunTest create_wallet "${VPortB} --specified 0" 's#.*: \(.*\)#\1#g'
AddressB="${TestRegMatch}"

RunTest create_blockchain "${VPortA} --address ${AddressA}"

RunTest send "${VPortA} --amount 50000000 --from ${AddressA} --to ${AddressB} --mine"

RunTest get_version "${VPortA} --address ${AddressA}"
RunTest get_balance "-v3 ${VPortA} --address ${AddressB}"
RunTest get_balance "-v3 ${VPortA} --address ${AddressA}"

RunTest all "${VPortA} --address ${AddressA}" &
sleep 1

RunTest sync "${VPortB} --address ${AddressB}"
RunTest all "${VPortB} --address ${AddressB}" &
sleep 1

RunTest get_version "${VPortB} --address ${AddressB}"
RunTest send "${VPortB} --amount 1 --from ${AddressB} --to ${AddressA}"
sleep 1
killall blockchain
