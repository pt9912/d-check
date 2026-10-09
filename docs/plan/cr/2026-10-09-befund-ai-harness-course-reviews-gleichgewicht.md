# Befund von `ai-harness-course` — `reviews` im Gleichgewichtszustand eines Kurs-Repos

**Absender:** `ai-harness-course` · **Eingegangen:** 2026-10-09
**Richtung:** zweite eingehende Antwort auf den ausgehenden
[CR](2026-10-09-cr-ai-harness-course-reviews-zusage.md)
([`MR-036`](../../../harness/conventions.md#mr-036)); die erste ist die
[Antwort](2026-10-09-antwort-ai-harness-course-reviews-zusage.md).
**Berührt:** [`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
**Stand:** **angenommen** — Punkt 1 umgesetzt (slice-271, kommt mit dem nächsten
Release), Punkt 2 geplant
(slice-272), Punkt 3 in der nächsten Release-Prep. Entscheidung je Punkt unten.

---

## Wortlaut

Gemessen beim Absender mit `v0.85.0`. Kennungen des Absenders sind nicht
verlinkt; die Satzzeichen der Reproduktions-Zeile sind geglättet, sonst
unverändert.

> Beim Umsetzen eures CR (Kurs-Welle 161) sind drei Fälle aufgefallen. In ihnen
> meldet reviews etwas Falsches oder etwas nicht. Zwei davon machen die
> empfohlene Konfiguration (match: name, require-promises: true) in jedem
> Kurs-Repo regelmäßig rot, obwohl nichts fehlt. Der Kurs hält Welle 161
> deshalb zurück, bis diese Fälle geklärt sind.
>
> 1. Nach dem Archivieren ist das Gate rot, obwohl nichts fehlt.
>
> Der Kurs archiviert abgeschlossene Slices in zwei Formen. Volltext und
> Review-Report wandern ins Archiv, an der Stelle des Slice bleibt ein Stub
> ohne DoD:
>
> - mit Wellen: Stub unter done/&lt;welle-id&gt;/slice-&lt;Kennung&gt;.md,
>   Volltexte in done/&lt;welle-id&gt;/archiv.zip;
> - ohne Wellen: Stub flach als done/slice-&lt;Kennung&gt;.md, daneben
>   done/slice-&lt;Kennung&gt;-archiv.zip.
>
> Sind alle abgeschlossenen Slices archiviert, also im normalen Ruhezustand,
> misst ihr:
>
> | Aufbau | Ergebnis |
> |---|---|
> | mit Wellen: nur Stub unter done/welle-1/ + Ergebnisnotiz flach | review-missing „leere Pruefmenge: 0 Kandidat(en) … fail-closed", auch ohne require-promises |
> | ohne Wellen: nur flacher Stub, require-promises: true | review-missing „keine Review-Zusage unter 1 Kandidat(en)" |
> | ohne Wellen: dazu skip-pattern: '(?m)^&gt; \\*\\*ARCHIVIERT' | review-missing „leere Pruefmenge … fail-closed" |
>
> Alle drei Läufe sind rot, und in allen drei fehlt nichts: Die Reports liegen
> im Archiv, geprüft wurde vor dem Archivieren. Es gibt keine Konfiguration,
> die das grün lässt. skip-pattern nimmt die Stubs zwar heraus, macht die
> Menge damit aber leer, und eine leere Menge ist fail-closed.
>
> Vorschlag: Eine Kandidatenmenge, die nur deshalb leer ist (oder keine Zusage
> trägt), weil skip-pattern Dateien ausgenommen hat, ist kein
> Fail-Closed-Fall. Fail-closed bleibt der Fall „gar keine slice-*.md
> vorhanden", also Fehlkonfiguration oder falscher Pfad. Unterscheidbar wird
> das über die Zahl der übersprungenen Dateien: übersprungen > 0 und Rest = 0
> ⇒ still, oder höchstens ein Hinweis.
>
> 2. match: name deckt einen Präfix-Slice mit.
>
> Bei den Slices slice-cache und slice-cache-warmup mit nur einem Report
> 2026-10-09-slice-cache-warmup.md meldet ihr 0 Befunde. slice-cache hat
> keinen eigenen Report. Der Abgleich „Basisname, gefolgt von einem Zeichen,
> das weder Buchstabe noch Ziffer ist" lässt - als Grenze zu, Slug-Kennungen
> bestehen aber selbst aus -. Gegenprobe: slice-cachex deckt slice-cache
> korrekt nicht.
>
> Vorschlag: Ein Report deckt nur den längsten Slice-Basisnamen, der in seinem
> Namen passt. Oder ein Report, der mehrere Slices treffen würde, wird als
> mehrdeutig gemeldet.
>
> 3. Zur Kenntnis, kein Fehler: Auch ein ungehakter Punkt „- [ ] Review
> durchgeführt" zählt als Zusage. Das ist vertretbar und sollte im Handbuch
> stehen; der Kurs benennt es in seiner Vorlage.
>
> Was der Kurs tut:
>
> - Der Pin auf v0.85.0 geht mit match: name im Beispiel jetzt rein. Das
>   Beispiel archiviert nicht, dort ist das Gate damit wirklich scharf; vorher
>   lief es leer.
> - Welle 161, also Vorlage mit require-promises, Modul-10-Absatz und
>   Team-Sim-Gruppe s31, wartet auf eure Antwort zu Punkt 1. Danach würde die
>   Vorlage zusätzlich skip-pattern für Stubs setzen.
> - Punkt 2 führt der Kurs bis zu einem Fix als benannte Grenze.
>
> Reproduktion: Für jeden Fall reichen eine .d-check.yml mit modules: [reviews]
> und den genannten Schlüsseln, eine Slice-Datei mit „- [x] Review
> durchgeführt, Report unter docs/reviews/ liegt vor" bzw. ein Stub mit
> „&gt; \*\*ARCHIVIERT\*\* — Volltext: …" und ein leeres oder befülltes
> docs/reviews/.

---

## Einordnung (am Code gemessen, 2026-10-09)

Nachgestellt mit dem Produktstand `v0.85.0` in einem Wegwerf-Baum, je Fall
eine `.d-check.yml` und die genannten Dateien:

- **Fall 1 trifft zu, alle drei Zeilen.** Mit Wellen ist die Menge leer, auch
  mit `recursive` und `skip-pattern`; ohne Wellen meldet `require-promises`
  den Stub als Kandidaten ohne Zusage, mit `skip-pattern` die leere Menge.
- **Fall 1 ist nicht auf `reviews` beschränkt.** `planning.closure` meldet
  auf demselben Baum `closure-note-missing`, eine `structure`-Regel über
  `done/slice-*.md` `section-missing` — beide, weil die Leere nach dem Abzug
  von `skip-pattern` zählt. So war es beim Einführen der Schlüssel bewusst
  festgelegt; es übersah den Ruhezustand eines archivierenden Repos.
- **Fall 2 trifft zu.** `slice-cache` bleibt mit dem Report von
  `slice-cache-warmup` still. Die Spezifikation nennt die Präfix-Deckung als
  Grenze; sie ist damit benannt, aber nicht gewollt.
- **Fall 3 trifft zu** und ist gewollt: Die Zusage hängt am Text des
  DoD-Punkts, nicht am Haken. Die Spezifikation sagt es, das Handbuch nicht.

---

## Entscheidung je Punkt

1. **Leere Menge nach `skip-pattern` — angenommen, deklariert, in allen drei
   Modulen** (slice-271). Der Default bleibt fail-closed. Ein neuer Schlüssel
   `skip-allows-empty: true` neben `skip-pattern` — in `reviews`,
   `planning.closure` und je `structure`-Regel — erklärt eine Menge, die erst
   `skip-pattern` leert, zum Ruhezustand; dann ist sie kein Befund.
   `require-promises` zählt ohnehin nur die Kandidaten, die `skip-pattern`
   übrig lässt. Fail-closed bleibt die Menge, aus der `skip-pattern` nichts
   genommen hat. Deklariert statt per Default still, weil die stille Form den
   Umzugs-Wächter abschaltete: ein flacher Stub neben Volltexten, die in ein
   nicht gelesenes Unterverzeichnis umgezogen sind, meldete dann nichts mehr.
   Für die Vorlage des Kurses heißt das: `skip-pattern` **und**
   `skip-allows-empty: true` setzen, mit Wellen dazu `recursive: true` — sonst
   liest `reviews` die Stubs unter `done/<welle-id>/` gar nicht und die
   Menge bleibt leer, ohne dass `skip-pattern` etwas genommen hat.
2. **Präfix-Deckung unter `match: name` — angenommen, längster Name
   gewinnt** (slice-272). Ein Report deckt nur den längsten Slice-Basisnamen,
   der in seinem Namen passt; gezählt werden dabei alle Slice-Dateien unter
   `done/`, auch die übersprungenen Stubs, damit ein archivierter längerer
   Name nicht den kürzeren deckt.
3. **Ungehakte Zusage — zur Kenntnis.** Das Handbuch sagt es mit der
   nächsten Release-Prep.

Beide Fixes gehen mit dem nächsten Release hinaus; die Release-Notiz nennt
sie.
