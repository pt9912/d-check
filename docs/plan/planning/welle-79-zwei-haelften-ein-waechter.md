# Welle welle-79-zwei-haelften-ein-waechter: Baseline v5.7.0 + das eigene Prädikat der Listen-Hälfte

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-79-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Adoptions-Pflege der Baseline
plus die eine Produkt-Anpassung, die sie fällig macht).

**Verantwortlich:** pt9912. **Datum:** 2026-08-21.

---

## 1. Welle-Ziel

Die Baseline von `v5.6.0` auf **`v5.7.0`** heben (Kurs-Welle 81 „Zwei
Hälften, ein Wächter", Tag 2026-08-21; Bundle-Delta fünf Dateien — drei
Regelwerks-Dateien +5/−3, zwei Template-Spiegel derselben Regeln) und die
eine Produkt-Konsequenz ziehen. Beide Stufen-Änderungen sind
Landungen **eigener Upstream-Notizen**: modul-10 führt `klasse` als sechstes
Output-Feld (die 5-vs-6-Drift aus der welle-78-Wiedervorlage), und modul-06
fasst §Offene Wellen als **zwei unabhängige Aussagen** — die Liste folgt den
Dateien, der Ruhe-Marker folgt dem Anspruch und steht **zusätzlich** zur
Liste; gewächtert ist nur die Marker-Hälfte, die Listen-Bijektion „braucht
ein eigenes Prädikat". Damit ist der seit der slice-108-Closure benannte
Grenz-Zustand „Welle offen, `in-progress/` leer" Baseline-**Normalfall**,
und die W3-Kopplung des `planning`-Moduls (Aktiv-Status ⟺ Datei-Zahl)
widerspricht dem Modell, das sie stützen soll. Die Welle liefert den Bump
([slice-110](in-progress/slice-110-baseline-v570-bump.md)) und genau das
benannte eigene Prädikat — als **opt-in `planning.waves.mode: many`** nach
dem formalen CR des Konsumenten ai-harness-course (2026-08-21), Default
byte-identisch
([slice-111](open/slice-111-wave-drift-zwei-haelften.md)).

## 2. Trigger (Welle startet)

Anstoß des Auftraggebers („Jetzt gibt es v5.7.0", 2026-08-21); der Kurs-Tag
v5.7.0 samt `lab-regelwerk.zip`-Release-Asset ist verifiziert (2026-08-21);
der formale Konsumenten-CR liegt vor (mit team-sim-Messung s04a–s04d,
11/11 PASS); WIP-Slot frei (welle-78 geschlossen, Ruhe-Zustand, `open/`
war leer).

## 3. Closure-Trigger (Welle schließt)

- Beide Slices in `done/` —
  [slice-110](in-progress/slice-110-baseline-v570-bump.md) und
  [slice-111](open/slice-111-wave-drift-zwei-haelften.md).
- Pin `v5.7.0` vendored, `--verify` offline grün **und** `--check-latest`
  ohne Currency-/Content-Drift-Befund; kein **lebender** Verweis nennt mehr
  `baseline/v5.6.0`, die eingefrorenen sind quell-skopiert getombstoned
  (`make doc-check` belegt beides).
- Der Delta-Audit trägt **je geänderter Regel** eine Antwort (zwei Regeln —
  der Umfang folgt dem Delta, kein Stufen-Ritual).
- Unter `mode: many` misst `wave-drift` die **Kennungs-Bijektion** (beide
  Richtungen als Testfälle belegt rot, die drei unter `one` roten
  baseline-legitimen Zustände belegt grün); der Default ist byte-identisch
  belegt; **dieses Repo läuft nach dem Digest-Backfill selbst auf `many`**
  mit grünem Gate.
- **Release Minor** (v0.62.0) samt Digest-Backfill — der Release-Punkt ist
  hier **vor** dem Schnitt entschieden (konsumentensichtbar additiv),
  anders als bei welle-78, wo er ehrlich offen blieb.
- `make fullbuild` grün.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-110](in-progress/slice-110-baseline-v570-bump.md) | Bump: Vendoring v5.7.0, Pin-Nachtrag, Verweis-Hebung + Tombstones, Zwei-Hälften-Prosa der Roadmap, Delta-Audit |
| [slice-111](open/slice-111-wave-drift-zwei-haelften.md) | Produkt: `planning.waves.mode` (`one`\|`many`) — Kennungs-Bijektion als opt-in nach Konsumenten-CR; CR-Commit + [ADR-0055](../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)-Fortschreibung, Release v0.62.0 |

## 5. Abhängigkeiten

- Kurs-Repo-Tag `v5.7.0` mit Release-Asset `lab-regelwerk.zip` (liegt vor).
- **slice-110 vor slice-111 (bindend):** die ADR-Fortschreibung zitiert die
  **vendorte** v5.7.0-Formulierung, nicht das Kurs-Repo — dieselbe
  Netzlos-Disziplin wie beim welle-78-Stufen-Audit.
- Das Materialisierungs-Skript
  (`tools/harness/fetch-baseline-cache.sh`,
  [`MR-023`](../../../harness/conventions.md#mr-023)-Layout) nimmt ein
  explizites Tag-Argument — der erste Vendor-Lauf braucht den neuen Pin
  nicht.

## 6. Out-of-Scope für diese Welle

- **Keine Änderung an `planning-drift`**, kein zweiter Grund-Code, keine
  Default-Änderung, keine Festlegung der Block-Form (Tabelle vs. Liste) —
  CR §6 wörtlich.
- **Kein Mehr-Wellen-Betrieb dieses Repos** — die Bijektion macht ihn
  prüfbar, die Nutzung ist ein eigener Roadmap-Entscheid.
- **Kein Retrofit** eingefrorener Artefakte (immutable ADRs, `done/`-Slices,
  Review-Reports) — Tombstones statt Umschreiben; die
  [ADR-0055](../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)-Fortschreibung
  ist **kein** Retrofit (Proposed, `## Geschichte`).
