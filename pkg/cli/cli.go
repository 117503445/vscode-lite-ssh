package cli

import (
	"regexp"

	"github.com/117503445/goutils"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
	"github.com/rs/zerolog/log"
)

type SshNode struct {
	Host string `help:"node host"`
	Port int    `help:"node port"`
	User string `help:"node user"`
	Pri  string `help:"node private key path"`
}

var Cli struct {
	Nodes map[string]SshNode `help:"nodes"`
}

func cfgCheck() {
	if len(Cli.Nodes) == 0 {
		log.Fatal().Msg("nodes is empty")
	}

	// name must only contain [a-zA-Z0-9_]
	isValidName := func(name string) bool {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, name)
		return matched
	}

	for name, node := range Cli.Nodes {
		logger := log.With().Str("node", name).Logger()
		if name == "" {
			logger.Fatal().Msg("node name is empty")
		} else {
			if !isValidName(name) {
				logger.Fatal().Msg("node name is invalid")
			}
		}

		if node.Host == "" {
			logger.Fatal().Msg("node host is empty")
		}
		if node.Port == 0 {
			logger.Warn().Msg("node port is empty, use default 22")
		}
		if node.User == "" {
			logger.Warn().Msg("node user is empty, use default root")
		}

		if node.Pri == "" {
			logger.Fatal().Msg("node pri is empty")
		} else {
			if !goutils.FileExists(node.Pri) {
				logger.Fatal().Msg("node pubkey path is not exists")
			}
		}
	}
}

func cfgSetDefault() {
	newNodes := make(map[string]SshNode)
	for name, node := range Cli.Nodes {
		newNodes[name] = node
		if node.Port == 0 || node.User == "" {
			port := 22
			if node.Port != 0 {
				port = node.Port
			}

			user := "root"
			if node.User != "" {
				user = node.User
			}

			newNode := SshNode{
				Host: node.Host,
				Port: port,
				User: user,
				Pri:  node.Pri,
			}
			newNodes[name] = newNode
		}
	}
	for name, node := range newNodes {
		Cli.Nodes[name] = node
	}
}

func CfgLoad() {
	kong.Parse(&Cli, kong.Configuration(kongtoml.Loader, "/workspace/config.toml"))
	cfgCheck()
	cfgSetDefault()
	log.Info().Interface("cli", Cli).Msg("")
}
