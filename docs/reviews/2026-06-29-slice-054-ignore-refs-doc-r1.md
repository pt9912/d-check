# Review — slice-054 (`codepaths.ignore-refs`) · R1 (Doc-first-Kohärenz / Harness-Konformität / Ehrlichkeit)

## Kopf-Metadaten

- **Review-Lauf:** R1 (unabhängiger Reviewer). Schwerpunkt: Doc-first-Kohärenz,
  Harness-Prozess-Konformität, Ehrlichkeit der Entscheidung — **nicht** Code-Logik
  (separater Code-Reviewer R2).
- **Datum:** 2026-06-29
- **Reviewer-Skill:** `.harness/skills/reviewer.md` v1.2.0 (Output-Schema, Kategorien-Anker,
  kein Finding ohne Failure-Szenario, kein Lösungsvorschlag im Befund, REFUTED nur mit
  Zitat, Negativbefund-Pflicht).
- **Gegenstand:** Working-Tree-Änderungen für slice-054 (vor Commit), 18 M / 1 D / 2 ??.
- **Eingangs-Kontext:** Slice slice-054; Anforderung DC-FA-CODE-001 (`spec/lastenheft.md`,
  `spec/spezifikation.md` §DC-FA-CODE-001.a + §2); ADRs 0025 (Proposed), 0024 (Accepted),
  0016; Hard Rules `AGENTS.md` §3.5.
- **Rollen-Abgrenzung:** Kein Verifier (Gate-Läufe nicht meine Rolle — `make ci` +
  `completeness-check` liefen laut Auftrag grün); kein Stil-Polizist; REFUTED nur mit Zitat.

## Findings

### MEDIUM-1 — Index annotiert ADR-0024 als (teil-)superseded, während ADR-0025 noch `Proposed` ist

- **Kategorie:** MEDIUM
- **Quelle:** AGENTS §3.5 / slice-054 §3d / slice-053-R1-F-4-Klasse / Ehrlichkeit der Entscheidung
- **Pfad:** `docs/plan/adr/README.md` (ADR-0024-Zeile) vs. slice-054 §3d
- **Befund:** Die ADR-0024-Index-Zeile trug bereits die Status-Annotation „Accepted
  (Skript-Stabilitäts-Teilentscheidung superseded durch ADR-0025)", obwohl ADR-0025
  `Proposed` ist und der Slice-Plan §3d genau diese Annotation („Geschichte + Index, bei
  Closure — wie slice-053-R1-F-4") als Closure-Schritt terminiert. Die Annotation behauptet
  einen ratifizierten Supersede, der noch nicht stattgefunden hat — exakt die
  slice-053-R1-F-4-Klasse, die dieser Slice selbst vermeiden will. **Failure-Szenario:**
  Wird ADR-0025 nicht angenommen/umgearbeitet, verzeichnet der Index ADR-0024 als von einer
  nie akzeptierten ADR superseded; ein Auditor liest mitten im Slice einen vollzogenen
  Supersede statt eines Vorschlags. (Die ADR-0024-/ADR-0016-Bodies sind korrekt unangetastet
  — §3.5 ist insoweit gewahrt; der Verstoß betrifft nur die vorgezogene Index-Annotation.)
- **Verifizierbar:** nein (kein Gate; Plan-vs-Working-Tree-Diff manuell).

### LOW-1 — §2-Schema und Ventil-Prosa nennen nur „existenz-/anker", unterschlagen die `repo-escape`-Unterdrückung

- **Kategorie:** LOW
- **Quelle:** DC-FA-CODE-001 (Spec-Treue der Messmethode)
- **Pfad:** `spec/spezifikation.md` §2-Schema `codepaths.ignore-refs` + Referenz-Ventil-Absatz
  vs. Schritt 5 und `internal/hexagon/core/rules/codepaths.go`
- **Befund:** §2-Schema und Referenz-Ventil-Absatz zählen nur Existenz und Anker auf;
  Schritt 5 ordnet das `ignore-refs`-Match jedoch **vor** den Escape-Zweig, ADR-0025 sagt
  „`anchor-missing`/`repo-escape` … entfallen", der Code prüft `ignored(rel, …)` vor
  `if escaped`. Damit unterdrückt das Ventil auch `repo-escape`, was die zwei
  zusammenfassenden Stellen nicht erwähnen. **Failure-Szenario:** Ein Auditor liest die
  kanonische §2-Schema-Zeile und schließt, `repo-escape`-Schutz bleibe für
  `ignore-refs`-Pfade erhalten; tatsächlich wird auch ein die Wurzel verlassender
  Tombstone-Pfad stillgestellt.
- **Verifizierbar:** nein (Spec-internes Lese-Delta).

### INFO-1 — `ignore-refs` deckt nur die `codepaths`-Achse; ADR-0025s „generisch aufgelöst" benennt die `links`-Rest-Falle nicht

- **Kategorie:** INFO
- **Quelle:** ADR-0025 / Ehrlichkeit / Konsistenz `codepaths`↔`links`
- **Pfad:** `docs/plan/adr/0025-codepaths-ignore-refs.md` (Konsequenz „generisch aufgelöst")
  + Beleg `docs/reviews/2026-06-28-slice-052-immutable-r2.md`
- **Befund:** Das Ventil wirkt nur im Modul `codepaths` (Inline-Code-Pfade). ADR-0025
  formulierte die Auflösung als „generisch", ohne die `links`-Achse auszunehmen. Im selben
  Diff musste ein Markdown-**Link** auf das gelöschte Skript im Review-Record zu Inline-Code
  entlinkt werden, weil `links` (`target-missing`) ihn sonst meldet. **Failure-Szenario
  (latent):** Zitiert eine **immutable** Doku einen entfernten Pfad als Markdown-Link statt
  Inline-Code, feuert `links target-missing` an einer uneditierbaren Datei weiter;
  `ignore-refs` kann das nicht decken. Aktuell rein latent (kein überlebendes Link-Ziel).
- **Verifizierbar:** ja (konstruierter Negativtest: immutable Doc mit Markdown-Link auf einen
  `ignore-refs`-Pfad → `links target-missing` bleibt).

### INFO-2 — Rückwirkende Bearbeitung eines historischen Review-Records, um `links` grün zu halten

- **Kategorie:** INFO
- **Quelle:** Reviewer-Skill (Ablage „Nie überschreiben") / Record-Integrität
- **Pfad:** `docs/reviews/2026-06-28-slice-052-immutable-r2.md`
- **Befund:** Zwei Markdown-Links auf `tools/adr-immutable-check.sh` im abgeschlossenen
  R2-Record vom 2026-06-28 wurden zu reinem Inline-Code entlinkt (das `docs/reviews/**`-
  `exempt-paths` deckt `codepaths` ab, nicht `links`). Inhaltlich bleibt die Aussage erhalten,
  aber ein datiertes Review-Artefakt wird nachträglich verändert, um ein Gate eines
  Folge-Slices grün zu halten. **Failure-Szenario:** Wer den Record als unveränderten
  Zeitstempel-Beleg heranzieht, findet eine still nachgezogene Form.
- **Verifizierbar:** nein (git-Diff-Beobachtung).

### INFO-3 — ADR-0016s immutable Geschichte behauptet weiter „das Skript bleibt pfad-stabiler Fallback im Baum"

- **Kategorie:** INFO
- **Quelle:** Traceability / Doc-Drift
- **Pfad:** `docs/plan/adr/0016-adr-immutable-gate.md` (Geschichte-Zeile 2026-06-29) vs. slice-054 §3d
- **Befund:** ADR-0016s Geschichte-Zeile sagt „das Skript bleibt pfad-stabiler Fallback im
  Baum" — nach dieser Slice falsch. ADR-0016 ist immutable, die bestehende Zeile darf nicht
  editiert werden (korrekt unangetastet); §3d plante aber nur eine Annotation an ADR-0024.
  **Failure-Szenario:** Ein Auditor, der den Skript-Lebenszyklus an der ursprünglichen
  Gate-ADR (0016) verfolgt, liest „bleibt … im Baum" ohne Entfernungs-Notiz.
- **Verifizierbar:** nein (Doc-Lese-Trail).

## Negativbefunde (geprüft, ohne blockierenden Befund)

- **Supersede-Korrektheit/Scope:** ADR-0025 scopt sauber **nur** die „Skript bleibt
  pfad-stabil"-Teilentscheidung von ADR-0024 und benennt explizit, dass VCS-Port, Modul
  `vcs` und der Dogfood-Ersatz des `adr-check`-Gates „unberührt" bleiben. Index-Quellenspalte
  und ADR-0025-`Bezug` stimmen überein. Ehrlich.
- **AGENTS §3.5 (immutable):** Die Bodies von ADR-0024 und ADR-0016 sind **nicht** geändert;
  der „Skript bleibt pfad-stabil"-Core von ADR-0024 bleibt wörtlich stehen. Der Supersede
  läuft regelkonform über eine **neue** ADR; Geschichte-Append korrekt auf Closure terminiert.
- **Spec-Schritt-Reihenfolge:** Schritt 5 (`ignore-refs` **vor** `repo-escape`/
  `codepath-missing`/`anchor-missing`) ist mit ADR-0025 und dem Code konsistent.
- **Ehrlichkeit „Akzeptanztest im Modul `vcs`":** Belegt — `vcs_test.go` (`TestVCSModified`)
  deckt die sieben Selbsttest-Klassen des abgelösten Skripts ab. Kein Garantie-Überclaim.
- **Versions-Kohärenz:** `lastenheft.md` 0.34.0 == §7-Zeile == Roadmap == Slice; auch
  `spezifikation.md` §7 trägt eine slice-054-Zeile. Kohärent.
- **Lastenheft-Akzeptanzkriterium/Out-of-Scope:** DC-FA-CODE-001 wird **erweitert** (kein
  neuer DC-ID → keine Waise); neues „Ventil (ignore-refs)"-Kriterium deckt Happy + Negative
  + Default-leer; Out-of-Scope kohärent.
- **Referenz-Richtung (SDP)/Marker-Ehrlichkeit:** Kein neuer Abwärts-Prosa-Verweis; slice-054
  nur in §7-Verweisspalte + ADR-0025-Geschichte (ausgenommen). Keine neuen Provenance-Marker.
- **Config-Surface (slice-053-Lehre):** `ignore-refs` durchgängig in `--print-config`,
  `--suggest-config` (bewusst leer/auskommentiert), Benutzerhandbuch und `.d-check.yml`.
- **„pfad-stabil"-Entschlackung (Doc-Rollen-Trennung):** In den operativen Dokumenten (AGENTS
  §4, `harness/README.md`, `Makefile`, `.githooks/pre-commit`, `completeness-check.sh`) sind
  alle „pfad-stabiler Fallback"-Hinweise entfernt/aktualisiert; verbliebene Nennungen liegen
  nur in den immutablen ADR-Bodies (korrekt unangetastet). Rollenkonform.
- **Premise-Verifikation:** `codepaths.roots` enthält `tools`; die immutablen ADR-0016/0024
  zitieren `tools/adr-immutable-check.sh` in Inline-Code — der Trap ist real, das Ventil greift.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (MEDIUM-1) |
| LOW | 1 (LOW-1) |
| INFO | 3 (INFO-1/2/3) |

## Verdikt: NICHT-MERGE-FÄHIG (bedingt — ein offener MEDIUM, trivial schließbar)

Kein HIGH; Doc-first-Kohärenz, Supersede-Ehrlichkeit, §3.5-Immutabilität (Bodies),
Versions-Kohärenz, Spec-Schritt-Reihenfolge, der `vcs`-Akzeptanztest-Beleg und die
Config-Surface sind sauber. Es bleibt MEDIUM-1 (vorgezogene Index-Annotation). Trivial
schließbar.

## Einarbeitung (Implementation, 2026-06-29)

- **MEDIUM-1 — behoben:** Die ADR-0024-Index-Annotation wurde zurückgenommen (Index trägt
  wieder schlicht „Accepted"); die Teil-Supersede-Annotation (Geschichte + Index) ist in
  slice-054 §3d nun **explizit erst bei Closure** terminiert (mit Verweis auf R1-MEDIUM-1).
- **LOW-1 — behoben:** §2-Schema und Referenz-Ventil-Absatz nennen jetzt alle drei
  unterdrückten Grund-Codes (`codepath-missing`/`repo-escape`/`anchor-missing`).
- **INFO-1 — behoben:** ADR-0025 schränkt den Claim auf die `codepaths`-Achse ein und führt
  die `links`-Rest-Falle als Konsequenz **und** Re-Evaluierungs-Trigger; slice-054 §4
  dokumentiert sie.
- **INFO-2 — dokumentiert:** slice-054 §4 hält die rückwirkende Entlinkung des Review-Records
  fest (Inhalt unverändert, nur Link-Form nachgezogen; `links` hat kein referenz-weites
  Ventil).
- **INFO-3 — eingeplant:** slice-054 §3d ergänzt einen `## Geschichte`-Append an die immutable
  ADR-0016 (erlaubter Anhang) bei Closure, der die Skript-Entfernung verzeichnet.
