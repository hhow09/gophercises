package writeaheadlog

import (
	"bufio"
	"encoding/binary"
	"errors"
	"os"
	"sync"
)

// bigendian: order in which the "big end" (most significant value in the sequence) is stored first, at the lowest storage address.
var enc = binary.BigEndian

var (
	// ErrNotFound is returned when an entry is not found.
	ErrNotFound = errors.New("not found")
)

const (
	// byte size of len mark
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.RWMutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name()) // file info
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// returns
// n: total bytes written
// offset: write offset
func (s *store) Append(p []byte) (n uint64, offset uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	offset = s.size // record offset

	// 1. write record length to buffer (lenWidth)
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	// 2. write record
	nn, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	// 3. total written = nn + lenWidth
	nn += lenWidth
	s.size += uint64(nn)
	return uint64(nn), offset, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if pos > s.size {
		return nil, ErrNotFound
	}
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}
