# Re-Review (Bestätigung): slice-085 (Etappe B) — Korrektur-Verifikation — 2026-08-02

**Review-Art:** Fokussierte Bestätigungs-Re-Review (nur die vier Korrekturen aus
`docs/reviews/2026-08-02-slice-085-etappe-b-review.md` + Regressions-Sichtung).
**Gegenstand:** `docs/plan/planning/in-progress/slice-085-etappe-b-modul-delta.md`
§3 (Korrektur-Commit `ac0ec1f`). **Modell:** Opus 4.8. **Verfahren:** read-only,
je Korrektur am Quelltext belegt; keine Repo-Datei außer diesem Report geändert.

---

## Korrektur-Verifikation

- **F-1 / D-5 + C-6 — BESTÄTIGT.** D-5 (Z. 71) trennt jetzt korrekt: **Slice**-`Status:`
  = Lifecycle-Verzeichnis-Dublette → entfernen; **NB** das MR-`Status: Accepted` =
  „andere Achse (Akzeptanz, template-vorgeschrieben) → C-6, **nicht** entfernen".
  C-6 (Z. 59) führt `Status: Accepted` (Template Z. 21) als zu **ergänzendes** Feld,
  „Akzeptanz-Achse, kein Lifecycle-Zustand". Template belegt beide Achsen: Z. 8
  „der Zustand ist die Verzeichnis-Position, kein Status-Feld" (Lifecycle) vs. Z. 21
  `- **Status:** Accepted` (Akzeptanz). Scheinbarer Widerspruch korrekt aufgelöst.

- **F-2 / D-8 — BESTÄTIGT.** D-8 (Z. 74) sagt jetzt „führt die Body-Sektion
  `## Re-Evaluierungs-Trigger` bereits in **27/46** ADRs … → **nahe konform**",
  ältere ohne Sektion „**immutable** (grandfathered)". Zahl nachgezählt:
  `grep -rlE '^##+ .*Re-Eval' docs/plan/adr/[0-9]*.md` = **27**, Gesamt-ADRs = **46**.
  Stimmt exakt. Die falsche Ist-Aussage („ADRs ohne Trigger") ist getilgt.

- **F-3 / D-11 — BESTÄTIGT.** D-11 (Z. 77) neu ergänzt: „**AGENTS.md gegen
  `AGENTS.template.md` angleichen**", Quelle `modul-09-implementierung` §AGENTS.md-Regeln,
  Ziel `AGENTS.md` §1 + **slice-083 §2.7-D3** (Schritt gedeckt). §3-Spot-Check-Eröffnung
  (Z. 47–48) nachgezogen: „keine neue Pflicht … **außer** dem AGENTS.md-Angleich
  (`modul-09` §AGENTS.md-Regeln → D-11)" — der frühere Selbstwiderspruch ist geheilt.

- **F-4 / C-7 — BESTÄTIGT.** C-7 (Z. 60) zitiert jetzt kanonisch
  `modul-13-quality-gates` §Guard-Härtung (**Z. 189**, „auch
  `grundlagen-durchsetzungsschicht`" nur nachrangig). Quelle bestätigt: modul-13
  Z. 165 `### Guard-Härtung`, Z. 189 „Die Grenz-Zeile wird mitgezogen".

## Regressions-Sichtung

- **Keine Regression.** Kopfzahl konsistent: §3-Kopf (Z. 45) „**19 Findings**
  (8 → Etappe C, 11 → Etappe D)"; gezählt §3.1 = 8 C-Zeilen, §3.2 = 11 D-Zeilen = 19.
  `make gates` **grün** (274 Dateien, 0 Befunde; `links`/`anchors`/`ids` inklusive →
  kein gebrochener MR-Link/Anker durch die Edits). Keine falsche Zahl, kein
  widersprüchlicher Satz eingeschleppt.
- **Unveränderte Kern-Funde halten (Plausibilität):** C-3/§3.3 (Historie-Provenance-
  Revocation, Doppelbeleg) und §3.4-Korrekturen (Fork-Klassifikation, Status-Herkunft)
  stehen unverändert und bleiben in sich schlüssig; D-5/C-6-Trennung ist mit §3.4
  jetzt konsistent statt widersprüchlich.

---

## Verdikt

**Abnahmereif.** Alle vier Befunde (3 MEDIUM / 1 LOW) korrekt und quell-belegt
eingearbeitet, keine Regression, Zählungen konsistent, Gates grün.
