// Package schemas embeds the kernel's JSON Schemas (draft 2020-12), one file per core type
// (spec/v2/03-architecture.md, "Core types"). The schema_version below is recorded in every run.
package schemas

import "embed"

// Version is the schema_version recorded in run.json.
const Version = "2.0.0-alpha.1"

// FS holds every <type>.schema.json, keyed by the lowercased type name.
//
//go:embed *.schema.json
var FS embed.FS
