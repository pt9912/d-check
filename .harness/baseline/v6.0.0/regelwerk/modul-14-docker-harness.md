## Modul 14 — Docker Harness

<!-- Quelle: [05-betrieb/modul-14-docker-harness.md](https://github.com/pt9912/ai-harness-course/blob/v6.0.0/kurs/de/05-betrieb/modul-14-docker-harness.md) -->

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

### Zwei Formen des Reproduzierbarkeits-Ankers

Der Hash adressiert *ein* Image; er trägt, solange dieses Image noch existiert
— in einer Registry oder im lokalen Cache. Ein Tag, den jeder Build
überschreibt, ist keine Adresse, sondern eine Notiz.

| Form | Was den Lauf adressiert | Bedingung |
|---|---|---|
| **Archiv** | der Digest des gebauten Images | das Image wird aufbewahrt — gepusht oder nachweislich vorgehalten |
| **Rezept** | Commit der Quellen **und** die gepinnten Digests der Eingangs-Images | beim Build wird nichts installiert, jede Eingabe ist digest-gepinnt |

Die Rezept-Form braucht kein Lager: Wer den Commit hat, baut die Umgebung neu.
Ihr Preis ist, dass der Digest des *gebauten* Images nicht ihr Griff sein kann
— Datei-Zeitstempel wandern in die `COPY`-Schichten, ein frischer Checkout
stempelt sie neu, und die Normalisierungs-Schalter der Build-Werkzeuge räumen
das nicht zuverlässig aus. `harness/image-hash.txt` hält dann fest, *welches*
Image einen Lauf gemacht hat; ein Wiederholungs-Schlüssel ist es nicht. Wer die
Archiv-Form will, muss das Image aufbewahren — wer das nicht tut, hat die
Rezept-Form und benennt sie besser auch so.

### Besitz der Belege eines containerisierten Gates

**`nonroot` endet nicht am Runtime-Image.** Den Arbeitsbaum auf deiner Platte
berührt die Toolchain-Stage, in der die Gates laufen. Läuft sie als root über
einem beschreibbaren Mount (`-v "$(pwd)":/src`), gehören Build-Verzeichnisse,
Coverage-Dateien und Werkzeug-Caches danach `root:root` — der Mensch, der den
Lauf angefordert hat, kann sie nicht einmal löschen. **Wem gehören die Belege,
die ein containerisierter Gate schreibt?** Wer sie nicht beantwortet,
beantwortet sie faktisch mit *root*.

Die Antwort dieses Regelwerks steht im nächsten Abschnitt: **kein Mount**. Wer
trotzdem mountet, wählt zwischen zwei Preisen — `:ro` plus Umleitung alles
Schreibenden (trägt nur, solange die Prüfung nichts in den Baum schreiben
*muss*) oder `--user $(id -u):$(id -g)` (hebt die `nonroot`-Zusage des Images
zur Laufzeit auf und verlangt Werkzeuge, die ohne eigenes `HOME` auskommen).

Der Testfall kostet nichts: `ls -l` auf das Build-Verzeichnis nach dem ersten
Gate-Lauf.

### Der Prüflauf ist hermetisch — kein Mount

Ein Bind-Mount reicht dem Container den Arbeitsbaum **von außen** herein. Das
kostet drei Dinge: Der Prüfgegenstand ist nicht Teil des Images, also sagt kein
Digest, *was* geprüft wurde; ein schreibender Lauf verändert den Baum, und die
Besitzfrage stellt sich überhaupt erst; und Docker legt Mountpunkte host-seitig
als root an — auch ein `--tmpfs` hinterlässt eine root-eigene Hülle, die keine
UID-Option wegräumt.

**Die Quellen wandern beim Build ins Image, die Ergebnisse kommen über `stdout`
heraus.**

**Zwei Wege, die Prüfung auszulösen — und was der zweite braucht.** Liegt das
Werkzeug im Image, ruft man es per `docker run` auf: jeder Lauf prüft frisch.
Zieht das Werkzeug seine Abhängigkeiten dagegen beim Build (Maven, Gradle,
NuGet), ist die **Gate-Stage selbst** das Gate — `docker build --target
lint-gate`. Dann aber zwei Griffe, sonst ist es kein Gate mehr:

| Griff | Warum |
|---|---|
| `--no-cache-filter <stage>` | führt genau diese Stage neu aus, während `repo` und die Werkzeug-Layer gecacht bleiben. Ohne ihn wiederholt ein gecachtes Grün nur, dass *dieser Stand* schon einmal durchlief — der Gate urteilt nicht, er erinnert sich. |
| **kein** `-q` | mit `-q` zeigt ein roter Gate nur *„exit code: 1"*. Die Befunde des Werkzeugs stehen im Build-Log, und genau die braucht der, der sie beheben soll. |

| Eigenschaft | Wirkung |
|---|---|
| Quellen per `COPY`, kein Mount | der Baum ist während des Laufs unerreichbar — ein Gate *kann* ihn nicht ändern |
| Rückweg über `stdout`, host-seitig ausgepackt | der schreibende Prozess ist der Host: die Belege gehören dem, der sie angefordert hat |
| Gate-Stage und Beleg-Stage getrennt | das Urteil bricht den Build, der Report entsteht trotzdem ([Modul 13](modul-13-quality-gates.md#gate-und-beleg--zwei-rollen-derselben-prüfung)) |
| `export` erbt von `repo` | ein roter Gate macht das Werkzeug nicht unbaubar, mit dem man ihn untersucht |

**Der Preis, offen benannt:** Jede Quelländerung verlangt einen Rebuild — der
Layer-Cache trägt ihn, aber ein Lauf gegen ungespeicherte Änderungen ist nicht
möglich. Und der Rückweg löst **Ausgaben**, nicht Eingaben: Ein erneuertes
Lock-File ist eine Quelle und muss zurück in den Baum.

**Die Ausnahme, und ihre Bedingung:** Ein Werkzeug, das seine Recipe-Form
selbst mitbringt, mountet in der Regel read-only. Das ist zulässig, solange es
**nur liest** — der Besitz-Schaden entsteht am Schreiben. Wer es hermetisch
will, bindet das Fragment nicht ein und schreibt die Recipe aus; das ist eine
Abweichung von der Werkzeug-Form und gehört als `MR-<NNN>` deklariert.

### Devcontainer/Compose-Kriterium

Devcontainer für IDE-Setup (Sprache-Server, Debugger-Anschluss). Compose
für Lauf- und CI-Vertrag. Beides parallel, wenn das Team mehrere IDEs
nutzt. Faustregel: Compose ist *Pflicht* (CI-Vertrag), Devcontainer ist
*Komfort*. Wer mit Devcontainer beginnt, baut sich eine zweite Toolchain
ohne die erste.

