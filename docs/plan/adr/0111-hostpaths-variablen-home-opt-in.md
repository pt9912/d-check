# ADR-0111: `hostpaths` erkennt Variablen-Home-Pfade über einen opt-in-Schlüssel

**Status:** Proposed

**Datum:** 2026-10-29

**Autor:** pt9912

**Bezug:** [`DC-FA-HOST-001`](../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)
(das Modul `hostpaths`);
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Byte-Identität ohne den Schlüssel);
[ADR-0110](0110-trace-default-erkennt-namens-slices.md) (derselbe Zuschnitt: ein
Feldbefund, eine Erkennungsform, ein bisher grüner Lauf bleibt grün);
slice-280 <!-- d-check:status-provenance -->.

**Schärft:** [`DC-FA-HOST-001.a`](../../../spec/spezifikation.md#dc-fa-host-001a--host-pfad-erkennung)
Schritte 2 und 4 und die Schema-Zeilen `hostpaths.variables` und
`hostpaths.exempt-targets`.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`hostpaths` meldet Maschinen-Layout-Leaks: absolute Pfade unter einem Host-Präfix
und Home-relative Pfade mit Tilde. `$HOME/Dev/x/y.log`, `${HOME}/Dev/z` und
`"$HOME"/Dev/w` melden nichts. Die Ursache liegt im Unix-Muster: dessen
Wortgrenzen-Vorbedingung sieht vor `/Dev…` das `E` von `HOME`. Weil `hostpaths`
das einzige Modul ist, das Inline-Code mitliest, ist das ein blinder Fleck genau
dort, wo solche Pfade leben.

Eine frühere Fassung hatte die Variablen-Formen ausgeschlossen, weil kein belegter
Anlass bestand (null Vorkommen im Bestand). Der Anlass ist jetzt belegt, der
Ausschluss fällt für die genannten Namen.

Der Wunsch des Feldes nennt zwei Dinge, die nicht zugleich gehen: die Variable im
Default erkennen, und einen bisher grünen Lauf nicht ungefragt rot werden lassen.

## Entscheidung

### 1. Ein neuer opt-in-Schlüssel `hostpaths.variables`

Eine Liste von Shell-Variablen-Namen (`[A-Za-z_][A-Za-z0-9_]*`, ohne `$`). Ist sie
abwesend oder leer, bleibt der Befundsatz byte-identisch. Das Konfigurations-Gerüst
schlägt `[HOME]` vor; ein leerer oder anders geformter Name ist Exit 2.

### 2. Die Erkennung folgt der Tilde

`$NAME/…`, `${NAME}/…` und `"$NAME"/…`, dann ein erstes Segment, das weder mit `.`
beginnt noch `/` ist — die Punkt-Segment-Regel der Tilde, damit
`$HOME/.config/x` still bleibt. Die Wortgrenze ist die der Tilde, ergänzt um `$` und
Backslash (kein `$$HOME`, keine maskierte Variable). Nicht erkannt werden die nackte
Variable, ein längerer Name (`$HOMEDIR/x`), Ausdrucks-Formen (`${HOME:-/tmp}/x`) und
ein Zeichen zwischen `}` und `/`.

### 3. Gemeldet wird einmal, in voller Form

Der Pfad wird mit der Variablen gemeldet, wie geschrieben; das schließende `"` von
`"$HOME"/x` gehört nicht dazu, damit `exempt-targets` einheitlich greift. Ein
Unix-Treffer, der innerhalb eines Variablen-Treffers beginnt, entfällt.

### 4. `exempt-targets` akzeptiert `$`-Globs mit gelisteter Variable

Ein Glob, das mit `$NAME/` oder `${NAME}/` einer in `variables` genannten Variablen
beginnt, ist gültig. Ein `$`-Glob ohne passende Variable träfe nie und ist Exit 2 —
die Rand-Disziplin des Lastenhefts: ein Glob, das nie trifft, ist ein
Konfigurationsfehler.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| Default `[HOME]` an | schließt die Lücke ohne Handlung des Adopters | ein bisher grüner Lauf wird ungefragt rot; widerspricht der Byte-Identität ohne Schlüssel |
| Gewählt: opt-in, abwesend = aus | bestehende Läufe bleiben byte-identisch; das Gerüst führt den Vorschlag | die Lücke bleibt für Repos, die den Schlüssel nicht setzen |
| `$HOME` als Sonderfall im Unix-Muster | kein neuer Schlüssel | verteilt eine Variablen-Lexik in ein Muster, das Präfix-Namen kennt; weitere Variablen wären neue Sonderfälle |
| Auch `%USERPROFILE%` und `$XDG_*` | breiter | nicht belegt; `%NAME%` ist eine andere Syntax und bleibt ausgeschlossen |

**Fitness-Funktion:**

- Die Tests in `internal/hexagon/core/rules` (`hostpaths_variables_test.go`) fahren
  Happy, Boundary, Negativliste, Mehrfach-Variablen, Verdrängung und Ventil; die
  Config-Tests in `configyaml_test.go` den Rand. Der `$HOME`-Test wird ohne die
  Variablen-Erkennung rot (leere Funde), der Verdrängungs-Test ohne die Verdrängung
  (ein zweiter Unix-Fund innerhalb des Variablen-Pfads) — jede Stelle einzeln gebrochen und gesehen.

## Konsequenzen

- Ein Repo ohne den Schlüssel sieht keine Änderung; das Repo selbst schaltet ihn
  nicht ein (es führt `$HOME` in eingefrorenen Dokumenten).
- **Grenze:** Die Hervorhebung am Rand verhält sich wie bei der Tilde —
  `**$HOME/x/y**` meldet mit den Stern-Zeichen. Eine Variable, die ein Werkzeug erst
  zur Laufzeit auflöst, ist kein Layout-Leak; die Erkennung kann das nicht
  unterscheiden und meldet jede gelistete.
- **Grenze:** Die Verdrängung eines Unix-Treffers greift nur, wenn vor dem
  Präfix-Pfad innerhalb des Variablen-Pfads ein Zeichen steht, das die Wortgrenze des
  Unix-Musters erfüllt (etwa `=`). Ist das erste Segment des Variablen-Pfads selbst
  ein Präfix-Name, trifft das Unix-Muster ohnehin nicht, weil davor das letzte Zeichen
  des Variablen-Namens steht.
- Weitere Variablen sind ein Eintrag in der Liste, keine neue Anforderung.

## Re-Evaluierungs-Trigger

- Ein Adopter meldet einen Fund, den die Variablen-Erkennung zu Unrecht liefert, und
  ein `exempt-targets`-Glob ist keine zumutbare Abhilfe.
- Eine zweite Variable (`$USERPROFILE`, `$XDG_*`) wird belegt — dann ist die
  Punkt-Segment-Regel für sie zu prüfen.
- Das Repo schaltet `hostpaths.variables` für sich selbst ein.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-29 | Entwurf zu slice-280 nach dem Feldbefund (GitHub-Issue 12). Status Proposed, bis die Umsetzung geprüft ist. |
