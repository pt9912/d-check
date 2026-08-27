# Slice slice-158: Der `citations`-Scan sieht Inline-Code nicht

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:**
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
(die Anforderung); [ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
(das Modul); [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
(die geteilte Lexik und ihre gescopten Ausnahmen);
[slice-152](../next/slice-152-citations-scharfschalten.md) (der Anlass).

**Berührte Spec-Stellen:**
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
und [`DC-FA-CITE-001.a`](../../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) Schritt 1.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

[`DC-FA-CITE-001.a`](../../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) sagt über den Marker-Scan: *„Arbeitet auf den rohen Zeilen
(fence-aware wie die übrigen Module)."* Das ist eine **Zusage**, kein Versehen —
und sie hat eine Folge, die erst beim Scharfschalten sichtbar wird: **die
Dokumentation der Direktive ist selbst ein Fund**. Wer die Syntax in
Inline-Code schreibt — `<!-- d-check:cite <pfad>:<von>-<bis> -->` —, erzeugt
einen malformten Marker, und das Modul bricht fail-closed über den ganzen Lauf.

**Gemessen** ([slice-152](../next/slice-152-citations-scharfschalten.md)):
**72** Vorkommen des Markers in **20** getrackten Dateien, davon **70
außerhalb** eines Fenced-Blocks. Neun der zwanzig sind eingefrorene
Review-Reporte, dazu ein `done/`-Slice und zwei `Accepted`-ADRs — alle
unantastbar. Eine Doku-Konvention „Syntax nur noch in Fences" ist damit keine
teure Option, sondern unmöglich, und ein Ventil hat das Modul nicht: es führt
**keinen einzigen** Konfigurations-Schlüssel.

**Die Frage dieses Slice ist deshalb eine Vertrags-Frage**, keine
Implementierungs-Frage: soll der Marker-Scan Inline-Code überspringen wie jedes
andere prosa-lesende Modul — und wenn ja, ist das eine Schärfung der
vorhandenen Anforderung oder eine neue?

## 2. Vorgehen

1. **Die Vertrags-Frage zuerst.** [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
   trennt *„andere Antwort"* (Defekt) von *„andere Frage"* (per ADR gescopt) und
   führt `versions`, `pins`, `immutable` als gescopte Ausnahmen. `citations`
   steht dort **nicht** — die Rohzeilen-Lesart ist in der Spec zugesagt, aber
   nirgends als *andere Frage* begründet. Ob das eine Lücke oder eine
   ungeschriebene Absicht ist, gehört entschieden, bevor Code entsteht.
2. **Den Preis beider Antworten messen.** Überspringt der Scan Inline-Code,
   verschwindet die Selbst-Fundstelle — aber auch jede **echte** Direktive, die
   jemand in Inline-Code setzt. Ob es solche gibt, ist zu zählen, nicht
   anzunehmen.
3. Trägt die Änderung: [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in) in
   [`spec/lastenheft.md`](../../../../spec/lastenheft.md) (Akzeptanzkriterium,
   Versions-Bump, Historie), [`DC-FA-CITE-001.a`](../../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) Schritt 1 in der Spezifikation,
   eine ADR mit `Schärft:`-Feld, dann Code und Tests.
4. **Die Fail-closed-Frage mitentscheiden**, weil sie an derselben Stelle
   hängt: ein malformter Marker nimmt heute den **ganzen Lauf** mit. Bei einem
   Modul im inneren Loop ist das eine andere Zumutung als bei einem
   Closure-Gate. Entweder bleibt es so — dann steht die Begründung im Vertrag —
   oder es wird ein Befund wie jeder andere.
5. Bewusstes Brechen je gedeckter Form, **Ursache gelesen**; Rückbau grün.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Scharfschalten.** Das ist [slice-152](../next/slice-152-citations-scharfschalten.md),
  und es wartet auf dieses Ergebnis.
- **Keine Auszeichnung von Zitaten.** Auch das gehört zu slice-152.
- **Keine Änderung an den drei gescopten Ausnahmen** aus
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) —
  `versions`, `pins`, `immutable` beantworten andere Fragen und bleiben, wie
  sie sind.

## 4. Definition of Done

- [ ] Die Vertrags-Frage ist entschieden: Schärfung, neue Anforderung oder
      bewusster Verzicht — mit Begründung gegen
      [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md).
- [ ] Bei Änderung: Lastenheft, Spezifikation, ADR, Code und Tests hängen
      zusammen; kein Stratum verweist abwärts.
- [ ] Die Zahl der **echten** Direktiven in Inline-Code ist gezählt, nicht
      angenommen.
- [ ] Die Fail-closed-Entscheidung steht im Vertrag, nicht nur im Code.
- [ ] Ein konstruierter Verstoß je gedeckter Form, **Ursache gelesen**.
- [ ] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Lexik-Änderung wirkt über das Modul hinaus.** Wer den Marker-Scan auf
  die geteilte Antwort umstellt, ändert eine Zusage, die andere Module teilen —
  und die Gegenprobe muss zeigen, dass sich für sie **nichts** ändert. —
  **Ausgang:** *(bei Closure)*
- **Die Selbst-Fundstelle ist ein Sonderfall, der wie ein Allgemeinfall
  aussieht.** Dass ausgerechnet die Doku der Direktive stolpert, verführt zu
  einer Regel, die nur diesen Fall löst. Der Bestand entscheidet, nicht der
  Anlass ([`BEO-011`](../observations.md)). — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Vertrags-Frage einen
Auftraggeber-Entscheid verlangt — dann bleibt die Lücke benannt, und
[slice-152](../next/slice-152-citations-scharfschalten.md) wartet weiter.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-011`](../observations.md) — die Regel gehört aus dem Bestand, nicht aus
  dem Anlass; [`BEO-017`](../observations.md) — ein rotes Gate muss vom
  geprüften Grund kommen.

Slice-ID: slice-158. Betroffene IDs:
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in).
Module: `citations`. Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Schärfung eines vorhandenen Vertrags an
einem vorhandenen Modul.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
