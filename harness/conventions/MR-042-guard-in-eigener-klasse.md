# MR-042 — Der Wächter läuft in der Klasse, die er durchsetzt (löst MR-041 auf)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon verlangt für die
  Durchsetzungsschicht die Härtung als neuen Eintrag
  ([`modul-13-quality-gates.md` §Guard-Härtung](../../.harness/baseline/v6.0.0/regelwerk/modul-13-quality-gates.md))
  und beschreibt den Befehls-Guard als Artefakt
  ([`grundlagen-durchsetzungsschicht.md` §Das vollständige Artefakt-Set](../../.harness/baseline/v6.0.0/regelwerk/grundlagen-durchsetzungsschicht.md)).
  *Womit* er geschrieben ist, sagt er nicht — die Form-Frage tritt die Rangliste
  an diesen Speicher ab.
- **Datum:** 2026-08-27
- **Geltungsbereich:**
  [`.claude/hooks/pretooluse-command-guard.sh`](../../.claude/hooks/pretooluse-command-guard.sh),
  [`tools/harness/extract-command.awk`](../../tools/harness/extract-command.awk),
  [`tools/harness/guard-probe.sh`](../../tools/harness/guard-probe.sh),
  das `guard-probe`-Target im [`Makefile`](../../Makefile),
  [`AGENTS.md`](../../AGENTS.md) §3.1 und §4,
  [`harness/README.md`](../README.md).
- **Adaption:** Der Wächter ist in `bash` geschrieben und holt das Kommando aus
  der Hook-Eingabe über einen zeichenweisen POSIX-`awk`-Scanner. Seine einzige
  Host-Abhängigkeit ist damit `awk` — eines der POSIX-Standardwerkzeuge, die
  [`AGENTS.md`](../../AGENTS.md) §3.1 ohnehin als Host-Klasse nennt. Die Sperre
  für `node`, `deno` und `bun` aus [`MR-041`](../conventions.md#mr-041) gilt
  unverändert weiter; dieser Eintrag trägt sie fort.

  **Fail-closed ist die Voreinstellung, nicht die Ausnahme:** fehlt `awk`, fehlt
  der Extraktor, oder meldet er Parse-Zweifel, blockiert der Wächter. Damit das
  auch ohne Host-`PATH` gilt, kommt er ohne jedes weitere externe Programm aus —
  Pfad, Eingabe und Ausgabe laufen über bash-Builtins; ein fehlendes `dirname`
  oder `cat` hätte ihn mit Exit 127 und **ohne** Ausgabe enden lassen, also ohne
  Block.

  **Drei Formen führen zum Zweifel, und jede war einmal ein stilles
  Durchwinken:** ein Zeichen, das außerhalb einer Zeichenkette nichts in JSON zu
  suchen hat (der Befehl käme **abgeschnitten** zurück); ein zweiter String im
  selben Container ohne Trenner (der zweite **überschriebe** den ersten); ein
  `\u`-Escape im Wert **oder im Schlüssel** (der Schlüssel träfe verkürzt den
  falschen Pfad, der Wert bliebe ungelesen). Alle drei enden mit Exit 3. Das ist
  die gefährlichste Klasse, weil das Ergebnis wie ein Urteil aussieht.

  **Die Proben sind ein `make`-Target** ([`AGENTS.md`](../../AGENTS.md) §4,
  `make guard-probe`), kein Gate: der Wächter ist eine Werkzeug-Einstellung, und
  ein Gate darüber wäre eine Zusage, die ein Lauf ohne dieses Werkzeug nicht
  hält. Ohne wiederholbare Proben wäre seine Zusage aber eine Erinnerung.

  **Die quote-blinde Falsch-Positiv-Klasse aus
  [`MR-040`](../conventions.md#mr-040) ist geerbt, gemessen und bleibt** — sie
  steht als Probe, nicht als Fehler. Ein Wächter, der Daten von Befehlen sicher
  unterscheiden wollte, wäre ein Parser.

  **Umfangs-Grenze — was der Wächter NICHT sieht, gemessen an dieser Fassung.**
  Vier Klassen erreichen ein gelistetes Werkzeug, ohne dass die Segmentierung es
  als Segment-Kopf sieht; alle vier sind **geerbt**, keine ist neu:

  | Klasse | Beispiel, das durchläuft |
  |---|---|
  | Shell-Schlüsselwort als Kopf | `if true; then pip install x; fi` · `for f in a; do go build; done` · `! pip install x` |
  | Wrapper außerhalb der Präfix-Liste | `nohup pip install x` · `timeout 5 go build` · `stdbuf -o0 pip install x` |
  | wort-interne Splices | `p"i"p install x` · `g\o build` |
  | escapte Quotes in der Verschachtelung | `bash -c "bash -c \"bash -c pip\""` |

  Die Präfix-Liste ist eine **Liste** und damit von Natur aus unvollständig; die
  Schlüsselwort-Klasse verlangte einen Grammatik-Begriff, den ein Stolperdraht
  nicht führt. Beides bleibt so und steht hier, statt still zu gelten
  ([`modul-13-quality-gates.md` §Guard-Härtung](../../.harness/baseline/v6.0.0/regelwerk/modul-13-quality-gates.md):
  die Grenz-Zeile wird mitgezogen).
- **Begründung:** [`MR-041`](../conventions.md#mr-041) benannte die
  Inkonsistenz und ließ sie stehen: eine Regel, deren Durchsetzung außerhalb
  ihrer eigenen Klasse steht, ist nur so lange glaubwürdig, wie das dasteht.
  Sie steht jetzt nicht mehr, weil es die Inkonsistenz nicht mehr gibt.

  **Warum das kein Edit an [`MR-041`](../conventions.md#mr-041) ist:** dessen
  Auflösungs-Trigger war für den **ganzen** Eintrag formuliert, obwohl der
  Eintrag zwei Dinge trägt — die eingelöste Inkonsistenz und die weiter geltende
  Sperre. Ein Move nach `done/` ohne Nachfolger nähme der Sperre ihren Träger.
  Der Kanon löst genau das mit einem neuen Eintrag
  ([`modul-02-harness-bootstrap.md` §Freshness-Audit](../../.harness/baseline/v6.0.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)).
- **Löst auf:** [`MR-041`](../conventions.md#mr-041)
- **Auflösungs-Trigger:** der Kanon schreibt der Durchsetzungsschicht selbst
  eine Werkzeug-Klasse vor. Dann gilt seine.
