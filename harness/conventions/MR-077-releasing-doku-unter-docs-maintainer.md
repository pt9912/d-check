# MR-077 — Die Releasing-Doku liegt unter docs/maintainer/, Rang 6 der Source Precedence umfasst beide Verzeichnisse

- **Datum:** 2026-10-09
- **Geltungsbereich:** [`docs/maintainer/`](../../docs/maintainer/releasing.md),
  Rang 6 in [`AGENTS.md`](../../AGENTS.md) §2 und
  [`harness/README.md`](../README.md#source-precedence) §Source precedence
- **Ersetzt-Baseline-Regel:** Baseline-Regelwerk
  [`grundlagen-source-precedence.md` §Source Precedence](../../.harness/baseline/v6.17.0/regelwerk/grundlagen-source-precedence.md#source-precedence),
  Rang 6 — `docs/user/*.md` als Ort der Betriebs-, Quality-, Releasing- und
  Runbook-Sichten
- **Adaption:** Die Releasing-Doku
  ([`docs/maintainer/releasing.md`](../../docs/maintainer/releasing.md))
  liegt getrennt von der Doku für Anwender. Rang 6 nennt beide Verzeichnisse:
  `docs/user/` (Operations) und `docs/maintainer/` (Releasing). Der Rang selbst
  bleibt derselbe — die Datei steht über `README.md` und unter der Roadmap, wie
  vorher.
- **Grenze:** Kein Gate liest die Rangtabelle gegen die Ablage; eine weitere
  Datei unter `docs/maintainer/` ist durch diesen Eintrag gerankt, eine unter
  einem dritten Verzeichnis nicht.
- **Begründung:** Auftraggeber-Wunsch: wer das Werkzeug benutzt, sucht unter
  `docs/user/`, wer es veröffentlicht, unter `docs/maintainer/`.
- **Auflösungs-Trigger:** `docs/maintainer/` ist leer — die Releasing-Doku
  ist nach `docs/user/` zurückgekehrt —, oder die Baseline legt die
  Releasing-Sicht selbst unter `docs/maintainer/`. Nennt die Baseline einen
  anderen Ort, ist das ein Anlass zur Neubewertung, keine Auflösung.
