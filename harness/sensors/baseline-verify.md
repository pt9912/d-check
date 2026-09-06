# `make baseline-verify` — prüft den committeten vendorten Baseline-Bestand auf Unversehrtheit

## Vertrag

Integritätsprüfung von `.harness/baseline/<tag>/{regelwerk,templates}/` gegen
`SHA256SUMS`, via
[`fetch-baseline-cache.sh --verify`](../../tools/harness/fetch-baseline-cache.sh).
**Drei Fragen, alle nötig:**

1. **`sha256sum -c`** erkennt **geänderte** und **gelöschte** Dateien.
2. **Manifest-Deckung** zusätzlich **eingelegte** — ohne sie passierte eine
   untermengige, in sich konsistente `SHA256SUMS` grün. Gezählt wird das
   **ganze Tag-Verzeichnis** ohne das Manifest, nicht nur die zwei Bäume: eine
   als Geschwister abgelegte Datei liegt in keinem Baum und in keiner
   Manifest-Zeile und bliebe sonst unsichtbar.
3. **Alias-Auflösung** unter `.claude/rules/` — eine andere Achse: Ein Symlink
   dorthin bindet denselben Pin, wird aber von keinem Modul gescannt und steht
   in keiner Manifest-Zeile; beim Bump bräche er still
   ([`MR-055`](../conventions/MR-055-symlink-als-pin-traeger.md)). Geprüft wird
   **rekursiv und dotfile-bewusst**.

## Grenze — was das Grün nicht abdeckt

1. **Geprüft wird die Auflösung, nicht das Ziel** — ein Alias auf ein
   Verzeichnis passiert. Permanent.
2. **Ein fehlendes `.claude/rules/` ist von „hier gibt es keine Aliase" nicht
   unterscheidbar** — wer die Aliase löscht statt sie umzuhängen, hat einen
   grünen Lauf. Permanent.
3. **Die dritte Frage läuft nach den beiden ersten und akkumuliert nicht.**
4. **Ein Symlink überlebt nicht jedes Dateisystem** — `core.symlinks=false`
   macht Textdateien daraus.

Ihre Proben fährt `make baseline-probe` (neun Fälle, netzlos).

## Ausgabe und Ausgänge

Netzlos — der `--verify`-Pfad ruft kein Netz-Werkzeug —, **fail-closed**.

## Bindung

Bestandteil von `make gates`.
[`MR-011`](../conventions/done/MR-011-baseline-pin-release-tag.md)-Kette ·
[`MR-021`](../conventions/MR-021-vendored-verweise-pin-gebunden.md) ·
[`MR-055`](../conventions/MR-055-symlink-als-pin-traeger.md)
