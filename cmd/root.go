package cmd

import (
	"github.com/spf13/cobra"
)

func GetRootCmd(args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "kubeschema",
		Short:             "generate json schema of kubernetes resources",
		DisableAutoGenTag: true,
		Args:              cobra.ArbitraryArgs,
	}
	rootCmd.SetArgs(args)
	rootCmd.AddCommand(NewCrdCmd(args))
	rootCmd.AddCommand(NewIndexCmd(args))
	rootCmd.AddCommand(NewDumpCmd(args))
	// rootCmd.AddCommand(NewOpenapiCmd(args))
	return rootCmd
}
