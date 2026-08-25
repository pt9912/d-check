# MR-035 — Ausgehende Change Requests an die Baseline liegen im Repo

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon kennt den **eingehenden** Change
  Request (Vertragsänderung am Lastenheft) und den Konsumenten-CR als Vorgang,
  aber keine Ablage für den **ausgehenden** — die Verzeichnisfrage ist damit
  eine Form-Frage, und die tritt die Rangliste an diesen Speicher ab
  ([`grundlagen-harness-dateien.md` §Konventionsspeicher](../../.harness/baseline/v5.11.0/regelwerk/grundlagen-harness-dateien.md)).
- **Datum:** 2026-08-23
- **Geltungsbereich:** `docs/plan/cr/` — Change Requests, die dieses Repo an die
  **adoptierte Baseline** richtet.
- **Adaption:** Ausgehende CRs liegen als datierte Datei unter `docs/plan/cr/`,
  Form `<YYYY-MM-DD>-cr-<gegenstand>.md`. Sie sind **kein** Spec-Stratum, keine
  ADR und kein Slice: sie entscheiden nichts in diesem Repo, sie bitten
  woanders um eine Entscheidung.

  **Warum überhaupt eine Ablage:** Der vorige Konsumenten-CR dieses Repos wurde
  geschrieben, gesendet und **nicht aufbewahrt**. Seine Annahme ist heute nur
  noch über das CHANGELOG der Gegenseite belegbar — und was genau gebeten wurde,
  welche Punkte abgelehnt wurden und mit welcher Begründung, steht nirgends.
  Genau diese Fragen sind später die wertvollsten: sie sagen, welche Form beim
  nächsten Mal trägt.

  **Was die Ablage nicht ist:** kein Register und keine Liste. Es gibt keinen
  Index, keinen Zähler und keine Pflicht, den Ausgang nachzutragen — wer das
  will, hängt es an den Slice, der den CR schreibt.
- **Ausgelöst durch Baseline-Stand:** v5.11.0
- **Auflösungs-Trigger:** der Kanon benennt selbst einen Ruheort für ausgehende
  Change Requests. Dann gilt seiner.
