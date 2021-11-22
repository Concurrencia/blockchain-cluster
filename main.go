package main

import (
	"blockchain/blockchain"
	"blockchain/services"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var localhostReg string
var localhostNot string
var localhostServices string
var localhostUpdate string
var remotehost string

var bitacoraAddr []string

func Test() {
	fmt.Println("")
}

func main() {
	defer os.Exit(0)

	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Println("El Nodo utiliza hasta 10 puertos para todos los servicios")
	fmt.Println("Si escoge nodo 0 se utilizaran los puerto  9000 -> 9005")
	fmt.Println("Si escoge nodo 1 se utilizaran los puerto  9010 -> 9015")
	fmt.Println("Si escoge nodo 2 se utilizaran los puerto  9020 -> 9025")
	fmt.Println("El cluster constara de 3 nodos")
	fmt.Printf("Ingrese numero del nodo 0 or 1 or 2: ")

	numNodo, _ := bufferIn.ReadString('\n')
	numNodo = strings.TrimSpace(numNodo)
	localhostReg = fmt.Sprintf("localhost:90%s0", numNodo)
	localhostNot = fmt.Sprintf("localhost:90%s1", numNodo)
	localhostServices = fmt.Sprintf("localhost:90%s2", numNodo)
	localhostUpdate = fmt.Sprintf("localhost:90%s3", numNodo)

	chain := blockchain.InitBlockChain(numNodo)
	defer chain.Database.Close()
	blockchain.PrintAllChain(chain)

	go activarServicioRegistro()
	go activarServicioUpdateChain(chain)

	if numNodo == "0" {
		remotehost = ""
	} else {
		num, _ := strconv.Atoi(numNodo)
		num--
		numNodo = strconv.Itoa(num)
		remotehost = fmt.Sprintf("localhost:90%s0", numNodo)
		registrarSolicitud()
	}

	go procesarNotificaciones()

	services.ActivarServicios(chain, localhostServices)

}

func activarServicioRegistro() {
	ln, _ := net.Listen("tcp", localhostReg)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejadorRegistro(con)
	}
}

func manejadorRegistro(con net.Conn) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	ident, _ := bufferIn.ReadString('\n')
	ident = strings.TrimSpace(ident)

	bitacoraBytes, _ := json.Marshal(append(bitacoraAddr, localhostNot))
	fmt.Fprintln(con, string(bitacoraBytes))

	ident2, _ := bufferIn.ReadString('\n')
	ident2 = strings.TrimSpace(ident2)

	bitacoraBytes, _ = json.Marshal(append(services.BitacoraAddrUpdate, localhostUpdate))
	fmt.Fprintln(con, string(bitacoraBytes))

	comunicarTodos(ident, ident2)

	bitacoraAddr = append(bitacoraAddr, ident)
	services.BitacoraAddrUpdate = append(services.BitacoraAddrUpdate, ident2)

	fmt.Println(bitacoraAddr)
	fmt.Println(services.BitacoraAddrUpdate)
}

func comunicarTodos(ident, ident2 string) {
	for _, addr := range bitacoraAddr {
		notificar(addr, ident, ident2)
	}
}

func notificar(addr, ident, ident2 string) {
	con, _ := net.Dial("tcp", addr)
	defer con.Close()
	fmt.Fprintln(con, ident)
	fmt.Fprintln(con, ident2)
}

func procesarNotificaciones() {
	ln, _ := net.Listen("tcp", localhostNot)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejadorNotificacionesEnviadas(con)
	}
}

func manejadorNotificacionesEnviadas(con net.Conn) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	ident, _ := bufferIn.ReadString('\n')
	ident = strings.TrimSpace(ident)

	bitacoraAddr = append(bitacoraAddr, ident)

	ident2, _ := bufferIn.ReadString('\n')
	ident2 = strings.TrimSpace(ident2)

	services.BitacoraAddrUpdate = append(services.BitacoraAddrUpdate, ident2)

	fmt.Println(bitacoraAddr)
	fmt.Println("update:", services.BitacoraAddrUpdate)
}

func registrarSolicitud() {

	con, _ := net.Dial("tcp", remotehost)
	defer con.Close()
	fmt.Fprintln(con, localhostNot)

	bufferIn := bufio.NewReader(con)
	bitacoraNodo, _ := bufferIn.ReadString('\n')
	var bitacoraTemp []string
	json.Unmarshal([]byte(bitacoraNodo), &bitacoraTemp)

	bitacoraAddr = bitacoraTemp

	fmt.Fprintln(con, localhostUpdate)

	bitacoraNodo, _ = bufferIn.ReadString('\n')
	var bitacoraTempUpdate []string

	json.Unmarshal([]byte(bitacoraNodo), &bitacoraTempUpdate)
	services.BitacoraAddrUpdate = bitacoraTempUpdate
	fmt.Println(bitacoraAddr)
	fmt.Println("update:", services.BitacoraAddrUpdate)
}

func activarServicioUpdateChain(chain *blockchain.BlockChain) {
	ln, _ := net.Listen("tcp", localhostUpdate)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manjerChainUpdate(con, chain)
	}
}

func manjerChainUpdate(con net.Conn, chain *blockchain.BlockChain) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	strBlock, _ := bufferIn.ReadString('\n')
	strBlock = strings.TrimSpace(strBlock)

	block := blockchain.Block{}
	err := json.Unmarshal([]byte(strBlock), &block)
	blockchain.Handle(err)

	if chain.GetBlock(block.Hash) != nil {
		chain.UpdateBlock(block.Hash, &block)
	} else {
		iter := chain.Iterator()
		lastBlock := iter.Next()
		if len(lastBlock.PrevHash) == 0 {
			block.PrevHash = lastBlock.Hash
		}
		chain.AddCreatedBlock(&block)
	}
}
