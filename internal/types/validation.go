package types

type RepositoryValidation struct {
	RepoName                string
	HasValidKebabCaseNaming bool
	HasReadme               bool
	HasCodeOwners           bool
	HasEditorConfig         bool
	HasBuildGradle          bool
	HasMavenBuild           bool
}
