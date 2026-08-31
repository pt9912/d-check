package cli

// Generator für das include-bare d-check.mk (DC-FA-CLI-010, slice-038).
// Read-only: das Fragment wird auf stdout ausgegeben, nie geschrieben
// (DC-QA-03). Analog zum statischen --print-config-Gerüst.

import (
	"fmt"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// version ist die Release-Version. Beim Tag-Build via
// `-ldflags "-X .../cli.version=<tag>"` eingebettet (Dockerfile build-Stage,
// gespeist aus der Release-Pipeline); Default für lokale/Dev-/Gate-Builds.
// Quelle des Image-Refs in --print-mk — das Binary kennt seine Version,
// nicht seinen eigenen Digest (Henne-Ei; Digest via DCHECK_IMAGE-Override).
var version = "0.0.0-dev"

// makefileFragment erzeugt das d-check.mk: version-gepinnter, per
// DCHECK_IMAGE/DCHECK_DIGEST überschreibbarer Image-Ref plus dreizehn
// `##`-annotierte Targets (doc-check/doc-trace/doc-complete/doc-doctor/
// doc-repair/doc-immutable/doc-commits/doc-planning/doc-tracked/doc-targets/
// doc-structure/doc-usage/doc-help) und die TRACE_FLAGS-Variable. Das Template hat genau
// SIEBEN fmt-Verben — das %s der Version + je ein %s der vcs-/commits-/
// planning-/tracked-/targets-/structure-Disable-Flags; sonst KEIN '%' (sed
// statt awk-printf im doc-help-Recipe), sonst bräche fmt.Sprintf. Deterministisch (hängt nur an der eingebetteten Version +
// dem Modulsatz), read-only.
func makefileFragment() string {
	return fmt.Sprintf(mkTemplate, version, disableAllExcept("vcs"), disableAllExcept("commits"), disableAllExcept("planning"), disableAllExcept("tracked"), disableAllExcept("targets"), disableAllExcept("structure"))
}

// disableAllExcept liefert "--disable <m>"-Flags für alle Module außer keep,
// abgeleitet aus model.ValidModules (trackt den Modulsatz automatisch). So läuft
// ein fokussiertes Target NUR das Modul keep — sonst über-feuerten die Doc-Module
// des Konsumenten auf Nicht-Ziel-Befunde des Arbeitsbaums (vgl. das fokussierte
// `make adr-check`/`make trace-check` von d-check selbst).
func disableAllExcept(keep string) string {
	var flags []string
	for _, m := range model.ValidModules() {
		if m != keep {
			flags = append(flags, "--disable "+m)
		}
	}
	return strings.Join(flags, " ")
}

const mkTemplate = "# d-check.mk — erzeugt von: d-check --print-mk (DC-FA-CLI-010).\n" +
	"#\n" +
	"# Einbinden: \"include d-check.mk\" im eigenen Makefile; eine eigene\n" +
	"# .d-check.yml danebenlegen. Keine Recipe-/Skript-Kopie — der Image-Pin\n" +
	"# lebt in d-check.\n" +
	"#\n" +
	"# Für strikte Reproduzierbarkeit den Digest aus den Release-Notes pinnen —\n" +
	"# direkt über DCHECK_IMAGE oder bequemer über DCHECK_DIGEST (sticht den Tag):\n" +
	"#   DCHECK_DIGEST = sha256:<digest>\n" +
	"#\n" +
	// Der Zeiger steht im KOPF, weil das Fragment in ein fremdes Repo reist:
	// dort ist die Einleitung der einzige Ort, an dem er mitfaehrt.
	"# Benutzerhandbuch (aufgabenorientiert, deutsch):\n" +
	"#   " + handbuchURL + "\n" +
	"DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v%s\n" +
	"DCHECK_DIGEST ?=\n" +
	"# TRACE_FLAGS: optionale Flags für die RTM-Targets (z. B. --json).\n" +
	"TRACE_FLAGS ?=\n" +
	"\n" +
	"# Ein gesetzter DCHECK_DIGEST sticht den Tag von DCHECK_IMAGE.\n" +
	"ifeq ($(strip $(DCHECK_DIGEST)),)\n" +
	"DCHECK_REF := $(DCHECK_IMAGE)\n" +
	"else\n" +
	"DCHECK_REF := ghcr.io/pt9912/d-check@$(DCHECK_DIGEST)\n" +
	"endif\n" +
	"\n" +
	".PHONY: doc-check\n" +
	"doc-check: ## Doku-Referenzen prüfen (Befund-Gate)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF)\n" +
	"\n" +
	".PHONY: doc-trace\n" +
	"doc-trace: ## Requirements Traceability Matrix auf stdout (advisory, DC-FA-CLI-009)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --trace $(TRACE_FLAGS)\n" +
	"\n" +
	".PHONY: doc-complete\n" +
	"doc-complete: ## Vollständigkeits-Gate: Requirements-Waise ⇒ Exit 1 (DC-FA-CLI-011)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --trace --require-complete $(TRACE_FLAGS)\n" +
	"\n" +
	".PHONY: doc-doctor\n" +
	"doc-doctor: ## erklärende Diagnose mit Fix-Kandidaten (DC-FA-CLI-007)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --doctor\n" +
	"\n" +
	".PHONY: doc-repair\n" +
	"doc-repair: ## Reparatur-Patch (unified diff) auf stdout, git-apply-rein (DC-FA-CLI-008)\n" +
	"\t@docker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --repair\n" +
	"\n" +
	".PHONY: doc-immutable\n" +
	"doc-immutable: ## Doc-/ADR-Immutabilität via git-Diff (Modul vcs); RANGE=base..head oder STAGED=1 (DC-FA-VCS-001)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --enable vcs %s $(if $(STAGED),--staged,--range $(RANGE))\n" +
	"\n" +
	".PHONY: doc-commits\n" +
	"doc-commits: ## Commit-Message-Traceability via Modul commits; RANGE=base..head (DC-FA-COMMITS-001)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --enable commits %s --range $(RANGE)\n" +
	"\n" +
	".PHONY: doc-planning\n" +
	"doc-planning: ## Planning-Lifecycle-Konsistenz (Roadmap <-> in-progress) via Modul planning; hermetisch, ohne Range (DC-FA-PLAN-001)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --enable planning %s\n" +
	"\n" +
	".PHONY: doc-tracked\n" +
	"doc-tracked: ## Getrackt-Status aufloesbarer Referenz-Ziele via Modul tracked; braucht .git im Mount, ohne Range (DC-FA-TRK-001)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --enable tracked %s\n" +
	"\n" +
	".PHONY: doc-targets\n" +
	"doc-targets: ## Deklarations-Konsistenz Doku<->Build-Targets via Modul targets; hermetisch, ohne Range (DC-FA-TGT-001)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --enable targets %s\n" +
	"\n" +
	".PHONY: doc-structure\n" +
	"doc-structure: ## Struktur-Invarianten innerhalb der Dokumente via Modul structure; hermetisch, ohne Range (DC-FA-STRUCT-001)\n" +
	"\tdocker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --enable structure %s\n" +
	"\n" +
	// --help schreibt auf stderr und endet mit Exit 0 — eine Umleitung von
	// stdout fängt nichts (gemessen: 0 Byte stdout, 2488 Byte stderr). Das
	// Recipe-Echo ist unterdrückt, weil die Ausgabe die Nutzlast ist, nicht
	// ein Befund — dieselbe Klasse wie doc-repair und doc-help.
	".PHONY: doc-usage\n" +
	"doc-usage: ## Aufruf und Optionen von d-check selbst (--help)\n" +
	"\t@docker run --rm --network none -v \"$(CURDIR):/repo:ro\" $(DCHECK_REF) --help\n" +
	"\n" +
	".PHONY: doc-help\n" +
	"doc-help: ## diese Liste der doc-*-Targets\n" +
	"\t@grep -hE '^doc-[a-z-]+:.*## ' $(MAKEFILE_LIST) | sort | sed -E 's/:.*## /  /'\n"
