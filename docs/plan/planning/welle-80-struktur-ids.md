# Welle welle-80-struktur-ids: Struktur-IDs nach Baseline — die Umkehr von MR-027

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-80-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Konventions-Pflege: Rückkehr auf
den Baseline-Default der Spec-Straten-Adressierung; kein Produkt-Release).

**Verantwortlich:** pt9912. **Datum:** 2026-08-22.

---

## 1. Welle-Ziel

Die Spec-Straten tragen **Struktur-IDs nach Baseline-Default**: `SPEC-<NNN>`
für jede technische Festlegung ohne eigene Anforderung in
`spec/spezifikation.md` (§2-Schemas als Überschrift-Kennung; §3 Defaults,
§4 Grund-Codes, §6 externe Verträge als Kennung **je Tabellenzeile** neben dem
fachlichen Schlüssel) und `ARC-<NNN>` für jede Komponente/Rolle und jeden
externen Berührungspunkt in `spec/architecture.md` — **fortlaufend je Datei,
Lücken werden nicht nachbelegt**, der Link trägt den Abschnitt, der Text die
Kennung. Neue ADRs nennen im `Schärft:`-Feld die Kennung, wo das Zielelement
eine trägt; die `Accepted`-ADRs davor bleiben auf ihren `§`-Ankern
(immutabel). Die deklarierte Abweichung
[`MR-027`](../../../harness/conventions.md#mr-027) ist damit **aufgelöst durch
Baseline-Konformität** — Auftraggeber-Entscheid 2026-08-22 (D1 Vollvergabe
im Bestand · D2 eine Kennung je §4-Zeile · D3 Konsument `ids` + `structure` +
`diagrams` opt-in · D4 Accepted-Bestand bleibt, Proposed zieht nach); der
mechanische Trigger „kein eindeutiger §-Anker" ist **nicht** abgewartet
worden, der Kandidat „§Kern"-Zeile in
[ADR-0012](../adr/0012-kern-paketschnitt-model-rules-app.md) wird bei der
Vergabe als Beleg gemessen, nicht behauptet. Das **Mehr** gegenüber den
Slice-DoDs: ein **Gate-Konsument**, der die Vergabe repo-weit lebendig hält —
ohne ihn wäre die Welle genau die „Pflege ohne Gegenwert", die [`MR-027`](../../../harness/conventions.md#mr-027) benannt
hat. Weil der `pre-commit`-Hook jeden Commit an einen grünen `doc-check`
bindet, wird das Gate **im Vergabe-Commit** scharf (vorher am Bestand gemessen
rot — das ist der Beleg), nicht davor.

## 2. Trigger (Welle startet)

Freigabe des Wellen-Plans durch den Auftraggeber mit den vier Entscheiden
D1–D4 (2026-08-22) — eingetreten; slice-112 in `done/` (eingetreten);
Inventur-Zahlen werden **je Slice frisch gemessen**, nicht aus der Vorschau
übernommen (Stand 2026-08-22: 46 ADRs mit `Schärft:`-Feld, 55 `Accepted`,
2 `Proposed`; Spezifikation 44 `###`-Sektionen, davon 5 kennungslos in §2;
Architektur ohne jede Kennung).

## 3. Closure-Trigger (Welle schließt)

- Alle vier Slices in `done/`; `make fullbuild` grün (Exit explizit).
- `make doc-check` grün **mit** den `ids`-Mustern `SPEC-\d{3}`/`ARC-\d{3}`,
  der `structure`-Regel „jede §2-Überschrift der Spezifikation trägt eine
  Kennung" und `diagrams` opt-in auf `spec/architecture.md` — am eigenen
  Bestand null Befunde.
- `matrix`-Messung: keine `ARC-*`-Nennung in `spec/spezifikation.md`
  (Abwärtsverweis Spezifikation → Architektur) — gemessen, nicht gesetzt.
- [`MR-027`](../../../harness/conventions.md#mr-027) per `git mv` in
  `conventions/done/`, Index-Zeile „aufgelöst durch Baseline-Konformität";
  [`MR-000`](../../../harness/conventions.md#mr-000--baseline-aussage)-Aussage
  nachgezogen; ADR-Index-Konvention trägt die Zwei-Formen-Regel.
- Ergebnisnotiz `welle-80-results.md` mit Register-Lese-Schritt.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [slice-113](open/slice-113-struktur-id-konvention.md) | Konvention zuerst: [`MR-027`](../../../harness/conventions.md#mr-027) aufgelöst (Baseline-Konformität), [`MR-000`](../../../harness/conventions.md#mr-000--baseline-aussage) nachgezogen, `ids.patterns` `SPEC-\d{3}`/`ARC-\d{3}` mit Wortgrenzen (grün am Bestand — noch keine Kennung existiert), ADR-README-Konvention (Kennung wo vorhanden, sonst §-Anker; Accepted-Bestand bleibt), AGENTS §5-Satz | [`DC-FA-ID-001`](../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids), [`MR-000`](../../../harness/conventions.md#mr-000--baseline-aussage)/[`MR-027`](../../../harness/conventions.md#mr-027) |
| [slice-114](open/slice-114-spec-vergabe-spezifikation.md) | `SPEC-*`-Vergabe in `spec/spezifikation.md` (§2 Überschriften, §3/§4/§6 je Zeile, fortlaufend) **und** die `structure`-Regel für §2 im selben Commit (vorher am Bestand gemessen rot); §7-Historie-Zeile; der [ADR-0012](../adr/0012-kern-paketschnitt-model-rules-app.md)-„§Kern"-Kandidat gemessen | [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in), Spezifikation §2–§6 |
| [slice-115](open/slice-115-arc-vergabe-architektur.md) | `ARC-*`-Vergabe in `spec/architecture.md` (§1 Kästen mit Komponenten-Tabelle, §2 Rollen nennen dieselben, §3 extern und §5 Fehlermodelle setzen fort) + `diagrams` opt-in auf die Architektur-Datei + `matrix`-Messung | Architektur §1–§5, Modul `diagrams` |
| [slice-116](open/slice-116-adr-neuzugangs-regel.md) | ADR-Neuzugangs-Regel + Erstanwendung: die beiden `Proposed`-ADRs ziehen `Schärft:`/`Bezug:` auf Kennungen, Slice-Kopf-Feld „Berührte Spec-Stellen" nutzt Kennungen, Reviewer-Anker | ADR-Index, `.harness/skills/reviewer.md` |

## 5. Abhängigkeiten

- Blockiert: nichts Geplantes.
- Wird blockiert von: nichts — konventions-intern; keine Produkt-Änderung,
  kein Release, keine Konsumenten-Kopplung (die RTM zählt `ARC-*` nicht, das
  `trace`-Anforderungsmuster zieht `SPEC-*` nicht mit). Reihenfolge innerhalb:
  113 vor 114/115 (Konvention vor Vergabe), 116 nach 114/115 (Erstanwendung
  braucht Kennungen).

## 6. Out-of-Scope für diese Welle

- **Kein Retarget bestehender `Accepted`-ADRs** — ihre `Schärft:`-Felder sind
  immutabel; die Zwei-Formen-Welt (alt §-Anker, neu Kennung) ist erklärte
  Konvention, kein Defekt.
- **Keine Produkt-Änderung:** `--suggest-config`/`--print-config` schlagen
  keine `SPEC-`/`ARC-`-Muster vor (wäre ein Change Request); ein neuer
  `structure`-Schlüssel ebenso; kein Tag.
- **Keine Struktur-IDs in Commit-Botschaften** (`commits.id-patterns`
  unverändert — Struktur-IDs gehören nicht in die Traceability).
- **Keine Lastenheft-Änderung:** §3 deklariert das projektspezifische
  Anforderungs-Präfix; Struktur-IDs sind Baseline-feste Form und werden im
  Konventionsspeicher deklariert. Zeigt die Umsetzung, dass §3 sie
  ausschließt, ist das ein eigener CR-Commit außerhalb der Welle.
- **Kein Bereichssegment, kein `DC-SPEC-*`** — die Straten-Tabelle schreibt
  `SPEC-<NNN>`/`ARC-<NNN>` fest.
