package config

type Config interface {
	LoadConfig(filePath string) (Config, error)
}
