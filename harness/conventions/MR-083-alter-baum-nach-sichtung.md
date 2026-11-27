# MR-083 — Der alte Baseline-Baum fällt nach der Sichtung der Unterschiede, und erst danach wird gezählt (Nachtrag zu MR-021)

- **Datum:** 2026-11-27
- **Geltungsbereich:** jede Hebung des Baseline-Pins ab der nächsten: das
  Entfernen von `.harness/baseline/<alt-tag>/` und die Feststellung, was den alten
  Pin noch nennt. Die Hebung von [`MR-081`](../conventions.md#mr-081) bleibt, wie
  sie lief.
- **Ersetzt-Baseline-Regel:**
  [`modul-02-harness-bootstrap.md` §Freshness-Audit](../../.harness/baseline/v7.0.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2),
  der Satz <!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-02-harness-bootstrap.md:268-269 -->
  „Das alte Verzeichnis fällt erst, wenn der Review durch ist."
  „Review" meint dort den Review des neuen Stands im Freshness-Audit, kein
  Code-Review des Hebungs-Diffs; zu ihm zählt derselbe Abschnitt den Formvergleich
  **und** den Gang durch die Adaptions-Liste. Hier genügt für das Entfernen die
  Sichtung der Unterschiede (unten); der Gang durch die Adaptions-Liste ist keine
  Bedingung dafür, und dieser Eintrag regelt ihn nicht.
- **Adaption:** Nach dem Vendoring des neuen Stands bleibt der alte Baum liegen.
  Dann, in dieser Reihenfolge:

  1. die Unterschiede beider Bäume **sichten**;
  2. **erst danach** den alten Baum entfernen
     ([`MR-021`](../conventions.md#mr-021) (1));
  3. **erst danach** feststellen, was den alten Pin noch nennt.

  **Gesichtet** ist, wenn eine committete Datei jede Datei nennt, die
  `diff -rq .harness/baseline/<alt-tag> .harness/baseline/<neu-tag>` meldet
  (ohne `SHA256SUMS`; eine eindeutige Kurzform genügt, bei einer Umbenennung
  einer der beiden Namen), und je Datei sagt, was daraus für dieses Repo folgt.
  „Nichts" ist eine Folgerung, ein Verweis auf ein späteres Dokument ist keine.
  Dateien, die nur Versionsstempel und Quell-Verweis ändern, dürfen eine Zeile
  teilen. Die Entscheidungen des Auftraggebers über die Folgerungen sind keine
  Bedingung. Ob eine Folgerung stimmt, ist Lektüre.

  **Die Feststellung** wird am Stand nach dem Entfernen gemessen (`make doc-check`,
  `make baseline-verify`, `git grep` nach dem alten Tag) und steht bei dem Vorgang,
  der entfernt (Commit-Botschaft oder Closure-Notiz). Ein früher committeter
  Pin-Eintrag trägt sie nicht und wird dafür nicht nachgetragen.

  **Abgeschlossen** ist die Hebung erst, wenn `ls .harness/baseline/` nur den
  gepinnten Tag zeigt: Was die Hebung trägt (die Welle, ohne Welle der Slice),
  geht vorher nicht nach `done/`.

  Die Link-Pin-Regel aus [`MR-021`](../conventions.md#mr-021) gilt weiter, samt
  [`MR-051`](../conventions.md#mr-051) und [`MR-055`](../conventions.md#mr-055).
  Auch (1) gilt; nur sein Zeitpunkt wandert: „darum ist (1) Pflicht-Teil des Bumps"
  heißt Pflicht-Teil der Hebung, nicht ihres Pin-Commits.
- **Begründung:** Bei der Hebung von [`MR-081`](../conventions.md#mr-081) fiel der
  alte Baum im Bump, und die Menge der Dateien, die den alten Pin nennen, wurde
  dort festgestellt. Die Sichtung kam danach
  ([`slice-296`](../../docs/plan/planning/done/slice-296-baseline-v700-neuerungen-sichtung.md));
  ihr Liefergegenstand nennt den alten Pin, stand nicht in der Menge, und ein
  Closure-Trigger von [`wave-98`](../../docs/plan/planning/wave-98-baseline-v700.md)
  §3 musste um diese Datei erweitert werden
  ([`BEO-ALL/closure-trigger-misst-was-ein-spaeterer-slice-veraendert`](../../docs/plan/planning/observations/BEO-ALL/closure-trigger-misst-was-ein-spaeterer-slice-veraendert/observation.md)).
  Nur das Löschen zu verschieben genügt nicht: Eine am Bump gezählte Menge kennt
  nicht, was die Sichtung danach schreibt.

  **Entfernt werden muss weiterhin.** Gemessen am 2026-11-27 an einem Wegwerf-Klon
  von `7e024045`, alter Baum aus `git archive 738d2bca`, die Verweise des Repos
  schon auf dem neuen Stand: Mit beiden Bäumen enden `make gates`,
  `make baseline-probe` und `make verify-closure-notes` mit 0. Ein dort neu
  eingefügter Link in `README.md` auf den alten Baum bleibt unter `make doc-check`
  ohne Befund und ist nach dem Entfernen `target-missing`; ein auf den alten Baum
  zurückgehängter Alias unter `.claude/rules/` besteht `make baseline-verify`
  (Exit 0) und fällt nach dem Entfernen (Exit 2). Solange beide liegen, belegt ein
  grüner Lauf den Nachzug also nicht.

  **Grenzen.** Kein Gate prüft Reihenfolge, Sichtung oder Entfernen; ein
  liegengebliebener alter Baum bleibt grün. Einen Sensor dafür gibt es bewusst
  nicht: Zwei Bäume sind während der Sichtung der erlaubte Zustand, den er rot
  färbte; den Rest trägt die Bedingung *Abgeschlossen*. `make baseline-verify`
  prüft nur den gepinnten Tag; den alten Baum prüft von Hand
  `tools/harness/fetch-baseline-cache.sh --verify <alt-tag>`. **Nicht gemessen:**
  der Vendoring-Lauf bei liegendem alten Baum (nach Lektüre leert
  [`fetch-baseline-cache.sh`](../../tools/harness/fetch-baseline-cache.sh) nur die
  Verzeichnisse des Ziel-Stands), die CI und der Zustand vor dem Nachzug der
  Verweise.
- **Auflösungs-Trigger:** permanent, solange eine Hebung den alten Baum entfernt
  ([`MR-021`](../conventions.md#mr-021) (1)).
