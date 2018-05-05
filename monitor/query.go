package monitor

//Query defines an executable query with name
type Query struct {
	Name string
	Body string
	Exec func(string, string) (*Hits, error)
}
