# Slice slice-102: Wellen-Lifecycle-Invariante — Roadmap-Abschnitte gegen Wellen-Dateien

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** bei Start zu eröffnen (ein Slice in Arbeit verlangt eine benannte
aktive Welle — die Zwei-Zustands-Kopplung aus
[`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)).

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(dritte Fähigkeit desselben Moduls),
[ADR-0028](../../adr/0028-planning-lifecycle-modul.md) (das Modul und sein
hermetischer Schnitt). **Anlass:** Auftraggeber-Beobachtung, dass die
Wellen-Dateien ihren Zustand genauso über das Verzeichnis tragen wie die Slices
— und damit gegen die Roadmap prüfbar sind.

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Das Modul `planning` prüft heute die **Slice**-Ebene (Ruhe-Marker ⟺ kein
`slice-*` in `in-progress/`). Dieser Slice zieht dieselbe Invariante eine Ebene
höher: die **Wellen**-Abschnitte der Roadmap gegen die Wellen-Dateien.

## 2. Warum das eine echte Invariante ist — und der Beleg

Eine Wellen-Datei trägt ihren Zustand im **Ort**, exakt wie ein Slice. Das
Baseline-Wellen-Template sagt es wörtlich: geplante Wellen bekommen **noch keine
Datei**, sie stehen in der Roadmap unter *Nächste Wellen* „und nirgends sonst —
**zwei Positionen, nicht drei**". Damit sind vier Aussagen maschinell
entscheidbar:

| # | Aussage |
|---|---|
| 1 | §Aktuelle Welle nennt eine Welle-ID ⟺ **genau eine** flache Wellen-Datei existiert, und beide tragen dieselbe Kennung |
| 2 | §Aktuelle Welle trägt den Ruhe-Marker ⟺ **keine** flache Wellen-Datei |
| 3 | Jede Zeile in §Nächste Wellen nennt eine Welle **ohne** Datei — weder flach noch im Ruheort |
| 4 | Jede Zeile in §Abgeschlossene Wellen nennt eine Welle **mit** Datei im Ruheort |

**Der Beleg stammt aus dem eigenen Repo, aus zwei aufeinanderfolgenden
Wellen-Closures — beide Male mit grünen Gates.** Bei der Closure von welle-68
wie von welle-69 stand über mehrere Commits hinweg der **Ruhe-Marker** in
§Aktuelle Welle, während die flache Wellen-Datei noch danebenlag: die Roadmap
behauptete „keine laufende Welle", das Verzeichnis sagte „eine läuft". Das ist
Aussage 2, verletzt, zweimal, unbemerkt.

Beide Male fiel es erst auf, als der Auftraggeber nachfragte — nicht durch ein
Gate. Genau die Klasse, für die dieses Modul gebaut wurde: eine Aussage in der
Roadmap, die das Verzeichnis widerlegt.

## 3. Abnahme-Punkte

1. **Alle vier Aussagen oder nur die ersten zwei?** Aussagen 1 und 2 sind das
   direkte Wellen-Pendant zur bestehenden Slice-Invariante. Aussagen 3 und 4
   brauchen zusätzlich das **Parsen von Tabellenzeilen** und eine
   Kennungs-Extraktion — mehr Vertragsfläche und die Frage, woran eine Zeile
   ihre Welle-ID trägt. Zu entscheiden, ob 3/4 in denselben Slice gehören.
   — **Ausgang: alle vier, in einem Slice** (beim Wellen-Schnitt entschieden).
   Die Abwägung hat sich gedreht, weil die Tabellenzeilen-Lexik seit welle-74
   **entdriftet** vorliegt: sie zählt nur außerhalb von Fences und braucht hier
   nur noch eine **Spalten-Adresse**, keinen Neubau.
2. **Woran erkennt das Modul eine Wellen-Datei?** Vorschlag: ein Glob analog
   `planning.slice-glob` (etwa `planning.wave-glob`, Default `welle-*.md`) plus
   das Verzeichnis, in dem flache Wellen liegen. **Achtung — dieselbe Falle wie
   bei der Closure-Fähigkeit:** die Ergebnis-Notizen (`welle-*-results.md`)
   liegen im Ruheort und matchen dasselbe Muster; Plan-Datei und Ergebnis-Notiz
   müssen unterscheidbar bleiben.
   — **Ausgang: zwei Globs und zwei Rollen, von der Messung erzwungen.**
   `waves.glob` (Default `welle-*.md`) benennt das **Plan-Dokument**,
   `waves.results-glob` (Default `welle-*-results.md`) die **Ergebnisnotiz**;
   die zweite Menge wird von der ersten **abgezogen**. Und die Rollen sind
   nicht austauschbar: die Aussagen 1/2 fragen nach dem **Plan-Dokument**
   (es liegt flach, solange die Welle läuft), die Aussage 4 nach der
   **Ergebnisnotiz** — gegen das Plan-Dokument gemessen meldet sie 19-mal falsch
   (§3a). Verglichen wird über das **Zahlen-Präfix** `welle-<n>`: die Zeile trägt
   die volle Kennung (`welle-74-geteilte-lexik-raender`), die Notiz den kurzen
   Namen (`welle-74-results.md`).
3. **Ein Grund-Code oder mehrere?** Die vier Aussagen haben verschiedene
   Reparaturen (Roadmap nachziehen · Datei verschieben · Vorschau-Zeile
   entfernen · Ergebnis-Notiz nachtragen). Nach der in
   [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md)
   festgehaltenen Begründung spricht das für Trennung — und die
   Befund-Deduplikation über (Datei, Zeile, Regel, Ziel, Grund) verlangt sie
   sogar, wenn zwei Verletzungen dieselbe Roadmap-Zeile treffen.
   — **Ausgang: vier Grund-Codes für vier Reparaturen.** `wave-drift`
   (Aussagen 1+2, das direkte Pendant zu `planning-drift`, beide Richtungen in
   **einer** Meldung wie dort) · `wave-preview-exists` (Aussage 3: eine
   Vorschau-Zeile nennt eine Welle, die schon eine Datei hat — „drei Positionen
   statt zwei“) · `wave-results-missing` (Aussage 4) · `wave-unregistered`
   (die Gegenrichtung: Ergebnisnotiz im Ruheort ohne Zeile — die Richtung aus
   **BEO-001**). Die Trennung ist nicht Geschmack: `wave-results-missing` und
   `wave-preview-exists` können dieselbe Roadmap-Zeile treffen, und die
   Deduplikation über (Datei, Zeile, Regel, Ziel, Grund) ließe sie sonst
   zusammenfallen.
4. **Verhältnis zur bestehenden Fähigkeit.** Beide lesen denselben
   `planning.heading`-Abschnitt und denselben Marker. Die Slice-Invariante prüft
   `hasActive == hasSlices`, die Wellen-Invariante `hasActive == hasWave` —
   zusammen ergibt das eine Dreier-Kopplung. Zu prüfen, ob sie sich widersprechen
   können und was dann gilt.
   — **Ausgang: sie können sich widersprechen, und beide melden.** `hasActive`
   ist **eine** Größe, aus **einer** Quelle (dem `planning.heading`-Block); sie
   wird gegen zwei Verzeichnis-Zustände geprüft, gegen `hasSlices` und gegen
   `hasWave`. Sind beide verletzt, entstehen zwei Befunde mit **verschiedenen**
   Grund-Codes und verschiedenen Reparaturen — das ist gewollt, nicht doppelt.
   Ein Widerspruch **zwischen** den Prüfungen ist ausgeschlossen, weil beide
   dieselbe linke Seite lesen; die Aussage „die Roadmap sagt aktiv“ kann nicht
   zugleich wahr und falsch sein. **Und das ist selbst eine Lexik-Frage** nach
   [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md): die
   dritte Fähigkeit ruft dieselbe Aktiv-Status-Bestimmung auf, statt sie ein
   zweites Mal zu beantworten.

## 3a. Messung: was der Bestand zu den vier Aussagen sagt

Methode wie in den Vorgänger-Slices: der **reale** Bestand dreier Planungs-Bäume
(`d-check`, `a-check` und das Beispiel-Repo des Kurses, das die Baseline-Form
vorführt).

| Aussage | Messung |
|---|---|
| 1 + 2 (aktive Welle ⟺ genau eine flache Datei) | heute in allen drei Bäumen konsistent — auch in `a-check`: dessen scheinbare Verletzung war ein **Artefakt des Default-Markers** der Probe-Konfiguration („**Keine.**“ statt „Keine aktive Welle“, Review F-8); mit konsument-gerechtem Marker gilt der Block als ruhend, und keine flache Datei liegt |
| 3 (Vorschau-Zeile ohne Datei) | **nicht messbar wie formuliert** — siehe unten |
| 4 (Abschluss-Zeile mit Datei) | gegen das **Plan-Dokument**: 19 Verletzungen über zwei Bäume. Gegen die **Ergebnisnotiz**: `d-check` 15/15 sauber, Kurs-Beispiel 1/1 sauber, `a-check` 11 Zeilen ohne Notiz |
| 4-Gegenrichtung (Datei ohne Zeile) | heute überall grün — in der welle-73-Closure **dreimal verletzt** und dort geheilt (**BEO-001**) |

**Zwei Entwurfsannahmen sind damit widerlegt:**

1. **Aussage 4 meint die Ergebnisnotiz, nicht das Plan-Dokument.** Gegen das
   Plan-Dokument gemessen meldet sie 19-mal — und jedes Mal zu Unrecht: die
   älteren Wellen sind geschlossen worden, als es die Konvention des flachen
   Wellendokuments noch nicht gab. Das **verpflichtende** Artefakt einer
   geschlossenen Welle ist die Ergebnisnotiz; sie verlangt die
   Closure-Prozedur, und nur sie ist über den ganzen Bestand vorhanden.
2. **Aussage 3 braucht eine Spalten-Adresse und greift nur bei ID-Zeilen.** Zwei
   der drei Bäume schreiben in §Nächste Wellen **Namen** („Chronologie-Ordnung in
   Tabellen“), das Baseline-Beispiel schreibt **Kennungen**
   (`welle-3-skalierung`). Eine geplante Welle hat in der gelebten Praxis dieses
   Repos noch **keine** Kennung — sie bekommt sie bei der Eröffnung. Ein
   Token-Scan über die ganze Zeile ist zudem sofort falsch: die Trigger-Spalte
   einer Vorschau-Zeile **darf** andere Wellen nennen, und genau das tut die
   eigene Roadmap seit heute (die Zeile zur Chronologie-Ordnung nennt
   `welle-75` als Trigger). Die Aussage greift also nur, wenn die **Welle-Spalte**
   eine Kennung trägt — dann aber scharf: das ist der Fall, vor dem die Baseline
   warnt („zwei Positionen, nicht drei“).

**Und ein Befund über den eigenen Bestand hinaus:** `a-check` verletzt Aussage 4
**elffach** (elf Register-Zeilen ohne Ergebnisnotiz — robust, mit
konsument-gerechtem Marker nachgemessen). Der **zwölfte** Befund des ersten
Probe-Laufs war dagegen ein **Artefakt der Probe-Konfiguration**: `a-check`
schreibt seinen Ruhe-Marker als „**Keine.**“, der Default-Marker matcht nicht,
und der Block galt fälschlich als aktiv (Review F-8 — die Messung hatte den
eigenen Konfigurations-Fehler als Bestands-Befund gelesen). Die Fähigkeit
findet dort also elf echte Rückstände; sie bleibt opt-in, und ihre Einführung
dort ist ein eigener Schritt samt konsument-gerechtem Marker.

## 4. Definition of Done

- [ ] Abnahme-Punkte entschieden; Change Request an
      [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
      (dritte Fähigkeit) + Algorithmus + Schema + begleitende ADR, falls die
      Dreier-Kopplung eine eigene Entscheidung braucht.
- [ ] Implementierung + Tests je Aussage, fail-closed; **Realdatenbeleg**: die
      beiden belegten Fenster aus den Closures von welle-68 und welle-69 wären
      rot gewesen.
- [ ] `make gates` grün; Release als **Minor** (opt-in, ohne die neuen
      Schlüssel byte-identisch — d-check findet danach mehr).

## 5. Risiken / offene Punkte

- **Der Move-Commit wird enger.** Wellen-Closure heißt heute: Roadmap ändern
  **und** Datei verschieben. Mit der Invariante müssen beide im **selben**
  Commit liegen — dieselbe Bündelung, die
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)
  für Slices schon vorschreibt. — **Ausgang:** offen; vermutlich eine Ergänzung
  der Adaption statt eines neuen Mechanismus.
- **Plan-Datei vs. Ergebnis-Notiz** teilen den Namensraum (Abnahme-Punkt 2).
  — **Ausgang:** offen; ohne saubere Trennung entstehen Falschbefunde im Ruheort.
- **Aussagen 3/4 brauchen Tabellen-Parsing**, das dieses Modul bisher nicht
  kennt. — **Ausgang:** offen bis Abnahme-Punkt 1; ein eigener Slice ist die
  Rückfallebene.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei. Unabhängig von den
`structure`- und Closure-Strängen.

**Rückführungen:** `in-progress` → `next`, falls Abnahme-Punkt 1 die Aussagen 3/4
als eigenen Slice ausweist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001** und
  **BEO-002**. **BEO-001 ist verwandt und muss abgegrenzt werden:** dort geht es
  um „existiert eine Datei, die niemand registriert?" (Artefakt ⇒ Register),
  hier um „stimmt die Roadmap-Aussage mit dem Verzeichnis überein?" (Aussage ⇔
  Zustand). Aussage 4 dieses Slice berührt BEO-001 allerdings unmittelbar — wer
  sie baut, sollte prüfen, ob das Register-Gate damit teilweise erledigt ist.
  **BEO-002** betrifft ihn ebenfalls: eine dritte Fähigkeit im selben Modul
  berührt Modul-Doku, Handbuch-§6, Config-Schema und drei CLI-Enumerationen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Change Request und Algorithmus zuerst, der
Go-Code liefert sie. Kein Brownfield: es wird kein undokumentierter Bestand
inventarisiert, sondern eine bislang nur konventionell geltende Invariante
mechanisiert.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
