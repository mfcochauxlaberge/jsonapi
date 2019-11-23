package commands

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	beSimple *bool
	showDeps *bool
)

func init() {
	beSimple = cmdVersion.Flags().BoolP("simple", "", false, "simply print the version")
	showDeps = cmdVersion.Flags().BoolP("deps", "", false, "show dependencies")
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the version number",
	Run: func(cmd *cobra.Command, args []string) {
		bi, _ := debug.ReadBuildInfo()

		if *beSimple {
			fmt.Printf("%s \n", bi.Main.Version)
			return
		}

		fmt.Printf("%s %s\n", bi.Main.Path, bi.Main.Version)

		if *showDeps {
			fmt.Printf("\n")
			for _, mod := range bi.Deps {
				fmt.Printf("%s %s\n", mod.Path, mod.Version)
			}
		}
	},
}
