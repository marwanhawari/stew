package cmd

import (
	stew "github.com/marwanhawari/stew/lib"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/errs"
	"github.com/marwanhawari/stew/lib/ui/progress"
	"github.com/marwanhawari/stew/lib/ui/prompt"
)

// Search is executed when you run `stew search`
func Search(cliInput string) {
	rt := errs.Strip(config.Initialize())
	sp := progress.Spinner(rt)

	err := stew.ValidateCLIInput(cliInput)
	errs.MaybeExit(err)

	err = stew.ValidateGithubSearchQuery(cliInput)
	errs.MaybeExit(err)

	sp.Start()
	githubSearch, err := stew.NewGithubSearch(rt.Config, cliInput)
	sp.Stop()
	errs.MaybeExit(err)

	if len(githubSearch.Items) == 0 {
		errs.MaybeExit(stew.NoGithubSearchResultsError{SearchQuery: githubSearch.SearchQuery})
	}

	formattedSearchResults := stew.FormatSearchResults(githubSearch)

	githubProjectName, err := prompt.Select(rt, "Choose a GitHub project:", formattedSearchResults)
	errs.MaybeExit(err)

	searchResultIndex, _ := stew.Contains(formattedSearchResults, githubProjectName)

	Browse(githubSearch.Items[searchResultIndex].FullName)
}
