# Review slice-161 — Sechs gemeldete Pin-Rückstände heben

**Gegenstand:** [slice-161](../plan/planning/done/slice-161-sechs-pins-heben.md), Stand des Bump-Commits.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Schließbar nach Nacharbeit.** Die Substanz trägt — alle sechs Pins zeigen auf
genau das, was ihr Tag-Kommentar behauptet, alle zwölf Achsen melden `ok`, jedes
Gate ist grün. Die Nacharbeit ist **Text und Verfahren, kein Code**.

## Was der Reviewer selbst nachgemessen hat und was hält

- **Alle sechs Pins korrekt.** Die zwei Action-SHAs gegen `git ls-remote --tags`,
  die vier Image-Digests gegen die Registry. Die von
  [`AGENTS.md`](../../AGENTS.md) §3.9 benannte Grenze (SHA gegen Kommentar)
  verdeckt hier **keinen** Defekt.
- **Alle zwölf Achsen einzeln gefahren** — je Exit 0, je `ok`.
- **Die semgrep-Parität auf die Ziffer bestätigt:** beide Fassungen melden
  *„Ran 55 rules on 50 files: 0 findings"*.
- **Der Major-Sprung weiter geprüft als die Botschaft:** Release-Notes v7.0.0
  **und** v7.0.1, dazu der `action.yml`-Diff und der Quelltext des
  Sicherheits-Helfers. Die einzige Schnittstellenänderung ist ein neuer Input
  mit Default `false`; `node24` lief schon in v6.0.2. **Kein verschwiegener
  zweiter Bruch.**
- **Gates:** zehn Glieder plus `image-test`, je einzeln bestätigt, Coverage
  94,90 %.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | *„Der Nachtlauf war seit slice-142 dauerrot"* ist gemessen falsch: es gab **genau einen** roten Lauf, und er war beim Commit rund zwei Stunden alt — die drei Läufe davor waren grün. Der Slice-Plan formuliert es korrekt (*„beim ersten Lauf"*); die Botschaft machte daraus eine Dauerzustands-Behauptung |
| F-2 | MEDIUM | Die a-check-Paritäts-Aussage steht auf **einer** Probe, wo [ADR-0029](../plan/adr/0029-arch-check-via-a-check.md) eine **Matrix je Verbotszweig** verlangt und [`harness/README.md`](../../harness/README.md) sie als lebende Zusage trägt. Der Reviewer hat die Matrix gefahren — sieben Verbotszweige, drei Allow-Gegenproben, byte-gleiche Befundzeilen: die Aussage hält **inhaltlich**, aber das Repo hatte den Beleg nicht |
| F-3 | MEDIUM | Das `a-check.mk`-Fragment wurde nicht per `--print-mk` neu erzeugt, obwohl sein eigener Kopf das vorschreibt. Drei Neuerungen blieben unbeachtet: `DOCKER ?= docker`, ein Target `a-check-graph`, die erweiterte `.PHONY`-Zeile |
| F-4 | MEDIUM | Der Pin-Spiegel-Zensus zählt zwei lebende Fundstellen; es sind **drei** — `tools/harness/pin-freshness.sh` erklärte die v-Normalisierung mit demselben Beispielwert wie `harness/README.md`, eine Zeile tiefer nicht mitgezogen |
| F-5 | MEDIUM | Das zweite §5-Risiko ist **eingetreten** und blieb unbemerkt: der CA-Trust-Store des ausgelieferten Images hat sich geändert (216 591 → 224 449 Byte). Von libc ist dagegen nichts betroffen — `distroless/static` bringt keine mit |
| F-6 | LOW | *„Je Pin eine eigene Entscheidung"* ist für zwei von sechs nicht belegt: `docker/login-action` trägt keinen Satz, die zwei Digest-Nachzüge nur *„neuer Bau desselben Tags"* |
| F-7 | LOW | [ADR-0010](../plan/adr/0010-semgrep-hermetisches-gate.md) §Entscheidung nennt die gepinnte Fassung im **Indikativ** und ist seit dem Bump überholt. Der Weg ist die `## Geschichte`-Zeile, nicht eine neue ADR — die Botschaft warf *„eingefroren und weiter wahr"* und *„eingefroren und jetzt falsch"* in einen Topf |
| F-8 | LOW | Die dritte §5-Frage (Kadenz) ist unbeantwortet, blockiert die Closure aber **nicht**: der DoD-Zweig *„nicht mehr dauerrot"* ist erfüllt, und die Lücke steht bereits als benannte Grenze im Workflow-Kopf |

**Die Einordnung der übersprungenen Spiegel-Treffer ist sonst korrekt** — der
Reviewer hat jeden einzeln nachgesehen, einschließlich der Gegenprobe, dass die
zwei Mindest-Versions-Aussagen in `.a-check.yml` in v0.17.0 unverändert tragen.
§3 ist nicht verletzt.

## Erledigung

- **F-1** in die Closure-Notiz, weil der Commit bereits gepusht war — Amenden
  hätte veröffentlichte Historie umgeschrieben.
- **F-2** als Beleg in der Closure-Notiz, mit Attribution: die Matrix ist die
  Messung dieses Reviews, nicht des Implementers.
- **F-3** als **Abgrenzung** im Fragment-Kopf: keine der drei Neuerungen ist
  adoptiert, und beide Gründe stehen jetzt dort.
- **F-4** durch **Entfernen** statt Heben — das Beispiel war ein vierter
  Spiegel, der beim nächsten Bump wieder drifteten würde.
- **F-5** als Risiko-Ausgang `eingetreten` und als der eigentlich
  nutzersichtbare Teil des `CHANGELOG`-Eintrags (142 → 150 Wurzel-Zertifikate,
  eigens nachgemessen).
- **F-6** durch den fehlenden Satz und eine ehrlichere Formel.
- **F-7** als `## Geschichte`-Zeile in ADR-0010.
- **F-8** als Ausgang `weiter offen` mit
  [slice-164](../plan/planning/done/slice-164-nachtlauf-kadenz.md).
