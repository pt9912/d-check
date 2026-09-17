# Slice slice-219: Das Release, das `CO-001` auflöst

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
— und [`CO-001`](../../carveouts/done/CO-001-vcs-range-stiller-skip.md), dessen
Auflösung dieser Slice trägt.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Der Fix aus [slice-218](../done/slice-218-go-git-pack-namenskonvention.md)
erreicht das **publizierte** Image, und
[`CO-001`](../../carveouts/done/CO-001-vcs-range-stiller-skip.md) wird aufgelöst.
Solange das nicht geschehen ist, fährt jeder Konsument mit gepinntem
`v0.75.0` weiterhin den stillen Pfad — **der Fix im Quellstand hilft ihm
nicht.**

**Warum das ein eigener Vorgang ist und nicht der Rest von slice-218.**
Ein Carveout braucht einen Folge-Slice, der ihn **überlebt**
(Baseline-Regelwerk `modul-07-carveouts.md`). slice-218 schließt, sobald der
Fix committet ist — er kann die Auflösung also nicht tragen, und ein
Carveout, dessen Folge-Slice vor ihm schließt, hat faktisch keinen. Genau das
hat der unabhängige Review an der ersten Fassung von `CO-001` beanstandet.

**Abgrenzung — drei Punkte, jeder mit Grund:**

1. **Kein Produkt-Code.** Der Fix ist in slice-218 geliefert und dort
   gemessen; dieser Slice liefert ihn **aus** — *es wäre ein anderer Vorgang*,
   den Fix hier nachzubessern.
2. **Kein Vorziehen des Release-Zeitpunkts.** Wann geschnitten wird, ist eine
   Entscheidung des Auftraggebers, nicht dieses Plans; der Slice beschreibt,
   *was* dann zu tun ist — *Bestand bleibt bewusst stehen*.
3. **Keine weiteren Carveouts.** Es gibt heute genau einen; sollte bis dahin
   ein zweiter entstehen, löst dieser Slice ihn **nicht** mit —
   *Schicht-Abgrenzung*, sonst wächst er unbemerkt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** Ein Release nach `docs/user/releasing.md` ist veröffentlicht und
      am **gezogenen** Image verifiziert (Digest, OCI-Label, Smoke). `v0.76.1`,
      GHCR-Digest `sha256:1470ecdcaa686a5ef4513dee9b0ae522586f54b87d568b06fc6b5b2741b633b3`,
      `make ci` (inkl. `image-test`, Smoke eingeschlossen) grün vor dem Tag.
- [x] **(2)** Die Auflösung ist **gemessen, nicht angenommen**: dasselbe
      Probe-Repo-Muster mit partiell unsichtbarem Pack, gegen das
      **publizierte** Image gefahren — vorher `0 Befund(e)`/Exit 0, danach
      Exit 2. **Präzisierung ggü. dem Erstentwurf dieses Punkts:** `CO-001`
      nennt inzwischen **drei** Ausprägungen (unsichtbares BASE-Blob,
      BASE-Tree mit Pendant, BASE-Tree ohne Pendant) statt der hier
      ursprünglich genannten zwei — slice-220 hat die dritte geschlossen und
      eine vierte, zuvor fehldiagnostizierte (unlesbarer HEAD-Tree)
      zusätzlich aufgedeckt und behoben. **Alle vier** gegen das gezogene
      `v0.76.1`-Image gemessen, alle brechen mit Exit 2 ab (Belege in
      `CO-001`s Auflösungs-Trigger).
- [x] **(3)** `CO-001` ist aufgelöst: Verifikations-Haken abgehakt, Datei per
      reinem `git mv` nach `docs/plan/carveouts/done/`, Index in
      `docs/plan/carveouts/README.md` und die Bindung-Spalte in
      `harness/README.md` §Sensors nachgezogen.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

1. Release-Prep nach `docs/user/releasing.md` §Release-Prep in **einem**
   Commit; die Version trägt einen **Patch**- oder **Minor**-Bump je nachdem,
   was bis dahin sonst noch aufgelaufen ist.
2. `make ci` grün, Tag, Pipeline beobachten, `docker pull` als Beweis.
3. Die Auflösungs-Probe aus DoD (2) gegen das publizierte Image fahren —
   **vor** dem Abhaken der Verifikations-Liste in `CO-001`.
4. Carveout-Move und Index-Nachzug.

## 4. Trigger

**Beanspruchung:** slice-218 ist geschlossen **und** der Auftraggeber gibt das
Release frei.

**Rückführung nach `open/`** (`in-progress→open`): wenn die Auflösungs-Probe
gegen das publizierte Image den stillen Pfad **weiterhin** zeigt — dann trägt
der Fix nicht, und der Carveout bleibt aktiv.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `CO-001` liegt in `done/`, `make gates` grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Carveout altert, ohne dass etwas meldet.** Kein Sensor prüft, wie
  lange `CO-001` schon aktiv ist; die kanonische Frist misst in **Wellen**, und
  dieses Repo arbeitet wellenlos — eine benannte Lücke des Kanons
  (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).
  Der einzige Träger ist dieser Slice in `open/`. — **Ausgang:** entfallen.
  Der Carveout ist jetzt aufgelöst, bevor die Lücke praktisch relevant wurde
  (angelegt 2026-09-08, aufgelöst 2026-09-17 — neun Tage, kein Wellen-Maß
  nötig). Die benannte Kanon-Lücke selbst bleibt bestehen, betrifft aber
  keinen aktiven Carveout mehr.
- **Ein Release kann aus anderen Gründen fällig werden und den Carveout
  „nebenbei" auflösen**, ohne dass jemand die Probe aus DoD (2) fährt. Dann
  wäre der Carveout formal offen, obwohl der Defekt weg ist — oder schlimmer,
  er würde geschlossen, ohne dass die Auflösung gemessen wurde. — **Ausgang:**
  entfallen. Genau umgekehrt eingetreten: Dieser Slice **ist** das Release,
  das die Auflösung trägt, und die Probe aus DoD (2) wurde vor dem Abhaken der
  Verifikations-Liste gefahren (Plan-Schritt 3), nicht übersprungen.

## 7. Closure-Notiz

**Geliefert.** Release `v0.76.1` ist veröffentlicht (GHCR-Digest
`sha256:1470ecdcaa686a5ef4513dee9b0ae522586f54b87d568b06fc6b5b2741b633b3`,
Docker-Hub-Spiegel gleichgeprüft) und trägt slice-220s Klassen-Fix.
[`CO-001`](../../carveouts/done/CO-001-vcs-range-stiller-skip.md) ist
aufgelöst — alle vier Ausprägungen seiner Tabelle (nicht nur die zwei aus
diesem Slice-Plans Erstentwurf) wurden gegen das **gezogene** Image
gemessen, nicht nur den Quellstand angenommen; alle vier brechen jetzt mit
Exit 2 ab.

**Ein Nachzug, den weder slice-220 noch seine zwei Review-Runden fanden:**
`harness/sensors/adr-check.md` und `trace-check.md` beschrieben nach
slice-220s Feature-Commit weiterhin den abgelösten Diff-Mechanismus (drei
statt vier Ausprägungen, zwei nicht mehr existierende Testnamen, `CO-001`
als „weiterhin offen"). Der `git mv` von `CO-001` machte die Staleness
sichtbar (gebrochene Link-Tiefe im Doc-Check), nicht ein gezielter
Vergleich — beide Guide-Dateien sind jetzt nachgezogen. Registriert als
[`guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md):
`AGENTS.md` §6 Schritt 7 nennt „öffentlicher Vertrag", ohne zu sagen, dass
eine Guide-Doku wie `harness/sensors/*.md` denselben Rang trägt wie
Spec/Handbuch/README — ein Implementer, der nur die im Slice-Kopf genannte
Spec-Stelle prüft, lässt sie unbemerkt zurück.

**Review** ([`docs/reviews/2026-09-17-slice-219-release-loest-co-001-review-r1.md`](../../../reviews/2026-09-17-slice-219-release-loest-co-001-review-r1.md)):
**0 HIGH · 0 MEDIUM · 0 LOW · 0 INFO**, Verdikt „Freigegeben". Der Reviewer
verifizierte den Digest unabhängig (eigener `docker pull`, `gh release
view`, OCI-Label, Docker-Hub-Spiegel-Gleichheit) und fuhr alle vier
Ausprägungen selbst gegen fünf eigene Probe-Repos (die vier Defekte plus
eine Kontrolle mit einer echten, unversteckten Verletzung) — die Kontrolle
meldete korrekt `1 Befund(e)`/`core-drift-vcs`/Exit 1, alle vier Defekte
Exit 2.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

Dieses Repo führt **drei** Prüfungen — die zwei kanonischen und, als
Adaption, den Nachtlauf-Stand ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Die drei Vorprüfungen sind bei der Beanspruchung am 2026-09-17 gelesen.**

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Zwei** Sub-Areas: `*` (Repo-Default — Release-Prosa in `README*.md`,
`docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `version.md`) und
`docs/plan/carveouts/` (die Auflösung selbst: `git mv` nach `done/`,
Index-Nachzug). Keine trägt eine eigene Modus-Deklaration in
`harness/conventions.md` und fällt damit unter den Default `*`. Beide
Greenfield, wie der Rest des Produkts.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, 41 Verzeichnisse). **Keine**
Beobachtung trifft `docs/plan/carveouts/`, `docs/user/` oder die
Release-Prosa als Sub-Area; die einschlägigen Einträge aus slice-220
(`fix-schliesst-pfad-nicht-klasse`, `racily-clean-git-fixture`) betreffen
den bereits geschlossenen Code-Umbau, nicht diesen Release-Slice.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-17 gelesen: `image-scan.yml` grün
(2026-09-16T08:37:54Z). `upstream-drift.yml` **rot** (2026-09-17T05:38:51Z)
— laut eigener Meldung eine **planmäßige** Fremd-Release-Benachrichtigung
([`MR-051`](../../../../harness/conventions.md#mr-051)), keine unerwartete;
betrifft gepinnte Fremd-Bestände, nicht diesen Release-Slice.

**Modus-Begründung:** alle berührten Sub-Areas GF — kein Begründungsblock
nötig.
