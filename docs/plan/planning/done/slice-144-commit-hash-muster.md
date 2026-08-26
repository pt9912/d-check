# Slice slice-144: Commit-Hashes in den Spec-Straten — ein Muster mit vertretbarer Falsch-Positiv-Last

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [slice-138](../done/slice-138-matrix-wellen-klasse.md); [`AGENTS.md`](../../../../AGENTS.md) §3.4; [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix).

**Berührte Spec-Stellen:** — (Konfigurations-Profil; eine Anforderung wächst nur, falls das Muster eine neue Fähigkeit braucht).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

§3.4 verbietet den Spec-Straten fünf Referenz-Kategorien; drei sind gedeckt.
Der **Commit-Hash** wäre als Token-Klasse ausdrückbar — die Mechanik existiert —,
aber ein Muster über Hex-Zeichenketten träfe jedes Wort, das wie ein Hash
aussieht. Was fehlt, ist **Präzision**, nicht die Fähigkeit.

Der Slice beantwortet genau eine Frage: **gibt es ein Muster, dessen
Falsch-Positiv-Last am eigenen Bestand vertretbar ist?** Wenn nein, ist das das
Ergebnis, und die Kategorie bleibt ausgewiesen.

## 2. Vorgehen

0. **Die Schwelle steht, bevor gemessen wird** — festgehalten am 2026-08-26,
   vor dem ersten Lauf, weil „vertretbar" sonst nachträglich an das Ergebnis
   angepasst wird ([`BEO-011`](../observations.md), §7):

   - **Null Falsch-Positive** im Bestand der drei Spec-Straten. Ein Gate in
     `gates`, das legitime Sätze meldet, erzwingt Umformulierung — und §3
     verbietet die Ausnahme-Liste als Ersatz für Präzision.
   - **Positivkontrolle:** ein konstruierter echter Commit-Hash wird gemeldet,
     die Ursache gelesen, der Rückbau grün.
   - **Die Falsch-Negativ-Klasse ist zu benennen**, nicht zu minimieren: welche
     echten Hashes das Muster nicht fängt, gehört ins Ergebnis.

   Ein Kandidat, der die erste Bedingung nur mit Ausnahmen erfüllt, ist
   durchgefallen.

1. **Kandidaten-Muster formulieren** (Mindestlänge, Wortgrenzen, Ausschluss von
   Inline-Code und Links) und je Kandidat am Bestand messen.
2. **Beide Fehlerrichtungen zählen**, nicht nur die Treffer: was würde gemeldet,
   das legitim ist — und was bliebe unentdeckt?
3. Nur scharfschalten, wenn die Messung trägt; sonst §3.4s Ausweisung schärfen.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Ausnahme-Liste als Ersatz für Präzision.** Ein Muster, das nur mit
  einer langen Ausnahme-Liste grün wird, ist das falsche Muster.
- **Keine Ausweitung auf andere Dokumentklassen.**

## 4. Definition of Done

- [x] Vier Kandidaten, je ein Lauf **mit dem Produkt**, beide Fehlerrichtungen:
      Treffer im Bestand (**null**, alle vier) und Positivkontrolle — die erst
      die Kandidaten trennt (§9).
- [x] **Scharfgeschaltet** als `matrix`-Klasse `commit-hash`; konstruierter
      Verstoß am scharfen Gate rot, der Befund nennt den Hash selbst, Rückbau
      grün.
- [x] §3.4 sagt jetzt, was gilt — und nach dem Review auch, was **nicht** gilt:
      ADRs sind über die Link-Prüfung gedeckt, nicht als Token; die Klasse sagt
      mehr zu, als ihr Name sagt; drei Grenzen sind benannt.
- [x] `make gates` Exit 0 (zehn Glieder), `make fullbuild` Exit 0; unabhängiger
      Review ([Report](../../../reviews/2026-08-26-slice-144-commit-hash-muster-review.md)),
      blockierend mit **einem HIGH**, fünf MEDIUM, vier LOW — alle elf
      eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Ein Hex-Muster ist ein Heuristik-Wächter.** Genau die Sorte, die welle-84
  ausgeschlossen hat — hier ist sie zulässig **nur**, wenn die Messung sie trägt.
  — **Ausgang:** *entfallen — die Messung trägt sie, aber knapper als zuerst
  beschrieben.* Im Bestand der drei Straten meldet die Klasse null. Ihre
  Toleranz ruht auf einer **Eigenschaft der Straten**, nicht auf der Präzision
  des Musters: deren Datums- und Versions-Schreibweisen tragen Trennzeichen und
  fallen an der Wortgrenze auseinander. Das steht jetzt so da — die erste
  Fassung nannte eine engere Risiko-Klasse, als das Muster hat.
- **„Vertretbar" ist ein Urteil.** Die Schwelle gehört vorher benannt, nicht
  nachträglich an das Ergebnis angepasst. — **Ausgang:** *entfallen — und zwar
  belegbar.* Die drei Bedingungen stehen in einem **eigenen Commit**, der nur
  die Slice-Datei anfasst und vor jeder Messung liegt. Der Review hat die
  Reihenfolge unabhängig geprüft und die Bedingungen als deckungsgleich mit
  denen bestätigt, gegen die entschieden wurde. **Eine Schwelle in einem eigenen
  Commit ist eine Tatsache; in einer Botschaft wäre sie eine Behauptung.**

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls kein Muster die Messung besteht — dann ist die Ausweisung das Ergebnis.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Konfigurations-Profil (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-011`](../observations.md) ist zentral: die Schwelle für „vertretbar" gehört **vor** die Messung.

Slice-ID: slice-144. Betroffene IDs: [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix). Module: `matrix`, Konfigurations-Profil.
Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Messung an bestehender Modul-Mechanik.

## 9. Closure-Notiz (nach `done/`)

Geliefert: die vierte der fünf Kategorien aus
[`AGENTS.md`](../../../../AGENTS.md) §3.4 ist gedeckt — eine `matrix`-Klasse
`commit-hash` **ohne Zieldateien**, deren Gegenstand eine Zeichenkette ist.

**Null Befunde heißt nicht, dass es greift.** Alle vier Kandidaten meldeten im
Bestand null. Erst die Positivkontrolle trennt sie: das gewählte Muster meldet
Kurz- **und** Langform, das 40-Zeichen-Muster nur die Langform. Ohne diese Probe
wäre *„0 Befunde"* von *„die Klasse tut nichts"* nicht zu unterscheiden gewesen
— und das ist die Gestalt, in der ein stiller Grün-Pfad entsteht.

**Die Schwelle stand in einem eigenen Commit, vor jeder Messung.** Das ist der
Teil dieses Slice, der am billigsten war und am meisten getragen hat: Was in
einer Commit-Botschaft eine Behauptung wäre, ist als eigener Commit eine
Tatsache in der Historie. Der Review hat die Reihenfolge unabhängig geprüft.

**Die Risiko-Klasse war trotzdem zu eng benannt.** Ich schrieb *„reine
a–f-Wörter ab sieben Zeichen"*. Das Muster trifft **jede** 7- bis 40-stellige
Zeichenkette aus `[0-9a-f]` — auch rein dezimale: eine CI-Lauf-Nummer, ein
Kompaktdatum, ein Hex-Farbwert. Alle drei nachgestellt, alle drei gemeldet.
Tragbar bleibt die Klasse, aber aus einem anderen Grund als dem, den ich
aufgeschrieben hatte.

**Zwei Zahlen standen unter der falschen Überschrift.** Der Kommentar sagte
*„Gemessen mit dem Produkt, nicht per grep"* — und nannte darunter zwei Zahlen
aus einem eigenen Rohtext-Skript. Sie sind aus der Konfiguration entfernt; was
dort steht, ist die Zusage, nicht die Messgeschichte. Das ist
[`BEO-009`](../observations.md) Richtung (b), an einem Tag zum vierten Mal.

**Ein Kommentar trug, was §3.7 verbietet:** eine Slice-Nummer und zwei Zeilen
Lauf-Historie. Der Reviewer hat es als HIGH gemeldet, und zu Recht — die
Geschwister-Kommentare lösen ihre Herkunft über
[`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
und die Adaptionen auf. Dieser tut es jetzt auch.

**§3.4 behauptete ein Gate, das es so nicht gibt.** Der Satz sagte, `matrix`
halte ADRs, Slices, Wellen und Commit-Hashes als **Token**. Die Klasse `adr`
trägt gar kein `token`: ein bares `ADR-<NNNN>` im Spec-Körper erzeugt keinen
Befund — ADRs sind über die **Link**-Prüfung gedeckt. Nachgestellt und im Text
getrennt.

**Drei Grenzen fehlten.** Ein Hash im **Linktext oder in der Link-URL** bleibt
stumm (die bare URL nicht); der **Provenance-Marker** auf derselben rohen Zeile
nimmt auch diese Klasse aus, obwohl §3.4 für die Straten keine solche Ausnahme
kennt; und die Obergrenze 40 hält nur den **vollen** 64-Steller heraus, nicht
die Hausschreibweise mit gekürztem Digest.

**Und eine Vertrags-Lücke, halb geschlossen.** Meine Begründung *„matrix
verlangt nur `name`"* stammte aus der Implementierung. Die Schema-Zeile führte
`paths` ohne Default — wie ein Pflichtfeld —, und
[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
beschreibt Klassen *„über Pfad-Muster"*. Die **technische** Spezifikation ist
nachgezogen (`paths` optional, eine Klasse ohne Pfade ist reines Token-Ziel,
ohne beides inert), mit Historie-Zeile und ohne Verhaltens-Delta. Die
**vertragliche** Beschreibung ist es nicht — das ist keine Implementierungs-,
sondern eine Vertrags-Frage und liegt als
[slice-153](../done/slice-153-lastenheft-token-ziel-klasse.md) vor, mit Kennung
statt als Absatz.

**Register:** [`BEO-009`](../observations.md) auf Zähler **6**.
