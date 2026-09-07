# Slice slice-211: Obermengen-Nachweis gegen `markdownlint` MD013

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [eingehender CR *Zeilenlänge als prüfbare Größe*](../../cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md)
(Frage 2, Auftraggeber-Entscheid 2026-09-07: erst der Obermengen-Nachweis),
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das Modul, das eine Bedingung bekäme),
[`MR-046`](../../../../harness/conventions.md#mr-046) (keine vierte Toolchain
nebenbei — geprüft, siehe §6).

**Berührte Spec-Stellen:** — *(keine; der Slice misst und entscheidet, er
ändert keine Anforderung. Fällt der Entscheid auf „bauen", ist das ein
Folge-Slice mit eigener Spec-Berührung.)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die offene **Frage 2** des Zeilenlängen-CR beantworten — ist
`markdownlint` MD013 eine **Obermenge** des repo-eigenen Mittels? — und mit
dem Ergebnis den CR entscheiden. Die Form des Nachweises gibt der Kanon vor
und sie ist ausdrücklich **nicht** der Datenblatt-Vergleich: je Verstoßklasse
ein Break-Test mit **beiden** Sensoren nebeneinander, plus der unveränderte
Bestand, auf dem beide schweigen müssen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die dauerhafte Aufnahme von `markdownlint`.** Sie wäre ein eigener
  Vorgang mit [`MR-046`](../../../../harness/conventions.md#mr-046)-Nachfolger,
  ADR, `make`-Target und Deklaration in
  [`AGENTS.md`](../../../../AGENTS.md) §4 und
  [`harness/README.md`](../../../../harness/README.md) — und sie entschiede die
  Frage vorweg, die dieser Slice beantworten soll. Wer das Werkzeug aufnimmt,
  bevor der Nachweis vorliegt, dreht die Reihenfolge um.
- **Eine `max-line-chars`-Bedingung im Produkt.** Sie hängt am Entscheid;
  fällt er auf „bauen", ist sie ein Folge-Slice mit
  Lastenheft-Anforderung und Akzeptanzkriterien.
- **Jede Änderung an [`.d-check.yml`](../../../../.d-check.yml).** Der
  `forbid-pattern`-Weg wird **gemessen**, nicht scharfgeschaltet — sonst
  entstünde ein Gate, dessen Berechtigung derselbe Slice erst prüft.
- **Ein Retrofit des Bestands.** Die 524 überlangen Tabellenzeilen und 281
  Fließtext-Zeilen bleiben, wie sie sind; sie sind der **Messgegenstand**,
  nicht die Arbeit.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** Der Nachweis liegt vor: **je Verstoßklasse** ein Break-Test mit
      beiden Sensoren nebeneinander **und** der unveränderte Bestand, auf dem
      beide schweigen — jeweils mit **echter Ausgabe**, nicht mit behauptetem
      Ergebnis.
- [x] **(2)** Der Nachweis ist **reproduzierbar dokumentiert**: das
      vollständige Kommando samt Image-Digest steht im CR. Ohne das ist er
      eine Selbstauskunft und der Kanon-Anspruch nicht erfüllt.
- [x] **(3)** Der CR ist **entschieden**: alle drei Fragen beantwortet,
      `Stand:`-Zeile gesetzt, Begründung je Frage.
- [x] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Zuerst die Form des Gegenstands** — die Regel aus
[`AGENTS.md`](../../../../AGENTS.md) §5 (`seit slice-210`) gilt diesem Slice
als erstem, denn er ist eine Messung.

**Was zählt als *Verstoßklasse*?** Ein Paar aus **Eingabe-Form** und
**erwartetem Verhalten**, das der CR in seinen Akzeptanzkriterien nennt — nicht
jede denkbare Markdown-Konstruktion. Sechs sind es:

| # | Eingabe | vom CR gefordert |
|---|---|---|
| K1 | Fließtext-Zeile über der Schwelle | Befund |
| K2 | Abschnitt aus lauter kurzen Zeilen, Summe über der Schwelle | kein Befund |
| K3 | Zeile in einem Fenced Block | kein Befund |
| K4 | Tabellenzeile | kein Befund |
| K5 | unteilbares Token (lange URL) | kein Befund |
| K6 | Inline-Code-Spanne | kein Befund |

**Was zählt als *Obermenge*?** Drei Teile, alle drei einzeln nachzuweisen
(Baseline-Regelwerk `modul-11-verification.md` §Fitness Function ohne
Standard-Tool): dieselbe **Kandidaten-Menge** (welche Dateien werden geprüft),
dieselben **Bedingungen**, dieselbe **Schwelle, wie die Anforderung sie setzt**
— nicht die Vorbelegung des Werkzeugs. Fehlt einer der drei, ist MD013 keine
Obermenge, auch wenn es alle sechs Klassen richtig meldet.

**Der Skopus-Unterschied ist vorab bekannt und gehört gemessen, nicht
weggerechnet:** MD013 urteilt **datei-weit**, `structure` **abschnitts-weit**.
Ein Vergleich, der das übergeht, zählt zwei verschiedene Gegenstände
gegeneinander.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Proben-Dateien unter dem Scratchpad | neu (repo-extern) | je Klasse eine Datei; sie gehören nicht ins Repo — sie sind Messmittel, kein Artefakt |
| [der CR](../../cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md) | update | Nachweis, Kommando mit Digest, Entscheid |

## 4. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei. Der Auftraggeber-Entscheid
vom 2026-09-07 („erst der Obermengen-Nachweis") liegt vor und steht im CR.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt sich, dass der Nachweis nur mit einem
  **dauerhaften** Träger reproduzierbar ist — etwa weil das Kommando ohne
  `make`-Target und Konfigurationsdatei nicht wiederholbar bleibt —, dann ist
  die Toolchain-Aufnahme die eigentliche Arbeit und dieser Slice zu klein
  geschnitten.
- `in-progress` → `open` (blockiert): Lässt sich kein digest-gepinntes
  markdownlint-Image ohne Netz-Bau beschaffen, ruht der Slice bis zum
  Entscheid über die Bezugsform.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) zu **jeder** der sechs
Klassen steht die echte Ausgabe beider Sensoren im CR, dazu der Lauf über den
unveränderten Bestand; (b) der CR trägt eine `Stand:`-Zeile mit Entscheid und
`make gates` ist grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **[`MR-046`](../../../../harness/conventions.md#mr-046) ist geprüft und bleibt
  gültig — aber nur unter einer Bedingung.**
  Der Eintrag sagt, eine vierte Toolchain entstehe nicht nebenbei, und sein
  Auflösungs-Trigger ist genau diese Slice-Planung. — **Ausgang:** entfallen —
  er bleibt gültig, und die Bedingung ist geprüft statt behauptet: **kein
  Target, keine Deklaration in [`AGENTS.md`](../../../../AGENTS.md) §4 oder
  [`harness/README.md`](../../../../harness/README.md), kein `make`-Aufruf.**
  **Eine Zusage der ersten Fassung trifft dagegen nicht mehr zu, und der Review
  hat es gefunden:** Sie sagte *„kein Pin in einer getrackten Datei"* — der
  Image-Digest steht jetzt im CR, und der ist getrackt. **Das ist gewollt und
  kein Rückfall:** DoD (2) verlangt einen reproduzierbaren Nachweis, und
  reproduzierbar heißt hier, dass der Digest lesbar dasteht. Was
  [`MR-046`](../../../../harness/conventions.md#mr-046)
  meint, ist ein Pin, den ein **Lauf** zieht — eine Zahl in einem Zeitdokument
  ist keine Toolchain. Die ungenaue Formulierung ist der Befund, nicht der
  Digest.
- **MD013 kann Obermenge sein und trotzdem nicht einsetzbar.** Die drei Fragen
  des CR sind unabhängig: Frage 2 („Obermenge?") kann **ja** ergeben und
  Frage 3 („rechtfertigt der Nutzen den Schritt?") trotzdem **nein** — die
  Kosten einer Node-Toolchain sind keine Eigenschaft von MD013. Wer aus einem
  Ja auf Frage 2 ein Ja auf Frage 3 schließt, hat eine Frage übersprungen.
  — **Ausgang:** \<offen\>
- **Die Klassen-Liste stammt aus dem Anlass, nicht aus einer Inventur**
  ([`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  9×). Die sechs Klassen sind die Akzeptanzkriterien des CR; eine siebte, die
  niemand aufgeschrieben hat, fällt durch den Nachweis. Das ist die Grenze des
  Nachweises und gehört in den Entscheid. — **Ausgang:** \<offen\>
- **Der Vergleich misst zwei Werkzeuge mit verschiedenem Skopus.** MD013
  urteilt datei-weit, `structure` abschnitts-weit; die Zahl der Befunde ist
  deshalb **nicht** vergleichbar, nur ihr Verhalten je Klasse. Die eigene
  Vormessung hat das bereits gezeigt: zwei Abschnitts-Selektoren lieferten über
  derselben Schwelle 342 gegen 645 Befunde. — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice ändert genau **ein**
getracktes Artefakt — das CR-Dokument. Die Proben liegen außerhalb des Repos.
`tools/harness/` ist **nicht** berührt: §1 schließt Target und Deklaration
ausdrücklich aus, und ohne sie gibt es dort nichts anzufassen. Die
Sub-Area-Wahl ist damit trivial und wird trotzdem notiert, weil die Prüfung
unbedingt ist.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **37** Verzeichnisse über beide
Kürzel). **Fünf** Einträge sind einschlägig — und der erste ist es, weil
dieser Slice selbst eine Messung ist:

- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (3×, seit slice-210 *verkörpert* in [`AGENTS.md`](../../../../AGENTS.md) §5)
  — **erste Anwendung der eigenen frischen Regel.** §3 schreibt deshalb die
  Form von *Verstoßklasse* und von *Obermenge* aus, **bevor** gemessen wird.
  Wäre das unterblieben, hätte der Slice die Regel gebrochen, die sein direkter
  Vorgänger geschrieben hat.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (11×, Stand *gemischt*) — der Skopus-Unterschied ist genau diese Klasse:
  MD013 urteilt datei-weit, `structure` abschnitts-weit. Wer die **Zahl** der
  Befunde vergleicht statt das **Verhalten je Klasse**, misst zwei verschiedene
  Mengen und sagt über eine aus.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (8×, Ausgang *geplant*) — der Entscheid am Ende darf nur behaupten, was
  gelaufen ist. Ein „MD013 deckt das ab", das aus dem Datenblatt statt aus der
  Ausgabe stammt, ist genau dieser Eintrag; der Kanon verbietet den
  Datenblatt-Vergleich ausdrücklich.
- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (9×) — §6 führt es als Risiko: Die sechs Klassen sind die
  Akzeptanzkriterien des CR, also der Anlass. Eine siebte, die niemand
  aufgeschrieben hat, fällt durch den Nachweis.
- [`begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md)
  (2×) — für den Entscheid: Die drei CR-Fragen sind unabhängig, und ein Ja auf
  Frage 2 begründet kein Ja auf Frage 3. Wer sie koppelt, trifft womöglich die
  richtige Entscheidung mit dem falschen Grund.

**Keiner der fünf erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-07 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). Nichts zu tun.

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch, und ausnahmsweise **gegen** den Slice
  gerichtet: [`MR-046`](../../../../harness/conventions.md#mr-046) regelt
  genau den Fall, den dieser Slice berührt, und sein Auflösungs-Trigger ist
  diese Planung. Er bleibt gültig — aber nur, solange nichts im Bestand bleibt
  (§6).
- **Phase-Reife:** Phase 5 für die Slice-Mechanik, **Phase 2** für den
  Gegenstand: Ein Obermengen-Nachweis nach der Kanon-Form ist in diesem Repo
  noch nie geführt worden. Es gibt keinen Präzedenzfall, an dem sich die Form
  ablesen ließe — nur den Kanon-Absatz selbst.
- **Evidenz-/Diskrepanz-Risiko:** **mittel.** Am Bestand ist nichts zu
  inventarisieren; das Risiko sitzt in der Messung selbst — verschiedener
  Skopus, eine Klassen-Liste aus dem Anlass, und ein Werkzeug, das dieses Repo
  noch nie gefahren hat.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.
