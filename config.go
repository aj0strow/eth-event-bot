package main

type BotOptions struct {
	Sqlite   *SqliteOptions
	Infura   *InfuraOptions
	Contract []*ContractOptions
	Telegram *TelegramOptions
}

type SqliteOptions struct {
	Connection string
}

type InfuraOptions struct {
	APIKey  string
	Network string
}

func (infura *InfuraOptions) Endpoint() string {
	return "https://" + infura.Network + ".infura.io/" + infura.APIKey
}

type ContractOptions struct {
	TruffleJSON string
	Address     string
}

type TelegramOptions struct {
	Token  string
	ChatID string
}
