# Slice slice-145: Ein Sensor auf Pfad-Token in der Architektur-Sicht

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.4 (zweiter Auflösungs-Trigger); [`MR-033`](../../../../harness/conventions.md#mr-033); [slice-136](welle-84/slice-136-agents-34-klaerung.md).

**Berührte Spec-Stellen:** — (abhängig von der Wegwahl; ein Produkt-Delta bräuchte Bump und ADR).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

§3.4 verbietet der Architektur-Sicht Modul-Pfade — eine Verschärfung gegenüber
der Baseline, geführt als [`MR-033`](../../../../harness/conventions.md#mr-033).
Die Regel zerfällt in zwei Hälften: *Rollen statt Technologie* ist Urteil,
**Modul-Pfad ja/nein** ist ein detektierbarer Zustand. Gemessen: die Sicht trägt
heute **null** solcher Token.

Kein heutiges Modul trägt die Prüfung — `codepaths` **findet** solche Token,
verbietet sie aber nicht.

## 2. Vorgehen

0. **Die Schwelle steht, bevor gemessen wird** — festgehalten am 2026-08-26,
   vor dem ersten Lauf und vor jedem Blick in die Sicht
   ([`BEO-011`](../observations.md); dieselbe Form wie in
   [slice-144](../done/slice-144-commit-hash-muster.md), wo sie getragen hat):

   - **Null Falsch-Positive** in `spec/architecture.md`. Die Sicht ist die
     einzige bewachte Datei; ein Befund dort, der kein Modul-Pfad ist, kippt den
     Kandidaten.
   - **Positivkontrolle:** ein eingesetzter Code-Modul-Pfad wird gemeldet, die
     Ursache gelesen, der Rückbau grün. Ohne sie ist „null Befunde" von „die
     Regel tut nichts" nicht zu unterscheiden.
   - **Die Falsch-Negativ-Klasse ist zu benennen**, nicht zu minimieren.
   - **Keine dritte Mechanik** (§3): die Wahl fällt zwischen der aus
     [slice-143](../done/slice-143-structure-abschnitts-skopus.md) und der aus
     [slice-144](../done/slice-144-commit-hash-muster.md), und sie wird
     **begründet**, nicht bequem getroffen.

1. **Ort und Form klären**: ein Muster-Verbot (dieselbe Fähigkeit, die
   [slice-144](../done/slice-144-commit-hash-muster.md) und
   [slice-143](../done/slice-143-structure-abschnitts-skopus.md) berühren) oder ein
   Ventil an `codepaths`. Die drei Slices teilen eine Frage — das gehört
   gesehen, bevor drei Mechaniken entstehen.
2. Am Bestand messen; die Sicht startet grün, andere Dateien womöglich nicht.
3. Konstruierter Verstoß mit gelesener Ursache; Rückbau grün.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine dritte Muster-Mechanik.** Fällt die Entscheidung in slice-143 oder
  slice-144 anders, folgt dieser Slice ihr — er erfindet nichts Eigenes.
- **Keine Ausweitung auf die anderen Spec-Straten.** Die Spezifikation **darf**
  Pfade führen.

## 4. Definition of Done

- [x] Die Wegwahl ist abgeglichen — **im zweiten Anlauf**. Der erste folgte
      [slice-143](slice-143-structure-abschnitts-skopus.md) und war falsch; der
      zweite folgt [slice-144](slice-144-commit-hash-muster.md), und der Grund
      steht in §9.
- [x] **Sieben Formen** einzeln gefahren: vier gemeldet (Inline-Code, blank,
      `cmd/`, Link-Ziel), zwei still als benannte Grenze (Großschreibung,
      Fenced Block), dazu die Gegenprobe, dass die Commit-Hash-Regel in der
      Sicht weiter greift. Rückbau je grün.
- [x] §3.4 trägt die eingelöste Hälfte samt ihren Grenzen und dem **Preis der
      eigenen Klasse**; die Urteils-Hälfte bleibt *permanent*.
- [x] `make gates` Exit 0 (zehn Glieder), `make fullbuild` Exit 0; unabhängiger
      Review ([Report](../../../reviews/2026-08-26-slice-145-pfad-token-sicht-review.md)),
      blockierend mit **einem HIGH** — alle neun eingearbeitet, und der HIGH hat
      den Entwurf ersetzt, nicht nur korrigiert.

## 5. Abnahme-Punkte / Risiken

- **Drei Slices, eine Frage.** Wer sie einzeln löst, baut drei Antworten auf
  dieselbe Lexik-Frage — genau der Defekt, den dieses Repo an fremden Skripten
  ablehnt. — **Ausgang:** *eingetreten, und beinahe unbemerkt.* Der erste
  Entwurf folgte [slice-143](slice-143-structure-abschnitts-skopus.md) und
  erfüllte die Vorgabe *dem Buchstaben nach* — dieselbe Mechanik, keine dritte.
  Er war trotzdem falsch: `structure` liest den **bereinigten** Text, und die
  Form, um die es geht, steht in diesem Repo zu **869 von 901** Fällen in
  Inline-Code. Die Vorgabe „keine dritte Mechanik" schützt vor Wildwuchs, nicht
  vor der **falschen** der beiden vorhandenen. Am Ende trägt die Frage genau
  eine Mechanik — die von
  [slice-144](slice-144-commit-hash-muster.md).

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Wegwahl ein Produkt-Delta verlangt, das in einen der verwandten Slices gehört.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-012`](../observations.md), weil der Slice auf eine Verschärfung aufsetzt, deren Grundlage eine Lesart ist.

Slice-ID: slice-145. Betroffene IDs: — (abhängig von der Wegwahl). Module: `codepaths` bzw. Muster-Verbot.
Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Sensor auf eigenem Bestand.

## 9. Closure-Notiz (nach `done/`)

Geliefert: Die detektierbare Hälfte von §3.4s Sprachfreiheit ist gedeckt — die
Sicht ist eine **eigene `matrix`-Klasse**, und ein Modul-Pfad in ihr ist ein
`matrix-forbidden`-Befund.

**Der erste Entwurf deckte 3,5 Prozent.** Er folgte
[slice-143](slice-143-structure-abschnitts-skopus.md) und benutzte
`structure.forbid-pattern`. `structure` liest den **bereinigten** Abschnitts-Text
— und **869 von 901** solcher Pfad-Token stehen in diesem Repo in Inline-Code.
Ein Modul-Pfad in Backticks passierte den vollen `doc-check` mit Exit 0.

**Schlimmer als die Lücke war, wie ich sie beschrieben habe.** Im Kommentar und
in §3.4 stand, die Sicht dürfe einen Pfad *„zitieren, ohne ihn zu führen"*. Diese
Unterscheidung kennt weder §3.4 noch
[`MR-033`](../../../../harness/conventions.md#mr-033) noch der Kanon — ich habe
sie erfunden, um eine **Schwäche des Werkzeugs als Absicht** zu lesen. Das ist
die teuerste Form von Harness-Lüge: nicht ein Gate, das schweigt, sondern eine
Zusage, die das Schweigen rechtfertigt.

**Meine Begründung für die Wegwahl widerlegte sich selbst.** Ich verwarf
`matrix` mit dem Argument, ein Pfad-Token in *einer* Datei sei keine Referenz
zwischen Dokumentklassen. Dasselbe gilt für Commit-Hashes — und die habe ich
noch am selben Tag als `matrix`-Klasse gebaut
([slice-144](slice-144-commit-hash-muster.md)). Der echte Blocker war ein
anderer, und er stand bei mir nur als ungemessener Nebensatz: die
Klassen-Zuordnung ist **First-Match**.

**Gemessen statt abgewogen — daran hing die Lösung.** Nimmt man die Sicht in
eine **eigene Klasse**, deklariert **vor** den übrigen Straten, dann greift
`matrix`, und `matrix` liest die rohen Zeilen. Der Preis: die Sicht fällt aus
der klasseninternen Richtungs-Regel; ihre Abwärts-Kante und die vier geteilten
Verbote sind als **explizite** Regeln nachgebaut. Die Richtungs-Gegenprobe zeigt
keine Regression.

**Was den Umbau ausgelöst hat, war eine Rückfrage.** Der Review hatte den HIGH
gemeldet; ich hatte daraufhin die Beschreibung korrigiert und die Lücke in einen
Folge-Slice geschoben. Erst die Frage *„kann man das nicht besser machen"* hat
mich messen lassen, statt zu verwalten — und die Messung stand in zwanzig
Minuten. **Eine Lücke zu dokumentieren ist billiger als sie zu schließen, und
genau darum ist es die verdächtigere Antwort.**

**Ein Werkzeug-Fehler gehört dazu:** Beim Rückbau einer Sonde nahm ein
`git checkout` auf die Konfigurationsdatei meine noch nicht committeten
Korrekturen mit. Wiederhergestellt — aber eine Sonde, die dieselbe Datei
anfasst wie die laufende Arbeit, braucht ihre eigene Sicherung, nicht `git`.

**Kein Registereintrag.** Die Klasse *„Blindstelle als Erlaubnis gelesen"* ist
eine Ausprägung von [`BEO-011`](../observations.md) — eine Regel aus dem Anlass
statt aus dem Bestand — und dort mit Zähler 4 geführt; sie ein zweites Mal zu
benennen wäre die Umformulierung, vor der derselbe Eintrag warnt.
