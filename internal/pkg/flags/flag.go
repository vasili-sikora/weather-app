package flags

import "flag"

type flags struct {
	Path string
}

func Parse() *flags {
	configPath := flag.String("config", "./config/config.yaml", "path to config")

	flag.Parse()

	return &flags{
		Path: *configPath,
	}
}
