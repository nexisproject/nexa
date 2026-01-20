// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-20, by liasica

package gen

import (
	"path/filepath"
	"strings"

	"nexis.run/nexa/cmd/nexa/internal/base"
)

// RenderDao 生成 Dao 代码
func RenderDao(cfg *base.Config, name string) (b []byte, err error) {
	var module string
	module, err = base.GetModule(cfg.RootDir)
	if err != nil {
		return
	}

	pkg := strings.ReplaceAll(strings.ToLower(filepath.Base(cfg.DaoPath)), "-", "_")

	vars := &base.DaoTemplateVariables{
		CommonTemplateVariables: base.CommonTemplateVariables{
			Module:  module,
			Year:    base.GetYear(),
			Date:    base.GetDate(),
			Package: pkg,
		},
		EntPkgImport: base.GetPkgImport(module, cfg.RootDir, cfg.EntPath),
		NameLower:    strings.ToLower(name),
		Name:         name,
		OrmClient:    cfg.OrmClient,
	}

	return RenderTemplate("dao.tmpl", vars)
}
