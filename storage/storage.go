package storage

import (
	"github.com/tendermint/iavl"
	tmdb "github.com/tendermint/tm-db"
)

type Storage struct {
	tree *iavl.MutableTree
}

func NewStorage(dataDir string) (*Storage, error) {
	lvldb, err := tmdb.NewGoLevelDB("treedb", dataDir)
	if err != nil {
		return nil, err
	}

	var sto Storage
	tree, err := iavl.NewMutableTree(lvldb, 0)
	if err != nil {
		return nil, err
	}
	sto.tree = tree

	return &sto, nil
}

func (sto *Storage) Set(k, v []byte) {
	sto.tree.Set(k, v)
}

func (sto *Storage) Get(k []byte) []byte {
	_, v := sto.tree.Get(k)
	return v
}
func (sto *Storage) State() []byte {
	return sto.tree.Hash()
}

func (sto *Storage) Commit() ([]byte, error) {
	h, _, err := sto.tree.SaveVersion()
	return h, err
}
