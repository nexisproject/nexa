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

// DaoTemplateVariables 定义 DAO 模板变量
type DaoTemplateVariables struct {
	*CommonTemplateVariables

	EntPkgImport string
	NameLower    string
	Name         string
	OrmClient    string
}

// EchoCtxTemplateVariables 定义 echo Context 模板变量
type EchoCtxTemplateVariables struct {
	*CommonTemplateVariables

	Name string
}
