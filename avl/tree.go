package avl

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	Height int
	Prev   *Node
	Next   *Node
}

type Tree struct {
	root *Node
	// for prev next
	mark Node
}

type Less func(a *Node, b *Node) bool
type Cmp func(n *Node) int

func (n *Node) Reset() {
	*n = Node{}
}

// constructor
func New() *Tree {
	t := &Tree{}
	t.Init()
	return t
}

func (t *Tree) Init() {
	t.root = nil
	t.mark.Prev = &t.mark
	t.mark.Next = &t.mark
}

// iterators
func (t *Tree) Begin() *Node {
	return t.mark.Next
}

func (t *Tree) End() *Node {
	return &t.mark
}

// update
func (t *Tree) Insert(n *Node, less Less) {
	if *n != (Node{}) {
		panic("node not clean")
	}

	if t.root == nil {
		t.root = n
		listInsert(&t.mark, n)
	} else {
		t.root = insert(t.root, n, less)
	}
}

// insert n before ref
func (t *Tree) InsertAt(ref *Node, n *Node) {
	if *n != (Node{}) {
		panic("node not clean")
	}

	if t.root == nil {
		// case 1, insert to empty tree, ref must be end
		if ref != &t.mark {
			panic("bad ref")
		}
		t.root = n
	} else if ref != &t.mark && ref.Left == nil {
		// case 2, insert to ref's left
		t.root = linkUpdated(t.root, ref, n, &ref.Left)
	} else {
		// case 3, insert to ref's left subtree
		// case 4, insert to end
		t.root = linkUpdated(t.root, ref.Prev, n, &ref.Prev.Right)
	}

	// list
	listInsert(ref, n)
}

func (t *Tree) Remove(n *Node) {
	t.root = remove(t.root, n)
	n.Reset()
}

func (t *Tree) Clear() {
	t.root = nil
	t.mark.Reset()
}

// query
func (t *Tree) Find(cmp Cmp) *Node {
	return find(t.root, cmp)
}

// test
func (t *Tree) Verify() {
	end := verify(nil, t.root, t.Begin())
	if end != t.End() {
		panic("list not end")
	}
	if t.Begin().Prev != t.End() {
		panic("list not close")
	}
	if t.mark.Left != nil || t.mark.Right != nil || t.mark.Parent != nil {
		panic("bad mark")
	}
}

func listInsert(at *Node, n *Node) {
	prev := at.Prev
	prev.Next = n
	n.Prev = prev
	n.Next = at
	at.Prev = n
}

func height(n *Node) int {
	if n == nil {
		return -1
	}
	return n.Height
}

func max(a int, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func updateHeight(n *Node) {
	n.Height = 1 + max(height(n.Left), height(n.Right))
}

func setRight(n *Node, r *Node) {
	n.Right = r
	if n.Right != nil {
		n.Right.Parent = n
	}
}

func setLeft(n *Node, l *Node) {
	n.Left = l
	if n.Left != nil {
		n.Left.Parent = n
	}
}

func rotateLeft(n *Node) *Node {
	r := n.Right
	setRight(n, r.Left)
	updateHeight(n)
	setLeft(r, n)
	updateHeight(r)
	return r
}

func rotateRight(n *Node) *Node {
	l := n.Left
	setLeft(n, l.Right)
	updateHeight(n)
	setRight(l, n)
	updateHeight(l)
	return l
}

func leanLeft(n *Node) *Node {
	if height(n.Right) > height(n.Left) {
		n = rotateLeft(n)
	}
	return n
}

func leanRight(n *Node) *Node {
	if height(n.Left) > height(n.Right) {
		n = rotateRight(n)
	}
	return n
}

func fix(n *Node) *Node {
	if height(n.Left) == height(n.Right)+2 {
		setLeft(n, leanLeft(n.Left))
		n = rotateRight(n)
	} else if height(n.Right) == height(n.Left)+2 {
		setRight(n, leanRight(n.Right))
		n = rotateLeft(n)
	} else {
		updateHeight(n)
	}
	return n
}

func insert(root *Node, n *Node, less Less) *Node {
	cur := root
	for {
		if less(n, cur) {
			if cur.Left == nil {
				listInsert(cur, n)
				return linkUpdated(root, cur, n, &cur.Left)
			} else {
				cur = cur.Left
			}
		} else {
			if cur.Right == nil {
				listInsert(cur.Next, n)
				return linkUpdated(root, cur, n, &cur.Right)
			} else {
				cur = cur.Right
			}
		}
	}
}

func removeMin(root *Node, n *Node) *Node {
	p := n.Parent
	if p == nil {
		return n.Right
	}
	return linkUpdated(root, p, n.Right, &p.Left)
}

func remove(root *Node, n *Node) *Node {
	// for later use
	prev := n.Prev
	next := n.Next

	// detach from list
	prev.Next = n.Next
	next.Prev = n.Prev

	p := n.Parent
	var updated *Node
	if n.Left == nil {
		// replace n with right subtree
		updated = n.Right
	} else if n.Right == nil {
		// replace n with left subtree
		updated = n.Left
	} else {
		// borrow next from right subtree
		r := n.Right
		r.Parent = nil
		r = removeMin(r, next)
		// set up next
		setRight(next, r)
		setLeft(next, n.Left)
		// replace n with next
		updated = fix(next)
	}

	if p == nil {
		if updated != nil {
			updated.Parent = nil
		}
		return updated
	} else if n == p.Left {
		return linkUpdated(root, p, updated, &p.Left)
	} else {
		return linkUpdated(root, p, updated, &p.Right)
	}
}

func linkUpdated(root *Node, p *Node, c *Node, link **Node) *Node {
	for {
		// p is not nil, c may be nil
		newp := p.Parent
		h := height(p)

		// link c to p and adjust p to newc
		*link = c
		if c != nil {
			c.Parent = p
		}
		newc := fix(p)

		if newc == p && height(newc) == h {
			// p was not rotated and p's height not updated
			return root
		}
		if newp == nil {
			// newc is new root
			newc.Parent = nil
			return newc
		}

		// next linkd
		var newlink **Node
		if p == newp.Left {
			newlink = &newp.Left
		} else {
			newlink = &newp.Right
		}

		p, c, link = newp, newc, newlink
	}
}

func find(n *Node, cmp Cmp) *Node {
	for n != nil {
		r := cmp(n)
		if r < 0 {
			n = n.Left
		} else if r > 0 {
			n = n.Right
		} else {
			return n
		}
	}
	return nil
}

func verify(parent *Node, n *Node, cur *Node) *Node {
	if n == nil {
		return cur
	}

	if n.Parent != parent {
		panic("parent mismatch")
	}

	if height(n) != 1+max(height(n.Left), height(n.Right)) {
		panic("bad height")
	}

	cur = verify(n, n.Left, cur)
	if cur != n {
		panic("tree mismatch with list")
	}
	if cur.Next.Prev != cur {
		panic("bad list")
	}
	cur = cur.Next
	cur = verify(n, n.Right, cur)

	return cur
}
