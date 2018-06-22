# Ethereum Event Bot

Install `eth-event-bot` using the go programming language.

```
$ go get github.com/aj0strow/eth-event-bot
```

See the configuration file format below. The bot pulls all the latest events and stores them into a sqlite database for easy access and query abilities. 

## Config File

Example file.

```toml
[sqlite]
connection = "file::memory:?mode=memory"

[infura]
network = "rinkeby"
apikey = "insert_api_key"

[[contract]]
truffleJson = "./contracts/Contract1.json"
address = "0xa7a4f43c7c174a5a758fce548582413040Ab4134"

[[contract]]
truffleJson = "./contracts/Contract2.json"
address = "0x631b6C1A40AB37DA459e823C27aFb73b7f984e0e"
```

Further notes:

* See [sqlite connection string](https://github.com/mattn/go-sqlite3#connection-string) to format the connection string to store the database as a file on your computer. 

* See [sqlite browser app](https://sqlitebrowser.org/) for a way to easily view data in a sqlite database table. 

## License

Ethereum Event Bot is free as in libre software. See the `LICENSE.txt` file. 
