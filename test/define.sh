#!/bin/bash

declare -A global
global['filename']=$1
global['testcount']=0
logfile=`tee_log.sh blockchain`

if [[ "$2" == "debug" ]]; then
	NeedDebug=1
else
	NeedDebug=0
fi

function RunTest() {
	subCommand=$1
	parameter=$2
	regularExpr=$3
	directDebug=$4

	(( global['testcount']++ ))

	command="blockchain $subCommand $parameter"
	echo -e "\n====={ [${global['filename']}] TEST${global['testcount']} }=====" | tee -a "$logfile"
	echo -e "[TEST]: $command 2>&1" |\
		ack --flush --passthru --color --color-match "bold blue" "\[(TEST)\].*" |\
		tee -a last.color.log | tee -a "$logfile"
	if [[ "$directDebug" == "debug" ]]; then
		echo r $subCommand $parameter | xsel -b -i
		echo Ctrl+Shift+V to paste
		dlv debug main.go
		return
	fi
	$command 2>&1 | tee /tmp/a |\
		ack --flush --passthru --color --color-match "underline bold red" "(\[ERROR\]|NotImplement|.*panic).*" |\
		ack --flush --passthru --color --color-match "bold cyan" "\[(INFO)\].*" |\
		ack --flush --passthru --color --color-match "bold black" "\[(DEBUG)\].*" |\
		ack --flush --passthru --color --color-match "bold yellow" "\[(WARN)\].*" |\
		ack --flush --passthru --color --color-match "underline bold red on_green" "\[(TRACE)\].*" |\
		tee -a "$logfile"
	res=`cat /tmp/a`
	echo -en "$res" | grep "\[ERROR\]" >/dev/null
	rescode="$?"
	if [[ "$rescode" == "1" ]]; then
		echo [PASS]: $command |\
			ack --flush --passthru --color --color-match "bold green" "\[(PASS)\].*" |\
			tee -a last.color.log | tee -a "$logfile"
		if [[ "$regularExpr" != "" ]]; then
			TestRegMatch=$(echo -n "$res" | sed "$regularExpr")
			if [[ "$?" != "0" ]]; then
				echo "RE:{$regularExpr}" | tee -a "$logfile"
			fi
			regularExpr=''
		fi
	else
		echo [FAIL]: $command |\
			ack --flush --passthru --color --color-match "underline bold red" "\[FAIL\].*" |\
			tee -a last.color.log | tee -a "$logfile"
		if [[ "$NeedDebug" == "1" ]]; then
			echo r $subCommand $parameter | xsel -b -i
			echo Ctrl+Shift+V to paste
			dlv debug main.go
		fi
		kill -9 $pid > /dev/null 2>&1
		exit 1
	fi
}
