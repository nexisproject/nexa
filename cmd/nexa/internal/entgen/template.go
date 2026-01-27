// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-27, by liasica

package entgen

import _ "embed"

var (
	//go:embed template/meta.tmpl
	TemplateMeta string

	//go:embed template/new_schema.tmpl
	TemplateNewSchema string

	//go:embed template/soft_delete.tmpl
	TemplateSoftDelete string

	//go:embed template/upsert.tmpl
	TemplateUpsert string
)
