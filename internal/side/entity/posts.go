package entity

type Posts []*Post

func (p Posts) Slice(i, j int) Posts {
	if p.At(i) == nil {
		return p
	}
	if p.At(j) == nil {
		return p[i:]
	}
	return p[i:j]
}

func (p Posts) At(i int) *Post {
	if p.LastIndex() < 0 {
		return nil
	}
	if p.LastIndex() < i {
		return nil
	}
	return p[i]
}

func (p Posts) First() *Post {
	return p.At(0)
}

func (p Posts) Second() *Post {
	return p.At(1)
}

func (p Posts) Last() *Post {
	return p.At(p.LastIndex())
}

func (p Posts) Penultimate() *Post {
	return p.At(p.LastIndex() - 1)
}

func (p Posts) LastIndex() int {
	return len(p) - 1
}

var EmptyPosts = make(Posts, 0)
