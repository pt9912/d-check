# Slice slice-286: Anforderung, Spezifikation und ADR für das Modul `pinsync`; Entscheide festschreiben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-94-pinsync-modul`](../welle-94-pinsync-modul.md).

**Bezug:** neue Anforderung des Moduls `pinsync` (die Kennung vergibt dieser Slice
beim Schreiben ins Lastenheft; Empfehlung: Bereichs-Kürzel `PSYNC`, Nummer 001),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(Scope), [`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei).

**Berührte Spec-Stellen:** `lastenheft.md` §3 (Bereichs-Kürzel-Liste und die neue
Anforderung), `spezifikation.md` §4 (Grund-Codes) und der Abschnitt der
Modul-Konfigurationen; `architecture.md` bleibt unberührt (Rolle „Regel-Modul“ besteht).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Anforderung des Moduls `pinsync` steht im Lastenheft, ihr Verfahren
und ihre Grund-Codes in der Spezifikation und die Entscheidung in einer
`Accepted`-ADR — und die Entscheide, die das Issue offen lässt, sind **vor dem
ersten Code** festgeschrieben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Konfigurations-Schema im Programm, Tests des Moduls** — gehören
  `slice-287`; dieser Slice schreibt, was der Code
  einlösen soll, und die Anforderung steht vor dem Code, damit die Reviews gegen sie
  prüfen können.
- **Der Selbstlauf im eigenen Profil** — gehört
  `slice-288`; er braucht das fertige Modul.
- **CHANGELOG, READMEs, Handbuch-Kopf** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17), kein Feature-Commit.
- **Aktualitätsprüfung und Normalisierung der Werte** (`strip-prefix`) — vom
  Auftraggeber ausgeschlossen (2026-10-31); die Anforderung nennt beides in ihrem
  Out-of-Scope, damit es nicht still hineinwächst.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Die Anforderung steht im Lastenheft (Bereichs-Kürzel in der Liste ergänzt,
      Version gehoben, Historie-Zeile nach
      [`MR-032`](../../../../harness/conventions.md#mr-032)) mit Akzeptanzkriterien
      je Befund-Klasse, den Rändern mit Exit 2 und ihrem Out-of-Scope.
- [ ] Die Spezifikation führt das Verfahren (Leseweise der Ausdrücke, Zählregel, Reihenfolge der
      Befunde), das Konfigurations-Schema und die Grund-Codes; jede Aussage über die Grenzen
      (Eingaben außerhalb der Scan-Menge, [`AGENTS.md`](../../../../AGENTS.md) §3.8) ist
      am Gegenstand gezählt oder gemessen, nicht aus der Beschreibung übernommen.
- [ ] Die ADR (Entscheide samt verglichenen Alternativen) ist `Accepted`, im
      [ADR-Index](../../adr/README.md) geführt und trägt `## Re-Evaluierungs-Trigger`.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | neue Anforderung, Kürzel in der Liste, Version und Historie |
| `spec/spezifikation.md` | update | Verfahren, Konfigurations-Schema, Grund-Codes |
| ADR im Verzeichnis `docs/plan/adr/` (nächste freie Nummer, Titel „pinsync-modul“) | neu | Entscheidung samt Alternativen |
| `docs/plan/adr/README.md` | update | ADR-Index ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 4) |

**Festzuschreibende Entscheide** (Empfehlung zuerst; der Auftraggeber bestätigt oder
ändert, bevor die Anforderung entsteht):

1. **Name und Kürzel:** Modul `pinsync`, Bereichs-Kürzel `PSYNC` (`PIN` und `VER` sind
   belegt: Content-Pin und Versions-Pin-Konsistenz gegen `version.md`).
2. **Leseweise:** die Datei wird zeilenweise gelesen (`\n`, ein `\r` am Zeilenende
   gehört nicht zur Zeile); der Ausdruck läuft je Zeile (RE2, `^`/`$` bedeuten
   Zeilenanfang und -ende), der Wert ist die Fanggruppe. Kein mehrzeiliger Modus.
3. **„Genau ein Treffer“:** gezählt werden **Zeilen mit Treffer**; null ⇒ Befund
   „kein Treffer“, mehr als eine ⇒ Befund „mehrdeutig“. Ein Treffer in einer Zeile genügt
   (kein Zählen mehrerer Vorkommen innerhalb einer Zeile).
4. **Eine Fanggruppe:** ein Ausdruck mit null oder mehr als einer Fanggruppe (benannte
   mitgezählt, nicht fangende `(?:…)` nicht) ist ein Eingabefehler (Exit 2). Eine Gruppe,
   die nichts fängt (`(\S*)` bei leerem Wert), liefert den leeren Wert — der ist ein Befund
   „kein Treffer“, kein Vergleichswert.
5. **Ein Pin mit nur einer Quelle** ist ein Konfigurationsfehler (Exit 2): ohne zweite
   Quelle gibt es nichts zu vergleichen, und ein stilles Grün wäre ein Gate ohne Prüfung.
6. **Pfade der Quellen:** relativ zur Scan-Wurzel, ohne `..`, ohne Glob, nicht absolut,
   wie `metrics.excluded` ([`DC-FA-MET-001`](../../../../spec/lastenheft.md#dc-fa-met-001--kennzahl-eines-planungsdokuments-gegen-messung-aus-git-modul-metrics-opt-in));
   die Dateien werden **nicht** von `scan.roots`/`scan.ignore` eingeschränkt — die
   Zusage gilt für diese benannten Eingaben, nicht für den Rest des Repos (§3.8). Ein
   Symlink als Quelle wird nach der Regel des Moduls `targets` behandelt; die Spezifikation
   nennt, welche das ist, nach Lesen des Codes.
7. **Reihenfolge der Befunde:** nach Pin-Name, dann nach Quelle in der Reihenfolge der
   Konfiguration (Determinismus, [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
8. **Befund-Ort:** Datei und Zeile der Quelle, soweit sie einen Treffer hat; bei
   „ungleiche Werte“ die Quellen mit ihren Werten in der Meldung.
9. **Grund-Codes:** `pinsync-mismatch`, `pinsync-no-match`, `pinsync-ambiguous`
   (Empfehlung; die Spezifikation führt die endgültige Liste und zählt sie am Code).

**Gegenprobe vor dem Schreiben (Schritt 19 des Workflow-Skeletts):** die Ausdrücke der
Anforderungsbeispiele laufen am eigenen Bestand; die Zählung steht beim Schreiben im
Plan-Nachtrag (`grep -cE '<ausdruck>' Dockerfile Makefile`, erwartet je 1) und die
Negativfälle (Kommentarzeile mit demselben Text, Wert in einem `FROM`, zwei Zeilen mit
Treffer, `v`-Präfix, CRLF) stehen als Testfälle für
`slice-287`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Welle `welle-94-pinsync-modul` ist offen und das
WIP-Limit ist frei (`in-progress/` trägt nur die Roadmap).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Entscheide 1–9
  brauchen mehr als eine Review-Sitzung, oder die Spezifikation verlangt einen
  Parse-Port, den die Zeilen-Leseweise nicht braucht.
- `in-progress` → `open` (blockiert): der Auftraggeber entscheidet einen der Punkte 1–6
  nicht; ohne Entscheid beginnt kein Code.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Lastenheft, Spezifikation und ADR liegen vor und `make gates` ist
grün), Review- und Verifikations-Report liegen unter `docs/reviews/`, Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Zeilen-Leseweise passt nicht zu einem Pin, der über mehrere Zeilen läuft (ein
  `ARG` mit Zeilenfortsetzung `\`) — **Ausgang:** *entfallen* oder *weiter offen*, nach der
  Gegenprobe am eigenen Bestand entschieden; die Anforderung nennt die Grenze
  ausdrücklich (Out-of-Scope: mehrzeilige Werte).
- Der Name `pinsync` kollidiert in Lexik und Handbuch mit `pins` und `versions` —
  **Ausgang:** *entfallen*, sobald die Anforderung das Abgrenzungs-Wort („Gleichheit
  mehrerer Dateien“) im ersten Satz trägt und das Reviewer-Urteil es bestätigt.
- Ein stilles Grün über einer Quelle, die ein Symlink auf eine Datei außerhalb des Repos
  ist (§3.8) — **Ausgang:** *eingetreten* nur, wenn die Spezifikation die Symlink-Regel
  nicht nennt; sonst *entfallen* mit Verweis auf die Regel.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln
(das Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde).

Entsteht bei der Closure.

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
Lastenheft, Spezifikation, ADR; keine Sub-Area mit eigener Modus-Deklaration
(`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Treffer:
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — die Quelldateien sind Eingaben
außerhalb der Scan-Menge, Entscheid 6),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(verkörpert in Schritt 19 des Workflow-Skeletts — die Negativliste in §3 ist seine
Anwendung),
[`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
(ohne formgültigen Ausgang; das Modul deckt davon nur die Wertgleichheit zweier Dateien und
löst den Eintrag **nicht** auf) und
[`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
(offen — jede Randbedingung der Anforderung braucht in
`slice-287` einen Test). Der Zähler folgt aus den Dateien
der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-10-31: `upstream-drift` grün (jüngster Lauf 2026-10-31T06:31:58Z),
`image-scan` grün (2026-10-31T09:41:43Z). Nichts zu lesen, nichts ausstehend.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
