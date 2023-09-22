package main

import (
	"crypto-price/api"
	"crypto-price/datatypes"
	"crypto-price/util"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {

	walletURL, err := util.GetFileUrl("wallet")
	if err != nil {
		panic(err)
	}

	var profile datatypes.Profile
	if err := profile.GetFromJsonFile(); err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "gm",
		Usage: "are you rich or broke today?",
		Action: func(cCtx *cli.Context) error {
			if profile.Token == "" {
				fmt.Println("Hello and welcome. To use gm please set your currency api token.")
				fmt.Println("You can get one for free at https://currencyapi.com/")
				fmt.Println("Then run `gm set token <your token>`")
				return nil
			}

			var (
				greeting string
			)

			fiatSymbol := profile.GetFiatSymbol()

			if time.Now().Hour() < 12 {
				greeting = "Good morning"
			} else if time.Now().Hour() < 18 {
				greeting = "Good afternoon"
			} else {
				greeting = "Good evening"
			}

			fmt.Printf("\n\n\n\n%s %s, let's check your crypto on this fine day.\n\n", greeting, profile.Name)

			exchangeRate, err := api.GetUsdGbpExchangeRate()
			if err != nil {
				return err
			}

			var wg sync.WaitGroup

			totalWorth := 0.0

			var wallet datatypes.Wallet
			if err := wallet.GetFromJsonFile(); err != nil {
				return err
			}

			for symbol, amount := range wallet {
				wg.Add(1)
				go func(s string, a float64) {
					defer wg.Done()
					r, err := api.GetRate(s)

					if err != nil {
						panic(err)
					}

					fmt.Printf("You have %s%.2f worth of %s\n", fiatSymbol, r.Price*a*exchangeRate, s)
					totalWorth += r.Price * a * exchangeRate
				}(symbol, amount)
			}

			wg.Wait()

			defer fmt.Printf("\nTotal: %s%.2f at %s\n", fiatSymbol, totalWorth, time.Now().Format("15:04:05"))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "get",
				Subcommands: []*cli.Command{
					{
						Name:    "profile",
						Aliases: []string{"p"},
						Action: func(cCtx *cli.Context) error {
							fmt.Printf("Name: %s\n", profile.Name)
							fmt.Printf("Token: %s\n", profile.Token)
							fmt.Printf("Base Currency: %s\n", profile.BaseCurrency)
							return nil
						},
					},
					{
						Name:    "wallet",
						Aliases: []string{"w"},
						Action: func(cCtx *cli.Context) error {
							var wallet datatypes.Wallet
							if err := wallet.GetFromJsonFile(); err != nil {
								return err
							}

							for currency, amount := range wallet {
								fmt.Printf("%s: %.2f\n", currency, amount)
							}

							return nil
						},
					},
				},
			},
			{
				Name: "set",
				Subcommands: []*cli.Command{
					{
						Name:    "token",
						Aliases: []string{"t"},
						Usage:   "sets the currency api token",
						Action: func(cCtx *cli.Context) error {
							return profile.SetTokenAndSave(cCtx.Args().First())
						},
					},
					{
						Name:    "currency",
						Aliases: []string{"c"},
						Usage:   "sets the base currency, e.g. GBP, USD, EUR",
						Action: func(cCtx *cli.Context) error {
							return profile.SetCurrencyAndSave(cCtx.Args().First())
						},
					},
					{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "sets the name of the user",
						Action: func(cCtx *cli.Context) error {
							return profile.SetNameAndSave(cCtx.Args().First())
						},
					},
				},
			},
			{
				Name:    "buy",
				Aliases: []string{"b"},
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "set"},
				},
				Usage: "buy a crypto currency",
				Action: func(cCtx *cli.Context) error {
					amount := cCtx.Args().First()
					currency := strings.ToUpper(cCtx.Args().Get(1))

					if amount == "" || currency == "" {
						return fmt.Errorf("please specify an amount and currency")
					}

					var wallet datatypes.Wallet
					if err := wallet.GetFromJsonFile(); err != nil {
						return err
					}

					// check if currency exists in map
					if _, ok := wallet[currency]; !ok {
						wallet[currency] = 0
					}

					currentAmount := wallet[currency]

					amountFloat, err := strconv.ParseFloat(amount, 64)
					if err != nil {
						return err
					}

					if amountFloat < 0 {
						return fmt.Errorf("amount must be positive")
					}

					if cCtx.Bool("set") {
						wallet[currency] = amountFloat
					} else {
						wallet[currency] = currentAmount + amountFloat
					}
					b, err := json.Marshal(wallet)
					if err != nil {
						return err
					}

					os.WriteFile(walletURL, b, 0644)

					fmt.Println("Previous holdings of", currency, ":", currentAmount)
					fmt.Println("New holdings of", currency, ":", wallet[currency])
					return nil

				},
			},
			{
				Name:    "sell",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "all"},
				},
				Usage: "sell a crypto currency",
				Action: func(cCtx *cli.Context) error {
					amount := cCtx.Args().First()

					var currency string
					if cCtx.Bool("all") {
						currency = strings.ToUpper(cCtx.Args().Get(0))
					} else {
						currency = strings.ToUpper(cCtx.Args().Get(1))
					}

					if currency == "" {
						return fmt.Errorf("please specify a currency")
					}

					if amount == "" && !cCtx.Bool("all") {
						return fmt.Errorf("please specify an amount")
					}

					var wallet datatypes.Wallet
					if err := wallet.GetFromJsonFile(); err != nil {
						return err
					}

					// check if currency exists in map
					if _, ok := wallet[currency]; !ok {
						wallet[currency] = 0
					}

					currentAmount := wallet[currency]

					if cCtx.Bool("all") {
						delete(wallet, currency)
					} else {
						amountFloat, err := strconv.ParseFloat(amount, 64)
						if err != nil {
							return err
						}

						if amountFloat < 0 {
							return fmt.Errorf("amount must be positive")
						}

						if currentAmount < amountFloat {
							return fmt.Errorf("you do not have enough %s to sell", currency)
						}
						wallet[currency] = currentAmount - amountFloat
					}
					b, err := json.Marshal(wallet)
					if err != nil {
						return err
					}

					os.WriteFile(walletURL, b, 0644)

					fmt.Println("Previous holdings of", currency, ":", currentAmount)
					fmt.Println("New holdings of", currency, ":", wallet[currency])
					return nil
				},
			},
			{
				Name:  "reset",
				Usage: "reset your wallet",
				Action: func(cCtx *cli.Context) error {
					var wallet datatypes.Wallet
					if err := wallet.Reset(); err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
