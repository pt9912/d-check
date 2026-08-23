# Slice slice-130: Historie-Form auf vier Spalten, und unsere Strenge deklarieren

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-83-baseline-v5110-migration](../welle-83-baseline-v5110-migration.md)
(Etappe C, geschnitten vom Delta-Audit).

**Bezug:** Baseline-Regelwerk
[`grundlagen-source-precedence.md`](../../../../.harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md)
(Kurs-Welle 90: CR-Pflicht am Lastenheft-Status, Tatsachenberichtigung,
zurückgezogene Anforderungen) und die Vorlage
`.harness/baseline/v5.11.0/templates/spec/lastenheft.template.md`;
[`spec/lastenheft.md`](../../../../spec/lastenheft.md) §7 Historie.

**Berührte Spec-Stellen:** `spec/lastenheft.md` §7 (Form der Historie-Tabelle),
Kopf-Feld `**Status:**` — keine Anforderung ändert ihre Aussage.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Kurs-Welle 90 bindet die **CR-Pflicht an den Lastenheft-Status** und gibt der
Historie eine vierte Spalte. Zwei Deltas gegen unseren Bestand, beide vom
Delta-Audit belegt:

- **Die Lastenheft-Vorlage trägt `| Version | Datum | Änderung | Verweis |`** —
  die vierte Spalte nennt den externen CR-Vorgang, bei einer
  **Tatsachenberichtigung** ein `—`. Unser Lastenheft führt drei.
  **Nur das Lastenheft:** `spec/spezifikation.md` §7 führt **zwei** Spalten und
  trifft damit die vendorte `spezifikation.template.md` exakt. Der Grund steht in
  [`grundlagen-referenz-richtung.md`](../../../../.harness/baseline/v5.11.0/regelwerk/grundlagen-referenz-richtung.md):
  eine `Verweis`-Spalte trägt nur, **was sonst nirgends im Repo steht** — beim
  Vertrag der externe CR, der kein anderes Zuhause hat; die Technik verankert
  ihre Aufwärts-Bezüge schon im Körper, und dieselbe Kopplung ein zweites Mal zu
  führen erzeugt keine Information, sondern eine zweite Fassung, die driftet.
- **Vor `Accepted` verlangt der Kanon *nichts*:** *„frei änderbar, ohne Change
  Request, ohne Historie-Zeile"*. Unser Lastenheft steht auf **`Draft`** — wir
  fahren Versions-Bumps und Historie-Zeilen also **strenger als verlangt**.
  Das ist eine legitime Wahl, aber sie ist heute **undeklariert**: ein Leser
  kann nicht unterscheiden, ob unsere Disziplin Pflicht oder Vorsatz ist.

## 2. Vorgehen

1. Die vierte Spalte **einführen — nur im Lastenheft**; Bestandszeilen
   bekommen `—`, weil es für sie keinen externen Vorgang gibt. Die
   Spezifikations-Historie bleibt bei zwei Spalten.
2. **Die eigene Strenge deklarieren** — als Adaption im Konventionsspeicher
   (`MR-`Eintrag) oder als Satz im Lastenheft-Kopf; **im Slice entscheiden**,
   nicht vorab. Kriterium: eine Form-Frage gehört in den Konventionsspeicher,
   eine Vertrags-Aussage ins Lastenheft.
3. Prüfen, ob die `structure`-Regel über `## 7. Historie` (Chronologie-Monotonie
   auf Spalte 1) von der neuen Spaltenzahl berührt ist.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Status-Wechsel des Lastenhefts.** Ob es von `Draft` auf `Accepted`
  geht, ist eine Auftraggeber-Entscheidung mit Vertragswirkung, kein Nachzug.
- **Kein Rückbau der bisherigen Strenge.** Der Kanon *erlaubt* weniger; er
  verbietet nicht mehr.
- **Keine Anwendung der Zurückgezogen-Regel** — dieses Repo hat keine
  entfallene Anforderung.

## 4. Definition of Done

- [x] Die **Lastenheft**-Historie trägt die vierte Spalte, alle 95 Bestandszeilen
      `—`; die Spezifikations-Historie ist **unverändert** bei zwei Spalten.
- [x] Die eigene Strenge ist deklariert — **an beiden Orten, und das war nicht
      die geplante Wahl**: eine Kopf-Notiz im Lastenheft (der Leser stellt die
      Frage dort) **und**
      [`MR-032`](../../../../harness/conventions.md#mr-032). Der erste Anlauf
      hielt den MR für entbehrlich; der Review hat das widerlegt, siehe §9.
- [x] Die `structure`-Regel auf `## 7. Historie` läuft unverändert grün — sie
      prüft die Monotonie über Spalte 1, und die Schlüsselspalte hat sich nicht
      bewegt. Ebenso ungerührt bleibt `matrix`: die neue Spalte steht
      durchgehend auf `—` und trägt kein Token.
- [x] `make gates` Exit 0 (acht Gates, 465 Dateien, 0 Befunde), `make fullbuild`
      Exit 0; unabhängiger Review
      ([Report](../../../reviews/2026-08-23-slice-130-historie-form-review.md)),
      blockierend, beide Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Eine Spalte in einer chronologischen Bestandstabelle zu ergänzen berührt
  die `table-order`-Regel.** — **Ausgang:** *nicht eingetreten, und geprüft
  statt angenommen.* Der Gate-Lauf über die geänderte Tabelle ist grün; die
  Regel liest Spalte 1, die unberührt blieb.
- **„Strenger als der Kanon" kann als Verstoß gelesen werden**, wenn die
  Deklaration fehlt oder am falschen Ort steht. — **Ausgang:** *eingetreten, in
  der schärferen Form.* Nicht der Ort war das Problem, sondern dass ich die
  Deklaration für ausreichend hielt, wo der Kanon einen Eintrag im
  Konventionsspeicher **zuweist**. Ohne ihn wäre die Abweichung für den
  Freshness-Audit unsichtbar geblieben — ein Leser hätte sie gefunden, die
  Maschine nicht.
- **Nachträglich aufgenommen, weil im Slice gefunden:** eine `Verweis`-Spalte
  wieder einzuführen berührt eine Accepted-ADR, die sie einst gestrichen hat.
  — **Ausgang:** *aufgelöst, ohne neue Entscheidung.* Siehe §9.

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-129](../done/slice-129-baseline-v5110-delta-audit.md)
in `done/`.

**Rückführungen:** `in-progress` → `next`, falls die Ortswahl für die
Deklaration eine Auftraggeber-Entscheidung verlangt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Spec-Straten (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **`BEO-002`**
  für die Ränder der Historie-Form — und ausdrücklich für die Frage, welches
  Stratum sie überhaupt trägt (die erste Fassung dieses Slice hat sie falsch
  beantwortet);
  **`BEO-011`** für jede Aussage darüber, welche Tabellen die Form tragen.

Slice-ID: slice-130. Betroffene IDs:
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
Spec-Straten, Konventionsspeicher. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Form-Nachzug an der adoptierten Vorlage.

## 9. Closure-Notiz (nach `done/`)

Geliefert: die Lastenheft-Historie trägt die kanonische vierte Spalte, alle 95
Bestandszeilen stehen auf `—`, die Spezifikations-Historie ist unberührt bei
zwei Spalten, und unsere Strenge ist deklariert — im Lastenheft-Kopf, wo der
Leser die Frage stellt, und in
[`MR-032`](../../../../harness/conventions.md#mr-032), wo die Maschine sie
findet.

**Der Slice hat seine eigene Regel im selben Commit gebrochen.** Der
Feature-Commit hat das Lastenheft geändert — vierte Spalte, Kopf-Notiz — **ohne
Versions-Bump und ohne Historie-Zeile**, während er in derselben Datei
deklarierte, dass dieses Repo beides seit der ersten Fassung führt. Kein Gate
hat das gemeldet, denn keines kann es: die Kopplung „Lastenheft geändert ⇒
Historie-Zeile" ist nirgends mechanisiert. Gefunden habe ich es beim Einarbeiten
der Review-Befunde, nicht beim Schreiben.

**Der Review hat einen Selbstwiderspruch im Abstand von zwei Stunden
gefunden.** Ich hatte den `MR-`Eintrag mit einer Abgrenzung abgelehnt — *„eine
Freiheit nicht zu nutzen ist keine Abweichung"*. [`MR-031`](../../../../harness/conventions.md#mr-031), am selben Tag
geschrieben, sagt wörtlich das Gegenteil: *wer verschärft, weicht ab, auch wenn
er nur mehr verlangt.* Und der Kanon **weist diese Entscheidung ausdrücklich
zu**: *„Welche Stelle der Version steigt, entscheidet das Repo und gehört in den
Adaptions-Block von `harness/conventions.md`."* Der Freshness-Audit liest genau
diese Liste — ohne Eintrag wäre die Abweichung für ihn unsichtbar geblieben. Ein
Leser hätte sie gefunden, die Maschine nicht.

**Der zweite Befund traf die Bedingung, unter der eine Berichtigung überhaupt
ohne Change Request auskommt.** Die Zeilen 0.65.1 und 0.65.2 waren
Tatsachenberichtigungen und trugen nur *„redaktionell"*. Der Kanon nennt zwei
Bedingungen, und *„als solche ausgewiesen"* ist die erste davon. Nachgetragen.

**Eine ADR-Frage, die im Slice aufkam und keine neue Entscheidung brauchte.**
[ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md) hat
2026-08-02 genau diese Spalte **gestrichen**, und
[`.d-check.yml`](../../../../.d-check.yml) begründete die daran hängende
`matrix`-Verengung mit dem Satz *„die Verweis-Spalte entfernt"*. Dieser Satz
wurde durch diesen Slice falsch. Die Auflösung liegt aber nicht in einer neuen
ADR, sondern im Baseline-Text selbst: er trägt für die Historie **dieselbe
Begründung** — eine superseded ADR in einer Protokollzeile ist unreparierbar —
und zieht daraus nur einen anderen Schluss über das Mittel. Die
Spalte darf existieren, weil sie **ausschließlich** den externen Change Request
aufnimmt; ADR- und Slice-Zeiger bleiben dort verboten. Substanz und Fitness
Function jener ADR stehen unverändert; die immutable ADR bleibt unberührt,
der Config-Kommentar trägt jetzt den Grund, der wirklich trägt.

**Und die dritte Berichtigung ging gegen mich selbst.** Meine Begründung dafür,
warum die Spezifikation bei zwei Spalten bleibt — *„sie kennt keinen
Change-Request-Vorgang"* — stammte aus meiner Konstruktion, nicht aus dem Kanon.
Der sagt etwas anderes und Präziseres: *eine Verweis-Spalte trägt nur, was sonst
nirgends im Repo steht*; die Technik verankert ihre Aufwärts-Bezüge bereits im
Körper, eine zweite Fassung in der Historie driftet nur. Die Aussage war die
`BEO-011`-Klasse in ihrer unangenehmsten Lage: in einer versionierten
Vertragsdatei, in genau der Zeile, die Korrektheit protokollierte. Berichtigt an
allen drei Fundstellen — die Zeile gehört diesem Slice, deshalb an Ort und
Stelle statt über eine weitere Berichtigungs-Zeile.

**Offen und benannt:** Die `Verweis`-Spalte bleibt durchgehend `—`, solange das
Lastenheft auf `Draft` steht — ohne begonnene CR-Pflicht gibt es keinen externen
Vorgang, den sie nennen könnte. [`MR-032`](../../../../harness/conventions.md#mr-032) trägt den Auflösungs-Trigger; der
Status-Wechsel selbst ist eine Auftraggeber-Entscheidung mit Vertragswirkung und
kein Nachzug.
