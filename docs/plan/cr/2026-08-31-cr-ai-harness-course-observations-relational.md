# Eingehender CR aus `ai-harness-course` — teamfähige BEO-Ablage relational prüfen

**Absender:** ai-harness-course (Adopter von d-check; zugleich Quelle der
adoptierten Baseline).
**Eingegangen:** 2026-08-31, über den Auftraggeber.
**Gegenstand:** eine opt-in Fähigkeit für ein **verzeichnisbasiertes**
Beobachtungs-Register — Identität aus Pfad und Inhalt, Evidence-Paarung,
Invalidierungen als Differenzmenge, Alias-Auflösung, abgeleiteter Zähler mit
3×-Gate, Diagnosesicht.
**Nachtrag:** [Antwort des Absenders auf unsere vier Rückfragen](2026-08-31-antwort-ai-harness-course-observations-relational.md)
— sie zieht Teile des Antrags zurück und ist zusammen mit ihm zu lesen.

Dieses Dokument hält den CR **wie empfangen** fest. Die Bewertung steht nicht
hier — ein CR-Dokument trägt Bitte und Beleg, nicht die Antwort darauf.

---

## Anlass (des Absenders)

Das Beobachtungs-Register liegt bisher in einer gemeinsamen Datei mit
fortlaufender Kennung, gespeichertem Zähler und mehreren Slice-Belegen je Zeile.
Bei mehreren Branches entstehen daraus drei Nebenläufigkeitsprobleme:

1. Offene Branches kennen nur den zuletzt gemergten Zählerstand.
2. Zwei Branches können dasselbe Phänomen unter verschiedenen freien Nummern
   anlegen und konfliktfrei mergen.
3. Zwei isoliert grüne Branches können gemeinsam den Übergang von 1× auf 3×
   erzeugen, ohne dass einer der beiden Branch-Läufe die fällige Folgeaktion
   verlangt.

Der Entwurf ersetzt Nummer und gespeicherten Zähler durch getrennte Artefakte
unter `observations/BEO-<SUB-AREA>/<slug>/` mit `observation.md`, `state.md`,
`evidence/<slice-id>.md` und `invalidations/<slice-id>.md`. Die kanonische
Kennung ist der Pfad; der Zähler wird aus den gültigen Evidence-Dateien
**abgeleitet** statt gespeichert.

**Ausdrücklich nicht Teil des Antrags:** die Immutabilität von Namespace, Slug,
`observation.md` und bestehenden Evidence-Dateien — die deckt das vorhandene
Modul `vcs` über `--range`/`--staged` bereits ab
([`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)).

## Beantragter Vertrag (sechs Teile)

1. **Identität aus Pfad und Inhalt** — Namespace `BEO-<SUB-AREA>`, Slug 3–80
   Zeichen lowercase-Kebab-Case; die in `observation.md` deklarierte Kennung
   entspricht der aus dem Pfad abgeleiteten; die Herkunfts-Sub-Area ist im
   konfigurierten Sub-Area-Bestand deklariert; jede kanonische Kennung kommt
   genau einmal vor.
2. **Evidence und Slice-Closure sind beidseitig gepaart** — Dateiname trägt
   genau eine Slice-Kennung, der Slice existiert im `done/`-Bestand, die
   Evidence verweist auf dessen Closure, und die Closure referenziert dieselbe
   Kennung zurück. Einseitige Referenzen zählen nicht und sind Befunde; je Slice
   und Kennung höchstens ein Beleg.
3. **Invalidierungen bilden eine Differenzmenge** — ein falscher Beleg wird
   nicht gelöscht, sondern zurückgenommen; verwaiste oder doppelte
   Invalidierungen sind Befunde; gültige Evidence = Evidence − Invalidierungen.
4. **Aliase werden zu einem kanonischen Ziel aufgelöst** — Ziel existiert, Graph
   azyklisch, jede Kette endet an genau einer nicht-Alias-Kennung, ein Alias
   nimmt keine Evidence an, Belege werden am Ziel nach Slice-Kennung vereinigt.
5. **Der abgeleitete Zähler erzwingt den Steering-Loop-Ausgang** — unter 3 ist
   `open` zulässig; ab 3 verlangt der Vertrag einen Zustand mit auflösbarem
   Aktions- bzw. Verkörperungsziel. Eine Kennung mit ≥ 3 gültigen Belegen ohne
   gültigen Ausgang lässt den Lauf fehlschlagen.
6. **Ableitung im Diagnoseformat sichtbar machen** — je kanonischer Kennung
   Zustand, Anzahl gültiger und invalidierter Belege, Alias-Mitglieder und Ziel.

Der Lauf implementiert **keine** Merge-Queue; der Absender setzt voraus, dass
die CI den aktuellen bzw. synthetischen Merge-Stand vorlegt.

## Abnahme durch Mutationsproben

Der Antrag trägt sechzehn Proben mit Erwartung, darunter: zwei gültige Belege ⇒
grün mit Zähler 2 · dritter Beleg ohne Aktion ⇒ rot · dritter Beleg mit totem
Aktionsziel ⇒ rot · Evidence ohne Rückreferenz der Closure ⇒ rot ·
Invalidierung ohne Evidence ⇒ rot · gültig invalidierter Beleg ⇒ grün, Zähler
sinkt · Alias auf fehlende Kennung ⇒ rot · Alias-Zyklus ⇒ rot · Evidence unter
einem Alias ⇒ rot · Kennung weicht vom Pfad ab ⇒ rot · **zwei semantisch
gleiche Beobachtungen mit verschiedenen Slugs ⇒ grün, bewusst menschliche
Grenze** · fehlender Root bei aktiver Konfiguration ⇒ rot · vorhandener, aber
leerer Root ⇒ grün mit ausgewiesener Null.

## Abgrenzung (des Absenders)

Nicht beantragt: semantische Erkennung gleichbedeutender Slugs · Beurteilung,
ob ein Vorfall fachlich wahr ist · Beurteilung der Qualität einer Folgeaktion ·
Einrichtung oder Steuerung einer Merge-Queue · Änderung des vorhandenen
`vcs`-Vertrags · automatische Migration bestehender Kennungen · ein
allgemeiner, frei programmierbarer Regel-Interpreter.
