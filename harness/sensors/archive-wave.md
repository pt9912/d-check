# `make archive-wave` — bewegt geschlossene Zeitdokumente ins Archiv

## Vertrag

Setzt Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur
Schritt 4 um. **Kein Gate** — ein bewusster, von Hand ausgelöster Vorgang, der
den Repo-Zustand *verändert*, statt über ihn zu urteilen.

**Drei Modi, wechselseitig ausschließend** — genau einer von `WELLE`, `SLICE`,
`REVIEW`; `APPLY=1` ist bei allen dreien optional.

| Modus | Was er einsammelt | Wohin |
| --- | --- | --- |
| `WELLE=<id>` | die Slices einer geschlossenen Welle (über ihr `**Welle:**`-Feld) und deren Review-Reports | `done/<welle-id>/archiv.zip` |
| `SLICE=<id>` | einen einzelnen **wellenlosen** `done/`-Slice | `docs/plan/planning/done/wellenlos/` |
| `REVIEW=<dateiname>` | einen **eigenständigen** Review-Report ohne Slice-Partner | `docs/reviews/archiv/` <!-- d-check:ignore (existiert erst nach mindestens einer Anwendung, seit slice-199) --> |

**Was in allen Modi passiert:** Der Volltext wird durch einen **Stub** ersetzt,
und repo-weite Verweise werden nachgezogen.

**Wo die Modi sich unterscheiden, und warum:**

- **Review-Reports einer Welle oder eines Slice bekommen keinen Stub** — ihre
  Identität kommt vom Slice bzw. von der Welle, nicht von ihnen selbst.
- **Ein eigenständiger Review-Report bekommt einen** (seit slice-198), weil er
  selbst der abgeschlossene Vorgang ist. Trägt sein Dateiname eine
  `slice-<NNN>`-Kennung, wird er abgelehnt — er gehört in den `SLICE`-Modus.
- **Ein Slice mit echter Wellen-Zugehörigkeit wird im `SLICE`-Modus
  abgelehnt** (seit slice-196; Regelwerk §Wann Arbeit eine Welle braucht,
  *„ohne Wellen tut es die Slice-Closure selbst"*). Alle wellenlosen Archive
  teilen sich **ein** Verzeichnis statt eines Unterverzeichnisses je Slice
  (seit welle-89).

**Sicherer Default:** Ohne `APPLY=1` wird **nichts** geschrieben — nur der
geplante Umfang angezeigt.

**Eigenständiges Werkzeug** unter `tools/archive-wave/` mit eigenem `go.mod`
und eigenem `Dockerfile`: portabel für jedes Repo mit demselben
Planning-Layout, **kein** Import aus d-checks internen Paketen. Seine
Testsuite ist deshalb `make archive-wave-test` und **nicht** Teil von
`make test`.

## Grenze — was der Lauf nicht leistet

1. **Tote externe Verweise werden im `SLICE`-Modus gemeldet, nicht behoben.**
   Für Verweise auf die zu löschenden Review-Reports gibt es kein Move-Ziel;
   das Werkzeug nennt sie und lässt sie stehen. Permanent — wohin sie zeigen
   sollen, ist ein Urteil.
2. **Der Trocken-Lauf zeigt den Umfang, nicht das Ergebnis.** Was ohne
   `APPLY=1` angezeigt wird, ist die geplante Menge; ob der Stub danach die
   richtige Form hat, sagt erst `make doc-check` auf dem geschriebenen Stand.
3. **Die Vollständigkeit des Archivs bezeugt nur der Commit.** Ob wirklich
   alles eingesammelt wurde, was zur Welle gehörte, prüft kein Sensor — das
   ist die benannte Grenze des Kanons (Regelwerk `modul-06`, Schritt 4) und
   der Grund, warum die Operation in einem Werkzeug steckt und nicht in
   Handarbeit.

## Bindung

**Kein Gate** — weder in `gates` noch in `ci`. Die Commit-Granularität der
drei Modi ist geregelt: [`MR-059`](../conventions/MR-059-wellen-archiv-stub-move.md)
(Welle) · [`MR-062`](../conventions/MR-062-wellenloser-slice-archiv-move.md)
(Einzel-Slice) · [`MR-063`](../conventions/MR-063-eigenstaendiger-review-archiv-move.md)
(eigenständiger Review) · [`MR-064`](../conventions/MR-064-buendelung-slice-archiv-move.md)
(Bündelung mehrerer Einzel-Slice-Moves).
