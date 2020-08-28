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
	"github.com/stellar/go/network"
	proto "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"strconv"
)

var (
	client = horizon.DefaultTestNetClient
)

func printBalances(pair *keypair.Full) error {
	state("Retrieving balances for %s", pair.Address())

	request := horizon.AccountRequest{AccountID: pair.Address()}
	account, err := client.AccountDetail(request)
	if err != nil {
		return printError(err)
	}

	fmt.Println("[*]")
	for _, balance := range account.Balances {
		name := balance.Code
		if name == "" {
			name = "XLM"
		}
		fmt.Printf("  %s (%s)\n", balance.Balance, name)
	}

	return nil
}

func signAndSend(sender *keypair.Full, tx *txnbuild.Transaction) (resp proto.Transaction, err error) {
	// Sign the transaction to prove you are actually the person sending it.
	signedTx, err := tx.Sign(network.TestNetworkPassphrase, sender)
	if err != nil {
		printErrorMsg("Failed signing: %s", err)
		return
	}

	// And finally, send it off to Stellar!
	resp, err = client.SubmitTransaction(signedTx)
	if err != nil {
		printErrorMsg("Failed submitting: %s", err)
		p := horizon.GetError(err)
		fmt.Printf("  Info: %s\n", p.Problem)
		fmt.Printf("  Extras: %s\n",
			p.Problem.Extras["result_codes"])
		return
	}

	fmt.Println("[*]")
	fmt.Println("  Ledger:", resp.Ledger)
	fmt.Println("  Hash:", resp.Hash)
	return
}

func sendMoney(source *keypair.Full, destination string, amount string, asset txnbuild.Asset) (string, error) {
	assetName := "XLM"
	if !asset.IsNative() {
		assetName = asset.GetCode()
	}

	state("Sending %s %s to %s", amount, assetName, destination)

	// Make sure destination account exists
	request := horizon.AccountRequest{AccountID: destination}
	_, err := client.AccountDetail(request)
	if err != nil {
		return "", printError(err)
	}

	// Load the source account
	request = horizon.AccountRequest{AccountID: source.Address()}
	sourceAccount, err := client.AccountDetail(request)
	if err != nil {
		return "", printError(err)
	}

	// Build transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(), // Use a real timeout in production!
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: destination,
					Amount:      amount,
					Asset:       asset,
				},
			},
			Memo: txnbuild.MemoText("dank memes"),
		},
	)

	if err != nil {
		printErrorMsg("Failed building tx: %s", err)
		return "", err
	}

	resp, err := signAndSend(source, tx)
	if err != nil {
		return "", err
	}

	return resp.Hash, nil
}

func requireTrustline(issuer *keypair.Full, asset txnbuild.Asset) error {
	state("Requiring trustline approvals from %s for %s", issuer.Address(), asset.GetCode())

	// Load the source account
	request := horizon.AccountRequest{AccountID: issuer.Address()}
	sourceAccount, err := client.AccountDetail(request)
	if err != nil {
		return printError(err)
	}

	// First, the receiving account must trust the asset if it never has done so.
	for _, balance := range sourceAccount.Balances {
		if balance.Code == asset.GetCode() {
			fmt.Println("[*]\n  Skipping: already exists.")
			return nil
		}
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.ChangeTrust{
					Line:  asset,
					Limit: "2000",
				},
			},
			Memo: txnbuild.MemoText(fmt.Sprintf("i trust u %s bb", asset.GetCode())),
		},
	)

	// The usual song & dance w/ signing, submitting, etc.
	if _, err := signAndSend(issuer, tx); err != nil {
		return err
	}

	return nil
}

func createTrustline(from *keypair.Full, to *keypair.Full, asset txnbuild.Asset) error {
	state("Opening trustline to %s for %s", from.Address(), asset.GetCode())

	// Load the source account
	request := horizon.AccountRequest{AccountID: from.Address()}
	sourceAccount, err := client.AccountDetail(request)
	if err != nil {
		return printError(err)
	}

	// First, the receiving account must trust the asset if it never has done so.
	for _, balance := range sourceAccount.Balances {
		if balance.Code == asset.GetCode() {
			fmt.Println("[*]\n  Skipping: already exists.")
			return nil
		}
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.ChangeTrust{
					Line:  asset,
					Limit: "2000",
				},
			},
			Memo: txnbuild.MemoText(fmt.Sprintf("i trust u %s bb", asset.GetCode())),
		},
	)

	// The usual song & dance w/ signing, submitting, etc.
	if _, err := signAndSend(from, tx); err != nil {
		return err
	}

	return nil
}

func main() {
	accountToLoad := flag.String("load", "", "load account from json file instead of creating")
	accountToPay := flag.String("dest", "", "send 'money' to this account")
	amount := flag.Int("amount", 42, "amount of money to send")
	token := flag.String("asset", "XLM", "type of token to send")
	requireTrust := flag.Bool("require", false, "should we need to approve trustlines?")

	flag.Parse()

	var keys *keypair.Full = nil

	// New account flow: create and fund via Friendbot.
	if accountToLoad == nil || *accountToLoad == "" {
		keys = createAccount()
		fundAccount(keys)
	} else {
		keys = loadAccount(*accountToLoad)
	}

	if keys == nil {
		return
	}

	printBalances(keys)

	if accountToPay == nil || *accountToPay == "" {
		return
	}

	target := loadAccount(*accountToPay)
	if target == nil {
		return
	}

	var asset txnbuild.Asset
	if token != nil && *token == "XLM" {
		// Native assets are simple.
		asset = txnbuild.NativeAsset{}
	} else {
		// In this case, we need to create a trustline from the dest to the
		// source for our new asset, first, and only THEN send the money.
		asset = txnbuild.CreditAsset{Code: *token, Issuer: keys.Address()}

		var err error
		if *requireTrust {
			err = requireTrustline(keys, asset)
		}

		err = createTrustline(target, keys, asset)
		if err != nil {
			return
		}
	}

	sendMoney(keys, target.Address(), strconv.Itoa(*amount), asset)
}
