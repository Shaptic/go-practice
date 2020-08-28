package main

/*
 * I want this one to be a little larger.
 *
 * Here are some of the goals:
 *		- create accounts from the command line
 *		- save / load credentials from a file
 *		- initiate transfers between accounts
 *		- view transaction histories
 */

import (
	"encoding/json"
	"fmt"
	"github.com/stellar/go/keypair"
	"io/ioutil"
	"net/http"
	"os"
)

type SerializedAccount struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

func fundAccount(pair *keypair.Full) error {
	pubkey := pair.Address()
	state("Funding %s", pubkey)

	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + pubkey)
	if err != nil {
		return printError(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return printError(err)
	}

	fmt.Println("[*]")

	respJson := make(map[string]interface{})
	json.Unmarshal(body, &respJson)
	fmt.Printf("  Hash: %s\n", respJson["hash"].(string))

	return nil
}

func createAccount() *keypair.Full {
	state("Generating address")

	pair, err := keypair.Random()
	if err != nil {
		fmt.Printf("[-]\n")
		fmt.Printf("  Error: %s\n", err)
		return nil
	}

	fmt.Println("[*]")
	fmt.Printf("  public key: %s\n", pair.Address())
	fmt.Printf("  private key: %s\n", pair.Seed())

	counter := 1
	filename := "./accounts/account.json"
	for {
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			break
		}

		filename = fmt.Sprintf("./accounts/account-%d.json", counter)
		counter++
	}

	state(fmt.Sprintf("Saving to %s", filename))

	f, err := os.Create(filename)
	if err != nil {
		printError(err)
		return nil
	}

	defer f.Close()

	data, _ := json.Marshal(SerializedAccount{Public: pair.Address(), Secret: pair.Seed()})

	f.Write(data)
	f.Sync()

	fmt.Println("[*]")
	return pair
}

func loadAccount(filename string) *keypair.Full {
	state("Loading account from %s", filename)

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// try prepending the directory name, even though we should probably
		// force people to give us paths that make sense...
		filename = fmt.Sprintf("accounts/%s", filename)
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			printErrorMsg("File '%s' does not exist.", filename)
			return nil
		}
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		printError(err)
		return nil
	}

	var result SerializedAccount
	err = json.Unmarshal(data, &result)
	if err != nil {
		printErrorMsg("Decoding '%s' failed:\n  %s", filename, err)
		return nil
	}

	kp, err := keypair.ParseFull(result.Secret)
	if err != nil {
		printError(err)
		return nil
	}

	fmt.Println("[*]")
	return kp
}
