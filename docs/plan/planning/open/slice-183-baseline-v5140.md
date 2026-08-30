# Slice slice-183: Die Baseline steht auf v5.14.0

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist ein **Upstream-Release**, keine Welle.

**Bezug:**
[`MR-021`](../../../../harness/conventions.md#mr-021) (In-Repo-Verweise sind
pin-gebunden — die Pfad-Hälfte des Bumps);
[`MR-051`](../../../../harness/conventions.md#mr-051) (die `cite`-Spannen, die
zweite pin-gebundene Größe, samt den drei Fallunterscheidungen);
[`MR-055`](../../../../harness/conventions.md#mr-055) (die Symlink-Aliase, die
denselben Pin binden und von keinem Modul gescannt werden);
[`MR-039`](../../../../harness/conventions.md#mr-039) (ein zitierter Wortlaut
wird nicht rückwirkend umgeschrieben — der Ausgang für ein echtes Delta).

**Berührte Spec-Stellen:** — (Adoptions-Stand einer externen Konvention; keine
Produkt-Anforderung berührt).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-30.

---

## 1. Ziel

**Der vendorte Baseline-Baum steht auf `v5.12.0`, upstream stehen drei
Releases weiter.** Gemessen mit `make baseline-freshness`: `v5.13.0`,
`v5.13.1`, `v5.14.0` — und der Content am gepinnten Tag ist **unverändert**
(`Bytes == vendored SHA256SUMS`). Es ist also ein reiner **Currency**-Rückstand,
kein Drift: nichts ist falsch, aber die verkörperte Form beruft sich auf eine
Fassung, die drei Releases alt ist.

**Der Herausgeber hat den Inhalt für `v5.13.0` angekündigt** — fünf
Regelwerk-Dateien mit echtem Inhalt (`grundlagen-durchsetzungsschicht`,
`grundlagen-begriffe`, `modul-13`, `modul-14`, `modul-02`), dazu zwei Punkte,
die uns direkt betreffen: der Abschnitt `§Referenz-Implementierung` heißt jetzt
`§Das vollständige Artefakt-Set`, und die `MR`-Vorlage hat ihr `Status:`-Feld
verloren. **Was `v5.13.1` und `v5.14.0` tragen, wissen wir nicht** — das ist
der erste Schritt, nicht eine Annahme dieses Plans.

## 2. Vorgehen

1. **Re-vendorn und das Delta sichten.** `fetch-baseline-cache.sh v5.14.0`,
   dann `make baseline-verify` (Integrität **und** Manifest-Deckung **und**
   Alias-Auflösung). Erst danach steht fest, was die Schritte 3 und 4 zu tun
   haben — der Plan nimmt es nicht vorweg.
2. **Alle Pfad-Verweise ziehen** ([`MR-021`](../../../../harness/conventions.md#mr-021)):
   gemessen **76** Dateien mit `baseline/v5.12.0`, dazu die **zwei** Symlinks
   unter `.claude/rules/` ([`MR-055`](../../../../harness/conventions.md#mr-055)) —
   sie binden denselben Pin, werden von keinem Modul gescannt und brächen
   still. `make baseline-verify` fängt genau das.
3. **Die `cite`-Spannen neu ankern**
   ([`MR-051`](../../../../harness/conventions.md#mr-051)). Gemessen sind es
   **rund 33 echte Direktiven** in den vendorten Baum — die rohe Zahl 133 zählt
   Erwähnungen mit, darunter Doku-Beispiele wie `<pfad>`. Verteilung: 25 in
   `modul-05`, je 2 in `modul-09` und `grundlagen-durchsetzungsschicht`, je 1 in
   `modul-02`, `modul-10`, `modul-11` und den Templates. **Unterschieden wird
   nach Grund-Code, nicht nach Gefühl:** `citation-mismatch` mit dem Wortlaut
   anderswo in derselben Datei ⇒ Spanne nachziehen; Wortlaut nirgends mehr ⇒
   [`MR-039`](../../../../harness/conventions.md#mr-039), Direktive entfernen,
   Delta im Bump-Eintrag festhalten; `citation-out-of-range` ⇒ **zuerst** prüfen,
   ob der Abschnitt in eine andere Datei gewandert ist.
4. **Der Adaptions-Review durch die Liste, nicht durch den Diff.** Der Kanon
   verlangt ihn ausdrücklich, und zwar je Eintrag mit einem von **fünf**
   Ausgängen (gegenstandslos · bleibt gültig · teilweise überholt · Bezug
   entfallen · Schärfung). Gegenstand sind die **33 lebenden** `MR`-Einträge;
   die 23 in `conventions/done/` sind aufgelöst und nicht Gegenstand.
   **Ein Ausgang ist schon bekannt:** [`MR-048`](../../../../harness/conventions.md#mr-048)
   zitiert `§Referenz-Implementierung`; der Abschnitt heißt ab `v5.13.0`
   `§Das vollständige Artefakt-Set`. Damit schließt sich zugleich der
   `F-7`-Befund aus dem slice-155-Review.
5. **Die Bestands-Stichprobe fahren**, die
   [`AGENTS.md`](../../../../AGENTS.md) §1 auch bei aktuellem Pin verlangt —
   sie hängt nicht am Delta.
6. `make gates` und `make fullbuild`; **Review** und **Verifikation** als
   getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Entfernen der `Status:`-Felder.** Die `MR`-Vorlage des Herausgebers
  hat das Feld verloren, und unsere eigene Regel
  ([`AGENTS.md`](../../../../AGENTS.md) §3.7: der Zustand ist die
  Verzeichnis-Position) spricht fürs Entfernen. Es sind **56** Dateien; zusammen
  mit den 76 des Bumps wäre der Diff nicht mehr in einer Sitzung prüfbar. Der
  Herausgeber sagt ausdrücklich, dass daraus nichts folgt — also eigener Slice.
- **Keine Antwort auf die `--print-mk`-Frage.** Der Herausgeber fragt, ob eine
  **mount-freie** Recipe-Form ins Werkzeug gehört; ein Adopter führt sie als
  Adaption mit **unserem** Werkzeug als Auflösungs-Trigger. Das ist eine
  Produkt-Entscheidung mit ADR und Lastenheft-Bump, kein Bump-Nebenprodukt.
  **Dieser Slice macht sie aber erst beurteilbar** — die Kanon-Stelle, auf die
  sie sich beruft (`modul-14 §Der Prüflauf ist hermetisch`), liegt in unserem
  Baum heute **nicht** vor.
- **Kein Nachziehen der Wortlaute in `done/`.** Ein zitierter Wortlaut wird
  nicht rückwirkend umgeschrieben
  ([`MR-039`](../../../../harness/conventions.md#mr-039)); [`MR-051`](../../../../harness/conventions.md#mr-051)s
  Geltungsbereich nennt die **lebenden** Dokumente.

## 4. Definition of Done

- [ ] Der Pin steht auf `v5.14.0`: vendorter Baum re-vendored,
      `make baseline-verify` grün (Integrität, Manifest-Deckung, Alias-Auflösung),
      alle Pfad-Verweise und die zwei Symlinks gezogen.
- [ ] **Jede `cite`-Direktive ist entschieden**, nicht nur grün: je Direktive
      steht fest, ob sie nachgezogen, umgehängt oder nach
      [`MR-039`](../../../../harness/conventions.md#mr-039) entfernt wurde. Die
      Zahl der drei Fälle steht in der Commit-Botschaft.
- [ ] **Der Adaptions-Review ist gefahren und dokumentiert:** je lebendem
      `MR`-Eintrag einer der fünf Ausgänge. **Keine Treffer sind ebenfalls eine
      Antwort** und werden notiert.
- [ ] [`MR-048`](../../../../harness/conventions.md#mr-048) zeigt auf
      `§Das vollständige Artefakt-Set`; der `F-7`-Befund aus slice-155 ist als
      geschlossen vermerkt.
- [ ] **Das Delta ist gelesen, nicht angenommen:** was `v5.13.1` und `v5.14.0`
      gegenüber `v5.13.0` tragen, steht im Slice — der Herausgeber hat nur
      `v5.13.0` angekündigt.
- [ ] `make gates` und `make fullbuild` grün (Exit explizit); **unabhängiger
      Review**; **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Der Bump-Alarm ist gewollt, aber er unterscheidet nicht.**
  [`MR-051`](../../../../harness/conventions.md#mr-051) nennt den Preis: eine
  Direktive meldet **auch dann**, wenn nur die Zeilennummern gewandert sind. Bei
  33 Direktiven ist das Arbeit ohne inhaltlichen Anlass, und die Versuchung ist,
  sie mechanisch nachzuziehen statt den Grund-Code zu lesen. — **Ausgang:**
  *(bei Closure)*
- **Der Adaptions-Review ist Urteil, kein `grep`.** Fünf Ausgänge je Eintrag,
  33 Einträge — die Gefahr ist nicht, einen falsch zu entscheiden, sondern alle
  reflexhaft auf „bleibt gültig" zu setzen, weil das der Normalfall ist. —
  **Ausgang:** *(bei Closure)*
- **Drei Releases auf einmal.** Angekündigt ist nur `v5.13.0`; `v5.13.1` und
  `v5.14.0` sind unbekannt. Ein Sprung über drei Fassungen macht das Delta
  größer als den Diff, den ein Reviewer in einer Sitzung prüft. — **Ausgang:**
  *(bei Closure)*
- **Der Bump macht eine fremde Adaption beurteilbar, die auf uns zeigt.** Die
  `--print-mk`-Frage hängt an `modul-14 §Der Prüflauf ist hermetisch`; solange
  der Abschnitt nicht im Baum liegt, wäre jede Antwort darauf ein Urteil aus dem
  Zitat eines Dritten ([`BEO-012`](../observations.md)). — **Ausgang:**
  *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice; `v5.14.0` ist upstream verfügbar (gemessen mit
`make baseline-freshness`).

**Rückführungen:** `in-progress` → `open`, falls das Delta einen Ausgang
**gegenstandslos** oder **teilweise überholt** ergibt, dessen Rückbau selbst ein
Slice ist — dann trägt dieser Slice den Bump und der Rückbau folgt getrennt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.harness/baseline/` (vendorte Fremd-Konvention) und
  `harness/` (der Konventionsspeicher). Beide fallen unter den Default `*` =
  **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration). Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-30, höchste Kennung
  `BEO-024`): [`BEO-012`](../observations.md) — ein Zitat über seinen
  Geltungsbereich hinaus: der Adaptions-Review liest 33 Fremd-Regeln und
  entscheidet je Eintrag, ob sie eine Adaption ablösen; das ist genau die Lage,
  in der eine Reichweite überdehnt wird; [`BEO-002`](../observations.md) — eine
  Semantik-Änderung wird nur im Körper nachgezogen, die Ränder bleiben stehen:
  ein Bump ist die Rand-Operation schlechthin, und die Symlinks sind der Rand,
  den kein Modul scannt; [`BEO-013`](../observations.md) — ein Wächter, der
  nichts mehr fängt: eine `cite`-Direktive, deren Ziel verschwunden ist, wird
  nach [`MR-039`](../../../../harness/conventions.md#mr-039) entfernt und hinterlässt einen Wortlaut ohne Prüfung.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-30T06:08:17Z,
  `image-scan.yml` 2026-08-30T09:16:25Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-183. Betroffene IDs:
[`MR-021`](../../../../harness/conventions.md#mr-021),
[`MR-039`](../../../../harness/conventions.md#mr-039),
[`MR-048`](../../../../harness/conventions.md#mr-048),
[`MR-051`](../../../../harness/conventions.md#mr-051),
[`MR-055`](../../../../harness/conventions.md#mr-055). Module: `links`,
`anchors`, `citations`. Gates: `make baseline-verify`, `make gates`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Ein Adoptions-Stand und sein Konventionsspeicher;
kein Produkt-Code, keine Reconciliation. Das **Evidenz-Risiko** ist die einzige
Achse mit Substanz: der vendorte Baum ist Fremd-Inhalt, und ob eine Adaption
noch trägt, entscheidet sein Delta — nicht unser Bestand.

## 9. Closure-Notiz (nach `done/`)
