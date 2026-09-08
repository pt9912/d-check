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
  ([`BEO-ALL/pin-bump-mirrors-ungated`](../../docs/plan/planning/observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)).
  **Die Form der Zählung steht neben der Zahl**, sonst ist sie beliebig:
  *Vorkommen* meint jedes einzelne Auftreten der Zeichenkette (`grep -ro`),
  nicht die Zeile; *Pfad-Verweis* ein Vorkommen mit unmittelbar
  vorausgehendem `baseline/` — die **weite** Form, weil die enge
  (`.harness/baseline/…`) die relativen Schreibweisen übersieht.

  Gemessen am Vorzustand: 108 Dateien nannten `v6.5.0` mit **315** Vorkommen
  auf 302 Zeilen. **Lebend** waren 52 Dateien mit 87 Pfad-Verweisen und 30
  reinen Versionsnennungen, darunter **sieben** Release-/Tree-/Blob-URLs;
  **eingefroren** 28 Dateien mit 169 Vorkommen, dazu die 29 Nennungen im
  vendorten Baum, die mit ihm verschwanden. **Stehen geblieben** sind die
  bewussten Vergangenheits-Aussagen — der `ignore-refs`-Tombstone, das
  wörtliche Fremdzitat eines Adopters im CR und
  [`MR-067`](../conventions.md#mr-067), der die v6.5.0-Hebung **ist**.

  **Drei Fehler der mechanischen Ersetzung sind aufgetreten und behoben**, die
  ersten beiden von der Beobachtung vorhergesagt: eine **Über-Hebung** (der
  Geltungsbereich von `MR-067` sagte plötzlich `v6.6.0`, obwohl der Eintrag die
  v6.5.0-Hebung beschreibt), eine **übersehene Form** (die relative
  Schreibweise `../baseline/v6.5.0/` <!-- d-check:ignore (der Baum ist mit diesem Eintrag entfernt) --> in
  [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), die das
  Muster `.harness/baseline/v6.5.0/` nicht traf) und eine **stehengebliebene
  URL** in [`harness/README.md`](../README.md), die auf den `v6.6.0`-Baum
  zeigte und im selben Satz das `v6.5.0`-Release-Asset nannte — vom Review
  gefunden, nicht vom Lauf.

  **`ignore-refs` wächst um einen elften Eintrag**
  ([`MR-069`](../conventions.md#mr-069)) — **bemessen an den Befunden, nicht an
  den Verweisen**: Ein Lauf ohne ihn meldet **4** Befunde, alle vier
  `target-missing` auf dieselbe Datei, die `README`-Vorlage des alten Baums.
  Die 35 eingefrorenen Pfad-Verweise insgesamt sagen darüber nichts — ein
  Verweis feuert, wenn ein Modul ihn auflöst, und die übrigen sind
  `d-check:cite`-Direktiven in ausgenommenen Verzeichnissen oder stehen in
  Inline-Code. Beide ausgenommenen Artefakte sind **Lauf-Belege**: slice-217
  maß das Zellen-Padding gegen die Tabellen **jener** Vorlage. Ein Lift zeigte
  auf einen anderen Gegenstand als den gemessenen — die v6.6.0-Vorlage führt
  40 Tabellenzeilen statt 39 —, ohne dass es auffiele.

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
