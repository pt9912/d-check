# Verifikations-Bericht — slice-286 (DoD-/Plan-Konformität)

**Rolle:** Verifier (gegen DoD §2, Plan §1/§3 samt Plan-Nachtrag/§6 und Spec; nicht Reviewer, nicht Validator)
**Gegenstand:** Commit `bf4a52a8` (Lastenheft 0.111.0, Spezifikation `DC-FA-PSYNC-001.a`, ADR-0113, ADR-Index, Plan slice-286/slice-287, Review-Report R1)
**Vorgänger:** `docs/reviews/2026-10-31-slice-286-code-r1.md` (Vorfassung `ab4d974c..3a67a138`, Commit existiert nicht mehr)
**Datum:** 2026-10-31

## Selbst gefahrene Sensoren (echte Ausgabe)

- `make gates` → Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, **Exit 0**. Darin `d-check: 1230 Datei(en) geprüft, 0 Befund(e)`, semgrep `Ran 55 rules on 73 files: 0 findings`.
- `make adr-check` (Range `HEAD~1..HEAD`, 1 Commit) → `0 Befund(e)`, **Exit 0**.
- `make review-coverage` → `0 Befund(e)`.
- Arbeitsbaum nach den Läufen unverändert (`git status --short` leer).

## Code-Aussagen am Gegenstand gelesen

- `targets.go`: `fsys.ReadFile(mf)` (Z. 96) und `fsys.ReadFile(file)` (Z. 238) ohne `Kind`; nur der Glob-Treffer-Zweig (Z. 341–377, `KindSymlink`) ist fail-closed. Die Aussage in Spezifikation und ADR („folgt einem wörtlich eingetragenen Symlink still") stimmt.
- `fs.Adapter.Kind`: `os.Lstat`; Symlink/Verzeichnis/Datei/fehlend; ein Lstat-Fehler außer „nicht vorhanden" (ENOTDIR) wird als Fehler zurückgegeben. Stimmt mit Spezifikation Schritt 3 („auch ein Datei-Segment mitten im Pfad, das `Kind` als Fehler meldet").
- `ValidLiteralPath` (`model/config.go:1301`): vier Fehlerklassen wie genannt; `strings.TrimSuffix(p, "/")` lässt `Dockerfile/` durch (siehe Abweichung A-1).
- ADR-0113 Bezug: Verweis auf ADR-0112 trägt jetzt nur noch „Pfad-Prüfung" und führt die Symlink-Haltung auf `DC-FA-MET-001.a` zurück; dort steht sie (Spezifikation Z. 3404).

## Review-Findings aus R1

| Finding | Urteil | Beleg |
|---|---|---|
| F-1 MEDIUM (ADR-Geschichte überdehnt Bestätigung) | behoben | Geschichte-Zeile nennt: 1–9 bestätigt, Punkt 6 weicht ab, 4/5/8 geschärft, 10 neu. Rest-Bemerkung: „im Bericht an den Auftraggeber benannt" verweist auf einen Bericht, der im Repo nicht liegt; die Zeile behauptet aber nicht mehr, dass die Abweichung bestätigt sei. |
| F-2 MEDIUM (Treffer-Definition) | behoben | Lastenheft („erste Passung … nicht leere Fanggruppe"), Spezifikation Schritt 5, ADR Entscheid 3 gleichlautend; Lastenheft-Kriterium „Treffer" (`a(b*)` auf `a ab` ⇒ kein Treffer) herleitbar aus Schritt 5. |
| F-3 MEDIUM (Reihenfolge Eingabefehler) | behoben | Lastenheft: „zuerst die Konfiguration aller Pins … dann das Lesen"; Spezifikation Schritt 2 und 10 identisch; Kriterium „Eingabefehler, Reihenfolge" (Pin a Lesefehler, Pin b Ausdruck ⇒ Ursache von b). |
| F-4 MEDIUM (Kriterien zu Randbedingungen) | behoben, mit Restlücke (L-1) | neue Kriterien: Bezugsmenge, Scan-Menge, Reihenfolge der Eingabefehler, weitere Eingabefehler (Symlink-Segment, Verzeichnis, Datei-Segment, gleiches Paar, doppelter Name, Modul ohne Pin), Treffer. |
| F-5 MEDIUM (Fitness „Selbstlauf") | behoben | Zeile trägt „wirkt erst, wenn das Modul registriert ist und das eigene Profil den Block trägt". |
| F-6 LOW (Arbeitsbaum, Hardlink) | behoben | Lastenheft: „Gelesen wird der Arbeitsbaum … Der Symlink-Schutz erfasst Symlinks, nicht Hardlinks und nicht Bind-Mounts." |
| F-7 LOW (Randformen `ValidLiteralPath`) | **halb**: ENOTDIR behoben (Lastenheft + Schritt 3); abschließendes `/` (`Dockerfile/`) weder in Spezifikation noch im Plan-Nachtrag benannt (A-1) |
| F-8 LOW (ADR-0112-Zitat) | behoben | s. o. |
| F-9 LOW (Fremdplan nicht deklariert) | behoben | §3-Tabelle von slice-286 führt die slice-287-Zeile; Diff berührt exakt die dort genannte Datei. |
| F-10 INFO | behoben bzw. benannt | `<file>=?` für Quelle ohne genau einen Treffer festgelegt; „pinsync-Block bei nicht aktivem Modul" und „zwei Ausdrücke derselben Zeile" stehen im Plan-Nachtrag für slice-287. |
| F-11 INFO (ADR vor dem Code Accepted) | ohne benannten Ausgang | im Plan-Nachtrag nicht erwähnt; DoD verlangt `Accepted`, daher bewusst, aber nicht ausgeschrieben (A-2, unkritisch). |

## Prüfpunkte (1)–(7)

1. **Widerspruchsfreiheit Lastenheft · Spezifikation · ADR:** erfüllt. Treffer-Definition, Reihenfolge der Eingabefehler, Exit-Codes (Eingabefehler 2, Befund 1), Bezug des Vergleichs (erste Quelle mit genau einem Treffer), Befund-Ort (Zeile 1 ohne Treffer, `target` = Pin), Grund-Codes, Bezugsmenge: überall gleich. Hinweis ohne Widerspruch: die Spezifikation führt „Quelle ohne `file`/`pattern`" als eigenen Config-Fehler, das Lastenheft deckt ihn über „leerer Pfad" und „keine Fanggruppe" ab.
2. **Kriterien ↔ Verfahren:** weitgehend erfüllt. Offen (L-1): kein Kriterium für (a) drei Quellen, von denen mehrere abweichen (Spezifikation Schritt 7 „jede weitere Quelle"), (b) „Vergleich entfällt bei weniger als zwei Quellen mit genau einem Treffer", (c) Pin ohne Namen, (d) ungültiges UTF-8/Lesefehler außer fehlend. Das Randformen-Kriterium ist als Sammel-Satz („verhält sich … wie oben beschrieben") schwach. Dafür trägt der Plan die Folgepflicht (`BEO-ALL/spec-randbedingung-ohne-test`, Plan-Nachtrag: jede Listenzeile wird in slice-287 am Code gezählt und belegt).
3. **Code-Aussagen:** erfüllt (s. o.), bis auf `Dockerfile/` (A-1).
4. **ADR:** `Status: Accepted`; Index-Zeile vorhanden (README, gleicher Titel/Datum/Status); `## Re-Evaluierungs-Trigger` vorhanden (a–c an Alternativen gebunden). erfüllt.
5. **Version/Historie nach MR-032:** 0.110.3 → 0.111.0; Historie-Zeile mit Datum, Anforderungs-ID, Verweis `GitHub-Issue 14`; Kürzel `PSYNC` in der §3-Liste; Modul-Aufzählung bewusst erst mit Registrierung (Historie-Zeile und slice-287-Tabelle benennen es). Spezifikation: Historie-Zeile 2026-10-31, Schema-Tabelle `pinsync.*`. erfüllt. Spec-Straten enthalten keine ADR-/Slice-Kennung im Körper (doc-check/matrix grün).
6. **Plan verspricht nichts, was nicht im Commit steht:** erfüllt. DoD 1–3 `[x]` und im Commit belegt, 4–9 `[ ]`. Plan-Nachtrag (Symlink-Abweichung, Schärfungen 4/5/8, Tabelle erst mit Code) stimmt mit Lastenheft/Spezifikation/ADR überein; Gegenprobenzahlen des Nachtrags hat R1 reproduziert (nicht erneut gefahren). §1-Abgrenzung eingehalten: kein Code, kein Selbstlauf, keine Release-Prep, kein `strip-prefix`.
7. **`make gates` / `make adr-check`:** grün, siehe oben.

## Urteil je DoD-Punkt

1. Anforderung im Lastenheft (Kürzel, Version, Historie, Kriterien je Befund-Klasse, Exit-2-Ränder, Out-of-Scope): **erfüllt**
2. Spezifikation (Verfahren, Schema, Grund-Codes; Grenzen am Gegenstand gezählt): **erfüllt** (Grund-Codes im Verfahren, Tabelle §4 bewusst erst mit Code; Ausnahme A-1)
3. ADR `Accepted`, im Index, mit Re-Evaluierungs-Trigger: **erfüllt**
4. `make gates` grün: **erfüllt** (Exit 0, eigene Läufe; im Plan noch `[ ]`, Haken fehlt — nachzuziehen)
5. Review mit Report: **offen** — R1 liegt vor (zur Vorfassung); ein Review des neu gebildeten Commits fehlt, Haken zu Recht leer
6. Closure-Notiz mit Lerneintrag: **offen** (§7 „Entsteht bei der Closure")
7. Beobachtungs-Register fortgeschrieben: **offen**
8. Risiken §6 mit Ausgang: **offen** (Risiko 3 „Symlink": Regel ist genannt ⇒ *entfallen* möglich; Risiko 1 Zeilenfortsetzung: Out-of-Scope (4) genannt, Gegenprobe am Bestand steht im Nachtrag)
9. Drei Paarungen: **offen** (Welle-Closure)

## Abweichungen im Wortlaut

- **A-1 (LOW):** Lastenheft/Spezifikation führen „leeres Segment" als Eingabefehler, `ValidLiteralPath` lässt aber `Dockerfile/` durch (`strings.Split(strings.TrimSuffix(p, "/"), "/")`, `model/config.go:1308`). R1-Finding F-7 (erste Hälfte) hat weder in Spezifikation noch im Plan-Nachtrag für slice-287 einen benannten Ausgang. Empfehlung: in den Nachtrag für slice-287 aufnehmen (Test mit `Dockerfile/`; Verhalten festlegen).
- **A-2 (INFO):** R1-F-11 („ADR vor dem Code `Accepted`; Korrektur nur per Folge-ADR") ist nirgends als bewusst angenommen vermerkt.
- **A-3 (INFO):** ADR-Geschichte: „im Bericht an den Auftraggeber benannt" — der Bericht liegt nicht im Repo; der Commit-Text benennt die Abweichung, der Auftraggeber hat Punkt 6 laut Plan aber nicht bestätigt.
- **A-4 (LOW, L-1):** Randbedingungen ohne eigenes Kriterium, s. Prüfpunkt 2; nicht blockierend, weil slice-287 sie am Code zählt.
- **A-5 (Kosmetik):** DoD-Haken 4 (`make gates`) steht im Plan noch offen, obwohl der Lauf grün ist.

## Verdikt

Keine DoD-Verletzung. Die Punkte 1–3 sind bestätigt, 4 ist durch eigenen Lauf belegt (Haken nachziehen), 5–9 sind erwartungsgemäß offen. Alle fünf MEDIUM-Findings aus R1 sind in `bf4a52a8` tatsächlich behoben; von den LOW/INFO-Findings hat nur der `Dockerfile/`-Teil von F-7 (A-1) und F-11 (A-2) keinen benannten Ausgang.
