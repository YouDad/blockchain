FLAG="--port 10000 -v3"

echo -e "\n[TEST]: blockchain create_wallet $FLAG 2>&1"
res=`blockchain create_wallet $FLAG 2>&1`
if [[ "$?" == "0" ]]; then
	echo [PASS]: blockchain create_wallet $FLAG 2>&1
	address=`echo -n $res | sed 's/.*: \(.*\)/\1/g'`
else
	echo [FAIL]: blockchain create_wallet $FLAG 2>&1
	if [[ "$1" == "debug" ]]; then
		echo r create_wallet $FLAG
		dlv debug main.go
	fi
	exit 1
fi

echo -e "\n[TEST]: blockchain create_blockchain $FLAG 2>&1"
FLAG="${FLAG} --address ${address}"
blockchain create_blockchain $FLAG 2>&1
if [[ "$?" == "0" ]]; then
	echo [PASS]: blockchain create_blockchain $FLAG 2>&1
else
	echo [FAIL]: blockchain create_blockchain $FLAG 2>&1
	if [[ "$1" == "debug" ]]; then
		echo r create_blockchain $FLAG
		dlv debug main.go
	fi
	exit 1
fi

echo -e "\n[TEST]: blockchain get_balance $FLAG 2>&1"
blockchain get_balance $FLAG 2>&1
if [[ "$?" == "0" ]]; then
	echo [PASS]: blockchain get_balance $FLAG 2>&1
else
	echo [FAIL]: blockchain get_balance $FLAG 2>&1
	if [[ "$1" == "debug" ]]; then
		echo r get_balance $FLAG
		dlv debug main.go
	fi
	exit 1
fi
