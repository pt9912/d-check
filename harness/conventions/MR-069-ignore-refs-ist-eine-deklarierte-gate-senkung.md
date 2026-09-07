# MR-069 — Das `ignore-refs`-Ventil ist eine deklarierte Gate-Senkung, und es wächst mit jedem Bump

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:**
  [`grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt](../../.harness/baseline/v6.5.0/regelwerk/grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt) (Absatz *„Ein einfrierendes Artefakt nennt ein prozess-bewegtes bei seiner Kennung"* — eine Fett-Zeile, keine Überschrift, deshalb der Abschnitts-Anker)
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

  **Was der Kanon jetzt zusätzlich sagt, ist die Aussage dieses Eintrags:** Das
  Ventil ist keine neutrale Konfiguration, sondern eine **Gate-Senkung**. Aber
  eine **enge**, und die erste Fassung dieses Eintrags hat sie zu breit
  bemessen: Sie behauptete, hinter dem Ventil melde `links` nichts mehr, auch
  nicht über einen Defekt ohne Baseline-Bezug. **Gemessen ist das falsch** —
  ein `in:`/`refs:`-Paar nimmt genau die Referenzen auf das **Ziel-Glob** aus,
  nicht die Datei; ein toter Nicht-Baseline-Link in derselben ADR meldet
  weiterhin `target-missing` (Bruch-Test, unabhängiger Review).

  **Die Senkung ist damit präzise benennbar:** In der genannten Datei wird
  **jeder** Verweis auf den entfernten Baum stumm — auch ein **neuer**, den
  jemand künftig einträgt. Das ist der ganze Preis, und er ist klein. Er wird
  aber **25-fach** bezahlt und wächst mit jedem Bump.

- **Grenze — drei, und sie sind der Grund für den Auflösungs-Trigger:**

  **(1) Das Ventil wächst mit jedem Bump.** Jede Hebung, die eine Adresse in
  einem einfrierenden Artefakt zurücklässt, braucht einen weiteren Eintrag.

  **(2) Ein Teil davon skopiert nichts mehr.** Gemessen: **16** der Einträge
  in `ignore-refs` nennen eine `in:`-Datei, die es **nicht mehr gibt** — sie
  wurde archiviert, nachdem der Eintrag geschrieben war. Sie sind wirkungslos
  und niemand meldet sie; die 25 Baseline-Skopen oben sind die **deklarierten**,
  nicht durchgehend die **wirksamen**. Das Aufräumen ist ein eigener Vorgang und
  gehört nicht in diesen Eintrag.

  **(3) Die Auflösung nach vorn trägt nur eine der vier Klassen.** Der Kanon
  löst das Problem, indem einfrierende Artefakte die **Kennung statt der
  Adresse** schreiben — Tag plus Pfad in Inline-Code statt Link. Für
  Review-Reports steht diese Regel jetzt im Reviewer-Skill. Für die **18
  `Accepted`-ADRs**, die die Mehrheit des Bestands stellen, gibt es **keinen
  Träger**: Eine ADR entsteht ohne Skill, aus einer Vorlage, und die adoptierte
  Vorlage führt die Zitier-Form nicht. Solange das so ist, erzeugt die gelebte
  ADR-Praxis weiter Ventil-Einträge, und der Auflösungs-Trigger unten kann
  **nicht feuern**. Das ist eine benannte Lücke, kein Versehen — und der
  nächstliegende Träger wäre eine Ergänzung der ADR-Vorlagen-Nutzung, nicht
  dieser Eintrag.

  Der **Bestand** wird nicht nachgezogen; er ist eingefroren.
- **Begründung:** Ohne diesen Eintrag steht die Gate-Senkung als
  Konfigurations-Detail da, verteilt über 25 Kommentare, und niemand sieht ihr
  Wachstum. Die Begründungslast, die der Kanon verlangt, ist eine Aussage über
  das **Ganze**, nicht über den einzelnen Eintrag.
- **Auflösungs-Trigger:** **beide** Bedingungen — ein Träger für die ADR-Klasse existiert (Grenze 3), **und** kein neuer Eintrag über drei aufeinanderfolgende
  Pin-Hebungen — dann trägt die Kennung-statt-Adresse-Form, und dieser Eintrag
  beschreibt nur noch einen ruhenden Bestand.
