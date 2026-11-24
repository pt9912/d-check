# Architect-Verdikt — slice-292 (`vcs` folgt einer Umbenennung, Modulkern), Plan gegen Entscheidungslage

**Rolle:** Architect (Modul 8), frischer Kontext, **vor** Code. Rolleninhaber: pt9912.
**Gegenstand:** [slice-292](../plan/planning/open/slice-292-vcs-umbenennung-modulkern.md) im Stand nach dem
vierten Nachzug, gegen [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
(Lastenheft 0.114.3, 30 Kriterien, Z. 2390–2419),
§[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
und [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) (`Proposed`).
**Vorgänger:** [Verdikt zu welle-96](2026-11-23-welle-96-architect-verdikt.md),
[Verifikation von slice-291](2026-11-24-slice-291-vcs-umbenennung-verify.md).
**Datum:** 2026-11-24.
**Geändert wurde nichts außer diesem Report.** Keine Folge-Entscheidung ist nötig; ADR-0115 bleibt
unverändert.

## Gesamturteil

Der Plan trägt. Dateitabelle, Port-Vertrag und die Reihenfolge aus Entscheidung 14 stimmen mit ADR und
Anforderung; die heikelsten Fälle aus R3 und R4 (Tausch über Zwischennamen, freigemachter Name, Rotation,
zwei Quellen, Löschen und Neuanlegen) sind richtig gefasst, samt Schlüsselkombination und
Commit-Bedingung. **Keine Entscheidung fehlt**, die ADR-0115 nicht schon trifft: die Teilung von Kette und
Gleichsetzung zwischen Adapter und Kern steht in Entscheidung 13, und die Form der Abbildung folgt daraus
(F-6, nur im Plan festzuhalten). Offen sind **Lücken im Plan**: sechs Kriterien haben keinen oder nur einen
halben Testfall (F-1, F-2), DoD und Abgrenzung widersprechen der Dateitabelle an zwei Stellen (F-3), und
die Angleichung der Index-Enumeration ändert ein Ergebnis **ohne** den Schalter, was Welle und Plan nicht
sagen (F-4). Alles **V1**.

**Braucht keine Entscheidung des Auftraggebers.** Kein Befund steht gegen das Lastenheft oder gegen einen
bestätigten Entscheid; F-4 folgt aus der Spezifikation (Rang 2), die schon gilt, und F-5 ist ein
Altbestand der Spezifikation, den der Planner dem Implementer zuweist.

Verdikt-Kürzel (Modul 8, Konflikt-Pfad): **V1** = Entscheidung gilt, Plan hat falsch behauptet/ausgelassen ·
**V2** = Folge-Entscheidung nötig · **V3** = Lockerung legitim, aber undokumentiert.

## Prüffrage 2 — Kriterien gegen Testfälle

Zugeordnet nach den Testfällen in slice-292 §3 (Z. 136–194). Kriterium 1–5 sind die bestehenden;
sie tragen die heutigen Tests von `vcs_test.go`, `vcs_pfad_nachzug_test.go` und `git_test.go`, und Risiko 1
(Z. 232–238) bindet sie an „unverändert grün“.

| Nr. | Kriterium (Lastenheft Z.) | Testfall in §3 | Stand |
|---|---|---|---|
| 1–5 | Happy Path, Modul-aus, Negative, fail-closed, Pfad-Nachzug ohne Schalter (2390–2394) | bestehende Tests, Risiko 1 | gedeckt |
| 6 | Happy Path (Umbenennung folgen) (2395) | Z. 136–137, nur „ein Commit“ | **halb** — „in einem anderen Commit“ fehlt (F-1) |
| 7 | gelöscht und neu angelegt (2396) | Z. 138 | gedeckt |
| 8 | Umbiegen (2397) | Z. 139 | gedeckt |
| 9 | Anker, Linktext (2398) | Z. 140 | gedeckt |
| 10 | umbenennen und ändern (2399) | Z. 141 | **halb** — „Inhaltsänderung in einem späteren Commit ⇒ kein Befund“ mit Nachzug A → B fehlt; Z. 142 fährt eine Kette, nicht diesen Fall (F-1) |
| 11 | Kette (2400) | Z. 142, 145, 177–179, 189 | **halb** — „X → Y und Y → Z auf zwei Seiten eines Merges ⇒ setzt sich nicht fort“ fehlt; Z. 146 ist ein anderer Fall (zwei Enden) (F-1) |
| 12 | Mehrdeutigkeit (2401) | Z. 143, 151–152, 173, 177 | gedeckt |
| 13 | Merge-Seite löscht, legt später neu an (2402) | Z. 177–178 | gedeckt |
| 14 | Ende gelöscht und neu angelegt (2403) | Z. 174 | gedeckt |
| 15 | Zwischenpfad neu belegt (2404) | Z. 175 | gedeckt |
| 16 | seiten-gleich (2405) | Z. 149–150 | gedeckt |
| 17 | beide Schlüssel (2406) | — | **fehlt** — A → B mit anderem Namen, Umbiegen von A auf eine andere Datei mit dem Namen von A; Z. 167 und 169 sind Kriterium 19 und 20 (F-1) |
| 18 | Negativ-Kriterien mit beiden Schlüsseln (2407) | nur einzelne Fälle | **fehlt** als Regel — die Negativ-Fixtures laufen nicht in beiden Kombinationen (F-2) |
| 19 | Umbiegen auf den alten Namen (2408) | Z. 167–168 | gedeckt |
| 20 | Umbenennung ohne Namenswechsel (2409) | Z. 169–170 | gedeckt |
| 21 | Name freigemacht und neu belegt (2410) | Z. 156 | gedeckt; Schlüssel nicht genannt (F-2) |
| 22 | zwei Quellen auf ein Ende (2411) | Z. 164–166 | gedeckt |
| 23 | Rotation (2412) | Z. 157–158 | gedeckt |
| 24 | Löschen und Neuanlegen im selben Commit (2413) | Z. 171–172 | gedeckt |
| 25 | Art des Eintrags, Verzeichnis (2414) | Z. 180–182 | gedeckt |
| 26 | immutable Datei umbenannt (2415) | Z. 183 | gedeckt |
| 27 | Basis kein Vorfahr (2416) | Z. 147–148 | gedeckt |
| 28 | `--staged` (2417) | Z. 184–185 | gedeckt |
| 29 | fail-closed (Zwischen-Commit) (2418) | Z. 186, nur „Commit der Range“ | **halb** — Baum eines Eltern-Commits vor BASE und geänderter Teilbaum fehlen, ebenso „ohne Schlüssel Exit 0“ (F-1) |
| 30 | ohne Schalter byte-identisch (2419) | Z. 136–137, Welle-Probe | **halb** — „kein Commit zwischen BASE und HEAD wird gelesen“ hat keinen Test; trägt sich mit dem fehlenden Teil von 29 (F-1) |

## Befunde

### F-1 — Sechs Kriterien haben keinen oder nur einen halben Testfall — **V1**

- Am Gegenstand: Tabelle oben, Kriterium 6, 10, 11, 17, 29, 30.
- Fehlende Testfälle, je einer genügt:
  - **6:** Move in einem Commit, Nachzug in einem späteren, Range über beide ⇒ kein Befund.
  - **10:** reiner Move A → B, Inhaltsänderung an B in einem späteren Commit, Nachzug A → B ⇒ kein Befund.
  - **11:** X → Y auf der einen, Y → Z auf der anderen Seite eines Merges (keine Seite Vorfahr der anderen),
    Nachzug X → Z ⇒ Befund. Das ist der Fall, der eine ordnungsfreie Verkettung über alle Kanten
    (ADR-0115 Entscheidung 11) rot macht; ohne ihn bleibt die Abstammungs-Bedingung ungetestet.
  - **17:** beide Schlüssel, A → B mit anderem Dateinamen, Umbiegen von A auf eine andere Datei mit dem
    Dateinamen von A ⇒ Befund; mit `ignore-link-targets` allein kein Befund.
  - **29/30:** ein unlesbarer Baum (a) eines Zwischen-Commits, (b) des Eltern-Commits vor BASE, an dem eine
    Merge-Seite abzweigt, (c) ein geänderter Teilbaum ⇒ mit Schlüssel Exit 2, **ohne** Schlüssel Exit 0.
    Die Hälfte „ohne Schlüssel Exit 0“ belegt zugleich Kriterium 30 („kein Commit zwischen BASE und HEAD
    wird gelesen“). Dieselben Fälle führt ADR-0115 §Fitness Function.
- **Plan-Änderung:** die fünf Testfälle in §3 aufnehmen.

### F-2 — Kriterium 18 verlangt die Negativ-Fälle in beiden Schlüssel-Kombinationen; der Plan fährt nur ausgewählte — **V1**

- Kriterium 18 (Z. 2407) sagt für **jedes** Negativ-Kriterium „ohne `vcs.ignore-link-targets`“ zu, dass es
  mit beiden Schlüsseln gilt, mit der einen Ausnahme gleichen Dateinamens ohne teilnehmendes Ziel. Der Plan
  hat einzelne Fälle mit beiden Schlüsseln (Z. 157–170), keine Regel für die übrigen; die Fixtures
  Z. 138–141, 143, 151–156, 173–181 nennen keine Kombination.
- Der billigste Weg, der real trägt: jedes Negativ-Fixture läuft tabellengetrieben ein zweites Mal mit
  beiden Schlüsseln, Erwartung nach Kriterium 18. Ein eigener Testfall je Kriterium wäre Doppelarbeit.
- **Plan-Änderung:** §3 vor die Testfall-Liste ein Satz: „Jeder Negativ-Fall läuft mit `follow-renames`
  allein und mit beiden Schlüsseln; Erwartung nach Kriterium ‚Boundary (Negativ-Kriterien mit beiden
  Schlüsseln)‘.“ Bei Z. 156 (freigemachter Name) die Kombination nennen wie im Kriterium.

### F-3 — DoD und Abgrenzung widersprechen der Dateitabelle — **V1**

- DoD-Punkt 3 (Z. 79–80) verlangt, dass die **READMEs** den Schlüssel nennen; §3 Z. 115 legt sie in die
  Release-Prep (Verdikt welle-96 F-10, [`AGENTS.md`](../../AGENTS.md) §5 Regel 17).
- §1 Z. 52–55 schließt „Anforderung, Verfahren … ändern“ aus; §3 Z. 118 und der Kopf (Z. 21–25) ändern den
  Nebensatz in Schritt 2 der Spezifikation (Verifikation V-1). Die Änderung ist richtig, die Abgrenzung
  nennt sie nicht als Ausnahme.
- §1 Z. 58 „Die Grenze ‚gleichnamige andere Datei‘ — Welle-Out-of-Scope“ gilt seit Entscheidung 10 nur
  noch für Ziele, die nicht teilnehmen (so in Welle §6 und Lastenheft Z. 2350–2355).
- DoD-Punkt 1 (Z. 70–72) zählt nicht auf, was ADR-0115 Entscheidung 9, 12 und 13 dem Adapter zuweisen:
  die Bedingung an das Ende (in HEAD, nicht in BASE, eine Quelle), „Basis kein Vorfahr ⇒ leer“ und die
  Index-Angleichung.
- **Plan-Änderung:** DoD-Punkt 3 ohne READMEs; §1 erster Punkt „außer dem Nebensatz zur Index-Enumeration
  (V-1), der mit dem Code fällt“; §1 dritter Punkt „für Ziele, die an keiner Umbenennung teilnehmen“;
  DoD-Punkt 1 um die drei Punkte ergänzen.

### F-4 — Die Index-Angleichung ändert ein Ergebnis ohne den Schalter; Welle und Plan sagen „nichts Neues“ — **V1**

- Nach Lesart des Codes (nicht gefahren): `AllPaths` nimmt im gestagten Modus jeden Index-Eintrag
  ([git.go](../../internal/adapter/driven/git/git.go) Z. 95–98). Ersetzt ein gestagter Gitlink eine immutable
  Datei an ihrem Pfad, liegt der Pfad in beiden Mengen, `vcsModified` liest ihn über `fileFromIndex`
  (Z. 331–356) als Blob einer Commit-Kennung, und der Lauf endet mit einem Fehler (Exit 2). Nach der
  Angleichung fehlt der Pfad in der Index-Menge, `vcsDeleted` meldet `core-drift-vcs` (Exit 1) — dieselbe
  Antwort wie im `--range`-Modus, wo `walkTree` den Gitlink überspringt (Z. 173–175). Für jeden anderen
  Gitlink unter `vcs.paths` bleibt die Ausgabe gleich (er war und bleibt frei).
- Das ist eine Angleichung an Rang 2 und keine Lockerung — von Exit 2 zu Exit 1, der Befund wird sichtbar.
  Aber es ist eine Änderung **ohne** `vcs.follow-renames`. Der Welle-Closure-Trigger sagt „Ohne den
  Schalter nichts Neues … darf keinen bisher grünen oder roten Lauf verschieben“
  ([welle-96](../plan/planning/welle-96-vcs-folgt-umbenennung.md) Z. 61–65), Kriterium 30 nennt nur
  `--range`.
- **Verdikt:** gilt; der Plan muss die Abweichung benennen, damit Welle-Probe und Verifier sie nicht als
  Bruch lesen.
- **Plan-Änderung:** §3 Z. 108 und Testfall Z. 187 um den Fall „gestagter Gitlink ersetzt eine immutable
  Datei: vorher Exit 2, jetzt `core-drift-vcs` wie im `--range`-Modus“ ergänzen; Welle §3 den Satz
  „Ohne den Schalter nichts Neues“ um diese eine, gewollte Abweichung im gestagten Modus einschränken
  (die Blackbox-Fixtures treffen sie nicht).

### F-5 — Die Spezifikation sagt, ein durch einen Gitlink ersetztes immutables Datei bleibe unbemerkt; der Code meldet es — **V1** (Altbestand)

- §DC-FA-VCS-001.a Schritt 2 ([spezifikation.md](../../spec/spezifikation.md) Z. 2003–2007): „mit dem
  benannten Preis, dass eine immutable Datei, die durch einen Gitlink **ersetzt** wird, unbemerkt bleibt“.
  Nach Lesart des Codes fehlt der Pfad dann in der HEAD-Menge, und `vcsDeleted`
  ([vcs.go](../../internal/hexagon/core/rules/vcs.go) Z. 92–105) meldet „gelöscht oder umbenannt“ — im
  `--range`-Modus heute, im gestagten nach F-4. Die Abweichung liegt in der sicheren Richtung.
- Nicht gemessen; ein Test entscheidet es. Der Satz steht im selben Absatz wie der Nebensatz, den der Slice
  ohnehin ändert.
- **Plan-Änderung:** ein Testfall „immutable Datei durch Gitlink ersetzt ⇒ `core-drift-vcs`, in beiden
  Modi“; bestätigt er den Code, fällt der Satz im selben Commit wie der Nebensatz (Plan-Nachtrag zu Z. 118,
  Kopf „Berührte Spec-Stellen“). Kein neuer Grund-Code, keine Anforderung berührt.

### F-6 — Form der Port-Methode und Lesart des Schlüssels: entschieden, im Plan nicht festgehalten — **V1**

Keine neue Entscheidung: ADR-0115 Entscheidung 13 teilt Kette, Mehrdeutigkeit und die Bedingung an das Ende
dem Adapter zu, das Gleichsetzen dem Kern. Daraus folgt, was der Plan nennen sollte, damit der Review dagegen
prüfen kann:

- **Form:** die Methode liefert die Abbildung der **teilnehmenden** Umbenennungen als Zuordnung Quelle →
  Ende (repo-relativ, `/`-getrennt), schon nach den Bedingungen von Entscheidung 9 gefiltert; leer, nicht
  Fehler, bei „Basis kein Vorfahr“ und im gestagten Modus ohne HEAD.
- **Invariante, die der Kern prüfen darf:** keine Quelle ist zugleich Ende. Sie folgt aus „Ende nicht in
  BASE“ (jede Quelle liegt in BASE); dann ist die Gleichsetzung im Kern eine reine Nachschlage-Funktion —
  Quelle ⇒ ihr Ende, Ende ⇒ es selbst —, und die Iterationsreihenfolge der Zuordnung kann die Ausgabe nicht
  berühren ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Auflöser:** heute baut `CheckVCS` die Vereinigung der Pfad-Bäume nur mit `ignore-link-targets`
  (vcs.go Z. 53–56); mit `follow-renames` allein braucht er sie ebenso.
- **Konfigurations-Rand:** `follow-renames` wird gelesen wie `ignore-link-targets`
  ([configyaml.go](../../internal/adapter/driven/configyaml/configyaml.go) Z. 163) — dieselbe Lesart für
  einen falschen Typ und für einen Schlüssel ohne Wert.
- **Plan-Änderung:** §3 Zeile Port (Z. 107) um Form und Invariante, Zeile Regel (Z. 111) um den Auflöser,
  Zeile Konfigurations-Rand (Z. 110) um „wie `ignore-link-targets`“.

### F-7 — Risiken in §6: aktuell, eines fehlt — **V1**

- Risiko 1 (Auflöser) und 2 (racily-clean-Fixtures) sind aktuell; Risiko 2 gilt wegen der vielen Fixtures
  mit gleichen Blobs (Kopien, Rotation, Tausch) eher stärker als bei der Planung.
- Risiko 3 (Laufzeit) misst „100 Commits mit und ohne Schalter“. Seit Entscheidung 13 liest der Schalter
  auch Bäume von Eltern-Commits vor BASE; die Messung trägt das mit, wenn die Range einen Merge mit einer
  vor BASE abzweigenden Seite enthält — sonst misst sie den teuren Fall nicht. Kleine Schärfung, kein neues
  Risiko.
- Es fehlt das Risiko aus F-4 als benannter Ausgang: die gewollte Abweichung ohne Schalter muss am Ende
  belegt und nicht als Bruch gelesen werden.
- **Plan-Änderung:** Risiko 3 „die Range enthält einen Merge, dessen eine Seite vor BASE abzweigt“; ein
  Risiko zu F-4, Ausgang *entfallen*, wenn Testfall und Welle-Satz nach F-4 stehen.

## Bestätigt, mit akzeptiertem Negativ

- **Dateitabelle gegen den Port-Vertrag (Prüffrage 1).** Ein einziger Implementierer außerhalb des Adapters:
  der Test-Doppelgänger in `internal/hexagon/core/rules/vcs_test.go` (gesucht über alle `_test.go`;
  `commits_test.go` nutzt denselben Typ im selben Paket). Kein weiterer Spiegel des Schlüssels: die Suche
  nach `IgnoreLinkTargets` trifft im Code nur `model/config.go`, `configyaml`, `config_template.go`,
  `vcs.go` und Tests — alle in der Tabelle. Die Abdeckungs-Dateien sind zwei
  (`abdeckung-tests.md`, `abdeckung-e2e.md`); die zweite ändert sich nur, wenn `tools/image-test.sh` eine
  Phase bekommt, und der Plan sieht keine vor — die Probe am Image trägt die Welle.
- **Reihenfolge aus Entscheidung 14 (Prüffrage 4).** DoD-Punkt 3 (Z. 82–88) und §3 Z. 113 binden den
  Profil-Schalter an den Commit, der ADR-0115 `Accepted` hebt; die Zeiger in `AGENTS.md` §3.5 und
  `harness/README.md` stehen dabei.
- **R3- und R4-Fälle (Prüffrage 3).** Tausch über Zwischennamen (Z. 151–155) mit richtiger Begründung und
  ohne Rot-Behauptung für die BASE-Bedingung; freigemachter Name (Z. 156) und die nachsichtige Hälfte mit
  beiden Schlüsseln (Z. 159–163) wie Lastenheft Z. 2338–2349; Rotation (Z. 157–158) in beiden
  Kombinationen; zwei Quellen (Z. 164–166) mit Rot-Beleg; Löschen und Neuanlegen im selben bzw. in einem
  späteren Commit (Z. 171–179) mit der Commit-Bedingung. Kein Testfall widerspricht 0.114.3.
- **Reihenfolge-Unabhängigkeit.** Z. 191–194 fährt einen Wechsel der Kurz-SHA-Ordnung, nicht nur zehn
  gleiche Läufe — die Verifikation V-2 ist eingelöst.
- **Rückführung „Blob-Kennungen“ (Z. 211–213).** Bleibt kaum erreichbar (Kennungen stehen in Baum- und
  Index-Einträgen); stehen lassen ist billiger als streichen.

## Nicht geprüft

- Kein Code gefahren; F-4 und F-5 stehen auf gelesenem Code.
- Die Testfälle sind gegen den Text gelesen, nicht gegen gebaute Fixtures.
- Gates: nur `make doc-check` über diesen Report.

## Zusammenfassung

| Befund | Verdikt | Träger der Korrektur |
|---|---|---|
| F-1 sechs Kriterien ohne oder mit halbem Testfall | V1 | Plan 292 §3 |
| F-2 Negativ-Fälle in beiden Schlüssel-Kombinationen | V1 | Plan 292 §3 |
| F-3 DoD/Abgrenzung gegen Dateitabelle | V1 | Plan 292 §1, §2 |
| F-4 Index-Angleichung ändert ein Ergebnis ohne Schalter | V1 | Plan 292 §3, Welle §3 |
| F-5 Spezifikation: Gitlink-Ersatz „unbemerkt“ | V1 (Altbestand) | Plan 292 §3, Kopf |
| F-6 Form der Port-Methode, Auflöser, Schlüssel-Lesart | V1 | Plan 292 §3 |
| F-7 Risiken | V1 | Plan 292 §6 |

Kein V2, kein V3: ADR-0115 bleibt unverändert. Keine Entscheidung des Auftraggebers nötig.

## Plan-Änderungen für den Planner

1. §3 Testfälle aus F-1 aufnehmen (Kriterium 6, 10, 11, 17, 29/30).
2. §3 vor der Testfall-Liste: jeder Negativ-Fall läuft mit `follow-renames` allein und mit beiden
   Schlüsseln, Erwartung nach Kriterium 18; Z. 156 nennt die Kombination (F-2).
3. DoD-Punkt 3 ohne READMEs; DoD-Punkt 1 um Bedingung an das Ende, „Basis kein Vorfahr ⇒ leer“ und
   Index-Angleichung (F-3).
4. §1: Ausnahme für den Nebensatz in Schritt 2 (V-1); Grenze „gleichnamige andere Datei“ nur für Ziele ohne
   Umbenennung (F-3).
5. §3 Z. 108 und Z. 187: gestagter Gitlink ersetzt eine immutable Datei — vorher Exit 2, jetzt
   `core-drift-vcs`; Welle §3 „Ohne den Schalter nichts Neues“ um diese Abweichung einschränken (F-4).
6. §3: Testfall „immutable Datei durch Gitlink ersetzt, beide Modi“; bestätigt er den Code, fällt der Satz
   Z. 2006 der Spezifikation mit dem Nebensatz (Plan-Nachtrag, Kopf) (F-5).
7. §3 Zeilen Port, Regel, Konfigurations-Rand: Form der Abbildung, Invariante, Auflöser, Lesart des
   Schlüssels (F-6).
8. §6: Risiko 3 auf eine Range mit vor BASE abzweigender Merge-Seite; neues Risiko zu F-4 (F-7).
