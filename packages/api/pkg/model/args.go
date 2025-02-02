package model

const MAX_LIMIT = 100

type Order struct {
	Order     *string `in:"query=order"`
	Ascending *bool   `in:"query=ascending"`
}

type Pagination struct {
	Limit  *int `in:"query=limit"`
	Offset *int `in:"query=offset"`
}

func (p *Pagination) Validate() {
	if p.Limit == nil || *p.Limit > MAX_LIMIT {
		defaultLimit := MAX_LIMIT
		p.Limit = &defaultLimit
	}
	if p.Offset == nil {
		defaultOffset := 0
		p.Offset = &defaultOffset
	}
}
