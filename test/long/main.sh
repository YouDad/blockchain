#!/bin/bash
make
cd test/long
source ../define.sh $0 $1

VPort9999="-v3 --port 9999"
VPort7777="-v3 --port 7777"

RunTest create_wallet "${VPort9999}" 's#.*: \(.*\)#\1#g'
Address9999="${TestRegMatch}"

RunTest create_wallet "${VPort7777}" 's#.*: \(.*\)#\1#g'
Address7777="${TestRegMatch}"

RunTest create_blockchain "${VPort9999} --address ${Address9999}"

RunTest all "${VPort9999}" &
sleep 1

RunTest sync "${VPort7777}"

killall_blockchain() {
	sleep 3600
	killall blockchain
}
killall_blockchain &

RunTest mining "${VPort7777} --address ${Address7777} --speed 1"

RunTest get_version "${VPort7777}"
RunTest get_version "${VPort9999}"

rm -rf *.db *.dat
