package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"go-repository-checker/internal/types"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Color codes
const (
	Reset = "\033[0m"
	Green = "\033[32m"
	Red   = "\033[31m"
)

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

func hasGradleBuild(files []*github.RepositoryContent) bool {
	for _, content := range files {
		if content.GetType() == "file" && strings.EqualFold(content.GetName(), "build.gradle") {
			return true
		}
	}

	return false
}

func hasMavenBuild(files []*github.RepositoryContent) bool {
	for _, content := range files {
		if content.GetType() == "file" && strings.EqualFold(content.GetName(), "pom.xml") {
			return true
		}
	}
	return false
}



func handleScan(org *string, repo *string, token *string) {

	if *org == "" || *token == "" {
		fmt.Println("The --org and --token flags are required for the scan command")
		flag.Usage()
		os.Exit(1)
	}

	// clients definition
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if *repo != "" {
		scanRepository(ctx, ts, tc, client, repo, org)
	} else {
		scanOrganisation(ctx, ts, tc, client, org)
	}
}

func handleReport(output *string, format *string, token *string) {
	// Validate required flags
	if *output == "" || *token == "" {
		fmt.Println("The --output and --token flags are required for the report command")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Generating report at %s in %s format with token: %s\n", *output, *format, *token)

}

func scanRepository(ctx context.Context, ts oauth2.TokenSource, tc *http.Client, client *github.Client, repo *string, org *string) {

	var repositoryContents, error = getRepositoriesContents(ctx, client, *org, *repo, "")

	if error != nil {
		fmt.Printf("Error fetching repository: %v\n", error)
		return
	}

	repoValid := types.RepositoryValidation{
		RepoName:                *repo,
		HasReadme:               hasReadme(repositoryContents),
		HasCodeOwners:           hasCodeOwners(repositoryContents),
		HasEditorConfig:         hasEditorConfig(repositoryContents),
		HasValidKebabCaseNaming: hasValidKebabCaseNaming(*repo),
		HasMavenBuild:           hasMavenBuild(repositoryContents),
		HasBuildGradle:          hasGradleBuild(repositoryContents),
	}

	
	// Print rows
	fmt.Printf("Repository name 		%s\n", repoValid.RepoName)
	fmt.Printf("Has kebab-case naming 	\t%s\n", formatResult(repoValid.HasValidKebabCaseNaming))
	fmt.Printf("Has README.MD 			%s\n", formatResult(repoValid.HasReadme))
	fmt.Printf("Has CODEOWNERS 			%s\n", formatResult(repoValid.HasCodeOwners))
	fmt.Printf("Has .editorconfig 		%s\n", formatResult(repoValid.HasEditorConfig))
	fmt.Printf("Has build.gradle 		%s\n", formatResult(repoValid.HasBuildGradle))
	fmt.Printf("Has pom.xml 			%s\n", formatResult(repoValid.HasMavenBuild))
}

// formatResult formats the boolean result with colors
func formatResult(result bool) string {
	if result {
		return fmt.Sprintf("%s[PASSED]%s", Green, Reset)
	}
	return fmt.Sprintf("%s[FAIL]%s", Red, Reset)
}

func scanOrganisation(ctx context.Context, ts oauth2.TokenSource, tc *http.Client, client *github.Client, org *string) {
	// couters
	itemsCounter := 0

	// validation types
	var repoValidations []types.RepositoryValidation

	opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 10}}

	fmt.Println("Standard validation routine started")

	// Loop through all pages of results
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, *org, opt)

		if err != nil {
			fmt.Printf("Error fetching repositories: %v\n", err)
			return
		}

		for _, repo := range repos {
			repositoryContents, error := getRepositoriesContents(ctx, client, *org, repo.GetName(), "")
			if error != nil {
				fmt.Printf("Error fetching content of the repository: %v\n with error: %v\n", repo.GetName(), error)
				continue
			}
			repoName := repo.GetName()
			repoValid := types.RepositoryValidation{
				RepoName:                repoName,
				HasReadme:               hasReadme(repositoryContents),
				HasCodeOwners:           hasCodeOwners(repositoryContents),
				HasEditorConfig:         hasEditorConfig(repositoryContents),
				HasValidKebabCaseNaming: hasValidKebabCaseNaming(repoName),
				HasMavenBuild:           hasMavenBuild(repositoryContents),
				HasBuildGradle:          hasGradleBuild(repositoryContents),
			}
			repoValidations = append(repoValidations, repoValid)
		}

		// add number of items to counter
		itemsCounter += len(repos)
		if itemsCounter%10 == 0 {
			fmt.Printf("#num of processed repositories: %d\n", itemsCounter)
		}

		// Check if there are more pages to fetch
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}

func main() {

	// Define flags for the scan command
	var scanCmd = flag.NewFlagSet("scan", flag.ExitOnError)
	var org = scanCmd.String("org", "", "Specify the organization name (required)")
	var repo = scanCmd.String("repo", "", "Specify a specific repository to scan")
	var token = scanCmd.String("token", "", "Specify the authentication token (required)")

	// Define flags for the report command
	var reportCmd = flag.NewFlagSet("report", flag.ExitOnError)
	output := flag.String("output", "", "Specify the output file for the report (required)")
	format := flag.String("format", "json", "Specify the report format (choices: json, html, csv; default: json)")

	// Parse the flags
	flag.Parse()

	// Check the command and handle accordingly
	if len(os.Args) < 2 {
		fmt.Println("scan or report subcommand is required")
		flag.Usage()
		os.Exit(1)
	}

	// check top-level command
	command := os.Args[1]

	switch command {
	case "scan":
		scanCmd.Parse(os.Args[2:])
		handleScan(org, repo, token)
	case "report":
		reportCmd.Parse(os.Args[1:])
		handleReport(output, format, token)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		flag.Usage()
		os.Exit(1)
	}

}
