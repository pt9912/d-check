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

## Gegenprobe (R1-Nachtrag)

| Feld | Wert |
|---|---|
| **Gegenstand** | Nachprüfung der Einarbeitung von F-1…F-10, unabhängig gegen die Dateien auf Platte (frisch gelesen, nicht aus dem Gedächtnis) |
| **Commit** | `4e9404a` |
| **Datum** | 2026-07-19 |
| **Methode** | `spec/spezifikation.md` `DC-FA-SRC-001.a` Schritt 2-5 + §2-Schema, `spec/lastenheft.md` `DC-FA-SRC-001` + §7, `docs/plan/adr/0046-sources-upstream-content-drift.md`, `slice-080-sources-modul.md` zeilenweise gegengelesen |

### Befund-Status

- **F-1 (Manifest-Sortier-/Parität-Widerspruch) — GESCHLOSSEN.** `spec/spezifikation.md:1758-1778` (Schritt 4) sortiert jetzt **nach `<pfad>`** (byteweise, `LC_ALL=C`, ausdrücklich „**nicht** nach der ganzen Zeile") und streicht die falsche Byte-Parität: „folgt **konzeptionell** dem … `SHA256SUMS`-Muster … **nicht** byte-identisch zur unsortierten `sha256sum <glob>`-Ausgabe". Lastenheft (`:2010-2017`) und ADR (`:47-52`) tragen exakt dieselbe Aussage (pfad-sortiert, eigenständig kanonisiert, keine `sha256sum`-Byte-Parität). Der Sortier-nach-Zeile-vs-Parität-Widerspruch ist aufgelöst; lastenheft ↔ spec ↔ ADR konsistent.
- **F-2 (`<name>`-/`<pfad>`-Form undefiniert) — GESCHLOSSEN.** `spec/spezifikation.md:1763-1770`: `<pfad>` = **voller Zip-interner Pfad** normalisiert (Backslashes → `/`, führendes `./` und `/` entfernt, ausdrücklich **kein** Basisname → „verschachtelte Verzeichnisse bleiben im Pfad, daher keine Basisnamen-Kollision"); Verzeichnis-Einträge (Name endet auf `/`) raus; zwei Leerzeichen zwischen `<hex>` und `<pfad>`; je Zeile `\n`-terminiert. Der hash-bestimmende Punkt ist jetzt eindeutig; zwei konforme Implementierungen müssen denselben Manifest-Hash liefern.
- **F-3 (`sha256`-Case) — GESCHLOSSEN.** Schritt 2 („`sha256` **case-insensitiv** … beide zu Kleinbuchstaben normalisiert"), Schritt 5 (Vergleich case-insensitiv, Ist-Hash in Kleinbuchstaben als „schreibweisen-stabile" Re-Pin-Vorlage), §2-Constraint `spec/spezifikation.md:1986` („**case-insensitiv** verglichen") und Lastenheft („`sha256` wird **case-insensitiv** geführt"). Der Falsch-Drift durch Großschreibung ist ausgeschlossen.
- **F-4 (2xx-kein-Zip) — GESCHLOSSEN.** Schritt 3 nennt explizit „(bei `unpack: zip`) eine 2xx-Antwort, die **kein gültiges Zip** ist (Nicht-Zip / HTML-Fehlerseite / abgeschnittenes Archiv)" als `source-unreachable`; die Boundary-AK im Lastenheft ist entsprechend erweitert. Der zuvor offene Ausgang ist definiert.
- **F-5 (Größen-/Zip-Bomben-Limit) — GESCHLOSSEN.** Schritt 3 (Body **≤ 64 MiB**, > fünf Redirects) und Schritt 4 (Entpack-Gesamtgröße **≤ 256 MiB**, Eintragszahl **≤ 10 000**) → `source-unreachable`; deckungsgleich in ADR `:63-67` und Lastenheft. Die Verteidigungslinie steht im Vertrag.
- **F-6 (Boundary-AKs) — WEITGEHEND GESCHLOSSEN (Rest-Nit).** Lastenheft trägt neu „**Boundary (kein Doppelbefund / repo-intern)**" und „**Boundary (Marker-Bindung)**". Rest siehe F-12: die Marker-Bindungs-AK prüft nur den „kein vorausgehender `http(s)`-Link → inert"-Fall, nicht die Disambiguierung bei **mehreren** Links je Zeile (die pins-AK prüft beide).
- **F-7 (Config-Pin `line`) — GESCHLOSSEN.** Schritt 5: `file` = `.d-check.yml`, `line` = „die Zeile des `url`-Feldes des Eintrags (die YAML-Dekodierung führt sie; sonst `1`)". Zeilenwahl und Fallback sind deterministisch festgelegt.
- **F-8 (Marker-Keyword `archive`→`zip`) — NICHT GESCHLOSSEN, mit neuer Inkonsistenz — siehe F-11.**
- **F-9 (inerter Marker dokumentiert) — GESCHLOSSEN.** Schritt 2: „Ein **wohlgeformter** `source-pin` an einem repo-internen oder Nicht-`http(s)`-Link ist **inert** … die Direktive ist dort bewusst wirkungslos … **nicht** fail-closed." Die Ergonomie-Falle ist explizit.
- **F-10 (Redirect-Politik) — GESCHLOSSEN.** Schritt 3: „Redirects werden wie `external` bis zu **fünf** verfolgt, gehasht wird der Inhalt der **finalen** Antwort"; > fünf → `source-unreachable`.

### Neue / offene Findings

#### F-11 · MEDIUM · Konsistenz Lastenheft ↔ Spec (normative Direktiven-Syntax) · Regression aus dem F-8-Fix
- **pfad:** `spec/spezifikation.md:1729` (Schritt 1) · `docs/plan/adr/0046-sources-upstream-content-drift.md:36` — beide vs. `spec/lastenheft.md:2003-2004`, `docs/plan/adr/0046-…:53-56`, `slice-080-sources-modul.md:25`
- **befund:** Der F-8-Fix (Marker-Keyword `archive` → `zip`) ist nur teilweise propagiert. Lastenheft (`:2003-2004`: `<!-- source-pin: zip sha256:<hex> -->`), Slice (`:25`: `[zip]`) und die ADR-**Entscheidung** (`:53-56`: „das Archiv-Keyword ist explizit (Marker `zip`, parallel zum Config-Wert `unpack: zip`)") nennen `zip`. Aber die **normative** Spec `DC-FA-SRC-001.a` Schritt 1 (`:1729`) zeigt weiterhin `<!-- source-pin: [archive] sha256:<hex> -->`, und der ADR-**Marker-Beispiel-Punkt** (`:36`) ebenso — die ADR widerspricht damit **sich selbst** (`:36` `archive` vs `:53-56` `zip`). Der Marker-Parser wird aus Spec §.a Schritt 1 gebaut: er akzeptiert `archive`, während ein Nutzer laut Lastenheft `zip` schreibt — der Archiv-Pin würde als schlüsselwort-los (Einzeldatei-Hash eines Zip) gelesen und liefe auf Falsch-`source-drift`/`source-unreachable`.
- **verifizierbar:** ja — ein Marker-Parser-Test mit `source-pin: zip …` gegen die Spec-§.a-Fassung schlägt fehl (Keyword nicht erkannt); der Widerspruch ist ein direkter Datei-Diff `spec/lastenheft.md:2004` vs `spec/spezifikation.md:1729`.

#### F-12 · INFO · Akzeptanzkriterien-Restlücke (Marker-Bindung bei mehreren Links je Zeile) · `DC-FA-SRC-001`
- **pfad:** `spec/lastenheft.md` (AK „Boundary (Marker-Bindung)") · Vergleich `spec/lastenheft.md:1548` (pins-AK)
- **befund:** Die neue Marker-Bindungs-AK deckt nur „kein vorausgehender `http(s)`-Link → inert". Der von `spec/spezifikation.md:1727-1731` zugesicherte „wie `pins`"-Fall — mehrere Links je Zeile, jeder Marker bindet an den ihm **unmittelbar** vorausgehenden Link — hat weiterhin kein Kriterium (die pins-AK `:1548` prüft genau diese Disambiguierung). Eine Fehl-Bindung an den falschen Link auf einer Mehrfach-Link-Zeile bliebe ungetestet.
- **verifizierbar:** ja — der fehlende Boundary-Test wäre genau der, der eine Bindungs-Regression fangen würde.

#### F-13 · INFO · Manifest-Tie-Break bei doppeltem Zip-Eintrags-Pfad · `DC-QA-02`
- **pfad:** `spec/spezifikation.md:1766-1768` (Schritt 4, Sortierung „nach `<pfad>`")
- **befund:** Das Zip-Format erlaubt **zwei** Einträge mit identischem Namen. Bei identischem normalisiertem `<pfad>` und unterschiedlichem Inhalt ist die Sortierung „aufsteigend nach `<pfad>`" **kein** totaler Ordnungsschlüssel (Tie) — die Reihenfolge der beiden Zeilen und damit der Manifest-Hash wären unbestimmt. Pathologischer Rand, aber gegen die Determinismus-Zusage; ein sekundärer Tie-Break (z. B. nach `<hex>`) würde ihn schließen.
- **verifizierbar:** ja — ein Zip mit zwei gleichnamigen Einträgen zeigt den nicht determinierten Tie.

### Kategorie-Summary (Nachtrag)

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | F-11 |
| LOW | 0 | — |
| INFO | 2 | F-12, F-13 |

Geschlossen aus dem Erst-Report: F-1, F-2, F-3, F-4, F-5, F-6 (Rest F-12), F-7, F-9, F-10.

### Gesamt-Verdikt (Nachtrag) — ACCEPT-WITH-NITS

Der ursprüngliche BLOCK ist **aufgehoben**: F-1/F-2 sind am Kern — dem byte-genauen Content-Manifest — sauber gelöst (Pfad-Sortierung statt Ganz-Zeilen-Sortierung, voller normalisierter Zip-Pfad statt undefiniertem `<name>`, korrigierte konzeptionell-statt-byte-Parität), und lastenheft ↔ spec ↔ ADR sind dazu konsistent. Die MEDIUMs F-3/F-4/F-5 und LOW/INFO F-7/F-9/F-10 tragen real und deckungsgleich über die Straten; F-6 ist bis auf eine AK-Restlücke geschlossen.

**Vor Code zu beheben (der eine echte Nit, F-11):** Der Keyword-Rename `archive` → `zip` ist unvollständig propagiert — die **normative** Spec `DC-FA-SRC-001.a` Schritt 1 (`spec/spezifikation.md:1729`) und der ADR-Marker-Beispiel-Punkt (`docs/plan/adr/0046-…:36`) tragen noch `[archive]`, während Lastenheft, Slice und die ADR-Entscheidung `zip` sagen; die ADR ist dadurch in sich widersprüchlich. Das ist kein Design-Offenpunkt (der Zielwert ist eindeutig `zip`), sondern ein mechanischer Zwei-Datei-Abgleich — aber ein echter Lastenheft-↔-Spec-Widerspruch an der Nutzer-Direktive, der den Marker-Parser falsch erden würde; er muss vor der Implementierung des Marker-Parsers geschlossen werden. F-12/F-13 sind mitzunehmen (AK-Restlücke, Tie-Break), blockieren nicht. Kein neuer BLOCK.
