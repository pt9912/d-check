# Slice slice-103: Dieselbe Klasse, andere Lexiken — Absatzbildung, Anker-Auflösung, git-Revisionen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Bezug:** [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in).
Vorgeschichte: [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md) und die drei
Review-Runden an slice-101.

**Autor:** pt9912. **Datum:** 2026-08-10.

---

## 1. Ziel

Drei Befunde aus der dritten Review-Runde an slice-101 abarbeiten, die dort
ausdrücklich **nicht** hingehörten: sie betreffen andere Module mit eigenem
Vertrag, und sie sind älter als der Fence-Wächter. Gemeinsam ist ihnen die
Klasse, die slice-101 für die **Fence**-Lexik geschlossen hat.

## 2. Die Klasse

*Eine geteilte Lexik driftet an den Rändern, weil jeder Konsument sie selbst
vorbereitet.* Bei den Fences waren es fünf Stellen, die dieselbe Zeile
verschieden trimmten, und zwei Automaten, die dieselbe Frage verschieden
beantworteten. Jede für sich vertretbar — zusammen ein stiller Grün-Pfad.

slice-101 hat das für die Fence-Lexik geschlossen: ein geteiltes Trimm-Prädikat,
beide Schluss-Lesarten geteilt und ausgewertet, je Konsument eine Assertion
gegen Wieder-Divergenz. Die drei Befunde hier zeigen dieselbe Bauform in
**anderen** Lexiken.

## 3. Die drei Befunde

1. **Absatzbildung in `citations`** (Review R-2, HIGH). Das Modul gruppiert
   Absätze selbst, und ein Fence ist dort **keine** Grenze. Dieselbe Datei
   liefert Exit 0 statt des zugesagten fail-closed Exit 2, nur weil zwischen
   Direktive und Zitat ein Code-Block statt einer Leerzeile steht.
2. **Anker-Auflösung in `headingSection`** (Review R-3, MEDIUM). Sie beantwortet
   die Anker-Frage roh: ein HTML-Anker **innerhalb** eines Fence erfüllt die
   fail-closed-Bedingung von `versions.current-from`, während `anchors` im
   selben Lauf `anchor-missing` meldet. Zwei Module, dieselbe Frage, zwei
   Antworten.
3. **git-Revisionen als dritte unerreichbare Eingabe-Achse** (Review R-7, LOW).
   `vcs` rechnet die fence-empfindliche Section-Maske auf git-Blobs, die kein
   scannendes Modul je sieht — ein Fixture belegt ein falsches
   `core-drift-vcs`. Das ist die **dritte** Wiederholung von „Modul-Grenze nur
   auf der Quell-Achse gedacht" (slice-101 Review F-3, N-2, R-7).

## 4. Abnahme-Punkte

1. **Erst messen, dann entscheiden** — wie bei slice-101. Wie viele Dokumente im
   Ökosystem lösen die drei Fälle heute aus? Die Fence-Messung dort drehte die
   Entscheidung (776 Dateien, null Vorkommen ⇒ latent statt aktiv); ohne Zahl
   ist die Reichweite jeder Variante Spekulation.
2. **Ein Slice oder drei?** Die Klasse ist gemeinsam, die Verträge sind es
   nicht. Zu entscheiden nach der Messung.
3. **Reparieren oder melden?** slice-101 hat den Zustand gemeldet statt die
   Paarung zu reparieren. Ob das hier trägt, ist offen: bei `citations` und
   `headingSection` geht es um eine **falsche Antwort**, nicht um einen
   unentscheidbaren Zustand.
4. **Die dritte Wiederholung.** „Modul-Grenze nur auf der Quell-Achse" hat mit
   R-7 die Schwelle des Beobachtungs-Registers erreicht. Zu entscheiden ist,
   welche Form sie verkörpert — die Register-Regel sagt: verkörpern statt
   weiterzählen.

## 5. Definition of Done

- [ ] Bestandsmessung für alle drei Fälle über die drei Repos.
- [ ] Abnahme-Punkte 1–4 entschieden, Vertragsanpassung geliefert.
- [ ] Je Fall ein mutations-echter Test; die Gegenprobe über eine **Dateikopie**,
      nicht über `git checkout` (die Lehre aus slice-101).
- [ ] `make gates` grün; SemVer-Einordnung begründet, **beide** Richtungen
      genannt.

## 6. Risiken / offene Punkte

- **Drei Verträge in einem Slice** könnten den Schnitt sprengen. Abnahme-Punkt 2
  entscheidet das nach der Messung, nicht vorher.
- **Der `vcs`-Fall ist LOW und teuer.** Die git-Achse für ein scannendes Modul
  erreichbar zu machen, ist ein anderer Umbau als ein geteiltes Prädikat.
  Möglich, dass hier nur die Grenze benannt gehört.

## 7. Trigger

**Start** (`open` → `next`): nach der Closure von slice-101 und wenn ein
WIP-Slot frei ist. Keine Kopplung an einen Release.

## 8. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** bei der Planung erneut lesen — die
  Ziel-Achsen-Klasse aus Abnahme-Punkt 4 steht dann voraussichtlich im Register.

## 9. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — pro Fall wird zuerst die Zusage formuliert
(welche Antwort gilt?), dann geliefert.

## 10. Closure-Notiz (nach `done/`)

_Ausstehend._
