package setup

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/db"
)

var (
	cmdLong    = `Initialize the tl configuration file and database.`
	cmdExample = `
  # Setup tl
  tl setup --jira-url https://your-domain.atlassian.net --username yourusername --token yourapitoken`
	jiraURL  string
	username string
	token    string
)

func NewCommand(r func() db.RepositoryInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "setup",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "Setup tl database and configuration",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := r().AutoMigrate()
			if err != nil {
				return err
			}

			viper.Set("jira_base_url", strings.TrimSpace(jiraURL))
			viper.Set("jira_username", strings.TrimSpace(username))
			viper.Set("jira_api_token", strings.TrimSpace(token))

			err = viper.WriteConfig()
			if err != nil {
				return fmt.Errorf("failed to write config file: %v", err)
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Setup complete.")
			return nil
		},
	}

	cmd.Flags().StringVar(&jiraURL, "jira-url", "", "Jira URL")
	_ = cmd.MarkFlagRequired("jira-url")
	cmd.Flags().StringVar(&username, "username", "", "Jira username")
	_ = cmd.MarkFlagRequired("username")
	cmd.Flags().StringVar(&token, "token", "", "Jira API token")
	_ = cmd.MarkFlagRequired("token")

	return cmd
}
