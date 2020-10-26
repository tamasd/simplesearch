package disjoint_set

type StringItem struct {
	parent *StringItem
	data   string
	size   uint
}

func NewStringItem(data string) *StringItem {
	i := &StringItem{
		data: data,
		size: 1,
	}
	i.parent = i

	return i
}

func (i *StringItem) String() string {
	return i.data
}

func (i *StringItem) Find() *StringItem {
	item := i

	for item.parent != item {
		item, item.parent = item.parent, item.parent.parent
	}

	return item
}

func StringUnion(x, y *StringItem) {
	xRoot := x.Find()
	yRoot := y.Find()

	if xRoot == yRoot {
		return
	}

	if xRoot.size < yRoot.size {
		xRoot, yRoot = yRoot, xRoot
	}

	yRoot.parent = xRoot
	xRoot.size = xRoot.size + yRoot.size
}
