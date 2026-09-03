# Slice slice-192: Der Übergangs-Wächter bindet jetzt Review- und Register-Deckung

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-86](../welle-86-closure-uebergang-durchsetzen.md) — fünfter
Slice, nachträglich durch den Trigger-Audit (Modul 6 Schritt 2) geschnitten.

**Bezug:** [ADR-0082](../../adr/0082-uebergangswaechter-reviews-observations.md)
(die Entscheidung), [ADR-0081](../../adr/0081-reviews-modul.md) §Re-Evaluierungs-Trigger
(der eingelöste Trigger), slice-175 <!-- d-check:status-provenance --> (der
Bindepunkt, den dieser Slice erweitert).

**Berührte Spec-Stellen:** — (Verdrahtung bestehender Fähigkeiten an einen
bestehenden Bindepunkt; keine neue `DC-FA`-Anforderung, kein neuer
Grund-Code).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-03.

---

## 1. Ziel

Der Übergangs-Wächter aus slice-175 <!-- d-check:status-provenance -->
prüft ab jetzt auch Beobachtungs-Register-Deckung und
Review-Report-Deckung — nicht nur DoD-Häkchen. Damit sind erstmals **alle
vier** Vorbedingungen von welle-86 §1 real beim lokalen `mv`-Commit
durchgesetzt, nicht nur zwei.

## 2. Vorgehen

1. **Ursache lesen, nicht vermuten:** welle-86s eigene Trigger-Prüfung (vier
   reale Proben) zeigte, dass `.d-check.closure.yml` keinen
   `planning.observations`-Block trägt und `reviews` in keinem erzwungenen
   Lauf aktiv ist.
2. **`.d-check.closure.yml` um `planning.observations` ergänzen**
   (Register + Verzeichnis, identisch zum Hauptprofil).
3. **`.d-check.closure.yml` um einen `reviews`-Block ergänzen**
   (`done-dir`, `reviews-dir`, dieselbe Bestands-Ausnahme wie im
   Hauptprofil).
4. **`make verify-closure-notes` aktiviert zusätzlich `--enable reviews`.**
   Kein neues Modul, keine neue Prüf-Logik.
5. **Beide zuvor bestandenen Proben erneut fahren** — jetzt erwartet:
   abgewiesen. Ergebnis in §9.
6. **Zwei Review-Report-Lücken der eigenen Session schließen:** die
   Verschärfung deckte auf, dass slice-173 und slice-175 selbst — beide mit
   `[x] … unabhängiger Review` in ihrer DoD — keinen Report unter
   `docs/reviews/` haben. Die Reviews **fanden statt** (unabhängige
   Sub-Agenten-Läufe, Befunde eingearbeitet, in den jeweiligen
   Closure-Notizen dokumentiert) — nur nie als Artefakt festgehalten. Zwei
   Reports werden jetzt aus den tatsächlichen Befunden dieser Läufe
   nachgezogen (Modul-10-Form), nicht rückwirkend erfunden.
7. **[ADR-0082](../../adr/0082-uebergangswaechter-reviews-observations.md)**
   dokumentiert die Entscheidung; **[ADR-0081](../../adr/0081-reviews-modul.md)**
   bekommt einen `## Geschichte`-Eintrag zum eingelösten Trigger (immutabel,
   kein Kern-Edit).

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Aufnahme von `reviews` in die unconditional `gates`-Zusammensetzung**
  — [ADR-0082](../../adr/0082-uebergangswaechter-reviews-observations.md)
  Entscheidung 6: der Bindepunkt bleibt gezielt auf `done/`-Transitionen.
- **Kein Nachrüsten des übrigen `done/`-Bestands** mit fehlenden
  Review-Reports — nur die zwei, die diese Session selbst erzeugt hat,
  werden nachgezogen. Der Rest bleibt in der Bestands-Ausnahme aus
  [ADR-0081](../../adr/0081-reviews-modul.md).
- **Kein Folge-ADR mit `supersedes` auf [ADR-0081](../../adr/0081-reviews-modul.md)**
  — die Entscheidung dort war richtig und benannte selbst die Bedingung,
  unter der sie sich ändert.

## 4. Definition of Done

- [x] `.d-check.closure.yml` trägt `planning.observations` und einen
      `reviews`-Block, beide deckungsgleich mit dem Hauptprofil.
- [x] `make verify-closure-notes` aktiviert `--enable reviews`; der reale
      Bestand bleibt bei null Befunden (Bestands-Ausnahme greift
      unverändert in beiden Profilen).
- [x] **Zwei Proben real erneut gefahren:** eine zitierte, nicht
      registrierte `BEO-999` wird jetzt lokal abgewiesen
      (`observation-unregistered`); ein Test-Kandidat mit Review-Zusage
      ohne Report wird jetzt lokal abgewiesen (`review-missing`). Beide
      Ausgaben in §9.
- [x] Drei Review-Reports nachgezogen (`docs/reviews/2026-09-03-slice-173-…`,
      `…-slice-175-…`, `…-slice-192-…` — der dritte, weil der eigene
      `mv`-Versuch vom eigenen Wächter abgewiesen wurde, siehe §9), aus den
      tatsächlichen Befunden der bereits gelaufenen unabhängigen Reviews.
- [x] [ADR-0082](../../adr/0082-uebergangswaechter-reviews-observations.md)
      trägt Kontext, Entscheidung, mindestens drei verglichene
      Alternativen, Konsequenzen, Fitness Function,
      Re-Evaluierungs-Trigger.
- [x] [ADR-0081](../../adr/0081-reviews-modul.md) trägt einen
      `## Geschichte`-Eintrag zum eingelösten Re-Evaluierungs-Trigger.
- [x] welle-86s eigener Plan trägt den Nachtrag (§3, §4-Tabelle).
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register
      fortgeschrieben (oder Fehlanzeige begründet); jedes Risiko aus §5 mit
      Ausgang; die drei Paarungen laufen bei der Welle-Closure.

## 5. Abnahme-Punkte / Risiken

- **Die Config-Duplikation zwischen Haupt- und Closure-Profil wächst um
  zwei Blöcke** — wächst die Bestands-Ausnahme künftig, müssen beide
  Stellen nachgezogen werden, und kein Sensor hält das. — **Ausgang:**
  *weiter offen* → kein neuer `BEO`-Eintrag (die Grenze ist bereits im
  Datei-Kopf von `.d-check.closure.yml` benannt, seit dessen Anlage;
  dieselbe akzeptierte Kosten-Klasse, kein neuer Fall).
- **Die zwei nachgezogenen Review-Reports könnten als rückwirkend
  konstruiert wirken, nicht als echte Läufe.** — **Ausgang:** *entfallen* —
  beide Reports geben die tatsächlichen Befunde der bereits gelaufenen
  Sub-Agenten-Reviews wieder (dieselben Kategorien, Dateien, Zeilen wie in
  den jeweiligen Closure-Notizen bereits dokumentiert); kein Befund ist neu
  erfunden.
- **Eine künftige Erweiterung der `reviews`-Bestands-Ausnahme im
  Hauptprofil, die im Closure-Profil vergessen wird, bliebe unsichtbar, bis
  ein echter Übergang sie träfe.** — **Ausgang:** *weiter offen* →
  benannt als Re-Evaluierungs-Trigger in
  [ADR-0082](../../adr/0082-uebergangswaechter-reviews-observations.md)
  („ein dritter Bindepunkt, an dem dieselbe Lücken-Klasse auftritt").

## 6. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — WIP-Limit frei.

**Rückführungen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls die
  Config-Verdrahtung mehr als die zwei benannten Blöcke berühren müsste —
  trat nicht ein.
- `in-progress` → `open` (blockiert): keine Bedingung erkannt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.d-check.closure.yml` und `Makefile` (derselbe
  Bindepunkt wie slice-175) fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area). Die Regel, die diesen Schritt
  vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-09-03, höchste
  Kennung `BEO-027`): keine Beobachtung trifft speziell diesen
  Config-Bindepunkt über das hinaus, was jeder Slice ohnehin berührt — keine
  Treffer. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): unverändert seit
  slice-175 (`upstream-drift.yml` zwei planmäßige `VERALTET`-Meldungen,
  `image-scan.yml` grün). Keiner der beiden Funde berührt diesen Slice.
  **Dieser Block trägt bewusst keine `cite`-Direktive** — sein Ziel ist
  eine Repo-Adaption, kein Baseline-Abschnitt
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-192. Betroffene IDs: keine neue `DC-FA`-Anforderung.
Module: `planning`, `reviews` (unverändert, nur neu gebunden). Gates:
`make test`, `make verify-closure-notes`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die berührte Sub-Area fällt unter den
Default: Doc führt, Code folgt. Es entsteht keine neue Fähigkeit; die
Konventions-Dichte ist hoch (dieselbe Zwei-Profil-Duplikation, die
`.d-check.closure.yml`s eigener Kopf bereits für `planning.heading`/`marker`
begründet).

## 9. Closure-Notiz (nach `done/`)

- **Was hat funktioniert:** Die vier realen Proben aus welle-86s eigener
  Trigger-Prüfung waren der richtige Mechanismus, um die Lücke zu finden —
  eine Behauptung „vier Vorbedingungen sind jetzt durchgesetzt" hätte sie
  nicht gezeigt. Nach der Umsetzung lieferten beide erneut gefahrenen Proben
  die erwarteten Grund-Codes: die Probe mit der zitierten, nicht
  registrierten Test-Kennung meldete `observation-unregistered`, die Probe
  mit fehlendem Report meldete `review-missing` — jeweils mit
  `make: *** [Makefile:348: verify-closure-notes] Fehler 1`, kein neuer
  Commit entstand (`git log` bestätigt).
- **Was ging anders als geplant:** Der erste Anlauf von Probe C schlug aus
  dem falschen Grund fehl — mein eigenes Test-Fixture schrieb
  „unabhaengiger" (ASCII-Transliteration) statt „unabhängiger" und wurde
  deshalb gar nicht als Kandidat erkannt; der Commit ging fälschlich durch.
  Sofort per `git reset --soft` zurückgenommen und mit korrektem Umlaut
  wiederholt. Zweitens: die neue ADR verletzte beim ersten Schreiben
  `matrix-forbidden` achtfach (Bare-Token-Referenzen auf `welle-86` und
  `slice-175`) — behoben mit `<!-- d-check:status-provenance -->`, diesmal
  auch für die `welle`-Klasse als Bare-Token (anders als die reine
  Link-Form, die keinen Marker-Ausweg hat). Drittens, und am lehrreichsten:
  der erste `git mv`-Versuch dieses Slice nach `done/` wurde vom **eigenen,
  gerade geschärften** Wächter abgewiesen — `slice-192` selbst trägt eine
  `unabhängiger Review`-DoD-Zusage ohne persistierten Report unter
  `docs/reviews/` (derselbe Fund wie bei `slice-173`/`slice-175`, nur diesmal
  von der eigenen Neuerung sofort gefangen statt erst später gemessen).
  Behoben durch denselben Nachzug: ein dritter Report
  ([2026-09-03-slice-192-…](../../../reviews/2026-09-03-slice-192-uebergangswaechter-nachliefern-review.md))
  aus den tatsächlichen Befunden des bereits gelaufenen Reviews (ein
  MEDIUM: Doku-Rückstand, in Commit `d606ec9` behoben).
- **Steering-Loop-Eintrag:** keiner verkörpert — alle drei Beobachtungen
  oben sind Einzelfälle dieser Session.
- **Beobachtungs-Register (`../observations.md`):** keine Beobachtung
  angefallen, die eine neue Kennung oder einen Zähler-Schritt rechtfertigt.
- **Folge-Slices:** keine.
- **Risiken aus §5:** siehe §5, je Zeile ein Ausgang.
- **Drei Paarungen:** nicht hier geprüft — Repo mit Wellen-Betrieb, prüft
  die Closure von welle-86.
