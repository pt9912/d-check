# Review-Report: slice-116 — ADR-Neuzugangs-Regel + Erstanwendung an den `Proposed`-ADRs

**Datum:** 2026-08-22 · **Review-Art:** Doku-/Konventions-Review (geprüft gegen Slice-Plan slice-116 §2/§3/§5, Wellendokument welle-80 D4 + §3 Closure-Trigger + §6 Out-of-Scope, Hard Rules AGENTS §3.4/§3.5/§3.7 und §5, MR-000 Vergabe- und Adressierungs-Aussage, MR-015, MR-025, ADR-Index-Konventionen, Baseline v5.7.0 `templates/docs/plan/adr/NNNN-titel.template.md` + `templates/docs/plan/planning/slice.template.md` + `regelwerk/grundlagen-source-precedence.md` §ID-Schema als Klammer + `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice) mit eigenen Gegenproben am gebauten Image
**Gegenstand:** Commit `8b8fc1d` (Range `20319cd..8b8fc1d`) — `Schärft:`-Feld in `docs/plan/adr/0050-fence-unclosed-in-spans.md` um die Grund-Code-Festlegung erweitert, `Schärft:`-Feld in `docs/plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md` erstmals angelegt, je eine Geschichte-Zeile, Beispiel beider Adressierungs-Formen im ADR-Index `docs/plan/adr/README.md`, neuer MEDIUM-Anker in `.harness/skills/reviewer.md` (1.5.0 → 1.6.0), neuer Bullet in `AGENTS.md` §5; **vor** der Closure, kein Release, kein Produkt-Code
**Skill:** `.harness/skills/reviewer.md` @ 1.6.0 · **Modell-ID:** `claude-opus-5[1m]`
**Besonderheit:** der Skill ist in diesem Lauf **zugleich Werkzeug und Gegenstand** — geprüft und angewandt wurde die neue Fassung 1.6.0 (Kategorien-Anker inkl. des neuen MEDIUM-Eintrags, Anti-Pattern, sechsfeldriges Output-Schema mit `klasse`, Negativbefund-Pflicht, Ablage).
**Eingangs-Kontext:** Slice-Plan `docs/plan/planning/in-progress/slice-116-adr-neuzugangs-regel.md`; Wellendokument `docs/plan/planning/welle-80-struktur-ids.md`; die drei geschlossenen Vorgänger `docs/plan/planning/done/slice-113-struktur-id-konvention.md`, `docs/plan/planning/done/slice-114-spec-vergabe-spezifikation.md`, `docs/plan/planning/done/slice-115-arc-vergabe-architektur.md`; `harness/conventions.md` §MR-000 und §Aufgelöste Adaptionen, `harness/conventions/MR-015-agents-md-routet.md`, `harness/conventions/MR-025-spiegel-vor-dem-editieren.md`, `harness/conventions/done/MR-027-struktur-id-verzicht.md`; `spec/spezifikation.md` §2 (SPEC-005), §4 (SPEC-030, SPEC-045–SPEC-048), §DC-FA-SPAN-001.a, §DC-FA-PLAN-001.a W1–W5, §7 Historie; `.d-check.yml` (ids-Muster, structure-Regeln, matrix-Klassen); `docs/plan/planning/observations.md` (BEO-002, BEO-004, BEO-007). Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle). Kein `make`-Target im echten Repo (dort läuft ein Gate-Lauf parallel) — alle Produkt-Proben liefen als Image-Lauf gegen eine `.git`-freie Baum-Kopie außerhalb des Repos, Exit je Lauf explizit in eine Datei umgeleitet und gelesen (Arbeitsregel BEO-007); Rückbauten per Dateikopie, `md5sum` gegen das Original geprüft.

## Findings

### F-1 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** ADR-Index-Konvention (`docs/plan/adr/README.md` ist Rang 4 der Source Precedence, AGENTS §2) / Slice-Plan §2 Schritt 3 und DoD-Punkt 2 („Beispiel beider Formen")
- **pfad:** `docs/plan/adr/README.md:18-21`
- **befund:** Das neu eingefügte Beispiel steht in einer einfachen Backtick-Spanne und setzt darin zwei per Backslash geschützte Backticks. CommonMark kennt diesen Schutz nicht: Code-Span-Delimiter binden stärker als Backslash-Escapes, die Spannen paaren daher an den falschen Stellen. Gerendert entsteht kein Beispiel, sondern Bruchstück-Code — und der Schaden endet nicht am Beispiel: die beiden erklärenden Halbsätze („die Kennung steht im Linktext, das Link-Ziel ist der Abschnitt" und „der Abschnitt ohne Kennung, weil die ADR vor der Vergabe `Accepted` wurde") landen selbst **innerhalb** von Code-Spannen, das zweite Beispiel rendert sein `**Schärft:**` als **Fettschrift** statt literal, und `<NNN>` steht als roher HTML-Tag außerhalb jeder Spanne. Wer den Index rendert liest, findet an dieser Stelle die Konvention nicht mehr; wer die Zeile als Vorlage kopiert, übernimmt die kaputte Auszeichnung. Kein Modul sieht das: die Backtick-Parität des Absatzes ist gerade, `spans` schweigt, der Lauf bleibt bei 0 Befunden.
- **verifizierbar:** nein — kein Gate-Lauf bestätigt ihn (Negativ-Probe 1 belegt das Schweigen); belegt ist er per Render-Orakel in Negativ-Probe 6 gegen die CommonMark-Vorrangregel für Code-Spannen.
- **klasse:** backslash-escape-im-code-span-zerreisst-das-beispiel

### F-2 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** `MR-025` (Spiegel *Config-Schema* = `spec/spezifikation.md` §2) / Baseline `templates/docs/plan/adr/NNNN-titel.template.md` (das Feld ist die „Aufwärts-Deklaration der Änderungskopplung: wer diese ADR ändert, zieht von hier die betroffenen Spec-Stellen nach") / Präzedenz ADR-0025
- **pfad:** `docs/plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md:11-14`
- **befund:** Das erstmals angelegte `Schärft:`-Feld nennt zwei der drei Spec-Stellen, die diese ADR verbindlich macht: den Algorithmus-Abschnitt (W1–W5) und die vier Grund-Code-Zeilen. Das **§2-Schema fehlt**, obwohl die ADR sieben Konfigurations-Schlüssel festlegt — `planning.waves.dir`, `done-dir`, `mode`, `glob`, `results-glob`, `next-heading`, `closed-heading` — und ihre Entscheidung 6 die Werte-Menge samt Exit-2-Verhalten des Schlüssels `planning.waves.mode` bestimmt. Die Spezifikation selbst führt die Kopplung in ihrer §7-Historie zweimal wörtlich („um die Schritte **W1–W5** + §2-Schema (`planning.waves.*`) erweitert", „§2-Schema + §4-Zeile nachgezogen"); der Abschnitt trägt seit slice-114 die Kennung SPEC-005, ist also adressierbar, und die Präzedenz existiert im Bestand (ADR-0025 nennt „die Algorithmus-Sektion … **und das §2-Schema**"). Wer künftig Entscheidung 6 ändert — der Re-Evaluierungs-Trigger „ein Adopter braucht eine dritte Kardinalitäts-Semantik" steht in derselben Datei — liest von hier zwei Nachzieh-Stellen statt drei; das Schema-Feld ist genau die Stelle, an der die Fortschreibung dieser ADR bisher zweimal gelandet ist.
- **verifizierbar:** nein — kein Gate misst die Vollständigkeit eines `Schärft:`-Feldes (`ids`/`anchors` prüfen die Auflösung des Genannten, nicht das Fehlende); belegt gegen `spec/spezifikation.md` §7 (Einträge 2026-08-16 und 2026-08-21) und §2-Zeilen `planning.waves.*`, Negativ-Probe 7.
- **klasse:** schaerft-feld-nennt-algorithmus-und-code-aber-nicht-das-config-schema

### F-3 · LOW

- **kategorie:** LOW
- **quelle:** Bestands-Konvention der ADR-Geschichte (10 von 10 mehrzeiligen Geschichte-Listen aufsteigend) / Baseline `templates/docs/plan/adr/NNNN-titel.template.md` §Geschichte (Tabelle `Proposed` vor `Accepted`, also aufsteigend)
- **pfad:** `docs/plan/adr/0050-fence-unclosed-in-spans.md:179-182`
- **befund:** Die neue Geschichte-Zeile ist **oben** eingefügt, vor dem `Proposed`-Eintrag der Datei; die Datei war vor dem Commit aufsteigend (2026-08-09, 2026-08-09, 2026-08-10, 2026-08-10) und ist danach die **einzige** ADR im Bestand, deren Geschichte-Liste nicht aufsteigt. Im selben Commit ist dieselbe Art Zeile in ADR-0055 **unten** angehängt worden — zwei Richtungen in einer Änderung. Der Leser der Datei findet die Entstehungs-Zeile („Proposed") jetzt an zweiter Stelle unter einer Nachtrags-Zeile; der nächste Anhänger hat in derselben Datei zwei widersprüchliche Muster vor sich. Die Chronologie-Monotonie ist im Repo mechanisiert (`structure`, `table-order`), aber nur für die sechs deklarierten **Tabellen** — Listen-Geschichten der ADRs fallen nicht darunter.
- **verifizierbar:** nein — die `structure`-Chronologie-Regeln adressieren `spec/lastenheft.md` §7, `spec/spezifikation.md` §7, zwei Roadmap-Abschnitte, `version.md` §Verlauf und `docs/user/benutzerhandbuch.md` §11, keine ADR-Datei (Negativ-Probe 1 bleibt entsprechend stumm); belegt per Datums-Messung über alle ADR-Geschichte-Listen, Negativ-Probe 8.
- **klasse:** historie-eintrag-gegen-die-richtung-der-eigenen-liste

### F-4 · LOW

- **kategorie:** LOW
- **quelle:** `MR-000` (Vergabe-/Adressierungs-Aussage: „ADRs mit Status `Accepted` **vor welle-80** adressieren weiter per `§`-Anker") / ADR-Index-Konvention (gleicher Wortlaut) / Slice-Plan §5 Risiko 2
- **pfad:** `.harness/skills/reviewer.md:60-63`
- **befund:** Die Ausnahme des neuen Ankers lautet „**Nicht** zu melden ist die alte Form in `Accepted`-ADRs: sie sind immutabel und bleiben auf ihren `§`-Ankern". Die beiden Quellen, auf die der Anker verweist, binden dieselbe Ausnahme an einen **Zeitpunkt** („Accepted **vor welle-80**"); im Anker fehlt diese Schranke. Damit deckt die Ausnahme dem Wortlaut nach jede ADR im Zustand `Accepted` — auch eine, die nach welle-80 entsteht und ihren Zustand innerhalb desselben Slice erreicht (das Immutabilitäts-Gate erlaubt genau das: eine neu angelegte Datei hat keine BASE-Fassung, der Status-Übergang ist zulässig). Der Anker schließt damit die Richtung, vor der Slice-Plan §5 warnt (Scheinbefunde am Alt-Bestand), sauber ab, öffnet aber die Gegenrichtung: ein Neuzugang mit `§`-Zeiger fällt unter die Ausnahme und wird nicht gemeldet — die Regel erodiert still an genau den Dateien, für die der Anker geschrieben wurde.
- **verifizierbar:** nein — Guide-Wortlaut; belegt gegen `harness/conventions.md` §MR-000 und `docs/plan/adr/README.md:8-17`, Negativ-Probe 9.
- **klasse:** ausnahme-ohne-die-zeitschranke-ihrer-quelle

### F-5 · LOW

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §1 („trägt Hard Rules und Pointer … und **dupliziert deren Inhalt nicht** — sonst entsteht Drift") / `MR-015` (AGENTS routet, spiegelt nicht) / `MR-025` (Spiegel *Autoritäts-Doku*)
- **pfad:** `AGENTS.md:246-250`
- **befund:** Der neue Bullet gibt den Feld-Text der Baseline-Vorlage nahezu wörtlich wieder — „die Kennung nennen, wo das Zielelement eine trägt, sonst den Abschnitt … `—`, wenn der Slice keine Spec-Stelle berührt. Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie" steht so in `.harness/baseline/v5.7.0/templates/docs/plan/planning/slice.template.md:21-26` und in `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice —, nennt seine Quelle aber nicht. Die drei benachbarten Slice-Kopf-Bullets tun es (`Verantwortlich:` → „Baseline v5.5.0, template-forward"; `Status:` → „Baseline-`slice.template.md`"; Vorprüfungen → „Baseline-Regelwerk Modul 5/6"). Beim nächsten Pin-Bump führt der Freshness-/Delta-Audit über die genannten Baseline-Stellen; dieser Bullet trägt keinen Anker, über den eine geänderte Feld-Semantik ihn erreichen würde — die Kopie bleibt stehen und wirkt weiter autoritativ.
- **verifizierbar:** nein — Autoritäts-Doku-Drift; kein Gate liest Bullet-Texte gegen die vendorte Baseline (Negativ-Probe 1 bleibt stumm). Belegt per Textvergleich, Negativ-Probe 10.
- **klasse:** briefing-kopiert-baseline-wortlaut-ohne-quellen-anker

## Negativbefunde (geprüft, ohne Befund)

Alle Produkt-Läufe: `docker run --rm --network none -v "$SCRATCH":/repo:ro d-check:latest`, Ausgabe je Lauf in eine Datei, Exit explizit gelesen. `$SCRATCH` ist eine `.git`-freie Kopie des Arbeitsbaums außerhalb des Repos.

1. **Sollform-Lauf (P1).** Unveränderte Baum-Kopie ⇒ `EXIT=0`, `d-check: 423 Datei(en) geprüft, 0 Befund(e)`. Die Zahl der Commit-Botschaft ist damit unabhängig reproduziert; alle zehn konfigurierten Module laufen mit. Ohne Befund.

2. **Sind die neuen `Schärft:`-Anker gate-gedeckt (P2)?** Stille ist nur dann eine Aussage, wenn das Gate den Fall sehen könnte: der Anker der SPEC-045-Zeile in ADR-0055 wurde auf `#4-grund-und-fehler-codes` (ein Bindestrich zu wenig) mutiert ⇒ `EXIT=1`, `docs/plan/adr/0055-…:13 … anchor-missing`. Alle fünf neuen Ziel-Anker (viermal §4, einmal §DC-FA-PLAN-001.a; in ADR-0050 §4 und §DC-FA-SPAN-001.a) lösen im Sollform-Lauf auf. Ohne Befund.

3. **Kennungs-Linkpflicht der neuen Kennungen (P3).** Der Link um SPEC-046 wurde entfernt (nackte Kennung in Inline-Code) ⇒ `EXIT=1`, `docs/plan/adr/0055-…:13 SPEC-046 id-unlinked`. Rückbau per Dateikopie, `md5sum` identisch mit der Repo-Fassung, erneuter Lauf ⇒ `EXIT=0` / 0 Befunde. Die Form „Kennung im Linktext, Link-Ziel = Abschnitt" ist damit sowohl erfüllt als auch erzwungen. Ohne Befund.

4. **Liegt der geänderte Skill im Prüf-Scope (P4)?** Der MR-000-Anker in `.harness/skills/reviewer.md` wurde gebrochen ⇒ `EXIT=1`, `.harness/skills/reviewer.md:63 … anchor-missing`. `scan.roots: ["."]` schließt `.harness/skills/` ein (nur `.harness/baseline/**` und `.harness/cache/**` sind ausgenommen); der neue Anker-Link ist also gate-gedeckt. Ohne Befund.

5. **ADR-Immutabilität, AGENTS §3.5 (Leitfrage 1).** Eigene Zählung: `docs/plan/adr/` trägt 57 nummerierte ADRs; `grep '^\*\*Status:\*\*'` liefert 54 × `Accepted`, 1 × `Accepted — wieder aufgenommen …` (= 55 `Accepted`) und 2 × `Proposed`. `git diff --name-only 20319cd..8b8fc1d -- docs/plan/adr/` nennt genau drei Dateien: die beiden `Proposed`-ADRs und den Index. **Keine `Accepted`-ADR ist berührt.** Der Kern-Eingriff (neues bzw. erweitertes `Schärft:`-Feld) findet nur in `Proposed`-Dateien statt — dort ist er zulässig, das Gate greift ausweislich seiner Konfiguration (`vcs.immutable-when: '^\*\*Status:\*\* Accepted'`) erst ab `Accepted`. Beide Status-Zeilen stehen unverändert auf `Proposed`, kein Entscheidungs-, Alternativen- oder Konsequenz-Absatz ist angefasst; die Nachtrags-Zeilen stehen in beiden Dateien im Anhang `## Geschichte`. Ohne Befund (die **Richtung** der Geschichte-Liste ist Gegenstand von F-3).

6. **Render-Orakel für F-1 (Leitfrage 2, Form).** `sed -n '18,21p' docs/plan/adr/README.md | pandoc -f commonmark -t html` ⇒ `EXIT=0`. Host-`pandoc` diente hier ausschließlich als **Lese-Orakel** (keine Installation, kein Gate, keine `make`-Kette; die Aussage folgt unabhängig aus der CommonMark-Vorrangregel „Code-Span-Backticks binden stärker als Backslash-Escapes"). Ausgabe:

```text
<p>Beispiel beider Formen: neu
<code>**Schärft:** [\</code>SPEC-<NNN>`](…#4-grund--und-fehler-codes)<code>— die Kennung steht im Linktext, das Link-Ziel ist der Abschnitt; alt und unverändert</code><strong>Schärft:</strong>
`spec/architecture.md`
§2<code>— der Abschnitt ohne Kennung, weil die ADR vor der Vergabe</code>Accepted`
wurde.</p>
```

   Zum Vergleich: dieselben Zeilen im Sollform-Lauf des Produkts sind befundfrei — das ist der Grund, warum F-1 als „verifizierbar: nein" geführt wird.

7. **Kennung ↔ Grund-Code und Vollständigkeit von ADR-0050 (Leitfrage 2, Inhalt).** §4 der Spezifikation führt SPEC-030 = `fence-unclosed` (Modul `spans`) sowie SPEC-045 = `wave-drift`, SPEC-046 = `wave-preview-exists`, SPEC-047 = `wave-results-missing`, SPEC-048 = `wave-unregistered` (Modul `planning`) — exakt die Codes, die ADR-0050 Entscheidung 1 bzw. ADR-0055 Entscheidung 4 festlegen; keine Kennung zeigt auf eine fremde Zeile. Für **ADR-0050** ist das Feld vollständig: seine Entscheidung 4 sagt ausdrücklich „**Kein neuer Config-Schlüssel**", und die §7-Historie der Spezifikation verzeichnet zu diesem Vorgang nur §DC-FA-SPAN-001.a Schritt 3 und die §4-Zeile — beide werden genannt. Für **ADR-0055** ist die Lage anders (F-2). Die Wellen-Schritte W1–W5 existieren im adressierten Abschnitt (`spec/spezifikation.md:1761-1824`). Ohne Befund außer F-2.

8. **Richtung der Geschichte je Datei (Leitfrage 1, gemessen statt angenommen).** Über alle 57 ADRs: 11 Dateien tragen eine Geschichte-Liste mit ≥ 2 Datums-Einträgen; 10 davon sind aufsteigend, eine nicht — `0050-fence-unclosed-in-spans.md`, und zwar erst seit diesem Commit (`git show 20319cd:…` liefert für dieselbe Datei 08-09, 08-09, 08-10, 08-10, also aufsteigend). Die Annahme „neueste oben" trägt in diesem Bestand nicht. Siehe F-3.

9. **Reviewer-Anker: Kategorie, Präzision, Versions-Spiegel (Leitfrage 3).** Der Anker steht **innerhalb** des MEDIUM-Listenpunkts (`reviewer.md:56-63`), also in der Kategorie, die Slice-Plan §2 Schritt 3 verlangt; er ist als Prüf-Frage formuliert („Ein neuer Zeiger nur auf ‚§N‘ ist ein Finding") und damit anwendbar — dieser Report ist der Beleg, F-2 ist über ihn entstanden. Kopf-Zeile `**Version:** 1.6.0 · **Datum:** 2026-08-22` ist mit dem Commit-Datum konsistent. **Versions-Spiegel gesucht** (`grep -rn 'reviewer\.md'` ohne `docs/reviews/` und `.harness/baseline/`): `harness/README.md:61` (Guides-Tabelle) beschreibt den Skill **ohne** Versionsangabe, `harness/conventions/MR-028-baseline-v570.md:13` nennt die Datei ohne Version, `CHANGELOG.md:1606` ebenso; keine Stelle spiegelt „1.5.0". Die im Slice-Plan §2 Schritt 5 offen gestellte Frage („Guides-Zeile … Version?") ist damit beantwortet: **nichts nachzuziehen**. Ohne Befund; die Wortlaut-Grenze der Ausnahme ist F-4.

10. **AGENTS §5-Bullet gegen Vorlage und gelebten Haus-Stil (Leitfrage 4).** Inhaltlich **kein Widerspruch**: die Baseline-Vorlage schreibt Kennung → sonst Abschnitt → sonst `—` und die Aufwärts-Richtung, der Bullet sagt dasselbe; der Bestand in `docs/plan/planning/` (sechs Dateien mit dem Feld) ist konform — slice-112 nennt die Verfeinerungs-Kennung als Link, slice-114/115 nennen Abschnitte (zum Zeitpunkt ihrer Beanspruchung trugen die Zielelemente keine Kennung), slice-113 und slice-116 tragen `—`. Der Bullet dupliziert **nicht** den ADR-Satz aus slice-113 (anderer Gegenstand: Slice-Kopf statt ADR-Feld) und nicht die MR-000-Vergabe-Aussage. Verbleibende Beobachtung: der fehlende Quellen-Anker, F-5.

11. **MR-025-Spiegel der ADR-Feld-Konvention (Leitfrage 5).** `grep -rn 'Schärft'` über alle Live-Dokumente: die Konvention lebt in `docs/plan/adr/README.md` (Index-Kopf, in diesem Commit erweitert), `harness/conventions.md` §MR-000 (slice-113), `AGENTS.md:139` §3.4 (nur die Richtungs-Aussage — bleibt korrekt), `spec/architecture.md:21` (Richtungs-Aussage — korrekt), `harness/conventions/MR-006-referenzrichtung-matrix.md` (Richtung), `harness/conventions/done/MR-027-struktur-id-verzicht.md` (aufgelöst, Vergangenheits-Aussage — unverändert korrekt) und `.harness/skills/reviewer.md` (neu). Die Baseline-Vorlage trägt die Ziel-Form und ist derivativ, also zu Recht unberührt (Slice §3). **Kein fehlender Spiegel gefunden.** Ohne Befund.

12. **Wellen-Closure-Bedingungen, welle-80 §3 (Leitfrage 6).** Mechanisch nachgemessen: `.d-check.yml` führt beide `ids`-Muster (`\bSPEC-\d{3}\b` → `spec/spezifikation.md`, `\bARC-\d{3}\b` → `spec/architecture.md`), die `structure`-Regel „jede `###`-Sektion in §2 trägt eine SPEC-Kennung" und `diagrams` opt-in auf die Architektur-Sicht; der Sollform-Lauf ist am eigenen Bestand befundfrei (Probe 1). `grep -c 'ARC-[0-9]' spec/spezifikation.md` ⇒ `0` — die verlangte `matrix`-Messung (kein Abwärtsverweis Spezifikation → Architektur) hält. `MR-027` liegt in `harness/conventions/done/` und trägt in der Tabelle §Aufgelöste Adaptionen die Zeile „Baseline-Konformität … die Vergabe-Aussage trägt MR-000"; MR-000 trägt die Vergabe- **und** die Adressierungs-Aussage samt Zwei-Formen-Satz; der ADR-Index trägt die Zwei-Formen-Regel. **Offen bleiben genau die Closure-Akte selbst:** slice-116 in `done/` (folgt), `make fullbuild` grün, und die Ergebnisnotiz `welle-80-results.md` mit Register-Lese-Schritt (existiert noch nicht — weder flach noch in `done/`). Keine der vier Bedingungen ist **unerfüllbar** oder still liegengeblieben. Ohne Befund.

13. **Slice-Plan §5, beide Risiken (Leitfrage 7).** Risiko 1 („ADR-0055 ist Gegenstand laufender Fortschreibungen — der Nachzug darf keine Entscheidung umschreiben"): der Diff berührt in ADR-0055 ausschließlich die neuen Kopf-Zeilen 11–14 und die Geschichte-Zeilen 184–187; Entscheidungen 1–6, Alternativen, Konsequenzen und Trigger sind byte-identisch. Risiko 2 („MEDIUM-Anker könnte Alt-ADR-Zitate als Befund lesen"): der Anker nennt die Zwei-Formen-Regel ausdrücklich und verweist auf MR-000 — die Scheinbefund-Richtung ist geschlossen (die Gegenrichtung ist F-4). Ohne Befund.

14. **Geprüft und bewusst nicht gemeldet.** (a) In ADR-0055 folgt das neue `**Schärft:**` ohne Satzzeichen und ohne Leerzeile auf die letzte `Bezug:`-Zeile, sodass der gerenderte Kopf-Absatz „… Tabellenzeilen-Lexik `DC-FA-TGT-001` **Schärft:** die Wellen-Schritte …" lautet. Der gesamte Kopf-Block ist in **allen** ADRs dieses Repos ein einziger Absatz (Haus-Stil, kein Baseline-Bruch), und keine Konvention verlangt das Satzzeichen — nach dem Anti-Pattern „kein Stil-Polizist" kein Finding. (b) Die Platzhalter-Form des Index-Beispiels (`SPEC-<NNN>` statt einer echten Kennung) ist sachlich richtig: eine dreistellige Kennung im Index wäre über `ids` selbst linkpflichtig geworden — die Begründung der Commit-Botschaft ist gegen das aktive Muster geprüft und trägt. (c) `CHANGELOG.md` ist nicht gepflegt worden — korrekt, der Commit ändert nichts Nutzersichtbares (AGENTS §5).

15. **Dieser Report gegen das Gate (P6).** Der fertige Report wurde in die Baum-Kopie gelegt und mitgeprüft ⇒ `EXIT=0`, `d-check: 424 Datei(en) geprüft, 0 Befund(e)` — die Ablage-Regel des Skills (`docs/reviews/<YYYY-MM-DD>-<gegenstand>.md`, neue Datei statt Überschreiben) ist eingehalten, und der Report bricht `make doc-check` nicht. Ohne Befund.

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | F-1, F-2 |
| LOW | 3 | F-3, F-4, F-5 |
| INFO | 0 | — |
| **Summe** | **5** | |

## Verdikt

**Nicht abnahmereif ohne Klärung der beiden MEDIUM** — HIGH und MEDIUM blockieren nach Skill-Regel typischerweise, und hier trifft der Regelfall zu: F-1 macht ausgerechnet die **Beispiel-Zeile**, mit der die Welle ihre neue Adressierungs-Form dem nächsten Autor zeigt, im gerenderten Index unlesbar und zieht zwei Sätze der Konvention mit hinein; F-2 lässt das erstmals angelegte `Schärft:`-Feld genau die Spiegel-Klasse aus, für die dieses Repo eine eigene Adaption führt (MR-025 §Config-Schema) und für die im ADR-Bestand eine Präzedenz existiert. Beide sind Doku-Korrekturen ohne Produkt-Wirkung und liegen in `Proposed`-ADRs bzw. im Index — also vor der Closure klärbar, ohne die Immutabilitäts-Grenze zu berühren.

Der Kern des Slice trägt: **keine `Accepted`-ADR ist angefasst** (55 `Accepted` / 2 `Proposed`, gezählt statt geglaubt), beide Status-Zeilen stehen unverändert, die Nachträge stehen im Geschichte-Anhang, alle fünf neuen Kennungs-Zeiger lösen auf einen existierenden Anker auf und entsprechen der Form „Kennung im Linktext, Link-Ziel = Abschnitt" der Baseline, der Reviewer-Anker steht in der richtigen Kategorie und hat in diesem Lauf bereits gearbeitet, und die vier Closure-Bedingungen der Welle sind bis auf die Closure-Akte selbst erfüllt. Die drei LOW sind Nachzieh-Kandidaten und blockieren für sich genommen nicht.
