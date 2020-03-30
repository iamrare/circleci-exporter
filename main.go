package main

import (
	"encoding/json"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type collector struct{}

func (c collector) Describe(ch chan<- *prometheus.Desc) {
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	d := CallAPI()

	for _, m := range d {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(prometheus.BuildFQName("circleci", "deploys", "per_day"), "deploys per day", []string{"date"}, nil),
			prometheus.CounterValue,
			float64(m.Deploys),
			m.Date,
		)
	}
}

type Response struct {
	NextPageToken string  `json:"next_page_token"`
	Items         []Items `json:"items"`
}

type Items struct {
	Id          string    `json:"id"`
	Status      string    `json:"status"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	StoppedAt   time.Time `json:"stopped_at"`
	CreditsUsed int       `json:"credits_used"`
}

type Count struct {
        Date    string
        Deploys int
}

func CallAPI() []Count {
	url := os.Getenv("URL")
	apiResponse := Response{}

	req, _ := http.NewRequest("GET", url, nil)
	token := "Basic " + os.Getenv("AUTH_TOKEN")
	req.Header.Set("Authorization", token)

	client := new(http.Client)
	resp, _ := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	jsonErr := json.Unmarshal(body, &apiResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	c := []Count{}
	for _, m := range apiResponse.Items {
		i := contains(c, m.CreatedAt.Format("01-02-2006"))
		if i == -1 {
			c = append(c, Count{m.CreatedAt.Format("01-02-2006"), 1})
		} else {
			c[i].Deploys = c[i].Deploys + 1
		}
	}

	return c
}

func contains(s []Count, e string) int {
	for i, a := range s {
		if a.Date == e {
			return i
		}
	}
	return -1
}

var addr = flag.String("listen-address", ":9179", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()

	var c collector
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
