# Slice slice-150: Ein Zitat der Baseline ist pin-gebunden wie ein Link

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../welle-85-baseline-v5120-migration.md)
— **Etappe C**, geschnitten vom Delta-Audit in
[slice-149](../done/slice-149-baseline-v5120-delta-audit.md).

**Bezug:** [`MR-021`](../../../../harness/conventions.md#mr-021) (in-Repo-Verweise
sind pin-gebunden — der Eintrag, den dieser Slice schärfen würde);
[`MR-033`](../../../../harness/conventions.md#mr-033) (der eine Bestandsfall);
[`BEO-008`](../observations.md) (vierte Spiegel-Klasse, seit
[slice-148](../done/slice-148-baseline-v5120-vendoring.md) im Register).

**Berührte Spec-Stellen:** — (Harness-Regeltext; ein Produkt-Delta nur, falls
die Messung eine mechanische Form trägt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

[`MR-021`](../../../../harness/conventions.md#mr-021) bindet **Links** auf den
vendorten Baum an den Pin: beim Bump wandern sie mit. Ein **Zitat** des
vendorten Textes wandert nicht — der Pfad daneben wird gehoben, löst sauber auf,
und der zitierte Wortlaut existiert am neuen Ziel nicht mehr. Beide Hälften sind
für sich in Ordnung, und kein Gate sieht die Kombination.

Bestandsfall: [`MR-033`](../../../../harness/conventions.md#mr-033) zitiert
zweimal die Fassung, die der eigene Konsumenten-CR hat ändern lassen.

**Die Vorfrage ist keine Bauentscheidung, sondern eine Regelfrage:** Ein
`Accepted`-Eintrag wird nach dem Kanon **nie überschrieben**. Ist die Korrektur
eines veralteten Zitats ein Überschreiben — oder die Pflege einer
pin-gebundenen Referenz, die [`MR-021`](../../../../harness/conventions.md#mr-021)
ohnehin verlangt? Davon hängt ab, ob der Bestandsfall repariert werden **darf**.

## 2. Vorgehen

1. **Die Regelfrage zuerst und begründet beantworten** — mit dem Kanon, nicht
   mit der Bequemlichkeit. Der Geltungsbereich von
   [`MR-021`](../../../../harness/conventions.md#mr-021) nennt heute
   ausdrücklich **Links**; ihn auf Zitate zu lesen, wäre genau die
   Reichweiten-Dehnung, die [`BEO-012`](../observations.md) führt.
2. Fällt sie für die Pflege: ein **neuer** Eintrag, der
   [`MR-021`](../../../../harness/conventions.md#mr-021) **schärft** (Titel
   trägt das nach Kanon-Form) — nicht dessen Änderung.
3. Erst dann den Bestandsfall reparieren.
4. **Messen, ob eine mechanische Form trägt:** die Prüfung ist *Zitat gegen
   sein eigenes Link-Ziel*, nicht gegen einen Korpus — der Korpus-Test ist in
   [slice-148](../done/slice-148-baseline-v5120-vendoring.md) nachweislich
   gescheitert. Ob das Produkt das heute hergibt (`citations` ist ein Modul
   dieses Repos), ist zu **messen**.
5. Die Bump-Prozedur um den Schritt ergänzen, falls er trägt.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung an [`MR-021`](../../../../harness/conventions.md#mr-021)
  selbst.** Einträge werden nicht überschrieben.
- **Keine Reparatur eingefrorener Zitate.** `done/` und Review-Reporte zitieren
  den Stand ihrer Zeit; der Bestandsfall ist der **lebende** Eintrag.
- **Kein Gate ohne Messung.** Ein Wächter auf Zitate hat eine
  Falsch-Positiv-Last, und die gehört gezählt, bevor er gebaut wird.

## 4. Definition of Done

- [x] Die Regelfrage ist mit Kanon-Fundstelle beantwortet — **im zweiten
      Anlauf**. Der Kanon entscheidet sie:
      [`modul-02` §Freshness-Audit](../../../../.harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)
      hält *„bestehende werden nicht rückwirkend umgeschrieben"* und
      *„**Rückbau ist ein neuer Eintrag, kein Edit**"* fest (§9).
- [x] **Nicht** die Pflege-Antwort — der erste Entwurf ging diesen Weg und ist
      aufgelöst ([`MR-038`](../../../../harness/conventions.md#mr-038) in `conventions/done/`).
- [x] Der Bestandsfall ist **ausgewiesen statt repariert**:
      [`MR-033`](../../../../harness/conventions.md#mr-033) steht wortgleich wie am 23.08.; die Ablage-Form
      der Feststellung ist [`MR-039`](../../../../harness/conventions.md#mr-039), der **beide** veralteten
      Zitate mit alter und neuer Fassung führt.
- [x] Die mechanische Form ist gemessen: sie **existiert** (Modul `citations`,
      seit `v0.50.0`, Zitat gegen die eigene Quelle) und ist **heute nicht
      aktivierbar** — zehn Markdown-Dateien tragen den Direktiven-Marker, der
      Lauf bricht an zwei Stellen des Algorithmus. Geschnitten als
      [slice-152](../open/slice-152-citations-scharfschalten.md).
- [x] `make gates` Exit 0 (zehn Glieder), `make fullbuild` Exit 0; unabhängiger
      Review ([Report](../../../reviews/2026-08-26-slice-150-pin-gebundene-zitate-review.md)),
      blockierend mit **einem HIGH**, sieben MEDIUM und sechs LOW — alle
      vierzehn eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Die bequeme Antwort ist die Pflege-Antwort.** Sie erlaubt die Reparatur und
  spart einen Eintrag. Genau darum ist sie die verdächtige — die Begründung muss
  aus dem Kanon kommen, nicht aus dem Ergebnis. — **Ausgang:** *eingetreten,
  und zwar vollständig.* Ich habe die bequeme Antwort gewählt, sie mit einer zu
  engen Kanon-Suche gestützt (*„der Kanon entscheidet es nicht"* — er
  entscheidet es) und die daraus gebaute Regel im selben Commit gebrochen. Der
  Review hat es gefunden; die Auflösung ist
  [`MR-038`](../../../../harness/conventions.md#mr-038) → [`MR-039`](../../../../harness/conventions.md#mr-039). **Das Risiko war wörtlich
  richtig formuliert und hat trotzdem nicht geschützt** — es zu notieren ist
  nicht dasselbe, wie es zu befolgen.
- **Ein Zitat-Wächter misst Prosa gegen Prosa.** Umbrüche, Auszeichnung und
  Auslassungszeichen machen jede naive Gleichheit falsch; in
  [slice-148](../done/slice-148-baseline-v5120-vendoring.md) hat dieselbe
  Prüfung mehrere verworfene Fassungen gebraucht, bis sie trug — die dort
  genannte Zahl **sechs** zählt Dokument/Quell-Paare, nicht Anläufe. —
  **Ausgang:** *eingetreten — aufgefangen von
  [slice-152](../open/slice-152-citations-scharfschalten.md).* Die Messung
  bestätigt das Risiko und entschärft es zugleich: das Produkt normalisiert
  bereits Whitespace und prüft gegen die **eigene Quelle** statt gegen einen
  Korpus. Was bleibt, ist die Wartungslast der Zeilennummern und die
  Fail-closed-Semantik — beides steht als Risiko in slice-152.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Regelfrage eine Klärung des
Kanons braucht — dann ist sie ein CR-Kandidat, kein Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Konventionsspeicher (GF), Produkt-Module (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-008`](../observations.md) ist der Anlass;
  [`BEO-012`](../observations.md) für die Frage, wie weit
  [`MR-021`](../../../../harness/conventions.md#mr-021) trägt.

Slice-ID: slice-150. Betroffene IDs: — (kein `DC-`-Bezug, solange kein
Produkt-Delta entsteht). Module: Konventionsspeicher. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Regelfrage am eigenen Konventionsspeicher.

## 9. Closure-Notiz (nach `done/`)

Geliefert: die Regelfrage beantwortet, der Bestandsfall ausgewiesen, die
mechanische Form gemessen — und ein Rückbau der eigenen ersten Antwort.

**Der Kanon entscheidet die Frage. Ich habe zu eng gesucht.** Der erste Entwurf
stand auf *„keine Stelle sagt, ob ein Zitat zum geschützten Kern gehört"*,
geprüft an zwei Dateien.
[`modul-02` §Freshness-Audit](../../../../.harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)
sagt es: *„bestehende werden nicht rückwirkend umgeschrieben"*, *„**Rückbau ist
ein neuer Eintrag, kein Edit**"*, und — der Satz, auf den es ankommt — *„Die
alte Zeile ist die historisch korrekte Aussage über den damaligen Zustand."*

**Die Regel wurde in ihrem eigenen Einführungsfall gebrochen.**
[`MR-033`](../../../../harness/conventions.md#mr-033) trug **zwei** veraltete Zitate, nicht eines. Der Entwurf
sagte *„ergänzt, nicht ersetzt"* und wurde auf eines angewandt; das zweite habe
ich ersatzlos entfernt. Die Commit-Botschaft dazu behauptete *„Nichts entfernt,
nichts umgeschrieben"* — bei einem Diff von zwölf Einfügungen und **vier
Löschungen**. Das ist [`BEO-009`](../observations.md) Richtung (a): eine
behauptete Nicht-Änderung, die nicht stattfand.

**Der Rückbau folgt der Form, die der Kanon dafür vorsieht.**
[`MR-033`](../../../../harness/conventions.md#mr-033) ist wortgleich wiederhergestellt.
[`MR-038`](../../../../harness/conventions.md#mr-038) ist **aufgelöst, nicht korrigiert** — ein Eintrag wird
nicht überschrieben, auch der eigene falsche nicht.
[`MR-039`](../../../../harness/conventions.md#mr-039) tritt an seine Stelle und dreht die Mechanik um: der
zitierende Eintrag bleibt unangetastet, der **Bump-Eintrag** hält fest, welches
Zitat seit wann anders lautet.

**Das Risiko stand wörtlich im Plan und hat nicht geschützt.** §5 sagte: *„Die
bequeme Antwort ist die Pflege-Antwort … die Begründung muss aus dem Kanon
kommen, nicht aus dem Ergebnis."* Genau so ist es gekommen. Ein Risiko zu
notieren ist nicht dasselbe, wie es zu befolgen — und diese Lehre ist teurer als
der Befund, den sie beschreibt.

**Die mechanische Form war die ganze Zeit gebaut.** Das Modul `citations`
prüft ein ausgezeichnetes Zitat gegen **die von ihm zitierte Quelle**,
whitespace-normalisiert — genau die Form, an der meine eigenen Messungen in
diesem Strang viermal gescheitert sind. Ich hatte gefragt, „ob das Produkt das
hergibt". Es gibt es her, seit `v0.50.0`.

**Es läuft nur nicht an, und der Grund ist älter als dieser Slice.** Zehn
Markdown-Dateien tragen den Direktiven-Marker außerhalb von Fences; der Lauf
bricht an **zwei** Stellen des Algorithmus (15 malformte Marker, zwei
wohlgeformte ohne folgendes Zitat). Vier der zehn Dateien sind eingefroren.
**Neu ist das nicht:** der Design-Review des Moduls führt denselben Blocker
seit dem 18.07. als INFO-Befund. Mein erster Zensus nannte sechs Dateien und
einen Bruchpfad — er hatte nur die editierbaren Fundstellen gesehen.

**Register:** [`BEO-009`](../observations.md) auf Zähler **5**.
