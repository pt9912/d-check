# Slice slice-106: Baseline-Bump v5.0.0 → v5.6.0 — Etappe A (Vendoring, Pin, Verweis-Hebung)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-78-baseline-v560-migration (zugeordnet bei der Eröffnung).

**Bezug:** [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
(Vorgänger-Hebung; Präzedenz für Layout und Tombstone-Muster),
[`MR-021`](../../../../harness/conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden)
(vendored-Verweise pin-gebunden),
[`MR-022`](../../../../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)
(Currency-Audit). Neue Adaption **MR-026** (Pin-Hebung auf <!-- d-check:ignore -->
v5.6.0) entsteht mit diesem Slice. Kein `DC-*`-Bezug — reine Harness-/Adoptions-Arbeit, kein
Produkt-Code.

**Autor:** pt9912. **Datum:** 2026-08-21.

---

## 1. Ziel

Die adoptierte Baseline von `v5.0.0` auf **`v5.6.0`** heben (Kurs-Tag
2026-08-16; sechs additive Stufen, 20 Regelwerks-Dateien, +902/−152 Zeilen —
nichts entfernt): das self-contained Bundle `lab-regelwerk.zip@v5.6.0` committet
vendored unter `.harness/baseline/v5.6.0/` (beide Bäume
`{regelwerk,templates}` + `SHA256SUMS`), der Pin in `harness/conventions.md`
§Baseline auf v5.6.0, alle **lebenden** Verweise auf den neuen Pfad, die
**eingefrorenen** getombstoned.

## 2. Vorgehen

1. **Vendorn:** `tools/harness/fetch-baseline-cache.sh v5.6.0` (explizites
   Tag-Argument — der Pin wird erst im selben Bogen umgestellt), danach
   `--verify` offline gegen das frische `SHA256SUMS`.
2. **Pin heben:** `harness/conventions.md` §Baseline auf v5.6.0; neue Adaption
   **MR-026** als Datei `harness/conventions/MR-026-baseline-v560.md` <!-- d-check:ignore -->
   (Nachtrag zu [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout),
   dieselbe Bauform) + Index-Zeile.
3. **Alten Baum entfernen** (`.harness/baseline/v5.0.0/`) — Präzedenz [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout):
   **ein** Pin, **eine** netzlose Lese-Form; zwei parallele Bäume wären eine
   zweite Quelle, die altert. Abnahme-Punkt 1 unten.
4. **Verweise heben** (Bestandsmessung beim Slice-Start: 29 Dateien außerhalb
   des vendorten Baums nennen `baseline/v5.0.0`):
   - **lebend** → retargeten auf `v5.6.0/` (u. a. `AGENTS.md` §1,
     `harness/README.md` §Guides, `harness/conventions.md`, aktive
     `MR-*`-Dateien, `docs/plan/planning/README.md`, die
     Regeln-dieser-Sektion-Links der Roadmap);
   - **eingefroren** → `ignore-refs`-Tombstones in `.d-check.yml` erweitern
     (dieselbe Klasse wie die drei v1.4.0-Tombstones aus [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)): immutable
     `Accepted`-ADRs (mind. 0047/0048), `done/`-Slices, Review-Reports —
     **nicht** editieren, quell-skopiert ausnehmen.
5. **Gates:** `make doc-check` als engster Sensor (Links/Anker über den
   neuen Baum), dann `make gates`; `--check-latest` als Currency-Gegenprobe
   (erwartet: Pin aktuell, Content am Tag unverändert).

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Konformitäts-Abgleich** der eigenen Dokumente gegen die neuen Regeln —
  das ist Etappe B ([slice-107](../open/slice-107-baseline-v560-delta-audit.md), Lesen)
  und Etappe C (Umsetzen, Slices nach B-Befund).
- **Keine inhaltliche Übernahme** neuer Regelwerks-Konzepte (Team-Fähigkeit,
  Reconciliation-Register, …) — erst lesen, dann entscheiden.

## 4. Definition of Done

- [ ] `.harness/baseline/v5.6.0/{regelwerk,templates}` + `SHA256SUMS`
      committet; `--verify` offline grün; der v5.0.0-Baum ist entfernt
      (bzw. der Abnahme-Punkt anders entschieden und begründet).
- [ ] Pin in §Baseline auf v5.6.0; MR-026 als Datei + Index-Zeile <!-- d-check:ignore -->
      (Bauform [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)).
- [ ] Kein lebender Verweis nennt mehr `baseline/v5.0.0`; die eingefrorenen
      sind quell-skopiert getombstoned — `make doc-check` grün belegt beides.
- [ ] `make gates` grün; unabhängiger Review; **kein Release** (der Harness
      ist nicht das Produkt — Präzedenz welle-67: Releases trugen nur die
      Produkt-Slices).

## 5. Abnahme-Punkte / Risiken

1. **v5.0.0-Baum entfernen vs. behalten.** Vorschlag: entfernen (Präzedenz
   [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout), „was hier steht, liest jeder Agentenlauf"); der Preis sind mehr
   Tombstones als bei v1.4.0 (die eingefrorene Verweis-Menge ist seither
   gewachsen). Behalten hieße: zwei Regelwerks-Stände netzlos nebeneinander,
   und jeder künftige Leser muss wissen, welcher gilt.
2. **Immutable ADRs mit v5.0.0-Pfaden** (0047/0048 u. a.) dürfen nicht
   editiert werden — der Tombstone ist die einzige gate-ehrliche Antwort;
   die Tombstone-Liste wächst und gehört im Config-Kommentar begründet.
3. **Der Bundle-Inhalt könnte vom Kurs-Arbeitsstand abweichen** — vendored
   wird das **Release-Asset** am Tag, nicht der Arbeitsbaum; `--check-latest`
   Teil B ist die Authentizitäts-Gegenprobe.

## 6. Trigger

**Start** (`open` → `in-progress`): Freigabe des Auftraggebers (2026-08-21,
Reihenfolge-Entscheid „erst BEO-005, dann Migration") **und** WIP-Slot frei —
beide mit der Eröffnung von welle-78 erfüllt.

**Rückführungen:** `in-progress` → `next`, falls das Bundle-Layout von v5.6.0
vom v5.0.0-Layout abweicht und das Materialisierungs-Skript einen eigenen
Anpassungs-Slice braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Mechanik (`tools/harness/`, GF via
  [`MR-004`](../../../../harness/conventions.md#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild))
  und Harness-Doku (Repo-Default GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-21: keine
  unverkörperte Beobachtung offen; BEO-002/003/004 verkörpert, BEO-001/005
  gestrichen): **BEO-002**/[`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten) wirkt als Regel — die Pin-Hebung ändert
  eine deklarierte Menge (den Baseline-Stand) mit Spiegeln in `AGENTS.md`,
  `harness/README.md`, beiden Konventions-Sektionen und dem
  Materialisierungs-Skript-Kommentar; die Liste steht in §2. Die in
  slice-090 notierte **upstream**-Beobachtung (Baseline-interne
  5-vs-6-Finding-Feld-Drift) ist beim v5.6.0-Drift-Audit in Etappe B
  wiedervorzulegen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — adoptierte Konvention, konventionsgetragene
Hebung nach dokumentierter Präzedenz.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
