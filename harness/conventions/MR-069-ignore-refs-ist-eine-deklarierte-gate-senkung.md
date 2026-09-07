# MR-069 — Das `ignore-refs`-Ventil ist eine deklarierte Gate-Senkung, und es wächst mit jedem Bump

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:**
  [`grundlagen-harness-dateien.md` §Ein einfrierendes Artefakt nennt ein prozess-bewegtes bei seiner Kennung](../../.harness/baseline/v6.5.0/regelwerk/grundlagen-harness-dateien.md)
  — keine Abweichung, sondern die **Einlösung** der dort verlangten
  Begründungslast: *„ein Ausnahme-Ventil im Prüfbereich, also eine
  Gate-Senkung mit eigener Begründungslast"*.
- **Datum:** 2026-09-07
- **Geltungsbereich:** die `ignore-refs`-Einträge in
  [`.d-check.yml`](../../.d-check.yml), die auf entfernte
  `.harness/baseline/<tag>/`-Bäume zeigen
- **Adaption:** Dieses Repo hält Links auf **entfernte** Baseline-Bäume über
  ein Ventil offen, statt sie nachzuziehen. Gemessen: **25** Einträge über
  **zehn** entfernte Tags, dahinter **28** Dateien — 18 `Accepted`-ADRs, sechs
  aufgelöste Konventions-Einträge, drei `done/`-Slices und ein gesendeter CR.

  **Der Grund je Klasse ist derselbe und trägt:** Alle vier Klassen zitieren
  den Stand ihrer Zeit. Eine `Accepted`-ADR darf im Kern nicht angefasst werden
  ([`AGENTS.md`](../../AGENTS.md#35-adrs-sind-nach-accepted-immutable) §3.5,
  maschinell gehalten), die übrigen wären als Lauf-Beleg verfälscht. Der alte
  Baum bleibt über die git-Historie prüfbar
  ([`MR-052`](../conventions.md#mr-052)).

  **Was der Kanon jetzt zusätzlich sagt, ist die eigentliche Aussage dieses
  Eintrags:** Das Ventil ist keine neutrale Konfiguration, sondern eine
  **Gate-Senkung** — im Prüfbereich hinter dem Ventil meldet `links` nichts
  mehr, auch nicht über einen Defekt, der nichts mit der Baseline zu tun hat.
  Deshalb ist jeder Eintrag **so eng wie möglich** zu skopieren: der jüngste
  ist auf **eine Datei** skopiert, nicht auf `docs/plan/adr/**` — die breite
  Form hätte jeden künftigen Link in den toten Baum still gedeckt.

- **Grenze — und sie ist der Grund, warum dieser Eintrag einen Auflösungs-Trigger
  hat und nicht „permanent" ist:** Das Ventil **wächst mit jedem Bump**, und
  jede Erweiterung senkt ein Stück mehr. Die Auflösung nach vorn steht im
  Kanon: Wer in einem einfrierenden Artefakt die **Kennung statt der Adresse**
  schreibt — Tag plus Pfad in Inline-Code statt Link —, erzeugt keinen neuen
  Eintrag mehr. Der Bestand wird **nicht** nachgezogen; er ist eingefroren.
- **Begründung:** Ohne diesen Eintrag steht die Gate-Senkung als
  Konfigurations-Detail da, verteilt über 25 Kommentare, und niemand sieht ihr
  Wachstum. Die Begründungslast, die der Kanon verlangt, ist eine Aussage über
  das **Ganze**, nicht über den einzelnen Eintrag.
- **Auflösungs-Trigger:** kein neuer Eintrag über drei aufeinanderfolgende
  Pin-Hebungen — dann trägt die Kennung-statt-Adresse-Form, und dieser Eintrag
  beschreibt nur noch einen ruhenden Bestand.
