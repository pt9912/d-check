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

**Verantwortlich:** pt9912.

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
4. **Die Meldung trägt die Reparatur, nicht nur den Zustand.** `section-forbidden`
   sagt „hier steht etwas Verbotenes"; die Meldung muss sagen, *was zu tun ist*:
   den Haken setzen oder den Slice zurückführen. **Das Schema kann das heute
   nicht** — `structure` kennt keine regel-eigene Meldung, und der
   `--doctor`-Klartext hängt am Grund-Code, nicht an der Regel. Dieser Punkt
   ist der Grund für die Rückführung (§6); die Fähigkeit liefert
   [slice-177](../done/slice-177-structure-hint.md), und diese Regel
   wird ihr erster Konsument.
5. `make gates` und `make fullbuild`; **Review** und **Verifikation** als
   getrennte Läufe; Closure.


**Vorarbeit, am 2026-08-29 gemessen und hier festgehalten, damit die nächste
Beanspruchung sie nicht wiederholt:**

- **Das Überschriften-Muster trägt.** `^#{2,3} [0-9]+\. Definition of Done`
  deckt alle sieben im Bestand vorkommenden Formen (`## 2.`/`## 3.`/`## 4.`/
  `## 5.`, dazu die Zusätze `(vorläufig)` und `(R1 eingearbeitet)`) — **null**
  `section-missing` über den `done/`-Bestand. Das ist der Unterschied zu
  [`MR-049`](../../../../harness/conventions.md#mr-049), dessen Ausnahme 107
  `section-missing` abfangen musste.
- **Die Bestandszahl stimmt mit der Planung überein:** ohne Ausnahme **37**
  `section-forbidden`, mit ihr **null** — beides im echten Closure-Profil
  gefahren, nicht in einer Probe-Konfiguration.
- **Bei der zweiten Beanspruchung nachgeprüft, nicht wiederholt:** der Bestand
  ist von 170 auf **171** `done/`-Slices gewachsen (slice-177 kam hinzu), die
  Trefferzahl bleibt **37**, die Nummern bleiben dieselben, und es gibt weiter
  **kein** `section-missing`. Der seither geschlossene Slice trägt
  seine Haken gesetzt — er hätte die Zahl sonst erhöht.
- **Die Treffer streuen nicht beliebig:** 025–104 (34 Stück), dann 160, 168,
  169, 170. Eine feste Ziffern-Ausnahme trägt also, sie braucht drei Muster
  statt zwei: `slice-0??-*`, `slice-1[0-6]?-*`, `slice-170-*`. Die Grenze liegt
  damit bei [slice-171](../done/slice-171-vorpruefungen-belegen.md) — dem
  ersten Slice, der unter der korrigierten Praxis geschlossen wurde.
- **Die Positiv-Probe läuft:** ein geöffneter Haken in slice-171, also
  außerhalb der Ausnahme, wird mit Datei und Zeile gefangen (Exit 1); zurück
  gesetzt wieder Exit 0.
- **Der Anlassfall bleibt ausgenommen, und das ist entschieden.** slice-168,
  -169 und -170 tragen genau den offenen Review-Haken, der diese Welle
  ausgelöst hat. Sie nachträglich zu haken behauptete einen Review, den es
  nicht gab (§3); ihr Befund steht im Register und in der Closure-Notiz von
  slice-171, nicht in einem stillen Haken.


**Beim ersten Bau gefunden und hier aufgehoben — zwei Blindstellen, beide
gemessen.** Sie sind der Grund für die zweite Rückführung (§6):

- **Ein einzelner zusätzlicher Backtick im DoD-Listenblock schaltet den Wächter
  ganz ab.** Die absatzweite Inline-Code-Paarung
  ([ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md))
  schluckt den offenen Haken. An derselben Datei gefahren: offener Haken ⇒ **1**
  Befund, Exit 1; derselbe Haken plus ein Backtick weiter oben ⇒ **0** Befunde,
  Exit 0. **Die Exposition ist real:** `slice-061` und `slice-076` tragen heute
  ungerade Backtick-Zahlen in ihrem DoD-Abschnitt (25 bzw. 45). Der Preis war in
  [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md) für einen **Platzhalter** ausgewiesen; hier zahlt ihn eine
  **Vorbedingung** des Übergangs.
- **`* [ ]` und `+ [ ]` laufen still durch.** Beide sind gültige
  CommonMark-Task-Items; gemessen je **0** Befunde. Das Muster muss alle drei
  Bullet-Formen decken (`[-*+] \[ \]`), nicht nur die im Repo übliche.

**Die Substanz der zurückgenommenen Adaption, damit sie mit der Regel
zurückkommt:**

- **Die Grenze der Ausnahme ist ein Datum, kein Aufräum-Rest:** sie endet bei
  [slice-171](../done/slice-171-vorpruefungen-belegen.md), dem ersten Slice nach
  der korrigierten Praxis. Damit ist sie überprüfbar statt verhandelbar.
- **Der Anlassfall liegt bewusst *in* der Ausnahme.** slice-168, -169 und -170
  tragen genau den offenen Review-Haken, der
  [welle-86](../welle-86-closure-uebergang-durchsetzen.md) ausgelöst hat — der
  Wächter wird sie nie melden. Ein nachträglich gesetzter Haken behauptete einen
  Review, den es nicht gab. Wer die Ausnahme liest, ohne das zu wissen, hält sie
  für Bequemlichkeit; deshalb gehört es in die Adaption.
- **Zweite Instanz der Form aus
  [`MR-049`](../../../../harness/conventions.md#mr-049)** — damit ist sie die
  Haus-Antwort auf „neue Closure-Bedingung über gewachsenem Bestand". Anders als
  dort fängt die Ausnahme hier **kein** `section-missing`.
- **Grenze:** ein Haken ist eine Selbstauskunft; die Lücke wandert von
  *unsichtbar* nach *behauptet*. Und der Befund zeigt auf die
  **Abschnitts-Überschrift**, nicht auf die Zeile des Hakens — bei neun offenen
  Haken in `slice-045` gibt es einen Befund, und der Leser sucht selbst.

**Eine Design-Frage ist offen und muss vor dem nächsten Bau beantwortet
werden.** Der Kopfkommentar von [`.d-check.closure.yml`](../../../../.d-check.closure.yml)
führt genau diese Regel als **bewusst verworfen**, mit einem sachlichen Grund:
*ein abgeschlossener Slice darf eine offene Box tragen, wenn die **Welle** sie
einlöst*. Der Fall ist belegt — `slice-094` und `slice-104` tragen je einen
offenen Release-Haken mit dem ausdrücklichen Vermerk „Wellen-Trigger, nicht
Slice-Trigger". Für den Altbestand löst die Ausnahme das mit; für **jeden
künftigen** Slice wäre er rot. Der naheliegende Ersatzweg: was die Welle
einlöst, gehört in ihren **Closure-Trigger** (welle-86 §3 führt ihn bereits so)
und nicht in die DoD eines Slice, der es nicht einlösen kann. Das ist ein
Vorschlag, kein Entscheid — er kehrt eine dokumentierte Ablehnung um, und der
Kopfkommentar gehört mit ihm ersetzt.

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
- [ ] Die neue Bedingung ist **zugestellt**, bevor sie blockiert:
      [`AGENTS.md`](../../../../AGENTS.md) §5 und die Sensors-Tabelle in
      [`harness/README.md`](../../../../harness/README.md) nennen sie
      ([`BEO-022`](../observations.md)).
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

**Start** (`open` → `in-progress`): WIP-Limit frei, beansprucht am 2026-08-29 —
[slice-171](../done/slice-171-vorpruefungen-belegen.md) ist geschlossen.

**Rückführung eingetreten am 2026-08-29** (`in-progress` → `open`): nicht aus
dem unten genannten Grund — die Bestandsmessung stimmt (§2). Der Sensor war
fertig und beidseitig gemessen, als §2 Punkt 4 auf eine Produkt-Grenze traf:
eine regel-eigene Meldung gibt es nicht. Die Alternative wäre gewesen, den
Punkt stillschweigend fallenzulassen; stattdessen wartet dieser Slice auf
[slice-177](../done/slice-177-structure-hint.md).

**Zweite Beanspruchung am 2026-08-29** (`open` → `in-progress`): die
Vorbedingung ist geliefert — `structure[].hint` existiert seit
[slice-177](../done/slice-177-structure-hint.md), und der Nachtlauf-Stand ist
neu gelesen (§7). Die gemessene Vorarbeit in §2 gilt unverändert; sie wurde
nicht wiederholt, sondern gegen den heutigen Bestand nur nachgeprüft.

**Zweite Rückführung am 2026-08-29** (`in-progress` → `open`): erneut nicht aus
dem unten genannten Grund. Der Wächter war gebaut, beidseitig gemessen und
zugestellt, als Review und Verifikation zwei **stille** Blindstellen fanden
(§2): ein Backtick schaltet ihn ganz ab, zwei Bullet-Formen laufen durch.
Entscheid des Auftraggebers: erst die Produkt-Frage lösen — eine Bedingung, die
den **rohen** Abschnittstext liest —, dann diesen Slice damit schließen. Die
Regel, ihre Adaption und die Zustellung sind zurückgenommen; die Messungen
bleiben in §2.

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
  `BEO-023`): [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr
  fängt: hier eine Stufe früher, eine Regel **ohne** Wächter, und die
  Bestands-Ausnahme ist der künftige Kandidat für dieselbe Klasse;
  [`BEO-011`](../observations.md) — Regel aus dem Anlass: die Dringlichkeit
  stammt aus drei Slices, die Regel selbst aus 37 gemessenen Treffern;
  [`BEO-015`](../observations.md) — ein offener Punkt bekommt bei der Closure
  einen Ausgang, den es nicht gibt: dieselbe Familie urteilsfreier
  Closure-Prüfungen; [`BEO-022`](../observations.md) — eine Regel tritt in Kraft,
  bevor ihre Zustellung existiert: die neue Bedingung am Übergang gehört in
  `AGENTS.md` §5 und in die Sensors-Tabelle, sonst blockiert sie einen Autor,
  der nicht weiß, warum; [`BEO-023`](../observations.md) — ein Wächter, der nie
  fangen konnte: dieser Slice **ist** ein Wächter, und seine DoD verlangt die
  Umkehr-Probe (die Regel gegen den Bestand rot, mit Ausnahme grün) statt eines
  grünen Laufs als Beleg. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-29T07:34:35Z,
  `image-scan.yml` 2026-08-29T10:07:43Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-172. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`MR-049`](../../../../harness/conventions.md#mr-049). Module: `structure`,
`planning`. Gates: `make gates`, `make verify-closure-notes`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Eine Konfigurationsregel über eine vorhandene
Produkt-Fähigkeit; kein Produkt-Code, kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)
