# MR-055 — Ein Symlink auf den vendored Baum ist ein pin-gebundener Träger (Nachtrag zu MR-021)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon kennt den vendored Baum als
  **präsente, gepinnte Referenz**
  ([`modul-02-harness-bootstrap.md` §Freshness-Audit](../../.harness/baseline/v5.18.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2))
  und die Verzeichniskonvention, die ihn beherbergt
  ([`grundlagen-harness-dateien.md`](../../.harness/baseline/v5.18.0/regelwerk/grundlagen-harness-dateien.md#verzeichniskonvention)).
  Dass **In-Repo-Verweise** auf diesen Baum an den Pin gebunden sind, ist
  dagegen keine Kanon-Aussage, sondern die Adaption
  [`MR-021`](../conventions.md#mr-021) — dieser Eintrag ist ihr Nachtrag, nicht
  der des Kanons. Der Kanon kennt keinen **Symlink** als Träger, weil er kein
  Werkzeug voraussetzt, das Dateien über einen Alias einliest.
- **Datum:** 2026-08-30
- **Geltungsbereich:** **Symlinks** im Baum unterhalb von `.claude/rules/` —
  rekursiv, Punkt-Namen eingeschlossen. **Nicht** Markdown-Links auf denselben
  Baum: die führt [`MR-021`](../conventions.md#mr-021).

## Adaption

**[`MR-021`](../conventions.md#mr-021) bindet *alle Markdown-Links* auf
`.harness/baseline/<tag>/` an den Pin und überlässt die Menge dem Zensus der
Bump-Prozedur.** Ein **Symlink** ist kein Markdown-Link: ein Zensus, der nach
`](.harness/baseline/` sucht, findet ihn nicht. Er bindet denselben Pin und
bricht beim Bump genauso — nur unbemerkt.

**Die Lücke war real und ist gemessen.** Weder `links` noch `anchors` sehen
diese Dateien: der Scanner folgt Symlinks **nicht** in die Prüfmenge (an einem
Probe-Repo gefahren — eine *echte* Datei unter `.claude/rules/` wird gescannt
und ihr toter Link gemeldet, ein *Symlink* nicht). Und `make baseline-verify`
prüfte die gepinnten **Dateien**, nicht Aliase darauf: `sha256sum -c` und die
Manifest-Deckung bleiben grün, während der Alias ins Leere zeigt.

**Die Antwort ist ein Sensor, keine Prozedur-Zeile.** `--verify` prüft als
**dritte** Frage, dass jeder Symlink unterhalb von `.claude/rules/` auflöst —
**rekursiv und dotfile-bewusst**, weil ein flacher Glob genau die Klasse
zurückholte, gegen die die Frage steht. Er läuft in `make gates`; seine Proben
fahren als `make baseline-probe`.

**Warum die Prüfung dort und nicht in einem Modul:** `baseline-verify` ist der
Ort, der die Integrität des gepinnten Bestands verantwortet. Ein Alias **in**
diesen Bestand gehört zu seiner Oberfläche. Ein Produkt-Modul dafür zu bauen
hieße, dem Werkzeug eine Frage zu geben, die es nur in **diesem** Repo hat.

## Begründung

**Ein pin-gebundener Träger ohne Sensor ist die Bauform, die dieses Repo sonst
überall benennt.** Der Bump-Zensus findet, wonach er sucht; ein Träger, den er
nicht kennt, überlebt jede Prozedur-Zeile. Die Alternative — den Symlink in die
Bump-Prozedur zu schreiben — hätte dieselbe Lücke gelassen: sie hängt daran,
dass jemand sie liest. Und der Ausfall wäre **still**: vier tote Aliase, ein
grüner Lauf, und eine Zustellung, die niemand mehr bekommt.

## Grenzen

Sechs, ausgeschrieben statt implizit:

1. **Geprüft wird die Auflösung, nicht das Ziel.** Ein Alias, der auf eine
   Datei **außerhalb** des gepinnten Baums zeigt, löst auf und passiert. Die
   Zusage lautet „kein toter Alias", nicht „jeder Alias zeigt in den Pin" — die
   Implementierung prüft daher **jeden** Symlink dort, nicht nur die in den Pin.
2. **Ein fehlendes oder leeres `.claude/rules/` meldet nicht.** Wer die Aliase
   beim Bump **löscht** statt sie umzuhängen, hat einen grünen Lauf und eine
   verschwundene Zustellung. Der Sensor kann „gibt es hier nicht" nicht von
   „ist weg" unterscheiden, und ein Pflicht-Verzeichnis wäre für Adopter ohne
   Zustellung falsch.
3. **Ein Alias auf ein Verzeichnis passiert.** `readlink -e` gilt Verzeichnissen
   ebenso; die Frage lautet „löst auf", nicht „ist eine Datei".
4. **Die dritte Frage läuft nach den beiden ersten.** Bricht `sha256sum -c` oder
   die Manifest-Deckung ab, wird sie nicht mehr gestellt — die Befunde
   akkumulieren nicht.
5. **Ein Symlink überlebt nicht jedes Dateisystem.** Auf einem Checkout ohne
   Symlink-Unterstützung (`core.symlinks=false`) entstehen reguläre Dateien mit
   dem Pfad als Inhalt; der Sensor sieht dann keinen Symlink und schweigt, und
   die Zustellung ist eine Textzeile. Das ist eine Eigenschaft des Trägers, die
   kein Sensor im Repo abfangen kann.
6. **Der Sensor sagt nichts über den Zweck.** Dass die verlinkten Dateien
   tatsächlich in den Kontext geladen werden, ist eine Werkzeug-Eigenschaft und
   von hier aus nicht prüfbar — beobachtet, nicht gewächtert.

## Auflösungs-Trigger

**Wenn der Scanner Symlinks in die Prüfmenge aufnimmt.** Dann gilt für sie
[`MR-021`](../conventions.md#mr-021) wie für jeden anderen Verweis, `links`
meldet den toten Alias, und dieser Eintrag ist Bestandsschutz ohne Gegenstand.

**Wenn Aliase außerhalb von `.claude/rules/` entstehen.** Der Geltungsbereich
nennt einen Baum, weil es heute genau einen gibt. Ein zweiter Ort ist kein
Anlass, die Regel zu dehnen, sondern sie zu prüfen.
