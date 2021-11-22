package services

import (
	"blockchain/blockchain"
	"blockchain/models"
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

func createConsult(con net.Conn, chain *blockchain.BlockChain, bufferIn *bufio.Reader) {

	blockchain.PrintAllChain(chain)
	userHash, _ := bufferIn.ReadString('\n')

	consult, _ := bufferIn.ReadString('\n')
	consult = strings.TrimSpace(consult)

	newConsult := models.Consultation{}
	err := json.Unmarshal([]byte(consult), &newConsult)
	blockchain.Handle(err)

	bytes, _ := hex.DecodeString(userHash)
	block := chain.GetBlock(bytes)

	if block == nil {
		fmt.Fprintln(con, "nil")
	} else {
		fmt.Fprintln(con, "")
		id := len(block.Data.Consultations) + 1
		newConsult.ID = id
		newConsult.UserID = hex.EncodeToString(block.Hash)
		block.Data.Consultations = append(block.Data.Consultations, newConsult)
		chain.UpdateBlock(block.Hash, block)

		byteInfo, _ := json.Marshal(newConsult)
		fmt.Fprintln(con, string(byteInfo))

		sendChainUpdate(block)
	}
}

func getAllConsults(con net.Conn, chain *blockchain.BlockChain, bufferIn *bufio.Reader) {

	consultation := []models.Consultation{}

	iter := chain.Iterator()
	for {
		block := iter.Next()

		consultation = append(consultation, block.Data.Consultations...)
		if len(block.PrevHash) == 0 {
			break
		}
	}

	fmt.Println(consultation)

	byteInfo, _ := json.Marshal(consultation)
	fmt.Fprintln(con, string(byteInfo))
}

func getAllConsultsByUserId(con net.Conn, chain *blockchain.BlockChain, bufferIn *bufio.Reader) {

	blockchain.PrintAllChain(chain)
	msg, _ := bufferIn.ReadString('\n')
	fmt.Println("msg:", msg)

	bytes, _ := hex.DecodeString(msg)
	block := chain.GetBlock(bytes)

	if block == nil {
		fmt.Fprintln(con, "nil")
	} else {
		fmt.Fprintln(con, "")
		byteInfo, _ := json.Marshal(block.Data.Consultations)
		fmt.Fprintln(con, string(byteInfo))
	}
}
