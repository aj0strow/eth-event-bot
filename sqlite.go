package main

import (
	"database/sql"
)

func initDatabase(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS events (
  network text NOT NULL,
  contract_name text NOT NULL,
  address text NOT NULL,
  block_hash text NOT NULL,
  block_number int NOT NULL,
  transaction_hash text NOT NULL,
  transaction_index int NOT NULL,
  log_index int NOT NULL,
  event_name text NOT NULL,
  arg0 text,
  arg1 text,
  arg2 text,
  arg3 text,
  arg4 text,
  arg5 text,
  arg6 text,
  arg7 text,
  arg8 text,
  arg9 text
);`)
	return err
}

func queryLastBlock(db *sql.DB, network string, address string) (sql.NullInt64, error) {
	row := db.QueryRow(`SELECT block_number FROM events WHERE network = ? AND address = ? ORDER BY block_number DESC LIMIT 1`, network, address)
	var blockNumber sql.NullInt64
	err := row.Scan(&blockNumber)
	if err == sql.ErrNoRows {
		return sql.NullInt64{}, nil
	} else if err != nil {
		return sql.NullInt64{}, err
	}
	return blockNumber, nil
}

type sqlEvent struct {
	Network          string
	ContractName     string
	Address          string
	BlockHash        string
	BlockNumber      int64
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
	EventName        string
	Arguments        []string
}

func queryEventExists(db *sql.DB, event sqlEvent) (bool, error) {
	row := db.QueryRow(`SELECT 1 FROM events WHERE network = ? AND transaction_hash = ? AND log_index = ?`, event.Network, event.TransactionHash, event.LogIndex)
	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func insertEvent(db *sql.DB, event sqlEvent) error {
	exists, err := queryEventExists(db, event)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	args := packArgs(event.Arguments)
	_, err = db.Exec(`
		INSERT INTO events (network, contract_name, address, block_hash, block_number, transaction_hash, transaction_index, log_index, event_name, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.Network,
		event.ContractName,
		event.Address,
		event.BlockHash,
		event.BlockNumber,
		event.TransactionHash,
		event.TransactionIndex,
		event.LogIndex,
		event.EventName,
		args[0],
		args[1],
		args[2],
		args[3],
		args[4],
		args[5],
		args[6],
		args[7],
		args[8],
		args[9],
	)
	return err
}

func packArgs(args []string) [10]sql.NullString {
	flatArgs := [10]sql.NullString{}
	for i := 0; i < len(flatArgs) && i < len(args); i++ {
		flatArgs[i] = sql.NullString{
			Valid:  true,
			String: args[i],
		}
	}
	return flatArgs
}

func unpackArgs(flatArgs [10]sql.NullString) []string {
	var args []string
	for _, flatArg := range flatArgs {
		if flatArg.Valid {
			args = append(args, flatArg.String)
		}
	}
	return args
}
