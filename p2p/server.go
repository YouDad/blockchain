package p2p

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"net"

	coin_core "github.com/YouDad/blockchain/app/coin/core"
	"github.com/YouDad/blockchain/core"
	"github.com/YouDad/blockchain/log"
)

const (
	protocol      = "udp"
	nodeVersion   = 1
	commandLength = 12
)

var (
	nodeAddress     string
	miningAddress   string
	KnownNodes      = []string{"39.107.64.93:9999"}
	blocksInTransit = [][]byte{}
	mempool         = make(map[string]coin_core.Transaction)
)

type RWriter struct {
	Read  []byte
	Write []byte
}

type addr struct {
	AddrList []string
}

type block struct {
	Block []byte
}

type getdata struct {
	Type string
	ID   []byte
}

type inv struct {
	Type  string
	Items [][]byte
}

type tx struct {
	Transaction []byte
}

type version struct {
	Version    int
	BestHeight int
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(rw *RWriter) []byte {
	return rw.Read[:commandLength]
}

func sendAddr(rw *RWriter) {
	nodes := addr{KnownNodes}
	payload := gobEncode(nodes)
	rw.Write = append(commandToBytes("addr"), payload...)
}

func sendBlock(rw *RWriter, b *core.Block) {
	data := block{b.Serialize()}
	payload := gobEncode(data)
	rw.Write = append(commandToBytes("block"), payload...)
}

func sendInv(rw *RWriter, kind string, items [][]byte) {
	inventory := inv{kind, items}
	payload := gobEncode(inventory)
	rw.Write = append(commandToBytes("inv"), payload...)
}

func sendGetGenesis(rw *RWriter) {
	log.Println("sendGetGenesis")
	rw.Write = commandToBytes("getgenesis")
}

func sendGetBlocks(rw *RWriter) {
	rw.Write = commandToBytes("getblocks")
}

func sendGetData(rw *RWriter, kind string, id []byte) {
	payload := gobEncode(getdata{kind, id})
	rw.Write = append(commandToBytes("getdata"), payload...)
}

func SendTx(tnx *coin_core.Transaction) {
	data := tx{tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)
	for _, node := range KnownNodes {
		if nodeAddress != node {
			conn, err := net.Dial(protocol, node)
			if err != nil {
				log.Panic(err)
			}
			conn.Write(request)
			conn.Close()
		}
	}
}

func sendTx(rw *RWriter, tnx *coin_core.Transaction) {
	data := tx{tnx.Serialize()}
	payload := gobEncode(data)
	rw.Write = append(commandToBytes("tx"), payload...)
}

func sendVersion(rw *RWriter, bc *coin_core.CoinBlockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight})
	rw.Write = append(commandToBytes("version"), payload...)
}

func handleAddr(rw *RWriter) []string {
	var buff bytes.Buffer
	var payload addr

	buff.Write(rw.Read[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(KnownNodes))
	return KnownNodes
}

func handleBlock(rw *RWriter, bc *coin_core.CoinBlockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(rw.Read[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := core.DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(rw, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := coin_core.NewUTXOSet()
		defer UTXOSet.Close()
		UTXOSet.Reindex()
	}
}

func handleInv(rw *RWriter, bc *coin_core.CoinBlockchain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(rw.Read[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(rw, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "genesis" {
		genesisBlock := core.DeserializeBlock(payload.Items[0])
		bc.Blockchain = core.CreateBlockchainFromGenesis(genesisBlock)
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(rw, "tx", txID)
		}
	}
}

func handleGetGenesis(rw *RWriter, bc *coin_core.CoinBlockchain) {
	log.Println("handleGetGenesis")
	genesisBlock := bc.GetGenesis()
	sendInv(rw, "genesis", [][]byte{genesisBlock})
}

func handleGetBlocks(rw *RWriter, bc *coin_core.CoinBlockchain) {
	blocks := bc.GetBlockHashes()
	sendInv(rw, "block", blocks)
}

func handleGetData(rw *RWriter, bc *coin_core.CoinBlockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(rw.Read[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(rw, block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(rw, &tx)
		// delete(mempool, txID)
	}
}

func handleTx(rw *RWriter, bc *coin_core.CoinBlockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(rw.Read[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := coin_core.DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	// if nodeAddress == KnownNodes[0] {
	//     for _, node := range KnownNodes {
	//         if node != nodeAddress && node != payload.AddFrom {
	//             sendInv(node, "tx", [][]byte{tx.ID})
	//         }
	//     }
	// } else {
	if len(mempool) >= 2 && len(miningAddress) > 0 {
	MineTransactions:
		var txs []*coin_core.Transaction

		for id := range mempool {
			tx := mempool[id]
			if bc.VerifyTransaction(&tx) {
				txs = append(txs, &tx)
			}
		}

		if len(txs) == 0 {
			fmt.Println("All transactions are invalid! Waiting for new ones...")
			return
		}

		cbTx := coin_core.NewCoinbaseTX(miningAddress, "")
		txs = append(txs, cbTx)

		newBlock := bc.MineBlock(txs)
		UTXOSet := coin_core.NewUTXOSet()
		defer UTXOSet.Close()
		UTXOSet.Reindex()

		fmt.Println("New block is mined!")

		for _, tx := range txs {
			txID := hex.EncodeToString(tx.ID)
			delete(mempool, txID)
		}

		for _, node := range KnownNodes {
			if node != nodeAddress {
				sendInv(rw, "block", [][]byte{newBlock.Hash})
			}
		}

		if len(mempool) > 0 {
			goto MineTransactions
		}
	}
	// }
}

func handleVersion(rw *RWriter, bc *coin_core.CoinBlockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(rw.Read[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(rw)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(rw, bc)
	}

	// if !nodeIsKnown(conn.RemoteAddr().String()) {
	//     KnownNodes = append(KnownNodes, conn.RemoteAddr().String())
	// }
}

func handleConnection(read []byte, remoteAddr *net.UDPAddr, listener *net.UDPConn, bc *coin_core.CoinBlockchain) {
	var rw RWriter
	rw.Read = read
	command := bytesToCommand(rw.Read[:commandLength])
	log.Printf("[from %s]:Received %s command\n", remoteAddr.String(), command)

	switch command {
	case "addr":
		handleAddr(&rw)
	case "block":
		handleBlock(&rw, bc)
	case "inv":
		handleInv(&rw, bc)
	case "getgenesis":
		handleGetGenesis(&rw, bc)
	case "getblocks":
		handleGetBlocks(&rw, bc)
	case "getdata":
		handleGetData(&rw, bc)
	case "tx":
		handleTx(&rw, bc)
	case "version":
		handleVersion(&rw, bc)
	default:
		fmt.Println("Unknown command!")
	}
	if len(rw.Write) > 0 {
		log.Println("Write", remoteAddr.String(), len(rw.Write))
		listener.WriteToUDP(rw.Write, remoteAddr)
	}
}

// StartServer starts a node
func StartServer(nodeID, nodeIP, minerAddress string) {
	listenAddress := fmt.Sprintf("0.0.0.0:%s", nodeID)
	nodeAddress = fmt.Sprintf("%s:%s", nodeIP, nodeID)
	miningAddress = minerAddress
	ln, err := net.ListenUDP(protocol, &net.UDPAddr{IP: net.IPv4zero, Port: 9999})
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	var bc *coin_core.CoinBlockchain
	if !core.IsBlockchainExists() {
		conn, err := net.Dial(protocol, KnownNodes[0])
		if err != nil {
			log.Panic(err)
		}
		var rw RWriter
		sendGetGenesis(&rw)
		conn.Write(rw.Write)
		conn.Close()

		data := make([]byte, 1024)
		n, remoteAddr, err := ln.ReadFromUDP(data)
		if err != nil {
			log.Panic(err)
		}
		handleConnection(data[:n], remoteAddr, ln, bc)
	} else {
		bc = coin_core.NewBlockchain()
	}

	if listenAddress != KnownNodes[0] {
		conn, err := net.Dial(protocol, KnownNodes[0])
		if err != nil {
			log.Panic(err)
		}
		var rw RWriter
		sendVersion(&rw, bc)
		conn.Write(rw.Write)
		conn.Close()
	}

	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := ln.ReadFromUDP(data)
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(data[:n], remoteAddr, ln, bc)
	}
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func nodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}

	return false
}
