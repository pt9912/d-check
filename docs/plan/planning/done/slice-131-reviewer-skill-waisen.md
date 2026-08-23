# Slice slice-131: Der Reviewer-Skill trägt Regeln, die nirgends gerankt stehen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-83-baseline-v5110-migration](welle-83-baseline-v5110-migration.md)
(Etappe C, geschnitten vom Delta-Audit).

**Bezug:** Baseline-Regelwerk
[`grundlagen-source-precedence.md`](../../../../.harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md)
§Vollständigkeit (Kurs-Welle 94) und
[`modul-09-implementierung.md` §AGENTS.md-Regeln](../../../../.harness/baseline/v5.11.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)
(*„Jede Hard Rule liegt in zwei Quadranten"* — **Durchsetzungs**-Verdopplung,
nicht Verortung; die Verortung regelt §Vollständigkeit);
[`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md);
[`AGENTS.md`](../../../../AGENTS.md) §3.

**Berührte Spec-Stellen:** — (Harness-Dateien; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der Vollständigkeits-Zensus aus [slice-129](../done/slice-129-baseline-v5110-delta-audit.md)
hat **fünf** Fundorte ergeben statt des einen bekannten. Der schwerste ist der
**Reviewer-Skill**: mehrere seiner Kategorien-Anker
sind **dort entstanden** („neuer HIGH-Eintrag seit 1.5.0", „neuer MEDIUM-Anker
seit 1.9.0"), statt einen in einer gerankten Quelle stehenden Ablauf
auszubuchstabieren.

Der Kanon erlaubt Artefakten außerhalb der Rangliste ausdrücklich, zu
**verweisen**, **auszuführen** und einen dort gerankten Ablauf
**auszubuchstabieren** — aber nichts **festzulegen**, was nicht dort steht. Ein
HIGH-Anker, der eine Prüf-Pflicht *einführt*, legt fest.

**Zwei kleinere Fundorte reiten mit**, weil sie dieselbe Frage stellen und je
zwei Zeilen sind: der `FOCUS_DISABLE`-Kommentar im
[`Makefile`](../../../../Makefile) (*„Spiegelt die `.d-check.yml`-modules-Liste;
wächst die dort, hier nachziehen"* — zugleich die Hälfte von
[`BEO-010`](../observations.md)) und die SHA-Pin-Konvention in den beiden
Workflow-Köpfen. Beide legen fest, beide stehen nirgends gerankt.

**Was dieser Slice nicht behauptet:** dass alle vierzehn normativ wirkenden
Stellen Waisen sind. Mehrere buchstabieren `AGENTS.md` §3.7 aus und sind damit
zulässig. Die Trennung ist Urteilsarbeit **je Anker** — genau die Prüffrage des
Kanons, und sie ist der Inhalt dieses Slice.

## 2. Vorgehen

1. **Je Kategorien-Anker eine Antwort:** buchstabiert er eine gerankte Regel aus
   (zulässig, Quelle nennen) oder legt er fest (Waise)? Ergebnis ist eine
   Tabelle, kein Urteil im Fließtext.
2. **Waisen umziehen, nicht löschen** — nach `AGENTS.md`, mit Herkunfts-Anker
   (`modul-09`: Hard Rules aus dem Steering Loop tragen ihn). Der Skill behält
   den Anker als **Verweis**.
3. Prüfen, ob dabei `AGENTS.md` zur Sammelstelle wächst. Der Wächter ist die
   **Prüffrage des Kanons** selbst — *steht die Aussage auch in einer gerankten
   Quelle?* —, denn genau dann ist ein Zuzug legitim: eine Waise steht nirgendwo
   sonst. Je Zuzug ist das zu belegen, nicht zu behaupten.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Rückbau von Ankern.** Was festlegt, wandert; es verschwindet nicht.
- **Kein Gate für die Anker.** Ob eine Review-Regel mechanisierbar ist, ist
  eine eigene Frage — und nach `modul-09` bleibt sie ohne Gate **halb
  durchgesetzt**, was hier zu **benennen** und nicht zu heilen ist.
- **Nicht `CLAUDE.md`** ([slice-127](../done/slice-127-claude-md-pointer.md))
  und nicht das Workflow-Skelett (im Audit als vermutlich konform eingestuft).

## 4. Definition of Done

- [x] Je Anker eine Antwort mit Quelle — die Tabelle in §9 deckt **alle** 18
      Aussagen, auch die Abschnitte, die der erste Anlauf nicht dokumentiert
      hatte (Anti-Pattern, Kontext-Eskalation, Output-Schema, Negativbefunde,
      Ablage, Eingangs-Kontext). **Eine Zuordnung war falsch** und ist
      berichtigt (§9).
- [x] Jede Waise steht in `AGENTS.md`; der Skill und die zwei Workflow-Köpfe
      verweisen. **Herkunfts-Anker nur, wo die Regel wirklich aus dem Steering
      Loop kam** (`seit welle-73`, `seit welle-82`) — §3.9 bekommt bewusst
      **keinen**, sie ist Bestand seit den Workflows selbst.
- [x] `AGENTS.md` ist nicht zur Sammelstelle geworden: für jeden der drei Zuzüge
      ist belegt, dass die Aussage im ganzen Repo nirgendwo sonst steht — der
      unabhängige Review hat das gegengeprüft. **Der Wächter, den §2 dafür nannte,
      war der falsche** (§9).
- [x] `make gates` Exit 0 (acht Gates, 466 Dateien, 0 Befunde); unabhängiger
      Review ([Report](../../../reviews/2026-08-23-slice-131-waisen-zensus-review.md)),
      blockierend, alle fünf Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Die Trennung „ausbuchstabieren vs. festlegen" ist Urteil, kein `grep`.**
  — **Ausgang:** *eingetreten, und zwar in die zu weite Richtung.* Zwei der
  vier vermuteten Waisen waren keine: die Skill-Anker in ihrer Mehrheit und der
  `FOCUS_DISABLE`-Kommentar. Beide Male hat die Nachprüfung am Kanon den
  eigenen Befund widerlegt, nicht bestätigt. Der zu enge Schnitt ist nicht
  eingetreten — der Review hat keine übersehene Waise gefunden.
- **Der Skill ist selbst das Werkzeug des Reviews.** — **Ausgang:**
  *eingetreten, benannt, nicht geheilt.* Der unabhängige Review lief gegen den
  Skill in seiner **neuen** Fassung 1.10.0 — das geänderte Instrument prüfte
  seine eigene Änderung. Tragbar war das nur, weil die Änderung an den zwei
  MEDIUM-Ankern reine **Zeiger** waren: keine Kategorie, keine Prüffrage, kein
  Anti-Pattern hat sich bewegt. Bei einer inhaltlichen Skill-Änderung wäre das
  keine gültige Prüfung mehr.

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-129](../done/slice-129-baseline-v5110-delta-audit.md)
in `done/`.

**Rückführungen:** `in-progress` → `next`, falls die Prüfung ergibt, dass die
Anker in ihrer Mehrheit ausbuchstabieren — dann ist es kein Umzugs-, sondern ein
Deklarations-Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Dateien (GF), Review-Infrastruktur (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **`BEO-011`**
  ist zentral — dieser Slice produziert eine Vollständigkeits-Aussage über die
  Anker, und genau diese Form ist in welle-82 achtmal gekippt. **`BEO-002`**
  für die Ränder jedes umgezogenen Ankers (Skill-Version, Index, Verweise).

Slice-ID: slice-131. Betroffene IDs: — (Harness-Dateien; keine Anforderung,
keine ADR, keine Adaption). Module: Harness-Dateien, Review-Infrastruktur. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Umzug nach kanonischer Ortsregel.

## 9. Closure-Notiz (nach `done/`)

Geliefert: ein Zensus über **18** Aussagen, **drei** Waisen nach `AGENTS.md`
umgezogen, ein eigener Befund widerlegt, und der Reviewer-Skill steht auf 1.10.0
und verweist, statt allein zu tragen.

**Die Prämisse dieses Slice war in ihrer Mehrheit falsch — und das ist sein
wertvollstes Ergebnis.** §1 nannte den Reviewer-Skill den *schwersten* Fundort,
weil seine Kategorien-Anker „dort entstanden" seien. Der Kanon sagt das
Gegenteil: `modul-10` §Ziel-Form **weist** dem Skill die repo-konkrete
Klassifikation ausdrücklich zu — *„die HIGH-Liste muss mindestens zwei
repo-spezifische Regeln nennen, die ein generischer Skill nicht abdeckt"*. Alle
sieben HIGH-Anker buchstabieren aus, einer davon nach `AGENTS.md` §3.7, das von
sich aus sagt: *„Der Reviewer-Skill trägt den HIGH-Anker dazu."* Wer den Slice
ohne diese Gegenprobe gefahren hätte, hätte einen kanonisch **vorgeschriebenen**
Inhalt als Regelverstoß behandelt.

**Das Kriterium, das die Waisen trennt — und woher es kommt.** Der Kanon stellt
nur die Prüffrage (*steht die Aussage auch in einer gerankten Quelle?*), nicht
das Kriterium. Es folgt aus zwei Sätzen zusammen: `grundlagen-source-precedence.md`
§Vollständigkeit verbietet Artefakten außerhalb der Rangliste das **Festlegen**,
und `modul-10` §Ziel-Form **weist** dem Skill bestimmte Blöcke als operative
Pflichtteile zu. Wo der Kanon einen Block zuweist, ist dessen Inhalt Ausführung
eines Auftrags; wo eine Aussage darüber hinaus **andere Rollen bindet** — den
Implementer, den Botschafts-Schreiber, jeden, der einen Workflow anfasst —,
steht sie außerhalb jedes Auftrags und braucht ein geranktes Zuhause. Diese
Ableitung stand im ersten Anlauf nicht da; der Review hat sie zu Recht als
unbelegt beanstandet — genau die `BEO-011`-Form, die §7 als Risiko benannt hatte.

**Der Zensus:**

| Aussage im Artefakt | Verdikt | Gerankte Fundstelle |
|---|---|---|
| HIGH · Stilles-Grün-Pfad in Gate/Gate-Skript | buchstabiert aus | `AGENTS.md` §4 (*„Halluzinierte Gates sind die häufigste Form von Harness-Lüge"*), §6.8 |
| HIGH · Korrektheitsfehler in Kern-Modulen | buchstabiert aus | Baseline `modul-10` §Finding-Kategorien (HIGH = Korrektheits-Verstoß) — **berichtigt**, siehe unten |
| HIGH · Verstoß gegen die [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)-Import-Regeln | buchstabiert aus | [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) (Rang 4) |
| HIGH · Gate-Suppression ohne ADR | buchstabiert aus | `AGENTS.md` §3.2, §3.6 |
| HIGH · Netzzugriff außerhalb `external` | buchstabiert aus | [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (Rang 1) |
| HIGH · Kommentar trägt keine der fünf Klassen | buchstabiert aus | `AGENTS.md` §3.7 — sagt selbst: *„Der Reviewer-Skill trägt den HIGH-Anker dazu."* |
| HIGH · Zustandsfeld trägt Chronik | buchstabiert aus | `AGENTS.md` §3.7 (Zustandsfelder), samt beider Ausnahmen |
| MEDIUM · Botschaft verallgemeinert über die Messung hinaus | **Waise ⇒ umgezogen** | neu: [`AGENTS.md`](../../../../AGENTS.md) §5 |
| MEDIUM · Spec-Treue-/Konsistenz-Lücke, Alt-Tool-Differenz, fehlende Negativtests | buchstabiert aus | [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools) (Rang 1); der Rest ist Review-Handwerk aus `modul-10` |
| MEDIUM · Referenz-Richtung — Marker-Ehrlichkeit | buchstabiert aus | [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix) (Rang 1) + Baseline §Referenz-Richtung |
| MEDIUM · Modul-Grenze auf der Ziel-Achse | **Waise ⇒ umgezogen** | neu: [`AGENTS.md`](../../../../AGENTS.md) §3.8 |
| MEDIUM · Adressierungs-Form eines Neuzugangs | buchstabiert aus | `docs/plan/adr/README.md` §Konventionen (Rang 4) + [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage) |
| LOW / INFO | buchstabiert aus | Baseline `modul-10` §Finding-Kategorien |
| Kontext-Eskalation (Gate-/Sicherheitspfad, dritte Wiederholung) | buchstabiert aus | Baseline `modul-10` §Pflege (Steering-Loop) |
| Anti-Pattern (fünf Punkte) | buchstabiert aus | Baseline `modul-10` §„Was dieser Skill NICHT macht", §Regeln gegen typische Fehlannahmen, Output-Schema-Feld `verifizierbar` |
| Output-Schema · Negativbefunde · Ablage | buchstabiert aus | Baseline `modul-10` §Ziel-Form, §Reviewer berichtet auch, was er nicht gefunden hat |
| Eingangs-Kontext (Pflicht-Eingaben) | buchstabiert aus, **mit einer Lücke** | Baseline `modul-10` §Ziel-Form — der Block nannte eine Pflicht-Eingabe **weniger** als der Kanon; nachgetragen |
| `Makefile` · `FOCUS_DISABLE` spiegelt die Modul-Liste | **keine Waise** | `AGENTS.md` §3.7 Klasse **Kopplung** — die Pflege-Folge *ist* die Kopplung |
| Workflow-Köpfe · alle `uses:` SHA-gepinnt | **Waise ⇒ umgezogen** | neu: [`AGENTS.md`](../../../../AGENTS.md) §3.9 |

**Berichtigt:** Die Botschaft des Feature-Commits ordnete „Korrektheitsfehler in
Kern-Modulen" der `AGENTS.md` §6 Schritt 8 zu. Das ist falsch — jener Schritt
regelt die **Ehrlichkeit des Agenten-Berichts**, nicht die Korrektheit des
Produkts. Die Fundstelle ist die Baseline. Beide bis dahin quellenlosen
HIGH-Anker tragen ihre Fundstelle jetzt **im Skill-Text**: damit zeigt die Datei
selbst, dass sie ausbuchstabiert, statt dass eine Commit-Botschaft es behauptet.

**Zum dritten Mal an einem Tag habe ich eine Quelle für mehr zitiert, als sie
sagt.** §2 nannte [`MR-015`](../../../../harness/conventions.md#mr-015) als
Wächter gegen „`AGENTS.md` wird zur Sammelstelle". Sein `Geltungsbereich` nennt
`AGENTS.md` **§1**, zwei Pin-Einträge und die Konventions-Quellen-Sektion — §3
und §5 stehen nicht darin, und sein Gegenstand ist die Baseline-Pointer-Drift.
Nur sein **Titel** klingt allgemein. Der Wächter ist in Wahrheit die Prüffrage
selbst: ein Zuzug ist legitim, **weil** die Aussage nirgendwo sonst steht — und
genau das ist für alle drei belegt und gegengeprüft.

**Ein Nebenfund in die Gegenrichtung.** Beim Nachprüfen der Zensus-Reichweite
fiel auf, dass der `Eingangs-Kontext`-Block **weniger** verlangte als der Kanon:
*„vorherige Findings am gleichen Modul"* fehlte. Das ist keine Waise, sondern
eine Lücke gegen die Baseline — die Vollständigkeits-Prüfung schaut auf beide
Richtungen, wenn man sie lässt. Nachgetragen, mit dem Grund: ohne sie sieht der
Reviewer denselben Fehler zum zweiten Mal als ersten.

**Offen und benannt:** Keine der drei neuen Regeln hat ein Gate. §3.8 und die
§5-Zeile sind Urteilsfragen und bleiben es (`modul-09`: halb durchgesetzt, hier
gesagt statt übertüncht). §3.9 ist der einzige der drei, der mechanisierbar
wäre, und trägt deshalb als einziger einen **auflösenden** Trigger statt
*permanent* — ein Sensor auf `uses:`-Pins löst sie ab.
