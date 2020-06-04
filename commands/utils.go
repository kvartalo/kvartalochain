package commands

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"

	"kvartalochain/chain"
	"kvartalochain/storage"

	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

func loadTendermint(dbpath string) (*nm.Node, *storage.Storage, *badger.DB) {
	fmt.Println("PATH", dbpath)
	db, err := storage.NewStorage(dbpath)
	if err != nil {
		logger.Error("failed to open storage db: %v", err)
		os.Exit(1)
	}

	archiveDb, err := badger.Open(badger.DefaultOptions(dbpath))
	if err != nil {
		logger.Error("failed to open badger db: %v", err)
		os.Exit(1)
	}
	// defer db.Close()

	app := chain.NewKvartaloApplication(db, archiveDb)

	node, err := newTendermint(app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
	return node, db, archiveDb
}

// func newTendermint(app abci.Application, configFile string) (*nm.Node, error) {
func newTendermint(app abci.Application) (*nm.Node, error) {
	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	var err error
	logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log level")
	}
	// read private validator
	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)
	// read node key
	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, errors.Wrap(err, "failed to load node's key")
	}
	// create node
	node, err := nm.NewNode(
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new Tendermint node")
	}
	return node, nil
}
