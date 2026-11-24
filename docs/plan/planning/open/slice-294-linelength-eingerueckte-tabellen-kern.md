# Slice slice-294: Kern und Tests — `linelength` nimmt eine `|`-Zeile mit bis zu drei Leerzeichen davor aus

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-97-linelength-eingerueckte-tabellen`](../welle-97-linelength-eingerueckte-tabellen.md).

**Bezug:** [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
in der Fassung aus [`slice-293`](../done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md),
[`ADR-0054`](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) (Entscheidung 4: je Konsument
eine Assertion), die Folge-ADR aus [`slice-293`](../done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md)
(`Proposed` bis zur Closure dieses Slice), [`ADR-0114`](../../adr/0114-linelength-modul.md) (nur
die Teile, die die Folge-ADR nicht ablöst).

**Berührte Spec-Stellen:** `—`. Die Spezifikation ändert
[`slice-293`](../done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md); dieser Slice liest
§[`DC-FA-LLEN-001.a`](../../../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength)
Schritt 4 als Vorgabe. Trifft der Code auf einen Rand, den sie anders beschreibt, ist das ein
Befund zur Spezifikation und eine Plan-Änderung **vor** dem Code (§4, Rückführung), keine stille
Angleichung. Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-11-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Mit `linelength.exempt-tables` an nimmt das Modul eine `|`-Zeile aus, vor der ein bis drei
Leerzeichen stehen — und nur diese; jedes der sechs Kriterien aus
[`slice-293`](../done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md) §3 hat einen Test, der
ohne die Regel aus dem richtigen Grund rot läuft, und die Folge-ADR ist nach Review und
Verifikation `Accepted`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Image-Probe und der Vorher-/Nachher-Vergleich** — sie sind der Closure-Trigger der Welle
  ([§3](../welle-97-linelength-eingerueckte-tabellen.md)), nicht diese DoD
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 10); der Slice endet am grünen `make gates`.
- **`tableRowLine` selbst weiten** — entschieden dagegen: der Helfer wohnt daneben und schließt
  `tableRowLine` ein ([`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4,
  Alternative B; [Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-1). Die
  Weitung verschöbe `targets`, `planning.waves` und `structure` gegen ihre Spalte-0-Zusagen.
- **`> |`, eine Listen-Erkennung, ein neuer Schlüssel** — Welle §6; der Konfigurations-Rand
  (`configyaml.go`) bleibt unberührt, weil kein Schlüssel dazukommt.
- **`tools/image-test.sh` und `harness/sensors/image-test.md`** — der Fixture dort trägt die
  Tabellenzeile mit `|` in Spalte 0 (Z. 10) und keine eingerückte; seine Erwartungen ändern sich
  nicht. Der Image-Beleg des Anlasses ist die Probe der Welle.
- **CHANGELOG, READMEs, Version, Release** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 17), kein Feature-Commit.
- **Das eigene Profil des Repos** — schaltet `linelength` nicht ein
  ([`ADR-0114`](../../adr/0114-linelength-modul.md), Konsequenzen).

Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Erkennung:** [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
      (Kriterien 1 bis 4 und 6 aus [`slice-293`](../done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md)
      §3) ist im Modulkern erfüllt: die Regel nimmt genau die Zeilen aus, die dort stehen, und keine
      weitere; jeder Test läuft ohne die Regel **aus dem richtigen Grund** rot (Mutation je Test,
      §3), und die Closure-Notiz nennt je Test die Mutation und die gelesene Meldung.
- [ ] **Kopplung:** Kriterium 5 — der Kopplungs-Test hält die Frage „Tabellenzeile?“ je Konsument fest,
      **in beiden Richtungen** (`linelength` nimmt die eingerückte Zeile aus, `targets`,
      `planning.waves` und `structure` lesen sie weiter nicht als Tabellenzeile), und die geteilte
      Antwort `tableRowLine` steht im Diff unverändert. Die Assertionen sind **zwei**
      ([`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4): die drei
      Konsumenten von `tableRowLine` auf **einer** Antwort, die eingerückte Zeile eingeschlossen, und
      `linelength` auf der Antwort des Helfers. Der Kopf-Kommentar des Kopplungs-Tests („… GLEICH
      beantwortet …, eine eingerueckte Zeile ist keine Spalte-0-Tabellenzeile“) ändert sich mit.
- [ ] **Doku und Entscheidung:** das Handbuch §4.28 samt Modul-Tabelle, der Kommentar der
      `--print-config`-Vorlage und die Kommentare am Code nennen die Regel mit ihrer Bedingung; die
      Abdeckungs-Dateien sind aus den Tests abgeleitet (`make abdeckung`); die Folge-ADR ist
      **`Accepted`, erst nach Review und Verifikation** ([`BEO-ALL/adr-accepted-bevor-der-mechanismus-steht`](../observations/BEO-ALL/adr-accepted-bevor-der-mechanismus-steht/observation.md)).
      **Im selben Commit**, der [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)
      `Accepted` macht: eine Zeile in `## Geschichte` von [`ADR-0114`](../../adr/0114-linelength-modul.md)
      („Entscheidung 5 (b) teil-superseded durch [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)“) und im Index die Zeile von [`ADR-0114`](../../adr/0114-linelength-modul.md) mit
      „Accepted (teil-superseded: [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md))“; die Status-Zeile von [`ADR-0114`](../../adr/0114-linelength-modul.md) bleibt `Accepted`
      ([Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-6).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/rules/markdown.go` | update | **Ein neuer benannter Helfer** neben `tableRowLine` (Name vergibt der Implementer): er **ruft `tableRowLine` und erweitert sie** um null bis drei Leerzeichen (U+0020) davor, statt die Spalte-0-Antwort ein zweites Mal zu schreiben ([`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4). `tableRowLine` bleibt zeichengleich. Der Kommentar nennt die zweite Frage („eingerückt, bis drei Leerzeichen“) und die ADR als Herkunft. Für die eingerückte Zeile verlangt der Helfer selbst, dass sie außerhalb von Fenced-Code liegt (Kriterium 3) |
| `internal/hexagon/core/rules/linelength.go` | update | die Ausnahme (b) ruft **nur den Helfer** (nicht zusätzlich `tableRowLine`: mit zwei Aufrufen bliebe die Mutation zu Kriterium 6 grün); der Dateikopf (`GRENZE, ausgesprochen`) nennt die neue Antwort samt ihrer Grenze (ab vier Leerzeichen und mit Tab gemessen, weil die Lexik keinen Listenkontext hergibt) |
| `internal/hexagon/core/rules/linelength_test.go` | update | Kriterien 1 bis 4 und 6, Happy/Boundary/Negative, dazu die zwei Grenzen aus dem Punkt „Tabellenzeile heißt …“ ([Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-4): eine Prosa-Zeile mit drei Leerzeichen und `\|` ohne Tabelle ist ausgenommen (Grenze 1), `- \| a \| b \|` bleibt gemessen (Grenze 2); Grenze 3 (vier Leerzeichen) deckt Kriterium 2 — die Mutationen stehen unten |
| `internal/hexagon/core/rules/lexikon_kopplung_test.go` | update | Kriterium 5: **zwei Assertionen** ([`ADR-0054`](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 4, [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4) — `targets`, `planning.waves` und `structure` auf einer Antwort (`TestTabellenzeilenFrageHatEineAntwort`; die eingerückte Zeile eingeschlossen), `linelength` auf der Antwort des Helfers, in beiden Richtungen. **Der Kopf-Kommentar des Tests ändert sich mit:** er sagt heute „GLEICH beantwortet … eine eingerueckte Zeile ist keine Spalte-0-Tabellenzeile“ für alle vier |
| `docs/user/benutzerhandbuch.md` | update | §4.28: die YAML-Zeile zu `exempt-tables`, die Beschreibung der Ausnahme (b) und der Satz zur eingerückten `\|`-Zeile; die Zeile des Moduls in der Modul-Tabelle (§6). Zählung unten |
| `internal/adapter/driving/cli/config_template.go` | update | der Kommentar zu `exempt-tables` in der `--print-config`-Vorlage, falls der Wortlaut aus [`slice-293`](../done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md) ihn verlangt |
| `docs/user/abdeckung-tests.md` | `make abdeckung` | die neuen Testnamen; `make test` hält die Datei gegen ihre Ableitung ([ADR-0104](../../adr/0104-test-nachweise-entlasten-in-der-rtm.md)) |
| [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) und ADR-Index | `**Status:**`-Übergang | `Proposed` → `Accepted` **nach** Review und Verifikation |
| [`ADR-0114`](../../adr/0114-linelength-modul.md) und ADR-Index | Geschichte-Zeile, Index-Zeile | **im selben Commit** wie der Übergang von [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md): eine Zeile in `## Geschichte` („Entscheidung 5 (b) teil-superseded durch [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)“), im Index „Accepted (teil-superseded: [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md))“ wie bei [`ADR-0002`](../../adr/0002-distribution-ghcr-image.md) und [`ADR-0013`](../../adr/0013-pr-ci-und-traceability-gate.md). Die Status-Zeile bleibt `Accepted`: `vcs.head-allow` lässt nur `Accepted` oder `Superseded by ADR-NNNN` zu, und die Ablösung ist nur teilweise ([Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-6) |

**Rot-Beleg je Test — der Vorzustand ohne die Regel, ausgeschrieben**
([`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md)).
Die Mutation ist ein Vorschlag; der Implementer führt sie aus und trägt die gelesene Meldung in die
Closure-Notiz ein. Ein Test, dessen Eingabe den Zweig der Regel nicht erreicht, bliebe auch ohne sie
grün:

| Kriterium | Mutation | Erwartung |
|---|---|---|
| 1 (drei Leerzeichen) | `linelength.go` ruft für (b) wieder `tableRowLine` statt des Helfers | der Test findet einen Befund an der Zeile mit drei Leerzeichen, will keinen |
| 1 (die Grenze 3) | die Schwelle „höchstens drei“ wird „höchstens zwei“ | der Fall mit drei Leerzeichen wird rot, der mit einem bleibt grün |
| 2 (vier Leerzeichen, Tab, U+00A0) | die Schwelle wird „höchstens vier“ bzw. ein Tab zählt als Leerzeichen | der Fall mit vier Leerzeichen bzw. mit Tab wird rot — je Fall einzeln, damit nicht nur die erste Seite der Bedingung variiert ([`BEO-ALL/negativtest-deckt-nur-eine-regelseite`](../observations/BEO-ALL/negativtest-deckt-nur-eine-regelseite/observation.md)) |
| 2 (`> \|`, `a \| b`, Text mit Einrückung) | der Helfer trimmt auch `>` bzw. verlangt kein `\|` am Anfang | der jeweilige Fall wird rot |
| 3 (im Zaun) | die Prüfung gegen `proseLineSet` entfällt im Helfer für die eingerückte Zeile | der Zaun-Fall findet keinen Befund, will einen |
| Grenze 1 (Prosa-Zeile mit drei Leerzeichen und `\|`, ohne Tabelle) | der Helfer verlangt eine Trennzeile oder eine zweite `\|`-Zeile | die Prosa-Zeile wird gemessen und liefert einen Befund, der Test will keinen |
| Grenze 2 (`- \| a \| b \|`) | der Helfer trimmt auch eine Listenmarke vor dem `\|` | die Zeile wird ausgenommen, der Test will einen Befund |
| 4 (Bezugsmenge) | `count.Lines++` rückt vor die Ausnahme | die Bezugsmengen-Zeile nennt eine Zeile zu viel |
| 5 (Kopplung) | `tableRowLine` selbst nimmt die eingerückte Zeile an | die Assertion für `targets`, `planning.waves` und `structure` wird rot |
| 6 (Rest unverändert) | der Helfer ruft `tableRowLine` nicht mehr (ist `linelength` auf den Helfer allein gestellt, trifft er die Spalte 0 dann nicht mehr) | der Fall mit `\|` in Spalte 0 wird rot — mit **zwei** Aufrufen in `linelength.go` bliebe er grün, weil `tableRowLine` die Spalte 0 weiter trifft; deshalb ruft `linelength` nur den Helfer |

**Spiegel-Zählung im Handbuch und in der Vorlage** (Form vor Zahl, Regel 14): die Beschreibung der
Ausnahme (b) steht im Handbuch an **drei** Stellen und in der Modul-Tabelle an **einer**, in der
Vorlage an **einer** —

```bash
git grep -n -E 'exempt-tables|eingerückte `\|`-Zeile|Zeilen, die mit `\\?\|` beginnen' -- docs/user/benutzerhandbuch.md internal/adapter/driving/cli/config_template.go
```

— am 2026-11-24: Handbuch Z. 1679 (YAML-Kommentar), 1706 (Beschreibung), 1713 (Satz zur eingerückten
Zeile), 3135 (die Modul-Tabelle in §6), Vorlage Z. 314. **Unverändert bleiben** Z. 1725 (nennt
`exempt-tables:` ohne Wert ⇒ Exit 2) und Z. 3369 (die Historie-Zeile zu v0.92.0; die neue
Historie-Zeile gehört der Release-Prep). Die READMEs tragen dieselbe Beschreibung und gehören der
Release-Prep. Der Implementer zählt vor dem Editieren **neu**.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-293`](../done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md)
liegt in `done/` (die Anforderung steht im Lastenheft, die Folge-ADR als `Proposed`), der
Auftraggeber hat die vier Entscheide der Welle bestätigt, das Architect-Verdikt zu Frage 4 liegt
unter `docs/reviews/`, und das WIP-Limit ist frei (`in-progress/` trägt nur die Roadmap).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Erkennung verlangt eine Änderung an
  `tableRowLine` oder eine Listen-Erkennung, die über den einen neuen Helfer hinausgeht — beides
  trüge mehr als die sechs Kriterien.
- `in-progress` → `open` (blockiert): der Code trifft auf einen Rand, den die Anforderung oder die
  Spezifikation anders beschreibt; die Berichtigung ist dann ein Befund zur Spec und eine
  Plan-Änderung vor dem Code ([`BEO-ALL/spezifikation-und-code-weichen-am-rand-ab`](../observations/BEO-ALL/spezifikation-und-code-weichen-am-rand-ab/observation.md)).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Erkennung, Kopplung, Doku und Entscheidung, `make gates` grün), Review- und
Verifikations-Report liegen unter `docs/reviews/`, Closure-Notiz mit Lerneintrag geschrieben.
**Review-Vorgabe der Welle:** höchstens zwei Runden, eingearbeitet werden nur HIGH und MEDIUM; jede
Nachbesserung zitiert die Stelle, die sie trägt
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 18). Was die **Welle** einlöst — der Vorher-/Nachher-Lauf
und die Probe am gebauten Image — steht in
[ihrem Closure-Trigger](../welle-97-linelength-eingerueckte-tabellen.md) und nicht in dieser DoD
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 10).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der neue Helfer liegt in der geteilten Lexik; ein Zugriff dort verschiebt andere Module
  ([`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md),
  2×; [`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md),
  verkörpert) — **Ausgang:** *entfallen*, wenn der Diff von `markdown.go` nur Neues neben
  `tableRowLine` zeigt und der Kopplungs-Test beide Richtungen trägt; sonst *eingetreten* mit einer
  Evidenz-Datei im Register.
- Ein Test läuft grün, ohne die Regel zu fahren
  ([`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md),
  2×) — **Ausgang:** *entfallen*, wenn die Closure-Notiz je Test die Mutation und die gelesene
  Meldung nennt und jede rot war; sonst *eingetreten* mit einer Evidenz-Datei (dann 3×).
- Eine Randbedingung der Anforderung (Grenze 3 und 4, Tab, U+00A0, Zaun, Bezugsmenge) hat keinen
  Test ([`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md),
  2×) — **Ausgang:** *entfallen*, wenn die Verifikation zu jeder Grenze der Anforderung einen Test
  findet, der sie umdreht; sonst *eingetreten* mit einer Evidenz-Datei (dann 3×).
- Die Negativfälle fehlen, weil die Regel zwei Seiten hat (Einrückung **und** erstes Zeichen)
  ([`BEO-ALL/negativtest-deckt-nur-eine-regelseite`](../observations/BEO-ALL/negativtest-deckt-nur-eine-regelseite/observation.md),
  1×) — **Ausgang:** *entfallen*, wenn je Seite ein eigener Fall rot wird (siehe die Tabelle in §3);
  sonst *eingetreten* mit einer Evidenz-Datei.
- Die Bezugsmenge zählt nach der Änderung anders (weniger geprüfte Zeilen bei eingerückten
  Tabellenzeilen), und eine Konsumenten-Probe liest die Zahl — **Ausgang:** *entfallen*, wenn
  `make gates` grün ist und Kriterium 4 die Zahl festhält; die Abweichung ist gewollt und steht in
  der Probe der Welle als eigene Zeile.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) — Modulkern und geteilte
Lexik unter `internal/hexagon/core/rules/`, Tests, Handbuch; keine Sub-Area mit eigener
Modus-Deklaration (die Deklaration führt nur `*` und `tools/harness/`, und `tools/harness/` bleibt
unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand 2026-11-24,
`ls docs/plan/planning/observations/BEO-ALL/`), Sub-Area `*`. Zähler = Dateien unter `evidence/`.
Treffer:
[`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)
(offen, 2× — Risiko 1 in §6; der Helfer liegt in der geteilten Lexik),
[`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md)
(verkörpert in welle-74, 6× — der Kopplungs-Test, Kriterium 5),
[`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md)
(offen, 2× — Risiko 2; die Mutation je Test in §3),
[`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
(offen, 2× — Risiko 3),
[`BEO-ALL/negativtest-deckt-nur-eine-regelseite`](../observations/BEO-ALL/negativtest-deckt-nur-eine-regelseite/observation.md)
(offen, 1× — Risiko 4; die Regel hat zwei Seiten),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(verkörpert im Workflow-Skelett, 7× — die Negativliste: Text mit Einrückung ohne `|`, `a | b`, `> |`,
Tab, U+00A0),
[`BEO-ALL/spezifikation-und-code-weichen-am-rand-ab`](../observations/BEO-ALL/spezifikation-und-code-weichen-am-rand-ab/observation.md)
(offen, 1× — die Rückführung in §4),
[`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)
(offen, 1× — die Grenze „ab vier Leerzeichen gemessen“ steht im Code-Kommentar mit ihrer Bedingung),
[`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
(verkörpert, 4× — die Spiegel-Zählung in §3 nimmt Handbuch und Vorlage auf) und
[`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
(verkörpert — die Botschaft behauptet nur, was der Diff trägt). **Kein Treffer:**
[`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md)
(offen, 2×) — `linelength` misst die rohe Zeile und liest keine Inline-Code-Spanne. Der Zähler folgt
aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-24: `image-scan` **grün** (jüngster Lauf 2026-11-23T12:51:15Z).
`upstream-drift` **rot** (jüngster Lauf 2026-11-23T08:29:02Z), und dieser Lauf liegt **vor** den
Pin-Hebungen vom 2026-11-23; sein Log war über `gh run view <id> --log` nicht abrufbar (leere
Ausgabe) — ob die Meldung planmäßige Fremd-Drift ist, ist damit **nicht gelesen**. Beim Start dieses
Slice neu lesen. Berührt dieses Modul nicht (die Läufe messen Pins und Images, nicht `linelength`).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
