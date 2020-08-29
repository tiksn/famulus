package main

import (
	"fmt"
	"log"
	"os"

	config "github.com/tiksn/famulus/internal/app/famulus"
	root "github.com/tiksn/famulus/pkg/famulus/cmd"
)

func main() {
	c, err := config.ParseDefaultConfig()
	if err != nil {
		log.Fatalln(err)
	}

	expandedArgs := []string{}
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	rootCmd := root.NewCmdRoot(&c)
	cmd, _, err := rootCmd.Traverse(expandedArgs)
	if err != nil || cmd == rootCmd {
	}

	rootCmd.SetArgs(expandedArgs)

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		fmt.Fprintln(os.Stderr, cmd.UsageString())
	}
}
