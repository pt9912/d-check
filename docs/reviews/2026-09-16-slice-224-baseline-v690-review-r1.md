# Review-Report: slice-224 — 2026-09-16 (R1)

**Review-Art:** Code-Review — geprüft wird der Diff gegen Slice-Plan,
Hard Rules und die zitierten `MR`-Einträge.

**Gegenstand:** Commit-Kette `87571ec3^..HEAD` (`87571ec3`, `a03e8b1e`,
`1024d52e`, `133ca6f3`, `f12be7c2`).

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-16

> **Zitier-Form.** Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb Kennung statt Adresse: `slice-NNN` statt Lifecycle-Pfad,
> `make <target>` als Token statt Link auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v6.9.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- Slice-Plan `slice-224` (Fassung `in-progress/`, Stand des Prüfzeitpunkts)
- `MR-011`-Pin-Serie · `MR-021` · `MR-039` · `MR-051` · `MR-055` · `MR-069` ·
  `MR-070` · neu `MR-072`; bewegt: `MR-071` (→ `harness/conventions/done/`)
- `DC-FA-PLAN-001`, `DC-QA-03`
- `AGENTS.md` §3.1 · §3.3 · §3.7 · §3.8 · §4 · §5 · §6
- Beobachtungs-Register: `BEO-ALL/pin-bump-mirrors-ungated` (7×),
  `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes` (5×),
  `BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle`,
  `BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage` (2×, unter Schwelle)
- Precedent: `slice-222` (`docs/plan/planning/done/slice-222-baseline-v660-bump.md`)
  und dessen Review `docs/reviews/2026-09-08-slice-222-baseline-v660-review-r1.md`
- Baseline `v6.9.0` · `.harness/baseline/v6.9.0/`

---

## Messformen dieses Laufs (vor den Zahlen, `AGENTS.md` §5)

- **Vorkommen** = jede Fundstelle der Zeichenkette `v6.6.0`
  (`grep -rIo 'v6\.6\.0' | wc -l`).
- **Datei** = jede Datei mit mindestens einer Fundstelle
  (`grep -rlI 'v6\.6\.0' | wc -l`).
- Vergleichsstand „vor der Ersetzung" ist `a03e8b1e` (der letzte Commit vor
  dem Materialisieren des neuen Baums), ausgecheckt über `git archive`.
- Vergleichsstand „Regel-Delta" ist `diff -rq … -I '<!-- Quelle:'` zwischen
  dem archivierten `a03e8b1e`-Baum und dem HEAD-Baum je `regelwerk/` und
  `templates/`.

---

## Findings

### F-1 — Ein neunter lebender Fundort des alten Pins ist weder retargetet noch als eingefroren deklariert

- `kategorie`: MEDIUM
- `quelle`: MR-072 §Adaption (Vollständigkeits-Behauptung „87 Dateien …, 28
  davon der vendorte Baum … 59 Dateien liegen außerhalb" / „acht eingefrorene
  Lauf-Belege … die übrigen 51 lebenden Dateien sind retargetet") ·
  `AGENTS.md` §5 (*„Eine Commit-Botschaft oder Closure-Notiz behauptet nicht
  mehr, als die Arbeit trägt"*) · `BEO-ALL/pin-bump-mirrors-ungated`
- `pfad`: `.d-check.yml:235`
- `befund`: Die Zeile trägt seit slice-222 den Kommentar *„Die v6.6.0-Vorlage
  fuehrt 40 Tabellenzeilen statt 39"* (Begründung für den elften
  `ignore-refs`-Eintrag). Diese Datei liegt außerhalb des vendorten Baums,
  wurde im pre-bump-Stand (`a03e8b1e`) korrekt als eine der 87 Dateien mit
  einer `v6.6.0`-Fundstelle mitgezählt, ist aber **weder** eine der acht
  deklarierten Frozen-Dateien (MR-071 selbst, ein `done/`-Slice, zwei
  Review-Reports, drei Register-Belege, ein CR-Zitat) **noch** unter den 51
  retargeteten. `git show 133ca6f3 -- .d-check.yml` liefert keinen Diff — die
  Datei wurde im Bump-Commit gar nicht angefasst. Die Vollständigkeits-Rechnung
  „8 eingefroren + 51 retargetet = 59" geht damit an dieser Stelle nicht auf:
  Ein neunter Fund wurde weder retargetet noch benannt, obwohl er in einer
  eindeutig lebenden Konfigurationsdatei liegt. Inhaltlich ist die Zeile
  vermutlich zu Recht unverändert (sie beschreibt eine historische
  Tabellenzeilen-Zählung des `v6.5.0`→`v6.6.0`-Übergangs, keine aktuelle
  Pfad-Referenz) — aber genau das hätte die Frozen-Liste sagen müssen, statt
  es stillschweigend auszulassen.
- `verifizierbar`: ja — `grep -rn 'v6\.6\.0' --exclude-dir=.git .` gegen HEAD
  zeigt die Zeile neben den acht deklarierten Frozen-Stellen; `git show
  133ca6f3 -- .d-check.yml` zeigt einen leeren Diff.
- `klasse`: `spiegel-sweep-datei-weder-retargetet-noch-als-eingefroren-deklariert`

### F-2 — Eine unbeteiligte Slice-Plan-Datei trägt jetzt eine falsche Gegenwarts-Aussage über den Pin

- `kategorie`: LOW
- `quelle`: keine unmittelbare Hard Rule dieses Slice — latente
  Wartungsfalle, benachbart zu `BEO-ALL/pin-bump-mirrors-ungated`
- `pfad`: `docs/plan/planning/open/slice-223-commit-zerlegung-ausnahmen-aufloesen.md:148`
- `befund`: §4 *Trigger* sagt *„slice-222 ist geschlossen (der Pin steht auf
  `v6.6.0`, und die Vorlage ist die Grundlage dieses Entscheids)"* — mit dem
  Bump auf `v6.9.0` ist die Gegenwarts-Aussage *„der Pin steht auf `v6.6.0`"*
  jetzt falsch, sofern jemand sie beim künftigen Beanspruchen von slice-223
  als aktuelle Bedingung liest statt als Datumsanker vom 2026-09-08. Der
  Fund liegt außerhalb dessen, was slice-224 an sich selbst bindet (kein
  Abgrenzungspunkt nennt slice-223), und die Datei wurde im Diff auch nicht
  berührt — daher LOW statt MEDIUM: kein Verstoß dieses Slice, aber eine
  Spur, die der nächste Bump oder die Beanspruchung von slice-223 mitnehmen
  sollte.
- `verifizierbar`: ja — `grep -n 'v6\.6\.0' docs/plan/planning/open/slice-223-commit-zerlegung-ausnahmen-aufloesen.md`.
- `klasse`: `stale-pin-aussage-in-fremdem-slice-plan`

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Prüffrage 1 (Gate wird still grün) | geprüft, ohne Befund — keine Gate-Skripte/Workflows in diesem Diff berührt |
| Prüffrage 2 (Kern-Modul meldet falsch) | geprüft, ohne Befund — kein Go-Code im Diff |
| Prüffrage 3 (Hexagon-Import-Verstoß) | nicht anwendbar — kein Code geändert |
| Prüffrage 4 (Suppression ohne ADR) | geprüft, ohne Befund |
| Prüffrage 5 (Netzzugriff außerhalb `external`) | geprüft, ohne Befund — `fetch-baseline-cache.sh` ist einer der vier deklarierten Netz-Werkzeuge außerhalb `gates` (`AGENTS.md` §3.1) |
| Prüffrage 6 (Kommentar ohne der fünf Klassen) | geprüft, ohne Befund — die einzige inhaltlich geänderte Inline-Kommentarzeile (`.d-check.closure.yml:184`) ist eine Kopplungs-Aussage mit Herkunfts-Anker `slice-208`, nur der Tag wurde retargetet |
| Prüffrage 7 (Zustandsfeld erzählt Chronik) | geprüft, ohne Befund — `MR-072`/`MR-071` `**Status:**`-Felder tragen Zustand und Beleg, keine Chronik; die Roadmap-Ruhe-Marker-Commits (`1024d52e`) sind reine `git mv` |
| Prüffrage 8 (Botschaft verallgemeinert über Messung) | siehe F-1 — der einzige Fund |
| Prüffrage 9 (Quelle über Geltungsbereich hinaus zitiert) | geprüft, ohne Befund — `MR-025`, `MR-069`, `MR-070` werden je mit ihrem `Geltungsbereich` gelesen (Plan §1 Abgrenzung 4, §6); `MR-025`s Geltungsbereich (Grund-Code/Config-Schlüssel/Schwellenwert) trifft auf einen reinen Pin-Bump ohnehin nicht zu und wird vom Plan auch nicht zitiert |
| Prüffrage 10 (Messmethode klafft gegen Spec) | nicht anwendbar — keine Spec-Stelle berührt (Slice-Kopf: „—") |
| Prüffrage 11 (zwei Module, dieselbe Eingabeklasse, unterschiedlich behandelt) | nicht anwendbar — kein Modul-Code geändert |
| Prüffrage 12 (Erkennung < Alt-Tool-Familie) | nicht anwendbar |
| Prüffrage 13 (fehlender Negativtest zu neuem Vertrag) | nicht anwendbar — kein neuer öffentlicher Vertrag |
| Prüffrage 14 (Provenance-Marker begründet statt zeigt) | geprüft, ohne Befund — keine `d-check:status-provenance`-Marker im Diff |
| Prüffrage 15 (Modul liest ungescannte Eingabe) | nicht anwendbar — kein Modul-Code geändert |
| Prüffrage 16 (neues `Schärft:`/`Bezug:` nur „§N") | geprüft, ohne Befund — keine neuen `Schärft:`/`Bezug:`-Felder im Diff (`git show 133ca6f3` enthält keine hinzugefügte Zeile dieser Form) |
| Zahlen-Stichprobe „87 Dateien / 181 Vorkommen" (Vorzustand) | nachgerechnet gegen `a03e8b1e` — exakt bestätigt (87/181), ebenso 28 Dateien/29 Vorkommen im vendorten Baum und 59 Dateien außerhalb |
| Delta-Messung „16 von 26 Regelwerk-Dateien / 15 von 28 Templates" | nachgerechnet mit `diff -rq -I '<!-- Quelle:'` gegen den archivierten `a03e8b1e`-Baum — exakt bestätigt |
| `MR-072`/`MR-071`-Tabellenzeilen in `harness/conventions.md` | geprüft — `MR-072` korrekt in §Aktive Adaptionen mit Anchor `mr-072--…`, `MR-071` korrekt nach §Aufgelöste Adaptionen mit Nachfolger-Verweis auf `MR-072`; §Baseline zeigt auf `v6.9.0`/`MR-072` |
| Relative Links in `harness/conventions/done/MR-071-baseline-v660.md` | geprüft — alle vier verschobenen relativen Links (`pin-bump-mirrors-ungated`, `MR-067`, `reviewer.md`, `harness/README.md`, `MR-069`, `MR-051`) tragen die um eine Ebene erhöhte `../`-Tiefe; `make doc-check` (0 Befunde) bestätigt die Auflösung |
| `MR-056`-Zitat-Delta (`MR-039`-Konformität) | geprüft — der wörtlich zitierte Quellsatz in `MR-056` ist byte-identisch geblieben (`diff` gegen `87571ec3^`); nur die umgebende Prosa und der `d-check:cite`-Ausschluss wurden ergänzt, die Direktive selbst korrekt entfernt |
| Verbleibende `d-check:cite`-Direktiven auf `v6.6.0` | geprüft — `grep -rn 'd-check:cite.*v6\.6\.0'` liefert repo-weit **0** Treffer |
| Evidence-Dateien-Form (`pin-bump-mirrors-ungated`, `mechanical-id-rewrite-misses-frozen-classes`) | geprüft — Kopfzeilen `**Vorgang:**`/`**Fund:**` vorhanden, relative Links korrekt tief; die mehrsätzige `**Fund:**`-Prosa weicht von der Baseline-Vorlage („ein Satz") ab, folgt aber genau der seit mehreren Slices gelebten Haus-Form (siehe `slice-222`/`slice-217`-Belege in denselben Verzeichnissen) — kein neuer Bruch, kein Finding |
| Zähler-Stand vs. Dateizahl in `evidence/` | geprüft — `pin-bump-mirrors-ungated/evidence/` enthält 7 Dateien (`slice-106/110/117/148/189/222/224`, Anspruch „7×" bestätigt); `mechanical-id-rewrite-misses-frozen-classes/evidence/` enthält 5 Dateien (`slice-195/202/207/222/224`, Anspruch „5×" bestätigt) |
| MR-013-Bündelung (Lifecycle-Move + Feature-Commit) | geprüft — `MR-071`-Move (`git mv` + Link-Tiefen-Fix) ist im selben Commit wie der große Spiegel-Sweep gebündelt (`133ca6f3`, Similarity 89 %); identisch zum Präzedenzfall `slice-222`/`9c5f4d63`. Kein neuer Verstoß gegen `AGENTS.md` §3.3/`MR-013`, da die Rename-Detection (89 % > 50 %-Schwelle) unangetastet bleibt und die Bündelung dem etablierten Muster folgt |
| `make gates` | grün — `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `make doc-check` separat: `783 Datei(en) geprüft, 0 Befund(e)`; `make baseline-verify` separat: `fetch-baseline-cache: verify ok (54 Dateien, vollständig)` |
| Acht `.claude/rules/`-Symlinks (`MR-055`) | geprüft — alle acht zeigen auf `.harness/baseline/v6.9.0/regelwerk/…` |
| Vier Release-/Tree-URLs, sieben bare Versionsnennungen | Stichprobe gezogen (`AGENTS.md`, `harness/README.md`, `harness/conventions.md`, `.d-check.closure.yml`, `MR-021`, `docs/plan/planning/observations/README.md`, `docs/plan/planning/in-progress/roadmap.md`) — alle auf `v6.9.0` retargetet, keine Ausnahme gefunden |
| `AGENTS.template.md` §4 unverändert (Abgrenzung 2) | geprüft mit `diff -I '<!-- Quelle:'` gegen den `a03e8b1e`-Baum — §4-Zeiger-Form unverändert, einzige Änderungen sind die Release-URL und zwei `welle-<NN>`→`welle-<Kennung>`-Umbenennungen fernab von §4 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:**
`spiegel-sweep-datei-weder-retargetet-noch-als-eingefroren-deklariert` ·
`stale-pin-aussage-in-fremdem-slice-plan`

## Verdikt

**Merge-blockierend:** nein — ein MEDIUM-Finding (F-1) ist vor der Closure
zu klären (Frozen-Liste bzw. Retargeting-Liste um `.d-check.yml:235` ergänzen
oder explizit als neunte Frozen-Stelle benennen), blockiert aber keinen
funktionalen Vertrag: Alle zehn Gates sind grün, die Baseline-Integrität
(`baseline-verify`), alle vier Spiegel-Klassen bis auf die eine Lücke, und
die `MR-039`/`MR-070`/`MR-013`-Disziplin sind eingehalten. Das LOW-Finding
(F-2) betrifft eine fremde, unberührte Slice-Datei und ist kein Blocker.

**Übergabe:** Beide Findings gehen an den Implementer. F-1 gehört in die
Closure-Notiz §7 dieses Slice (Ergänzung der Frozen-/Retargeting-Bilanz) und
von dort in den Zähler bei `BEO-ALL/pin-bump-mirrors-ungated` (achtes
Auftreten einer gate-blinden Spiegel-Lücke, hier speziell: Vollständigkeit
der Frozen-Deklaration selbst). Dieser Report ist ein Lauf-Beleg und ersetzt
keine Verifikation.
