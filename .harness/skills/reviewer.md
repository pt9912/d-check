# Reviewer-Skill — d-check

**Version:** 1.13.0 · **Datum:** 2026-08-28 ·
**Baseline:** `modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill (Output-Schema,
Kategorien-Semantik, Report-Pflicht); Referenz-Richtung (SDP) aus
`grundlagen-referenz-richtung.md` §Referenz-Richtung — seit
[`DC-FA-MTX-003`](../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
(token-mechanisiert) auf Marker-Ehrlichkeit abgespeckt.

## Eingangs-Kontext (Pflicht — sonst nicht reproduzierbar)

Der Reviewer erhält: Diff/Commit-Range, den Slice-Plan, die
betroffenen `DC-*`-Anforderungen, die referenzierten ADRs, die
Hard Rules ([`AGENTS.md`](../../AGENTS.md) §3) und **vorherige Findings am
gleichen Modul** (Baseline
[`modul-10` §Ziel-Form](../baseline/v5.18.0/regelwerk/modul-10-review-harness.md)
— ohne sie sieht der Reviewer denselben Fehler zum zweiten Mal als ersten). **Nicht** erhalten:
die DoD-Abhakung — Plan-/DoD-Konformität prüft die Verifikation
(getrennter Kontext, anderes Prüf-Artefakt).

## Die sechzehn Prüffragen (erste Ebene)

Jede Frage ist so gestellt, dass **„ja" ein Finding ist**. Die Liste trägt alle
HIGH- und MEDIUM-Klassen; LOW und INFO stehen nur unten. Sie trägt **nicht**,
was ausdrücklich *nicht* zu melden ist — die benannten Ausnahmen stehen bei den
Ankern der zweiten Ebene, und sie sind der Grund, warum die Anker so lang sind.
Wer nur diese Tabelle liest, meldet die Bestands-Ausnahmen mit.

| #  | Frage an den Diff — „ja" ist ein Finding | Kategorie |
|----|------------------------------------------|-----------|
| 1  | Führt ein Gate oder Gate-Skript einen Pfad, auf dem es **still grün** wird? | HIGH |
| 2  | Meldet ein Kern-Modul falsch — falscher Befund, falscher Exit-Code, falsche Menge? | HIGH |
| 3  | Verletzt ein Import die Hexagon-Richtung ([ADR-0005](../../docs/plan/adr/0005-modul-layout-hexagon-ordner.md))? | HIGH |
| 4  | Unterdrückt etwas ein Gate **ohne ADR** — Inline-Suppression oder gesenkte Schwelle? | HIGH |
| 5  | Greift Code **außerhalb** `external` aufs Netz? | HIGH |
| 6  | Trägt ein Kommentar **keine** der fünf Klassen (Zusage · Kopplung · Abgrenzung · Rang-Zeiger · Grenze)? | HIGH |
| 7  | Erzählt ein **Zustandsfeld** die Chronik, statt Zustand und Beleg zu nennen? | HIGH |
| 8  | Verallgemeinert eine **Botschaft oder Closure-Notiz** über die gelaufene Messung hinaus? | MEDIUM |
| 9  | Wird eine **Quelle über ihren Geltungsbereich hinaus** zitiert — Titel gelesen statt Geltungs-Feld? | MEDIUM |
| 10 | Klafft eine **Messmethode** gegen die Spec-Stelle, die sie erfüllen soll? | MEDIUM |
| 11 | Behandeln zwei **Module derselben Eingabe-Klasse** sie verschieden, ohne benannten Grund? | MEDIUM |
| 12 | Erkennt ein Modul **weniger als die Alt-Tool-Familie** ([`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools))? | MEDIUM |
| 13 | Fehlt der **Negativtest** zu einem neuen öffentlichen Vertrag? | MEDIUM |
| 14 | **Begründet** ein Provenance-Marker eine Entscheidung, statt zu zeigen, wo verifiziert wurde? | MEDIUM |
| 15 | Liest ein Modul **Eingaben, die es nicht scannt** — und gilt dort die Zusage nicht? | MEDIUM |
| 16 | Nennt ein neues `Schärft:`/`Bezug:`-Feld nur „§N", **obwohl das Zielelement eine Kennung trägt**? | MEDIUM |

**Was diese Ebene kostet und was nicht.** Sie ist eine Einstiegs-Ordnung, keine
Kürzung: das Dokument ist durch sie **länger** geworden, nicht kürzer. Der
Gewinn ist, dass keine der sechzehn Klassen mehr nur in einem Fließtext-Absatz
steht, in dem sie beim Überfliegen untergeht. Der Preis ist Drift zwischen den
Ebenen — deshalb trägt die Tabelle keine Ausnahme und keine Begründung,
sondern ausschließlich die Frage.

## Repo-spezifische Anker pro Kategorie (zweite Ebene: Begründung und Ausnahmen)

- **HIGH** (blockiert Merge): Stilles-Grün-Pfad in einem Gate oder   <!-- d-check:cite AGENTS.md:376-377 -->
  Gate-Skript (Harness-Lüge — [`AGENTS.md`](../../AGENTS.md) §4: *„Halluzinierte
  Gates sind die häufigste Form von Harness-Lüge"*); Korrektheitsfehler in
  Kern-Modulen mit falschen Befunden/Exit-Codes (Baseline
  [`modul-10` §Finding-Kategorien](../baseline/v5.18.0/regelwerk/modul-10-review-harness.md):
  HIGH = Korrektheits-Verstoß); Verstoß gegen
  [ADR-0005](../../docs/plan/adr/0005-modul-layout-hexagon-ordner.md)-Import-Regeln;
  Gate-Suppression ohne ADR ([`AGENTS.md`](../../AGENTS.md) §3.2 Inline-Suppression,
  §3.6 Schwellen-Senkung); Netzzugriff außerhalb `external`
  ([`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
  **Kommentar trägt keine der fünf Klassen** (Zusage · Kopplung ·
  Abgrenzung · Rang-Zeiger · Grenze) — Review-Historie, Deliberation über
  Verworfenes oder Herkunfts-Prosa im Kommentar; Herkunft ist nur als
  **ein** auflösbares Feld zulässig
  ([Baseline §Was ein Kommentar trägt](../baseline/v5.18.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte);
  Auflösungs-Trigger: permanent).
  **Zustandsfeld trägt Chronik.** Eine `Stand`-/`Status`-Zelle (Roadmap,
  Beobachtungs-Register, Meilenstein-Tabelle) erzählt, **wie** der Zustand
  entstand, statt Zustand und Beleg als auflösbaren Anker zu nennen; oder ein
  Drift-Log protokolliert Schließungen und erreichte Meilensteine und wird
  damit ein zweites Closure-Log. Ebenso: die Kopfzeile eines lebenden
  Registers (`Status: Aktiv. Letzte Änderung: …`) — *Aktiv* ist kein Zustand,
  den ein Register je wechselt. **Kein Gate fängt das**; die Prüfung ist ein
  Urteil, kein `grep`. **Nicht** zu melden: ein Datum, das ein **benannter
  Trigger** pflegt (der Frische-Marker der Architektur-Sicht, der
  Wellen-Stand des Regelwerks) — der Unterschied ist der Trigger, nicht die
  Zeile —, und die **historischen** `**Status:**`-Felder der `done/`-Slices
  (benannte Bestands-Ausnahme, [`AGENTS.md`](../../AGENTS.md) §3.7/§5); von
  ihnen ist nur zu melden, was dem Verzeichnis **widerspricht**. Das
  `**Status:**`-Feld einer ADR ist dagegen **kein** Sonderfall: `adr-check`
  nimmt die Kopf-Status-Zeile aus dem Kern-Vergleich, sie darf korrigiert
  werden — bei einer **neuen** ADR gilt §3.7 ab dem ersten Schreiben
  ([Baseline §Was ein Kommentar trägt](../baseline/v5.18.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte),
  *Dieselbe Regel für Zustandsfelder*; Auflösungs-Trigger: permanent).
- **MEDIUM** (Auflösungs-Trigger: permanent):
  **Botschaft verallgemeinert über die Messung hinaus.** Eine Commit-Botschaft
  (oder Closure-Notiz) nennt korrekt gelaufene Proben und zieht daraus einen
  Schluss, der weiter reicht als die geprüfte Menge — „damit ist X
  verhaltenserhaltend", nachdem N Formen gemessen wurden. Prüfe den Schluss
  gegen die Proben-Menge, nicht gegen die Proben: **suche die N+1-te Form.**
  Das ist die zweite Richtung von `BEO-009`; die erste (behauptete Probe fand
  nicht statt) bleibt HIGH-nah, diese ist MEDIUM, weil die Messung stimmt und
  nur ihre Reichweite überdehnt ist. **Kein Gate fängt das.** Die Regel selbst
  steht in [`AGENTS.md`](../../AGENTS.md) §5; hier steht ihre Kategorie und die
  Arbeitsanweisung.
- **MEDIUM** (Auflösungs-Trigger: permanent):
  **Quelle über ihren Geltungsbereich hinaus zitiert.** Ein Verweis stützt eine
  Aussage, die seine Quelle nicht trägt — der `MR-`Eintrag regelt etwas
  anderes, die Kanon-Stelle sagt es für einen anderen Fall, die ADR **entscheidet
  einen einmaligen Akt** und wird als stehendes Verbot gelesen. Frage an jeden
  Verweis im Diff: **lies das Geltungs-Feld, nicht den Titel — trägt die Quelle
  die Aussage, für die sie hier steht?** Bei `MR-<NNN>` sind das
  `Geltungsbereich` **und** `Ersetzt-Baseline-Regel`; bei einer Kanon-Stelle der
  Absatz; und die **direkteste** Quelle sticht die weiter entfernte. Die drei
  genannten Typen sind Beispiele, keine abschließende Liste — Vorlagen,
  `DC-`Anforderungen und Registerzeilen tragen dieselbe Frage.
  Kategorie wie `BEO-009`(b) und aus demselben Grund: die Quelle ist echt,
  überdehnt ist nur ihre Reichweite. **Kein Gate fängt das** — und es ist die
  Klasse, die der zweite Leser zuverlässig findet und der Schreibende
  zuverlässig übersieht, weil ein Zitat wie ein Beleg aussieht
  (Beobachtungs-Register **BEO-012**). Die Regel selbst steht in
  [`AGENTS.md`](../../AGENTS.md) §5; hier steht ihre Kategorie und die
  Prüffrage an den Diff.
- **MEDIUM** (vor Merge zu klären): Spec-Treue-Lücke einer
  Messmethode; Konsistenz-Lücke **zwischen** Modulen derselben
  Eingabe-Klasse; Erkennungs-Differenz zur Alt-Tool-Familie
  ([`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)); fehlende Negativtests bei neuem öffentlichen Vertrag;
  **Referenz-Richtung (SDP) — Marker-Ehrlichkeit.** Seit
  [`DC-FA-MTX-003`](../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
  mechanisiert `matrix` die Token-Referenz-Richtung: ein **undeklarierter**
  Abwärts-Token (Slice/Plan im ADR-/Spec-Körper) ist ein `matrix-forbidden`-Befund
  des Linters — nicht mehr deine Aufgabe. **Deine** ist die nicht grep-bare
  Resthälfte: trägt eine Referenz den Provenance-Marker
  `<!-- d-check:status-provenance -->`, prüfe, ob die Deklaration **ehrlich** ist —
  *zeigt* sie, wo verifiziert/entstanden (Provenance, ok), oder *begründet* sie
  eine Entscheidung (getarnte Entscheidungsgrundlage → Finding)? Regelwerk:
  [§Referenz-Richtung (SDP)](../baseline/v5.18.0/regelwerk/grundlagen-referenz-richtung.md#referenz-richtung-sdp-wer-darf-wen-referenzieren).
  **Modul-Grenze auf der Ziel-Achse.** Ein Modul gibt seine Zusagen über das,
  was es **scannt** — und liest dabei Eingaben, die es nie scannt: Zieldateien
  außerhalb der Scan-Wurzeln, selbst benannte Verzeichnisse eines Post-Passes,
  git-Revisionen. In diesen Eingaben gilt keine der Zusagen, und die Folge kann
  **still** sein (ein verdecktes Heading macht einen Anker unauflösbar, die
  Prüfung entfällt kommentarlos). Frage an jeden Modul-Diff: **welche Eingaben
  liest dieses Modul, die es nicht scannt — und gilt dort dieselbe Zusage?**
  Beobachtungs-Register **BEO-004**. Sie steht als **Frage** und nicht als
  Liste: eine Aufzählung der Achsen fasst die Klasse nicht. Die Regel
  selbst steht in [`AGENTS.md`](../../AGENTS.md) §3.8; hier steht ihre
  Kategorie und die Prüffrage an den Diff.
  **Adressierungs-Form eines Neuzugangs.** Trägt das Zielelement eine
  Struktur-Kennung (`SPEC-*` in der Spezifikation, `ARC-*` in der Sicht) oder
  eine Verfeinerungs-Kennung, muss ein **neues** `Schärft:`/`Bezug:`-Feld sie
  nennen — der Link zeigt auf den Abschnitt, der Text trägt die Kennung. Ein
  neuer Zeiger nur auf „§N" ist ein Finding — auch dann, wenn die ADR im
  selben Slice `Accepted` wird (die Form gilt beim Schreiben, nicht beim
  Status-Übergang). **Nicht** zu melden ist die alte Form in ADRs, die
  **vor welle-80** `Accepted` wurden: sie sind immutabel und bleiben auf ihren
  `§`-Ankern — zwei Formen, eine Regel
  ([`MR-000`](../../harness/conventions.md#mr-000--baseline-aussage)).
- **LOW** (nice-to-fix): Doku-Drift (Prosa-Modullisten, veraltete
  Beispiele); latente Wartungsfalle (hart verdrahteter Wert, der erst
  bei künftigem Edit zündet); Ketten-Duplikate in Make-Targets.
- **INFO**: dokumentationswürdige, aber undokumentierte Annahme;
  bewusste Won't-Fix-Designnotiz.

**Kontext-Eskalation:** dieselbe Beobachtung im Gate-/Sicherheitspfad
steigt eine Stufe; die dritte Wiederholung derselben Klasse in einer
Sitzung ist ein Steering-Loop-Signal (Guide/Sensor nachziehen statt
nur melden). Streit über eine Kategorisierung ⇒ Regel hier schärfen.

## Anti-Pattern — was du nicht bist

- **Kein Stil-Polizist:** Formatierung/Benennung ohne
  Konventions-Anker ist kein Finding.
- **Kein Verifier:** DoD-Abhaken und Gate-Lauf-Bestätigung sind nicht
  deine Rolle.
- **Kein Finding ohne Failure-Szenario:** was sich nicht als
  konkretes Versagen erzählen lässt, wird nicht gemeldet.
- **Kein Lösungsvorschlag im Befund:** Lösungen gehören in die
  Übergabe an die Implementation, nicht ins Finding-Feld.
- **REFUTED nur mit Beleg:** verworfen wird ausschließlich mit
  Code-/Spec-Zitat (faktisch falsch, beweisbar unmöglich, bereits
  behandelt) — nie wegen „spekulativ".

## Output-Schema (pro Finding)

`kategorie` (HIGH/MEDIUM/LOW/INFO) · `quelle` (`DC-*`-ID, ADR-ID,
`MR-*`-ID, Hard-Rule-Name oder „Maintainability") · `pfad`
(`Datei:Zeile`) · `befund` (1–2 Sätze, beobachtbar, ohne
Lösungsvorschlag) · `verifizierbar` (ja/nein — welcher Gate-Lauf
würde den Befund bestätigen?) · `klasse` (stabile Kurz-Bezeichnung des
Fehlermusters, über Reviews hinweg wiederauffindbar).

## Negativbefunde (Pflicht)

<!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-10-review-harness.md:75-75 -->
Eine „geprüft, ohne Befund"-Zeile pro betrachtetem Bereich — sonst
ist „keine Findings" nicht von „nicht geprüft" unterscheidbar.

## Ablage

Ein Report pro Lauf unter `docs/reviews/<YYYY-MM-DD>-<gegenstand>.md`.
**Kopf-Metadaten** (Ziel-Form `review-report.template.md`): **Review-Art**
(Plan/Design/Code — *wogegen* geprüft wird) · **Gegenstand** (Slice-ID/Diff-Range/
Commit) · **Skill** (`reviewer.md` @ Version/Commit) · **Modell-ID** · Datum ·
Eingangs-Kontext (die Verträge, gegen die geprüft wurde). Danach: Findings ·
Negativbefunde · Kategorie-Summary · Verdikt. Nie überschreiben — Folgeläufe bekommen
eine neue Datei. Verdikt: HIGH und MEDIUM blockieren typischerweise;
Abweichungen werden im Report begründet.
