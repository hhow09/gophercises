package writeaheadlog

import (
	"fmt"
	"os"
	"path"

	api "github.com/hhow09/gophercises/writeaheadlog/api/v1"
	"google.golang.org/protobuf/proto"
)

type segment struct {
	store                  *store
	index                  *index
	baseOffset, nextOffset uint32
	config                 Config
}

func newSegment(dir string, baseOffset uint32, config Config) (*segment, error) {
	s := &segment{
		baseOffset: baseOffset,
		config:     config,
	}
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.store, err = newStore(storeFile); err != nil {
		return nil, err
	}
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.index, err = newIndex(indexFile, config); err != nil {
		return nil, err
	}
	if offset, err := s.index.GetOffset(); err != nil {
		s.nextOffset = baseOffset
	} else {
		s.nextOffset = baseOffset + offset + 1
	}
	return s, nil
}

// append record to the store file and update index
func (s *segment) Append(record *api.Record) (offset uint32, err error) {
	cur := s.nextOffset
	record.Offset = cur
	p, err := proto.Marshal(record)
	if err != nil {
		return 0, err
	}
	_, pos, err := s.store.Append(p)
	if err != nil {
		return 0, err
	}
	if err = s.index.Write(
		// index offsets are relative to base offset
		s.nextOffset-s.baseOffset,
		pos,
	); err != nil {
		return 0, err
	}
	s.nextOffset++
	return cur, nil
}

func (s *segment) Read(off uint32) (*api.Record, error) {
	pos, err := s.index.Read(off - s.baseOffset)
	if err != nil {
		return nil, err
	}
	p, err := s.store.Read(pos)
	if err != nil {
		return nil, err
	}
	var record api.Record
	err = proto.Unmarshal(p, &record)
	return &record, err
}

// close the index and store file discriptor
func (s *segment) Close() error {
	if s.index != nil {
		if err := s.index.Close(); err != nil {
			return err
		}
	}
	if s.store != nil {
		if err := s.store.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (s *segment) IsFull() bool {
	return s.store.size >= s.config.Segment.MaxStoreBytes ||
		s.index.size >= s.config.Segment.MaxIndexBytes
}

// remove the underlying index and store file
func (s *segment) Remove() error {
	if err := s.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.index.Name()); err != nil {
		return err
	}
	if err := os.Remove(s.store.Name()); err != nil {
		return err
	}
	return nil
}
