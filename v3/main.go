// API client for readthedocs.org.
package readthedocs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// APIs default base URL.
const BaseURLV3 = "https://readthedocs.org/api/v3"

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

// Create the API client using the default URL, providing the authentication key.
func NewClient(apiKey string) *Client {
	return NewClientWithURL(apiKey, BaseURLV3)
}

// Create the API client using a custom URL, providing the authentication key.
func NewClientWithURL(apiKey string, baseUrl string) *Client {
	return &Client{
		BaseURL: baseUrl,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type projects struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"` // TODO: atm we just do the first 1000, we should follow these
	Previous string    `json:"previous"`
	Results  []Project `json:"results"`
}

type Project struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Language struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"language"`
	ProgrammingLanguage struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"programming_language"`
	Repository struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"repository"`
	DefaultVersion string `json:"default_version"`
	DefaultBranch  string `json:"default_branch"`
	SubprojectOf   string `json:"subproject_of"`
	TranslationOf  string `json:"translation_of"`
	URLs           struct {
		Documentation string `json:"documentation"`
		Home          string `json:"home"`
	} `json:"urls"`
	Tags  []string `json:"tags"`
	Users []struct {
		Username string `json:"username"`
	}
	ActiveVersions map[string]string `json:"active_versions"`
}

type Repository struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type CreateProject struct {
	Name                string     `json:"name"`
	Repository          Repository `json:"repository"`
	Homepage            string     `json:"homepage,omitempty"`
	ProgrammingLanguage string     `json:"programming_language,omitempty"`
	Language            string     `json:"language,omitempty"`
	Organization        string     `json:"organization,omitempty"`
	Teams               string     `json:"teams,omitempty"`
}

type CreateUpdateProject struct {
	CreateProject
	DefaultVersion        string `json:"default_version,omitempty"`
	DefaultBranch         string `json:"default_branch,omitempty"`
	AnalyticsCode         string `json:"analytics_code,omitempty"`
	AnalyticsDisabled     bool   `json:"analytics_disabled"`
	ShowVersionWarning    bool   `json:"show_version_warning"`
	SingleVersion         bool   `json:"single_version"`
	ExternalBuildsEnabled bool   `json:"external_builds_enabled"`
}

func (c *Client) GetProjects(ctx context.Context) ([]Project, error) {
	projects := projects{}
	err := c.sendRequest(ctx, "GET", "/projects/?limit=1000", nil, &projects)
	return projects.Results, err
}

func (c *Client) GetProject(ctx context.Context, slug string) (Project, error) {
	project := Project{}
	err := c.sendRequest(ctx, "GET", fmt.Sprintf("/projects/%s/", slug), nil, &project)
	return project, err
}

func (c *Client) DeleteProject(ctx context.Context, slug string) error {
	return c.sendRequest(ctx, "DELETE", fmt.Sprintf("/projects/%s/", slug), nil, nil)
}

func (c *Client) CreateProject(ctx context.Context, createProject CreateUpdateProject) (string, error) {
	project := Project{}
	createProjectJson, err := json.Marshal(createProject.CreateProject)
	if err != nil {
		return "", err
	}
	err = c.sendRequest(ctx, "POST", "/projects/", createProjectJson, &project)
	if err != nil {
		return "", err
	}
	// API requires a create then patch to set all values
	return project.Slug, c.UpdateProject(ctx, project.Slug, createProject)
}

func (c *Client) UpdateProject(ctx context.Context, slug string, updateProject CreateUpdateProject) error {
	UpdateProjectJSON, err := json.Marshal(updateProject)
	if err != nil {
		return err
	}
	return c.sendRequest(ctx, "PATCH", fmt.Sprintf("/projects/%s/", slug), UpdateProjectJSON, nil)
}

func (c *Client) sendRequest(ctx context.Context, method string, url string, body []byte, result interface{}) error {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.apiKey))

	req = req.WithContext(ctx)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(bodyBytes))
	}

	if result != nil {
		if err = json.NewDecoder(res.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}
