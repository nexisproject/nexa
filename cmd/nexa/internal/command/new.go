// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-19, by liasica

package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"nexis.run/nexa/cmd/nexa/internal/base"
	"nexis.run/nexa/cmd/nexa/internal/gen"
)

func NewCmd() (*cobra.Group, *cobra.Command) {
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
		newDaoCmd(),
	)

	return g, cmd
}

func newDaoCmd() (cmd *cobra.Command) {
	var (
		force bool
	)

	cmd = &cobra.Command{
		Use:               "dao [names]",
		Short:             "新建数据访问对象模板",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Example: examples(
			"nexa new dao User",
			"nexa new dao User --force",
		),
		Args: func(_ *cobra.Command, names []string) error {
			for _, name := range names {
				if !base.StringIsUpperStart(name) {
					return base.ErrNameMustStartWithUpper
				}
			}
			return nil
		},
		Run: func(_ *cobra.Command, names []string) {
			cfg, err := base.GetConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var b []byte
			for _, name := range names {
				filename := filepath.Join(cfg.RootDir, cfg.DaoPath, strings.ToLower(name+".go"))

				b, err = gen.RenderDao(cfg, name)
				if err != nil {
					fmt.Printf("[DAO] %s 生成失败: %s\n", name, err.Error())
					os.Exit(1)
				}

				err = base.MkdirAll(filepath.Dir(filename))
				if err != nil {
					fmt.Printf("[DAO] %s 目录创建失败: %s\n", name, err.Error())
					os.Exit(1)
				}

				_, err = os.Stat(filename)
				if err == nil && !force {
					fmt.Printf("[DAO] %s 已存在，使用 --force 参数覆盖\n", name)
					continue
				}

				err = os.WriteFile(filename, b, 0644)
				if err != nil {
					fmt.Printf("[DAO] %s 文件写入失败: %s\n", name, err.Error())
					os.Exit(1)
				}

				fmt.Printf("[DAO] %s 创建成功: %s\n", name, filename)
			}
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "覆盖已存在的文件")

	return
}
