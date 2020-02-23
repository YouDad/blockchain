FLAG="--port 10000"

res=`blockchain create_wallet $FLAG`
echo -n $res | grep "^Your new address: [0-9A-Za-z]*$"
if [[ "$?" == "0" ]]; then
	echo [PASS]: blockchain create_wallet
	address=`echo -n $res | sed 's/.*: \(.*\)/\1/g'`
else
	echo [FAIL]: blockchain create_wallet
	echo create_wallet $FLAG
	dlv debug main.go
	exit 1
fi

FLAG="${FLAG} --address ${address}"
blockchain create_blockchain $FLAG 2>&1 |\
	ag --passthrough --color --color-match "5;31" "Done!"
if [[ "$?" == "0" ]]; then
	echo [PASS]: blockchain create_blockchain
else
	echo [FAIL]: blockchain create_blockchain
	rm -f blockchain10000.db
	echo r create_blockchain $FLAG
	dlv debug main.go
	exit 1
fi
