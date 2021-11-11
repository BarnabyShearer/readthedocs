// Example of listing projects
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BarnabyShearer/readthedocs/v3"
)

func main() {
	client := readthedocs.NewClient(os.Getenv("READTHEDOCS_TOKEN"))
	projects, err := client.GetProjects(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		fmt.Printf("%v\n", project)
	}
}
