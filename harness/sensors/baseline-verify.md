# `make baseline-verify` — prüft den committeten vendorten Baseline-Bestand auf Unversehrtheit

## Vertrag

Prüft den committeten vendorten Baseline-Bestand auf Unversehrtheit, über
[`fetch-baseline-cache.sh --verify`](../../tools/harness/fetch-baseline-cache.sh).
Die drei Fragen, die alle nötig sind, legt
[`SPEC-092`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) fest.

## Grenze — was das Grün nicht abdeckt

0. **Der Lauf beweist innere Konsistenz, nicht Echtheit** — und das ist die
   Grenze, die alle folgenden überwiegt. Geprüft wird der Baum gegen das
   `SHA256SUMS`, **das mit ihm kam**. Wer eine Datei ändert **und** das Manifest
   nachzieht, bekommt einen grünen Lauf. Gemessen: dieselbe Änderung meldet ohne
   nachgezogenes Manifest `GESCHEITERT` (Exit 1), mit nachgezogenem Manifest
   `verify ok (54 Dateien, vollständig)` (Exit 0).
   **Das ist kein Mangel des Gates, sondern die Grenze eines netzlosen
   Vergleichs:** Ohne Netz gibt es nichts, wogegen die Echtheit zu prüfen wäre.
   Sie hält [`make baseline-freshness`](baseline-freshness.md) mit
   `--check-latest` — dieselbe Manipulation ⇒ Exit 4,
   `UPSTREAM-CONTENT-DRIFT`. **Der Träger ist Netz, fail-open und kein Gate**:
   Die Zusage *„das Original-Bundle wird unverändert verwendet"* hängt damit am
   Nachtlauf, nicht am inneren Loop. Das Skript sagt es im Kopf
   (*„Integrität ist nicht Aktualität"*); hier stand es bisher nicht.
   Permanent — solange `gates` netzlos bleibt.
1. **Geprüft wird die Auflösung, nicht das Ziel** — ein Alias auf ein
   Verzeichnis passiert. Permanent.
2. **Ein fehlendes `.claude/rules/` ist von „hier gibt es keine Aliase" nicht
   unterscheidbar** — wer die Aliase löscht statt sie umzuhängen, hat einen
   grünen Lauf. Permanent.
3. **Die dritte Frage läuft nach den beiden ersten und akkumuliert nicht.**
4. **Ein Symlink überlebt nicht jedes Dateisystem** — `core.symlinks=false`
   macht Textdateien daraus.

Ihre Proben fährt `make baseline-probe` (neun Fälle, netzlos) — **die
Echtheits-Grenze (0) hat keine**, denn sie ist keine Eigenschaft der
Alias-Auflösung; ihr Bruch-Test steht in `slice-212`.

## Ausgabe und Ausgänge

Netzlos — der `--verify`-Pfad ruft kein Netz-Werkzeug —, **fail-closed**.

## Bindung

Bestandteil von `make gates`.
[`MR-011`](../conventions/done/MR-011-baseline-pin-release-tag.md)-Kette ·
[`MR-021`](../conventions/MR-021-vendored-verweise-pin-gebunden.md) ·
[`MR-055`](../conventions/MR-055-symlink-als-pin-traeger.md)
