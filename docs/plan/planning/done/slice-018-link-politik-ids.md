# Slice slice-018: Konfigurierbare Link-Politik für `ids` (`link-policy`)

**Status:** done.

**Welle:** welle-08-linkpolitik (per Roadmap-Fortschreibung; Start bei
Priorisierung durch den Auftraggeber — Ausgangsbefund im Dialog
2026-06-13: ein verlinktes [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) neben einem nur-Code-Span
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) in derselben DoD-Liste, vom grünen Gate nicht gefunden).

**Bezug:** [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(Change Request — neue konfigurierbare `link-policy`),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Schema-Erweiterung, Vollvalidierung),
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(Zeilen-Marker `d-check:ignore` — Geltungsbereich auf `ids` erweitert),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Abwärtskompatibilität: ohne `link-policy` byte-identisch),
[`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
(Muster der Fleet-Entdeckung über die Schwester-Repos).

**Autor:** pt9912. **Datum:** 2026-06-13.

---

## 1. Ziel

Das Modul `ids` erzwingt heute nur „keine *versehentlich* nackte ID im
Fließtext"; eine ID in Backticks (Code-Span) ist per Design
linkpflichtfrei. Dadurch ist der Sensor **blind** dafür, ob ein
Code-Span eigentlich ein Link sein sollte — „gut verlinkt" ist kein
gemessenes Property, sondern Glückssache menschlicher Aufmerksamkeit
(Ausgangsbefund: [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) in slice-017).

Ziel: „gut verlinkte Markdown-Dokumente" wird ein **im `.d-check.yml`
konfigurierbares, gemessenes Property**. Pro `ids`-Muster wählt das
Repo die Link-Politik: `prose` (heute, Default) oder `always` (auch
Code-Span-Vorkommen müssen im Linktext stehen). Zwei konfigurierbare
Ventile fangen die legitimen Nicht-Links, die die Kalibrierung (§6)
sichtbar gemacht hat.

## 2. Definition of Done

- [x] **Lastenheft-Change-Request** [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) (Version-Bump 0.8.0,
  Historie-Zeile): neue konfigurierbare `link-policy: prose|always`
  (Default `prose` = heutiges Verhalten, opt-in fürs Gating —
  Abwärtskompatibilität), `exempt-paths` (Glob-Liste pro Muster) und
  die Erweiterung des `d-check:ignore`-Zeilen-Markers von
  `codepaths`-only auf `ids` (gleiche Begründung wie `codepaths`:
  illustrative Beispiel-IDs); drei Akzeptanzkriterien
  (Happy/Boundary/Negative) + Out-of-Scope. Die Abwägung
  opt-in-Gating vs. Entdeckung ist im Anforderungstext festgehalten.
- [x] **Spezifikation** §[`DC-FA-ID-001.a`](../../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung) fortgeschrieben: `always`
  prüft zusätzlich Kennungs-Vorkommen *innerhalb* von Inline-Code-Spans
  (Wiederverwendung der `inlineSpansByLine`-Mechanik aus
  [`DC-FA-CODE-001.a`](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)); ein solches Vorkommen ist linkpflichtfrei nur,
  wenn der Code-Span im Linktext eines Markdown-Links liegt
  (`[` `` `ID` `` `](ziel)`), im `target` des Musters steht, in einer
  `exempt-paths`-Datei liegt oder die Zeile den `d-check:ignore`-Marker
  trägt. Schema-Tabelle §`.d-check.yml` um
  `ids.patterns[].link-policy` und `.exempt-paths`; **kein** neuer
  Grund-Code (weiterhin `id-unlinked`).
- [x] **Implementierung** im Modul `ids` + Config-Layer; die drei
  Akzeptanzkriterien als Tests; Config-Validierung (`link-policy` nur
  `prose|always`, sonst Exit 2; `exempt-paths`-Globs spiegeln die
  `scan.ignore`-Constraints).
- [x] **Abwärtskompatibilitäts-Beleg** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  Configs ohne `link-policy` verhalten sich byte-identisch (Eigenlauf
  und ≥1 Schwester-Repo vor/nach, identische Ausgabe-Hashes).
- [x] **Doku der Option** (Auftraggeber-Anspruch — die Option muss
  *sichtbar* sein, sonst kehrt „nie davon erfahren" durch die
  Hintertür zurück): nutzersichtbarer Abschnitt unter `docs/user/`
  („Linkdichte erzwingen — `link-policy: always`": dass es die Option
  gibt, wie man sie einschaltet, die Ventile `exempt-paths`/`ignore`,
  die Entdeckungs-vs.-Gating-Logik).
- [x] **Fleet-Entdeckungs-Pflicht** ([`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Muster):
  `always`-Lauf über alle `ids`-nutzenden Repos (heute d-check,
  u-boot, b-trace), Befunde pro Repo in der Closure-Notiz dokumentiert
  — Entdeckung hängt nicht am Opt-in des einzelnen Repos.
- [x] **Dogfooding:** d-checks eigene `.d-check.yml` setzt
  `link-policy: always` (mit `exempt-paths` für CHANGELOG + reviews);
  der `always`-Befundsatz (Kalibrierung: ~110 Kern-Treffer) wird zu
  echten Links bzw. begründeten `d-check:ignore`-Markern. Der bereits
  vorbereitete slice-017-Fix ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) verlinkt) wird hier gefaltet.
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz mit Steering-Loop-Lerneintrag (der Sensor-Blindfleck
  „grün ≠ gut verlinkt" → konfigurierbare Politik).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) | update | Change Request [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) (0.8.0): `link-policy`, `exempt-paths`, `ignore`-Erweiterung, AKs, Historie |
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | §[`DC-FA-ID-001.a`](../../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung) (`always`-Algorithmus, Ventile), Schema-Tabelle `.d-check.yml` |
| `internal/hexagon/core/rules/ids.go` | update | `always`: Vorkommen in Inline-Code-Spans prüfen, Linktext-/`exempt-paths`-/`ignore`-Ausnahmen |
| `internal/hexagon/core/` (Config-Typen/Parsing/Validierung) | update | `IDPattern.LinkPolicy`, `.ExemptPaths`; Validierung |
| `internal/hexagon/core/rules/ids_test.go` | update | drei AKs + Abwärtskompatibilität |
| `docs/user/` (neuer Abschnitt) | update/neu | Option sichtbar dokumentieren (Auftraggeber-Anspruch) |
| [`.d-check.yml`](../../../../.d-check.yml) | update | Dogfooding: `link-policy: always` + `exempt-paths` |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | 0.x-Eintrag |

## 4. Trigger

Priorisierung durch den Auftraggeber (Dialog 2026-06-13). Der
Ausgangsbefund (slice-017: [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) Code-Span neben verlinktem
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)) ist ein Steering-Loop-Signal („3× → Lücke"): der
`ids`-Sensor misst nicht das Ziel „gut verlinkt".

## 5. Closure-Trigger

DoD vollständig inkl. Fleet-Entdeckungs-Lauf, Dogfooding-Sweep grün
und nutzersichtbarer Doku der Option + Closure-Notiz mit Lerneintrag.

## 6. Kalibrierungs-Befund (Vorlauf, 2026-06-13)

Methode (ohne neuen Code, faithful): `ids`-Befunde auf Original vs. auf
einer Kopie mit neutralisierten Inline-Backticks (außerhalb Fences) —
die Differenz ist exakt der `always`-Befundsatz. Drei Repos nutzen
`ids` mit Mustern (die anderen haben es nur *nicht konfiguriert* — kein
Ausschluss).

| Repo | neue `always`-Befunde |
|---|---|
| d-check | 155 |
| u-boot | 9 |
| b-trace | 2 |

Triage der d-check-155: **39** literal-schwer (CHANGELOG 28 + Reviews
11 → `exempt-paths`), **6** Spec-Beispiel-IDs (`ADR-0042`/`ADR-0099` — <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
fiktive Illustrationen → `d-check:ignore`), **~110** Kern (echte
Inkonsistenzen wie AGENTS↔harness/README und Meta-Diskussion → Links).
Folgerungen, die die Regel-Form bestimmen:

1. Default `prose` (opt-in) — **nur** Abwärtskompatibilität für die
   anderen Repos, **kein** Urteil über Befundqualität (die ~110 sind
   echte Arbeit, kein Lärm).
2. `always` braucht konfigurierbare `exempt-paths` (39/155).
3. Spec-Beispiel-IDs sind eine irreduzible Klasse, die `exempt-paths`
   nicht lösen kann → Zeilen-Marker `d-check:ignore` für `ids`
   (konsistent mit `codepaths`' Beispiel-Pfad-Begründung).
4. Opt-in-Gating ohne Sichtbarkeit reproduziert den Blindfleck →
   Fleet-Entdeckungs-Pflicht + sichtbare Doku der Option.

## 7. Risiken und offene Punkte

- **Beispiel-ID-Ventil:** `d-check:ignore` auf `ids` zu erweitern
  weicht „für andere Module keine Suppression" auf. Begründung: gleiche
  Kategorie wie `codepaths` (illustrative Beispiele). Beim Review
  bestätigen oder durch Fenced-Code-Variante ersetzen.
- **Dogfooding-Sweep-Umfang** (~110): kann Folge-`id-unlinked`-Befunde
  in Spec-Beispiel-Kontexten erzeugen; iterativ wie slice-015/016
  („Kalibrierung in zwei Iterationen").
- **Determinismus:** `exempt-paths`-Glob-Auswertung muss
  reihenfolge-/plattformstabil sein ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default
der Modus-Tabelle).

## 9. Closure-Notiz (nach `done/`)

**Umsetzung:** Vertrag + Spezifikation (Commit `569c786`),
`ids`-Implementierung + Tests (`ba3deca`), Dogfooding-Sweep + Doku +
Closure (dieser Commit). `make gates` grün (Coverage 95,0 %).

**Fleet-Entdeckungs-Lauf (real, nicht Approximation):** `always` über
die drei `ids`-Repos — **d-check 170, u-boot 2, b-trace 2**. Die
Vorlauf-Approximation (§6: 155/9/2) überschätzte leicht (Backtick-Entfernen
verschmolz Schein-IDs); der span-basierte Echtlauf ist präziser. u-boot
und b-trace sind in *aktiver* Doku praktisch sauber (Treffer nur in
`done/`/CHANGELOG) — der Sensor lohnt sich, ohne zu lärmen.

**Dogfooding-Sweep (Auftraggeber-Entscheid: alle 77 `done/` fixen, kein
`done/`-Exempt):** 130 Befunde nach `exempt-paths` (CHANGELOG + Reviews)
aufgelöst — ~127 Inline-Code-Kennungen als Links ausgeführt
(Requirement-Kennungen → Lastenheft-Anker, `.a`-Sektions-Referenzen →
Spez-Sektionsanker), Beispiel-IDs (`ADR-0042`/`ADR-0099`) und zwei <!-- d-check:ignore (Beispiel-IDs, fiktiv) -->
literale Zitate mit `d-check:ignore`. Deterministischer Rewriter
(kanonische Ziele aus bestehenden Links geerntet, Slug-Fallback für
unverlinkte Sektionen); der `anchors`-/`matrix`-/`links`-Sensor
verifizierte jeden erzeugten Link (0 Befunde, keine verbotene
spec→adr-Kante). **Abwärtskompatibilität** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
b-trace ohne `always` byte-identisch zu v0.4.0.

- **Was hat funktioniert:** Den `always`-Befundsatz *aus dem echten Tool*
  zu sweepen und vom eigenen `anchors`-Sensor verifizieren zu lassen —
  ein falscher Anker hätte das Gate gerötet. Der Sweep musste nicht
  „auf Verdacht" korrekt sein; das Gate war das Netz.
- **Anders als geplant:** Die `.a`-Sektions-Referenzen waren eine
  eigene Klasse (Requirement-ID als Teilstring eines Sektions-Tokens) —
  der Rewriter wurde token-basiert mit getrennten Schlüsseln
  (`ID` vs. `ID.a`), sonst hätte er Sektions-Verweise auf
  Requirement-Anker fehlgeleitet.
- **Steering-Loop-Lerneintrag (geschärfte Regel + neuer Sensor):** Der
  `ids`-Sensor maß „keine *versehentlich* nackte ID", nicht das Ziel
  „gut verlinkt" — ein Code-Span verbarg stillschweigend fehlende
  Links (Blindfleck `grün ≠ gut verlinkt`). Geschärft: „gut verlinkt"
  ist jetzt ein *konfigurierbares, gemessenes* Property
  ([`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
  `link-policy: always`). Und: Opt-in fürs Gating, aber **Entdeckung
  als Bringschuld** (Fleet-Lauf) — „opt-in" darf nicht „unsichtbar"
  bedeuten.
- **Folge-Slices:** keine; u-boot/b-trace könnten `always` selbst
  aktivieren (ihre aktiven Doku-Befunde sind ~0).
