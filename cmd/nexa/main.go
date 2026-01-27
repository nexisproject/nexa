// Copyright (C) micros. 2025-present.
//
// Created at 2025-02-10, by liasica

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"nexis.run/nexa/cmd/nexa/internal/base"
	"nexis.run/nexa/cmd/nexa/internal/command"
)

var (
	Version   = "0.1.0"
	BuildTime string
	Hash      string
)

func getVersion() string {
	return fmt.Sprintf("%s.%s (built at %s)", Version, Hash, BuildTime)
}

func main() {
	var (
		configFile string
	)

	cmd := cobra.Command{
		Use:               "nexa",
		Short:             "NEXA 框架实用工具",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Version:           getVersion(),
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			// 初始化变量
			err := base.InitializeConfig(configFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	configGroup, configCommand := command.ConfigCmd()
	createGroup, createCommand := command.NewCmd()
	entGroup, entCommand := command.EntCmd()

	cmd.AddGroup(
		configGroup,
		createGroup,
		entGroup,
	)

	cmd.AddCommand(
		configCommand,
		createCommand,
		entCommand,
	)

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", ".nexa.yaml", "配置文件")

	err := cmd.Execute()
	if err != nil {
		fmt.Printf("command execution failed: %v\n", err)
		os.Exit(1)
	}
}
