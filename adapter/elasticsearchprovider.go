package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MaibornWolff/elcep/monitor"
)

//ElasticSearchProvider for Queries
type ElasticSearchProvider struct {
	URL string
}

//ElasticResponse defines as response from a query
type ElasticResponse struct {
	Took     int  `Json:"took"`
	TimedOut bool `Json:"timed_out"`
	Shards   struct {
		Total      int `Json:"total"`
		Successful int `Json:"successful"`
		Skipped    int `Json:"skipped"`
		Failed     int `Json:"failed"`
	} `Json:"_shards"`
	Hits struct {
		Total    float64       `Json:"total"`
		MaxScore float64       `Json:"max_score"`
		Hits     []interface{} `Json:"hits"`
	} `Json:"hits"`
}

//ExecRequest on elasticsearch host
func (provider *ElasticSearchProvider) ExecRequest(querypath string, request string) (*monitor.Hits, error) {
	url := fmt.Sprintf(provider.URL + querypath)
	elResponse := &ElasticResponse{}

	req, err := http.NewRequest("GET", url, bytes.NewBufferString(request))
	req.Header.Set("Content-Type", "application/Json")

	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(elResponse); err != nil {
		log.Println(err)
	}

	return &monitor.Hits{
		Total:   elResponse.Hits.Total,
		Results: elResponse.Hits.Hits,
	}, err
}
