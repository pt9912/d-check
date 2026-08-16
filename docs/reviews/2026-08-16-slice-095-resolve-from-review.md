# Review-Report: slice-095 — `links.resolve-from` (ortsfeste Verweise) — 2026-08-16

**Review-Art:** Code — geprüft wird der ausgelieferte Diff gegen die Verträge,
die er selbst anlegt (Anforderung samt fünf Akzeptanzkriterien, Algorithmus
Schritt 6, §2-Schema, §4-Zeile, `--doctor`-Klartext), gegen die referenzierten
ADRs und gegen die Messung in §3a des Slice-Plans samt Retro-Beleg. Nicht
geprüft wird die DoD-Abhakung (getrennter Kontext, Verifikation).

**Gegenstand:** Commit-Range `eb72435..cd79619` — vier Commits
(Wellen-Eröffnung `50ef29f` · Messung `6cc677d` · Vertrag `e079229` ·
Implementierung `cd79619`); Arbeitsbaum-Stand `cd79619` (= HEAD).

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-fable-5 · **Datum:** 2026-08-16

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  (Lastenheft-Fassung 0.60.0, die fünf neuen resolve-from-Akzeptanzkriterien),
  [`DC-FA-LINK-002`](../../spec/lastenheft.md#dc-fa-link-002--symlink-ablehnung)
  (Symlink-Vorrang, „genau ein Befund pro Linkziel"),
  [`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
  (das mitgeltende Ventil),
  [`DC-QA-01`](../../spec/lastenheft.md#dc-qa-01--performance),
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
- §[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 6 (neu), das §2-Schema `links.resolve-from`, die §4-Zeile
  `link-position-dependent`
- [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md) (Proposed —
  Entscheidungen 1–5, Konsequenzen, Re-Evaluierungs-Trigger),
  [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
  (Schnitt-Kriterium, Ventil-Semantik)
- [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel einer Semantik-Änderung),
  [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md)
  (die Invariante, die hier maschinell wird),
  [`AGENTS.md`](../../AGENTS.md) §3 Hard Rules, [`CLAUDE.md`](../../CLAUDE.md)
- Der Slice-Plan
  [slice-095](../plan/planning/in-progress/slice-095-links-resolve-from.md)
  (§2 dreifacher Anlass, §3 Abnahme-Punkte, §3a Messung samt Retro-Beleg) und
  das Wellendokument
  [welle-76](../plan/planning/welle-76-ortsfeste-verweise.md)

**Läufe dieses Reviews.** Alle Fixtures, Probe-Configs und Mutations-Kopien in
einem Scratch-Verzeichnis außerhalb des Repos (danach gelöscht); keine
Probe-Config wurde im Repo abgelegt, keine Repo-Datei außer diesem Report
geändert. Gefahren: `make test` als Baseline (grün, Exit 0) plus **fünf**
Mutationsläufe (je über Dateikopie mutieren → `make test` → Sicherungskopie
zurückschreiben, Ergebnis am **Exit-Code** abgelesen, nie per `git checkout`);
`make build` aus HEAD, danach rund fünfzehn Fixture-Läufe gegen das
HEAD-Image (netzlos, read-only); die Retro-Replikation über `git archive` des
Planungs-Baums vor der welle-69-Eröffnung in ein Scratch-Fixture; der
Konsistenz-Lauf über den eigenen Baum (383 Dateien, 0 Befunde); zwei
Zeit-Läufe mit/ohne `resolve-from`-Block; abschließend `make gates` über den
Baum samt diesem Report. `git status --short` ist am Ende leer bis auf diesen
Report.

---

## Findings

### F-1 — Ein `dirs`-Verzeichnis mit Tippfehler (oder außerhalb der Scan-Wurzeln) schaltet die Quellen-Rolle still ab — der Fehlzustand ist von Konsistenz nicht unterscheidbar

- `kategorie`: HIGH (Basis MEDIUM, Eskalation: dieselbe Beobachtung im
  Gate-Pfad — die Gruppe ist in `.d-check.yml` scharfgeschaltet und läuft in
  `make doc-check`/`make gates`; identische Klasse wie der blockierende
  F-1 des slice-102-Reviews, der zur fail-closed-Heilung in Fassung 0.59.1
  führte)
- `quelle`: [`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  · fail-closed-Disziplin des Hauses (Lastenheft-Historie 0.59.1:
  „Zeiger-Disziplin"; §Schritt C2/W1-Analogie im Planning-Modul) ·
  Reviewer-Anker „Modul-Grenze auf der Ziel-Achse" (die `dirs`/`fixed-dirs`
  sind selbst benannte Verzeichnisse, die das Modul liest, aber nie scannt)
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:522`
  (`applyResolveFrom` — validiert Form, nicht Existenz);
  `internal/hexagon/core/rules/links_resolvefrom.go:45`
  (`resolveFromGroupOf` — exakter Verzeichnis-Vergleich)
- `befund`: Kein Rand-Check bindet die konfigurierten Orte an den realen Baum.
  Scratch-Kopie des eigenen Repos, ein Zeichen Unterschied in der Gruppe
  (`in-progres` statt `in-progress`): 383 Dateien, **0 Befunde, Exit 0** —
  byte-gleich zur korrekten Konfiguration, weil der heutige Bestand
  positions-stabil verweist. Eine Probe-Datei mit präfixlosem Nachbar-Verweis
  in `docs/plan/planning/in-progress/` zeigt die Folge: mit korrekter Gruppe
  **1 Befund**, mit Tippfehler **0 Befunde** — das wandernde Verzeichnis hat
  seine Quellen-Rolle kommentarlos verloren, und genau die Klasse, die der
  Slice maschinell machen soll, bleibt dort dauerhaft unsichtbar. Dieselbe
  Stille tritt ein, wenn ein `dirs`-Verzeichnis außerhalb von `scan.roots`
  liegt oder von `scan.ignore` geprunt wird (seine Dateien werden nie
  gescannt); der Implementierungs-Commit dokumentiert, dass die eigenen
  Test-Fixtures in genau dieses Loch fielen („lagen ausserhalb der
  Default-Scan-Wurzeln und meldeten 0"). Die drei zugesagten
  Exit-2-Rand-Checks decken nur Form-Fehler; der wahrscheinlichste
  Konfigurations-Fehler — ein Pfad, den es nicht gibt — sieht wie Konsistenz
  aus.
- `verifizierbar`: ja — Repo-Kopie mit Tippfehler-Gruppe plus Probe-Datei,
  Image-Lauf: Exit 0/0 Befunde statt Exit 1/1 Befund (in diesem Review
  ausgeführt).
- `klasse`: config-zeiger-fail-open

### F-2 — Die drei Exit-2-Config-Zusagen sind vollständig ungetestet; der Rückbau einer Zusage überlebt `make test`

- `kategorie`: MEDIUM (Reviewer-Anker „fehlende Negativtests bei neuem
  öffentlichen Vertrag")
- `quelle`: §2-Schema `links.resolve-from` in
  [`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  (Schritt 6, Exit-2-Zusagen) · Haus-Muster der Negativtabellen
  (`TestDecode_WavesFehler` u. a. in
  `internal/adapter/driven/configyaml/configyaml_test.go`)
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:522`
  (`applyResolveFrom`); im Diff existiert **kein** configyaml-Test zu
  `resolve-from`
- `befund`: Die Zusagen wirken im Produkt (drei Probe-Configs: `dirs` mit
  einem Ort, absoluter Pfad, Verzeichnis in zwei Gruppen — jeweils Exit 2 mit
  Klartext), aber keine ist bewacht: Mutations-Stichprobe M-C (den
  `dirs < 2`-Check ersatzlos entfernt) kompiliert und übersteht `make test`
  mit **Exit 0**. Für die Schwester-Fähigkeit `planning.waves` existiert die
  Negativtabelle im selben Adapter-Testfile; für den neuen öffentlichen
  Vertrag hier existiert nichts — jeder künftige Refactor kann die Zusagen
  lautlos verlieren.
- `verifizierbar`: ja — die M-C-Mutation über eine Dateikopie, `make test`,
  Exit-Code (in diesem Review ausgeführt).
- `klasse`: vertrags-zusage-unbewacht

### F-3 — Die Ist-Ort-Vorbedingung (kein Doppelbefund mit `target-missing`) steht nur im Code und im Changelog — der normative Algorithmus verlangt wörtlich das Gegenteil

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 6 ·
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (alle Spiegel derselben Semantik)
- `pfad`: `spec/spezifikation.md:895` („von **jedem** Ort der Gruppe …
  Löst es von mindestens einem Ort nicht auf … ⇒ Grund-Code") gegen
  `internal/hexagon/core/rules/links_resolvefrom.go:65` (Vorbedingung,
  `return nil` bei fehlendem Ist-Ort-Ziel); `CHANGELOG.md:22`
- `befund`: Der Ist-Ort ist Mitglied der Gruppe (`dirs` ∪ `fixed-dirs`).
  Nach dem Wortlaut von Schritt 6 erzeugt ein am Ist-Ort fehlendes Ziel also
  `link-position-dependent` **zusätzlich** zu `target-missing` aus Schritt 4.
  Der Code tut bewusst das Gegenteil (Vorbedingung: Ist-Ort-Auflösung, sonst
  kein Befund), der Changelog verspricht es („meldet weiter nur
  `target-missing`, kein Doppelbefund"), ein Test pinnt es
  (`TestResolveFromKeinDoppelbefundBeiTargetMissing`) — aber die technisch
  verbindliche Vertragsfläche sagt es nirgends, und keines der fünf
  Akzeptanzkriterien deckt den Fall (das Negative-Kriterium setzt „am Ist-Ort
  existiert das Ziel" voraus, statt den anderen Fall zu entscheiden). Wer
  Schritt 6 implementiert oder nachprüft, kommt zu einem anderen Befundsatz
  als das Produkt.
- `verifizierbar`: ja — Fixture mit fehlendem Ist-Ort-Ziel: das Produkt
  liefert genau ein `target-missing`; der Spec-Wortlaut verlangt zwei Befunde.
- `klasse`: vertragsflaeche-luecke

### F-4 — [`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus) widerspricht jetzt aktiv: „das Ventil unterdrückt **nur** die Existenz-/Anker-Prüfung …, keine anderen Befunde" — Schritt 6 und der Code unterdrücken damit auch `link-position-dependent`

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
  (Lastenheft = ranghöchste Quelle) ·
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
- `pfad`: `spec/lastenheft.md:949` (die abschließende Aufzählung der
  Ventil-Wirkung) gegen `spec/spezifikation.md:901` („Das Ventil `ignore-refs`
  gilt wie in Schritt 4") und
  `internal/hexagon/core/rules/links_resolvefrom.go:32` (`refIgnored` →
  `continue`)
- `befund`: Die Anbindung des Ventils ist gewollt (der Slice, das
  Wellendokument und
  [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md) benennen
  `ignore-refs` ausdrücklich als das Ventil für absichtlich ortsgebundene
  Verweise) und getestet (`TestResolveFromVentil`). Aber die
  Ventil-Anforderung selbst zählt ihre Wirkung abschließend auf
  („unterdrückt **nur** die Existenz-/Anker-Prüfung des genannten Ziels,
  keine anderen Befunde") — ein neuer, dort nicht genannter Befund-Typ wird
  jetzt mit unterdrückt. Bei strenger Lesart der ranghöchsten Quelle ist das
  Produktverhalten vertragswidrig; der Diff erweitert den geteilten Vertrag,
  ohne seine Fläche anzufassen. Genau die Fläche also widerspricht aktiv —
  kein Release-Prep-Aufschub.
- `verifizierbar`: ja — Fixture mit `ignore-refs`-Eintrag auf ein
  positionsabhängiges Ziel: 0 Befunde (Produkt) gegen den Wortlaut der
  Anforderung.
- `klasse`: vertragsflaeche-widerspruch

### F-5 — Die `--print-config`-Vorlage führt `links.resolve-from` nicht — die einzige Option des Default-Moduls `links` ist unsichtbar

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CLI-005`](../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
  („es dokumentiert die verfügbaren Module und **Optionen** als Kommentare");
  Präzedenz: Lastenheft-Historie 0.33.0 nennt fehlende Template-Einträge eine
  „Harness-Ehrlichkeits-Lücke", `planning.waves` wurde nach Review-Befund in
  die Vorlage nachgezogen
- `pfad`: `internal/adapter/driving/cli/config_template.go` (kein
  `links:`-Abschnitt; das Kopf-Kommentar nennt links optionslos aktiv)
- `befund`: `links` hat mit diesem Diff erstmals einen eigenen
  Konfigurations-Schlüssel — die Vorlage, deren Vertrag die Sichtbarkeit der
  Optionen ist, zeigt ihn nicht. Ein Konsument, der die Fähigkeit über
  `--print-config` sucht (der dokumentierte Adoptions-Pfad), findet sie
  nicht; die Fähigkeit entstand auf Konsumenten-CR.
- `verifizierbar`: ja — `--print-config`-Lauf des HEAD-Images enthält keinen
  `resolve-from`-Eintrag.
- `klasse`: config-surface-luecke

### F-6 — „Der Retro-Beleg reproduziert den historischen Schaden zeichengenau" — repliziert wird die Zahl, nicht der Schaden: 4 der 19 realen Brüche liegen außerhalb der Fähigkeit, 4 der 19 Retro-Befunde sind real nicht gebrochen

- `kategorie`: MEDIUM
- `quelle`: Messaussagen-Ehrlichkeit (Präzedenz: Lastenheft-Historie 0.59.1
  „eine eigene Messaussage war falsch") ·
  [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md) §Kontext ·
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (dieselbe Aussage auf allen Flächen)
- `pfad`: `docs/plan/planning/in-progress/slice-095-links-resolve-from.md:103`
  („zeichengenau", „mit der erwarteten Verteilung … die vier
  zurückbleibenden Geschwister den Rest"); ADR-0056 §Kontext („retro
  reproduziert … mit der erwarteten Verteilung über Verschobenen und
  zurückbleibende Geschwister"); Commit-Message `cd79619` („mit identischer
  Verteilung")
- `befund`: Die Replikation dieses Reviews bestätigt die Retro-**Zahl** exakt
  (19 `link-position-dependent`, Verteilung 7/4/3/3/2 über den verschobenen
  Slice und die vier Geschwister). Der **reale** Bruch der welle-69-Eröffnung
  setzt sich aber anders zusammen: der Fix-Commit `2a94a408` zog 7 Verweise
  im verschobenen Slice, **8** in den vier Geschwistern (3/2/2/1) und **4**
  in einem Review-Report nach — ebenfalls 19. Die Mengen überlappen zu 15:
  die 4 Review-Report-Brüche kann die Fähigkeit strukturell nie melden
  (Quelle ortsfest, das **Ziel** wanderte), und 4 der Retro-Befunde sind
  präfixlose Geschwister-Verweise, deren Ziele damals nicht wanderten (real
  nicht gebrochen — als latente Positionsabhängigkeit gleichwohl korrekt
  gemeldet). „Exakt die Zahl" stimmt; „zeichengenau", „identische/erwartete
  Verteilung [des realen Bruchs]" stimmt nicht — und die dahinter liegende
  Fähigkeits-Grenze (Target-Move bei ortsfester Quelle bleibt unsichtbar)
  wird auf keiner Fläche gesagt.
- `verifizierbar`: ja — Retro-Lauf (in diesem Review repliziert: 19,
  7/4/3/3/2) gegen die Verweis-Fixes des Commits `2a94a408`
  (7 + 3 + 2 + 2 + 1 + 4).
- `klasse`: messaussage-praezision

### F-7 — Der welle-76-Closure-Trigger und die DoD-Zeile des Slice versprechen „die drei belegten Move-Brüche wären vor dem Move rot gewesen" — zwei der drei Fälle kann weder die Fähigkeit noch die scharfgeschaltete Gruppe melden

- `kategorie`: MEDIUM
- `quelle`: Messaussagen-/Trigger-Ehrlichkeit (wie F-6) · der Slice-Plan §2
  (die drei Fälle) gegen die Semantik aus
  [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md)
  Entscheidung 2 (nur Dateien in `dirs` sind Quellen)
- `pfad`: `docs/plan/planning/welle-76-ortsfeste-verweise.md:39`
  („die drei belegten Move-Brüche wären mit aktivem `resolve-from` **vor**
  dem Move rot gewesen");
  `docs/plan/planning/in-progress/slice-095-links-resolve-from.md:117`
  (DoD: „die beiden oben belegten Fälle wären vor dem Move rot gewesen")
- `befund`: Von den drei §2-Fällen ist nur der dritte (welle-69, 19er-Bruch)
  von der Fähigkeit erreichbar — und auch er nur zu 15/19 (F-6). Fall 1
  (slice-093-Closure: Links der **Review-Reports** auf den Slice) hat
  ortsfeste Quellen in `docs/reviews/` — kein `dirs`-Mitglied, der Bruch kam
  vom Target-Move; ein `link-position-dependent` ist dort per Konstruktion
  unmöglich. Fall 2 (das **Wellendokument** verwies auf den Ruheort und brach
  beim eigenen Move nach `done/`) hat seine Quelle **flach** unter
  `docs/plan/planning/` — auch kein `dirs`-Mitglied der scharfgeschalteten
  Gruppe, und als Gruppe nicht wohlgeformt konfigurierbar (der flache Ort
  ist Eltern-, nicht Geschwister-Verzeichnis des Ruheorts). Der
  Closure-Trigger der Welle ist damit, wörtlich genommen, nicht einlösbar;
  die DoD-Zeile des Slice benennt ausgerechnet die zwei nicht erreichbaren
  Fälle. Die nächste Closure dieser Welle wird die Klasse erneut vorführen:
  die Verweise dieses Reports und des Wellendokuments auf
  `in-progress/slice-095-…` brechen beim Move — ungemeldet.
- `verifizierbar`: ja — strukturell am Config-Stand (`.d-check.yml`-Gruppe)
  und den Quellorten der Fälle 1–2; ein Retro-Lauf gegen den Stand vor der
  slice-093-Closure würde 0 Befunde auf den betroffenen Verweisen zeigen.
- `klasse`: trigger-nicht-einloesbar

### F-8 — Symlink-Ziel am Ist-Ort: die Vorbedingung prüft nur „fehlt/escaped" — eine Referenz erhält zwei Befunde (`symlink` + `link-position-dependent`)

- `kategorie`: LOW
- `quelle`: [`DC-FA-LINK-002`](../../spec/lastenheft.md#dc-fa-link-002--symlink-ablehnung)
  („pro Linkziel entsteht genau ein Befund") · die eigene Klassen-Definition
  („am Ist-Ort grün", Commit `cd79619`)
- `pfad`: `internal/hexagon/core/rules/links_resolvefrom.go:68` (Vorbedingung
  akzeptiert `KindSymlink` als aufgelöst)
- `befund`: Scratch-Fixture: präfixloser Verweis auf einen Nachbar-Symlink in
  einem `dirs`-Verzeichnis ⇒ **zwei** Befunde derselben Referenz (`symlink`
  aus Schritt 5, `link-position-dependent` aus Schritt 6). Am Ist-Ort ist
  aber nichts grün — Schritt 5 meldet bereits; die Reparatur-Ansage der
  neuen Meldung („Pfad präfixieren") ist für diesen Fall die falsche. Für
  das fehlende Ziel wurde derselbe Doppelbefund bewusst ausgeschlossen; für
  den Symlink-Fall nicht.
- `verifizierbar`: ja — das Fixture dieses Reviews (2 Befunde, Exit 1).
- `klasse`: doppelbefund-randfall

### F-9 — Dateien in einem Unterverzeichnis eines `dirs`-Eintrags sind still keine Quellen; die Grenze steht auf keiner Vertragsfläche und in keinem Test

- `kategorie`: LOW
- `quelle`: §[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 6 („für jede gescannte Datei, deren Verzeichnis in `dirs` liegt" —
  mehrdeutig zwischen Listen-Mitgliedschaft und Teilbaum)
- `pfad`: `internal/hexagon/core/rules/links_resolvefrom.go:48` (exakter
  Vergleich `path.Clean(d) == path.Clean(dir)`); `spec/spezifikation.md:893`
- `befund`: Scratch-Fixture: eine Datei in einem Unterverzeichnis eines
  `dirs`-Eintrags mit präfixlosem Nachbar-Verweis ⇒ 0 Befunde. Der exakte
  Vergleich ist in sich konsistent (auch die Hypothesen-Konstruktion
  `path.Join(ort, path.Base(file))` setzt flache Moves voraus), aber die
  Zusage sagt die Grenze nicht, kein Akzeptanzkriterium und kein Test pinnt
  sie — ein Konsument mit gruppierten Unterordnern glaubt sich geprüft und
  ist es still nicht (dieselbe Stille-Richtung wie F-1, ohne Gate-Bezug).
- `verifizierbar`: ja — das Fixture dieses Reviews (0 Befunde, Exit 0).
- `klasse`: zusagen-grenze-still

### F-10 — Zwei kleine Config-Rand-Abweichungen: `..`-Substring statt `..`-Segment, und das gruppeninterne Duplikat trägt die „mehrere Gruppen"-Meldung

- `kategorie`: LOW
- `quelle`: §2-Schema/Schritt-6-Exit-2-Zusagen in
  [`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  („einem Verzeichnis absolut oder mit `..`-**Segment**")
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:533`
  (`strings.Contains(d, "..")`) und `:540` (`gesehen` über alle Gruppen
  **und** innerhalb einer Gruppe)
- `befund`: Ein legaler Verzeichnisname mit `..` als Substring (etwa
  `a..b`) wird mit Exit 2 verworfen, obwohl die Zusage nur das `..`-Segment
  nennt — fail-closed, aber strenger als der Vertrag. Ein Duplikat innerhalb
  **derselben** Gruppe bricht korrekt mit Exit 2, aber mit der Meldung
  „dirs-Mitglied mehrerer Gruppen" — die Diagnose benennt den falschen
  Fehler.
- `verifizierbar`: ja — zwei Probe-Configs gegen das Image.
- `klasse`: config-rand-diagnose

### F-11 — `fsys.Kind` wird vor dem `escaped`-Check aufgerufen: ein Lstat außerhalb der Repo-Wurzel bei flacheren `fixed-dirs` und `../`-Zielen

- `kategorie`: LOW
- `quelle`: Konsistenz zur bestehenden Prüfung (Schritt 5 der Spezifikation:
  „außerhalb liegende Komponenten sind nicht prüfbar"; `CheckLinks` ruft
  `Kind` nie auf einem escaped-Pfad) · Rand von
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  (lesend, aber außerhalb des geprüften Repos)
- `pfad`: `internal/hexagon/core/rules/links_resolvefrom.go:76` (`Kind(rel)`
  vor der `escaped`-Auswertung in Zeile 77);
  `internal/adapter/driven/fs/fs.go:23` (`filepath.Join(Root, rel)` folgt
  `../` über die Wurzel hinaus)
- `befund`: Löst ein Ziel von einem hypothetischen Ort aus lexikalisch aus
  der Wurzel heraus (erreichbar mit einem `fixed-dirs`-Ort geringerer Tiefe
  plus `../`-Ziel), wird der Escape-Pfad dennoch zuerst per Lstat auf dem
  Host geprüft — ein Metadaten-Zugriff außerhalb des Repos, den die übrigen
  Link-Prüfungen bewusst vermeiden. Verhaltensneutral (die
  `escaped`-Bedingung dominiert das Ergebnis), aber ein Disziplin-Bruch
  gegenüber der eigenen Modul-Familie; in der eigenen Konfiguration
  (gleichtiefe Geschwister) unerreichbar.
- `verifizierbar`: nein — kein bestehendes Gate sieht den Zugriff; sichtbar
  nur per Syscall-Trace oder Code-Lektüre.
- `klasse`: escape-disziplin

### F-12 — Determinismus- und Normalisierungs-Zweige sind unbewacht: Sortierung der Divergenz-Meldung und `path.Clean`-Vergleich überleben ihre Rückbauten

- `kategorie`: LOW
- `quelle`: [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)
  (byte-identische Ausgabe) · Mutations-Echtheit der Commit-Behauptung
  („sieben tödliche Rückbauten")
- `pfad`: `internal/hexagon/core/rules/links_resolvefrom.go:92`
  (`sort.Strings`) und `:48` (`path.Clean`-Vergleich);
  `internal/hexagon/core/rules/links_resolvefrom_test.go:63`
  (`TestResolveFromDivergierendeZiele` prüft nur die Anzahl, nie den
  Meldungs-Inhalt)
- `befund`: Mutation M-D (Sortierung entfernt) und M-E (`path.Clean` aus dem
  Gruppen-Vergleich entfernt) übersteht `make test` jeweils mit Exit 0. Die
  Sortierung ist die einzige Verteidigung der Divergenz-Meldung gegen die
  Map-Iterations-Ordnung — unbewacht wäre ihr Rückbau ein realer
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)-Defekt
  (wechselnde Ziel-Reihenfolge in der Meldung); der `path.Clean`-Vergleich
  ist die einzige Normalisierung, die eine Config mit Schrägstrich-Suffix
  am Leben hält (der Adapter speichert `dirs` unnormalisiert). Die
  Kern-Zweige selbst sind bewacht (M-A/M-B sterben, siehe Negativbefunde);
  die Behauptung „sieben Rückbauten" ist für diese zwei Ränder nicht
  eingelöst.
- `verifizierbar`: ja — die zwei Mutationsläufe (in diesem Review
  ausgeführt, je Exit 0).
- `klasse`: mutation-luecke-rand

### F-13 — Mehrere Flächen referenzieren „[`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links) §Ortsfeste Verweise" — einen Abschnitt, den das Lastenheft nicht führt; die Beschreibung der Anforderung erwähnt die Erweiterung mit keinem Wort

- `kategorie`: LOW
- `quelle`: [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  · Lastenheft als abnahmebindende Fläche
- `pfad`: `spec/lastenheft.md:987` (Beschreibung unverändert — die fünf
  neuen Akzeptanzkriterien hängen ohne einen Satz Anforderungstext an);
  Referenzen „§Ortsfeste Verweise" in `CHANGELOG.md:10`, `.d-check.yml:25`,
  `internal/hexagon/core/rules/links_resolvefrom.go:12`,
  `internal/hexagon/core/model/config.go:11`
- `befund`: Der zitierte Abschnitt existiert nur als Titel von Schritt 6 der
  Spezifikation, nicht im Lastenheft; und die Beschreibung von
  [`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  beschreibt weiterhin ausschließlich die Ist-Ort-Prüfung. Wer die
  abnahmebindende Quelle liest, findet die neue Zusage nur implizit in
  Given-When-Then-Form; wer den Verweisen folgt, sucht einen Abschnitt, den
  es nicht gibt.
- `verifizierbar`: ja — Textabgleich der Flächen.
- `klasse`: vertragsflaeche-referenz

### F-14 — `roadmap.md` wohnt dauerhaft in `in-progress/` und ist damit Quelle einer Gruppe, in der sie nie wandert — eine latente Falsch-Positiv-Klasse ohne Ausdrucksmittel in der Config

- `kategorie`: INFO
- `quelle`: dokumentationswürdige, undokumentierte Annahme (Reviewer-Anker
  INFO) · [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md)
  Entscheidung 2
- `pfad`: `.d-check.yml:37` (`in-progress` als `dirs`-Mitglied);
  `docs/plan/planning/in-progress/roadmap.md`
- `befund`: Die Gruppen-Semantik kennt nur Verzeichnisse; die Roadmap ist
  eine ortsfeste Datei **in** einem wandernden Verzeichnis. Heute grün (ihre
  Verweise sind positions-stabil), aber jeder künftige kurze Nachbar-Verweis
  der Roadmap würde als positionsabhängig gemeldet, obwohl die Datei nie
  wandert — das Gegenstück zur 108er-Messung im Kleinen. Als Ausweg bleibt
  nur das ziel-achsige Ventil oder Verweis-Disziplin; keine Fläche sagt das.
- `verifizierbar`: ja — Probe-Verweis in der Roadmap meldet (konstruierbar).
- `klasse`: quelle-ortsfest-im-wandernden-ort

### F-15 — Ablage-Kosmetik der neuen Konstante: `ReasonLinkPositionDependent` steht unter dem Kommentar der Struktur-/Wellen-Codes; die §4-Zeile zwischen den planning-Codes

- `kategorie`: INFO
- `quelle`: Maintainability (der Lockstep-Test verriegelt Mengen, nicht
  Gruppierung)
- `pfad`: `internal/hexagon/core/model/finding.go:51` (unter dem
  Kommentar „Struktur-Codes fuer die Abschnitts-Findung …");
  `spec/spezifikation.md:2582` (links-Code zwischen `closure-note-*` und
  `wave-*`)
- `befund`: Der erklärende Block-Kommentar beschreibt die neue Konstante
  nicht mehr korrekt, und die §4-Tabelle bricht ihre Modul-Gruppierung.
  Konsistenz Code↔Spec ist durch den Lockstep-Test gesichert; verwirrend
  wird es erst beim nächsten Edit an dieser Stelle.
- `verifizierbar`: nein — reine Lektüre.
- `klasse`: ablage-kosmetik

---

## Negativbefunde (geprüft, ohne Befund)

- **Auflösungs-Kern und Determinismus:** Die Meldung der
  Auflösbarkeits-Hälfte nennt den **ersten** nicht auflösenden Ort in
  Config-Reihenfolge (`orte` = `dirs` dann `fixed-dirs`, frühes Return);
  die Divergenz-Meldung sortiert die Ziel-Menge vor dem Join — beide
  Meldungen sind unabhängig von der Map-Iteration
  ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)); die
  Befund-Form (file/line/target/rule, `target` = roher Verweis wie in
  Schritt 4, „ein Befund je Referenz") ist vertragskonform und
  modul-konsistent. Ohne Befund (Testlücke separat: F-12).
- **Byte-Identität ohne Block:** `CheckResolveFrom` kehrt bei leerer
  Gruppen-Liste vor jeder Arbeit zurück; Lauf-Vergleich über den eigenen
  Baum mit/ohne Block: beide 383 Dateien, 0 Befunde;
  `TestResolveFromInertOhneGruppen` pinnt es. Ohne Befund.
- **Ventil-Anschluss:** `refIgnored` wird mit der Quelldatei und dem
  **Ist-Ort**-aufgelösten Ziel geprüft — exakt „wie in Schritt 4"
  (Vertragswortlaut), inklusive `in:`-Quellskopus (matcht die reale
  Quelldatei, nie den hypothetischen Ort); `TestResolveFromVentil` deckt
  den Kernfall. Verhaltens-Parität bestätigt; der Wortlaut-Widerspruch der
  Ventil-**Anforderung** ist F-4. Sonst ohne Befund.
- **Anker-Fragmente und Dekodierung:** Fragment-Strip in `localTarget` und
  `positionDependent` doppelt, aber identisch (erstes `#`); Auflösung beider
  Pfade über dasselbe `ResolveTarget` (volle Prozent-Dekodierung,
  Wurzel-Interpretation führender `/`) — „dieselbe Vorverarbeitung, dieselbe
  Prozent-Dekodierung" aus
  [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md)
  Entscheidung 5 ist eingelöst. Ohne Befund.
- **Escape-Ränder:** Ein am Ist-Ort escapender Verweis wird vor dem Ventil
  ausgefiltert und erhält nie `link-position-dependent` — korrekt, Schritt 4
  meldet `repo-escape` (Klassen-Definition „am Ist-Ort grün"). Ein nur an
  einem hypothetischen Ort escapender Verweis zählt als „löst nicht auf" —
  korrekt im Sinn der Modul-Zusage (existieren **und** innerhalb der
  Wurzel). Ohne Befund (der Neben-Lstat ist F-11, der Symlink-Rand F-8).
- **Mutations-Echtheit der Kern-Zweige (Stichprobe):** M-A (Ist-Ort-
  Vorbedingung entfernt) stirbt an
  `TestResolveFromKeinDoppelbefundBeiTargetMissing`, M-B (`fixed-dirs` aus
  der Ort-Menge entfernt) an `TestResolveFromFixedDirZaehltAlsOrt` — beide
  `make test` Exit 2 mit benanntem Test-FAIL, kein Compile-Artefakt. Die
  Rückbau-Behauptung des Commits hält für die geprüften Kern-Zweige. Ohne
  Befund.
- **Exit-2-Rand im Produkt:** Alle drei zugesagten Config-Rand-Fälle
  (eine-Ort-Gruppe, absoluter Pfad, Verzeichnis in zwei Gruppen) brechen im
  Image-Lauf mit Exit 2 und Klartext samt Gruppen-Index. Funktional ohne
  Befund (Testabdeckung: F-2; Meldungs-Ränder: F-10).
- **Retro-Beleg, Zahl und Verteilung der Messung:** `git archive` des
  Planungs-Baums vor der welle-69-Eröffnung, Probe-Config mit der
  `.d-check.yml`-Gruppe, HEAD-Image: **exakt 19** `link-position-dependent`,
  Verteilung 7/4/3/3/2 — die Messaussage des Slice §3a ist als
  Produkt-Messung reproduzierbar. Ohne Befund (die Deutung als
  „zeichengenaue" Schadens-Reproduktion ist F-6).
- **Heutiger Bestand und Laufzeit
  ([`DC-QA-01`](../../spec/lastenheft.md#dc-qa-01--performance)):**
  Scharfgeschaltete Gruppe über den eigenen Baum: 383 Dateien, 0 Befunde;
  Zeit-Läufe mit/ohne Block 1,58 s / 1,57 s (je ein Kontroll-Lauf — im
  Rauschen, konsistent zur Drei-Lauf-Messung des Commits). Ohne Befund.
- **Wiring, Schnitt und [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):**
  Aufruf unter `applies("links", file)` (Scope-/Modul-Bindung wie die
  Ist-Ort-Prüfung); `links_resolvefrom.go` importiert nur Standard-Lib,
  `model` und `port/driven` (ADR-0005-Schnitt, `make gates` inkl.
  arch-check grün); keine neue Netz- oder Schreib-Tür, alle Läufe
  read-only/netzlos. Ohne Befund (F-11 als Lese-Rand notiert).
- **Lockstep der Befund-Flächen:** `AllReasons`, `reasonTexts`, §4-Zeile und
  `--doctor`-Klartext tragen den neuen Code deckungsgleich (Mengen-Test
  `TestAllReasonsDeckungGegenSpezifikationGrundCodes` + Klartext-Deckung);
  §4-Zeile, Changelog und ADR sagen dieselbe Zwei-Hälften-Semantik. Ohne
  Befund (Gruppierungs-Kosmetik: F-15).
- **Referenz-Richtung (SDP)/Marker-Ehrlichkeit:**
  [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md) trägt
  keinen Provenance-Marker und keinen Slice-Token im Körper (die
  Geschichte-Zeile ist ausgenommener Abschnitt); das Lastenheft umschreibt
  das ADR-Kriterium statt abwärts zu verlinken. Ohne Befund.
- **Commit-Grenzen und Hard Rules:** Vier Commits in der erwarteten
  GF-Reihenfolge (Welle → Messung → Vertrag → Implementierung), der
  Lifecycle-Move der Eröffnung bündelt die gekoppelten Verweise
  ([`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md)-Form);
  keine Inline-Suppression, keine Gate-Lockerung, keine ADR-Mutation im
  Range. Ohne Befund.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 6 | F-2, F-3, F-4, F-5, F-6, F-7 |
| LOW | 6 | F-8, F-9, F-10, F-11, F-12, F-13 |
| INFO | 2 | F-14, F-15 |

## Verdikt

**REQUEST-CHANGES.** Der Kern ist solide: Auflösungs-Semantik, Determinismus,
Ventil-Parität, Byte-Identität, die Retro-**Zahl** und die Laufzeit-Zusage
halten jeder Nachmessung stand, und die bewachten Kern-Zweige sterben an
ihren Tests. Blockierend ist die Kombination aus F-1 (der wahrscheinlichste
Config-Fehler erzeugt dauerhaft stilles Grün im eigenen Gate-Pfad — die
Klasse, die in Fassung 0.59.1 im Nachbar-Modul fail-closed geheilt wurde,
ist hier neu entstanden) und den Vertragsflächen F-3/F-4 (der normative
Algorithmus verlangt wörtlich einen anderen Befundsatz; die ranghöchste
Ventil-Fläche widerspricht aktiv). F-6/F-7 sind Text-, keine Code-Arbeit,
gehören aber vor das Release: die begleitende ADR ist noch Proposed, die
Messaussagen sind jetzt billig zu präzisieren und nach der Closure immutabel
falsch.

**Release-Empfehlung:** Minor-Release **nach** Nachzug — (1) F-1 fail-closed
am Config-Rand oder als Befund (Zeiger-Disziplin wie 0.59.1), (2) F-3/F-4
als Vertrags-Nachzug im noch offenen Slice, (3) F-5 Template-Zeile,
(4) F-6/F-7 Formulierungs-Präzisierung samt ausdrücklicher Fähigkeits-Grenze
(Target-Move bei ortsfester Quelle bleibt unsichtbar), (5) F-2 als
Negativtabelle nach Haus-Muster. Die LOW/INFO-Punkte sind mit dem
Release-Prep oder der Closure-Notiz abtragbar. Handbuch/READMEs/
`operations.md` waren erklärtermaßen noch nicht nachzuziehen und sind kein
Befund dieses Reports.
