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
	"flag"
	"fmt"
	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"io/ioutil"
	"net/http"
	"os"
)

type KeyPair struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

func state(format string, args ...interface{}) {
	fmt.Printf("[ ] %s ...\r", fmt.Sprintf(format, args...))
}

func printErrorMsg(format string, args ...interface{}) {
	fmt.Println("[-]")
	fmt.Printf("  %s\n", fmt.Sprintf(format, args...))
}

func printError(err error) error {
	fmt.Println("[-]")
	fmt.Printf("  Error: %s\n", err)
	return err
}

func printBalances(pair *keypair.Full) error {
	state("Retrieving balances for %s", pair.Address())
	client := horizon.DefaultTestNetClient

	request := horizon.AccountRequest{AccountID: pair.Address()}
	account, err := client.AccountDetail(request)
	if err != nil {
		return printError(err)
	}

	fmt.Println("[*]")
	for _, balance := range account.Balances {
		fmt.Printf("  %s (%s)\n",
			balance.Balance, balance.Type)
	}

	return nil
}

func fundTestnetAccount(pair *keypair.Full) error {
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
	fmt.Println(string(body))

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

		filename = fmt.Sprintf("./account-%d.json", counter)
		counter++
	}

	state(fmt.Sprintf("Saving to %s", filename))

	f, err := os.Create(filename)
	if err != nil {
		printError(err)
		return nil
	}

	defer f.Close()

	data, _ := json.Marshal(KeyPair{Public: pair.Address(), Secret: pair.Seed()})

	f.Write(data)
	f.Sync()

	fmt.Println("[*]")
	return pair
}

func loadAccount(filename string) *keypair.Full {
	state("Loading account from %s", filename)

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		printErrorMsg("File '%s' does not exist.", filename)
		return nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		printError(err)
		return nil
	}

	var result KeyPair
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

func main() {
	accountToLoad := flag.String("load", "", "load account from json file instead of creating")

	flag.Parse()

	var keys *keypair.Full = nil

	// New account flow: create and fund via Friendbot.
	if accountToLoad == nil || *accountToLoad == "" {
		keys = createAccount()
		fundTestnetAccount(keys)
	} else {
		keys = loadAccount(*accountToLoad)
	}

	if keys != nil {
		printBalances(keys)
	}
}