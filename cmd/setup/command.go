package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/db"
)

var (
	cmdLong    = `Initialize the tl configuration file and database.`
	cmdExample = `
  # Initialize tl
  tl init`
)

func NewCommand(r func() db.RepositoryInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "init",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "Initialize tl",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := r().InitRepository()
			if err != nil {
				return err
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter the Jira URL: ")
			url, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			fmt.Print("Enter the Jira username: ")
			username, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			fmt.Print("Enter the Jira API token: ")
			token, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			viper.Set("jira_base_url", strings.TrimSpace(url))
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

	return cmd
}
