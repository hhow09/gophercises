package writeaheadlog

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupDir(i int) string {
	dir, _ := ioutil.TempDir("", fmt.Sprintf("segment-test-%d", i))
	return dir
}

func randomRecord() Record {
	rand.Seed(time.Now().UnixNano())
	return Record{
		Value: []byte(fmt.Sprintf("hello world %d", rand.Int())),
	}
}

func TestSegment(t *testing.T) {
	dir := setupDir(0)
	defer os.RemoveAll(dir)

	c := Config{}
	c.Segment.MaxStoreBytes = 1024
	c.Segment.MaxIndexBytes = totalWidth * 3

	s, err := newSegment(dir, 16, c)
	require.NoError(t, err)
	require.Equal(t, uint32(16), s.nextOffset, s.nextOffset)
	require.False(t, s.IsFull())

	for i := uint32(0); i < 3; i++ {
		want := randomRecord()
		off, err := s.Append(want)
		require.NoError(t, err)
		require.Equal(t, 16+i, off)

		got, err := s.Read(off)
		require.NoError(t, err)
		require.Equal(t, want.Value, got.Value)
	}

	want := randomRecord()
	_, err = s.Append(want)
	require.Equal(t, io.EOF, err)

	// maxed index
	require.True(t, s.IsFull())

	c.Segment.MaxStoreBytes = uint64(len(want.Value) * 3)
	c.Segment.MaxIndexBytes = 1024
	fmt.Println("h", c.Segment.MaxStoreBytes)

	s, err = newSegment(dir, 16, c)
	require.NoError(t, err)
	// maxed store
	require.True(t, s.IsFull())

	// remove segement and recreate new one
	err = s.Remove()
	require.NoError(t, err)
	s, err = newSegment(dir, 16, c)
	require.NoError(t, err)
	require.False(t, s.IsFull())
}
