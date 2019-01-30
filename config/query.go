package config

//Query all the options for a single query
type Query map[string]interface{}

func (query *Query) Name() string {
	return (*query)["Name"].(string)
}

func (query *Query) QueryText() string {
	return (*query)["Query"].(string)
}

func CreateQuery(name string, query string) Query {
	q := make(map[string]interface{})
	q["Name"] = name
	q["Query"] = query
	return q
}

func (query *Query) isValid() bool {
	_, ok1 := (*query)["Name"].(string)
	_, ok2 := (*query)["Query"].(string)
	return ok1 && ok2
}
