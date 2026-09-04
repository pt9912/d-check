# Review slice-167 — Der Kanal erreichte die Fundklasse nicht, für die er gebaut ist

**Gegenstand:** [slice-167](../plan/planning/done/slice-167-dependabot.md), Stand des Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-28. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Schließbar nach Nacharbeit.** Zwei HIGH, vier MEDIUM, drei LOW. Die tragende
Messung — die Traceability-Kollision — hält exakt wie behauptet und ist vom
Reviewer reproduziert. Die **Begründungen um sie herum** halten an zwei Stellen
nicht: die `docker`-Ausschluss-Begründung ist gegen `dependabot-core`
nachweislich falsch, und der Kanal erreicht in der ausgelieferten Konfiguration
die Fundklasse nicht, für die er gebaut wurde.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | Die `docker`-Ausschluss-Begründung ist **gegen die Quelle falsch — beide Hälften**. Der Docker-Updater hebt Tag **und** Digest gemeinsam (`resolved_digest_for`), und er stellt die Frage *„anderer Bau desselben Tags"* sehr wohl (`digest_requirement_up_to_date?`) — für einen nicht-vergleichbaren Tag wie `nonroot` sogar ohne Unterdrückung. Dritte, ungenannte Tatsache: `FROM golang:${GO_VERSION}@sha256:…` parst Dependabot ohnehin nicht |
| F-2 | **HIGH** | **Der Kanal erreicht die vierzehn Befunde nicht — und das ist messbar, nicht „nicht gemessen"**: 13 der 14 lagen `// indirect`; Version-Updates filtern indirekte Deps ohne `allow` heraus; und beide Repo-Schalter sind aus (`automated-security-fixes` `enabled:false`, `vulnerability-alerts` HTTP 404) |
| F-3 | MEDIUM | Die ADR-0067-Zeile steht bei `README.md:107` — **hinter** `## Konventionen`, nicht in der Tabelle. Die Datei verlangt an Z. 87 selbst *„am Ende der Tabelle"*. `doc-check` grün, weil der Link auflöst |
| F-4 | MEDIUM | Die konfigurierten Labels `dependencies`/`ci` **existieren im Repo nicht**; gesetzte Labels **ersetzen** die Defaults, unbekannte werden still ignoriert ⇒ PRs ganz ohne Label. Die Bauform, vor der das Schwester-Repo ausdrücklich warnt |
| F-5 | MEDIUM | *„sechs `uses:`-Einträge"* — der eigene Sensor sagt **acht** (sieben extern, **drei** verschiedene Actions). Die Zahl trifft keine der drei Größen |
| F-6 | MEDIUM | §3.6 wird als **Verbot** zitiert, wo es eine **Verfahrenspflicht** ist: es verlangt eine ADR für eine Lockerung, verbietet sie nicht. Alternative B wäre dokumentiert zulässig gewesen (`BEO-012`) |
| F-7 | LOW | *„die vierzehn lagen in indirekten Abhängigkeiten"* — **dreizehn** von vierzehn; `go-git` steht im direkten `require`-Block |
| F-8 | LOW | Folgepflicht über `ignore`-Einträge, die es nicht gibt — die Konfiguration enthält **null** `ignore`-Schlüssel; Risiko 2 des Slice-Plans ist gegenstandslos |
| F-9 | LOW | *„bei einem einzigen direkten Dependency"* — es sind **zwei** |
| F-10 | LOW | Keine Doku-Oberfläche außerhalb `AGENTS.md` §5 und der ADR; auch das Repo-Setting, an dem die halbe Zusage hängt, ist nirgends genannt |

**Negativbefunde — was der Reviewer zu brechen versuchte und nicht brechen
konnte.** Die tragende Messung ist **reproduziert**: ohne Kennung
`commit-untraceable`/Exit 2, mit Kennung Exit 0 — auch bei einer **gruppierten,
mehrzeiligen** Dependabot-Botschaft und bei einem Dependabot-**Merge**-Commit
(via `exempt-pattern`). `trace-check` prüft die **ganze** Message, nicht nur den
Betreff — die Präfix-Lösung ist damit hinreichend, nicht notwendig.
**Der Präfix-Mechanismus ist belegt, nicht behauptet:** `pr_name_prefixer.rb`
führt die **schließende eckige Klammer** ausdrücklich in seiner Zeichenklasse,
das SchemaStore-Limit von 50 Zeichen ist mit 22/20 eingehalten, und die
GitHub-Doku nennt *„closing bracket"*. `make workflow-pins` übersteht einen
Dependabot-Action-Bump — empirisch am Schwester-Repo: SHA **und** Tag-Kommentar
werden gemeinsam ersetzt. `commits.exempt-pattern` ist unverändert; der
`d-check:ignore`-Marker in `AGENTS.md` ist per Fixture-Gegenprobe **lebendig**
und verdeckt keine echte Linkpflicht. ADR-0067 hält vollständig gegen die
Baseline-Vorlage. Die Schwester-Repo-Zuschreibungen stimmen, ebenso *„fünf
Digest-Achsen"*, *„drei Frische-Achsen"* und die 9/4/1-Aufteilung.

## Erledigung

Alle neun Befunde sind eingearbeitet; die zwei HIGH sind **eigens nachgemessen**
statt übernommen.

- **F-2** `allow: dependency-type: all` ergänzt — für `gomod` dokumentiert
  unterstützt, nachgeschlagen statt angenommen. Die zweite Hälfte ist ein
  **Repo-Setting** und steht jetzt als Vorbedingung in
  [`releasing.md`](../user/releasing.md); die ADR nimmt die *„nicht
  gemessen"*-Aussage per `## Geschichte` zurück.
- **F-1** Der Ausschluss steht jetzt auf
  [ADR-0011](../plan/adr/0011-digest-pins-build-gate-images.md) §Entscheidung 4
  — einer **Policy** statt eines technischen Bruchs, den es nicht gibt. Beide
  widerlegten Sätze sind per `## Geschichte` zurückgenommen.
- **F-3** Zeile in die Tabelle bewegt — **und der Fehler zum zweiten Mal
  gemacht** (ADR-0066 ebenso). Deshalb mechanisiert: eine `structure`-Regel
  verbietet ADR-Tabellenzeilen unter `## Konventionen`, in beide Richtungen
  gegengeprobt. Die **erste** Fassung der Regel feuerte nicht — der `^`-Anker
  zielte ins Leere, weil `forbid-pattern` den Abschnittstext als Ganzes prüft;
  das steht als zweite Grenze im Kommentar.
- **F-4** `labels:` gestrichen, damit Dependabot sein `dependencies` selbst
  anlegt.
- **F-5**, **F-7**, **F-9** Zahlen auf die gemessenen Größen gesetzt.
- **F-6** Der §3.6-Verweis ist zurückgenommen; das Sachargument trägt allein.
- **F-8** Folgepflicht und Risiko auf die **Ausschlüsse** gezogen, die es
  wirklich gibt.
- **F-10** `harness/README.md` führt den Kanal — samt der Hälfte, die keine
  Datei ist.
