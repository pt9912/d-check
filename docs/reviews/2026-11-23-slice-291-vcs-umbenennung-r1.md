# Review-Report: slice-291 — `vcs` folgt einer Umbenennung, Anforderung, Spezifikation und ADR (R1) — 2026-11-23

**Review-Art:** Design (Diff gegen Slice-Plan, ADR und Hard Rules, Maintainability;
`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). Nicht gegen die DoD — die
prüft danach der Verifier.

**Gegenstand:** `slice-291`, Commit-Range `c0959fff..bfbe9be0` (vier Commits, nur lokal auf `main`):
`7dbcc240` (Plan-Änderung und Zählungen), `f0622065` (Lastenheft
[`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
Version 0.114.0), `a80cf8f6` (Spezifikation
§[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
Schritt 2 und 4, Schema-Zeile `vcs.follow-renames`), `bfbe9be0`
([ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), vom Architect angeglichen, `Proposed`).

**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`d1c37c14`).
**Modell:** claude-opus-5-5 · **Datum:** 2026-11-23.

**Eingangs-Kontext:**

- Plan `slice-291` §1–§3, §6 samt Plan-Änderung und Plan-Nachtrag (Zählungen am Code)
- [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), [ADR-0024](../plan/adr/0024-vcs-immutable-gate.md),
  [ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md),
  [ADR-0109](../plan/adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md)
- Vorbefunde am Modul: das Architect-Verdikt
  [welle-96](2026-11-23-welle-96-architect-verdikt.md) (F-1 bis F-10) und die Reviews R1–R8 zu
  `slice-267` (Pfad-Nachzug, `vcs.ignore-link-targets`)
- [`AGENTS.md`](../../AGENTS.md) §3.4, §3.5, §3.6, §3.8, §5 (Regeln 13, 14, 15, 16);
  [`harness/conventions.md`](../../harness/conventions.md), vor allem [`MR-032`](../../harness/conventions.md#mr-032)
- Lastenheft §5 (globale Out-of-Scope-Punkte),
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
- Gelesen am Gegenstand (Stand `bfbe9be0`): `internal/adapter/driven/git/git.go` (`AllPaths`,
  `walkTree`, `CommitMessages`, `ancestors`, `isAncestor`), `internal/hexagon/core/rules/vcs.go`
  (`CheckVCS`, `vcsModified`), `internal/hexagon/core/rules/vcs_link_targets.go`
  (`linkTargetResolver`, `pathTree`), `internal/hexagon/port/driven/vcs.go`
- Bekannt und **nicht** Befund dieses Slice: Spezifikation Schritt 2 (Z. 1996–1998), `AllPaths` nehme
  im Index nur drei Arten — als Plan-Änderung in `slice-292` geführt (`dd111da8`)

---

## Findings

### M-1 — Die neuen Kriterien lassen den zweiten Schlüssel offen; für dieselbe Eingabe sagen zwei Stellen Verschiedenes zu

- **Kategorie:** MEDIUM
- **Quelle:** [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  (Akzeptanzkriterien gegen die eigene Grenze); Plan `slice-291` §6 Risiko 1
- **Pfad:** `spec/lastenheft.md` · „**Negative (gelöscht und neu angelegt):** Given `vcs.follow-renames: true` und eine Range“
  (Z. 2348), ebenso Z. 2349 („Negative (Umbiegen)“), 2351, 2353, 2356, 2359 (zweiter Teil); und
  „Allgemein: ein Lauf, der ohne den Schlüssel keinen Befund trägt, trägt mit ihm keinen.“ (Z. 2354)
- **Befund:** Die Negativ-Kriterien nennen nur `vcs.follow-renames: true` und sagen `core-drift-vcs`
  zu, ohne `vcs.ignore-link-targets` auszuschließen oder einen anderen Dateinamen zu verlangen.
  Mit beiden Schlüsseln und einem Ziel gleichen Dateinamens in anderem Verzeichnis (A gelöscht, B neu
  mit anderem Inhalt; oder Umbiegen auf C) normiert die Spezifikation (Schritt 4, Z. 2134–2136:
  „jedes andere wie mit `vcs.ignore-link-targets` allein“) beide Seiten auf denselben Dateinamen, und
  die Grenze der Anforderung bestätigt es (Z. 2315–2316: „bleibt die Grenze ‚gleicher Name, andere
  Datei‘, wie sie ist“) — kein Befund. Umgekehrt verallgemeinert der „Allgemein“-Satz in Z. 2354 über
  das Given seines Kriteriums hinaus und widerspricht dann „Boundary (beide Schlüssel)“ (Z. 2355:
  ohne den Schlüssel kein Befund, mit ihm `core-drift-vcs`). Dass ein offenes Given beide Werte
  umfasst, ist die Lesart, die der Slice selbst anlegt: Das alte Kriterium wurde ausdrücklich auf
  „ohne `vcs.follow-renames`“ eingeschränkt (Z. 2346), weil es sonst neben den neuen gestanden hätte.
- **Failure-Szenario:** In `slice-292` schreibt der Implementer die Negativ-Tests aus den Kriterien und
  parametrisiert sie über `vcs.ignore-link-targets`; mit gleichem Dateinamen läuft der Test gegen
  spezifikationskonformen Code rot, oder der Code wird an das Kriterium angepasst und bricht damit die
  Grenze „gleicher Name, andere Datei“ und das Kriterium „Boundary (Pfad-Nachzug)“. Der Verifier hat
  zwei Rang-1-Stellen, die für dieselbe Eingabe Befund und Nicht-Befund zusagen. Das ist Risiko 1 des
  Plans: **eingetreten**.
- **Warum nicht HIGH:** Es gibt noch keinen Code, und die geteilte Eingabe verlangt beide Schlüssel
  samt gleichem Dateinamen — kein stilles Grün im Gate, sondern ein Widerspruch im Vertrag.
- **Verifizierbar:** nein (kein Gate liest Kriterien gegeneinander); nachlesbar durch Nebeneinanderlegen
  von Z. 2348–2359 mit Z. 2312–2316 und Spezifikation Z. 2134–2136.
- **Klasse:** kriterium-ohne-bedingung-auf-zweiten-schluessel
- **Träger:** Implementer (Lastenheft)

### M-2 — Die Grenze „beide Schlüssel“ ist ohne ihre Bedingung aufgeschrieben: bei einer Umbenennung, die den Dateinamen behält, gilt sie nicht

- **Kategorie:** MEDIUM
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 13 (Grenze gegen den Gegenstand prüfen);
  Prüffrage 18
- **Pfad:** `spec/lastenheft.md` · „Mit beiden Schlüsseln ist ein Nachzug von einer umbenannten Datei auf eine“
  (Z. 2312–2316); `spec/spezifikation.md` · „Schlüsseln ist ein Nachzug von einer umbenannten Datei auf eine **andere**“
  (Z. 2153–2156)
- **Befund:** Beide Straten sagen ohne Einschränkung, mit beiden Schlüsseln sei ein Nachzug von einer
  umbenannten Datei auf eine andere Datei mit ihrem alten Namen Drift. Nach der Normalform derselben
  Spezifikation (Z. 2134–2136: teilnehmendes Ziel → Marke, **Dateiname** des Endes, Anker) gilt das nur,
  wenn das Ende einen anderen Dateinamen trägt. Eine Umbenennung, die nur das Verzeichnis wechselt
  (`docs/x.md` → `spec/x.md`, nach der Definition Z. 2263–2267 eine Umbenennung), normiert auf `x.md`;
  ein Umbiegen auf `other/x.md` normiert ebenfalls auf `x.md` — kein Befund. ADR-0115 §Konsequenzen
  trägt die Bedingung („zwei Ziele gleichen Dateinamens, deren Enden verschieden heißen“), das
  Kriterium Z. 2355 auch („mit anderem Dateinamen“) — die beiden Grenz-Sätze haben sie verloren.
- **Failure-Szenario:** Ein Adopter schaltet beide Schlüssel ein, weil die Anforderung zusagt, damit
  sei die Grenze „gleichnamige andere Datei“ für umbenannte Ziele geschlossen, und migriert ein Layout
  per Verzeichnis-Wechsel (Dateinamen bleiben). Ein Link in einer `Accepted`-ADR, der dabei auf eine
  andere gleichnamige Datei umgebogen wird, passiert das Gate — entgegen der benannten Zusage.
- **Warum nicht HIGH:** Das Grün ist dasselbe, das `vcs.ignore-link-targets` allein schon trägt und als
  Grenze nennt; neu ist nur, dass der Text es für diesen Fall bestreitet.
- **Verifizierbar:** ja, ab `slice-292` — ein Test mit beiden Schlüsseln, Umbenennung `d1/x.md` →
  `d2/x.md` und Umbiegen auf `d3/x.md` meldet bei spezifikationskonformem Code keinen Befund.
- **Klasse:** grenz-aussage-ohne-ihre-bedingung
- **Träger:** Implementer (Lastenheft, Spezifikation)

### M-3 — Ein Pfad, der nach seiner Umbenennung in der Range neu belegt wird, wird trotzdem gleichgesetzt; die Grenze nennt das nur für Merges

- **Kategorie:** MEDIUM
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 9 und 11;
  [`AGENTS.md`](../../AGENTS.md) §3.8; Prüffrage 18
- **Pfad:** `spec/spezifikation.md` · „einem Ende und existiert es in HEAD (bzw. im Index), gilt Quelle → Ende;“
  (Z. 2085) und „Teil hat ein Pfad, der Quelle einer Umbenennung der Abbildung ist“ (Z. 2138);
  `spec/lastenheft.md` · „die Zusage ist Pfad-Identität entlang der Historie“ (Z. 2306) und
  „Trägt die Range einen Merge und belegt eine Seite einen Pfad neu“ (Z. 2323);
  `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „Quelle einer Umbenennung der zusammengesetzten Abbildung“ (Z. 108)
- **Befund:** Die Abbildung prüft am Ende nur, ob der Pfad des Endes in HEAD **existiert**, und setzt
  jeden Pfad gleich, der Quelle oder Ende ist — unabhängig davon, ob er in HEAD noch die umbenannte
  Datei trägt. Zwei lineare Folgen ohne Merge ergeben damit ein Grün, das ohne den Schlüssel rot wäre:
  (a) `a.md` → `b.md`, später `b.md` gelöscht und als **andere** Datei neu angelegt; ein Nachzug
  `a.md` → `b.md` geht durch, obwohl `b.md` nicht mehr aus `a.md` hervorgeht. (b) `a.md` → `b.md` →
  `d.md`, später ein neues `b.md`; `b.md` ist Quelle von `b.md` → `d.md` und wird auf `d.md` gesetzt,
  ein Umbiegen `a.md` → `b.md` (neue Datei) geht durch. Fall (b) entsteht erst mit der Angleichung in
  `bfbe9be0` („zusammengesetzte Abbildung, nicht der Einzelschritte“): nach der Einzelschritt-Lesart
  wäre `b.md` Ende und Quelle zugleich und nähme nicht teil — für eine Kette ohne Neubelegung ändert
  die Angleichung nichts. Die benannten Grenzen führen die Neubelegung nur für den Merge-Fall
  (Z. 2323); die Lesart, im linearen Fall sehe das Modul sie, liegt damit nahe, und Z. 2306 sagt
  „Pfad-Identität entlang der Historie“ zu, die hier unterbrochen ist.
- **Failure-Szenario:** In einem Push mit `vcs.follow-renames` wird `spezifikation.md` nach
  `specification.md` umbenannt, später `specification.md` nach `spec.md`, und ein neues
  `specification.md` entsteht für etwas anderes. Eine `Accepted`-ADR, deren Link von
  `spezifikation.md` auf `specification.md` umgebogen wird, zeigt danach auf das neue Dokument —
  `make adr-check` bleibt grün, ohne den Schlüssel wäre es rot.
- **Warum nicht HIGH:** Noch kein Code, und beide Folgen verlangen eine Neubelegung eines
  Umbenennungs-Pfads in derselben Range; bei der vorgetäuschten Umbenennung (Z. 2307–2308) nimmt die
  Anforderung eine verwandte Lücke ausdrücklich hin.
- **Verifizierbar:** ja, ab `slice-292` — ein Adapter-Fixture mit Folge (a) oder (b) liefert nach dem
  Wortlaut der Spezifikation die Abbildung `a.md` → `b.md` bzw. `b.md` → `d.md`.
- **Klasse:** neubelegung-eines-umbenennungs-pfads
- **Träger:** Architect (Entscheidung 9 in der angeglichenen Fassung), danach Implementer
  (Lastenheft, Spezifikation)

### M-4 — „Keinen Lauf rot, der ohne ihn grün ist“ bestreitet den eigenen fail-closed-Fall

- **Kategorie:** MEDIUM
- **Quelle:** [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  Kriterium „fail-closed (Zwischen-Commit)“; Prüffrage 18
- **Pfad:** `spec/spezifikation.md` · „`vcs.ignore-link-targets` macht `vcs.follow-renames` deshalb keinen Lauf“
  (Z. 2143–2145); `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „macht der Schalter keinen Lauf rot, der ohne ihn grün ist“
  (Z. 113) und „ohne `ignore-link-targets` mit Schalter nie röter als ohne“ (Z. 161)
- **Befund:** Mit dem Schlüssel liest das Modul jeden Commit und jeden geänderten Teilbaum der Range
  und bricht mit Exit 2 ab, wenn eines davon nicht lesbar ist, „auch wenn BASE und HEAD lesbar sind“
  (Spezifikation Z. 2090–2096, Lastenheft Z. 2360). Ohne den Schlüssel liest es diese Objekte nicht und
  bleibt grün. Ein solcher Lauf ist mit dem Schlüssel rot, ohne ihn grün — die Zusage in Z. 2143–2145
  und ADR Z. 113/161 stimmt nur für Befunde. Das Lastenheft formuliert es richtig („keinen Befund“,
  Z. 2354); die Spezifikation und die ADR haben „Befund“ zu „rot“ erweitert.
- **Failure-Szenario:** Ein Adopter liest die Spezifikation, schaltet den Schlüssel als risikofrei ein und
  fährt die CI auf einem partiellen Klon (`--filter=tree:0`), in dem BASE und HEAD vorliegen, Teilbäume
  der Zwischen-Commits aber nicht; das Gate fällt mit Exit 2, obwohl die zitierte Zusage das
  ausschließt.
- **Warum nicht HIGH:** Der Fehler ist laut (Exit 2), nicht still, und das Lastenheft trägt die richtige
  Fassung.
- **Verifizierbar:** ja, ab `slice-292` — das Fixture „unlesbares Objekt eines Zwischen-Commits“
  (ADR-0115 §Fitness Function) liefert ohne Schlüssel Exit 0, mit ihm Exit 2.
- **Klasse:** zusage-uebersieht-eigenen-fail-closed-fall
- **Träger:** Implementer (Spezifikation), Architect (ADR, `Proposed`)

### I-1 — Lastenheft definiert den Schritt anders als Spezifikation und ADR

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `spec/lastenheft.md` · „Ein Schritt ist jeder Nicht-Merge-Commit, der von HEAD“ (Z. 2268) und
  die Historie-Zeile 0.114.0 („je Nicht-Merge-Commit gegen seinen Eltern-Commit“)
- **Befund:** Spezifikation (Z. 2059–2062) und ADR-0115 Entscheidung 2 sagen „genau ein
  Eltern-Commit; ein Commit ohne Eltern-Commit trägt keine Umbenennung bei“. Das Lastenheft nennt
  „jeder Nicht-Merge-Commit … gegen seinen Eltern-Commit“, was einen Wurzel-Commit einschließt, der
  keinen hat. Das Ergebnis ist gleich (vor einem Wurzel-Commit verschwindet kein Pfad); die
  Angleichung in `bfbe9be0` hat nur die ADR erfasst.
- **Verifizierbar:** nein
- **Klasse:** stratum-definition-ungleich-formuliert

### I-2 — `Schärft:` nennt Schritt 1, die Spezifikation führt Range und Vorfahr in Schritt 2

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „Schritt 1 (Range und Vorfahren), Schritt 2“ (Z. 25)
- **Befund:** Der Diff ändert Schritt 1 der Spezifikation nicht; „Basis kein Vorfahr“ und die
  Commit-Menge stehen in Schritt 2 (Z. 2055–2066). Das Feld zeigt auf einen Abschnitt, den die ADR
  nicht schärft.
- **Verifizierbar:** nein
- **Klasse:** schaerft-feld-zeigt-auf-unberuehrten-schritt

### I-3 — „1–14 vom Auftraggeber bestätigt“ steht neben angeglichenem Text

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „1–14 vom Auftraggeber bestätigt“ (Z. 76)
- **Befund:** `bfbe9be0` ändert den Wortlaut der Entscheidungen 1, 2, 3, 9 und 13 nach der Bestätigung;
  die Geschichte-Zeile nennt eine Angleichung, keine erneute Bestätigung. Für 9 ist die Angleichung
  nicht nur Wortlaut (M-3, Fall b). Ob sie eine Neu-Bestätigung braucht, sagt weder ADR noch Plan.
- **Verifizierbar:** nein
- **Klasse:** bestaetigung-ueber-geaenderten-wortlaut

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Risiko 1: altes Kriterium „Boundary (Pfad-Nachzug)“ neben den neuen | Das eingeschränkte Kriterium selbst kollidiert mit keinem neuen: es setzt `vcs.ignore-link-targets: true` **ohne** `vcs.follow-renames`, jedes neue Kriterium setzt `vcs.follow-renames: true`; der Satz „**Ohne** `vcs.ignore-link-targets` meldet schon der reine Nachzug“ steht im Given „ohne `vcs.follow-renames`“. Der Widerspruch liegt bei den neuen Kriterien (M-1) |
| Absatz, Grenzen und ADR zum anderen Namen und zu „dieselbe Datei“ | geprüft. Absatz Z. 2242–2255 und Spezifikation Z. 2173–2180 auf „ohne `vcs.follow-renames`“ eingeschränkt, ADR Entscheidung 5/9/10 gleichlautend; Rest M-2, M-3 |
| Umbenennung der immutablen Datei selbst | geprüft, ohne Befund. Lastenheft Z. 2301–2302, Kriterium Z. 2357, Out-of-Scope Z. 2388–2390, Spezifikation Z. 2022–2024 (Kandidaten-Mengen unberührt) und ADR Entscheidung 6 sagen dasselbe; der Code bildet die Kandidaten aus `AllPaths` und der Mengen-Differenz, `vcsDeleted` meldet |
| Out-of-Scope §5 des Lastenhefts | geprüft, ohne Befund. Kein Punkt getroffen: gelesen werden Commits, Bäume und Blob-Kennungen, keine Nicht-Markdown-Datei wird geprüft, nichts repariert; die Historie-Zeile hält es fest |
| Out-of-Scope-Satz der Anforderung | geprüft, ohne Befund. Ähnlichkeit bleibt ausgeschlossen, Gleichheit ist der Mechanismus; widerspricht weder Absatz noch Kriterien |
| Grenz-Aussagen gegen den Code (Regel 13), selbst nachgezählt am Stand `bfbe9be0` | geprüft, ohne Befund. Datei-Arten in `walkTree`: drei (`case filemode.Regular, filemode.Executable, filemode.Symlink`, Z. 173; `Dir` Z. 165; Gitlink fällt durch). Index im gestagten Modus: jeder Eintrag (Z. 95–98). Commit-Menge `CommitMessages`: Ausschluss BASE samt Vorfahren, Nicht-Merge `<= 1` Eltern (Z. 217), Sortierung nach Kurz-SHA (Z. 226). `AllPaths` prüft keine Abstammung, nur den Leerfall (Z. 115–124); `isAncestor(` einmal gerufen (Z. 402) und definiert (Z. 503). Ziele ohne Auflösung in `linkTargetResolver`: sechs Bedingungen (leer, `/`, `://`, `?` in Z. 84; außerhalb des Repos, nicht im Baum in Z. 88) — die Spezifikation nennt genau diese sechs (Z. 2124–2127). Port: fünf Methoden. Alle acht Kommandos des Plan-Nachtrags reproduziert |
| fail-closed | geprüft. Lastenheft Z. 2295–2296, Kriterium Z. 2360, Spezifikation Z. 2090–2096 und ADR Entscheidung 13 einig: jedes nötige Objekt, Exit 2, auch bei lesbarem BASE und HEAD; „Basis kein Vorfahr“ ohne Exit 2 an allen drei Stellen gleich. Rest M-4 |
| Referenz-Richtung ([`AGENTS.md`](../../AGENTS.md) §3.4) | geprüft, ohne Befund. In den hinzugefügten Zeilen beider Spec-Straten kein `ADR-`, kein `slice-`, kein `welle-`, kein Commit-Hash, kein Link nach `docs/plan/adr/` (grep über die `+`-Zeilen leer); „Begründung in begleitender ADR“ ist die gelebte Form. `spec/architecture.md` unberührt |
| Determinismus und opt-in ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)) | geprüft, ohne Befund. Ohne Schlüssel byte-identisch und kein Commit zwischen BASE und HEAD gelesen (Kriterium Z. 2361, Schema-Zeile, Spezifikation Z. 2051–2052); am Code trägt die Zusage heute: `AllPaths` liest nur die Vorfahren von BASE und die beiden Bäume. Reihenfolge-Unabhängigkeit und „zehnmal byte-gleich“ (Architect-Verdikt F-8) stehen in Kriterium und Spezifikation |
| [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) | geprüft, ohne Befund. Nur lokale git-Objekte, kein Schreiben; der Absatz „Strikt opt-in“ ist unverändert |
| Scan-Achse ([`AGENTS.md`](../../AGENTS.md) §3.8) | geprüft. Benannt ist, was über die Scan-Menge hinaus gelesen wird (jeder Commit der Range, geänderte Teilbäume, Index; nur Blob-Kennungen, kein Inhalt — Lastenheft Z. 2291–2297, Spezifikation Z. 2090–2099) und welche Zusage dort gilt (Inhalts-Gleichheit, nicht „derselbe Gegenstand“). Rest M-3: wo die Zusage „Pfad-Identität entlang der Historie“ bricht, ist nur für Merges benannt |
| Version und Historie nach [`MR-032`](../../harness/conventions.md#mr-032) | geprüft, ohne Befund. 0.113.2 → 0.114.0 (neuer Schlüssel, Minor wie 0.108.0 für `hostpaths.variables`); Zeile oben in der Tabelle, Verweis-Spalte `—`, das Issue im Änderungstext. „Fünfzehn Kriterien“ nachgezählt: Z. 2347–2361 |
| Historie der Spezifikation | geprüft, ohne Befund. Datumszeile oben, nennt Schritt 2, Schritt 4 und §2, keine Abwärts-Referenz |
| ADR-Form und Status ([`AGENTS.md`](../../AGENTS.md) §3.5) | geprüft, ohne Befund. `Proposed`, also nicht immutabel; `## Geschichte` um eine Zeile ergänzt, Provenance-Marker zeigen auf `slice-291` als Herkunft, nicht als Begründung; `## Re-Evaluierungs-Trigger` vorhanden; Rest I-2, I-3 |
| Abgrenzung des Plans (§1) | geprüft, ohne Befund. Der Diff berührt vier Dateien: Plan, Lastenheft, Spezifikation, ADR. Kein Code, kein `.d-check.yml`, kein README/CHANGELOG/Handbuch |
| Commit-Disziplin | geprüft, ohne Befund. Jede Botschaft nennt `slice-291` und `DC-FA-VCS-001`, `bfbe9be0` auch `ADR-0115`. Die Plan-Änderung `7dbcc240` liegt vor dem Inhalt. Regel 15: „vollständige Liste der Ziele, die nicht auflösen“ trägt (sechs von sechs), „fünfzehn Kriterien“ trägt; die Aufzählung in `f0622065` führt „Merge“ als eigenes Kriterium, es steht in „Boundary (Kette)“ — die Summe stimmt |
| Plan-Änderung vor dem Inhalt | geprüft, ohne Befund. Die drei zusätzlichen Stellen liegen in schon geführten Abschnitten (Schritt 2, Schritt 4, §2) und stehen im Plan, bevor `a80cf8f6` sie schreibt |

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 4 | M-1, M-2, M-3, M-4 |
| LOW | 0 | — |
| INFO | 3 | I-1, I-2, I-3 |

## Verdikt

**Nachbessern vor Closure.** Kein HIGH. Die vier MEDIUM blockieren: zwei Stellen des Vertrags
sagen für dieselbe Eingabe Verschiedenes zu (M-1), eine Grenze ist ohne ihre Bedingung
aufgeschrieben (M-2), eine Neubelegung in linearer Historie ergibt ein falsches Grün, das keine Grenze
nennt (M-3), und eine Zusage übersieht den eigenen fail-closed-Fall (M-4). Plan-Risiko 1 ist damit
**eingetreten** (M-1), nicht entfallen; das Register
`BEO-ALL/plan-entscheid-steht-gegen-lastenheft` steht bei 2×, ob M-1 seine Klasse trifft (Plan gegen
Lastenheft oder Lastenheft gegen sich selbst), entscheidet der Planner. M-3 berührt Entscheidung 9
in der Fassung von `bfbe9be0` und geht deshalb zuerst an den Architect. Die Grenz-Aussagen gegen den
Code stimmen, die Referenz-Richtung, Opt-in, Versionssprung und Commit-Disziplin sind sauber.
