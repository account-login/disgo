package avl

import (
	"math/rand"
	"testing"
)

func BenchmarkInsertSequencial(b *testing.B) {
	tree := New()
	nodes := make([]Data, b.N)
	for i := 0; i < b.N; i++ {
		nodes[i].val = int64(i)
		tree.Insert(&nodes[i].node, less)
	}
}

func BenchmarkInsertRandom(b *testing.B) {
	tree := New()
	nodes := make([]Data, b.N)
	for i := 0; i < b.N; i++ {
		nodes[i].val = rand.Int63()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(&nodes[i].node, less)
	}
}

func BenchmarkRemoveSequencial(b *testing.B) {
	tree := New()
	nodes := make([]Data, b.N)
	for i := 0; i < b.N; i++ {
		nodes[i].val = int64(i)
		tree.Insert(&nodes[i].node, less)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Remove(&nodes[i].node)
	}
}

func BenchmarkRemoveRandom(b *testing.B) {
	tree := New()
	nodes := make([]Data, b.N)
	for i := 0; i < b.N; i++ {
		nodes[i].val = rand.Int63()
		tree.Insert(&nodes[i].node, less)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Remove(&nodes[i].node)
	}
}

func BenchmarkFind(b *testing.B) {
	tree := New()
	nodes := make([]Data, b.N)
	for i := 0; i < b.N; i++ {
		nodes[i].val = int64(i)
		tree.Insert(&nodes[i].node, less)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := tree.Find(func(n *Node) int {
			val := DataOf(n).val
			if int64(i) == val {
				return 0
			} else if int64(i) < val {
				return -1
			} else {
				return 1
			}
		})
		if n != &nodes[i].node {
			panic("find")
		}
	}
}
