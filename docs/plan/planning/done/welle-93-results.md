# Welle welle-93 — Modul `metrics`: Kennzahl eines Planungsdokuments gegen Messung aus git — Closure-Notiz

**Welle:** welle-93
**Abschluss:** 2026-10-31
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — Ergebnis, nicht Tätigkeit.

- **slice-282** (Anforderung, Spezifikation und ADR):
  [`DC-FA-MET-001`](../../../../spec/lastenheft.md#dc-fa-met-001--kennzahl-eines-planungsdokuments-gegen-messung-aus-git-modul-metrics-opt-in)
  im Lastenheft (Fassung 0.110.3), Verfahren und fünf Grund-Codes in der Spezifikation,
  [ADR-0112](../../adr/0112-metrics-modul.md) Accepted. Neun Entscheide des Auftraggebers
  vor dem Code festgelegt, darunter Block-Syntax (Fenced-Block, Zeilen `schluessel: wert`,
  kein YAML), Herkunft der Basis und Formel-Form.
- **slice-283** (Modulkern): Block lesen, Obergrenze gegen die ganzzahlige Formel,
  Schwellenregel, leerer Grund; fünf Grund-Codes `metrics-limit-mismatch`,
  `metrics-reason-empty`, `metrics-cut-reason-missing`, `metrics-over-limit`,
  `metrics-estimate-mismatch`; Eingabefehler als Exit 2.
- **slice-284** (Messung aus git): `LineChanges` im VCS-Port, go-git-Adapter ohne
  git-Binary im Image, Messung gegen die Obergrenze, Hinweise (Marke ausstehend,
  nicht gezählte Dateien, Pfad ohne Vorkommen).
- **slice-285** (Tabellenzählung): die Schätzung aus der einen Tabelle des Plans,
  Prüfung 2, Hinweis „Marke ausstehend“ mit berechneter Schätzung — die Hälfte, die der
  Issue-Autor am 2026-10-30 in die Anforderung zurückholte (GitHub-Issue 13).
- Release **v0.90.0** (außerhalb der Welle, Release-Prep nach Welle-Plan §6): GHCR und Docker
  Hub tragen denselben Index-Digest.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Die Fragen an den Auftraggeber im Issue statt im Chat:** der Autor hat sie dort
  beantwortet, und die Antworten sind Beleg am Gegenstand; zwei davon (Tabellenzählung
  doch im Modul, „Prosa-Form“ entfällt) änderten den Plan, ehe Code entstand.
- **Architect-Verdikt und Plan-Review vor dem ersten Slice:** beides fand Spiegel und
  Bestandsangaben, die der Plan falsch trug (drei MEDIUM im Plan-Review), bevor ein
  Test rot wurde.
- **Review in frischem Kontext, danach Verifier, je Slice:** jeder der vier Slices hatte
  MEDIUM-Findings, die der Autor nicht sah; der Verifier hat die Einarbeitung am
  Dateistand geprüft statt am Bericht.
- **Die Paritäts-Probe durch den Eigentümer des Anlass-Repos mit dem veröffentlichten
  Image:** sie fand keine Abweichung im Urteil, und die fünf benannten Unterschiede
  standen alle schon in der Spezifikation.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile möglichst mit der Konsequenz, die
daraus schon gezogen wurde.

- **Ein vierter Slice kam hinzu:** der Autor kehrte die Aufteilung um (die
  Tabellenzählung gehört ins Modul); slice-285 entstand per Plan-Änderung, die
  Reihenfolge blieb fachlich bindend.
- **„Dieselben Leser wie die übrigen Module“ galt für Zäune nur für eine von zwei
  Leser-Familien** (slice-282 Review R1/R2, slice-285 Review F-1): ein eigenes Muster
  neben dem einen Öffner-Test des Produkts. Konsequenz: das Modul nutzt die geteilten
  Zaun- und Zellen-Leser; die tolerantere Leseweise steht als Abweichung in der Anforderung
  (Beleg im Register, siehe unten).
- **„Binärdatei“ war im Adapter „Patch ohne Chunks“** (slice-284 F-1): leere Datei und
  Submodul wurden als binär gemeldet. Konsequenz: `File.IsBinary` und eine eigene
  Submodul-Behandlung; die Abweichungsklassen gegenüber `git diff --numstat` stehen in der
  Zusage. **Bekannte Grenze, nicht behoben:** die Zeilenzählung von go-git kann bei
  verschobenen Blöcken von der `git`-Zählung abweichen. Sie ist benannt, und die Probe
  des Autors hat keinen Plan gefunden, der davon abhängt.
- **Plan-Abgrenzung im Lauf berührt:** slice-285 änderte `spec/lastenheft.md` über den
  Fassungsstand hinaus, ohne vorherige Plan-Änderungs-Zeile (Verifier, slice-285).
  Die Zeile steht nachträglich in der Plan-Tabelle des Slice; einen Register-Beleg trägt
  das nicht, es ist ein Einzelfall dieser Welle.
- **Die Paritäts-Probe ist nicht Plan für Plan nachlesbar.** Das Anlass-Repo ist für
  diese Planung nicht lesbar, und der Autor hat die Pläne anonymisiert berichtet: die Tabelle
  unten gruppiert nach Stand und Variante und nennt den Stand des Anlass-Repos nur als
  *vier Zweig-Stände mit sieben Plänen in Arbeit*, nicht als Commit. Das trägt die Aussage
  „dasselbe Urteil“ am Bericht des Eigentümers, nicht an einem von hier reproduzierten Lauf.
- **Der Digest-Pin im Handbuch folgt dem Release-Tag erst im Nachgang** (eigener Commit
  nach der Pipeline) — wie bei jedem Release dieses Repos.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier stehen **nur** Beobachtungen, die im Register
3× erreicht haben.

Kein Eintrag über der Schwelle, der noch auf einen Ausgang wartet. Gelesen wurde das
Register vollständig: Einträge mit Stand `offen` und mindestens drei Belegen gibt es
keine. Gezählt, nicht neu verkörpert:

- `BEO-ALL/module-promise-only-on-scan-axis` steht bei 4×, ist verkörpert
  (`AGENTS.md` §3.8, Reviewer-Skill-Anker) und bekommt einen Beleg (slice-282): das
  Modul liest Tabellen, Pläne und git-Revisionen, die es nicht scannt, und die Grenze
  steht in der Anforderung statt in der Annahme.
- `BEO-ALL/eigene-menge-gemessen-fremde-behauptet` steht bei 18×, Stand differenziert nach
  Instanz, und bekommt einen Beleg (slice-284): die Messung zählt die eigene Menge, die
  Gleichheit mit `git` war an einer von Hand ermittelten Zahl behauptet.
- `BEO-ALL/gemeinsame-leser-sind-zwei-familien` ist **neu** und steht bei 2×, offen
  (slice-282, slice-285). Tritt sie ein drittes Mal auf, gehört die Frage *welche zwei
  Leser-Familien hat „derselbe Leser“?* in den Reviewer-Skill.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`../observations/`](../observations/).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*: nichts Neues.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur, das
Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine. Offen bleibt keine Arbeit der Welle. Der Abbau des Programms im Anlass-Repo
ist ein Vorgang jenes Repos (Welle-Plan §6), die abweichende Zeilenzählung bei
verschobenen Blöcken eine benannte Grenze der Anforderung, kein Slice.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle vier Slices liegen in `done/welle-93/`; `make fullbuild` Exit 0 am Stand `32c69c0b` —
  Image-Hash `sha256:d013d1753cac032c46d5e300b4859c1e0051d078c727ab0a72c4fbbef0385ded`
  (`make gates`, `make image-test`, `make bench`, `make completeness-check` und
  `make verify-closure-notes` darin grün; Vollständigkeit: 56 Anforderungen, 0 Waisen).
- **Ohne den Block nichts Neues:** `make blackbox-probe REF=v0.89.0`: byte-identisch über
  20 Vergleiche (stdout, stderr und Exit getrennt, Kanarienlauf vorher und nachher).
- **Betriebsbedingungen am gebauten Image:** `make image-test` (in `fullbuild`) fährt die
  Messung unter `:ro`, mit fremdem Benutzer und ohne `.git` (Exit 2) sowie die
  Tabellen-Schätzung — beide Fälle grün.
- Replay-Lauf gegen Golden Set: entfällt, d-check führt keinen (deterministischer Kern,
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Coverage gesamt: 94,8 % bei Schwelle 93 % (`make coverage-gate`), offene Carveouts: keine.

**Der Anlass trägt — Tabelle Gegenstand → Urteil des Moduls → Urteil des Programms.**
Quelle ist der Abnahme-Kommentar des Eigentümers des Anlass-Repos vom 2026-10-31
([Issue 13](https://github.com/pt9912/d-check/issues/13#issuecomment-6539067315)); Image
`ghcr.io/pt9912/d-check@sha256:d8ce395beb78c7f0b79dbdbfa6d294365463e30c572a94daa1d23ab63c72e1db`
(v0.90.0), Aufruf mit `-config <datei> -enable metrics`, Mount schreibgeschützt, netzlos;
Gegenstand: vier Zweig-Stände des Anlass-Repos, sieben Pläne in Arbeit, je mit dem Modul und
dem bisherigen Programm auf demselben Klon.

| Gegenstand | Läufe | Urteil des Moduls | Urteil des Programms | gleich |
|---|---|---|---|---|
| Stand wie er ist | 4 Stände | Exit 0, Hinweise mit Zahlen (z. B. Schätzung 7815 gegen Tabelle erwartet 8055, Messung 6540) | Exit 0, dieselben Hinweise und Zahlen | ja |
| Messung über der Marke | 4 Stände | `metrics-over-limit` | `metrics-over-limit` | ja |
| Marke ≠ Formel | 4 Stände | `metrics-limit-mismatch` | `metrics-limit-mismatch` | ja |
| Schwelle ohne `schnitt-grund` | 4 Stände | `metrics-cut-reason-missing` | `metrics-cut-reason-missing` | ja |
| Ausnahme-Eintrag mit leerem Grund | 4 Stände | `metrics-reason-empty` | `metrics-reason-empty` | ja |
| Basis gleich `HEAD`, abweichende Schätzung | 4 Stände | `metrics-estimate-mismatch` bei Messung 0 | `metrics-estimate-mismatch` | ja |
| Eingabefehler (Tabelle fehlt oder mehrdeutig, Basis unbekannt, Schlüssel doppelt, Pflichtfeld fehlt, Zahl mit führender Null, Pfad mit `..`, Klasse unbekannt, Nr. leer oder doppelt, Tabelle ohne Zeilen, vier und sechs Zellen) | 13 Fälle | Exit 2, Ursache benannt | Exit 2, gleiche Ursache | ja |
| Plan ohne Block | 1 Fall | Exit 0, „nicht gemessen“ | Exit 0, „nicht gemessen“ | ja |

Exit-Code in 24 von 24 Variantenläufen gleich, Befund-Art und Zahlen je Plan gleich. Fünf
Abweichungen, alle in der Anforderung genannt oder vereinbart: ungeschütztes `|` im Code-Span
einer Tabellenzelle (Modul toleranter), Zaun-Einrückung und fehlender Schlussstrich (Modul
toleranter), „Prosa-Form“ wie „kein Block“, der zusätzliche Hinweis „Pfad existiert weder an
der Basis noch an HEAD“, Abbruch beim ersten Eingabefehler statt Prüfung aller Pläne. Der
Eigentümer nimmt das Modul ab.

**Trigger-Audit (Schritt 2), je Klasse geprüft:**

- **Carveout:** keine aktiven (`docs/plan/carveouts/` trägt nur `done/` und die README).
- **Bootstrap-aware Gate:** `coverage-gate` ist kalibriert (Schwelle 93 %, gemessen 94,8 %,
  Senkung nur per ADR); kein Hochschalt-Trigger gefeuert, keine neue Stufe.
- **ADR:** keine Bedingung feuerte in der Welle. Gelesen wurden die berührten:
  [ADR-0112](../../adr/0112-metrics-modul.md) trägt seinen Re-Evaluierungs-Trigger — die
  Fixture-Gleichheit der Zeilen-Summe mit `git diff --numstat` hält am Adapter; die
  benannte Grenze bei verschobenen Blöcken ist der Beobachtungspunkt für Trigger (b);
  [ADR-0024](../../adr/0024-vcs-immutable-gate.md) (Port-Vorbild) bleibt gültig, der Port wuchs
  nur um `LineChanges`. Entscheid 7 der ADR sagt „längengleich“, Lastenheft und Spezifikation
  „mindestens gleiche Länge“; die ADR ist Accepted und immutabel, die Spezifikation führt.
- **Hard Rule:** in `AGENTS.md` §3 tragen alle einen Auflösungs-Trigger *permanent* oder
  benennen die ungedeckte Kategorie; §3.8 wurde in slice-282 angewandt (nicht gescannte
  Eingaben des Moduls benannt). Nichts feuerte, keine Zeile wird entfernt.

**Die drei Paarungen:** (a) Anker — kein Steering-Loop-Eintrag mit `liegt in`-Feld, keine
Paarung auszulösen; (b) Folge-Slice — keiner genannt; (c) Register — die zitierten
Beobachtungen existieren als Verzeichnisse und tragen alle Belege.
