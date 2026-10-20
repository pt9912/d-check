# MR-079 — Der Stop-Hook übernimmt HEAD-Stempel und Fail-closed-Trap von ai-harness-init v0.10.0 und bleibt streng (schärft MR-076)

- **Datum:** 2026-10-20
- **Geltungsbereich:** [`stop-require-gates.sh`](../../.claude/hooks/stop-require-gates.sh),
  [`record-gates.sh`](../../tools/harness/record-gates.sh),
  [`selbstpruefung.sh`](../../tools/harness/selbstpruefung.sh) und die versionierte
  Datei [`.harness/stop-gate-streng`](../../.harness/stop-gate-streng)
- **Ersetzt-Baseline-Regel:** — *(keine; schärft
  [`MR-076`](../conventions.md#mr-076), Härtung am Handoff-Gate nach Baseline-Regelwerk
  `modul-13-quality-gates.md` §Guard-Härtung)*
- **Adaption:** `record-gates` schreibt neben den Inhalts-Hash die aufgelöste SHA
  von `HEAD` nach `.harness/state/gates-passed.head`. Der Stop-Hook kennt zwei
  Modi. **Streng** — jedes Turn-Ende mit ungedecktem Inhalt blockt, auch ohne neuen
  `HEAD` — ist er, wenn die Datei `.harness/stop-gate-streng` existiert oder
  `STOP_GATE_STRENG` genau `1` ist. **Locker** — ein Turn-Ende, dessen `HEAD` dem
  gestempelten gleicht, geht auch mit ungedeckter Änderung frei — ist er nur ohne
  beides. **Dieses Repo trägt die Datei und bleibt streng:** der lockere Modus ist
  eine Lockerung des Handoff-Gates ([`MR-004`](../conventions.md#mr-004),
  [`MR-005`](../conventions.md#mr-005) wählten den inhaltsbasierten Nachweis bewusst),
  und `AGENTS.md` §3.6 verlangt dafür eine ADR, die es nicht gibt. Übernommen ist
  die Mechanik, nicht der Default. Darüber hinaus endet jeder unerwartete Fehler des
  Hooks mit Exit 2 (blockierend) statt mit einem anderen Code; die benannten Lagen
  aus [`MR-076`](../conventions.md#mr-076) (Hash, Nachweis, `git status`) blocken weiter
  mit Grund. Die Selbstprüfung fährt im Wegwerf-Klon alle Zweige: streng mit der
  Datei, locker ohne sie, streng mit der Variable, und der Commit ohne Nachweis blockt
  in beiden; fehlt die Datei im Klon, bricht sie ab.
- **Grenze:** Die Datei hält den strengen Modus nur, solange sie existiert: wer sie
  löscht, lockert das Gate, ohne dass ein Gate es meldet — gehalten wird sie von der
  Selbstprüfung (`make selbstpruefung`), die kein Teil von `make gates` ist. Ein
  frischer Klon ohne Nachweis und mit sauberem Baum geht weiter frei (die CI ist dort
  das Netz), und die Rückkehr zur gestempelten `HEAD`-SHA nach einem Commit gilt nur im
  lockeren Modus als „kein neuer Commit". Der Modus bindet nur Claude Code; wer die
  Hooks nicht lädt, umgeht beides.
- **Begründung:** Gemessen im Wegwerf-Repo gegen den Stand vor der Übernahme, neun
  Lagen im strengen Modus (frischer Klon, Änderung ohne Nachweis, passender Nachweis,
  Änderung nach dem Nachweis mit und ohne Commit, unlesbarer Nachweis, fehlendes
  Hash-Skript, kaputter Index, `stop_hook_active`): alt und neu entscheiden gleich.
  Mit entferntem strengen Zweig gibt der Hook „Änderung nach dem Nachweis, kein neuer
  `HEAD`" frei — der Rot-Beleg. Außerhalb eines Repositorys endete der alte Hook mit
  Exit 1 (nicht blockierend), der neue mit Exit 2; ohne `git` im `PATH` mit Exit 127,
  der neue blockt mit Grund.
- **Auflösungs-Trigger:** ein Beschluss, den Stop-Hook zu lockern — dann ist eine ADR
  geschrieben, die Datei `.harness/stop-gate-streng` fällt, und dieser Eintrag wird
  durch den Nachfolger ersetzt.
