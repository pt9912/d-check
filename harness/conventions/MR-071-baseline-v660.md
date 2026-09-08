# MR-071 — Baseline-Pin-Hebung auf `v6.6.0` (dreizehnter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(Pin-Fortschreibung im Bundle-Layout des
  Vorgängers; keine inhaltliche Abweichung)*
- **Datum:** 2026-09-08
- **Geltungsbereich:** §Baseline, pin-gebundene Verweise,
  `.harness/baseline/v6.6.0/`
- **Adaption:** Der vendorte Bestand steht auf
  [`v6.6.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.6.0).

  **Der Delta, gemessen und nicht geschätzt** (`diff -I '<!-- Quelle:'` — die
  Herkunftszeile in Zeile 3 jeder Regelwerk-Datei trägt den Tag und meldete
  sonst *jede* Datei als geändert): 54 Dateien vorher wie nachher, keine neu,
  keine entfallen. **Sechs** von 26 Regelwerk-Dateien tragen ein inhaltliches
  Delta (`README.md`, `grundlagen-harness-dateien.md`,
  `modul-02-harness-bootstrap.md`, `modul-09-implementierung.md`,
  `modul-13-quality-gates.md`, `modul-15-observability.md`), dazu **fünf**
  Templates (`.d-check.yml`, `AGENTS.template.md`, `Makefile`,
  die `README`-Vorlage, `conventions.template.md`).

  **Die Schlagzeile des Delta ist eine Streichung:** `AGENTS.template.md`
  führt die Gates-Tabelle in §4 **nicht mehr** — *„Der Gate-Index steht
  **einmal**, in `harness/README.md` §Sensors … Diese Datei führt die Liste
  nicht."* Die Template-`.d-check.yml` zieht nach und setzt für das Modul
  `targets` sowohl `doc-tables` als auch `authority` auf
  `harness/README.md`. **Die Übernahme dieser Form ist Sache des
  Adoptions-Slice, nicht dieses Eintrags** — hier steht nur, dass sie
  ansteht.

  **Vier Spiegel-Klassen, drei davon gate-blind**
  ([`BEO-ALL/pin-bump-mirrors-ungated`](../../docs/plan/planning/observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)):
  108 Dateien nannten `v6.5.0` mit 299 Vorkommen. Retargetet wurden 77
  Pfad-Verweise in 52 lebenden Dateien plus fünf Release-/Tree-URLs und sieben
  Prosa-Nennungen; **stehen geblieben** sind die eingefrorenen (28 Dateien,
  166 Vorkommen) sowie fünf bewusste Vergangenheits-Aussagen — der
  `ignore-refs`-Tombstone, das wörtliche Fremdzitat eines Adopters im CR und
  dreimal [`MR-067`](../conventions.md#mr-067), der die v6.5.0-Hebung **ist**.

  **Zwei Fehler der mechanischen Ersetzung sind dabei aufgetreten und
  behoben** — beide von der Beobachtung vorhergesagt: eine **Über-Hebung**
  (der Geltungsbereich von `MR-067` sagte plötzlich `v6.6.0`, obwohl der
  Eintrag die v6.5.0-Hebung beschreibt) und eine **übersehene Form** (die
  relative Schreibweise `../baseline/v6.5.0/` <!-- d-check:ignore (der Baum ist mit diesem Eintrag entfernt) --> in
  [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), die das
  Muster `.harness/baseline/v6.5.0/` nicht traf).

  **`ignore-refs` wächst um einen elften Eintrag**
  ([`MR-069`](../conventions.md#mr-069)): zwei eingefrorene Lauf-Belege aus
  slice-217 zitieren die entfernte `README`-Vorlage des alten
  Baums. Ein Lift machte ihre Aussage still falsch — die v6.6.0-Vorlage führt
  gar keine §4-Tabelle mehr, gegen die slice-217 gemessen hatte.

  **Die `d-check:cite`-Spannen brauchten kein Neu-Ankern**
  ([`MR-051`](../conventions.md#mr-051)) — vier Direktiven zeigen in geänderte
  Dateien, und `citations` ist für alle grün. Das ist gemessen, nicht
  vorausgesetzt: Die Regel bleibt, dass ein Bump sie neu ankert, wenn sie sich
  verschieben.
- **Begründung:** Der Nachtlauf meldete den neuen Release
  (`make baseline-freshness`, Exit 3), und der gepinnte Tag war upstream
  inhaltlich unverändert — die Hebung ist damit eine reine Fortschreibung,
  kein Reparatur-Fall.
- **Auflösungs-Trigger:** die nächste Pin-Hebung.
