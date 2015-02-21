package main

import (
	"math"
	"fmt"
	"testing"
)

func TestRing(t *testing.T) {
	ring := NewRing(3)

	node1 := &Node{}
	node1.Priority = 128
	node1.Name = "node1"
	node1.Id = 0
	ring.AddNodes([]*Node{node1})
	
	nc := 300000
	nnodes := make([]*Node,nc)
	for i := 0; i<nc; i++ {
		node := &Node{}
		node.Priority = 128
		node.Name = fmt.Sprintf("node%d", i+2)
		node.Id = uint32(i) * 0xFF
		nnodes[i] = node
	}
	ring.AddNodes(nnodes)

	node7 := &Node{}
	node7.Priority = 128
	node7.Name = "node99999999"
	node7.Id = 0xFFFFFFFF
	ring.AddNodes([]*Node{node7})
	
	nodec := make(map[string]int)
	for k:=uint32(0); k<=0xFFFFFF; k++ {
		key := (k << 8) + k&0xff;
		nodes := ring.GetNodesForKey(key)
		for _, node := range nodes {
			nodec[node.Name] = nodec[node.Name]+1
		}
	}
	sum := 0
	max := 0
	maxi := ""
	min := 0xFFFFFF
	mini := ""
	for k,val := range nodec {
		sum += val
		if val > max {
			max = val
			maxi = k
		}
		if val < min {
			min = val
			mini = k
		}
	}
	avg := float64(sum)/float64(len(nodec))
	fmt.Printf("Min: %d (%s)\n", min, mini)
	fmt.Printf("Average: %f\n", avg)
	fmt.Printf("Max: %d (%s)\n", max, maxi)
	mdevsum := 0.0
	for _,val := range nodec {
		mdevsum += math.Pow(float64(val) - avg,2)
	}
	mdevsum = math.Sqrt(mdevsum / float64(len(nodec)))
	fmt.Printf("Mdev: %f\n", mdevsum)
	if max-min > 10 {
		t.Errorf("Key distribution is too unbalanced. max-min == %d", max-min)
	}

}
