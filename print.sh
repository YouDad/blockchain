#!/bin/bash

GroupNumber=2
make
rm .a

function print() {
	port=$1
	group=$2
	echo blockchain -v3 -g$GroupNumber --port $port print --group $group >> .a
	blockchain -v3 -g$GroupNumber --port $port print --group $group 2>> .a
	echo >> .a

}

for (( k=0; k<GroupNumber; k++ )); do
	for (( i=1; i<=GroupNumber; i++ )); do
		print $(( i*1111 )) $k
		for (( j=1; j<=3; j++ )); do
			print $(( i*1100+j )) $k
		done
	done
done
nvim .a
