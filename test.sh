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
	res=$($command 2>&1)
	rescode="$?"
	echo "$res"
	if [[ "$rescode" == "0" ]]; then
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

RunTest create_wallet "${VPort}" 's#.*: \(.*\)#\1#g'
FromAddress="${TestRegMatch}"

RunTest create_wallet "${VPort}" 's#.*: \(.*\)#\1#g'
ToAddress="${TestRegMatch}"

# RunTest list_address "${VPort}"

RunTest create_blockchain "${VPort} --address ${FromAddress}"

# RunTest get_version "${VPort}"

# RunTest get_balance "${VPort} --address ${FromAddress}"

RunTest send "${VPort} --amount 1 --from ${FromAddress} --to ${ToAddress} --mine"
