# Verifikation — slice-248: matrix-Klassen `aussen` + `adaptionsblock` adoptieren, Bestands-Nachzug

- **Rolle:** Verifier (Modul 11) — Frage: „Bauen wir es richtig?" (DoD + Spec + Plan)
- **Gegenstand:** `slice-248` · Range `5a45b0b2..HEAD` — Slice-relevant `42e264c3` (Claim, reiner Move), `95f40dc1` (Implementierung), `6a5e6bed` (R1-Einarbeitung)
- **Datum:** 2026-09-29 · **Modell-ID:** glm-5.3-flash (Z.ai)
- **Eingangs-Kontext:** DoD-Bestätigung des Implementers samt Sensor-Belegen; Slice-Plan; R1-Report (`2026-09-29-slice-248-matrix-nachzug-r1.md`); Messgrundlage slice-246 (40 Befunde)
- **Vorbehalt:** §7 (Closure-Notiz) ist leer und die §6-Risiken stehen auf offen — laut Auftrag der Stand **vor** der Closure. Die Punkte DoD 6 und Nachtlauf-Notierung sind **Closure-Schuld**, keine DoD-Verletzung.

---

## DoD-Prüfung

### DoD 1 — Beide Klassen samt Regeln in `.d-check.yml` — **BESTÄTIGT**

- `adaptionsblock` (`.d-check.yml:558-560`): `paths: ["harness/conventions.md", "harness/conventions/**"]`, `token: 'MR-\d{3}'` — steht **vor** `aussen`.
- `aussen` (`.d-check.yml:568-569`): `paths: ["**"]` — **letzte** Klasse (First-Match), die slice-246-Lektion (zuvor vertauschte Befund-Labels) ist umgesetzt.
- Vier Regeln je `allow: false` (`.d-check.yml:572-575`): `spec-straten→aussen`, `sicht→aussen`, `spec-straten→adaptionsblock`, `sicht→adaptionsblock`.
- Begründungen am eigenen Bestand mit ADR-0097-Bezug (Zeilen 555-557, 561-567, 600-605) — §3.7-Form.
- Wirksamkeit der Regeln nachgewiesen (Gegenprobe, s. DoD 2): die Kanten feuern mit korrekter Klassifikation — `spec-straten → aussen` 8× und `spec-straten → adaptionsblock` 3×, genau nach First-Match-Ordnung.

### DoD 2 — Alle 40 gemessenen Befunde entschieden, Lauf über der Restmenge grün — **BESTÄTIGT**

Eigene Messungen (jeweils `docker run --rm --network none`, Module `matrix` + `ids`, übrige `--disable`):

| Lauf | Gegenstand | Ergebnis |
|---|---|---|
| Ist-Baum (`6a5e6bed`) | produktiv | **917 Dateien, 0 Befunde, Exit 0** |
| Vorher-Baum (`42e264c3`) + neue `.d-check.yml` (Worktree-Kopie) | bewusstes Brechen | **11 Befunde `matrix-forbidden`, Exit 1** — genau die 11 entfernten Links, mit gelesener Ursache (Kante + Fundstelle) |
| Ist-Baum ohne die drei `matrix.exempt-paths`-Neuzugänge | bewusstes Brechen (Status) | **11 Befunde `matrix-inactive`, Exit 1** — ADR-Index ×7, `CHANGELOG.md` ×3, `docs/user/releasing.md` ×1 |

Verteilung der 40 Befunde, im Diff nachgezählt:

- **11 lebende Verweise entfernt** — Diff `42e264c3..HEAD` über `spec/`: 3 MR-Links (`MR-032`, `MR-013`, `MR-022`), 4 `AGENTS.md`-Links (2× Anchor §3.3/§3.8 im Lastenheft, 2× bare im Spezifikation), 2 Baseline-Kopf-Zitate (nun Text-Form mit Version, beide Straten), 1 `harness/README.md`, 1 `packaging/dockerhub/README.md`.
- **17 Historie-Links + 1 nacktes MR-Token** — grep-geprüft: `MR-[0-9]{3}` und outbound-Referenzen stehen nur noch ab `## 7. Historie` (`spec/lastenheft.md:3967`, `spec/spezifikation.md:3545`), null Treffer im lebenden Text; abgedeckt durch `exclude-sections: [Geschichte, "7. Historie"]` (`.d-check.yml:614`). Die drei verbleibenden Testfall-Nennungen im Lastenheft (Zeilen 314, 1285, 1288) sind Inline-Code-Fixtures des ids-Moduls, keine Verweise.
- **11 Status-Fälle** — Break-Test oben: ohne Ausnahme rot (11), mit Ausnahme grün (0). Der file-weite Zuschnitt verliert keine zweite Prüfung: sämtliche 11 Break-Test-Befunde sind `matrix-inactive` aus genau den drei Dateien, keine Link-/Token-Prüfung wird mitgenommen.

Summe 11 + 18 + 11 = 40. Die Restmenge ist produktiv grün, und das Grün ist nicht Inertia: beide Break-Tests zeigen, dass jede betroffene Regel ohne die Fallentscheidung aus dem richtigen Grund rot läuft.

### DoD 3 — Status-Seiteneffekt per ADR-0097 geklärt — **BESTÄTIGT**

- [ADR-0097](../plan/adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md) trägt `**Status:** Accepted`, Datum, Betroffen-Feld und zwei Re-Evaluierungs-Trigger (§Re-Evaluierungs-Trigger, Zeilen 75-83).
- Die drei funktionalen Verweiser sind file-weit in `matrix.exempt-paths` aufgenommen (`.d-check.yml:624`); die ADR-Behauptung „file-weit ist kein Verlust" ist im Break-Test belegt (alle Befunde aus genau diesen Dateien sind Status-Fälle).
- Autor-Frage im Geschichte-Anhang notiert (§3.5-konform: kein Nachtrag in den frozen Körper).

### DoD 4 — `make gates` grün — **BESTÄTIGT**

- Eigener Lauf: Abschlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` — alle zehn gebundenen Targets (`record-gates` als letzter Schritt; die Zeile wird bei jedem Teilausfall nicht erreicht).
- `git status` leer (porcelain 0 Zeilen) vor und nach allen Messungen; die Break-Test-Config wurde byte-identisch restauriert.
- Ergänzend die beiden Range-Gates, die nicht in `gates` hängen: `make adr-check` Exit 0, `make trace-check` Exit 0.

### DoD 5 — Review-Report unter `docs/reviews/` — **BESTÄTIGT**

`docs/reviews/2026-09-29-slice-248-matrix-nachzug-r1.md` liegt vor (3 MEDIUM, 3 LOW, 1 INFO); die Einarbeitung ist unten geprüft.

### DoD 6 — Closure-Notiz + Risiko-Ausgänge — **OFFEN (stand-bedingt)**

§7 ist leer, die beiden §6-Risiken (Baseline-Zitat-Spannungsfeld; Lastenheft-Berührung/MR-032) stehen auf *(offen)*. Laut Auftrag Zustand vor der Closure — **Schuld vor dem `git mv`**, kein Befund. Ebenso offen: die §8-Zusage, den Nachtlauf-Stand in §7 zu notieren.

---

## R1-Folgepunkte (Commit `6a5e6bed`)

| R1 | Behauptung | Verifikation |
|---|---|---|
| M-1 (24 statt 11 entfernt) | Korrektur im ADR-0097-Geschichte-Anhang | **BESTÄTIGT** — Anhang korrigiert auf 11 (3 MR, 4 Agenten-Briefing, 2 Baseline-Köpfe, je 1 README/Packaging), benennt 17 Historie-Links + 1 Token als ausgenommen; deckt sich mit meiner Diff-Zählung |
| M-2 (ADR-0047-Verhältnis fehlt) | Bezug im Geschichte-Anhang | **BESTÄTIGT** — Anhang nennt die Umkehrung von ADR-0047 (Entscheidung 1) für die aussen/adaptionsblock-Lage und hält die ADR-0047-Form für den übrigen Bestand aufrecht. Form: der frozen Körper bleibt unangetastet, der Anhang ist der §3.5 erlaubte Ort — ein `Supersedes:`-Feld wäre ein unzulässiger Körper-Edit gewesen |
| M-3 (Abweichung undeklariert) | MR-0098 als Datei + Index-Zeile | **DATEI BESTÄTIGT — INDEX-ZEILE DEFEKT, neuer Befund V-1** |
| I-1 (fremde shallow-Klon-Grenze) | Notiz im Anhang | **BESTÄTIGT** — Anhang nennt die Grenze ausdrücklich als fremde, die nicht in die Konsequenzen gehört |
| L-2 (Plan §1 abgebrochen) | Vervollständigung | **BESTÄTIGT** — Diff `6a5e6bed` ergänzt die file-weite ADR-0097-Entscheidung und die docs/reviews-Begründung |
| L-3 (Index-Selbstbeleg) | Umstellung | **BESTÄTIGT** — Belege-Spalte der ADR-0097-Zeile trägt nun `DC-FA-MTX-003` + `MR-006` |
| L-1 (Autor-Feld fehlt) | Notiz | **BESTÄTIGT** — Anhang dokumentiert den Zustand und die §3.5-Sperre gegen Nachtragen |

### MR-0098-Datei (M-3-Hälfte „Deklaration")

`harness/conventions/MR-0098-matrix-nimmt-spec-historie-aus-der-referenz-richtung.md` trägt Status `Accepted`, `Ersetzt-Baseline-Regel` (Regel 5 des `grundlagen-referenz-richtung.md` — „in keinem Abschnitt, auch nicht in seiner Historie" — mit Zitat), Geltungsbereich (`matrix.exclude-sections`, heading-genau `7. Historie`, ADR-Geschichte separat), Adaption, Begründung (v5.11.0-CR-Zulässigkeit, 18 frozen Verweise gemessen) und einen **Auflösungs-Trigger**. Die Deklaration selbst ist vollständig.

---

## Neuer Befund V-1 (MEDIUM) — MR-0098 trägt keine Index-Zeile; dangling Fragment im Adaptions-Block

- **Pfad:** `harness/conventions.md:139`
- **Befund:** Die Zeile nach dem MR-070-Eintrag lautet `— **ergänzt** um die Historie-Ausnahme, die der Kanon in der Lastenheft-Historie kanonisch zulässt (v5.11.0); Begründung und Messung in [ADR-0097](…)`. Sie ist **keine Tabellenzeile** (beginnt nicht mit `|`), trägt keine MR-0098-Kennung und keinen Link auf `conventions/MR-0098-….md`. Der Kopf des Adaptions-Blocks verlangt „Eine Zeile je Datei in `harness/conventions/`" — MR-0098 ist als Datei vorhanden, im Index aber nicht vorhanden. Das Fragment liest sich als Rest einer gerissenen Zeile (die letzte Spalte einer vorgesehenen Vier-Spalten-Zeile).
- **Wirkung:** Der Index ist der Ort, den jeder Lauf ohne Öffnen der Eintragsdatei liest („Geltungsbereich und Ersetzt-Baseline-Regel stehen hier, damit ein Agent ohne Öffnen entscheiden kann, ob der Eintrag ihn betrifft"). Für einen Index-Leser existiert die von R1-M-3 geforderte Ausnahme **nicht** — die Deklaration ist in der Datei, aber an der lesenden Stelle stumm. Kein Sensor hält die Index-Zeile gegen die Verzeichnisliste: `make gates` bleibt grün (eigener Lauf), der Bruch ist still — dieselbe Klasse wie ein §3.8-Stillschweigen.
- **Behobung vor der Closure:** eine vollwertige Tabellenzeile für MR-0098 ergänzen (MR · Titel · Geltungsbereich · Ersetzt-Baseline-Regel) und das Fragment in sie aufnehmen.

---

## Plan-vs-Code-Diff

- **Alle vier Plan-Positionen (§3) adressiert:** `.d-check.yml` (Klassen + Regeln + Ausnahmen), drei Spec-Straten (Link-Nachzug, 11 entfernt + Historie-Ausnahme), Status-Fälle (ADR-0097 + exempt-paths). Die vierte Position „Testdatei — .d-check.yml-Beispiele im Spec-Testkorpus prüfen" löst sich als verifizierter **No-op** auf: die beiden Matrix-Beispiele der Spec (`spec/spezifikation.md:229-242`, `:3229-3248`) nutzen fiktionale Klassennamen (`spec-straten`, `contract`, `review`, …) und enumerieren die Repo-Klassen nicht; der configyaml-Validator des `make doc-check`-Laufs bleibt grün. Kein Testartifikat trägt die Prüfung — die Grün-Bindung läuft über den produktiven Lauf.
- **Keine Abgrenzungs-Verletzung:** die drei §1-Ausschlüsse (ids-Weite, exclude-sections-Scoping, stillschweigende Status-Senkung) sind eingehalten — die Status-Ausnahmen sind je §3.6-Fall per ADR-0097 (nicht exempt-paths-Freifeld — die gewählte Mechanik *ist* matrix.exempt-paths, aber ADR-gebunden), `docs/reviews/**` bleibt ohne Ausnahme, die ids-Weite wurde nicht adoptiert.
- **Rückführungs-Trigger (§4) nicht ausgelöst:** der Bestands-Nachzug blieb in einer Sitzung; die Befund-Zahl (40) hat die Session-Grenze nicht sichtbar überschritten.
- **Commit-Zerlegung (§3.3):** `42e264c3` reiner Move-Commit (Claim), Inhalt in `95f40dc1`/`6a5e6bed` — Zerlegung gewahrt.

---

## Fazit

| DoD-Punkt | Verdikt |
|---|---|
| 1 Klassen samt Regeln | bestätigt |
| 2 40 Befunde entschieden, Restmenge grün | bestätigt (eigene Messung + 2 Break-Tests) |
| 3 Status-Seiteneffekt per ADR-0097 | bestätigt |
| 4 `make gates` grün | bestätigt (eigenes Lauf) |
| 5 Review-Report | bestätigt |
| 6 Closure-Notiz + Risiko-Ausgänge | offen — Closure-Schuld vor dem `git mv` |

**Neu:** V-1 (MEDIUM) — MR-0098 ohne Index-Zeile (dangling Fragment in `harness/conventions.md:139`).

Die Mechanik ist vollständig belegt und aus dem richtigen Grund grün. Vor der Closure sind zu erledigen: V-1 beheben, §7 mit Lerneintrag und Risiko-Ausgängen füllen (inkl. Nachtlauf-Notierung §8), dann der reine Move-Commit.
