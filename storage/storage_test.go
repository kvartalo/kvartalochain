package storage

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	tmpDir, err := ioutil.TempDir("./", "tmpTest")
	require.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	sto, err := NewStorage(tmpDir)
	assert.Nil(t, err)

	sto.Set([]byte("test0"), []byte("value0"))
	assert.Equal(t, []byte("value0"), sto.Get([]byte("test0")))

	sto.Commit()
	assert.Equal(t, "c778aafd61b926abbfb8a8d6c7d8727bcbc069207a67ba9a26fefe71cd155ae5", hex.EncodeToString(sto.State()))
	assert.Equal(t, []byte("value0"), sto.Get([]byte("test0")))
}
