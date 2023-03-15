package writeaheadlog

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	api "github.com/hhow09/gophercises/writeaheadlog/api/v1"
)

type Log struct {
	mu sync.RWMutex

	Dir    string
	Config Config

	activeSegment *segment
	segments      []*segment
}

func NewLog(dir string, c Config) (*Log, error) {
	if c.Segment.MaxStoreBytes == 0 {
		c.Segment.MaxStoreBytes = 1024
	}
	if c.Segment.MaxIndexBytes == 0 {
		c.Segment.MaxIndexBytes = 1024
	}
	l := &Log{
		Dir:    dir,
		Config: c,
	}

	return l, l.setup()
}

// read the exisitng file into segements in order to Read()
func (l *Log) setup() error {
	files, err := ioutil.ReadDir(l.Dir)
	if err != nil {
		return err
	}
	offSetMap := map[uint32]int{}
	for _, file := range files {
		offset, err := parseFileOffset(file)
		if err != nil {
			continue
		}
		offSetMap[offset] += 1
	}
	var baseOffsets []uint32

	for off, count := range offSetMap {
		if count == 2 {
			baseOffsets = append(baseOffsets, off)
		}
	}

	sort.Slice(baseOffsets, func(i, j int) bool {
		return baseOffsets[i] < baseOffsets[j]
	})

	for i := 0; i < len(baseOffsets); i++ {
		if err = l.newSegment(baseOffsets[i]); err != nil {
			return err
		}
	}
	if l.segments == nil {
		if err = l.newSegment(l.Config.Segment.InitialOffset); err != nil {
			return err
		}
	}
	return nil
}

func (l *Log) newSegment(off uint32) error {
	s, err := newSegment(l.Dir, off, l.Config)
	if err != nil {
		return err
	}
	l.segments = append(l.segments, s)
	l.activeSegment = s
	return nil
}

// append record to log
// return the (offset, error)
func (l *Log) Append(record *api.Record) (uint32, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	off, err := l.activeSegment.Append(record)
	if err != nil {
		return 0, err
	}
	if l.activeSegment.IsFull() {
		err = l.newSegment(off + 1)
	}
	return off, err
}

// read the record given offset
func (l *Log) Read(off uint32) (*api.Record, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	idx := sort.Search(len(l.segments), func(i int) bool { return l.segments[i].baseOffset >= off })
	if idx >= len(l.segments) {
		return nil, fmt.Errorf("offset out of range: %d", off)
	}
	if l.segments[idx].nextOffset <= off {
		return nil, fmt.Errorf("offset out of range: %d", off)
	}

	return l.segments[idx].Read(off)
}

// close the file discriptor of index and store
func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, segment := range l.segments {
		if err := segment.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Log) MinOffset() (uint32, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.segments[0].baseOffset, nil
}

func (l *Log) MaxOffset() (uint32, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	off := l.segments[len(l.segments)-1].nextOffset
	if off == 0 {
		return 0, nil
	}
	return off - 1, nil
}

// return a reader of all records
func (l *Log) Reader() io.Reader {
	l.mu.RLock()
	defer l.mu.RUnlock()
	readers := make([]io.Reader, len(l.segments))
	for i, segment := range l.segments {
		segment.store.Flush()
		readers[i] = &fileReader{segment.store, 0}
	}
	return io.MultiReader(readers...)
}

type fileReader struct {
	*store
	off int64
}

// implement io.Reader
func (o *fileReader) Read(p []byte) (int, error) {
	n, err := o.ReadAt(p, o.off)

	o.off += int64(n)
	return n, err
}

// remove segements which offset less than given lowest offset.
func (l *Log) Truncate(lowest uint32) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	var segments []*segment
	for _, s := range l.segments {
		if s.nextOffset <= lowest+1 {
			if err := s.Remove(); err != nil {
				return err
			}
			continue
		}
		segments = append(segments, s)
	}
	l.segments = segments
	return nil
}

func parseFileOffset(file fs.FileInfo) (uint32, error) {
	offStr := strings.TrimSuffix(
		file.Name(),
		path.Ext(file.Name()),
	)
	off, err := strconv.ParseUint(offStr, 10, 0)
	return uint32(off), err
}
