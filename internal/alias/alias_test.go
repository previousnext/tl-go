package alias

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetAlias(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &AliasStorage{
		UserConfigDir: func() (string, error) { return tmpDir, nil },
	}
	aliasName := "foo"
	aliasValue := "bar"
	err := storage.SetAlias(aliasName, aliasValue)
	assert.NoError(t, err, "SetAlias failed")
	val, err := storage.GetAlias(aliasName)
	assert.NoError(t, err, "GetAlias failed")
	assert.Equal(t, aliasValue, val, "expected alias value")
}

func TestGetAliasNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &AliasStorage{
		UserConfigDir: func() (string, error) { return tmpDir, nil },
	}
	_, err := storage.GetAlias("doesnotexist")
	assert.Error(t, err, "expected error for missing alias")
	assert.Equal(t, "alias doesnotexist not found", err.Error())
}

func TestLoadAliases_FileDoesNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &AliasStorage{
		UserConfigDir: func() (string, error) { return tmpDir, nil },
	}
	aliases, err := storage.LoadAliases()
	assert.NoError(t, err, "LoadAliases failed")
	assert.Len(t, aliases, 0, "expected empty aliases")
}

func TestSaveAndLoadAliases(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &AliasStorage{
		UserConfigDir: func() (string, error) { return tmpDir, nil },
	}
	aliases := map[string]string{"a": "1", "b": "2"}
	err := storage.SaveAliases(aliases)
	assert.NoError(t, err, "SaveAliases failed")
	loaded, err := storage.LoadAliases()
	assert.NoError(t, err, "LoadAliases failed")
	assert.Len(t, loaded, 2, "unexpected loaded aliases count")
	assert.Equal(t, "1", loaded["a"])
	assert.Equal(t, "2", loaded["b"])
}

func TestSaveAliasesCreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &AliasStorage{
		UserConfigDir: func() (string, error) { return tmpDir, nil },
	}
	configDir := filepath.Join(tmpDir, "tl")
	err := os.RemoveAll(configDir)
	assert.True(t, err == nil || os.IsNotExist(err), "RemoveAll failed: %v", err)
	aliases := map[string]string{"x": "y"}
	err = storage.SaveAliases(aliases)
	assert.NoError(t, err, "SaveAliases failed")
	aliasFile := filepath.Join(configDir, "aliases.yml")
	_, err = os.Stat(aliasFile)
	assert.NoError(t, err, "aliases file not created")
}

func TestLoadAliasesInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &AliasStorage{
		UserConfigDir: func() (string, error) { return tmpDir, nil },
	}
	configDir := filepath.Join(tmpDir, "tl")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err, "MkdirAll failed")
	aliasFile := filepath.Join(configDir, "aliases.yml")
	err = os.WriteFile(aliasFile, []byte("not: yaml: here: ["), 0644)
	assert.NoError(t, err, "WriteFile failed")
	_, err = storage.LoadAliases()
	assert.Error(t, err, "expected error for invalid YAML")
}

func TestFindAliasesFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	storage := &AliasStorage{
		UserConfigDir: func() (string, error) { return tmpDir, nil },
	}
	path, err := storage.findAliasesFilePath()
	assert.NoError(t, err, "findAliasesFilePath failed")
	assert.Equal(t, filepath.Join(tmpDir, "tl"), filepath.Dir(path), "unexpected path")
}
