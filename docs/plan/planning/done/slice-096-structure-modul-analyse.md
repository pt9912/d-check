# Slice slice-096: Modul `structure` — Analyse, Modul-Schnitt und Ablöse-Pfad

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** [welle-69-structure-schnitt](welle-69-structure-schnitt.md),
eröffnet am 2026-08-09 — dieser Strang trägt eine Closure-Bedingung **jenseits**
seiner DoD: die Analyse allein löst kein Adopter-Skript ab. Umsetzung,
Paritäts-Beleg und der Ablöse-Pfad sind Folge-Slices.

**Dieser Slice läuft ZUERST** (Auftraggeber-Entscheid 2026-08-09, nach dem
Backlog-Schnitt-Review): [slice-094](../done/slice-094-closure-zaehl-paritaet.md),
[slice-097](../done/slice-097-closure-glob-entkopplung.md) und
[slice-098](../done/slice-098-closure-note-placeholder.md) schärfen alle dieselbe
Fähigkeit, die hier neu geschnitten wird. Sie zuerst einzeln auszuliefern hieße,
dreimal eine Semantik zu versprechen, die dieser Slice gerade neu definiert —
und sie danach zu migrieren. Die ursprüngliche Reihenfolge (094 zuerst) ist
damit **umgekehrt**.

**Bezug:** **Change Request** aus dem Schwester-Repo a-check (CR 1 seiner
Werkzeug-Abdeckungs-Analyse). Berührt
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
und [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) —
deren Closure-Fähigkeit ist der **Spezialfall** des beantragten Moduls.
Formvorbild für den Ablöse-Pfad:
[ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Entscheiden und festhalten, wie d-check **Struktur-Invarianten innerhalb eines
Dokuments** bekommt (bisher deckt der Modulsatz nur **Referenz**-Invarianten ab)
— und wie die in v0.52.0 ausgelieferte, planning-lokale Closure-Fähigkeit darin
aufgeht, **ohne** Config-Bruch bei Adoptern.

## 2. Die Kollision, die zuerst geklärt gehört

Der Antrag beschreibt `require-section` · `non-empty` · fence-treues
`min-sentences` · `forbid-pattern` — das ist, bis auf die Verallgemeinerung über
**beliebige Dokumentklassen**, genau die Fähigkeit, die
[slice-093](../done/slice-093-closure-note-gate.md) als
`planning.closure.*` ausgeliefert hat. Zwei Mechanismen für dieselbe Frage will
dieses Repo nicht.

Gemessen (2026-08-09, gegen eine Kopie des a-check-Bestands mit 76 Slices):
v0.52.0 deckt **eines** der drei beantragten Skripte, und das mit Kalibrierung —
das Adopter-Muster braucht wegen der fehlenden RE2-Lookahead-Unterstützung eine
positiv formulierte Ausschluss-Alternative, und die Floskel-Prüfung ist dort
zeilen-verankert, hier Teilstring. Nicht gedeckt sind die abschnitts-treue
Task-Zählung und die benannten Pflicht-Bausteine.

### Die Messung, je Prüfung eine Aussage

Nicht je Skript, sondern je **Prüfung** — ein Skript ist eine Datei, keine
Aussage. Die drei Skripte tragen zusammen **elf** Prüfungen (Stand 2026-08-09,
480 Zeilen Shell):

| # | Prüfung | Dokumentklasse | Aussage | Beleg |
|---|---|---|---|---|
| 1 | Closure-Abschnitt vorhanden | Slice-Plan | **gedeckt** | `closure-note-missing`, gegen 76 Adopter-Slices grün gefahren |
| 2 | **genau einer**, nicht mehrere | Slice-Plan | **nicht gedeckt** → Abnahme-Punkt 3 schließt die Lücke | d-check liest laut Spezifikation den **ersten** Treffer |
| 3 | Abschnitt nicht leer | Slice-Plan | **gedeckt** | Sonderfall von `min-sentences` (≥ 1) |
| 4 | kein Template-Platzhalter | Slice-Plan | **nach Kalibrierung** | das vermessene Skript prüft **fünf literale Phrasen** — die deckt die Floskel-Liste heute ab. Der weitergehende Fall (beliebige `<platzhalter>`) stammt aus einem **anderen** Antrag und ist hier nicht Beleg |
| 5 | keine Floskel | Slice-Plan | **nach Kalibrierung** | Teilstring statt zeilen-verankert ⇒ nur eindeutige Phrasen aufnehmbar |
| 6 | Mindest-Satzzahl | Slice-Plan | **nach Kalibrierung** | Inline-Code und Satzende-Form weichen ab; eigener Antrag |
| 7 | Obergrenze der DoD-Punkte **im Abschnitt** | Slice-Plan | **nicht gedeckt** | keine Zähl-Fähigkeit über Abschnitts-Elemente |
| 8 | Lerneintrag nennt **eine von drei** Formen | Slice-Plan | **nicht gedeckt** | **keine** Marken-Frage: die Formen stehen mehrheitlich *innerhalb* des Textlaufs ⇒ `require-pattern` |
| 9 | Dateiname gibt die Slice-Nummer her | Slice-Plan | **außerhalb** | siehe Abnahme-Punkt 1 — keine Struktur **innerhalb** eines Dokuments |
| 10 | vier fette Pflicht-Marken je Anforderung | Anforderung | **nicht gedeckt** | `**Happy Path:**` u. a. **am Zeilenanfang**, nicht als Teilstring |
| 11 | Stichtags-/Namens-Ausnahme | beide | **außerhalb** | die Pfad-Ausnahme deckt sie **nicht**: der gelebte Fall ist eine Zahlen-Schwelle bzw. eine Kennungs-Liste *innerhalb* einer Datei ⇒ Nicht-Ziel (Abnahme-Punkt 5) |

**2 gedeckt · 3 nach Kalibrierung · 4 nicht gedeckt · 2 außerhalb** (Summe 11).
Die Guard-Aussagen des Moduls (Leerlauf, Mehrdeutigkeit) sind in dieser Zählung
**nicht** enthalten — sie stammen nicht aus den vermessenen Skripten, sondern
kommen mit dem Vertrag hinzu.

Die Zeile, die den Schnitt entscheidet, ist **10**: sie betrifft eine **andere
Dokumentklasse** (Anforderungen, nicht Slice-Pläne) und dieselbe Frageform.

## 3. Abnahme-Punkte

1. **Modul-Schnitt.** → **Entschieden 2026-08-09: neues Modul `structure`**,
   eigenes Bereichskürzel `STRUCT`, als **Liste** von Regeln über Datei-Globs.

   **Das Kriterium ist entschieden, nicht gewählt.**
   [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md) hat es
   festgelegt: querschnittlich ⇒ neues Kürzel, Einzelmodul ⇒ bestehende
   Anforderung ändern. Zeile 10 der Messung entscheidet die Frage: die
   Pflicht-Marken einer **Anforderung** in
   [`spec/lastenheft.md`](../../../../spec/lastenheft.md) sind dieselbe
   Frageform wie die Substanz einer Closure-Notiz, aber eine **andere
   Dokumentklasse**. Sie unter einem Schlüssel namens `planning.closure`
   abzulegen, wäre eine Lüge über den Gegenstand — `planning` prüft den
   Planning-Lifecycle, nicht das Lastenheft.

   **Die Alternative ist geprüft und verworfen:** `planning.closure` um einen
   Datei-Glob zu erweitern, würde technisch reichen (Abnahme-Punkt von
   [slice-097](../done/slice-097-closure-glob-entkopplung.md) liefert den
   Kandidaten-Filter ohnehin). Sie scheitert am Namen, nicht an der Technik —
   und ein Modul, dessen Name über seinen Gegenstand täuscht, ist genau die
   Sorte Harness-Lüge, die dieses Repo mechanisch bekämpft.

   **Die Grenze des Moduls, aus Zeile 9 gewonnen** (und damit die Antwort auf
   das Sammelbecken-Risiko in §5): `structure` prüft die Form **innerhalb** eines
   Dokuments. Eine Konvention über den **Dateinamen** ist keine solche Form,
   auch wenn sie im selben Adopter-Skript steht — sie ist eine Aussage über das
   Verzeichnis, nicht über den Text. Das ist ein **Nicht-Ziel** und gehört
   ausdrücklich in den Vertrag, nicht in ein späteres „passt schon irgendwie".

   **Konsequenz für die Modul-Liste:** `structure` wäre das 20. Regelmodul, und
   das erste, das keine **Referenz**-Invariante prüft. Damit erhält der bisher
   nie ausgesprochene Satz „d-check prüft, ob ein Dokument korrekt auf andere
   zeigt" eine benannte Erweiterung: **und ob es selbst richtig gebaut ist.**
   Das gehört in die Einleitung des Lastenhefts, nicht nur in die Modul-Liste.
2. **Ablöse-Pfad für die Closure-Fähigkeit.** → **Entschieden 2026-08-09: es
   wird nichts superseded, und kein Grund-Code wird ersetzt.**

   Die Frage war falsch gestellt, und ein Blick in den eigenen Vertrag klärt
   sie: die Spezifikation führt die Grund-Codes als **„stabil,
   maschinenlesbar"** (§4), im Unterschied zur `message`, die ausdrücklich
   *nicht* stabilitätsgarantiert ist. `closure-note-*` gegen `section-*`
   auszutauschen wäre damit kein Alias-Problem, sondern der **Bruch einer
   zugesagten Fläche** — jede Konsumenten-CI, die auf den Code filtert, bricht.
   Das ist keine Abwägung, das ist ausgeschlossen.

   **Der Pfad ist deshalb umgekehrt zur ersten Annahme:**
   - `structure` entsteht **neben** der Closure-Fähigkeit und führt seine
     eigenen Codes.
   - [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
     bleibt unverändert bestehen — samt Aktiv-Status-Invariante **und** ihren
     drei Codes. Auch [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md)
     wird **nicht** superseded.
   - Was zusammengeführt wird, ist die **Semantik**: die Closure-Fähigkeit wird
     in der Spezifikation als **Preset** über denselben Struktur-Regeln
     definiert (gleiche Abschnitts-Bestimmung, gleiche Fence-Behandlung, gleiche
     Zählung). Damit können die beiden nicht auseinanderlaufen, ohne dass ein
     Test es merkt — die Doppelung liegt allein in der **Config-Oberfläche**,
     nicht in der Mechanik.

   Das ist dieselbe Form wie bei
   [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md), aber
   aus dem umgekehrten Grund: dort blieb der alte Schlüssel Alias, **weil** die
   Codes gleich waren; hier bleibt die alte Anforderung ganz stehen, **weil** die
   Codes verschieden sind. Der Preis ist zwei Config-Wege für verwandte Fragen —
   der Gegenwert ist ein Konsumenten-Vertrag, der hält.

3. **Kardinalität mehrerer Abschnitte.** → **Entschieden 2026-08-09: ja, mit
   eigenem Grund-Code — in beiden Oberflächen.**

   Das Adopter-Skript meldet Mehrdeutigkeit (Messzeile 2), und die
   Aktiv-Status-Prüfung **desselben Moduls** ist an genau dieser Stelle längst
   fail-closed: eine mehrfach vorkommende kanonische Überschrift ist dort ein
   Befund, kein „nimm die erste". Dass die Closure-Fähigkeit still den ersten
   Treffer nimmt, war keine Entscheidung, sondern ein Versäumnis — und ein
   stiller: ein zweiter, stehengebliebener Abschnitt ist der **typische**
   Platzhalter-Fall.

   `structure` bekommt `section-ambiguous`, die Closure-Fähigkeit additiv
   `closure-note-ambiguous`. Ein eigener Code, kein Sammelbefund, nach der
   Begründung aus [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md):
   drei Codes für drei Reparaturen — hier ist die Reparatur „den überzähligen
   Abschnitt entfernen", nicht „den fehlenden schreiben" und nicht „Substanz
   ergänzen". Additiv ⇒ SemVer-Minor, kein Bruch.

   **Folge für [slice-094](../done/slice-094-closure-zaehl-paritaet.md):** die
   dort bewusst verengte Paritäts-Zusage kann wieder auf volle Deckung gehen,
   sobald dieser Code liegt.

4. **Semantik der fehlenden Bausteine.** → **Entschieden 2026-08-09 — und die
   Messung hat einen Baustein mehr ergeben, als der Antrag nennt.**

   - **Zählung im Abschnitt** (`max-tasks`): gezählt werden Task-Items
     (`- [ ]`/`- [x]`) **innerhalb** des Abschnitts, nach derselben
     Fence-Bereinigung wie die Satzzählung. Genau daran ist die Skript-Variante
     gescheitert: sie zählte dateiweit und musste den Abschnitts-Schnitt selbst
     nachbauen.
   - **Benannte Marken — nach dem Review neu gefasst, weil ich die Form
     angenommen statt gemessen hatte.** Gemessen an beiden Repos: die Marke steht
     als **Listen-Item** (`- **M:**`, 108×), bare (`**M:**`, 44×) und
     **qualifiziert** (`- **M (Zusatz):**`). Meine erste Fassung („nach führendem
     Whitespace") hätte die Listen-Form ausgeschlossen und damit **jede**
     Anforderung des eigenen Lastenhefts rot gemeldet. Die Zusage lautet jetzt:
     hervorgehobener Textlauf am Zeilen-Anfang **nach optionalem Listen-Marker**.
     Ein `Boundary` im Fließtext genügt weiterhin nicht.
   - **Messzeile 8 ist gar keine Marken-Frage.** Die gesuchten Lerneintrag-Formen
     stehen mehrheitlich **innerhalb** des Textlaufs, nicht an seinem Anfang —
     eine Marken-Alternative hätte sie nicht gefunden. Der Vertrag bekommt
     stattdessen `require-pattern` als Spiegelbild von `forbid-pattern`; die
     zunächst geplante Form „mindestens eine von" **entfällt ersatzlos**.

5. **Grandfathering.** → **Entschieden 2026-08-09: `exempt-paths` je Regel — und
   ausdrücklich keine Stichtags-Mechanik.**

   Der Antrag verweist auf die bestehende Ausnahme-Mechanik, und die reicht für
   Pfad-Ausnahmen. Sie reicht **nicht** sauber für den gelebten Fall des
   Adopters („erst ab Slice 52"): eine Zahlen-Schwelle im **Dateinamen** ist als
   Glob nur umständlich und fehleranfällig ausdrückbar.

   Trotzdem **kein** eigener Mechanismus. Eine Stichtags-Regel zu lernen hieße,
   die Kennungs-Konvention des Adopters zu interpretieren — d-check müsste aus
   einem Dateinamen eine Ordnung lesen. Das ist dieselbe Grenze, die
   [`DC-FA-TRK-001`](../../../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)
   schon einmal gezogen hat (Index-Wahrheit statt eines zweiten
   `gitignore`-Interpreters): **kein zweiter Regel-Interpreter im Werkzeug.**
   Die Grenze gehört benannt, damit der nächste Antrag sie nicht neu verhandeln
   muss.

## 4. Definition of Done

- [x] Abnahme-Punkte 1–5 entschieden und begründet; Change Request in
      [`spec/lastenheft.md`](../../../../spec/lastenheft.md) formuliert
      (Bereichskürzel, Akzeptanzkriterien-Trio je Grund-Code) + begleitende ADR
      mit dem Supersede-/Alias-Pfad.
- [x] **Abdeckungs-Messung** je beantragter Prüfung: gedeckt / gedeckt nach
      Kalibrierung / nicht gedeckt — mit Beleg, nicht per Lektüre. Die
      Paritäts-Fixtures liegen im Antragsteller-Repo vor und werden beigezogen.
- [x] Folge-Slices geschnitten und als Dateien in `open/` angelegt:
      [slice-099](../open/slice-099-structure-modul.md) (Implementierung, in
      **einem** Slice — die zunächst geplante Aufteilung scheiterte an der
      Release-Grenze) und [slice-101](../done/slice-101-fence-unbalanciert.md)
      (der beim Review gefundene **ausgelieferte** Silent-Grün-Pfad im
      Fence-Automaten, unabhängig von diesem Antrag). **Kein Alias-Slice** — es
      wird nichts abgelöst (Abnahme-Punkt 2).

      **Entscheidung über die drei bereits liegenden Slices:** alle drei
      **bleiben eigenständig**, keiner geht im Modul auf. Der Grund ist derselbe
      wie in Abnahme-Punkt 2 — `structure` ersetzt die Closure-Fähigkeit nicht,
      es steht daneben, also bleiben auch ihre Schärfungen dort, wo sie wirken.
      - [slice-094](../done/slice-094-closure-zaehl-paritaet.md) (Zähl-Parität)
        wirkt nach der Preset-Kopplung auf **beide** Oberflächen und wird damit
        wichtiger, nicht überflüssig. Seine verengte Paritäts-Zusage kann
        zurückgenommen werden, sobald `closure-note-ambiguous` aus
        [slice-099](../open/slice-099-structure-modul.md) liegt.
      - [slice-097](../done/slice-097-closure-glob-entkopplung.md)
        (`closure.glob`) betrifft die Kandidaten-Menge **der Closure-Fähigkeit**;
        `structure` hat mit seinem eigenen Datei-Glob einen anderen Weg und löst
        das dort nicht mit.
      - [slice-098](../done/slice-098-closure-note-placeholder.md)
        (Platzhalter-Erkennung) bleibt an der Closure-Fähigkeit. Die allgemeine
        Entsprechung ist `forbid-pattern`, das **ohne** die drei
        Falsch-Positiv-Ausschlüsse auskommt — wer ein Muster selbst schreibt,
        verantwortet es. Genau die gehärtete Erkennung ist der Wert des Presets.

      **Wellen-Zuordnung:** 099 und 100 gehören zu dieser Welle; 094/097/098
      bilden die Folge-Welle, deren Zuschnitt damit entschieden ist.

## 5. Risiken / offene Punkte

- **Contract-Churn bei Adoptern:** `planning.closure` ist seit v0.52.0
  ausgeliefert. Ein Supersede ohne Alias wäre ein Bruch nach wenigen Tagen.
  — **Ausgang: entfallen.** Abnahme-Punkt 2 hat ergeben, dass gar nicht
  superseded wird: die Grund-Codes sind laut Spezifikation §4 stabilitäts-
  garantiert, also bleibt die bestehende Anforderung samt Codes stehen und wird
  nur semantisch als Preset über die gemeinsamen Struktur-Regeln definiert. Es
  gibt keinen Umstiegs-Zeitpunkt, an dem etwas brechen könnte.
- **Der Modul-Schnitt könnte zu breit geraten.** „Struktur-Invarianten" ist eine
  Kategorie, keine Prüfung; ohne scharfe Grenze wächst `structure` zum
  Sammelbecken. — **Ausgang: eingetreten und begrenzt.** Die Messung hat mit
  Zeile 9 (Dateinamen-Konvention) sofort einen Kandidaten geliefert, der im
  selben Adopter-Skript steht und **nicht** hineingehört. Die Grenze ist als
  Nicht-Ziel in Abnahme-Punkt 1 formuliert: Form **innerhalb** eines Dokuments,
  nicht Aussagen über seinen Ort. Sie gehört so in den Vertrag.
- **Fremd-Repo-Abhängigkeit:** die Paritäts-Fixtures liegen nicht hier.
  — **Ausgang:** offen; beizuziehen, nicht nachzubauen.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe. **Keine Vorbedingung** — dieser
Slice ist der erste des Strangs. Die frühere Fassung hängte ihn an
[slice-094](../done/slice-094-closure-zaehl-paritaet.md) („erst die Zähl-Parität, dann
messen"); der Schnitt-Review hat gezeigt, dass das die Reihenfolge verkehrt: 094
sagt Deckungsgleichheit einer Semantik zu, die **dieser** Slice gerade neu
definiert, und weil ausgeliefert wird, was zuerst fertig ist, wäre der Vertrag
dann schon draußen (F-3).

**Rückführungen:** `in-progress` → `open`, falls die Messung ergibt, dass der
Antrag mehrere unabhängige Module beschreibt statt eines.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Analyse berührt `spec/` und `docs/plan/`; die Folge-Slices
  zusätzlich `internal/`. Alle unter dem Repo-Default GF
  (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001** und — seit
  der Re-Review dieses Slice — **BEO-002** (eine Semantik-Änderung wird nur im
  Dokumentkörper nachgezogen, ihre Ränder bleiben stehen). BEO-002 ist in genau
  diesem Slice zweimal eingetreten und dort geheilt; die Belege werden bei seiner
  Closure ergänzt. **Zu BEO-001 bewusst festgehalten, weil die Verwechslung
  naheliegt:** `structure` deckt es **nicht** ab. Dort geht es um eine Referenz **zwischen** Dokumenten (existiert
  eine Datei, die niemand registriert?), hier um die Form **innerhalb** eines
  Dokuments. Wer BEO-001 in diesem Slice erledigt glaubt, lässt die Lücke offen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — der Slice liefert Analyse und Zusage
(Change Request + ADR), der Code folgt in den Folge-Slices. Kein Brownfield: die
abzulösenden Skripte liegen in einem **anderen** Repo und werden nicht
rückdokumentiert, sondern durch eine eigene, spec-first formulierte Fähigkeit
ersetzt.

## 9. Closure-Notiz (nach `done/`)

Geliefert: der Schnitt für ein 20. Regelmodul `structure`, der Change Request
(Lastenheft 0.51.0), die Algorithmus-Sektion, [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md)
`Accepted` und zwei Folge-Slices. Alle fünf Abnahme-Punkte sind entschieden und
begründet, die Abdeckung ist je Prüfung gemessen statt gelesen.

- **Lernsignal — drei Reviews, und der teuerste Befund war, dass der Vertrag
  seine eigene Begründung nicht ausdrücken konnte.** Der Modul-Schnitt wurde an
  **einer** Messzeile entschieden (die Pflicht-Bausteine einer Anforderung); der
  Vertrag sagte dann „genau ein Abschnitt je Regel" und machte Mehrfachtreffer
  zum Abbruch — Anforderungen sind aber **wiederkehrende** Abschnitte *einer*
  Datei. Die Lehre ist allgemein: **wer eine Entscheidung an einem Beispiel
  festmacht, muss das Beispiel danach gegen den fertigen Vertrag halten.** Ich
  hatte es gegen die *Absicht* gehalten.
- **Zweites Lernsignal — zweimal die Form angenommen statt gemessen.** Die
  Marken-Verankerung war in zwei aufeinanderfolgenden Fassungen falsch, beide
  Male aus demselben Grund. Erst die Messung (108 Listen-Item, 44 bare, dazu
  qualifizierte) zeigte, dass meine Fassung **jede** Anforderung des eigenen
  Lastenhefts rot gemeldet hätte. Bei der Schwelle und der Floskel-Liste hatte
  ich in slice-093 gemessen; hier nicht — **der Reflex ist nicht stabil, er
  braucht die Regel.**
- **Architektur-Beobachtung: eine Zusage kann eine andere unerreichbar machen.**
  Ich habe die Nicht-Stabilität der Befund-`message` als Grund benutzt, den
  Supersede *auszuschließen* — und im selben Dokument dieselbe `message` zum
  *Träger* der Unterscheidung zwischen sechs Bedingungen gemacht. Dazu kam, dass
  die Befund-Deduplikation über (Datei, Zeile, Regel, Ziel, Grund) zwei
  Verletzungen ohnehin zusammenfallen ließe. **Wer eine Eigenschaft als Argument
  benutzt, darf sie zwei Absätze später nicht als Werkzeug verwenden.**
- **Der Nicht-Supersede war der richtige Ausgang, aber aus dem Grund, den ein
  Blick in den eigenen Vertrag lieferte** — nicht aus Geschmack: Grund-Codes sind
  „stabil, maschinenlesbar" zugesagt. Die erste Fassung dieses Slice hatte den
  Supersede noch als Formsache behandelt.

**Reviews:** drei unabhängige Frischkontext-Läufe —
[Vertrag](../../../reviews/2026-08-09-slice-096-vertrag-review.md) und
[Umsetzbarkeit](../../../reviews/2026-08-09-slice-096-umsetzbarkeit-review.md)
(zusammen 4 HIGH / 19 MEDIUM), danach eine
[bestätigende Re-Review](../../../reviews/2026-08-09-slice-096-re-review.md)
(14 geheilt / 12 teilweise / 5 offen, dazu 12 neue Befunde). Alle eingearbeitet.
Dass [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md) bis zuletzt
`Proposed` war, hat zwei Korrekturrunden überhaupt erst möglich gemacht —
`Accepted` hätte vier Fehler eingefroren.

**Beobachtungs-Register:** [BEO-002](../observations.md) auf 2× erhöht, Beleg
ergänzt. Die Beobachtung ist **in diesem Slice** zweimal eingetreten (nach der
Mehrdeutigkeits-Härtung blieb „erster Treffer" an drei Stellen stehen; nach der
großen Überarbeitung referierten Spezifikations-Historie, ADR-Index,
ADR-Konsequenzen und das Entscheidungs-Protokoll weiter die verworfene Fassung)
— der Zähler zählt sie als **einen** Beleg, weil die Beleg-Form Slice-Kennungen
verlangt. Die Klasse ist damit dichter, als die Zahl zeigt.

**Anschluss:** [slice-099](../open/slice-099-structure-modul.md) (Implementierung)
und [slice-101](../done/slice-101-fence-unbalanciert.md) — Letzterer ist der beim
Review gefundene **ausgelieferte** stille Grün-Pfad im Fence-Automaten und läuft
**vor** 099, damit das neue Modul ihn nicht erbt.

**Risiken aus §5:** Contract-Churn **entfallen** (es wird nichts superseded);
Sammelbecken-Grenze **eingetreten und begrenzt** (die Dateinamen-Konvention war
der Kandidat, sie ist Nicht-Ziel); Fremd-Repo-Fixtures **weiter offen**, in die
DoD von [slice-099](../open/slice-099-structure-modul.md) übernommen.
