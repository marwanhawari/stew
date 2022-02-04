package cmd

import (
	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

func Search(cliInput string) {
	sp := constants.LoadingSpinner

	err := stew.ValidateCLIInput(cliInput)
	stew.CatchAndExit(err)

	err = stew.ValidateGithubSearchQuery(cliInput)
	stew.CatchAndExit(err)

	sp.Start()
	githubSearch, err := stew.NewGithubSearch(cliInput)
	sp.Stop()
	stew.CatchAndExit(err)

	if len(githubSearch.Items) == 0 {
		stew.CatchAndExit(stew.NoGithubSearchResultsError{SearchQuery: githubSearch.SearchQuery})
	}

	formattedSearchResults := stew.FormatSearchResults(githubSearch)

	githubProjectName, err := stew.PromptSelect("Choose a GitHub project:", formattedSearchResults)
	stew.CatchAndExit(err)

	searchResultIndex, _ := stew.Contains(formattedSearchResults, githubProjectName)

	Browse(githubSearch.Items[searchResultIndex].FullName)

}
