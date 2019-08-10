package avl

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	val  int64
	node Node
	tail string
}

func DataOf(n *Node) *Data {
	return (*Data)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) - unsafe.Offsetof(Data{}.node)))
}

func verifyMonotonic(tree *Tree) {
	for node := tree.Begin().Next; node != tree.End(); node = node.Next {
		if less(node, node.Prev) {
			panic("not monotonic")
		}
	}
}

func TestBasic(t *testing.T) {
	tree := New()
	tree.Verify()
	assert.Equal(t, tree.Begin(), tree.End())

	{
		d1 := &Data{val: 1, tail: "asdf"}
		tree.InsertAt(tree.End(), &d1.node)
		tree.Verify()
	}
	assert.Equal(t, int64(1), DataOf(tree.Begin()).val)
	assert.Equal(t, tree.End(), tree.Begin().Next)

	tree.Remove(tree.Begin())
	assert.Equal(t, tree.Begin(), tree.End())
}

func less(a *Node, b *Node) bool {
	return DataOf(a).val < DataOf(b).val
}

var ttt = "\t\t\t\t\t\t\t\t\t\t\t\t\t"

func dump(n *Node, level int) {
	s := ttt[:level]
	if n == nil {
		log.Print(s + "NIL")
	} else {
		d := DataOf(n)
		s += fmt.Sprintf("%p %v %v", &d.node, d.val, d.node.Height)
		log.Print(s)
		dump(n.Left, level+1)
		dump(n.Right, level+1)
	}
}

func TestTmp(t *testing.T) {
	ints := []int64{63, 35, 27, 33, 28}
	hist := map[int64]*Data{}
	tree := New()
	for _, val := range ints {
		if val > 0 {
			d := &Data{val: val}
			hist[val] = d
			t.Logf("insert %v", d.val)
			tree.Insert(&d.node, less)
		} else {
			d := hist[-val]
			t.Logf("remove %v", d.val)
			tree.Remove(&d.node)
		}

		dump(tree.root, 0)
		tree.Verify()
	}
}

func TestRandomInsertRemove(t *testing.T) {
	threshold := 2000
	p1 := 0.5
	p2 := 0.45
	N := 1 * 1000 * 1000

	tree := New()
	var nodes []*Data
	for i := 0; i < N; i++ {
		r := rand.Float64()
		p := p1
		if len(nodes) > threshold {
			p = p2
		}

		if r < p || len(nodes) == 0 {
			// insert a random node
			d := &Data{val: rand.Int63()}
			nodes = append(nodes, d)
			//t.Logf("insert %v %p", d.val, &d.node)
			tree.Insert(&d.node, less)
		} else {
			// remove a random node
			idx := rand.Intn(len(nodes))
			d := nodes[idx]
			//t.Logf("remove %v %p", d.val, &d.node)
			tree.Remove(&d.node)
			nodes[idx] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		}

		//dump(tree.root, 0)
		tree.Verify()
		verifyMonotonic(tree)
		if i%1000 == 0 {
			t.Log(i, len(nodes))
		}
	}
}

func TestRandomInsertAtRemove(t *testing.T) {
	threshold := 2000
	p1 := 0.5
	p2 := 0.45
	N := 1 * 1000 * 1000

	tree := New()
	var nodes []*Node
	for i := 0; i < N; i++ {
		r := rand.Float64()
		p := p1
		if len(nodes) > threshold {
			p = p2
		}

		if r < p || len(nodes) == 0 {
			// insert at random node
			nodes = append(nodes, tree.End())
			at := nodes[rand.Intn(len(nodes))]
			val := int64(0)
			if at == tree.End() {
				val = math.MaxInt32
			} else if at == tree.Begin() {
				val = math.MinInt32
			} else {
				val = (DataOf(at.Prev).val + DataOf(at).val) / 2
			}
			d := &Data{val: val}
			nodes[len(nodes)-1] = &d.node
			//t.Logf("insert %v %p at %p", d.val, &d.node, at)
			tree.InsertAt(at, &d.node)
		} else {
			// remove a random node
			idx := rand.Intn(len(nodes))
			d := DataOf(nodes[idx])
			//t.Logf("remove %v %p", d.val, &d.node)
			tree.Remove(&d.node)
			nodes[idx] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		}

		//dump(tree.root, 0)
		tree.Verify()
		verifyMonotonic(tree)
		if i%1000 == 0 {
			t.Log(i, len(nodes))
		}
	}
}
