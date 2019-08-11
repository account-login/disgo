package avl

import (
	"math/rand"
	"testing"
)

func setupSequencialNodes(N int) (nodes []Data, indexes []int) {
	nodes = make([]Data, N)
	indexes = make([]int, N)
	for i := 0; i < N; i++ {
		nodes[i].val = int64(i)
		indexes[i] = i
	}
	rand.Shuffle(N, func(i, j int) {
		nodes[i].val, nodes[j].val = nodes[j].val, nodes[i].val
		indexes[int(nodes[i].val)] = i
		indexes[int(nodes[j].val)] = j
	})
	return
}

func BenchmarkInsertSequencial(b *testing.B) {
	tree := New()
	nodes, indexes := setupSequencialNodes(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(&nodes[indexes[i]].node, less)
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
	nodes, indexes := setupSequencialNodes(b.N)
	for i := 0; i < b.N; i++ {
		tree.Insert(&nodes[i].node, less)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Remove(&nodes[indexes[i]].node)
	}
}

func BenchmarkRemoveSequencialReverse(b *testing.B) {
	tree := New()
	nodes, indexes := setupSequencialNodes(b.N)
	for i := 0; i < b.N; i++ {
		tree.Insert(&nodes[i].node, less)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Remove(&nodes[b.N-indexes[i]-1].node)
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

func cmpHelper(i int64, n *Node) int {
	val := DataOf(n).val
	if i == val {
		return 0
	} else if i < val {
		return -1
	} else {
		return 1
	}
}

func BenchmarkFindSequencial(b *testing.B) {
	tree := New()
	nodes := make([]Data, b.N)
	for i := 0; i < b.N; i++ {
		nodes[i].val = int64(i)
		tree.Insert(&nodes[i].node, less)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := tree.Find(func(n *Node) int {
			return cmpHelper(int64(i), n)
		})
		_ = n
	}
}

func BenchmarkFindRandom(b *testing.B) {
	tree := New()
	nodes := make([]Data, b.N)
	vals := make([]int64, b.N)
	for i := 0; i < b.N; i++ {
		val := rand.Int63()
		nodes[i].val = val
		vals[i] = val
		tree.Insert(&nodes[i].node, less)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := tree.Find(func(n *Node) int {
			return cmpHelper(vals[i], n)
		})
		_ = n
	}
}
