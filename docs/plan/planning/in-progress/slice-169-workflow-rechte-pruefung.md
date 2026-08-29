# Slice slice-169: Der Wächter prüft die Existenz des Ziels, nicht seine Rechte

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`tools/harness/workflow-pins.sh`](../../../../tools/harness/workflow-pins.sh)
(der Wächter); [`AGENTS.md`](../../../../AGENTS.md) §3.9 (die Hard Rule, die er
trägt) und §3.8 (*„Ein Modul verspricht nur über das, was es scannt"* — die
Regel, an der er scheitert);
[ADR-0068](../../adr/0068-lokale-workflow-referenzen-ohne-pin.md) (die Ausnahme,
die die Existenz-Prüfung an die Stelle des Pins setzt);
[`BEO-004`](../observations.md) (die Beobachtungs-Klasse).

**Berührte Spec-Stellen:** — (Harness-Mechanik, kein Produkt-Vertrag).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**Das Gate war grün, und das Release lief trotzdem nicht an.**

Gemessen am Tag-Push von `v0.66.0` (Run 33235586073): `conclusion:
startup_failure`, **null Jobs**, kein Log, keine Check-Runs. Die Ursache ist
eine Rechte-Vererbung — [`release.yml`](../../../../.github/workflows/release.yml)
führt `permissions: {}` im Kopf, der Job `hub-description` erbt das und ruft
[`hub-description.yml`](../../../../.github/workflows/hub-description.yml) über
eine lokale Referenz auf, dessen Job `sync` `contents: read` verlangt. Ein aufgerufener Workflow kann nur bekommen, was der
aufrufende Job selbst hat; GitHub lehnt den **ganzen Lauf vor dem ersten Job**
ab.

**`make workflow-pins` lief währenddessen grün** und meldete
*„8 uses:-Einträge geprüft, davon 1 lokale Referenz(en) ohne Pin-Pflicht"*.

**Das ist genau die Klasse, die dieses Repo als Hard Rule führt.**
[`AGENTS.md`](../../../../AGENTS.md) §3.8: *ein Modul verspricht nur über das,
was es scannt* — und die Frage, die es beantworten muss, ist *„welche Eingaben
liest es, die es nicht scannt — und gilt dort dieselbe Zusage?"*. Der Wächter
liest die `uses:`-**Zeile**; die referenzierte **Datei** ist eine Eingabe, die
er öffnet (`[ -f "$target" ]`) und nie inhaltlich ansieht. Im Register steht
die Klasse als [`BEO-004`](../observations.md).

**Was [ADR-0068](../../adr/0068-lokale-workflow-referenzen-ohne-pin.md) zusagt
und was nicht.** Sie ersetzt für die lokale Referenz den Pin-Check durch die
Existenz-Frage und schließt: *„Die Ausnahme lässt damit **keinen Eintrag
ungeprüft**."* Das ist über den **Pin** wahr und über die **Lauffähigkeit**
nicht. Die Rechte-Frage kommt in der ADR nirgends vor — sie war beim Schreiben
nicht im Blick, nicht verworfen.

## 2. Vorgehen

1. **Die Prüfung gehört an den Ort, der die lokale Referenz ohnehin auflöst** —
   in den `./`-Zweig von `workflow-pins.sh`, dort wo heute `[ -f "$target" ]`
   steht. Der Wächter öffnet die Zieldatei damit zum zweiten Mal, aber zum
   ersten Mal liest er sie.
2. **Zwei Bedingungen, beide urteilsfrei:**
   - **Deklarationspflicht.** Verlangt das Ziel irgendwo `permissions`, muss
     der aufrufende **Job** selbst ein `permissions:` tragen — die stille
     Vererbung vom Workflow-Kopf ist der Fall, der ausgefallen ist.
   - **Scope-Vergleich.** Für jeden Scope, den das Ziel verlangt, führt der
     Aufrufer ihn mindestens so hoch (`none` < `read` < `write`).
3. **Fail-closed bei jeder Form, die der Wächter nicht sicher liest.** `permissions: read-all` /
   `write-all`, Flow-Style, Anker — was nicht eindeutig zerlegbar ist, wird
   **gemeldet**, nicht übersprungen. Ein Wächter, der bei unbekannter Form
   schweigt, ist derselbe stille Grün-Pfad wie der, den dieser Slice schließt.
4. **Kein YAML-Parser als vierte Toolchain.** `bash` + `grep`/`awk` tragen die
   Block-Form; was darüber hinausgeht, ist ein **Entscheid**, kein Werkzeug
   nebenbei ([`MR-046`](../../../../harness/conventions.md#mr-046)). Trägt die
   Block-Form nicht weit genug, ist das der Befund — dann wird der Umfang
   geschnitten, nicht die Toolchain geweitet.
5. **Proben vor Schärfung.** Der Wächter bekommt Proben nach dem Muster von
   `make guard-probe`: der reale Ausfall-Fall (`{}` gegen `contents: read`), der
   grüne Fall (deklariert und ausreichend), der zu enge Fall (`read` gegen
   `write`) und die unlesbare Form. **Rot gemessen, bevor die Regel scharf
   geht** — sonst ist unbewiesen, dass sie greift.
6. **Retro-Probe gegen den Ausfall.** Die Regel muss den Stand **vor**
   [`4681835`](../../../../.github/workflows/release.yml) rot melden. Eine Regel,
   die den Anlassfall nicht fängt, ist keine.
7. ADR mit Fitness Function; `AGENTS.md` §4 (Target-Beschreibung) und §3.9
   nachziehen; `make gates`; Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Prüfung fremder Workflow-Referenzen** (`owner/repo/...@ref`). Deren
  Inhalt liegt nicht im Repo; die Rechte-Frage wäre dort Netz.
- **Kein Produkt-Modul.** Das ist repo-eigene CI-Hygiene, keine
  Doku-Konsistenz — `d-check` prüft Dokumentation gegen Zustand, nicht
  GitHub-Actions-Semantik. Ein `workflows`-Modul wäre eine andere Anforderung
  und ein eigener Entscheid.
- **Keine Simulation der GitHub-Auflösung.** Der Wächter prüft **eine**
  benannte Fehlerklasse, nicht ob ein Workflow läuft. Was er nicht deckt,
  gehört als Grenze in die ADR.
- **Kein Umbau der bestehenden Pin-Prüfung.** `uses-pin-missing` /
  `uses-pin-untagged` / `uses-local-missing` bleiben, wie sie sind.

## 4. Definition of Done

- [ ] Der Wächter meldet die **stille Vererbung**: ein Job mit lokaler Referenz
      ohne eigenes `permissions:`, dessen Ziel Rechte verlangt, ist ein Befund
      mit eigenem Grund-Code.
- [ ] Der Wächter meldet den **zu engen Aufrufer**: ein Scope, den das Ziel
      höher verlangt, als der Aufrufer ihn führt.
- [ ] **Unlesbare Formen melden**, statt zu schweigen — und die Grenze steht in
      der ADR, nicht nur im Skript-Kommentar.
- [ ] **Retro gemessen:** die Regel meldet den Stand vor `4681835` rot und den
      heutigen grün; beide Ausgaben stehen in der Commit-Botschaft.
- [ ] Proben nach dem Muster von `make guard-probe`, mit Erwartung und Ergebnis
      je Fall.
- [ ] ADR mit Fitness Function und benannten Grenzen; `AGENTS.md` §4 und §3.9
      nachgezogen.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Die Regel wird aus EINEM Fall gezogen** ([`BEO-011`](../observations.md)).
  Der Bestand trägt genau **eine** lokale Workflow-Referenz; alles, was der
  Wächter über die Klasse behauptet, ist an einem Exemplar geeicht. Der Retro-
  Test gegen den echten Ausfall ist die Gegenprobe, aber er deckt nur die eine
  Richtung. — **Ausgang:** *(bei Closure)*
- **Eine Härtung am Rand kippt die Fehlerpolitik, und die Proben sehen es
  nicht** ([`BEO-018`](../observations.md)). Der Wächter ist Teil von
  `make gates`; wird er bei unlesbarer Form fail-closed, kann eine harmlose
  Schreibweise künftig den ganzen inneren Loop rot machen. — **Ausgang:** *(bei
  Closure)*
- **YAML ohne YAML-Parser bleibt eine Näherung.** Einrückungs-Zerlegung mit
  `awk` trägt die Block-Form; Anker, Flow-Style und Mehrfach-Dokumente nicht.
  Die Grenze ist benennbar — die Versuchung, sie durch „noch ein `sed`" zu
  verschieben, ist der eigentliche Punkt
  ([`MR-046`](../../../../harness/conventions.md#mr-046)). — **Ausgang:** *(bei
  Closure)*
- **Der Wächter deckt eine Fehlerklasse, nicht die Lauffähigkeit.** Ein grüner
  Lauf sagt danach „diese eine Klasse liegt nicht vor" — wer mehr hineinliest,
  hat wieder das Versprechen aus [ADR-0068](../../adr/0068-lokale-workflow-referenzen-ohne-pin.md),
  das dieser Slice korrigiert. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei (slice-168 geschlossen).

**Rückführungen:** `in-progress` → `open`, falls sich zeigt, dass die
Block-Form-Zerlegung den realen Bestand nicht trägt — dann ist der Befund ein
anderer (Werkzeug-Entscheid nach [`MR-046`](../../../../harness/conventions.md#mr-046)),
und er gehört vor die Implementierung.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `tools/harness/` — in
  [`harness/conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration
  als **Greenfield** deklariert (adoptierte Harness-Mechanik, konventionsgetragen).
  Der Slice ändert ein Gate-Skript dieser Sub-Area; die Wahl trägt.
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29):
  [`BEO-004`](../observations.md) — die Modul-Grenze nur auf der Quell-Achse
  gedacht: **genau dieser Fall**, der Wächter liest eine Eingabe, über die er
  nichts verspricht;
  [`BEO-011`](../observations.md) — die Regel aus dem Anlass statt aus dem
  Bestand: es gibt **eine** lokale Referenz im Repo;
  [`BEO-018`](../observations.md) — eine Härtung am Rand kippt die
  Fehlerpolitik, und die Proben sehen es nicht;
  [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr fängt, bleibt
  stehen.
- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Nachtläufe
  melden **gruen** — `upstream-drift.yml` zuletzt 2026-08-28T12:25:19Z,
  `image-scan.yml` 2026-08-28T15:25:09Z. **Beide jüngsten Läufe stammen vom
  Vortag**, und das ist genau die benannte Grenze des Schritts: er liest den
  jüngsten Lauf, nicht sein Alter. Der geplante Lauf des 28. war ausgefallen;
  ob der des 29. gefeuert hat, sagt diese Ausgabe **nicht**.

Slice-ID: slice-169. Betroffene IDs: [ADR-0068](../../adr/0068-lokale-workflow-referenzen-ohne-pin.md),
[`BEO-004`](../observations.md). Module: — (Harness-Skript). Gates: `make gates`,
`make workflow-pins`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield)** — `tools/harness/` ist als Greenfield deklariert: adoptierte
Harness-Mechanik, konventionsgetragen. Der Slice erweitert ein vorhandenes
Gate-Skript um eine Bedingung; kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)
