// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-19, by liasica

package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"nexis.run/nexa/cmd/nexa/internal/base"
	"nexis.run/nexa/cmd/nexa/internal/gen"
)

func NewCmd() (*cobra.Group, *cobra.Command) {
	var (
		force bool
	)

	g := &cobra.Group{
		ID:    "new",
		Title: "新建代码命令",
	}

	cmd := &cobra.Command{
		Use:               "new",
		Short:             "新建代码模板",
		GroupID:           g.ID,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	cmd.AddCommand(
		newDaoCmd(force),
		newEchoctxCmd(force),
	)

	cmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "覆盖已存在的文件")

	return g, cmd
}

func newDaoCmd(force bool) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:               "dao [names]",
		Short:             "新建数据访问对象模板",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Example: examples(
			"nexa new dao User",
			"nexa new dao User --force",
		),
		Args: isUpperStartArgs,
		RunE: func(_ *cobra.Command, names []string) (err error) {
			var g *gen.Gen
			g, err = gen.New()
			if err != nil {
				return
			}

			for _, name := range names {
				err = g.Generate(gen.PackageDao, name, force, func(g *gen.Gen, c *base.CommonTemplateVariables) any {
					return &base.DaoTemplateVariables{
						CommonTemplateVariables: c,
						EntPkgImport:            base.GetPkgImport(g.Module, g.Config.RootDir, g.Config.EntPath),
						NameLower:               strings.ToLower(name),
						Name:                    name,
						OrmClient:               g.Config.OrmClient,
					}
				})
				if err != nil {
					return
				}

				fmt.Printf("[DAO] %s 创建成功\n", name)
			}

			return
		},
	}

	return
}

func newEchoctxCmd(force bool) (cmd *cobra.Command) {
	return &cobra.Command{
		Use:               "echoctx [names]",
		Short:             "新建数据访问对象模板",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Example: examples(
			"nexa new echoctx Rider",
			"nexa new echoctx Rider --force",
		),
		Args: isUpperStartArgs,
		RunE: func(_ *cobra.Command, names []string) (err error) {
			var g *gen.Gen
			g, err = gen.New()
			if err != nil {
				return
			}

			for _, name := range names {
				err = g.Generate(gen.PackageEchoctx, name, force, func(g *gen.Gen, c *base.CommonTemplateVariables) any {
					return &base.EchoCtxTemplateVariables{
						CommonTemplateVariables: c,
						Name:                    name,
					}
				})
				if err != nil {
					return
				}

				fmt.Printf("[Echo Context] %s 创建成功\n", name)
			}

			return
		},
	}
}
