package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/centum/blockonomics"
)

func main() {
	var (
		token,
		addr,
		txid,
		account,
		currency,
		description,
		tag string
		amount      float64
		reset       bool
		invoiceLive time.Duration
	)

	if len(os.Args) <= 1 {
		usage()
	}
	cmd := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	switch cmd.Name() {
	case "addr_mon_list":
		cmd.StringVar(&token, "token", "", "access token to API Blockonomics")
	case "addr_mon_add":
		cmd.StringVar(&token, "token", "", "access token to API Blockonomics")
		cmd.StringVar(&addr, "addr", "", "Whitespace separated list of bitcoin addresses/xpubs")
		cmd.StringVar(&tag, "tag", "", "tag name")
	case "addr_mon_del":
		cmd.StringVar(&token, "token", "", "access token to API Blockonomics")
		cmd.StringVar(&addr, "addr", "", "Whitespace separated list of bitcoin addresses/xpubs")
	case "addr_new":
		cmd.StringVar(&token, "token", "", "access token to API Blockonomics")
		cmd.StringVar(&account, "account", "", "get address for account")
		cmd.BoolVar(&reset, "reset", false, "reset prev address")
	case "balance":
		cmd.StringVar(&addr, "addr", "", "Whitespace separated list of bitcoin addresses/xpubs")
	case "searchhistory":
		cmd.StringVar(&addr, "addr", "", "Whitespace separated list of bitcoin addresses/xpubs")
	case "tx_detail":
		cmd.StringVar(&txid, "txid", "", "transaction id")
	case "invoice":
		cmd.StringVar(&addr, "addr", "", "Invoice for address")
		cmd.Float64Var(&amount, "amount", 0, "Invoice amount")
		cmd.StringVar(&currency, "currency", "USD", "Invoice currency")
		cmd.StringVar(&description, "description", "", "Invoice description")
		cmd.DurationVar(&invoiceLive, "live", 1*time.Hour, "Invoice live time")
	default:
		usage()
	}
	_ = cmd.Parse(os.Args[2:])

	api := blockonomics.NewClient(token, blockonomics.WithTimeout(30*time.Second))

	switch cmd.Name() {
	case "addr_mon_list":
		dump(api.AddrMonList())

	case "addr_mon_add":
		dump(api.AddrMonitor(addr, tag))

	case "addr_mon_del":
		dump(api.AddrMonDelete(addr))

	case "addr_new":
		dump(api.NewAddress(account, reset))

	case "balance":
		dump(api.Balance(addr))

	case "searchhistory":
		dump(api.SearchHistory(addr))

	case "tx_detail":
		dump(api.TxDetail(txid))

	case "invoice":
		dump(api.Invoice(addr, amount, currency, description, time.Now().Add(invoiceLive)))
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		usage()
	}

}

func dump(v interface{}, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
	d, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(d))
}

func usage() {
	fmt.Printf("usage: %s <command> [<args>]", os.Args[0])
	os.Exit(-1)
}
