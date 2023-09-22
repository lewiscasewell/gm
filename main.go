package main

import (
	"crypto-price/api"
	"crypto-price/datatypes"
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
	var profile datatypes.Profile
	b, err := os.ReadFile("profile.json")
	if err != nil {
		newProfile := []byte(`{
			"name": "Chief",
			"token": "",
			"baseCurrency": "GBP"
		}`)

		os.WriteFile("profile.json", newProfile, 0644)
	}

	if err := json.Unmarshal(b, &profile); err != nil {
		return
	}

	app := &cli.App{
		Name:  "gm",
		Usage: "are you rich or broke today?",
		Action: func(cCtx *cli.Context) error {
			if profile.Token == "" {
				return fmt.Errorf("please set your currency api token")
			}
			var greeting string

			if time.Now().Hour() < 12 {
				greeting = "Good morning"
			} else if time.Now().Hour() < 18 {
				greeting = "Good afternoon"
			} else {
				greeting = "Good evening"
			}

			fmt.Printf("\n\n\n\n%s %s, let's check your crypto on this fine day.\n\n", greeting, profile.Name)

			exchangeRate := api.GetUsdGbpExchangeRate()

			var wg sync.WaitGroup

			total := 0.0

			var wallet datatypes.Wallet
			b, err := os.ReadFile("wallet.json")
			if err != nil {
				newEmptyWallet := []byte(`{}`)
				if err := os.WriteFile("wallet.json", newEmptyWallet, 0644); err != nil {
					return err
				}

				if err := json.Unmarshal(newEmptyWallet, &wallet); err != nil {
					return err
				}

				b, err = os.ReadFile("wallet.json")
				if err != nil {
					return err
				}
			}

			if err := json.Unmarshal(b, &wallet); err != nil {
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

					fmt.Printf("You have £%.2f worth of %s\n", r.Price*a*exchangeRate, s)
					total += r.Price * a * exchangeRate
				}(symbol, amount)
			}

			wg.Wait()

			defer fmt.Printf("\nTotal: £%.2f at %s\n", total, time.Now().Format("15:04:05"))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "get",
				Subcommands: []*cli.Command{
					{
						Name: "profile",
						Action: func(cCtx *cli.Context) error {
							fmt.Printf("Name: %s\n", profile.Name)
							fmt.Printf("Token: %s\n", profile.Token)
							fmt.Printf("Base Currency: %s\n", profile.BaseCurrency)
							return nil
						},
					},
					{
						Name: "wallet",
						Action: func(cCtx *cli.Context) error {
							var wallet datatypes.Wallet
							b, err := os.ReadFile("wallet.json")
							if err != nil {
								newEmptyWallet := []byte(`{}`)
								os.WriteFile("wallet.json", newEmptyWallet, 0644)
								wallet = make(datatypes.Wallet)
							} else {
								if err := json.Unmarshal(b, &wallet); err != nil {
									return err
								}
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
						Name:  "token",
						Usage: "sets the currency api token",
						Action: func(cCtx *cli.Context) error {
							profile.Token = cCtx.Args().First()

							b, err := json.Marshal(profile)
							if err != nil {
								return err
							}

							os.WriteFile("profile.json", b, 0644)
							return nil
						},
					},
					{
						Name:  "currency",
						Usage: "sets the base currency, e.g. GBP, USD, EUR",
						Action: func(cCtx *cli.Context) error {
							newBaseCurrency := strings.ToUpper(cCtx.Args().First())

							if newBaseCurrency != "GBP" && newBaseCurrency != "USD" && newBaseCurrency != "EUR" {
								return fmt.Errorf("base currency must be GBP, USD or EUR, not %s", newBaseCurrency)
							}

							profile.BaseCurrency = newBaseCurrency

							b, err := json.Marshal(profile)
							if err != nil {
								return err
							}

							os.WriteFile("profile.json", b, 0644)
							return nil
						},
					},
					{
						Name:  "name",
						Usage: "sets the name of the user",
						Action: func(cCtx *cli.Context) error {
							profile.Name = cCtx.Args().First()

							b, err := json.Marshal(profile)
							if err != nil {
								return err
							}

							os.WriteFile("profile.json", b, 0644)
							return nil
						},
					},
				},
			},
			{
				Name:  "short",
				Usage: "complete a task on the list",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "serve", Aliases: []string{"s"}},
					&cli.BoolFlag{Name: "option", Aliases: []string{"o"}},
					&cli.StringFlag{Name: "message", Aliases: []string{"m"}},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("serve:", cCtx.Bool("serve"))
					fmt.Println("option:", cCtx.Bool("option"))
					fmt.Println("message:", cCtx.String("message"))
					return nil
				},
			},
			{
				Name:    "buy",
				Aliases: []string{"b"},
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "hard"},
				},
				Usage: "buy a crypto currency",
				Action: func(cCtx *cli.Context) error {
					amount := cCtx.Args().First()
					currency := strings.ToUpper(cCtx.Args().Get(1))

					if amount == "" || currency == "" {
						return fmt.Errorf("please specify an amount and currency")
					}

					var wallet datatypes.Wallet
					b, err := os.ReadFile("wallet.json")
					if err != nil {
						newEmptyWallet := []byte(`{}`)

						os.WriteFile("wallet.json", newEmptyWallet, 0644)
						wallet = make(datatypes.Wallet)

						if err := json.Unmarshal(newEmptyWallet, &wallet); err != nil {
							return err
						}

						b, err = os.ReadFile("wallet.json")
						if err != nil {
							return err
						}
					}

					if err := json.Unmarshal(b, &wallet); err != nil {
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

					if cCtx.Bool("hard") {
						wallet[currency] = amountFloat
					} else {
						wallet[currency] = currentAmount + amountFloat
					}
					b, err = json.Marshal(wallet)
					if err != nil {
						return err
					}

					os.WriteFile("wallet.json", b, 0644)

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
					b, err := os.ReadFile("wallet.json")
					if err != nil {
						newEmptyWallet := []byte(`{}`)

						os.WriteFile("wallet.json", newEmptyWallet, 0644)
						wallet = make(datatypes.Wallet)

						if err := json.Unmarshal(newEmptyWallet, &wallet); err != nil {
							return err
						}

						b, err = os.ReadFile("wallet.json")
						if err != nil {
							return err
						}
					}

					if err := json.Unmarshal(b, &wallet); err != nil {
						return err
					}

					// check if currency exists in map
					if _, ok := wallet[currency]; !ok {
						wallet[currency] = 0
					}

					currentAmount := wallet[currency]

					if cCtx.Bool("all") {
						wallet[currency] = 0
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
					b, err = json.Marshal(wallet)
					if err != nil {
						return err
					}

					os.WriteFile("wallet.json", b, 0644)

					fmt.Println("Previous holdings of", currency, ":", currentAmount)
					fmt.Println("New holdings of", currency, ":", wallet[currency])
					return nil

				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}