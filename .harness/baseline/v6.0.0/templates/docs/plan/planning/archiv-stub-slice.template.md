# slice-<NNN> — <Titel>

> **Template-Hinweis.** Vorlage für den gekürzten Stub, der beim Archivieren
> einer Welle an der Stelle des Slice-Volltexts liegen bleibt
> (`docs/plan/planning/done/<welle-id>/`). Kopiere, ersetze Platzhalter,
> lösche diesen Block. Für den Welle-Plan gilt `archiv-stub-welle.template.md`.

Regeln dieses Artefakts: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur (Modul 6), Schritt 4 — was archiviert wird, was
liegen bleibt, in welcher Form, und dass die Ergebnisnotiz vollständig **und
flach** bleibt.

**Der Stub trägt keine Abschnittsüberschriften.** Ein Stub hat keine, ein
ungekürzter Plan hat sie — daran ist die Kürzung form-prüfbar. Wer `## …` im
Stub stehen lässt, hat nicht gekürzt, sondern nur markiert.

**`Welle:` und `Archiviert mit:` sind zwei Tatsachen, kein Widerspruch.**
Das erste Feld nennt die *Zugehörigkeit*, das zweite die *Einsammlung*. Ein
Slice ohne Wellen-Zugehörigkeit behält `ohne Welle` und nennt im zweiten Feld
die Welle, deren Closure ihn archiviert hat.

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten.

Der Stub traegt vier Dinge und sonst nichts: Identitaet, Archiv-Zeiger,
Zustand, und die Kennungen, die den Vorgang ueberlebt haben. Lerneintrag,
Risiko-Ausgaenge, DoD-Tabelle und Abnahme gehen ins Archiv — sie stehen
ohnehin dort, wo sie gelesen werden: im Beobachtungs-Register, als ADR, als
Folge-Slice. Genau die nennt `Hervorgegangen:`.

Review-Reports bekommen KEINEN Stub — sie liegen im Archiv unter ihrem Slice.
-->

> **ARCHIVIERT** — Volltext:
> `unzip -p done/<welle-id>/archiv.zip <pfad-im-archiv>`

**Welle:** <welle-id | ohne Welle>
**Archiviert mit:** <welle-id> · **Geschlossen:** <JJJJ-MM-TT>
**Hervorgegangen:** <BEO-*, ADR-*, Folge-Slice — oder `— keine —`>
