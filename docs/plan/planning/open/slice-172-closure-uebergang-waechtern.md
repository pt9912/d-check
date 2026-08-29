# Slice slice-172: Der Closure-Übergang wird gewächtert

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-86](../welle-86-closure-uebergang-durchsetzen.md) — der
Closure-Trigger der Welle beobachtet **mehr**, als diese DoD belegt: dass die
Vorbedingungen am **Übergang** greifen, nicht nur im Zustand danach.

**Bezug:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das Werkzeug: `forbid-pattern` über einen Abschnitt);
[ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md) und
[`MR-049`](../../../../harness/conventions.md#mr-049) (die Präzedenz: die
urteilsfreie Hälfte einer Closure-Regel als `structure`-Regel, mit
Bestands-Ausnahme);
[ADR-0069](../../adr/0069-zellenlaenge-als-strukturbedingung.md) (dieselbe
Bauform ein Slice später).

**Berührte Spec-Stellen:** — (Konfigurationsregel über vorhandene Fähigkeiten;
keine Produkt-Anforderung berührt).

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**Der Übergang nach `done/` trägt eine Bedingung, die kein Sensor hält.**

Das Regelwerk ist an dieser Stelle unmissverständlich:

<!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:33-34 -->
> DoD-Häkchen und Closure-Notiz
> sind die Bedingung dafür, dass die Datei überhaupt nach `done/` darf.

**Gehalten wird davon heute eine Hälfte.** `verify-closure-notes` prüft die
Closure-Notiz (Abschnitt vorhanden, Substanz, keine Floskel) und die
Risiko-Ausgänge. Die **DoD-Häkchen** prüft nichts: `.d-check.closure.yml`
verlangt vom `## N. Definition of Done`-Abschnitt nur, dass er **nicht leer**
ist — nicht, dass seine Haken gesetzt sind.

**Der Bestand ist gemessen, nicht geschätzt:** von 169 `done/`-Slices tragen
**37** mindestens einen offenen Haken. Bei **20** davon ist es die
Review-Zeile — und **87 von 95** Slices mit Review-Zusage haben tatsächlich
einen Report in `docs/reviews/`, acht nicht. Die Konvention wird also gelebt;
was fehlt, ist der Wächter.

**Der Anlassfall ist frisch und liegt in diesem Repo:** slice-168, -169 und
-170 gingen mit offenem Review-Haken nach `done/`. Kein Gate hat es gemeldet,
weil keines danach sieht. Das ist die Bauform, die dieses Repo als
[`BEO-013`](../observations.md) führt — nur eine Stufe früher: nicht ein
Wächter, der nichts mehr fängt, sondern eine Regel, die nie einen hatte.

**Warum das der billigste Schnitt ist:** Der offene Haken ist **urteilsfrei**
erkennbar — `- [ ]` steht da oder nicht. Und er ist ein **Stellvertreter** für
mehr: Wo der Review-Haken offen ist, hat der Review nicht stattgefunden oder
niemand hat ihn quittiert; beides gehört gesehen. Was der Sensor **nicht**
sagt, ist, ob der Review taugte — dieselbe Arbeitsteilung wie beim
Beobachtungs-Register: Mensch urteilt, Maschine prüft Deckung.

## 2. Vorgehen

1. **Eine `structure`-Regel im Closure-Profil**, nach dem Muster von
   [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md):
   `forbid-pattern` auf das Task-Item-Muster über den
   `## N. Definition of Done`-Abschnitt von `done/slice-*`.
2. **Bestands-Ausnahme mit fester Ziffernzahl**, exakt wie
   [`MR-049`](../../../../harness/conventions.md#mr-049) sie für die
   Drei-Ausgänge-Regel führt — und aus demselben Grund, den der Auftraggeber
   benannt hat: **das Regelwerk wurde mehrfach gehoben, und die Dokumente sind
   nicht durchgängig nachgezogen.** Ein Befund auf einem Slice, der nach
   damaliger Form korrekt war, wäre kein Befund, sondern Lärm.
3. **Vor dem Scharfschalten rot messen**, nicht behaupten: die Regel gegen den
   heutigen Bestand fahren, die Trefferzahl mit den 37 abgleichen, und mit
   Ausnahme auf null bringen. Weicht die Zahl ab, ist das Muster falsch — nicht
   der Bestand.
4. **Der Grund-Code trägt die Reparatur, nicht den Zustand.** `section-forbidden`
   sagt „hier steht etwas Verbotenes"; die Meldung muss sagen, *was zu tun ist*:
   den Haken setzen oder den Slice zurückführen.
5. `make gates` und `make fullbuild`; **Review** und **Verifikation** als
   getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Review-Report-Deckung.** Die Prüfung „jeder `done/`-Slice mit
  Review-Zusage hat einen Report in `docs/reviews/`" ist gemessen tragfähig
  (87/95), aber sie ist eine **Deckung zwischen zwei Mengen** — das kann keine
  `structure`-Regel. Es wäre eine neue Fähigkeit und damit eine neue
  Anforderung: eigener Slice.
- **Keine `BEO`-Deckungsprüfung.** `modul-06` benennt die maschinelle Hälfte
  selbst (zitierte `BEO-<NNN>` hat eine Registerzeile; jede Zeile trägt einen
  Beleg). Dieselbe Klasse wie oben, eigener Schnitt.
- **Keine Reihenfolge-Prüfung über git** („Closure-Commit nach
  Review-Report-Commit"). Der VCS-Port existiert, aber das ist eine
  Historien-Aussage und gehört zur `vcs`-Familie, nicht hierher.
- **Kein Nachrüsten der 37 Bestands-Slices.** Sie sind Belege ihrer Zeit; ein
  nachträglich gesetzter Haken behauptet einen Review, den es nicht gab.
- **Keine Aussage über Review-*Qualität*.** Der Sensor prüft ein Häkchen.

## 4. Definition of Done

- [ ] Eine `structure`-Regel im Closure-Profil meldet einen offenen DoD-Haken
      in `done/slice-*` mit eigener, reparatur-benennender Meldung.
- [ ] **Retro gemessen:** die Regel meldet ohne Ausnahme die erwartete
      Bestandszahl (37 Dateien) und mit Bestands-Ausnahme **null**; beide
      Ausgaben stehen in der Commit-Botschaft.
- [ ] Die Bestands-Ausnahme ist als Adaption geführt (oder an
      [`MR-049`](../../../../harness/conventions.md#mr-049) angehängt) und nennt
      den Grund: gehobenes Regelwerk, nicht nachgezogene Dokumente.
- [ ] Die drei Folge-Kandidaten aus §3 sind als solche benannt — nicht
      stillschweigend weggelassen.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Ein Haken ist eine Selbstauskunft, und der Sensor prüft nur sie.** Wer den
  Haken setzt, ohne dass ein Review stattfand, passiert das Gate. Der Sensor
  verschiebt die Lücke von „unsichtbar" nach „behauptet" — das ist besser, aber
  es ist keine Prüfung des Reviews. — **Ausgang:** *(bei Closure)*
- **Die Bestands-Ausnahme ist eine Grandfathering-Klausel und altert.** Sie
  nimmt eine feste Nummernspanne heraus; jeder neue Slice fällt unter die
  Regel. Wächst die Ausnahme je wieder, ist das der Befund
  ([`BEO-013`](../observations.md)). — **Ausgang:** *(bei Closure)*
- **Der Sensor macht den inneren Loop rot, bevor jemand schließen will.**
  `verify-closure-notes` läuft in `make fullbuild`, nicht in `gates` — die neue
  Regel gehört an dieselbe Stelle, sonst meldet sie beim Arbeiten an einem
  laufenden Slice. — **Ausgang:** *(bei Closure)*
- **Die Regel wird aus einem Anlass gezogen** ([`BEO-011`](../observations.md)):
  drei Slices einer Sitzung. Gegen den Bestand gemessen ist sie es nicht — 37
  Treffer sind ein Bestand —, aber die *Dringlichkeit* stammt aus dem Anlass. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — heute hält
[slice-171](../in-progress/slice-171-vorpruefungen-belegen.md) den Slot.

**Rückführungen:** `in-progress` → `open`, falls die Bestandsmessung eine
andere Zahl als 37 liefert oder die Treffer über die Nummernspanne streuen —
dann trägt eine feste Ziffern-Ausnahme nicht, und der Schnitt ist ein anderer.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `docs/plan/planning/` (Planning-Form) und die
  Gate-Konfiguration `.d-check.closure.yml`. Beide fallen unter den Default
  `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration). Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29, höchste Kennung
  `BEO-021`): [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr
  fängt: hier eine Stufe früher, eine Regel **ohne** Wächter, und die
  Bestands-Ausnahme ist der künftige Kandidat für dieselbe Klasse;
  [`BEO-011`](../observations.md) — Regel aus dem Anlass: die Dringlichkeit
  stammt aus drei Slices, die Regel selbst aus 37 gemessenen Treffern;
  [`BEO-015`](../observations.md) — ein offener Punkt bekommt bei der Closure
  einen Ausgang, den es nicht gibt: dieselbe Familie urteilsfreier
  Closure-Prüfungen. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:224-225 -->
  > **Keine Treffer sind ebenfalls eine
  > Antwort** und werden notiert.

- **Nachtlauf-Stand lesen:** entfällt in `open/` — der Block entsteht
  **spätestens bei der Beanspruchung** (`open→in-progress`), weil ein zum
  Planungszeitpunkt gelesener Stand bis dahin veraltet wäre.

Slice-ID: slice-172. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`MR-049`](../../../../harness/conventions.md#mr-049). Module: `structure`,
`planning`. Gates: `make gates`, `make verify-closure-notes`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Eine Konfigurationsregel über eine vorhandene
Produkt-Fähigkeit; kein Produkt-Code, kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)
