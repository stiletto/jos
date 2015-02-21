// +build ignore
package main

/*import (
	"errors"
	"fmt"
    "sync"
    "sort"
)

type Node struct {
    Priority uint8
    Name string
    Id uint32
}

type Ring struct {
    sync.RWMutex
    nodes []*Node
    copies int
}

func NewRing(copies int) *Ring {
    return &Ring{
        nodes: make([]*Node,0,32),
        copies: copies,
    }
}

var NodeCollision = errors.New("Node Id collision")

func (ring *Ring) AddNode(node *Node) (error) {
    ring.Lock()
    defer ring.Unlock()
    
    
    i := sort.Search(len(ring.nodes), func(i int) bool { return ring.nodes[i].Id >= node.Id })
    if i < len(ring.nodes) && ring.nodes[i].Id == node.Id {
        return NodeCollision
    } else {
        ring.nodes = append(ring.nodes, nil)
        copy(ring.nodes[i+1:], ring.nodes[i:])
        ring.nodes[i] = node
	}
	return nil
}

func (ring *Ring) GetIdNodes(Id uint32) []*Node {
	result := make([]*Node, ring.copies)
	var divr uint32
	nodeCount := len(ring.nodes)
	if nodeCount > 2 {
		divr = uint32(uint64(0x100000000)/uint64(len(ring.nodes)))
	} else {
		divr = 0xFFFFFFFF
	}
	
	for copy := 0; copy<ring.copies; copy++ {
		nodenum := Id + uint32((uint64(copy) << 32) / uint64(ring.copies))
		nodenum = (nodenum / divr) % uint32(nodeCount)
		result[copy] = ring.nodes[nodenum]
	}
	return result
}

type Range [2]uint32

func (ring *Ring) GetNodeRanges(NodeId uint32) []Range {
}

func main() {
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
