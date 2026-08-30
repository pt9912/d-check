# MR-055 — Ein Symlink auf den vendored Baum ist ein pin-gebundener Träger (Nachtrag zu MR-021)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon kennt die Verzeichniskonvention
  und die Pin-Bindung
  ([`grundlagen-harness-dateien.md` §Verzeichniskonvention](../../.harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md#verzeichniskonvention));
  er kennt keinen Symlink als Träger, weil er kein Werkzeug voraussetzt, das
  Dateien über einen Alias einliest. Diese Adaption ergänzt den Träger, nicht
  die Bindung.
- **Datum:** 2026-08-30
- **Geltungsbereich:** **Symlinks** unterhalb von `.claude/rules/`, deren Ziel
  in `.harness/baseline/<tag>/` liegt. **Nicht** Markdown-Links auf denselben
  Baum — die führt [`MR-021`](../conventions.md#mr-021).

## Adaption

**[`MR-021`](../conventions.md#mr-021) bindet *alle Markdown-Links* auf
`.harness/baseline/<tag>/` an den Pin und überlässt die Menge dem Zensus der
Bump-Prozedur.** Ein **Symlink** ist kein Markdown-Link: ein Zensus, der nach
`](.harness/baseline/` sucht, findet ihn nicht. Er bindet denselben Pin und
bricht beim Bump genauso — nur unbemerkt.

**Die Lücke war real und ist gemessen.** Weder `links` noch `anchors` sehen
diese Dateien: der Scanner folgt Symlinks **nicht** in die Prüfmenge (an einem
Probe-Repo gefahren — eine *echte* Datei unter `.claude/rules/` wird gescannt
und ihr toter Link gemeldet, ein *Symlink* nicht). Und
`make baseline-verify` prüfte die gepinnten **Dateien**, nicht Aliase darauf:
`sha256sum -c` und die Manifest-Deckung bleiben grün, während der Alias ins
Leere zeigt.

**Die Antwort ist ein Sensor, keine Prozedur-Zeile.** `--verify` prüft seit
diesem Eintrag als **dritte** Frage, dass jeder Symlink unter `.claude/rules/`
auflöst; er läuft in `make gates`. Eine Prozedur hätte dieselbe Lücke gelassen,
die dieses Repo sonst überall benennt: sie hängt daran, dass jemand sie liest.

**Warum die Prüfung dort und nicht in einem Modul:** `baseline-verify` ist der
Ort, der die Integrität des gepinnten Bestands verantwortet. Ein Alias **in**
diesen Bestand gehört zu seiner Oberfläche. Ein Produkt-Modul dafür zu bauen
hieße, dem Werkzeug eine Frage zu geben, die es nur in **diesem** Repo hat.

## Grenze

**Geprüft wird die Auflösung, nicht das Ziel.** Ein Symlink unter
`.claude/rules/`, der auf eine Datei **außerhalb** des gepinnten Baums zeigt,
löst auf und passiert. Die Zusage lautet „kein toter Alias", nicht „jeder Alias
zeigt in den Pin".

**Und der Sensor sagt nichts über den Zweck.** Dass die verlinkten Module
tatsächlich in den Kontext geladen werden, ist eine Werkzeug-Eigenschaft und
von hier aus nicht prüfbar — sie ist beobachtet, nicht gewächtert.

## Auflösungs-Trigger

**Wenn der Scanner Symlinks in die Prüfmenge aufnimmt.** Dann gilt für sie
[`MR-021`](../conventions.md#mr-021) wie für jeden anderen Verweis, `links`
meldet den toten Alias, und dieser Eintrag ist Bestandsschutz ohne Gegenstand.

**Wenn Aliase außerhalb von `.claude/rules/` entstehen.** Der Geltungsbereich
nennt ein Verzeichnis, weil es heute genau eines gibt. Ein zweiter Ort ist
kein Anlass, die Regel zu dehnen, sondern sie zu prüfen.
