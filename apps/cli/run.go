package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mehmetalidsy/madget-cli/internal/archive"
	"github.com/mehmetalidsy/madget-cli/internal/config"
	"github.com/mehmetalidsy/madget-cli/internal/integrity"
	"github.com/spf13/cobra"
)

type resolveResponse struct {
	Package  string `json:"package"`
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
	Tarball  string `json:"tarball_url"`
}

func Run() error {
	var cfgPath string
	var registryURL string
	var token string

	rootCmd := &cobra.Command{
		Use:   "madget",
		Short: "Madget v2 CLI",
	}

	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", ".madget/config.json", "Config file path")
	rootCmd.PersistentFlags().StringVar(&registryURL, "registry", "http://localhost:8080", "Registry API URL")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Publisher token")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize local Madget config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.LocalConfig{
				RegistryURL: registryURL,
			}
			return config.Save(cfgPath, cfg)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Store publisher token in config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return fmt.Errorf("token is required: pass --token")
			}
			cfg, _ := config.Load(cfgPath)
			cfg.RegistryURL = registryURL
			cfg.Token = token
			return config.Save(cfgPath, cfg)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "publish <MadGet.xml> <tarball>",
		Short: "Publish package using MadGet.xml manifest and tarball",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			return publishPackage(cfg, args[0], args[1])
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "install <name>@<range>",
		Short: "Resolve, download, verify and extract package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			spec := args[0]
			name, versionRange, err := parseInstallSpec(spec)
			if err != nil {
				return err
			}

			resolveURL := fmt.Sprintf("%s/v1/packages/%s/resolve?range=%s", cfg.RegistryURL, name, versionRange)
			resp, err := http.Get(resolveURL)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 300 {
				return fmt.Errorf("resolve failed with status %d", resp.StatusCode)
			}

			var resolved resolveResponse
			if err := json.NewDecoder(resp.Body).Decode(&resolved); err != nil {
				return err
			}
			tarballURL := resolved.Tarball
			if len(tarballURL) > 0 && tarballURL[0] == '/' {
				tarballURL = strings.TrimRight(cfg.RegistryURL, "/") + tarballURL
			}

			tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s.tgz", resolved.Package, resolved.Version))
			if err := archive.Download(tarballURL, tmpPath); err != nil {
				return err
			}
			if err := integrity.VerifySHA256(tmpPath, resolved.Checksum); err != nil {
				return err
			}

			targetDir := filepath.Join("vendor", resolved.Package, resolved.Version)
			if err := archive.ExtractTarGz(tmpPath, targetDir); err != nil {
				return err
			}

			fmt.Printf("Installed %s@%s to %s\n", resolved.Package, resolved.Version, targetDir)
			return nil
		},
	})

	return rootCmd.Execute()
}
