# Welle 70 — Unbalancierter Fence als stiller Grün-Pfad — Closure-Notiz

**Welle:** welle-70-fence-lexik
**Abschluss:** 2026-08-10
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Der Grund-Code `fence-unclosed`** als dritte Artefakt-Klasse des Moduls
  `spans` ([`DC-FA-SPAN-001`](../../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in),
  Lastenheft 0.52.3, [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md)
  `Proposed`), ausgeliefert mit **v0.53.0** — Digest
  `sha256:0cbe2d54…9424`, Pipeline-Lauf 31362555433.
- **Die Fence-Lexik als geteiltes Prädikat.** `TrimFenceIndent`, `FenceRun` und
  `FenceCloses` leben an einer Stelle; alle fünf Konsumenten speisen sich daraus,
  jeder mit eigener Assertion gegen Wieder-Divergenz.
- **Ein §7-Eintrag im Benutzerhandbuch** (1.44), der den Befund benutzbar macht:
  ab der gemeldeten Zeile rückwärts lesen, und die zwei Fälle, die wie Fehlalarm
  aussehen und keiner sind.

## Was hat funktioniert?

- **Erst messen, dann entscheiden.** Die Bestandsmessung über 776 Dateien in drei
  Repos hat die Entscheidung gedreht — der als „Wurzel" geplante Kandidat wäre
  wirkungslos gewesen **und** löst den belegten Fall nicht. Ohne die Zahl wäre
  die falsche Variante gebaut worden.
- **Drei Review-Runden statt einer.** Jede war blockierend, jede fand dieselbe
  Klasse an einer neuen Stelle. Nach Runde eins hätte der Slice abnahmereif
  ausgesehen — und einen offenen Silent-Grün-Pfad im Vollständigkeits-Gate
  gehabt.
- **Die Mutations-Gegenprobe** hat in jeder Runde einen echten Testfehler
  gefunden: einen nicht verdrahteten Aufruf, eine ungetestete Lesart, drei
  unassertierte Konsumenten.

## Was ging anders als geplant?

- **Das geplante Ende (2026-08-10) hielt**, aber der Weg dahin nicht: geplant war
  ein Slice mit einem Review, geliefert wurden drei Review-Runden und zwei
  Heilungs-Commits. Kein Closure-Kriterium war berührt — die Welle war nie in
  Gefahr, nur länger als gedacht.
- **Die Mutations-Gegenprobe war beim ersten Anlauf selbst kaputt** (`git
  checkout` statt Dateikopie). Sie hat dabei die Implementierung im Arbeitsbaum
  gelöscht und drei Ergebnisse aus dem falschen Grund als „rot" gemeldet. Das
  Werkzeug, das die Tests prüft, braucht dieselbe Sorgfalt wie die Tests.
- **Zwei Zusagen mussten ehrlicher werden statt repariert.** Die gemeldete Zeile
  heißt jetzt **Fundstelle** (welche Öffnung fehlt, ist nicht entscheidbar), und
  die Grenze wuchs über drei Runden bis zu `pins`, wo die Folge still ist.
- **Eine Verhaltensänderung am ausgelieferten Gate** kam ungeplant dazu: die
  Trimm-Vereinheitlichung ändert das Closure-Gate aus v0.52.0 in **beide**
  Richtungen. Steht in CHANGELOG, Lastenheft-Delta und SemVer-Begründung — nach
  einem Review-Befund, nicht von selbst.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Bestandsmessung vor Varianten-Entscheid** hat sich zum zweiten Mal getragen
  (nach slice-096). Als erster Abnahme-Punkt in
  [slice-103](slice-103-geteilte-lexik-raender.md) übernommen.
- **Mutations-Gegenprobe über eine Dateikopie**, nie über `git checkout`, und mit
  Abbruch, wenn die Ersetzung nicht greift. Ebenfalls in slice-103 als
  DoD-Punkt.

## Beobachtungs-Register (Zeiger)

- **BEO-003** neu: *eine geteilte Lexik driftet an den Rändern, weil jeder
  Konsument sie selbst vorbereitet.* Für die Fence-Lexik geschlossen, für andere
  offen.
- **BEO-004** neu bei Zähler **3**: *eine Modul-Grenze wird nur auf der
  Quell-Achse gedacht.* Dreimal in den drei Review-Runden dieser Welle gefunden —
  die Verkörperungs-Schwelle ist erreicht, die Form ist noch zu entscheiden.

## Folge-Slices

- [slice-103](slice-103-geteilte-lexik-raender.md) — dieselbe Klasse in
  anderen Lexiken (`citations`-Absatzbildung, Anker-Auflösung, git-Revisionen).
  Aus der dritten Review-Runde geschnitten, ausdrücklich **nicht** in dieser
  Welle erledigt: eigene Verträge, eigene Module.
- [slice-099](slice-099-structure-modul.md) — der Trigger der
  Folge-Welle ist mit dieser Closure **eingetreten**: der Fence-Defekt ist
  behoben, das neue Modul erbt ihn nicht mehr.

## Verifikation

- **Closure-Trigger erfüllt:** der eine Slice dieser Welle
  ([slice-101](slice-101-fence-unbalanciert.md)) liegt in `done/`.
- `make fullbuild` grün; `make ci` vor dem Tag grün (Gates + Image-Test).
- Release **v0.53.0** gebaut und nach GHCR gepusht, Pipeline-Lauf 31362555433
  erfolgreich, Digest `sha256:0cbe2d54…9424`.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts; keine
  stehengebliebene Gate-Reifestufe; von den ADR-Re-Evaluierungs-Triggern ist
  keiner eingetreten — [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md)
  nennt drei (Falsch-Positive, reale Fehlpaarungen im Bestand, ein weiteres
  innerhalb eines Abschnitts messendes Modul), und alle drei sind offen: die
  dritte Review-Runde fand auf 796 realen Dokumenten null Falsch-Positive, der
  Bestand trägt weiterhin keine gemischten Fence-Längen, und kein neues Modul
  ist hinzugekommen.
