# Write Ahead Log
- write ahead log implemtation with go and protobuf.
- index is used to allow `O(logN)` search.

## API
- `Append()`
- `Read(offset uint32)`
- `Close()`

## Basic Usage
```go
os.Mkdir(log_dir, mode) // log directory
log, err := NewLog(log_dir, Config{})
defer log.Close()
record := &api.Record{
		Value: []byte("hello world"),
	}
off, err := log.Append(append)
read, err := log.Read(off) // {Value: "hello world!", Offset: 0}
```

## Components
- Record: the data stored in file.
- Store: the file store records in.
- Index: the file we store index entries in.
- Segment: the abstraction that ties a store and an index together.
- Log: the abstraction that ties all the segments together.

### Hierarchy
- Log = N * Segment
- Segment = Store + Index
- Store = File = M * Record

## Example
### a `Record`
```json
{"Offset": 15, "Value": "Hello World"}
```

### a Index entry
```
|  15        300     |
  uint32    uint64
  4 byte    8 byte
  offset    position
```

## Relative Index Offset
- in order to save storage (`uint32`), the index entry offset is the relative offset to `baseOffset` of index file.
- index file name: `{baseOffset}.index`


## [gommap](http://labix.org/gommap)
- [mmap(2) — Linux manual page](https://man7.org/linux/man-pages/man2/mmap.2.html)
- directly work with memory mapped files
- ref: [Discovering and exploring mmap using Go](https://brunocalza.me/discovering-and-exploring-mmap-using-go/)

## Development
1. [install protobuf](https://grpc.io/docs/protoc-installation/)
2. generate protobuf go code `make comple`

## Ref
- [Distributed Services with Go](https://www.oreilly.com/library/view/distributed-services-with/9781680508376/) `Chapter 3. Write a Log Package`
- [tidwall/wal](https://github.com/tidwall/wal)
- [Patterns of Distributed Systems: Write Ahead Log](https://martinfowler.com/articles/patterns-of-distributed-systems/wal.html)