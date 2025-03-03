package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bcp-innovations/hyperlane-cosmos/tests/simapp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed sample_config/*
var sampleConfig embed.FS

func InitSampleChain() *cobra.Command {
	return &cobra.Command{
		Use:   "init-sample-chain",
		Short: "Initializes a dummy chain which can be directly started",
		Run: func(cmd *cobra.Command, args []string) {
			destPath := viper.GetString("home")
			if destPath == "" {
				destPath = simapp.DefaultNodeHome
			}

			if !CheckDirIsDefault(destPath) {
				panic("directory already initialized.")
			}

			if err := os.Remove(filepath.Join(destPath, "config", "app.toml")); err != nil {
				panic(err)
			}
			if err := os.Remove(filepath.Join(destPath, "config", "client.toml")); err != nil {
				panic(err)
			}
			if err := os.Remove(filepath.Join(destPath, "config", "config.toml")); err != nil {
				panic(err)
			}

			if err := os.Remove(filepath.Join(destPath, "config")); err != nil {
				panic(err)
			}
			if err := os.Remove(filepath.Join(destPath, "data")); err != nil {
				panic(err)
			}

			if err := copyEmbedToDisk(sampleConfig, "sample_config", destPath); err != nil {
				panic(err)
			}

			fmt.Printf("Initialized sample chain. Run ./hypd start --home %s\n", destPath)
		},
	}
}

// CheckDirIsDefault checks if the provided directory is the default directory.
// Unfortunately, cosmos or comet always creates the default directory which we want
// to override.
func CheckDirIsDefault(path string) bool {
	dir, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	if len(dir) != 2 {
		return false
	}

	if !(dir[0].IsDir() && dir[0].Name() == "config" && dir[1].IsDir() && dir[1].Name() == "data") {
		return false
	}

	configDir, err := os.ReadDir(filepath.Join(path, "config"))
	if err != nil {
		panic(err)
	}

	if len(configDir) != 3 {
		return false
	}

	if !(configDir[0].Name() == "app.toml" && configDir[1].Name() == "client.toml" && configDir[2].Name() == "config.toml") {
		return false
	}

	return true
}

func copyEmbedToDisk(embedFS embed.FS, embedDir, targetDir string) error {
	err := fs.WalkDir(embedFS, embedDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(embedDir, path)
		destPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, os.ModePerm)
		}

		// Copy file
		srcFile, err := embedFS.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})

	return err
}
