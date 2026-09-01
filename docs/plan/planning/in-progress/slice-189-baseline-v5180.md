# Slice slice-189: Die Baseline steht auf v5.18.0

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist ein **Upstream-Release**, keine Welle.

**Bezug:**
[`MR-021`](../../../../harness/conventions.md#mr-021) (In-Repo-Verweise sind
pin-gebunden — die Pfad-Hälfte des Bumps);
[`MR-051`](../../../../harness/conventions.md#mr-051) (die `cite`-Spannen,
die zweite pin-gebundene Größe);
[`MR-055`](../../../../harness/conventions.md#mr-055) (die Symlink-Aliase,
die denselben Pin binden);
[`MR-013`](../../../../harness/conventions.md#mr-013) (Lifecycle-Move-Commit
bündelt gekoppelte Verweise — Gegenstand der zweiten Hälfte dieses Slice);
[`MR-057`](../../../../harness/conventions.md#mr-057) (meldet die
Kollision mit dem Kanon, jetzt auflösbar).

**Berührte Spec-Stellen:** — (Adoptions-Stand einer externen Konvention;
keine Produkt-Anforderung berührt).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-01.

---

## 1. Ziel

**Der vendorte Baseline-Baum steht auf `v5.15.0`, upstream stehen drei
Releases weiter.** Gemessen mit `make baseline-freshness`: `v5.16.0`,
`v5.17.0`, `v5.18.0` — Content am gepinnten Tag ist **unverändert**
(`Bytes == vendored SHA256SUMS`). Reiner Currency-Rückstand, kein Drift.

**Das Delta ist gelesen, nicht angenommen** — netzlos aus dem Kurs-Klon
(`git diff v5.15.0 v5.18.0 -- lab/regelwerk lab/templates`): von 51
Bundle-Dateien sind 42 unverändert, 1 trägt nur den Versions-Stempel
(`regelwerk/README.md`), **8** tragen echten Inhalt —
`grundlagen-traceability.md` (+11), `modul-05-planning-harness.md` (+6),
`modul-06-roadmap.md` (+56/-7, der größte Einzel-Diff), `modul-08-agentenrollen.md`
(+6/-5), `modul-10-review-harness.md` (+8), `templates/README.md` (+4/-1) und
zwei **neue** Templates (`archiv-stub-slice.template.md`,
`archiv-stub-welle.template.md`). Tags↔Kurs-Wellen (Top-Eintrag je Tag, aus
`git show <tag>:CHANGELOG.md`): `v5.16.0` = Welle 109, `v5.17.0` = Welle 110,
`v5.18.0` = Welle 111.

**Zwei inhaltlich unabhängige Stränge im Delta:**

1. **Ein neuer Zeitdokumente-Archiv-Mechanismus** (Wellen 109 + 111): Modul-6
   §Wellen-Closure-Prozedur bekommt einen neuen Schritt 4 (*„Zeitdokumente
   archivieren"*, alte Schritte 4/5 werden 5/6) — geschlossene Slices und
   Review-Reports wandern beim nächsten Wellen-Closure in ein `done/<welle-id>/`-Archiv,
   an ihrer Stelle bleibt bei Slices ein gekürzter **Stub** (Volltext-Zeiger,
   Identität, `Hervorgegangen:`-Kennungen), Review-Reports wandern **ohne**
   Stub vollständig. Welle 111 korrigiert Welle 109 an einer Stelle: das
   Archivieren ist **kein Nachrüst-Zwang und kein Verbot** — ein bewusster,
   optionaler Vorgang je Repo.
2. **Kurs-Welle 110 löst die in [`MR-057`](../../../../harness/conventions.md#mr-057)
   gemeldete Kollision auf** — mit genau dem Beispiel unseres eigenen
   CR: *„Ein Adopter hat eine Kollision gemeldet statt sie aufzulösen"*.
   `grundlagen-traceability.md` §Herkunfts-Anker bekommt die fehlende
   Bedingung: *„Beide Commits gehören in denselben Push. Zwischen ihnen ist
   das Repo kurz rot — zulässig, solange dieser Zwischenstand nicht die
   Spitze eines Push wird."*

## 2. Vorgehen

1. **Re-vendorn und das Delta sichten.** `fetch-baseline-cache.sh v5.18.0`,
   dann `make baseline-verify` (Integrität, Manifest-Deckung,
   Alias-Auflösung).
2. **Alle Pfad-Verweise ziehen** ([`MR-021`](../../../../harness/conventions.md#mr-021)):
   gemessen bei der Anlage 44 lebende Dateien mit `baseline/v5.15.0`
   (47 inklusive der drei eingefrorenen Aussagen über die Vergangenheit),
   dazu die vier Symlinks unter `.claude/rules/`
   ([`MR-055`](../../../../harness/conventions.md#mr-055)).
3. **Die `cite`-Spannen neu ankern**
   ([`MR-051`](../../../../harness/conventions.md#mr-051)). Gemessen bei der
   Anlage: 28 lebende Direktiven außerhalb der eingefrorenen Verzeichnisse.
   Da der reale Regelwerk-Diff nur fünf Dateien mit echtem Inhalt trifft
   (§1), betrifft eine Zeilenverschiebung höchstens Direktiven **in genau
   diesen fünf** — der Rest verschiebt sich nicht und ist beim Ausführen
   trotzdem gegen den Datei-Diff zu prüfen, nicht anzunehmen.
4. **Der Adaptions-Review durch die Liste, nicht durch den Diff** — alle 33
   lebenden `MR`-Einträge, fünf Ausgänge je Eintrag. **Ein Ausgang ist schon
   bekannt:** [`MR-013`](../../../../harness/conventions.md#mr-013)/[`MR-057`](../../../../harness/conventions.md#mr-057)
   — siehe §2a.
5. **Die Bestands-Stichprobe fahren** (`AGENTS.md` §1, hängt nicht am Delta).
6. `make gates` und `make fullbuild`; **Review** und **Verifikation** als
   getrennte Läufe; Closure.

## 2a. MR-013/MR-057: Auflösung, nicht Ablösung

**Die naheliegende Lesart wäre falsch.** Kurs-Welle 110 endet mit *„Die Regel
gewinnt damit ihre fehlende Bedingung, und der Adopter kann seine Adaption
auflösen, statt ihr einen Nachfolger zu geben"* — das liest sich wie ein
Auftrag, [`MR-013`](../../../../harness/conventions.md#mr-013)s Bündelung
ersatzlos zu streichen. Nachgemessen am eigenen Text des Eintrags zeigt sich:
die Kollision und ihre Auflösung treffen nur **zwei** der drei gebündelten
Fälle.

- **Slice-Lifecycle-Move und Beanspruchung** (§Begründung des Eintrags): Der
  ursprüngliche Anlass — sichtbar bei slice-040, 2026-06-21 — war eine
  **Push-CI**, die auf dem reinen Move-Commit rot lief (`target-missing` +
  `make planning-check`). Genau das ist Kurs-Welle 110s Fall: ein
  Zwischen-Commit, der zur geprüften Spitze eines Push wurde. Die jetzt
  zitierbare Bedingung **bestätigt** unsere Praxis, statt sie zu widerlegen —
  wir pushen beide Commits ohnehin nie getrennt.
- **MR-/Wellen-Lifecycle-Move** trägt im selben Eintrag eine **andere,
  eigene** Begründung: *„eine … wandernde Datei trägt relative Verweise, die vom
  neuen Ort eine Ebene tiefer auflösen müssen — ein byte-reiner Move-Commit
  wäre `doc-check`-rot."* Das ist kein Push-Tip-Risiko, sondern eine
  **lokale** Zusage: `make hooks`s `pre-commit`-Hook lässt seit welle-79
  keinen roten Gate-Exit passieren, und `doc-check` (Modul `links`) meldet
  `target-missing`, sobald irgendein Dokument den alten (jetzt eine Ebene zu
  flachen) Pfad noch referenziert — **unabhängig davon, ob und wann gepusht
  wird.** Kurs-Welle 110 adressiert das nicht; sie löst eine Bedingung für
  Push-Sichtbarkeit, nicht für lokale Commit-Zulässigkeit.

**Ausgang: Auflösung der Meldung, keine Praxis-Änderung.** Die Kollision war
eine **Vermischung zweier Begründungen** unter einem Eintrag, nicht ein
Widerspruch, der eine Seite zum Weichen zwingt. Der Folge-Eintrag (§4)
trennt beide Begründungen explizit und zitiert Kurs-Welle 110 nur dort, wo
sie zutrifft.

## 2b. Ergebnis der Bestands-Stichprobe (ausgeführt 2026-09-01)

**Gezogen:** [`modul-11-verification.md` §Fitness Function ohne Standard-Tool](../../../../.harness/baseline/v5.18.0/regelwerk/modul-11-verification.md#fitness-function-ohne-standard-tool-modul-11).
**Auswahl nach Kanon-Kriterium** — rotierend gegen die Stichproben der beiden
vorigen Bumps (`modul-14` beim `v5.6.0`-Bump, `modul-07` bei slice-183);
`modul-11` trägt in diesem Delta keinen Diff und hatte noch keine Stichprobe.

**Befund: konform, keine Diskrepanz.** Der Kanon beschreibt den Bau eines
selbstgebauten Sensors für eine ADR-Aussage ohne Standardwerkzeug
(Operationalisieren → Sensor-Schicht nach Kosten → Skript+Gate verdrahten →
inferentielle Schicht für Semantik). d-checks
[`verify-closure-notes`](../../../../AGENTS.md) ist genau dieses Muster: ein
Make-Target-Sensor (Modul `planning`/`structure`/`spans`) für die
strukturelle Hälfte, mit [`.harness/skills/closure-note-reviewer.md`](../../../../.harness/skills/closure-note-reviewer.md)
als der vom Kanon selbst benannten inferentiellen Schicht für die
Floskel-Erkennung — derselbe Prompt-Anker, den `modul-11` als Referenz
verlinkt. Beide Quadranten des „Hard Rule in zwei Quadranten"-Musters sind
besetzt: `AGENTS.md` §5 nennt die Closure-Notiz-Pflicht (feedforward), das
Gate prüft sie (feedback), und jeder Slice läuft `make fullbuild` (das
`verify-closure-notes` trägt) vor der Closure-Meldung — die vom Kanon
verlangte Selbstprüfung vor der „fertig"-Meldung.

**Nicht geprüft, weil nicht einschlägig:** die „Obermengen-Nachweis"-Pflicht
für den Fall, dass ein Standardwerkzeug den eigenen Sensor später ablöst —
kein solcher Kandidat ist derzeit sichtbar.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Adoption des Zeitdokumente-Archivs** (Wellen 109/111, neuer
  Schritt 4 der Wellen-Closure-Prozedur, zwei neue Stub-Templates). Welle 111
  macht es ausdrücklich optional; d-check fährt heute überwiegend wellenlos
  ([`docs/plan/planning/observations.md`](../observations.md) — kein Register-Treffer
  zu häufigem Wellen-Betrieb), der Adoptions-Aufwand (Modul 6/8/10 plus zwei
  Vorlagen) steht in keinem Verhältnis zum Nutzen für den aktuellen Betrieb.
  Eigener Folge-Slice, sobald ein Wellen-Closure ansteht, an dem sich das
  Archiv erproben lässt.
- **Kein Nachziehen der Wortlaute in `done/`.** Ein zitierter Wortlaut wird
  nicht rückwirkend umgeschrieben
  ([`MR-039`](../../../../harness/conventions.md#mr-039)).

## 4. Definition of Done

- [x] Der Pin steht auf `v5.18.0`: vendorter Baum re-vendored,
      `make baseline-verify` grün (Integrität, Manifest-Deckung,
      Alias-Auflösung, 53 Dateien), alle Pfad-Verweise und die vier Symlinks
      gezogen (44 lebende `v5.15.0`-Vorkommen gehoben, drei eingefrorene
      stehen gelassen).
- [x] **Jede `cite`-Direktive ist entschieden**, nicht nur grün: je Direktive
      steht fest, ob sie nachgezogen, umgehängt oder nach
      [`MR-039`](../../../../harness/conventions.md#mr-039) entfernt wurde.
      28 lebende Direktiven: 8 in den fünf real geänderten Dateien (6
      nachgezogen um +5/+6 Zeilen, 2 unverändert bestätigt — vor der
      jeweiligen Einfügestelle), 20 unverändert (Datei-Stempel-only oder
      unberührt), 0 entfernt.
- [x] **Der Adaptions-Review ist gefahren und dokumentiert:** je lebendem
      `MR`-Eintrag einer der fünf Ausgänge. Alle 33 lebenden Einträge geprüft:
      **alle 33 bleiben gültig** — kein Abschnitt umbenannt, keine Regel
      entfallen, keine Ergänzung berührt eine bestehende Zeile (siehe
      [`MR-058`](../../../../harness/conventions.md#mr-058)).
- [x] **Die Kollision aus §2a ist aufgelöst, nicht abgelöst:** 
      [`MR-058`](../../../../harness/conventions.md#mr-058) trennt die
      Push-CI-Begründung (Slice-/Beanspruchungs-Move) von der lokalen
      `doc-check`-Begründung (MR-/Wellen-Lifecycle-Move), zitiert Kurs-Welle
      110 nur für die erste, und markiert
      [`MR-057`](../../../../harness/conventions.md#mr-057) als aufgelöst
      (`git mv` nach `conventions/done/`).
      [`MR-013`](../../../../harness/conventions.md#mr-013) selbst bleibt
      inhaltlich unverändert in Kraft — kein Praxiswechsel, siehe §2a.
- [x] **Das Delta ist gelesen, nicht angenommen:** was `v5.16.0`, `v5.17.0`
      und `v5.18.0` tragen, steht im Slice (§1) — inklusive des
      Archiv-Mechanismus, den dieser Slice bewusst nicht adoptiert (§3).
- [ ] `make gates` und `make fullbuild` grün (Exit explizit); **unabhängiger
      Review**; **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.
      `make gates` grün: zehn Gates, 649 Dateien, 0 Befunde (2026-09-01).
      Review und Verifikation stehen noch aus.

## 5. Abnahme-Punkte / Risiken

- **Die Auflösung aus §2a ist Urteil, kein `grep`.** Die Versuchung ist,
  Kurs-Welle 110 pauschal als *„Adaption auflösen"* zu lesen und die
  MR-/Wellen-Lifecycle-Bündelung ersatzlos zu streichen — das bräche
  `make hooks`s `pre-commit`-Hook beim nächsten MR-Move (`doc-check`-rot vor
  jedem Commit). — **Ausgang: entfallen.** [`MR-058`](../../../../harness/conventions.md#mr-058)
  behält die Bündelung ausdrücklich bei und benennt den eigenständigen,
  lokalen Grund dafür — kein Praxiswechsel ist eingetreten.
- **Der Bump-Alarm bei kleinem Delta verleitet zur Unterschätzung.** Nur 8
  von 51 Dateien mit echtem Inhalt heißt nicht, dass der Adaptions-Review
  oberflächlich laufen darf — die Zahl der `MR`-Einträge (33) bleibt
  unverändert groß. — **Ausgang: entfallen.** Alle 33 Einträge sind einzeln
  gegen ihr `Ersetzt-Baseline-Regel`-Feld geprüft (nicht pauschal
  übernommen); das Ergebnis „alle bleiben gültig" ist durch die Prüfung
  begründet, nicht durch die kleine Delta-Zahl vorweggenommen.
- **Der optionale Archiv-Mechanismus könnte trotzdem einen `MR`-Eintrag
  berühren**, wenn ein bestehender Eintrag (z. B. zur Wellen-Closure-Prozedur)
  auf die jetzt verschobenen Schritt-Nummern 4/5 zeigt. — **Ausgang:
  entfallen.** Nachgemessen (`grep` über alle lebenden `MR`-Dateien nach
  „Wellen-Closure-Prozedur"/„Schritt 4"/„Schritt 5"): kein Treffer außer im
  neuen [`MR-058`](../../../../harness/conventions.md#mr-058) selbst.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice; `v5.18.0` ist upstream verfügbar (gemessen mit
`make baseline-freshness`).

**Rückführungen:** `in-progress` → `open`, falls das Delta einen Ausgang
**gegenstandslos** oder **teilweise überholt** ergibt, dessen Rückbau selbst
ein Slice ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.harness/baseline/` (vendorte Fremd-Konvention) und
  `harness/` (der Konventionsspeicher). Beide fallen unter den Default `*` =
  **Greenfield** ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Stand 2026-09-01, höchste Kennung
  `BEO-025`): [`BEO-009`](../observations.md) (**7**, siebte Instanz war
  slice-183 selbst — hier gilt dieselbe Lehre: eine im Vorgehen-Abschnitt als
  „gemessen" markierte Zahl gegen das tatsächliche DoD-Ergebnis prüfen, bevor
  der Slice schließt); [`BEO-012`](../observations.md) (**12**) — ein Zitat
  über seinen Geltungsbereich hinaus, einschlägig für §2a: Kurs-Welle 110
  gilt nur der Push-CI-Begründung, nicht der `doc-check`-Begründung, und die
  Auflösung muss diese Grenze halten, nicht pauschal zitieren; [`BEO-008`](../observations.md)
  (**4**) — die Spiegel-Klassen einer Pin-Hebung; [`BEO-002`](../observations.md)
  (**7**) — die Liste entsteht aus einem `grep` nach dem alten Wortlaut, nicht
  aus dem Gedächtnis; [`BEO-013`](../observations.md) (**1**) — die
  Bestands-Stichprobe des Freshness-Audits gehört gefahren, auch wenn der Pin
  aktuell ist; [`BEO-025`](../observations.md) (**1**) — ein Bump zerfällt in
  mehrere Commits, Zuordnung vor dem Editieren klären.

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — **bei der
  Beanspruchung neu gelesen:** `upstream-drift.yml` meldet weiter **ROT**,
  jetzt mit jüngstem Lauf 2026-09-01T05:56:48Z (ein zweiter Lauf seit der
  Anlage dieses Slice — derselbe Grund: unser Pin ist hinter `v5.18.0`
  zurück), `image-scan.yml` grün. Der Nachtlauf wird grün, wenn dieser Slice
  geschlossen ist. **Dieser Block trägt bewusst keine `cite`-Direktive** —
  sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-189. Betroffene IDs:
[`MR-013`](../../../../harness/conventions.md#mr-013),
[`MR-021`](../../../../harness/conventions.md#mr-021),
[`MR-039`](../../../../harness/conventions.md#mr-039),
[`MR-051`](../../../../harness/conventions.md#mr-051),
[`MR-055`](../../../../harness/conventions.md#mr-055),
[`MR-057`](../../../../harness/conventions.md#mr-057). Module: `links`,
`anchors`, `citations`. Gates: `make baseline-verify`, `make gates`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Ein Adoptions-Stand und sein
Konventionsspeicher; kein Produkt-Code, keine Reconciliation. Das
**Evidenz-Risiko** ist die einzige Achse mit Substanz: der vendorte Baum ist
Fremd-Inhalt, und ob eine Adaption noch trägt, entscheidet sein Delta — nicht
unser Bestand.

## 9. Closure-Notiz (nach `done/`)
