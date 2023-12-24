package main

import (
	"fmt"
	"os"

	"github.com/imroc/kubeschema/cmd"
)

func main() {
	rootCmd := cmd.GetRootCmd(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
