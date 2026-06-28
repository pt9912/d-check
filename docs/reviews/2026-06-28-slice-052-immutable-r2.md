# Review — slice-052 (Modul `immutable`) · R2 (Doc-first-Kohärenz / Harness-Konformität / Ehrlichkeit)

## Kopf-Metadaten

- **Gegenstand:** Commit `b62a520` — `feat(immutable): Modul immutable — Immutabilitäts-Pin gegen Core-Drift (DC-FA-IMM-001, ADR-0023, slice-052)`
- **Review-Lauf:** R2 (unabhängiger Reviewer). Schwerpunkt laut Auftrag: **Doc-first-Kohärenz, Harness-Prozess-Konformität, Ehrlichkeit der Entscheidung** — **nicht** Code-Logik (separater Code-Reviewer).
- **Datum:** 2026-06-28
- **Reviewer-Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) v1.2.0 (Output-Schema, Kategorien-Anker, Negativbefund-Pflicht, MEDIUM-Anker „Referenz-Richtung (SDP) — Marker-Ehrlichkeit").
- **Eingangs-Kontext:** Diff `b62a520`; Slice [`slice-052`](../plan/planning/done/slice-052-immutable-modul.md); Anforderung `DC-FA-IMM-001`; [ADR-0023](../plan/adr/0023-immutable-core-pin.md) (Proposed); Schwester-Gate [ADR-0016](../plan/adr/0016-adr-immutable-gate.md)/[`tools/adr-immutable-check.sh`](../../tools/adr-immutable-check.sh); Grenzen [ADR-0008](../plan/adr/0008-reparatur-ableitbarkeit.md)/[ADR-0020](../plan/adr/0020-content-pin-fence-ausnahme.md); Hard Rules [`AGENTS.md`](../../AGENTS.md) §3.
- **Rollen-Abgrenzung:** Kein Verifier (Gate-Läufe nicht meine Rolle — `make gates` + doc-check liefen laut Auftrag grün); kein Stil-Polizist; kein Finding ohne Failure-Szenario; REFUTED nur mit Zitat.
- **Repo-Stand bei Review:** `HEAD == b62a520`; kein `v0.32.0`-Tag; jüngster Release `v0.31.0` (elf Module).

---

## Findings

### F-1 — README behauptet Release-/GHCR-Zustand vor dem Release (MEDIUM)

- **kategorie:** MEDIUM
- **quelle:** Source Precedence ([`AGENTS.md` §2](../../AGENTS.md#2-kanonische-quellen-source-precedence): Roadmap Rang 5 > README Rang 7) + Harness-Ehrlichkeit („kein Erfolg ohne echte Gate-/Release-Ausgabe"; Maintainability)
- **pfad:** `README.md:6-11` (insb. „sind im GHCR-Image" Z. 8 und der CHANGELOG-Zeiger Z. 9-11)
- **befund:** Der Feat-Commit setzt die Statuszeile auf „**Status: released** — alle **zwölf** Regelmodule (…, `immutable`, …) **sind im GHCR-Image**" und „die jeweils jüngsten Änderungen (zuletzt das opt-in-Modul `immutable`) führt die [CHANGELOG.md]". Am Review-Commit ist beides unwahr: es existiert kein `v0.32.0`-Tag (`HEAD == b62a520`, jüngster Release `v0.31.0` = elf Module), und `CHANGELOG.md` trägt als obersten Eintrag `[0.31.0]` ohne `immutable`/`0.32.0`. Die höherrangige Roadmap (§Aktuelle Welle) sagt für genau diesen Stand „Code, zwei Reviews und Release v0.32.0 **folgen**". Die Modul-Zähler-Zeile wurde bisher konventionell **erst im Release-Prep-Commit** gebumpt (zehn→elf in `f1f2117` „docs(release): Release-Prep v0.29.0 … README"), nicht im Feat-Commit; `b62a520` zieht die Release-Behauptung in den Feat-Commit vor.
- **Failure-Szenario:** Ein Nutzer liest README an diesem Stand, zieht das jüngste GHCR-Image (`v0.31.0`) und findet `immutable` **nicht** — die Behauptung „alle zwölf … sind im GHCR-Image" trifft erst nach dem noch ausstehenden `v0.32.0`-Push zu; ebenso verweist der Text auf einen CHANGELOG-Eintrag, der noch fehlt.
- **verifizierbar:** ja — `git tag | grep 0.32.0` (leer) + Modul-Liste des publizierten Images zeigt elf Module; `grep -nE '^## ' CHANGELOG.md | head -1` zeigt `[0.31.0]`. Kein Gate fängt es (doc-check prüft Links/Anker/Token, nicht Release-Wahrheit).

### F-2 — ADR-0023 stellt den Backend-Trade einseitig als „nur schwächer" dar (INFO)

- **kategorie:** INFO
- **quelle:** ADR-0023 (Ehrlichkeit der Entscheidung) / `DC-FA-IMM-001`
- **pfad:** `docs/plan/adr/0023-immutable-core-pin.md:35-41` (Absatz „Zwei Backends, bewusst koexistent") + Alternativen-Tabellenzeile „Content-Pin `immutable` (gewählt)"
- **befund:** ADR-0023 charakterisiert den Pin-Backend gegenüber `adr-check` ausschließlich als „**schwächer** … neu-pinn-bar … Reviewer als Boden". Unerwähnt bleibt die **Gegen-Friktion**: die `**Status:**`-Zeile liegt immer im gehashten Core und ist per `immutable.exclude-sections` (nimmt **Abschnitte**, keine **Zeilen** aus) nicht ausnehmbar — ein vom Schwester-Skript still erlaubter Status-Übergang (Supersede) löst beim Pin-Backend `core-drift` aus und verlangt **Neu-Pinnen**. Diese Stelle, an der der Pin-Backend **strenger/lauter** ist als das Skript, ist nur in [`slice-052` §4](../plan/planning/done/slice-052-immutable-modul.md) und im Lastenheft-Out-of-Scope (`DC-FA-IMM-001`) offengelegt, nicht in der ADR-Trade-Darstellung.
- **Failure-Szenario (mild, da im Artefakt-Set offengelegt):** Wer ADR-0023 isoliert liest, unterschätzt die Adoptions-Friktion an `Accepted`-ADRs (Supersede ⇒ Neu-Pinnen). Kein Widerspruch zum bindenden Lastenheft — dort steht es im Out-of-Scope; reine Schwerpunkt-Asymmetrie in der ADR.
- **verifizierbar:** nein (Prosa-Kohärenz; kein Gate). Nicht-blockierend; als INFO geführt, weil das Artefakt-Set (Lastenheft + Slice) die Grenze ehrlich trägt.

---

## Negativbefunde (geprüft, ohne blockierenden Befund)

- **Versions-/ID-Kohärenz:** Header `**Version:** 0.32.0` (`spec/lastenheft.md:3`) deckt sich mit der §7-Historie-Zeile `0.32.0` und mit Roadmap/Slice-DoD („Lastenheft 0.32.0", „Versions-Bump 0.32.0"). Die sieben Verweise auf die neue Anforderung tragen **byte-identisch** den Anker `#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in` (ADR-0023 Bezug, ADR-Index, Spezifikation ×2, Roadmap, README, Slice); Umlaut-Slug `immutabilitäts` analog zum bestehenden `dc-fa-pin-001…`-Muster, doc-check-Anker grün. Die `…-001.a`-Spezifikations-Sektion + die ADR-`Schärft:`-Zielanker (`#dc-fa-imm-001a--…-immutable`) korrespondieren.
- **Ein-Begriff-Ehrlichkeit:** Grund-Code **`core-drift`** identisch über Lastenheft, Spezifikation (§4-Tabelle + `.a`), ADR-0023, ADR-Index, Slice, Roadmap, README, Handbuch §6. Marker **`<!-- immutable: sha256:… -->`** konsistent (README/Handbuch nutzen die Kurzform `<!-- immutable: … -->` analog `pins`). Schema-Key **`immutable.exclude-sections`** identisch in Spezifikation §2, Lastenheft, ADR, Slice.
- **Zwei-Backend-Story (Marker-Ehrlichkeit-Geist):** Die schwächere Garantie ist durchgängig benannt — „neu-pinn-bar", „Reviewer als Boden", „dieselbe Disziplin wie `pins`/`versions`" (ADR-0023, Lastenheft-Zweck, Slice §2/§4). Keine Stelle suggeriert eine **härtere** Garantie als das Modul liefert. `adr-check` wird durchgängig als **„unangetastet"/„koexistent"** geführt (ADR-0023 „bleibt unangetastet", Slice §2 „unangetastet", §Bezug) — **nirgends** „ersetzt". Die git-Form ist konsistent als „vertagt/Out-of-Scope/eigener nicht-hermetischer Port/eigene Anforderung" benannt (ADR-0023, Lastenheft Out-of-Scope, Spezifikation `.a`, Slice §2). ADR-0023-Beschreibung des Skript-Cores („Datei ohne `## Geschichte` und ohne die Status-Zeile") deckt sich mit [`tools/adr-immutable-check.sh`](../../tools/adr-immutable-check.sh) (`core()` + Status-Zeilen-Strip).
- **Status-Zeile/Supersede-Grenze:** Ehrlich und vollständig in [`slice-052` §4](../plan/planning/done/slice-052-immutable-modul.md) („die ausgeschlossene Status-Zeile … ist eine Zeile, kein Abschnitt … Supersede verlangt Neu-Pinnen; die feinere ‚nur-Status-Zeile-strippen'-Semantik des Skripts ist mögliche Folge-CR") und im Lastenheft-Out-of-Scope. Die Spec verspricht **nicht** die Skript-Semantik „nur Geschichte-Anhang + Status-Übergang erlaubt" für das Modul. (Asymmetrische Platzierung in ADR-0023 → F-2.)
- **Source Precedence / Referenz-Richtung (SDP):** Die Spec-Straten referenzieren **abwärts nichts** — der Lastenheft-Out-of-Scope und die Spezifikations-`.a` sagen bewusst „in **eigener/begleitender ADR** festgehalten" **ohne** ADR-Nummer/-Link; keine ADR-/Wellen-/Slice-Token in Lastenheft- oder Spezifikations-Prosa. ADR-0023 nennt `slice-052` **nur** in `## Geschichte` (ausgenommen); sein `Schärft:` zeigt **aufwärts** auf die Spezifikation, die Spezifikation verlinkt **nicht** zurück. Die §7-Verweis-Spalte `| slice-052 |` setzt die bestehende Historie-Tabellen-Konvention fort (wie 0.31.0/0.30.0/0.29.0); kein **neuer** Abwärts-Verweis in Prosa. Keine `<!-- d-check:status-provenance -->`-Marker neu gesetzt. doc-check (`matrix`/`ids`) grün.
- **Out-of-Scope-Vollständigkeit:** `DC-FA-IMM-001` Out-of-Scope deckt sauber ab: (a) VCS-/git-historienbasierte Form (`core(BASE)` vs. `core(HEAD)`, eigener nicht-hermetischer Port), (b) Pinnen/Neu-Pinnen durch das Werkzeug (`--bless` als eigene Anforderung, berührt `DC-FA-CLI-008`), (c) Status-Übergangs-Semantik, (d) Default-on/Pflicht-Pin, (e) mehrere Hash-Algorithmen.
- **ADR-Index & Status:** ADR-Index-Zeile (`docs/plan/adr/README.md`) entspricht ADR-0023 (Titel, **Proposed**, 2026-06-28, Bezüge `DC-FA-IMM-001`/0020/0016/0008). ADR-0023 ist `Proposed` ⇒ ADR-Immutable-Gate greift noch nicht; konsistent.
- **Roadmap-Flip / Planning-Konformität:** §Aktuelle Welle „welle-41-immutable — aktiv" mit `slice-052` in `in-progress/` ist konsistent mit `make planning-check` (aktive Welle ⇔ Slice in `in-progress/`); „Zuletzt abgeschlossen welle-40" korrekt nach `done/` verlinkt. Roadmap (Planning-Doc) darf ADR/Slice referenzieren (kein Spec-Stratum) — `DC-FA-IMM-001`/ADR-0023-Links zulässig.
- **Modul-Zähl-Konsistenz (Code/Spec):** `DC-FA-CLI-002`-Liste, Glossar (Lastenheft) und `model.validModules()` führen `immutable` als 12. Modul; keine verbliebene „elf"-Prosa außerhalb der Versionshistorie (`grep` leer). Default-aus überall benannt (opt-in pro Datei + striktes opt-in Modul).

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 (`README.md:6-11`) |
| LOW | 0 | — |
| INFO | 1 | F-2 (`docs/plan/adr/0023-immutable-core-pin.md:35-41`) |

---

## Verdikt

**Bedingt mergebar — 1 MEDIUM (F-1) vor Merge zu klären.**

Die Doc-first-Substanz ist kohärent und ehrlich: Versions-, Anker-, Grund-Code-,
Marker- und Schema-Konsistenz tragen durchgängig; die Zwei-Backend-Story benennt
die **schwächere** Pin-Garantie offen, stellt `adr-check` als **unangetastet/koexistent**
dar (nie „ersetzt") und vertagt die git-Form sauber als eigene Anforderung; die
Spec-Straten halten die Referenz-Richtung (SDP) ein (keine Abwärts-ADR/-Slice-Token
in Prosa); Out-of-Scope deckt git-Form, Neu-Pinnen und Status-Semantik ab.

**F-1** ist als MEDIUM (statt LOW-„Doku-Drift") eingestuft, weil es nicht bloß eine
Modul-Listen-Drift ist, sondern eine **Release-Zustands-Behauptung vor dem Release**
(„released … sind im GHCR-Image" + CHANGELOG-Zeiger), die der **höherrangigen**
Roadmap („Release v0.32.0 folgen") und der gelebten Release-Prep-Konvention
(Zähler-/CHANGELOG-Bump erst im Release-Prep-Commit, vgl. `f1f2117`) widerspricht —
der Harness-Wert „grün ist der Boden, nicht die Decke" gilt für die README-Behauptung
ebenso. Die Auflösung (Zähler-/Release-Zeile in den Release-Prep-Commit verschieben
**oder** den Release nachziehen) liegt bei der Implementation, nicht im Finding.

**F-2** ist nicht-blockierend (INFO): die Status-Übergangs-Friktion ist im bindenden
Lastenheft + Slice offengelegt; lediglich die ADR-Trade-Darstellung ist einseitig.
