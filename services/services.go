package services

import (
	"blockchain/blockchain"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

var BitacoraAddrUpdate []string

func ActivarServicios(chain *blockchain.BlockChain, localhost string) {
	ln, _ := net.Listen("tcp", localhost)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go procesarConexion(con, chain)
	}
}

func procesarConexion(con net.Conn, chain *blockchain.BlockChain) {
	defer con.Close()

	fmt.Println("Bitacora Update", BitacoraAddrUpdate)

	bufferIn := bufio.NewReader(con)
	msg, _ := bufferIn.ReadString('\n')
	msg = strings.TrimSpace(msg)

	switch msg {
	case "createUser":
		createUser(con, chain, bufferIn)
	case "getAllUsers":
		getAllUsers(con, chain, bufferIn)
	case "getUserById":
		getUserById(con, chain, bufferIn)
	case "getUserByEmailAndPassword":
		getUserByEmailAndPassword(con, chain, bufferIn)
	case "getAllConsults":
		getAllConsults(con, chain, bufferIn)
	case "createConsult":
		createConsult(con, chain, bufferIn)
	case "getAllConsultsByUserId":
		getAllConsultsByUserId(con, chain, bufferIn)
	}
}

func sendChainUpdate(block *blockchain.Block) {
	for _, addr := range BitacoraAddrUpdate {
		con, _ := net.Dial("tcp", addr)
		defer con.Close()
		byteInfo, _ := json.Marshal(block)
		fmt.Fprintln(con, string(byteInfo))
	}

}
