package alias

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Storage struct {
	Aliases map[string]string `yaml:"aliases"`
}

// AliasStorage holds config dir function for testability
// and provides methods for alias management.
type AliasStorage struct {
	UserConfigDir func() (string, error)
}

func NewAliasStorage() *AliasStorage {
	return &AliasStorage{
		UserConfigDir: os.UserConfigDir,
	}
}

func (a *AliasStorage) findAliasesFilePath() (string, error) {
	userConfigDir, err := a.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("error locating user config dir: %w", err)
	}
	aliasPath := filepath.Join(userConfigDir, "tl", "aliases.yml")
	return aliasPath, nil
}

func (a *AliasStorage) GetAlias(alias string) (string, error) {
	aliases, err := a.LoadAliases()
	if err != nil {
		return "", err
	}
	if _, ok := aliases[alias]; !ok {
		return "", fmt.Errorf("alias %s not found", alias)
	}
	return aliases[alias], nil
}

func (a *AliasStorage) SetAlias(alias, value string) error {
	aliases, err := a.LoadAliases()
	if err != nil {
		return err
	}
	aliases[alias] = value
	return a.SaveAliases(aliases)
}

func (a *AliasStorage) SaveAliases(aliases map[string]string) error {
	path, err := a.findAliasesFilePath()
	if err != nil {
		return err
	}
	configDir := filepath.Dir(path)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating user config dir %s: %w", configDir, err)
		}
	}
	storage := Storage{
		Aliases: aliases,
	}
	yamlData, err := yaml.Marshal(&storage)
	if err != nil {
		return fmt.Errorf("error while marshaling: %w", err)
	}
	err = os.WriteFile(path, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing aliases file %s: %w", path, err)
	}
	return nil
}

func (a *AliasStorage) LoadAliases() (map[string]string, error) {
	path, err := a.findAliasesFilePath()
	if err != nil {
		return map[string]string{}, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}, fmt.Errorf("error reading aliases file %s: %w", path, err)
	}
	var storage Storage
	err = yaml.Unmarshal(yamlFile, &storage)
	if err != nil {
		return map[string]string{}, fmt.Errorf("error unmarshalling aliases file %s: %w", path, err)
	}
	return storage.Aliases, nil
}

// ResolveAlias returns the resolved issue key if the input is an alias, otherwise returns the input.
func ResolveAlias(input string) string {
	aliasStorage := NewAliasStorage()
	aliases, err := aliasStorage.LoadAliases()
	if err == nil {
		if resolvedKey, ok := aliases[input]; ok {
			return resolvedKey
		}
	}
	return input
}
