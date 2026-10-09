# `make nightly-state` — liest den Ausgang der beiden Nachtläufe und sagt, ob er gelesen werden muss

## Vertrag

**Lese-Schritt, kein Gate.** Der Nachtlauf meldet korrekt und **an niemanden**:
Der Job fällt rot aus und ist nur in der Actions-Übersicht sichtbar. Dieses
Target liest seinen Ausgang über die GitHub-API und hängt damit an einem
Moment, den es schon gibt — der **Slice-Planung**, als dritte Vorprüfung
([`MR-053`](../conventions/MR-053-dritte-vorpruefung-nachtlauf.md)).

**Adressat zuerst, Takt danach:** Ein neuer Moment hätte dieselbe Verwaisung
eine Ebene höher erzeugt. Werkzeug ist `curl`, nicht `gh` — die Netz-Targets
tragen diese Erwartung ohnehin, `gh` wäre eine neue.

**Rausch-Unterscheidung — in der Ausgabe, nicht im Werkzeug:** Ob eine rote
Meldung planmäßig ist (Fremd-Release; Zitat-Spanne nach einem Bump,
[`MR-051`](../conventions/MR-051-cite-spannen-beim-bump.md)), entscheidet das
Werkzeug **nicht**. Es gibt bei `failure` einen festen Hinweis aus, die Ausgabe
des Nachtlaufs zu lesen; dort steht der Unterschied, nicht in der Farbe.

**Kein Benachrichtigungs-Kanal, mit Grund:** Jeder ohne Fremd-Dienst
verfügbare hängt an Watch-Einstellungen einzelner Nutzer — das ist keine
Repo-Zusage — oder erzeugt eine neue Artefaktklasse.

## Grenze — was das Grün nicht abdeckt

1. **In einer Pause liest niemand** — der Schritt greift nur bei einer
   Slice-Planung. Permanent, das ist der Preis des vorhandenen Moments.
2. **Gelesen wird der jüngste Lauf, nicht sein Alter** — ein abgeschalteter
   Nachtlauf meldete weiter `gruen`, und ein Lauf **vor** einer Behebung
   meldet weiter rot. Heilbar durch eine Altersprüfung, bisher nicht gebaut.
3. **Bei privatem Repo oder umbenanntem Workflow** ist der `SKIP` von einer
   Netzstörung nicht unterscheidbar. Permanent.
4. **Der Repo-Slug ist ein Default**, kein Fund (`NIGHTLY_REPO`,
   `NIGHTLY_WORKFLOWS` überschreiben ihn) — in einem Fork meldete es sonst den
   Nachtlauf des Originals.
5. **Der feste Hinweis bei `failure` passt nicht auf jeden Nachtlauf** — er
   nennt die planmäßigen Meldungen der Drift-Achsen und erscheint auch beim
   CVE-Scan. Ein `NIGHTLY_WORKFLOWS` aus lauter Leerraum liest nichts und
   meldet nichts.

## Ausgabe und Ausgänge

**Immer Exit 0**, fail-open; die Urteile je Ausgang legt
[`SPEC-101`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
fest. Der Ausgang steht in der **Ausgabe**: Wer sie nicht liest, hat den
Schritt nicht getan, und das soll kein Exit-Code verdecken. Netzlos prüfbar
über `--parse <datei>` und `--selftest`.

## Bindung

kein Gate — liest, urteilt nicht.
[`MR-053`](../conventions/MR-053-dritte-vorpruefung-nachtlauf.md) ·
[`SPEC-101`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
