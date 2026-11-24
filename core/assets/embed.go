// Package assets embeds the default pet (spritesheet + pet.json).
// Only these two files are embedded; the other folders here are source/work files.
package assets

import _ "embed"

//go:embed spritesheet.png
var Spritesheet []byte

//go:embed pet.json
var PetJSON []byte
