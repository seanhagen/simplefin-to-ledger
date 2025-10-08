/*
Copyright © 2025 Sean Patrick Hagen <sean.hagen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/fang"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/spf13/cobra"
)

type SimpleFinOrg struct {
	Domain string
	Name   string
	SF_URL string `json:"sfin-url"`
	URL    string
	ID     string
}

type SimpleFinTX struct {
	ID           string
	Posted       int
	Amount       string
	Description  string
	Payee        string
	Memo         string
	TransactedAt int
}

type SimpleFinHolding struct {
	ID            string
	Created       int
	Currency      string
	CostBasis     string
	Description   string
	MarketValue   string
	PurchasePrice string
	Shares        string
	Symbol        string
}

type SimpleFinAccount struct {
	Org              SimpleFinOrg
	ID               string
	Name             string
	Currency         string
	Balance          string
	AvailableBalance string
	BalanceDate      int
	Transactions     []SimpleFinTX
	Holdings         []SimpleFinHolding
}

type SimpleFinData struct {
	Errors   []string
	Accounts []SimpleFinAccount
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simplefin-to-ledger",
	Short: "Download transactions from SimpleFin.org into Ledger",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ArbitraryArgs,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("args:\n%v", args)

		file, err := os.Open(args[0])
		if err != nil {
			return err
		}

		data := SimpleFinData{}

		if err := json.NewDecoder(file).Decode(&data); err != nil {
			return err
		}

		// if err := json.UnmarshalRead(file, &data); err != nil {
		// 	return err
		// }

		for _, acct := range data.Accounts {
			fmt.Printf("Account: %s\n", acct.Name)
			if len(acct.Holdings) > 0 {
				fmt.Printf("Holdings: \n")
				for _, hold := range acct.Holdings {
					fmt.Printf(
						"\tShares: %s\n\tMarket Value: %s\n\n",
						hold.Shares,
						hold.MarketValue,
					)
				}
			}

			if len(acct.Transactions) > 0 {
				fmt.Printf("Transactions: \n")
				for _, tx := range acct.Transactions {
					fmt.Printf(
						"\tDesc: %s\n\tAmount: %s\n\tPayee: %s\n",
						tx.Description,
						tx.Amount,
						tx.Payee,
					)
				}
			}
			fmt.Printf("\n----------\n\n")
		}

		dbFile, err := xdg.DataFile("simplefin-to-ledger/data.db")
		if err != nil {
			return err
		}

		var version string
		db, _ := sql.Open("sqlite3", fmt.Sprintf("file:%s", dbFile))
		db.QueryRow(`SELECT sqlite_version()`).Scan(&version)

		fmt.Printf("SQLite3 version: %s\n", version)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}

	// err := rootCmd.Execute()
	// if err != nil {
	// 	os.Exit(1)
	// }
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.simplefin-to-ledger.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
