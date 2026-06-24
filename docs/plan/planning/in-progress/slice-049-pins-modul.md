# Slice slice-049: Modul `pins` — Content-Pin gegen inhaltlichen Drift (Idee 2)

**Status:** in-progress (offen, welle-38-pins).

**Welle:** welle-38-pins (Trigger: Auftraggeber-Idee — „wenn man einen
Markdown-Link auf eine Datei/Absatz setzt, könnte d-check den Drift bemerken";
nach Spike bestätigt).

**Bezug:** Idee 2 dieser Session-Kette. Führt eine **neue PIN-Anforderung** im
Lastenheft ein (Modul `pins`, Content-Pin/`link-stale`) plus einen
**begleitenden ADR** (Fence-Öffnung für den Ziel-Span, Mechanik-Präzedenz
[ADR-0018](../../adr/0018-diagram-fence-ausnahme.md)). Verteilung wie der Rest
des Werkzeugs (gepinntes Image, kein kopiertes Skript —
[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Linie).

**Autor:** pt9912. **Datum:** 2026-06-24.

---

## 1. Ziel

Strukturellen Drift fängt d-check schon: Ziel-Datei/-Anker weg →
`target-missing`/`anchor-missing`. **Neu:** *inhaltlicher* Drift — der Link löst
weiter auf, aber der **Inhalt** des Ziel-Spans hat sich seit dem Verlinken
geändert; die zitierende Aussage daneben ist evtl. veraltet (das klassische
„stale citation"-Problem). Ein optionaler **Content-Pin** (Hash des Ziel-Spans
zum Verlink-Zeitpunkt) macht das entscheidbar: d-check re-hasht den aufgelösten
Span und vergleicht; Mismatch → Befund `link-stale`. Gleiche Disziplin wie der
Image-Digest-Pin (pinnen → bei Änderung neu segnen), auf Doku-Querverweise
angewandt.

## 2. Entscheidungen

- **Entscheidbar statt semantisch.** Hash-Vergleich ist deterministisch + offline
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
  „ist das Zitat *sinngemäß* noch richtig?" wäre nicht gatebar und bleibt
  Out-of-Scope. Dieselbe Linie wie [ADR-0018](../../adr/0018-diagram-fence-ausnahme.md)
  (Existenz/Referenz, nicht Grammatik).
- **Ziel-Span:** ganze Datei (Link ohne Anker) oder Heading-Section (Anker → bis
  zur nächsten gleich-/höherrangigen Überschrift; Section-Logik aus
  `matrix`/`anchors` wiederverwendbar). Absatz-Ebene ist in Markdown nicht stabil
  adressierbar → draußen.
- **Roher Span inkl. Fenced-Code hashen** (Drift in Code-Beispielen ist ein
  Zielfall) — bewusste, eng begrenzte Ausnahme von der Fence-Opazität; Präzedenz
  [ADR-0018](../../adr/0018-diagram-fence-ausnahme.md) → eigener ADR.
- **Normalisierung whitespace-/reflow-invariant** (Wort-Inhalt, nicht
  Byte-Layout) — Spike 2026-06-24: 0/87 kosmetische Trips, das Rauschen
  materialisiert sich nur ohne Normalisierung.
- **opt-in pro Link:** nur Links **mit** Pin werden geprüft (der Pin ist die
  bewusste „hier zitiere ich Inhalt"-Markierung); strikt opt-in Modul, default-off
  byte-identisch.
- **diagnose-only v1:** `link-stale` liefert keinen `--repair`-Hunk — Re-Pinnen
  ist menschliche Annahme der Drift, kein eindeutig ableitbarer Fix; ein
  explizites `--bless` ist eine spätere CR.
- **Pin-Ablage inline:** `<!-- dpin: sha256:<hash> -->` direkt nach dem Link
  (lokal, git-diff-freundlich) — dieselbe Mechanik wie der `d-check:ignore`-Marker.
- **Eigenes Modul `pins`** (nicht mit `versions` verschmolzen): andere Mechanik
  (Span-Hash vs. Wert-Gleichheit), je einzeln opt-in und testbar.

## 3. Definition of Done

- [x] **Spec:** neue PIN-Anforderung im Lastenheft (neues Bereichskürzel `PIN` in
  §3, Versions-Bump + §7-Historie) + begleitender ADR (Fence-Öffnung) +
  spezifikation `.a`-Sektion + Grund-Code `link-stale` (§4) + `pins` als gültiges
  Modul in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)/Glossar
  + ADR-Index; doc-first vor Code. Die `.a` legt **deterministisch** fest (R1):
  (a) **Marker-Bindung** — der Pin bindet an den unmittelbar (nur Whitespace)
  vorausgehenden Link derselben Zeile; nicht eindeutig zuordenbare Marker sind
  inert; (b) **nicht auflösbares Ziel** — `pins` wertet nur auflösbare Links, der
  strukturelle Befund bleibt bei `links`/`anchors` (kein eigener `pins`-Befund,
  auch im pins-only-Lauf); (c) `pins` respektiert den Modul-Scope
  ([`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)).
- [ ] **Code:** Modul `pins` (Marker-Erkennung mit deterministischer Bindung,
  Ziel-Span-Auflösung, whitespace-Normalisierung, Hash + Vergleich, `link-stale`).
  Tests: Happy/Reflow-Boundary/Negative/Modul-aus; **Marker-Ambiguität** (zwei
  Links/Zeile, Marker zwischen Links, Marker auf Folgezeile); **Ziel-weg**
  (pins-only → kein `link-stale`; mit `links` aktiv → `target-missing`, kein
  Doppelbefund); **`pins.scope`** — Befunde nur im effektiven Modul-Scope.
- [ ] `make gates` grün; unabhängiges Review; Closure (Move nach `done/` +
  Roadmap-Flip, [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Normalisierung ist ein Vertrag** (definiert, was „inhaltliche Änderung"
  heißt) — Spike-kalibriert, exakt zu dokumentieren; konservativ.
- **Pin-Erstellung/-Erneuerung** ist menschlich (read-only-Kernvertrag); ein
  `--bless`-Emissionsmodus berührt den Reparatur-Vertrag und ist spätere CR.
- **Fence-Öffnung breiter begründen** als bei `diagrams` (ganzer Span statt
  gelistete Diagramm-Fences) — gehört in den ADR; opt-in default-off hält die
  Abwärtskompatibilität.

## 5. Trigger

Auftraggeber: „Idee 2 als Slice angehen" (2026-06-24), nach Spike + Design-
Diskussion; Reihenfolge D nach Idee 1 (slice-048 done, Release v0.28.0).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

_folgt mit dem Abschluss des Slice._
