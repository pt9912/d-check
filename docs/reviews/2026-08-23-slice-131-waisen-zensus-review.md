# Review-Report: slice-131 — Reviewer-Skill-Waisen-Zensus — 2026-08-23

**Review-Art:** Code-Review (Doku-/Konventions-Diff gegen Kanon und Slice-Plan,
Modul 10 §Drei Review-Arten) · **Gegenstand:** Commit `61de07d` — der
Feature-Commit von slice-131 (`docs(harness): drei Waisen nach AGENTS.md,
Reviewer-Skill 1.10.0 verweist (MR-015, slice-131)`). Der übergebene Diff-Range
`5530573..61de07d` umfasst zusätzlich den vorgelagerten Lifecycle-Move-Commit
`ea0daf5` (reiner `git mv` + Roadmap-Flip nach `MR-013`); dieser Report prüft
gezielt `61de07d` als den inhaltlichen Feature-Commit, wie im Auftrag benannt.

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`61de07d`, dieselbe Version,
die dieser Commit selbst einführt) · **Modell-ID:** `claude-sonnet-5` ·
**Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-131-reviewer-skill-waisen.md`
  (§1–§9)
- `AGENTS.md` §3 (insb. §3.7, neu §3.8/§3.9), §5, §6
- `.harness/skills/reviewer.md` (Version 1.9.0 vor, 1.10.0 nach diesem Commit)
- `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `Makefile`
  (`FOCUS_DISABLE`-Block)
- Baseline-Kanon `.harness/baseline/v5.11.0/regelwerk/`:
  `grundlagen-source-precedence.md` §Vollständigkeit/§Prüffrage,
  `modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill,
  `modul-09-implementierung.md` §AGENTS.md-Regeln,
  `grundlagen-traceability.md` §Herkunfts-Anker,
  `grundlagen-harness-dateien.md` §Was ein Kommentar trägt,
  `grundlagen-begriffe.md` (Kernbegriff „Harness-Lüge")
- `docs/plan/planning/observations.md` (`BEO-004`, `BEO-009`, `BEO-010`,
  `BEO-011`) und die zugehörigen Closure-Belege
  `docs/plan/planning/done/welle-73-results.md`,
  `docs/plan/planning/done/welle-82-results.md`
- `harness/conventions.md` §MR-015 und
  `harness/conventions/MR-015-agents-md-routet.md` (Geltungsbereich-Feld)
- Git-Historie der Skill-Datei (`git log -S`) zur Verifikation der
  Herkunfts-Anker
- Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext)

**Vom Reviewer selbst gefahren** (nur Lesekommandos, kein `make`,
keine Dateiänderung): `git show`/`git diff` auf `5530573..61de07d`,
`git log -S` auf Versionsstrings der Skill-Datei zur Herkunfts-Verifikation,
Volltext-`grep` über Baseline-Regelwerk, `AGENTS.md`, `harness/conventions.md`,
`spec/`, `docs/user/`, `README.md`, `harness/README.md` nach den drei
zugezogenen Aussagen (Duplikat-Suche), Byte-Zählung der SHA-Pins in beiden
Workflow-Dateien.

**Verdikt: blockierend** — kein HIGH, drei MEDIUM, zwei LOW, kein INFO.

---

## Findings

### F-1 — Das Trenn-Kriterium des Zensus ist im Bestand nicht verankert; die Ableitung, die es tragen würde, fehlt im Commit

- **kategorie:** MEDIUM
- **quelle:** `grundlagen-source-precedence.md` §Vollständigkeit/§Prüffrage
  („Steht jede Aussage dieser Datei auch in einer gerankten Quelle? Nein ⇒
  Waise") · `modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill
  („Operative Pflichtteile") · Beobachtungs-Register `BEO-011`
- **pfad:** Commit-Botschaft `61de07d`, Absatz 3 („Das Kriterium, das die drei
  Waisen trennt: bindet die Aussage ÜBER die Reviewer-Rolle hinaus? …")
- **befund:** Der Wortlaut „bindet über die Reviewer-Rolle hinaus ⇒ braucht
  ein gerangtes Zuhause; bindet nur den Reviewer innerhalb eines vom Kanon
  zugewiesenen Blocks ⇒ buchstabiert aus" steht **nicht** in der Baseline —
  weder in `grundlagen-source-precedence.md` (dort lautet die Prüffrage
  wörtlich anders: „Steht jede Aussage dieser Datei auch in einer gerankten
  Quelle?") noch in `modul-10-review-harness.md`. Die Ableitung, die dieses
  Kriterium tragen würde — dass `modul-10`s „Operative Pflichtteile"
  (Kontext-Eingang, Klassifikation, Anti-Pattern, Output-Schema, Pflege) dem
  Reviewer-Skill genau diese Blöcke als eigenen, in der adoptierten Baseline
  bereits gerankten Ablauf zuweist, und dass nur **innerhalb** dieser Blöcke
  formuliertes Reviewer-eigenes Verhalten „ausbuchstabiert" —, wird im Commit
  nicht gezogen; das Kriterium steht als Behauptung da. Es hält einer
  Nachprüfung stand (an den Nachbarn getestet, siehe Negativbefunde), aber
  genau das ungeprüfte Erscheinen eines am Anlass gebildeten Kriteriums ist
  die von `BEO-011` benannte Gefahrenklasse (b) — und der Slice selbst nennt
  `BEO-011` in §7 als zentrales Risiko dieses Laufs.
- **verifizierbar:** teilweise — kein Gate prüft Kriterien-Herkunft; die
  fehlende Zitat-Kette ist Text-gegen-Text verifizierbar (Commit-Botschaft
  gegen `grundlagen-source-precedence.md`/`modul-10-review-harness.md`).
- **klasse:** trenn-kriterium-ohne-kanon-zitat

### F-2 — Zwei der sieben als „ausbuchstabierend" erklärten HIGH-Anker tragen keine Quelle im Skill-Text; die im Commit behauptete Zuordnung passt für einen davon inhaltlich nicht

- **kategorie:** MEDIUM
- **quelle:** Slice-Plan §4 DoD-1 („Je Anker eine Antwort mit Quelle; die
  Tabelle ist vollständig **belegt**, nicht behauptet, `BEO-011`") ·
  `AGENTS.md` §6 Schritt 8 (Text bei Zeile 344-345)
- **pfad:** `.harness/skills/reviewer.md:20-22` (Stilles-Grün-Pfad,
  Korrektheitsfehler, Gate-Suppression ohne ADR)
- **befund:** Von den sieben HIGH-Bullets tragen fünf ein Inline-Zitat
  (`(Harness-Lüge)`, `ADR-0005`, `DC-QA-03`, zweimal die
  „Was-ein-Kommentar-trägt"-Baseline-Referenz). „Korrektheitsfehler in
  Kern-Modulen mit falschen Befunden/Exit-Codes" und „Gate-Suppression ohne
  ADR" tragen **keines**; ihr Zuhause steht nur in der Commit-Botschaft
  („AGENTS.md §3.2/§3.6, §4/§6.8"). Für „Gate-Suppression ohne ADR" ist die
  Zuordnung zu §3.2 (Suppression-Verbot)/§3.6 (Gates nicht ohne ADR lockern)
  inhaltlich naheliegend. Für „Korrektheitsfehler in Kern-Modulen mit
  falschen Befunden/Exit-Codes" passt die implizierte Zuordnung zu §6.8
  jedoch nicht: §6.8 verlangt ehrliches **Berichten** eines Agenten über
  gelaufene Gates („keine Erfolgsmeldung ohne Gate-Ausführung … ein
  behaupteter Exit-Code ist keiner") — das ist eine Aussage über
  Agent-Reporting-Disziplin, keine über Korrektheitsbugs im Produkt
  `d-check` selbst, die falsche Befunde/Exit-Codes produzieren. Der Zensus
  behauptet „alle sieben HIGH-Anker … haben ein gerantes Zuhause", ohne für
  diesen einen einen belastbaren Beleg zu zeigen — dieselbe DoD-Formel
  („belegt, nicht behauptet") wird an der eigenen Stelle verfehlt, an der sie
  eingefordert wird.
- **verifizierbar:** ja — Volltext-Vergleich `AGENTS.md` §6 Schritt 8 gegen
  die Bedeutung von „Korrektheitsfehler in Kern-Modulen".
- **klasse:** zensus-behauptet-statt-belegt

### F-3 — MR-015 wird als Wächter gegen „AGENTS.md wird Sammelstelle" zitiert, sein eigenes Geltungsbereich-Feld deckt aber nur §1

- **kategorie:** MEDIUM
- **quelle:** `harness/conventions/MR-015-agents-md-routet.md`
  (Geltungsbereich: „`AGENTS.md` §1, `MR-011`, `MR-012`, §Adoptierte
  Konventions-Quellen") · Slice-Plan §2 Schritt 3 („Prüfen, ob dabei
  `AGENTS.md` zur Sammelstelle wächst (`MR-015`: sie **routet**, spiegelt
  nicht)")
- **pfad:** Commit-Subjekt `61de07d` (`MR-015, slice-131`); Slice-Plan §2
  Schritt 3
- **befund:** `MR-015` ist die Auflösung einer Pointer-Drift beim
  Baseline-Bundle-Verweis in `AGENTS.md` §1 (welche Rohdatei/Version die
  Konventionsquelle nennt) — sein Geltungsbereich-Feld nennt ausdrücklich nur
  §1, `MR-011`, `MR-012` und §Adoptierte Konventions-Quellen, nicht §3/§5.
  Der Slice-Plan und das Commit-Subjekt zitieren `MR-015` aber als das
  Instrument, gegen das geprüft wird, ob die drei neuen Hard-Rule-Zuzüge
  `AGENTS.md` zur inhaltlichen Sammelstelle machen — eine andere Frage als
  die, die `MR-015` regelt. Der tatsächlich tragende Anker für diese Prüfung
  ist `grundlagen-source-precedence.md` §Vollständigkeit plus `AGENTS.md`s
  eigener Kopfsatz („dupliziert deren Inhalt nicht"), keiner von beiden wird
  im Commit genannt. `make trace-check` prüft nur die ID-**Präsenz**, nicht
  die semantische Passung — der Fehlgriff bleibt ungefangen, bis ihn ein
  Review liest.
- **verifizierbar:** ja — Geltungsbereich-Feld von `MR-015` gegen die im
  Commit behauptete Verwendung.
- **klasse:** id-zitat-mit-scope-mismatch

### F-4 — Zensus-Botschaft dokumentiert keine Prüfung von Anti-Pattern, Kontext-Eskalation, Output-Schema, Negativbefunde, Ablage, Eingangs-Kontext, LOW und INFO

- **kategorie:** LOW
- **quelle:** Slice-Plan §1 („Der Vollständigkeits-Zensus … hat fünf
  Fundorte ergeben") und §4 DoD-1 (Vollständigkeitsanspruch)
- **pfad:** Commit-Botschaft `61de07d` (gesamt)
- **befund:** Die Commit-Botschaft argumentiert ausschließlich über die
  HIGH- und die zwei MEDIUM-Anker (plus `FOCUS_DISABLE` und die
  Workflow-Köpfe). Sie erwähnt an keiner Stelle, dass Anti-Pattern
  (`.harness/skills/reviewer.md:109-121`, insb. den nicht kanonisch belegten
  Begriff „REFUTED"), Kontext-Eskalation (Zeile 104-107, insb. „dieselbe
  Beobachtung im Gate-/Sicherheitspfad steigt eine Stufe" — im Kanon nicht
  wörtlich auffindbar), Output-Schema, Negativbefunde, Ablage,
  Eingangs-Kontext sowie die LOW-/INFO-Kategorien geprüft wurden. Eine
  eigene Nachprüfung dieser Bereiche (siehe Negativbefunde unten) findet dort
  zwar keine Waise, aber der Zensus selbst weist diese Prüfung nirgends aus —
  bei einem DoD-Punkt, der ausdrücklich „vollständig belegt, nicht behauptet"
  verlangt, bleibt der Beleg für diese Bereiche aus.
- **verifizierbar:** ja — Volltext-Suche der Commit-Botschaft nach den
  Abschnittsnamen.
- **klasse:** zensus-scope-nicht-dokumentiert

### F-5 — Neue AGENTS.md-§5-Zeile führt die Notation „§6.8" ein, die von der etablierten Form „§6 Schritt N" abweicht

- **kategorie:** LOW
- **quelle:** `harness/conventions/MR-031-schritt-3-benennen.md` (verwendet
  „`AGENTS.md` §6 Schritt 3" für denselben nummerierten Block) ·
  `AGENTS.md` §6 (nummerierte Liste 1.–8., keine `### 6.x`-Unterüberschriften)
- **pfad:** `AGENTS.md:321`
- **befund:** `AGENTS.md` verwendet die Form `§X.Y` sonst ausschließlich für
  echte nummerierte Unterüberschriften (`### 3.1`–`### 3.9`). §6 ist eine
  einfache nummerierte Liste ohne solche Unterüberschriften; der bestehende
  Repo-Präzedenzfall (`MR-031`) referenziert denselben Workflow-Schritt
  deshalb als „§6 Schritt 3", nicht als „§6.3". Die neue Zeile schreibt
  stattdessen „(§6.8)" — eine Notation, die im ganzen Repo vor diesem Commit
  nirgends vorkommt und mit der Unterüberschriften-Form kollidiert, ohne
  einen auflösbaren Anker (`#68-…`) zu haben.
- **verifizierbar:** ja — `grep -rn "§6\.8\|§6 Schritt"` über `AGENTS.md`
  und `harness/conventions/`.
- **klasse:** notation-inkonsistent-mit-praezedenz

## Negativbefunde

- geprüft, ohne Befund: `.github/workflows/ci.yml`/`release.yml` — alle
  `uses:`-Einträge (`actions/checkout`, `docker/login-action`) tragen einen
  vollen 40-stelligen Commit-SHA mit Tag-Kommentar; die neue `AGENTS.md`-§3.9
  stimmt mit dem tatsächlichen Bestand überein, und es gibt keine dritte
  Workflow-Datei
- geprüft, ohne Befund: Herkunfts-Anker `(seit welle-73)` an `AGENTS.md`
  §3.8 — belegt in `docs/plan/planning/done/welle-73-results.md:106`
  („BEO-004 stand seit welle-70 bei 3 und ist in dieser Closure [welle-73]
  verkörpert"); der Anker nennt korrekt die Welle der **Verkörperung**, nicht
  die des Schwellen-Übertritts, wie `grundlagen-traceability.md` verlangt
- geprüft, ohne Befund: Herkunfts-Anker `(seit welle-82)` an der neuen
  `AGENTS.md` §5-Zeile — der MEDIUM-Anker „Botschaft verallgemeinert über die
  Messung hinaus" wurde per `git log -S` nachweislich in Commit `a07e732`
  (Version 1.9.0) eingeführt; dieser Commit gehört zu slice-123, das
  `docs/plan/planning/done/welle-82-results.md:23` als Teil von welle-82
  führt — die Welle-Nummer stimmt trotz einer widersprüchlichen
  Nebenbemerkung an anderer Stelle desselben (vorbestehenden, nicht von
  diesem Commit berührten) Dokuments
- geprüft, ohne Befund: `AGENTS.md` §3.9 trägt bewusst **keinen**
  Herkunfts-Anker — korrekt nach `grundlagen-traceability.md`
  §Herkunfts-Anker („Ab Einführung, kein Nachrüsten … `seit unbekannt` wäre
  eine Harness-Lüge"), da die SHA-Pin-Konvention nicht aus dem Steering Loop
  stammt, sondern Bestand seit den Workflow-Dateien selbst ist
- geprüft, ohne Befund: `AGENTS.md` §3.8 gibt den Inhalt von `BEO-004`
  (Scan-Menge vs. gelesene, nicht gescannte Eingaben; die Folge kann still
  sein) inhaltlich und nahezu wortgleich zum vorbestehenden
  Reviewer-Skill-Text wieder — kein Bedeutungsverlust beim Umzug
- geprüft, ohne Befund: die neue `AGENTS.md`-§5-Zeile gibt `BEO-009`
  Richtung (b) („verallgemeinert über die Messung hinaus") korrekt wieder
  und grenzt sie sauber gegen Richtung (a) ab (Verweis auf §6.8 für „Probe
  muss gelaufen sein")
- geprüft, ohne Befund: die Widerlegung des `FOCUS_DISABLE`-Befunds ist
  wortgenau belegt — „Spiegelt die `.d-check.yml`-modules-Liste; wächst die
  dort, hier nachziehen" beantwortet exakt die Kanon-Frage der Klasse
  **Kopplung** aus `grundlagen-harness-dateien.md:84` („Was muss ich
  mitändern, wenn ich das hier ändere?"); die verbleibende Lücke ist
  korrekt als bereits registriertes `BEO-010` benannt statt verschwiegen
- geprüft, ohne Befund: keine der drei zugezogenen Aussagen (§3.8, §5-Zeile,
  §3.9) stand vor diesem Commit in `spec/`, den ADRs, `docs/user/`,
  `README.md`, `harness/README.md` oder `harness/conventions.md` — Repo-weite
  Grep-Suche liefert null Treffer außerhalb des Reviewer-Skills bzw. der
  Workflow-Kommentare; `AGENTS.md` wird durch den Zuzug selbst nicht zur
  wörtlichen Sammelstelle
- geprüft, ohne Befund: die neuen `AGENTS.md`-Begründungen zitieren
  ausschließlich Wellen (`seit welle-73`, `seit welle-82`), keine
  Slice-Nummern — konform mit dem Slice-Verweis-Verbot in `AGENTS.md`/
  `harness/README.md`
- geprüft, ohne Befund: die geänderten Workflow-Kommentare
  (`ci.yml:16`, `release.yml:19`) sind nach der Änderung Rang-Zeiger-Klasse
  („Wo wohnt die Norm, deren Umsetzung das hier ist?") und damit konform zu
  `AGENTS.md` §3.7 — vorher trugen sie eine unverankerte
  Policy-Behauptung, jetzt einen auflösbaren Verweis
- geprüft, ohne Befund: die Versionshebung des Reviewer-Skills 1.9.0→1.10.0
  ist additiv (zwei angehängte Sätze, keine Streichung/Überschreibung),
  Datum aktualisiert — konform `modul-10-review-harness.md` §Pflege
  („Skill-Datei selbst wird nicht überschrieben, sondern versioniert")
- geprüft, ohne Befund: keine Kommentar-Klassen-Verstöße (§3.7) in den
  geänderten Code-/Config-Zeilen dieses Commits (die einzigen Code-/
  Config-Änderungen sind die zwei Workflow-Kommentarzeilen, siehe oben)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** `trenn-kriterium-ohne-kanon-zitat` ·
`zensus-behauptet-statt-belegt` · `id-zitat-mit-scope-mismatch` ·
`zensus-scope-nicht-dokumentiert` · `notation-inkonsistent-mit-praezedenz`

## Verdikt

**Merge-blockierend:** ja — drei MEDIUM offen. Kein Finding widerlegt die
Kern-Prämisse des Slice (drei echte Waisen, saubere Umzüge, Skill verweist
korrekt); F-1 bis F-3 betreffen die **Belegkette** des Zensus selbst, nicht
seine Schlussfolgerung — genau die Sorgfaltsebene, die der Slice mit seinem
eigenen `BEO-011`-Risikohinweis in §7 als kritisch benennt. F-1: das
Trenn-Kriterium hält einer Nachprüfung stand, steht aber unbelegt im Commit.
F-2: zwei der sieben „ausbuchstabierenden" HIGH-Anker tragen keine
nachprüfbare Quelle, einer davon (Korrektheitsfehler → §6.8) ist inhaltlich
falsch zugeordnet. F-3: `MR-015` wird außerhalb seines eigenen
Geltungsbereichs als Sammelstelle-Wächter zitiert. F-4/F-5 sind
Genauigkeits-Nachbesserungen ohne Substanzrisiko.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan bei
Plan-Defekt); die Finding-Klassen gehen zusätzlich in die Slice-Closure §7
und von dort in den Zähler. Dieser Report selbst ist ein Lauf-Beleg —
DoD-/Spec-Konformität prüft der Verifier separat.
