# MR-067 — Baseline-Pin-Hebung auf `v6.5.0` (zwölfter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(Pin-Fortschreibung im Bundle-Layout des
  Vorgängers; keine inhaltliche Abweichung)*
- **Datum:** 2026-09-07
- **Geltungsbereich:** §Baseline, pin-gebundene Verweise,
  `.harness/baseline/v6.5.0/`
- **Adaption:** Der vendorte Bestand steht auf
  [`v6.5.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.5.0).
  **Der Sprung übergeht `v6.4.0`** — beide Releases erschienen zwischen zwei
  Slices, und es gibt keinen Grund, den Zwischenstand zu materialisieren: Das
  Bundle ist self-contained, der Delta wird gegen den **alten Pin** gemessen,
  nicht gegen jede Zwischenstufe.

  **Der Delta, gemessen und nicht geschätzt:** 55 Dateien vorher wie nachher,
  keine neu, keine entfallen. **Zwölf** Markdown-Dateien tragen ein
  inhaltliches Delta — davon fünf die Umsetzung des ausgehenden Change Requests
  dieses Repos (Slice-Plan §1 *Ziel und Abgrenzung*, §8 mit unbedingtem Kopf)
  und sieben neue Regeln, darunter die **RTM** als Kanon-Begriff und die
  Unterscheidung **einfrierendes gegen lebendes Artefakt**. Die Liste steht im
  auslösenden Slice; die **Adoption** ist ein eigener Vorgang und nicht Teil
  dieser Hebung.

  **Die Messmethode des Vorgängers reichte nicht, und das ist der Ertrag
  dieser Hebung.** Die Lehre aus der letzten Pin-Hebung lautete *„Bundle-Delta
  nur mit `diff -I` messen"* — das filtert Versionsnummern und Daten. Damit
  standen **27** Dateien als geändert da; erst `diff -w -B -I` zeigt, dass
  **fünfzehn** davon reines Tabellen-Padding waren. Die größte Einzeldatei
  (81 geänderte Zeilen) enthielt **keine** Regel-Änderung. Wer nur `-I` fährt,
  gibt dem Folge-Slice fünfzehn Dateien zu beurteilen, die nichts sagen.

- **Grenze — eine `Accepted`-ADR wird nicht nachgezogen.** Alle pin-gebundenen
  Verweise wandern mit, auch die in eingefrorenen Dokumenten: Der Pfad ändert
  sich, die Aussage nicht. **Eine Klasse kann das nicht:** Eine `Accepted`-ADR
  ist nach [`AGENTS.md`](../../../AGENTS.md#35-adrs-sind-nach-accepted-immutable)
  §3.5 im Kern unantastbar, und `make adr-check` hat den Versuch gemeldet
  (`core-drift-vcs`). Sie behält ihren Link in den entfernten Baum und bekommt
  stattdessen einen Tombstone-Eintrag in
  [`.d-check.yml`](../../../.d-check.yml) — dieselbe Auflösung wie beim Vorgänger,
  und aus demselben Grund.

  **Der Kanon beendet diese Ausnahme mit genau diesem Release:** `v6.5.0`
  verlangt für einfrierende Artefakte die **Kennung statt der Adresse** — eine
  Baseline-Stelle als Tag plus Pfad in Inline-Code statt als Link. Wer das
  adoptiert, hat beim übernächsten Bump keinen Tombstone mehr nötig. Die
  Adoption liegt im Folge-Slice.


- **Zensus und cite-Konto** — die Zahlen, gegen die der nächste Bump prüfen
  kann: **51** lebende Dateien tragen pin-gebundene Verweise (Briefing,
  Harness-Einstieg, Konventionen samt aktiven Einträgen, beide Skills, zwei
  Spec-Straten, das Prüf-Profil, die Planungs-Indizes). **14**
  `d-check:cite`-Direktiven leben außerhalb des vendorten Baums; **vier**
  mussten neu geankert werden, **zehn** hielten. Dazu **acht** Symlinks unter
  `.claude/rules/` und **drei** URL-Formen mit dem Tag (Release-Download,
  Tree, `@`-Kurzform).

  **Die eingefrorenen Zitate sind NICHT in diesem Konto** — und das ist der
  Punkt: Die `d-check:cite`-Direktiven der `done/`-Slices nennen den Tag, gegen
  den sie geschrieben wurden, und bleiben dort. Ein erster Anlauf hat sie
  mitgehoben; der Review hat es gemessen und zurückgenommen. Der Tombstone-Block
  in [`.d-check.yml`](../../../.d-check.yml) hält die Regel fest, damit der nächste
  Bump nicht dieselbe Bewegung macht.

- **Zitat-Delta nach [`MR-039`](../../conventions.md#mr-039):** Vier Spannen sind
  gewandert, weil das Delta ihre Umgebung umschrieb — die zwei
  Vorprüfungs-Sätze aus `modul-05` (223-224/229-229 → 268-269/274-274, der §1-
  und §8-Umbau schob sie), die Raten-Zeile aus `modul-09` (166 → 175, dieselbe
  Ursache) und die Change-Request-Zeile aus `grundlagen-begriffe` (47 → 48, die
  neue RTM-Zeile davor). **Kein Wortlaut hat sich geändert**, nur die Zeile;
  `citations` hätte einen Wortlaut-Wechsel gemeldet und tat es nicht.

- **Adaptions-Review:** Das Delta hat die Abschnitte umgeschrieben, die drei
  aktive Einträge als ihre Baseline-Regel nennen — geprüft und alle drei
  tragen weiter: [`MR-031`](../../conventions.md#mr-031) (Schritt 3 verlangt
  Benennen) und [`MR-066`](../../conventions.md#mr-066) (was „nicht still
  weiterschieben" verlangt) präzisieren Regeln, die der Kanon unverändert
  führt; [`MR-035`](../../conventions.md#mr-035) (ausgehender CR) bleibt gültig,
  weil der Kanon weiterhin nur den **eingehenden** CR kennt. **Von 39 aktiven
  Auflösungs-Triggern hat genau einer gefeuert:** der von
  [`MR-065`](MR-065-baseline-v631.md), und dieser Eintrag löst ihn ab.
- **Begründung:** Der Nachtlauf meldete die Currency-Achse rot; der
  Content-Drift am gepinnten Tag war grün. Die Hebung schließt die eine, ohne
  die andere zu berühren.
- **Löst auf:** [`MR-065`](MR-065-baseline-v631.md)
- **Ausgelöst durch Baseline-Stand:** v6.5.0
- **Auflösungs-Trigger:** die nächste Pin-Hebung.
