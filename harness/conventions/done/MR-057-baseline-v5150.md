# MR-057 — Baseline-Pin-Hebung auf v5.15.0 (achter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-08-31
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../../AGENTS.md), [`harness/README.md`](../../README.md), den
  aktiven `MR-*`-Dateien, [`.harness/skills/reviewer.md`](../../../.harness/skills/reviewer.md),
  der Prüf-Konfiguration, den Spec-Straten und den Planning-Docs; dazu die vier
  Aliase unter `.claude/rules/`
  ([`MR-055`](../../conventions.md#mr-055))
- **Adaption:** Der Baseline-Pin ist von `v5.12.0` auf **`v5.15.0`** gehoben —
  die von [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, achter Nachtrag der Serie; ersetzt
  [`MR-037`](../../conventions.md#mr-037) nach dessen eigenem Auflösungs-Trigger.
  Kein Layout-Wechsel: dasselbe self-contained Bundle, dasselbe
  Materialisierungs-Skript, unverändertes Pfadschema.

  **Vier Tags, und der Inhalt je Tag ist gelesen statt angenommen** — die
  Herausgeber-Ankündigung nannte nur `v5.13.0`. Gemessen im Klon des
  Kurs-Repos, netzlos, `git log <tag>..<tag> -- kurs/de`:

  | Tag | Kurs-Wellen | Gegenstand |
  |---|---|---|
  | `v5.13.0` | 99 · 100 · 101 | der Spiegel zeigt nicht nach draußen · zwei Konsumenten lesen mit · **der Prüflauf verliert den Mount** |
  | `v5.13.1` | 102 | *kein* Delta in `kurs/de` — geändert ist die Konfigurations-Vorlage: „ein Beispiel, das nie traf" |
  | `v5.14.0` | 103 | zwei Rot-Quellen, ein Prinzip (Tiefe **und** Version eines Verweises) |
  | `v5.15.0` | 105 · 106 | das Bereichssegment bekommt einen Ort · **die Schwelle bekommt ihre drei Ausgänge** |

  **Das Bundle-Delta, gezählt statt geschätzt:** von **52** Dateien sind **20**
  unverändert und **32** geändert. Von den 32 tragen **16** ausschließlich den
  Versions-Stempel (≤ 2 Zeilen), **eine** ist das Manifest `SHA256SUMS`, **eine**
  ist `regelwerk/README.md` mit Stempel plus `**Stand:**`-Zeile. Bleiben
  **14** mit echtem Regel-Inhalt: zehn Regelwerk-Dateien und vier Vorlagen.

  **Die Lehre aus [`MR-037`](../../conventions.md#mr-037) hat sich diesmal
  ausgezahlt.** Dort war eigens vermerkt, dass die zwei Nicht-Markdown-Vorlagen
  des Bundles außerhalb der `*.md`-Delta-Schleife liegen und getrennt zu prüfen
  sind. Diesmal ist **`templates/.d-check.yml` um 31 Zeilen geändert** — eine
  Schleife über `*.md` hätte die Datei mit dem größten Adopter-Bezug übersehen.
  `templates/Makefile` ist unverändert.

  **Die Spiegel-Klassen aus [`BEO-008`](../../../docs/plan/planning/observations.md),
  je mit Zahl:** **129** lebende Vorkommen von `v5.12.0` vor der Hebung —
  **104** Pfad-Verweise, **5** Release-/Tree-URLs, **20** nackte Nennungen; dazu
  **130** eingefrorene. Gehoben sind **95**, stehen geblieben **34**: die drei
  Tombstone-Zeilen der Prüf-Konfiguration samt Kommentar, die zehn immutablen
  ADRs [`ADR-0069`](../../../docs/plan/adr/0069-zellenlaenge-als-strukturbedingung.md)–[`ADR-0078`](../../../docs/plan/adr/0078-erklaerte-leermenge-mit-zahl.md), das CR-Antwort-Dokument, der eigene Slice-Text,
  `BEO-017`, die Tabelle in
  [`MR-039`](../../conventions.md#mr-039), dieser Eintrag selbst und zwei Pläne in
  `open/`, deren Sätze Aussagen über ihren Planungszeitpunkt sind.
  **Die Messvorschrift gehört zur Zahl:** gezählt wurden **Vorkommen** (nicht
  Dateien) in den von `git ls-files` geführten
  `*.md`/`*.yml`/`*.sh`/`*.mk`/`Makefile`-Dateien, **vor** der Hebung, unter
  Ausschluss der vier eingefrorenen Präfixe (`docs/plan/planning/done/`,
  `docs/reviews/`, `harness/conventions/done/`, `.harness/baseline/`).

  **Die vier Aliase unter `.claude/rules/` hängen an keinem Modul.** Sie binden
  denselben Pin, werden von keinem Scan erfasst und wären still gebrochen;
  `make baseline-verify` hat alle vier gemeldet, bevor irgendetwas committet
  war — die Zusage aus [`MR-055`](../../conventions.md#mr-055), zum ersten Mal an
  einem echten Bump geprüft.

  **Ein Abschnitt ist umbenannt**, und das trifft drei lebende Einträge:
  `grundlagen-durchsetzungsschicht.md` §Referenz-Implementierung heißt seit
  `v5.13.0` **§Das vollständige Artefakt-Set**.
  [`MR-042`](../../conventions.md#mr-042), [`MR-043`](../../conventions.md#mr-043) und
  [`MR-048`](../../conventions.md#mr-048) tragen den Namen im **Linktext** — ohne
  Anker, deshalb sieht ihn kein Modul. Der Name ist nachgezogen; der zitierte
  **Wortlaut** bleibt unangetastet
  ([`MR-039`](../../conventions.md#mr-039)). Damit ist zugleich der `F-7`-Befund
  aus dem slice-155-Review geschlossen — er hatte denselben Namen in der
  Gegenrichtung gemeldet.

  **Zwei `cite`-Spannen sind neu geankert, keine entfernt.** `MR-035` von
  `grundlagen-begriffe.md:45-45` auf `:46-46` (eine Glossarzeile darüber
  eingefügt), `MR-043` von `grundlagen-durchsetzungsschicht.md:101-103` auf
  `:100-102` (eine Zeile darüber entfallen). Für jede übrige Direktive ist die
  **Position gegen den Datei-Diff** geprüft, nicht der Gate-Ausgang: `MR-031`
  und `MR-039` zitieren vor jeder Änderung ihrer Datei, `MR-005` ebenso; die
  restlichen zeigen in Dateien, deren einziges Delta der Versions-Stempel ist.
  Eine Spanne, die zufällig wieder wortgleich trifft, wäre kein Beleg dafür,
  dass sie dieselbe Regel zitiert.

  **Der Adaptions-Review ist durch alle 33 lebenden Einträge gelaufen:**
  **31 bleiben gültig** — die Deltas sind fast durchweg Ergänzungen
  (`modul-13` +21, `grundlagen-source-precedence` +5, `grundlagen-begriffe` +1
  ändern keine bestehende Zeile), und die wenigen geänderten Zeilen liegen in
  Abschnitten, die kein Eintrag als Basis nennt. Darunter fallen
  [`MR-042`](../../conventions.md#mr-042) und
  [`MR-043`](../../conventions.md#mr-043): ihr Linktext ist oben bereits als
  Nachzug des `F-7`-Befunds erklärt, eine **eigene** Antwort dieses
  Adaptions-Reviews ist das nicht. **Zwei tragen eine neue Antwort dieses
  Durchgangs:** [`MR-048`](../../conventions.md#mr-048) — sein veralteter
  Linktext war vor diesem Bump **nicht** bekannt (anders als bei MR-042/043
  über `F-7`; die Historie zeigt MR-048 seit seiner Anlage mit dem
  damals korrekten Namen, der erst durch die Umbenennung stale wurde) und ist
  damit ein eigenständiger Fund dieses Reviews, kein Nachzug eines bekannten
  Befunds — und [`MR-013`](../../conventions.md#mr-013) ist **teilweise
  überholt** — dazu unten.

  **Die Kollision, benannt statt aufgelöst.** Kurs-Welle 103 hat
  [`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln](../../../.harness/baseline/v5.18.0/regelwerk/grundlagen-traceability.md#herkunfts-anker-für-steering-loop-regeln)   <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/grundlagen-traceability.md:85-86 -->
  eine Regel für genau den Move gegeben, den
  [`MR-013`](../../conventions.md#mr-013) regelt:
  *„**Der `git mv` zieht die Pfad-Berichtigung nach sich**, als eigener Commit
  nach dem Umzug"* — die Fettung steht so in der Quelle; die Betonung liegt für
  diesen Eintrag dagegen auf *nach*, und das wird hier daneben gesagt statt im
  Zitat gesetzt.
  `MR-013` bündelt sie umgekehrt **in** den Move-Commit, weil die PR-/Push-CI
  den Push-Tip prüft und der ein Zwischen-Commit sein kann — ein roter
  Zwischenstand ist hier nicht folgenlos. Beide Begründungen stehen; die
  Baseline rangiert höher. Aufgelöst wird das **nicht in diesem Eintrag**:
  ein bestehender Eintrag wird nicht überschrieben, und die Antwort ist ein
  Nachfolge-Eintrag samt Nachzug in [`AGENTS.md`](../../../AGENTS.md) §3.3. Der
  Widerspruch ist damit gemeldet, nicht stillschweigend nach einer Seite
  entschieden.
- **Begründung:** Ein Adopter, der seine Baseline nicht auf einen Tag pinnt,
  auditiert gegen ein bewegliches Ziel; der Pin macht den Stand zitierbar und
  die Abweichung benennbar. Dass er **fortgeschrieben** wird statt zu altern,
  ist die Bedingung dafür, dass der Freshness-Audit etwas zu vergleichen hat.
  Diese Hebung hat zusätzlich einen Anlass, den die Serie bisher nicht kannte:
  sie ist die **Vorbedingung eines anderen Slice** —
  [slice-188](../../../docs/plan/planning/open/slice-188-register-gegen-neuen-kanon.md)
  kann die Beleg- und Ausgangs-Regeln erst zitieren, seit sie im vendorten Baum
  stehen.
- **Löst auf:** [`MR-037`](../../conventions.md#mr-037)
- **Ausgelöst durch Baseline-Stand:** v5.15.0
- **Auflösungs-Trigger:** der Kurs veröffentlicht einen neuen Release-Tag; dann
  Fortschreibung durch den nächsten Nachtrag zu
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt).
