package nnkademlia

import (
	"fmt"
	"math/big"
	"time"
)

const (
	MAX_BUCKET_SIZE  = 4
	MAX_BUCKET_COUNT = 160
)

type Node struct {
	ID           *big.Int
	routingTable *RoutingTable
	hashStore    *Store
	nodeInFind   map[*big.Int]bool
	valueInFind  map[string]bool
}

func NewNode(name string) *Node {
	node := &Node{
		ID:          IdFromString(name),
		hashStore:   NewStore(),
		nodeInFind:  make(map[*big.Int]bool),
		valueInFind: make(map[string]bool),
	}

	node.routingTable = NewRoutingTable(MAX_BUCKET_SIZE, MAX_BUCKET_COUNT, node)

	return node
}

func (n *Node) String() string {
	return fmt.Sprintf("id: %s, rt: {%s}, hs: {%s}", n.ID.Text(16), n.routingTable.String(), n.hashStore.String())
}

func (n *Node) Cmp(other *Node) bool {
	return n.ID.Cmp(other.ID) == 0
}

func (n *Node) GetBucketFromNode(other *Node) int {
	distance := DistanceBetween(n.ID, other.ID)
	bitString := fmt.Sprintf("%b", distance)

	return len(bitString) - 1
}

func (n *Node) GetBucketFromId(other *big.Int) int {
	distance := DistanceBetween(n.ID, other)
	bitString := fmt.Sprintf("%b", distance)

	return len(bitString) - 1
}

func (n *Node) ping(node *Node) []*Node {
	n.routingTable.Add(node)
	return n.routingTable.FindNearest(node)
}

func (n *Node) Ping(node *Node) {
	nodes := node.ping(n)
	n.routingTable.Add(node)
	for _, nn := range nodes {
		n.routingTable.Add(nn)
	}
}

func (n *Node) store(key string, value string, origin *big.Int) bool {
	return n.hashStore.Add(key, value, origin)
}

func (n *Node) Store(key string, value string, node *Node) bool {
	return node.store(key, value, n.ID)
}

func (n *Node) findValue(key string) (string, bool) {
	if isInFind := n.valueInFind[key]; isInFind {
		return "", false
	}
	n.valueInFind[key] = true
	defer func() {
		go func() {
			time.Sleep(1 * time.Second)
			n.valueInFind[key] = false	
		}()
	}()

	val, ok := n.hashStore.Find(key)
	if ok {
		return val, true
	}

	nearest := n.routingTable.FindNearest(n)

	for _, nn := range nearest {
		val, ok = nn.findValue(key)
		if ok {
			return val, true
		}
	}
	return "", false
}

func (n *Node) FindValue(key string) (string, bool) {
	return n.findValue(key)
}

func (n *Node) findNode(id *big.Int) ([]*Node, bool) {
	if isInFind, ok := n.nodeInFind[id]; ok && isInFind {
		return nil, false
	}
	n.nodeInFind[id] = true
	defer func() {
		n.nodeInFind[id] = false
	}()
	nearest := n.routingTable.FindNearest(&Node{ID: id})
	for _, nn := range nearest {
		if nn.ID.Cmp(id) == 0 {
			return nearest, true
		}
	}

	for _, nn := range nearest {
		rnearest, ok := nn.findNode(id)
		if ok {
			return rnearest, true
		}
	}
	return nil, false
}

func (n *Node) FindNode(id *big.Int) ([]*Node, bool) {
	return n.findNode(id)
}
