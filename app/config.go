package main

// Config struct to store parsed arguments
type Config struct {
	Dir        string
	DbFileName string
}

func (c *Config) DbFilePath() string {
	return c.Dir + "/" + c.DbFileName
}
