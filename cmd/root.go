/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/cmd/create"
	"github.com/previousnext/tl-go/cmd/delete"
	"github.com/previousnext/tl-go/cmd/list"
	"github.com/previousnext/tl-go/cmd/setup"
	"github.com/previousnext/tl-go/cmd/show"
	"github.com/previousnext/tl-go/cmd/update"
	"github.com/previousnext/tl-go/internal/db"
)

var cfgFile string
var dbFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tl",
	Short: "A command line too for logging time.",
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/tl/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&dbFile, "db", "", "db file (default is ~/.local/share/tl/db.sqlite)")

	// Hide the help command.
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Hide the completions command.
	rootCmd.CompletionOptions = cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	}

	// We need to lazy initialize the repository to ensure that the
	// configuration is loaded before we create the repository.
	r := func() db.RepositoryInterface {
		return db.NewRepository(viper.GetString("db_file"))
	}

	rootCmd.AddCommand(setup.NewCommand(r))
	rootCmd.AddCommand(create.NewCommand(r))
	rootCmd.AddCommand(show.NewCommand(r))
	rootCmd.AddCommand(list.NewCommand(r))
	rootCmd.AddCommand(update.NewCommand(r))
	rootCmd.AddCommand(delete.NewCommand(r))
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
