# Consistent Hashing

## Requirement
implement `ConsistentHashInterface` defined in `interfaces.go`
```golang
type ConsistentHashInterface interface {
	Add(node AnyType)
	AddWithReplicas(node AnyType, replicas int)
	AddWithWeight(node AnyType, weight int)
	Get(key AnyType) (AnyType, bool)
	Remove(node AnyType)
}
```


## Murmur Hash
[The murmur3 hash function: hashtables, bloom filters, hyperloglog](https://www.sderosiaux.com/articles/2017/08/26/the-murmur3-hash-function--hashtables-bloom-filters-hyperloglog/)


## Reference
- [Kevin Wan: Consistent Hash Algorithm and Go Implementation](https://faun.pub/consistent-hash-algorithm-and-go-implementation-a5a01d84845a)
- [zeromicro/go-zero](https://github.com/zeromicro/go-zero/blob/master/core/hash/consistenthash.go)
