# Slice slice-206: Modul `mentions` — Implementierung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
(Anforderung), [ADR-0084](../../adr/0084-mentions-eigenes-modul.md) (Bauform
und Verdikt-Semantik). Beide entstanden in slice-205 <!-- d-check:status-provenance -->; dieser Slice
liefert, was dort ausdrücklich aus dem Umfang genommen wurde.

**Berührte Spec-Stellen:**
[`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
— **keine Änderung**; der Slice setzt sie um. Berührt werden zusätzlich
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
(die Modulliste nennt `mentions` bereits) und
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(neuer Konfigurations-Block).

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Das Modul `mentions` bauen: Soll-Menge aus Pfad-Globs gegen Ist-Menge aus
Dokumenten, Befund `artifact-unmentioned` für jedes Mitglied, das in keinem
Dokument vorkommt. Umfang und Verdikt-Semantik sind entschieden; dieser Slice
verhandelt sie nicht neu.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Kern-Regel + Modul-Registrierung | neu | die Prüfung selbst, samt `validModules()`-Eintrag |
| Konfigurations-Schema | update | `mentions.artifacts`, `mentions.documents`, `mentions.match` |
| Tests | neu | das Akzeptanzkriterien-Trio der Anforderung als Fitness Function |
| [`docs/user/`](../../../user/) | update | Handbuch und Betriebsdoku nennen das neue Modul |

**Der Kalibrierungs-Auftrag ist der eigentliche Gehalt.** slice-205 schließt
mit dem Befund, dass d-checks eigene Gegenprobe auf der richtigen Menge grün
läuft — es gibt **kein eigenes Rauschen**, an dem sich Schwellen und
Ausnahme-Klassen justieren ließen. Dieser Slice muss deshalb an einem
**Fremd-Bestand** messen und das Ergebnis als solches ausweisen. Die drei in
slice-205 benannten Mengen sind der Startpunkt, nicht das Ziel.

## 3. Ausdrücklich NICHT in diesem Slice

- **Eine zweite Quell-Form der Soll-Menge** (Literale statt Pfade). Out-of-Scope
  (2) der Anforderung; sie braucht ihren eigenen Beleg.
- **Die Aufnahme in `make gates`.** Eine neue Modul-Klasse startet als
  eigenständiger Fokus-Lauf — dieselbe Vorsicht wie bei `review-coverage` und
  `trace-check`. Die Aufnahme ist ein späterer, eigener Entscheid.
- **Jede Änderung an [`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) oder [ADR-0084](../../adr/0084-mentions-eigenes-modul.md).** Zeigt die
  Implementierung, dass eine Festlegung nicht trägt, ist das ein Folge-Entscheid
  und keine stille Anpassung.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Das Modul `mentions` existiert, ist opt-in und erfüllt alle
      Akzeptanzkriterien von
      [`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
      als Tests — einschließlich beider fail-closed-Fälle und des
      byte-identischen Default-Laufs.
- [ ] **(2)** Ein Fokus-Target fährt es, und die Doku (`docs/user/`, Sensors-
      Tabelle) nennt es — mit `kein Gate` in der Zeile, solange es nicht in
      `gates` läuft.
- [ ] **(3)** Ein **Kalibrierungs-Beleg** an einem Fremd-Bestand liegt vor: je
      Mengen-Wahl die Zahl der Funde **und** das Urteil, wie viele davon Mängel
      sind. Ohne diesen Beleg ist die Ausnahme-Klasse geraten.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Befund-Ort ist festgelegt, aber ungewohnt.** `file` ist der
  Artefakt-Pfad, `line` der Platzhalter `1` — das Artefakt wird nie geöffnet.
  Sechs Module führen `Line: 1` bereits als Platzhalter; ob die Kombination aus
  nicht geöffneter Datei **und** Platzhalter-Zeile in `--doctor` und `--repair`
  sinnvoll erscheint, ist ungeprüft. — **Ausgang:** <offen>
- **Die Bezugsmenge muss in zwei Ausgabe-Formen erscheinen.** stderr im
  Default, `summary`-Felder unter `--json`/`--yaml`. Die zweite Hälfte
  erweitert ein Struct, das heute zwei Felder trägt; ob das rückwärtskompatibel
  bleibt, entscheidet sich an den vorhandenen Konsumenten. — **Ausgang:** <offen>
- **Kein eigenes Rauschen — geerbt aus slice-205.** Die Ausnahme-Klasse wird am
  Fremd-Bestand justiert, nicht am eigenen. Das weicht von der gelebten Praxis
  ab und ist der Grund für DoD (3). — **Ausgang:** <offen>

## 6. Trigger

**Start** (`open` → `in-progress`): slice-205 liegt in `done/`, Anforderung und
ADR stehen. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt sich, dass die Kalibrierung am
  Fremd-Bestand eine eigene Mess-Arbeit ist, wird sie ein eigener Slice vor der
  Implementierung.
- `in-progress` → `open` (blockiert): Trägt eine Festlegung aus
  [ADR-0084](../../adr/0084-mentions-eigenes-modul.md) nicht, ruht der Slice
  bis zum Folge-Entscheid — stillschweigend angepasst wird sie nicht.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a) das
Fokus-Target läuft grün und die Akzeptanzkriterien stehen als Tests; (b) der
Kalibrierungs-Beleg aus DoD (3) liegt vor.

## 7. Vorgelagert (vor der Modus-Begründung)

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:223-224 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice berührt Produkt-Code
(Kern-Regel, Konfigurations-Schema, Tests) und die Betriebsdoku — beides liegt
unter dem Default. `tools/harness/` ist **nicht** berührt: das Modul ist
Produkt, kein Harness-Werkzeug. Eine Ausdifferenzierung wäre künstlich, weil
alle drei Achsen des Inklusionskriteriums auf dieselbe Antwort zeigen.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:229-229 -->

> **Offene Beobachtungen sichten.**

Register durchgegangen (gemergter Stand, **35** Verzeichnisse — nachgezählt, nicht aus dem Vorgänger-Plan übernommen). **Fünf** Einträge
betreffen diesen Gegenstand — und diesmal ist nach beidem gesucht: nach dem
**Gegenstand** des Slice und nach dem, was er **anfasst**. Genau diese
Zweiteilung fehlte in slice-205 und erzeugte dort einen Fehl-Ausgang.

- [`module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
  (1×) — der einschlägigste. Das Modul liest **zwei** Eingaben und scannt nur
  eine: Ist-Dokumente als Text, Soll-Artefakte nur als Pfade aus dem
  Verzeichnisbaum. Die Anforderung schreibt diese Grenze bereits aus; der Slice
  muss sie im **Modul-Kommentar** wiederholen, wo der nächste Leser sie sucht.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (7×) — unmittelbar einschlägig für die Tests: Ein Test, dessen Fixture den
  Zustand gar nicht erzeugt, über den sein Name spricht, ist grün und von einem
  echten Wächter nicht zu unterscheiden. Beide fail-closed-Fälle sind genau
  solche Kandidaten — ihr Erfolgsfall ist der Abbruch, und ein Fixture, das
  versehentlich eine nicht-leere Menge liefert, macht den Test still nutzlos.
- [`spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
  (1×) — die Anforderung trägt **neun** Akzeptanzkriterien, nicht nur das Trio.
  Jedes einzelne braucht seinen Test, auch die beiden, die nur eine Ausgabe-Form
  betreffen.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (9×, mit slice-205 zuletzt gewachsen) — der Kalibrierungs-Beleg aus DoD (3)
  ist genau die Stelle, an der dieser Eintrag wieder zuschlagen kann: Wer am
  Fremd-Bestand misst und über das Modul spricht, muss beide Mengen benennen.
- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (1×, aus slice-205) — der Eintrag ist jung und trifft diesen Slice unmittelbar:
  Das Modul **ist** eine Zählmethode. Sein Ableiter — vor der Messung die Form
  des Gegenstands ausschreiben und die Trefferliste stichprobenweise gegen sie
  halten — gilt dem Kalibrierungs-Beleg wortwörtlich.

**Was ich geprüft und ausgeschlossen habe:**
[`modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
(2×) liegt nahe, weil der Slice `validModules()` erweitert. Sein
`Sub-Area`-Feld nennt `.d-check.yml`, `Makefile`, Gate-Doku, und seine drei
benannten Spiegel sind die der **Profil**-Modulliste. `mentions` kommt nicht ins
Profil (opt-in, eigener Fokus-Lauf wie `reviews`), `FOCUS_DISABLE` bleibt
unberührt, und die Lastenheft-Aufzählung trägt es bereits. **Ob der Eintrag
trotzdem greift, entscheidet sich am Ende an dem, was der Slice wirklich
angefasst hat** — nicht jetzt an einer Vermutung. In slice-205 wurde genau
diese Zuordnung zweimal zu schnell getroffen.

Keiner der fünf erreicht mit diesem Slice die Schwelle von 3× erstmalig.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-06 gelesen. `image-scan.yml` **grün**.
`upstream-drift.yml` **ROT**, und diesmal aus einem neuen und wichtigen Grund:
Die Baseline führt upstream **`v6.4.0`** — den angekündigten MINOR-Bump, der
die beiden angenommenen Bitten des ausgehenden CR einlöst
([Antwort](../../cr/2026-09-06-antwort-ai-harness-course-slice-formluecken.md)).
Alle fünf Versions-Achsen und der Content-Drift am gepinnten Tag sind
**grün**; rot ist allein die Currency-Achse. **Das berührt diesen Slice
nicht** — er fasst kein Planning-Artefakt an, dessen Form sich mit `v6.4.0`
ändert. Es berührt den **nächsten**: Die Pin-Hebung ist ein eigener Slice, und
mit ihr löst sich d-checks Slice-Haus-Form auf.

## 8. Sub-Area-Modus-Begründung

**Modus:** `*` ist **GF** (Greenfield, Repo-Default).

- **Konventionen-Dichte:** hoch. Der Modul-Schnitt ist durch
  [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) und
  [ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md) verankert und
  wird von `make arch-check` gehalten; die Bauform dieses Moduls ist zusätzlich
  in [ADR-0084](../../adr/0084-mentions-eigenes-modul.md) entschieden. Es gibt
  in diesem Slice keine offene Struktur-Frage.
- **Phase-Reife:** Phase 5 für die berührten Sektionen — Anforderung,
  Entscheid, Modul-Muster und Test-Layout liegen alle vor. Der Slice folgt
  einem sechsfach gelebten Muster (`reviews` als jüngstes).
- **Evidenz-/Diskrepanz-Risiko:** niedrig für den Code, **erhöht für den
  Beleg**. Der Code-Bestand kann von der Doku nicht divergieren, weil er noch
  nicht existiert. Die fünf gesichteten Beobachtungen zeigen aber, wo das
  Risiko wirklich liegt: nicht in der Implementierung, sondern in den
  **Aussagen über sie** — Tests, die nicht messen, was ihr Name sagt, und ein
  Kalibrierungs-Beleg ohne benannte Menge.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.
## 9. Closure-Notiz (nach `done/`)

<wird vor dem `git mv` nach `done/` gefüllt>
