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

func ActivarServicioCreateUser(chain *blockchain.BlockChain) {
	ln, _ := net.Listen("tcp", "localhost:9000")
	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go createUser(con, chain)
	}
}

func ActivarServicioGetAllUsers(chain *blockchain.BlockChain) {
	ln, _ := net.Listen("tcp", "localhost:9001")
	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go getAllUsers(con, chain)
	}
}

func ActivarServicioGetUserById(chain *blockchain.BlockChain) {
	ln, _ := net.Listen("tcp", "localhost:9002")
	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go getUserById(con, chain)
	}
}

func ActivarServicioGetUserByEmailAndPassword(chain *blockchain.BlockChain) {
	ln, _ := net.Listen("tcp", "localhost:9003")
	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go getUserByEmailAndPassword(con, chain)
	}
}

func createUser(con net.Conn, chain *blockchain.BlockChain) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	msg, _ := bufferIn.ReadString('\n')
	msg = strings.TrimSpace(msg)

	newUser := models.User{}

	err := json.Unmarshal([]byte(msg), &newUser)
	blockchain.Handle(err)

	block := chain.AddBlock(newUser)
	blockchain.PrintAllChain(chain)

	hashID := hex.EncodeToString(block.Hash)
	fmt.Fprintln(con, hashID)
}

func getAllUsers(con net.Conn, chain *blockchain.BlockChain) {
	defer con.Close()

	users := []models.UserResponse{}

	iter := chain.Iterator()
	for {
		block := iter.Next()
		user := models.UserResponse{
			ID:            hex.EncodeToString(block.Hash),
			Email:         block.Data.Email,
			Password:      block.Data.Password,
			Consultations: block.Data.Consultations,
		}

		users = append(users, user)
		if len(block.PrevHash) == 0 {
			break
		}
	}

	fmt.Println(users)

	byteInfo, _ := json.Marshal(users)
	fmt.Fprintln(con, string(byteInfo))
}

func getUserById(con net.Conn, chain *blockchain.BlockChain) {
	defer con.Close()
	blockchain.PrintAllChain(chain)
	bufferIn := bufio.NewReader(con)
	msg, _ := bufferIn.ReadString('\n')
	fmt.Println("msg:", msg)

	bytes, _ := hex.DecodeString(msg)

	block := chain.GetBlock(bytes)
	if block != nil {
		user := models.UserResponse{
			ID:            hex.EncodeToString(block.Hash),
			Email:         block.Data.Email,
			Password:      block.Data.Password,
			Consultations: block.Data.Consultations,
		}

		byteInfo, _ := json.Marshal(user)
		fmt.Fprintln(con, string(byteInfo))
	} else {
		byteInfo, _ := json.Marshal(models.User{})
		fmt.Fprintln(con, string(byteInfo))
	}

}

func getUserByEmailAndPassword(con net.Conn, chain *blockchain.BlockChain) {
	defer con.Close()
	blockchain.PrintAllChain(chain)
	bufferIn := bufio.NewReader(con)
	email, _ := bufferIn.ReadString('\n')
	email = strings.TrimSpace(email)
	password, _ := bufferIn.ReadString('\n')
	password = strings.TrimSpace(password)
	user := models.UserResponse{}
	iter := chain.Iterator()
	for {
		block := iter.Next()
		if block.Data.Email == email && block.Data.Password == password {
			user = models.UserResponse{
				ID:            hex.EncodeToString(block.Hash),
				Email:         block.Data.Email,
				Password:      block.Data.Password,
				Consultations: block.Data.Consultations,
			}
			fmt.Println("Found")
			break
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	byteInfo, _ := json.Marshal(user)
	fmt.Fprintln(con, string(byteInfo))
}
