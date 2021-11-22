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

func createUser(con net.Conn, chain *blockchain.BlockChain, bufferIn *bufio.Reader) {

	msg, _ := bufferIn.ReadString('\n')
	msg = strings.TrimSpace(msg)

	newUser := models.User{}

	err := json.Unmarshal([]byte(msg), &newUser)
	blockchain.Handle(err)

	block := chain.AddBlock(newUser)
	blockchain.PrintAllChain(chain)

	hashID := hex.EncodeToString(block.Hash)
	fmt.Fprintln(con, hashID)

	sendChainUpdate(block)
}

func getAllUsers(con net.Conn, chain *blockchain.BlockChain, bufferIn *bufio.Reader) {

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

func getUserById(con net.Conn, chain *blockchain.BlockChain, bufferIn *bufio.Reader) {

	blockchain.PrintAllChain(chain)
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

func getUserByEmailAndPassword(con net.Conn, chain *blockchain.BlockChain, bufferIn *bufio.Reader) {

	blockchain.PrintAllChain(chain)
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
