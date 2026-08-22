# Slice slice-113: Struktur-ID-Konvention zuerst — MR-027 aufgelöst, `ids`-Muster, ADR-Regel

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-80-struktur-ids](../welle-80-struktur-ids.md) (zugeordnet bei
der Eröffnung).

**Bezug:**
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(Kennungs-Linkpflicht — die neuen Muster laufen über dieselbe Mechanik),
[`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)
(ID-Schema-Aussage, trägt heute den Verzicht),
[`MR-027`](../../../../harness/conventions.md#mr-027) (die aufzulösende
Abweichung), Baseline-Regelwerk `grundlagen-source-precedence.md` §ID-Schema als
Klammer (Straten-Tabelle, §Vergabe „fortlaufend je Datei", „Der Link trägt den
Abschnitt, der Text die Kennung") und `modul-03-spec.md` §Zwei Kennungs-Arten;
Auftraggeber-Entscheide D1–D4 vom 2026-08-22.

**Berührte Spec-Stellen:** — (dieser Slice ändert Konvention und Konfiguration,
noch keine Spec-Zeile; die Vergabe folgt in slice-114/115).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die Regel steht, bevor die erste Kennung vergeben wird: der Konventionsspeicher
erklärt `SPEC-<NNN>`/`ARC-<NNN>` als vergeben (Baseline-Default, fortlaufend je
Datei, keine Lücken-Nachbelegung, Link auf den Abschnitt, Text = Kennung),
[`MR-027`](../../../../harness/conventions.md#mr-027) ist per `git mv`
aufgelöst, [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)
trägt den Verzicht nicht mehr, der ADR-Index nennt die Zwei-Formen-Regel
(Kennung, wo das Zielelement eine trägt; sonst §-Anker; `Accepted`-Bestand
bleibt unverändert), und `.d-check.yml` kennt die beiden Muster in
`ids.patterns` — **grün am Bestand**, weil noch keine dreistellige
`SPEC-`/`ARC-`-Kennung im gescannten Baum existiert (gemessen:
`.harness/baseline/**` liegt außerhalb des Scans, `docs/reviews/**` ist
exempt, der einzige Alt-Bestand sind zweistellige Diagramm-Beispiele in einem
`done/`-Slice).

## 2. Vorgehen

1. **Messen:** `git grep` nach `\b(SPEC|ARC)-[0-9]{3}\b` über den Baum außer
   `.harness/baseline` — erwartet null Treffer außerhalb exempt-Pfade; das
   Ergebnis steht in der Closure-Notiz.
2. **Konvention:** [`MR-027`](../../../../harness/conventions.md#mr-027) per `git mv` nach `harness/conventions/done/`
   (Move-Commit trägt die Link-Tiefen-Fixes der bewegten Datei,
   [`MR-013`](../../../../harness/conventions.md#mr-013)); Index-Zeile aus
   „Aktive Adaptionen" nach „Aufgelöste Adaptionen": *aufgelöst durch
   Baseline-Konformität (Struktur-ID-Vergabe, welle-80)* — die
   `<a id>`-Anker der Zeile wandern mit. [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)-Adaption: der Satz
   „Struktur-IDs werden nicht vergeben" wird zur Vergabe-Aussage
   (Baseline-Default, fortlaufend je Datei, Anker-Form).
3. **Gate-Konsument, Teil 1:** `.d-check.yml` `ids.patterns` um
   `\bSPEC-\d{3}\b` → `spec/spezifikation.md` und `\bARC-\d{3}\b` →
   `spec/architecture.md` (`link-policy: always`, dieselben `exempt-paths`);
   Kommentar in der Config trägt die Zusage (Wortgrenzen — sonst maskiert
   `ARC-01` eine dreistellige Kennung, der Präzedenz-Befund des `diagrams`-Reviews).
   Probe: `make doc-check` grün, und eine konstruierte Datei mit nackter
   dreistelligen `SPEC`-Kennung ⇒ `id-unlinked` (Gegenprobe, nicht committet).
4. **ADR-Index-Konvention** (`docs/plan/adr/README.md` Kopf): `Schärft:` nennt
   die Kennung, wo das Zielelement eine trägt, sonst den §-Anker; ADRs mit
   Status `Accepted` vor welle-80 bleiben auf §-Ankern. **AGENTS §5**: ein
   Satz zur Vergabe (nur beim Spec-Schreiben, nie ad hoc; nicht in
   Commit-Botschaften).
5. **Spiegel-Liste** ([`MR-025`](../../../../harness/conventions.md#mr-025),
   per `grep` nach dem alten Wortlaut „Struktur-IDs" / „nicht vergeben"):
   `harness/conventions.md` ([`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage), Index), [`harness/conventions/done/MR-027-struktur-id-verzicht.md`](../../../../harness/conventions/done/MR-027-struktur-id-verzicht.md),
   `AGENTS.md` §5, `docs/plan/adr/README.md`, `.d-check.yml`-Kommentar,
   `harness/README.md` (Sensors-Zeile `doc-check` nennt die Kennungs-Linkpflicht
   — prüfen, ob die Muster-Aufzählung dort steht).
6. Unabhängiger Review (Konvention + Config), Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Kennung wird vergeben** — Spezifikation und Architektur bleiben
  unverändert (slice-114/115).
- **Keine `structure`-Regel** (sie wäre rot am Bestand und blockierte über den
  `pre-commit`-Hook jeden Commit bis zur Vergabe — sie kommt mit slice-114).
- **Kein Produkt-Code**, keine Änderung an `--suggest-config`/`--print-config`.
- **Kein Eintrag in `commits.id-patterns`** (Struktur-IDs gehören nicht in die
  Traceability — Baseline).

## 4. Definition of Done

- [x] [`MR-027`](../../../../harness/conventions.md#mr-027) in `conventions/done/`, Index umgezogen, [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage) nachgezogen;
      `make doc-check` grün (Anker der Index-Zeile lösen weiter auf).
- [x] `ids.patterns` trägt `SPEC-\d{3}`/`ARC-\d{3}` mit Wortgrenzen;
      Negativ-Gegenprobe (nackte Kennung ⇒ `id-unlinked`) dokumentiert.
- [x] ADR-Index-Konvention + AGENTS §5-Satz; Spiegel-Liste abgehakt.
- [x] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Linkpflicht `always` wirkt sofort repo-weit, auch im Inline-Code:** jede
  nackte dreistellige `SPEC`-Kennung in Slices, Wellen, Reviews (außer exempt) wird Befund — auch Beispiel-Kennungen in Inline-Code; das Ventil dafür ist der Zeilen-Marker `d-check:ignore`. Das
  ist gewollt (Konsument), muss aber in der Konvention stehen, damit niemand
  es als Fehlalarm liest. — **Ausgang:** **eingetreten, adressiert** — die
  Wirkung ist real (drei Slice-Texte dieser Welle trugen Beispiel-Kennungen
  und wurden rot, bevor ein Commit durchging; im Bestand traf es eine
  Handbuch-Zeile mit einer fremden Kennung). Deklariert in
  [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage) und
  im Config-Kommentar, Ventil ist der Zeilen-Marker.
- **Zwei `<a id>`-Anker an der [`MR-027`](../../../../harness/conventions.md#mr-027)-Index-Zeile** werden von eingefrorenen
  Verweisen genutzt; sie wandern in die Tabelle „Aufgelöste Adaptionen" mit,
  sonst `anchor-missing`. — **Ausgang:** entfallen — beide Anker sind
  mitgewandert und je einmal vorhanden; alle eingefrorenen Verweise lösen auf
  (Review-Negativ-Probe).
- **Regex-Härte:** `\d{3}` ohne Wortgrenzen matcht in einer vierstelligen
  Kennung das Präfix. — **Ausgang:** **eingetreten und begrenzt** — die
  Wortgrenzen tragen gegen die Ziffern-Seite (gemessen: kein Befund an einer
  vierstelligen Kennung), **nicht** gegen Bindestrich-Präfixe: eine fremde
  Kennung meldet ihr Suffix. Als Grenze im Config-Kommentar benannt, Ventil
  ist der Zeilen-Marker; die Kennungs-Form eines Fremd-Repos ist kein
  Gegenstand dieses Musters.

## 6. Trigger

**Start** (`open` → `in-progress`): Wellen-Eröffnung (eingetreten) — dieser
Slice ist der erste der Welle.

**Rückführungen:** `in-progress` → `next`, falls sich zeigt, dass die
Linkpflicht `always` am Bestand rot ist (dann erst Bestand klären, dann
Muster).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-/Konventions-Doku (`harness/`, GF),
  Prüf-Profil `.d-check.yml` (GF), ADR-Index und AGENTS (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-006/008/009
  bei 2 — BEO-008 (Pin-Hebungs-Spiegel drei Klassen) ist **einschlägig** für
  den MR-Move: Pfad-Verweise (gate-gedeckt), Tag-URLs, Prosa-Pins — der
  Drei-Klassen-Zensus aus [`MR-028`](../../../../harness/conventions.md#mr-028) wird angewandt; BEO-002 wirkt als [`MR-025`](../../../../harness/conventions.md#mr-025)
  (Spiegel per `grep` nach dem alten Wortlaut, §2 Schritt 5).

Slice-ID: slice-113. Betroffene IDs:
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage),
[`MR-027`](../../../../harness/conventions.md#mr-027). Module: Konventionsspeicher,
Prüf-Profil, ADR-Index, AGENTS. Gates: `make doc-check` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Konventions- und Konfigurationsänderung
auf Baseline-Default; kein Legacy-Import.

## 9. Closure-Notiz (nach `done/`)

**Geliefert:** Die Regel steht vor der ersten Kennung. Der Move-Commit löste
die Abweichung auf (Datei nach `conventions/done/`, drei Link-Tiefen-Fixes,
Index-Zeile von „Aktive" nach „Aufgelöste Adaptionen" samt beider
`<a id>`-Anker), der Inhalts-Commit trägt die Vergabe-Aussage in
[`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)
(fortlaufend je Datei, keine Lücken-Nachbelegung, kein Bereichssegment, Link
auf den Abschnitt, nicht in Commit-Botschaften, `Accepted`-Bestand bleibt auf
`§`-Ankern), die beiden `ids`-Muster mit Wortgrenzen, die ADR-Index-Regel und
den AGENTS-Satz. Der Konsument wirkt sofort: gemessen an einer temporären
Probe meldet er nackte Kennungen und schweigt bei vierstelligen.

**Review** ([Report](../../../reviews/2026-08-22-slice-113-struktur-id-konvention-review.md)):
APPROVE — 0 HIGH, 0 MEDIUM, 3 LOW, 5 INFO, vierzehn Negativ-Proben und vier
Image-Gegenproben auf einer Scratch-Kopie. Alle drei LOW eingearbeitet: das
Bruch-Beispiel im Config-Kommentar beschrieb den Präzedenzfall eines anderen
Moduls statt des eigenen Musters; der `commits`-Kommentar nannte „die drei
`ids`-Muster" ohne die Abgrenzung, dass Struktur-IDs keine
Traceability-Frage beantworten; und der Ventil-Überblick des Handbuchs schrieb
dem Zeilen-Marker „nur `codepaths`" zu, obwohl Spezifikation, FAQ und Code ihn
seit der Kennungs-Linkpflicht für beide Module kennen — eine Bestandsdrift,
die genau dieser Slice benutzt hat.

**Was ging anders als geplant:** Der Plan sagte „grün am Bestand". Das stimmte
für den Bestand, nicht für die eigenen Planungs-Texte: drei Slice-Dokumente
dieser Welle nannten Beispiel-Kennungen in Prosa und Inline-Code und wurden im
Moment des Scharfschaltens rot — der `pre-commit`-Hook hielt jeden Commit auf,
bis sie auf Platzhalter-Form standen. Eine Regel, die Kennungen linkpflichtig
macht, macht auch das Reden über Kennungen linkpflichtig; wer sie einführt,
muss zuerst die eigene Sprache darüber ändern. Zweitens hat die Grenze der
Wortgrenzen ein reales Vorkommen getroffen: eine fremde Kennung im Handbuch
meldet ihr Suffix — die Wortgrenze schützt gegen Ziffern, nicht gegen
Bindestrich-Präfixe.

- **Steering-Loop-Eintrag:** Sensor geschärft: die Kennungs-Linkpflicht deckt
  jetzt die Struktur-IDs beider Spec-Straten — liegt in
  `.d-check.yml §ids.patterns`. Kein Auslöser aus dem Register; die Regel folgt
  dem Baseline-Default.
- **Beobachtungs-Register (`../observations.md`):** keine Beobachtung
  angefallen — BEO-008 (Drei-Klassen-Zensus einer Hebung) wurde angewandt und
  war leer außer der Pfad-Klasse; die Klasse ist nicht erneut eingetreten.
- **Folge-Slices:** [slice-114](../in-progress/slice-114-spec-vergabe-spezifikation.md)
  (Trigger mit dieser Closure eingetreten, inzwischen beansprucht) und
  [slice-115](../open/slice-115-arc-vergabe-architektur.md) als Datei in
  `open/`; [slice-116](../open/slice-116-adr-neuzugangs-regel.md) wartet auf
  beide.
- **Risiken aus §6:** alle mit Ausgang (§5) — zwei eingetreten und adressiert,
  eines entfallen.
- **Drei Paarungen:** Wellen-Slice — die Paarungen prüft die Welle-Closure.
