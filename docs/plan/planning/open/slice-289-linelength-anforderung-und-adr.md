# Slice slice-289: Anforderung, Spezifikation und ADR für das Modul `linelength`; Entscheide festschreiben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-95-linelength-modul`](../welle-95-linelength-modul.md).

**Bezug:** neue Anforderung des Moduls `linelength` (die Kennung vergibt dieser Slice
beim Schreiben ins Lastenheft; Empfehlung: Bereichs-Kürzel `LLEN`, Nummer 001),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(Scope), [`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
[`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope);
Abgrenzung gegen [`DC-FA-FILE-001`](../../../../spec/lastenheft.md#dc-fa-file-001--zeilen--und-byte-obergrenzen-einer-ganzen-datei-modul-file-opt-in).

**Berührte Spec-Stellen:** `lastenheft.md` §3 (Bereichs-Kürzel-Liste und die neue
Anforderung), `spezifikation.md` §1 (neues Verfahren), §2 (Schema `linelength.*`),
§4 (Grund-Code — die Tabellenzeile kommt mit dem Code); `architecture.md` bleibt
unberührt (die Rolle „Regel-Modul“ besteht). Der Verweis zeigt **aufwärts**: Die Spec
nennt diesen Slice nie.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-11-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Anforderung des Moduls `linelength` steht im Lastenheft, ihr Verfahren,
ihr Schema und ihr Grund-Code in der Spezifikation und die Entscheidung in einer
`Accepted`-ADR — und die Entscheide unten sind **vor dem ersten Code** vom
Auftraggeber bestätigt und festgeschrieben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Konfigurations-Rand, Tests des Moduls** — gehören
  [`slice-290`](../open/slice-290-linelength-modulkern.md); dieser Slice schreibt, was der Code
  einlösen soll, damit die Reviews dort gegen eine Anforderung prüfen.
- **Die Zeile in der Grund-Code-Tabelle (§4) und die Modul-Aufzählung in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) und
  im Glossar** — kommen mit dem Code in
  [`slice-290`](../open/slice-290-linelength-modulkern.md): die Spiegel-Tests halten beide gegen
  das registrierte Modul und wären vorher rot (bei `pinsync` gemessen).
- **Andere Dateiarten** — vom Auftraggeber ausgeschlossen (2026-11-22); die Anforderung
  nennt es in ihrem Out-of-Scope, damit es nicht über `file` oder einen Glob-Schalter
  still hineinwächst.
- **CHANGELOG, READMEs, Handbuch** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17), kein Feature-Commit.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Die Anforderung steht im Lastenheft (Bereichs-Kürzel in der Liste ergänzt,
      Version gehoben, Historie-Zeile nach
      [`MR-032`](../../../../harness/conventions.md#mr-032)) mit Akzeptanzkriterien je
      Entscheid (Zähleinheit, jede Ausnahme, Grenzwert genau auf der Schwelle), den
      Rändern mit Exit 2 und ihrem Out-of-Scope.
- [ ] Die Spezifikation führt das Verfahren (Leseweise, Zählregel, Ausnahmen, Reihenfolge
      der Befunde, Bezugsmenge), das Schema `linelength.*` und den Grund-Code im Verfahren;
      jede Aussage über die Grenzen (welche Eingaben das Modul liest, [`AGENTS.md`](../../../../AGENTS.md)
      §3.8) ist am Gegenstand gezählt oder gemessen, nicht aus der Beschreibung übernommen.
- [ ] Die ADR (Entscheide samt verglichenen Alternativen, darunter „Bedingung in `file`“ und
      „Bedingung in `structure`“) ist `Accepted`, im [ADR-Index](../../adr/README.md) geführt
      und trägt `## Re-Evaluierungs-Trigger`.
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
| `spec/spezifikation.md` | update | Verfahren, Schema `linelength.*`, Grund-Code im Verfahren |
| ADR im Verzeichnis `docs/plan/adr/` (nächste freie Nummer, Titel „linelength-modul“) | neu | Entscheidung samt Alternativen |
| `docs/plan/adr/README.md` | update | ADR-Index ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 4) |
| Plan von [`slice-290`](../open/slice-290-linelength-modulkern.md) | update | die Testfälle, die aus den bestätigten Entscheiden folgen, und die zwei mit dem Code kommenden Spiegel-Stellen |

**Festzuschreibende Entscheide** — Empfehlungen; offen bis zum Entscheid des
Auftraggebers, ohne ihn beginnt kein Code:

1. **Eigenes Modul statt Bedingung.** Modul `linelength`, Bereichs-Kürzel `LLEN`
   (`LINK` ist belegt). Verworfen: eine dritte Schwelle im Modul `file` — es liest jede
   Dateiart roh über eigene Globs, die Markdown-Ausnahmen (Fenced-Code, Tabellen)
   brächen seine Zusage „unabhängig von der Dateiart“; eine Bedingung in `structure` —
   sie ist an einen Abschnitt gebunden und meldet auf der Überschriftszeile, eine
   Zeilenlänge gilt der ganzen Datei und gehört auf ihre Zeile.
2. **Eingabemenge:** die Markdown-Dateien der Scan-Menge (`scan.roots`/`scan.ignore`,
   [`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln))
   mit modul-lokalem `scope`
   ([`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope))
   und `exempt-paths` wie die anderen Prosa-Module. Das Modul liest keine Datei außerhalb
   dieser Menge — die §3.8-Frage hat damit die Antwort „keine“, und die Spezifikation
   sagt es.
3. **Zähleinheit:** Unicode-Codepunkte der rohen Zeile ohne Zeilenende (`\n`, ein `\r`
   davor zählt nicht) — dieselbe Einheit wie `cell-max-chars`
   (`utf8.RuneCountInString`). Ein Tab zählt als ein Zeichen. Verworfen: Bytes (ein
   Umlaut zählte doppelt) und Anzeigebreite (braucht eine Breiten-Tabelle, East-Asian
   Wide, Emoji — nicht deterministisch über Unicode-Versionen hinweg). Ungültiges UTF-8
   zählt je Byte als ein Zeichen (Verhalten von `RuneCountInString`), die Anforderung
   nennt es.
4. **Schwelle und Schlüsselname:** `linelength.max-chars` (Pflicht bei aktivem Modul,
   ≥ 1, sonst Exit 2); geprüft wird `Länge > max-chars`, eine Zeile genau auf der
   Grenze ist befundfrei. Bewusst **nicht** `max-lines`/`max-length`: der Name
   `max-lines` trägt schon zwei Bedeutungen
   ([`BEO-ALL/max-lines-namenskollision-structure-file`](../observations/BEO-ALL/max-lines-namenskollision-structure-file/observation.md)),
   eine dritte nahe Form vergrößerte die Verwechslung.
5. **Ausnahmen, je einzeln schaltbar, Default „ausgenommen“:** (a) Zeilen in
   Fenced-Code-Blöcken samt ihrer Zäune (`exempt-code-blocks`), (b) Zeilen einer
   Pipe-Tabelle (`exempt-tables`), (c) Zeilen, deren Teil ab der Grenze kein Leerzeichen
   enthält — eine lange URL, ein langer Link, ein langer Inline-Code-Pfad lässt sich
   nicht umbrechen (`exempt-unbreakable`). Überschriften zählen mit. Die Erkennung von
   Fenced-Code und Tabellenzeile nimmt das Modul aus der geteilten Lexik
   (`proseLineSet`, `tableRowLine`), kein eigenes Muster daneben
   ([`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)).
   Front-Matter ist keine eigene Ausnahme (d-check kennt ihn heute nicht als Form).
6. **Befund:** je überlanger Zeile ein Befund, Ort = Datei und Zeile, Meldung nennt
   gemessene Länge und Grenze; Grund-Code `line-chars-exceeded` (Empfehlung, analog
   `file-lines-exceeded`/`file-bytes-exceeded`). Reihenfolge nach Pfad, dann Zeile
   ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
7. **Bezugsmenge:** die Vollständigkeits-Zeile nennt geprüfte Dateien und geprüfte
   Zeilen **nach** Abzug der Ausnahmen, damit ein grünes Ergebnis nicht mehr behauptet,
   als es misst (Baseline-Regelwerk `modul-13-quality-gates.md` §Hard Rule, *Ein Gate
   ohne seine Grenze*).
8. **Leere Kandidaten-Menge:** aktives Modul und keine Markdown-Datei in der Menge ist
   Exit 2 (fail-closed wie `structure`/`file`), kein stilles Grün.

**Gegenprobe vor dem Schreiben (Schritt 19 des Workflow-Skeletts):** die Zählregel läuft
am eigenen Bestand (`awk '{ if (length($0) > 120) n++ } END { print n }'` über die
gescannten `.md`, in einer UTF-8-Locale; Zahl mit und ohne Fenced-Code steht im
Plan-Nachtrag) — das belegt die Out-of-Scope-Aussage „eigenes Profil nicht scharf“ mit
einer Zahl statt einer Vermutung. Die Negativfälle (Zeile genau auf der Grenze,
Umlaut/Emoji, `\r\n`, Tab, nicht geschlossener Fence bis Dateiende, Tabellenzeile ohne
Trennzeile, URL mit Leerzeichen davor) stehen als Testfälle für
[`slice-290`](../open/slice-290-linelength-modulkern.md).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Welle `welle-95-linelength-modul` ist offen und das
WIP-Limit ist frei (`in-progress/` trägt nur die Roadmap).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Entscheide 1–8 brauchen
  mehr als eine Review-Sitzung, oder die Ausnahme (c) verlangt eine eigene Link-Lexik,
  die die geteilte nicht hergibt.
- `in-progress` → `open` (blockiert): der Auftraggeber entscheidet einen der Punkte
  1, 3 oder 5 nicht; ohne Entscheid beginnt kein Code.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Lastenheft, Spezifikation und ADR liegen vor und `make gates` ist grün),
Review- und Verifikations-Report liegen unter `docs/reviews/`, Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Ausnahme (c) „kein Leerzeichen ab der Grenze“ ist zu grob: eine Zeile mit langem
  Link **und** umbrechbarem Text davor wird ausgenommen, obwohl sie kürzbar wäre —
  **Ausgang:** offen bis zur Closure.
- Ein vergessener Schluss-Fence macht für die geteilte Lexik den Rest der Datei zu Code
  (die Fence-Hälfte von
  [`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md),
  2×); nimmt das Modul Fenced-Code über dieselbe Lexik aus, fällt jede überlange Zeile
  danach still aus der Prüfung — ein dritter Beleg wäre die Schwelle. Die
  Inline-Code-Hälfte trifft das Modul nicht: es misst die rohe Zeile —
  **Ausgang:** offen bis zur Closure.
- Der Name `linelength` liest sich neben `file.max-lines` und `structure.max-lines` als
  vierte Längen-Achse und wird verwechselt — **Ausgang:** offen bis zur Closure.

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

*(Wird bei der Closure gefüllt.)*

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand
2026-11-22), Sub-Area `*`. Treffer:
[`BEO-ALL/max-lines-namenskollision-structure-file`](../observations/BEO-ALL/max-lines-namenskollision-structure-file/observation.md)
(offen, 1× — Entscheid 4 meidet einen dritten nahen Schlüsselnamen; Risiko 3 in §6),
[`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md)
(offen, 2× — Risiko 2 in §6; ein dritter Beleg wäre die Schwelle),
[`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)
(offen, 2× — Entscheid 5 verlangt die geteilte Lexik),
[`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
(offen, 1× — jede Randbedingung der Anforderung braucht in
[`slice-290`](../open/slice-290-linelength-modulkern.md) einen Test) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — Entscheid 2 beantwortet die
Frage). Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-22: **beide rot**. `upstream-drift` (jüngster Lauf
2026-11-22T07:40:16Z): lokal nachgeprüft fünf veraltete Pins — Go 1.27.3 → 1.27.4,
Semgrep 1.183.0 → 1.186.0, a-check 0.23.3 → 0.24.0, `docker/login-action` 4.7.0 → 4.8.0,
Digest von `golang:1.27.3`. `image-scan` (2026-11-22T11:29:01Z): behebbare
CRITICAL/HIGH im publizierten Image, die Kennungen stehen im Job-Summary des Laufs (das
Log war über `gh` nicht abrufbar). Beides ist planmäßige Fremd-Drift, berührt dieses
Modul nicht und gehört als eigene Pin-Hebung vor das nächste Release
(Welle-Out-of-Scope).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
