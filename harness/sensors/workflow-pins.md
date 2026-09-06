# `make workflow-pins` — hält die Deklarations-Form aller `uses:`-Referenzen in den Workflows

## Vertrag

Via Modul `workflows` (Image, dogfood). Zwei Referenz-Klassen, zwei Zusagen:

- **Fremde Referenz:** voller 40-stelliger Commit-SHA (`uses-pin-missing`) plus
  Tag-Kommentar dahinter (`uses-pin-untagged`).
- **Lokale Referenz** (`./…`): das Ziel existiert (`uses-local-missing`) **und**
  der aufrufende **Job** führt die Rechte, die es verlangt
  (`uses-local-perms-undeclared`, `uses-local-perms-narrow`) — ein Job ohne
  eigenes `permissions:` kann nichts weitergeben, was er nicht deklariert.

Unlesbares YAML **meldet** (`workflow-unparsable`), statt übersprungen zu
werden. Die Referenzen kommen aus dem **YAML-Baum**, nicht aus einer
Textsuche; das Verzeichnis ist konfigurierbar (`workflows.dir`).

## Grenze — was das Grün nicht abdeckt

1. **Geprüft ist die Form des Pins, nicht seine Gültigkeit** — ob der SHA
   existiert und den Commit bezeichnet, den der Tag-Kommentar behauptet, wäre
   Netz. Permanent; die Gegenprobe fahren die Action-Pin-Freshness-Achsen.
2. **Als Pin gilt ein 40-stelliger git-SHA** — kürzere Formen erkennt er nicht
   als Pin.
3. **Er deckt eine Deklarations-Klasse, nicht die Lauffähigkeit** eines
   Workflows. Die Rechte-Zerlegung ist eine Näherung über die YAML-Block-Form,
   keine Parser-Zusage.
4. **Er liest die Ziele lokaler Referenzen, die er nicht scannt** — dieselbe
   Parse-Zusage gilt dort (`AGENTS.md` §3.8).

## Ausgabe und Ausgänge

Netzlos, hermetisch, **fail-closed auch bei leerer Prüfmenge**. Liest `.yml`
**und** `.yaml`; die Zahl der lokalen Referenzen steht in der Erfolgsmeldung,
statt stillschweigend übergangen zu werden.

## Bindung

Bestandteil von `make gates`.
[ADR-0072](../../docs/plan/adr/0072-workflows-modul.md) ·
[`DC-FA-WF-001`](../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)
