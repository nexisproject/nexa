// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-27, by liasica

package command

import (
	"github.com/spf13/cobra"

	"nexis.run/nexa/cmd/nexa/internal/entgen"
)

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

	cmd.AddCommand(
		entNewCmd(),
		entGenerateCmd(),
	)

	return g, cmd
}

func entNewCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:               "new",
		Short:             "新建 ent schema",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Args:              isUpperStartArgs,
		Example: examples(
			"nexa ent new Example",
		),
		RunE: func(_ *cobra.Command, names []string) error {
			g, err := entgen.New()
			if err != nil {
				return err
			}

			return g.New(names)
		},
	}

	return
}

func entGenerateCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:               "generate",
		Short:             "根据 ent schema 生成代码",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Example: examples(
			"nexa ent generate",
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			g, err := entgen.New()
			if err != nil {
				return err
			}

			return g.Generate()
		},
	}

	return
}
