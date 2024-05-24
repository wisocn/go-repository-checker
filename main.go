package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type RepositoryValidation struct {
	RepoName                string
	HasValidKebabCaseNaming bool
	HasReadme               bool
	HasCodeOwners           bool
	HasEditorConfig         bool
}

// isValidKebabCase checks if the repository name adheres to kebab-case.
func hasValidKebabCaseNaming(name string) bool {
	match, _ := regexp.MatchString(`^[a-z0-9]+(-[a-z0-9]+)*$`, name)
	return match
}

func hasEditorConfig(files []*github.RepositoryContent) bool {
	for _, content := range files {
		if content.GetType() == "file" && strings.EqualFold(content.GetName(), ".ediorconfig") {
			return true
		}
	}
	return false
}

func hasReadme(files []*github.RepositoryContent) bool {
	for _, content := range files {
		if content.GetType() == "file" && strings.EqualFold(content.GetName(), "README.md") {
			return true
		}
	}
	return false
}

func hasCodeOwners(files []*github.RepositoryContent) bool {
	for _, content := range files {
		if content.GetType() == "file" && strings.EqualFold(content.GetName(), "CODEOWNERS") {
			return true
		}
	}

	return false
}

func getRepositoriesContents(ctx context.Context, client *github.Client, organisation, repo, path string) ([]*github.RepositoryContent, error) {
	_, directoryContents, _, err := client.Repositories.GetContents(ctx, organisation, repo, path, nil)
	if err != nil {
		return nil, err
	}

	if directoryContents == nil {
		return nil, fmt.Errorf("no contents found at path: %s", path)
	}

	return directoryContents, nil
}

func generateReport(validations []RepositoryValidation) {
	fmt.Println("Generating report")

	// Read the HTML template from an external file
	templateFile := "report-template.html"
	tmplContent, err := os.ReadFile(templateFile)
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return
	}

	// Create an HTML file
	file, err := os.Create("report.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Parse and execute the template
	tmpl, err := template.New("report").Parse(string(tmplContent))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = tmpl.Execute(file, validations)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Report generated successfully.")

}

func main() {

	// check if args set
	if len(os.Args) < 3 {
		fmt.Println("Parameters missing: <github-access-token> <github-organisation>")
		fmt.Println("Example go run main.go xXe5iY6wsBTE0Ziu3Ln5ZT... tutuka")
		return
	}

	// github token input
	// github org name input
	token := os.Args[1]
	org := os.Args[2]

	// clients definition
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// validation types
	var repoValidations []RepositoryValidation

	opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 10}}

	fmt.Printf("Standard validation routine started")

	// Loop through all pages of results
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			fmt.Printf("Error fetching repositories: %v\n", err)
			return
		}

		for _, repo := range repos {
			repositoryContents, error := getRepositoriesContents(ctx, client, org, repo.GetName(), "")
			if error != nil {
				fmt.Printf("Error fetching content of the repository: %v\n with error: %v\n", repo.GetName(), error)
				continue
			}
			repoName := repo.GetName()
			repoValid := RepositoryValidation{
				RepoName:                repoName,
				HasReadme:               hasReadme(repositoryContents),
				HasCodeOwners:           hasCodeOwners(repositoryContents),
				HasEditorConfig:         hasEditorConfig(repositoryContents),
				HasValidKebabCaseNaming: hasValidKebabCaseNaming(repoName),
			}
			repoValidations = append(repoValidations, repoValid)
		}

		// Check if there are more pages to fetch
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	generateReport(repoValidations)

}
