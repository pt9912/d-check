# MR-015 — Auflösung der MR-012-Pointer-Drift: AGENTS.md routet, spiegelt nicht mehr

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`grundlagen-harness-dateien.md` §Template-Schichtung](../../.harness/baseline/v5.18.0/regelwerk/grundlagen-harness-dateien.md#template-schichtung--was-der-rumpf-trägt-und-was-der-kommentar)
- **Datum:** 2026-06-22
- **Geltungsbereich:** [`AGENTS.md`](../../AGENTS.md) §1,
  [`MR-011`](../conventions.md#mr-011--baseline-auf-release-tag-gepinnt),
  [`MR-012`](../conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011),
  §[Adoptierte Konventions-Quellen](../conventions.md#adoptierte-konventions-quellen)
- **Adaption:** Nachtrag zu
  [`MR-012`](../conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011), das den
  Baseline-Pin u. a. „in den gespiegelten Pointern in AGENTS.md und
  harness/README.md" hob. Mit der zip-only-Umstellung von `AGENTS.md` §1
  (Commit `f46326a`, 2026-06-22 — Lese-Form ausschließlich `lab-regelwerk.zip`)
  wurde der `agents-regelwerk.md`-Link **aus AGENTS.md entfernt**: AGENTS.md
  **routet** dort für Quelldatei und Stand auf diese Datei
  (§[Adoptierte Konventions-Quellen](../conventions.md#adoptierte-konventions-quellen) /
  §[Baseline](../conventions.md#baseline)), statt die Raw-URL zu spiegeln. Die von
  [`MR-011`](../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)/[`MR-012`](../conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011)
  gepinnte `agents-regelwerk.md`-Raw-URL lebt damit nur noch in
  §[Adoptierte Konventions-Quellen](../conventions.md#adoptierte-konventions-quellen) und in
  [`harness/README.md`](../README.md) §Guides. Beide bleiben als
  Vergangenheits-Aussage **unverändert** (eigener Eintrag, analog
  [`MR-008`](../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)).
- **Begründung:** Die Pointer-Liste in
  [`MR-012`](../conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011) war am 2026-06-18
  korrekt; die spätere zip-only-Korrektur machte den AGENTS.md-Teil
  navigatorisch stale (wer „den gespiegelten AGENTS.md-Pointer" sucht, findet
  ihn nicht mehr). Eine undeklarierte Pointer-Verschiebung in einem
  MR-gepinnten Bereich ist dieselbe stille-Setzung-Klasse, die das Pinning
  verhindern soll; der Nachtrag stellt die Provenienz eindeutig wieder her. Die
  zip-Änderung selbst brauchte kein eigenes MR (reine Lese-Form-/Wortlaut-
  Korrektur, Nutzer-Entscheid 2026-06-22); ihre **Wirkung** auf die gepinnte
  Pointer-Liste wird hier nachgezogen.
- **Auflösungs-Trigger:** permanent (Provenienz).
