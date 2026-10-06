# Eingehender Change Request — `targets.makefiles` nimmt Glob-Muster an

**Absender:** Adopter `ai-harness-init` · **Eingegangen:** 2026-10-06
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Modul `targets`, Schlüssel `targets.makefiles`)
**Stand:** **eingegangen**; Entscheidung und Umsetzung trägt `slice-250`.

**Ablage-Hinweis.** Ein **eingehender** CR ist die dritte Klasse neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf). Die
Datei liegt hier aus demselben Grund wie ihre Vorgänger in diesem
Verzeichnis: Was gebeten wurde und wie entschieden wird, soll im Repo
stehen und nicht nur im Vorgang. Im Wortlaut sind Kennungen verlinkt
(Linkpflicht des eigenen Gates), sonst unverändert.

---

## Wortlaut

> **CR aus ai-harness-init: targets.makefiles nimmt Glob-Muster an**
>
> Absender: ai-harness-init (Adopter von d-check, emittiert d-check in
> gebootstrappte Ziel-Repos). Gegenstand: [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in), Modul targets,
> Schlüssel targets.makefiles. Gemessen am Stand: v0.79.0 (gepinnt beim
> Absender), Quelltext gegengelesen bis v0.80.0-1.
>
> ### Anlass
>
> Ein von ai-harness-init gebootstrapptes Ziel-Repo verteilt seine
> make-Targets auf mehrere Dateien:
>
> - das root-Makefile mit drei eigenen Targets (help, gates, record-gates),
> - elf Fragmente unter `harness/mk/*.mk` mit je 2 bis 6 Targets, eingebunden
>   per `include harness/mk/*.mk`,
> - `d-check.mk` mit 26 Targets.
>
> Die Menge der Fragmente ist nicht fest. Sie wächst mit `add-lang <sprache>`
> und mit eigenen Fragmenten des Adopters. Ein Glob ist die einzige Angabe,
> die mit dem Bestand mitwächst.
>
> ### Ist-Verhalten
>
> collectMakefileRules (`internal/hexagon/core/rules/targets.go`, Zeile 67
> ff.) liest jeden Eintrag von makefiles wörtlich per fsys.ReadFile(mf). Ein
> Eintrag `harness/mk/*.mk` wird als Dateiname gelesen, scheitert und bricht
> fail-closed ab:
>
> das Modul targets kann das Makefile "harness/mk/*.mk" nicht lesen
> ([`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in), fail-closed): …
>
> applyTargets (configyaml.go) prüft nur, dass der Pfad relativ ist (kein /,
> kein ..). Ein Glob passiert die Validierung und scheitert erst im Lauf.
>
> Folge beim Adopter: Wer alle Fragmente prüfen will, muss jede Datei einzeln
> listen und die Liste bei jedem neuen Fragment nachziehen. Eine vergessene
> Datei macht deren Targets still ungeprüft. Das ist die Lücke, die ein Glob
> schließt.
>
> ### Bitte
>
> targets.makefiles nimmt neben wörtlichen Pfaden auch Glob-Muster an, mit
> derselben Glob-Semantik, die d-check bereits an anderen Schlüsseln führt
> (z. B. structure[].files, die Closure- und Mentions-Globs).
>
> ### Akzeptanzkriterien
>
> 1. Expansion: Ein Eintrag mit Glob-Zeichen expandiert gegen die
>    Repo-Wurzel. Jede Treffer-Datei wird wie ein wörtlicher Eintrag gelesen,
>    Fundstellen nennen den echten Dateipfad.
> 2. Wörtliche Pfade unverändert: Ein Eintrag ohne Glob-Zeichen verhält sich
>    wie heute, eine fehlende Datei bricht fail-closed ab.
> 3. Leerer Glob: Ein Glob ohne Treffer ist ein Befund oder ein Abbruch, nicht
>    still. Welches von beiden, entscheidet d-check. Bitte nicht still inert,
>    denn dann prüft das Modul unbemerkt nichts.
> 4. Doppelte Treffer: Eine Datei, die wörtlich und per Glob erfasst wird,
>    zählt ihre Regeln einmal. Kein doppelter Befund, kein falsches
>    gate-phantom.
> 5. Validierung: Glob-Einträge unterliegen derselben Pfad-Regel wie heute
>    (relativ, kein /, kein ..). Ein Muster, das aus der Wurzel herausführen
>    würde, wird beim Laden abgewiesen.
> 6. Rotes Gegenbeispiel: Eine Regel in einer per Glob erfassten Datei ohne
>    Doku-Zeile meldet gate-undocumented mit Fundstelle in dieser Datei.
>
> ### Abgrenzung
>
> - doc-tables und authority: Glob ist dort nicht Gegenstand dieses Antrags.
> - Regel-Extraktion: Die Regel-Erkennung (makefileRuleRe) und die Semantik
>   von exempt-targets bleiben unverändert.
> - include-Auflösung: Make-include-Direktiven auswerten ist ausdrücklich
>   nicht gewünscht. Der Glob in der Konfiguration ist die explizite,
>   prüfbare Angabe.
>
> ### Offene Fragen an d-check
>
> 1. Leerer Glob: Befund (mit welchem Grund-Code) oder Abbruch mit Exit 2?
> 2. Soll --print-config bzw. die Konfig-Vorlage (config_template.go:254,
>    heute makefiles: [Makefile]) ein Glob-Beispiel tragen?
