# MR-065 — Baseline-Pin-Hebung auf v6.3.1 (elfter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-09-06
- **Geltungsbereich:** [§Baseline](../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../AGENTS.md), [`harness/README.md`](../README.md), den
  aktiven `MR-*`-Dateien, den Spec-Straten, den Planning-Docs und den beiden
  Skills; dazu die acht Baseline-Aliase unter `.claude/rules/`
  ([`MR-055`](../conventions.md#mr-055)).
- **Adaption:** Der Baseline-Pin ist von `v6.0.0` auf **`v6.3.1`** gehoben —
  die von [`MR-011`](../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, elfter Nachtrag der Serie; ersetzt
  [`MR-060`](done/MR-060-baseline-v600.md) nach dessen eigenem
  Auflösungs-Trigger. Kein Layout-Wechsel: dasselbe self-contained Bundle,
  dasselbe Materialisierungs-Skript, unverändertes Pfadschema.

  **Vier Tags, Inhalt je Tag gelesen statt angenommen** — die Release-Bodies
  taugen dafür nicht, sie sind für alle vier dieselbe Boilerplate über das ZIP;
  gelesen wurden die Bundles selbst:

  | Tag | Gegenstand |
  |---|---|
  | `v6.1.0` | `modul-07`: Carveout-Auflösung setzt die Bindung-Spalte in §Sensors zurück · `modul-10`: Review-**Deckung** ist mechanisierbar, während die Kategorisierung inferential bleibt — nennt d-checks Modul `reviews` als Werkzeug-Beispiel · `modul-13`: feuert der Auflösungs-Trigger, ist die Entfernung der Hard-Rule-Zeile ein **DoD-Punkt** des auslösenden Slice |
  | `v6.2.0` | `modul-05`: der Review-Report zählt nicht als Liefer-Punkt · `slice.template.md`: neuer DoD-Punkt *Review durchgeführt, Report liegt vor* (aus vier Closure-Pflichten werden fünf) · `templates/.d-check.yml`: auskommentierter `reviews`-Block |
  | `v6.3.0` | **Der große Strang:** `harness/sensors/<target>.md` — eine Datei je Gate, sobald sein Vertrag mehr als **einen Satz** braucht; die Tabellenzeile bleibt Index, die Target-Zelle wird zum **Link**. Trägt `grundlagen-harness-dateien` (+81), `grundlagen-begriffe` (Artefakt-Zeile), `README.template` (+47) und die **neue** Vorlage `templates/harness/sensors/gate.template.md`. Dazu `modul-13`: die **dritte Lage** — genannt, aber kein Gate; `kein Gate` gehört **in die Zeile**, nicht in Prosa daneben |
  | `v6.3.1` | `modul-13`: *Vorhanden ≠ behauptet* geschärft (ein reales Target **nicht** zu versprechen ist keine Lüge) und die Grenz-Regel — ein Gate ohne seine Grenze behauptet zu viel; eine Vollständigkeits-Zeile (*N Dateien geprüft*) nennt das **Kommando**, nicht die eingefrorene Zahl |


  **Grenze der Tabelle, benannt:** Vendored liegt nur der `v6.3.1`-Endstand.
  Der **Inhalt** jeder Zeile ist gegen den Aggregat-Diff belegt, die
  **Zuordnung zum einzelnen Tag** dagegen stammt aus den Zwischen-Bundles, die
  im Repo nicht liegen — sie ist aus dem Repo heraus nicht nachvollziehbar,
  sondern nur über einen erneuten Netz-Abruf der vier Release-Assets.

  **Das Bundle-Delta, gezählt statt geschätzt:** von **54** Dateien tragen
  **13** geänderten Inhalt, **eine** ist neu
  (`templates/harness/sensors/gate.template.md`), der Rest ist unverändert.

  **Die Messung selbst war die erste Falle.** Ein naiver Datei-Diff meldet
  **alle** Bundle-Dateien als geändert — jede trägt in Zeile 3 eine
  `<!-- Quelle: … blob/<tag>/… -->`-Herkunftszeile, die den Tag nennt. Wer diese
  Zahl berichtet, behauptet ein Delta von 53 statt der gemessenen 13. Der Fall
  ist als eingetretene Instanz von
  [`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../../docs/plan/planning/observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  notiert; das echte Delta entsteht erst mit `diff -I '<!-- Quelle:'`.

  **Adoptions-Entscheid je Strang — die Hebung ist nicht die Adoption.**
  Alle vier Stränge sind auf **Regelwerk-Ebene adoptiert** (der Kanon gilt ab
  diesem Pin), ihre **Umsetzung in der verkörperten Form** ist der Folge-Slice
  `slice-203`. Das ist keine Vertagung aus Bequemlichkeit, sondern eine
  Reihenfolge-Bedingung: die neuen Vorlagen liegen erst **nach** diesem Bump im
  Repo. Der `sensors/`-Strang trifft dabei unmittelbar d-checks eigene
  §Sensors-Tabelle, deren Zellen längst absatzlang sind — der Kanon nennt genau
  das den Fund, nicht die Ausnahme.

  **Der Zensus** ([`MR-021`](../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden)):
  **51** lebende Dateien gehoben — `AGENTS.md`, `harness/README.md`,
  `harness/conventions.md`, **32** aktive `MR-*`-Einträge, beide Skills, beide
  berührten Spec-Straten, vier Planning-Docs und die acht Aliase. Beide von
  [`MR-060`](done/MR-060-baseline-v600.md) benannten
  Übersetzungsfehler-Klassen sind erneut aufgetreten und wurden gefangen: der
  `../baseline/` <!-- d-check:ignore (Muster-Praefix, benennt die Fehlerklasse) -->-relative Pfad in `.harness/skills/` und der Release-Link ohne
  `.harness/baseline/`-Segment. **Eine dritte Klasse kam hinzu** — der
  gehobene **Link-Text** neben dem gehobenen Link-**Ziel**: `**Stand:**
  [`v6.0.0`](…/releases/tag/v6.3.1)` ist ein in sich widersprüchlicher
  Zwischenstand, den kein Sensor meldet.

  **Was ausdrücklich NICHT gehoben wurde**, obwohl es `v6.0.0` nennt: die
  historischen Aussagen. Zwei Kommentare in
  [`config.go`](../../internal/hexagon/core/model/config.go) (*„vor/seit
  `v6.0.0`"* über den Wechsel der Register-Form), die
  Verzeichnis-Form-Notiz in
  [`observations/README.md`](../../docs/plan/planning/observations/README.md),
  der Tombstone-Kommentar der **vorigen** Migration in `.d-check.yml` und die
  Titel-/Geltungsbereich-Angaben von
  [`MR-060`](done/MR-060-baseline-v600.md) selbst. Sie bleiben wahr; ein
  mechanischer Rewrite hätte sie still verfälscht
  ([`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../../docs/plan/planning/observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)).
  Die Gegenprobe gilt auch rückwärts: **zwei** Vorkommen waren trotz
  Versions-Nennung **Live-Zeiger** und wurden gehoben — die
  Kanon-Sektions-Annotation in der Roadmap und das Pin-Beispiel in
  [`MR-021`](../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden).

  **cite-Direktiven** ([`MR-051`](../conventions.md#mr-051)): **16** lebende
  Direktiven in 11 Dateien. **3 neu verankert** — `grundlagen-begriffe.md`
  `:46`→`:47`, `grundlagen-harness-dateien.md` `:187`→`:190`,
  `modul-10-review-harness.md` `:75`→`:82` —, **0 entfallen**, **13**
  unverändert bestätigt. Alle drei Verschiebungen liegen in Dateien mit
  echtem Delta und wurden vom Modul `citations` **gemeldet**, nicht
  vorausberechnet; die Neuankerung folgte dem Zieltext.
  **`.d-check.yml`-Tombstone: EIN Glob-Eintrag, nicht fünf** — und der Grund
  ist enger, als die Entscheidung aussehen lässt. Gemessen trägt nur
  [ADR-0083](../../docs/plan/adr/0083-beobachtungsregister-verzeichnis-modus.md)
  einen **Markdown-Link** auf den entfernten `v6.0.0`-Baum, und nur solche
  liest das Modul `links`. Die übrigen eingefrorenen Klassen nennen den Tag
  sehr wohl **mit auflösendem Pfad**: `done/slice-200` und `done/slice-201`
  tragen je zwei `d-check:cite`-Direktiven dorthin, `conventions/done/MR-057`
  und `MR-058` je eine. Sie bleiben still, weil `citations.scope` diese
  Verzeichnisse ausnimmt und `links` keine HTML-Kommentare liest — **nicht**,
  weil dort kein Pfad stünde. Ein Glob mehr deckte dort also nichts, was nicht
  ohnehin ausgenommen ist; vier auf Vorrat wären unbelegte Ausnahmen.

  **Der Adaptions-Review ist durch alle 38 aktiven Einträge gelaufen**
  (37 vorbestehende plus dieser): **37 bleiben unverändert gültig, einer trägt
  eine überholte Feld-Aussage weiter.** Kein Abschnittsname ist umbenannt,
  keine bestehende Regel entfällt oder widerspricht dem Delta — alle vier
  Stränge sind additiv. Die Ausnahme ist
  [`MR-051`](../conventions.md#mr-051): sein `Geltungsbereich` sagt, in den
  eingefrorenen Verzeichnissen stehe **keine** `cite`-Direktive — die vier in
  `done/slice-200`/`slice-201` widerlegen das. Die **Regel** trägt
  unverändert; überholt ist nur die Feld-Aussage, und zwar bereits durch
  [`MR-054`](../conventions.md#mr-054), was die Index-Zeile in
  [`conventions.md`](../conventions.md#mr-051) schon ausweist. Der Eintrag
  bleibt deshalb unangetastet — ein `MR` wird nicht nachträglich
  umgeschrieben —, aber „gültig" heißt hier *die Regel gilt*, nicht *jedes
  Feld stimmt noch*. Eigens geprüft, weil ihr Auflösungs-Trigger genau darauf
  zeigt:
  [`MR-035`](../conventions.md#mr-035) und
  [`MR-036`](../conventions.md#mr-036) bleiben **aktiv** — der Kanon benennt
  auch in `v6.3.1` keinen Ruheort für den ausgehenden Change Request;
  `grundlagen-begriffe.md` hat seine Artefakt-Tabelle erweitert, aber um
  `harness/sensors/` <!-- d-check:ignore (Kanon-Form aus v6.3.1, im Repo erst mit slice-203) -->, nicht um eine CR-Ablage.

- **Begründung:** Ein Adopter, der seine Baseline nicht auf einen Tag pinnt,
  auditiert gegen ein bewegliches Ziel; der Pin macht den Stand zitierbar
  und die Abweichung benennbar. Dass er **fortgeschrieben** wird statt zu
  altern, ist die Bedingung dafür, dass der Freshness-Audit etwas zu
  vergleichen hat.
- **Löst auf:** [`MR-060`](done/MR-060-baseline-v600.md)
- **Ausgelöst durch Baseline-Stand:** v6.3.1
- **Auflösungs-Trigger:** der Kurs veröffentlicht einen neuen Release-Tag;
  dann Fortschreibung durch den nächsten Nachtrag zu
  [`MR-011`](../conventions.md#mr-011--baseline-auf-release-tag-gepinnt).
