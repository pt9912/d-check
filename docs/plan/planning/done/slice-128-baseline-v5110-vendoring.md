# Slice slice-128: Etappe A — Bundle `v5.11.0` vendoren und den Pin heben

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-83-baseline-v5110-migration](../welle-83-baseline-v5110-migration.md)
(zugeordnet bei der Eröffnung).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011) (Pin auf
Release-Tag) und die Kette bis
[`MR-029`](../../../../harness/conventions.md#mr-029);
[`MR-021`](../../../../harness/conventions.md#mr-021) (pin-gebundene Verweise);
[`MR-023`](../../../../harness/conventions.md#mr-023) (Bundle-Layout);
Baseline-Regelwerk
[`modul-02-harness-bootstrap.md` §Freshness-Audit](../../../../.harness/baseline/v5.11.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2).

**Berührte Spec-Stellen:** — (Harness-Vendoring und Konventionsspeicher; keine
Anforderung, kein Spec-Stratum, kein Produkt-Code).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der vendorte Baum trägt `v5.9.0` (Kurs-Welle 86). Diese Etappe hebt ihn auf
`v5.11.0` (Kurs-Welle 94) — **nur den Pin und seine Verweise**, ohne eine
einzige Regel anzuwenden. Was der neue Stand inhaltlich verlangt, beantwortet
[slice-129](../open/slice-129-baseline-v5110-delta-audit.md); wer beides mischt, kann
später nicht mehr trennen, ob eine Änderung aus dem Bump oder aus dem Audit kam.

## 2. Vorgehen

1. **Bundle materialisieren** nach `.harness/baseline/v5.11.0/`
   (`{regelwerk,templates}/` + `SHA256SUMS`) über
   `tools/harness/fetch-baseline-cache.sh`; Integrität **offline** verifizieren
   (`--verify`). Kein Handanlegen an den entpackten Bäumen.
2. **Pin heben** in [`harness/conventions.md`](../../../../harness/conventions.md)
   §Baseline und §Adoptierte Konventions-Quellen; neuer `MR-`Eintrag als
   nächster Schritt der Pin-Serie, Index-Zeile ergänzt, Vorgänger nach
   `conventions/done/` mit seiner Nachfolger-Zeile.
3. **Pin-gebundene Verweise heben** ([`MR-021`](../../../../harness/conventions.md#mr-021)) —
   und zwar nach der **Drei-Klassen-Prüfung** aus `BEO-008`, nicht nur per
   Pfad-`grep`: Pfad-Verweise, Release-/Tree-**URLs** mit dem Tag, und
   **Prosa-/Ellipsen-Pins**. Dazu die Gegenprobe, ob eine Stelle über die
   **Gegenwart** oder über die **Vergangenheit** spricht — eine gehobene
   Vergangenheits-Aussage ist ein neuer Fehler.
4. **Alt-Baum entfernen** und prüfen, ob eingefrorene Verweise darauf zeigen
   (immutable ADRs, `done/`-Slices) — falls ja, quell-skopiert über
   `ignore-refs` ausnehmen statt die eingefrorene Doku zu editieren.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Regel des neuen Stands anwenden.** Auch nicht die
  Vollständigkeits-Zusage aus Kurs-Welle 94, obwohl wir eine Verletzung bereits
  kennen — sie gehört [slice-127](../next/slice-127-claude-md-pointer.md).
- **Kein Delta-Audit.** Das ist Etappe B.
- **Keine Zwischenstufe über `v5.10.0`.**

## 4. Definition of Done

- [x] `.harness/baseline/v5.11.0/` liegt vollständig (51 Dateien); `SHA256SUMS`
      offline verifiziert — Vendor-Lauf, `--verify` gegen den alten **und** den
      neuen Pin, je Exit 0. Der Reviewer hat zusätzlich `--check-latest` in
      **beiden** Teilen grün gefahren: `v5.11.0` ist der neueste Tag, und die
      vendorten Bytes sind identisch mit dem Upstream-Asset am Tag.
- [x] Pin gehoben, [`MR-030`](../../../../harness/conventions.md#mr-030) angelegt,
      beide Index-Tabellen nachgezogen. **Nach Review korrigiert:** die
      aufgelöste Zeile stand zwischen [`MR-026`](../../../../harness/conventions.md#mr-026) und [`MR-027`](../../../../harness/conventions.md#mr-027), weil der Anker den
      [`MR-028`](../../../../harness/conventions.md#mr-028)-**Link in der Zeile darüber** traf statt deren eigene.
- [x] Verweis-Hebung über alle drei Klassen belegt, samt
      Vergangenheits-Gegenprobe — **mit einem Miss, siehe §5.** Der Reviewer
      hat alle 34 Pfade und 5 URLs einzeln auf Zeitform geprüft: **keine**
      Über-Hebung.
- [x] Alt-Baum entfernt; kein Verweis läuft ins Leere — **mit einem benannten
      blinden Fleck:** das grüne Gate belegt das nicht (§5).
- [x] `make gates` Exit 0 (acht Gates, 458 Dateien, 0 Befunde, Coverage
      94,80 %); unabhängiger Review
      ([Report](../../../reviews/2026-08-23-slice-128-baseline-v5110-review.md)),
      Verdikt blockierend, 0 HIGH · 4 MEDIUM · 3 LOW · 1 INFO, alle acht
      eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **`BEO-008` steht bei Zähler 3 und ist genau hier einschlägig.** — **Ausgang:**
  **eingetreten**, und zwar an *derselben Datei:Zeile* wie in slice-106 und
  slice-110: `harness/README.md:60` behielt die
  `releases/download/`-URL auf `v5.9.0`, während das Link-Ziel derselben Zeile
  gehoben wurde. Die Ursache ist diesmal benennbar und nicht bloß „übersehen":
  der **Zensus** fand die Datei, das **Hebe-Skript** führte für Klasse 2 eine
  Datei-**Liste** statt eines Musters — und `harness/README.md` stand nicht
  darauf. Ein Zensus, der eine Stelle findet, und ein Werkzeug, das sie nicht
  anfasst, sind zwei verschiedene Dinge; nur das erste hatte ich geprüft.
- **Die Über-Hebung ist die zweite Richtung derselben Klasse.** — **Ausgang:**
  *nicht eingetreten.* Alle 34 Pfade und 5 URLs sind vom Reviewer **einzeln**
  auf Zeitform geprüft; die stehen gelassenen Vergangenheits-Aussagen sind
  korrekt abgegrenzt. Die Trennung hat also gehalten — die Lücke lag auf der
  anderen Seite.
- **Zwei Minors auf einmal.** — **Ausgang:** *nicht eingetreten.* Im ganzen
  Repo existiert kein einziger `v5.10.0`-Verweis; die Annahme ist geprüft,
  nicht geglaubt.

## 6. Trigger

**Start** (`open` → `in-progress`): sofort — die Welle ist eröffnet,
`in-progress/` frei.

**Rückführungen:** `in-progress` → `next`, falls die Integritätsprüfung des
Bundles fehlschlägt — dann ist es kein Vendoring-, sondern ein Upstream-Problem.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Vendoring (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **`BEO-008`**
  ist die zentrale — Zähler 3, Schwelle erreicht, und ihre benannte
  mechanische Form ist seit slice-122 **baubar** (`versions.patterns`). Ob
  dieser Slice sie baut, ist ein eigener Entscheid; ihn hier zu treffen wäre
  bequem und falsch. **`BEO-002`** für die Ränder der Pin-Hebung.
  **`BEO-011`** für jede Vollständigkeits-Aussage über die gehobenen Verweise.

Slice-ID: slice-128. Betroffene IDs:
[`MR-011`](../../../../harness/conventions.md#mr-011),
[`MR-021`](../../../../harness/conventions.md#mr-021). Module:
Harness-Vendoring, Konventionsspeicher. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Vendoring nach etablierter, viermal
gefahrener Prozedur.

## 9. Closure-Notiz (nach `done/`)

Etappe A ist geliefert: Bundle `v5.11.0` vendored und offline verifiziert, Pin
auf den neuen Tag gehoben, [`MR-030`](../../../../harness/conventions.md#mr-030) als sechster Nachtrag der Serie, Verweise
über drei Klassen gezogen, Alt-Baum entfernt. Der Kanon aus Kurs-Welle 94 —
Antwort auf einen Konsumenten-CR dieses Repos — liegt damit **gepinnt** vor.

**`BEO-008` ist zum vierten Mal eingetreten, und zum dritten Mal an derselben
Zeile.** Das ist die Lehre dieses Slice, und sie ist präziser als „nochmal
passiert": Ich habe den Zensus **zweimal** gefahren und beim zweiten Mal sogar
die Lücke der ersten Fassung gefunden (relative Pfade). Trotzdem blieb eine
URL stehen — weil zwischen Zensus und Hebung ein **Werkzeug** liegt, das ich
nicht gegen den Zensus gehalten habe: das Skript hob Klasse 2 über eine
handgeschriebene Datei-Liste, der Zensus lief über ein Muster. **Der Ableiter
ist nicht „gründlicher zählen", sondern: die Hebung gegen den Zensus
gegenprüfen, nicht gegen die Erwartung.** Ein `grep` nach der alten Form
**nach** der Hebung hätte sie in einer Sekunde gefunden; ich habe ihn nur für
Klasse 1 gefahren.

**Die zweite Lehre ist eine neue Ausprägung der Wellen-Klasse.** Ich habe zwei
Stellen mit der Begründung *„das ist nicht geprüft"* offen gelassen — und die
Prüfung lag im **selben Commit**: `modul-06-roadmap.md` liegt in der
23er-Gruppe der reinen Versions-Stempel, der zitierte Abschnitt ist
byte-identisch. Das ist nicht „Aussage aus dem Anlass statt aus dem Bestand",
sondern etwas Schärferes: **eine Aussage gegen die eigene Messung.** Die Daten
waren da, die Behauptung war das Gegenteil. Beide Stellen sind jetzt gehoben.

**Der dritte Punkt ist ein blinder Fleck, kein Fehler — und er wiegt schwerer
als die ersten beiden.** Ich habe „`457 Dateien, 0 Befunde` nach dem Entfernen
des Alt-Baums" als Beleg dafür genommen, dass kein Verweis ins Leere läuft, und
das mit Verzeichnis-vs-Datei-Pfaden begründet. Am Modul gemessen ist die
Begründung falsch: `codepaths` löst Verzeichnisse sehr wohl auf; der wahre
Grund ist `classifyCodepath`
([`internal/hexagon/core/rules/codepaths.go`](../../../../internal/hexagon/core/rules/codepaths.go)),
das `.harness/…`-Pfade **gar nicht als Pfad erkennt**. Die Zahl belegt die
Abwesenheit toter Inline-Code-Verweise also **nicht**. Die Schlussfolgerung
hält (es gibt nichts auszunehmen, weil nichts geprüft wird), die Begründung
nicht — und das ist genau der Unterschied zwischen einem Beleg und einem
grünen Gate. Was hält: die **Link**-Hälfte feuert, und `anchors` löst in den
`scan.ignore`-Baum hinein auf.

**Deklariert statt verschwiegen:** der Move-Commit von [`MR-029`](../../../../harness/conventions.md#mr-029) nach `done/`
bündelt mehr als [`AGENTS.md`](../../../../AGENTS.md) §3.3 für den
MR-Lifecycle-Move erlaubt („alles Übrige bleibt Commit 2"). Folgenlos — die
Rename-Detection greift bei 75 % —, aber es ist eine undeklarierte Präzedenz,
und undeklarierte Präzedenzen werden zitiert.

**Ein Folge-Punkt, hier nicht entschieden:** `BEO-008` steht jetzt bei vier
Vorfällen, und seine benannte mechanische Form ist seit
[slice-122](slice-122-versions-musterliste.md) **baubar** — `versions.patterns`
trägt ein zweites Muster-Quellen-Paar, das den Baseline-Tag in URLs und Prosa
gegen den §Baseline-Pin hielte. Drei der vier Vorfälle wären davon gefangen
worden. Das zu bauen ist ein eigener Entscheid mit eigener Messung; ihn hier
zu treffen wäre bequem und falsch — aber ihn *nicht zu benennen* wäre nach vier
Vorfällen unehrlich.
