# MR-079 — Der Stop-Hook übernimmt HEAD-Stempel und Fail-closed-Trap von ai-harness-init v0.10.0 und bleibt streng, ohne Schalter (schärft MR-076)

- **Datum:** 2026-10-20
- **Geltungsbereich:** [`stop-require-gates.sh`](../../.claude/hooks/stop-require-gates.sh),
  [`record-gates.sh`](../../tools/harness/record-gates.sh),
  [`selbstpruefung.sh`](../../tools/harness/selbstpruefung.sh)
- **Ersetzt-Baseline-Regel:** — *(keine; schärft
  [`MR-076`](../conventions.md#mr-076), Härtung am Handoff-Gate nach Baseline-Regelwerk
  `modul-13-quality-gates.md` §Guard-Härtung)*
- **Adaption:** `record-gates` schreibt neben den Inhalts-Hash die aufgelöste SHA
  von `HEAD` nach `.harness/state/gates-passed.head`. Der Stop-Hook kennt zwei
  Modi. **Streng** — jedes Turn-Ende mit ungedecktem Inhalt blockt, auch ohne neuen
  `HEAD` — ist er **ohne jeden Schalter**. **Locker** — ein Turn-Ende, dessen `HEAD`
  dem gestempelten gleicht, geht auch mit ungedeckter Änderung frei — ist er nur mit
  der Datei `.harness/stop-gate-locker`; `STOP_GATE_STRENG` genau `1` hebt sie auf.
  **Dieses Repo legt die Datei nicht an und bleibt streng:** der lockere Modus ist
  eine Lockerung des Handoff-Gates ([`MR-004`](../conventions.md#mr-004),
  [`MR-005`](../conventions.md#mr-005) wählten den inhaltsbasierten Nachweis bewusst),
  und `AGENTS.md` §3.6 verlangt dafür eine ADR, die es nicht gibt. **Abweichung vom
  Upstream:** `ai-harness-init` v0.10.0 ist ohne Schalter **locker** und wird durch
  die Datei `.harness/stop-gate-streng` streng; hier ist es umgekehrt. Der Grund ist
  die Richtung des Ausfalls: eine Datei, die das Gate hält, lockert es still, wenn
  sie fehlt (`rm`, `git stash`, ein Branch-Wechsel, ein toter Symlink) — eine Datei,
  die lockert, kann nur absichtlich entstehen. Darüber hinaus endet jeder
  unerwartete Fehler des Hooks mit Exit 2 (blockierend) statt mit einem anderen
  Code; die benannten Lagen aus [`MR-076`](../conventions.md#mr-076) (Hash, Nachweis,
  `git status`) blocken weiter mit Grund. Die Selbstprüfung fährt im Wegwerf-Klon
  alle Zweige: streng ohne Schalter, locker mit der Datei, streng mit der Variable
  trotz der Datei, und der Commit ohne Nachweis blockt; liegt die Datei im Klon, bricht
  sie ab.
- **Grenze:** Wer `.harness/stop-gate-locker` anlegt, lockert das Gate, ohne dass ein
  Gate es meldet — gehalten wird das von der Selbstprüfung (`make selbstpruefung`), die
  kein Teil von `make gates` ist, und vom Review der Änderung; die Datei ist aber ein
  sichtbarer Eintrag im Diff, kein stiller Verlust. Ein frischer Klon ohne Nachweis und
  mit sauberem Baum geht weiter frei (die CI ist dort das Netz). Der Modus bindet nur
  Claude Code; wer die Hooks nicht lädt, umgeht beides.
- **Begründung:** Gemessen im Wegwerf-Repo gegen den Stand vor der Übernahme, neun
  Lagen (frischer Klon, Änderung ohne Nachweis, passender Nachweis, Änderung nach dem
  Nachweis mit und ohne Commit, unlesbarer Nachweis, fehlendes Hash-Skript, kaputter
  Index, `stop_hook_active`): alt und neu entscheiden gleich. Mit entferntem strengen
  Zweig gibt der Hook „Änderung nach dem Nachweis, kein neuer `HEAD`" frei — der
  Rot-Beleg. Der Review der ersten Fassung (eine Schalter-Datei `.harness/stop-gate-streng`
  hält den strengen Modus) fand den Mangel: im Scratch-Repo gab der Hook mit
  gelöschter Datei trotz ungedeckter Änderung frei. Außerhalb eines Repositorys endete
  der alte Hook mit Exit 1 (nicht blockierend), der neue mit Exit 2; ohne `git` im
  `PATH` mit Exit 127, der neue blockt mit Grund.
- **Auflösungs-Trigger:** ein Beschluss, den Stop-Hook zu lockern — dann ist eine ADR
  geschrieben, die Datei `.harness/stop-gate-locker` wird angelegt, und dieser Eintrag
  wird durch den Nachfolger ersetzt.
