#!/bin/bash
source test/define.sh $0 $1

VPortA="-v3 --port 10002"

RunTest create_wallet "${VPortA}" 's#.*: \(.*\)#\1#g'
AddressA="${TestRegMatch}"

RunTest list_address "${VPortA}"
