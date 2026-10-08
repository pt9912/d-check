# slice-258: Black-Box-Vorher/Nachher-Probe als make-Target

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
(zwölf Belege, viermal in Folge am Modul `targets`); der Sensor-Kandidat ist
in slice-251 benannt und in slice-252 als offene Auftraggeber-Entscheidung
geführt; Auftraggeber-Freigabe 2026-10-07.

**Berührte Spec-Stellen:** — *(Harness-Werkzeug, keine Spec-Aussage)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Ein make-Target vergleicht das Verhalten des Werkzeugs **vor** und
**nach** einer Änderung mechanisch: Es baut aus einer Git-Referenz ein
Vorher-Image (`git archive <ref>`, Dockerfile und Build-Defaults dieses Stands),
nimmt das aktuelle Image als Nachher und fährt beide über dieselben Fixtures
und Ausgabeformen. Verglichen werden **stdout, stderr und Exit getrennt** —
zusammengeführte Streams erzeugten in slice-252 Schein-Abweichungen. Das
Ergebnis ist eine Liste der Abweichungen je Fixture und Form, oder die
Aussage „byte-identisch über N Vergleiche".

Gegenstand ist die Zusage „ohne Schalter byte-identisch", die fünf der
letzten Slices am Modul `targets` von Hand gemessen haben — jedes Mal mit den
Fällen, an die der Autor dachte. Der Sensor fährt dieselben Fälle jedes Mal.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Gate** — das Target urteilt nicht über den Repo-Zustand, sondern
  meldet den Unterschied zweier Stände; ob ein Unterschied gewollt ist,
  entscheidet der Slice, der ihn erzeugt. Es steht als Werkzeug im Index.
- **Die Lesart einer Regel** — ein Unterschied, den niemand erwartet, findet
  der Sensor; dass eine erwartete Gleichheit die falsche Zusage ist, findet er
  nicht (slice-252 hat das festgehalten). Das bleibt beim fremden Leser.
- **Fixtures für jedes Modul** — der Bestand an Fixtures wächst mit den
  Slices, die sie brauchen; dieser Slice legt die Mechanik und eine
  Grundmenge an.
- **Pflicht im Workflow** (`AGENTS.md` §6, Slice-Vorlage) — erst wenn das
  Target sich in Slices bewährt hat; eine Pflicht ohne Erfahrung wäre eine
  Regel aus dem Anlass.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `tools/blackbox-probe.sh` und `make blackbox-probe REF=<ref>`: baut das <!-- d-check:ignore (Datei entsteht mit diesem Slice) -->
      Vorher-Image, fährt Vorher und Nachher über die Fixtures und die
      Ausgabeformen, vergleicht stdout, stderr und Exit getrennt; Exit 0 bei
      Gleichheit, 1 bei Abweichung (mit Liste), 2 bei gescheitertem Lauf;
      fail-closed bei leerer Fixture- oder Formenmenge.
- [ ] Grundmenge an Fixtures unter `tools/blackbox-probe/fixtures/` plus das <!-- d-check:ignore (Datei entsteht mit diesem Slice) -->
      Repo selbst; Bruch-Test: eine bewusst geänderte Meldung wird als
      Abweichung gemeldet, der unveränderte Stand als byte-identisch.
- [ ] `harness/README.md` (Werkzeug-Zeile, `kein Gate`) und
      `harness/sensors/blackbox-probe.md` (Vertrag, Grenzen); `make gates` <!-- d-check:ignore (Datei entsteht mit diesem Slice) -->
      grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/blackbox-probe.sh` | neu | Vorher-Image, Läufe, getrennter Vergleich | <!-- d-check:ignore (Datei entsteht mit diesem Slice) -->
| `tools/blackbox-probe/fixtures/` | neu | Grundmenge | <!-- d-check:ignore (Datei entsteht mit diesem Slice) -->
| `Makefile` | update | Target `blackbox-probe` |
| `harness/README.md`, `harness/sensors/blackbox-probe.md` | update / neu | Werkzeug-Zeile, Vertrag | <!-- d-check:ignore (Datei entsteht mit diesem Slice) -->

*(Plan-Änderung vor dem Code-Commit: `.d-check.yml` `scan.ignore` nimmt die Fixtures aus dem Dogfooding-Scan — sie tragen absichtlich kaputte Links und nackte Kennungen und sind Eingaben der Probe, kein Doku-Vertrag; bisher Geprüftes fällt dadurch nicht heraus. Dazu `tools/blackbox-probe/README.md` eine Ebene über den Fixtures, damit sie weiter geprüft wird. Nach R1: ein Lauf zählt nur mit einem Exit des Werkzeugs (0, 1, 2) und ohne Docker-Fehlermeldung, sonst ist die Probe gescheitert (F-1); `VERSION` ist auf beiden Seiten fest `0.0.0-dev` (F-4); `make clean` räumt das Vorher-Image ab (F-6); die Sensor-Datei nennt die Grenzen, `harness/sensors/doc-check.md` das dritte Ventil (F-2, F-3, F-5). Nach der Verifikation: ein Lauf zählt nur mit dem Lebenszeichen des Werkzeugs (Exit 1 mit Befund auf stdout, Exit 2 mit `d-check:` auf stderr) statt am Wortlaut einer Docker-Meldung (V-1); die Sensor-Datei nennt den Neubau von `:latest` und „leer heißt Default" (V-2, V-3). Nach R2: das Lebenszeichen gilt für jeden Exit gleich (stdout nicht leer oder eine `d-check:`-Zeile auf stderr — gemessen über sieben Formen, darunter `--repair`), und ein Fall, der auf beiden Seiten mit Exit 2 endet, ist kein Vergleich — die Probe scheitert dann (R2-1, R2-3); die Sensor-Datei nennt die Restgrenze und den Kommandozeilen-Fall von `VERSION` (R2-2, R2-4). Nach R3: die Klassifikation über die Ausgabe des Werkzeugs entfällt — sie wies echte Läufe ab (`-h`, Panic) und verschluckte Unterschiede bei beidseitigem Exit 2 (R3-1, R3-2). Statt dessen ein **Kanarienlauf** vor und nach allen Läufen auf beiden Images (Fixture `sauber`: Exit 0 und genau eine geprüfte Datei) belegt Daemon, Mounts und gesehenen Inhalt; Exit 125/126/127 bricht ab; alles andere wird verglichen, auch Exit 2. Nach R4: auch das Repo wird als lesbare Kopie unter demselben Arbeitsverzeichnis gemountet — eine Mount-Quelle für alle Läufe, die der Kanarienlauf mitbelegt (R4-1); der Kanarienlauf prüft nur Exit-Codes (`sauber` 0, `links` 1) statt einer Zusammenfassungs-Zeile (R4-2); jeder Exit außer 0, 1, 2 bricht ab. Nach der zweiten Verifikation: die Repo-Kopie überspringt gelöschte, noch getrackte Dateien (V2-1), die Kanarien-Meldung unterscheidet Umgebung und Exit-Vertrag des Nachher-Stands (V2-2).)*

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
keine bestehende Zusage wird geändert; neu ist nur das Target. Der
Gate-Index (`make gate-consistency`) verlangt die Werkzeug-Zeile.

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): das Vorher-Image lässt sich aus
  `git archive` nicht mit denselben Build-Args bauen (etwa weil sich die
  Build-Mechanik zwischen den Ständen ändert) und braucht eine eigene
  Versions-Logik.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- **Laufzeit** — zwei Image-Builds plus Läufe über Fixtures und Formen. —
  **Ausgang:** *(offen)*
- **Rauschen** — Ausgaben, die zwischen zwei Läufen desselben Stands
  schwanken (Zeit, Pfade), meldeten Abweichungen ohne Verhaltensänderung;
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) sagt Determinismus zu, gemessen wird es hier zum ersten Mal über
  zwei Images. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** eine berührte Sub-Area: das
Werkzeug des Repos unter dem Default `*` (`ALL`); deklariert. `tools/harness/`
(`HARN`) ist nicht berührt — das Skript liegt unter `tools/`.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
— der Anlass dieses Slice. Dazu
[`BEO-ALL/fix-commit-ausserhalb-review-range`](../observations/BEO-ALL/fix-commit-ausserhalb-review-range/state.md)
(2×): ein Fix, der mehr schreibt als seine Befunde, bekommt eine eigene
Review-Runde.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-08 — `upstream-drift` und `image-scan` grün.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
