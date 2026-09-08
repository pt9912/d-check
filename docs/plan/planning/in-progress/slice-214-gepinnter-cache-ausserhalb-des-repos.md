# Slice slice-214: Was einen gepinnten Cache außerhalb des Repos prüft — und wann

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md) (hermetisches
semgrep-Gate), [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)
(Digest-Pins), [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette
(`baseline-verify` — die **andere** Hälfte derselben Frage),
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
(4×, verkörpert — der Anlass kam aus seinem Review).

**Berührte Spec-Stellen:** — *(keine; der Slice beschreibt eine bestehende
Eigenschaft und ändert kein Verhalten)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Für **jedes gepinnte Fremd-Artefakt** beantworten, **was seine
Unversehrtheit prüft und zu welchem Zeitpunkt** — und die Antwort dort
hinschreiben, wo sie fehlt. Der Anlass ist der `semgrep`-Regel-Cache: Er liegt
als **einziger** außerhalb des Repos, und die Frage nach ihm ist im Review von
slice-212 gestellt und nicht beantwortet worden.

**Die Achse ist neu und die Antwort vermutlich nicht überall dieselbe.**
`baseline-verify` prüft **jeden Lauf** und kann die **Echtheit** nicht beweisen
(slice-212). Beim `semgrep`-Cache ist es genau umgekehrt: Der Bezug über einen
**git-Commit-Pin** bindet die Echtheit stärker als jedes Manifest — ein
Commit-SHA ist ein Hash über den Baum —, aber **nach** dem Holen prüft ihn
nichts mehr; die Bedingung ist `[ ! -d "$RULES_DIR/$RULES_SUBSET" ]`, also
Existenz eines Verzeichnisses. **Zwei Artefakte, zwei entgegengesetzte
Lücken** — das ist der Grund für diesen Slice und nicht nur für eine Zeile.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Re-Verifikation einbauen.** Ob der `semgrep`-Cache bei jedem Lauf
  gegen seinen Commit-Pin geprüft werden soll, ist ein **Entscheid mit ADR**
  (Laufzeit gegen Sicherheit, und [ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)
  nennt den einmaligen Bezug ausdrücklich als Eigenschaft). Dieser Slice
  beschreibt den Ist-Zustand.
- **Der `ignore-refs`-/`exempt-paths`-CR.** Er liegt unentschieden in
  [`docs/plan/cr/`](../../cr/) und wartet auf eine Antwort des Absenders;
  ein anderer Vorgang.
- **Eine erneute Inventur der `## Grenze`-Abschnitte.** Sie liegt in
  [slice-212](../done/slice-212-grenzen-liste-nennt-ihre-groesste-luecke.md);
  hier wird **eine** Achse über **wenige** Artefakte gemessen, nicht der
  ganze Bestand noch einmal.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** Die **Inventur** liegt vor: je gepinntem Fremd-Artefakt **Ort**
      (im Repo / außerhalb), **Bindung** (Digest · Commit-SHA · Manifest),
      **Prüfzeitpunkt** (jeder Lauf · einmalig · nie) — aus der Konfiguration
      bzw. dem Skript gelesen, nicht aus der Prosa darüber
      ([`AGENTS.md`](../../../../AGENTS.md) §5, `seit slice-213`).
- [x] **(2)** Wo die Antwort in der Sensor-Beschreibung fehlt, steht sie dort —
      mit dem Zeitpunkt, nicht nur mit der Bindung.
- [x] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Zuerst die Form des Gegenstands** ([`AGENTS.md`](../../../../AGENTS.md) §5,
`seit slice-210`), denn DoD (1) ist eine Messung.

**Was zählt als *gepinntes Fremd-Artefakt*?** Etwas, das (a) **nicht** in
diesem Repo entsteht, (b) an einer **festen Kennung** bezogen wird
(Digest, Commit-SHA, Tag+Manifest), und (c) in einen **Lauf** eingeht — nicht
jede Abhängigkeit überhaupt. Go-Module fallen damit heraus: Sie sind über
`go.sum` gebunden und werden von der Toolchain bei **jedem** Bau geprüft; sie
sind kein offener Punkt, sondern die Referenz-Antwort.

**Was zählt als *Prüfzeitpunkt*?** Der Moment, in dem etwas die Bindung
**tatsächlich nachrechnet** — nicht der, in dem sie in einer Datei steht. Ein
Digest im `Dockerfile` ist eine Deklaration; nachgerechnet wird sie von Docker
beim **Pull**, und beim Lauf aus dem lokalen Store.

**Der vermutete Bestand ist klein** — vier Klassen, und die Zahl gehört
geprüft statt vorausgesetzt: vier digest-gepinnte `FROM`-Zeilen im
[`Dockerfile`](../../../../Dockerfile), das digest-gepinnte a-check-Image in
[`a-check.mk`](../../../../a-check.mk), das digest-gepinnte semgrep-Image plus
sein **git-gepinntes Regelset**, und der vendorte Baseline-Baum. **Nur eines
davon liegt außerhalb des Repos.**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`harness/sensors/semgrep.md`](../../../../harness/sensors/semgrep.md) | update | der Anlass: Bindung stark, Prüfzeitpunkt einmalig |
| weitere Sensor-Dateien | update, wo die Antwort fehlt | Ergebnis von DoD (1) |

### Die Inventur (DoD 1)

**Gelesen aus Konfiguration und Skript**, nicht aus der Prosa darüber
([`AGENTS.md`](../../../../AGENTS.md) §5, `seit slice-213`). **Sieben**
Artefakte erfüllen die Form aus §3; die achte Zeile ist die Referenz-Antwort
und zählt nicht mit.

| Artefakt | Ort | Bindung | Prüfzeitpunkt |
|---|---|---|---|
| `golang`-Basis (`deps`) | außerhalb, Docker-Store | Digest im [`Dockerfile`](../../../../Dockerfile) | **jeder Image-Bau** (Docker rechnet den Digest nach) |
| `golangci-lint`-Image | außerhalb, Docker-Store | Digest im `Dockerfile` | **jeder Image-Bau** |
| `distroless`-Runtime | außerhalb, Docker-Store | Digest im `Dockerfile` | **jeder Image-Bau** |
| `a-check`-Image | außerhalb, Docker-Store | Digest in [`a-check.mk`](../../../../a-check.mk) | **jeder Lauf** |
| `semgrep`-Image | außerhalb, Docker-Store | Digest in [`tools/semgrep.sh`](../../../../tools/semgrep.sh) | **jeder Lauf** |
| **`semgrep`-Regelset** | **außerhalb, unter dem Nutzer-Cache** | **git-Commit-SHA** | **einmalig** — danach nur Verzeichnis-Existenz |
| vendorte Baseline | **im Repo**, `.harness/baseline/<tag>/` | `SHA256SUMS` (kommt mit dem Baum) | **jeder `make gates`** — Echtheit nur im Nachtlauf |
| *(Referenz: Modul-Abhängigkeiten)* | *außerhalb, Modul-Cache* | *Prüfsummen-Datei des Moduls* | *jeder Bau, über den read-only-Schalter im `Dockerfile`* |

**Zwei Artefakte fallen aus der Reihe, und zwar in entgegengesetzte
Richtungen.**

**Das `semgrep`-Regelset ist das einzige mit einmaligem Prüfzeitpunkt.** Der
Bezug ist der stärkste im ganzen Bestand — ein git-Commit-SHA ist ein Hash über
den **Baum**, nicht über eine mitgelieferte Liste —, aber er wird genau einmal
eingelöst. Danach prüft das Skript nur noch, ob das Verzeichnis **existiert**.
Wer den Cache danach ändert, wird von nichts bemerkt; er liegt außerhalb des
Repos, also sieht ihn auch kein Gate, das den Arbeitsbaum prüft.

**Die vendorte Baseline ist der Gegenfall** und in slice-212 gemessen: jeder
Lauf, aber die Bindung kann die Echtheit nicht beweisen, weil das Manifest mit
dem Baum kommt. **Zusammen ergeben die beiden die Achse dieses Slice** —
*stark gebunden, selten geprüft* gegen *schwach gebunden, oft geprüft*.

**Was gemessen ist und was nicht — die Grenze der Inventur.** Die
Docker-Zeilen tragen *„Docker rechnet den Digest nach"*, und **das ist in
diesem Repo nicht gemessen**, sondern die dokumentierte Eigenschaft des
Werkzeugs; dasselbe gilt für die Modul-Prüfsummen. Gemessen sind hier **Ort
und Bindung** (aus den Dateien gelesen) und der Prüfzeitpunkt **dort, wo ein
eigenes Skript ihn setzt** — beim `semgrep`-Regelset und bei der vendorten
Baseline. **Wo der Beleg von einem fremden Werkzeug kommt, steht das da**,
statt als eigene Messung aufzutreten; §6 führt genau das als Risiko.

## 4. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.
[slice-213](../done/slice-213-grenzen-liste-braucht-fremden-leser.md) liegt in
`done/` — die Regel, nach der DoD (1) misst, ist damit geschrieben.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Ergibt die Inventur, dass mehrere
  Artefakte eine **Verhaltens**-Änderung bräuchten statt einer Beschreibung,
  ist das ein ADR-Vorgang und dieser Slice zu klein geschnitten.
- `in-progress` → `open` (blockiert): Zeigt sich, dass der Prüfzeitpunkt eines
  Artefakts von hier aus **nicht feststellbar** ist (fremde Werkzeug-Innerei),
  ruht der Slice bis zum Entscheid, ob eine Messung oder eine benannte Grenze
  die Antwort ist.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) zu **jedem** Artefakt der
Inventur stehen Ort, Bindung und Prüfzeitpunkt im Slice, jeweils aus
Konfiguration oder Skript gelesen; (b) `make gates` ist grün und die
Sensor-Beschreibungen tragen, was ihnen fehlte.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Antwort könnte lauten: alles in Ordnung, nur unbeschrieben.** Dann ist
  der Slice eine Doku-Ergänzung und kein Fund — und das wäre ein **gutes**
  Ergebnis, das als solches dastehen muss. Wer eine Lücke sucht, findet eine;
  die Inventur ist gegen die Konfiguration zu führen, nicht gegen die
  Erwartung. — **Ausgang:** \<offen\>
- **Der Prüfzeitpunkt fremder Werkzeuge ist von hier aus schwer zu belegen.**
  Dass Docker einen Digest beim Pull nachrechnet, ist bekannt, aber nicht in
  diesem Repo gemessen; dass `go.sum` bei jedem Bau greift, ebenso. **Wo der
  Beleg fehlt, gehört das gesagt statt behauptet** — sonst ist die Inventur
  eine Aufzählung von Vermutungen mit Tabellen-Rahmen. — **Ausgang:** \<offen\>
- **Vierter Slice in Folge an Grenzen-Beschreibungen.** slice-212, slice-213
  und dieser berühren dieselbe Familie. Die Gefahr ist nicht Wiederholung,
  sondern **Selbstbezug**: ein Harness, der nur noch sich selbst beschreibt.
  Der Unterschied hier ist der Gegenstand — eine **Supply-Chain**-Frage nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.1, die zufällig in einer
  Grenzen-Zeile landet. — **Ausgang:** \<offen\>

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

**Zwei** Sub-Areas, und das ist neu gegenüber den drei Vorgängern:

- **`*`** (Repo-Default) — die Sensor-Beschreibungen unter
  [`harness/sensors/`](../../../../harness/sensors/), die geändert werden.
- **`tools/harness/`** (Kürzel `HARN`,
  [`MR-004`](../../../../harness/conventions.md#mr-004)) — **gelesen, nicht
  geändert.** Der Gegenstand von DoD (1) sind die Skripte und
  Konfigurationen, die den Bezug herstellen ([`tools/semgrep.sh`](../../../../tools/semgrep.sh),
  [`tools/harness/fetch-baseline-cache.sh`](../../../../tools/harness/fetch-baseline-cache.sh),
  [`Dockerfile`](../../../../Dockerfile), [`a-check.mk`](../../../../a-check.mk)).
  **Die Sub-Area wird geführt, weil sie den Gegenstand trägt, nicht weil sie
  ein Diff bekommt** — genau das verlangt das Inklusionskriterium, und §1
  schließt Verhaltens-Änderungen dort aus.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **38** Verzeichnisse über beide
Kürzel). Für `HARN` steht **kein** Eintrag im Register — das ist ebenfalls eine
Antwort und wird notiert. **Vier** Einträge unter `ALL` sind einschlägig:

- [`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
  (4×, seit slice-213 verkörpert) — **der Anlass**, und seine Regel gilt diesem
  Slice unmittelbar: DoD (1) liest Konfiguration und Skript, **nicht** die
  Prosa darüber. Genau daran ist die Vermutung *„der `semgrep`-Cache ist
  ungeprüft"* schon einmal zerbrochen — er ist per Commit-SHA gebunden.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (14×, Stand *gemischt*) — **der gefährlichste hier.** Die Inventur hat eine
  **Menge** (welche Artefakte zählen?) und eine **Aussage je Mitglied**
  (Prüfzeitpunkt). slice-212 hat gezeigt, dass die Gefahr in der zweiten sitzt:
  Dort war die Menge trivial und die Antwort je Datei falsch.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (8×, Ausgang *geplant*) — §6 führt es als Risiko: Dass Docker einen Digest
  beim Pull nachrechnet, ist **bekannt**, nicht in diesem Repo **gemessen**.
  Eine Inventur, die Vermutungen in Tabellenzellen schreibt, behauptet
  Prüfungen.
- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (4×, verkörpert) — deshalb schreibt §3 vorher aus, was als *gepinntes
  Fremd-Artefakt* und was als *Prüfzeitpunkt* zählt.

**Keiner der vier erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-08 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). **Beide sind für diesen Slice einschlägig, nicht
Routine:** Sie sind die einzigen Träger, die gepinnte Fremd-Artefakte gegen
upstream halten — und ihr Grün ist von gestern. **Genau das ist die
Prüfzeitpunkt-Frage, die dieser Slice stellt**, angewandt auf seine eigene
Vorprüfung.

**Modus-Begründungsblock.** Beide berührten Sub-Areas GF — ein Block je
Sub-Area, der zweite kurz, weil dort nichts geändert wird.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch — die Form der Sensor-Dateien gibt die vendorte
  `gate.template.md` vor, und seit slice-213 sagt
  [`AGENTS.md`](../../../../AGENTS.md) §5, wie eine Grenze zu prüfen ist.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** **mittel.** Der Bestand liegt offen; das
  Risiko sitzt in der Antwort je Artefakt, und §6 führt es.
- **Reconciliation-Aufwand:** keiner (GF).

### Sub-Area: `tools/harness/`

- **Modus:** GF ([`MR-004`](../../../../harness/conventions.md#mr-004)).
- **Konventions-Dichte:** hoch — die Skripte sind über
  [`MR-004`](../../../../harness/conventions.md#mr-004),
  [`MR-005`](../../../../harness/conventions.md#mr-005) und
  [`MR-042`](../../../../harness/conventions.md#mr-042) konventionsgetragen.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig** — die Sub-Area wird **gelesen**,
  nicht geändert; eine Diskrepanz zwischen Skript und Doku ist genau der Fund,
  den der Slice sucht, und kein Risiko seines Vorgehens.
- **Reconciliation-Aufwand:** keiner (GF).
