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

## Grenze — was das Grün nicht abdeckt

1. **Geprüft wird das lokal gebaute Image, nicht das publizierte** — was
   Anwender ziehen, ist Gegenstand von [`image-scan`](image-scan.md), und das
   ist eine andere Frage mit anderem Bindepunkt.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Kriterien erfüllt |
| 2 | Mount fehlt — mit Hinweis, kein stilles Grün |

## Bindung

Bestandteil von `make ci` und `make fullbuild`.
[`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image) ·
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)
