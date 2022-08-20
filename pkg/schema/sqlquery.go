package schema

type SqlQuery struct {
	Query  string `json:"query"`
	Params []any  `json:"params"`
}
