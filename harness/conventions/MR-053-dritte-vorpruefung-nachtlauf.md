# MR-053 — Die Slice-Planung trägt eine dritte Vorprüfung: den Nachtlauf-Stand

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon kennt **zwei** `Vorgelagert`-Blöcke
  (Sub-Area-Wahl prüfen, offene Beobachtungen sichten) und **keinen**
  Nachtlauf — er kann ihn nicht vorsehen. Dieser Eintrag ergänzt einen dritten
  Block; die Form-Frage tritt die Rangliste an diesen Speicher ab.
- **Datum:** 2026-08-27
- **Geltungsbereich:** jeder Slice-Plan dieses Repos, §Vorgelagert. Der Block
  entsteht **spätestens bei der Beanspruchung** (`open`→`in-progress` bzw.
  `next`→`in-progress`) — wie das `**Verantwortlich:**`-Feld; ein Plan in
  `open/` trägt ihn noch nicht. Kein
  Retrofit: die geschlossenen Slices bleiben, wie sie sind.
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

  **Kein Benachrichtigungs-Kanal — und das ist eine Erwägung, keine Messung.**
  Gefahren ist keiner; behauptet wird deshalb nichts über ihre Zustellung.
  **Drei** Kandidaten sind betrachtet, und die Liste ist nicht vollständig:

  | Kandidat | Warum nicht |
  |---|---|
  | Actions-Benachrichtigung | hängt an der **Watch-Einstellung** eines Kontos, nicht am Repo. In einem Ein-Personen-Repo ist dieser Konto-Inhaber derselbe wie der Adressat oben — das Argument ist formal richtig und praktisch schwach |
  | Issue je rotem Lauf | neue Artefakt-Klasse mit eigenem Rausch- und Pflege-Problem. **Eine mildere Form — ein einziges, fortgeschriebenes Issue — ist nicht betrachtet worden** |
  | Status-Badge im `README` | keine Konto-Einstellung, keine neue Artefakt-Klasse, kein Geheimnis — aber sichtbar nur auf der **gerenderten** GitHub-Seite, nicht im lokalen Lauf, in dem die Slice-Planung stattfindet |

  **Ausdrücklich keine Kandidaten:** `CODEOWNERS` und Repository-Rulesets
  steuern PR-Review-Zuweisung, keine Actions-Ausgänge.

  **Keiner der drei löst das Kernproblem** — dass jemand hinsehen muss — besser
  als ein Lese-Schritt an einem Moment, den es schon gibt. Das ist der Grund;
  „gemessen" wäre er nicht.

  **Grenze, benannt statt wegerklärt:** greift die Vorprüfung nur bei einer
  Slice-Planung, dann liest in einer **Pause** niemand. Das ist der Rest, den
  diese Regel nicht schließt.
- **Begründung:** Ein Wächter ohne Leser ist derselbe verwaiste Sensor, gegen
  den er gebaut wurde — nur eine Ebene höher. Die billigste Reparatur ist, ihn
  an einen Moment zu hängen, den es schon gibt, statt einen neuen zu erfinden.
- **Auflösungs-Trigger:** ein Rot bleibt über mehrere Slice-Planungen hinweg
  ungelesen. Dann trägt der Lese-Schritt nicht, und die Kanal-Frage ist neu zu
  stellen — mit dem dann gemessenen Grund.
