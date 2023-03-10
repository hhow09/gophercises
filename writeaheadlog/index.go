package writeaheadlog

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

const (
	offWidth   uint64 = 4
	posWidth   uint64 = 8
	totalWidth        = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())
	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}
	// creates a new mapping in the virtual address space of the calling process.
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}
	return idx, nil
}

// write to offset and position and update the size
func (i *index) Write(offset uint32, pos uint64) error {
	if i.size+totalWidth > uint64(len(i.mmap)) {
		return io.EOF
	}
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], offset)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+totalWidth], pos)
	i.size += uint64(totalWidth)
	return nil
}

// takes in an offset and returns the associated record’s position in the store.
func (i *index) Read(offset uint32) (pos uint64, err error) {
	if i.size == 0 {
		return 0, io.EOF
	}
	idxPos := uint64(offset) * totalWidth
	if idxPos+totalWidth > i.size {
		return 0, io.EOF
	}
	// out = enc.Uint32(i.mmap[idxPos : idxPos+offWidth])
	pos = enc.Uint64(i.mmap[idxPos+offWidth : idxPos+totalWidth])
	return pos, nil
}

// returns the last offset from index file
func (i *index) GetOffset() (offset uint32, err error) {
	if i.size == 0 {
		return 0, io.EOF
	}
	idxPos := (i.size/totalWidth - 1) * totalWidth
	if idxPos+totalWidth > i.size {
		return 0, io.EOF
	}
	offset = enc.Uint32(i.mmap[idxPos : idxPos+offWidth])
	return offset, nil
}

func (i *index) Name() string {
	return i.file.Name()
}

func (i *index) Close() error {
	if len(i.mmap) == 0 {
		return nil
	}
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	return i.file.Close()
}
