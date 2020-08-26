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
	"flag"
	"fmt"
	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
)

type KeyPair struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
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

func main() {
	accountToLoad := flag.String("load", "", "load account from json file instead of creating")

	flag.Parse()

	var keys *keypair.Full = nil

	// New account flow: create and fund via Friendbot.
	if accountToLoad == nil || *accountToLoad == "" {
		keys = createAccount()
		fundAccount(keys)
	} else {
		keys = loadAccount(*accountToLoad)
	}

	if keys != nil {
		printBalances(keys)
	}
}
