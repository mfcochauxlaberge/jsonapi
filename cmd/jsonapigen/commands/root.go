package commands

import (
	"fmt"
	"os"
	"reflect"


	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jsonapigen",
	Short: "Code generator for mfcochauxlaberge/jsonapi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Reading Go files...\n")

		vals := []reflect.Value{}

		// vals, err := reader.ReadFile([]byte(""))
		// if err != nil {
		// 	panic(err)
		// }

		fmt.Printf("The following %d structs will be considered\n", 0)
		for range vals {
		}

		fmt.Printf("Generating Go files...\n")
		fmt.Printf("Done.\n")
	},
}

// Execute ...
func Execute() {
	rootCmd.AddCommand(
		cmdVersion,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
