package client

import (
	"os"
	"os/user"
	"path"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/base64"
	"time"
	"net/http"
)

type account struct {
	Id string `json:"id"`
	Secret string `json:"secret"`
}

type client struct {  }

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

// func (c client) Test() {
	// var sig = c.getSignature("GET", "/v1/user/self", 123)
	// println(sig)
// }

func (c client) getSignature(method string, path string, time int64, key string) string {
	var req = c.getCanonicalRequest(method, path, time)
	res := hmac.New(sha1.New, []byte(key))
	res.Write([]byte(req))
	sha := hex.EncodeToString(res.Sum(nil))
	return sha
}

func (c client) getCanonicalRequest(method string, path string, time int64) string {
	return fmt.Sprintf("%s %s %d", method, path, time)
}

func (c client) CurrentUser() (string, error) {
	var account, _ = c.getAccountConfig()

	location,_ := time.LoadLocation("GMT")
	var now = time.Now().In(location)
	var time = now.UnixNano() / int64(time.Second)
	var method = "GET"
	var path = "/v1/user/self"
	var signature = c.getSignature(method, path, time, account.Secret)

	var auth = fmt.Sprintf("%s:%s", account.Id, signature)
	basic :=  base64.StdEncoding.EncodeToString([]byte(auth))

	var url = fmt.Sprintf("https://rest.websupport.sk%s", path)

	httpClient := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basic))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Date", now.Format("20060102T150405Z"))

	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func GetClient() *client {
	client := client{}
	return &client
}
