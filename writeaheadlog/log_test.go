package writeaheadlog

import (
	"io/ioutil"
	"os"
	"testing"

	api "github.com/hhow09/gophercises/writeaheadlog/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestLog(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T, log *Log,
	){
		"append and read a record succeeds": testAppendRead,
		"offset out of range error":         testOutOfRangeErr,
		"init with existing segments":       testInitExisting,
		"reader":                            testReader,
		"truncate":                          testTruncate,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "log-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			c := Config{}
			c.Segment.MaxStoreBytes = 32
			log, err := NewLog(dir, c)
			require.NoError(t, err)
			defer log.Close()

			fn(t, log)
		})
	}
}

func testAppendRead(t *testing.T, log *Log) {
	append := randomRecord()
	off, err := log.Append(append)
	require.NoError(t, err)
	require.Equal(t, uint32(0), off)

	read, err := log.Read(off)
	require.NoError(t, err)
	require.Equal(t, append.Value, read.Value)
}

func testOutOfRangeErr(t *testing.T, log *Log) {
	read, err := log.Read(1)
	require.Nil(t, read)
	require.Error(t, err)
}

func testInitExisting(t *testing.T, o *Log) {
	append := randomRecord()
	for i := 0; i < 3; i++ {
		_, err := o.Append(append)
		require.NoError(t, err)
	}
	require.NoError(t, o.Close())

	off, err := o.MinOffset()
	require.NoError(t, err)
	require.Equal(t, uint32(0), off)
	off, err = o.MaxOffset()
	require.NoError(t, err)
	require.Equal(t, uint32(2), off)

	n, err := NewLog(o.Dir, o.Config)
	require.NoError(t, err)

	off, err = n.MinOffset()
	require.NoError(t, err)
	require.Equal(t, uint32(0), off)
	off, err = n.MaxOffset()
	require.NoError(t, err)
	require.Equal(t, uint32(2), off)
}

func testReader(t *testing.T, log *Log) {
	append := randomRecord()
	off, err := log.Append(append)
	require.NoError(t, err)
	require.Equal(t, uint32(0), off)

	reader := log.Reader()
	b, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	read := &api.Record{}
	err = proto.Unmarshal(b[lenWidth:], read)
	require.NoError(t, err)
	require.Equal(t, append.Value, read.Value)
}

func testTruncate(t *testing.T, log *Log) {
	append := randomRecord()
	for i := 0; i < 3; i++ {
		_, err := log.Append(append)
		require.NoError(t, err)
	}

	err := log.Truncate(1)
	require.NoError(t, err)

	_, err = log.Read(0)
	require.Error(t, err)
}
