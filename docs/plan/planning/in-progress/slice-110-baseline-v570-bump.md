# Slice slice-110: Baseline-Bump v5.6.0 → v5.7.0 (Vendoring, Pin, Verweis-Hebung, Zwei-Hälften-Prosa)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** welle-79-zwei-haelften-ein-waechter (zugeordnet bei der Eröffnung).

**Bezug:** [`MR-021`](../../../../harness/conventions.md#mr-021)
(vendored-Verweise pin-gebunden),
[`MR-023`](../../../../harness/conventions.md#mr-023) (Bundle-Layout,
Materialisierungs-Skript),
[`MR-026`](../../../../harness/conventions.md#mr-026) (der zu ersetzende Pin
v5.6.0).

**Autor:** pt9912. **Datum:** 2026-08-21.

---

## 1. Ziel

Die Baseline von `v5.6.0` auf **`v5.7.0`** heben (Kurs-Welle 81 „Zwei Hälften,
ein Wächter", Tag 2026-08-21). Das Delta ist klein und vollständig gelesen:
drei Regelwerks-Dateien, +5/−3 Zeilen — der README-Stand, die
**modul-06-Neufassung von §Offene Wellen** (der Ruhe-Marker steht **zusätzlich**
zur Liste, nicht an ihrer Stelle; gewächtert ist nur die Marker-Hälfte, in
beide Richtungen; die Liste ist Ableitung ohne Wächter — eine Bijektion
„braucht ein eigenes Prädikat") und das **modul-10-Output-Schema mit `klasse`
als sechstem Feld**. Beide inhaltlichen Änderungen sind Landungen eigener
Upstream-Notizen (welle-78: 5-vs-6-Feld-Drift; Kurs-Session 2026-08-21:
Offene-Wellen-Widerspruch). Der Slice liefert Vendoring, Pin-Nachtrag,
Verweis-Hebung, den Prosa-Nachzug der Roadmap auf die Zwei-Hälften-Lesart und
den Delta-Audit über die zwei geänderten Regeln.

## 2. Vorgehen

1. **Vendoring:** `tools/harness/fetch-baseline-cache.sh v5.7.0` (explizites
   Tag-Argument — der Lauf braucht den neuen Pin nicht), danach `--verify`
   offline und `--check-latest` beidseitig (Currency **und** Content). Den
   `v5.6.0`-Baum entfernen — ein Pin, eine Wahrheit (Präzedenz slice-106).
   Arbeitsregeln aus dem Register: vor jedem pfad-selektiven Commit
   `git status` auf vorgestagten Index prüfen (BEO-006), Gate-Läufe mit
   explizitem Exit statt Pipe (BEO-007).
2. **Pin + Nachtrag:** `harness/conventions.md` §Baseline auf `v5.7.0`; der
   neue MR-Eintrag ersetzt [`MR-026`](../../../../harness/conventions.md#mr-026)
   in der Nachtrags-Kette ([`MR-011`](../../../../harness/conventions.md#mr-011)
   → [`MR-023`](../../../../harness/conventions.md#mr-023) →
   [`MR-026`](../../../../harness/conventions.md#mr-026) → neu); die Ersetzte
   wandert nach `conventions/done/` samt Zeile im Aufgelöste-Register und
   Link-Tiefen-Fix; Index-Zeile mit Kennungs-Anker für den Neuzugang.
3. **Verweis-Hebung:** alle **lebenden** pin-gebundenen Verweise auf
   `baseline/v5.6.0` retargeten. Zensus beim Schnitt (2026-08-21):
   `AGENTS.md` (3 Links), `harness/README.md` (4), `harness/conventions.md`
   (2), neun MR-Dateien (je 1),
   [`MR-027`](../../../../harness/conventions.md#mr-027) (2),
   `.harness/skills/reviewer.md`
   (2), `docs/plan/planning/README.md` (1), Roadmap (1) — beim Vollzug
   **frisch messen** (seit dem Zensus entstehen Fundstellen, u. a. dieser
   Slice selbst). **Eingefrorene** Fundstellen (`done/slice-109`,
   [`MR-024`](../../../../harness/conventions.md#mr-024) in
   `conventions/done/`) quell-skopiert tombstonen.
4. **Zwei-Hälften-Prosa-Nachzug:** die Sektionsregel der Roadmap unter
   §Offene Wellen trägt noch die v5.6.0-Lesart („… trägt der Abschnitt
   *stattdessen* den deklarierten Ruhe-Marker") — Neufassung nach v5.7.0:
   die Liste folgt den Dateien (je offener Welle ein Zeiger), der Marker
   folgt dem Anspruch und steht **zusätzlich**, beides zugleich ist der
   Normalfall nach der Eröffnung. Paraphrase-Disziplin beachten (der
   Marker-Wortlaut darf in der Sektions-Prosa nicht wörtlich stehen —
   Substring-Match des Wächters). **Nicht hier:** die Beschreibungen des
   **ausgelieferten Gate-Verhaltens** (AGENTS §3.3-Atomicity-Hinweis, die
   Grenz-Kommentare der Prüf-Profile) — die bleiben wahr, bis slice-111 das
   Produkt umbaut, und wandern dort.
5. **Delta-Audit** (je geänderter Regel eine Antwort, wie der Stufen-Audit —
   nur zwei Regeln statt sechs Stufen): **modul-06 §Offene Wellen** —
   Roadmap-Prosa nachgezogen (dieser Slice, Schritt 4); der
   **Produkt-Widerspruch** ist benannt und geht an slice-111 (die
   W3-Kopplung Aktiv-Status ⟺ Datei-Zahl gegen das Zwei-Hälften-Modell).
   **modul-10 §Output-Schema** — konform: `reviewer.md` 1.5.0 führt
   `klasse` als sechstes Feld seit slice-090; Fundstelle belegen, keine
   Handlung.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Produkt-Verhalten:** der `wave-drift`-Umbau ist slice-111 (eigene
  ADR, eigener Release-Punkt). Dieser Slice ändert Harness-Doku und
  Vendoring, keinen ausgelieferten Code.
- **Kein Retrofit** eingefrorener Artefakte (immutable ADRs, `done/`-Slices,
  Review-Reports) — Tombstones statt Umschreiben.
- **Kein Sechs-Stufen-Audit:** v5.6.0 → v5.7.0 ist eine Stufe mit zwei
  gelesenen Regel-Änderungen; der Audit-Umfang folgt dem Delta.

## 4. Definition of Done

- [ ] Pin `v5.7.0` vendored (beide Bäume + `SHA256SUMS`), `--verify` offline
      grün, `--check-latest` ohne Currency-/Content-Drift-Befund.
- [ ] `v5.6.0`-Baum entfernt; kein **lebender** Verweis nennt
      `baseline/v5.6.0`, eingefrorene sind quell-skopiert getombstoned
      (`make doc-check` belegt beides).
- [ ] Pin-Nachtrag Accepted (Index-Zeile + Kennungs-Anker);
      [`MR-026`](../../../../harness/conventions.md#mr-026) in
      `conventions/done/` samt Aufgelöste-Zeile, Links tiefen-korrigiert.
- [ ] Die Roadmap-Sektionsregel §Offene Wellen trägt die Zwei-Hälften-Lesart
      (Marker zusätzlich, paraphrasiert).
- [ ] Delta-Audit: je geänderter Regel eine Antwort mit Fundstelle.
- [ ] `make gates` grün (Exit explizit geprüft), GUARD vor jedem Commit.
- [ ] Unabhängiger Review vor der Closure.

## 5. Abnahme-Punkte / Risiken

- **Bundle-Inhalt kann wachsen:** die Kurs-Welle 81 enthält auch
  team-sim-Commits; ob `lab-regelwerk.zip` dadurch mehr Dateien trägt,
  entscheidet der `--verify`-Lauf — die Dateizahl wird **gemessen** notiert,
  nicht aus v5.6.0 (51) fortgeschrieben.
- **Marker-Substring-Falle:** der Prosa-Nachzug (Schritt 4) editiert genau
  den Block, den der Wächter liest — nach jedem Edit Gate-Lauf mit
  explizitem Exit.
- **Kettenzustand der MR-Nachträge:**
  [`MR-023`](../../../../harness/conventions.md#mr-023) bleibt
  (Layout-Aussage lebt),
  [`MR-026`](../../../../harness/conventions.md#mr-026) wandert (reine
  Pin-Fortschreibung, ersetzt) — dieselbe Auflösung wie
  [`MR-024`](../../../../harness/conventions.md#mr-024) in slice-108, nicht
  dieselbe wie [`MR-023`](../../../../harness/conventions.md#mr-023) in
  slice-106.

## 6. Trigger

**Start** (`open` → `in-progress`): Eröffnung der Welle
welle-79-zwei-haelften-ein-waechter (Kurs-Tag v5.7.0 samt Release-Asset
verifiziert 2026-08-21; Anstoß des Auftraggebers am selben Tag).

**Rückführungen:** `in-progress` → `next`, falls das Bundle-Layout von
v5.7.0 vom [`MR-023`](../../../../harness/conventions.md#mr-023)-Layout
abweicht und das Skript einen eigenen Vorab-Fix braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

Slice-ID: slice-110. Betroffene IDs: keine `DC-*` (Harness, kein Produkt);
[`MR-021`](../../../../harness/conventions.md#mr-021)/[`MR-023`](../../../../harness/conventions.md#mr-023)/[`MR-026`](../../../../harness/conventions.md#mr-026)
und der neue Pin-Nachtrag. ADRs: keine. Module:
`.harness/baseline/`, `harness/conventions*`, `AGENTS.md`,
`harness/README.md`, `.harness/skills/reviewer.md`, Roadmap,
`docs/plan/planning/README.md`. Gates: `make gates` (mit `doc-check`),
GUARD-Image-Lauf vor jedem Commit.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — adoptierte Konvention,
konventionsgetragene Doku-Hebung; kein Legacy-Import.

## 9. Delta-Audit v5.6.0 → v5.7.0 (je geänderter Datei eine Antwort)

Das Delta der Stufe umfasst drei Regelwerks-Dateien (`git diff
v5.6.0..v5.7.0 -- lab/regelwerk/` im Kurs-Repo, nachvollzogen am
vendorten Baum): der Audit folgt dem Delta, kein Stufen-Ritual.

| Datei / geänderte Regel | Urteil | Beleg |
|---|---|---|
| `README.md` (Stand-Zeile „Kurs-Welle 81 · 2026-08-21") | n. a. | reine Stand-Fortschreibung, keine Regel |
| `modul-06-roadmap.md` §Offene Wellen (Zwei-Hälften-Fassung: Marker **zusätzlich** zur Liste, nur die Marker-Hälfte gewächtert) | **angepasst** (Doku) + Produkt-Folge benannt | die Roadmap-Sektionsregel trägt die v5.7.0-Lesart (dieser Slice, Schritt 4; Marker paraphrasiert); der Produkt-Widerspruch — `wave-drift` hält den Aktiv-Status gegen die Datei-Zahl — geht als Umsetzung des Konsumenten-CR an slice-111 (`planning.waves.mode`) |
| `modul-10-review-harness.md` §Output-Schema (`klasse` als sechstes Feld) | **konform** | [`reviewer.md`](../../../../.harness/skills/reviewer.md) 1.5.0 führt `kategorie · quelle · pfad · befund · verifizierbar · klasse` seit slice-090 — keine Handlung; die Baseline hat damit die eigene Upstream-Notiz (5-vs-6-Feld-Drift, welle-78-Wiedervorlage) geschlossen |

## 10. Closure-Notiz (nach `done/`)

*Wird bei der Closure geschrieben (Struktur nach `closure.heading-pattern`).*
