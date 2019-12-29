package handlers

import (
	"fmt"
	"log"

	"github.com/vinisoftcom/ws-cli/client"
)

const LOGIN = 0
const LOGOUT = 1
const ISLOGGEDIN = 2

type AuthHandler struct {
	Command int
	Id      string
	Secret  string
}

func GetAuthCommantId(login bool, logout bool, isLoggedIn bool) int {
	switch {
	case login:
		{
			return LOGIN
		}
	case logout:
		{
			return LOGOUT
		}
	case isLoggedIn:
		{
			return ISLOGGEDIN
		}
	default:
		{
			return -1
		}
	}
}

func (auth AuthHandler) Run() {
	switch auth.Command {
	case LOGIN:
		{
			error := client.GetClient().Login(auth.Id, auth.Secret)
			if error != nil {
				log.Fatal(error)
			}
		}
	case LOGOUT:
		{
			client.GetClient().Logout()
		}
	case ISLOGGEDIN:
		{
			if client.GetClient().IsLoggedIn() {
				fmt.Println("User is logged in")
			} else {
				fmt.Println("User is not logged in")
			}
		}
	default:
		{
			fmt.Printf("Bad command called")
		}
	}
}
