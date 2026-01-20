// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-20, by liasica

package base

type CommonTemplateVariables struct {
	Module  string
	Year    string
	Date    string
	Package string
}

type DaoTemplateVariables struct {
	CommonTemplateVariables

	EntPkgImport string
	NameLower    string
	Name         string
	OrmClient    string
}
