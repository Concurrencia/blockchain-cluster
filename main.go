package main

import (
	"blockchain/blockchain"
	"blockchain/services"
	"os"
)

func main() {

	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	blockchain.PrintAllChain(chain)

	// User services
	go services.ActivarServicioCreateUser(chain)
	go services.ActivarServicioGetAllUsers(chain)
	go services.ActivarServicioGetUserById(chain)
	go services.ActivarServicioGetUserByEmailAndPassword(chain)

	// Consults service
	go services.ActivarServicioGetAllConsults(chain)
	go services.ActivarServicioCreateConsults(chain)
	go services.ActivarServicioGetAllConsultsByUserId(chain)

	procesarNotificaciones()
}

func procesarNotificaciones() {
	for {

	}
}
