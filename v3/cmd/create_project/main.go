// Example of creating a project
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BarnabyShearer/readthedocs/v3"
)

var nameFlag = flag.String("name", "", "Name of the project.")
var repositoryFlag = flag.String("repository", "", "URL of the repository.")

func main() {
	flag.Parse()
	if *nameFlag == "" || *repositoryFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	client := readthedocs.NewClient(os.Getenv("READTHEDOCS_TOKEN"))
	slug, err := client.CreateProject(context.Background(), readthedocs.CreateUpdateProject{
		CreateProject:         readthedocs.CreateProject{Name: *nameFlag, Repository: readthedocs.Repository{URL: *repositoryFlag, Type: "git"}},
		DefaultVersion:        "latest",
		DefaultBranch:         "main",
		AnalyticsCode:         "",
		AnalyticsDisabled:     false,
		ShowVersionWarning:    true,
		SingleVersion:         false,
		ExternalBuildsEnabled: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	project, err := client.GetProject(context.Background(), slug)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", project)
}
