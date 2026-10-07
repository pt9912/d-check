# `make image-test` — prüft die Distributions-Akzeptanzkriterien gegen das lokal gebaute Image

## Vertrag

Die Akzeptanzkriterien von
[`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image) gegen
das lokal gebaute Image (`tools/image-test.sh`):

- Befund-Ausgabe und Exit-Code **nativ vs. Container byte-identisch** — das ist
  die Determinismus-Zusage an der Distributionsgrenze;
- der read-only-Mount ist vollständig;
- ein **fehlender** Mount endet mit Exit 2 und einem Hinweis, nicht mit einer
  stillen Leermenge.

`make image-test` prüft die Variante der Host-Plattform (`$(IMAGE):latest`),
`make image-test-arm64` dieselben Kriterien gegen die `linux/arm64`-Variante
(`$(IMAGE):arm64`), Binary und Container auf dieser Plattform. Vor dem Vergleich
wird die ELF-Maschine des extrahierten Binaries gegen die verlangte Plattform
geprüft — ein Bild, das still die Host-Variante liefert, prüfte sonst die
falsche ([ADR-0102](../../docs/plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)).

## Grenze — was das Grün nicht abdeckt

1. **Geprüft wird das lokal gebaute Image, nicht das publizierte** — was
   Anwender ziehen, ist Gegenstand von [`image-scan`](image-scan.md), und das
   ist eine andere Frage mit anderem Bindepunkt. Dass das publizierte Binary
   das geprüfte ist, belegt im Release-Pfad `make image-publish` (Binary je
   Plattform aus dem gepushten Index gegen das geprüfte).
2. **`make image-test-arm64` braucht binfmt/QEMU für `arm64` auf dem Host.**
   Fehlt es, bricht der Lauf mit Hinweis ab (Exit 1), statt still zu
   überspringen. Die Release-Pipeline richtet es ein; lokal ist es Sache des
   Hosts.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Kriterien erfüllt |
| 1 | ein Kriterium verfehlt, die Binary-Plattform falsch oder das Binary auf dem Host nicht ausführbar — die Meldung nennt, welches |

## Bindung

`make image-test`: Bestandteil von `make ci` und `make fullbuild`.
`make image-test-arm64`: Schritt in `release.yml`, nicht in `ci`.
[`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image) ·
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) ·
[ADR-0102](../../docs/plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)
