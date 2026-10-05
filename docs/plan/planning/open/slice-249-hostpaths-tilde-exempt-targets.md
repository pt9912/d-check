# slice-249: `hostpaths` erkennt Home-relative Pfade — samt Ziel-Ventil `exempt-targets`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst; das
Release danach ist ein Meilenstein-Vorgang, keine Wellen-Bedingung.

**Bezug:** [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)
(Scope — die Anforderung wird erweitert, kein neues Kürzel: Einzelmodul-Frage
nach dem Schnitt-Kriterium von [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md));
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Byte-Identität ohne den neuen Schlüssel — gilt für das Ventil, **nicht** für
die Tilde-Erkennung, die bewusst schärft); Vorbild der Konfigurations-Weitung
[ADR-0058](../../adr/0058-konfigurations-flaechen-additiv-weiten.md).

**Berührte Spec-Stellen:** `spezifikation.md` [§DC-FA-HOST-001.a](../../../../spec/spezifikation.md#dc-fa-host-001a--host-pfad-erkennung) (Muster,
Ventil), [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--d-checkyml)
(Schema-Zeilen `hostpaths.*`).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ausgangslage (gemessen):** heute trifft das Unix-Muster eine Tilde-Angabe
nur zufällig — steht hinter `~/` ein Präfix-Name, meldet es den Rest ab dem
zweiten `/` (die Tilde ist keine Wortgrenzen-Sperre), sonst gar nichts. Ein
Probelauf über diesen Plan meldete genau so einen abgeschnittenen Treffer.

**Ziel:** `hostpaths` meldet Home-relative Pfade (`~/` gefolgt von einem
ersten Segment, das **nicht** mit `.` beginnt) als Maschinen-Layout-Leak, in
voller Form und einmal — das Unix-Muster sperrt dafür `~` als Vorgänger;
Werkzeug-Konventionen unter `~/.<name>` bleiben still, und repo-spezifische
Ausnahmen trägt ein Ziel-Ventil `hostpaths.exempt-targets` (Globs über den
gefundenen Pfad, Unix- und Tilde-Funde).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Zeilen-Marker `d-check:ignore` für `hostpaths`** — anderer Vorgang mit
  anderer Abwägung: das Lastenheft schließt ihn ausdrücklich aus, und das
  Ziel-Ventil deckt den Anlass (wiederkehrende legitime Pfade) zentral statt
  verstreut. Der Satz „keinen Opt-out-Marker" bleibt für den Zeilen-Marker
  stehen.
- **`~user/…` (Tilde mit Benutzername)** — Schicht-/Erkennungsgrenze: selten,
  und in Prosa nicht verlässlich von anderer Tilde-Verwendung zu trennen. Wird
  als benannte Grenze ins Lastenheft geschrieben, nicht erkannt.
- **`$HOME`, `%USERPROFILE%` und andere Variablen-Formen** — kein belegter
  Anlass (Bestand gemessen: null Vorkommen); ohne Fall wäre es Vorsorge.
- **Ventil für Windows-/UNC-Funde** — Bestand bleibt bewusst stehen: beide
  Muster sind laut Lastenheft fest, und Backslash ist im Glob (Go
  `path.Match`) das Escape-Zeichen — ein Ventil dort bräuchte eine eigene
  Lexik.
- **Benutzerhandbuch, README, CHANGELOG** — Release-Prep
  (`AGENTS.md` §5 Regel 17), nicht dieser Slice.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Liefer-Punkte (3):

- [ ] Tilde-Erkennung: Lastenheft (Erweiterung, Versions-Bump +
      Historie-Zeile nach [`MR-032`](../../../../harness/conventions.md#mr-032)),
      Spezifikation [§DC-FA-HOST-001.a](../../../../spec/spezifikation.md#dc-fa-host-001a--host-pfad-erkennung), Kern-Regel; Tests Happy (`~/.claude/…`
      still), Negative (Tilde vor einem Nicht-Punkt-Segment gemeldet, in voller Form statt des heute abgeschnittenen Präfix-Treffers), Boundary (`~/` allein,
      `~user/`, URL-Tilde, Fence) — der Negative-Test lief ohne die Änderung
      aus dem richtigen Grund rot (Bewusstes Brechen, Modul 11).
- [ ] Ziel-Ventil `hostpaths.exempt-targets`: Schema-Zeile, Validierung am
      Config-Rand (leeres/ungültiges Glob ⇒ Exit 2), Kern-Prüfung, `--print-config`-Gerüst;
      ohne den Schlüssel byte-identischer Befundsatz für Unix-Funde.
- [ ] `hostpaths.prefixes`-Eintrag mit `~` ⇒ Exit 2 mit Hinweis auf die
      eigene Tilde-Erkennung (heute still wirkungslos: er würde zu `/~/`).
- [ ] ADR für die Entscheidungen (Punkt-Regel, Ventil auf das Ziel statt
      Zeilen-Marker, Schärfung ohne Byte-Identität); ADR-Index ergänzt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | die Anforderung erweitert (Tilde, Ventil, Grenze `~user`), Akzeptanzkriterien, Bump + Historie |
| `spec/spezifikation.md` | update | Host-Pfad-Erkennung Muster + Ventil; Schema-Zeilen; Historie |
| `internal/hexagon/core/rules/hostpaths.go` | update | Tilde-Muster, Ventil |
| `internal/hexagon/core/model/config.go` | update | `ExemptTargets` |
| `internal/adapter/driven/configyaml/configyaml.go` | update | Schlüssel, Validierung (Glob, `~`-Präfix) |
| `internal/adapter/driving/cli/config_template.go` | update | Gerüst-Zeile |
| `internal/hexagon/core/rules/hostpaths_test.go`, `configyaml_test.go` | update | Happy/Negative/Boundary nach den Akzeptanzkriterien der Anforderung |
| `docs/plan/adr/0098-…`, `docs/plan/adr/README.md` | neu / update | Entscheidungen |

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
Lastenheft (die Anforderung aus dem Bezug); Spezifikation (Host-Pfad-Erkennung), Schema-Zeile
`hostpaths.prefixes`, Grund-Code-Zeile `hostpath-forbidden`;
`--print-config`-Gerüst; Benutzerhandbuch §5/§6 und `operations.md`
(Release-Prep); `suggest.go` (Modulset — nur Aktivierung, keine Semantik,
unberührt).

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; slice-248 ist
geschlossen, `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): zeigt der Bestand (eigenes Repo,
  Schwester-Repos) so viele Tilde-Funde, dass die Punkt-Regel die
  Fehlalarme nicht trägt und eine zweite Erkennungs-Achse nötig wird, wird
  der Zuschnitt auf das Ventil verengt.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Die Tilde-Erkennung schärft ein bestehendes Modul: Konsumenten mit
  `~/<Verzeichnis>/…` in Prosa werden nach dem Update rot. Das ist der
  Zweck, aber kein additiver Schritt — CHANGELOG und Handbuch müssen es
  sagen (Release-Prep). — **Ausgang:** *(offen)*
- Die Wortgrenze der Tilde: `~` steht in Prosa auch als „ungefähr"
  (`~5 %`). Ohne folgenden `/` greift das Muster nicht, die Tests müssen den
  Fall dennoch tragen. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das Produkt
samt Spec-Stratum (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** zwei offene Einträge der
Sub-Area treffen den Gegenstand —
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/state.md)
(ein neues Muster braucht Negativfälle: `~5 %`, `~user/`, URL-Tilde — in
§2 als Boundary aufgenommen) und
[`BEO-ALL/semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/state.md)
(die Spiegel in §3 vorab gelistet).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-05 — `image-scan` grün; `upstream-drift` rot mit drei
Fremd-Meldungen (Baseline v6.14.0 verfügbar, semgrep 1.179.0, golang-Digest
unter 1.27.1 neu gebaut). Keine berührt diesen Gegenstand; je eigener
Pin-Vorgang.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
