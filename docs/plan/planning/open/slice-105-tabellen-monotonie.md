# Slice slice-105: Chronologie-Monotonie als siebte `structure`-Bedingung

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-77-chronologie-ordnung (zugeordnet bei der Eröffnung).

**Bezug:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md) (Modul-Schnitt),
[ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
(Schnitt-Kriterium: Einzelmodul-Frage ⇒ bestehende Anforderung ändern),
[ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
(geteilte Lexik bindet ihre Konsumenten); Anlass **BEO-005** im
[Beobachtungs-Register](../observations.md); begleitende ADR entsteht mit diesem
Slice.

**Autor:** pt9912. **Datum:** 2026-08-21.

---

## 1. Ziel

Eine chronologische Tabelle kippt still ihre Richtung (**BEO-005**, am selben Tag
an drei Tabellen gefunden): wer eine Zeile anfügt, schaut auf die Zeile daneben
statt auf die Regel, und danach führt dieselbe Tabelle zwei gegenläufige Blöcke.
Kein Gate liest Reihenfolge. Dieser Slice macht die Ordnung maschinell: ein
**typisierter Monotonie-Vergleich der Schlüsselspalte** als **siebte Bedingung**
im Modul `structure` — dieselbe Eingabe-Achse (eine Regel benennt Dateien und
Abschnitt), darum Erweiterung von
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
statt neues Modul oder neues Kürzel
([ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)-Kriterium).

**Der Typ ist Pflicht, kein Komfort** (gemessen, nicht behauptet): ein naiver
String-Vergleich meldet drei **korrekt** sortierte Tabellen rot
(`0.10.0 → 0.9.0`, `v0.10.0 → v0.9.0`, `1.9 → 1.10`); der typisierte Vergleich
findet am Stand **vor** der Heilung alle drei gekippten Tabellen
(14 · 6 · 7 Verletzungen) und am Stand danach null.

## 2. Die drei Entscheidungen, die keine Details sind

Aus dem Register-Eintrag übernommen; die Begründungen trägt die begleitende ADR.

1. **Rohe Abschnitts-Zeilen.** Die Bereinigung aus Schritt 5 des Algorithmus
   leert Inline-Code — und die Schlüsselzelle von `version.md` steht genau dort
   (`` `v0.60.0` ``, die aktuelle Zeile zusätzlich mit HTML-Anker). Die Bedingung
   liest deshalb die **rohen** Zeilen des Abschnitts — eine **ausdrückliche,
   vertraglich benannte Ausnahme** von „alle Bedingungen arbeiten auf diesem
   Text“. Fence-Bewusstsein bleibt: Tabellenzeilen in Fenced-Code deklarieren
   nichts (geteilte Lexik, siehe 2.).
2. **Zell-Adresse über die geteilte Tabellen-Lexik.** Die Bedingung braucht als
   einzige eine Zell-Adresse: Tabellenzeilen erkennen, Kopf-/Trennzeile
   überspringen, Spalte indizieren. Genau diese Antworten existieren —
   `tableRowLine` (geteilt: `targets`, `planning.waves`) sowie Trennzeilen-/
   Kopfzeilen-Erkennung und Zell-Splitting (heute privat in `planning.waves`).
   Sie ein zweites Mal zu bauen wäre BEO-003 in Reinform. Mit `structure` hat
   die Tabellen-Lexik ihren **dritten** Konsumenten und bekommt einen
   **Kopplungs-Test** (Form aus
   [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md):
   dieselbe Eingabe durch alle Konsumenten derselben Frage, Fehlschlag bei
   abweichender Antwort) statt weiterer Einzel-Assertionen.
3. **Richtung je Regel, Befund statt stillem Übersprung.** Das Handbuch-§11 ist
   als einzige Bestandstabelle **aufsteigend** — die Richtung ist
   Regel-Konfiguration, kein Modul-Default. Eine Zelle, die sich nicht
   typisieren lässt (oder eine Typ-Mischung in der Schlüsselspalte), ist ein
   **Befund**, kein stilles Auslassen — sonst schaltete ein Tippfehler in einer
   Zelle die Prüfung der restlichen Tabelle wortlos ab (dieselbe fail-closed-
   Disziplin wie beim Leerlauf-Befund des Moduls).

**Vorschlag der Vertragsform** (endgültige Namen entscheidet die ADR):

- `structure[].key-order`: `desc` | `asc` — setzt die Bedingung scharf; anderer
  Wert ⇒ Exit 2.
- `structure[].key-column`: int ≥ 1 (Default 1) — 1-basierte Schlüsselspalte;
  explizit < 1 oder gesetzt ohne `key-order` ⇒ Exit 2.
- **Typen** (geschlossene Menge): ISO-Datum (`JJJJ-MM-TT`) und Punkt-Version
  (optionales `v`-Präfix, ≥ 2 numerische Segmente, segmentweise numerisch
  verglichen; `1.10 > 1.9`, `0.60.2 > 0.60.0`). Getypt wird der **erste**
  passende Token der rohen Zelle — so überleben Inline-Code-Backticks und der
  wandernde `<a id>`-Anker der `version.md`.
- **Monotonie nicht-strikt** (gleiche Schlüssel erlaubt — mehrere Releases an
  einem Tag sind der Normalfall), geprüft je **zusammenhängender** Tabelle im
  Abschnitt, über benachbarte Datenzeilen.
- **Zwei neue Grund-Codes** (die Reparaturen sind verschieden):
  `section-order-broken` (Zeile bricht die Monotonie; `line` = die brechende
  Zeile) und `section-key-untyped` (Zelle nicht typisierbar bzw. Typ-Mischung;
  `line` = die Zelle).

## 3. Vorgehen

1. **Tabellen-Lexik konsolidieren, nicht kopieren:** Trennzeilen-/
   Kopfzeilen-Erkennung und Zell-Splitting aus `planning_waves.go` an den
   geteilten Ort (`markdown.go`) heben; `planning.waves` konsumiert sie von
   dort. Kopplungs-Test über alle drei Konsumenten.
2. **Bedingung implementieren** (Schritt-Erweiterung in
   §`DC-FA-STRUCT-001.a`): rohe Abschnitts-Zeilen, fence-bewusste
   Tabellenzeilen-Auswahl, Typisierung, Richtungs-Vergleich, zwei Grund-Codes.
3. **Grund-Codes im Lockstep:** `AllReasons()`, `--doctor`-Klartexte,
   Spezifikation §4 — im **selben** Commit.
4. **Messen, dann selbst aktivieren:** je Kandidaten-Tabelle ein Lauf (die sechs
   Bestandstabellen aus §5), Aufnahme nur bei 0 Befunden bzw. mit begründetem
   Nein; Bindepunkt siehe Abnahme-Punkt 2.

## 3a. Spiegel dieser Semantik-Änderung ([`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten))

**Vor** dem ersten Editor aufgeschrieben; die Grund-Code-Menge wächst um zwei.
Nach der slice-099-Lehre aus dem Repo abgeleitet (`grep` nach dem letzten
Bedingungs-Zugang `require-pattern` und nach `section-marker-missing`), nicht
aus dem Gedächtnis.

| Spiegel | berührt? | was genau |
|---|---|---|
| Anforderung [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) | **ja** | Bedingungs-Tabelle +1 Zeile, Roh-Lese-Ausnahme, fail-closed-Aufzählung (+3 Config-Ränder), neue Akzeptanzkriterien, Out-of-Scope (Typ-Menge geschlossen) |
| Algorithmus [`DC-FA-STRUCT-001.a`](../../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure) | **ja** | Schritt 1 (Config-Ränder), Schritt 5 (Roh-Lese-Ausnahme benennen), Schritt 6 (siebte Bedingung samt Zell-Adresse und Typisierung) |
| §2-Config-Schema | **ja** | zwei neue `structure[]`-Schlüssel |
| §4-Grund-Code-Tabelle | **ja, zwei Zeilen** | `section-order-broken`, `section-key-untyped` |
| `AllReasons()` / `reasonTexts()` (Doctor-Klartexte) | **ja, zwei** | im **selben** Commit wie §4 (Lockstep) |
| `--print-config`-Vorlage | **prüfen** | nur falls das `structure`-Gerüst Bedingungs-Schlüssel zeigt |
| `--suggest-config` | **nein** | Modul-Menge unverändert (kein neues Modul) |
| `--print-mk`-Target-Liste | **nein** | `doc-structure` existiert; Target-Zahl unverändert |
| Benutzerhandbuch | **ja** | §5-`structure`-Block (Bedingungs-Liste/Config-Beispiel) + §11-Zeile (Release-Prep; **chronologisch unter die letzte** — ab diesem Slice maschinell) |
| README (beide Sprachen) | **nein** (Feat) | Modul-Zahl/-Liste unverändert; Versions-Pin ist Release-Prep |
| `operations.md` | **nein** | Modul-Enumerationen unverändert |
| `.d-check.yml` | **ja** | Selbst-Aktivierung nach Messung (anders als bei slice-099): `structure` in die `modules`-Liste + Regel-Liste für die sechs Bestandstabellen |
| DC-QA-03-Modullisten-Go-Test | **prüfen** | `structure` ist hermetisch — die Netzlos-Liste muss die neue `modules`-Zeile tragen |
| `.d-check.closure.yml` | **prüfen** | ausdrücklich messen statt „unberührt“ behaupten (slice-099-Lehre: genau diese Zeile war falsch) |
| `AGENTS.md` / `harness/README.md` | **nur falls** sich eine Gate-Beschreibung ändert | `doc-check`-Zeile bleibt; kein neues Target |
| Lastenheft §7 + Versions-Bump | **ja** | 0.60.2 → 0.61.0 |
| Spezifikation §7 | **ja** | neue Datums-Zeile **oben** (absteigend) |
| ADR-Index | **ja** | begleitende ADR |
| CHANGELOG / `version.md` | **ja** (Release-Prep) | Minor-Notiz; `version.md`-Zeile + Anker-Wanderung |

## 4. Definition of Done

- [ ] Siebte Bedingung vollständig: typisierter Vergleich (Datum/Version),
      Richtung je Regel, rohe Zell-Lesung (fence-bewusste Zeilen-Auswahl),
      Kopf-/Trennzeilen übersprungen, nicht-strikte Monotonie je Tabelle,
      Befund je Bruch-Zeile, Befund für untypisierbare Zelle/Typ-Mischung,
      fail-closed Config-Ränder (Exit 2).
- [ ] Tabellen-Lexik geteilt statt verdreifacht: Trennzeilen-/Kopfzeilen-/
      Zell-Antwort an einem Ort, **Kopplungs-Test** über alle drei Konsumenten
      (`targets`, `planning.waves`, `structure`); Befundsatz der beiden
      Bestands-Konsumenten unverändert (Belegt durch bestehende Tests).
- [ ] **Zwei** neue Grund-Codes im Lockstep (`AllReasons()`, §4,
      Doctor-Klartexte); jedes neue Akzeptanzkriterium als Test.
- [ ] **Retro-Beleg mit dem Produkt:** am Stand vor der Heilung
      (welle-73-Ära) melden die drei gekippten Tabellen 14 · 6 · 7
      Verletzungen, am heutigen Stand null; die naive Gegenprobe (String-
      Vergleich rot auf korrekt Sortiertem) ist als Testfall festgehalten.
- [ ] Selbst-Aktivierung nach Messung (sechs Bestandstabellen, Handbuch
      aufsteigend); `make gates` grün; unabhängiger Review; Release als
      **Minor** (opt-in; ohne die neuen Schlüssel byte-identisch).

## 5. Die sechs Bestandstabellen (Mess- und Aktivierungs-Kandidaten)

| Tabelle | Schlüsselspalte | Richtung |
|---|---|---|
| `spec/spezifikation.md` §7 Historie | 1 (Datum) | absteigend |
| `spec/lastenheft.md` §7 Historie | 1 (Version) | absteigend |
| Roadmap §Historische Trigger-Verschiebungen | 1 (Datum) | absteigend |
| Roadmap §Abgeschlossene Wellen | 2 (Abschluss-Datum) | absteigend |
| `version.md` §Verlauf | 1 (Version, in Inline-Code) | absteigend |
| Benutzerhandbuch §11 | 1 (Handbuch-Version) | **aufsteigend** |

## 6. Abnahme-Punkte / Risiken

1. **Vertragsnamen** (`key-order`/`key-column`, Code-Namen) — Entscheid in der
   begleitenden ADR, Vorschlag oben.
2. **Bindepunkt der Selbst-Aktivierung:** `.d-check.yml` (Inner-Loop — die
   Zeilen entstehen in Feat-/Release-Prep-Commits, dort soll es rot werden)
   statt Closure-Profil. Gegenprüfung beim Aktivieren: kein
   `planning-check`-Nebeneffekt (das Modul läuft nur unter `--enable structure`,
   also in `doc-check`... **zu verifizieren**: `doc-check` fährt die
   `modules`-Liste — `structure` gehört hinein).
3. **Geteilte Lexik-Hebung berührt ausgeliefertes Verhalten** (`planning.waves`
   seit v0.59.0, `targets` seit v0.44.x): reine Code-Bewegung, der
   Kopplungs-Test plus Bestandstests sind die Absicherung (BEO-003-Klasse).
4. **Zwei getrennte Tabellen im selben Abschnitt** werden je für sich geprüft —
   zwei einzeln sortierte, gegenläufige Tabellen sind damit **nicht** erkennbar
   (benannte Grenze; der Anlassfall war ein Bruch **innerhalb** einer Tabelle).
5. **BEO-004-Frage** (welche Eingaben liest die Bedingung, die das Modul nicht
   scannt?): keine neuen Dateien — dieselben Kandidaten, aber **roher** Text
   derselben Datei; genau darum steht die Ausnahme im Vertrag, nicht nur im
   Code.

## 7. Trigger

**Start** (`open` → `in-progress`): Freigabe des Auftraggebers **und** WIP-Slot
frei — beide mit der Eröffnung von welle-77 erfüllt (2026-08-21). Der
Vorschau-Trigger der Welle („nach welle-75: die Tabellen-Lexik hat ihren dritten
Konsumenten verdient") ist seit 2026-08-16 eingetreten.

**Rückführungen:** `in-progress` → `next`, falls die Lexik-Hebung einen eigenen
Refactor-Slice verlangt (Befundsatz-Änderung bei `targets`/`planning.waves`).

## 8. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide
  unter dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-21: BEO-002/003/004
  verkörpert, BEO-005 offen): **BEO-005** ist der Anlass dieses Slice — die
  Welle entscheidet bei der Closure über Streichung (mechanisiert für die
  aktivierten Tabellen) oder Verbleib. **BEO-002** wirkt verkörpert als
  [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  (§3a, aus dem Repo abgeleitet). **BEO-003** wirkt als Auftrag: die
  Tabellen-Lexik wird geteilt statt verdreifacht, mit Kopplung statt Aufzählung
  ([ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)).
  **BEO-004** wirkt als Frage — beantwortet in §6 Punkt 5 (Roh-Text-Ausnahme als
  benannte Vertragsgrenze).

## 9. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die Anforderung existiert; dieser Slice
erweitert sie um eine gemessen belegte Bedingung. Doc führt (Lastenheft/
Spezifikation/ADR zuerst), Code folgt.

## 10. Closure-Notiz (nach `done/`)

_Ausstehend._
