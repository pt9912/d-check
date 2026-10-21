# Change Request an `ai-harness-init` — Der Extraktor des Tool-Call-Wächters lässt zwei Parse-Zweifel durch

**Absender:** d-check (Adopter) · **Datum:** 2026-10-21
**Richtung:** ausgehend ([`MR-035`](../../../harness/conventions.md#mr-035)); Ziel ist
hier der Werkzeug-Generator `ai-harness-init`, nicht die Baseline — die Ablage
folgt der Form von [`MR-035`](../../../harness/conventions.md#mr-035), weil kein anderer Ort für ausgehende Bitten besteht
**Ziel:** Datei extract-command.awk im Werkzeug-Repo (Pfade dort:
harness/tools/ und internal/emit/templates/enforce/; hier ohne Inline-Code, weil
sie in diesem Repo nicht liegen), dazu die Probenmenge des Wächters
**Gemessen gegen:** `ai-harness-init` bei `v0.10.0-136-g2b6830c6`; Vergleichsbaum
ist ein damit erzeugtes Repo (`v0.10.0`)
**Stand:** **beantwortet** — angenommen und umgesetzt, siehe
[Antwort](2026-10-21-antwort-ai-harness-init-guard-extraktor-parse-zweifel.md);
ausgeliefert erst mit dem nächsten Release von `ai-harness-init`.

---

## Befund — zwei Eingaben, bei denen der Extraktor nicht fail-closed ist

Der Wächter des erzeugten Repos zieht das Feld `tool_input.command` mit einem
`awk`-Extraktor und soll bei Parse-Zweifel blocken. Zwei Eingaben, die kein
wohlgeformtes JSON sind, lösen keinen Zweifel aus:

| Eingabe (`stdin` des Wächters) | Extraktor-Ausgabe | Verdikt |
| --- | --- | --- |
| `{"tool_input":{"command":"pip x""make gates"}}` | `make gates` (Exit 0) | **durchgelassen**, obwohl `pip x` im Feld stand |
| `{"tool_input":{"command":"a" garbage }}` | `a` (Exit 0) | durchgelassen |

Im ersten Fall geht der erste String verloren: Zwei Strings ohne Trenner sind
kein gültiges JSON, der Extraktor liefert aber den zweiten und meldet Erfolg.
Der Wächter prüft dann `make gates` und lässt durch.

**Gemessen** mit `printf '%s' '<Eingabe>' | awk -f <extract-command.awk>` gegen
beide Fassungen. Der d-check-Extraktor
([`tools/harness/extract-command.awk`](../../../tools/harness/extract-command.awk))
beendet beide Eingaben mit Exit 3 (Parse-Zweifel), der Wächter blockt.

**Schwere, ehrlich:** Claude Code liefert wohlgeformtes JSON; die Lücke ist über
den regulären Aufrufer nicht erreichbar. Sie wiegt als Verletzung der
Zusage *fail-closed bei Parse-Zweifel*, die der Wächter in seinem Kopfkommentar
selbst abgibt — nicht als beobachteter Umgehungsweg.

## Bitte

1. Der Extraktor meldet Parse-Zweifel (Exit ≠ 0) auch bei **zwei Strings ohne
   Trenner** und bei **Zeichen außerhalb eines Strings** hinter dem Wert des
   Feldes. Die d-check-Fassung deckt beides.
2. Die Probenmenge des erzeugten Repos führt diese Fälle. Die Proben in
   [`tools/harness/guard-probe.sh`](../../../tools/harness/guard-probe.sh),
   Abschnitt *Fail-closed am Extraktor*, sind eine Vorlage.

## Nicht gefordert — und warum

- **Die Interpreter-Sperre** (`python`, `node`, `perl`, …): `ai-harness-init`
  sperrt sprachabhängig über `tools/harness/blocked/*`; ohne Sprachfragment
  bleiben `go` und `python` offen. Das ist die Entscheidung des Generators, kein
  Mangel. Ein erster Vergleich hielt das fälschlich für eine Lücke.
- **Brace-Group, Sub-Shell-Rekursion (`bash -c`, `-lc`), `\u`-Escape im
  Befehl oder im Schlüssel:** gemessen, blocken bei beiden gleich — der
  Upstream-Extraktor liest den escapten Schlüssel als `command`, das ist
  konservativ.
- **Der Exit-Riegel** (Block über JSON **und** Nicht-Null-Exit,
  [`MR-044`](../../../harness/conventions/MR-044-guard-zwei-kanaele.md)): der
  Generator antwortet mit `{"decision":"block"}` bei Exit 0. Das wirkt über die
  Abwärtskompatibilität des Aufrufers; ob ein zweiter Kanal gewollt ist, ist
  Sache des Generators und hier nur genannt, nicht erbeten.
- **Leeres `PATH`:** der Vergleich dazu ist nicht belegt (die Probe lief mit
  gelöschtem `PATH` schon vor dem Wächter ins Leere). Dazu steht hier nichts.
