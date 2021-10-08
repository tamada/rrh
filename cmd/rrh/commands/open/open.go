package open

import (
	"fmt"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "open the folder or web page of the given repositories",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, performOpen)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&openOpts.webpageFlag, "browser", "b", false, "open the web page of the repository")
	return cmd
}

type openOptions struct {
	webpageFlag bool
}

var openOpts = &openOptions{}

func convertToRepositoryURL(url string) (string, error) {
	str := strings.TrimPrefix(url, "git@")
	str = strings.TrimSuffix(str, ".git")
	index := strings.Index(str, ":")
	if index < 0 {
		return "", fmt.Errorf("%s: unrecognized git repository url", url)
	}
	host := str[0:index]
	return "https://" + host + "/" + str[index+1:], nil
}

func convertURL(url string) (string, error) {
	if strings.HasPrefix(url, "git@") {
		convertedURL, err := convertToRepositoryURL(url)
		if err != nil {
			return "", err
		}
		url = convertedURL
	}
	if strings.HasPrefix(url, "https") && strings.HasSuffix(url, ".git") {
		url = strings.TrimSuffix(url, ".git")
	}
	return url, nil
}

func generateWebPageURL(repo *rrh.Repository) (string, error) {
	if len(repo.Remotes) == 0 {
		return "", fmt.Errorf("%s: remote repository not found", repo.ID)
	}
	return convertURL(repo.Remotes[0].URL)
}

func execOpen(repo *rrh.Repository) (string, error) {
	if openOpts.webpageFlag {
		return generateWebPageURL(repo)
	}
	return repo.Path, nil
}

func performEach(arg string, db *rrh.Database) error {
	repo := db.FindRepository(arg)
	if repo == nil {
		return fmt.Errorf("%s: repository not found", arg)
	}
	path, err := execOpen(repo)
	if err != nil {
		return err
	}
	return open.Start(path)
}

func performOpen(c *cobra.Command, args []string, db *rrh.Database) error {
	el := common.NewErrorList()
	for _, arg := range args {
		err := performEach(arg, db)
		el = el.Append(err)
	}
	return el.NilOrThis()
}
