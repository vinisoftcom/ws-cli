package handlers

import "fmt"

const LIST = 0
const DETAIL = 1

type UserHandler struct {
	Command int
	UserId  string
}

func GetUserCommantId(list bool, detail bool) int {
	switch {
	case list:
		{
			return LIST
		}
	case detail:
		{
			return DETAIL
		}
	default:
		{
			return -1
		}
	}
}

func (u UserHandler) Run() {
	switch u.Command {
	case LIST:
		{
			fmt.Println("List")
		}
	case DETAIL:
		{
			fmt.Println("Detail")
		}
	}
}
