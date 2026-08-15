# Slice slice-101: Unbalancierter Fence verschluckt still den Rest des Abschnitts

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** [welle-70-fence-lexik](../done/welle-70-fence-lexik.md), eröffnet am
2026-08-09. Die Welle bündelt nur diesen Slice — nicht wegen eines Mehr
gegenüber der DoD, sondern weil `make planning-check` einen Slice in Arbeit ohne
benannte aktive Welle nicht zulässt (die Zwei-Zustands-Kopplung aus
[`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)).

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Note-Struktur), [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md)
(die Fence-Lexik und ihre bewusst offen gelassene Grenze).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Einen **ausgelieferten stillen Grün-Pfad** schließen: ein ungeschlossener oder
verschachtelter Fenced-Code-Block verschluckt in v0.52.0 alles hinter sich, und
die Bedingungen der Closure-Note-Struktur laufen darüber grün.

## 2. Der Beleg

Nachgestellt am 2026-08-09 gegen das veröffentlichte Image, mit einer wörtlich
konfigurierten Floskel hinter einem ungeschlossenen Fence:

```text
Floskel hinter ungeschlossenem Fence  → 0 Befund(e), Exit 0
dieselbe Floskel ohne den Fence       → closure-note-boilerplate, Exit 1
```

Ein Autor, der einen Code-Block nicht schließt, schaltet damit unbemerkt die
Prüfung des restlichen Abschnitts ab. Das ist die schwerste Befund-Klasse dieses
Repos — ein Gate, das grün meldet, ohne geprüft zu haben.

**Der Defekt ist älter als die Closure-Fähigkeit.** Der Fence-Automat toggelt
naiv; [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) hat den
längenabgeglichenen Fence-Schluss ausdrücklich offen gelassen („naiver-Toggle-vs-
strikter-Schluss bewusst offen"). Was damals eine vertretbare Grenze war, ist mit
einer Bedingung, die **innerhalb** eines Abschnitts misst, zu einem Silent-Grün
geworden.

## 3. Bestandsmessung (Abnahme-Punkt 2, erledigt)

Gemessen am 2026-08-09 mit d-checks **eigener** `FenceToggle`-Lexik über drei
Repos — das eigene und die zwei, die die offenen Change Requests gestellt haben:

| Repo | Markdown-Dateien | ungerade Fence-Zahl | gemischte Fence-Längen | `~~~`-Fences |
|---|---|---|---|---|
| d-check | 347 | 0 | 0 | 0 |
| a-check | 184 | 0 | 0 | 0 |
| ai-harness-course | 245 | 0 | 0 | 0 |
| **Summe** | **776** | **0** | **0** | **0** |

Drei Aussagen folgen daraus, und sie drehen die Ausgangslage:

1. **Der Defekt ist latent, nicht aktiv.** Kein einziges Dokument im Ökosystem
   löst ihn heute aus. Der Reproduktionsfall in §2 ist konstruiert — was ihn
   nicht harmloser macht, aber anders einordnet.
2. **Die befürchtete Reichweite von Variante (b) ist empirisch widerlegt.** Der
   längenabgeglichene CommonMark-Schluss hätte auf diesen 776 Dateien **null**
   Wirkung: es gibt weder gemischte Fence-Längen noch `~~~`-Fences, also auch
   keine Fehlpaarung, die er korrigieren könnte. Er wäre kein Aufräum-Projekt —
   er wäre ein No-op mit Vertragsfläche.
3. **Das eigentliche Risiko ist nicht der Bestand, sondern der nächste Autor.**
   Wer künftig einen Fence nicht schließt, schaltet still die Prüfung des
   restlichen Abschnitts ab. Gesucht ist also keine bessere Paarung, sondern
   dass der unbalancierte Zustand **überhaupt gemeldet** wird.

## 4. Abnahme-Punkte

1. **Wie weit geht der Fix?** → **Entschieden 2026-08-09: (c)** — der neue
   Grund-Code. Die Messung hat die drei Kandidaten neu geordnet:
   - **(a) nur die Closure-/Struktur-Bedingungen** behandeln den offenen Fence
     lokal — deckt den belegten Fall, lässt die Klasse aber in jedem anderen
     Modul stehen.
   - **(b) längenabgeglichener Fence-Schluss im geteilten Automaten** — die
     Wurzel nach CommonMark, aber laut Messung ohne jede Wirkung auf den
     realen Bestand. Löst zudem **nicht** den belegten Fall: ein Fence, der
     **gar nicht** geschlossen wird, bleibt auch mit strengerer Paarungsregel
     offen.
   - **(c) ein eigener Grund-Code für den unbalancierten Fence** — macht den
     stillen Zustand laut. Deckt den belegten Fall direkt und kostet auf dem
     gemessenen Bestand nichts.
   (b) ist nach der Messung **nicht** die Wurzel dieses Befundes, sondern eine
   verwandte, folgenlose Frage — und die von
   [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) offen gelassene
   Grenze bleibt bewusst offen, jetzt mit einem Wächter davor.
2. **Wo wohnt der neue Befund?** → **Entschieden 2026-08-09: im Modul `spans`**,
   als Erweiterung von
   [`DC-FA-SPAN-001`](../../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
   — dieselbe Frageform, eine Ebene höher, also kein neues Kürzel
   ([ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)-Kriterium).
   Nicht im `planning`-Modul: der offene Fence ist
   kein Planning-Thema, sondern ein **Markdown-Artefakt**. Das Modul `spans`
   sagt genau das zu — „ungeschlossene Code-Spans … kippen die Backtick-Parität
   des restlichen Absatzes". Ein unbalancierter **Fence** ist dieselbe Aussage
   eine Ebene höher: eine Öffnung ohne Schluss, die alles Folgende umdeutet.
   Ein Befund im `planning`-Modul wäre näher am Fundort, aber falsch
   einsortiert — der nächste Konsument derselben Lexik fände ihn dort nicht.
3. **Absatz oder Datei?** → **Entschieden 2026-08-09: Datei**, Befund an der
   **Öffnungszeile**, genau einer je Datei. `span-unclosed` misst je Absatz; für
   einen Fence ist das nicht übertragbar, weil er selbst eine Absatzgrenze
   **ist** — absatzweise gemessen wäre er per Definition nie ungeschlossen. Die
   Öffnungszeile ist der **nächstgelegene** Ort; dass sie nicht immer die
   Reparaturstelle trifft, hat erst die Umsetzung gezeigt (§6) — die Entscheidung
   für die Öffnungszeile bleibt davon unberührt, ihre Begründung ist korrigiert.

## 5. Definition of Done

- [x] **Bestandsmessung** (§3): 776 Dateien über drei Repos, null Vorkommen —
      der Defekt ist latent, und Variante (b) ist empirisch als folgenlos
      belegt.
- [x] Abnahme-Punkte 1–3 entschieden; Vertragsanpassung geliefert (Lastenheft
      0.52.0: dritte Artefakt-Klasse + zwei Akzeptanzkriterien; Spezifikation
      §[`DC-FA-SPAN-001.a`](../../../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
      Schritt 3) samt
      [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md) `Proposed`.
- [x] Der oben belegte Fall meldet; Test mutations-echt. End-to-End gegen das
      gebaute Image: `done/slice-001-x.md:7` mit `fence-unclosed` neben dem
      `closure-note-thin` — der Wächter macht sichtbar, was die Prüfung
      unsichtbar gemacht hatte.
- [x] **Unabhängiger Code-Review** (Frischkontext) — merge-blockierend mit
      2 HIGH, 3 MEDIUM, 2 LOW; alle nachvollzogen und behoben bzw. bewusst
      dokumentiert, siehe §6.
- [x] **Bestätigende Re-Review** (Frischkontext) — erneut blockierend: 1 HIGH,
      4 MEDIUM, 1 INFO. Der Kern war echt geheilt, die Heilung hörte aber an
      derselben Stelle zu früh auf wie die Erstfassung. Alle nachgezogen,
      siehe §6.
- [x] **Dritte Review-Runde** (Frischkontext) — Leitfrage war nicht „sind die
      Befunde abgehakt", sondern „ist die **Klasse** geschlossen". Antwort: für
      die Fence-Lexik **ja**, belegt über alle Aufrufstellen im Repo; die
      Invariante war aber an drei von fünf Konsumenten nicht assertiert (HIGH)
      und ließ sich einzeln zurückdrehen, ohne dass ein Test rot wurde. Jetzt
      hat jeder Konsument seine eigene. Zwei Befunde in **anderen** Lexiken
      (`citations`-Absatzbildung, `headingSection`) und einer auf der
      git-Revisions-Achse (`vcs`) gehören ausdrücklich **nicht** in diesen
      Slice — eigener Vertrag, eigener Folge-Slice.
- [x] **Mutations-Gegenprobe, zweiter Anlauf.** Sechs Rückbauten, alle rot.
      Der erste Anlauf war methodisch kaputt: er setzte nach jeder Mutation
      per `git checkout` zurück — also auf HEAD statt auf den Arbeitsstand,
      wodurch die folgenden Ersetzungen ins Leere griffen und aus dem falschen
      Grund rot wurden. Über eine Dateikopie wiederholt, fiel ein echter
      Testfehler auf: die naive Lesart konnte ersatzlos entfallen, ohne dass
      ein Test rot wurde — in allen Fällen waren **beide** Lesarten
      gleichzeitig offen. Der Fall, in dem nur die Parität kippt, fehlte.
- [ ] `make gates` + `make verify-closure-notes` grün; Release als **Minor**
      (d-check findet danach mehr).

## 6. Risiken / offene Punkte

- **Reichweite von Variante (b).** — **Ausgang: entfallen.** Die Messung zeigt
  null betroffene Dateien; die Änderung wäre auf dem realen Bestand wirkungslos
  und löst den belegten Fall ohnehin nicht.
- **Der Defekt ist ausgeliefert.** Bis zum Fix ist die Zusage der
  Closure-Note-Struktur schwächer als dokumentiert. — **Ausgang:** offen; die
  Messung entschärft ihn (null Vorkommen im Ökosystem), hebt ihn aber nicht auf.
- **Wer nur `planning` aktiviert, sieht den Befund nicht** — und, beim
  End-to-End-Beleg **geschärft**: auch mit aktivem `spans` nur, wenn die Datei im
  **Scan-Scope** liegt. Der Review hat die Grenze ein zweites Mal geweitet: neben
  Post-Pässen über selbst benannte Verzeichnisse fallen auch **Zieldateien**
  außerhalb der Scan-Wurzeln heraus, aus denen Module lesen (`matrix` den Status,
  `anchors` die Slugs, `diagrams` und `versions` ihre deklarierten Quellen).
  — **Ausgang: eingetreten und benannt** — die Grenze steht in der Anforderung
  und in [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md), nicht nur im
  Slice. Sie zu **schließen** ist eigene Arbeit und nicht Teil dieses Slice.
- **Der Wächter bewachte nur eine von zwei Lesarten** (Review, HIGH). d-check
  trägt zwei Schluss-Regeln; die Umsetzung wertete nur den naiven Toggle aus,
  sodass der Tabellen-Leser unbewacht blieb — belegt an einem
  Vollständigkeits-Gate, das Exit 0 über eine ungedeckte Anforderung hinter
  einem offenen Fence meldete. — **Ausgang: behoben.** Beide Lesarten werden
  ausgewertet, und die strenge ist aus ihrem Konsumenten herausgezogen und
  geteilt, damit keine zweite Kopie danebenliegt.
- **Der Wächter trimmte anders als das Bewachte** (Review, HIGH): `TrimSpace`
  unicode-weit gegen `TrimLeft` auf Space/Tab. Eine mit U+00A0 eingerückte
  Fence-Zeile kippte damit **nur** die Parität des Wächters und machte ihn blind
  für den echten offenen Fence dahinter — null Befunde, Exit 0.
  — **Ausgang: behoben — aber erst im zweiten Anlauf.** Die Re-Review fand
  dieselbe Divergenz im Modul `planning` (HIGH): dort trimmten **beide**
  Fence-Stellen weiter unicode-weit, sodass der Anlassfall dieses Slice
  unverändert reproduzierbar blieb — Floskel hinter einer U+00A0-eingerückten
  Fence-Zeile, `planning` still grün **und** `spans` still grün. Ich hatte die
  Instanz geheilt, nicht die Klasse. Jetzt trägt die Trimmung ein geteiltes
  Prädikat, das alle vier Konsumenten rufen; beide `planning`-Stellen haben
  eine eigene Assertion.
- **Die Heilung selbst brachte eine Regression** (Re-Review, MEDIUM): das
  ersetzte `TrimSpace` hatte das CR einer CRLF-Zeile mit entfernt, `TrimRight`
  auf Space/Tab nicht — das rohe CR landete im Befund-Ziel und überschrieb im
  Terminal die Befundzeile. Der Erst-Report hatte CRLF ausdrücklich als in
  Ordnung geprüft; die Eigenschaft ging beim Heilen verloren.
  — **Ausgang: behoben**, mit Assertion.
- **Zusagen ohne Assertion** (Re-Review, MEDIUM): der Vorrang der strengen
  Lesart ließ sich tauschen, ohne dass ein Test rot wurde.
  — **Ausgang: behoben.**
- **Die gemeldete Zeile ist nicht immer die Reparaturstelle** (Review, LOW).
  Grundsätzlich nicht bestimmbar, welche von mehreren Öffnungen fehlt.
  — **Ausgang: nicht behebbar, benannt.** Anforderung, Spezifikation und ADR
  nennen sie **Fundstelle**. Die Re-Review zeigte, dass die erste Fassung dieser
  Klausel zu eng war: sie band die Einschränkung an die Paritätskippung, aber
  auch die strenge Lesart zeigt daneben, sobald eine längere Fence-Zeile eine
  kürzere Öffnung geschlossen hat. Jetzt gilt sie für beide.
- **Die Grenzen-Aufzählung war unvollständig** (Re-Review, MEDIUM): `codepaths`
  und `pins` fehlten. Bei `pins` ist die Richtung **still** — ein offener Fence
  verdeckt die Überschrift, der Anker wird unauflösbar, und die Drift-Prüfung
  entfällt kommentarlos statt zu melden. — **Ausgang: benannt**, samt dem
  eingerückten Code-Block, der in keiner Lesart modelliert ist.

## 7. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei. **Vor**
[slice-099](slice-099-structure-modul.md) sinnvoll — sonst erbt das neue Modul
den bekannten stillen Grün-Pfad über die geteilte Mechanik.

**Rückführungen:** `in-progress` → `open`, falls die Bestandsmessung zeigt, dass
Variante (b) eine eigene Sanierung nach sich zieht.

## 8. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**; andere
  Klasse, nichts zu berücksichtigen. **Kandidat für einen neuen Eintrag,** von
  den beiden Reviews geschärft: „eine geteilte Lexik driftet an den Rändern,
  weil jeder Konsument sie selbst vorbereitet". Nicht die Lexik-Grenze war das
  Problem, sondern dass vier Stellen dieselbe Zeile verschieden trimmten und
  drei Automaten dieselbe Frage verschieden beantworteten — jede für sich
  vertretbar, zusammen ein stiller Grün-Pfad. Bei der Closure einzutragen.

## 9. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die Korrektur wird zuerst als Zusage
formuliert (welches Verhalten gilt bei unbalanciertem Fence?), dann geliefert.
Kein Brownfield: es wird kein undokumentierter Bestand inventarisiert, sondern
eine dokumentierte Grenze neu bewertet.

## 10. Closure-Notiz (nach `done/`)

Geliefert ist der Grund-Code `fence-unclosed` als dritte Artefakt-Klasse des
Moduls `spans`, ausgeliefert mit **v0.53.0**
(Digest `sha256:0cbe2d54…9424`, Pipeline-Lauf 31362555433). Der belegte
Reproduktionsfall meldet: die Floskel hinter einem offenen Fence, die in
v0.52.0 still durchlief, erzeugt jetzt einen Befund an der Öffnungszeile.

**Die Bestandsmessung hat die Entscheidung gedreht.** Vor der Messung stand der
längenabgeglichene CommonMark-Schluss als „die Wurzel" im Plan. 776 Dateien über
drei Repos zeigten null unbalancierte Fences, null gemischte Fence-Längen, null
Tilde-Fences — er wäre wirkungslos gewesen. Und beim Aufschreiben fiel auf, dass
er den belegten Fall **gar nicht löst**: er korrigiert Fehlpaarungen, aber ein
Fence, der nie geschlossen wird, bleibt unter jeder Paarungsregel offen. Ich
hatte ihn als Wurzel dieses Befundes bezeichnet; er war die Wurzel eines anderen.

**Drei Review-Runden, und jede fand dieselbe Klasse an einer neuen Stelle.**
Das ist die eigentliche Lehre dieses Slice. Runde eins: der Wächter wertete nur
eine der zwei Schluss-Lesarten des Produkts aus, sodass der Tabellen-Leser
unbewacht blieb — belegt an einem Vollständigkeits-Gate mit Exit 0 über einer
ungedeckten Anforderung. Runde zwei: dieselbe Trim-Divergenz stand unverändert
im Modul `planning`, also genau dort, wo der Anlassfall gemessen wird; ich hatte
die Instanz geheilt, nicht die Klasse. Runde drei: die Invariante war an drei von
fünf Konsumenten nicht assertiert und ließ sich einzeln zurückdrehen, ohne dass
ein Test rot wurde.

Erst danach war die Fence-Lexik als Klasse geschlossen — ein geteiltes
Trimm-Prädikat, beide Schluss-Lesarten geteilt und ausgewertet, je Konsument eine
eigene Assertion gegen Wieder-Divergenz.

**Die Mutations-Gegenprobe war beim ersten Anlauf selbst kaputt.** Sie setzte
nach jeder Mutation per `git checkout` zurück — also auf HEAD statt auf den
Arbeitsstand. Die folgenden Ersetzungen griffen ins Leere und wurden aus dem
falschen Grund rot; nebenbei war die Implementierung im Arbeitsbaum gelöscht.
Über eine Dateikopie wiederholt, mit Abbruch bei nicht greifender Ersetzung,
förderte sie sofort einen echten Testfehler zutage. Sechs Rückbauten am Ende,
alle rot.

**Zwei Dinge wurden ehrlicher statt repariert.** Die gemeldete Zeile heißt
**Fundstelle**: welche von mehreren Öffnungen fehlt, ist grundsätzlich nicht
entscheidbar. Und die Grenze wuchs über drei Runden von „wer nur `planning`
aktiviert" über die Ziel-Achse bis zu `pins`, wo die Folge **still** ist — ein
offener Fence verdeckt die Überschrift, der Anker wird unauflösbar, und die
Drift-Prüfung entfällt kommentarlos. Sie zu schließen ist eigene Arbeit.

**Bewusst nicht hier erledigt:** dieselbe Klasse in **anderen** Lexiken —
`citations` bildet Absätze selbst, `headingSection` beantwortet die Anker-Frage
roh, `vcs` rechnet auf git-Blobs, die kein scannendes Modul sieht. Eigene
Verträge, eigener Slice: [slice-103](../open/slice-103-geteilte-lexik-raender.md),
mit Bestandsmessung als erstem Abnahme-Punkt — die Lehre dieses Slice.

**Fürs Register:** die Klasse „eine geteilte Lexik driftet an den Rändern, weil
jeder Konsument sie selbst vorbereitet" ist als **BEO-003** eingetragen. Der
`vcs`-Befund aus Runde drei ist zugleich die **dritte** Wiederholung von
„Modul-Grenze nur auf der Quell-Achse gedacht" (**BEO-004**) — damit ist die
Verkörperungs-Schwelle erreicht.
