#!/bin/bash
make
rm -rf *.log
source test/define.sh $0 $1

RunTest all "-v3 --port 9999" &
sleep 1

for ((i=2;i<=8;i++)); do
	RunTest sync "-v3 --port $((i*1111))"
	RunTest all "-v3 --port $((i*1111))" &
done

for (( i=2;i<=9;i++)); do
	for (( j=1;j<=3;j++)); do
		sleep 1
		RunTest sync "-v3 --port $(( i*1000+j ))"
		RunTest list_address "-v3 --port $(( i*1000+j ))" 's#.*[0-9:.]\{15,15\} \(.*\)#\1#g'
		RunTest mining "-v3 --port $(( i*1000+j )) --address ${TestRegMatch}" &
	done
done
