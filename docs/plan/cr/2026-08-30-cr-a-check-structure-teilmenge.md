# Eingehender CR aus `a-check` — `structure`: die geprüfte Menge deklarierbar machen

**Absender:** a-check (Adopter von d-check).
**Eingegangen:** 2026-08-30, über den Auftraggeber.
**Gegenstand:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
— zwei Optionen an bestehenden Regeln.

Dieses Dokument hält den CR **wie empfangen** fest. Die Bewertung steht nicht
hier, sondern im Slice, der ihn aufnimmt — ein CR-Dokument trägt Bitte und
Beleg, nicht die Antwort darauf.

---

## Anlass, vom Absender gemessen

`structure` wendet `max-tasks` und die Abschnitts-Auswahl auf die
**vorgefundene** Menge an. Der Absender verlangt an zwei Stellen eine
**erklärte Teilmenge**:

1. **Die Größen-Regel zählt nicht alle Task-Items.** Ein Lauf von
   `max-tasks: 3` gegen die Slice-Pläne des Adopters (`v0.67.0`,
   `--enable structure`) liefert **9 Befunde** `section-oversized`. Acht davon
   sind die Slices 108–115 — jeder Slice, der seit Inkrafttreten der
   Konstanten-Regel geschrieben wurde: je **7** Task-Items, davon **2–3**
   Liefer-Punkte. Der Bestand ist regelkonform; der Zähler misst das Falsche.
2. **Grandfathering lebt innerhalb einer Datei.** Wer
   `require-all: [Happy, Boundary, Negative]` über `section-pattern`
   durchsetzt, muss die bei Einführung bestehenden Anforderungen ausnehmen —
   dort **19**. Sie sind Abschnitte **einer** Datei; `exempt-paths` ist
   datei-granular und kann sie nicht erreichen. Die Alternative wäre, den
   Bestand umzuschreiben — das träfe die Form statt der Substanz.

**Nicht adopter-spezifisch.** Die Slice-Vorlage der Baseline liefert eine DoD
mit **neun** Checkboxen aus, von denen **sechs** pro Slice konstant sind, und
schreibt eine Zeile darüber selbst: *„Gate-Läufe und die vier Closure-Pflichten
darunter zählen nicht mit."* Wer die Vorlage benutzt **und** die Größen-Regel
prüft, bekommt zwangsläufig Falsch-Positive — auf jedem neuen Slice, während
der Altbestand grün bleibt. Der Sensor wird über die Zeit unbrauchbar, nicht
sofort.

## Vertrag, wie beantragt

Zwei Optionen an bestehenden Regeln; ohne sie byte-identisches Verhalten. Keine
neuen Grund-Codes — beide **verkleinern** nur die geprüfte Menge.

<!-- d-check-test:not-config: Antragstext des Absenders, kein .d-check.yml-Input -->
```yaml
structure:
  - files: "docs/plan/planning/**/slice-*.md"
    section-pattern: '^## .*(DoD|Definition of Done)'
    max-tasks: 3
    tasks-ignore-pattern: '(make gates|Closure-Notiz|Beobachtungs-Register|Risiko aus)'
    #   Task-Items, die dieses RE2 treffen, zaehlen fuer max-tasks NICHT mit.
    #   Ohne den Schluessel: alle Items zaehlen — heutiges Verhalten.

  - files: "spec/lastenheft.md"
    section-pattern: '^### AC-'
    sections: each
    require-all: [Happy, Boundary, Negative]
    exempt-sections: '^AC-[A-Z]+-0(0[1-9]|1[0-9])\b'
    #   Abschnitte, deren UEBERSCHRIFTSTEXT dieses RE2 trifft, prueft DIESE Regel
    #   nicht — Geschwister von exempt-paths, eine Granularitaetsstufe tiefer.
    #   Der Stichtag steht damit in der Konfiguration statt in einem Skript.
```

**Warum keine eigenen Module.** Beides ist dieselbe Frage mit erklärter
Grundmenge, keine neue Frage. `tasks-ignore-pattern` gehört zu `max-tasks` wie
`order-column` zu `order`; `exempt-sections` ist das Geschwister von
`exempt-paths` eine Stufe tiefer. Als Optionen bleibt der Default unberührt.

**Fence-Treue gilt weiter.** Beide Muster dürfen Code-Blöcke und Inline-Code
nicht sehen — sonst kippt genau die Eigenschaft, die den Modul-Weg gegenüber
der Skript-Variante ausgezeichnet hat. Ein Adopter, der über sein eigenes
Regelwerk schreibt, zitiert seine Konstanten-Begriffe ständig in Backticks.

## Paritäts-Mutations-Beleg, vom Absender vorgelegt

| Probe | Erwartung |
|---|---|
| Slice mit 3 Liefer-Punkten + 4 Konstanten, `max-tasks: 3` **mit** `tasks-ignore-pattern` | grün (heute rot — die acht gemessenen Fälle) |
| derselbe Slice mit **vier** Liefer-Punkten | rot |
| `tasks-ignore-pattern` abwesend | rot — heutiges Verhalten unverändert |
| grandfatherte Anforderung ohne `Boundary`-Marke, von `exempt-sections` getroffen | grün |
| neue Anforderung ohne `Boundary`-Marke, nicht getroffen | rot |

## Abgrenzung des Absenders

**Nicht Teil des Antrags:** Grandfathering **ab einer Nummer** als
Werkzeug-Begriff — das ist über Globs bzw. RE2 ausdrückbar. Ebenfalls nicht:
abschnitts-**übergreifende** Bedingungen („jedes Risiko aus §6 trägt einen
Ausgang in §7"). Das ist eine andere Frage und braucht einen eigenen Antrag.
