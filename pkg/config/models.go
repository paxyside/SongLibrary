package config

type config struct {
	Database    dbConfig     `json:"database"`
	Server      serverConfig `json:"server"`
	Credentials credentials  `json:"credentials"`
	External    external     `json:"external"`
}

type dbConfig struct {
	URI string `json:"uri"`
}

type serverConfig struct {
	ServerUrl string `json:"server_url"`
	Host      string `json:"host"`
	Port      string `json:"port"`
}

type credentials struct {
	ApiKey string `json:"api_key"`
}

type external struct {
	ExtApiUrl string `json:"ext_api_url"`
}

var Config config

func (c *config) DatabaseURI() string {
	return c.Database.URI
}

func (c *config) ServerURI() string {
	return c.Server.ServerUrl
}

func (c *config) ServerHost() string {
	return c.Server.Host
}

func (c *config) ServerPort() string {
	return c.Server.Port
}

func (c *config) ApiKey() string {
	return c.Credentials.ApiKey
}

func (c *config) ExternalApiUrl() string {
	return c.External.ExtApiUrl
}
