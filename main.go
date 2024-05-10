package main

import (
    "context"
    "fmt"
    "regexp"
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
	"os"
)

// isValidKebabCase checks if the repository name adheres to kebab-case.
func isValidKebabCase(name string) bool {
    match, _ := regexp.MatchString(`^[a-z0-9]+(-[a-z0-9]+)*$`, name)
    return match
}

func main() {

    // check if args set
    if len(os.Args) < 3 {
        fmt.Println("Parameters missing: <github-access-token> <github-organisation>")
        fmt.Println("Example go run main.go xXe5iY6wsBTE0Ziu3Ln5ZT... tutuka")
        return
    }

    // github token
    token: = os.Args[1]
        // github org name  
    org: = os.Args[2]

    ctx: = context.Background()

    ts: = oauth2.StaticTokenSource( & oauth2.Token {
        AccessToken: token
    }, )
    tc: = oauth2.NewClient(ctx, ts)

    client: = github.NewClient(tc)

    opt: = & github.RepositoryListByOrgOptions {
        ListOptions: github.ListOptions {
            PerPage: 10
        },
    }

    // all repo counters
    allRepoCounter: = 0
    kebabCaseRepositoryCounter: = 0
    nonStandardRepositoryCounter: = 0

    // Loop through all pages of results
    for {
        repos, resp, err: = client.Repositories.ListByOrg(ctx, org, opt)
        if err != nil {
            fmt.Printf("Error fetching repositories: %v\n", err)
            return
        }

        allRepoCounter = allRepoCounter + len(repos)

        for _, repo: = range repos {
            repoName: = repo.GetName()
            validCase: = isValidKebabCase(repoName)

            if validCase {
                kebabCaseRepositoryCounter++;
            } else {
                nonStandardRepositoryCounter++;
            }

            fmt.Printf("Repository: %s has valid repository name: %t\n", repoName, validCase)
        }

        // Check if there are more pages to fetch
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }

    // break-line for formatting 
    fmt.Println()
    fmt.Printf("all respos counter :%d\n", allRepoCounter)
    fmt.Printf("kebab-case repos counter: :%d\n", kebabCaseRepositoryCounter)
    fmt.Printf("non-standard repos counter: :%d\n", nonStandardRepositoryCounter)

}