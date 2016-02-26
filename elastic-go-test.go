// Copyright 2013 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"encoding/json"

	"github.com/mattbaird/elastigo/lib"
)


type MyUser struct {
	Name string
	Age  int
}

type OrderType struct {
	Doc DocType
}

type DocType struct {
//	Date string
	Code string
}


var (
	host *string = flag.String("host", "pre.elasticsearch1.oneboxtickets.com", "Elasticsearch Host")
)

func main() {
	c := elastigo.NewConn()
	log.SetFlags(log.LstdFlags)
	flag.Parse()

	// Trace all requests
/*	c.RequestTracer = func(method, url, body string) {
		log.Printf("Requesting %s %s", method, url)
		log.Printf("Request body: %s", body)
	}
*/
	fmt.Println("host = ", *host)
	// Set the Elasticsearch Host to Connect to
	c.Domain = *host

	// Index a document
//	_, err := c.Index("testindex", "user", "docid_1", nil, `{"name":"bob"}`)
//	exitIfErr(err)

	// Index a doc using a map of values
//	_, err = c.Index("testindex", "user", "docid_2", nil, map[string]string{"name": "venkatesh"})
//	exitIfErr(err)

	// Index a doc using Structs
//	_, err = c.Index("testindex", "user", "docid_3", nil, MyUser{"wanda", 22})
//	exitIfErr(err)

	// Search Using Raw json String
	searchJson := `{
	    "query": {
		  "nested": {
			"path": "doc.products",
			"query": {
			  "filtered": {
				"filter": {
				  "bool": {
				   "must": [
					  {
						"term": { "doc.products.eventId": 2627}
					  }
					]
				  }
				}
			  }
			}
		  }
	  }
	}`

	// Elasticsearch Search
	out, err := c.Search("onebox-order/couchbaseDocument", "", map[string]interface{} {"size" : 10}, searchJson)
	if len(out.Hits.Hits) > 0	{

		var order OrderType
		var orders []*OrderType

		for i := 0; i < out.Hits.Total-1; i++ {

			err := json.Unmarshal(*out.Hits.Hits[i].Source, &order)
			if err != nil {
				panic(err)
			}
			o := new(OrderType)
			orders = append(orders, o)
			fmt.Printf("order: %v\n", order)
		}
	}
	exitIfErr(err)

}
func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}


