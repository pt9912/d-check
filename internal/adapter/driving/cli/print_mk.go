package cli

// Generator für das include-bare d-check.mk (DC-FA-CLI-010, slice-038).
// Read-only: das Fragment wird auf stdout ausgegeben, nie geschrieben
// (DC-QA-03). Analog zum statischen --print-config-Gerüst.

import "fmt"

// version ist die Release-Version. Beim Tag-Build via
// `-ldflags "-X .../cli.version=<tag>"` eingebettet (Dockerfile build-Stage,
// gespeist aus der Release-Pipeline); Default für lokale/Dev-/Gate-Builds.
// Quelle des Image-Refs in --print-mk — das Binary kennt seine Version,
// nicht seinen eigenen Digest (Henne-Ei; Digest via DCHECK_IMAGE-Override).
var version = "0.0.0-dev"

// makefileFragment erzeugt das d-check.mk: version-gepinnter, per
// DCHECK_IMAGE überschreibbarer Image-Ref plus ein `doc-check`-Target.
// Deterministisch (hängt nur an der eingebetteten Version), read-only.
func makefileFragment() string { return fmt.Sprintf(mkTemplate, version) }

const mkTemplate = "# d-check.mk — erzeugt von: d-check --print-mk (DC-FA-CLI-010).\n" +
	"#\n" +
	"# Einbinden: \"include d-check.mk\" im eigenen Makefile; eine eigene\n" +
	"# .d-check.yml danebenlegen. Keine Recipe-/Skript-Kopie — der Image-Pin\n" +
	"# lebt in d-check.\n" +
	"#\n" +
	"# DCHECK_IMAGE ist überschreibbar. Für strikte Reproduzierbarkeit den\n" +
	"# Digest aus den Release-Notes pinnen:\n" +
	"#   DCHECK_IMAGE = ghcr.io/pt9912/d-check@sha256:<digest>\n" +
	"DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v%s\n" +
	"\n" +
	".PHONY: doc-check\n" +
	"doc-check:\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_IMAGE)\n"
