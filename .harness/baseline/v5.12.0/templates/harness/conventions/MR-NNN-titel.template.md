# MR-\<NNN\> — \<Titel der Adaption\>

> **Template-Hinweis.** Vorlage für **einen** Adaptions-Eintrag. Kopiere nach
> `harness/conventions/MR-<NNN>-<titel>.md`, ersetze alle `<Platzhalter>` und
> lösche diesen Block. Ein Eintrag je Datei — die Index-Tabelle in
> [`../conventions.md`](../conventions.md) bekommt eine Zeile dazu.
> Ist der Auflösungs-Trigger eingetreten, wandert die Datei per `git mv` nach
> `done/`; der Zustand ist die Verzeichnis-Position, kein Status-Feld
> (Baseline-Regelwerk `grundlagen-harness-dateien.md`
> §harness/conventions.md als Konventionsspeicher).

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger;
`Löst auf` und `Ausgelöst durch Baseline-Stand` nur, wenn dieser Eintrag einen
früheren ablöst. `Ersetzt-Baseline-Regel` nennt **genau eine** Regel der
Baseline, an deren Stelle dieser Eintrag tritt — als Link mit
Abschnitts-Anker in die vendored Fassung; ein Datei-Link benennt keine Regel.
Wer mehrere Regeln ersetzen will, schreibt mehrere Einträge. Ein Eintrag, der
keine benannte Regel ersetzt, ist ein **Fork**, keine Adaption.

- **Status:** Accepted
- **Datum:** <Datum>
- **Geltungsbereich:** <Dateien / Module / Sub-Areas in DIESEM Repo — z. B. „`harness/README.md` §Source precedence und `AGENTS.md` §Kanonische Quellen">
- **Ersetzt-Baseline-Regel:** <genau eine Regel der Baseline, als Link mit Anker — z. B. [`grundlagen-referenz-richtung.md` §Spec-Straten](../../.harness/baseline/<tag>/regelwerk/grundlagen-referenz-richtung.md#spec-straten-mehr-als-ein-spec-dokument)>
- **Adaption:** <was stattdessen gilt — z. B. „Source-Precedence-Tabelle führt keinen eigenen Rang für `spec/spezifikation.md`: acht statt neun Ränge">
- **Begründung:** <warum, idealerweise mit Praxis-Bezug — z. B. „reines Policy-Repo, in dem keine eigenen technischen Festlegungen entstehen">
- **Auflösungs-Trigger:** <Trigger oder "permanent" — z. B. „sobald das Repo eigene technische Festlegungen trägt">
- **Löst auf:** [`MR-<NNN>`](../conventions.md#mr-<NNN>) *(nur, wenn dieser Eintrag einen früheren ablöst — sonst Zeile weglassen; der Verweis geht auf die Index-Zeile, nicht auf die Eintrags-Datei: die wandert bei Auflösung nach `done/` und ein Pfad-Link bricht genau dann)*
- **Ausgelöst durch Baseline-Stand:** <tag> *(Pflicht zusammen mit „Löst auf" — welcher Baseline-Stand die Ablösung ausgelöst hat)*
