# Welle 85 — Vier Kurs-Wellen, alle vier die eigene Bitte, und viermal eine Zahl ohne ihre Vorschrift — Closure-Notiz

**Welle:** welle-85-baseline-v5120-migration
**Abschluss:** 2026-08-26
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief · Steering-Loop · Register-Lese-Schritt.

Der Baseline-Pin steht auf `v5.12.0`. Der vendorte Baum ist gewechselt (51
Dateien, Manifest offline verifiziert), alle **lebenden** Verweise sind gehoben,
[`MR-030`](../../../../harness/conventions.md#mr-030) ist aufgelöst,
[`MR-037`](../../../../harness/conventions.md#mr-037) trägt die Hebung, und der
Alt-Baum ist entfernt.

**Das Delta ist gemessen, nicht gezählt:** von 52 Bundle-Dateien unterscheiden
sich 29 — eine ist das Manifest, 22 tragen ausschließlich den Versions-Stempel,
eine trägt Stempel plus `**Stand:**`-Zeile. **Fünf** tragen echten Regel-Inhalt,
und sie decken sich exakt mit der Lieferzusage der Kurs-Antwort.

**Alle vier Kurs-Wellen sind die Antwort auf den Konsumenten-CR dieses Repos.**
Zwei tragen Handlung, zwei sind belegt folgenlos — die Einzelheiten stehen im
Closure-Body von [slice-149](../done/slice-149-baseline-v5120-delta-audit.md).

Zwei Slices der Welle sind geschlossen ([slice-148](../done/slice-148-baseline-v5120-vendoring.md),
[slice-149](../done/slice-149-baseline-v5120-delta-audit.md)), dazu die vom Audit
geschnittene Etappe C-1 ([slice-150](../done/slice-150-pin-gebundene-zitate.md)).

## Was hat funktioniert?

**Der Frische-Sensor hat die Welle ausgelöst, nicht ein Mensch.**
`make baseline-freshness` meldete den neuen Release mit Exit 3 und bestätigte im
selben Lauf, dass der gepinnte Stand upstream unverändert ist — eine reine
Fortschreibung, kein Nachholen stiller Drift.

**Messen vor Schreiben hat zweimal getragen.** Das Wellen-Ziel entstand **nach**
dem Delta, nicht davor; und der Freshness-Audit hat alle 16 Adaptionen gefragt,
bevor er auf eine verengte.

**Die Reviews haben in jedem Slice etwas gefunden, das kein Gate sieht** — und
in zweien davon genau die Zusage, die ich am gründlichsten geprüft zu haben
glaubte.

## Was ging anders als geplant?

**Viermal war eine Zahl richtig gemessen und falsch gerahmt.** Der Zensus der
vierten Spiegel-Klasse nannte zwei Dokumente statt drei, weil der **Skopus**
nicht dabeistand; die dritte Spiegel-Klasse nannte „24" ohne
**Messvorschrift** (der Reviewer kam auf 60, 22 und 14); die Zitat-Zusage in
einem Adaptions-Eintrag war breiter als ihre Messung; und der `citations`-Zensus
nannte sechs Dateien statt zehn und einen Bruchpfad statt zwei. **Eine Zahl ohne
ihre Vorschrift ist keine Messung, sondern eine Behauptung mit Ziffern.**

**Ein Filter hat die falsche Frage gestellt.** Der Freshness-Audit fragte, wohin
der `Ersetzt-Baseline-Regel`-Zeiger zielt; der Kanon fragt, ob die neue Fassung
regelt, **wofür der Eintrag angelegt wurde**. Für vier Einträge konnte der
Filter strukturell nie „betroffen" liefern.

**Und einmal wurde eine Regel in ihrem eigenen Einführungsfall gebrochen.** Der
erste Entwurf zur Zitat-Frage stand auf einer zu engen Kanon-Suche, sagte
*„ergänzt, nicht ersetzt"* — und entfernte im selben Commit eines von zwei
veralteten Zitaten. Der Rückbau folgt der Form, die der Kanon dafür vorsieht:
[`MR-038`](../../../../harness/conventions.md#mr-038) ist **aufgelöst, nicht
korrigiert**; [`MR-039`](../../../../harness/conventions.md#mr-039) tritt an seine
Stelle.

## Steering-Loop-Einträge

- **Der Pin-Bump hat eine vierte Spiegel-Klasse bekommen** — den *zitierenden*
  Verweis. [`BEO-008`](../observations.md) führt sie samt der einzigen tragfähigen
  Messform: Zitat gegen die Zeilen, die zwischen den **zwei gepinnten Bäumen**
  verschwunden sind. Ein Korpus-Test isoliert sie nicht.
- **[`MR-039`](../../../../harness/conventions.md#mr-039)** hält fest, wohin das
  Delta gehört: in den **Bump-Eintrag**, nicht in das zitierende Dokument.
- **[`MR-037`](../../../../harness/conventions.md#mr-037)** trägt die
  Spiegel-Klassen **mit Zahl und Messvorschrift** — die Form, die den vier
  Rahmungs-Fehlern dieser Welle abgeht.
- **Nicht** mechanisiert: die Klasse *behauptete Prüfung, die nicht stattfand*
  ([`BEO-009`](../observations.md), Zähler 5) und *rotes Gate aus dem falschen
  Grund* ([`BEO-017`](../observations.md), Zähler 3). Beides sind Urteile über
  Aussagen, kein `grep`. Das ist Ergebnis, nicht Rest.

## Beobachtungs-Register (Zeiger)

Lese-Schritt über die Bewegungen dieser Welle:

| Eintrag | Stand | Was daraus folgt |
|---|---|---|
| [`BEO-008`](../observations.md) | 3 → **4** | Vierte Klasse benannt; die mechanische Form existiert im Produkt (`citations`) und ist blockiert — [slice-152](../next/slice-152-citations-scharfschalten.md) |
| [`BEO-009`](../observations.md) | 3 → **5** | Beide neuen Instanzen Richtung (a); dabei die eigene **Beleg-Form-Lücke** des Eintrags benannt statt geglättet |
| [`BEO-017`](../observations.md) | neu, **3** | Der Kanon trägt die Regel seit `v5.12.0` — auf unseren CR hin; keine Regelzeile hier, aber eine Prozedur |
| [`BEO-012`](../observations.md) | 4 | Unverändert in dieser Welle; die Feedforward-Hälfte wartet in [slice-147](../done/slice-147-reviewer-anker-reichweite.md) |

## Folge-Slices

**Zwei Slices werden aus der Welle herausgelöst**, weil ihr Gegenstand die
Pin-Hebung nicht mehr berührt — sie sind ab jetzt **wellenlos**
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht: ihre
Closure-Bedingung geht nicht über die eigene DoD hinaus):

- [slice-151](../done/slice-151-urteilsfreie-haelfte-voll.md) — die urteilsfreie
  Hälfte so weit ziehen, wie `modul-05` sie seit `v5.12.0` benennt. Der Kurs hat
  das ausdrücklich als unsere Entscheidung und **kein Konformitätsthema**
  bezeichnet.
- [slice-152](../next/slice-152-citations-scharfschalten.md) — `citations`
  scharfschalten. Der Blocker ist älter als diese Welle (Design-Review des
  Moduls, 2026-07-18).

**Der einzige offene Punkt aus der Migration selbst ist keiner:** die zwei
veralteten Zitate in [`MR-033`](../../../../harness/conventions.md#mr-033) bleiben
dort stehen — sie sind die historisch korrekte Aussage über den Stand, gegen den
der Eintrag geschrieben wurde, und ihr Delta steht in
[`MR-039`](../../../../harness/conventions.md#mr-039).

## Verifikation

- `make gates` Exit 0 (zehn Glieder) auf jedem Slice-Commit dieser Welle.
- `make fullbuild` Exit 0 (fünf Closure-Glieder).
- `bash tools/harness/fetch-baseline-cache.sh --verify` Exit 0 — 51 Dateien,
  Manifest über beide Bäume vollständig.
- CI grün auf jedem gepushten Stand.
- Drei unabhängige Reviews, alle blockierend: ein HIGH, sechzehn MEDIUM, fünfzehn
  LOW — jeder Befund nachgeprüft und eingearbeitet.
