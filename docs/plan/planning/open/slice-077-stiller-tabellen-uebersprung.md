# Slice slice-077: Stiller Tabellen-Übersprung — eine irrelevante Tabelle verschluckt die nächste

**Status:** open (Eingang, auf Wellen-Einplanung wartend). **Erfasst 2026-07-17
als Messbefund; die Regel ist bewusst NICHT entschieden.**

**Bezug:** betrifft
[`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 5 und — über den geteilten Reader —
[`DC-FA-XREF-001.a`](../../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency).
Berührt den Guard aus
[ADR-0037](../../adr/0037-trace-tabellenquellen-nullmengen-guard.md).
**Kein ADR — bewusst:** die tragende Regel ist offen. Sie wird benannt, wenn sie
belegt ist, nicht vorher (Lehre aus
[slice-074](slice-074-kommentar-suffix-tabellenzeilen.md) §2).

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

Eine **irrelevante** Tabelle (ihr Header bindet keine konfigurierte Rolle)
verschluckt die unmittelbar folgende **relevante** Tabelle samt Header,
Trennzeile und allen Datenzeilen — und weil die verschluckende Tabelle
irrelevant ist, meldet `extractTable` nichts. **Die Anforderungen der zweiten
Tabelle existieren nicht, und kein Gate sagt es.**

Gemessen gegen das **ausgelieferte** `ghcr.io/pt9912/d-check:v0.45.1` <!-- d-check:ignore (historischer Beleg gegen das damals ausgelieferte Image; darf nicht auf die aktuelle Version mitwandern) -->, Fixture
`fx-s`, **ohne jeden Marker**:

```markdown
| ID | Anforderung | Notiz |          ← relevante Tabelle 1, gedeckt
|---|---|---|
| F-1 | Alpha | ok |
                                       ← Leerzeile
| Werkzeug | Zweck | Stand |          ← IRRELEVANTE Tabelle, gleiche Breite
|---|---|---|
| foo | bar | baz |
| ID | Anforderung | Notiz |          ← relevante Tabelle 3: als Datenzeile gefressen
|---|---|---|
| F-2 | Beta | ok |
| F-3 | Gamma | ok |
```

```text
v0.45.1 --trace --require-complete   ⇒  1 Anforderung(en), 0 Waise(n).   EXIT=0
                                        (F-2 und F-3 sind echte Waisen)
```

Ein Lauf, der zwei echte Waisen verschweigt und **grün** endet. Vorbedingung ist
allein: gleiche Breite, keine Leerzeile. Kein Marker, keine Direktive, keine
Sonderform.

## 2. Entscheidungen / Regel

**Keine.** Das ist der Punkt dieses Slices.

Was **feststeht** (gemessen, 2026-07-17):

- **Der Defekt ist ausgeliefert** und älter als
  [slice-074](slice-074-kommentar-suffix-tabellenzeilen.md) — er braucht keinen
  Ignore-Marker. slice-074s R3-F-1 **verbreitert** dieses Loch, es entsteht dort
  nicht. v0.45.1s korrektes Verhalten an der Marker-Variante war **Zufall**: die
  kaputte Zeile setzte den Header-Scan versehentlich neu auf.
- **Ein GFM-Parser löst es nicht.** goldmark v1.8.4 sieht `fx-s` **exakt** wie
  der heutige Reader: zwei Tabellen, die zweite mit fünf Datenzeilen. **GFM gibt
  uns recht** — ohne Leerzeile beginnt keine neue Tabelle. Die Frage ist
  **Policy**, nicht Grammatik.
- **Die naheliegende Regel ist widerlegt.** „Eine Tabelle endet, wo eine neue
  beginnt (`isNewTableHeader` unbedingt)" zerreißt Fixture `fx-t`: eine
  Datenzeile aus lauter `---` sähe wie Header+Trennzeile aus, die laufende
  Tabelle endete davor, und ihre Anforderungen verschwänden **lautlos**. v0.45.1
  liest `fx-t` heute korrekt (3 Anforderungen).

Was **offen** ist: die Regel. Kandidaten, alle unbelegt:

| Kandidat | Idee | Bekanntes Problem |
|---|---|---|
| Trennzeile im Body ⇒ Exit 2 | eine Trennzeile mitten in einer Tabelle ist ein Dokumentdefekt, **unabhängig von Relevanz** | trifft `fx-t` (heute grün ⇒ dann Exit 2); braucht ein Signal, das auch für **irrelevante** Tabellen greift — `badLine` tut das nicht |
| Anforderungs-förmige Zeile in irrelevanter Tabelle ⇒ laut | „hier steht etwas, das wie eine Anforderung aussieht, und ich lese es nicht" | heuristisch; `id-pattern` gegen Zellen einer Tabelle, die per Definition nicht gebunden ist |
| Relevanz-Prüfung pro Zeile statt pro Tabelle | Header-Bindung erneut versuchen, wenn eine Zeile wie ein Header aussieht | das ist `isNewTableHeader` unbedingt ⇒ an `fx-t` widerlegt |

**Bedingung an jede künftige Regel:** sie muss an `fx-s` **und** `fx-t` **und**
`fx-m` gleichzeitig bestehen, gegen das ausgelieferte Image gemessen, und per
Mutation gepinnt sein. Fünf Fassungen in slice-074 haben je einen Zweig geprüft
und auf die Klasse geschlossen; dreimal war das Ergebnis ein stilles Grün.

## 3. Definition of Done

- [x] **Befund erfasst** samt Reproduktion gegen das ausgelieferte Image.
- [ ] **Tragende Regel** benannt und begründet ⇒ **eigene ADR**.
- [ ] **Spezifikation:** [`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
  Schritt 5 geschärft.
- [ ] **Fixtures zuerst, rot:** `fx-s` (der Befund), `fx-t` (die Gegenprobe, die
  die naheliegende Regel killt), `fx-m` (slice-074s Marker-Variante).
- [ ] **Mutations-Härte:** jede neue Grenze kippt einen Test — gemessen.
- [ ] **Realdatenbeleg** + **unabhängiger, kontext-getrennter Review vor** dem
  Release; `make gates`/`make ci` grün.

## 4. Risiken / offene Punkte

- **Der Defekt bleibt bis dahin ausgeliefert.** Er ist von den drei offenen
  Tabellen-Defekten der einzige, der heute **still Waisen verschweigt** —
  slice-074 ist laut (Exit 2), [slice-076](../done/slice-076-markdown-lexik-commonmark.md)
  ist Blindheit ohne Falschaussage. Dieser hier sagt „0 Waise(n)", während zwei
  existieren.
- **Praxis-Häufigkeit unbekannt.** In den 522 gemessenen Realdateien kam die
  Vorbedingung **nicht** vor (grid-gyms `spec/` ist sauber). Der Befund ist
  konstruiert-reproduzierbar, nicht real beobachtet — das senkt die Dringlichkeit,
  nicht die Gültigkeit.
- **Kopplung an slice-074:** dessen Toleranz **verbreitert** dieses Loch. Wer
  slice-074 wieder aufnimmt, muss diesen Slice mitdenken, sonst wird R3-F-1 zum
  vierten Mal reproduziert.
- **Die Verlockung ist eine schnelle Regel.** Sie ist der Grund, warum dieser
  Slice ohne ADR erfasst wird.

## 5. Trigger

Isolierungs-Messung während der Rücknahme von
[slice-074](slice-074-kommentar-suffix-tabellenzeilen.md) (2026-07-17): auf die
Frage, ob R3-F-1 die Klasse **einführt** oder eine bestehende **verbreitert**,
wurde `fx-s` ohne Marker gegen v0.45.1 gefahren — und das ausgelieferte Image
verschwieg zwei echte Waisen. Der Befund war in drei unabhängigen Reviews und
fünf Implementierungsanläufen nicht aufgetaucht, weil alle die **Marker**-Frage
untersuchten.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Die Spezifikation führt, der Code folgt. Hier besonders: die
Regel gehört **zuerst** in
[`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 5 und in eine ADR — der Code ist die letzte Station, nicht die erste.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend._
