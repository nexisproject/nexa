// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-27, by liasica

package command

import "github.com/spf13/cobra"

func EntCmd() (*cobra.Group, *cobra.Command) {
	g := &cobra.Group{
		ID:    "ent",
		Title: "ent 相关命令",
	}

	cmd := &cobra.Command{
		Use:               "ent",
		Short:             "ent 相关命令",
		GroupID:           g.ID,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	cmd.AddCommand()

	return g, cmd
}

func entNewCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:               "new",
		Short:             "新建 ent schema",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Args:              isUpperStartArgs,
		RunE:              func(_ *cobra.Command, _ []string) (err error) {},
	}
}

func entGenerateCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:               "generate",
		Short:             "生成 ent",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Args:              isUpperStartArgs,
		RunE:              func(_ *cobra.Command, _ []string) (err error) {},
	}
}
