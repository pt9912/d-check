# MR-052 — Ein historisches Baseline-Zitat ist prüfbar: der alte Baum liegt in der git-Historie (schärft MR-039)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon regelt, dass ein bestehender
  Eintrag nicht rückwirkend umgeschrieben wird
  ([`modul-02-harness-bootstrap.md` §Freshness-Audit](../../.harness/baseline/v5.18.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)).
  Ob das damalige Zitat **wörtlich** war, ist eine andere Frage, und der Kanon
  stellt sie nicht.
- **Datum:** 2026-08-27
- **Geltungsbereich:** wörtliche Zitate der Baseline in **allen lebenden**
  Dokumenten dieses Repos, die einen **früheren** Pin zitieren — also genau die
  Menge, für die [`MR-039`](../conventions.md#mr-039) den Wortlaut einfriert.
- **Adaption:** [`MR-039`](../conventions.md#mr-039) friert den zitierten
  Wortlaut ein, *„auch wenn er nicht mehr gilt"*, und begründet das mit der
  historisch korrekten Aussage über den damaligen Zustand — ein Grund, der von
  der **Zugänglichkeit** der Quelle unabhängig ist. Der Eintrag selbst nimmt
  also nichts an. Die **Praxis** tat es: [slice-152](../../docs/plan/planning/done/slice-152-citations-scharfschalten.md)
  legte drei Abweichungen mit *„am Repo nicht entscheidbar"* ab, weil der alte
  Pin nicht mehr vendored sei. Das ist strukturell falsch: der Bump entfernt
  `.harness/baseline/<alt-tag>/` aus dem **Arbeitsbaum**, nicht aus dem
  **Repo**. Der Lösch-Commit hat einen Vorgänger, und in dem steht der alte
  Baum vollständig.

  **Folge, und sie ist eine Verschärfung, keine Lockerung:** ein historisches
  Zitat ist **prüfbar**. Wer behauptet, es sei am Repo nicht entscheidbar, ob
  ein Zitat damals ungenau transkribiert oder seither gedriftet ist, hat den
  falschen Baum gelesen. Der Wortlaut bleibt trotzdem stehen — daran ändert
  diese Schärfung nichts; sie ändert nur, dass die **Feststellung** möglich ist
  und deshalb auch verlangt werden darf.

  **Bestand bei Einführung, gemessen:** die Delta-Tabelle in
  [`MR-039`](../conventions.md#mr-039) und der Eintrag
  [`MR-033`](../conventions.md#mr-033) geben **dasselbe** `v5.11.0`-Zitat
  verschieden wieder:

  | Stelle | Wiedergabe |
  |---|---|
  | Quelle (`v5.11.0`, aus der git-Historie) | *„…`spec/architecture.md` referenziert Modul-Pfade, aber **keine** Wellen…"* |
  | [`MR-039`](../conventions.md#mr-039) Delta-Tabelle | **wörtlich** |
  | [`MR-033`](../conventions.md#mr-033) | Hervorhebung **verschoben** — sie fettet den ersten Halbsatz und lässt die des zweiten weg |

  Beide bleiben, wie sie sind: der eine ist wörtlich, der andere ist ein
  eingefrorener Wortlaut nach
  [`MR-039`](../conventions.md#mr-039). Festgehalten ist **welcher** — das war
  vorher nicht entscheidbar, weil niemand die Quelle gesucht hat.
- **Begründung:** Eine Praxis, die auf „nicht mehr feststellbar" ausweicht,
  verdeckt genau dann etwas, wenn die Feststellung doch möglich ist. Der Preis
  der Schärfung ist gering — ein `git show <lösch-commit>^:<alter-pfad>` —, der
  Ertrag ist, dass „nicht entscheidbar" keine zulässige Antwort mehr ist, wo sie
  bloß bequem wäre.
- **Auflösungs-Trigger:** entfällt, sobald ein Bump den alten Baum nicht mehr
  aus dem Arbeitsbaum entfernt — dann ist die Quelle ohnehin da. Oder sobald
  die git-Historie beschnitten wird; dann greift wieder die Annahme von
  [`MR-039`](../conventions.md#mr-039).
