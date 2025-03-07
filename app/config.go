package main

// Config struct to store parsed arguments
type Config struct {
	Dir        string
	DbFileName string
	Port       uint
}

func (c *Config) DbFilePath() string {
	return c.Dir + "/" + c.DbFileName
}
