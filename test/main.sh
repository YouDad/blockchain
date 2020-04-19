#!/bin/bash
make
make clean
rm -rf *.log
source test/define.sh $0 $1

Pre="-v3 -g2 --port"

declare -a Address

for (( i=1; i<=2; i++ )); do
	RunTest create_wallet "$Pre $(( i*1111 )) --specified $(( i-1 ))" 's#.*: \(.*\)#\1#g'
	sleep 0.1
	Address["$(( i*1111 ))"]="${TestRegMatch}"
done

for (( i=1; i<=2; i++ )); do
	for (( j=1; j<=3; j++ )); do
		RunTest create_wallet "$Pre $(( i*1100+j )) --specified $(( i-1 ))" 's#.*: \(.*\)#\1#g'
		sleep 0.1
		Address["$(( i*1100+j ))"]="${TestRegMatch}"
	done
done

for (( i=1; i<=2; i++ )); do
	RunTest create_blockchain "$Pre $(( i*1111 )) --address ${Address[i*1111]}"
	sleep 0.1
	RunTest send_test "$Pre $(( i*1111 )) --from ${Address[i*1111]}" &
	sleep 0.1
done

for (( i=1; i<=2; i++ )); do
	for (( j=1; j<=3; j++ )); do
		RunTest sync "$Pre $(( i*1100+j )) --address ${Address[i*1100+j]}"
		sleep 0.1
		RunTest mining "$Pre $(( i*1100+j )) --address ${Address[i*1100+j]}" &
		sleep 0.8
	done
done
