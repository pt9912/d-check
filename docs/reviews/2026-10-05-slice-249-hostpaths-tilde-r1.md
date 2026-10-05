# Review R1 — slice-249: `hostpaths` erkennt Home-relative Pfade, Ziel-Ventil `exempt-targets`

- **Review-Art:** Code — geprüft gegen den Slice-Plan (`slice-249`, §1 Abgrenzung,
  §3 Spiegel-Liste), `ADR-0098` (Proposed), `ADR-0044`/`ADR-0058` (Bezug),
  `MR-025`/`MR-032`, Hard Rules `AGENTS.md` §3.1/§3.2/§3.7/§3.8, §5 Regel 15/17
  (Maintainability; DoD-Abhakung ist nicht Gegenstand dieses Reviews)
- **Gegenstand:** `slice-249` · Commit `21fde701` (Implementierung: Kern-Regel,
  Config-Rand, Gerüst, Lastenheft 0.94.0, Spezifikation, `ADR-0098`, ADR-Index)
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-05
- **Eingangs-Kontext:** Slice-Plan (`slice-249`); `DC-FA-HOST-001` (Lastenheft)
  und `DC-FA-HOST-001.a`/`SPEC-005`/`SPEC-027` (Spezifikation); `ADR-0098`;
  `MR-025` (Spiegel-Tabelle), `MR-032`; Beobachtungs-Register
  `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet` und
  `BEO-ALL/semantic-change-body-only-edges-stale` (im Plan §8 gesichtet);
  vorherige Findings am Modul `hostpaths`: keine Reports unter `docs/reviews/`
  außerhalb des Archivs.
- **Proben:** temporäre Probe-Testdatei im Paket `rules`, gefahren per
  `make test` — einmal gegen den Stand `21fde701`, einmal mit
  `hostpaths.go` aus `21fde701^` (Vorher-Verhalten). Arbeitsbaum danach
  per `git checkout` + `rm` wiederhergestellt, `git status` leer.

Probe-Ergebnisse (gemeldete Targets, Default-Präfixliste, ohne Ventil):

```text
Eingabe                                  vorher (21fde701^)     nachher (21fde701)
~~/home/alice/x~~                        ["/home/alice/x~~"]    []
Pfad~/home/a/b                           ["/home/a/b"]          []
~/home/x                                 ["/home/x"]            ["~/home/x"]
~~alt ~/src/x~~                          []                     ["~/src/x~~"]
*~/src/e*                                []                     ["~/src/e*"]
_~/src/x_                                []                     []
(~/src/x) "~/src/y" '~/src/z'            []                     ["~/src/x" "~/src/y" "~/src/z"]
~/My Documents/x                         []                     ["~/My"]
siehe ~/                                 []                     []
~/~/q                                    []                     ["~/~/q"]
~C:\Users\a                              ["C:\\Users\\a"]       ["C:\\Users\\a"]
exempt-targets [home/**, src/**] auf
"/home/alice/x und ~/src/q"              —                      ["/home/alice/x" "~/src/q"]
```

---

## Findings

### H-1 — Die Tilde-Sperre im Unix-Muster lässt bisher gemeldete Host-Pfade still verschwinden

- **Kategorie:** HIGH
- **Quelle:** `ADR-0098` Entscheidung 2 · `DC-FA-HOST-001` · Prüffrage 2 (falsche Menge)
- **Pfad:** `internal/hexagon/core/rules/hostpaths.go:73` (Vorbedingung des
  Unix-Musters), `internal/hexagon/core/rules/hostpaths.go:27` (Vorbedingung
  `tildeRE`); Behauptung in `docs/plan/adr/0098-hostpaths-home-relativ-und-ziel-ventil.md:58-61`
- **Befund:** Steht ein `~` unmittelbar vor einem Präfix-Pfad, ohne dass davor
  eine Wortgrenze liegt, greift jetzt **keines** der beiden Muster: das
  Unix-Muster sperrt die Tilde, das Tilde-Muster sperrt den Vorgänger
  (zweite Tilde bzw. Buchstabe). Gemessen: ein per GFM-Strikethrough
  (`~~…~~`) markierter absoluter Host-Pfad und ein an ein Wort geklebter
  Tilde-Präfix-Pfad wurden vorher gemeldet und liefern jetzt null Befunde
  (Probe-Zeilen 1–2). `ADR-0098` sagt für genau diese Musteränderung „die
  Fundstelle bleibt dieselbe", der Lastenheft-Historie-Eintrag 0.94.0 nennt die
  Änderung nur eine Schärfung — tatsächlich lockert sie das Modul für diese
  Klasse, ohne dass die Grenze irgendwo benannt ist. Der Test „Pfad~/x klebt"
  prüft nur den Nicht-Präfix-Fall und kann das nicht zeigen.
- **Verifizierbar:** ja — `make test` mit einem Test über
  Strikethrough-Präfix-Pfad (erwartet ein Befund) wird rot.
- **Klasse:** `wortgrenzen-sperre-verliert-bestandsbefunde`

### M-1 — `--doctor`-Klartext und `--print-config`-Kopf sagen weiter „absolut"

- **Kategorie:** MEDIUM
- **Quelle:** `MR-025` (Spiegel „Klartexte" und „Emittierte Vorlage") ·
  `BEO-ALL/semantic-change-body-only-edges-stale`
- **Pfad:** `internal/hexagon/core/app/diagnose.go:134`;
  `internal/adapter/driving/cli/config_template.go:78`
- **Befund:** `SPEC-027` wurde auf „host-lokaler absoluter **oder
  Home-relativer** Pfad" erweitert, aber der Klartext von `hostpath-forbidden`
  in `--doctor` lautet unverändert „Host-lokaler absoluter Pfad", und die
  Kopfzeile des `hostpaths`-Blocks im `--print-config`-Gerüst — in **derselben**
  Datei, deren nächste Zeilen der Commit änderte — sagt „host-lokale absolute
  Pfade". Ein Konsument, der einen Tilde-Befund mit `--doctor` liest, bekommt
  einen Text, der auf den Fund nicht passt. Die Spiegel-Liste in §3 des Plans
  führt die Klartexte nicht; `MR-025` führt sie ausdrücklich, und „der Spiegel
  ist die Stelle, nicht die Datei" trifft die Gerüst-Zeile wörtlich.
- **Verifizierbar:** nein — kein Gate vergleicht Klartext mit Grund-Code-Tabelle.
- **Klasse:** `semantic-change-spiegel-klartext-vergessen`

### L-1 — Das Ventil akzeptiert Globs, die per Konstruktion nie treffen, während der Präfix-Rand genau das ablehnt

- **Kategorie:** LOW
- **Quelle:** `ADR-0098` Entscheidung 3/4 · Maintainability
- **Pfad:** `internal/adapter/driven/configyaml/configyaml.go:1837`
  (`validateSegmentGlobs` für `hostpaths.exempt-targets`)
- **Befund:** `tracked.exempt-targets` sind repo-relative Pfade ohne führenden
  Schrägstrich; `hostpaths.exempt-targets` trägt denselben Schlüsselnamen, matcht
  aber den **gemeldeten** Pfad mit führendem `/` bzw. `~`. Ein Nutzer, der nach
  dem Vorbild des Nachbarmoduls ein Glob ohne führendes Zeichen schreibt, oder
  eines für Windows-Funde, erhält Exit 0 und ein wirkungsloses Ventil
  (letzte Probe-Zeile: beide Funde bleiben). `ADR-0098` begründet Exit 2 für
  einen Tilde-Präfix gerade damit, dass er „still wirkungslos" wäre — dieselbe
  Klasse bleibt am Nachbar-Schlüssel offen. Fehlrichtung ist rot (Befund bleibt),
  nicht still grün, daher LOW.
- **Verifizierbar:** ja — `make test` mit einem Decode-Test, der ein Glob ohne
  führendes `/` oder `~` als Fehler erwartet.
- **Klasse:** `ventil-glob-ohne-trefferchance-akzeptiert`

### L-2 — Markdown-Hervorhebung an der Pfadgrenze: Endzeichen wandern ins Target, `_` schließt aus

- **Kategorie:** LOW
- **Quelle:** `DC-FA-HOST-001.a` Schritt 3 (Normalisierung) · Maintainability
- **Pfad:** `internal/hexagon/core/rules/hostpaths.go:27`, `:40`
- **Befund:** Das Restpfad-Ende kennt weder `~` noch `*`; ein Tilde-Pfad
  innerhalb von Strikethrough oder `*`-Hervorhebung wird mit den Markern
  gemeldet (Probe: Targets enden auf `~~` bzw. `*`), und ein exaktes
  Ventil-Glob für den eigentlichen Pfad trifft ihn dann nicht. Umgekehrt
  macht `_`-Hervorhebung den Pfad still, weil `_` in der Vorbedingung steht.
  Beides gilt seit jeher für das Unix-Muster; neu ist, dass das
  Tilde-Muster die Klasse erbt, ohne dass Tests oder Spezifikation sie
  nennen (`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`).
- **Verifizierbar:** ja — `make test` mit den Probe-Zeilen.
- **Klasse:** `restpfad-endzeichen-markdown-emphasis`

### I-1 — Lastenheft-Vorbedingung „keine Tilde" gilt nicht für die Windows-Muster

- **Kategorie:** INFO
- **Quelle:** `DC-FA-HOST-001` · Prüffrage 18
- **Pfad:** `spec/lastenheft.md:1733-1735`
- **Befund:** Das Lastenheft formuliert die Vorbedingung pauschal für alle
  Muster, jetzt einschließlich „keine Tilde"; `windowsDriveRE`/`windowsUNCRE`
  schließen nur Wortzeichen aus — ein Laufwerkspfad direkt hinter `~` wird
  gemeldet (Probe). Die Spezifikation (Schritt 2, „Vorbedingung hier: kein
  Wortzeichen davor") ist präzise; die Ungenauigkeit im Lastenheft bestand für
  „URL-/Pfadzeichen" schon vorher und wird durch den Zusatz nur verbreitert.
- **Verifizierbar:** nein.
- **Klasse:** `lastenheft-vorbedingung-pauschaler-als-muster`

### I-2 — Release-Prep-Spiegel, die nach dem Commit Falsches oder Unvollständiges sagen

- **Kategorie:** INFO
- **Quelle:** `AGENTS.md` §5 Regel 17 (planmäßig ausgenommen, §1 des Plans) · `MR-025`
- **Pfad:** `README.md:60`, `README.de.md:59` (Modul-Liste „host-local absolute
  paths"); `docs/user/benutzerhandbuch.md:2623` (Modul-Tabelle „host-lokale
  absolute Pfade"); `docs/user/benutzerhandbuch.md:1475-1476` (§5-Beispiel nennt
  nur den Schrägstrich als Präfix-Verbot, kein `exempt-targets`);
  `CHANGELOG.md` (Schärfung, `ADR-0098` §Konsequenzen verlangt den Hinweis).
  `docs/user/operations.md` trägt keine `hostpaths`-Semantik — ohne Befund.
- **Befund:** Zur Auffindbarkeit in der Release-Prep; kein Befund gegen diesen
  Commit.
- **Verifizierbar:** nein.
- **Klasse:** `release-prep-spiegel-offen`

---

## Negativbefunde

- **Doppelbefunde Unix/Tilde:** geprüft — ein Tilde-Pfad mit Präfix-Namen als
  erstem Segment wird genau einmal, in voller Form gemeldet (Probe); ein
  Präfix-Pfad mit innerer Tilde fällt nur dem Unix-Muster zu. Ohne Befund.
- **Tilde-Negativfälle aus dem Auftrag:** Tilde nach öffnender Klammer bzw.
  einfachem/doppeltem Anführungszeichen — korrekt gemeldet und am Schließzeichen
  abgeschnitten; nackte Tilde mit Schrägstrich am Zeilenende — still;
  Leerzeichen im Pfad — abgeschnitten am Leerzeichen, wie beim Unix-Muster
  (bestehende, im Lastenheft nicht beanspruchte Grenze). Strikethrough: siehe
  H-1/L-2.
- **Byte-Identität ohne neuen Schlüssel:** geprüft gegen `ADR-0098` („gilt für
  den Schlüssel, nicht für das Muster") — das Ventil ist bei leerer Liste ein
  No-op (`exemptTarget` über null Globs). Die Musteränderung ist als gewollt
  deklariert; ihre ungewollte Hälfte ist H-1.
- **Hexagon-Richtung (`ADR-0005`):** `rules` importiert nur `model`, `configyaml`
  nur `model`. Ohne Befund.
- **Suppression (`AGENTS.md` §3.2):** keine `//nolint` im Diff. Ohne Befund.
- **Netz (`DC-QA-03`):** kein Netzzugriff hinzugekommen. Ohne Befund.
- **Kommentare (`AGENTS.md` §3.7):** `tildeRE`-, `CheckHostpaths`-,
  `unixHostpathRE`-, `applyHostpaths`-, `HostpathsConfig`-Kommentar, die
  Gerüst-Zeile und die drei Test-Kommentare tragen Zusage, Kopplung oder Grenze;
  keine Review-Historie, keine Slice-Nummer, kein Mess-Label. Ohne Befund.
- **Modul-Grenze auf der Ziel-Achse (`AGENTS.md` §3.8):** `hostpaths` liest nur
  die gescannten Dateien; das Ventil vergleicht Text, liest keine weitere
  Eingabe. Ohne Befund.
- **Lastenheft-Bump und Historie (`MR-032`):** 0.93.4 → 0.94.0 (Minor bei neuen
  Akzeptanzkriterien), Historie-Zeile vorhanden, vierspaltig, ohne
  Abwärts-Verweis. Ohne Befund.
- **ADR-Form:** `Status` Proposed (Übergang mit der Closure, laut Auftrag
  Absicht); `Schärft:` nennt `SPEC-005` und die Verfeinerungs-Kennung
  `DC-FA-HOST-001.a`; Re-Evaluierungs-Trigger vorhanden; Index-Zeile ergänzt;
  Provenance-Marker im `Bezug:` zeigt den Anlass, begründet nichts.
  Ohne Befund.
- **Zitierte Quellen (Prüffrage 9):** `ADR-0044` trägt das Schnitt-Kriterium
  „Einzelmodul-Erweiterung → bestehende Anforderung ändern" (Zeile 46);
  `ADR-0058` wird nur für die Form der Schlüssel-Weitung zitiert, nicht für die
  Musteränderung. Ohne Befund.
- **Abgrenzung §1:** kein Zeilen-Marker, keine Benutzername-Tilde, keine
  Variablen-Formen, kein Windows-Ventil, kein Handbuch/README/CHANGELOG im
  Diff. Ohne Befund.
- **Commit-Botschaft (`AGENTS.md` §5 Regel 15):** die Belege (Rot vor der
  Kern-Änderung) sind als gelaufen beschrieben, nicht verallgemeinert; die
  Gate-Aussage ist Verifier-Gegenstand. Die einzige überdehnte Aussage steht in
  `ADR-0098`, nicht in der Botschaft — siehe H-1.

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | H-1 |
| MEDIUM | 1 | M-1 |
| LOW | 2 | L-1, L-2 |
| INFO | 2 | I-1, I-2 |

Wiederkehrende Finding-Klasse für die Closure: `semantic-change-spiegel-klartext-vergessen`
(M-1) gehört zu `BEO-ALL/semantic-change-body-only-edges-stale`; H-1 und L-2 zu
`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`.

## Verdikt

**Blockiert.** H-1 und M-1 sind vor der Closure zu klären. H-1 ist
spezifikationskonform (die Spezifikation führt die Tilde in der Vorbedingung),
widerspricht aber der Konsequenz-Aussage von `ADR-0098`. Ob der Verlust behoben
oder als Grenze benannt und die ADR-Aussage korrigiert wird, ist eine
Entscheidung — bei Widerspruch des Implementers läuft sie als Konflikt-Pfad über
den Architect.
