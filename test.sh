FLAG="--port 10000"

res=`blockchain create_wallet $FLAG`
echo -n $res | grep "^Your new address: [0-9A-Za-z]*$"
if [[ "$?" == "0" ]]; then
	echo [PASS]: blockchain create_wallet
	address=`echo -n $res | sed 's/.*: \(.*\)/\1/g'`
else
	echo [FAIL]: blockchain create_wallet
	exit 1
fi

FLAG="${FLAG} --address ${address}"
res=`blockchain create_blockchain $FLAG`
echo $res
