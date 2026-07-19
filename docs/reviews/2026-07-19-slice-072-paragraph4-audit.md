# Frisches §4-Audit gegen den Benutzerhandbuch-Standard (slice-072-Neuvermessung)

**Datum:** 2026-07-19 · **Auditor:** unabhängiger, adversarialer Doku-Reviewer ·
**Gegenstand:** `docs/user/benutzerhandbuch.md`, `## 4. Aufgaben` (aktuell Z. 146–973) ·
**Kriterium:** `docs/user/benutzerhandbuch-standard.md` §2 („Aufgaben statt Funktionen")
und §5 („Schritt-für-Schritt": Ausgangssituation → Ziel → Voraussetzungen →
nummerierte Schritte → erwartetes Ergebnis → Fehler-Hinweise) ·
**Vergleichspunkt (nicht Wahrheit):** slice-072-Audit B-1…B-8 vom 2026-07-17
(DoD in `docs/plan/planning/in-progress/slice-072-handbuch-aufgabenorientierung.md`).

## Methodik / Beleglage

- §4 vollständig gelesen (Z. 146–973), §5 (`trace`-Referenz Z. 1379–1534) und
  §6 (Modul-Tabelle) als Duplikations-Gegenprobe.
- Zeitliche Zuordnung per `git`: Handbuch-Stand zur Audit-Zeit = Commit `d3a8757`
  (letzte slice-071-Doku-Änderung vor dem 2026-07-17-Audit). Post-Audit haben die
  Slices 073–081 (v0.44.1 … v0.51.1) das Handbuch angefasst; **keine** davon war
  slice-072-Remediation — slice-072 steht weiter „In Arbeit", alle DoD-Haken offen.
- Belegter Befund: §4-Kapitelnummerierung und -Titel sind seit `d3a8757` stabil
  (§4.1–§4.13 identisch), d. h. die B-Item-Referenzen mappen 1:1 auf den aktuellen
  Text. Zwei relevante Detail-Belege sind unten je Item benannt.

## Kurzverdikt je B-Item

| Item | Status | Aktuelle Stelle | Ein-Satz-Begründung |
|---|---|---|---|
| B-1 | OFFEN | Z. 333–354 | `order`/`direction`-Block öffnet produkt-innensichtig („Trägt eine Klasse zusätzlich `order` …, prüft d-check auch …"), Titel ist ein Schlüsselname, Feld-/Fail-closed-Enumeration statt Handlungsfolge. |
| B-2 | OFFEN | Z. 356–382 | `token`-Block öffnet mit „`matrix` sieht standardmäßig nur Markdown-Links" — Modul-Innensicht statt Lesersituation. |
| B-3 | OFFEN (schlimmer) | Z. 595–662 vor „Vorgehen" Z. 664 | ~67 Z. Grammatik-/Feldreferenz stehen weiter **vor** dem Vorgehen und duplizieren §5 (Z. 1395–1408); der Block ist seit dem Audit um zwei weitere Referenz-Absätze gewachsen. |
| B-4 | OFFEN | Z. 776–805 | Modalität als Modul beschrieben („klassifiziert der opt-in `modality`-Block …"), Mechanik + Fail-closed-Liste inline, dupliziert §5 (Z. 1454–1473) — nicht als Aufgabe formuliert. |
| B-5 | OFFEN (proliferiert) | Z. 699–718 (+ 618, 645–646, 652–653, 661–662, 774, 807) | Die „Versionsgrenze"-Prosa steht weiter in der Aufgabe; zusätzlich sind post-Audit ≥ 5 weitere inline „Bis vX / Ab vX"-Notizen in §4.12 dazugekommen. |
| B-6 | OFFEN (schlimmer) | Z. 585–915 (~330 Z.) | §4.12 ist weiter **eine** Sektion mit ~8 Themen; post-Audit drei neue Blöcke; doppelte WAISE-Definition (Z. 589–591 und 671–673) besteht. Das vom Audit genannte promovierte H2 existiert real nicht (Beispielausgabe im Fence, s. u.). |
| B-7 | TEILWEISE | Z. 208–211 (Rahmen) / 226–235 (Rest) | Der §4.4-Außenrahmen (Ziel/Voraussetzung) ist gut — aber er war **schon zur Audit-Zeit byte-identisch** so; die vom Audit gemeinte Fähigkeits-Inventar-Passage vor der Entscheidungsfrage ist unverändert (nicht remediiert). |
| B-8 | BEREITS ERLEDIGT | Z. 274, 535, 585 | Die genannten Kapitel-Titel folgen bereits dem von slice-072 selbst gebilligten Muster „Ziel + Schlüssel in Klammern"; kein reiner Flag-Titel mehr auf §4.x-Ebene. |

**Zählung:** 6 OFFEN (B-1, B-2, B-3, B-4, B-5, B-6), 1 TEILWEISE (B-7), 1 ERLEDIGT (B-8).

## Details je B-Item

### B-1 — §4.7 `order`/`direction` (Z. 333–354) — OFFEN

Der Sub-Block trägt den Titel „Richtung innerhalb einer Klasse (`order`/`direction`)"
(Schlüsselname, kein Nutzerziel) und öffnet mit „Trägt eine Klasse zusätzlich
`order` (…) und `direction: no-downward`, prüft d-check auch *klasseninterne*
Referenzen …". Das ist Produkt-Innensicht (was d-check tut, wenn ein Feld gesetzt
ist), nicht Ausgangslage/Ziel des Lesers. Der Schluss (Z. 351–354) enumeriert
Konfigurationsfehler-Bedingungen (Exit 2), die §5 als Referenz führt. Text
gegenüber Audit-Zeit unverändert. Standard §2/§5 verletzt.

### B-2 — §4.7 `token` (Z. 356–382) — OFFEN

Öffnet mit „`matrix` sieht standardmäßig nur **Markdown-Links**." — Modul-Innensicht.
Titel „Token-Referenzen (`token`) und Provenance-Marker" ist ein Feature-Name.
Unverändert seit Audit. Standard §2 verletzt.

### B-3 — §4.12 „Unterstützte Definitionssyntax" (Z. 595–662) — OFFEN, verschärft

Die Grammatik-/Feldreferenz steht weiter **vor** dem „Vorgehen" (Z. 664) — inzwischen
~67 Zeilen. Sie dupliziert §5: die Tabellen-Binding-Regeln (Z. 635–646, „`table.id-column`
und genau eine von `table.text-column` oder `table.text-columns` sind Pflicht …")
stehen wortnah nochmal in §5 (Z. 1395–1408). **Belegt neu seit Audit:** in `d3a8757`
enthielt §4.12 die Zeichenketten „Tabellengrenze" und „Direktiven-Marker" nicht;
beide Absätze (Z. 648–653, 655–662) sind post-Audit angehängt worden — der von B-3
markierte Block ist also nicht gekürzt, sondern gewachsen. Standard §2 (Referenz
gehört in §5) und §5-Reihenfolge verletzt.

### B-4 — §4.12 Modalität (Z. 776–805) — OFFEN

„Modalität (`trace.requirements.modality`)" öffnet mit „Tragen Ihre Anforderungen
RFC-2119-Modalität …, klassifiziert der opt-in `modality`-Block jede Anforderung
anhand …". Das ist die Fähigkeitsbeschreibung, nicht die vom Audit gewünschte
Aufgabe („Nur MUSS-Anforderungen gaten, SOLLTE/KANN advisory lassen"). Die
Klassifikations-Mechanik (Z. 792–798) und die Fail-closed-Liste (Z. 803–805)
liegen inline und duplizieren §5 (Z. 1454–1473). Standard §2 verletzt.

### B-5 — §4.12 Versionsgrenze (Z. 699–718) — OFFEN, proliferiert

Der „Versionsgrenze"-Blockquote („Bis v0.42.0 … Ab v0.43.1 …") steht weiter in der
Aufgabe statt als Fehlerbild (§7) / Versionsaussage (§11). Verschärfend sind
post-Audit weitere inline-Versionsgrenzen in §4.12 dazugekommen: „Ab v0.43.1 liest
… `format: table`" (Z. 618), „(Bis v0.46.0 verlangte d-check drei Bindestriche …)"
(Z. 645–646), „(Bis v0.47.0 verschwand eine so verschluckte relevante Tabelle …)"
(Z. 652–653), „(Bis v0.48.0 brach eine solche Zeile mit Exit 2 ab …)" (Z. 661–662),
„Ab v0.46.0." (Z. 774), „Ab v0.43.1 ist die native Tabellenquelle …" (Z. 807). Das
B-5-Muster ist damit nicht mehr eine Stelle, sondern über die ganze Sektion verteilt.

### B-6 — §4.12 strukturell (Z. 585–915) — OFFEN, verschärft

§4.12 ist unverändert **eine** Sektion, jetzt ~330 Zeilen (Audit: ~300), und bündelt
mindestens acht Themen in einer Aufgabe: Heading-Syntax, native Tabellen, Tabellengrenze,
Direktiven-Marker, Vorgehen/Ergebnis, Versionsgrenze, andere Repo-Konvention (`trace`),
Coverage (`trace.coverage`), Range-Notation, Modalität, Tabellen-Migration und
Kreuzverweis-Konsistenz (`trace.cross-consistency`). Die **doppelte WAISE-Definition**
besteht: einmal im Ziel (Z. 589–591), einmal im Ergebnis (Z. 671–673).

**Korrektur zur Audit-Formulierung:** Das vom Audit genannte „promovierte
`## Kreuzverweis-Konsistenz`-H2" existiert im aktuellen Text **nicht** als echte
Dokument-Überschrift — die Zeile `## Kreuzverweis-Konsistenz` (Z. 883) steht
innerhalb eines Beispielausgabe-Fences (Z. 879–892) und ist Teil der gezeigten
d-check-Ausgabe. Ebenso ist `### F-1 — Repository-Struktur` (Z. 601) ein
Beispiel im Markdown-Fence, keine reale Zwischenüberschrift. Der strukturelle
Kern von B-6 (die Riesen-Sektion auftrennen) bleibt gültig; nur diese eine
Teilbegründung des Audits trifft nicht mehr.

### B-7 — §4.4 Fähigkeits-Inventar vs. Entscheidung (Z. 207–272) — TEILWEISE

Der Task-Hinweis („§4.4 scheint inzwischen Ziel/Voraussetzung/Vorgehen zu tragen")
ist zu prüfen — und irreführend: der §4.4-**Außenrahmen** (Ziel Z. 208–211,
Voraussetzung, Vorgehen, Ergebnis) ist gut, **war aber schon zur Audit-Zeit
byte-identisch** so (belegt: `git show d3a8757` zeigt denselben Rahmen). Er ist also
keine Nach-Audit-Verbesserung. Die Passage, die B-7 real meint, ist der
Harness-Vorlage-Sub-Block: er öffnet mit ~9 Zeilen Fähigkeits-Inventar („Es enthält
die kanonischen `ids`-Muster, die `matrix`-Referenzrichtung und das fixe
Standard-Modulset (`links`, `anchors`, …) samt einem repo-bewussten `planning`-Block
…", Z. 226–234) **bevor** die Entscheidungsfrage „Welche Quelle, hängt von Ihrer
Ausgangslage ab" (Z. 235) und die frisches-vs-bestehendes-Verzweigung (Z. 237/246)
kommt — dieser Teil ist gegenüber dem Audit unverändert. Ergebnis: TEILWEISE — der
Rahmen erfüllt den Standard, das Inventar-vor-Entscheidung ist unremediiert.
Severity niedrig/redaktionell (die Entscheidung *wird* getroffen, nur zu spät).

### B-8 — Flag-/Schlüssel-Titel vs. Ziel-Titel — BEREITS ERLEDIGT

Die vom Audit genannten Kapitel-Titel folgen bereits dem von slice-072 §2 selbst
gebilligten Muster „Titel benennt ein Nutzerziel, Flag/Schlüssel in Klammern":
§4.5 „Regelmodule zu- und abschalten" (reines Ziel), §4.11 „Maschinenlesbare Ausgabe
(`--json` / `--yaml`)" (Ziel + Flag), §4.12 „Anforderungs-Abdeckung prüfen —
Traceability-Matrix (`--trace`)" (Ziel + Flag). Beide, §4.5 und §4.11, trugen diese
Titel schon zur Audit-Zeit. Restposten sind ausschließlich die **inline** Sub-Block-
Titel in §4.7/§4.12, die mit Feature-Namen öffnen (`order`/`direction`, `token`,
`trace.coverage`, `modality`) — die zählen aber bereits unter B-1/B-2/B-4. Der
§4.4-Sub-Titel „Kennungs-Präfix (`--id-prefix`)" (Z. 258) ist ein Substantiv-plus-
Schlüssel statt einer Aufgabe (milde), fällt aber unter die B-7-Region. Auf §4.x-Ebene
ist B-8 erfüllt.

## Neue §4-Mängel (nicht in B-1…B-8)

- **N-1 — §4.12 „Tabellengrenze" (Z. 648–653), NEU post-Audit (v0.47.0).** Edge-Case-
  Grammatik in Funktions-Form („Folgen zwei Tabellen ohne Leerzeile aufeinander,
  beendet der Beginn der zweiten die erste, sofern …") plus Versions-/Fehlerbild
  („Bis v0.47.0 … still"). Gehört als Referenz nach §5, als Fehlerbild nach §7 —
  neue Instanz des B-3/B-5-Musters.
- **N-2 — §4.12 „Direktiven-Marker in einer Tabellenzeile" (Z. 655–662), NEU
  post-Audit (v0.48.0, slice-074).** Reader-Verhaltensbeschreibung + „(Bis v0.48.0
  brach eine solche Zeile mit Exit 2 ab …)". Gleiche Klasse wie N-1.
- **N-3 — §4.12 „Zugesagte Range-/Aufzählungs-Notationen (`ranges: true`)"
  (Z. 765–774), NEU post-Audit (v0.46.0).** Notations-Grammatik (welche Formen
  expandieren, Komma-Kurzform ⇒ Exit 2), abgeschlossen mit „Ab v0.46.0." — reine
  Referenz, gehört nach §5.
- **N-4 — Versions-Prosa über §4.12 verstreut.** Verallgemeinerung von B-5: neben dem
  „Versionsgrenze"-Blockquote trägt §4.12 mindestens sechs weitere inline
  „Bis vX / Ab vX"-Notizen (siehe B-5-Zeilenliste). Standard §9/§11 will
  Versions-/Changelog-Angaben gebündelt, nicht in Aufgaben.
- **N-5 (Beobachtung, kein harter Mangel) — neue Module ohne §4-Aufgabe.** `citations`
  (slice-079) und der `sources`-Content-Pin (slice-080) sind nur in §5/§6 dokumentiert;
  eine §4-Aufgabe „wortgleiche Zitate prüfen" bzw. „eine externe Quelle auf ihren
  Inhalt pinnen" fehlt. Das ist eher korrekte Zurückhaltung (fortgeschrittene opt-in-
  Features) als ein Funktions-statt-Aufgabe-Mangel — aber inkonsistent damit, dass
  `pins` in §5 (Z. 1239) eine aufgaben-betitelte Unterüberschrift bekam.
- **N-6 (mild) — §4.9/§4.11-Überlappung.** Beide erklären `--json`/`--yaml`; §4.9
  bettet zusätzlich JSON-Schemadetail (`reasonText`/`fixCandidate`, Z. 436–469) ein.
  In-Task-Referenz, akzeptabel, trägt aber zur §4-Masse bei.

## Was bewusst Referenz bleibt (kein Mangel, Standard §3 erlaubt §5/§6 als Referenz)

- **§5 vollständig** (`trace`/`coverage`/`modality`/`cross-consistency`-Feldreferenzen,
  Z. 1379–1534) und **§6** (Modul-Tabelle). Standard §3 Punkt 5 „Konfiguration". Legitim.
- **§4.9 „Die drei Ausgaben im Vergleich"** (Z. 471–477) — Entscheidungs-Hilfstabelle
  innerhalb der Aufgabe. Kein Mangel.
- **§4.10 „Was `--repair` behebt"** (Z. 502–506) — die Tabelle ist Teil der
  Reparatur-Entscheidung, in-Task-Referenz. Legitim.
- **§4.13 Target-Aufzählung** (Z. 929–934) — in-Task. Legitim.
- **Kurze YAML-Konfig-Beispiele „so erreichen Sie X"** (z. B. das minimale
  `order`/`direction`-Beispiel Z. 342–349) — aufgabengerecht; problematisch ist nur
  die umgebende feldweise Prosa, nicht das Beispiel selbst.
- **§4.12-Kreuzverweis-Block (`trace.cross-consistency`, Z. 841–914)** — das ist die
  slice-071-Schablone, die slice-072 §2 ausdrücklich als **Vorbild** benennt
  (Ausgangslage → Ziel → nummerierte Schritte → Ergebnis-Deutung → Hinweis). Kein
  Mangel, sondern das Modell für die übrige Auftrennung.

## Gesamt-Urteil: wie viel von slice-072 ist real noch zu tun

Praktisch alles. Von acht Befunden ist genau einer (B-8) erledigt und einer (B-7)
teilerledigt-aber-unremediiert; die sechs inhaltlich schweren (B-1, B-2, B-3, B-4,
B-5, B-6) sind offen. Der Schwerpunkt — der §4.12-Cluster B-3/B-4/B-5/B-6 — ist seit
dem Audit **größer** geworden: drei neue Funktions-Form-Blöcke (N-1/N-2/N-3) und
verstreute Versions-Prosa (N-4) sind dazugekommen, weil die Slices 074/075/076/077
ihre Fähigkeit wie gewohnt an die bestehende §4.12-Aufgabe angehängt haben — exakt
die in slice-072 §4 benannte Ursache. Der billigste dauerhafte Fix bleibt daher der
dort offene Designpunkt: die Release-Prep-Regel „neuer Handbuch-Abschnitt = eigene
Aufgabe mit Ziel/Vorgehen/Ergebnis, keine Anhängung" mit aufzunehmen, sonst erodiert
§4.12 nach jeder Umstellung erneut. Der geplante Aufwand von slice-072 ist eher zu
niedrig als zu hoch angesetzt.
</content>
</invoke>
