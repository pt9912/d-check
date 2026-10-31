# Welle welle-94: Modul `pinsync` — gleiche Versions-Pins in mehreren Dateien gegeneinander prüfen

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-10-31.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Ein Repo, das dieselbe Version an mehreren Stellen führt (etwa `GO_VERSION` im
`Dockerfile` und in einer eingebundenen `.mk`-Datei), prüft die Stellen mit
d-check gegeneinander: das neue opt-in-Modul `pinsync` fängt pro benanntem Pin je
Quelle (Datei und regulärer Ausdruck mit genau einer Fanggruppe) einen Wert und
meldet (1) ungleiche Werte und (2) einen Ausdruck, der in seiner Datei keinen oder
mehr als einen Treffer hat (GitHub-Issue 14). Der Lauf bleibt hermetisch: nur
Dateizugriff, kein Netz, kein git, kein Ausführen, keine Schreibzugriffe — wie das
Modul `targets`. Ein führendes `v` steuert der Ausdruck selbst (`v?` vor der
Fanggruppe), einen Schalter dafür gibt es nicht; geprüft wird Gleichheit, keine
Aktualität.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Das GitHub-Issue 14 liegt offen vor (`gh issue list --state open`), und der
  Auftraggeber hat seine Bearbeitung als Welle angefordert (2026-10-31).
- `origin/main` trägt nichts, was dem lokalen Hauptzweig fehlt (Abstand 0/4 nach
  `git fetch`: die vier Closure-Commits von welle-93 sind lokal und ungepusht), die
  Lifecycle-Verzeichnisse `next/` und `in-progress/` tragen keinen Slice, `open/`
  nur die drei dieser Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle sind `done/`.
- `make fullbuild` grün.
- **Das Mehr gegenüber den Slice-DoDs** (repo-weit, in keiner einzelnen DoD):
  - **Ohne den Block nichts Neues:** `make blackbox-probe` gegen das Vorher-Image
    der letzten Release-Version (`REF=v0.90.0`) zeigt für den Eigenlauf und ein
    Fixture ohne `pinsync`-Konfiguration **dieselbe** Ausgabe (stdout, stderr und
    Exit getrennt). Ein neues opt-in-Modul darf keinen bisher grünen Lauf
    verschieben; jeder Slice belegt seine Hälfte, dass keine die andere verschoben
    hat, belegt erst der gemeinsame Lauf.
  - **Der Anlass trägt am eigenen Bestand:** das Modul läuft im eigenen Profil über
    die zwei Pins, die das Repo heute doppelt führt (`GO_VERSION` und
    `GOLANGCI_LINT_VERSION`, je im `Dockerfile` und im `Makefile`), im Gate
    `make doc-check`; die Probe ist **beidseitig**: ein gehobener Wert in nur einer
    der beiden Dateien färbt den Lauf rot mit dem Grund-Code des Moduls, der
    unveränderte Bestand bleibt grün. Beobachtbar wird sie durch die Tabelle
    Probe → Ergebnis in `welle-94-results.md`.
  - **Die Betriebsbedingungen des Issues am gebauten Image:** der Lauf unter
    `:ro` und gegen ein Repo, das einem anderen Benutzer gehört, liefert dasselbe
    Ergebnis wie gegen ein eigenes beschreibbares Repo (`make image-test`).
- Das Trigger-Audit (seit slice-241) ist über die Welle-Closure vollzogen und in
  `welle-94-results.md` belegt.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-286`](done/slice-286-pinsync-anforderung-und-adr.md) | Anforderung, Spezifikation und ADR für das Modul `pinsync`; Entscheide festschreiben | neue Anforderung des Moduls (Kennung vergibt der Slice), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| [`slice-287`](done/slice-287-pinsync-modulkern.md) | Modulkern: Konfiguration lesen, Quellen fangen, Werte vergleichen, Grund-Codes | die Anforderung aus [`slice-286`](done/slice-286-pinsync-anforderung-und-adr.md) |
| [`slice-288`](open/slice-288-pinsync-eigenes-profil.md) | Eigenes Profil: die zwei Doppel-Pins von `Dockerfile` und `Makefile` im Gate, beidseitig belegt | die Anforderung aus [`slice-286`](done/slice-286-pinsync-anforderung-und-adr.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle.
- Wird blockiert von: keiner.
- **Reihenfolge der Slices ist fachlich bindend:** erst
  [`slice-286`](done/slice-286-pinsync-anforderung-und-adr.md) (die Anforderung und
  die Entscheide stehen, bevor Code entsteht), dann
  [`slice-287`](done/slice-287-pinsync-modulkern.md), dann
  [`slice-288`](open/slice-288-pinsync-eigenes-profil.md) (braucht das fertige
  Modul im Image). Das WIP-Limit 1 hält sie ohnehin nacheinander.
- **Entscheide, die [`slice-286`](done/slice-286-pinsync-anforderung-und-adr.md)
  vor dem Schreiben der Anforderung festschreibt** (das Issue legt sie nicht fest;
  der Slice benennt sie mit Empfehlung, ohne Entscheid beginnt kein Code): Modulname
  und Kürzel (Name im Issue „offen"), Leseweise der Ausdrücke (zeilenweise,
  RE2-Dialekt, `\r` am Zeilenende), was „genau ein Treffer" zählt (Zeilen mit
  Treffer), ob ein Pin mit nur einer Quelle ein Konfigurationsfehler ist, die
  Pfad-Eingrenzung der Quelldateien (Eingaben außerhalb der Scan-Menge,
  [`AGENTS.md`](../../../AGENTS.md) §3.8) und die Grund-Codes.
- **Entschieden mit dem Auftraggeber (2026-10-31):** kein Schalter `strip-prefix`
  (der Ausdruck lässt das `v` vor der Fanggruppe stehen), nur Gleichheit, kein Netz.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Release, Tag, Release-Prep und die Antwort an das Issue** (CHANGELOG,
  READMEs, Handbuch-Kopf, Digest-Pins, Schließen von Issue 14): eigener Vorgang
  nach der Welle ([`AGENTS.md`](../../../AGENTS.md) §5 Zeile 17) und eine Handlung
  des Auftraggebers.
- **Aktualitätsprüfung gegen eine Netz-Quelle** (ist `1.27.3` die neueste
  Go-Version?): ein anderes Modul mit Netzzugriff; die Nachtlauf-Targets
  `freshness-*` tragen sie heute.
- **Ein Schalter `strip-prefix` oder jede andere Normalisierung der Werte:** der
  Ausdruck trägt sie, ein zweiter Mechanismus wäre eine zweite Stelle für dieselbe
  Entscheidung.
- **Schreibende Wege** (den Wert in der zweiten Datei angleichen): ein Vorschlag
  wäre Ausgabe, nie Schreiben ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Die übrigen Spiegel einer Pin-Hebung** (Release-URLs, Prosa-Pins,
  zitierende Verweise des Baseline-Pins): ein anderer Gegenstand als Wertgleichheit
  zwischen zwei Dateien; das Register führt ihn als
  [`BEO-ALL/pin-bump-mirrors-ungated`](observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md),
  und diese Welle löst ihn nicht auf.
- **Pins in Digest-Form** (`@sha256:…` im `Dockerfile`): der Ausdruck kann sie
  fangen, aber welche Digests zusammengehören, ist Fachwissen des Repos; das eigene
  Profil (slice-288) nimmt sie nicht auf.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: entsteht bei der Closure als `welle-94-results.md`, Geschwister im Ruheort `done/`.
Zähler: das Beobachtungs-Register, eine Ebene über dem Ruheort.
