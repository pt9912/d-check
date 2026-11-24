# Trigger-Audit welle-96 — Zweige ADR und Hard Rule, Lese-Schritt des Registers (Architect)

**Rolle:** Architect (Modul 8, §Rollen-Sequenz für eine Welle: Closure-Schritt 2 und 3b, Planner → Architect →
Planner). Rolleninhaber: pt9912.
**Gegenstand:** [welle-96](../plan/planning/done/welle-96/welle-96-vcs-folgt-umbenennung.md) mit
[slice-291](../plan/planning/done/welle-96/slice-291-vcs-umbenennung-anforderung-und-adr.md) und
[slice-292](../plan/planning/done/welle-96/slice-292-vcs-umbenennung-modulkern.md), Range `v0.92.0..HEAD`
(60 Commits, ab der Eröffnung `e81aebd4`; die Teil-Range `6aca94b0^..HEAD` umfasst nur slice-292).
**Eingang:** Audit-Vorlage des Planners; die Verifikationen der beiden Slices
(slice-291 (Verifikation, im Archiv von welle-96),
slice-292 (Verifikation, im Archiv von welle-96)). Die repo-weite Verifikation der Welle entsteht
parallel und ist nicht Eingang.
**Nicht Gegenstand:** die Zweige Carveout und bootstrap-aware Gate (prüft der Planner).
**Datum:** 2026-11-24.
**Geändert wurde nichts:** keine ADR, keine Hard Rule, kein Register-Eintrag; dies ist nur das Verdikt.

## Gesamturteil

Kein Re-Evaluierungs-Trigger und kein Auflösungs-Trigger ist eingetreten. **Keine Folge-ADR**, **keine
Hard-Rule-Zeile** wird entfernt. Die einzige Änderung an der Gate-Konfiguration ist eine Lockerung der
Prüfregel des eigenen ADR-Gates (`follow-renames: true`); sie ist durch
[ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 7 getragen und landete im selben Commit wie
deren Übergang nach `Accepted` (`10e28121`, Entscheidung 14) — [`AGENTS.md`](../../AGENTS.md) §3.6 ist erfüllt.

Für den Register-Eintrag
[`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../plan/planning/observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md)
(3×) schlage ich **`verkörpert`** vor: eine Dokumentations-Regel 18 in [`AGENTS.md`](../../AGENTS.md) §5 mit
Volltext unter `harness/rules/dokumentations-regeln/`, *permanent*, `seit welle-96` (Teil 2).

**Braucht keine Entscheidung des Auftraggebers.** Die Regel richtet sich an Planner und Implementer und ändert
weder Anforderung noch Gate.

Gate-Lage über die Welle:
`git diff --stat v0.92.0..HEAD -- .d-check.yml Makefile .golangci.yml tools/coverage-gate.sh` ⇒
`.d-check.yml | 5 ++++-` (4 Einfügungen, 1 Löschung), sonst nichts; für `6aca94b0..HEAD` dasselbe. Der Diff
ist genau der `vcs`-Block: ein Kommentar-Satz und `follow-renames: true`. Keine Schwelle, keine Linter-Strenge,
kein `coverage-gate`; `make gates` ist am Ende von slice-292 grün (Verifikation slice-292, Punkt 4:
Coverage 94,90 % gegen die Schwelle 93 %).

## Teil 1 — ADR

| ADR | Trigger | eingetreten? | Beleg |
|---|---|---|---|
| [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) | (1) das Baseline-Regelwerk ändert die Ausnahme für den Referenz-/Pfad-Nachzug | **nein** | Pin `v6.18.1` über die Welle unverändert; `.harness/` hat in `v0.92.0..HEAD` keinen Diff |
| ADR-0115 | (2) ein Review findet einen Nachzug über eine Umbenennung, der die Aussage einer ADR verschoben hat | **nein** | die Reviews R1–R4 zu slice-291 und R1–R2 zu slice-292 betrafen Vertragstext, Tests und Code, keinen Nachzug im Bestand; der erste Lauf mit Schalter über `HEAD~100..HEAD` meldete 0 Befunde |
| ADR-0115 | (3) eine benannte Grenze meldet bei einem Adopter wiederholt einen Nachzug als Drift | **nein** | noch nicht veröffentlicht; kein Adopter-Lauf mit dem Schalter |
| ADR-0115 | (4) `make adr-check` über 100 Commits mit Schalter mehr als doppelt so lange wie ohne | **nein** | gemessen in der Verifikation slice-292 (Risiko 3): Median 1,14 s ohne, 1,56 s mit, Faktor **1,37**, über eine Range mit einer vor BASE abzweigenden Merge-Seite |
| ADR-0115 | (5) die Bibliothek bietet eine Umbenennungs-Erkennung, die Mehrdeutigkeit meldet und auf einem unlesbaren Objekt abbricht | **nein** | `go.mod` über die Welle unverändert (go-git v5.19.3) |
| [ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) | (1) Baseline ändert die Ausnahme · (2) ein Review findet einen Nachzug, der die Aussage verschoben hat | **nein** | wie ADR-0115 (1) und (2). ADR-0115 **erweitert** sie um einen eigenen Schlüssel; Entscheidungen 1–4 und die Wirkung von `ignore-link-targets` bleiben gültig, die Tests in `vcs_pfad_nachzug_test.go` laufen unverändert grün (Verifikation slice-292, Kriterium 5). Kein `supersedes` |
| [ADR-0024](../plan/adr/0024-vcs-immutable-gate.md) | (1) ADR-Vorlage ändert `## Geschichte` oder die Status-Zeile | **nein** | Baseline-Pin unverändert |
| ADR-0024 | (2) Bedarf, die git-Garantie in `make gates` zu führen | **nein** | `adr-check` bleibt am Commit-/Diff-Bindepunkt; `Makefile` unverändert |
| ADR-0024 | (3) eine zweite git-bedürftige Befund- oder Reparatur-Klasse (z. B. git-basierte Move-Erkennung) ⇒ nutzt denselben VCS-Port | **berührt, nicht eingetreten** | Die Welle baut eine git-basierte Umbenennungs-Erkennung, aber **keine neue Befund- oder Reparatur-Klasse**: kein neuer Grund-Code, `core-drift-vcs` bleibt diagnose-only, die Erkennung speist nur die Normalisierung des bestehenden Vergleichs. Sie folgt der Vorgabe des Triggers schon: derselbe VCS-Port (`Renames`), Zusage in `DC-FA-VCS-001`. Eine Reparatur über Move-Erkennung (ADR-0008) bleibt offen und fände die Port-Methode vor |
| [ADR-0109](../plan/adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md) | (1) Force-Push auf den Default-Branch erlaubt · (2) GitHub liefert eine auflösbare `before`-Basis · (3) der Rückfall meldet wiederholt Altlast | **nein** | keine Workflow- und Ruleset-Änderung in der Welle. ADR-0115 Entscheidung 12 hält ihre Annahme: eine Basis, die kein Vorfahr ist — wie sie der Rückfall `origin/<default>..HEAD` liefern kann —, ergibt keine Umbenennung und **keinen** Exit 2; der Push-Lauf wird durch den Schalter nicht rot |

**§3.6 und `.d-check.yml`:** keine Schwelle; eine Lockerung der Prüfregel, getragen von ADR-0115 Entscheidung 7
und angenommen im selben Commit, der sie einschaltet (`git show --stat 10e28121`: `.d-check.yml`, ADR-0115,
Index, `AGENTS.md`, `harness/README.md`). Die verbleibenden Grenzen stehen in ADR-0115 §Konsequenzen und
im Lastenheft.

## Teil 1 — Hard Rules in `AGENTS.md` §3

| Regel | Auflösungs-Trigger | eingetreten? |
|---|---|---|
| §3.1 Docker/make-only | keiner benannt (Bestand) | — |
| §3.2 Suppression-Verbot | permanent | nein |
| §3.3 `git mv` + Inhalt = zwei Commits | permanent | nein — angewandt: die Lifecycle-Moves der Welle lassen die bewegte Datei unverändert (`R100` in `c0959fff`, `2d28467e`, `6aca94b0`, `274c6ba4`; mitreisende Verweis-Dateien sind nach §3.3 zulässig) |
| §3.4 Architektur sprachfrei, Straten nie abwärts | Closure-Datum: keiner; Rollen statt Technologie: permanent | nein — `spec/architecture.md` unberührt; die neuen Spec-Zeilen nennen weder ADR noch Slice (Review R1 zu slice-291, Negativbefund) |
| §3.5 ADRs nach `Accepted` immutable | durchgesetzt (`make adr-check`), kein Auflösungs-Trigger | — Die Welle hat die Zeile **ergänzt**, nicht aufgelöst: der Pointer nennt ADR-0115 neben ADR-0103 (`10e28121`). Die Aussage („ein reiner Pfad-Nachzug, der nur Link-Ziele ändert“) bleibt wahr; der Nachzug auf den neuen Namen ist ein solcher. ADR-0115 wurde mit `10e28121` `Accepted`, `make adr-check` über `HEAD~1..HEAD`, `HEAD~100..HEAD` und `--staged` grün |
| §3.6 Gates nicht ohne ADR lockern | permanent | nein — die eine Lockerung ist getragen (oben) |
| §3.7 Kommentare tragen eine der fünf Klassen | keiner benannt | — Der neue Kommentar im `vcs`-Block ist eine Zusage samt `ADR-*`-Feld |
| §3.8 Ein Modul verspricht nur über seine Scan-Menge | permanent | nein — angewandt: ADR-0115 Entscheidung 6 und das Lastenheft („Was gelesen wird“) beantworten die Frage für die Bäume außerhalb der Scan-Menge |
| §3.9 Action-Referenzen SHA-gepinnt | durchgesetzt (`make workflow-pins`) | — keine Workflow-Änderung |

**0 offen** in beiden Klassen.

## Teil 2 — Lese-Schritt: `BEO-ALL/plan-entscheid-steht-gegen-lastenheft` (3×)

**Gelesen:** `observation.md`, `state.md` (offen, 3×) und die drei Belege
[`slice-289`](../plan/planning/observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/evidence/slice-289.md),
[`slice-290`](../plan/planning/observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/evidence/slice-290.md),
[`slice-291`](../plan/planning/observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/evidence/slice-291.md).

**Die Klasse:** Eine Festlegung im Slice-Plan steht absoluter oder anders da, als Lastenheft oder ADR sie tragen.
Weil der Code dem Plan folgt, sieht man die Abweichung nur noch im Review. Geschrieben haben sie **beide**
Rollen, die Pläne schreiben: der Planner zweimal (slice-289 Entscheid 9; slice-291 die Richtung von
Entscheidung 9 und ihre drei Fälle), der Implementer einmal (slice-290, Plan-Änderungs-Notiz). Gefunden hat sie
jedes Mal ein Review; die Abhilfe in slice-291 war, die Stelle des Lastenhefts zu **zitieren**, statt sie
nachzuerzählen.

**Vorschlag: `verkörpert`**, Zielort Hard Rule (Dokumentations-Regel in `AGENTS.md` §5). Begründung:

- **Entstehungsort.** Die Festlegung entsteht beim Schreiben am Plan, und zwar ohne Command: der Planner zieht
  Pläne nach Verdikten und Reviews von Hand nach, der Implementer schreibt Plan-Änderungs-Notizen in
  `implement-slice` Schritt 12–16. `AGENTS.md` ist der einzige Ort, den beide Rollen in **jedem** Lauf lesen
  (`CLAUDE.md` importiert es). Ein Satz in `.claude/commands/plan-welle.md` griffe nur bei der Eröffnung
  (Schritt 6), nicht bei den Nachzügen, aus denen zwei der drei Belege stammen. Ein Satz in
  `.claude/commands/implement-slice.md` griffe nur beim Implementer. Der Reviewer-Skill (`.harness/skills/reviewer.md`)
  fand die Fälle bereits ohne eigenen Anker; er sitzt hinter dem Entstehungsort.
- **Form der Nachbarn.** Die Regeln 13–16 in §5 sind dieselbe Art (Urteil beim Schreiben, aus dem Steering
  Loop, *permanent*), und das Muster aus [ADR-0096](../plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md)
  gibt die Form: eine kurze Tabellenzeile, Volltext unter `harness/rules/dokumentations-regeln/`. Die
  Volltext-Datei ist zugleich der Zielort, gegen den die Anker-Paarung prüft (`seit welle-96` irgendwo in einer
  Datei ohne Suffix).
- **Kleinste tragende Form:** eine Zeile und eine kurze Datei. Kein neuer Reviewer-Anker, kein Gate, kein
  Folge-Slice.

**Mechanisierbar?** Nur die Form, und die schon heute. Wer zitiert statt nacherzählt, kann das Zitat mit
`<!-- d-check:cite spec/lastenheft.md:<von>-<bis> -->` deklarieren. Das Modul `citations` hält das Zitat dann
gegen die Quelle: `citations.scope` nimmt nur `done/`, `docs/reviews/` und die Baseline aus, aktive Pläne also
nicht. Ändert sich das Lastenheft, wird der Plan rot statt still alt. **Nicht** mechanisierbar ist, ob ein Satz
eine Festlegung ist und ob die zitierte Stelle sie trägt. Das bleibt Urteil, deshalb *permanent*.
[`BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft`](../plan/planning/observations/BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft/observation.md)
ist die Nachbar-Klasse eine Ebene höher (Anforderung gegen Lastenheft §5). Sie wird hiervon nicht gedeckt und
bleibt offen.

**Greift der Satz nur bei einer Plan-Festlegung?** Ja, wenn er an den Plan gebunden ist („Eine Festlegung im
Slice-Plan …“). Er gilt dann nicht für Prosa in Spec, ADR oder Handbuch, die eigene Regeln haben (Regel 3, 13, 16).
Die Abgrenzung „Festlegung“ gegen „Beschreibung“ steht im Volltext.

**Wortlaut 1 — Tabellenzeile in `AGENTS.md` §5**, nach Zeile 17:

```markdown
| 18 | Eine Festlegung im Slice-Plan nennt die Stelle des Lastenhefts oder der ADR, die sie trägt, und zitiert sie, statt sie nachzuerzählen; trägt keine Stelle sie, ist sie eine Entscheidung, keine Notiz. | [Volltext](harness/rules/dokumentations-regeln/18-plan-festlegung-zitiert-ihre-quelle.md) |
```

**Wortlaut 2 — Volltext** `harness/rules/dokumentations-regeln/18-plan-festlegung-zitiert-ihre-quelle.md`:

```markdown
# Eine Festlegung im Slice-Plan zitiert die Stelle, die sie trägt

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Eine **Festlegung** im Slice-Plan ist jeder Satz, nach dem sich Code, Test oder eine Vorlage an den
Auftraggeber richten wird: ein Entscheid in §3, eine Plan-Änderungs-Notiz, ein Testfall mit erwartetem
Ergebnis, eine Zusammenfassung der Richtung einer Änderung. Wer eine schreibt — Planner wie Implementer —,
beantwortet vorher: **Welche Stelle des Lastenhefts oder der ADR trägt sie?** Diese Stelle wird genannt
(verlinkte Kennung) und **zitiert**, nicht nacherzählt; eine Nacherzählung rundet ab, macht absolut oder
verschiebt die Richtung, und der Code folgt dem Plan, nicht der Quelle. Ein Zitat lässt sich als
`<!-- d-check:cite <datei>:<von>-<bis> -->` deklarieren; dann hält das Modul `citations` es gegen die
Quelle, und eine spätere Änderung des Lastenhefts macht den Plan rot statt still alt. **Trägt keine Stelle
die Festlegung, ist sie eine Entscheidung und keine Notiz:** Sie geht an den Architect (ADR) oder an den
Auftraggeber (Lastenheft), bevor Code ihr folgt.

**Grenzen:** Ob ein Satz eine Festlegung ist und ob das Zitat sie trägt, ist ein Urteil; `citations`
prüft nur, dass ein deklariertes Zitat noch in der Quelle steht, nicht, dass eine Festlegung eines hat.
Die Regel gilt dem Slice-Plan; Spec, ADR und Handbuch folgen ihren eigenen Regeln (3, 13, 16). Ob eine
neue Anforderung einem globalen Ausschluss des Lastenhefts widerspricht, ist eine andere Frage und hier
nicht gedeckt. *(Hard Rule aus dem Steering Loop,
[`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../../../docs/plan/planning/observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md),
seit welle-96; Auflösungs-Trigger: permanent.)*
```

**Wortlaut 3 — `state.md`** und Steering-Loop-Eintrag in `welle-96-results.md`:

- `state.md`: `**Stand:** verkörpert — Dokumentations-Regel 18 (AGENTS.md §5, Volltext harness/rules/dokumentations-regeln/18-plan-festlegung-zitiert-ihre-quelle.md), seit welle-96.`
- Steering-Loop-Eintrag: Regel · `BEO-ALL/plan-entscheid-steht-gegen-lastenheft` · Belege slice-289, slice-290,
  slice-291 · liegt in `harness/rules/dokumentations-regeln/18-plan-festlegung-zitiert-ihre-quelle.md`.

**Hinweise für den Commit:**
- Die neue Tabellenzeile steht hinter Zeile 327 von `AGENTS.md`. Die einzige Zeilen-Zitation in die Datei
  (`.harness/skills/reviewer.md` Z. 62, `AGENTS.md:287-287`) liegt davor und verschiebt sich nicht.
- Die Kennungen im Volltext sind verlinkt. `welle-96` steht dort als Herkunfts-Anker, wie in den Regeln 13 und 16.

## Nicht geprüft

- Die Zweige Carveout und bootstrap-aware Gate (Planner).
- Die repo-weite Verifikation der Welle (`make fullbuild`, Welle-Probe am Image). Sie entsteht parallel; Teil 1
  stützt sich auf die Verifikationen der Slices.
- Die übrigen Beobachtungen der Welle (Review- und Verifikations-Befunde aus slice-292). Ihre Zählung ist
  Sache des Planners. Keine von ihnen berührt einen Trigger in Teil 1.
- Ob `make doc-check` die Wortlaute in Teil 2 grün liest, wenn sie eingefügt sind. Gelesen wird hier nur dieser
  Report; die Blöcke stehen in Code-Zäunen.
