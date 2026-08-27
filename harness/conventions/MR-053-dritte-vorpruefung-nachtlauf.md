# MR-053 — Die Slice-Planung trägt eine dritte Vorprüfung: den Nachtlauf-Stand

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`docs/plan/planning/slice.template.md`](../../.harness/baseline/v5.12.0/templates/docs/plan/planning/slice.template.md)
  §Vorgelagert — die Vorlage kennt **zwei** Blöcke (Sub-Area-Wahl prüfen,
  offene Beobachtungen sichten). Dieser Eintrag ergänzt einen dritten. Der
  Kanon kennt keinen Nachtlauf; er kann ihn nicht vorsehen.
- **Datum:** 2026-08-27
- **Geltungsbereich:** jeder neue Slice-Plan dieses Repos, §Vorgelagert. Kein
  Retrofit — die geschlossenen Slices bleiben, wie sie sind.
- **Adaption:** Dieses Repo betreibt einen Nachtlauf über **zwölf** gepinnte
  Fremd-Bestände ([`upstream-drift.yml`](../../.github/workflows/upstream-drift.yml)).
  Er meldet korrekt — und **an niemanden**: der Job fällt rot aus und ist nur
  in der Actions-Übersicht sichtbar. Der Workflow-Kopf trägt die Lücke seit
  seiner Erweiterung als benannte Grenze; sie ist damit alt, aber nicht
  geschlossen.

  **Die Regel benennt zuerst den Adressaten, dann den Takt.** Adressat ist der
  **Rolleninhaber der Implementer-Rolle beim Planen des nächsten Slice** —
  nicht „das Team" und nicht „wer hinsieht". Takt ist **jede Slice-Planung**,
  als dritter `**Vorgelagert**`-Block neben den zwei des Kanons. Der Moment ist
  gewählt, weil dort ohnehin ein Lese-Schritt stattfindet; ein neuer Moment
  hätte dieselbe Verwaisung eine Ebene höher.

  **Der Stand ist ohne Fremd-Werkzeug lesbar:** `make nightly-state` fragt die
  GitHub-API mit `curl` — derselben Erwartung, die die Netz-Targets ohnehin
  tragen —, ist fail-open und endet **immer** mit Exit 0. Der Ausgang steht in
  der **Ausgabe**, nicht im Code: wer sie nicht liest, hat den Schritt nicht
  getan, und das soll kein Exit-Code verdecken können.

  **Rausch-Unterscheidung, mitentschieden:** eine **planmäßige** Meldung wird
  anders behandelt als eine unerwartete. Planmäßig sind ein Fremd-Release
  (`VERALTET`/`ABWEICHEND` einer Frische-Achse) und eine verschobene
  Zitat-Spanne nach einem Baseline-Bump
  ([`MR-051`](../conventions.md#mr-051)). Der Unterschied steht in der
  **Ausgabe** des Laufs, nicht in seiner Farbe — deshalb verlangt die
  Vorprüfung das Lesen, nicht das Zählen.

  **Kein Benachrichtigungs-Kanal, und der Grund ist gemessen, nicht gescheut:**
  jeder ohne Fremd-Dienst verfügbare Kanal hängt entweder an den
  **Watch-Einstellungen einzelner Nutzer** — das ist keine Repo-Zusage, sondern
  eine Konto-Einstellung — oder erzeugt eine **neue Artefakt-Klasse** (ein
  Issue je rotem Lauf) mit eigenem Rausch-Problem und eigener Pflege. Beides
  ist teurer als der Lese-Schritt und löst das Kernproblem nicht: dass jemand
  hinsehen muss.

  **Grenze, benannt statt wegerklärt:** greift die Vorprüfung nur bei einer
  Slice-Planung, dann liest in einer **Pause** niemand. Das ist der Rest, den
  diese Regel nicht schließt.
- **Begründung:** Ein Wächter ohne Leser ist derselbe verwaiste Sensor, gegen
  den er gebaut wurde — nur eine Ebene höher. Die billigste Reparatur ist, ihn
  an einen Moment zu hängen, den es schon gibt, statt einen neuen zu erfinden.
- **Auflösungs-Trigger:** ein Rot bleibt über mehrere Slice-Planungen hinweg
  ungelesen. Dann trägt der Lese-Schritt nicht, und die Kanal-Frage ist neu zu
  stellen — mit dem dann gemessenen Grund.
