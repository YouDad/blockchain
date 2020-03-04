if [[ "$1" == "debug" ]]; then
	NeedDebug=1
else
	NeedDebug=0
fi

function RunTest() {
	subCommand=$1
	parameter=$2
	regularExpr=$3
	directDebug=$4

	if [[ "$_TestCount" == "" ]]; then
		_TestCount=0
	fi
	(( _TestCount++ ))

	command="blockchain $subCommand $parameter"
	echo -e "\n====={ TEST$_TestCount }====="
	echo -e "[TEST]: $command 2>&1"
	if [[ "$directDebug" == "debug" ]]; then
		echo r $subCommand $parameter | xsel -b -i
		echo Ctrl+Shift+V to paste
		dlv debug main.go
		return
	fi
	$command 2>&1 | tee /tmp/a |\
		ack --flush --passthru --color --color-match "underline bold red" "(\[ERROR\]|NotImplement).*"
	rescode="$?"
	res=`cat /tmp/a`
	if [[ "$rescode" == "1" ]]; then
		echo [PASS]: $command
		if [[ "$regularExpr" != "" ]]; then
			TestRegMatch=$(echo -n "$res" | sed "$regularExpr")
			if [[ "$?" != "0" ]]; then
				echo "RE:{$regularExpr}"
			fi
		fi
	else
		echo [FAIL]: $command
		if [[ "$NeedDebug" == "1" ]]; then
			echo r $subCommand $parameter | xsel -b -i
			echo Ctrl+Shift+V to paste
			dlv debug main.go
		fi
		exit 1
	fi
}

Port9999="--port 9999"
Port8888="--port 8888"
Port7777="--port 7777"
VPort9999="-v3 --port 9999"
VPort8888="-v3 --port 8888"
VPort7777="-v3 --port 7777"

killall blockchain

RunTest create_wallet "${VPort9999}" 's#.*: \(.*\)#\1#g'
Address9999="${TestRegMatch}"

RunTest create_wallet "${VPort8888}" 's#.*: \(.*\)#\1#g'
Address8888="${TestRegMatch}"

RunTest create_wallet "${VPort7777}" 's#.*: \(.*\)#\1#g'
Address7777="${TestRegMatch}"

# RunTest list_address "${VPort}"

RunTest create_blockchain "${VPort9999} --address ${Address9999}"

# RunTest get_version "${VPort}"

# RunTest get_balance "${VPort} --address ${FromAddress}"

RunTest send "${VPort9999} --amount 50000000 --from ${Address9999} --to ${Address7777} --mine"

RunTest get_version "${VPort9999}"
RunTest get_balance "-v3 ${Port9999} --address ${Address7777}"
RunTest get_balance "-v3 ${Port9999} --address ${Address9999}"

RunTest start_node "${VPort9999}" &
sleep 1
RunTest sync "${VPort7777}"
RunTest get_version "${VPort7777}"
RunTest send "${VPort7777} --amount 1 --from ${Address7777} --to ${Address8888}"
