# ADR-0054 — Geteilte Lexik bindet ihre Konsumenten; die Revisions-Achse bleibt eine benannte Grenze

**Status:** Proposed
**Datum:** 2026-08-16
**Autor:** pt9912
**Bezug:** [`DC-FA-CITE-001`](../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
[`DC-FA-VER-001`](../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-PIN-001`](../../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in),
[`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in);
Vorläufer [ADR-0050](0050-fence-unclosed-in-spans.md) (**nicht** superseded — dort
ging es um die Fence-Lexik selbst, hier um ihre Konsumenten); gescopte
Roh-Lesungen [ADR-0019](0019-versions-pin-fence-ausnahme.md) und
[ADR-0020](0020-content-pin-fence-ausnahme.md)

## Kontext

Die Fence-Lexik ist geteilt ([ADR-0050](0050-fence-unclosed-in-spans.md)): ein
Trimm-Prädikat, beide Schluss-Lesarten, je Konsument eine Assertion. Drei
Befunde derselben Review-Runde betrafen **andere** Lexiken und wurden bewusst
nicht dort behandelt.
Sie sind gemessen worden, und die Messung ist eindeutig — **alle drei sind
latent**: null von 18 `d-check:cite`-Direktiven, ein HTML-Anker in einem Fence
(ein Formbeispiel, das kein Konsument nachschlägt), null von 152 immutablen
Revisions-Blobs.

Latent heißt nicht harmlos: in zwei der drei Fälle **beantwortet ein Konsument
eine Lexik-Frage selbst und anders**, als dieselbe Frage an anderer Stelle im
Produkt beantwortet wird. Dieselbe Datei liefert dann je nach Trennzeichen zwei
entgegengesetzte Ergebnisse, und die Richtung, in der sie grün wird, ist die, in
der das Modul fail-closed abbrechen müsste.

Der dritte Fall hat eine andere Wurzel: dort ist die Lexik **nicht** uneins — die
Section-Maske ist fence-bewusst. Unerreichbar ist die **Eingabe**: `vcs` rechnet
auf git-Blobs, die kein scannender Wächter je sieht.

## Entscheidung

1. **Ein Konsument, der eine Lexik-Frage selbst beantwortet, ist ein Defekt —
   keine Variante.** Auch wenn seine Antwort für sich plausibel ist: zwei
   Antworten auf dieselbe Frage in einem Lauf sind ein stiller Grün-Pfad, den
   kein Gate sieht, weil jedes Modul für sich konsistent ist. Die Reparatur ist
   die **Übernahme der vorhandenen Antwort**, nicht eine zweite Implementierung.

2. **Die Trennlinie ist nicht „roh oder vorverarbeitet", sondern „andere Antwort
   oder andere Frage".** `versions` liest Pins auch in Fences
   ([ADR-0019](0019-versions-pin-fence-ausnahme.md)), `pins` hasht den rohen
   Ziel-Span einschließlich Fenced-Code
   ([ADR-0020](0020-content-pin-fence-ausnahme.md)), `immutable` normalisiert den
   rohen Core. Das sind **andere Fragen**, per ADR gescopt, und sie bleiben
   unangetastet. Wer dagegen fragt „ist das eine Überschrift", „ist das ein
   Anker", „ist das derselbe Absatz", bekommt die geteilte Antwort.

3. **Die git-Revisions-Achse bleibt eine benannte Grenze, kein Mechanismus.**
   Für Fall 3 wird **kein** Code geliefert: die Grenze wird im Vertrag benannt —
   und zwar in **beiden** Ausprägungen. Die gefährliche ist nicht die verschobene
   Maske (Falsch-Rot), sondern das **stille Grün**: liegt die Fence-Öffnung
   innerhalb des ausgenommenen Abschnitts, läuft die Ausnahme bis zum Dateiende,
   und eine reale Core-Änderung einer `Accepted`-ADR passiert das Gate mit Exit 0
   ohne Ausgabe. Weil an `vcs` selbst nichts zu beobachten ist, hängt der
   Trigger an der Arbeitsbaum-Fassung (unten).
   Begründung ist die Erreichbarkeit, nicht der Aufwand — bei den ersten beiden
   Fällen gibt es eine vorhandene richtige Antwort zu übernehmen, hier gäbe es
   nur einen neuen Wächter zu bauen, für einen Fall, der in 152 Revisions-Blobs
   nie eingetreten ist und für jeden bereits existierenden Commit ohnehin
   rückwirkend blind bliebe.

4. **Je Konsument eine Assertion.** Die Form aus
   [ADR-0050](0050-fence-unclosed-in-spans.md) gilt weiter: ein
   geteiltes Prädikat allein genügt nicht, sonst ist die Invariante nur ein
   Kommentar. Jeder reparierte Konsument bekommt einen Test, der die geteilte
   Antwort **an ihm** festnagelt.

## Alternativen

- **Alles reparieren, auch Fall 3.** Verworfen nach Entscheidung 3: der
  Mechanismus wäre ein neuer Wächter auf einer Eingabe-Achse, die keine der
  bestehenden Zusagen kennt — und der Anlassfall existiert im Bestand nicht.
- **Nichts reparieren, alle drei nur benennen** (die Antwort aus
  [ADR-0050](0050-fence-unclosed-in-spans.md): den Zustand melden statt die
  Paarung reparieren). Verworfen, weil dort ein
  **unentscheidbarer** Zustand vorlag; hier liegt eine **falsche Antwort** vor,
  zu der die richtige bereits im selben Binary existiert.
- **Die Lexik-Frage je Modul konfigurierbar machen.** Verworfen: das macht die
  Divergenz zur Funktion und verlagert die Entscheidung in die Config, wo sie
  niemand trifft.

## Konsequenzen

- Die Änderung wirkt in **beide** Richtungen, und die zweite ist die
  unangenehmere. **Mehr:** eine bisher grüne Direktive hinter einem Fenced-Block
  bricht fail-closed ab; ein Anker, der keiner ist, löst nicht mehr auf.
  **Weniger:** ein von einem Fence unterbrochenes Blockquote wird gekürzt gelesen
  und meldet sein `citation-mismatch` nicht mehr, und ein `dpin`, dessen Ziel-Anker
  nur im Fence steht, verliert seinen Drift-Schutz **kommentarlos** — das Modul
  schweigt zu unauflösbaren Zielen, das ist seine zugesagte Semantik und hier ihre
  unangenehme Folge. Gemessen betrifft beides heute **keinen** Konsumenten; die
  erste Fassung dieser Entscheidung behauptete „und weniger an keiner Stelle“,
  und das war als universelle Aussage falsch.
- Die drei gescopten Roh-Lesungen bleiben unberührt. Wer den Unterschied nicht
  mitliest, hält Entscheidung 1 für einen Widerspruch zu ADR-0019/0020 — deshalb
  steht er als Entscheidung 2 hier und nicht in einer Fußnote.
- Der Vertrag von [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  trägt die Grenze ab jetzt ausdrücklich. Sie
  ist damit auffindbar, statt in einem Review-Report zu liegen.

## Re-Evaluierungs-Trigger

- Ein Fall der Klasse tritt **im Bestand** ein (statt latent zu bleiben) — dann
  ist die Latenz-Begründung dieser Entscheidung verbraucht.
- **Beobachtbar gemachter Trigger für Entscheidung 3:** das Modul `spans` meldet
  ein `fence-unclosed` in einer Datei, die `vcs.paths` trifft. Dann ist der stille
  Pfad im Arbeitsbaum entstanden und über die Historie erreichbar geworden — an
  `vcs` selbst wäre nichts zu sehen gewesen, weil die stille Richtung per
  Konstruktion keine Ausgabe erzeugt. Ebenso, wenn das Immutabilitäts-Gate auf
  eine Klasse ausgeweitet wird, in der Revisionen routinemäßig unfertige
  Dokumente tragen.
- Eine **vierte** Stelle beantwortet eine Lexik-Frage selbst. Dann genügt die
  Aufzählung als Beleg nicht mehr, und die Frage ist, ob ein Gate die Klasse
  prüfen kann statt eines Reviews.

## Geschichte

- 2026-08-16: Proposed (`slice-103`, nach der Bestandsmessung).
