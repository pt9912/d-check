# Slice slice-282: Anforderung, Spezifikation und ADR für das Modul `metrics`; Bestands-Gegenprobe

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-93

**Bezug:** [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) (Determinismus der Formel und der Messung),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (kein Schreiben, kein Netz). Die Anforderung des Moduls
entsteht in diesem Slice und trägt erst danach eine Kennung; Vorbild für Form und
Reihenfolge sind [`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) und
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in). Eine begleitende ADR hält die Entscheide fest (Nummer bei der Anlage).

**Berührte Spec-Stellen:** `spec/lastenheft.md` (neue Anforderung, §3
Bereichskürzel-Liste, Historie) · `spec/spezifikation.md §2` (Schema
`metrics.*`) und §7 (Grund-Codes) · `spec/architecture.md` ([`ARC-005`](../../../../spec/architecture.md),
[`ARC-010`](../../../../spec/architecture.md) und die Fehlerzeile zum Adapter) — die Kennungen entstehen mit dem Slice.
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-30.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Anforderung des Moduls `metrics` steht im Lastenheft, ihre
Schnittstelle in der Spezifikation und ihre Entscheide in einer ADR — **bevor**
Code entsteht, und gegen die Pläne des Anlass-Repos gemessen statt aus dem Issue
geraten.

**Die acht Entscheide, die das Issue offenlässt** und die dieser Slice setzt (mit
Empfehlung; die Entscheide trifft der Auftraggeber, vor dem Schreiben der
Anforderung):

1. **Syntax des Blocks im Plan.** Empfehlung: ein Fenced-Block, dessen
   Info-String der konfigurierte Blockname ist, mit **einfacher Zeilen-Syntax**
   `Feld: Ganzzahl` bzw. `Feld: Text` je Zeile — **kein YAML** (siehe Entscheid 6);
   die Feldnamen sind ebenfalls Konfiguration (Issue: „Blockname … Feldnamen“).
   Maßgeblich ist die Form, die die Pläne des Anlass-Repos **schon haben** — sie wird
   gelesen, nicht erfunden. Trägt die Zeilen-Syntax sie nicht (Verschachtelung,
   Listen), fällt Entscheid 6 auf den Parse-Port.
2. **Woher die Basis `<basis>` kommt.** Empfehlung: ein Feld des Blocks (der Plan
   kennt seine Basis), `HEAD` ist fest. Nicht `--range` auf der Kommandozeile: ein
   Lauf prüft mehrere Pläne mit je eigener Basis.
3. **Form der Formel.** Empfehlung: linear, `Obergrenze = Konstante + Σ Faktor ·
   Eingabe`, mit **ganzzahliger** Arithmetik (feste Rundungsregel) — kein
   Ausdrucks-Parser, kein Gleitkomma ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)). Ob die Fachlogik des Anlass-Repos
   („Zustände“ im Issue) eine Verzweigung nach einem Zustandsfeld verlangt, sagt die
   Gegenprobe; wenn ja, ist sie Teil der Anforderung, sonst ausdrücklich nicht.
4. **Was `git diff --numstat` hier heißt.** Zu setzen: Einfügungen plus Löschungen
   je Datei, Pathspec-Semantik der Block-Pfade, Umgang mit Binärdateien (`-` in
   `numstat`) und mit Renames (Empfehlung: Rename-Erkennung aus, Messung also wie
   `--no-renames`; der Wert ist in der Spezifikation benannt). Die Messung läuft
   über den VCS-Port ohne git-Binary — die Gleichheit mit dem Kommando ist eine
   **Behauptung, die [`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) am Fixture belegt**, nicht dieser Slice.
5. **Grund-Codes und der Ausnahme-Eintrag.** Vorschlag: `metrics-limit-mismatch`,
   `metrics-over-limit`, `metrics-reason-empty`, `metrics-block-invalid`; ein
   Ausnahme-Eintrag mit Grund bei einer Messung **im** Rahmen ist kein Befund.
6. **Wo der Block gelesen wird.** Der Kern trägt keine YAML-Bibliothek
   ([`ARC-002`](../../../../spec/architecture.md), Verbotsspalte; `yaml.v3` ist in `.a-check.yml` auf die Adapter
   `configyaml`, `report`, `workflowyaml` beschränkt). Zwei Wege: (a) **Zeilen-Syntax**
   im Kern, kein YAML (Empfehlung — kein neuer Port, drei Slices bleiben); (b) ein
   **Parse-Port** nach dem Muster von `WorkflowParser`
   ([`ADR-0072`](../../adr/0072-workflows-modul.md)), dann wächst
   [`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) um Port und Adapter, und die Block-Lexik wird vorab ein eigener Slice.
7. **Fence-Behandlung.** Die geteilte Vorverarbeitung entfernt Fenced-Code für alle
   Module; das Modul braucht eine **auf `metrics` begrenzte** Ausnahme, die den Block
   liest, ohne das Verhalten anderer Module zu verschieben. Ein **ungeschlossener**
   Fence ist `metrics-block-invalid`, nicht still „ohne Block“.
8. **Wo die Zeilen gezählt werden.** Empfehlung: neue Methode am VCS-Port (der
   Adapter besitzt den Diff). Die Alternative — der Kern zählt selbst über `AllPaths`
   und `FileAt` — braucht keinen Port-Zuwachs, aber einen Zeilen-Diff im Kern; die
   ADR verwirft eine der beiden mit Begründung, und `spec/architecture.md`
   ([`ARC-005`](../../../../spec/architecture.md), [`ARC-010`](../../../../spec/architecture.md), Fehlerzeile) wird in diesem Slice entsprechend geschärft.

**Die Frage aus [`AGENTS.md`](../../../../AGENTS.md) §3.8 wird in der Anforderung beantwortet, nicht erst im Code:**
das Modul liest Eingaben, die es nie scannt — die Pfade des Blocks (außerhalb
der Scan-Wurzeln, nicht Markdown) und die git-Revisionen. Dort gilt die Zusage
der Scan-Menge nicht; die Grenze gehört in die Anforderung (wie bei
[`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in): „liest zwei Eingaben, scannt nur eine“).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Konfigurations-Leser, Gerüst, Spiegel der Modul-Registrierung** — Schicht-
  Abgrenzung: dieser Slice ist Spec und ADR, „kein Produkt-Code“ ist die wirksamste
  Selbstbindung ([`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) trägt die Registrierung).
- **Der Modulsatz in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) und die Glossar-Zeile „Regelmodul“** — der
  Spiegel-Test (`registry_mirror_test.go`) hält diese Orte in **beide** Richtungen
  gegen `model.ValidModules()`; stünde `metrics` hier schon drin, wäre `make gates`
  ab diesem Slice rot. Beide Stellen schreibt
  [`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) mit der Registrierung.
- **Der VCS-Port, die Messung und die Prüfung gegen die Obergrenze** —
  [`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) übernimmt sie; die Anforderung beschreibt sie hier schon vollständig,
  damit der Kern nicht ohne die Messung entworfen wird.
- **Die Antwort an das Issue und die Release-Prep** — ein anderer Vorgang
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17), außerdem Handlung des Auftraggebers.
- **Weitere Messgrößen und schreibende Wege** — eine neue Anforderung bzw. durch
  [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) ausgeschlossen; die Out-of-Scope-Liste der Anforderung nennt sie.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Lastenheft:** die Anforderung des Moduls mit Beschreibung (inkl. der
      §3.8-Grenze), Out-of-Scope, Akzeptanzkriterien, Bereichskürzel in §3, Version und
      Historie-Zeile ([`MR-032`](../../../../harness/conventions.md#mr-032)); die acht Entscheide aus §1 sind gesetzt.
      **Ohne** Modulsatz in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) und ohne Glossar-Zeile (§1).
- [ ] **Spezifikation:** Schema `metrics.*` (§2), Block-Syntax, Formel samt
      Rundungsregel, Messsemantik der Summe, Grund-Codes (§7) — jeder neue Grund-Code
      mit `--doctor`-Klartext-Vorgabe; dazu die Schärfung von `spec/architecture.md`
      ([`ARC-005`](../../../../spec/architecture.md), [`ARC-010`](../../../../spec/architecture.md), Fehlerzeile) für den Zuwachs am VCS-Port (Entscheid 8).
- [ ] **ADR** mit Alternativen-Tabelle zu den Entscheiden 1–8, Re-Evaluierungs-Trigger,
      `Schärft:`-Feld; ADR-Index nachgezogen ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeilen 4/5).
      **Status `Proposed` bis zur Closure von [`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle)** — dort wird die `numstat`-Zusage
      belegt oder korrigiert; erst danach `Accepted`. Sonst bräuchte eine Korrektur
      eine Folge-ADR ([`AGENTS.md`](../../../../AGENTS.md) §3.5).
- [ ] Bestands-Gegenprobe am Anlass-Repo dokumentiert (§3): Block-Syntax, Formel und
      die 16 Fälle aus dem Issue sind **gelesen und in der Anforderung belegt**, nicht
      nacherzählt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Benutzerhandbuch und `operations.md` **nicht** (Release-Prep, §1).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | neue Anforderung des Moduls; §3 Bereichskürzel-Liste ([`MR-025`](../../../../harness/conventions.md#mr-025): Spiegel vor dem Editieren auflisten); Version, Historie-Zeile. **Nicht:** Modulsatz in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl), Glossar-Zeile (→ [`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle)) |
| `spec/spezifikation.md` | update | Schema-Zeilen `metrics.*`, Algorithmus-Abschnitt zur Anforderung, Grund-Codes, Historie |
| `spec/architecture.md` | update | [`ARC-005`](../../../../spec/architecture.md), [`ARC-010`](../../../../spec/architecture.md) und Fehlerzeile nennen `metrics` und die Zeilen-Zählung am Adapter; ohne Modul-Pfade ([`MR-033`](../../../../harness/conventions.md#mr-033)) |
| `docs/plan/adr/` · `docs/plan/adr/README.md` | neu · update | ADR zu Entscheid 1–8 (Alternativen-Tabelle), Status `Proposed`, Index nachgezogen |

**Bestands-Gegenprobe vor dem Schreiben (nicht danach):**

- Die Pläne des Anlass-Repos lesen: Block-Syntax, Feldnamen, die Zahl der Fälle
  (Issue: 16 von 16 überschritten), und welche **Zustände** vorkommen. Ergebnis ist
  eine Tabelle Plan → Eingaben → Marke → Messung → Urteil, die der Slice-Closure-
  Notiz anhängt; sie ist zugleich die Vorlage der Paritäts-Probe in
  [`welle-93`](../welle-93-metrics-modul.md) §3.
- Der eigene Bestand dieses Repos trägt keinen solchen Block; die Probe „kein Block ⇒
  kein Befund“ folgt in [`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle).
- Das Modul zählt über den VCS-Port; ob die Summe mit `git diff --numstat
  --no-renames` übereinstimmt, ist **nicht** hier zu messen, sondern in
  [`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) am Fixture — hier nur als Zusage formuliert und als solche gekennzeichnet.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): welle-93 ist offen, das WIP-Limit ist frei, und
der Auftraggeber hat die offenen Entscheide aus [`welle-93`](../welle-93-metrics-modul.md) §5
beantwortet oder die Empfehlungen aus §1 (Entscheide 1–8) bestätigt, **und** das
Anlass-Repo ist vom Auftraggeber **benannt** (Name und Stand) und seine Pläne liegen
zum Lesen vor. Der Übergang landet auf dem Hauptzweig, vor der Arbeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Gegenprobe zeigt eine
  Fachlogik jenseits einer linearen Formel (Verzweigung nach Zustand, mehrere
  Messgrößen) — dann gehört die Formel-Sprache in einen eigenen Slice vor diesem.
  Ebenso, wenn die Block-Form des Anlass-Repos echtes YAML verlangt und Entscheid 6
  auf den Parse-Port fällt: dann ist die Block-Lexik ein eigener Slice davor.
- `in-progress` → `open` (blockiert — Carveout?): die Pläne des Anlass-Repos sind
  nicht verfügbar, oder der Auftraggeber entscheidet die Basis (Entscheid 2) so, dass
  sie nicht aus dem Plan stammt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (darunter `doc-check`, `adr-check`,
`gate-consistency`), Review-Report liegt vor, und die Gegenprobe-Tabelle (§3) steht
in der Closure-Notiz. Dazu die Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Fachlogik des Anlass-Repos passt nicht in eine lineare Formel; eine zu enge
  Anforderung zwingt das Anlass-Repo, sein Programm zu behalten, und das Modul wäre
  ohne Konsumenten — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Die Anforderung beschreibt die Messung per `numstat`-Gleichheit, bevor der Adapter
  steht; ein Unterschied zwischen reiner-Go-Diff und git-Kommando (Zeilenenden,
  Binärdateien, große Dateien) würde sie nachträglich ändern. Die Zusage ist
  deshalb als noch nicht belegt gekennzeichnet — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Die Pläne des Anlass-Repos ändern sich zwischen Gegenprobe und Paritäts-Probe;
  die Tabelle misst dann einen alten Stand — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) —
Spec und ADR; keine Sub-Area mit eigener Modus-Deklaration.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Treffer:
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(3×, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — das Modul liest Pfade und
git-Revisionen, die es nicht scannt; die Grenze steht in der Anforderung, §1) und
[`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(17×, gemischt — die Gegenprobe misst die Pläne des Anlass-Repos, nicht den Eigenbestand,
der keinen Block trägt; die Anforderung sagt, was am Fremdbestand gemessen ist und
was nur behauptet). Der Zähler folgt aus den Dateien der Closure.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
