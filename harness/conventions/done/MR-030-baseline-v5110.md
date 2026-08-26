# MR-030 — Baseline-Pin-Hebung auf v5.11.0 (sechster Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-08-23
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../../AGENTS.md), [`harness/README.md`](../../README.md), den
  aktiven `MR-*`-Dateien, [`.harness/skills/reviewer.md`](../../../.harness/skills/reviewer.md),
  den Spec-Straten und den Planning-Docs
- **Adaption:** Der Baseline-Pin ist von `v5.9.0` auf **`v5.11.0`** gehoben
  (Kurs-Tag vom 2026-08-23) — die von
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, sechster Nachtrag der Serie; ersetzt
  [`MR-029`](MR-029-baseline-v590.md) nach dessen eigenem
  Auflösungs-Trigger. Kein Layout-Wechsel: dasselbe self-contained Bundle
  (`lab-regelwerk.zip`, beide Bäume, `SHA256SUMS`, 51 Dateien — gemessen),
  dasselbe Materialisierungs-Skript, unverändertes Pfadschema
  `.harness/baseline/<tag>/{regelwerk,templates}/`.

  **Zwei Minors auf einmal, und das Delta ist gemessen statt gezählt:** von 52
  Bundle-Dateien unterscheiden sich **30**; **23 davon ändern ausschließlich
  den Versions-Stempel** (Quell-URL bzw. Pfad), eine ist das Manifest — und von
  den verbleibenden sechs trägt **eine im Rumpf ebenfalls nur Versions-Zeiger**
  (`regelwerk/README.md`: zwei URL-Stempel plus die `**Stand:**`-Zeile).
  **Fünf Dateien tragen also echten Regel-Inhalt**; dieselbe Subtraktion hat
  [`MR-029`](MR-029-baseline-v590.md) für ihre Hebung vorgenommen. Die
  sechs mit Umfang, die Stempel-Datei ausgewiesen:

  | Datei | Umfang |
  |---|---|
  | `regelwerk/grundlagen-source-precedence.md` | +75/−4 |
  | `regelwerk/grundlagen-referenz-richtung.md` | +30/−1 |
  | `regelwerk/grundlagen-durchsetzungsschicht.md` | +8/−1 |
  | `regelwerk/grundlagen-begriffe.md` | +7/−2 |
  | `templates/spec/lastenheft.template.md` | +22/−7 |
  | `regelwerk/README.md` | +3/−3 — **kein Regel-Inhalt**, nur Stempel |

  **Der größte Block ist die Antwort auf einen Konsumenten-CR dieses Repos.**
  Kurs-Welle 94 („Eine Rangliste ordnet, jetzt deckt sie auch ab") nennt ihn im
  Kurs-CHANGELOG als Auslöser; sie bringt die **Vollständigkeits-Zusage** in
  `grundlagen-source-precedence.md` und die Rolle der Werkzeug-Einstiegsdatei
  in `grundlagen-durchsetzungsschicht.md`. **Diese MR hebt den Pin, sie
  behauptet keine Konformität** — was der neue Stand von diesem Repo verlangt,
  beantwortet das Delta-Audit der Etappe B, und die erste bekannte Verletzung
  hat bereits ihren eigenen Slice.

  **Der Drei-Klassen-Zensus, fortgeschrieben für den Nachfolger.** MR-029 hat
  ihn als Checkliste hinterlassen; er wandert mit ihr nach `done/`, das nicht
  jeder Lauf liest, und steht deshalb hier erneut — mit dem, was diese Hebung
  ergänzt hat:

  | Klasse | Form | Fallstrick dieser Hebung |
  |---|---|---|
  | 1 — Pfade | `.harness/baseline/<tag>/…` **und relativ** (`../baseline/<tag>/…`) | ein `grep` auf die absolute Form übersieht die relative (`.harness/skills/reviewer.md`) |
  | 2 — URLs | `releases/tag/`, `releases/download/`, `tree/` | Link-**Text** und Link-**Ziel** driften auseinander; eine fünfte Fundstelle (`harness/README.md`) lag außerhalb der zuerst gewählten Datei-Liste |
  | 3 — Prosa/Ellipsen | `…/<tag>/…`, „Stand"-Angaben, Regel-**Fassungs**-Zitate | ein Fassungs-Zitat ist hebbar, **sobald** die zitierte Regel als unverändert gemessen ist — „nicht geprüft" ist keine Begründung, wenn die Messung im selben Bogen vorliegt |

  **Historische Verweise bleiben stehen.** `done/`-Slices, eingefrorene
  Review-Reports, ADRs und aufgelöste `MR-*`-Dateien tragen weiterhin
  `baseline/v5.0.0/`, `v5.6.0/`, `v5.7.0/`, `v5.9.0/` — sie sprechen über die
  **Vergangenheit**, und ein Pfad-`grep` kennt keine Zeitform. Gehoben werden
  ausschließlich Verweise, die den **gegenwärtig** gültigen Stand meinen. Das
  ist die zweite Richtung von `BEO-008` und der Grund, warum die Hebung eine
  Gegenprobe braucht und kein `sed` über den ganzen Baum ist.
- **Begründung:** Auftraggeber-Anstoß 2026-08-23; der Baseline-Default sticht
  die repo-lokale Adaption. Vendored wird das **Release-Asset am Tag**
  (`--check-latest` ist die Currency-/Authentizitäts-Gegenprobe), nicht der
  Kurs-Arbeitsbaum — der lag zum Zeitpunkt der Hebung bereits vor dem Tag.
- **Löst auf:** [`MR-029`](../../conventions.md#mr-029) *(der Verweis geht auf die
  Index-Zeile, nicht auf die Eintrags-Datei — die wandert bei Auflösung nach
  `done/`, und ein Pfad-Link bräche genau dann)*
- **Ausgelöst durch Baseline-Stand:** v5.11.0
- **Auflösungs-Trigger:** die nächste Pin-Hebung ersetzt diesen Eintrag durch
  ihren Nachfolger — wie [`MR-029`](MR-029-baseline-v590.md) durch diesen
  Eintrag ersetzt wurde.
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  bleibt daneben **aktiv** stehen: es trägt das Bundle-Layout, nicht den
  Pin-Wert.
