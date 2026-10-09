// Package dcheck trägt die Dokumente, die das Werkzeug mitliefert
// (DC-FA-CLI-013): Benutzerhandbuch und Spezifikation, eingebettet zum
// Build-Stand und damit byte-gleich zu ihren Quelldateien. Das Paket liegt im
// Modul-Root, weil go:embed nur Dateien unter dem eigenen Verzeichnis erreicht.
package dcheck

import _ "embed"

// HandbuchPfad und SpezifikationPfad nennen die Quelldateien im Repo.
const (
	HandbuchPfad      = "docs/user/benutzerhandbuch.md"
	SpezifikationPfad = "spec/spezifikation.md"
)

// Handbuch ist das Benutzerhandbuch zum Build-Stand.
//
//go:embed docs/user/benutzerhandbuch.md
var Handbuch string

// Spezifikation ist die Spezifikation zum Build-Stand.
//
//go:embed spec/spezifikation.md
var Spezifikation string
