package trie

type Trie interface {
	Add([]byte)
	Find([]byte) bool
}

type Radix struct {
	root *node
}

func NewRadix() *Radix {
	return &Radix{}
}

func (r *Radix) Add(value []byte) {
	if r.root == nil {
		r.root = newNode()
	}
	r.root.merge(value, 0)
}

func (r *Radix) Find(value []byte) bool {
	if r.root == nil {
		return len(value) == 0
	}

	return r.root.find(value, 0)
}

type node struct {
	children [2]*node
}

func newNode() *node {
	return &node{}
}

func (n *node) merge(value []byte, pos int) {
	if len(value) == 0 {
		return
	}

	n.current(value, pos, true).merge(next(value, pos))
}

func (n *node) find(value []byte, pos int) bool {
	if len(value) == 0 {
		return true
	}

	current := n.current(value, pos, false)
	if current == nil {
		return false
	}

	return current.find(next(value, pos))
}

func next(value []byte, pos int) ([]byte, int) {
	if pos == 7 {
		return value[1:], 0
	}

	return value, pos + 1
}

func (n *node) current(value []byte, pos int, ensure bool) *node {
	current := value[0] & (1 << (7 - pos))
	if current > 0 {
		current = 1
	}

	if ensure && n.children[current] == nil {
		n.children[current] = newNode()
	}

	return n.children[current]
}
