package monitor

import (
	"strings"
	"time"
)

//Query defines an executable query with name
type Query struct {
	Name   string
	filter string
	Exec   func(string) (*Hits, error)
}

//BuildBody builds a complete query
func (query *Query) BuildBody(size string, since time.Time) string {
	body := strings.Replace(queryTemplate, "<size>", size, 1)
	body = strings.Replace(body, "<timestamp>", formatTime(since), 1)
	return strings.Replace(body, "<query>", query.filter, 1)
}

func formatTime(timestamp time.Time) string {
	return timestamp.Format("2006-01-02 15:04:05")
}

const queryTemplate = `{
  "query": {
    "bool": {
      "must": {
        "query_string": {
          "query": "<query>"
        }
      },
      "filter": {
        "range": {
          "@timestamp": {
            "gte": "<timestamp>",
            "format": "yyyy-MM-dd HH:mm:ss"
          }
        }
      }
    }
  },
  "size":<size>
}`
