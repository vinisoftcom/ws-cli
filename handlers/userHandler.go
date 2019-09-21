package handlers

import "fmt"

type UserHandler struct {
	Command int
	UserId  string
}

func GetUserCommantId(list bool, detail bool) int {
	switch {
	case list:
		{
			return 0
		}
	case detail:
		{
			return 1
		}
	default:
		{
			return -1
		}
	}
}

func (u UserHandler) Run() {
	switch u.Command {
	case 0:
		{
			fmt.Println("List")
		}
	case 1:
		{
			fmt.Println("Detail")
		}
	}
}
