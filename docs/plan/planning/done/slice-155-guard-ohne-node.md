# Slice slice-155: Der Wächter läuft in der Klasse, die er durchsetzt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`MR-041`](../../../../harness/conventions.md#mr-041) (die benannte
Inkonsistenz und ihr Auflösungs-Trigger);
[`MR-040`](../../../../harness/conventions.md#mr-040),
[`MR-005`](../../../../harness/conventions.md#mr-005) (die Härtungs-Kette);
[`AGENTS.md`](../../../../AGENTS.md) §3.1.

**Berührte Spec-Stellen:** — (Durchsetzungsschicht; das Produkt bleibt
unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

Der Tool-Call-Wächter erzwingt, dass nur die in
[`AGENTS.md`](../../../../AGENTS.md) §3.1 genannte Host-Klasse benutzt wird —
und ist selbst in `node` geschrieben, das nicht dazugehört. Er blockiert seit
[`MR-041`](../../../../harness/conventions.md#mr-041) den Agenten-Aufruf von
`node`, ruft es aber weiter aus dem Hook.

**Eine Regel, deren Durchsetzung außerhalb ihrer eigenen Klasse steht, ist nur
so lange glaubwürdig, wie das dasteht.** Dieser Slice löst den
Auflösungs-Trigger ein: der Wächter läuft mit `bash` und den POSIX-Werkzeugen,
die Host-Abhängigkeit entfällt.

**Es ist gebaut, nur nicht hier.** Das Schwester-Repo `ai-harness-init` fährt
**denselben Wächter** in reinem `bash` + `awk` und führt „kein `node`/`jq`" dort
als Anforderung. Die JSON-Extraktion liegt in einem eigenen POSIX-awk-Skript —
ein zeichenweiser Scanner mit Tiefen- und Key-Stack, der Schlüssel von Werten
unterscheidet. Der Grund ist genau unser Fall: ein Kommando-String kann den
Schlüsselnamen **als Daten** enthalten, und ein Regex-Griff nähme den Treffer im
Freitext. **Diese Fassung ist zu übernehmen, nicht neu zu erfinden** — die
Übersetzung ist ein Abgleich, kein Entwurf.

**Was dort ausserdem gelöst ist und hier gelesen werden muss:** die Fassung des
Schwester-Repos überspringt Zuweisungs-, Wrapper- **und Brace-Group**-Präfixe
und trennt auch am einfachen `&`. Beide Lücken waren hier offen und sind
inzwischen geschlossen; die Übersetzung darf sie nicht wieder verlieren. Und sie
benennt die quote-blinde Falsch-Positiv-Klasse samt Gegenmittel — Inhalt mit
blockierten Wörtern gehört in eine **Datei**, nicht in die Kommandozeile.

**Ein zweites Problem am selben Ort, das dieser Slice mitnehmen kann:** Unser
Wächter antwortet mit `{"decision", "reason"}` auf oberster Ebene. Für
`PreToolUse` ist das die **veraltete** Form; die aktuelle ist
`hookSpecificOutput.permissionDecision` und funktioniert heute nur über eine
Abwärtskompatibilitäts-Abbildung. Der Nachbar-Hook des Schwester-Repos hält das
fest — **dort ist es ebenfalls offen** und als eigener Slice geführt. Ob beides
in einen Zug gehört, entscheidet dieser Slice; getrennt gehalten werden muss die
**Prüfung** von der **Antwortform**.

## 2. Vorgehen

1. **Die Aufgabe zerlegen**, bevor irgendetwas geschrieben wird: (a) das
   Kommando aus der JSON-Eingabe holen, (b) segmentieren, (c) Befehlsposition je
   Segment bestimmen, (d) Sub-Shell-Strings rekursiv prüfen. Nur (a) ist neu —
   (b) bis (d) sind bereits als Logik da und werden übersetzt, nicht neu
   erfunden.
2. **(a) ist der harte Teil und braucht eine Entscheidung, keine Bastelei.**
   Beliebiges JSON in `sed` zu parsen ist ein Parser durch die Hintertür. Die
   Eingabe-Form ist aber vom Harness bestimmt und schmal. Die tragfähige Form
   ist deshalb: **eng lesen, breit scheitern** — was nicht sicher extrahierbar
   ist, ist fail-closed.
3. **Die bestehenden Proben zuerst gegen die neue Fassung fahren**, nicht danach:
   Paketmanager, Host-Go, die Interpreter samt Versions-Suffix, `/usr/bin/`-Pfad,
   `sudo`-Präfix, Sub-Shell in Flag-Bündeln, und die Gegenkontrollen, die still
   bleiben müssen (`git commit -m "docs: pip"`, `docker run img npm test`,
   `make gates`).
4. **Die Falsch-Positiv-Klasse aus [`MR-040`](../../../../harness/conventions.md#mr-040)
   mitmessen** — die Segmentierung an `(` blockiert Kommandos, die einen
   Interpreter-Namen als *Daten* tragen. Ob die neue Fassung sie erbt, gehört
   gemessen und benannt, nicht stillschweigend geändert.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Lockerung.** Was heute blockiert, blockiert danach — die Prüfmenge ist
  die Untergrenze, nicht der Richtwert.
- **Kein zweiter Wächter daneben.** Am Ende trägt **eine** Datei die Prüfung.
- **Kein `jq` und kein anderer Fremd-Parser.** Das wäre dieselbe Abhängigkeit
  unter anderem Namen.

## 4. Definition of Done

- [x] Der Wächter läuft ohne `node`; `command -v node` kommt darin nicht mehr
      vor.
- [x] **Alle** Proben der bisherigen Fassung laufen gegen die neue, mit
      gelesener Meldung — Blocker wie Gegenkontrollen.
- [x] Die Extraktion ist **fail-closed** belegt: eine unlesbare oder
      unerwartete Eingabe blockiert.
- [x] Die Falsch-Positiv-Klasse ist gemessen und benannt — geerbt oder nicht.
- [x] [`AGENTS.md`](../../../../AGENTS.md) §3.1 und
      [`MR-041`](../../../../harness/conventions.md#mr-041)s Auflösungs-Trigger
      sind eingelöst; die Inkonsistenz-Zeile fällt.
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Wächter, der schlechter prüft als vorher, ist eine Verschlechterung mit
  besserem Gewissen.** Die Übersetzung darf keine Prüfung verlieren; die
  Proben-Menge entscheidet das, nicht der Eindruck. — **Ausgang: entfallen.**
  Der unabhängige Review hat die alte gegen die neue Fassung Zeile für Zeile
  gelesen und die Wortlisten mechanisch verglichen: `BLOCKED` (24 Wörter),
  `PREFIXES` (8) und `SHELLS` (5) sind byte-identisch, Segmentierung,
  Präfix-Überspringen, Basename-Kürzung, Versions-Suffix-Muster,
  `-c`-Bündel-Rekursion und Tiefenlimit äquivalent. **Eine** Verhaltensdifferenz
  ist gemessen: die alte Fassung trennte an vertikalem Tabulator und NBSP, die
  neue nicht — dafür trennt sie an `\r`, was die alte nicht tat. Der Verlust ist
  keiner, weil `bash` an VT/NBSP ebenfalls nicht trennt; das neue Modell ist dort
  **treuer**. Das ist die Messung, die [`BEO-011`](../observations.md) verlangt —
  sie trägt „nichts Wesentliches verloren", nicht „prüft dasselbe".
- **JSON in `bash` verleitet zur Bastelei.** Jede Zeile, die „meistens" richtig
  liegt, ist hier eine Sicherheits-Aussage. Fail-closed ist die einzige
  vertretbare Voreinstellung, und sie gehört belegt statt behauptet. —
  **Ausgang: eingetreten, und dreimal.** Der Review fand drei Eingaben, bei
  denen der Extraktor einen Befehl lieferte, der nicht dem entsprach, was das
  Werkzeug ausgeführt hätte: zwei Strings ohne Trenner (der zweite überschrieb
  den ersten), ein `\u`-Escape im Schlüssel (verkürzt dekodiert, Pfad verfehlt),
  und — außerhalb des Extraktors — ein fehlendes `PATH`, das den Wächter mit
  Exit 127 ohne Ausgabe enden ließ. Alle drei sind in `1fd1d6a` geschlossen und
  als Proben gebunden; keiner war über den Befehlsinhalt erreichbar, weil das
  Werkzeug die Hook-Eingabe serialisiert. **Der Ausgang ist trotzdem
  „eingetreten": die Zusage stand, bevor sie stimmte.**

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Extraktion in `bash` nicht
fail-closed hinzubekommen ist — dann bleibt die Inkonsistenz benannt stehen, und
das ist ein Ergebnis.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Durchsetzungsschicht (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-007`](../observations.md) für jeden Exit hinter einer Pipe;
  [`BEO-011`](../observations.md) für jede Aussage darüber, dass die neue
  Fassung „dasselbe" prüfe.

Slice-ID: slice-155. Betroffene IDs: — (kein `DC-`-Bezug;
Durchsetzungsschicht). Module: — . Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Übersetzung einer vorhandenen Prüfung in
die zugelassene Werkzeug-Klasse.

## 9. Closure-Notiz (nach `done/`)

Der Wächter läuft in der Klasse, die er durchsetzt: reines `bash`, das Kommando
aus einem zeichenweisen POSIX-awk-Scanner. `awk` ist die einzige externe
Abhängigkeit, und sie steht in der Host-Klasse, die §3.1 ohnehin nennt.
[`MR-041`](../../../../harness/conventions.md#mr-041) ist damit aufgelöst;
[`MR-042`](../../../../harness/conventions.md#mr-042) trägt die
`node`-Sperre fort, die dort mitdrinhing — ein Move nach `done/` ohne Nachfolger
hätte ihr den Träger genommen.

**Was der Slice über sein Ziel hinaus gefunden hat.** Die Vorlage des
Schwester-Repos hatte eine Lücke, die hier nicht mit übernommen wurde: der
Scanner lief über malformes JSON hinweg und lieferte einen **abgeschnittenen**
Befehl. Der unabhängige Review fand zwei weitere derselben Familie und eine
dritte im Wächter selbst. Alle vier eint dieselbe Form — der Wächter liefert
ein Ergebnis, das wie ein Urteil aussieht und keines ist. Das ist die teuerste
Klasse, die dieser Slice berührt hat, und sie war in keiner DoD-Zeile
vorhergesehen.

**Zwei Aussagen dieses Slice waren zu weit, beide gemessen korrigiert.** Die
Commit-Botschaft von `617fbfa` nennt **24** Proben; es waren **28**, und 24 ist
die Länge der `BLOCKED`-Liste — eine Zahl aus der Nachbarschaft statt aus der
Messung ([`BEO-009`](../observations.md) Richtung a). Die Botschaft steht so im
Verlauf und wird nicht umgeschrieben; die Zahl steht heute bei 33 und wird
gefahren, nicht behauptet. Und
[`MR-042`](../../../../harness/conventions.md#mr-042) sagte „einzige
Host-Abhängigkeit", während der Wächter noch `dirname` und `cat` rief — die
Zusage stimmt erst seit `1fd1d6a`.

**Was liegen bleibt und wo es steht.** Die Antwortform ist die veraltete, und
mit dem Wegfall des `exit 2`-Pfades läuft jetzt **jeder** Block darüber — der
Slice hatte die Frage ausdrücklich offengelassen, der Ausgang ist
[slice-156](../in-progress/slice-156-hook-antwortform.md). Die vier Umgehungs-Klassen
der Segmentierung (Shell-Schlüsselwort als Kopf, Wrapper außerhalb der
Präfix-Liste, wort-interne Splices, escapte Quotes) sind geerbt, gemessen und
als Grenz-Tabelle in
[`MR-042`](../../../../harness/conventions.md#mr-042) geführt; sie bleiben —
ein Wächter, der sie schlösse, wäre ein Parser, und die Regel hängt nicht an
ihm.

**Ein Nebenbefund gehört zum Werkzeug-Einstieg, nicht zum Wächter.** Der
Reviewer hat als Messung berichtet, was ihm vor seinem ersten eigenen
Tool-Aufruf im Kontext lag: `CLAUDE.md` **und** `AGENTS.md`. Damit ist die
Mechanik, die [`MR-043`](../../../../harness/conventions.md#mr-043) behauptet,
an einem Sub-Agenten belegt statt erschlossen — die Zusage hätte sonst genau
die Form gehabt, gegen die dieser Slice angetreten ist.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 526 Dateien, 0 Befunde),
`make guard-probe` (33 Proben, 0 Fehlschläge), `make planning-check` (Exit 0).
Ein unabhängiger Review ist gelaufen; seine fünf MEDIUM und sechs LOW sind in
`1fd1d6a` eingearbeitet, zwei davon als benannte Grenze statt als Fix.
