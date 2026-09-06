# `make guard-probe` — fährt den Tool-Call-Wächter gegen seine Proben

## Vertrag

Proben für den Wächter
([`pretooluse-command-guard.sh`](../../.claude/hooks/pretooluse-command-guard.sh)):
Paketmanager, die Host-Sprachtoolchain, Skript-Interpreter, Brace-Group und Sub-Shell
werden blockiert; **Gegenkontrollen** — legitime Aufrufe, die ein blockiertes
Wort tragen — müssen durchlaufen. Dazu die Fail-closed-Fälle: die des
Extraktors (malformes/abgeschnittenes JSON, `\u`-Escape in Wert und Schlüssel,
zwei Strings ohne Trenner, Müll außerhalb eines Strings) und der des Wächters
selbst (leeres `PATH`).

**Zwei Verdikte stehen neben `pass`/`block`:** `crash` — der Wächter gibt
nichts aus, was sonst von „erlaubt" nicht zu unterscheiden wäre — und `halb`:
er lehnt ab, aber ohne den Exit-Riegel des zweiten Kanals
([`MR-044`](../conventions/MR-044-guard-zwei-kanaele.md)).

## Grenze — was das Grün nicht abdeckt

1. **Die quote-blinde Falsch-Positiv-Klasse** und **die vier
   Umgehungs-Klassen der Segmentierung** sind geführt, nicht behoben —
   [`MR-042`](../conventions/MR-042-guard-in-eigener-klasse.md) trägt sie als
   Tabelle. Der Wächter ist ein **Stolperdraht, keine Sandbox**.
2. **Sein Gegenstand ist eine Werkzeug-Einstellung, keine Repo-Invariante** —
   deshalb ist dieses Target **kein Gate**, obwohl es fail-closed urteilt.
   Kein CI-Lauf ruft den Wächter; ein Lauf ohne dieses Werkzeug ist ungebunden.

## Bindung

kein Gate — **werkzeug-lokal**, bewusst nicht in `gates`. Ohne wiederholbare
Proben wäre die Zusage des Wächters eine Erinnerung.
[`MR-005`](../conventions/MR-005-haertung-gate-nachweis.md) ·
[`MR-040`](../conventions/MR-040-guard-skript-interpreter.md) ·
[`MR-042`](../conventions/MR-042-guard-in-eigener-klasse.md) ·
[`MR-044`](../conventions/MR-044-guard-zwei-kanaele.md)
