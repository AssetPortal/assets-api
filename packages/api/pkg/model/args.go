package model

type Order struct {
	Order     *string `in:"query=order"`
	Ascending *bool   `in:"query=ascending"`
}

type Pagination struct {
	Limit  *int `in:"query=limit"`
	Offset *int `in:"query=offset"`
}
