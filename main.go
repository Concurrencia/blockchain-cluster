package main

import (
	"blockchain/blockchain"
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
var localhostchain string
var localhostUpdate string
var remotehost string // para conectarse al nodo externo para la solicitud de registro a la red

var bitacoraAddr []string //todos los localhost + puerto notificacion

func main() {
	defer os.Exit(0)

	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Println("El Nodo utiliza hasta 10 puertos para todos los servicios")
	fmt.Println("Si escoge nodo 0 se utilizaran los puerto  9000 -> 9009")
	fmt.Println("Si escoge nodo 1 se utilizaran los puerto  9010 -> 9019")
	fmt.Println("Si escoge nodo 2 se utilizaran los puerto  9020 -> 9029")
	fmt.Println("El cluster constara de 3 nodos")
	fmt.Printf("Ingrese numero del nodo 0 or 1 or 2: ")
	numNodo, _ := bufferIn.ReadString('\n')
	numNodo = strings.TrimSpace(numNodo)
	localhostReg = fmt.Sprintf("localhost:90%s0", numNodo)
	localhostNot = fmt.Sprintf("localhost:90%s1", numNodo)

	go activarServicioRegistro()

	chain := blockchain.InitBlockChain(numNodo)
	defer chain.Database.Close()
	blockchain.PrintAllChain(chain)

	if numNodo == "0" {
		remotehost = ""
	} else {
		num, _ := strconv.Atoi(numNodo)
		num--
		numNodo = strconv.Itoa(num)
		remotehost = fmt.Sprintf("localhost:90%s0", numNodo)
		registrarSolicitud()
	}

	// User services
	// go services.ActivarServicioCreateUser(chain)
	// go services.ActivarServicioGetAllUsers(chain)
	// go services.ActivarServicioGetUserById(chain)
	// go services.ActivarServicioGetUserByEmailAndPassword(chain)

	// // Consults service
	// go services.ActivarServicioGetAllConsults(chain)
	// go services.ActivarServicioCreateConsults(chain)
	// go services.ActivarServicioGetAllConsultsByUserId(chain)

	procesarNotificaciones()

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

	comunicarTodos(ident)

	bitacoraAddr = append(bitacoraAddr, ident)

	fmt.Println(bitacoraAddr)
}

func comunicarTodos(ident string) {
	for _, addr := range bitacoraAddr {
		notificar(addr, ident)
	}
}

func notificar(addr, ident string) {
	con, _ := net.Dial("tcp", addr)
	defer con.Close()
	fmt.Fprintln(con, ident)
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

	fmt.Println(bitacoraAddr)
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

	fmt.Println(bitacoraAddr)
}
