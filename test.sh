if [[ "$1" == "debug" ]]; then
	NeedDebug=1
else
	NeedDebug=0
fi

function RunTest() {
	subCommand=$1
	parameter=$2
	regularExpr=$3

	if [[ "$_TestCount" == "" ]]; then
		_TestCount=0
	fi
	(( _TestCount++ ))

	command="blockchain $subCommand $parameter"
	echo -e "\n====={ TEST$_TestCount }====="
	echo -e "[TEST]: $command 2>&1"
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
			echo r $subCommand $parameter
			dlv debug main.go
		fi
		exit 1
	fi
}

VPort="-v3 --port 10000"

# RunTest create_wallet "${VPort}" 's#.*: \(.*\)#\1#g'
# FromAddress="${TestRegMatch}"
#
# RunTest create_wallet "${VPort}" 's#.*: \(.*\)#\1#g'
# ToAddress="${TestRegMatch}"
#
# RunTest list_address "${VPort}"
#
# RunTest create_blockchain "${VPort} --address ${FromAddress}"
#
# RunTest get_version "${VPort}"
#
# RunTest get_balance "${VPort} --address ${FromAddress}"
#
# RunTest send "${VPort} --amount 1 --from ${FromAddress} --to ${ToAddress} --mine"

VPort9999="-v3 --port 9999"
VPort7777="-v3 --port 7777"

RunTest create_wallet "${VPort9999}" 's#.*: \(.*\)#\1#g'
Address9999="${TestRegMatch}"

RunTest create_wallet "${VPort7777}" 's#.*: \(.*\)#\1#g'
Address7777="${TestRegMatch}"

# RunTest list_address "${VPort}"

RunTest create_blockchain "${VPort9999} --address ${Address9999}"

# RunTest get_version "${VPort}"

# RunTest get_balance "${VPort} --address ${FromAddress}"

# RunTest send "${VPort} --amount 1 --from ${FromAddress} --to ${ToAddress} --mine"

RunTest start_node "${VPort9999} --address ${Address9999}" &
sleep 5
RunTest start_node "${VPort7777} --address ${Address7777}"
