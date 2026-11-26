## Modul 14 — Docker Harness

<!-- Quelle: [05-betrieb/modul-14-docker-harness.md](https://github.com/pt9912/ai-harness-course/blob/v7.0.0/kurs/de/05-betrieb/modul-14-docker-harness.md) -->

### Kernidee (Modul 14)

Wenn lokal und CI nicht dasselbe Image benutzen, debuggst du den
Unterschied, nicht den Bug.

### Die Umgebungs-Zusage — Anforderung vor Träger

**Die Anforderung hängt nicht an Docker.** Hinter „lokal und CI benutzen
dasselbe Image“ steht eine Anforderung, die auch ohne Container gilt: die
**Umgebungs-Zusage**. *Die Umgebung, in der ein Lauf stattfindet, ist die
behauptete — und das lässt sich belegen.* Sie hat zwei Teile: Die Umgebung ist
**festgelegt** (gepinnt), und sie ist **adressierbar** (ein Beleg sagt, welche es
war). Die Zusage ist die Aussage; ein Sensor ist ihre **Prüfung**, nicht ihr
Bestandteil. Dass ein Projekt eine Umgebungs-Zusage braucht, steht als `LH-QA-*`
in seinem Lastenheft, mit dem Sensor als Messmethode — der Kurs lehrt die
Klasse, das Projekt trägt die Instanz. Das Präfix *Umgebungs-* steht immer mit:
Die Zusage des Lastenhefts an den Auftraggeber
([`source-precedence.md` §Spec-Stratifizierung](grundlagen-source-precedence.md#spec-stratifizierung))
ist etwas anderes.

**Welche Frage ein Pin beantwortet.** Die Pins dieses Kurses stehen über sechs
Module verteilt. Geordnet nach der Frage, nicht nach dem Werkzeug:

| Frage | Klasse | Träger (Beispiele) |
|---|---|---|
| *Welche Bytes laufen?* | Artefakt-Adresse | Image-Digest, Release-`SHA256SUMS`, vendorte Baseline ([Modul 2](modul-02-harness-bootstrap.md)) |
| *Welche Version von was?* | Versions-Pin | Lock-File, Installer-Pin, Base-Image-Digest, Tool-Pin |
| *Welche Eingaben und welcher Zufall?* | Lauf-Pin | Replay-Manifest: `inputs_ref`, `seed`, `determinism:` ([Modul 12](modul-12-replay-evaluierung.md)) |
| *Welches Modell, welcher Prompt?* | Modell-Pin | `model.version` (ohne Alias), `prompt_context`-Hash |
| *Welcher Zustand beim Start?* | Umgebungs-Pin | `runtime.env`, Zeitzone, Locale |
| *Hält der Pin noch?* | **Prüfung der Zusage** | Vergleich Soll ↔ Ist |

Die letzte Zeile ist kein weiterer Träger, sondern die Prüfung, dass die
anderen noch tragen. Bei den meisten Pins fehlt sie: Der Image-Digest wird
festgehalten (`harness/image-hash.txt`), aber nichts vergleicht ihn mit dem, was
gerade läuft.

**Zwei Träger.** Die Umgebungs-Zusage lässt sich auf zwei Wegen erfüllen: mit
Docker und Dockerfile, die die Umgebung aus gepinnten Eingaben bauen, oder mit
Tool-Pins und einem Tool-/Versions-Check, der feststellt, ob die vorhandene
Umgebung sie trägt. Das ist
der heutige Stand des Kurses, keine Abwägung zwischen beiden oder gegen weitere
Träger — Nix, devbox, asdf, mise oder eine VM behandelt der Kurs nicht. Die Wahl
trifft das Repo; es kann sie als ADR festhalten, mit den betrachteten Optionen
([Modul 4](modul-04-adrs.md)).

### Regeln gegen typische Fehlannahmen (Modul 14)

- **"FROM python:3 ist konkret genug."** — Nein. Ohne Digest (`FROM python:3.12.4-slim@sha256:…`) baust du jeden Monat einen anderen Container.
- **"Lock-Files sind nur für Python."** — Lock-Files gibt es für jede Sprache: `package-lock.json`, `go.sum`, `Cargo.lock`, `packages.lock.json` (mit Central Package Management), `pnpm-lock.yaml`, `poetry.lock`. Wer ohne Lock-File baut, baut nicht reproduzierbar.
- **"Docker-only ist Overkill für Tools."** — Tools driften am schnellsten. Genau dort lohnt eine Umgebungs-Zusage am meisten — Docker ist einer ihrer zwei Träger, nicht der einzige.
- **"Devcontainer ersetzt Compose."** — Nein. Devcontainer ist für *Entwickler-IDE-Setup*, Compose für *Lauf- und CI-Vertrag*. Sie ergänzen sich.
- **"Reproduzierbar heißt: im Container."** — Der Container ist ein Träger der Umgebungs-Zusage, nicht die Zusage. Ein Repo ohne Container trägt sie mit Tool-Pin und Tool-/Versions-Check — und im Container prüft derselbe Check, ob darin steckt, was der Pin behauptet.
- **"DevOps ist YAML schreiben — Container = Deployment."** — Verbreitet, weil Container historisch über die Deployment-Seite eingeführt wurden. In diesem Regelwerk ist der primäre Zweck eines Containers ein anderer: er ist ein **Träger der Umgebungs-Zusage** — derselbe Image-Hash garantiert dieselbe Toolchain auf jeder Maschine, im CI und in sechs Monaten. Deployment ist *eine* Anwendung dieses Trägers, nicht sein Hauptzweck. Bei einem Replay-Lauf gegen ein altes Golden Set ([Modul 12](modul-12-replay-evaluierung.md)) brauchst du den *Image-Hash von damals*, nicht das aktuelle Deployment. Wer das Bild "Container = Auslieferung" pflegt, hat keinen Hebel für *time-travel reproducibility* — und damit kein belastbares Replay.

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

- **Mindestkombination für Build-Reproduzierbarkeit:** Lock-File (sichert Abhängigkeits-Versionen) + Image-Hash (sichert Runtime-/Toolchain-Version). Ohne Lock-File driftet das Dependency-Tree, ohne Image-Hash driftet die Sprach-/Tool-Version. Folge: ein Replay-Manifest (Modul 12) referenziert *beide* — ohne Image-Hash lässt sich Modell-Drift nicht von Toolchain-Drift trennen; ohne Lock-File-Hash nicht von Dependency-Drift. Drei Drift-Quellen, drei Pins.
- **Drift-Klassen:** `FROM python:3` ⇒ Toolchain-Drift (Tag floatet, kein Digest); fehlendes `--frozen`/Lock-File ⇒ Dependency-Drift; `COPY . .` vor `pyproject.toml` ⇒ Layer-Cache-Drift (Cache invalidiert bei jedem Code-Change).
- **Drei Stage-Schnitte mit Härtung:** **deps** (gepinnte Base + Lock-File-Install gegen Toolchain-/Dependency-Drift) · **build** (`FROM deps`, Code-Kompilierung getrennt vom Cache-sensiblen Layer) · **runtime** (Distroless/nonroot, nur Artefakte kopiert — kleinere Angriffsfläche, kein Build-Layer im Image). Image-Hash macht den Schnitt erst messbar.
- **Warum `make gates` im Host-OS ohne Umgebungs-Zusage keine valide Gate-Ausführung ist:** Host-Toolchain ist nicht versionsgleich mit CI; Gate-Ergebnisse divergieren; Debugging erfolgt am Unterschied, nicht am Bug. Konsequenz: ohne Umgebungs-Zusage zwischen lokal und CI sind grüne lokale Gates *kein* Vertrag — sie sind eine private Information. Ein Tool-Pin mit Tool-/Versions-Check macht den Host-Lauf belegbar (Prüfprotokoll Werkzeug · Soll · Ist), isoliert ihn aber nicht.

### Zwei Formen des Belegs beim Docker-Träger

Der Hash adressiert *ein* Image; er trägt, solange dieses Image noch existiert
— in einer Registry oder im lokalen Cache. Ein Tag, den jeder Build
überschreibt, ist keine Adresse, sondern eine Notiz.

| Form | Was den Lauf adressiert | Bedingung |
| --- | --- | --- |
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
| --- | --- |
| `--no-cache-filter <stage>` | führt genau diese Stage neu aus, während `repo` und die Werkzeug-Layer gecacht bleiben. Ohne ihn wiederholt ein gecachtes Grün nur, dass *dieser Stand* schon einmal durchlief — der Gate urteilt nicht, er erinnert sich. |
| **kein** `-q` | mit `-q` zeigt ein roter Gate nur *„exit code: 1"*. Die Befunde des Werkzeugs stehen im Build-Log, und genau die braucht der, der sie beheben soll. |

| Eigenschaft | Wirkung |
| --- | --- |
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

### Der zweite Träger: Tool-Pin und Tool-/Versions-Check

**Pin legt fest, Check belegt.** Der Check allein erfüllt nur den Teil
*adressierbar*: Er stellt fest, welche Umgebung lief. *Festgelegt* ist sie erst
durch einen **Tool-Pin** — die Versionen der Werkzeuge sind deklariert oder
gelockt (Manifest, Lock-File, Versionsmanager). Das Paar trägt beide Teile der
Zusage, auf einem anderen Weg als Docker:

| | **Docker / Dockerfile** | **Tool-Pin + Tool-/Versions-Check** |
|---|---|---|
| Wirkungsweise | **konstruktiv** — baut die Umgebung aus gepinnten Eingaben | **feststellend** — prüft, ob die vorhandene Umgebung den Pin trägt |
| Festlegung durch | Digest, Lock-File, Stages | Tool-Pin |
| Beleg | Image-Digest, `harness/image-hash.txt` | Prüfprotokoll: Werkzeug, Soll, Ist |
| Isolation | ja (eigenes Dateisystem; kein Netz nur mit `--network none`) | nein — die Umgebung bleibt die des Hosts |
| Schwäche | Docker am Host ist selbst ungepinnt; ein Pin kann etwas anderes installieren, als er behauptet | prüft nur, was das Soll nennt; Drift **zwischen** zwei Prüfungen bleibt unbemerkt |
| Trägt, wenn … | das Repo Docker nutzen darf und die Toolchain im Image liegt | kein Container möglich ist (Embedded, Host-gebundene Toolchain, Agent-Umgebung) oder der Container **geprüft** werden soll |

**Sie schließen sich nicht aus.** Wo Docker festlegt, ist der Check die
**Prüfung** der Zusage — der Digest wird festgehalten, aber erst der Check
vergleicht die Werkzeuge darin mit dem Pin. Wo Docker nicht trägt, ist der Check **selbst der Träger**. Derselbe
Mechanismus hat zwei Rollen. Auch im Docker-only-Repo prüft er etwas, das der
Digest nicht sagt:

| Kontext | Was der Digest sagt | Was nur der Check fängt |
|---|---|---|
| **Host** (`docker`, `make`, `bash`, `git`, `sed`, `awk`, `grep`) | nichts — der Host liegt vor dem Image und ruft es auf | ein Werkzeug fehlt oder ist zu alt. Host-Werkzeuge sind **nicht versionsgleich** mit dem Image: macOS liefert `bash` 3.2 und `make` 3.81, von `sed`, `awk` und `grep` gibt es GNU-, BSD-, BusyBox- und weitere Varianten mit verschiedenem Verhalten |
| **Im Image** | *dass* es dieses Image ist | dass darin die Version steckt, die der Pin behauptet |
| **Repo ohne Container** | nichts — es gibt keinen Digest | der Check **ist** der Träger; ohne ihn gibt es keinen |

#### Zwei Prüfungen derselben Zusage

- **Stellen-Abgleich (statisch):** Derselbe Pin steht oft an mehreren Stellen —
  `go.mod` und `Dockerfile`, das Image-Tag im `Makefile` und im CI-Workflow. Die
  Stellen müssen übereinstimmen. Dafür braucht es kein einziges Werkzeug, nur den
  Repo-Text. Mit d-check ist das das Modul `pinsync` (ein Pin, mindestens zwei
  Quellen, je Quelle ein Ausdruck mit genau einer Fanggruppe; byte-genauer
  Vergleich).
- **Laufzeit-Gate (dynamisch):** Die installierten Werkzeuge tragen den Pin.
  Das läuft im Kontext, den es prüft — ein Skript des Repos. Ein Doku-Prüfer wie
  d-check startet keine Programme und sähe im eigenen Container nur sein eigenes
  Image.

#### Woher das Soll kommt

Das Gate vergleicht Ist gegen Soll. Steht das Soll in einer **zweiten** Datei neben
dem Pin, driften beide auseinander, und das Gate prüft seine eigene Abschrift.
Zwei Wege, beide zulässig:

- **Ableiten:** Das Soll steht in der Pin-Stelle (`Dockerfile`, Lock-File,
  `go.mod`), das Gate liest sie. Eine Quelle — aber nur, wo eine maschinell
  lesbare Pin-Stelle existiert. Für Host-Werkzeuge (`bash` ≥ …, `make` ≥ …) gibt es
  keine.
- **Manifest:** eine eigene Deklaration je Werkzeug. Sprachneutral und
  vollständig deklarierbar, aber eine zweite Quelle; sie driftet gegen die
  Pin-Stelle, solange kein Stellen-Abgleich beide gegeneinander hält.

Die **Mischform** leitet ab, wo eine Pin-Stelle existiert, und führt den Rest im
Manifest. Die Kriterien: Gibt es eine Pin-Stelle? Läuft ein Teil am Host? Steht
ein Stellen-Abgleich zur Verfügung? Die Wahl ist eine Adaption der Baseline und
steht als `MR-<NNN>` in `harness/conventions.md`; der Kurs gibt keinen Weg vor.

#### Prüfform je Werkzeug — kein Werkzeug bleibt ungeprüft

Jedes Werkzeug, das ein Skript des Repos aufruft, steht im **Inventar** und trägt
eine von vier Prüfformen; sie lassen sich kombinieren:

| Prüfform | Wann | Beispiel |
|---|---|---|
| **exakt** | das Repo **legt** das Werkzeug fest (Pin, Image, Lock) — jede Abweichung ist Drift | `go`, `golangci-lint` |
| **Mindestversion** | das Repo kann es nur **fordern** (Host) — exakt wäre auf einem fremden Host ein Fehlalarm | `make` ≥ … |
| **Fähigkeit** | Implementierungen unterscheiden sich; geprüft wird, was das Skript braucht, nicht die Nummer | `bash`: `declare -A`; `sed`: die Option, die das Skript nutzt |
| **Anwesenheit** | jede Version genügt | `git`, `curl` |

Die Hauptachse ist **festlegen oder fordern**. Die Fähigkeitsprobe ist der
einzige tragfähige Weg, eine Variante zu erkennen: Das Banner trägt nicht —
BusyBox-`sed` meldet „This is **not GNU** sed“ und besteht einen Test auf „GNU“.
Jeder Eintrag trägt einen **Grund**, besonders bei *Anwesenheit* („jede Version
genügt, weil …“); ohne ihn wird sie zur stillen Ausnahme. Eine Pauschale wie
„Coreutils bleiben ungeprüft“ gibt es nicht — `sed` und `awk` sind nicht einmal
Coreutils.

**Vergleichsregeln.** Versionen werden **segmentweise und numerisch** verglichen,
nie als Zeichenkette oder Dezimalzahl, und *exakt* heißt die ganze Zeichenkette,
nie ein Teilstring. Zwei Fehlerbilder, die richtig aussehen:

| Prüfung | Fehler |
|---|---|
| `case "$ist" in *"$soll"*)` | Teilstring: `go1.27.30` passt auf den Pin `1.27.3` — falsches Grün |
| Mindestversion per `awk` als Dezimalzahl | `4.10 < 4.4` ist dort wahr — falsches Rot |

#### Inventar, Stufen, Kontext

Ein Eintrag nennt **Werkzeug · Prüfform · Soll · Kontext · Grund**.

**Das Inventar ist eine Datei des Repos, die das Skript liest** — keine Tabelle in
der Sensor-Datei. Beim Weg *Manifest* ist es das Manifest selbst, und das Feld
*Soll* trägt den Wert. Beim Weg *Ableiten* nennt das Feld die Pin-Stelle statt
ihres Werts; schriebe es den Wert ab, wäre es die zweite Quelle aus §Woher das
Soll kommt. Für ein Werkzeug ohne Pin-Stelle (`bash` ≥ …, `make` ≥ …) trägt das
Inventar den Wert selbst — dort gilt Manifest oder Mischform. Die **Regel** des Checks — welche Prüfformen,
wie verglichen wird, was rot ist und was Warnung — legt das Repo in seiner
Spezifikation fest; das Inventar trägt die **Daten**, die das Skript nach dieser
Regel liest. Die Sensor-Datei verlinkt die Festlegung, nennt den Pfad des
Inventars und sagt, wie ein Lauf zu lesen ist ([`grundlagen-harness-dateien.md`](grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt)).

- **Rot ist der Standard.** Eine **Warnung** gibt es nur für einen Eintrag, der
  ausdrücklich `optional` heißt, mit Grund und der Angabe, **was ausfällt**
  („Tests mit echtem `sshd` werden übersprungen“). Eine Warnung ohne diese
  Angabe ist ein stiller Ausgang.
- **Kontext.** Der Check prüft, was dort läuft, wo er läuft. Wo das ist, legt der
  Kurs nicht fest; jeder Weg hat seine Lücke, und das Repo deklariert seine Wahl
  in derselben `MR-<NNN>`:

| Ausführungsort | Lücke |
|---|---|
| nur Host (Docker-only, dem Image wird über den Digest vertraut) | die Werkzeuge im Image sind adressiert, nicht verglichen |
| nur Image (erster Schritt im Image) | der Host bleibt ungeprüft; ein zu altes `make` zeigt sich erst beim Einlesen des Makefiles |
| Host und Image, zwei Läufe | keine der beiden — zwei Instanzen desselben Musters, je mit eigenem Kontext |
| ein Kontext ohne Docker | keine — der Host ist alles |
| zusätzlich der CI-Runner | — das Runner-Image ist ein weiterer Kontext |

#### Wo der Check läuft — und was das Skript leisten muss

Der Check ist der **erste Check in `make gates`**. Ein Rot dort ist ein
Umgebungs-Befund, kein Befund über den Code — die Diagnose-Reihenfolge aus
[Modul 12](modul-12-replay-evaluierung.md) (Rang 1:
Toolchain-Drift) als Vorabprüfung. Die Grenze dieses Orts: Ein zu altes `make`
scheitert schon beim Einlesen des Makefiles, bevor der Check dran ist; der Check
prüft nur, was *trotzdem* startete.

Das Laufzeit-Gate ist ein **Skript des Repos**; die Vorlage liefert das Muster,
kein Skript. Was das Skript leisten muss:

1. **Werkzeug fehlt → rot. Version weicht ab → rot. Alles stimmt → Exit 0 ohne
   Befund.** Rot nennt Werkzeug, Soll und Ist.
2. **Es prüft seine eigene Bash zuerst** — nur mit Builtins (`BASH_VERSINFO`) und
   in Syntax, die auch eine ältere Bash als die zugelassene noch einliest. Der
   Guard steht **vor jeder Funktionsdefinition und außerhalb jeder Funktion**:
   Bash liest Top-Level-Befehle einzeln, aber eine Funktion ganz. Steht davor eine
   Funktion mit neuerer Syntax — oder alles in `main() { … }` —, scheitert eine
   alte Bash am Syntaxfehler statt an der Diagnose. Die Untergrenze legt das Repo
   über die Prüfform *Mindestversion* fest.
3. **Seine eigenen Werkzeuge stehen im Inventar** — `sed`, `awk`, `grep`, mit denen
   es Pins parst. Empfohlen: Die ersten Proben nutzen nur Builtins (`command -v`,
   `case`, `[`); erst danach parst das Skript — sonst prüft man `sed`, indem man
   `sed` benutzt.
4. **Der Pin des Prüfwerkzeugs** (bei d-check `DCHECK_DIGEST`) gehört zum
   Stellen-Abgleich.
5. **Der Break-Test ist Pflicht:** ein Werkzeug entfernen, eine Version
   verschieben, eine Bash unterhalb der Mindestversion starten — jedes muss rot
   werden, Letzteres mit Diagnose; unverändert grün. Eine Prüfung, die richtig
   aussieht, ist keine geprüfte Prüfung. Sein Protokoll steht im Skriptkopf oder
   in der ADR des Checks, nicht in der Sensor-Datei: Es belegt, dass das Werkzeug
   richtig ist, und sagt nichts darüber, wie ein Lauf zu lesen ist.

**Einordnung.** Der Check ist ein **Gate** (computational feedback), keiner der
drei Bindepunkte der
[Durchsetzungsschicht](grundlagen-durchsetzungsschicht.md#drei-bindepunkte);
das Handoff-Gate prüft, dass die Gates liefen, und deckt ihn damit ab. Ein
Tool-Call-Gate wie `MR-004` ([Modul 13](modul-13-quality-gates.md#guard-härtung-wächter-reifen-in-wellen-modul-13))
verhindert den Aufruf am falschen Werkzeug; der Check stellt die Umgebung fest.
Er ergänzt `MR-004`, er ersetzt es nicht.

#### Ist das Inventar vollständig? Die Probe im leeren Kontext

Handgepflegt driftet das Inventar: Ein Skript ruft nach einem Jahr `jq` auf, das
nirgends steht, und das Gate bleibt grün. Den Beleg liefert eine **Probe**: `make
gates` läuft in einem Kontext, der genau das Inventar enthält — träger-neutral als
**`PATH`-Käfig**, ein Wegwerf-Verzeichnis mit Verweisen nur auf die deklarierten
Werkzeuge, `PATH` darauf gesetzt. Ein undeklariertes Werkzeug lässt den Lauf
scheitern, **sofern der Fehler bis zu `make` durchkommt**: Die Shell meldet es als
nicht gefunden, der Befehl endet mit 127, `make` bricht mit „Error 127“ ab und
endet selbst mit Exit 2.

- **Beleg mit Lücke, kein Beweis.** Absolute Pfade umgehen den Käfig, auch
  `#!/bin/bash` und die Rezept-Shell von `make` (`/bin/sh`); Builtins liegen nicht
  im `PATH`; was im Lauf nicht ausgeführt wird, ist nicht geprüft. Was der Käfig
  nicht fangen kann, steht trotzdem im Inventar, von Hand.
- **Kein stiller Fallback.** Ein Zweig `command -v awk || …` weicht im Käfig still
  aus. Nutzt ein Skript ein Werkzeug optional, führt das Inventar es als
  `optional` (oben), und das Skript meldet den Ausfall sichtbar. Ebenso still
  endet ein Skript, das Fehler nicht weiterreicht: Ohne `set -e` meldet es das
  fehlende Werkzeug, läuft weiter und endet grün — in einer Pipe ebenso ohne
  `pipefail`. Die Skripte reichen Fehler weiter, sonst belegt die Probe nichts.
  Beides sind Regeln für das Skript, geprüft im Review — kein Sensor fängt sie.
- **Werkzeuge über `PATH` aufrufen, nicht absolut** — daraus `#!/usr/bin/env bash`
  statt `#!/bin/bash`; eine neuere Bash liegt am Mac ohnehin nicht unter `/bin`.
  Eine Konvention, geprüft im Review.
- **Auf Anlass, nicht vor jedem Commit:** wenn das Inventar sich ändert, ein
  Skript dazukommt, und in CI — etwa als eigenes Target `make gates-probe`.
- **Auch im Image.** Ein übliches Basis-Image bringt Werkzeuge mit, die niemand
  deklariert hat (`debian:stable-slim`: `awk`, `tar`, `find`, `perl`). Ohne Käfig
  läuft ein Skript, das `awk` aufruft, dort grün — das Image allein ist kein
  Inventar-Beleg.
- **Je Kontext.** Die Probe belegt nur, was im geprüften Kontext und in den
  geprüften Targets lief; das Feld *Kontext* sagt, welcher Lauf welchen Teil des
  Inventars belegt.

**Ohne Probe** bleibt das Inventar handgepflegt und unbelegt. Die Mindestlage sind
dann die Pflichten des Skripts aus §Wo der Check läuft: Der Check prüft
wenigstens sich selbst — seine Bash und seine eigenen Werkzeuge — und hat seinen
Break-Test hinter sich.

Ziel-Form der Sensor-Beschreibung:
[`templates/harness/sensors/environment-check.template.md`](../templates/harness/sensors/environment-check.template.md).

### Devcontainer/Compose-Kriterium

Devcontainer für IDE-Setup (Sprache-Server, Debugger-Anschluss). Compose
für Lauf- und CI-Vertrag. Beides parallel, wenn das Team mehrere IDEs
nutzt. Faustregel: Compose ist *Pflicht* (CI-Vertrag), Devcontainer ist
*Komfort*. Wer mit Devcontainer beginnt, baut sich eine zweite Toolchain
ohne die erste.

