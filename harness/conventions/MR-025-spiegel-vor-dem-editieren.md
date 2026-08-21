# MR-025 — Semantik-Änderung: die Spiegel vor dem Editieren auflisten

- **Status:** Accepted
- **Ergänzt-Baseline-Regel:** [`modul-10-review-harness.md`](../../.harness/baseline/v5.7.0/regelwerk/modul-10-review-harness.md)
  (Review-Arten) — die Baseline kennt die Review-Arten, aber keine Vorab-Auflistung
  der Stellen, die eine Semantik spiegeln.
- **Datum:** 2026-08-15
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

  **Die Liste wird aus dem Repo abgeleitet, nicht aus dem Gedächtnis.** Bei einem
  neuen Modul, einem neuen Grund-Code oder einem neuen Target ist der billigste
  vollständige Ableiter ein `grep` nach dem **vorigen** Vertreter derselben Art
  (etwa dem zuletzt hinzugefügten Modulnamen) über den ganzen Baum. Bei der
  ersten Anwendung dieser Regel („slice-099“) fehlten vier Spiegel —
  `Makefile`, die Modul-Registry im Kern und zwei Anforderungs-Abschnitte —,
  und drei davon hätte genau dieses `grep` gefunden. Eine aus dem Kopf
  geschriebene Liste ist besser als keine, aber sie ist selbst ein Artefakt mit
  Lücken.

  **Der Spiegel ist die Stelle, nicht die Datei.** Eine bestätigende Re-Review
  desselben Slice fand eine **fünfte** Lücke, die der Datei-`grep` nicht gefunden
  hätte: die Modul-Enumeration in der emittierten Config-Vorlage lag in genau der
  Datei, die ohnehin bearbeitet wurde — sie stand deshalb nicht auf der Liste und
  blieb doch stehen. Wo eine Aufzählung eine Menge spiegelt, gehört sie an ihre
  Quelle **gebunden** (ein Test gegen die Registry), nicht in die Liste
  geschrieben. Die Liste ist der Notbehelf für alles, was sich nicht binden
  lässt.
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
