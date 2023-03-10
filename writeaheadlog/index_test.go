// START: intro
package writeaheadlog

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	idx, err := newIndex(f, c)
	require.NoError(t, err)
	_, err = idx.GetOffset()
	require.Equal(t, f.Name(), idx.Name())
	require.ErrorIs(t, err, io.EOF)

	entries := []struct {
		Off uint32
		Pos uint64
	}{
		{Off: 0, Pos: 0},
		{Off: 1, Pos: 10},
	}

	for _, ent := range entries {
		err = idx.Write(ent.Off, ent.Pos)
		require.NoError(t, err)

		pos, err := idx.Read(ent.Off)
		require.NoError(t, err)
		require.Equal(t, ent.Pos, pos)
	}

	// index and scanner should error when reading exceed existing entries
	_, err = idx.Read(uint32(len(entries)))
	require.ErrorIs(t, err, io.EOF)
	_ = idx.Close()

	// index should build its state from the existing file
	f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
	idx, err = newIndex(f, c)
	require.NoError(t, err)
	off, err := idx.GetOffset()
	require.NoError(t, err)
	require.Equal(t, uint32(1), off)
	pos, err := idx.Read(off)
	require.Equal(t, entries[1].Pos, pos)
	require.NoError(t, err)
}
