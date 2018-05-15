package monitor

import "strings"

//Query defines an executable query with name
type Query struct {
	Name   string
	filter string
	Exec   func(string) (*Hits, error)
}

//BuildBody builds a complete query
func (query *Query) BuildBody(size string) string {
	body := strings.Replace(queryTemplate, "<size>", size, 1)
	return strings.Replace(body, "<query>", query.filter, 1)
}

const queryTemplate = `{
    "query": {
        "query_string": {
            "query": "<query>"
        }
    },
    "size":<size>
}`
