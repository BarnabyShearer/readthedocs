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
var organizationFlag = flag.String("organization", "", "(OPTIONAL) ReadTheDocs for Business organization where the project should be created.")
var apiBaseUrlFlag = flag.String("base_url", "https://readthedocs.org/api/v3", "ReadTheDocs API base URL. Can be used to target the Read The Docs For Business API.")

func main() {
	flag.Parse()
	if *nameFlag == "" || *repositoryFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	client := readthedocs.NewClientWithURL(os.Getenv("READTHEDOCS_TOKEN"), *apiBaseUrlFlag)
	slug, err := client.CreateProject(context.Background(), readthedocs.CreateUpdateProject{
		CreateProject:         readthedocs.CreateProject{Name: *nameFlag, Repository: readthedocs.Repository{URL: *repositoryFlag, Type: "git"}, Organization: *organizationFlag},
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
