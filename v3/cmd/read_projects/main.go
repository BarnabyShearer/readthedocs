// Example of listing projects
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BarnabyShearer/readthedocs/v3"
)

var apiBaseUrlFlag = flag.String("base_url", "https://readthedocs.org/api/v3", "ReadTheDocs API base URL. Can be used to target the Read The Docs For Business API.")

func main() {
	flag.Parse()
	client := readthedocs.NewClientWithURL(os.Getenv("READTHEDOCS_TOKEN"), *apiBaseUrlFlag)
	projects, err := client.GetProjects(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		fmt.Printf("%v\n", project)
	}
}
