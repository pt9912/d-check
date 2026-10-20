# slice-277: Baseline-Pin auf `v6.18.1` heben — mechanisch, ohne Urteil über das Delta

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — reaktive Arbeit, die Closure-Bedingung ist die DoD
selbst.

**Bezug:** [`MR-074`](../../../../harness/conventions.md#mr-074) (die vorige
Pin-Hebung und ihre Lehre), [`MR-011`](../../../../harness/conventions.md#mr-011)
(die Pin-Kette), [`MR-021`](../../../../harness/conventions.md#mr-021)
(pin-gebundene Verweise), [`MR-039`](../../../../harness/conventions.md#mr-039)
(Zitat-Delta), [`MR-051`](../../../../harness/conventions.md#mr-051)
(`d-check:cite`-Spannen), [`MR-055`](../../../../harness/conventions.md#mr-055)
(Symlinks als Träger), [`MR-069`](../../../../harness/conventions.md#mr-069)
(`ignore-refs`), [`MR-070`](../../../../harness/conventions.md#mr-070)
(Frozen-Klassen vor mechanischer Ersetzung). Anlass: der Nachtlauf
`upstream-drift` meldet den neuen Release `v6.18.1` (Pin `v6.17.0`).

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Den committet vendorten Baseline-Bestand auf
[`v6.18.1`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.18.1)
heben (den aktuellen Tag; `v6.18.0` liegt dazwischen), den Pin in
`harness/conventions.md` §Baseline fortschreiben, alle lebenden pin-gebundenen
Referenzen retargeten und die Hebung als neuen MR-Eintrag führen — Frozen-Klassen
und Zitat-Delta gemessen am echten Vorzustand.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Adoption der inhaltlichen Deltas.** Jede ist eine Konventions- oder
  Produkt-Entscheidung mit eigener Begründung; der Bump ist mechanisch (anderer
  Vorgang). Der neue MR-Eintrag nennt das Delta vollständig nach Datei; Adoptieren
  gehört in eigene Folge-Slices nach Freigabe des Auftraggebers. Zu erwarten
  (aus einem Probelauf von `ai-harness-init` v0.10.0 gemessen, hier erst am echten
  Delta zu bestätigen): Vorlagen-Pfade nennen `.harness/baseline/<tag>/templates/…`,
  die Review-Report-Vorlage verlangt die volle Slice-Kennung im Dateinamen, die
  `.d-check.yml`-Vorlage trägt die `reviews`-Schlüssel `match: name`,
  `require-promises`, `skip-allows-empty`, `skip-pattern` und `recursive`.
- **Frozen-Dateien bleiben byte-stabil** — `done/`-Slices und -Wellen,
  `docs/reviews/`, `Accepted`-ADRs, eingehende CRs, `conventions/done/MR-*`,
  CHANGELOG- und Spec-Historie nennen den alten Pin als Vergangenheits-Aussage
  korrekt (Bestand bewusst).
- **`tools/harness/selbstpruefung.sh`** — sein Stand-Vermerk nennt die Version
  des **Schwester-Werkzeugs** `ai-harness-init`, nicht den Baseline-Pin
  (anderer Gegenstand).
- **Die Tool-Pins** (Go, golangci-lint, semgrep, a-check) und alle Dependabot-
  Hebungen — eigene Vorgänge mit eigener Messung.
- **Kein Release** — der Bump berührt keine Distribution-Fläche.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] **Re-Vendor:** Der vendorte Baum steht auf `v6.18.1` (`SHA256SUMS`
      generiert und verifiziert), der `v6.17.0`-Baum ist entfernt; §Baseline-Pin
      nachgezogen; der Vorgänger-Eintrag nach `conventions/done/`, neuer
      MR-Eintrag im Index.
- [ ] **Spiegel:** Alle lebenden pin-gebundenen Referenzen nennen `v6.18.1`
      (Pfad-Verweise, Release-/Tree-URLs, bare Nennungen in lebenden Dokumenten,
      Symlinks unter `.claude/rules/`); `d-check:cite`-Spannen neu geankert;
      Zitat-Delta und Frozen-Liste im MR-Eintrag gemessen.
- [ ] `make gates` grün; `make baseline-freshness` meldet den Pin aktuell.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8, kein
      Self-Review; Verifikation.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder
      „keine Beobachtung angefallen" in §7 notiert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/{v6.17.0 → v6.18.1}/` | neu + entfernen | Re-Vendor beider Bäume + `SHA256SUMS` über `tools/harness/fetch-baseline-cache.sh v6.18.1`; alter Baum weg |
| `harness/conventions.md` | update | §Baseline-Pin, §Adoptierte Konventions-Quellen, MR-Index |
| neue MR-Datei unter `harness/conventions/` | neu | Hebung-Eintrag (Delta nach Datei, Frozen-Liste, Zitat-Delta, Cite-Neu-Ankern) |
| die Datei des Vorgänger-Eintrags (`…-baseline-v6170.md`, [`MR-074`](../../../../harness/conventions.md#mr-074)) | move | reiner `git mv` nach `conventions/done/`, dann Link-Tiefen |
| lebende Dokumente, die `v6.17.0` nennen (`AGENTS.md`, `harness/README.md`, `harness/rules/*`, aktive `MR-*`, `.claude/agents/*`, `.claude/commands/*`, `.harness/skills/*`, `roadmap.md`, `planning/README.md`, `observations/README.md`, `spec/*`) | update | Tag-Swap je Datei, Mischfundstellen zeilenweise |
| `.claude/rules/*.md` (Baseline-Symlinks) | update | Aliase auf `v6.18.1` |
| `.d-check.yml` | update | nur, falls gemessen: nächste `ignore-refs`-Stufe |
| dieser Plan | update | die beiden `d-check:cite`-Spannen auf den neuen Baum |

**Plan-Änderung vor dem Code (2026-10-20, Schritt 15):** Der erste Re-Vendor-Lauf
bricht auf macOS an der Under-Copy-Barriere ab (`Quelle (54) != vendored (      54)`):
`tools/harness/fetch-baseline-cache.sh` vergleicht die zwei Zahlen als Strings
(`[ "$src_n" = "$dst_n" ]`), und BSD-`wc -l` füllt links mit Leerzeichen auf — dieselbe
Falle, die der Kommentar bei der Manifest-Deckung (`-eq statt =`) schon benennt.
Die Dateizahl stimmt (54 entpackt, 54 in der Quelle). Zusätzlich zur Tabelle:

| `tools/harness/fetch-baseline-cache.sh` | update | Zeile der Under-Copy-Barriere vergleicht numerisch (`-eq`); eigener Commit vor dem Re-Vendor, weil es Werkzeug ist und nicht Bestand |

Aussagen-Suche (Schritt 20): `grep -n 'wc -l' tools/harness/fetch-baseline-cache.sh` —
drei Zählstellen; Zeile 106 rechnet arithmetisch (`$(( … ))`), Zeile 218 vergleicht
mit `-eq`, nur Zeile 343 vergleicht als String. Kein Spiegel außerhalb des Skripts.

**Plan-Änderung (nachgetragen im Review-Vorlauf, Schritt 15 — hätte vor dem
Code-Commit stehen müssen):** Zusätzlich zur Tabelle trägt der Hauptcommit eine
Datei aus dem Frozen-Bestand: `slice-274` (am selben Tag geschlossen) enthält einen
Markdown-Link auf den `v6.17.0`-Baum (`modul-11` §Bewusstes Brechen). Nach dem
Entfernen des alten Baums meldete `doc-check` ihn als toten Link; die Alternative wäre
ein neuer `ignore-refs`-Eintrag ([`MR-069`](../../../../harness/conventions.md#mr-069)),
eine deklarierte Gate-Senkung. Der Link ist auf `v6.18.1` retargetet, der Text des
Slice unberührt; der MR-Eintrag führt es als einzige Ausnahme von der Frozen-Regel.
Dazu die Link-Tiefen von `MR-074` nach dem Move (eigener Commit) und die
Zitat-Spanne `modul-10` im Reviewer-Skill (87→99).

**Spiegel-Messung** ([`MR-070`](../../../../harness/conventions.md#mr-070),
[`MR-074`](../../../../harness/conventions.md#mr-074) „für die nächste Hebung"):
`git grep -l "v6\.17\.0"` über das **ganze** Repo, ausgenommen nur
`.harness/baseline/` (nicht `.harness/` — dort liegen Skills), in Frozen- und
lebende Klassen geteilt, **bevor** der erste Tag-Swap fällt; die Messung steht im
MR-Eintrag. Beim Anlegen: 52 Dateien außerhalb der Frozen-Verzeichnisse, des
Baums, des CHANGELOG und der CRs nennen den Tag.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice;
Priorisierung gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): Das Delta ändert die Struktur des Bundles
  (neue Bäume, entfallene Dateien, die verlinkt sind), sodass der Bump nicht mehr
  mechanisch ist. Wächst der Diff über eine Review-Sitzung, gilt die Form von
  [`MR-066`](../../../../harness/conventions.md#mr-066).
- `in-progress` → `open` (blockiert): Das Release-Bundle ist nicht ladbar oder
  sein `SHA256SUMS` stimmt nicht.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft; `make baseline-freshness` meldet den Pin aktuell.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Ein lebender Verweis zeigt auf einen Anker oder eine Datei, die es im neuen
  Baum nicht mehr gibt.** — **Ausgang:** offen bis Closure; Gegenprobe: 55 Dateien
  vorher wie nachher, jeder lebende Anker löst gegen `v6.18.1` auf.
- **Eine Mischfundstelle wird pauschal ersetzt und fälscht eine
  Vergangenheits-Aussage.** — **Ausgang:** offen bis Closure; zeilenweise entschieden.
- **Eine Spiegel-Klasse bleibt ungemessen** (Skills, Symlinks, Kommentare in
  Konfigurationen). — **Ausgang:** offen bis Closure; Messung über das ganze Repo
  ohne `.harness/baseline/`.

## 7. Closure-Notiz

<!-- Wird bei der Closure gefüllt, vor dem `git mv` nach `done/`. -->

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die Harness-Mechanik
(`tools/harness/`, `HARN`, GF) und die pin-gebundenen Verweise im ganzen Repo
(`*`, `ALL`); beide deklariert, beide GF.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; Treffer:

- [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  (8×) und
  [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (5×): genau diese Lage — die Spiegel-Messung schließt nur `.harness/baseline/`
  aus, die Frozen-Klassen werden vorher gelistet.
- [`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
  (18×): das Delta-Inventar im MR-Eintrag wird am echten Baum gezählt, nicht aus
  dem Probelauf übernommen.
- [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (20×): neu geankerte `d-check:cite`-Spannen werden gegen den Wortlaut des neuen
  Baums gelesen, nicht nur verschoben.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`image-scan` grün; `upstream-drift` rot, **planmäßig gelesen**: neuer
Baseline-Release `v6.18.1` — das ist der Gegenstand dieses Slice — und veraltete
Pins für Go, golangci-lint, semgrep und a-check (eigene Vorgänge, hier
ausgeschlossen).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
