## Modul 14 — Docker Harness

*Quelle: [05-betrieb/modul-14-docker-harness.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/05-betrieb/modul-14-docker-harness.md)*

### Kernidee (Modul 14)

Wenn lokal und CI nicht dasselbe Image benutzen, debuggst du den
Unterschied, nicht den Bug.

### Regeln gegen typische Fehlannahmen (Modul 14)

- **"FROM python:3 ist konkret genug."** — Nein. Ohne Digest (`FROM python:3.12.4-slim@sha256:…`) baust du jeden Monat einen anderen Container.
- **"Lock-Files sind nur für Python."** — Lock-Files gibt es für jede Sprache: `package-lock.json`, `go.sum`, `Cargo.lock`, `packages.lock.json` (mit Central Package Management, siehe `bess-ems`), `pnpm-lock.yaml`, `poetry.lock`. Wer ohne Lock-File baut, baut nicht reproduzierbar.
- **"Docker-only ist Overkill für Tools."** — Tools driften am schnellsten. Genau dort lohnt Docker am meisten.
- **"Devcontainer ersetzt Compose."** — Nein. Devcontainer ist für *Entwickler-IDE-Setup*, Compose für *Lauf- und CI-Vertrag*. Sie ergänzen sich.
- **"DevOps ist YAML schreiben — Container = Deployment."** — Verbreitet, weil Container historisch über die Deployment-Seite eingeführt wurden. In diesem Kurs ist der primäre Zweck eines Containers ein anderer: er ist **Reproduzierbarkeits-Anker** — derselbe Image-Hash garantiert dieselbe Toolchain auf jeder Maschine, im CI und in sechs Monaten. Deployment ist *eine* Anwendung dieses Ankers, nicht sein Hauptzweck. Bei einem Replay-Lauf gegen ein altes Golden Set ([Modul 12](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-12-replay-evaluierung.md)) brauchst du den *Image-Hash von damals*, nicht das aktuelle Deployment. Wer das Bild "Container = Auslieferung" pflegt, hat keinen Hebel für *time-travel reproducibility* — und damit kein belastbares Replay.

### Worked Example: vom einstufigen Dockerfile zur reproduzierbaren Multi-Stage-Pipeline

**Ausgangs-Dockerfile (Python, Anti-Beispiel):**

```dockerfile
FROM python:3
COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt
CMD ["python", "-m", "docsearch"]
```

Vier Zeilen, vier Drift-Quellen: Tag `python:3` zeigt jeden Monat auf
ein anderes Image; `requirements.txt` ist nicht aufgelöst (transitive
Versionen frei); `pip install` ohne Cache-Trennung baut bei jedem Code-
Change die Dependencies neu; das Runtime-Image enthält den Build-Layer
mit Quellcode und Compiler-Toolchain. Sechs Schritte bringen das in
einen Multi-Stage-Build, der lokal und in CI denselben Image-Hash
produziert.

**Schritt 1 — Base-Image mit Digest pinnen.** Tag-Floating ist die
unsichtbarste Drift, weil sie nichts ändert *außer* dass das Image
neu ist. Lösung: SHA-256-Digest dazu.

```dockerfile
FROM python:3.12.4-slim@sha256:9c7f4a9d0c1b2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f AS deps
```

Der Digest wird beim ersten erfolgreichen Lokal-Build von
`docker buildx imagetools inspect python:3.12.4-slim` ausgelesen und
festgeschrieben. Update-Pfad: bei Sprach-/Sicherheits-Update Digest
*bewusst* anheben — ein Commit, der nur die Digest-Zeile ändert.

**Schritt 2 — Lock-File trennen und vor dem Code in den Build-Kontext
holen.** Damit Dependency-Installation cache-freundlich wird (sie läuft
neu *nur* wenn `pyproject.toml` / `poetry.lock` sich ändert, nicht bei
jedem Code-Change):

```dockerfile
FROM python:3.12.4-slim@sha256:9c7f4a... AS deps
WORKDIR /src
COPY pyproject.toml poetry.lock ./
RUN pip install --no-cache-dir uv==0.4.0 && \
    uv pip install --system --no-cache --frozen .
```

Drei Disziplinen in dieser Stage: (a) Installer-Version (`uv==0.4.0`)
selbst pinnen, sonst ist das Installations-Tool die zweite
Drift-Quelle; (b) `--frozen` verbietet, dass uv beim Build neue
Versionen auflöst — Lock-File entscheidet, nicht Build; (c) noch *kein*
Code im Image — Layer-Cache greift, solange Lock unverändert.

**Schritt 3 — Build-Stage separieren.** Code-Kompilierung gehört nicht
ins Runtime-Image; sie braucht aber die Dependencies aus Stage 1.

```dockerfile
FROM deps AS build
COPY . .
RUN python -m compileall src/
```

`FROM deps` referenziert die vorherige Stage — `build` erbt die
installierten Pakete, ohne sie neu zu installieren. `compileall`
ist hier symbolisch für jede Sprach-spezifische Build-Aktion
(Bytecode-Vorgenerierung, Asset-Build, Typ-Stubs). In Go wäre es
`go build`, in Java `mvn package`.

**Schritt 4 — Distroless-Runtime-Stage mit nonroot.** Das Runtime-Image
trägt nur das, was zur Laufzeit *gebraucht* wird — keine Shell, kein
Paketmanager, keine Build-Toolchain. Angriffsfläche minus ~90 %.

```dockerfile
FROM python:3.12.4-slim@sha256:9c7f4a... AS runtime
WORKDIR /app
COPY --from=build /src/src/docsearch /app/docsearch
COPY --from=build /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
RUN useradd --uid 65532 --no-create-home nonroot
USER nonroot
ENTRYPOINT ["python", "-m", "docsearch"]
```

Für Sprachen mit eigenständigem Binär-Output (Go, Rust, statisch
gelinkte JVM-AOT) ist die noch härtere Variante
`gcr.io/distroless/static-debian12:nonroot@sha256:d093aa3e…` (auch
Distroless-Tags floaten — deshalb per Digest gepinnt) ohne
interpretierbares Runtime — siehe
[`../../lab/example/go/Dockerfile`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/go/Dockerfile)
als Vorbild.

**Schritt 5 — Image-Hash im Build-Output festhalten.** Damit das Image
in einem Replay-Manifest (Modul 12) referenzierbar wird:

```makefile
build:  ## LH-QA-03 — reproduzierbarer Build, Image-Hash erfasst
	docker buildx build \
		--platform linux/amd64 \
		--tag docsearch:welle-2 \
		--metadata-file build-metadata.json \
		--load .
	@jq -r '."containerimage.config.digest"' build-metadata.json > harness/image-hash.txt
	@cat harness/image-hash.txt
```

`build-metadata.json` enthält den exakten Manifest-Digest. Die
`harness/image-hash.txt` ist ein einzeiliges Beleg-Artefakt, das in
`harness/README.md` referenziert wird (siehe Vorlage in
[`/lab/templates/harness/README.template.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/templates/harness/README.template.md)).
Ohne diesen Schritt ist das Replay-Manifest in Modul 12 zur Hälfte
blind — der `image_hash`-Slot bleibt unbelegt.

**Schritt 6 — Bewusstes Brechen: Drift provozieren.** Ändere in einer
Kopie *eine* Zeile zurück auf den unsicheren Stand und messe die
Wirkung:

| Änderung | Erwartete Beobachtung |
|---|---|
| Digest weglassen (`FROM python:3.12.4-slim`) | Image-Hash ändert sich beim nächsten Lokal-/CI-Build, obwohl kein Code-Diff vorliegt. |
| `--frozen` aus Schritt 2 entfernen | uv löst beim Build neue Patch-Versionen auf; Lock-File und Image divergieren stillschweigend. |
| `COPY . .` *vor* `COPY pyproject.toml ./` ziehen | Dependency-Stage wird bei jedem Code-Change rebuilt; Build-Zeit explodiert, Cache wirkt nicht. |

Sechs Schritte, ein Image, drei Drift-Anker (Digest · Lock-File ·
Stage-Trennung). Vergleich:
[`../../lab/example/python/Dockerfile`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/python/Dockerfile)
und
[`../../lab/example/go/Dockerfile`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/go/Dockerfile)
— beide tragen den ID-Kommentar `LH-QA-03` im Header und folgen
demselben Drei-Stage-Schnitt mit sprach-spezifischen Anpassungen.

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

