package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/mattn/go-sqlite3"
	"math/big"
	"os"
)

const version = "0.0.2"

var (
	wantVersion bool
	configFile  string
)

func main() {
	// Parse flags.
	flag.BoolVar(&wantVersion, "version", false, "print version")
	flag.StringVar(&configFile, "C", "./config.toml", "path to configuration file")
	flag.Parse()

	// Setup command context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Print version.
	if wantVersion {
		printVersion()
	}

	// Read configuration file.
	config := &BotOptions{}
	_, err := toml.DecodeFile(configFile, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %s\n", err.Error())
		os.Exit(1)
	}

	// Open database connection.
	db, err := sql.Open("sqlite3", config.Sqlite.Connection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sqlite: %s\n", err.Error())
		return
	}
	defer db.Close()

	// Set up infura api.
	rpcUrl := config.Infura.Endpoint()
	ethRpc, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}
	defer ethRpc.Close()

	// Parse ABI definitions.
	var contracts []*Contract
	for _, cc := range config.Contract {
		contract, err := ParseTruffleBuild(cc.TruffleJSON)
		if err != nil {
			panic(err)
		}
		contract.Address = cc.Address
		contracts = append(contracts, contract)
	}
	for _, contract := range contracts {
		fmt.Printf("%s\n", contract.ContractName)
		for _, event := range contract.ABI.Events {
			fmt.Printf("  %s\n", event.String())
		}
	}

	// Initialize database.
	err = initDatabase(db)
	if err != nil {
		panic(err)
	}

	// Set up integrations.
	var integrations []Integration
	if config.Telegram != nil {
		integrations = append(integrations, &Telegram{
			Token:  config.Telegram.Token,
			ChatID: config.Telegram.ChatID,
		})
	}

	// Get logs for each contract.
	for _, contract := range contracts {
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(0),
			Addresses: []common.Address{
				common.HexToAddress(contract.Address),
			},
		}
		logs, err := ethRpc.FilterLogs(ctx, query)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: %d events\n", contract.ContractName, len(logs))
		for _, log := range logs {
			event, ok := contract.GetEvent(log.Topics[0])
			if !ok {
				continue
			}
			values, err := UnpackValues(event, log)
			if err != nil {
				panic(err)
			}
			args := make([]string, len(values))
			for i, v := range values {
				args[i] = PrintValue(v)
			}
			row := sqlEvent{
				Network:          config.Infura.Network,
				ContractName:     contract.ContractName,
				Address:          contract.Address,
				BlockHash:        log.BlockHash.String(),
				BlockNumber:      int64(log.BlockNumber),
				TransactionHash:  log.TxHash.String(),
				TransactionIndex: int64(log.TxIndex),
				LogIndex:         int64(log.Index),
				EventName:        event.Name,
				Arguments:        args,
			}
			err = insertEvent(db, row)
			if err != nil {
				panic(err)
			}
			for _, integration := range integrations {
				notification := sqlNotification{
					Network:         row.Network,
					TransactionHash: row.TransactionHash,
					LogIndex:        row.LogIndex,
					Platform:        integration.Platform(),
				}
				exists, err := queryNotificationExists(db, notification)
				if err != nil {
					panic(err)
				}
				if exists {
					continue
				}
				err = integration.Send(row)
				if err != nil {
					fmt.Printf("%s: %s\n", integration.Platform(), err.Error())
					continue
				}
				err = insertNotification(db, notification)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func printVersion() {
	fmt.Fprintf(os.Stdout, "eth-event-bot v%s\n", version)
	os.Exit(0)
}
