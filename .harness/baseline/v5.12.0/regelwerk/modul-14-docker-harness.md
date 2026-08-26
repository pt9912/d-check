## Modul 14 — Docker Harness

<!-- Quelle: [05-betrieb/modul-14-docker-harness.md](https://github.com/pt9912/ai-harness-course/blob/v5.12.0/kurs/de/05-betrieb/modul-14-docker-harness.md) -->

### Kernidee (Modul 14)

Wenn lokal und CI nicht dasselbe Image benutzen, debuggst du den
Unterschied, nicht den Bug.

### Regeln gegen typische Fehlannahmen (Modul 14)

- **"FROM python:3 ist konkret genug."** — Nein. Ohne Digest (`FROM python:3.12.4-slim@sha256:…`) baust du jeden Monat einen anderen Container.
- **"Lock-Files sind nur für Python."** — Lock-Files gibt es für jede Sprache: `package-lock.json`, `go.sum`, `Cargo.lock`, `packages.lock.json` (mit Central Package Management), `pnpm-lock.yaml`, `poetry.lock`. Wer ohne Lock-File baut, baut nicht reproduzierbar.
- **"Docker-only ist Overkill für Tools."** — Tools driften am schnellsten. Genau dort lohnt Docker am meisten.
- **"Devcontainer ersetzt Compose."** — Nein. Devcontainer ist für *Entwickler-IDE-Setup*, Compose für *Lauf- und CI-Vertrag*. Sie ergänzen sich.
- **"DevOps ist YAML schreiben — Container = Deployment."** — Verbreitet, weil Container historisch über die Deployment-Seite eingeführt wurden. In diesem Regelwerk ist der primäre Zweck eines Containers ein anderer: er ist **Reproduzierbarkeits-Anker** — derselbe Image-Hash garantiert dieselbe Toolchain auf jeder Maschine, im CI und in sechs Monaten. Deployment ist *eine* Anwendung dieses Ankers, nicht sein Hauptzweck. Bei einem Replay-Lauf gegen ein altes Golden Set ([Modul 12](modul-12-replay-evaluierung.md)) brauchst du den *Image-Hash von damals*, nicht das aktuelle Deployment. Wer das Bild "Container = Auslieferung" pflegt, hat keinen Hebel für *time-travel reproducibility* — und damit kein belastbares Replay.

### Multi-Stage-Build: die operativen Disziplinen (Modul 14)

Ein einstufiges `FROM python:3` / `COPY .` / `pip install`-Dockerfile hat
vier Drift-Quellen (floatender Tag · unaufgelöste Dependencies · kein
Cache-Schnitt · Build-Toolchain im Runtime-Image). Der Multi-Stage-Build,
der lokal und in CI denselben Image-Hash produziert, verlangt:

- **Base-Image per Digest pinnen** (`FROM …@sha256:…`), nicht per Tag —
  Tag-Floating ist die unsichtbarste Drift (ändert nichts *außer* dass
  das Image neu ist). Digest beim ersten Build aus
  `docker buildx imagetools inspect` auslesen; Update = *bewusster*
  Commit, der nur die Digest-Zeile anhebt.
- **Lock-File vor dem Code** in den Build-Kontext holen (Layer-Cache
  greift, solange das Lock unverändert ist); **Installer-Version selbst
  pinnen** (`uv==0.4.0` o. ä., sonst ist das Tool die zweite
  Drift-Quelle); **`--frozen`** verbietet Auflösung neuer Versionen beim
  Build — das Lock-File entscheidet, nicht der Build.
- **Stages trennen:** `deps` (gepinnte Base + Lock-Install) → `build`
  (`FROM deps`, Kompilierung getrennt vom Cache-sensiblen Layer) →
  `runtime` (Distroless/nonroot, nur Artefakte kopiert — keine Shell,
  kein Paketmanager, keine Build-Toolchain; Angriffsfläche minus ~90 %).
- **Image-Hash im Build-Output festhalten** (`docker buildx build
  --metadata-file …` → einzeiliges Beleg-Artefakt `harness/image-hash.txt`,
  referenziert in `harness/README.md`, Vorlage
  [`templates/harness/README.template.md`](../templates/harness/README.template.md)).
  Ohne ihn bleibt der `image_hash`-Slot des Replay-Manifests
  ([Modul 12](modul-12-replay-evaluierung.md)) blind — Modell-Drift lässt
  sich dann nicht von Toolchain-Drift trennen.

### Reproduzierbarkeits-Regeln: Drift-Klassen und Stage-Schnitte

- **Mindestkombination für Build-Reproduzierbarkeit:** Lock-File (sichert Abhängigkeits-Versionen) + Image-Hash (sichert Runtime-/Toolchain-Version). Ohne Lock-File driftet das Dependency-Tree, ohne Image-Hash driftet die Sprach-/Tool-Version. Folge: ein Replay-Manifest (Modul 12) referenziert *beide* — ohne Image-Hash lässt sich Modell-Drift nicht von Toolchain-Drift trennen; ohne Lock-File-Hash nicht von Dependency-Drift. Drei Drift-Quellen, drei Anker.
- **Drift-Klassen:** `FROM python:3` ⇒ Toolchain-Drift (Tag floatet, kein Digest); fehlendes `--frozen`/Lock-File ⇒ Dependency-Drift; `COPY . .` vor `pyproject.toml` ⇒ Layer-Cache-Drift (Cache invalidiert bei jedem Code-Change).
- **Drei Stage-Schnitte mit Härtung:** **deps** (gepinnte Base + Lock-File-Install gegen Toolchain-/Dependency-Drift) · **build** (`FROM deps`, Code-Kompilierung getrennt vom Cache-sensiblen Layer) · **runtime** (Distroless/nonroot, nur Artefakte kopiert — kleinere Angriffsfläche, kein Build-Layer im Image). Image-Hash macht den Schnitt erst messbar.
- **Warum `make gates` im Host-OS keine valide Gate-Ausführung ist:** Host-Toolchain ist nicht versionsgleich mit CI; Gate-Ergebnisse divergieren; Debugging erfolgt am Unterschied, nicht am Bug. Konsequenz: ohne Image-Hash-Vertrag zwischen lokal und CI sind grüne lokale Gates *kein* Vertrag — sie sind eine private Information.

### Devcontainer/Compose-Kriterium

Devcontainer für IDE-Setup (Sprache-Server, Debugger-Anschluss). Compose
für Lauf- und CI-Vertrag. Beides parallel, wenn das Team mehrere IDEs
nutzt. Faustregel: Compose ist *Pflicht* (CI-Vertrag), Devcontainer ist
*Komfort*. Wer mit Devcontainer beginnt, baut sich eine zweite Toolchain
ohne die erste.

