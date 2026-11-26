# Architect-Verdikt — welle-98 (Baseline `v7.0.0`, Namensregel), Plan gegen Entscheidungslage

**Rolle:** Architect (Modul 8), frischer Kontext, **vor** Code. Rolleninhaber: pt9912.
**Gegenstand:** [wave-98](../plan/planning/wave-98-baseline-v700.md) §5 (zwölf Fragen) und ihre vier Slice-Pläne
([slice-295](../plan/planning/done/slice-295-baseline-v700-bump.md),
[slice-296](../plan/planning/open/slice-296-baseline-v700-neuerungen-sichtung.md),
[slice-297](../plan/planning/done/slice-297-spezifikation-heisst-specification.md),
[slice-298](../plan/planning/in-progress/slice-298-neue-namen-wave-archive-anker.md)), Stand `e0b06112`.
**Datum:** 2026-11-26.
**Bindend, vom Auftraggeber am 2026-11-26 vorgegeben und hier nicht neu aufgerollt:** Dateinamen englisch, Prosa
deutsch, „ab jetzt, Bestand bleibt“ · Neuerungen nur sichten, keine ohne Einzelentscheid · Reihenfolge
295 → 297 → 298 → 296 · Frage 2 (`wave-98` erst in slice-298) · Frage 3 (Konfiguration statt Produkt, 38 Register-Zeilen
ungewächtert, benannt) · Frage 6 (Lastenheft nachziehen).
**Nummerierung:** die der Welle §5, nicht die der Anfrage.
**Geändert wurden die Pläne (Welle §5, slice-295, 297, 298) und der Entscheidungs-Bestand** (Folge-ADR
[ADR-0117](../plan/adr/0117-wellen-namen-ueber-die-konfiguration.md), `Proposed`, samt Index-Zeile). Kein Code, keine Spec,
kein Lastenheft, kein Eintrag unter `harness/conventions/`.

## Gesamturteil

Der Plan trägt; er ist klein genug und die kleinsten Antworten der Welle tragen fast alle. **Drei Stellen
brauchten mehr als die Planner-Antwort:**

- **Frage 3** ist eine **Senkung einer Prüfregel** (`planning` bewacht den Bestand nicht mehr). Der Plan führte sie
  nur in einem MR; `AGENTS.md` §3.6 verlangt eine ADR — **V3**, hier vollzogen als ADR-0117.
- **Frage 5**: die Kopplung von slice-297 an den Koexistenz-Eintrag aus slice-298 (Risiko 4) ist unnötig — die
  Form für das Ventil steht schon (ADR-0025, Vorgänger `docs/user/releasing.md`) — **V1**.
- **Frage 8** ist mit „nur Wellen“ **nicht klein**: sie ließe für wellenlose Slices und Reviews eine dauerhafte Abweichung
  von den Namen, die die Baseline nennt — **V1**, Antwort: alle drei Modi.

**Keine weitere Entscheidung des Auftraggebers nötig.** Wo er später entscheiden kann, steht es je Frage als
„größere Antwort“ mit ihrem Mehraufwand.

Verdikt-Kürzel (Modul 8, Konflikt-Pfad): **V1** = Entscheidung gilt, Plan hat falsch behauptet/ausgelassen ·
**V2** = Folge-Entscheidung nötig · **V3** = Lockerung legitim, aber undokumentiert · **E** = Frage beantwortet,
kein Konflikt.

## Befunde und Antworten

### Frage 1 — Form der Koexistenz-Adaption — **E**

**Ein eigener Eintrag**, Nummer **MR-082** (der Bump-Eintrag aus slice-295 ist **MR-081**; höchste vergebene war
MR-080). Geschrieben **im Lauf von slice-298 vom Implementer**, nicht jetzt: er nennt gemessene Zustände der
Konfiguration nach der Umstellung, die heute noch nicht gelten (`AGENTS.md` §5 Regel 13). Das Verdikt legt fest:

- **Ersetzt-Baseline-Regel:** `grundlagen-harness-dateien.md` §Verzeichniskonvention, Absatz *Dateinamen sind
  englisch* (Link mit Abschnitts-Anker in den vendorten Baum von `v7.0.0`).
- **Adaption:** je Gegenstand eine Zeile — `wave-<NN>` neben `welle-<NN>`, `archive.zip` neben `archiv.zip`,
  `seit wave-<NN>` neben `seit welle-<NN>` — mit der Auskunft, welcher Sensor **beide** Formen liest und welcher
  nur die neue. Die Senkung von `planning` wird **nicht nacherzählt**, der Eintrag zeigt auf
  [ADR-0117](../plan/adr/0117-wellen-namen-ueber-die-konfiguration.md) (Regel 18).
- **Auflösungs-Trigger:** der Bestand ist weg — `git ls-files | grep -c 'welle-[0-9]\|archiv\.zip'` ergibt 0
  **und** kein lebendes Dokument trägt `seit welle-<NN>`. **Nicht „permanent“:** permanent ist den Urteilsfällen
  vorbehalten (`AGENTS.md` §3.2, §3.4), und hier gibt es einen Zähler. Er feuert nur auf einen Beschluss des
  Auftraggebers; bis dahin gilt der Eintrag.
- Kein Eintrag je Gegenstand (drei Dateien, ein Trigger, nicht unabhängig auflösbar).

### Frage 2 — `wave-98` — **E**, Ablauf nach dem Bindenden

`welle-98` heißt erst in slice-298 `wave-98`, zusammen mit der Konfiguration. **Keine Lücke dazwischen:** bis dahin
bewacht `planning` die Welle über `welle-*`, danach über `wave-*`; Move, Roadmap-Zeiger und Konfiguration stehen im
selben Arbeitsbaum, der Hook prüft ihn. Der Plan-Zeiger „bis dahin unsichtbar“ galt für einen frühen Move.
**Plan-Änderung** (slice-298 §3): reiner `git mv` der Welle-Datei zuerst, dann der Nachzug **aller** Verweise
(Zählung wiederholen: nach slice-295 und slice-297 sind es mehr als die 14 Zeilen der Eröffnung — Closure-Notizen,
Reports, `**Welle:**`-Felder), die Namen `welle-98-results.md` in §3 (e) und §7 der Welle werden `wave-98-results.md`.
Die Closure-Artefakte tragen die neuen Namen: `wave-98-results.md`, `done/wave-98/archive.zip`, `seit wave-98`.

### Frage 3 — Reichweite der Sensor-Anpassung — **E**, Senkung **V3** (vollzogen: ADR-0117)

Die Konfiguration. Die 13 Zeilen des Plans sind richtig geschnitten (die Zählung wiederholt der Implementer). Der Preis
ist gemessen (38 Register-Zeilen, 38 Ergebnisnotizen) und nach `AGENTS.md` §3.6 eine **ADR**, kein Eintrag allein.
**ADR-0117** trägt die Senkung, ihre Grenze (`planning` allein; `matrix`, Closure-Profil und Werkzeug lesen beide Formen)
und die **größere Antwort** (das Modul liest mehrere Präfixe: acht Stellen `wavePrefix`, Konfigurations-Prüfung,
rund 100 Test-Zeilen, Lastenheft, Spezifikation, ADR, Handbuch, Release) als Wiedervorlage des Auftraggebers.
ADR-0117 wird bei der Closure von slice-298 `Accepted`; der `**Status:**`-Übergang ist dem ADR-Gate erlaubt.

### Frage 4 — Commit-Zerlegung der Link-Nachzüge — **E**

**Zulässig, und keine Ausnahme.** `AGENTS.md` §3.3 schützt die Rename-Erkennung der **bewegten Datei** (Begründung
dort: 50 %-Similarity); sie bleibt `R100`. Der Satz, andere Dateien dürften im Move-Commit mitreisen, solange die bewegte
Datei unverändert bleibt, steht dort beim Lifecycle-Übergang — die Begründung trägt dieselbe Lesart für den Rename.
Die Spezifikation selbst (161 Zeilen Selbst-Links) ist der **zweite** Commit: Fall 1 des Paragraphen, unverändert.
Der erste Commit ist allein nicht baubar (`go:embed`) und gehört mit dem zweiten in **denselben Push**; der Hook prüft
den Arbeitsbaum. **Kein neuer MR** wie MR-059/062/063: dort wich die Zerlegung ab (reine `D`/`A`-Paare), hier nicht.
**Notiert** in slice-297 §3 (Reihenfolge, Schritt 3) und in dieser Datei. **Akzeptiertes Negativ:** `AGENTS.md` §3.3
nennt den Fall nicht ausdrücklich; eine Zeile dort unterbleibt, weil er einmalig ist (Umbenennung einer
Spec-Datei, die `Accepted`-ADRs verlinken). Kehrt er wieder (zweite Umbenennung), gehört die Lesart in §3.3.

### Frage 5 — eingefrorene Link-Ziele — **E**, Kopplung **V1**

- **Spezifikation:** **Nachzug der Link-Ziele** in 64 ADRs und 39 eingefrorenen Dateien, im Move-Commit. Der Grund,
  warum hier nachgezogen wird und beim alten Baum nicht: der Link zeigte immer auf die **lebende** Datei, nach dem Nachzug
  tut er es weiter; die Bedeutung ändert sich nicht.
- **Inline-Code (104 Nennungen):** **ein** `codepaths.ignore-refs`-Eintrag `spec/spezifikation.md` mit der Grenze im
  Kommentar — Form und Vorgänger sind ADR-0025 und der Eintrag für `docs/user/releasing.md` (Kommentar „GRENZE: … bleibt
  das still“). Ein MR ist dafür **nicht** nötig. **V1:** slice-297 Risiko 4 koppelte den Ausgang an den Koexistenz-Eintrag
  aus slice-298 — eine Vorwärts-Abhängigkeit zu einem späteren Slice; sie entfällt, der Kommentar genügt.
- **Alter Baum (7 Dateien, fünf ADRs, slice-274, ein Architect-Verdikt):** **Ventil**, sieben datei-skopierte
  `links.ignore-refs`-Einträge in der Form von MR-069. Hier wäre ein Nachzug **falsch**: die Links zitieren den Stand ihrer Zeit,
  ~~und ein Anker kann im neuen Baum fehlen (eine Überschrift in `modul-14` entfällt)~~ **[Prämisse falsch, berichtigt und der Ausgang für slice-274 und das Verdikt zu welle-93 geändert: Nachtrag V1 am Ende dieser Datei]**. Der Bump-Eintrag nennt die sieben
  Einträge als Wachstum des Ventils (MR-069, Trigger „drei Hebungen ohne neuen Eintrag“ beginnt neu).
- **Größere Antwort:** `links.ignore-refs` je Klasse auch für die Spezifikation (3–4 Einträge statt 80 Zeilen Nachzug):
  spart den Eingriff in die Belege, lässt aber nach dem Rename tote Links in ihnen stehen. Nicht gewählt; kein Mehraufwand
  bei Wechsel außer der Konfiguration.
- Wortgrenze: ein Tausch greift nur nach `/`, `(` oder Zeilenanfang (zwei Dateinamen enden auf `spezifikation.md`).

### Frage 6 — Lastenheft Z. 4762 — **E**, vom Auftraggeber entschieden

Nachziehen. **Plan-Auftrag an den Implementer in slice-297** (Rang 1, nicht mein Akt):
`spec/lastenheft.md` Z. 4762 nennt `spec/specification.md`; `**Version:**` 0.115.1 → **0.115.2**; **eine** Historie-Zeile
nach MR-032 und der Form der übrigen („Redaktionell: Pfad der Messmethode zu DC-QA-01 folgt der Umbenennung der
Spezifikation; keine inhaltliche Änderung, kein DC-* ändert sich“, Verweis `—`). Kein weiterer Satz im Lastenheft.
Der Punkt in slice-297 ist damit **unbedingt**, die „bei nein“-Klausel entfällt.

### Frage 7 — Kommentarzeilen — **E**

**Nur was funktioniert:** `manual.go`, die zwei Tests, `tools/image-test.sh` Z. 163. Die **94 Kommentarzeilen**
(90 Go in 33 Dateien, `Dockerfile` 1, `Makefile` 1, `tools/*.sh` 2) bleiben. **Akzeptiertes Negativ, mit Grund:**
`AGENTS.md` §3.7 verlangt, dass ein Kommentar beschreibt, was da ist — diese Rang-Zeiger nennen künftig eine Datei,
die es nicht mehr gibt. Sie haben keine Wirkung auf Verhalten; kein Sensor liest sie; der Preis ist ein toter Zeiger, den
ein Leser mit dem Dateinamen sofort findet. Keine Folgepflicht (wer eine dieser Zeilen ohnehin anfasst, schreibt den neuen
Namen). **Größere Antwort:** ein mechanischer, kommentar-only Nachzug (ein Commit, `make test` und `make lint`; prüfbar
mit der Regel aus dem Closure-Trigger (d1)); Mehraufwand ein Commit und 36 Dateien im Diff. Die Konsumenten-Beispiele
(4 Zeilen) bleiben in jedem Fall. **Folge für den Plan:** der Trigger (d1) der Welle verliert die Kommentar-Klausel.

### Frage 8 — Reichweite von `archive.zip` — **V1**

**Alle drei Modi.** Die Baseline nennt nicht nur `done/<wave-id>/archive.zip`, sondern auch
`done/slice-<Kennung>-archive.zip` (`modul-05` §Lifecycle, `modul-10`). „Nur Wellen“ hinterließe für jeden neuen
wellenlosen Slice eine dauerhafte Abweichung von einem Namen, den die Baseline wörtlich führt — und der Eintrag zur
Koexistenz müsste sie tragen. **Regel:** der Welle-Modus **folgt dem Präfix der Kennung** (`welle-<NN>` →
`archiv.zip`, `wave-<NN>` → `archive.zip`; der Bestands-Lauf bleibt byte-gleich); der Slice- und der Review-Modus
schreiben **`-archive.zip`** (sie haben kein Präfix, das sie unterscheidet, jedes Archiv, das jetzt entsteht, ist neu; der
Review-Modus folgt dem Slice-Modus, die Baseline nennt ihn nicht wörtlich). Mehraufwand gegenüber „nur Wellen“: null
zusätzliche Stellen — die 14 Stellen und 16 Test-Zeilen des Plans zählen alle drei Modi. **Plan-Änderung:** slice-298
§1 (Ausschluss entfällt) und §2 (Liefer-Punkt Werkzeug).

### Frage 9 — alter Baum fällt im Bump — **E**

Die Praxis bleibt. **Die Abweichung vom Kanon ist schon deklariert:** MR-021 (1) legt fest, dass der Bump den alten
Baum **entfernt**, permanent. Der Bump-Eintrag nennt den **Vor-Bump-Commit** und zeigt auf MR-021 (1); ein eigener
Satz zur Abweichung wäre eine zweite Wahrheit.

### Frage 10 — Freshness-Audit der Adaptionen — **E**

Nein. Die Obergrenze von 17 Treffern bleibt Messung. **Akzeptiertes Negativ:** ob eine der 41 Adaptionen durch die 19
bewegten Regelwerk-Dateien betroffen ist, bleibt ungeprüft; die Wiederkehr ist nicht spurlos — der Wellenplan steht in
`done/`, und der nächste Major liest ihn. Die größere Antwort ist ein weiterer Slice mit je einer Frage pro Treffer.

### Frage 11 — Ort der Entscheidungsvorlage — **E**

`docs/plan/planning/baseline-v700-neuerungen.md`, flach. Kein Glob des Moduls `planning` trifft den Namen (`welle-*` und
`wave-*`); eine Datei neben der Welle, die die Closure überlebt, ist die kleinste Lösung. **Kein Eintrag für die
Klasse:** eine Vorlage zu einem Major ist kein wiederkehrender Gegenstand mit eigenem Sensor.

### Frage 12 — `stub.go`, zwei Kommentarzeilen — **E**

Im Bump nachziehen, wie geplant. Sie nennen die **Kopplung** an die Vorlagen-Namen (nicht einen Rang-Zeiger wie in Frage 7),
`make archive-wave-test` hält den Rest, kein Verhalten ändert sich. slice-298 fasst `stub.go` ohnehin an; der Nachzug bleibt
trotzdem in slice-295, weil dort die Namen entfallen.

## Plan-Änderungen (vollzogen)

Welle §3 (b) (c) (d1), §5 (zwölf Antworten statt zwölf Fragen); slice-295 (Nummer MR-081, Frage 9, Verdikt-Verweis);
slice-297 (Frage 4 Lesart, Frage 6 unbedingt, Risiko 4 entkoppelt, Frage 7 entschieden); slice-298 (Frage 1 Form und
Nummer, Frage 2 Ablauf, ADR-0117, Frage 8). slice-296 bleibt unverändert; sein Start-Trigger zu Frage 11 ist mit diesem
Verdikt erfüllt.

## Nicht entschieden, weil nicht nötig

Die Übernahme einer der übrigen Neuerungen (Entscheid 2 des Auftraggebers) — Gegenstand von slice-296.

## Nachtrag 2026-11-26 — Frage 5, alter Baum (Verdikt **V1**, Rolleninhaber: pt9912)

**Anlass:** Reviewer-Befund F-1 (MEDIUM) zu slice-295
([Report](2026-11-26-slice-295-baseline-v700-bump-r1.md)): die Begründung des Ventils trifft die Dateien nicht.

**Berichtigung der Prämisse (V1 — die Entscheidung gilt für fünf Dateien, ich habe für zwei falsch behauptet).** Der Satz
„ein Anker kann im neuen Baum fehlen (eine Überschrift in `modul-14` entfällt)" in Frage 5 war nicht gemessen und trifft
keine der sieben Dateien: keine verlinkt `modul-14`. **Gemessen** (Wegwerf-Klon im Scratchpad, in den sieben Dateien nur
`baseline/v6.18.1/` → `baseline/v7.0.0/` getauscht, die sieben Ventil-Einträge entfernt): `make doc-check` meldet
`1318 Datei(en) geprüft, 0 Befund(e)` — **alle sieben Ziele samt Anker lösen auf**. Damit hängt die Antwort nicht an der Auflösung,
sondern daran, ob die **zitierte Stelle wortgleich** ist (gemessen je Abschnitt, `awk` über die Überschrift bis zur nächsten):

| Datei(en) | Ziel in `v7.0.0` | zitierte Stelle gegenüber `v6.18.1` | Antwort |
|---|---|---|---|
| ADRs 0112–0116 | `modul-04` §Ziel-Form: ADR (MADR) | **nicht wortgleich** — der Abschnitt nennt die Vorlage jetzt `NNNN-title.template.md` statt `NNNN-titel.template.md` | **Ventil bleibt** |
| `slice-274` (Rot-Beleg) | `modul-11` §Bewusstes Brechen | wortgleich (Abschnitts-Prüfsumme identisch) | **Retarget**, Ventil-Eintrag entfällt |
| Verdikt zu welle-93 | `modul-05` §Ziel-Form: Slice | wortgleich (57 Zeilen, `diff` leer; „Schnitt nach Lieferwert" vorhanden) | **Retarget**, Ventil-Eintrag entfällt |

**Warum `slice-274` und das Verdikt jetzt doch retargetet werden (anders als in Frage 5 gesagt).** Es ist der Vorgänger
([MR-080](../../harness/conventions/done/MR-080-baseline-v6181.md), dort die Ausnahme von der Frozen-Regel für `slice-274`) mit
seiner eigenen Bedingung: ein Link-Ziel-Tausch ist nur dann kein verfälschtes Zitat, wenn die zitierte Stelle gelesen
wortgleich ist — bei beiden Dateien gemessen. Damit sinkt die Gate-Senkung ([MR-069](../../harness/conventions/MR-069-ignore-refs-ist-eine-deklarierte-gate-senkung.md))
um zwei Dateien, in denen künftige Links auf den alten Tag still würden; der Preis sind zwei Zeilen in zwei eingefrorenen
Dateien, derselbe Preis wie bei `slice-277` (reiner Link-Ziel-Wechsel auf eine lebende Datei). Der Trigger von MR-069 startet
**so oder so neu** (die fünf ADRs bleiben hinter dem Ventil); die Ersparnis ist die der zwei Dateien, nicht der Zähler.

**Warum die fünf `Accepted`-ADRs im Ventil bleiben, obwohl ein Pfad-Nachzug erlaubt wäre.** Gemessen: erlaubt ist er —
im selben Klon, die fünf Links auf `v7.0.0`, danach `make adr-check RANGE=738d2bca..HEAD`: `1312 Datei(en) geprüft, 0 Befund(e)`
([ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md)). **Kleiner ist er nicht**, aus drei Gründen: (1) der
zitierte Abschnitt ist nicht wortgleich (Tabelle), der Link würde einen Stand behaupten, dem die ADR nicht folgte —
`AGENTS.md` §3.5 schützt gerade das; (2) der Nachzug verschöbe nur die Wand: beim nächsten Bump stünden dieselben fünf
Links wieder vor derselben Wahl (MR-069 Grenze 3: die ADR-Klasse hat keinen Träger für Kennung statt Adresse), das Ventil
wird je Tag **einmal** bezahlt; (3) der Trigger von MR-069 braucht den Träger **und** drei Hebungen ohne Wachstum, der
Nachzug ließe ihn trotzdem nicht feuern. **Größere Antwort, nur benannt:** für die ADR-Klasse einen Träger für die
Kennung-statt-Adresse-Form schaffen (ADR-Vorlagen-Nutzung, MR-069 Grenze 3); dann entfielen künftige Ventil-Einträge. Nicht
Gegenstand von slice-295.

**Plan-Auftrag an den Implementer (slice-295; die Änderungen unten ändern nichts außer den benannten Stellen).**

1. **Retarget, je eine Zeile, nur das Link-Ziel** (`baseline/v6.18.1/` → `baseline/v7.0.0/`, Anker unverändert):
   `docs/plan/planning/done/wellenlos/slice-274-trace-requirements-als-liste-von-quellen.md` Z. 185 und
   `docs/reviews/2026-10-30-welle-93-architect-verdikt.md` Z. 114. Eigener Commit (Pfad-Nachzug), Kennung `slice-295, MR-081`.
2. **`.d-check.yml` Z. 409–412:** die zwei Einträge (je `- in:` und `refs:`) für diese zwei Dateien entfernen (die fünf ADR-Einträge bleiben);
   der Kommentar (Z. 392–399) lautet danach: *„Tombstone der Baseline-Hebung auf v7.0.0 (MR-081). Datei-skopiert, wie beim
   v6.9.0-Block: fünf `Accepted`-ADRs (0112 bis 0116) tragen je einen Markdown-Link auf den entfernten v6.18.1-Baum in ihrer
   Regeln-Zeile (modul-04 §Ziel-Form: ADR). Sie zitieren den Stand ihrer Zeit — der Abschnitt nennt die Vorlage in v7.0.0
   anders (NNNN-title.template.md), ein Nachziehen verfälschte das Zitat (MR-069, AGENTS.md §3.5); der alte Baum liegt in der
   git-Historie (MR-052). GRENZE: in diesen fünf Dateien wird jeder künftige Link auf .harness/baseline/v6.18.1/ still."*
3. **`harness/conventions/MR-081-baseline-v700.md`**, Absatz *`ignore-refs` wächst* (Z. 199–216) und die Verweise darauf
   (Z. 193 „die sieben Ventil-Einträge", Z. 159–162 „Frozen hinter dem Ventil (5 / 5)" und „byte-stabil (17 Dateien, 51 …)"):
   sieben → **fünf** Einträge, Bestand **35** statt 37 über unverändert **13** Tags (Kommando wie bisher, danach neu gemessen),
   Auflösungs-Trigger von MR-069 „beginnt neu" bleibt (die ADRs). Den Grund durch den Wortlaut oben ersetzen (Stand ihrer Zeit,
   **nicht wortgleich**; „ein Anker kann fehlen" und die Nennung von `modul-14` streichen). Den Satz zu `slice-274`
   (Z. 212–215, „diesmal bleibt die Datei unberührt …") ersetzen durch: *„`slice-274` und das Verdikt zu welle-93 sind wie bei
   MR-080 auf den neuen Baum retargetet: die zitierte Stelle ist gelesen wortgleich (modul-11 §Bewusstes Brechen, modul-05
   §Ziel-Form: Slice), und der Tausch hält die Gate-Senkung um zwei Dateien kleiner."* Der Absatz *Eine Ausnahme von der
   Frozen-Regel, nicht im Ventil* (Z. 217 ff.) nennt dann **drei** Ausnahmen (`slice-277`-Verifikation, `slice-274`, Verdikt
   `welle-93`; je 1 Datei, 1 Zeile, nur Link-Ziel); die Frozen-Zählung der Spiegel-Klassen (Z. 159–162) nennt zwei der 17
   byte-stabilen Dateien als Ausnahme (15 bleiben, je Klasse neu zählen, Kommando nennen). Die Treffermenge (10 Dateien)
   ändert sich nicht, nur ihr Klammertext „(die sieben Ventil-Einträge)" → fünf.
4. **Nebenspuren:** `slice-295` §3 Z. 69–72 und 140 („je eingefrorener Datei … gemessen 7") als Plan-Änderung nachtragen
   (fünf; zwei Retargets nach Verdikt-Nachtrag); `welle-98` §3 Z. 213 und §5 Z. 324–325 („sieben datei-skopierte Einträge")
   auf fünf. Das Verdikt bleibt Beleg und wird nicht umgeschrieben.
5. **Prüfen vor Handoff:** `make doc-check` Exit 0; `make adr-check RANGE=738d2bca..HEAD` Exit 0 (die ADRs bleiben
   byte-stabil); `git grep -c 'refs: \[\?"\?\.harness/baseline/v' .d-check.yml` = 35; Review-Report F-1 an den Reviewer zurück.

**Nicht entschieden, weil nicht nötig:** F-2 und F-3 des Reports (LOW, Sache des Implementers, kein Konflikt).

## Nachtrag 2026-11-26 (V2) — Beispiel-Kennung `wave-99` in ADR-0117 (Rolleninhaber: pt9912)

**Anlass.** Der Planner hat vor Commit 2 von slice-298 gemeldet: [ADR-0117](../plan/adr/0117-wellen-namen-ueber-die-konfiguration.md)
trägt in Z. 34 und 77 die Beispiel-Kennung `wave-99`. Mit der neuen `matrix`-Klasse `wave` (Token `wave-\d{2,}`) und dem Verbot
`adr → wave` wären das zwei `matrix-forbidden` — das Plan-Risiko 5 („0 Befunde nach der Konfigurationsumstellung") trägt dann nicht.

**Gemessen** (Wegwerf-Klon im Scratchpad, Stand `7de32a36`; Klasse `wave` mit `paths: docs/plan/planning/**/wave-*.md`, `token: 'wave-\d{2,}'`
und den drei Verboten `spec-straten`/`sicht`/`adr → wave` eingefügt):

| Zustand der ADR | `make doc-check` |
|---|---|
| unverändert (`wave-99`, Z. 34 und 77) | `1328 Datei(en) geprüft, 2 Befund(e)`, beide `matrix-forbidden` (Z. 34, 77), Exit 2 |
| Beispiel als Platzhalter `wave-<NN>` | `1328 Datei(en) geprüft, 0 Befund(e)`, Exit 0 |

Außer ADR-0117 trägt keine Datei unter `spec/` oder `docs/plan/adr/` eine `wave-<Ziffern>`-Kennung (`git grep -n -E 'wave-[0-9]{2,}' -- spec docs/plan/adr`
= die zwei Zeilen). Die übrigen Treffer des Plans (Welle, Verdikt, Plan, Test) liegen außerhalb der drei Quell-Klassen und sind kein Gegenstand der Regel.

**Verdikt: (a), die Entscheidung gilt und die Nennung wird umgeschrieben.** Die ADR ist `Proposed`; `AGENTS.md` §3.5 schützt erst ab `Accepted`.
Die Kennung `99` war eine Illustration, keine Aussage der Entscheidung: gemeint ist „eine neue Welle in der neuen Schreibweise" — `wave-<NN>` sagt dasselbe
und ist als Platzhalter keine Kennung des Token-Ziels. Der Inhalt der Entscheidung, die Senkung (V3) und die Fitness Function bleiben unverändert;
geändert sind genau die zwei Beispiel-Schreibweisen in der Tabelle und im Kontext-Absatz.

**Warum nicht (b) und (c).** (b) der Provenance-Marker (`<!-- d-check:status-provenance -->`, den die Meldung selbst anbietet) würde die Nennung als
Herkunfts-Zeiger deklarieren — sie ist aber ein Beispiel, kein Zeiger; die Deklaration wäre falsch, und der Marker nimmt je Zeile die ganze Klasse aus
(Grenze: eine echte Nennung neben dem Beispiel würde ebenfalls still). (c) ein engeres Verbot wäre eine Prüfregel-Senkung (`AGENTS.md` §3.6) und damit ADR-pflichtig —
für ein Beispiel, das sich ohne Verlust umschreiben lässt, die größere Antwort. **Größere Antwort, nur benannt:** ein Pfad-/Platzhalter-Träger in der
ADR-Klasse für Kennungs-Beispiele; ein einmaliger Fall, kein Gegenstand.

**Akzeptiertes Negativ.** Die Platzhalter-Form `wave-<NN>` ist von der Token-Erkennung nicht erfasst (kein `\d`). Wer künftig in einer ADR wieder eine konkrete
Kennung als Beispiel schreibt, bekommt den Befund — das ist das gewünschte Verhalten, nicht eine Lücke.

**Plan-Auftrag an den Planner (Plan-Datei von slice-298, nicht durch den Architect geändert).**

1. Risiko 5 und die Vorbedingung vor Commit 2: *„Vorbedingung erfüllt — ADR-0117 trägt seit dem Umschreiben (Nachtrag V2) in Z. 34 und 77 `wave-<NN>` statt `wave-99`;
   gemessen: `make doc-check` mit Klasse `wave` und Verbot `adr → wave` meldet 0 Befunde (1328 Dateien)."* Der Treffer `ADR-0117` in der Treffermenge
   (`git grep -l -E 'wave-[0-9]{2,}'`) entfällt dort damit: vier statt fünf Dateien.
2. Commit 2 verlangt von der ADR nichts mehr; die Prüfung „`make doc-check` 0 Befunde nach der Umstellung" bleibt DoD-Punkt des Implementers.
3. ADR-0117 bleibt `Proposed`; der Übergang nach `Accepted` bleibt bei der Closure von slice-298 (Plan §Ergebnis, ADR-Zeile der Änderungstabelle).
