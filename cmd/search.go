package cmd

import (
	"net/url"
	"strings"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

// Search is executed when you run `stew search`
func Search(cliInput []string) {
	sp := constants.LoadingSpinner

	if len(cliInput) == 0 {
		stew.CatchAndExit(stew.EmptyCLIInputError{})
	}
	for _, input := range cliInput {
		err := stew.ValidateGithubSearchQuery(input)
		stew.CatchAndExit(err)
	}

	searchQuery := url.QueryEscape(strings.Join(cliInput, " "))

	sp.Start()
	githubSearch, err := stew.NewGithubSearch(searchQuery)
	sp.Stop()
	stew.CatchAndExit(err)

	if len(githubSearch.Items) == 0 {
		stew.CatchAndExit(stew.NoGithubSearchResultsError{SearchQuery: githubSearch.SearchQuery})
	}

	formattedSearchResults := stew.FormatSearchResults(githubSearch)

	githubProjectName, err := stew.PromptSelect("Choose a GitHub project:", formattedSearchResults)
	stew.CatchAndExit(err)

	searchResultIndex, _ := stew.Contains(formattedSearchResults, githubProjectName)

	Install(githubSearch.Items[searchResultIndex].FullName)

}
