package command

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence"
	"github.com/spf13/cobra"
)

func NewRootCmd() (Cmd, error) {
	var (
		pluginDir string
		rootCmd   = &cobra.Command{
			Use:     "sqnc [--plugins DIR] [-h] [command]",
			Version: sequence.Semver(),
			Args:    cobra.NoArgs,
			PersistentPreRun: func(cmd *cobra.Command, _ []string) {
				for _, dir := range js.Unique(
					js.Filter(
						[]string{
							pluginDir,
							"/etc/sqnc/plugins",
							os.Getenv(EnvVarPlugins),
						}, func(s string, _ int, _ []string) bool {
							return s != ""
						},
					),
				) {
					if err := OpenPlugins(dir); err != nil {
						cmd.PrintErrln(err)
					}
				}
			},
		}
	)

	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ .Name }}{{ .Version }} %s\n", runtime.Version()),
	)

	var (
		pluginsValue = ""
		homeDir, err = os.UserHomeDir()
	)
	if err == nil {
		pluginsValue = filepath.Join(homeDir, ".sqnc/plugins")
	}

	flags := rootCmd.Flags()
	flags.StringVar(&pluginDir, "plugins", pluginsValue, "plugin directory")
	if err := rootCmd.MarkFlagDirname("plugins"); err != nil {
		return nil, err
	}

	for _, newChildCmd := range []func() (Cmd, error){
		NewRunCmd,
		NewPruneCmd,
	} {
		childCmd, err := newChildCmd()
		if err != nil {
			return nil, err
		}

		rootCmd.AddCommand(childCmd.(*cobra.Command))
	}

	return rootCmd, nil
}
