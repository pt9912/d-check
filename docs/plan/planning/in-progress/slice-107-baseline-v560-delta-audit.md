# Slice slice-107: Baseline-Bump v5.0.0 → v5.6.0 — Etappe B (Stufen-Audit)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-78-baseline-v560-migration.

**Bezug:** Präzedenz slice-085 (das v5.0.0-Modul-Delta-Audit: 18 Findings,
je mit Zuordnung); Grundlage ist der in
[slice-106](../done/slice-106-baseline-v560-vendoring.md) vendorte Baum. Kein
`DC-*`-Bezug — Lese-/Planungs-Arbeit.

**Autor:** pt9912. **Datum:** 2026-08-21.

---

## 1. Ziel

Das Regelwerks-Delta v5.0.0 → v5.6.0 **je Stufe gegen die Tag-Notizen** lesen
(nicht pauschal) und für jede neue oder geänderte Regel eine von drei
Antworten festhalten — **konform bereits** · **anzupassen** (mit Fundstelle
und Etappe-C-Kandidat) · **nicht anwendbar** (mit Begründung). Das Ergebnis
ist ein Findings-Register in diesem Slice, aus dem die Etappe-C-Slices
geschnitten werden.

Die sechs Stufen:

1. **v5.1.0** — §Vergabe unter §ID-Schema (wer vergibt die nächste Kennung).
2. **v5.2.0** — Straten-IDs, Bestands-Stichprobe, Reconciliation-Register.
3. **v5.3.0** — Kommentar-Regel (`grundlagen-harness-dateien.md`).
4. **v5.3.1/v5.4.0** — Korrekturen („Zwei Sensoren an derselben Aussage",
   „Zwei Kopien, zwei Antworten") + drei Regel-Ergänzungen.
5. **v5.5.0** — **Team-Fähigkeit** (größter Block: Rolleninhaber,
   Konflikt-Terminal, `team.md` dreistufiges SOLL; dazu `lab/team-sim/`).
6. **v5.6.0** — TA-7 nennt seine Wirkung (Hauptzweig-Regel, ein Absatz).

## 2. Leitfragen je Stufe

- Verlangt die Regel ein **Artefakt**, das d-check nicht führt (z. B.
  `team.md`, Reconciliation-Register)? Gilt sie auch für ein
  Ein-Operator-Repo, oder deklariert sie selbst ihre Grenze?
- Widerspricht eine **bestehende** d-check-Konvention der neuen Regel
  (Adaption nötig — `MR-*` — oder Anpassung)?
- Ist etwas, das d-check **bereits lebt**, jetzt Baseline-Default (dann ist
  ggf. eine bestehende Adaption auflösbar — Präzedenz: [`MR-018`](../../../../harness/conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)/019/020/022
  wurden bei v5.0.0 „Baseline-Stand"-aufgelöst)?
- **Wiedervorlage aus slice-090:** die upstream notierte 5-vs-6-Finding-
  Feld-Drift der Baseline — im v5.6.0-Stand behoben oder weiter offen?

## 3. Ausdrücklich NICHT in diesem Slice

Kein Editieren der eigenen Artefakte (außer diesem Slice selbst) — Etappe C
setzt um, dieser Slice liest und schneidet.

## 4. Definition of Done

- [ ] Je Stufe ein Abschnitt mit Findings-Tabelle (Regel · Antwort ·
      Fundstelle/Begründung); **kein** „pauschal konform" ohne je Regel eine
      Zeile (die welle-74-Lehre: eine Aufzählung, die vollständig heißt,
      braucht je Kandidat einen Negativbefund).
- [ ] Etappe-C-Schnitt: die „anzupassen"-Findings sind zu Slices gebündelt
      (Dateien in `open/`), die Wellendokument-§4-Tabelle ist nachgeführt
      (Roadmap-Drift-Log-Eintrag).
- [ ] Unabhängiger Review des Audits (Frischkontext gegen das vendored
      Delta); `make gates` grün (dieser Slice ändert nur Planungs-Doku).

## 5. Risiken

- **Die Team-Stufe könnte formal greifen, obwohl operativ ein
  Ein-Operator-Repo vorliegt** — dann ist die ehrliche Antwort eine
  deklarierte Adaption (`MR-*`) statt eines leeren Pflicht-Artefakts.
- **Additiv heißt nicht folgenlos:** auch eine neue Regel ohne
  Widerspruch kann ein stehendes Artefakt zur Pflicht machen (Präzedenz:
  das Beobachtungs-Register kam bei v5.0.0 genauso additiv).

## 6. Trigger

**Start** (`open` → `in-progress`):
[slice-106](../done/slice-106-baseline-v560-vendoring.md) in `done/` (der Audit liest
den **vendorten** Baum, nicht das Kurs-Repo) **und** WIP-Slot frei.

**Rückführungen:** `in-progress` → `next`, falls das Delta eine Vorfrage
aufwirft, die der Auftraggeber entscheiden muss (z. B. Team-Stufe adoptieren
vs. adaptieren).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-/Planungs-Doku, Repo-Default GF.
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-21: keine
  unverkörperte offen): **BEO-002**/[`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten) ist die Arbeits-Brille des Audits
  selbst — jede „anzupassen"-Zeile benennt ihre Spiegel gleich mit, statt
  sie Etappe C suchen zu lassen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Lese-Audit gegen die adoptierte
Konvention, dokumentierte Präzedenz (slice-085).

## 9. Stufen-Audit (Ergebnis, 2026-08-21)

Gelesen wurde der **vendorte** Baum (`.harness/baseline/v5.6.0/regelwerk/`),
je Stufe gegen die Tag-Notizen des Kurs-Repos; Antworten: **konform** (bereits
gelebt) · **anzupassen** (mit Etappe-C-Kandidat) · **n. a.** (mit Begründung).

### Stufe v5.1.0 — §Vergabe

| Regel | Antwort | Befund |
|---|---|---|
| §Vergabe: „Welche Form gilt, deklariert das Repo" — Bereichssegment oder dichte Nummern, die Wahl gehört in die ID-Schema-Deklaration, nicht in stille Gewohnheit | **anzupassen** → C-2 | d-check vergibt dicht (ein Schreiber) — gelebt, aber nirgends deklariert; [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)-ID-Schema-Aussage um die Vergabe-Deklaration ergänzen |
| Heading-Ebenen-Restrukturierung der Grundlagen-Dateien | konform | reine Baseline-interne Form; alle d-check-Anker lösen weiter auf (`make doc-check` grün über den vendorten Baum hinweg) |

### Stufe v5.2.0 — Straten-IDs · Bootstrap-Schärfungen · Replay-Form

| Regel | Antwort | Befund |
|---|---|---|
| Verfeinerungs-Suffix `<PREFIX>-FA-<NN>.<Buchstabe>` im Technik-Stratum | **konform** | seit jeher gelebt: jede Algorithmus-Sektion der Spezifikation heißt `DC-FA-…-001.a` |
| Struktur-IDs `SPEC-*`/`ARC-*` | **anzupassen** → C-2 | der `§`-Anker-Rückfallweg deckt nur Referenzen auf kennungslose **Elemente**, nicht den repo-weiten **Vergabe-Verzicht**: `modul-03` schreibt `SPEC-*` für genau die Sektionstypen der Spezifikation §2–§6 vor. Der Verzicht ist damit eine **Abweichung** und geht den Widerspruchs-Ausgang: deklarierte Adaption **MR-027** <!-- d-check:ignore --> (dichte `.a`-Verfeinerungen + `§`-Anker statt Struktur-IDs, mit Begründung und Auflösungs-Trigger) — **nicht** still in der [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)-Aussage („keine inhaltlichen Adaptionen"), die sonst falsch würde (Review-Befund F-4) |
| Slice-Kopf normativ: betroffene Technik-/Sicht-ID, ersatzweise Spec-`§` | konform | das `**Bezug:**`-Feld nennt `DC-*`-IDs samt `.a`-Sektionen; die Sicht trägt keine Kennungen, der `§`-Rückfallweg deckt sie |
| Freshness-Audit Teil „Bestands-Stichprobe" (rotierend ein delta-freier Abschnitt je Sync — prüft die [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)-Aussage über nie übernommene Alt-Regeln) | **anzupassen** → C-6 | `--check-latest` deckt Currency + Content-Drift, nicht die Stichprobe; für diesen Bump ist sie einmal auszuführen. Delta-freie Menge v5.0.0..v5.6.0 (per Kurs-Diff verifiziert): `grundlagen-klassifikation`, `modul-00/01/14/15/16`; Kandidat: **`modul-14-docker-harness.md`** (der erste Anlauf nannte `modul-07` — das änderte sich in v5.2.0; Review-Befund F-2) |
| Gate-Fragment `d-check.mk` als Konsumenten-Baustein (Schritt 2) | konform | d-check ist der Producer (`--print-mk`); die eigene Nutzung läuft über die Makefile-Targets |
| Reconciliation-Register `reconciliation.md` (`RC-<NNN>`, Schritt 8) + Bestands-Inventur | **n. a.** | BF-spezifischer Bootstrap-Schritt („Diskrepanz-Schock"); d-check ist durchgängig GF, es gab nie eine Inventur mit Funden — ein leeres Pflicht-Artefakt anzulegen widerspräche „jedes Artefakt hat einen Konsumenten" |
| `grundlagen-bootstrap`: **Migration ≠ BF** (eine Baseline-Hebung ist keine Brownfield-Markierung) | konform | genau so gelebt: diese Welle läuft als GF-Migrations-Welle ohne BF-Markierung |
| `grundlagen-traceability`: `ARC-*` zählt **nicht** als Traceability-Bezug (geänderte Regel) | n. a. | d-check vergibt keine `ARC-*`; die RTM zählt `DC-*`-IDs — nichts zu ändern, nichts zu verlieren |
| Golden Set / `manifest.yaml` (Replay-Ziel-Form, `modul-12` +223) | **n. a.** | modul-12 skopiert selbst: „Modell" = der **nicht-deterministische Kern**; d-checks Kern ist per Vertrag deterministisch ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)) — das Closure-Kriterium „Replay grün" füllt hier wie bisher `make fullbuild` |

### Stufe v5.3.0 — Kommentar-Regel

| Regel | Antwort | Befund |
|---|---|---|
| §Was ein Kommentar trägt (fünf Klassen, Adressaten-/Zeitform-Test, drei Verbots-Klassen) | konform gelebt, **mit gemessener Rest-Klasse** → C-3 | die Disziplin (Funktion statt Review-Historie) ist etabliert — aber die Stichprobe des Reviews fand nackte **Review-Finding-Tokens** als Herkunfts-Feld (u. a. `cli.go`, `planning.go`, dazu die in dieser Welle selbst geschriebenen in `structure_tableorder.go`/`markdown.go`): repo-weit mehrdeutig, keine der drei zugelassenen Herkunfts-Formen. C-3 räumt die Fundstellen (Review-Befund F-9) |
| `modul-09`: **Emittierte Artefakte tragen keinen Anker** (Erzeuger-Kontext reist nicht mit) | **konform** | vom Review gegengeprüft: das emittierte `d-check.mk`-Fragment trägt nur `DC-FA-*`-Verweise, keine Slice-/Befund-Nummern des Erzeugers |
| **Träger-Pflicht:** Briefing (AGENTS §3) **plus** HIGH-Eintrag „Kommentar trägt keine der fünf Klassen" im Reviewer-Skill | **anzupassen** → C-3 | beides fehlt als verankerter Träger: AGENTS führt keine Kommentar-Hard-Rule, `reviewer.md` 1.4.0 keinen entsprechenden HIGH-Anker |

### Stufe v5.3.1/v5.4.0 — Korrekturen und drei Regel-Ergänzungen

| Regel | Antwort | Befund |
|---|---|---|
| MR-Index-Adressierung: von außen `conventions.md#mr-<NNN>` — der Anker trägt die **Kennung**, nicht den Titel (Titel-Anker brechen bei Umformulierung) | **anzupassen** → C-4 | d-checks Index-Zeilen tragen **Voll-Slug**-Anker (deklarierte Migrations-Schuld); Baseline-Form: kurzer Kennungs-Anker je Zeile **zusätzlich**, Alt-Slugs bleiben (genau die im slice-106-Review gestreifte Klasse: die [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)-Titelzelle wurde umformuliert, nur der eingefrorene Anker rettete die Verweise) |
| „Ein selbstgebautes Gate ist auf Zeit gebaut" — Obermenge-Nachweis je Verstoßklasse (Kandidaten-Menge · Bedingungen · Schwelle wie die ADR) | **konform** | gelebte Praxis aller Skript-Ablösungen (u. a. Paritäts-Matrizen bei arch-check und den Modul-Ablösungen) |
| Cross-Reference-Trigger normativ nur volatil→stabil über alle Straten | konform | maschinell kodiert (`matrix`, [`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)) |
| Sektions-Scope des „liegt in"-Felds (Anker-Paarung nur in den zwei Register-Sektionen) | konform | entspricht der gelebten Wellen-Closure-Praxis |
| **Widerspruchs-Ausgang** des Freshness-Audits (Regel-Ergänzung: ein Fund geht den Weg jeder Diskrepanz — übernehmen im nächsten Slice, behalten als deklarierte Abweichung/`MR`, oder Carveout; mehrere Funde treffen die [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)-Aussage) | **konform als Verfahren — in diesem Audit angewandt** | genau dieser Ausgang entscheidet den Struktur-ID-Fall oben (MR-027 statt stiller MR-000-Lüge) <!-- d-check:ignore -->; Review-Befund F-1: die Zeile fehlte ausgerechnet für die Einordnungs-Regel der eigenen Abweichungs-Fälle |

### Stufe v5.5.0 — Team-Fähigkeit

| Regel | Antwort | Befund |
|---|---|---|
| **Roadmap-Struktur: §Aktuelle Welle → §Offene Wellen** (derivativ: die flachen Welle-Dateien; Ruhe-Marker „Nichts in Arbeit" bei leerem `in-progress/`, mit Wächter) + Closure-Schritt 5 ohne Beförderung | **anzupassen** → C-1 (größter Punkt) | d-checks Roadmap, [`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform) und die `planning`-Selbstkonfiguration folgen der alten Form; das **Produkt** deckt die neue per Config (`planning.heading`/`marker` sind konfigurierbar), die Drift-Kopplung (genau ein flaches Wellendokument bei Anspruch) bleibt im Ein-Wellen-Betrieb wahr. Baseline-Default sticht (Auftraggeber-Linie seit v5.0.0) ⇒ adoptieren, [`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform) dabei auflösen oder fortschreiben |
| Leseordnung in `harness/README.md` (3–5 geordnete Zeiger, Menschen-Hälfte; ausdrücklich **kein** Closure-Prüfpunkt) | **anzupassen** → C-5 | fehlt; billig zu ergänzen |
| Slice-Kopf-Feld `Verantwortlich:` (bei `open→next` gesetzt; Deklaration, kein Sensor) | **anzupassen** → C-3 (going-forward) | neue Slices tragen das Feld künftig; kein Retrofit (template-forward wie beim Status-Feld) |
| `next → in-progress` landet auf dem Hauptzweig, **vor** der Arbeit (TA-7-Kern) | **konform** | gelebte Eröffnungs-Praxis jeder Welle |
| WIP-Limit = 1 **pro Rolleninhaber**; Rolleninhaber-Begriff; Konflikt-Terminal = ADR (Folge-ADR braucht neue Evidenz) | konform / n. a. (operativ) | Ein-Operator-Betrieb: WIP 1 gelebt; das ADR-Terminal entspricht der gelebten Review-/ADR-Praxis; Mehr-Schreiber-Teile sind laut Baseline selbst „entworfen, nicht belegt" |
| Neue Hard Rules / neue Reviewer-HIGH-Einträge tragen ab Einführung Auflösungs-Trigger oder „permanent" | konform (going-forward) | bereits AGENTS-§5-Praxis für ADRs; gilt ab jetzt auch für die C-3-Ergänzungen |
| §Vergabe-Ergänzungen: die **Welle** zählt repo-weit dicht (kein Sub-Area-Segment); „lokal ableitbar" hat eine Grenze — die nächste Nummer liest Verzeichnis **und** offene Welle-Dateien (auch §4-Zeilen ohne Datei vergeben Nummern) | **konform** | gelebt und in dieser Welle belegbar: welle-78 §4 vergab slice-108/109 vor Datei-Existenz, und die Nummern-Ableitung las beide Quellen (Review-Befund F-5) |
| `modul-10`: bei Abweichung zwischen Rolleninhabern wird der **Skill geschärft**, nicht die mildere Lesart gewählt (Dissens ≠ Nicht-Determinismus) | konform (going-forward) | entspricht der gelebten Praxis (Reviewer-Skill wird bei Befunds-Klassen geschärft, zuletzt 1.4.0); operativ Ein-Operator |
| `lab/team-sim/` | **n. a.** | Lehr-Material des Kurses, keine Regelwerks-Pflicht; ein `team.md` verlangt das Regelwerk nicht (kein Vorkommen im vendorten Baum) |

### Stufe v5.6.0 — TA-7 nennt seine Wirkung

| Regel | Antwort | Befund |
|---|---|---|
| Die Regel trägt die **Wirkung** (Anspruch sichtbar vor der Arbeit); bei push-geschütztem Hauptzweig deklariert das Repo seinen Träger als `MR` | **konform** | der beobachtete Träger ist der Direkt-Commit auf `main` (jede Wellen-Eröffnung landet dort vor der Arbeit); ob ein Schutzregime existiert, ist netzlos nicht verifizierbar — belegt ist die gelebte Praxis, und sie erfüllt die Wirkung (Review-Befund F-11) |

### Wiedervorlage aus slice-090

Die upstream notierte **5-vs-6-Finding-Feld-Drift** besteht in v5.6.0 fort:
`modul-10` §Output-Schema zählt fünf Felder (ohne `klasse`), das
Report-Template sechs. d-check folgt dem Template (sechs, seit slice-090) —
keine d-check-Handlung; als Upstream-Notiz an den Kurs weitergegeben.

### Etappe-C-Schnitt

Zwei Slices aus den sechs „anzupassen"-Findings:

- **[slice-108](../open/slice-108-roadmap-offene-wellen.md)** — C-1: Roadmap
  auf §Offene Wellen (Form + `planning`-Config + [`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)-Entscheid).
- **[slice-109](../open/slice-109-v560-konventions-nachzuege.md)** — C-2
  (ID-Schema-Deklaration: Vergabe + Struktur-ID-Verzicht), C-3
  (Kommentar-Regel-Träger in AGENTS + `reviewer.md` 1.5.0;
  `Verantwortlich:`-Feld going-forward), C-4 (Kennungs-Anker im MR-Index),
  C-5 (Leseordnung), C-6 (Bestands-Stichprobe `modul-14-docker-harness.md`);
  C-2 umfasst nach Review-Befund F-4 auch die neue Adaption **MR-027** <!-- d-check:ignore -->
  (Struktur-ID-Verzicht als deklarierte Abweichung).

## 10. Closure-Notiz (nach `done/`)

_Ausstehend._
