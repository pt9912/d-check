# Welle 68 — Planning-/Roadmap-Harness vollenden — Closure-Notiz

**Welle:** welle-68-planning-roadmap-harness
**Abschluss:** 2026-08-09
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **`## Aktuelle Welle` erreicht die Template-Struktur-Form** (slice-092): eine
  **aktive** Welle trägt die Baseline-Felder (Welle-ID · Start · Geplantes Ende ·
  Closure-Trigger) ohne Ruhe-Marker, `planning-check` bleibt grün. Damit war der
  Abschnitt template-konform **ohne** Modul-Umbau; der Ruhe-Marker ist auf den
  wellenlosen Zustand zurückgeschnitten (Konventionsspeicher-Eintrag verfeinert).
- **Der Closure-Note-Qualitäts-Nachlauf ist mechanisiert** (slice-093,
  Etappe-D-Finding D-7). Das Modul `planning` trägt eine zweite Fähigkeit
  ([`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  geschärft): opt-in über `closure.dir` wird je abgeschlossenem Slice der
  Closure-Notiz-Abschnitt strukturell geprüft — drei Grund-Codes für drei
  verschiedene Reparaturen. Akzeptanzkriterien je Fall grün.
- **Neue Anforderung
  [`DC-FA-CLI-012`](../../../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)**
  (`--config <datei>`): Herkunfts-Umschaltung der Konfiguration, die zwei
  disjunkte Prüf-Profile im selben Repo möglich macht.
- **`make verify-closure-notes`** als zweiter Closure-Bindepunkt (in
  `fullbuild`, nicht in `gates`) plus der inferentielle Schwester-Skill
  `closure-note-reviewer.md`. Beide Ziel-Formen der Baseline, die d-check bis
  dahin gar nicht hatte, sind damit besetzt.
- **Release v0.52.0** (Digest `sha256:412a6fd3…662c`), Lastenheft 0.50.0,
  Handbuch 1.43, [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md)
  `Accepted`.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Vorab messen statt Schwellen raten.** Die Frage „welche Mindest-Satzzahl?"
  wurde nicht geschätzt, sondern am eigenen Bestand gemessen (92/92 Notizen,
  Minimum 5 Satzende-Zeichen außerhalb Code). Die Schwelle 4 ergab sich daraus —
  und der Beleg zeigte zugleich, dass **kein** Retrofit nötig ist. Dieselbe
  Messung hat die Floskel-Liste gefiltert: zwei naheliegende Phrasen hätten
  Falschbefunde erzeugt und blieben draußen.
- **Zwei Reviews mit getrennten Linsen, parallel und ohne Kenntnis voneinander.**
  Vertrag und Code fanden **verschiedene** Klassen; der Code-Review fand den
  einzigen HIGH, der Vertrags-Review die Selbstwidersprüche zwischen den
  Dokumenten. Ein einzelner Reviewer hätte vermutlich eine der beiden Hälften
  verloren.
- **Die eigenen Gates haben zweimal echte Fehler gefangen**, nicht nur formale:
  der `AllReasons`-↔-§4-Lockstep verhinderte eine Spec, die dem Code vorauseilt,
  und der Handbuch-Harness lehnte einen nicht replaybaren Ausgabeblock ab.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile möglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Der Modul-Schnitt hat eine öffentliche CLI-Anforderung nachgezogen.** Beim
  Zuschnitt war „Fähigkeit im `planning`-Modul statt neues Modul" eine reine
  Modul-Frage. Erst beim Entwurf zeigte sich: weil die Konfiguration
  konventionell aus **einer** Datei kommt, läuft alles im Modul dort mit, wo das
  Modul läuft — der gewünschte Closure-Bindepunkt war nur über `--config`
  erreichbar. Konsequenz: die Welle lieferte eine Anforderung mehr als geplant.
- **Der Slice-Umfang wuchs am Review.** Aus einem HIGH und zehn MEDIUM wurden
  eine umgekehrte Vertragszusage (kandidatenfreies Verzeichnis ist jetzt
  fail-closed, [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md)
  Entscheidung 8), eine symlink-feste Wurzel-Grenze und acht neue
  Akzeptanzkriterien. Das ist der Zweck des Reviews — planbar war es nicht.
- **Ein Rückstand aus einer früheren Welle fiel nebenbei auf:** eine `Accepted`-ADR
  stand sieben Tage nicht im ADR-Index. Nachgetragen, aber die **Ursache** bleibt
  offen — kein Gate deckt die Vollständigkeit eines Registers. Als `BEO-001` im
  Beobachtungs-Register geführt statt in diese Welle gezogen (WIP-Limit).
- **Das geplante Wellen-Ende (2026-08-05) verstrich.** Die Schätzung stammte aus
  der Zeit, als slice-093 noch als Doc-Arbeit gedacht war; tatsächlich wurde es
  Produkt-Code samt Release. Kein Closure-Kriterium war betroffen.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`).

- — keine — (kein Register-Eintrag hat 3× erreicht; `BEO-001` steht bei 1×).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`../observations.md`](../observations.md).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

- — keine — . Der offene Kandidat aus dieser Welle (`BEO-001`, Register-Gate für
  die Richtung „Artefakt ⇒ registriert") ist **bewusst kein** Folge-Slice: er
  steht als Beobachtung im Register und wird bei der nächsten Slice-Planung
  gesichtet. Ihn hier zu nennen, ohne die Datei anzulegen, wäre dieselbe Klasse
  wie ein halluziniertes Gate.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle Slices der Welle (092, 093) liegen in `done/`.
- `make fullbuild` grün — Image-Hash
  `sha256:c42dfc4ba1f8388afd8c595848e6d9a4ac9a1adeef13cb8bee0ab634d0272b11`.
- `make gates` grün: 330 Datei(en) / 0 Befund(e), Coverage 94,30 % (Schwelle
  93 %), keine offenen Carveouts.
- `make verify-closure-notes` grün: 301 Datei(en) / 0 Befund(e) — das **Mehr**
  dieser Welle gegenüber den Slice-DoDs, repo-weit über den `done/`-Bestand
  gefahren, einschließlich der Closure-Notiz des Slice, der das Gate gebaut hat.
- `make completeness-check`: 47 Anforderung(en), 0 Waise(n).
- `make bench`: Median 795 ms (Kriterium < 5000 ms).
- Release-Pipeline zum Tag `v0.52.0` grün; das veröffentlichte, digest-gepinnte
  Image gegen dieses Repo gegengeprüft.
- **Trigger-Audit** (Schritt 2, alle drei Artefaktklassen): keine offenen
  Carveouts; keine stehengebliebene Gate-Reifestufe (die Coverage-Schwelle ist
  kalibriert, Ist über Soll); kein ADR-Re-Evaluierungs-Trigger eingetreten —
  geprüft wurden insbesondere die von dieser Welle berührten
  [ADR-0028](../../adr/0028-planning-lifecycle-modul.md) (Roadmap-Struktur
  unverändert, weiterhin genau eine Roadmap, git-basierte Lifecycle-Prüfung
  bleibt Domäne des VCS-Ports) und
  [ADR-0026](../../adr/0026-completeness-in-product-gate.md) (Bindepunkt-Policy
  unverändert; dessen Silent-Green-Trigger für leere Prüf-Mengen ist durch den
  neuen Nullmengen-Guard eher bestätigt als in Frage gestellt).
