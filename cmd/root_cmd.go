package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/cmd/create"
	"github.com/previousnext/tl-go/cmd/delete"
	"github.com/previousnext/tl-go/cmd/list"
	"github.com/previousnext/tl-go/cmd/send"
	"github.com/previousnext/tl-go/cmd/setup"
	"github.com/previousnext/tl-go/cmd/show"
	"github.com/previousnext/tl-go/cmd/update"
	"github.com/previousnext/tl-go/internal/api"
	"github.com/previousnext/tl-go/internal/api/types"
	"github.com/previousnext/tl-go/internal/db"
)

var cfgFile string
var dbFile string

// version overridden at build time by:
//
//	-ldflags="-X github.com/previousnext/tl-go/cmd.version=$(git describe --tags --always)"
var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tl",
	Short:   "A command-line tool for logging time.",
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/tl/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&dbFile, "db", "", "db file (default is ~/.local/share/tl/db.sqlite)")

	// Hide the help command.
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Hide the completions command.
	rootCmd.CompletionOptions = cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	}

	// We need to lazy-initialise injected dependencies so that
	// viper has a chance to read in the config.
	repository := func() db.RepositoryInterface {
		return db.NewRepository(viper.GetString("db_file"))
	}
	jiraClient := func() api.JiraClientInterface {
		params := types.JiraClientParams{
			BaseURL:  viper.GetString("jira_base_url"),
			Username: viper.GetString("jira_username"),
			APIToken: viper.GetString("jira_api_token"),
		}
		httpClient := &http.Client{}
		return api.NewJiraClient(httpClient, params)
	}

	rootCmd.AddCommand(setup.NewCommand(repository))
	rootCmd.AddCommand(create.NewCommand(repository))
	rootCmd.AddCommand(show.NewCommand(repository))
	rootCmd.AddCommand(list.NewCommand(repository))
	rootCmd.AddCommand(update.NewCommand(repository))
	rootCmd.AddCommand(delete.NewCommand(repository))
	rootCmd.AddCommand(send.NewCommand(repository, jiraClient))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find userConfigDir directory.
		userConfigDir, err := os.UserConfigDir()
		cobra.CheckErr(err)
		configPath := filepath.Join(userConfigDir, "tl")

		err = os.MkdirAll(configPath, 0755)
		cobra.CheckErr(err)

		// Search config in user config + tl directory with name "config" (without extension).
		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if dbFile != "" {
		// Use the db file from the flag.
		viper.Set("db_file", dbFile)
	} else {
		// Find userDataDir directory.
		userDataDir, err := os.UserConfigDir()
		cobra.CheckErr(err)
		dataPath := filepath.Join(userDataDir, "tl")

		err = os.MkdirAll(dataPath, 0755)
		cobra.CheckErr(err)

		// Set the db file path.
		viper.Set("db_file", filepath.Join(dataPath, "db.sqlite"))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
