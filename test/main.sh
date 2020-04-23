#!/bin/bash
make
make clean
rm -rf *.log
source test/define.sh $0 $1

GroupNumber=1
Pre="-v3 -g$GroupNumber --port"

declare -a Address

for (( i=1; i<=GroupNumber; i++ )); do
	RunTest create_wallet "$Pre $(( i*1111 )) --specified $(( i-1 ))" 's#.*: \(.*\)#\1#g'
	sleep 0.1
	Address["$(( i*1111 ))"]="${TestRegMatch}"
done

for (( i=1; i<=GroupNumber; i++ )); do
	for (( j=1; j<=3; j++ )); do
		RunTest create_wallet "$Pre $(( i*1100+j )) --specified $(( i-1 ))" 's#.*: \(.*\)#\1#g'
		sleep 0.1
		Address["$(( i*1100+j ))"]="${TestRegMatch}"
	done
done

for (( i=1; i<=GroupNumber; i++ )); do
	RunTest create_blockchain "$Pre $(( i*1111 )) --address ${Address[i*1111]}"
	sleep 0.1
	RunTest send_test "$Pre $(( i*1111 )) --from ${Address[i*1111]}" &
	sleep 0.1
done

function mine() {
	i=$1
	j=$2
	RunTest sync "$Pre $(( i*1100+j )) --address ${Address[i*1100+j]}"
	sleep 2
	RunTest mining "$Pre $(( i*1100+j )) --address ${Address[i*1100+j]}" &
	sleep 3
}

for (( i=1; i<GroupNumber; i++ )); do
	for (( j=1; j<=3; j++ )); do
		mine $i $j
	done
done

for (( j=1; j<=2; j++ )); do
	mine $GroupNumber $j
done

RunTest sync "$Pre $(( GroupNumber*1100+3 )) --address ${Address[GroupNumber*1100+3]}"
sleep 2
RunTest mining "$Pre $(( GroupNumber*1100+3 )) --address ${Address[GroupNumber*1100+3]}"
