package client

import (
	"os"
	"os/user"
	"path"
	"encoding/json"
	"io/ioutil"
)

type account struct {
	Id string `json:"id"`
	Secret string `json:"secret"`
}

type client struct {
	id string
	secret string
}

func (c client) getConfigFilePath() string {
	user, err := user.Current()
	if err != nil {
		// TODO: Add error handling
		panic("Could not access home folder.")
	}

	return path.Join(user.HomeDir, ".ws-cli", ".account")
}

// Log the user in.
func (c client) Login(id string, secret string) error {
	var configFile = c.getConfigFilePath()

	account := &account{Id: id, Secret: secret}
	fileContents, err := json.Marshal(account)
	if err != nil {
		return err
	}

	var dir = path.Dir(configFile)

	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, os.ModePerm)
		return err
	}

	f, err := os.Create(configFile)
	if err != nil {
		return err
	}

	f.WriteString(string(fileContents))
	defer f.Close()

	return nil
}

func (c client) getAccountConfig() (account, error) {
	var filePath = c.getConfigFilePath()

	if _, err := os.Stat(filePath); err != nil {
	  acc := account{}
	  return acc, err
	}

	f, err := os.Open(filePath)
	if err != nil {
		acc := account{}
		return acc, err
	}
	defer f.Close()

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		acc := account{}
		return acc, err
	}
	contentsString := string(contents)

	var acc account
	errUn := json.Unmarshal([]byte(contentsString), &acc)
	if errUn != nil {
		acc := account{}
		return acc, errUn
	}

	return acc, nil
}

// Check whether user is logged in or not.
func (c client) IsLoggedIn() bool {
	var acc, err = c.getAccountConfig()
	if err != nil {
		return false
	}

	return len(acc.Id) > 0 && len(acc.Secret) > 0
}

// Log the user out.
func (c client) Logout() {
	var configFile = c.getConfigFilePath()

	if _, err := os.Stat(configFile); err == nil {
		os.Remove(configFile)
		// TODO: handle error
	}
}

func GetClient() *client {
	client := client{id: "id", secret: "secret"}
	return &client
}
