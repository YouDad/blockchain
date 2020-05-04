#!/bin/bash

GroupNumber=2
rm .a

function print() {
	port=$1
	group=$2
	echo blockchain -v3 -g$GroupNumber --port $port print --group $group >> .a
	blockchain -v3 -g$GroupNumber --port $port print --group $group 2>> .a
	echo >> .a

}

for (( i=1; i<=GroupNumber; i++ )); do
	print $(( i*1111 )) $(( GroupNumber-1 ))
	for (( j=1; j<=3; j++ )); do
		print $(( i*1100+j )) $(( GroupNumber-1 ))
	done
done
nvim .a
