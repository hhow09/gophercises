package hash

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// Placeholder is a placeholder object that can be used globally.

const (
	// top weight that one entry might set.
	TopWeight = 100

	minReplicas = 100
	prime       = 16777619
)

// assert it implments the interface
var _ ConsistentHashInterface = (*ConsistentHash)(nil)

type (
	Func func(data []byte) uint64

	ConsistentHash struct {
		hashFunc Func
		// the number of virtual nodes of the node
		replicas int
		// list of all virtual nodes
		keys []uint64
		// virtual node -> []physical nodes mapping
		// stores as slice to preventd hash conflict
		ring map[uint64][]interface{}
		// Physical node map, quickly determine if a node exists
		nodes map[string]PlaceholderType
		lock  sync.RWMutex
	}
)

var Placeholder PlaceholderType

// NewConsistentHash returns a ConsistentHash.
func NewConsistentHash() *ConsistentHash {
	return NewCustomConsistentHash(minReplicas, Hash)
}

// NewCustomConsistentHash returns a ConsistentHash with given replicas and hash func.
func NewCustomConsistentHash(replicas int, fn Func) *ConsistentHash {
	if replicas < minReplicas {
		replicas = minReplicas
	}

	if fn == nil {
		fn = Hash
	}

	return &ConsistentHash{
		hashFunc: fn,
		replicas: replicas,
		ring:     make(map[uint64][]any),
		nodes:    make(map[string]PlaceholderType),
	}
}

// adds the node with the number of replicas,
// the later call will overwrite the replicas of the former calls.
func (h *ConsistentHash) AddWithReplicas(node AnyType, replicas int) {
	h.Remove(node)

	// set the limit of replicas
	if replicas > h.replicas {
		replicas = h.replicas
	}

	nodeRepr := Repr(node)
	h.lock.Lock()
	defer h.lock.Unlock()
	h.addNode(nodeRepr)

	for i := 0; i < replicas; i++ {
		// Create virtual node
		hash := h.hashFunc([]byte(nodeRepr + strconv.Itoa(i))) // {nodeRepr}0, {nodeRepr}1
		h.keys = append(h.keys, hash)
		h.ring[hash] = append(h.ring[hash], node)
	}

	sort.Slice(h.keys, func(i, j int) bool {
		return h.keys[i] < h.keys[j]
	})

}

// AddWithWeight adds the node with weight, the weight can be 1 to 100, indicates the percent,
// the later call will overwrite the replicas of the former calls.
func (h *ConsistentHash) AddWithWeight(node any, weight int) {
	// don't need to make sure weight not larger than TopWeight,
	// because AddWithReplicas makes sure replicas cannot be larger than h.replicas
	replicas := h.replicas * weight / TopWeight
	h.AddWithReplicas(node, replicas)
}

// default Add
func (h *ConsistentHash) Add(node any) {
	h.AddWithReplicas(node, h.replicas)
}

func (h *ConsistentHash) Remove(node any) {
	nodeRepr := Repr(node)

	h.lock.Lock()
	defer h.lock.Unlock()

	if !h.containsNode(nodeRepr) {
		return
	}

	for i := 0; i < h.replicas; i++ {
		// hash of virtual node
		hash := h.hashFunc([]byte(nodeRepr + strconv.Itoa(i)))
		index := h.searchVirtualNode(hash)
		if index > -1 {
			h.keys = append(h.keys[:index], h.keys[index+1:]...)
		}
		h.removeRingNode(hash, nodeRepr)
	}

	// remove physical nodes
	h.removeNode(nodeRepr)
}

// Get returns the corresponding node from h base on the given v.
func (h *ConsistentHash) Get(v any) (any, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	if len(h.ring) == 0 {
		return nil, false
	}

	hash := h.hashFunc([]byte(Repr(v)))
	index := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	}) % len(h.keys)

	nodes := h.ring[h.keys[index]]
	switch len(nodes) {
	case 0:
		return nil, false
	case 1:
		return nodes[0], true
	default:
		innerIndex := h.hashFunc([]byte(innerRepr(v)))
		pos := int(innerIndex % uint64(len(nodes)))
		return nodes[pos], true
	}
}
func (h *ConsistentHash) addNode(nodeRepr string) {
	h.nodes[nodeRepr] = Placeholder
}

func (h *ConsistentHash) containsNode(nodeRepr string) bool {
	_, ok := h.nodes[nodeRepr]
	return ok
}

// return index of vitual node
func (h *ConsistentHash) searchVirtualNode(hash uint64) int {
	index := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})
	if index < len(h.keys) && h.keys[index] == hash {
		return index
	}
	return -1
}

func (h *ConsistentHash) removeNode(nodeRepr string) {
	delete(h.nodes, nodeRepr)
}

func (h *ConsistentHash) removeRingNode(hash uint64, nodeRepr string) {
	if nodes, ok := h.ring[hash]; ok {
		newNodes := nodes[:0]
		// should be reltivly small amount of nodes within same hash
		for _, x := range nodes {
			// keep the smae hash but actually different nodeRepr ones
			if Repr(x) != nodeRepr {
				newNodes = append(newNodes, x)
			}
		}
		if len(newNodes) > 0 {
			h.ring[hash] = newNodes
		} else {
			delete(h.ring, hash)
		}
	}
}

// can be interpreted as a serialization method to determine the node string value
// In case of a hash conflict, the key needs to be hashed again
// In order to reduce the probability of conflict, a prime time is appended to reduce the probability of conflict
func innerRepr(node any) string {
	return fmt.Sprintf("%d:%v", prime, node)
}
