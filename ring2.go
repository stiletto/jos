package main

import (
	"errors"
	//"fmt"
    "sync"
    "sort"
)

type Node struct {
    Priority uint8
    Name string
    Id uint32
    Start uint64
}

type Ring struct {
    sync.RWMutex
    nodes []*Node
    copies int
    nodeStep uint64
    copyStep uint64
}

func NewRing(copies int) *Ring {
	var rotator uint64 = 0
	if copies > 1 {
		rotator = uint64(1 << (8*7))/uint64(copies)
	}
	return &Ring{
        nodes: make([]*Node,0,32),
        copies: copies,
        copyStep: rotator,
        nodeStep: 0,
    }
}

var NodeCollision = errors.New("Node Id collision")

func (ring *Ring) AddNodes(nodes []*Node) (errors []error) {
    ring.Lock()
    defer ring.Unlock()
    errors = nil
    for ni, node := range nodes {
		i := sort.Search(len(ring.nodes), func(i int) bool { return ring.nodes[i].Id >= node.Id })
		if i < len(ring.nodes) && ring.nodes[i].Id == node.Id {
			if errors == nil {
				errors = make([]error, len(nodes))
			}
			errors[ni] = NodeCollision
		} else {
			ring.nodes = append(ring.nodes, nil)
			copy(ring.nodes[i+1:], ring.nodes[i:])
			ring.nodes[i] = node
		}
	}
	ring.recalcRanges()
	return
}

// never call this shit with unlocked ring
func (ring *Ring) recalcRanges() {
	var divr uint64 = 0
	nodeCount := len(ring.nodes)
	if nodeCount > 1 {
		divr = uint64(1 << (8*7))/uint64(nodeCount)
	}
	for i, node := range ring.nodes {
		node.Start = divr*uint64(i)
	}
	ring.nodeStep = divr
}

func (ring *Ring) GetNodesForKey(Key uint32) []*Node {
	result := make([]*Node, ring.copies)
	key := uint64(Key) << (8*3)
	for copy := 0; copy<ring.copies; copy++ {
		nodenum := key;
		nodenum = (nodenum / ring.nodeStep) % uint64(len(ring.nodes))
		result[copy] = ring.nodes[nodenum]
		key = (key+ring.copyStep) & ^(uint64(0xFF) << (8*7))
	}
	return result
}

type Range [2]uint32

/*func (ring *Ring) GetNodeRanges(NodeId uint32) []Range {
	for copy := 0; copy<ring.copies; copy++ {
		nodenum := Key
		nodenum = (nodenum / ring.nodeStep) % uint32(len(ring.nodes))
		result[copy] = ring.nodes[nodenum]
		Key += ring.ringCopyStep
	}
}*/

/*func main() {
	ring := NewRing(3)

	node1 := &Node{}
	node1.Priority = 128
	node1.Name = "node1"
	node1.Id = 0
	ring.AddNode(node1)
	
	for i := 2; i<=6; i++ {
		node := &Node{}
		node.Priority = 128
		node.Name = fmt.Sprintf("node%d", i)
		node.Id = uint32(i) * 0xFFFFFF
		ring.AddNode(node)
	}

	node7 := &Node{}
	node7.Priority = 128
	node7.Name = "node7"
	node7.Id = 0xFFFFFFFF
	ring.AddNode(node7)
	
	nodec := make(map[string]int)
	for k:=uint32(0); k<=0xFFFFFF; k++ {
		key := (k << 8);
		nodes := ring.GetIdNodes(key)
		for _, node := range nodes {
			nodec[node.Name] = nodec[node.Name]+1
		}
	}
	fmt.Printf("%#v\n", nodec)
}*/
