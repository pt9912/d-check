# slice-267: `vcs` lässt einen reinen Pfad-Nachzug in immutablen Dateien durch

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Baseline `v6.17.0` · `regelwerk/modul-04-adrs.md` §Nachzug ist
keine Überschreibung (ein Referenz-/Pfad-Nachzug in einer `Accepted`-ADR ist
keine inhaltliche Überschreibung); Anlass: der Umzug von
`docs/user/releasing.md` (slice-268) braucht ihn für zwei `Accepted`-ADRs;
Auftraggeber-Entscheid 2026-10-09.

**Berührte Spec-Stellen:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
`AGENTS.md` §3.5.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das Modul `vcs` erkennt eine Änderung an einer immutablen Datei, die
nur das **Ziel** eines Markdown-Links ändert und Linktext wie übrige Zeile
gleich lässt, als Pfad-Nachzug und meldet sie nicht — opt-in, ohne den
Schlüssel byte-identisch. Dieses Repo schaltet es für die ADRs ein;
`AGENTS.md` §3.5 nennt den Nachzug neben `## Geschichte` und dem
Status-Übergang.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Umzug von `releasing.md`** — slice-268.
- **Ein Nachzug, der den Linktext oder andere Worte der Zeile ändert** — das
  wäre Inhalt; die Grenze liegt beim Ziel in den Klammern.
- **Der Template-Feld-Nachzug** (zweiter Fall im Regelwerk) — anderer
  Gegenstand, kein Anlass.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] `vcs` lässt einen reinen Link-Ziel-Nachzug durch (Schlüssel, Lastenheft,
      Spezifikation, Tests, die ohne die Änderung aus dem richtigen Grund rot
      sind); eine Änderung am Linktext oder an der übrigen Zeile bleibt ein
      Befund.
- [x] Ohne den Schlüssel ist die Ausgabe unverändert; `.d-check.yml` schaltet
      ihn für die ADRs ein; `AGENTS.md` §3.5 nachgezogen; `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Kern-Regel `vcs` | update | Pfad-Nachzug erkennen |
| Konfig-Modell, YAML-Adapter, `--print-config`-Vorlage | update | Schlüssel |
| `spec/lastenheft.md`, `spec/spezifikation.md` | update | Anforderung und Verfeinerung |
| `.d-check.yml`, `AGENTS.md` §3.5, `harness/README.md`, `harness/sensors/adr-check.md` | update | Nutzung in diesem Repo und die Beschreibung des ADR-Gates |
| neue ADR (Index in `docs/plan/adr/README.md`) | neu | das ADR-Gate lockert eine Prüfregel (`AGENTS.md` §3.6) |

*(Plan-Änderung vor dem Code-Commit: Dass `make adr-check` in diesem Repo
einen Link-Ziel-Nachzug durchlässt, lockert eine Prüfregel — nach `AGENTS.md`
§3.6 braucht das eine ADR, auch wenn die Baseline den Nachzug erlaubt. Sie
trägt die Entscheidung, den Schlüssel für die ADRs dieses Repos einzuschalten,
und ihre Grenze.)*

*(Plan-Änderung vor dem Code-Commit, Auftraggeber-Entscheid nach R6: Ein
erkanntes Link-Ziel wird nicht mehr geleert, sondern auf `Dateiname#Anker`
normiert — und nur, wenn es im jeweiligen Stand als Datei oder Verzeichnis
auflöst (BASE gegen den BASE-Baum, HEAD gegen den HEAD-Baum). Grund: sechs
Review-Runden fanden je neue Markdown-Formen, in denen Inhaltstext als
Link-Ziel geleert wurde; ein Inhaltswort löst nicht als Datei auf, damit ist
die Klasse strukturell geschlossen statt Form für Form. Die Rückführung aus
§4 (Zeilen-Paarung) greift nicht: die Normierung braucht keine Paarung. Der
Slice wächst damit über eine Review-Sitzung, ohne zurückgeführt zu werden
([`MR-066`](../../../../harness/conventions.md#mr-066)): Grund — der
Gegenstand ist ein gelockertes Gate, ein halber Stand wäre schlechter als
keiner; Ersatz-Form der Prüfung — jede Runde misst am gebauten Image beide
Richtungen, Umgehung und unveränderten Bestand. [ADR-0103](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) bekommt einen
Geschichte-Anhang; Entscheidung und Anlass bleiben.)*

## 4. Trigger

**Start** (`next` → `in-progress`): slice-265 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): der Diff-Vergleich im Modul kennt keine
  Zeilen-Paarung, die einen Nachzug von einer Ersetzung unterscheidet — dann
  zuerst die Paarung.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit dem nächsten Release
hinaus.

## 6. Risiken und offene Punkte

- **Nachzug als Tarnung** — eine Link-Ziel-Änderung kann die Aussage einer
  Entscheidung verschieben, wenn das neue Ziel etwas anderes ist. Das Modul
  sieht nur die Form. — **Ausgang:** *weiter offen* — nach der Normierung auf
  auflösende Ziele bleibt ein Nachzug auf eine gleichnamige andere Datei oder
  auf eine, die nur noch in BASE existiert;
  [`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/state.md)
  (1×).

## 7. Closure-Notiz

- **Was hat funktioniert:** `vcs.ignore-link-targets` lässt einen reinen
  Pfad-Nachzug in `Accepted`-ADRs durch: Ein Link-Ziel, das in BASE oder HEAD
  als Datei oder Verzeichnis auflöst, wird auf Marke, Dateiname und Anker
  normiert; Inhaltstext löst nicht auf und bleibt Drift. Was ein Link ist,
  beantwortet die Erkennung des Moduls `links`, dazu Filter, die nur in
  Richtung Drift wirken. Am Bestand gemessen: der Umzug von
  `docs/user/releasing.md` samt Nachzug in [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) und der von
  `harness/conventions.md` (46 Links in 21 ADRs) gehen durch, ein Nachzug auf
  ein anderes Dokument, auf ein fehlendes Ziel und eine Linktext-Änderung sind
  Drift, das Löschen einer verlinkten Datei ohne ADR-Änderung bleibt still.
  Ohne den Schlüssel in 40 Läufen byte-identisch; die Verifikation fuhr zehn
  Mutationen, jede aus dem richtigen Grund rot.
- **Was ging anders als geplant:** Acht Review-Runden statt einer. R1 bis R6
  fanden je neue Markdown-Formen, in denen Inhaltstext als Ziel geleert wurde;
  jeder Fix schloss die gemeldete Form, einer (Code-Span-Überlappung) machte
  den Anlassfall selbst zur Drift (R4 H-1). Die Klasse schloss erst die
  Normierung auf auflösende Ziele (Auftraggeber-Entscheid nach R6). Ihr
  erster Stand löste BASE und HEAD gegen getrennte Bäume auf und meldete das
  Löschen einer verlinkten Datei als Drift in unveränderten ADRs — gefunden
  beim Bearbeiten von R7, gemessen (sechs falsche Befunde), behoben. Die
  zweite Plan-Änderung in §3 nennt noch die getrennten Bäume; gültig ist die
  Vereinigung (Verifikation V-1). Der Fix zu R8 M-1 (`00cc1f45`) schreibt nur
  seinen Befund und lief ohne eigene Review-Runde; die Verifikation hat ihn
  per Mutation bestätigt (V-2). Zwei Prozessfehler: eine gestagte
  ADR-Zeile reiste in den R2-Report-Commit mit, und die gepushte Zeile wurde
  ersetzt statt ergänzt (R3 M-2 bis M-4, append-only wiederhergestellt).
  [ADR-0103](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) war im Commit ihrer Entstehung `Accepted`; ihr Körper beschreibt
  die Leerung, vier Geschichte-Anhänge tragen die Differenz. Ein Ziel mit
  Query oder Prozent-Kodierung bleibt Drift und steht nur im Code-Kommentar
  (V-3).
- **Steering-Loop-Eintrag:** Workflow-Skelett, beide in
  `.claude/commands/implement-slice.md`: Schritt 19 geschärft — eine Frage,
  die das Produkt schon beantwortet, wird mit dessen Erkennung beantwortet
  (Auslöser `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`,
  viertes Auftreten nach der Verkörperung; kein mechanischer Sensor, weil die
  Vollständigkeit einer Negativliste ein Urteil ist), dazu Prüffrage 21 in
  `.harness/skills/reviewer.md`, liegt in `.harness/skills/reviewer.md`;
  Schritt 21 neu — vor jedem Commit `git diff --cached --stat` gegen die
  Botschaft halten (Auslöser `BEO-ALL/path-scoped-commit-carries-staged-rest`,
  slice-106, slice-108, slice-267, 3×), liegt in
  `.claude/commands/implement-slice.md`. Beide Ausgänge Auftraggeber-Entscheid.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-267.md` in
  [`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/state.md)
  (4×, verkörpert, geschärft),
  [`BEO-ALL/path-scoped-commit-carries-staged-rest`](../observations/BEO-ALL/path-scoped-commit-carries-staged-rest/state.md)
  (3×, verkörpert),
  [`BEO-ALL/fix-schliesst-pfad-nicht-klasse`](../observations/BEO-ALL/fix-schliesst-pfad-nicht-klasse/state.md)
  (2×),
  [`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/state.md),
  [`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md),
  [`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/state.md)
  und
  [`BEO-ALL/semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/state.md);
  neu
  [`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/state.md),
  [`BEO-ALL/adr-accepted-bevor-der-mechanismus-steht`](../observations/BEO-ALL/adr-accepted-bevor-der-mechanismus-steht/state.md)
  und
  [`BEO-ALL/range-leerfall-mit-lokalem-branchnamen`](../observations/BEO-ALL/range-leerfall-mit-lokalem-branchnamen/state.md)
  (je 1×; der letzte außerhalb des Gegenstands, ungeprüft).
- **Folge-Slices:** slice-268 (Umzug von `releasing.md`) — `make adr-check`
  lässt ihn jetzt durch; [ADR-0097](../../adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md) nennt den Pfad als Inline-Code, nicht als
  Link, das gehört in seinen Plan. Release v0.85.0 mit slice-263, slice-265,
  diesem Slice und slice-268.
- **Risiken aus §6:** eines weiter offen (Register, siehe §6). Trigger-Audit:
  kein Carveout, kein bootstrap-aware Gate; [ADR-0103](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) neu (Re-Evaluierungs-
  Trigger nicht eingetreten); keine Hard Rule mit eingetretenem Trigger.
  Nachtlauf-Stand ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — `.claude/commands/implement-slice.md` trägt
  `seit slice-267` in den Schritten 19 und 21, `.harness/skills/reviewer.md`
  in Prüffrage 21; (b) Folge-Slices — slice-268 liegt in `open/`;
  (c) Register — die zehn zitierten Beobachtungen existieren und tragen Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** geändert werden Produkt-Kern (`vcs`),
Konfig-Modell, YAML-Adapter, Konfig-Vorlage, Lastenheft, Spezifikation und die
Harness-Doku des ADR-Gates — alle unter dem Default `*` (`ALL`); deklariert.

**Spiegel vor dem Editieren** (Schritt 17; gemessen mit
`grep -n "Geschichte" AGENTS.md harness/README.md harness/sensors/adr-check.md .d-check.yml`
und `grep -rln` nach `VCSConfig` und der Kennung der Anforderung über Code und Doku): was als
erlaubte Änderung einer `Accepted`-ADR genannt ist, steht in `AGENTS.md`
§3.5, `harness/README.md` §Traceability rules, `harness/sensors/adr-check.md`
(Vertrag und Grenze) und im Kommentar über dem `vcs`-Block der `.d-check.yml`;
der Schlüssel selbst in Modell, Adapter, `--print-config`-Vorlage, Lastenheft
und Spezifikation. Das Modul `immutable` ist in diesem Repo nicht aktiv —
kein Spiegel hier. README und Handbuch ziehen die Release-Prep nach.

**Entwurf vor dem Code:** `vcs.ignore-link-targets` (bool, Default aus):
beim Vergleich des Core wird das Ziel jedes Inline-Links und jeder
Referenz-Definition auf eine leere Form gebracht — `[text](ziel)` wird
`[text]()`, `[label]: ziel` wird `[label]:`. Eine Änderung, die nur Ziele
ändert, ergibt denselben Core; jede Änderung am Linktext, an der übrigen Zeile
oder das Hinzufügen und Entfernen eines Links bleibt ein Befund.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** gelesen am 2026-10-09.
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/state.md)
(verkörpert als Schritt 19) — die Link-Erkennung ist ein Erkennungsmuster:
Negativliste vor dem Code (Link im Inline-Code, im Codeblock, Bild,
verschachtelte Klammern, Ziel mit Titel); [`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/state.md)
— die Zusage gilt der Core-Bildung, also auch dem gestagten Lauf;
[`BEO-ALL/review-fix-applied-only-at-cited-site`](../observations/BEO-ALL/review-fix-applied-only-at-cited-site/state.md)
(verkörpert als Schritt 20).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-09 aus dem jüngsten Lauf (`make nightly-state`) —
`upstream-drift` grün (2026-10-09 07:17 UTC, nach den Pin-Hebungen),
`image-scan` grün (2026-10-08 10:38 UTC).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
