# `make environment-check` — prüft, ob die Umgebung die deklarierte ist

> **Template-Hinweis.** Vorlage für die Sensor-Datei des Tool-/Versions-Checks.
> Kopiere nach `harness/sensors/environment-check.md` (heißt dein Target anders:
> nach dessen Namen), ersetze alle `<Platzhalter>` und lösche diesen Block. Die
> Vorlage liefert das **Muster, kein Skript**: Das Laufzeit-Gate schreibt das Repo
> selbst. Eine Index-Zeile in `harness/README.md` §Sensors gehört erst dazu, wenn
> das Skript existiert und sein Break-Test gelaufen ist — sonst behauptet sie ein
> Gate, das es nicht gibt. Diese Datei trägt Grenze und Bedienvertrag: Das
> Inventar ist eine eigene Datei, die das Skript liest; das Protokoll des
> Break-Tests steht im Skriptkopf oder in der ADR. Regeln und Herleitung:
> Baseline-Regelwerk `modul-14-docker-harness.md` §Der zweite Träger: Tool-Pin
> und Tool-/Versions-Check; Form der Sensor-Datei wie `gate.template.md`.

## Vertrag

Rot heißt: Die Umgebung trägt die Umgebungs-Zusage nicht — ein Werkzeug des
Inventars fehlt, weicht vom Soll ab oder hat die Fähigkeit nicht, die ein Skript
braucht. Ein **Umgebungs-Befund**, kein Befund über den Code. Der Check läuft als
**erster** in `make gates`.

Festlegung: `SPEC-<NNN>` in der Spezifikation — welche Prüfformen, wie
verglichen wird, was rot ist und was Warnung. Bindung der Anforderung:
`LH-QA-<…>` (die Umgebungs-Zusage des Projekts; dieser Check ist ihre
Messmethode).

## Inventar und Soll

**Inventar:** `<pfad>` — je Eintrag Werkzeug · Prüfform · Soll · Kontext · Grund.
Die Form der Datei wählt das Repo; das Skript liest sie nach der Regel aus
`SPEC-<NNN>`.

**Woher das Soll kommt:** `MR-<NNN>` in `harness/conventions.md` — <Weg 1,
ableiten: das Feld *Soll* nennt die Pin-Stelle, nicht ihren Wert · Weg 2,
Manifest: das Inventar trägt den Wert · Mischform. Für ein Werkzeug ohne
Pin-Stelle (`bash`, `make`) trägt das Inventar den Wert selbst>. Einen voreingestellten
Weg gibt es nicht. Hält ein Stellen-Abgleich Manifest und Pin-Stellen gegeneinander
(d-check `pinsync`), steht sein Block in `.d-check.yml`.

**Ausführungsort:** <nur Host · nur Image · Host und Image · ein Kontext ohne
Docker · zusätzlich CI-Runner> — deklariert in derselben `MR-<NNN>`.

## Grenze — was das Grün nicht abdeckt

1. **Drift zwischen zwei Läufen.** Der Check stellt die Umgebung zum Zeitpunkt
   des Laufs fest — permanent.
2. **Werkzeuge außerhalb des Inventars.** Was nicht im Inventar steht, prüft der
   Check nicht. <Belegt durch die Probe im leeren Kontext, Target
   `<make gates-probe>`, mit ihren Lücken: absolute Pfade (`#!/bin/bash`, die
   Rezept-Shell `/bin/sh`), Builtins, nicht ausgeführte Zweige · Ohne Probe:
   handgepflegt und unbelegt.>
3. **Kontexte, in denen er nicht läuft** — <Lücke des Ausführungsorts>.
4. **Zusätzlich Vorhandenes.** Der Check prüft, was das Soll nennt; globale
   Pakete und weitere Werkzeuge am Host bleiben unbemerkt. Er isoliert nicht.
5. <weitere>

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | jedes Werkzeug trägt sein Soll — kein Befund; fehlt ein `optional`-Eintrag, steht eine Warnzeile, die sagt, was ausfällt (eine Warnung ist kein Rot) |
| <n> | Umgebungs-Befund: je Werkzeug eine Zeile mit Werkzeug, Soll, Ist |

## Sperren

- **Bash zu alt** — der Guard steht als erster Top-Level-Befehl, vor jeder
  Funktionsdefinition und außerhalb jeder Funktion, nur mit Builtins
  (`BASH_VERSINFO`): Unterhalb der Mindestversion endet der Lauf mit einer
  Diagnose statt mit einem Syntaxfehler → eine neuere Bash in den `PATH` stellen.

## Bindung

`SPEC-<NNN>` · `LH-QA-<…>` · Inventar: `<pfad>` · Soll-Quelle und Ausführungsort:
`MR-<NNN>`.
