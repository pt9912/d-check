# MR-021 — In-Repo-Verweise auf das vendored Regelwerk sind pin-gebunden

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`grundlagen-harness-dateien.md` §Verzeichniskonvention](../../.harness/baseline/v5.11.0/regelwerk/grundlagen-harness-dateien.md#verzeichniskonvention)
- **Datum:** 2026-06-26
- **Geltungsbereich:** **alle** Markdown-Links auf `.harness/baseline/<tag>/…`
  in der Live-Doku — das Briefing, [`harness/README.md`](../README.md), dieser
  Konventionsspeicher samt seiner **aktiven** Eintrags-Dateien, der
  Reviewer-Skill und die lebenden Planungs-Dokumente. Eine Aufzählung einzelner
  Dateien stünde hier falsch: sie wächst mit jedem Eintrag und ist beim
  nächsten Bump veraltet — die Menge bestimmt der Zensus der Bump-Prozedur,
  nicht diese Zeile. Ebenso erfasst: die
  Baseline-Pin-Bump-Prozedur; Nachtrag zu
  [`MR-019`](../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)/[`MR-020`](../conventions.md#mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt)
- **Adaption:** Seit
  [`MR-019`](../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)
  ist das vendored Regelwerk ein in-repo auflösbares **Link-Ziel** — die
  Live-Doku verweist auf konkrete Regelwerk-Dateien (Lesestoff: Modul-/
  Grundlagen-Verweise) statt nur auf externe Kurs-URLs (die als **Provenienz**
  bleiben). Diese Links tragen den **konkreten** Pin (aktuell `…/v5.11.0/…`), nicht
  `<tag>` — sie sind damit **pin-gebunden**. Regel: Der Baseline-Pin-Bump-
  Drift-Audit ([`MR-020`](../conventions.md#mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt))
  (1) entfernt das alte `.harness/baseline/<alt-tag>/` und (2) zieht alle
  vendored-Pfad-Links auf den neuen Tag. Wird (2) vergessen, schlägt nach (1)
  `make doc-check` mit `target-missing` an — die Pin-Kopplung ist damit
  **gate-erzwungen**, kein stiller Drift. (Bliebe das alte Tag-Verzeichnis
  stehen, läge stiller Stale-Content vor; darum ist (1) Pflicht-Teil des Bumps.)
- **Begründung:** Nutzer-Entscheid 2026-06-26, das vendored Regelwerk als
  Lesestoff zu verlinken (§Guides-Lese-Form + Modul-13/14-Verweise +
  §Adoptierte-Aktuell-Link). Der Nutzen (klickbar, netzlos, offline auffindbar)
  hat als Preis die Pin-Bindung; die Regel macht den Preis explizit und delegiert
  die Durchsetzung an das vorhandene `links`-Gate, statt einen neuen Sensor zu
  bauen (Steering-Loop-Ökonomie: kein Gate für etwas, das ein vorhandenes Gate
  schon fängt).
- **Auflösungs-Trigger:** permanent, solange in-repo-Verweise auf das vendored
  Regelwerk bestehen.
