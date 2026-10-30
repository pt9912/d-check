# Welle welle-93: Modul `metrics` — Kennzahl eines Planungsdokuments gegen Messung aus git

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-10-30.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Ein Repo, das in seinen Slice-Plänen eine Größenmarke als maschinenlesbaren
Block führt, prüft sie mit d-check statt mit einem eigenen Programm: das neue
opt-in-Modul `metrics` hält pro Plan mit Block (1) die Obergrenze gegen die
konfigurierte Formel über den Eingaben des Blocks, (2) die Messung aus git
(`git diff --numstat`-Summe über die Pfade des Blocks) gegen die Obergrenze und
(3) einen Ausnahme-Eintrag auf einen nicht leeren Grund (GitHub-Issue 13). Der
Lauf bleibt netzlos, schreibt nicht ins geprüfte Repo und braucht kein
git-Binary im Image — wie [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
und [`DC-FA-COMMITS-001`](../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Das GitHub-Issue 13 liegt offen vor (`gh issue list --state open`), und der
  Auftraggeber hat seine Bearbeitung als Welle angefordert (2026-10-30).
- `origin/main` und der lokale Hauptzweig sind gleich (Abstand 0/0 nach
  `git fetch`), die Lifecycle-Verzeichnisse `next/` und `in-progress/` tragen
  keinen Slice, `open/` nur die drei dieser Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle sind `done/`.
- `make fullbuild` grün.
- **Das Mehr gegenüber den Slice-DoDs** (repo-weit, in keiner einzelnen DoD):
  - **Ohne den Block nichts Neues:** `make blackbox-probe` gegen das Vorher-Image
    der letzten Release-Version zeigt für den Eigenlauf und ein Fixture ohne
    `metrics`-Konfiguration **dieselbe** Ausgabe (stdout, stderr und Exit
    getrennt). Ein neues opt-in-Modul darf keinen bisher grünen Lauf verschieben;
    jeder Slice belegt seine Hälfte, dass keine die andere verschoben hat, belegt
    erst der gemeinsame Lauf.
  - **Der Anlass trägt:** das Modul fällt über die Pläne des Anlass-Repos (die
    Pläne, deren Block das Issue beschreibt) **dasselbe Urteil je Plan** wie dessen
    bisheriges Programm — Plan für Plan verglichen, nicht als Summe. Die Probe fährt der
    **Eigentümer des Anlass-Repos** mit dem veröffentlichten Image und meldet das
    Ergebnis im Issue (Zusage im Kommentar vom 2026-10-30); das Repo ist für diese
    Planung nicht lesbar. Verglichen werden je Plan Exit-Code und Inhalt der Befunde und
    Hinweise (welche Prüfung, welche Zahlen), nicht der Wortlaut; es gibt keine
    Arbeitsteilung mit dem Programm, es entfällt nach der Abnahme ersatzlos. Weicht ein
    Urteil ab, ist die Abweichung benannt und entweder in der Anforderung
    verankert oder behoben. Dieser Lauf braucht die Pläne aus dem Anlass-Repo; sie
    stellt der Auftraggeber bereit (siehe §5). **Beobachtbar** wird er durch die
    Tabelle Plan → Urteil des Moduls → Urteil des Programms in
    `welle-93-results.md`, mit dem benannten Stand des Anlass-Repos; ohne diese
    Tabelle ist der Punkt nicht eingetreten.
  - **Die Betriebsbedingungen des Issues am gebauten Image:** der Lauf unter
    `:ro` und gegen ein Repo, das einem anderen Benutzer gehört, liefert dasselbe
    Ergebnis wie gegen ein eigenes beschreibbares Repo; ohne `.git` endet er mit
    Exit 2.
- Das Trigger-Audit (seit slice-241) ist über die Welle-Closure vollzogen und in
  `welle-93-results.md` belegt.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-282`](done/welle-93/slice-282-metrics-anforderung-und-adr.md) | Anforderung, Spezifikation und ADR für das Modul `metrics`; Bestands-Gegenprobe | neue Anforderung des Moduls (Kennung vergibt der Slice), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| [`slice-283`](done/welle-93/slice-283-metrics-modulkern-formel-und-grund.md) | Modulkern: Block lesen, Obergrenze gegen Formel, leerer Grund (hermetisch) | die Anforderung aus [`slice-282`](done/welle-93/slice-282-metrics-anforderung-und-adr.md) |
| [`slice-284`](done/welle-93/slice-284-metrics-messung-aus-git.md) | Messung aus git: Port, Adapter und Prüfung gegen die Obergrenze | die Anforderung aus [`slice-282`](done/welle-93/slice-282-metrics-anforderung-und-adr.md), [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) (Port-Vorbild) |
| [`slice-285`](open/slice-285-metrics-tabellenzaehlung-und-schaetzung.md) | Tabellenzählung: Schätzung aus der Tabelle des Plans, Prüfung 2, Hinweis „Marke ausstehend“ | die Anforderung aus [`slice-282`](done/welle-93/slice-282-metrics-anforderung-und-adr.md) (Entscheid 9) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle.
- Wird blockiert von: keiner.
- **Reihenfolge der Slices ist fachlich bindend:** erst
  [`slice-282`](done/welle-93/slice-282-metrics-anforderung-und-adr.md) (die Anforderung
  und die offenen Entscheide stehen, bevor Code entsteht), dann
  [`slice-283`](done/welle-93/slice-283-metrics-modulkern-formel-und-grund.md), dann
  [`slice-284`](done/welle-93/slice-284-metrics-messung-aus-git.md), dann [`slice-285`](open/slice-285-metrics-tabellenzaehlung-und-schaetzung.md) (Prüfung 2 braucht die Messung). Alle vier schreiben eine
  Fassung ins Lastenheft oder die Spezifikation; das WIP-Limit 1 hält sie
  ohnehin nacheinander.
- **Offene Entscheide des Auftraggebers, die [`slice-282`](done/welle-93/slice-282-metrics-anforderung-und-adr.md)
  vor dem Schreiben der Anforderung braucht** (das Issue legt sie nicht fest):
  neun, darunter die Syntax des Blocks im Plan, woher die Basis `<basis>` stammt, die
  Formel-Form, wo der Block gelesen wird (der Kern trägt keine YAML-Bibliothek) und
  wo die Zeilen gezählt werden, und wie die Tabellenzählung des Anlass-Repos im Modul
  konfiguriert wird. Der Slice benennt sie mit Empfehlung; ohne Entscheid
  beginnt kein Code. Fällt Entscheid 6 auf den Parse-Port, wird die Block-Lexik ein
  eigener Slice vor [`slice-283`](done/welle-93/slice-283-metrics-modulkern-formel-und-grund.md), und die Slice-Tabelle in §4 wächst.
- **Die Angaben und die Rückmeldung des Anlass-Repos:** das Repo ist nicht lesbar. Die
  Gegenprobe stützt sich auf die anonymisierten Angaben im Issue (Kommentar vom
  2026-10-30, erledigt); die Paritäts-Probe fährt der Eigentümer mit dem veröffentlichten
  Image und meldet je Plan im Issue. Ohne diese Rückmeldung ist die Probe nicht
  fahrbar, und die Welle schließt nicht mit einer Behauptung, sondern bleibt offen.
  Die Rückfrage zu den Slices B und C ist beantwortet (erwarteter Produktcode in der
  Schätzung); sie gehen als Testfälle ein. Die Tabellenzählung ist Teil der Anforderung
  (Kommentar vom 2026-10-30, 14:49).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Release, Tag, Release-Prep und die Antwort an das Issue** (CHANGELOG,
  READMEs, Handbuch-Kopf, Digest-Pins, Schließen von Issue 13): eigener Vorgang
  nach der Welle ([`AGENTS.md`](../../../AGENTS.md) §5 Zeile 17) und eine Handlung
  des Auftraggebers. Die Welle liefert Code, Spec und Handbuch-Inhalt, keine
  Distribution; der Closure-Trigger hängt deshalb nicht an einer Antwort, die vor
  dem Release falsch wäre.
- **Eine feste Regel statt einer konfigurierbaren Formel:** die Faktoren, Zustände
  und Schwellen sind Fachlogik des Anlass-Repos (Abgrenzung des Issues). Das Modul
  liefert den Mechanismus; welche Marke ein Repo setzt, steht in dessen
  `.d-check.yml`.
- **Der Abbau des Programms im Anlass-Repo:** ein Vorgang in einem anderen Repo,
  nicht am Gegenstand dieser Welle. Die Welle belegt nur, dass das Modul dasselbe
  Urteil fällt.
- **Andere Messgrößen als die `numstat`-Summe** (Dateizahl, Commit-Zahl,
  Zyklomatik): das Issue verlangt Einfügungen plus Löschungen. Eine weitere Größe
  wäre eine neue Anforderung mit eigener Gegenprobe, kein Schalter dieses Moduls.
- **Schreibende Wege** (den Block fortschreiben, die Obergrenze eintragen):
  [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  schließt sie aus; ein Vorschlag wäre Ausgabe, nie Schreiben.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: bei der Closure — Zeiger auf `welle-93-results.md`, Geschwister im Ruheort `done/`.
Zähler: bei der Closure — Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort.
