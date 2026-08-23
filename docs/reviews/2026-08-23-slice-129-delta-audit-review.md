# Review-Report: slice-129 — Etappe B, Delta-Audit über acht Kurs-Wellen

**Datum:** 2026-08-23 · **Review-Art:** Doku-/Konventions-Review (geprüft gegen Slice-Plan slice-129 §1–§5, Wellendokument welle-83 §1/§3/§4/§6, Baseline-Regelwerk `grundlagen-source-precedence.md`/`grundlagen-durchsetzungsschicht.md`/`grundlagen-referenz-richtung.md`/`modul-09-implementierung.md`, Beobachtungs-Register `BEO-002`/`BEO-009`/`BEO-011`), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** `git diff c72f763..HEAD` — zwei Commits: `a5b9d90` (Lifecycle-Move `open/` → `in-progress/`), `df0ab90` (Delta-Audit selbst + Etappe-C-Schnitt `slice-130`/`slice-131` + Drift-Log-Zeile)
**Skill:** `.harness/skills/reviewer.md` @ 1.9.0 (`5331466`) · **Modell-ID:** `claude-sonnet-5`

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/done/slice-129-baseline-v5110-delta-audit.md`, besonders §3/§4/§4a/§5
- Wellendokument `docs/plan/planning/welle-83-baseline-v5110-migration.md`
- Neue Slices `docs/plan/planning/open/slice-130-lastenheft-historie-form.md`, `slice-131-reviewer-skill-waisen.md`
- `docs/plan/planning/in-progress/roadmap.md` §Historische Trigger-Verschiebungen (Drift-Log)
- Vendorter Kanon `.harness/baseline/v5.11.0/regelwerk/`: `grundlagen-source-precedence.md`, `grundlagen-durchsetzungsschicht.md`, `grundlagen-referenz-richtung.md`, `modul-09-implementierung.md`, `modul-07-carveouts.md`; Vorlagen `templates/spec/lastenheft.template.md`, `templates/spec/spezifikation.template.md`
- Schwester-Repo `ai-harness-course` (lokaler Klon, read-only): `CHANGELOG.md`, `git diff v5.9.0..v5.11.0`, Commits `b72b61d`/`355f8dd`/`79f2e5c`/`387677b`/`76813cd`/`e580c0e`/`725279a`/`3ce5982`
- `docs/plan/planning/observations.md` (`BEO-011`, `BEO-002`, `BEO-009`)
- `spec/lastenheft.md` §7, `spec/spezifikation.md` §7, `AGENTS.md` §4, `.d-check.yml` (`targets`-Block)
- Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext)

**Vom Reviewer selbst gefahren** (Exit je Lauf direkt gelesen, `BEO-007`):
`make gates` Exit 0 (acht Gates: `doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check`; 460 Datei(en), 0 Befund(e); Coverage 94,80 %) · `make gate-consistency` isoliert Exit 0 (460/0) · `grep -inE '\b(muss|müssen|darf|dürfen|soll|sollen|pflicht|verboten|verpflichtend|zwingend)\b' CLAUDE.md` Exit 1 (0 Treffer) · `git diff v5.9.0..v5.11.0 --stat`/`-- lab/regelwerk/`/`-- kurs/de/` im Kurs-Repo · Volltext-Lesung der neun Wellen-90/91/92/93/94-Commits im Kurs-Repo · Volltext-Vergleich `spec/lastenheft.md` §7 gegen `spec/spezifikation.md` §7 gegen beide vendorten Templates · Suche nach zusätzlichen Nicht-Rangliste-Kandidaten in `Makefile`, `.github/workflows/*.yml`, `tools/harness/*.sh`, `tools/*.sh`, `docs/plan/planning/README.md`.

**Verdikt: blockierend** — kein HIGH, vier MEDIUM, ein LOW, kein INFO.

---

## Findings

**F-1**

- **kategorie:** MEDIUM
- **quelle:** `grundlagen-source-precedence.md` §Vollständigkeit i. V. m. `lab/regelwerk/README.md` (Kurs-Repo, unter `.harness/baseline/v5.11.0/regelwerk/README.md` vendort) · `BEO-011`
- **pfad:** `docs/plan/planning/done/slice-129-baseline-v5110-delta-audit.md:86-90,99-100`
- **befund:** Die tragende Begründung für sechs der acht Wellen-Zeilen ist „kein Regelwerks-Diff" (wörtlich für Welle 92/93, sinngleich für 87/88/91/94-Gegenprobe), zusammengefasst in Z. 89–90: „Wellen ohne Regelwerks-Änderung sind damit **belegt** folgenlos für den Kanon — nicht vermutet." Das vendorte `regelwerk/README.md` selbst widerspricht der Prämisse, dass ein Regelwerks-Diff die maßgebliche Prüfgröße ist: „Es trägt **keine eigene Normativität**: maßgeblich für den *Inhalt* bleibt der Kurs unter `/kurs/de/`" und ausdrücklich „**Was dieses Verzeichnis NICHT ist.** Eine eigene Quelle der Wahrheit." Der Audit prüft damit gegen einen didaktik-freien *Extrakt*, nicht gegen die selbst benannte Quelle. Das Audit-Ergebnis hält in der Sache stand (eigene Gegenprobe: `kurs/de/` ändert sich in v5.9.0..v5.11.0 exakt in den drei Wellen 89/90/94, `lab/example/AGENTS.md` — die vom Kurs selbst benannte „gelehrte Form" für Adopter — bleibt über alle acht Wellen unverändert), aber der im Slice-Text **geführte** Beleg zeigt das nicht: Kein Wort in §4a prüft `kurs/de/` oder `lab/example/` direkt, und die Prämisse „Regelwerk ⇔ Kanon" wird nirgends gegen ihre eigene, gegenteilige Selbstauskunft gehalten.
- **verifizierbar:** teilweise — die Regelwerk-Selbstauskunft ist Zitat (`.harness/baseline/v5.11.0/regelwerk/README.md:17-25`); die Äquivalenz „Regelwerks-Diff ⇔ kurs/de-Diff" ist für diesen Zeitraum durch `git diff v5.9.0..v5.11.0 --stat -- kurs/de/` bestätigt, aber kein Gate prüft diese Äquivalenz strukturell für künftige Wellen.
- **klasse:** vollstaendigkeits-praemisse-nicht-am-kanon-selbst-geprueft

**F-2**

- **kategorie:** MEDIUM
- **quelle:** `BEO-011` (Vollständigkeits-Aussage aus dem Anlass statt aus dem Bestand) · Slice-Plan §4 DoD dritter Punkt („vollständig **belegt**, nicht behauptet, `BEO-011`") · Slice-Plan §7 (nennt `BEO-011` selbst als zentrales Risiko für genau diesen Zensus)
- **pfad:** `docs/plan/planning/done/slice-129-baseline-v5110-delta-audit.md:110-124` (§Welle 94 — der Vollständigkeits-Zensus)
- **befund:** Der Zensus nennt drei Fundorte (`CLAUDE.md`, `.harness/skills/reviewer.md`, `.claude/commands/implement-slice.md`) und schließt `.githooks/*`/`.claude/hooks/*` als Durchsetzung aus — ohne zu zeigen, dass diese vier Kandidaten das vollständige Suchfeld waren. Eine eigene Suche über nicht betrachtete Artefakt-Klassen liefert zwei weitere Kandidaten mit genau demselben Muster (Regel-Text ohne Rückbindung an eine gerankte Quelle, `harness/conventions.md` oder die Baseline): (1) `Makefile:203-209` legt fest, dass `FOCUS_DISABLE` bei jedem neuen `.d-check.yml`-Modul manuell nachgezogen werden **muss** („Spiegelt die `.d-check.yml`-modules-Liste; wächst die dort, hier nachziehen") — diese Pflicht steht in keiner der drei Rangquellen-Klassen (`AGENTS.md`, `harness/conventions.md`, `harness/README.md` kennen `FOCUS_DISABLE` nicht); (2) `.github/workflows/release.yml:18-19` und `ci.yml:15-16` legen die „Action-Pinning"-Pflicht fest (jeder `uses:`-Eintrag SHA-gepinnt mit Tag-Kommentar) — auch diese Regel ist in keinem ADR, `AGENTS.md` oder `harness/conventions.md` dokumentiert. Beides sind Regeln, denen ein Agent folgen muss, der das Makefile bzw. die Workflows ändert, und beide „legen fest" statt einen gerankten Ablauf auszubuchstabieren (Prüffrage aus `grundlagen-source-precedence.md`: „Steht jede Aussage dieser Datei auch in einer gerankten Quelle?"). Die vom Delta-Audit selbst benannte Sorge trifft damit den Zensus, den es führt: eine Vollständigkeits-Aussage wird aus den bereits bekannten Kandidaten (Anlass: `CLAUDE.md` aus slice-127, die zwei naheliegenden Skill-/Command-Dateien) gebildet, nicht aus einer belegten Durchsuchung des Bestands (`Makefile`-Kommentare, `.github/workflows/`, `tools/*.sh`, `docs/plan/planning/README.md` werden nirgends erwähnt — auch nicht als geprüft-und-verworfen).
- **verifizierbar:** ja — `Makefile:203-209`, `.github/workflows/release.yml:18-19`, `.github/workflows/ci.yml:15-16` sind Zitate; `grep -rln "FOCUS_DISABLE\|Action-Pinning" AGENTS.md harness/conventions.md docs/plan/adr/*.md` liefert 0 Treffer.
- **klasse:** unvollstaendiger-zensus-als-vollstaendig-gefuehrt

**F-3**

- **kategorie:** MEDIUM
- **quelle:** `BEO-009` Richtung (a) (behauptete Prüfung fand so nicht statt) · vendorte Vorlage `templates/spec/spezifikation.template.md` §7
- **pfad:** `docs/plan/planning/done/slice-129-baseline-v5110-delta-audit.md:108` und wortgleich `docs/plan/planning/open/slice-130-lastenheft-historie-form.md:31-32,41`
- **befund:** slice-129 sagt „Unsere Historie-Tabellen führen drei Spalten" (Plural, unspezifisch), slice-130 macht daraus explizit: „Unsere **beiden** Historie-Tabellen (Lastenheft, Spezifikation) führen drei Spalten" und plant unter „Vorgehen" Punkt 1, die vierte Spalte „in **beiden** Straten-Historien" einzuführen. Das ist an `spec/spezifikation.md:2790-2792` widerlegt: Die dortige `## 7. Historie`-Tabelle führt **zwei** Spalten (`| Datum | Änderung |`), nicht drei. Das ist keine Nachlässigkeit im Bestand, sondern die vom selben v5.11.0-Kanon vorgeschriebene Form: `templates/spec/spezifikation.template.md:122` zeigt exakt dieselben zwei Spalten und ausdrücklich „**kein** ADR- und kein Slice-Verweis" — die Spezifikation hat kein `Version`-Kopffeld (anders als das Lastenheft, `**Version:** 0.65.2`) und keinen CR-Prozess, dem eine `Verweis`-Spalte etwas hinzuzufügen hätte; Kurs-Welle 90 hat `spezifikation.template.md` denn auch **nicht** geändert (Diff v5.9.0..v5.11.0 berührt nur `lastenheft.template.md`, +29/−7). Der von diesem Audit geschnittene Folge-Slice 130 ist damit auf einer falschen Bestandsaufnahme aufgesetzt: Wörtlich umgesetzt würde er der Spezifikation eine `Verweis`-Spalte geben, die der Kanon für dieses Stratum explizit **nicht** vorsieht.
- **verifizierbar:** ja — `spec/spezifikation.md:2790-2792` vs. `.harness/baseline/v5.11.0/templates/spec/spezifikation.template.md:115-124` vs. `.harness/baseline/v5.11.0/templates/spec/lastenheft.template.md:123-154`, direkt gegeneinander gelesen.
- **klasse:** bestandsaufnahme-nicht-am-artefakt-selbst-geprueft

**F-4**

- **kategorie:** MEDIUM
- **quelle:** Referenz-Genauigkeit/Herkunfts-Anker (`grundlagen-traceability.md` §Herkunfts-Anker) · Kurs-`CHANGELOG.md` Welle 94, Review-Nacharbeit
- **pfad:** `docs/plan/planning/open/slice-131-reviewer-skill-waisen.md:13`
- **befund:** Der `Bezug:`-Kopf zitiert „`modul-09-implementierung.md` §AGENTS.md-Regeln **(Zwei-Quadranten-Regel)**" als Baseline-Beleg. Der Begriff ist ein Phantom: Das Kurs-`CHANGELOG.md` dokumentiert in der Review-Nacharbeit zu genau dieser Welle (94) ausdrücklich, dass „**Zwei-Quadranten-Regel**" eine „Wortprägung [war], die es im Korpus nicht gibt, und ihre Glosse gab die bestehende Regel falsch wieder … — gestrichen". Er kommt nirgends im vendorten `.harness/baseline/v5.11.0/regelwerk/*.md` als benannte Regel vor (eigene Volltextsuche über alle 24 Regelwerks-Dateien: 0 Treffer auf „Zwei-Quadranten-**Regel**"; die Wortfolge „zwei Quadranten" existiert nur unbenannt in `modul-09-implementierung.md:170` und `modul-11-verification.md:68` — dort geht es um Hard-Rule-**Durchsetzungsdualität** (AGENTS.md-Eintrag + Gate), nicht um die Frage, **wo** normativer Text stehen darf, die slice-131 tatsächlich bearbeitet und die korrekt in `grundlagen-source-precedence.md` §Vollständigkeit steht (dort im selben `Bezug:`-Feld separat zitiert). Ein Agent, der dem Zitat folgt, um die „Zwei-Quadranten-Regel" im Kanon nachzuschlagen, findet sie nicht — oder wendet versehentlich die tatsächliche Zwei-Quadranten-Regel (Durchsetzungsdualität) auf die falsche Frage an. Bemerkenswert am Fundort: slice-131 klassifiziert Waisen-Regeln (Text ohne Rückbindung an eine gerankte Quelle) und zitiert dabei selbst einen Begriff ohne Rückbindung.
- **verifizierbar:** ja — `grep -rn "Zwei-Quadranten" CHANGELOG.md` im Schwester-Repo `ai-harness-course` (Zeilen 37, 83–85) und `grep -rn "Quadranten" .harness/baseline/v5.11.0/regelwerk/*.md` (kein Treffer auf „…-Regel").
- **klasse:** zitat-auf-gestrichenen-begriff

**F-5**

- **kategorie:** LOW
- **quelle:** Maintainability / GFM-Tabellenform
- **pfad:** `docs/plan/planning/welle-83-baseline-v5110-migration.md:74-89` (§4 Slices in dieser Welle)
- **befund:** Die Tabelle „Slices in dieser Welle" trägt Kopfzeile + zwei Datenzeilen (`slice-128`, `slice-129`), dann eine Leerzeile (unverändert aus der Vorfassung übernommen, dort trennte sie Tabelle von Fließtext), dann die zwei neuen Zeilen für `slice-130`/`slice-131` **ohne eigene Kopf-/Trennzeile**. Nach GFM endet ein Tabellenblock an der ersten Leerzeile; die zweite Zeilenpaar-Gruppe hat keine vorangehende `| Spalte | Spalte |`/`|---|---|`-Definition und wird deshalb nicht als Fortsetzung derselben Tabelle gerendert. Sichtbare Folge: `slice-130`/`slice-131` erscheinen visuell von `slice-128`/`slice-129` abgetrennt und ohne die Spaltenköpfe „Slice | Rolle". `make gates` bleibt grün (kein `structure.table-*` ist auf diesen Abschnitt konfiguriert) — das Modul deckt sechs benannte chronologische Tabellen, diese Slice-Übersichtstabelle steht nicht darunter.
- **verifizierbar:** ja — Datei direkt gelesen; `.d-check.yml:258-277` listet die sechs konfigurierten `table-order`-Ziele, `welle-83…md` ist keines davon.
- **klasse:** tabellenblock-durch-leerzeile-gespalten

---

## Negativbefunde

1. **Scope-Treue gegen §3 gehalten.** Außerhalb der drei neuen/geänderten Slice-/Wellen-Dokumente und der Drift-Log-Zeile ändert der Diff nichts — kein Code, keine Konfiguration, kein Regelwerk, keine `spec/`-Datei. Das Anlegen von `slice-130.md`/`slice-131.md` und die Drift-Log-Zeile sind durch §2 Schritt 4 des Slice-Plans selbst gedeckt („Etappe C schneiden … mit Drift-Log-Eintrag") und entsprechen der Roadmap-eigenen Drift-Log-Definition („Slice oder Welle umgehängt oder neu geschnitten", `roadmap.md:125-126`) — beides ist Planung, keine Umsetzung. Kein Vorgriff auf slice-127 (nur ein mechanischer Pfad-Retarget infolge des Lifecycle-Moves von slice-129 selbst, in `a5b9d90`). Keine Sammel-Antwort: alle acht Wellen tragen eine eigene Zeile mit eigenem Beleg.
2. **„Acht Zeilen, keine Sammel-Zeile" — erfüllt.** Jede der acht Wellenzeilen (87–94) trägt eine individuelle, am jeweiligen Kurs-Diff verifizierte Begründung; keine der sechs „folgenlos"-Zeilen verweist pauschal auf eine andere.
3. **Welle-89-Konformitätsbeleg hält.** „Dieses Repo führt keine User Story" — `spec/` enthält ausschließlich `lastenheft.md`, `spezifikation.md`, `architecture.md` (kein `user-stories/`, keine Datei mit „story" im Namen). „`spec/spezifikation.md` deklariert ihr Stratum bereits über den Prozess" — der Kopf trägt wörtlich „**Rolle:** Technik-Stratum — fortschreibbar ohne Change Request" (`spec/spezifikation.md:7`).
4. **Welle-93-Konformitätsbeleg hält, real gegengeprüft.** `.d-check.yml:370-374` setzt `targets.authority: AGENTS.md` und `targets.doc-tables: [AGENTS.md, harness/README.md]` (beide Richtungen sind damit scharf; Kurs-Welle 93s eigene Korrektur zu Welle 92 zeigt, dass `gate-phantom`/`gate-undocumented` ohne `doc-tables` gar nicht liest). Eigener Lauf `make gate-consistency` isoliert: Exit 0, „460 Datei(en) geprüft, 0 Befund(e)". `AGENTS.md` §4 trägt eine erschöpfende Target-Tabelle (25 Zeilen), nicht nur die von Kurs-Welle 93 geforderten sieben Wurzel-Targets. Randnotiz ohne Befundcharakter: d-check erreicht das mit `exempt-targets: []` (vollständige Auflistung statt der von Welle 93 neu eingeführten namentlichen `doc-*`-Ausnahme) — ein anderer Mechanismus zum selben Ziel, kein Delta.
5. **Welle-87/88/91/92-„nicht anwendbar"/„konform"-Belege reproduziert.** `git diff v5.9.0..v5.11.0 --stat -- lab/regelwerk/` berührt exakt fünf Dateien (`README.md`, `grundlagen-begriffe.md`, `grundlagen-durchsetzungsschicht.md`, `grundlagen-referenz-richtung.md`, `grundlagen-source-precedence.md`); `README.md` (+3/−3 im Kurs-Repo, dort +1/−1 im vendorten Auszug) ändert ausschließlich die `**Stand:**`-Zeile. `modul-10-review-harness.md`, `modul-12-replay-evaluierung.md` stehen nicht im Diff (Welle 87/88 bestätigt). Wave-91/92-Commits (`e580c0e`, `725279a`) ändern ausschließlich `AGENTS.md`/`CLAUDE.md`/`.d-check.yml` des Kurs-Repos selbst, nicht `kurs/de/` und nicht `lab/example/`; beide Commit-Botschaften erklären „Regelwerk unberührt, Stand bleibt 90" bzw. sinngleich. Kein `team-sim`-Artefakt in diesem Repo (`find … -iname '*team-sim*'`: 0 Treffer, Welle 87 zusätzlich bestätigt).
6. **Zahlen der Commit-Botschaften, eigens nachgerechnet.** „fünf Regelwerks-Dateien" — 5 (bestätigt). „vier davon mit Regel-Inhalt" — 4 (`README.md` trägt nur die Stand-Zeile, die anderen vier je einen inhaltlichen Diff-Block). „460 Dateien, 0 Befunde" — `make gates` reproduziert exakt „460 Datei(en) geprüft, 0 Befund(e)" (targets- und planning-Modul, je isoliert). „acht Gates" — `make gates`-Ausgabe endet mit „`[gates] doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`", acht Namen. „sechs Wellen … belegt folgenlos" — Zeilen 87/88/89/91/92/93 tragen `nicht anwendbar` oder `konform, keine Handlung`, macht sechs; 90/94 tragen `Handlung nötig`, macht acht insgesamt. Alle Zahlen stimmen.
7. **Der grep-vs-Prüffrage-Nebenbefund ist reproduzierbar.** `grep -inE '\b(muss|müssen|darf|dürfen|soll|sollen|pflicht|verboten|verpflichtend|zwingend)\b' CLAUDE.md` liefert Exit 1, 0 Treffer — obwohl `CLAUDE.md` mit „Vor der Implementierung benennen: …" und „Bei Quellen-Konflikt: Konflikt melden …" zwei echte Vorschriften trägt (Imperativ-Infinitiv-Form ohne Modalverb). Die Folgerung „die Prüfung ist ein Urteil, kein `grep`" trifft — mit der in F-2 belegten Einschränkung — auch den in §4a geführten Zensus selbst: Er ist für die drei tatsächlich betrachteten Dateien ein Urteil (er liest ihren Inhalt und bewertet ihn einzeln), aber die Auswahl **welche** Dateien betrachtet werden, ist nirgends als Ergebnis einer systematischen Durchsuchung ausgewiesen — genau der Unterschied, den F-2 an zwei zusätzlichen Fundstellen zeigt.
8. **Welle-90-Kernaussagen (Lastenheft-Hälfte) halten.** `spec/lastenheft.md` Kopf trägt `**Status:** Draft` (nicht `Accepted`); die Historie-Tabelle (`spec/lastenheft.md` §7) führt `| Version | Datum | Änderung |`, drei Spalten, gegen `templates/spec/lastenheft.template.md:152` mit vier (`… | Verweis |`) — dieser Teil der Welle-90-Analyse ist korrekt (nur die Ausweitung auf „beide Straten" in F-3 nicht). Keine „zurückgezogene" Anforderung im Bestand (`grep -in zurückgezogen spec/lastenheft.md`: 0 Treffer) — slice-130 §3 nimmt das korrekt aus.
9. **Die in welle-83 vorab notierte adr-check/vcs-Vermutung zu Welle 90 hat sich beim Nachlesen des Kurs-Diffs nicht bestätigt** (`grundlagen-source-precedence.md`s Diff ist vollständig Lastenheft-/Spec-Stratifizierung, keine ADR-Immutabilität-Inhalte) — das Audit widerspricht dem nicht, erwähnt es aber auch nicht ausdrücklich als „geprüft und verworfen". Ohne Befundcharakter, da §1 des Slice-Plans die drei Vermutungen selbst als „Verdacht, nicht als Befund" einordnet und keine DoD-Pflicht daraus macht.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 4 | F-1, F-2, F-3, F-4 |
| LOW | 1 | F-5 |
| INFO | 0 | — |

## Verdikt

**Blockierend.** Kein HIGH, aber vier MEDIUM, die alle dieselbe Grundform teilen: eine im Audit **behauptete** Prüfung (Regelwerk als Prüfgröße, Vollständigkeits-Zensus, Spalten-Bestand beider Straten, Baseline-Beleg für einen Fachbegriff) hält bei eigenem Nachlesen am Ursprungsartefakt nicht in vollem Umfang, was der Fließtext suggeriert. Bemerkenswert: In drei von vier Fällen (F-1, F-2, F-4) bestätigt die tiefere Prüfung, dass die im Audit gezogene *Schlussfolgerung* in der Sache **richtig** bleibt (die sechs Wellen sind tatsächlich folgenlos; slice-131 zielt in der Sache auf die richtige Vollständigkeits-Frage) — der Mangel liegt im **Beleg**, nicht im Ergebnis. F-3 ist die Ausnahme: Dort trägt eine Aussage, die wörtlich in einen Folge-Slice übernommen wurde, einen Sachfehler, der bei unveränderter Umsetzung zu einer kanonwidrigen Änderung an `spec/spezifikation.md` führen würde.

Das Audit selbst benennt in §7 `BEO-011` als sein zentrales Risiko für genau die Vollständigkeits-Aussagen, die es produziert („acht Wellen abgedeckt", „alle normativen Artefakte gelistet") — F-2 zeigt, dass dieses selbst erkannte Risiko eingetreten ist: Der Zensus in §4a ist an genau der Stelle unvollständig, an der er Vollständigkeit für sich beansprucht. F-1 zeigt denselben Mechanismus auf der methodischen Ebene eine Stufe darunter (die Prüfgröße „Regelwerks-Diff" selbst ist ungeprüft gegen ihre eigene Nicht-Normativitäts-Erklärung übernommen).

Für die Closure: F-1/F-2 sind vor allem Beleg-Schärfung (Ergebnis bleibt), F-3 verlangt eine Korrektur an `slice-130` **vor** dessen Umsetzung (sonst entsteht dort ein neuer Sachfehler im Code/Doku-Diff), F-4 verlangt eine Korrektur des `Bezug:`-Felds in `slice-131`. F-5 ist kosmetisch und unabhängig entscheidbar.
