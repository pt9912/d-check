# ADR-0098: `hostpaths` erkennt Home-relative Pfade mit Punkt-Ausnahme und trägt ein Ziel-Ventil statt eines Zeilen-Markers

**Status:** Proposed

**Datum:** 2026-10-05

**Autor:** pt9912

**Bezug:** [`DC-FA-HOST-001`](../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)
(erweitert, Lastenheft 0.94.0); Schnitt-Kriterium (Einzelmodul-Frage ⇒
bestehende Anforderung ändern) aus
[ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md); Form der
Konfigurations-Weitung wie
[ADR-0058](0058-konfigurations-flaechen-additiv-weiten.md); Anlass ist
eine Nutzer-Frage (2026-10-05), slice-249 <!-- d-check:status-provenance -->.

**Schärft:**
[`spec/spezifikation.md` §DC-FA-HOST-001.a](../../../spec/spezifikation.md#dc-fa-host-001a--host-pfad-erkennung)
und die Schema-Zeilen `hostpaths.*` unter
[`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Modul `hostpaths` meldet absolute Pfade unter einem Host-Präfix. Eine
Home-relative Angabe — Tilde, Schrägstrich, Verzeichnis — kannte es nicht
als eigene Form. Gemessen traf das Unix-Muster sie trotzdem, aber nur
**zufällig**: Die Tilde stand nicht in der Wortgrenzen-Vorbedingung, also
meldete es bei einem Präfix-Namen als erstem Segment den Rest ab dem
zweiten Schrägstrich, abgeschnitten und ohne Tilde. Bei jedem anderen
Verzeichnisnamen (persönliche Layouts wie Quell- oder Projektordner) blieb
es still. Der Probelauf über den Plan dieses Vorgangs zeigte den
abgeschnittenen Treffer an einer Stelle, an der das Modul ihn nicht hätte
erzeugen sollen.

Zugleich ist die Tilde in Doku häufig legitim: Werkzeug-Konventionen
unter `~/.config`, `~/.cache`, `~/.claude` lauten auf jedem Host gleich und
leaken kein Layout. Der gemessene Bestand (eigenes Repo ohne vendorte
Baseline, Kurs, a-check) trägt genau zwei Tilde-Stellen, beide dieser Art.

Das Lastenheft sagte bisher „keinen Opt-out-Marker"; Auswege waren Fences
und die Präfixliste. Für eine Home-relative Erkennung trägt die Präfixliste
nicht: sie ist für die Wurzel gebaut, und ein Präfix-Name mit Tilde würde
zu Schrägstrich-Tilde-Schrägstrich und träfe nie, was gemeint ist — still,
weil der Config-Rand nur den Schrägstrich ablehnte.

## Entscheidung

1. **Home-relativ ist ein eigenes, festes Muster: Tilde, Schrägstrich,
   erstes Segment ohne führenden Punkt.** Das Punkt-Kriterium trennt
   Werkzeug-Konvention (host-unabhängig) von persönlichem Layout
   (host-gebunden) ohne Liste. Die Tilde mit Benutzername, die nackte
   Tilde mit Schrägstrich und Variablen-Formen bleiben außen vor —
   benannte Grenzen im Lastenheft.
2. **Ein Unix-Treffer innerhalb eines Home-relativen Treffers entfällt —
   kein anderer.** Ein Home-relativer Pfad wird genau einmal und in voller
   Form gemeldet; für bestehende Konsumenten ändert sich bei diesen Fällen
   der gemeldete Pfad (mit Tilde statt ohne), die Fundstelle bleibt. Eine
   pauschale Tilde-Sperre im Unix-Muster wäre einfacher, nimmt aber Funde
   weg, die keine Home-Verweise sind (durchgestrichener Pfad, am Wort
   klebender Pfad) — Messung im Review des Vorgangs.
3. **Ausnahmen trägt ein Ziel-Ventil `hostpaths.exempt-targets`, kein
   Zeilen-Marker.** Globs über den gemeldeten Pfad, `matchGlob` wie
   `scan.ignore`, segmentweise am Config-Rand validiert; ein Glob, das
   weder mit `/` noch mit `~` beginnt, ist Exit 2 (es träfe nie). Es gilt für Unix-
   und Home-relative Funde, nicht für Windows/UNC — dort ist der
   Backslash im Glob das Escape-Zeichen. Der Satz „keinen
   Opt-out-Marker" wird auf den **Zeilen-Marker** präzisiert: eine
   repo-weite Ausnahme steht an einer Stelle und ist im Review sichtbar,
   ein Marker an jeder Fundstelle nicht.
4. **Ein Präfix-Name mit Tilde ist ein Konfigurationsfehler (Exit 2)**,
   mit Hinweis auf die eigene Erkennung und das Ventil.

Die Erkennung **schärft** das Modul. Die Byte-Identitäts-Zusage
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)) gilt
für den neuen Schlüssel, **nicht** für das neue Muster — das ist der Zweck.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Nur Tilde vor einem Präfix-Namen** | kaum Fehlalarme | die Präfixliste ist für die Wurzel gebaut; unter dem Home-Verzeichnis passt praktisch nur ein Name, und der ist ein persönliches Layout — die Variante fängt fast nichts und sieht aus, als fange sie viel |
| **Jede Tilde, Ausnahmen nur per Ventil** | eine Regel, keine Sonderform | jedes Repo müsste dieselben Werkzeug-Konventionen ausnehmen — Pflicht-Konfiguration für den Normalfall; der gemessene Bestand wäre an beiden Stellen rot, eine davon in einer `Accepted`-ADR |
| **Zeilen-Marker `d-check:ignore` für `hostpaths`** | Ausnahme direkt an der Fundstelle | verstreut, nicht zentral reviewbar, verdeckt jeden Pfad der Zeile; das Lastenheft hatte ihn ausdrücklich ausgeschlossen |
| **Punkt-Regel + Ziel-Ventil** (gewählt) | Werkzeug-Konventionen still ohne Konfiguration; Rest zentral ausnehmbar | eine host-unabhängige Konvention ohne führenden Punkt wird gemeldet (benannte Grenze, Ausweg Ventil) |

## Konsequenzen

- Konsumenten mit Home-relativen Layout-Pfaden in Prosa werden nach dem
  Update rot; CHANGELOG und Benutzerhandbuch müssen das sagen.
- Die Spezifikation trägt ein viertes Muster und einen Ventil-Schritt;
  das Schema zwei Zeilen.
- Der `--print-config`-Gerüst-Block `hostpaths` nennt das Ventil.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test `TestHostpathsTilde` | Home-relativ einmal, voll; Punkt-Segment, nackte Tilde, Benutzername, URL, Fence still | `make test` |
| Go-Test `TestHostpathsTildeOhneUnixVerlust` | durchgestrichener und am Wort klebender absoluter Host-Pfad bleiben gemeldet | `make test` |
| Go-Test `TestHostpathsExemptTargets` | Ventil nimmt Unix-/Home-relative Funde aus, Windows bleibt | `make test` |
| Go-Test `TestDecode_HostpathsTildeUndExemptTargets` | Präfix mit Tilde und ungültiges Glob ⇒ Fehler | `make test` |
| `d-check` (Dogfooding, `hostpaths` aktiv) | der eigene Bestand bleibt grün | `make doc-check` / `make gates` |

## Re-Evaluierungs-Trigger

Ein Konsument meldet eine host-unabhängige Home-Konvention ohne führenden
Punkt, die so verbreitet ist, dass das Ventil zur Pflicht-Konfiguration
wird — dann eine feste Ausnahmeliste neben der Punkt-Regel erwägen. Ohne
das: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-05 | Angelegt als `Proposed`; `Accepted` erst mit der Closure des Vorgangs, nach Review und Verifikation |
| 2026-10-05 | Nach R1 (HIGH): Entscheidung 2 von der pauschalen Tilde-Sperre auf den Treffer-in-Treffer-Ausschluss umgestellt; Ventil-Globs müssen mit `/` oder `~` beginnen. Noch `Proposed`, Körper daher geändert statt angehängt |
