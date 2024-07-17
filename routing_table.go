package nnkademlia

import (
	"fmt"
	"math/big"
	"sort"
)

type Bucket struct {
	nodes []*Node
	max   int
}

func NewBucket(max int) *Bucket {
	return &Bucket{
		nodes: make([]*Node, 0, max),
		max:   max,
	}
}

func (b *Bucket) String() string {
	bucket := "{"
	for _, n := range b.nodes {
		bucket += fmt.Sprintf("%s, ", n.ID.Text(16))
	}
	bucket = bucket[:len(bucket) - 2]
	bucket += "}"
	return bucket
}

func (b *Bucket) Add(node *Node) {
	if len(b.nodes) < b.max {
		b.nodes = append([]*Node{node}, b.nodes...)
	} else {
		for i := len(b.nodes) - 1; i > 1; i-- {
			b.nodes[i] = b.nodes[i-1]
		}
		b.nodes[0] = node
	}
}

func (b *Bucket) GetById(id *big.Int) *Node {
	for i := 0; i < len(b.nodes); i++ {
		if b.nodes[i].ID.Cmp(id) == 0 {
			return b.nodes[i]
		}
	}
	return nil
}

func (b *Bucket) GetAll() []*Node {
	nodes := make([]*Node, 0)

	for _, n := range b.nodes {
		if n != nil {
			nodes = append(nodes, n)
		}
	}

	return nodes
}

type RoutingTable struct {
	owner          *Node
	Buckets        []*Bucket
	neighborsCount int
	maxBucketSize  int
}

func NewRoutingTable(maxBucketSize int, bucketsCount int, owner *Node) *RoutingTable {
	buckets := make([]*Bucket, bucketsCount)

	return &RoutingTable{
		owner:          owner,
		Buckets:        buckets,
		neighborsCount: 0,
		maxBucketSize:  maxBucketSize,
	}
}

func (rt *RoutingTable) String() string {
	buckets := ""
	for _, b := range rt.Buckets {
		if b == nil {
			buckets += "nil, "
			continue
		}
		buckets += fmt.Sprintf("%s, ", b.String())
	}
	return buckets[:len(buckets)-2]
}

func (rt *RoutingTable) Add(node *Node) {
	if rt.owner.Cmp(node) {
		return
	}

	index := rt.owner.GetBucketFromNode(node)
	if rt.Buckets[index] == nil {
		rt.Buckets[index] = NewBucket(rt.maxBucketSize)
	}

	rt.neighborsCount++
	rt.Buckets[index].Add(node)
}

func (rt *RoutingTable) GetById(id *big.Int) *Node {
	index := rt.owner.GetBucketFromId(id)
	if rt.Buckets[index] == nil {
		return nil
	}

	return rt.Buckets[index].GetById(id)
}

func (rt *RoutingTable) FindNearest(node *Node) []*Node {
	nearest := make([]*Node, 0, rt.maxBucketSize)

	index := rt.owner.GetBucketFromNode(node)
	if index < 0 {
		index = 0
	}

	if rt.Buckets[index] != nil {
		nearest = append(nearest, rt.Buckets[index].GetAll()...)
	}

	if len(nearest) < rt.maxBucketSize {
		for current := index - 1; current >= 0; current-- {
			bucket := rt.Buckets[current]
			if bucket != nil {
				nearest = append(nearest, bucket.GetAll()...)
			}
			if len(nearest) >= rt.maxBucketSize {
				break
			}
		}
	}

	if len(nearest) < rt.maxBucketSize {
		for current := index + 1; current < len(rt.Buckets); current++ {
			bucket := rt.Buckets[current]
			if bucket != nil {
				nearest = append(nearest, bucket.GetAll()...)
			}
			if len(nearest) >= rt.maxBucketSize {
				break
			}
		}
	}

	sort.Slice(nearest, func(i, j int) bool {
		toA := DistanceBetween(node.ID, nearest[i].ID)
		toB := DistanceBetween(node.ID, nearest[j].ID)
		return toA.Cmp(toB) == -1
	})

	breakIndex := rt.maxBucketSize
	if len(nearest) < breakIndex {
		breakIndex = len(nearest)
	}
	nearest = nearest[:breakIndex]

	return nearest
}

func (rt *RoutingTable) FindAll() []*Node {
	nodes := make([]*Node, 0)
	for _, b := range rt.Buckets {
		if b == nil {
			continue
		}
		nodes = append(nodes, b.GetAll()...)
	}
	return nodes
}
