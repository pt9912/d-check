# MR-019 — Regelwerk-Lese-Form committet statt gecacht (Nachtrag zu MR-017)

- **Status:** Accepted
- **Datum:** 2026-06-26
- **Geltungsbereich:** [§Adoptierte Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen)
  (Lokale Lese-Form), [`.d-check.yml`](../../../.d-check.yml) `scan.ignore`, das Skript
  [`tools/harness/fetch-baseline-cache.sh`](../../../tools/harness/fetch-baseline-cache.sh),
  [`AGENTS.md`](../../../AGENTS.md) §1, der neu committete Pfad `.harness/baseline/<tag>/`;
  Nachtrag zu
  [`MR-017`](../../conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)
- **Adaption:** Die **Regelwerk**-Lese-Form der adoptierten Baseline wird nicht
  mehr nur im gitignorierten Cache materialisiert, sondern **committet vendored**
  nach `.harness/baseline/<tag>/regelwerk/` (entpacktes `lab-regelwerk.zip`),
  zusammen mit einem committeten `.harness/baseline/<tag>/SHA256SUMS` über die
  vendorten Dateien; aktueller `<tag>` = `v1.4.0`.
  `tools/harness/fetch-baseline-cache.sh` schreibt das Regelwerk in diesen
  Vendor-Pfad, (re)generiert SHA256SUMS und verifiziert; der `--verify`-Modus
  prüft das committete Regelwerk **netzlos** gegen das Manifest
  (CI/Audit/frischer Checkout). Der Vendor-Pfad ist über
  `scan.ignore: [".harness/baseline/**", …]` in [`.d-check.yml`](../../../.d-check.yml)
  aus dem Dogfooding-Selbst-Scan ausgenommen — **dieselbe** Begründung wie für den
  Cache in
  [`MR-017`](../../conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen):
  Fremdinhalt (die Kurs-Docs referenzieren eigene `ADR-`/`MR-`-IDs und Modulpfade,
  die in *diesem* Repo nicht existieren); **keine** Gate-Lockerung im Sinne von
  [`AGENTS.md` §3.6](../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden),
  da Nicht-Repo-Inhalt ausgenommen wird, keine Repo-Doku Deckung verliert. Das
  **Template**-Set bleibt unverändert ephemerer Cache unter
  `.harness/cache/<tag>/templates/` (nur Adoptions-/Drift-Audit-Staging,
  [`MR-018`](../../conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)) —
  template-frei wird nur das Regelwerk vendored, nicht die Templates.
  [`MR-017`](../../conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)
  bleibt als Provenienz stehen; seine Cache-Aussage gilt fortan nur noch für die
  Templates.
- **Begründung:** Sichtbar geworden über die Steering-Loop-Verweis-Frage
  (2026-06-26): ein Pointer auf die kanonische Regelwerk-Definition (z. B.
  `grundlagen-klassifikation.md` §Steering Loop) löste in keiner Zielumgebung
  in-repo auf, weil die Lese-Form ephemer/gitignored war — auf frischem Checkout
  oder ohne Netz schlicht abwesend. Schärfer noch: d-check pinnt seine eigenen
  Release-Digests (`sha256:…`) und erkennt mit dem `pins`-Modul inhaltlichen Drift
  verlinkter Quellen, konsumierte aber seine **eigene** Baseline-Lese-Form per
  `curl` von einem Release-Asset **ohne Content-Hash** — der Pin hielt den Tag,
  nicht die Bytes. Das Vendoring schließt drei Lücken in einem: Präsenz (jeder
  Checkout hat das Regelwerk), Offline-Auditierbarkeit und Integrität/Provenienz
  (Bytes via SHA256SUMS + git-Historie statt unverifiziertem Netz-Fetch). Es
  bestätigt das Verkörperungs-Prinzip (Kurs Modul 0: Per-Lauf-/Schwellen-
  Relevantes gehört verkörpert, nicht extern nachgeladen) und stützt
  [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  (netzlose Lese-Form). Re-Vendor + neues Manifest sind ein bewusster Akt am
  Baseline-Pin-Bump, kein laufender Drift.
- **Auflösungs-Trigger:** permanent, solange die Baseline-Regelwerk-Lese-Form
  lokal vendored wird; der nächste Baseline-Pin-Bump re-vendored und erneuert
  SHA256SUMS.
