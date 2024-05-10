# GitHub Repository Name Checker

Simple github organisation repostiroy checker that fetches all the repositories within an organization to ensure they adhere to the standard naming convention of kebab-case (lowercase-words-separated-by-dashes). 

It's designed to help maintain consistency and clarity in repository naming across organization.

## Prerequisites

- Go (version 1.22 or higher) - probably lower version would work too but that was not tested.
- A GitHub Personal Access Token with `repo` scope

## Installation

Clone the repository and run:

```bash
git clone https://github.com/your-org/repo-name-checker.git
cd repo-name-checker
go run main.go <token> <organisation-name>
```

Results
```bash
all respos counter :315
kebab-case repos counter: :186
non-standard repos counter: :129
```