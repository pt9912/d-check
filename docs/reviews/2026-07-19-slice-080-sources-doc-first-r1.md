# Review-Report — slice-080 (Modul `sources`, Upstream-Content-Drift externer Quellen), doc-first-Stratum

| Feld | Wert |
|---|---|
| **Gegenstand** | slice-080 doc-first (kein Code): `spec/lastenheft.md` (`DC-FA-SRC-001`, Bereich `SRC` §3, `DC-QA-03`-Erweiterung, §7-Zeile 0.49.0, Glossar/`DC-FA-CLI-002`), `spec/spezifikation.md` (`DC-FA-SRC-001.a` + §2-Schema `sources[]` + §7-Zeile), `docs/plan/adr/0046-sources-upstream-content-drift.md` (NEU, `Proposed`), `docs/plan/planning/in-progress/slice-080-sources-modul.md` (NEU), `…/roadmap.md` (Wellen-Flip) |
| **Commit** | `b36dc58` |
| **Datum** | 2026-07-19 |
| **Reviewer** | R1-doc, unabhängig/adversarial (Subagent), kontext-getrennt |
| **Rolle** | Reviewer, **nicht** Verifier (kein DoD-Abhaken, kein Gate-Lauf-Bestätigen) |
| **Vergleichsmuster** | `DC-FA-EXT-001` (Netz-Vorbild), `DC-FA-PIN-001` (Content-Pin-Geschwister), `MR-022`/`fetch-baseline-cache.sh` (Tooling-Vorläufer) |

## Findings

### F-1 · MEDIUM · `DC-QA-02` (Determinismus) / Manifest-Vertrag · Spec ↔ Lastenheft ↔ ADR ↔ Vorläufer-Tooling
- **pfad:** `spec/spezifikation.md:1746-1748` (Schritt 4) · `spec/lastenheft.md:2011-2015` · `docs/plan/adr/0046-sources-upstream-content-drift.md` (Entscheidung „Zwei Quelltypen", Fitness Function)
- **befund:** Schritt 4 verlangt gleichzeitig zwei **unvereinbare** Eigenschaften des Content-Manifests: (a) die Zeilen `"<hex>  <name>\n"` werden `LC_ALL=C`-**sortiert** (→ Sortierschlüssel ist die **ganze Zeile**, also `<hex>` zuerst) und (b) das Ergebnis sei „byte-Muster-gleich zu `sha256sum regelwerk/*.md`". Bare `sha256sum regelwerk/*.md` gibt in **Datei-namen-Reihenfolge** (Glob) aus, **unsortiert** — genau so schreibt der Vorläufer die vendored Datei (`tools/harness/fetch-baseline-cache.sh:175`: `sha256sum regelwerk/*.md > SHA256SUMS`, kein `sort`). Das Skript sortiert **beide** Seiten erst zur Vergleichszeit (`:114` `sha256sum … | LC_ALL=C sort`, `:115` `LC_ALL=C sort < "$sums"`). Zeilen-Sortierung (hash-first) und Dateinamen-Reihenfolge fallen nur zusammen, wenn die Inhalts-Hashes zufällig in Dateinamen-Ordnung aufsteigen — im Allgemeinen falsch. Das Modul hasht **ein** kanonisches Manifest und muss sich für (a) **oder** (b) entscheiden; die Vorlage `SHA256SUMS` liegt unsortiert vor, die Reihenfolge-Invarianz (AK „Archiv-Determinismus") verlangt (a). Die Parität-Behauptung zu bare `sha256sum regelwerk/*.md` bzw. „spiegelt exakt das … `SHA256SUMS`-Muster" ist damit im Regelfall unwahr.
- **verifizierbar:** ja — ein Go-Test mit einem 2-Datei-Bundle, dessen Inhalts-Hash-Ordnung ≠ Dateinamen-Ordnung ist, zeigt Manifest-Bytes ≠ `sha256sum <glob>`-Bytes; die AK „Archiv-Determinismus" prüft nur Reorder-Invarianz und fängt die falsche Parität **nicht**.

### F-2 · MEDIUM · Spec-`.a`-Implementierbarkeit (Interpretationslücke im Kern-Algorithmus) · `DC-QA-02`
- **pfad:** `spec/spezifikation.md:1745-1746` (Schritt 4, `<name>`)
- **befund:** Der Manifest-Eintrag lautet `"<hex>  <name>\n"`, aber `<name>` ist **nicht definiert**. Für ein Zip mit **verschachtelten Verzeichnissen** ist offen, ob `<name>` der volle Zip-interne Pfad (`regelwerk/foo.md`), der Basisname (`foo.md`), mit/ohne führendes `./`, mit `/`- oder OS-Trenner ist. Der Manifest-Hash — die eigentliche Produkt-Ausgabe — hängt direkt davon ab: zwei konforme Implementierungen können bei identischem Zip-Inhalt **verschiedene** (je reihenfolge-invariante) Hashes liefern, ohne dass ein Akzeptanzkriterium sie trennt. Zusätzlich unbehandelt: zwei Einträge mit gleichem Basisnamen in verschiedenen Zip-Verzeichnissen (Kollision/Tie-Break bei Basisnamen-Wahl). Der Algorithmus ist an genau dem Punkt, der den Hash bestimmt, nicht ohne Raten codierbar.
- **verifizierbar:** ja — Determinismus-/Paritätstest gegen ein Zip mit `unterordner/datei.md`; ohne festgelegte `<name>`-Form ist das erwartete Manifest nicht eindeutig ableitbar.

### F-3 · MEDIUM · `DC-QA-02` (Determinismus) / fail-closed-Abgrenzung · fehlender Negativfall
- **pfad:** `spec/spezifikation.md:1733-1735` (Schritt 2) · `:1739-1748` (Schritt 3/4) · `spec/lastenheft.md:2031-2033`
- **befund:** Die `sha256`-Konstante (`sources[].sha256`, `spec/spezifikation.md:1953`: „genau 64 Hex-Zeichen") legt die **Groß-/Kleinschreibung nicht fest**; der errechnete Ist-Hash (Go `hex.EncodeToString`) ist stets **Kleinbuchstaben**. Ob der Vergleich (Schritt 5) case-normalisiert, steht nirgends. Ein wohlgeformter Pin in Großbuchstaben (64 gültige Hex-Zeichen) erzeugt damit einen **Falsch-`source-drift`**, und die als „Re-Pin-Vorlage" emittierte Ist-Hash-Form (Groß/Klein) ist ebenfalls unbestimmt — jeder Copy-Re-Pin churnt zwischen den Schreibweisen. Kein Akzeptanzkriterium deckt den Fall.
- **verifizierbar:** ja — ein Test mit korrektem Inhalt und großgeschriebenem Pin muss grün sein; ohne Normalisierungs-Klausel schlägt er als `source-drift` fehl.

### F-4 · MEDIUM · fail-closed-Vollständigkeit / Robustheit · fehlender Ausgang bei nicht-parsebarer Antwort
- **pfad:** `spec/spezifikation.md:1733-1748` (Schritt 2-4) · `spec/lastenheft.md:2031-2033`
- **befund:** fail-closed ist ausschließlich für **Direktive/Config-Fehler** (kein `sha256:<hex>`, fehlende `url`/`sha256`, unbekanntes `unpack`) auf **Exit 2** festgelegt. Der Fetch gelingt (HTTP 200), aber der Body ist bei `unpack: zip` **kein gültiges Zip** (Nicht-Zip an einer `.zip`-URL, HTML-Fehlerseite mit Status 200, abgeschnittenes Archiv) — dieser Ausgang hat **keinen** definierten Befund: `source-drift`? `source-unreachable`? Exit 2? Schritt 3 fängt nur Status ≥ 400/Timeout/Netzfehler, Schritt 4 setzt ein lesbares Zip voraus. Eine ganze reale Fehlerklasse ist unspezifiziert; der Implementierer muss den Grund-Code raten.
- **verifizierbar:** ja — Test gegen einen `httptest`-Server, der bei `unpack: zip` `200 text/html` liefert; der erwartete Befund/Exit ist aus dem Vertrag nicht ableitbar.

### F-5 · MEDIUM · Sicherheit/Robustheit (Fetch beliebiger URLs, Zip-Entpacken) · kein Größen-/Dekompressions-Limit
- **pfad:** `docs/plan/adr/0046-sources-upstream-content-drift.md` (Konsequenzen: „Netz-Download + Zip-Entpacken") · `spec/spezifikation.md:1739-1745` (Schritt 3/4) · `spec/lastenheft.md:2038-2045` (AKs)
- **befund:** `sources` lädt beliebige gepinnte `http(s)`-URLs **voll** herunter (Hash über den ganzen Body) und entpackt bei `unpack: zip`. Anders als `external` (nur HEAD/GET-Erreichbarkeit) muss `sources` die vollständige Antwort in den Speicher ziehen und dekomprimieren. Es ist **kein** Maximal-Antwort-/Dekompressions-Limit definiert — eine sehr große Antwort oder eine **Zip-Bombe** (hohe Dekompressionsrate) kann den Prozess OOM/aushungern lassen. Für ein opt-in-Netz-Modul, das per Design fremde Inhalte holt, fehlt die Verteidigungslinie im Vertrag (kein AK, keine Schema-Grenze wie `external.timeout-seconds`).
- **verifizierbar:** ja — ein `httptest`-Server, der ein hoch-komprimiertes Zip oder einen unbounded Body liefert, belegt fehlende Begrenzung; heute keine Zusicherung, die ein Test einfordern könnte.

### F-6 · MEDIUM · Akzeptanzkriterien-Vollständigkeit (Boundary/„kein Doppelbefund") · Vergleich `DC-FA-PIN-001`
- **pfad:** `spec/lastenheft.md:2017-2020` (Prosa „Nur externe http(s)-Ziele, kein Doppelbefund") · `:2038-2045` (Akzeptanzkriterien) · `spec/spezifikation.md:1727-1731` (Marker-Bindung)
- **befund:** Zwei im Körper zugesicherte Verhalten haben **kein** Akzeptanzkriterium: (1) der „kein Doppelbefund"/repo-interne-Ziel-Fall — ein `source-pin`-Marker an einem repo-internen Link bleibt `pins`/`links`-Domäne; (2) die **Marker-Bindung** (mehrere Links je Zeile, ein allein auf der Folgezeile stehender Marker ist inert). `DC-FA-PIN-001` verankert beide als eigene AKs (`spec/lastenheft.md:1548-1549`, „Boundary (Marker-Bindung)"/„Boundary (Ziel weg)"); `sources` übernimmt das Muster „wie `pins`" in Prosa, ohne die Prüfbarkeit nachzuziehen. Ein Verifier, der nur das AK-Set abprüft, testet die Bindungs-/Doppelbefund-Regel nicht.
- **verifizierbar:** ja — die fehlenden Boundary-Tests wären genau die, die eine Bindungs-/Skopus-Regression fangen; ihr Fehlen ist im AK-Block sichtbar.

### F-7 · LOW · Befund-Schema-Implementierbarkeit (Config-Pin `line`) · Novum „Befund zeigt in `.d-check.yml`"
- **pfad:** `spec/spezifikation.md:1750-1751` (Schritt 5, „die `.d-check.yml`-Herkunft") · Befund-Schema `:1770-1771` (`line integer ≥ 1`)
- **befund:** Für einen Config-Pin sagt Schritt 5 `file`/`line` = „die `.d-check.yml`-Herkunft", ohne festzulegen, **welche** Zeile (Eintrags-Beginn? `url`-Zeile?). Das Befund-Schema fordert `line ≥ 1` (kein Platzhalter 0), und `DC-QA-02` fordert Determinismus. Zudem wäre dies der **erste** Befund, dessen `file` in die Config zeigt (alle bestehenden config-getriebenen Module — `matrix`/`hostpaths`/`external` — melden am **Markdown**-Vorkommen); es setzt YAML-Node-Zeilenverfolgung voraus, die die heutige Struktur-Dekodierung evtl. nicht führt. Unterspezifiziert, aber implementierbar.
- **verifizierbar:** ja — ein Determinismus-/Golden-Test auf `file`/`line` eines Config-Pin-Befunds legt die heute offene Zeilenwahl bloß.

### F-8 · LOW · Konsistenz Marker ↔ Config (Keyword-Asymmetrie) · Zukunfts-Reibung
- **pfad:** `spec/lastenheft.md:2002-2005` (Marker `archive`) vs. `:2004-2005`/`spec/spezifikation.md:1954` (Config `unpack: none|zip`)
- **befund:** Für dieselbe Semantik „Archiv vs. Einzeldatei" nennt der Marker ein **präsenzbasiertes**, format-implizites Keyword `archive`, die Config ein **format-benanntes** `unpack: zip`. Der ADR-Re-Evaluierungs-Trigger nennt `tar`/`gz` als spätere Anforderung — dann ist `unpack: tar` sauber erweiterbar, das Marker-Keyword `archive` aber format-mehrdeutig und syntaktisch anzufassen. Kleine, aber echte Asymmetrie zweier Deklarations-Flächen desselben Vertrags.
- **verifizierbar:** nein (Design-/Konsistenz-Beobachtung; kein Verhaltens-Gate).

### F-9 · INFO · fail-closed-Asymmetrie (stiller inerter Marker) · `DC-FA-SRC-001`
- **pfad:** `spec/spezifikation.md:1736-1738` (Schritt 2)
- **befund:** Ein **wohlgeformter** `source-pin: sha256:<hex>` an einem **repo-internen** Link ist „kein `sources`-Kandidat" → **still inert** (kein Befund, kein Hinweis), während eine **malformte** Direktive (kein `sha256`) fail-closed auf Exit 2 geht. Wer den `sources`-Marker (statt `dpin`) am falschen Zieltyp setzt, erhält keinerlei Diagnose — die Drift-Prüfung, die er zu aktivieren glaubte, läuft schlicht nicht. Bewusste Abgrenzung (`pins`/`links`-Domäne), aber als undokumentierte Ergonomie-Falle festzuhalten.
- **verifizierbar:** nein (Design-Notiz; kein deterministischer Fehl-Ausgang).

### F-10 · INFO · Redirect-Politik für `sources` unspezifiziert · Vergleich `DC-FA-EXT-001`
- **pfad:** `spec/spezifikation.md:1739-1740` (Schritt 3, „GET; Timeout/Parallelität wie `external`")
- **befund:** Schritt 3 übernimmt von `external` **nur** Timeout/Parallelität. `external` folgt bis zu fünf Redirects (`spec/lastenheft.md:1207-1208`); ob `sources` Redirects folgt — und ob ein Redirect auf einen anderen Host für eine **inhalts-gepinnte** Quelle zulässig ist (Integritäts-Frage: der Inhalt kommt dann von anderswo) — ist offen. Bei einer Redirect-Kette ist unklar, ob `source-unreachable`, Folgen, oder Content-vom-Ziel-Host greift.
- **verifizierbar:** nein (im Vertrag heute nicht entscheidbar; Klärung nötig, nicht falsifizierbar).

## Negativbefunde (geprüft, ohne Befund)

- **Grund-Code-Vertagung (§4-Lockstep):** `spec/spezifikation.md` §4 listet `source-drift`/`source-unreachable` korrekt **noch nicht** (geprüft: die vorhandenen Zeilen decken `external-status`, `link-stale`, `citation-*`, `gate-phantom` etc.). Die Vertagungs-Formel „landen mit der Modul-Implementierung … (AllReasons-↔-§4-Lockstep)" (`:1760-1762`) ist wortgleich zum Präzedenzmuster von `targets` (`gate-phantom`/`gate-undocumented`, Lastenheft-Historie 0.38.0) und `citations`. Ein voreiliger §4-Eintrag ließe den Lockstep-Test rot laufen — die Vertagung ist korrekt und konsistent.
- **`DC-QA-03`-Amendment (Anforderung + Messmethode):** `spec/lastenheft.md:2144-2145` führt `sources` sauber als zweite Netz-Ausnahme („außer … `external` und `sources`"; Messmethode: „alle Module außer `external`, `sources` und `vcs` aktiv"; „`sources` fetcht nur die explizit gepinnten Quellen, nie den Markdown-Baum"). Der Wortlaut hält.
- **Keine hängengebliebenen „einzige Netz-Tür"-Behauptungen in Spec/Lastenheft:** Grep über beide Dateien findet nur die **bereits revidierten** Stellen `spec/spezifikation.md:1720` und `spec/lastenheft.md:2027` (jeweils „neben `external` — die einzige Netz-Tür"). Die noch offenen `README`/`operations.md`/`config_template.go`-Stellen sind laut Slice-DoD legitim **Code-Phase** — kein doc-first-Rückstand.
- **Stratifizierung (MR-006 / SDP):** ADR-0046-Körper referenziert nur aufwärts/horizontal (`DC-FA-*`, `DC-QA-*`, `MR-019`/`MR-022`, `ADR-0005`/`ADR-0012`); der einzige `slice-080`-Token steht in der Sektion **Geschichte** (`:111`), die `matrix.exclude-sections` (`Geschichte`) vom Token-Scan ausnimmt — kein `matrix-forbidden`-Risiko, reine Provenance (SDP-konform). Spec-Körper (`DC-FA-SRC-001.a`) nennt keine ADR/Slice; die `slice-080`/„begleitende ADR"-Erwähnung steht nur in §7-Historie (matrix-exkludiert, hausüblich). Lastenheft-Körper verweist nicht auf Spec-Anker/ADR-Nummern (nur generisch „in der Spezifikation"/„begleitender ADR" in der Historie).
- **Exit-Code-Konsistenz:** `source-drift`/`source-unreachable` → Exit 1, malform/Config → Exit 2 (`spec/lastenheft.md:2043-2045`) stimmen mit `pins`/`external`/`DC-FA-CLI-003` überein; kein `--repair`-Hunk (wie `link-stale`).
- **Roadmap-/Bereichs-/Glossar-Einträge:** `welle-63-sources`-Flip in `roadmap.md` samt Historien-Zeile, Bereich `SRC` in §3 (`spec/lastenheft.md:57`), `sources` in `DC-FA-CLI-002` (`:89`) und Glossar (`:2166`) sind vorhanden und konsistent.
- **Determinismus-Kern (bei geschlossenem Manifest-Vertrag):** Der Manifest-statt-Zip-Roh-Bytes-Ansatz ist grundsätzlich `DC-QA-02`-tauglich und reorder-invariant — die offenen Punkte (F-1/F-2) betreffen die **byte-genaue** Kanonisierung, nicht das Grundprinzip.

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 6 | F-1, F-2, F-3, F-4, F-5, F-6 |
| LOW | 2 | F-7, F-8 |
| INFO | 2 | F-9, F-10 |

## Verdikt — BLOCK (eng, auf F-1 + F-2)

Die Anforderung ist gut motiviert, sauber stratifiziert und im Netz-/Grund-Code-Design konsistent mit `external`/`pins`; das `DC-QA-03`-Amendment ist im doc-first-Stratum vollständig vollzogen. **Blockierend** ist jedoch der Kern des einzig wirklich neuen Mechanismus — das **Content-Manifest**: F-2 lässt die den Hash bestimmende `<name>`-Form offen (zwei konforme Implementierungen können bei gleichem Inhalt verschieden hashen), und F-1 stellt zwei **unvereinbare** Manifest-Eigenschaften nebeneinander (zeilen-sortiert/reorder-invariant **vs.** byte-gleich zum unsortierten `sha256sum regelwerk/*.md`/vendored `SHA256SUMS`), belegt durch den Vorläufer, der die Vorlage unsortiert schreibt und erst beim Vergleich beidseitig sortiert. Der Algorithmus ist an genau dem Punkt, der die Produkt-Ausgabe definiert, **nicht ohne Interpretationslücke** codierbar — die AK „Archiv-Determinismus" prüft nur Reorder-Invarianz und würde die Divergenz nicht fangen.

Vor Code sind zusätzlich zu klären (MEDIUM): die `sha256`-Case-Normalisierung (F-3, Falsch-Drift-Risiko), der Ausgang bei nicht-parsebarer Archiv-Antwort (F-4), ein Größen-/Dekompressions-Limit (F-5, Sicherheit) und die fehlenden Boundary-AKs für „kein Doppelbefund"/Marker-Bindung (F-6). F-7/F-8 sind vor der Umsetzung mitzuerledigen, F-9/F-10 sind festzuhalten. Der „beides × beides"-Zuschnitt selbst ist im ADR real begründet (je ein Nutzungsmuster) und **nicht** überdimensioniert — die Kosten liegen nicht im Scope, sondern in der noch offenen byte-genauen Kanonisierung.
