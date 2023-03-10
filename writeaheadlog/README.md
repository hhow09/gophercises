# Write Ahead Log

## Components
- Record: the data stored in our log.
- Store: the file store records in.
- Index: the file we store index entries in.
- Segment: the abstraction that ties a store and an index together.
- Log: the abstraction that ties all the segments together.

### Hierarchy
- Log = N * Segment
- Segment = Store + Index
- Store = File = M * Record


## gommap
http://labix.org/gommap

## Ref
- [Distributed Services with Go](https://www.oreilly.com/library/view/distributed-services-with/9781680508376/) `Chapter 3. Write a Log Package`
- [tidwall/wal](https://github.com/tidwall/wal)
- [Patterns of Distributed Systems: Write Ahead Log](https://martinfowler.com/articles/patterns-of-distributed-systems/wal.html)