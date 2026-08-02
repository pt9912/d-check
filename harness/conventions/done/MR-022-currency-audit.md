# MR-022 — Baseline-Currency-Audit-Modus (Nachtrag zu MR-019)

- **Status:** Accepted
- **Datum:** 2026-07-19
- **Geltungsbereich:** [`tools/harness/fetch-baseline-cache.sh`](../../../tools/harness/fetch-baseline-cache.sh)
  (neuer Modus `--check-latest`), [`AGENTS.md`](../../../AGENTS.md) §1; Nachtrag zu
  [`MR-019`](../../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)
- **Adaption:** `fetch-baseline-cache.sh` erhält einen dritten Modus
  `--check-latest` mit **zwei** Upstream-Prüfungen (beide Netz, informativ):
  **(A) Currency** — der in [§Baseline](../../conventions.md#baseline) gepinnte Tag gegen das
  **neueste stabile** Release
  (`https://api.github.com/repos/pt9912/ai-harness-course/releases/latest`;
  GitHub blendet dort Prereleases/Drafts aus, passend zur Re-Adopt-Semantik).
  **(B) Content-Drift am gepinnten Tag** — das Skript lädt `lab-regelwerk.zip`
  des **gepinnten** Tags und vergleicht dessen Bytes (dasselbe
  `sha256sum regelwerk/*.md`-Manifest) gegen das committete
  `.harness/baseline/<tag>/SHA256SUMS`; eine Abweichung heißt, der Tag wurde
  **verschoben** oder das Asset neu hochgeladen. Ausgang (schlimmster Fall):
  aktuell & authentisch → `exit 0`; neuerer Tag (`sort -V`) → `exit 3`
  (**Signal**, kein Fehler); Content-Drift am gepinnten Tag → `exit 4`
  (**Provenienz-Alarm**); nicht erreichbare Teile (Netz/API/Rate-Limit,
  fehlendes Werkzeug/Manifest) → **SKIP** je Teil (`exit 0`, sofern der andere
  Teil nicht `3`/`4` meldet). Der Modus ist bewusst **kein Gate**
  (`--network`-abhängig wie das Re-Vendoring) und bewusst **nicht fail-closed**
  (Gegenstück zu `--verify`, das netzlose Integrität prüft und sehr wohl
  fail-closed ist): ein nicht erreichbares Upstream darf keinen Lauf blockieren.
- **Begründung:** Übernommen aus dem Kurs-Beispiel
  `lab/example/tools/check_regelwerk_drift.py` (inhaltsbasierter Drift-Sensor der
  adoptierten Form-Quelle), auf d-checks Tag-Pin-Modell übersetzt. d-check pinnt
  einen Release-Tag und vendored ihn hash-verifiziert (`--verify`,
  [`MR-019`](../../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017));
  `--verify` prüft aber nur die **vendorten** Bytes gegen ihr **eigenes**
  Manifest — nie gegen Upstream. Es blieben zwei tote Winkel, die
  `--check-latest` schließt: **Currency** (Teil A — liegt der Pin hinter
  Upstream? Ein grüner Integritäts-Check verdeckte das) und **Authentizität des
  gepinnten Tags** (Teil B — d-check *setzt* auf Tag-Immutabilität, aber »Tag
  verschoben / Asset neu« bemerkt `--verify` nicht). Teil B **verifiziert** genau
  die Immutabilitäts-Annahme aus
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt) (prüfen statt vertrauen)
  und ist der Content-Hash-Drift-Kern des Kurs-Beispiels, auf den gepinnten Tag
  angewandt. Beides **ohne** den Grundsatz aus
  [`MR-019`](../../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)
  zu unterlaufen (»Re-Vendor + neues Manifest sind ein bewusster Akt am
  Baseline-Pin-Bump, kein laufender Drift«): der Modus **automatisiert nichts**,
  er meldet nur — der Re-Adopt bleibt der bewusste manuelle Akt
  ([`MR-020`](../../conventions.md#mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt) /
  [`MR-021`](../../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden)).
- **Auflösungs-Trigger:** permanent, solange die Baseline extern gepinnt und
  lokal vendored wird.
