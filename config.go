package main

type botConfig struct {
	Sqlite   sqliteConfig
	Infura   infuraConfig
	Contract []contractConfig
}

type sqliteConfig struct {
	Connection string
}

type infuraConfig struct {
	APIKey  string
	Network string
}

func (infura *infuraConfig) Endpoint() string {
	return "https://" + infura.Network + ".infura.io/" + infura.APIKey
}

type contractConfig struct {
	TruffleJSON string
	Address     string
}
