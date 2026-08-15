# MR-025 — Semantik-Änderung: die Spiegel vor dem Editieren auflisten

- **Status:** Accepted
- **Ergänzt-Baseline-Regel:** [`modul-10-review-harness.md`](../../.harness/baseline/v5.0.0/regelwerk/modul-10-review-harness.md)
  (Review-Arten) — die Baseline kennt die Review-Arten, aber keine Vorab-Auflistung
  der Stellen, die eine Semantik spiegeln.
- **Datum:** 2026-08-10
- **Geltungsbereich:** jede Änderung an einer zugesagten Semantik — Grund-Code,
  Algorithmus-Schritt, Config-Schlüssel, Schwellenwert, Erkennungs-Form.
- **Adaption:** **Vor** dem Editieren wird die Liste der Spiegel dieser Semantik
  aufgeschrieben und im Slice oder Commit sichtbar gemacht. Die Spiegel dieses
  Repos sind:

  | Spiegel | Ort |
  |---|---|
  | Anforderung | `spec/lastenheft.md` (Beschreibung **und** Akzeptanzkriterien) |
  | Algorithmus | `spec/spezifikation.md` (Schritt-Text) |
  | Config-Schema | `spec/spezifikation.md` §2 |
  | Grund-Code-Tabelle | `spec/spezifikation.md` §4 |
  | Klartexte | `AllReasons()` und `reasonTexts()` (`--doctor`) |
  | Emittierte Vorlage | `--print-config` / `--suggest-config` |
  | Nutzer-Doku | `docs/user/benutzerhandbuch.md` (Aufgabe · Modul-Tabelle · §5 · §11) |
  | Autoritäts-Doku | `AGENTS.md`, `harness/README.md` (Gate-Beschreibungen) |
  | Entscheidung | die zugehörige ADR (Entscheidung **und** Fitness Function) |

  Nicht jede Änderung berührt jeden Spiegel — aber **welche** sie berührt, wird
  entschieden, bevor der erste Editor aufgeht, nicht danach.
- **Begründung:** Die Klasse ist **dreimal** eingetreten
  ([`BEO-002`](../../docs/plan/planning/observations.md)) und **kein einziges Mal
  von einem Gate** gefunden worden — jedes Mal von einem Menschen oder einem
  unabhängigen Review. Sie ist gate-unsichtbar, weil die Spiegel **Prosa** sind:
  `links`, `anchors` und `ids` prüfen, ob eine Referenz auflöst, nicht ob eine
  Aussage noch gilt. Ein Rand, der eine widerrufene Fassung referiert, wirkt dabei
  so belastbar wie der Körper.

  Der Zähler weiterzuführen hätte nichts geändert: die Klasse wird nicht dadurch
  seltener, dass man sie öfter notiert. Die Auflistung **vor** dem Editieren ist
  die kleinste Form, die wirkt — sie kostet eine Minute und macht das Vergessen
  zu einer sichtbaren Auslassung statt zu einem stillen.
- **Auflösungs-Trigger:** wenn ein Gate die Spiegel-Konsistenz prüft (etwa ein
  Modul, das Aussagen statt Referenzen vergleicht), ist diese Regel überflüssig
  und wandert nach `done/`.
