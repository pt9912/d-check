# E2E-Abdeckung je Anforderung (Image-Test)

Abgeleitet aus `tools/image-test.sh` (`make image-test`): eine Zeile je
Phase, die unter ihrer Kopfzeile einen Anker `# abdeckung:` mit ihren
Kennungen trägt. Geprüft wird das gebaute Image, nativ gegen Container. Die
Datei schreibt `make abdeckung`; `make test` hält sie gegen das Skript.

**Grenze:** eine Deklaration, kein Beleg — die Zeile sagt, dass die Phase
die Anforderung prüfen soll, nicht, dass sie es tut.

| Kennung | Phase | Datei |
| --- | --- | --- |
| [`DC-FA-DIST-001`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `(1) Happy: nativ vs. Container byte-identisch` | [`tools/image-test.sh`](../../tools/image-test.sh) |
| [`DC-FA-DIST-001`](../../spec/lastenheft.md), [`DC-QA-03`](../../spec/lastenheft.md) | `(2) Boundary: read-only-Mount, sauberes Fixture → Exit 0` | [`tools/image-test.sh`](../../tools/image-test.sh) |
| [`DC-FA-DIST-001`](../../spec/lastenheft.md) | `(3) Negative: kein Mount → Exit 2 + Mount-Hinweis` | [`tools/image-test.sh`](../../tools/image-test.sh) |
| [`DC-FA-CLI-007`](../../spec/lastenheft.md), [`DC-FA-CLI-008`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `(4) Modi: --doctor und --repair nativ vs. Container` | [`tools/image-test.sh`](../../tools/image-test.sh) |
