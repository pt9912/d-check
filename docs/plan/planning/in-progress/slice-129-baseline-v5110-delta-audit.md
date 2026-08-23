# Slice slice-129: Etappe B — Delta-Audit über acht Kurs-Wellen, je Welle eine Antwort

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-83-baseline-v5110-migration](../welle-83-baseline-v5110-migration.md)
(zugeordnet bei der Eröffnung).

**Bezug:** Baseline-Regelwerk
[`modul-02-harness-bootstrap.md` §Freshness-Audit](../../../../.harness/baseline/v5.11.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)
(darunter die **Bestands-Stichprobe, die auch bei aktuellem Pin läuft**);
[`AGENTS.md`](../../../../AGENTS.md) §1 (breiterer Pflicht-Blick beim
Drift-Audit); [slice-128](../done/slice-128-baseline-v5110-vendoring.md) (liefert den
Baum, gegen den geprüft wird).

**Berührte Spec-Stellen:** — (das Audit **liest**; was es findet, schneidet
eigene Slices).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Acht Kurs-Wellen (87–94) liegen zwischen unserem alten und dem neuen Pin. Dieses
Audit beantwortet sie **einzeln**: konform ohne Handlung (mit Beleg), Handlung
nötig (mit Slice), oder nicht anwendbar (mit Begründung). **Keine Welle ohne
Zeile** — eine übersprungene ist eine stille Annahme, und stille Annahmen sind
die Klasse, an der die Vorgänger-Welle achtmal gescheitert ist.

Drei Wellen sind vorab als wahrscheinlich folgenreich markiert — **als
Verdacht, nicht als Befund**:

- **Kurs-Welle 90** („ab `Accepted` zählt jede Zeile") berührt vermutlich
  `make adr-check` und das Modul `vcs` — also unsere ADR-Immutabilität.
- **Kurs-Welle 93** („AGENTS.md §4 wird die Autorität über die Targets")
  berührt vermutlich das Modul `targets`, `make gate-consistency` und unsere
  eigene §4-Tabelle.
- **Kurs-Welle 94** („eine Rangliste ordnet, jetzt deckt sie auch ab") bringt
  die **Vollständigkeits-Zusage**. Eine Verletzung ist bereits bekannt und hat
  ihren Slice ([slice-127](../next/slice-127-claude-md-pointer.md)); die Frage
  hier ist, ob es **weitere** gibt — Skill-Dateien, emittierte Fragmente,
  `.claude/commands/`.

Die übrigen fünf sind damit **nicht** als folgenlos erklärt. Sie bekommen
dieselbe Zeile wie die drei.

## 2. Vorgehen

1. **Je Kurs-Welle 87–94 eine Zeile** mit Antwort und Beleg. Quelle ist das
   Kurs-CHANGELOG **und** der Regelwerks-Diff, nicht nur die Überschrift — eine
   Wellen-Überschrift ist eine Zusammenfassung, kein Vertrag.
2. **Die Bestands-Stichprobe fahren**, die das Freshness-Audit auch bei
   aktuellem Pin verlangt — sie prüft, ob der gelebte Bestand dem Regelwerk
   entspricht, nicht nur ob der Pin stimmt.
3. **Die Vollständigkeits-Zusage auf das eigene Repo anwenden** (lesend):
   welche Artefakte außerhalb der Rangliste tragen normativen Text? Die
   Prüffrage lautet *„Steht jede Aussage dieser Datei auch in einer gerankten
   Quelle?"* — Ergebnis ist eine **Liste**, keine Bereinigung.
4. **Etappe C schneiden:** aus den Handlungs-Zeilen werden Slices, mit
   Drift-Log-Eintrag in der Roadmap. Was zu groß ist, wird als eigene Welle
   ausgewiesen statt angehängt.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Umsetzung.** Das Audit liest und schneidet; es ändert weder Code noch
  Konfiguration noch Doku außerhalb seines eigenen Ergebnisses.
- **Kein Vorgriff auf [slice-127](../next/slice-127-claude-md-pointer.md).**
- **Keine Sammel-Antwort.** „Die übrigen fünf sind folgenlos" ist genau die
  Aussage, die dieses Audit verbietet.

## 4. Definition of Done

- [ ] Acht Zeilen, eine je Kurs-Welle 87–94, jede mit Antwort **und** Beleg;
      keine Sammel-Zeile.
- [ ] Bestands-Stichprobe gefahren und ihr Ergebnis benannt.
- [ ] Liste der Artefakte außerhalb der Rangliste, die normativen Text tragen —
      vollständig **belegt**, nicht behauptet (`BEO-011`).
- [ ] Etappe C geschnitten oder als Folge-Welle ausgewiesen; Drift-Log-Zeile
      gesetzt.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 4a. Ergebnis: acht Kurs-Wellen, acht Antworten

Quelle je Zeile ist **CHANGELOG und Regelwerks-Diff** (`git diff v5.9.0..v5.11.0
-- lab/regelwerk/`), nicht die Wellen-Überschrift. Der Diff berührt **fünf**
Regelwerks-Dateien; vier davon tragen Regel-Inhalt, die fünfte ist die
Stand-Zeile.

**Die Prämisse dahinter ist nachgetragen, nicht vorausgesetzt** — sie stand in
der ersten Fassung dieses Audits unbelegt da. Das vendorte Regelwerk sagt über
sich selbst: *„Es trägt **keine eigene Normativität**: maßgeblich für den
Inhalt bleibt der Kurs"*
([`regelwerk/README.md`](../../../../.harness/baseline/v5.11.0/regelwerk/README.md)).
„Kein Regelwerks-Diff" belegt also **nicht** aus sich heraus „folgenlos" — es
belegt es nur, wenn der **Kurs** für dieselben Wellen nichts trägt, was das
Regelwerk nicht abbildet. Für die Wellen 87–94 ist das unabhängig
gegengeprüft: `kurs/de/` und `lab/example/` folgen dem Regelwerks-Diff exakt.
Die Antwort bleibt, ihre Begründung ist jetzt vollständig.

| Kurs-Welle | Antwort | Beleg |
|---|---|---|
| 87 — Team-Sim in Modul-12-Form | **nicht anwendbar** | Regelwerks-Diff berührt `modul-12` nicht; dieses Repo führt keine team-sim |
| 88 — Vier unbelegte Aussagen, sieben Verdikte | **konform, keine Handlung** | `modul-10-review-harness.md` steht **nicht** im Diff — die Welle änderte Kurs-Inhalt, keine Regel |
| 89 — Form ist kein Beleg | **konform, keine Handlung** | `grundlagen-referenz-richtung.md` +29: Genre-Falle und die zwei Achsen (Gehalt · **Änderungs-Prozess**, letzterer entscheidet). Probe am Fall ist die User Story — dieses Repo führt **keine** (`spec/` trägt genau die drei Straten), und `spec/spezifikation.md` deklariert ihr Stratum bereits über den **Prozess** („fortschreibbar ohne Change Request") |
| 90 — Ab `Accepted` zählt jede Zeile | **Handlung nötig** | `grundlagen-source-precedence.md`: die CR-Pflicht hängt am **Lastenheft-Status**; die **Lastenheft**-Vorlage trägt jetzt vier Spalten (`Verweis`), unsere drei. Die **Spezifikation** ist davon **nicht** betroffen — siehe unten |
| 91 — Das Kurs-Repo sagt, wie an ihm gearbeitet wird | **nicht anwendbar** | Kurs-internes Repo-Briefing; die Welle selbst erklärt „Regelwerk unberuehrt". Trägt aber einen **Beleg** für uns: der Kurs begründet dort seine eigene Wurzel-Einstiegsdatei als „Werkzeug-Verkabelung, kein Harness-Konstrukt" |
| 92 — Zwei Gewohnheiten werden Invarianten | **nicht anwendbar** | Kurs-eigene `structure`-Regeln in dessen Prüf-Profil; kein Regelwerks-Diff |
| 93 — AGENTS.md §4 wird die Autorität über die Targets | **konform, keine Handlung** | kein Regelwerks-Diff; und dieses Repo fährt die Form bereits: `targets` mit `authority: AGENTS.md`, beide Richtungen (`gate-phantom`/`gate-undocumented`) |
| 94 — Eine Rangliste ordnet, jetzt deckt sie auch ab | **Handlung nötig** | `grundlagen-source-precedence.md` +77 (Vollständigkeits-Zusage, Prüffrage, Aufräum-Reihenfolge) und `grundlagen-durchsetzungsschicht.md` +7 (Rolle der Wurzel-Einstiegsdatei). Auslöser war ein Konsumenten-CR dieses Repos. Siehe unten |

### Welle 90 — was zu tun ist

Zwei Deltas, beide klein, beide belegt:

- **Die Lastenheft-Historie-Vorlage trägt eine vierte Spalte `Verweis`** (für
  den externen CR-Vorgang; bei einer Tatsachenberichtigung `—`). Unser
  Lastenheft führt drei.
  **Nach Review korrigiert — die erste Fassung sagte „beide Straten führen drei
  Spalten", und das war falsch:** `spec/spezifikation.md` §7 führt **zwei**
  (`Datum | Änderung`) und trifft damit die vendorte
  `spezifikation.template.md` **exakt**. Der Kanon gibt dem Technik-Stratum
  bewusst **keine** `Verweis`-Spalte — es kennt keinen Change-Request-Vorgang,
  den sie nennen könnte. Wer die Spalte dort ergänzte, führte eine Form ein,
  die der Kanon für dieses Stratum gerade **nicht** vorsieht.
- **Die CR-Pflicht beginnt erst ab Lastenheft-Status `Accepted`.** Unser
  Lastenheft steht auf **`Draft`** — vor `Accepted` ist es laut Kanon *„frei
  änderbar, ohne Change Request, ohne Historie-Zeile"*. Wir fahren also
  **strenger als der Kanon verlangt**. Das ist kein Verstoß, aber es ist heute
  **undeklariert**: ein Leser kann nicht unterscheiden, ob unsere
  Historie-Disziplin Pflicht oder Wahl ist.

### Welle 94 — der Vollständigkeits-Zensus

Der Kanon nennt die Prüfung ausdrücklich eine **Prüffrage**, kein `grep` — und
das ist keine Feinheit: ein Muster-Zensus über normative Wörter meldet für
`CLAUDE.md` **null** Treffer, obwohl dort **zwei** Regeln stehen, die es
nirgends sonst gibt (in slice-127 einzeln belegt). Die Prüfung ist ein Urteil.
Als Urteil geführt, ergibt sie **fünf** Fundorte statt des einen bekannten — drei davon fand die
erste Fassung, zwei erst der Review:

| Artefakt | Befund | Einordnung |
|---|---|---|
| `CLAUDE.md` | **zwei Waisen** — Meldepflicht bei Quellen-Konflikt, Benenn-Pflicht vor der Implementierung | bekannt, hat seinen Slice ([slice-127](../next/slice-127-claude-md-pointer.md)) |
| [`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md) | **zu prüfen, vermutlich mehrere.** Die Kategorien-Anker sind Regeln, denen der Reviewer folgt; mehrere sind dort *entstanden* („neuer HIGH-Eintrag seit 1.5.0", „neuer MEDIUM-Anker seit 1.9.0") statt einen gerankten Ablauf auszubuchstabieren | **neu** — eigener Slice |
| `.claude/commands/implement-slice.md` | **Grenzfall.** „Do not skip gates" und „Do not claim completion without command output" sind Regel-Sätze; der Kanon erlaubt dem Workflow-Skelett aber ausdrücklich, einen gerankten Ablauf **auszubuchstabieren** — und beide Sätze haben ihr Original in AGENTS.md §6 | **vermutlich konform**, im Slice zu entscheiden |

| [`Makefile`](../../../../Makefile) (Kommentar über `FOCUS_DISABLE`) | **Waise.** *„Spiegelt die `.d-check.yml`-modules-Liste; wächst die dort, hier nachziehen"* — eine Pflege-Pflicht, die nirgends gerankt steht. Sie ist zugleich die Hälfte von `BEO-010`, dort als gate-blind benannt | **neu** |
| `.github/workflows/{ci,release}.yml` (Kopf-Kommentare) | **Waise.** *„`uses:`-Einträge sind SHA-gepinnt mit Tag-Kommentar"* — eine Konvention ohne Rückhalt in `AGENTS.md`, Konventionsspeicher oder ADR | **neu** |

`.githooks/*` und die Hooks unter `.claude/hooks/` tragen **Durchsetzung**, kein
Regel-Original — sie stehen außerhalb dieser Frage.

**Der Zensus war in seiner ersten Fassung selbst unvollständig**, und zwar auf
die Weise, vor der [`BEO-011`](../observations.md) warnt: er nannte **drei**
Fundorte, hatte aber `Makefile`-Kommentare, `.github/workflows/` und
`tools/*.sh` nie durchsucht. Der Review hat zwei weitere gefunden. Das ist eine
**frische, ungezählte Instanz** der Klasse, die dieser Slice in §7 als sein
eigenes zentrales Risiko benennt — und sie steht hier, weil ein Zensus, der
seine eigene Lücke verschweigt, schlechter ist als keiner.

**Damit ist die Reihenfolge-Frage aus [slice-127](../next/slice-127-claude-md-pointer.md)
beantwortet:** `CLAUDE.md` ist **nicht** der einzige Fall. Wer nur ihn erledigt,
schließt einen von **fünf** Fundorten und lässt den größeren offen.

## 5. Abnahme-Punkte / Risiken

- **Ein Audit über acht Wellen verführt zur Sammel-Antwort.** Je länger die
  Liste, desto attraktiver „betrifft uns nicht". Genau dagegen steht die
  Ein-Zeile-je-Welle-Regel. — **Ausgang:** *(bei Closure)*
- **Die drei Verdachts-Wellen könnten das Audit dominieren** und die fünf
  übrigen zur Formsache machen. — **Ausgang:** *(bei Closure)*
- **Die Vollständigkeits-Liste ist selbst eine Vollständigkeits-Aussage** — und
  damit die Form, die in welle-82 achtmal gekippt ist. Sie braucht denselben
  zeilenweisen Beleg, den sie einfordert. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-128](../done/slice-128-baseline-v5110-vendoring.md)
in `done/` — das Audit liest den **neuen** Baum.

**Rückführungen:** `in-progress` → `next`, falls das Audit eine Produkt-
Konsequenz findet, die vor dem Rest des Audits gehört (etwa ein rotes Gate).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Konventionen (GF), Gate-Mechanik (GF),
  Nutzer-Doku (GF, nur lesend).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **`BEO-011`**
  ist die zentrale — dieses Audit produziert zwei Vollständigkeits-Aussagen
  (acht Wellen abgedeckt; alle normativen Artefakte gelistet), und genau solche
  Aussagen sind die Klasse. **`BEO-002`** für die Ränder jeder Regel, die das
  Audit als „konform" abhakt. **`BEO-009`** für jede Zahl im Ergebnis.

Slice-ID: slice-129. Betroffene IDs:
[`MR-021`](../../../../harness/conventions.md#mr-021). Module:
Harness-Konventionen, Gate-Mechanik. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Audit nach der Form der v5.6.0-Migration,
die dieselbe Aufgabe über sechs Stufen gelöst hat.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
