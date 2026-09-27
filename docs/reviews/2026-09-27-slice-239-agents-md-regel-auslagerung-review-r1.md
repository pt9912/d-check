# Review-Report: slice-239 — 2026-09-27

**Review-Art:** Code (gegen Plan + ADR + Konventionen — Maintainability;
Modul 10 §Drei Review-Arten).

**Gegenstand:** Commit `d60f1baf` (`feat(agents): AGENTS.md auf 400 Zeilen
gekuerzt, file[]-Gate aktiviert (slice-239)`).

**Skill:** `.harness/skills/reviewer.md` @ 1.16.0 (2026-09-07; die Version
selbst wurde von diesem Commit nicht gehoben, nur eine `d-check:cite`-Zeile
darin neu geankert).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-27.

**Eingangs-Kontext:**

- `slice-239` (`docs/plan/planning/in-progress/slice-239-agents-md-zeilenbudget-schwelle.md`)
- [ADR-0096](../plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md)
- [ADR-0088](../plan/adr/0088-file-modul-groessengrenzen.md) (Vorläufer, `file`-Modul)
- `DC-FA-FILE-001`
- `AGENTS.md` (Hard Rules, insbesondere §3.3, §3.6, §3.7, §5)
- Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice,
  `modul-04-adrs.md` §Ziel-Form: ADR, `grundlagen-traceability.md`
  §Herkunfts-Anker

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Zeile 11 der neuen `AGENTS.md` §5-Index-Tabelle lässt gegenüber dem Original-Bullet zwei Wortlaut-Teile weg (`stattdessen`; `(Baseline-\`slice.template.md\`)`) — die einzige der vier laut Plan „vollständig in der Tabelle" stehenden Regeln, die tatsächlich nicht wortgleich ist. Der Verweis auf `slice.template.md` als Quelle der `**Lifecycle:**`-Form existiert dadurch nirgends mehr in `AGENTS.md`. | ADR-0096 (§Entscheidung: „Kein Wortlaut geht verloren"); slice-239 §1 („Zusage … bleiben wortgleich, nur der Ort wechselt") | `AGENTS.md:315` | ja — Diff gegen `git show d60f1baf^:AGENTS.md` Zeilen 415–417 | Extraktion verliert Wortlaut trotz „vollständig"-Zusage |
| F-2 | LOW | Uneinheitliche Behandlung interner Abschnitts-Querverweise beim Verschieben: In `harness/rules/docker-make-only.md` wurde `(§4)` beim Kopieren stillschweigend zu einem expliziten Link auf `harness/README.md` §Sensors umgeschrieben (keine reine Kopie); in `harness/rules/kommentare-fuenf-klassen.md` blieb der analoge Fall `(§5)` dagegen unverändert und unqualifiziert stehen — in dieser Datei gibt es kein eigenes §5, der Verweis löst nur auf, wenn man weiß, dass „AGENTS.md §5" gemeint ist. Kein Gate prüft bloße Klammer-Referenzen dieser Form (`(§N)` ist kein Markdown-Link, `links`/`anchors` sehen es nicht). | Maintainability | `harness/rules/docker-make-only.md:53-54`, `harness/rules/kommentare-fuenf-klassen.md:32` | nein — reine Prosa-Referenz, kein Gate deckt `(§N)`-Klammerverweise | dangling section cross-reference nach Extraktion |
| F-3 | LOW | `.harness/skills/reviewer.md` wurde geändert (eine `d-check:cite`-Zeilenspanne von `366-367` auf `281-281` neu geankert, notwendig weil `AGENTS.md` geschrumpft ist), ist aber weder in der `Plan (vor Code)`-Tabelle §3 noch im Abgrenzungs-Abschnitt §1 des Slice-Plans als berührte Datei genannt. Die Änderung selbst ist korrekt und notwendig (bestätigt durch grünen `citations`-Lauf) — der Plan benennt sein eigenes Datei-Set nur unvollständig. | Maintainability | `docs/plan/planning/in-progress/slice-239-agents-md-zeilenbudget-schwelle.md` §3 (Plan-Tabelle) | ja — `git show d60f1baf --stat` zeigt die Datei, der Plan nennt sie nicht | Plan-Datei-Liste unvollständig ggü. tatsächlichem Diff |
| F-4 | LOW | ADR-0096s Fitness-Function-Tabelle nennt als Werkzeug `d-check --enable file` für beide `file`-Regeln. Tatsächlich ist `file` Teil der Default-`modules:`-Liste in `.d-check.yml`; `make doc-check`/`make gates` rufen `$(DCHECK_RUN)` ohne `--enable`/`--disable`-Flags auf (`Makefile:139`). Die Regel wird also über die Config aktiviert, nicht über einen expliziten CLI-Schalter — die Fitness-Function-Zeile beschreibt eine andere Aufrufform als die tatsächlich gebundene. | ADR-0096 (§Fitness Function) | `docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md:101-102` | ja — `grep -n "^doc-check:" -A2 Makefile` zeigt den Aufruf ohne `--enable file` | Fitness-Function-Beschreibung passt nicht zur tatsächlichen Bindung |

Keine HIGH- oder MEDIUM-Findings. Kein Gate-Lüge-Muster, kein Hexagon-Verstoß,
kein Netzzugriff außerhalb `external`, keine Inline-Suppression, keine
verletzte Kommentar-Klasse, kein bare `ADR-*`/`DC-*`/`MR-*` außerhalb
Ausnahme-Pfaden, kein Moduls-Scan-Achsen-Übergriff.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `AGENTS.md` §3.1/§3.7/§3.9 — Überschriften/Anker-Stabilität vor/nach Diff | geprüft, ohne Befund (identisch: `### 3.1 Docker/make-only`, `### 3.7 Kommentare tragen eine der fünf Klassen`, `### 3.9 GitHub-Action-Referenzen sind SHA-gepinnt`, `## 5. Dokumentations-Regeln`) |
| Wortlaut-Treue `harness/rules/docker-make-only.md` ggü. `AGENTS.md` §3.1 (vor Diff) | ein Abweichungspunkt gefunden (F-2), sonst wortgleich inkl. Markdown-Betonung und Link-Ziele |
| Wortlaut-Treue `harness/rules/kommentare-fuenf-klassen.md` ggü. `AGENTS.md` §3.7 (vor Diff) | ein Abweichungspunkt gefunden (F-2), sonst wortgleich |
| Wortlaut-Treue `harness/rules/github-action-sha-pin.md` ggü. `AGENTS.md` §3.9 (vor Diff) | geprüft, ohne Befund — byte-identisch bis auf Pfad-Tiefe der Links |
| Wortlaut-Treue der 13 `harness/rules/dokumentations-regeln/*.md`-Dateien ggü. den zugehörigen §5-Bullets (vor Diff) | ein Abweichungspunkt gefunden (F-1, Zeile 11 der Index-Tabelle selbst — nicht in einer der 13 Volltext-Dateien), die übrigen zwölf wortgleich geprüft |
| Bare `ADR-*`/`DC-*`/`MR-*`-Mentions in allen neuen `harness/rules/`-Dateien (außerhalb Markdown-Links) | geprüft, ohne Befund (`grep` zeigt ausschließlich `[…](…)`-Links; einzige Ausnahme ist der Platzhalter `<DC-ID>.<Buchstabe>`, unverändert aus dem Original übernommen) |
| §3.7 (fünf Kommentar-Klassen) / §3.4-Nachbarregel §3.8 (Modul-Scan-Achse) — Verstöße in den neuen Dateien | geprüft, ohne Befund; die beiden HTML-Kommentare in den neuen Dateien (`d-check:ignore`, `d-check:status-provenance`) sind unverändert aus dem Original bzw. folgen dem etablierten Provenance-Marker-Muster |
| ADR-0096 — interne Konsistenz (Kontext/Entscheidung/Konsequenzen/Alternativen) | geprüft, ein Befund (F-4, Fitness Function); Rest konsistent, `Bezug`/`Schärft`-Felder korrekt befüllt, Re-Evaluierungs-Trigger vorhanden |
| ADR-Index (`docs/plan/adr/README.md`) — neuer Eintrag | geprüft, ohne Befund — Zeile korrekt eingefügt, Tabellenspalten-Längen im `structure`-Gate grün |
| `.d-check.yml` `file:`-Block | geprüft, ohne Befund — zwei Regeln exakt für die beiden Auftraggeber-entschiedenen Ziel-Dateien, `max-lines: 400`, mit begründendem Kommentar samt Herkunfts-Anker |
| `Makefile` `FOCUS_DISABLE` | geprüft, ohne Befund — `--disable file` korrekt ergänzt, spiegelt `.d-check.yml`s `modules:`-Liste |
| `gate_consistency_test.go` `netlessDocModules()` | geprüft, ohne Befund — `"file"` korrekt ergänzt |
| `make gates` auf dem tatsächlich committeten Stand (`d60f1baf`, HEAD) | geprüft, grün: `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `doc-check` einzeln: „882 Datei(en) geprüft, 0 Befund(e)"; Coverage 94,60 % (Schwelle 93 %); Semgrep 0 Findings |
| `AGENTS.md`/`harness/README.md` Zeilenzahl nach Kürzung | geprüft, ohne Befund — 355 bzw. 233 Zeilen, beide unter der 400-Zeilen-Schwelle, deckt sich mit DoD-Behauptung |
| Anzahl referenzierender Dateien auf `AGENTS.md#3.x`/`§5` (Plan behauptet „~24") | geprüft, ohne Befund — grobe Nachzählung (präzisere Anker-Form) findet ~26 Treffer, im Toleranzbereich der Ungefähr-Angabe |
| `harness/conventions.md` — laut Plan „bleibt unangetastet" | geprüft, ohne Befund — Datei nicht im Diff |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Extraktion verliert Wortlaut trotz
„vollständig"-Zusage · dangling section cross-reference nach Extraktion ·
Plan-Datei-Liste unvollständig ggü. tatsächlichem Diff ·
Fitness-Function-Beschreibung passt nicht zur tatsächlichen Bindung.

## Verdikt

**Merge-blockierend:** nein — alle vier Findings sind LOW (kleine
Wortlaut-/Dokugenauigkeits-Lücken ohne Verhaltens- oder Gate-Auswirkung);
`make gates` ist auf dem tatsächlich committeten Stand grün, kein
Hard-Rule-Verstoß, kein stiller Grün-Pfad, keine Architektur-/
Suppression-/Netzverletzung gefunden.

**Übergabe:** Findings gehen an den Implementer zur Kenntnis/optionalen
Nachbesserung (kein Rollen-Konflikt, keine Sequenz nach Modul 8 nötig — alle
vier sind isolierte LOW-Befunde). Die Finding-Klassen gehören bei Closure in
§7 des Slice-Plans und von dort in den Zähler. Dieser Report ist ein
Lauf-Beleg und prüft **nicht** DoD-/Spec-Konformität — das ist Aufgabe des
Verifiers (anderer Kontext, anderes Prüf-Artefakt).
