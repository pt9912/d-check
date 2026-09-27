# Review-Report: slice-235 — 2026-09-27

**Review-Art:** Code — geprüft gegen Plan + Konventionen (Modul 10 §Drei
Review-Arten).

**Gegenstand:** `slice-235-rollen-agents-reviewer-verifier`, Commit-Range
`dc217813..58951d12`.

**Skill:** `.harness/skills/reviewer.md` @ Version 1.16.0 (2026-09-07).
**Modell:** claude-sonnet-5. **Datum:** 2026-09-27.

**Agent-Typ dieses Laufs — expliziter Befund, wie vom Auftraggeber
verlangt.** Dieser Review wurde zuerst mit `subagent_type: "reviewer"`
versucht — dem neuen Agent-Typ, den dieser Slice gerade anlegt
(`.claude/agents/reviewer.md`). Der Aufruf schlug fehl:
`Agent type 'reviewer' not found. Available agents: claude,
claude-code-guide, Explore, general-purpose, Plan, statusline-setup`. Der
Review lief stattdessen als generischer `general-purpose`-Typ mit
Ad-hoc-Auftrag im Prompt. Das bestätigt exakt das in §6 des Risiko-Abschnitts
benannte Risiko ("Ein neu angelegter Agent-Typ löst in der laufenden Sitzung
nicht auf") — die Discovery neuer `.claude/agents/*.md`-Dateien greift
offenbar erst bei Sitzungsstart, nicht innerhalb einer laufenden Sitzung.
Gehört in §7 der Slice-Datei als Ausgang dieses Risikos.

**Eingangs-Kontext:**

- Slice-Plan `slice-235-rollen-agents-reviewer-verifier.md` (insb. §1, §2)
- `.claude/agents/reviewer.md`, `.claude/agents/verifier.md` (neu)
- `.claude/commands/implement-slice.md` (Diff)
- `AGENTS.md` §3, §5, §6
- Baseline-Regelwerk `modul-08-agentenrollen.md`, `modul-11-verification.md`
- `harness/conventions/MR-048-gate-ueber-werkzeug-datei.md`
- `harness/README.md` §Sensors
- Referenz-Repo pg-change-feed, dort `.claude/agents/{reviewer,verifier}.md`
- Code: `internal/hexagon/core/rules/reviews.go`,
  `internal/hexagon/core/rules/reviews_test.go`

---

## Satz-für-Satz-Prüfung (Slice-DoD-Pflicht: „Steht diese Aussage auch in einer gerankten Quelle?")

Geprüft wurden alle Sätze aus `.claude/agents/reviewer.md`,
`.claude/agents/verifier.md` sowie dem `implement-slice.md`-Diff. Frontmatter
(`name`/`description`/`tools`) und reine Datei-/Abschnittspointer zähle ich
nicht als „Aussage" im Sinn der Prüffrage — sie verweisen, statt zu behaupten.

**`.claude/agents/reviewer.md`:**

| Satz (gekürzt) | Gerankte Quelle? | Befund |
|---|---|---|
| „Du bist der Reviewer (Modul 8/10)…" | ja | grounded |
| „Dein Anweisungssatz steht in der Skill-Datei … lies sie als Erstes …" | ja, via Modul 8 §Welche Rolle braucht welche Artefaktklasse (Reviewer = Skill-Artefaktklasse) | grounded |
| „Sie ist repo-gepflegt … Bei Abweichung gilt der Skill." | ja, Instanziierung der allgemeinen Source-Precedence-Klausel auf ein neues Artefakt | grounded (Instanz, keine neue Regel) |
| „Eingang: Diff/Commit-Range + Slice-Plan-Verweis … (AGENTS.md §6 Schritt 8)" | ja (verkürzte Fassung des Eingangs-Kontexts der Skill-Datei) | grounded |
| „Ausgang: Findings … als Report unter docs/reviews/." | ja | grounded |
| „Du prüfst den Diff gegen Plan, Entscheidungen und Hard Rules … nicht gegen die DoD … Zwei Fragen, zwei Antworten, zwei Kontexte." | ja, wörtlich `modul-08-agentenrollen.md` §Regeln gegen typische Fehlannahmen: „Reviewer prüft gegen Plan/ADR (Maintainability). Verification prüft gegen DoD/Spec …" | grounded |
| „Kein Self-Review. Du prüfst Arbeit, die du nicht geschrieben hast, in frischem Kontext." | ja, `AGENTS.md` §6 | grounded |
| **„Übernimm keine Einschätzung des Implementers ungeprüft — auch keine, die plausibel klingt: eine übernommene Einschätzung ist derselbe blinde Fleck, nur zweimal gezählt."** | **nein** — an keiner Stelle in `AGENTS.md`, im Baseline-Regelwerk oder in `harness/conventions.md` zu finden (geprüft: `modul-08-agentenrollen.md`, `modul-10-review-harness.md`, `grundlagen-*`) | **HIGH-Kandidat — siehe F-1** |
| „Ein Finding wird nicht herabgestuft, weil der Implementer widerspricht. Ab HIGH mit Rollen-Widerspruch — oder ab dem dritten gleichen Konflikttyp — läuft der Konflikt-Pfad … über den Architect." | ja, `modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz (Verdikt-Tabelle + „Wann nicht modellieren") | grounded |
| „Bei isolierten LOW/INFO-Findings ist die Sequenz Overkill; dort genügt Annahme oder Begründung." | ja, wörtlich dieselbe Stelle | grounded |
| Bullet-Liste „Deine repo-spezifischen Quellen" | ja, reine Pointer | grounded |

**`.claude/agents/verifier.md`:**

| Satz (gekürzt) | Gerankte Quelle? | Befund |
|---|---|---|
| „Du bist der Verifier (Modul 8/11)…" | ja | grounded |
| „Deine Frage ist „Bauen wir es richtig?" — gegen Plan und DoD. Das ist nicht die Frage des Reviewers … und nicht die des Validators („Bauen wir das Richtige?" …)" | ja, wörtlich `modul-08-agentenrollen.md` — **aber falsch verankert**, siehe F-2 | grounded inhaltlich, **Zitat-Anker falsch** |
| „Eingang: die DoD-Bestätigung des Implementers plus seine Sensor-Belege (AGENTS.md §6 Schritt 8)." | ja | grounded |
| „Ausgang: ein DoD-/Plan-Konformitätsbericht unter docs/reviews/, Dateiname `<YYYY-MM-DD>-slice-<NNN>-<kurz>-verify.md` …" | ja, Instanziierung des allgemeinen Reviewer-Ablage-Musters (`docs/reviews/<YYYY-MM-DD>-<gegenstand>.md`) auf einen neuen, im DoD geforderten Ablage-Entscheid | grounded (lizenzierte Entscheidung, siehe DoD §2) |
| „… `make review-coverage` prüft nur, ob irgendein Dateiname die slice-<NNN>-Kennung trägt, nicht Inhalt oder Suffix — ein Verifikations-Bericht deckt dieselbe Review-Zusage mit." | Tatsachenbehauptung über den Code, **verifiziert** (siehe Abschnitt „Verifier-Ablage-Entscheidung" unten) | grounded/korrekt |
| „Die häufigste Verifier-Lücke ist eine Behauptung ohne Bestätigung." | ja, wörtlich `modul-11-verification.md` §Begriffe | grounded |
| „Prüfe die Belege … nicht die Behauptung — und fahre Sensoren, deren Ausgabe du nicht siehst, selbst nach." | ja, Ausbuchstabierung von `AGENTS.md` §6 Schritt 8 („ein behaupteter Exit-Code ist keiner") | grounded |
| „Eine DoD-Verletzung ist eine Verifier-only-Klasse: unsichtbar für Tests und für das Review." | ja, wörtlich `modul-11-verification.md` §Begriffe | grounded |
| „Ein DoD-Punkt, der sich auf einen Test beruft, ist mit der Verlinkung allein nicht bestätigt. Zeig …, dass der genannte Test ohne den Fix aus dem richtigen Grund rot liefe." | ja, `modul-11-verification.md` §Bewusstes Brechen für DoD-Testbehauptungen | grounded |
| **„Was du NICHT bist: der Reviewer — er sieht den Diff, du siehst die Zusage. Und du bist nicht der Implementer: du reparierst nichts, du berichtest."** | Rollentrennung selbst ist gerankt (Modul 8), aber **dieser konkrete Wortlaut** („du reparierst nichts, du berichtest") steht an keiner geprüften Stelle | **HIGH-Kandidat — siehe F-1 (zweite Instanz)** |
| Bullet-Liste „Deine repo-spezifischen Sensor-Belege" (u. a. „zehn gebundenen Gate-Ziele") | ja, verifiziert gegen `harness/README.md` §Sensors und den tatsächlichen `make gates`-Lauf (10 Ziele: baseline-verify, workflow-pins, doc-check, lint, test, arch-check, coverage-gate, semgrep, gate-consistency, planning-check) | grounded/korrekt |

**`implement-slice.md`-Diff:**

| Zeile | Gerankte Quelle? | Befund |
|---|---|---|
| „Update docs, ADR index, and planning lifecycle if a public contract is touched" | ja, `AGENTS.md` §6 Schritt 7 | grounded |
| „(`AGENTS.md` §5 — CHANGELOG is Release-Prep, not the feature commit)" | ja, wörtlich `AGENTS.md` §5 | grounded |
| „Hand off to the `reviewer` agent type (no self-review, `AGENTS.md` §6 Schritt 8/Modul 8), then to the `verifier` agent type before closure." | ja, `AGENTS.md` §6 Schritt 8 und Modul 8 | grounded |

**Ergebnis der Satz-für-Satz-Prüfung:** zwei Sätze (einer je Agent-Datei)
bestehen die Prüffrage nicht und sind wortgleich aus dem Referenz-Repo
pg-change-feed (dort `.claude/agents/{reviewer,verifier}.md`) übernommen,
ohne dass eine gerankte d-check-Quelle sie trägt.

---

## Verifier-Ablage-Entscheidung (verifiziert am Code)

Behauptung in `verifier.md`: „`make review-coverage` prüft nur, ob
irgendein Dateiname die `slice-<NNN>`-Kennung trägt, nicht Inhalt oder
Suffix."

Geprüft in `internal/hexagon/core/rules/reviews.go`, Funktion
`hasMatchingReview` (Zeile 158–168): Sie iteriert über alle Dateien in
`reviewsDir`, filtert nur auf `Kind == driven.KindFile` und vergleicht
`sliceIDRE.FindString(e.Name) == id` — **keine** Prüfung auf „review" oder
„verify" im Namen, kein Inhalts-Zugriff. Der neue Test
`TestReviewsVerifierNamedReportSatisfiesCoverage`
(`reviews_test.go:225-240`) bestätigt das empirisch mit einer Datei
`…-slice-100-x-verify.md`. Die Behauptung ist **korrekt** — selbst
nachvollzogen, nicht nur am Test abgelesen. `make gates` bestätigt
zusätzlich, dass `reviews_test.go` grün läuft (Docker-Lauf, `go test ./...`
inkl. `internal/hexagon/core/rules`, siehe Gate-Belege unten).

---

## Spiegel-Prüfung (MR-025)

Behauptung im Plan: `AGENTS.md` §6 und `harness/README.md` §Guides nennen
`.claude/commands/implement-slice.md` selbst nicht, also müssen sie die
neuen Agent-Typen auch nicht nennen.

Selbst geprüft: `grep -n "implement-slice" AGENTS.md harness/README.md`
liefert **keinen** Treffer. `grep -n "reviewer\b" AGENTS.md` liefert genau
eine Fundstelle (§6, Zeile 569, Verweis auf die Skill-Datei, nicht auf einen
Agent-Typ). `harness/README.md` §Guides nennt nur die Skill-Dateien
(`reviewer.md`, `closure-note-reviewer.md`), keine `.claude/agents/`-Typen.
Die Behauptung ist **korrekt** — bestätigt, nicht nur vermutet.

---

## Vergleich mit dem Referenz-Repo

Die Auslassung der „Warum dieser Typ existiert — und woran er hängt"-Passage
(Rollen-Erfassung über den Agenten-Typnamen) ist konsequent: d-check führt
keinen entsprechenden Mechanismus, ein Verweis darauf wäre ein
Phantom-Anspruch. Die Auslassung von „Kein Pfeil ohne benennbares Artefakt …"
(pg-change-feed, `reviewer.md` Zeile 29–30) ist dagegen **unnötig streng** —
dieser Satz ist tatsächlich wörtlich in `modul-08-agentenrollen.md`
§Konflikt-Pfad verankert und hätte übernommen werden können; das ist aber
keine Regelverletzung, nur eine ausgelassene Gelegenheit (kein Finding).

Zwei andere Sätze wurden dagegen **trotz** fehlender Verankerung
übernommen (siehe Satz-für-Satz-Tabelle, F-1) — die Auslassungs-Disziplin
war also nicht durchgehend: An der Stelle, wo sie am prominentesten geprüft
wurde (Rollen-Erfassung), wurde sie korrekt angewandt; an zwei
unauffälligeren Stellen (Übernahme einer Implementer-Einschätzung /
Rollenabgrenzung „du reparierst nichts") nicht.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Zwei Sätze führen eine neue, in keiner gerankten Quelle verankerte Verhaltensregel ein und wurden unverändert aus dem Referenz-Repo pg-change-feed (dort `.claude/agents/{reviewer,verifier}.md`) übernommen: (a) „Übernimm keine Einschätzung des Implementers ungeprüft — auch keine, die plausibel klingt: eine übernommene Einschätzung ist derselbe blinde Fleck, nur zweimal gezählt."; (b) „Und du bist nicht der Implementer: du reparierst nichts, du berichtest." Der Slice-Plan (§1) verlangt ausdrücklich, dass jede Datei die Prüffrage „Steht diese Aussage auch in einer gerankten Quelle?" besteht — bei beiden ist die Antwort nein. | Slice-Plan §1 (Prüffrage) | `.claude/agents/reviewer.md:26` und `.claude/agents/verifier.md:28` | ja — Satz-für-Satz-Prüfung gegen `AGENTS.md`, Baseline-Regelwerk und `harness/conventions.md`, reproduzierbar per Grep (kein Treffer) | Werkzeug-Datei führt unbelegte Regel ein (Kanon: keine) |
| F-2 | MEDIUM | Der Verweis auf `modul-08-agentenrollen.md` §Welche Rolle braucht welche Artefaktklasse trägt nicht das Zitat, für das er in `verifier.md` steht: die Sätze „Bauen wir es richtig?"/„Bauen wir das Richtige?" stehen wörtlich unter der Überschrift §Rollen-Regeln (Modul 8), nicht unter §Welche Rolle braucht welche Artefaktklasse (Modul 8). Ein Leser, der dem Anker folgt, findet den Beleg nicht an der genannten Stelle. | Reviewer-Prüffrage #9 (Geltungsbereich/Anker einer Quelle) | `.claude/agents/verifier.md:12` | ja — `grep -n "^###" .harness/baseline/v6.9.0/regelwerk/modul-08-agentenrollen.md` zeigt §Rollen-Regeln bei Zeile 128, §Welche Rolle … bei Zeile 157; das Zitat steht bei Zeile 140–141 | falscher Abschnitts-Anker innerhalb derselben Baseline-Datei |

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.claude/agents/reviewer.md` — Frontmatter (`name`/`description`/`tools`) | geprüft, ohne Befund |
| `.claude/agents/verifier.md` — Frontmatter | geprüft, ohne Befund |
| `.claude/agents/reviewer.md` — Zitate auf `AGENTS.md` §3/§6 | geprüft, ohne Befund |
| `.claude/agents/verifier.md` — „zehn gebundenen Gate-Ziele"-Behauptung | geprüft gegen `harness/README.md` §Sensors und realen `make gates`-Lauf, ohne Befund |
| `internal/hexagon/core/rules/reviews.go` / `reviews_test.go` | geprüft — keine Produktionscode-Änderung, nur ein neuer Bestätigungstest; kein Verhalten geändert, ohne Befund |
| `.claude/commands/implement-slice.md`-Diff | geprüft Satz für Satz, ohne Befund |
| Hexagon-Import-Richtung ([ADR-0005](../../docs/plan/adr/0005-modul-layout-hexagon-ordner.md)) | nicht berührt (reine Doku-/Werkzeug-Dateien), ohne Befund |
| Netzzugriff außerhalb `external` | nicht berührt, ohne Befund |
| Kommentar-Fünf-Klassen (`AGENTS.md` §3.7) | `reviews.go`/`reviews_test.go` tragen nur den neuen Test-Kommentar (Zeile 225–230 in `reviews_test.go`); trägt Zusage + Herkunfts-Anker `slice-235`, keine Chronik/Deliberation, ohne Befund |
| Zustandsfelder (Roadmap/Register) | in diesem Diff nicht berührt, ohne Befund |
| MR-048 (Gate über Werkzeug-Datei) | kein neues Makefile-/`.d-check.yml`-Target hinzugefügt, das Anwesenheit statt Wohlgeformtheit einfordert — Diff bestätigt: keine Änderung an `Makefile`/`.d-check.yml`, ohne Befund |
| `CHANGELOG.md`/`docs/user/`/`README*.md` | unverändert in diesem Diff — konsistent mit dem behobenen Fehler in `implement-slice.md` Schritt 11, ohne Befund |
| Dateigröße (§6 Risiko „Wachstum") | `reviewer.md` 2327 Byte, `verifier.md` 2728 Byte — deutlich unter den in den Schwester-Repos beobachteten 17–28 KB, ohne Befund (Zahl gehört in §7 der Closure-Notiz, DoD-Pflicht) |
| Scan-Zugehörigkeit der Agent-Dateien zu `doc-check` | `scan.roots: ["."]`, `scan.ignore` nennt nur `.harness/baseline/**` und `.harness/cache/**` — die Behauptung „Agent-Dateien liegen im Scan-Bereich" ist korrekt, `make gates` (doc-check) lief grün über den vollen Baum, ohne Befund |

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Werkzeug-Datei führt unbelegte Regel ein
(Übernahme aus Schwester-Repo ohne Prüfung gegen gerankte Quelle) · falscher
Abschnitts-Anker innerhalb derselben Baseline-Datei.

---

## Gate-Belege (unabhängig gefahren)

- `make test` (Docker, `go test ./...`): grün, alle Pakete inkl.
  `internal/hexagon/core/rules` (enthält `reviews_test.go` und den neuen
  `TestReviewsVerifierNamedReportSatisfiesCoverage`).
- `make gates`: grün — Ausgabe endet mit „[gates] baseline-verify +
  workflow-pins + doc-check + lint + test + arch-check + coverage-gate +
  semgrep + gate-consistency + planning-check green"; Coverage-Gate meldet
  94.30 % (Schwelle 93 %); Semgrep 0 Findings über 65 Dateien/55 Regeln;
  `targets`/`planning`-Sonderläufe je 855 Dateien, 0 Befunde.

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (HIGH). Der Slice-Plan macht die
Satz-für-Satz-Prüfung selbst zur DoD-Bedingung; zwei Sätze bestehen sie
nicht. Reparatur ist klein (zwei Sätze streichen oder durch eine gerankte
Formulierung ersetzen bzw. explizit als repo-eigene, im Slice-Plan
begründete Ergänzung kennzeichnen) und ändert an der sonst soliden
Kontext-Trennungs-Architektur der beiden Agent-Dateien nichts.

**Übergabe:** Findings gehen an den Implementer. F-1 und F-2 gehen in die
Slice-Closure §7 und von dort in den Zähler. Dieser Report ist Lauf-Beleg
und ersetzt keine Verifikation — DoD-/Plan-Konformität (insb. die übrigen
DoD-Punkte: Beobachtungs-Register, Risiko-Ausgänge, Closure-Notiz) prüft der
Verifier separat.
