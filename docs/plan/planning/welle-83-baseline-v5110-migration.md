# Welle welle-83-baseline-v5110-migration: Acht Wellen Rückstand, und eine davon ist die Antwort auf unseren eigenen CR

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-83-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug. Ob die Welle mit einem Release
schließt, entscheidet das Delta-Audit — Harness-Migrationen sind
konsumentensichtbar nur dort, wo sie das Produkt berühren.

**Verantwortlich:** pt9912. **Datum:** 2026-08-23.

---

## 1. Welle-Ziel

Der Baseline-Pin steht auf
[`v5.9.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v5.9.0)
(Kurs-Welle **86**); der Kurs steht bei
[`v5.11.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v5.11.0)
(Kurs-Welle **94**). **Acht Wellen über zwei Minors** — das ist keine
Pin-Zeile, sondern eine Migration mit eigenem Audit.

| Kurs-Welle | Thema | Warum es uns angehen könnte |
|---|---|---|
| 87 | Team-Sim in Modul-12-Form | vermutlich folgenlos (kein team-sim in diesem Repo) |
| 88 | Vier unbelegte Aussagen, sieben Verdikte | Review-Form — berührt den Reviewer-Skill |
| 89 | **Form ist kein Beleg** | trifft die Klasse, an der welle-82 achtmal gescheitert ist |
| 90 | **Ab `Accepted` zählt jede Zeile** | berührt `adr-check` / Modul `vcs` und die ADR-Immutabilität |
| 91 | Das Kurs-Repo sagt, wie an ihm gearbeitet wird | Vorbild-Form, keine Regel für Adopter |
| 92 | Zwei Gewohnheiten werden Invarianten | zwei Form-Invarianten via `structure` — Konsumenten-Konfiguration |
| 93 | **AGENTS.md §4 wird die Autorität über die Targets** | berührt Modul `targets`, `gate-consistency` und unsere §4-Tabelle |
| 94 | **Eine Rangliste ordnet, jetzt deckt sie auch ab** | die **Vollständigkeits-Zusage**; Auslöser war unser eigener CR |

**Welle 94 ist die Antwort auf einen Konsumenten-CR dieses Repos.** Das Kurs-
CHANGELOG nennt ihn als Auslöser. Die Zusage, die daraus entstand — *jede
Regel, der ein Agent folgen muss, steht in einer gerankten Quelle oder im
Konventionsspeicher; Artefakte außerhalb dürfen verweisen und ausführen, aber
nichts festlegen* — gilt nach dem Bump **für uns selbst**, und wir wissen
bereits von **einer** Verletzung ([slice-127](done/slice-127-claude-md-pointer.md)).
Sie einzusammeln ist Folge der Welle, nicht ihr Inhalt.

**Die Reihenfolge ist nicht verhandelbar:** erst pinnen, dann dem Kanon folgen.
Andernfalls führte `AGENTS.md` Stand 94, während der Konventionsspeicher
`v5.9.0` pinnt — genau die Drift, die das Delta-Audit finden soll.

## 2. Trigger (Welle startet)

`v5.11.0` ist veröffentlicht (geprüft: der Tag trägt `Stand: Kurs-Welle 94` und
den Vollständigkeits-Absatz); welle-82 ist geschlossen, `in-progress/` frei.

## 3. Closure-Trigger (Welle schließt)

- **Bundle vendored und offline verifiziert:** `.harness/baseline/v5.11.0/`
  aus dem Release-Asset, `SHA256SUMS` geprüft — netzlos, ohne Handanlegen an
  den Bäumen.
- **Pin gehoben** als nächster Eintrag der Pin-Serie
  ([`MR-011`](../../../harness/conventions.md#mr-011)-Kette), **alle**
  pin-gebundenen Verweise retargetet
  ([`MR-021`](../../../harness/conventions.md#mr-021)) — und die
  Drei-Klassen-Prüfung aus `BEO-008` gefahren, nicht nur das Pfad-`grep`.
- **Delta-Audit:** je Kurs-Welle 87–94 **eine Antwort** — konform ohne
  Handlung (mit Beleg), Handlung nötig (mit Slice), oder nicht anwendbar (mit
  Begründung). Keine Welle ohne Zeile.
- **Etappe C abgearbeitet oder ausgewiesen:** was das Audit schneidet, ist
  entweder in dieser Welle erledigt oder als Folge-Welle benannt — nicht
  stillschweigend offen.
- `make fullbuild` grün (Exit explizit); Ergebnisnotiz `welle-83-results.md`
  mit Register-Lese-Schritt.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-128](done/slice-128-baseline-v5110-vendoring.md) | **Etappe A:** Bundle `v5.11.0` vendored, Pin-Hebung als `MR-`Eintrag, pin-gebundene Verweise gehoben, Alt-Baum entfernt |
| [slice-129](done/slice-129-baseline-v5110-delta-audit.md) | **Etappe B:** Delta-Audit über die Kurs-Wellen 87–94, je Welle eine Antwort; **schneidet Etappe C** |
| [slice-130](done/slice-130-lastenheft-historie-form.md) | **Etappe C-1** (vom Audit geschnitten): Historie-Form auf vier Spalten, und die eigene Strenge deklarieren — aus Kurs-Welle 90 |
| [slice-131](open/slice-131-reviewer-skill-waisen.md) | **Etappe C-2** (vom Audit geschnitten): die Waisen im Reviewer-Skill nach `AGENTS.md` umziehen — aus Kurs-Welle 94 |

**Etappe C ist vom Audit geschnitten**, nicht hier geraten — dieselbe Form wie
in der v5.6.0-Migration; nachgeführt mit einer Drift-Log-Zeile, wie bei der
Eröffnung angekündigt. Zwei Slices, aus den **zwei** Wellen mit
Handlungs-Antwort; die übrigen sechs sind belegt folgenlos.

**Der dritte Fundort des Vollständigkeits-Zensus hat bereits seinen Slice:**
[slice-127](done/slice-127-claude-md-pointer.md) ist in Arbeit. Damit ist auch seine offene Reihenfolge-Frage beantwortet —
`CLAUDE.md` ist **einer von drei** Fundorten, nicht der einzige.

## 5. Abhängigkeiten

- **Entblockt:** [slice-127](done/slice-127-claude-md-pointer.md) — die
  Regel, der er folgt, ist mit Etappe A gepinnt; er läuft seither.
- **Wird blockiert von:** nichts. Reihenfolge innerhalb: A vor B (das Audit
  liest den neuen Baum), C nach B.

## 6. Out-of-Scope für diese Welle

- **Kein Release, solange das Audit keine Produkt-Konsequenz findet.** Ein
  Harness-Bump ist für Konsumenten nur dort sichtbar, wo er Modul-Verhalten
  oder Konfiguration ändert.
- **slice-127 wird nicht mitgezogen.** Er ist die *erste Anwendung* der neuen
  Zusage und bleibt ein eigener Slice nach der Welle — sonst vermischen sich
  Pin und Konsequenz in einem Diff.
- **Keine Zwischenstufe über `v5.10.0`.** Gepinnt wird auf den aktuellen Tag;
  das Audit läuft über die Wellen, nicht über die Releases.
- **Keine vorgezogene Umsetzung** einer Kurs-Welle vor dem Pin. Der Kanon
  gilt ab dem gehobenen Pin, nicht ab dem Lesen.
