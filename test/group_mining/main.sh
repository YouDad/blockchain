#!/bin/bash
make
make clean
source test/define.sh $0 $1

VPortA="-v3 --port 9999"
VPortB="-v3 --port 8888"
VPortC="-v3 --port 8000"

RunTest create_wallet "${VPortA}" 's#.*: \(.*\)#\1#g'
AddressA="${TestRegMatch}"

RunTest create_wallet "${VPortB}" 's#.*: \(.*\)#\1#g'
AddressB="${TestRegMatch}"

RunTest create_wallet "${VPortC}" 's#.*: \(.*\)#\1#g'
AddressC="${TestRegMatch}"

RunTest create_blockchain "${VPortA} --address ${AddressA}"
RunTest create_blockchain "${VPortB} --address ${AddressB}"

RunTest all "${VPortA} --group 1" &
RunTest all "${VPortB} --group 1" &
sleep 1

# RunTest sync "${VPortC} --group 2"

RunTest mining "${VPortC} --address ${AddressC} --speed 1 --group 2"

killall blockchain
