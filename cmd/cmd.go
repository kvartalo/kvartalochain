package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"kvartalochain/endpoint"

	"github.com/pkg/errors"
	// log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	log "github.com/tendermint/tendermint/libs/log"
	"github.com/urfave/cli"
)

var config = cfg.DefaultConfig()
var logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))

var ServerCommands = []cli.Command{
	{
		Name:    "initNode",
		Aliases: []string{},
		Usage:   "initialize the database",
		Action:  cmdInitNode,
	},
	{
		Name:    "initChain",
		Aliases: []string{},
		Usage:   "initialize genesis",
		Action:  cmdInitChain,
	},
	{
		Name:    "start",
		Aliases: []string{},
		Usage:   "start the server",
		Action:  cmdStart,
	},
	{
		Name:    "info",
		Aliases: []string{},
		Usage:   "get info about the server",
		Action:  cmdInfo,
	},
}

func cmdInitNode(c *cli.Context) error {
	err := initNode(config)
	return err
}
func cmdInitChain(c *cli.Context) error {
	err := initGenesis(config)
	return err
}

func cmdStart(c *cli.Context) error {
	configFile := "tmp/config/config.toml"

	// read config
	config.RootDir = filepath.Dir(filepath.Dir(configFile))
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "viper failed to read config file")
	}
	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrap(err, "viper failed to unmarshal config")
	}
	if err := config.ValidateBasic(); err != nil {
		return errors.Wrap(err, "config is invalid")
	}

	node, db, archiveDb := loadTendermint(config.DBPath)

	go func() {
		apiservice := endpoint.Serve(db, archiveDb)
		apiservice.Run(":" + "3000")
		logger.Info("api server running at :" + "3000")
	}()

	fmt.Println("starting node")
	node.Start()
	defer func() {
		node.Stop()
		node.Wait()
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	os.Exit(0)

	return nil
}

func cmdInfo(c *cli.Context) error {
	logger.Info("info")

	return nil
}
