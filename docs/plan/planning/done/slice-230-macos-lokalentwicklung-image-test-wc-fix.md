# Slice slice-230: macOS-Lokalentwicklung — GNU-Annahmen in `fetch-baseline-cache.sh` (`wc`, `readlink`) + `image-test`-Rezept dokumentieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus — kein repo-weiter Beleg, den dieser Slice allein nicht liefert
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** — kein `DC-*`, keine aktive ADR. Reine Harness-Meta-Tooling-/
Doku-Korrektur ohne Berührung der `d-check`-Produkt-Spezifikation; Sub-Area
`tools/harness/` ist über [MR-004](../../../../harness/conventions/MR-004-gate-nachweis-mechanik.md)
konventionsgetragen (siehe `harness/conventions.md` §Modus-Deklaration).

**Berührte Spec-Stellen:** — (Harness-Tooling und `docs/user/`-Prozessdoku,
keine Spec-Stelle).

**Verantwortlich:** claude-sonnet-5.

**Autor:** claude-sonnet-5. **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Beim Release `v0.77.0` aufgetretene macOS-Lokalentwicklungs-
Reibungen beheben bzw. dokumentieren, damit sie beim nächsten Release nicht
erneut live entdeckt werden müssen: (a) einen echten Portabilitäts-Bug im
Datei-Anzahl-Vergleich von `tools/harness/fetch-baseline-cache.sh --verify`
beheben (String- statt Zahlenvergleich bricht mit dem rechtsbündigen
Padding von BSD-`wc -l`); (b) einen **zweiten**, beim Verifizieren von (a)
aufgedeckten Fund in derselben Funktion (`check_aliases`) derselben Datei
beheben: `readlink -e` ist ein GNU-spezifisches Flag, das BSD-`readlink`
(macOS) mit „illegal option" ablehnt — jeder gesunde Symlink unter
`.claude/rules/` wurde dadurch fälschlich als unauflösbar gemeldet; (c) das
im Release-Lauf erprobte Linux-Wrapper-Rezept für `make image-test` auf
macOS in `docs/user/releasing.md` festhalten. (b) ist eine Erweiterung
gegenüber der ursprünglichen Fassung dieses Plans — aufgedeckt beim
Verifizieren von (a) in derselben Datei/Funktion, keine neue,
unabhängige Fundstelle (siehe §7).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine automatisierte Makefile-Zielfunktion** (z. B. `make image-test-linux`),
  die den Linux-Wrapper-Container transparent aufsetzt — bewusst
  zurückgestellt: der Wartungsaufwand einer neuen, dauerhaften
  Automatisierungs-Fläche lohnt für einen Vorgang, der nur einmal pro
  Release relevant ist, nicht gegenüber einer dokumentierten
  Kopier-Anleitung (Auftraggeber-Entscheid im Lauf).
- **Die Behebung der Upstream-Freshness-Befunde** aus dem Nachtlauf-Stand
  (`make nightly-state`: `upstream-drift.yml` ROT seit 2026-09-18T05:27,
  `make freshness-semgrep` und `make go-base-digest` betroffen) — anderer
  Gegenstand (Versions-Pins gegen Upstream), keine Berührung mit der
  `wc`-Portabilität oder `image-test`. Gelesen, aber für diesen Slice nicht
  einschlägig.
- **Eine generelle Portabilitäts-Revision aller `tools/harness/*.sh`-Skripte**
  auf BSD/GNU-Kompatibilität — nur der eine, konkret aufgetretene Fund wird
  behoben; ein Vollaudit auf weitere GNU-Annahmen (`sed`, `grep -P`, …) wäre
  ein eigener, größerer Vorgang mit eigenem Anlass.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

- [x] `tools/harness/fetch-baseline-cache.sh --verify`: Datei-Anzahl-Vergleich
      numerisch (`-eq`) statt String-Vergleich (`=`) — funktioniert
      unabhängig von GNU-/BSD-`wc`-Padding.
- [x] `tools/harness/fetch-baseline-cache.sh` `check_aliases()`: `readlink -f`
      + explizites `[ -e ]` statt `readlink -e` — funktioniert auf GNU wie
      BSD `readlink`; `--selftest` (neun Proben) bleibt grün.
- [x] `docs/user/releasing.md`: Abschnitt mit dem Linux-Wrapper-Rezept für
      `image-test` auf macOS ergänzt (Docker-Socket-Mount, Begründung, warum
      nötig).
- [x] `make gates` grün — **mit Standard-`wc`** (ohne `PATH`-Override), als
      Beleg, dass der Fix trägt.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) —
      [`2026-09-18-slice-230-macos-lokalentwicklung-review-r1.md`](../../../reviews/2026-09-18-slice-230-macos-lokalentwicklung-review-r1.md),
      1 LOW-Finding (Kommentar-Präzision), behoben; kein HIGH/MEDIUM.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen
      (wellenlos — hier geprüft, nicht von einer Welle-Closure).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/fetch-baseline-cache.sh` (Manifest-Deckung) | update | `[ "$on_disk" = "$manifest" ]` → `[ "$on_disk" -eq "$manifest" ]`; behebt den BSD-`wc`-Padding-Bug (reproduziert: `on_disk="      54"` vs. `manifest="54"`, String-Vergleich falsch-negativ) |
| `tools/harness/fetch-baseline-cache.sh` (`check_aliases`) | update | `readlink -e` → `readlink -f` + `[ -e "$tgt" ]`; behebt die BSD-Ablehnung von `-e` (reproduziert: `--selftest` 2 von 9 Proben rot — beide gesunden Alias-Fälle — vor dem Fix, 9 von 9 grün danach) |
| `docs/user/releasing.md` | update | neuer Abschnitt „`image-test` auf macOS" mit dem Docker-Wrapper-Rezept (Linux-Container über den gemounteten Docker-Socket, `TMPDIR` auf einen host-aufgelösten Pfad umgelenkt) |

Kein neuer Testfall — `fetch-baseline-cache.sh` trägt bereits eine
Selbst-Probe (`--selftest`, neun Fälle) für `check_aliases`; sie deckt beide
gefixten Zeilen ab dem Fix hinweg auf (vorher 2 von 9 rot, nachher 9 von 9
grün — reales Vorher/Nachher, kein Mutationstest, aber derselbe Beleg-
Charakter). Der Manifest-Vergleich hat keine eigene Probe; Beleg ist der
reale `make gates`-Lauf mit Standard-`wc` (DoD-Punkt 3).

## 4. Trigger

**Start** (`next` → `in-progress`): sofort — kleine, bereits vollständig
verstandene Korrektur, direkt beansprucht (kein `next`-Zwischenschritt
nötig, WIP-Limit frei).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich beim
  Beheben zeigt, dass weitere `tools/harness/*.sh`-Skripte denselben
  `wc`-Bug tragen und ein Sammel-Fix nötig wird (Umfang würde die
  Abgrenzung aus §1 verletzen).
- `in-progress` → `open` (blockiert): falls `make gates` mit dem Fix aus
  einem anderen, hier nicht vorhergesehenen Grund rot bleibt.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report vorliegt + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Der numerische Vergleich (`-eq`) verlangt, dass beide Operanden reine
  Ziffernfolgen sind — bei `find | wc -l` und `grep -c .` immer der Fall,
  aber ein künftiger Umbau der Zählung dürfte das nicht stillschweigend
  ändern. **Ausgang:** entfallen — beide Quellen sind strukturell
  garantiert numerisch (Zeilenzahlen), kein realistisches Risiko.
- Die Doku-Ergänzung in `releasing.md` könnte beim nächsten
  Baseline-Regelwerk-Bump durch eine geänderte Abschnittsstruktur
  überschrieben/verschoben werden. **Ausgang:** entfallen — `releasing.md`
  ist repo-eigene Prozessdoku, kein vendorter Baseline-Text; ein Bump
  berührt sie nicht.

## 7. Closure-Notiz

- **Was hat funktioniert:** Beide Fixes wurden mit einer echten
  Vorher/Nachher-Probe verifiziert statt blind gepatcht: der `wc`-Bug durch
  direkte Reproduktion (`on_disk="      54"` vs. `manifest="54"`), der
  `readlink`-Bug durch `--selftest` gegen den alten Skriptstand
  (`git show HEAD:…`) — 2 von 9 rot vorher, 9 von 9 grün nachher. Der
  unabhängige Review hat beide Reproduktionen selbst nachvollzogen, nicht
  nur den Plan übernommen.
- **Was ging anders als geplant:** Der ursprüngliche Plan nannte nur den
  `wc`-Bug. Das Verifizieren dieses Fixes (DoD-Punkt 1) deckte einen
  zweiten, verwandten GNU-Annahme-Bug (`readlink -e`) in derselben Funktion
  derselben Datei auf — der Plan (§1, §3) wurde **vor** dem zweiten Fix
  erweitert, nicht nachträglich stillschweigend mitgenommen. Der
  unabhängige Review bestätigt in seinem Negativbefund zur
  Scope-/Abgrenzungs-Treue, dass die Erweiterung vorab deklariert war.
- **Steering-Loop-Eintrag:** gezählt, nicht verkörpert — Beobachtung steht
  bei 1×, unter der Schwelle.
- **Beobachtungs-Register (`../observations/`):** `BEO-HARN/gnu-coreutils-annahme-bricht-auf-bsd-host/`
  neu angelegt, Beleg `evidence/slice-230.md`.
- **Folge-Slices:** keine.
- **Risiken aus §6:** beide entfallen (Begründung dort).
- **Drei Paarungen:** Anker — kein `liegt in`-Feld verwendet (nichts
  verkörpert), kein Gegenstand. Folge-Slice — keiner genannt, kein
  Gegenstand. Register — `BEO-HARN/gnu-coreutils-annahme-bricht-auf-bsd-host/`
  existiert und trägt einen Beleg, geprüft ok.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Zwei berührte Sub-Areas: `tools/harness/`
(eigener Eintrag in der Modus-Deklaration, Kürzel `HARN`) und `docs/user/`
(fällt unter den Default `*`, Kürzel `ALL`). Beide sind bereits als
Sub-Areas deklariert (`harness/conventions.md` §Modus-Deklaration) — keine
Ausdifferenzierung nötig.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-HARN/`
enthält einen Eintrag (`check-latest-blind-before-pin`) — anderes Thema
(Freshness-Check-Blindstelle), keine Berührung mit `wc`-Portabilität.
`BEO-ALL/` enthält keinen Eintrag zu Host-Tool-Portabilität oder
`image-test`. Keine Treffer für den Gegenstand dieses Slice.

**Modus-Begründungsblock:** Beide berührten Sub-Areas sind GF — reiner
Bugfix (Skript) + Prozessdoku-Ergänzung, kein Bootstrap-relevanter Umbau.

### Sub-Area: `tools/harness/`

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — [MR-004](../../../../harness/conventions/MR-004-gate-nachweis-mechanik.md)
  trägt die Gate-Nachweis-Mechanik dieses Verzeichnisses konventionsseitig.
- **Phase-Reife:** Phase 5 (produktiv genutzt, seit mehreren Wellen stabil).
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — Doc führt, der Fix macht das
  Skript konformer zu seiner eigenen dokumentierten Absicht (POSIX-Host-
  Werkzeug-Klasse, `AGENTS.md` §3.1), keine neue Diskrepanz.
- **Reconciliation-Aufwand:** Keiner — kein Brownfield-Bestand hier.

### Sub-Area: `docs/user/`

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — `docs/user/releasing.md` selbst trägt
  bereits die Release-Prep-Struktur, in die dieser Abschnitt sich einfügt.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — reine Ergänzung, keine
  Änderung bestehender Aussagen.
- **Reconciliation-Aufwand:** Keiner.
