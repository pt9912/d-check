# Slice slice-172: Der Closure-Übergang wird gewächtert

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-86](../welle-86-closure-uebergang-durchsetzen.md) — der
Closure-Trigger der Welle beobachtet **mehr**, als diese DoD belegt: dass die
Vorbedingungen am **Übergang** greifen, nicht nur im Zustand danach.

**Bezug:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das Werkzeug: `max-open-tasks` über einen Abschnitt — seit
[slice-178](../done/slice-178-offene-tasks-roh.md));
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
   [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md), aber
   mit **`max-open-tasks: 0`** statt `forbid-pattern` über den
   `## N. Definition of Done`-Abschnitt von `done/slice-*`. **Das ist die
   Änderung gegenüber den beiden vorigen Anläufen**, und sie ist der Grund, aus
   dem dieser Slice zweimal zurückging: `forbid-pattern` liest den bereinigten
   Text und fiel an einem einzelnen Backtick auf null Befunde, und es deckte nur
   die Bullet-Form, die sein Autor aufschrieb. Beides ist mit
   [slice-178](../done/slice-178-offene-tasks-roh.md) erledigt —
   [`ADR-0074`](../../adr/0074-offene-tasks-auf-rohen-zeilen.md) zählt roh und
   über die Modul-Lexik, und der Befund steht auf **der Zeile des Hakens** statt
   auf der Abschnitts-Überschrift.
2. **Bestands-Ausnahme mit fester Ziffernzahl**, exakt wie
   [`MR-049`](../../../../harness/conventions.md#mr-049) sie für die
   Drei-Ausgänge-Regel führt — und aus demselben Grund, den der Auftraggeber
   benannt hat: **das Regelwerk wurde mehrfach gehoben, und die Dokumente sind
   nicht durchgängig nachgezogen.** Ein Befund auf einem Slice, der nach
   damaliger Form korrekt war, wäre kein Befund, sondern Lärm.
3. **Vor dem Scharfschalten rot messen**, nicht behaupten: die Regel gegen den
   heutigen Bestand fahren und mit Ausnahme auf null bringen. Weicht die Zahl
   ab, ist das Muster falsch — nicht der Bestand.
4. **Die Meldung trägt die Reparatur, nicht nur den Zustand.**
   `section-tasks-open` sagt „hier steht ein offener Haken"; die Meldung muss
   sagen, *was zu tun ist*: den Haken setzen oder den Slice zurückführen. Das
   leistet `structure[].hint` seit
   [slice-177](../done/slice-177-structure-hint.md) — diese Regel wird sein
   erster Konsument, und dass die Fähigkeit fehlte, war der Grund für die
   **erste** Rückführung (§6).
5. **`spans` gehört ins Profil, und das ist keine Kür.**
   [slice-178](../done/slice-178-offene-tasks-roh.md) hat den vergessenen
   Schluss-Fence als **deklarierte, nicht behobene** Grenze der Bedingung
   ausgewiesen: er blendet alles Folgende aus, und die rohe Lesung behebt das
   nicht. Gefangen wird er von `fence-unclosed`. Der Closure-Bindepunkt fährt
   `spans` seit [slice-180](../done/slice-180-closure-profil-spans.md) —
   **dieser Slice hängt also von einer Eigenschaft ab, die er nicht selbst
   herstellt**, und das gehört in seine DoD statt in eine Fußnote.
6. `make gates` und `make fullbuild`; **Review** und **Verifikation** als
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
- **Bei der dritten Beanspruchung neu gemessen, weil sich die Form geändert
  hat** (2026-08-30): mit `max-open-tasks: 0` statt `forbid-pattern` sind es
  **144 Befunde in denselben 37 Dateien**, mit der Ausnahme **null**. Die Zahl
  der *Dateien* ist unverändert, die der *Befunde* nicht — die neue Form meldet
  **je offenem Haken** statt je Abschnitt. Die drei Ausnahme-Muster aus der
  vorigen Zeile tragen unverändert, und die Grenze liegt weiter bei
  [slice-171](../done/slice-171-vorpruefungen-belegen.md): die seither
  geschlossenen Slices tragen ihre Haken gesetzt. **Das ist zugleich der Beleg,
  dass die korrigierte Praxis hält** — **acht** Dateien liegen in der Spanne
  171–182 (172 ist in Arbeit, 173–175 gibt es noch nicht), und keine trägt
  einen offenen Haken.
- **Bei der zweiten Beanspruchung nachgeprüft, nicht wiederholt:** der Bestand
  ist von 170 auf **171** `done/`-Slices gewachsen (slice-177 kam hinzu), die
  Trefferzahl bleibt **37**, die Nummern bleiben dieselben, und es gibt weiter
  **kein** `section-missing`. Der seither geschlossene Slice trägt
  seine Haken gesetzt — er hätte die Zahl sonst erhöht.
- **Die Treffer streuen nicht beliebig:** 025–104 (33 Stück), dann 160, 168,
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

**Die Design-Frage ist entschieden** (Auftraggeber, 2026-08-29): **was eine
Welle einlöst, gehört in ihren Closure-Trigger, nicht in die DoD eines Slice.**

Der Kopfkommentar von [`.d-check.closure.yml`](../../../../.d-check.closure.yml)
führt genau diese Regel als **bewusst verworfen**, mit einem sachlichen Grund:
*ein abgeschlossener Slice darf eine offene Box tragen, wenn die **Welle** sie
einlöst*. Der Fall ist belegt — `slice-094` und `slice-104` tragen je einen
offenen Release-Haken mit dem ausdrücklichen Vermerk „Wellen-Trigger, nicht
Slice-Trigger".

**Die Ablehnung war richtig für ihr Werkzeug und ist es nicht mehr für ihre
Form.** Damals gab es keinen Ort, an dem ein wellen-eingelöster Punkt sonst
hätte stehen können; heute gibt es ihn, und welle-86 §3 führt ihn bereits so.
Ein DoD-Punkt, den der Slice **selbst nicht abhaken kann**, zwingt ihn, mit
offenem Haken zu schließen — und macht damit den Haken als Zustandsfeld
unbrauchbar: er sagt dann nicht mehr „hier fehlt etwas", sondern „hier fehlt
vielleicht etwas". Genau diese Unschärfe ist der Grund, aus dem der Wächter
gebaut wird.

**Zwei Folgen, beide Teil dieses Slice:**

- Der **Kopfkommentar wird ersetzt**, wenn die Regel landet — er beantwortet
  dieselbe Frage sonst zweimal gegensätzlich. Solange die Regel nicht steht,
  ist er zutreffend und bleibt.
- Die **Form-Regel gilt ab sofort**, unabhängig vom Sensor, und steht deshalb
  in [`AGENTS.md`](../../../../AGENTS.md) §5 — ein Slice, der heute geschrieben
  wird, soll nicht in eine Form laufen, die morgen rot ist.

**Der Altbestand bleibt unberührt.** slice-094 und slice-104 liegen in der
Bestands-Ausnahme; sie nachträglich umzuschreiben hieße, ihre Lauf-Belege zu
fälschen.

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

- [x] Eine `structure`-Regel im Closure-Profil meldet einen offenen DoD-Haken
      in `done/slice-*` mit eigener, reparatur-benennender Meldung.
- [x] **`spans` läuft im selben Profil, und das ist belegt statt angenommen.**
      Der vergessene Schluss-Fence ist die deklarierte, nicht behobene Grenze
      von `max-open-tasks`; ohne `fence-unclosed` hängt dieser Wächter an einem
      einzelnen Zeichen. Beleg: eine Probe, in der ein offener Haken hinter
      einem vergessenen Fence steht — `structure` allein schweigt, das Profil
      meldet.
- [x] **Retro gemessen:** die Regel meldet ohne Ausnahme die erwartete
      Bestandszahl (144 Befunde in 37 Dateien) und mit Bestands-Ausnahme
      **null**; beide
      Ausgaben stehen in der Commit-Botschaft.
- [x] Die Bestands-Ausnahme ist als Adaption geführt (oder an
      [`MR-049`](../../../../harness/conventions.md#mr-049) angehängt) und nennt
      den Grund: gehobenes Regelwerk, nicht nachgezogene Dokumente.
- [x] Die drei Folge-Kandidaten aus §3 sind als solche benannt — nicht
      stillschweigend weggelassen.
- [x] Die neue Bedingung ist **zugestellt**, bevor sie blockiert:
      [`AGENTS.md`](../../../../AGENTS.md) §5 und die Sensors-Tabelle in
      [`harness/README.md`](../../../../harness/README.md) nennen sie
      ([`BEO-022`](../observations.md)).
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Ein Haken ist eine Selbstauskunft, und der Sensor prüft nur sie.** Wer den
  Haken setzt, ohne dass ein Review stattfand, passiert das Gate. Der Sensor
  verschiebt die Lücke von „unsichtbar" nach „behauptet" — das ist besser, aber
  es ist keine Prüfung des Reviews. — **Ausgang: weiter offen.** Nicht
  auflösbar: was der Sensor sieht, **ist** die Selbstauskunft. Der Review hat
  zwei weitere Wege gemessen, auf denen ein Haken unsichtbar bleibt — im
  **wohlgeformten** Fence (dort meldet auch `fence-unclosed` nichts) und hinter
  einem **vergessenen** Fence (dort meldet es). Alle drei stehen jetzt als
  benannte Grenzen in [`MR-056`](../../../../harness/conventions.md#mr-056) und
  an den drei Zustell-Orten. **Ein Sensor dagegen ist nicht denkbar**: ob ein
  Fence Beispiel oder Versteck ist, ist ein Urteil.
- **Die Bestands-Ausnahme ist eine Grandfathering-Klausel und altert.** Sie
  nimmt eine feste Nummernspanne heraus; jeder neue Slice fällt unter die
  Regel. Wächst die Ausnahme je wieder, ist das der Befund
  ([`BEO-013`](../observations.md)). — **Ausgang: weiter offen**, und die
  Verifikation hat ihn **beziffert**: die Ausnahme nimmt **169 von 177**
  `done/`-Slices heraus, unter der Regel stehen heute **acht** Dateien. Das ist
  die deklarierte Alterung, jetzt mit Zahl statt als Ahnung. Der
  Auflösungs-Trigger steht in [`MR-056`](../../../../harness/conventions.md#mr-056): wächst die Ausnahme, ist das der Befund.
- **Der Sensor macht den inneren Loop rot, bevor jemand schließen will.**
  `verify-closure-notes` läuft in `make fullbuild`, nicht in `gates` — die neue
  Regel gehört an dieselbe Stelle, sonst meldet sie beim Arbeiten an einem
  laufenden Slice. — **Ausgang: entfallen, mit Messung.** Die Regel steht im
  Closure-Profil, nicht im Hauptprofil; die Verifikation hat gegengeprüft, dass
  `make doc-check` mit einem offenen Haken bei **0 Befunden** bleibt. Der
  befürchtete Fall kann in dieser Verortung nicht eintreten — er käme nur
  zurück, wenn jemand die Regel ins Hauptprofil zöge, und das wäre eine eigene
  Entscheidung.
- **Die Regel wird aus einem Anlass gezogen** ([`BEO-011`](../observations.md)):
  drei Slices einer Sitzung. Gegen den Bestand gemessen ist sie es nicht — 37
  Treffer sind ein Bestand —, aber die *Dringlichkeit* stammt aus dem Anlass. —
  **Ausgang: weiter offen.** Die Unterscheidung bleibt richtig und ist nicht
  auflösbar: eine Regel kann ihren Anlass nicht loswerden. Was sich seit der
  Notiz geändert hat, stützt sie eher — die Bestandsmessung ist zweimal
  unabhängig reproduziert worden (144 Befunde in 37 Dateien), und die drei
  Anlassfälle liegen **in** der Ausnahme, tragen die Regel also gar nicht.

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

**Dritte Beanspruchung am 2026-08-30** (`open` → `in-progress`): **beide
Vorbedingungen der Rückführungen liegen vor.** Die erste — eine regel-eigene
Meldung — kam mit
[slice-177](../done/slice-177-structure-hint.md) (`structure[].hint`); die
zweite — eine Bedingung, die den **rohen** Abschnittstext liest — mit
[slice-178](../done/slice-178-offene-tasks-roh.md) (`max-open-tasks`). Damit
sind die beiden Blindstellen, an denen der zweite Bau scheiterte, keine
Blindstellen mehr: der Backtick schaltet nichts ab, und alle vier
Listen-Marker zählen, ohne dass ein Konfigurations-Autor sie kennen muss.

**Die Vorarbeit ist neu gemessen, nicht übernommen** — die Form hat sich
geändert, also ändern sich die Zahlen: 144 Befunde in denselben 37 Dateien,
mit Ausnahme null. Was **bleibt**, ist die Nummernspanne der Ausnahme; was
**dazukommt**, ist ein Befund je Haken statt je Abschnitt.

**Eine Abhängigkeit ist neu und steht in der DoD:** dieser Wächter hängt an
`spans` im selben Profil. Der vergessene Schluss-Fence ist die deklarierte,
nicht behobene Grenze von `max-open-tasks`; ohne `fence-unclosed` schaltet ein
einzelnes Zeichen die Vorbedingung ab.

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

- **Offene Beobachtungen sichten** (Register-Stand **2026-08-30**, höchste
  Kennung **`BEO-024`** — bei der dritten Beanspruchung nachgeholt: der
  Nachtlauf-Block war aktualisiert, dieser nicht, und `BEO-024` stand da schon
  im Register): [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr
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
  grünen Laufs als Beleg; [`BEO-024`](../observations.md) — **ein Zustell-Kanal
  hängt an der Arbeitsweise, die Regel aber am Inhalt**, und das trifft diesen
  Slice direkt: seine Zustellung liegt in `AGENTS.md` und in der
  Sensors-Tabelle, also in Dateien, die **jeder** Lauf liest — nicht in einem
  Kanal, der nur bei bestimmten Werkzeugwegen greift. Der Unterschied zu
  `BEO-022` ist die Richtung: dort fehlt der Kanal, hier greift er am falschen
  Kriterium. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — bei der dritten Beanspruchung neu gelesen,
  `upstream-drift.yml` zuletzt 2026-08-30T06:08:17Z, `image-scan.yml`
  2026-08-30T09:16:25Z. **Dieser Block trägt bewusst keine
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

**Geliefert.** Eine `structure`-Regel im Closure-Profil hält die zweite Hälfte
der Kanon-Bedingung für den Übergang nach `done/`: `max-open-tasks: 0` über den
`## N. Definition of Done`-Abschnitt, ein Befund **je Haken auf seiner Zeile**,
mit verfasstem Reparatur-Hinweis. Dazu
[`MR-056`](../../../../harness/conventions.md#mr-056), die Zustellung an drei
Orten und der **Ersatz** des Kopfkommentars, der dieselbe Regel bisher als
*bewusst verworfen* führte.

**Der dritte Anlauf, und diesmal lag es nicht an diesem Slice.** Beide
Rückführungen hatten Produkt-Vorbedingungen: eine regel-eigene Meldung
([slice-177](slice-177-structure-hint.md)) und eine Bedingung auf dem **rohen**
Abschnittstext ([slice-178](slice-178-offene-tasks-roh.md)). Beide liegen vor,
und damit sind die zwei Blindstellen des zweiten Baus keine mehr. **Die
Vorarbeit wurde neu gemessen statt übernommen** — die Form hat sich geändert,
also ändern sich die Zahlen: 144 Befunde in denselben 37 Dateien, mit Ausnahme
null.

**Was funktioniert hat.** Der `spans`-Punkt, den ich beim Beanspruchen in die
DoD gezogen habe, statt ihn als Fußnote zu führen. Isoliert gemessen ist die
Bedingung hinter einem vergessenen Fence **völlig blind** (0 Befunde, Exit 0);
im Repo-Profil sieht man das **nicht**, weil eine Nachbarregel zufällig
`section-empty` wirft. Ohne die isolierte Probe hätte ich eine Abhängigkeit für
gedeckt gehalten, die es nur zufällig ist.

**Was anders lief.** Der Review blockierte mit zwei HIGH, die Verifikation
urteilte konform mit sechs Präzisions-Befunden. Beide fanden dieselbe falsche
Zahl unabhängig. Das **Verhalten** hielt durchweg: 144/37 ohne Ausnahme, null
mit ihr, Positiv-Probe fängt mit Zeile und Hinweis, Ausnahme-Globs decken die
Trefferliste exakt und disjunkt (30 + 6 + 1 = 37), alle Gates grün. Rot waren
die **Aussagen darüber**:

1. **„Elf Closures" war aus der Spanne gerechnet, nicht aus dem Verzeichnis
   gelesen.** Es sind **acht** Dateien (171, 176–182); 173–175 existieren
   nirgends, und genau das steht in
   [welle-86](../welle-86-closure-uebergang-durchsetzen.md) §4 — auf das
   derselbe Absatz verlinkt. Die Zahl trug den einzigen Beleg dafür, dass die
   Grenze bei 171 richtig liegt.
2. **Der neue YAML-Kommentar referierte seine eigene Vorgänger-Fassung.** Sein
   Subjekt war ein Text, den derselbe Commit entfernt hat — Edit-Historie im
   Kommentar (§3.7). Die Herkunft gehört in die Commit-Botschaft, die
   Abgrenzung in den Kommentar.
3. **Der `hint` erscheint auch auf `section-missing`.** Gemessen: fehlt der
   DoD-Abschnitt, meldete die Regel *„offener DoD-Haken … Haken setzen"*, obwohl
   es keinen gibt. Er nennt jetzt die **Zusage** statt des Defekts und stimmt
   damit für beide Grund-Codes.
4. **Eine dritte Grenze fehlte an allen Zustell-Orten:** ein Haken **innerhalb**
   eines wohlgeformten Fenced-Blocks ist unsichtbar, und dort meldet auch
   `fence-unclosed` nichts.

**Und der Fehler, den ich beim Beheben eines anderen eingebaut habe.** Meine
Korrektur in [`harness/README.md`](../../../../harness/README.md) sagte, die
dortigen vier Grenzen gälten für `max-open-tasks` **nicht**. Für die
**Fence**-Hälften ist das messbar falsch — die Bedingung ist fence-treu wie
jede andere, und genau daraus folgt Befund 4. Nur die Inline-Code-Hälften
entfallen. Der Review hat es gemessen; das ist dieselbe Gestalt wie in
slice-180, wo eine Korrektur die zu korrigierende Zahl durch eine andere
falsche ersetzte.

**Steering-Loop-Einträge.** Zwei Zähler, keine neue Kennung:

- **[`BEO-020`](../observations.md)** (eigene Menge gemessen, fremde
  ausgesagt) — hier in zwei Ausprägungen im selben Slice: eine Zahl aus der
  **Nummernspanne** statt aus dem Verzeichnis, und eine Korrektur, die einen
  neuen Fehler einführte. **Prozedur, ergänzt:** eine Zahl über einen Bestand
  wird **gelesen**, nicht gerechnet — `ls | wc -l` statt „von 171 bis 182 sind
  elf".
- **[`BEO-002`](../observations.md)** (Semantik-Änderung nur im Körper
  nachgezogen, die Ränder bleiben stehen) — neue Ausprägung: **zwei Hälften
  desselben Pflicht-Blocks**. Bei der dritten Beanspruchung wurde der
  Nachtlauf-Teil von §7 aktualisiert, der Sichtungs-Teil nicht; der veraltete
  sah so frisch aus wie der frische. Und ausgerechnet die übersehene Kennung
  [`BEO-024`](../observations.md) handelt von Zustell-Kanälen — der Frage, die
  DoD-Punkt 6 dieses Slice beantwortet.

**Verifikation.** `make gates` Exit 0 (625 Dateien, 0 Befunde) ·
`make verify-closure-notes` Exit 0 (562 Dateien) · `make fullbuild` Exit 0,
Image-Hash `sha256:fb5c3b907001c300e9dd7cb11134ca2b10b9354a2c315ed7e436bdcc3b880f68` ·
beide Reports in [`docs/reviews/`](../../../reviews/). Die Verifikation hat den
ersetzten Kopfkommentar gegen den ausgepackten Baum von `e93d6a9` nachgerechnet:
die damals gezählten **32** Dateien sind **keine** still repariert worden, die
Differenz zu heute sind fünf Neuzugänge. Nicht benannt war dort nur der
**Einheitenwechsel** — 32 waren Dateien, 144 sind Haken.

**Offen und benannt.** Die Regel prüft den **Zustand am Ruheort**, nicht den
**Übergang**; ein Slice, der mit offenem Haken nach `done/` wandert, wird erst
vom nächsten `verify-closure-notes`-Lauf gesehen. Das ist der Gegenstand von
`slice-175`, und der Closure-Trigger von
[welle-86](../welle-86-closure-uebergang-durchsetzen.md) verlangt genau diesen
Beleg. Drei Grenzen bleiben und sind keine Lücken, sondern Eigenschaften: der
Haken ist eine Selbstauskunft, ein Haken im wohlgeformten Fence ist unsichtbar,
und ein vergessener Fence macht die Bedingung blind. **Ein vierter Punkt ist
neu und gehört einem eigenen Entscheid:** ein regel-eigener `hint` erscheint auf
Befunden, die gar nicht die Bedingung sind (`section-missing`,
`section-ambiguous`) — [ADR-0073](../../adr/0073-befund-erlaeuterung-fuer-menschen.md)
nimmt nur zwei solche Fälle aus.
