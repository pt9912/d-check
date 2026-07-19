# Review — slice-072 §4.12-Cluster-Auftrennung (uncommittete Arbeitsbaum-Änderung)

**Datum:** 2026-07-19 ·
**Auditor:** unabhängiger, adversarialer Reviewer (Reviewer-Skill v1.2.0) ·
**Gegenstand:** uncommittete Änderung an `docs/user/benutzerhandbuch.md`
(`git diff`, +130/-187 Z.) — Auftrennung des ~330-Z.-Monolithen §4.12 in vier
aufgabenorientierte Sektionen (§4.12 RTM/`--trace` · §4.13 `trace.coverage` ·
§4.14 `trace.requirements.modality` · §4.15 `trace.cross-consistency` · §4.16 =
alt §4.13 `--print-mk`); Referenz nach §5, „0 Anforderungen"-Fehlerbild nach §7,
inline-Versionsprosa entfernt, §5-/§11-Verweise nachgezogen ·
**Kriterium:** `docs/user/benutzerhandbuch-standard.md` §2/§5; slice-072-Plan §2
(„Referenz nach §5, Aufgabe nach §4 — kein Text zweimal"; „verifizierte Beispiele
bleiben verifiziert") ·
**Beleglage:** Diff und Vor-Stand (`git show HEAD:…`) zeilenweise verglichen;
Arbeitsbaum direkt gelesen; die beiden Test-Harnesse
(`handbook_examples_test.go`, `docexamples_test.go`) gelesen und ihre Invarianten
statisch (whitespace-normalisierte Substring-Prüfung, Fence-Extraktor +
Klassifikations-Replikation) sowie dynamisch (`make test`, Docker) belegt.

---

## Findings

### F-1 — §11-Changelog: drei Zeilen verweisen auf §4.12 für Regeln, die diese Änderung nach §5 verschoben hat

- **kategorie:** LOW
- **quelle:** Maintainability (Doku-Drift); slice-072-Plan §2 („Referenz gehört
  nach §5")
- **pfad:** `docs/user/benutzerhandbuch.md:1627` (Zeile 1.35), `:1628` (1.36),
  `:1629` (1.37)
- **befund:** Die Änderung hat die Sonderfälle der Tabellen-Grammatik
  (GFM-Trennzeile mit einem Bindestrich, Tabellengrenze am relevanten Header,
  Direktiven-Marker-Toleranz in Tabellenzeilen) aus §4.12 nach §5 verlagert
  (§4.12 Z. 618–620 nennt sie nur noch und verweist auf §5; die Regeln stehen in
  §5 Z. 1321–1329). Die Changelog-Zeilen 1.35/1.36/1.37 verweisen für **genau
  diese** Regeln weiterhin auf „§4.12". Im selben Diff wurden die parallelen
  Verweise für coverage/modality/cross/print-mk konsistent nachgezogen (1.25→§4.13,
  1.26→§4.14, 1.30→§4.15, print-mk-Zeilen→§4.16) — diese drei blieben zurück.
- **failure-szenario:** Ein Maintainer schlägt über Changelog-Zeile 1.36
  („Tabellengrenze am relevanten Header (§4.12)") nach, wo die Regel dokumentiert
  ist, springt nach §4.12 und findet dort nur den Weiterverweis „stehen in §5" —
  ein zusätzlicher Sprung; der genaueste Zeiger wäre §5, wie bei den anderen im
  selben Diff aktualisierten Zeilen.
- **verifizierbar:** nein — kein Gate prüft Prosa-„§4.x"-Verweise (der `anchors`-
  Gate greift nur bei Markdown-`](#slug)`-Links); manuell gegen §5 Z. 1321–1329
  zu bestätigen.

### F-2 — Vollständig entfallene Konsumenten-Empfehlung („Bestandsinvariante") ohne Ersatzort

- **kategorie:** INFO
- **quelle:** Content-Verlust-Achse (slice-072-Plan §2, §4)
- **pfad:** entfällt aus altem §4.12 (`HEAD:docs/user/benutzerhandbuch.md`, im
  `source: ""`-Fallback-Absatz)
- **befund:** Der HEAD-Satz „Für eine zusätzliche Bestandsinvariante kann ein
  Konsument die erkannte Gesamtzahl weiterhin gegen seine erwartete Zahl
  plausibilisieren." ist ersatzlos entfallen; die Zeichenketten
  `Bestandsinvariante`/`plausibilisieren`/`erwartete Zahl` kommen im Arbeitsbaum
  nirgends mehr vor. Anders als die vier ebenfalls verkürzten Verhaltensregeln
  (siehe Negativbefunde) hat dieser Hinweis **keinen** Zielort in §5/§7 bekommen.
- **failure-szenario:** Ein Konsument, der bisher dem dokumentierten Muster
  folgte (erkannte Gesamtzahl gegen eine erwartete plausibilisieren), findet die
  Empfehlung nicht mehr. Geringe Wirkung: es ist eine advisory Nutzungs-
  empfehlung, kein Tool-Verhalten; ihr Wegfall ist mit dem Trim-Ziel des Slice
  vereinbar.
- **verifizierbar:** nein — kein Gate; Vor-Stand-Vergleich belegt den Wegfall.

### F-3 — §4.12 trägt weiter einen Referenz-förmigen „Definitionssyntax"-Block vor „Vorgehen"

- **kategorie:** INFO
- **quelle:** Standard §5 (Reihenfolge Ausgangslage→Ziel→Vorgehen); slice-072
  DoD B-3
- **pfad:** `docs/user/benutzerhandbuch.md:595`–`620`
- **befund:** Zwischen **Voraussetzung** (Z. 592) und **Vorgehen** (Z. 622) steht
  weiterhin der Block „**Unterstützte Definitionssyntax:**" mit zwei
  Markdown-Beispielen (Heading-Anforderung positiv, Pipe-Tabelle negativ). Die
  Feld-/Grammatik-Referenz ist auf einen 3-Zeilen-Verweis nach §5 verkürzt
  (Z. 618–620) — der von B-3 kritisierte ~67-Z.-Block ist damit real erheblich
  entschärft —, aber der erklärende Syntax-Rahmen bleibt vor dem Vorgehen stehen.
- **failure-szenario:** Ein Leser, der nur die Handlungsfolge sucht, liest erst
  eine halbe Seite Definitions-Syntax. Vertretbar als „was der Leser einmal
  wissen muss (Anforderungen kommen aus Headings)", aber formal noch nicht die
  reine Ausgangslage→Ziel→Vorgehen-Öffnung wie in §4.13/§4.14.
- **verifizierbar:** nein — Ermessen gegen den Standard, kein Gate.

---

## Negativbefunde je Prüfachse (geprüft, ohne blockierenden Befund)

**Achse 1 — Content-Verlust bei der Verlagerung.** Die sieben in der Aufgabe
konkret benannten Verhaltensregeln sind erhalten (Vor-Stand vs. Arbeitsbaum
zeilenweise belegt):
- Escaped Pipes / Code-Span-Pipes „bleiben Teil derselben Zelle": erhalten in §5
  Z. 1289–1291 (generisch; der literale `\|`-Beispiel-Token entfiel, die Regel
  nicht).
- `table.modality-column` liefert exklusiv den Zelleninhalt: erhalten in §5
  Z. 1391–1393.
- Die zwei Nicht-Pipe-Migrationsalternativen: erhalten (verdichtet) in §5
  Z. 1315–1319.
- „`unknown` gatet nicht unter `require-levels: [must]`"-Warnung: erhalten in
  §4.14 Z. 716–719.
- Tabellengrenze-Regel „ein Header ohne gebundene Rolle beendet nicht": erhalten
  in §5 Z. 1323–1325 (die `| - | - |`-Illustration entfiel, die Regel nicht).
- Komma-Kurzform ⇒ Exit 2: erhalten in §5 Z. 1375–1379 (der Zusatz „Fehlermeldung
  nennt `..` und `/`" entfiel, die Regel nicht).
- Nullmengen-Guard-`source: ""`-Fallback: erhalten in §5 Z. 1331–1335 **und**
  §7 Z. 1542–1543.
Einziger vollständiger Wegfall ohne Ersatzort: die „Bestandsinvariante"-Empfehlung
(→ F-2, INFO). Die drei entfallenen Illustrationen/Detail-Phrasen (`\|`,
`| - | - |`, „Hinweis auf `..`/`/`") sind Beispiel-/Meldungs-Details, keine
verlorenen Verhaltensregeln.

**Achse 2 — „kein Text zweimal".** Keine wörtliche/nah-wörtliche Verdopplung
zwischen §4.13/§4.14/§4.15 und §5 gefunden: die Aufgaben tragen kurze
Ergebnis-Zusammenfassungen mit Zeiger „… stehen in §5", die Feldreferenzen leben
allein in §5. Grenzfall (kein Finding): §4.13 Z. 690–692 zählt die vier
fail-closed-Fälle als Vorschau auf, §5 Z. 1357–1358 führt sie als Referenz —
durch die explizite „stehen in §5"-Rahmung ein zulässiger Zeiger, keine
verdoppelte Prosa. Milde intra-§5-Redundanz (id-column/text-column-Pflicht in
Z. 1280–1281 und erneut Z. 1312–1315) liegt innerhalb der Referenz und außerhalb
der Slice-Regel.

**Achse 3 — Standard-Konformität der neuen Sektionen.** §4.13/§4.14 öffnen mit
**Ausgangslage → Ziel → Vorgehen → Ergebnis** aus der Lesersituation; §4.15 ist
die vom Slice §2 gebilligte slice-071-Schablone (Ausgangslage→Ziel→Schritt 1/2/3)
unverändert übernommen; §4.12 selbst öffnet ziel-zuerst (Z. 588–591). Kein
explizites „Voraussetzung" in §4.13/§4.14 — konsistent mit dem gebilligten
Vorbild §4.15 (das ebenfalls keines führt), daher kein Finding. Einzige Residuen:
§4.12 F-3 (INFO).

**Achse 4 — Test-Harness-Integrität.** `make test` (Docker) grün:
`…/driven/configyaml ok` und `…/driving/cli ok`.
- (a) Alle elf von `TestHandbook_TraceParsergrenzenDokumentiert` verlangten
  whitespace-normalisierten Teilstrings vorhanden; beide verbotenen Strings
  („**niemand** referenziert (Waise)", „mindestens eine Waise ⇒ Exit 1 statt 0")
  abwesend.
- (b) Die drei verankerten Ausgabeblöcke eindeutig und intakt: `# Requirements
  Traceability Matrix` (1×, §4.12 Z. 633), `ergab 0 Anforderungen` als `text`-
  Fence **genau 1×** (§7 Z. 1538; die zweite Fundstelle Z. 1528 ist Prosa-
  Überschrift, kein Fence — Uniqueness bleibt), `## Kreuzverweis-Konsistenz` (1×,
  §4.15 Z. 762). Alle `formTokens` in den verschobenen Blöcken erhalten. Eine
  Replikation des Fence-Extraktors + der `OutputBlocksClassified`-Logik ergab:
  10 gesweepte Blöcke, 0 unklassifiziert, der `--print-mk`-Block trägt weiter
  seinen `not-replayable`-Marker (Z. 819).
- (c) `--trace`/`--require-complete` weiter in `bash`-Blöcken dokumentiert
  (§4.12 Z. 624, §7 Z. 1533).
- (d) Alle achtzehn `yaml`-Blöcke sind Config-Domäne mit korrektem Marker-Status;
  `make test` validiert sie gegen `configyaml.Decode` grün.

**Achse 5 — Umnummerierungs-Vollständigkeit / Referenz-Korrektheit.** Die vier
§5-Prosa-Zeiger stimmen: §4.12 (Z. 1264), §4.13 (Z. 1347), §4.14 (Z. 1382),
§4.15 (Z. 1403). Die vom Auftrag benannten §11-Zeilen stimmen: print-mk→§4.16
(Z. 1596/1598/1611/1614), coverage→§4.13, modality→§4.14, cross→§4.15. **Kein**
Markdown-Anker-Link (`](#…)`) zeigt auf die umnummerierten Headings (nur der
unveränderte Eltern-Anker `#4-aufgaben`), der `anchors`-Gate bleibt also grün.
`docs/user/operations.md` referenziert `--trace`/`--print-mk`/`cross-consistency`
ausschließlich per Flag-/Schlüsselname, nicht per §4.1x-Nummer — keine Drift.
Ausnahme siehe F-1 (drei stale §4.12-Verweise in der Changelog-Historie, LOW).

**Achse 6 — §7-Fehlerbild-Korrektheit.** Der nach §7 verlagerte „0
Anforderungen"-Block (Z. 1528–1543) ist fachlich korrekt: `--trace
--require-complete` gegen eine explizit gesetzte, aber im Heading-Format leere
Quelle endet mit Exit 2; die gezeigte Meldung „… im Format headings ergab 0
Anforderungen" deckt sich mit der echten CLI-Ausgabe (durch das achte
Replay-Beispiel `tableOnlyRequirements` + `formTokens` E2E belegt). Ursache/Lösung
(„Quelldatei, `format` und `id-pattern` prüfen"; `source: ""`-Fallback) stimmen
mit dem Nullmengen-Guard aus §5 überein.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 1 | F-1 |
| INFO | 2 | F-2, F-3 |

---

## Verdikt

**ACCEPT.** Kein HIGH, kein MEDIUM — die typischen Blocker fehlen. Die
Kern-Invariante des Slice („verifizierte Beispiele bleiben verifiziert") ist durch
`make test` (Docker, grün) am schärfsten Sensor belegt; die Auftrennung folgt der
Ausgangslage→Ziel→Vorgehen→Ergebnis-Form, verlagert Referenz sauber nach §5 ohne
Text zu verdoppeln, und kein Verhalten ging verloren. Die drei stale
§4.12-Changelog-Verweise (F-1, LOW) sind reine Doku-Drift ohne Gate-Wirkung und
inkonsistent mit den im selben Diff korrekt nachgezogenen Nachbarzeilen — ein
optionaler Nachzug, kein Merge-Blocker. F-2/F-3 sind INFO.
