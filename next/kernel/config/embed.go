// Package defaults embeds config/defaults.toml so a single binary carries every threshold it uses.
package defaults

import _ "embed"

// TOML is the embedded defaults.toml, byte for byte. Its sha256 is the run's config_hash.
//
//go:embed defaults.toml
var TOML []byte
