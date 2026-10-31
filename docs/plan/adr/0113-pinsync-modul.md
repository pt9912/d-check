# ADR-0113: Das Modul `pinsync` vergleicht benannte Werte an mehreren Dateistellen zeilenweise

**Status:** Accepted

**Datum:** 2026-10-31

**Autor:** pt9912

**Bezug:** GitHub-Issue 14 (Modul `pinsync`) samt der Antwort des Auftraggebers vom
2026-10-31,
[`DC-FA-PSYNC-001`](../../../spec/lastenheft.md#dc-fa-psync-001--gleichheit-eines-wertes-an-mehreren-stellen-modul-pinsync-opt-in),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[ADR-0031](0031-targets-deklarations-konsistenz-modul.md) (das hermetische Muster, das `pinsync`
fortsetzt), [ADR-0112](0112-metrics-modul.md) (Pfad-Prüfung; dessen Verfahren [`DC-FA-MET-001.a`](../../../spec/spezifikation.md#dc-fa-met-001a--kennzahl-eines-planungsdokuments-gegen-messung-aus-git-metrics) weist einen Plan als Symlink ab),
slice-286 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-PSYNC-001.a`](../../../spec/spezifikation.md#dc-fa-psync-001a--gleichheit-eines-wertes-an-mehreren-stellen-pinsync)
(neues Modul-Verfahren).

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v6.18.1/regelwerk/modul-04-adrs.md).

---

## Kontext

Ein Repo führt dieselbe Version an mehreren Stellen: im `Dockerfile` als `ARG`, in einer
eingebundenen `.mk`-Datei als Variable. Wer eine hebt und die andere vergisst, merkt es erst,
wenn ein Bau auseinanderläuft. Das Issue verlangt ein opt-in-Modul, das benannte Werte an
konfigurierten Stellen auf Gleichheit prüft. Es lässt offen, wie das Modul heißt, wie die
Ausdrücke gelesen werden, was „genau ein Treffer“ zählt, ob ein Pin mit einer Quelle ein Fehler
ist, welche Pfade zulässig sind und wie die Grund-Codes heißen. Der Auftraggeber hat dazu
entschieden: kein Schalter `strip-prefix` (das `v` steuert der Ausdruck), nur Gleichheit, kein
Netz.

Zwei Randbedingungen gelten unabhängig davon: ein Lauf ohne `pinsync`-Konfiguration bleibt
byte-identisch, und das Modul schreibt nie ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Die eigene Prüfung des Repos hat einen Gegenstand: `GO_VERSION` und `GOLANGCI_LINT_VERSION`
stehen heute im `Dockerfile` und im `Makefile`, jeweils einmal.

## Entscheidung

Wir führen ein **eigenes opt-in-Modul `pinsync`** mit diesen Festlegungen (die Nummern dienen
dieser Liste als Bezug):

1. **Name und Kürzel:** Modul `pinsync`, Bereichskürzel `PSYNC`; `PIN` (Content-Pin) und `VER`
   (Versions-Pin gegen `version.md`) sind belegt, und das Modul meint etwas Drittes: die
   Gleichheit mehrerer Dateien.
2. **Leseweise:** zeilenweise (`\n`, ein `\r` am Zeilenende gehört nicht zur Zeile), RE2 je
   Zeile, der Wert ist die Fanggruppe; kein mehrzeiliger Modus.
3. **Ein Treffer ist eine Zeile:** gezählt werden Zeilen, deren **erste** Passung des
   Ausdrucks eine nicht leere Gruppe hat. Mehrere Vorkommen in einer Zeile sind eine Zeile.
4. **Genau eine Fanggruppe**, sonst Exit 2; eine leer fangende Gruppe ist kein Treffer, damit
   ein leerer Wert nie als Vergleichswert gilt.
5. **Ein Pin hat mindestens zwei Quellen**, sonst Exit 2: ein stilles Grün ohne Gegenstelle
   wäre ein Gate ohne Prüfung.
6. **Quellen sind literale, relative Pfade** ohne `..` und Glob, nicht auf `scan.roots`
   beschränkt. Ein **Symlink** — als Datei oder als Verzeichnis-Segment — ist Exit 2. `targets`
   liest eine wörtlich eingetragene Datei über `ReadFile` und folgt einem Symlink still; das
   wird nicht übernommen, weil die Quelle sonst eine Stelle sein dürfte, die kein Scan sieht.
7. **Reihenfolge:** Pins nach Namen, Quellen in Konfigurations-Reihenfolge; die Ausgabe
   sortiert wie bei jedem Modul.
8. **Befund-Ort:** Datei und Zeile der Quelle (Zeile 1 ohne Treffer), `target` ist der
   Pin-Name; bei ungleichen Werten meldet das Modul jede Quelle, die von der **ersten** Quelle
   mit genau einem Treffer abweicht, und die Meldung nennt alle Werte.
9. **Grund-Codes:** `pinsync-mismatch`, `pinsync-no-match`, `pinsync-ambiguous`.
10. **Bezugsmenge:** der Lauf trägt eine Zeile mit der Zahl der Pins und der gelesenen Quellen;
    sie sagt, was gelesen wurde, nicht, dass das Repo vollständig erfasst ist.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, ein Skript je Repo | kein Aufwand hier | genau der Zustand, den das Issue beklagt: jedes Repo baut sein eigenes `grep`/`awk` und pflegt es |
| B — Schalter `strip-prefix` für das `v` | ein Ausdruck für beide Schreibweisen | eine zweite Stelle für dieselbe Entscheidung; der Ausdruck kann `v?` ohnehin; der Auftraggeber schloss ihn aus |
| C — mehrzeiliger Modus (`(?s)` über die ganze Datei) | erreicht Werte über Zeilenfortsetzungen | `^`/`$` und „genau ein Treffer“ verlieren ihre einfache Bedeutung; ein Treffer über Zeilen hat keine Zeilennummer |
| D — Aktualitätsprüfung gegen eine Netz-Quelle | sagt auch, ob der Wert der neueste ist | ein anderes Modul mit Netzzugriff; die Nachtlauf-Targets `freshness-*` tragen es |
| E — Symlinks wie `targets` folgen | gleiches Verhalten wie dort | die Quelle könnte außerhalb des Repos liegen; die Zusage der Scan-Menge gälte dort nicht, und das Grün bliebe still |
| F — ein Pin darf eine Quelle haben | kein Konfigurationsfehler bei Platzhaltern | ein Gate ohne Gegenstelle prüft nichts und meldet Grün |
| **G — eigenes Modul, zeilenweise, Symlinks abgewiesen, Pin ab zwei Quellen** | kleine Fläche, hermetisch, jede Grenze ist eine Fehlerklasse statt eines stillen Grüns | Werte über Zeilenfortsetzungen sind nicht erreichbar; die Gleichheit ist Textgleichheit, `1.27` ≠ `1.27.0` |

## Konsequenzen

- Positiv: die zwei Doppel-Pins des Repos stehen im Gate `make doc-check`; ein weiteres Repo
  braucht nur eine Konfiguration. Das Modul braucht keinen neuen Port.
- Negativ: das Modul liest Dateien, die es nicht scannt — dort gilt die Zusage der Scan-Menge
  nicht, und die Anforderung benennt die Grenze. Ein Ausdruck, der die falsche Stelle trifft,
  fällt nur auf, wenn er mehrdeutig ist; ob er die **richtige** Stelle meint, bleibt die Aussage
  des Ausdrucks. Das Modul deckt von den Spiegeln einer Pin-Hebung nur die Wertgleichheit
  zweier Dateien und löst
  [`BEO-ALL/pin-bump-mirrors-ungated`](../planning/observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  nicht auf.
- Folgepflicht: jede Randbedingung der Anforderung (Eingabefehler, Lesen, Treffer, Vergleich,
  Reihenfolge) hat einen Test im Modulkern; das Modul wird zusammen mit der Modul-Aufzählung
  (Lastenheft, READMEs, Handbuch, `operations.md`) registriert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | die YAML-Bibliothek ist nur in den benannten Adaptern importierbar, nicht im Kern; der Kern liest Dateien nur über den Filesystem-Port | `make arch-check` |
| Spiegel-Test | `model.ValidModules()` stimmt mit Lastenheft, READMEs, Handbuch und `operations.md` überein | `make test` |
| Selbstlauf (wirkt erst, wenn das Modul registriert ist und das eigene Profil den Block trägt) | die zwei Doppel-Pins von `Dockerfile` und `Makefile` sind im eigenen Profil, beidseitig belegt | `make doc-check` |

## Re-Evaluierungs-Trigger

Die Entscheidung wird neu geprüft, wenn (a) ein Konsument einen Wert prüfen will, der über
mehrere Zeilen läuft (dann mehrzeiliger Modus, Alternative C), (b) ein zweites Repo eine
Normalisierung der Werte braucht, die der Ausdruck nicht ausdrücken kann (dann `strip-prefix`
oder eine Nachfolge, Alternative B), oder (c) ein Konsument Quellen über einen Symlink führen
muss (dann ein benannter, geprüfter Ausnahme-Weg statt Folgen, Alternative E).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-31 | Accepted: Entscheide 1 bis 9 des Plans vom Auftraggeber bestätigt (kein `strip-prefix`, nur Gleichheit, kein Netz); Punkt 6 (Symlink) weicht von der Empfehlung ab, die Punkte 4, 5 und 8 sind geschärft und Punkt 10 neu herausgelöst, im Bericht an den Auftraggeber benannt | slice-286 <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
