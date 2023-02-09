package hash

type Node AnyType

type ConsistentHashInterface interface {
	Add(node AnyType)
	AddWithReplicas(node AnyType, replicas int)
	AddWithWeight(node AnyType, weight int)
	Get(key AnyType) (AnyType, bool)
	Remove(node AnyType)
}
