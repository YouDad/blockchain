#!/bin/bash
make
source test/define.sh $0 $1

rm -rf *.db *.dat

RunTest create_wallet "-v3 --port 9999" 's#.*: \(.*\)#\1#g'
Address9999="${TestRegMatch}"

for (( i=2;i<=9;i++)); do
	for (( j=1;j<=3;j++)); do
		RunTest create_wallet "-v3 --port $(( i*1000+j ))" 's#.*: \(.*\)#\1#g'
	done
done

RunTest create_blockchain "-v3 --port 9999 --address ${Address9999}"
