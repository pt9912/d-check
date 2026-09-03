# Review-Report: <slice-NN | PR-Ref> — <YYYY-MM-DD>

> **Template-Hinweis.** Vorlage für einen Review-Report (das
> Übergabe-Artefakt Reviewer → Implementer, Modul 8/10). Kopiere
> nach `docs/reviews/<YYYY-MM-DD>-<slice-oder-diff-ref>.md`, ersetze
> `<Platzhalter>` und lösche diesen Block. Ein Report pro Lauf —
> Folgeläufe bekommen eine neue Datei, keine Überschreibung
> (Auditierbarkeit).

**Review-Art:** Plan | Design | Code — *wogegen* geprüft wird:
Plan-Review gegen Spec/ADR, Design-Review gegen Architektur,
Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** <Slice-ID / Diff-Range / Commit-Hash>

**Skill:** `.harness/skills/reviewer.md` @ <Version/Commit> · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad, existiert im Ziel-Repo ggf. nicht) -->
**Modell:** <Modell-ID> · **Datum:** <YYYY-MM-DD>

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- <Slice-Plan / Plan-Dokument>
- <aktive ADRs, z. B. ADR-<NNNN>>
- <berührte `LH-*`-IDs>
- `AGENTS.md` (Hard Rules)

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Felder unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
Baseline-Regelwerk `modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — <Kurztitel>

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: <ADR-ID, LH-ID, Hard-Rule-Name oder "Maintainability">
- `pfad`: <Datei:Zeile>
- `befund`: <1–2 Sätze, beobachtbar, ohne Lösungsvorschlag>
- `verifizierbar`: ja/nein — <welcher Gate-Lauf würde es bestätigen?>
- `klasse`: <stabile Kurz-Bezeichnung des Fehlermusters, z. B.
  „Tie-Break in sortierender Operation nicht dokumentiert">

## Negativbefunde

<!--
Eine Zeile pro betrachtetem Bereich. Ohne diesen Block ist "keine
Findings" nicht von "nicht geprüft" unterscheidbar (Modul 10
§Reviewer berichtet auch, was er nicht gefunden hat).
-->

- geprüft, ohne Befund: <Verzeichnis/Bereich>
- geprüft, ohne Befund: <Verzeichnis/Bereich>

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | <n> |
| MEDIUM | <n> |
| LOW | <n> |
| INFO | <n> |

**Finding-Klassen dieses Laufs:** <klasse-1> · <klasse-2>

<!--
Die Klassen-Zeile ist der Übergabepunkt in den Steering-Loop-ZÄHLER. Sie
wird bei der Slice-Closure (§7) ins Beobachtungs-Register eingetragen und
dort gezählt; bei 3x wird der Reviewer-Skill geschärft (Modul 10 §Pflege).

DIE KLASSEN-BEZEICHNUNG MUSS ÜBER LÄUFE HINWEG STABIL SEIN. Dieser Report
kennt das Register NICHT (er ist Lauf-Beleg) — die Zuordnung zur BEO-<NNN>
passiert erst bei der Slice-Closure und braucht den wiedererkennbaren Namen.
Ab dann zitiert die Closure die Kennung; dort ist die Bezeichnung nur noch
Label. Niemand muss alte Reports lesen: die Häufung steht im Register, nicht
in einem Archiv-Scan.
-->

## Verdikt

**Merge-blockierend:** ja | nein — HIGH und MEDIUM blockieren
typischerweise; eine Abweichung davon wird hier begründet, nicht
still entschieden.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
