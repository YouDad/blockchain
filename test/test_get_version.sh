#!/bin/bash
source test/define.sh $0 $1

VPortA="-v3 --port 10001"

RunTest create_wallet "${VPortA}" 's#.*: \(.*\)#\1#g'
AddressA="${TestRegMatch}"

RunTest create_blockchain "${VPortA} --address ${AddressA}"
RunTest get_version "${VPortA}"
