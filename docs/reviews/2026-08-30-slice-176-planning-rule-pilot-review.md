# Review slice-176 — Der Wächter hielt weniger, als „jeder" verspricht

**Gegenstand:** [slice-176](../plan/planning/in-progress/slice-176-planning-rule-pilot.md), Stand `0cfe876`.
**Datum:** 2026-08-30. **Reviewer:** unabhängiger Subagent, Skill `.harness/skills/reviewer.md`.
**Eigener Lauf:** siehe §Negativbefunde. **Alle Befunde sind an grünen Gates vorbeigelaufen.**

---

Repo unverändert (`git status` leer, HEAD `0cfe876`), alle Proben liefen in `/tmp` und sind entfernt.

---

# Review-Report — slice-176

- **Review-Art:** Code + Plan (gegen `AGENTS.md` §3/§4/§5, `harness/conventions.md`, Baseline `modul-05`/`modul-06`/`grundlagen-harness-dateien.md`)
- **Gegenstand:** slice-176, Commits `8879ff1` · `776c0dd` · `a8f2c6b` · `b0850a9` · `0cfe876` (HEAD `0cfe876`)
- **Skill:** `.harness/skills/reviewer.md` v1.13.0
- **Modell-ID:** claude-opus-5[1m]
- **Datum:** 2026-08-30
- **Eingangs-Kontext:** `AGENTS.md` §3.1/§3.4/§3.6/§3.7/§3.8/§5, `harness/conventions.md` §Modus-Deklaration, `MR-011`/`MR-021`/`MR-042`/`MR-043`/`MR-045`/`MR-055`, Baseline `grundlagen-harness-dateien.md`, `modul-02`, `modul-05`, `modul-06`, `observations.md` (BEO-008/010/012/020/022/023/024), `welle-86`, `tools/harness/fetch-baseline-cache.sh`
- **Eigene Läufe:** `make gates` (Exit 0), `make baseline-verify`, 11 Umkehr-/Grenz-Proben gegen den neuen Sensor in einer isolierten Kopie unter `/tmp`, eine Produkt-Messung (`d-check:latest` gegen einen eigenen Mount)

---

## Findings

### HIGH-1 — Symlinks in Unterverzeichnissen von `.claude/rules/` werden nie geprüft, obwohl „jeder" zugesagt ist

- **kategorie:** HIGH
- **quelle:** `MR-055`; `AGENTS.md` §3.8; Reviewer-Frage 1 (Stilles-Grün-Pfad im Gate-Skript)
- **pfad:** `tools/harness/fetch-baseline-cache.sh:137` (`for l in "$rules"/*`), Zusage in `harness/conventions/MR-055-symlink-als-pin-traeger.md:11-13` („Symlinks **unterhalb** von `.claude/rules/`") und `:33-35` („dass **jeder** Symlink unter `.claude/rules/` auflöst"), gleichlautend `harness/conventions.md:135` und `slice-176:…§2 Punkt 4`
- **befund:** Der Glob ist flach und dotfile-blind. Gemessen in einer isolierten Kopie: ein **toter** Symlink unter `.claude/rules/sub/` und ein toter Symlink `.claude/rules/.versteckt.md` passieren beide grün — der Lauf meldet `verify ok (51 Dateien, vollständig)`, Exit 0. Die Zusage lautet an drei Stellen „jeder Symlink unter `.claude/rules/`"; der Geltungsbereich sagt „unterhalb", der Code sagt „genau eine Ebene, ohne Punkt-Namen". Das ist die Klasse, die `baseline-verify` mit derselben Änderung schließen wollte, in dem Verzeichnis, das der Eintrag selbst als seinen Gegenstand nennt — und `baseline-verify` läuft in `make gates`.
  ```
  ### H: toter Symlink in UNTERVERZEICHNIS .claude/rules/sub/
  fetch-baseline-cache: verify ok (51 Dateien, vollständig)
    exit=0
  ### I: toter Symlink als DOTFILE .claude/rules/.versteckt.md
  fetch-baseline-cache: verify ok (51 Dateien, vollständig)
    exit=0
  ```
  Heutiger Bestand ist flach und ohne Punkt-Namen, der Pfad wird also **derzeit** nicht begangen; die Zusage gilt trotzdem schon.
- **verifizierbar:** ja — `make baseline-verify` gegen einen Baum mit `.claude/rules/sub/<toter-alias>` (oben gefahren, Exit 0)
- **klasse:** `sensor-scope-flacher-glob`

### MEDIUM-1 — Der Sensor kann nicht rot werden, wenn die Zustellung selbst verschwindet; die Grenze ist nicht benannt

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.8; `MR-055` §Grenze
- **pfad:** `tools/harness/fetch-baseline-cache.sh:136` (`if [ -d "$rules" ]`), `harness/conventions/MR-055-symlink-als-pin-traeger.md:44-48`
- **befund:** Gemessen: fehlt `.claude/rules/` ganz, ist es leer, oder sind alle vier Aliase entfernt, meldet der Lauf `verify ok`, Exit 0 — in allen drei Fällen.
  ```
  ### E: .claude/rules/ FEHLT ganz (Zustellung geloescht)     -> verify ok, exit=0
  ### G: alle vier Symlinks GELOESCHT (Verzeichnis bleibt)     -> verify ok, exit=0
  ### D: LEERES .claude/rules/                                 -> verify ok, exit=0
  ```
  `MR-055` §Grenze nennt **eine** Grenze („geprüft wird die Auflösung, nicht das Ziel") und nicht diese. Konkretes Versagen: beim Bump stehen vier tote Aliase rot da; wer sie **löscht** statt sie umzuhängen, hat einen grünen Lauf und eine verschwundene Zustellung, und der DoD-Punkt „Der Bump-Träger ist gewächtert" trägt diesen Ausgang nicht mehr.
- **verifizierbar:** ja — `make baseline-verify` nach `rm .claude/rules/*` (oben gefahren, Exit 0)
- **klasse:** `zusage-vakuum-nicht-benannt`

### MEDIUM-2 — Die neue dritte Frage steht in keiner Deklarations-Fläche

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.8, §4 („Halluzinierte Gates sind die häufigste Form von Harness-Lüge")
- **pfad:** `tools/harness/fetch-baseline-cache.sh:23-24` (Kopf-Block `Modi:`), `Makefile:189` (Hilfetext), `AGENTS.md:373`, `harness/README.md:90`
- **befund:** Alle vier Flächen beschreiben `--verify`/`make baseline-verify` weiterhin als „`sha256sum -c` **plus Manifest-Deckung**" bzw. „**Zwei Hälften, beide nötig**". Das Skript liest seit `a8f2c6b` eine Eingabe, die es nie scannt (`.claude/rules/`), mit einer eigenen Fehlerpolitik und einer eigenen Grenze — genau die Konstellation, für die §3.8 die Frage „welche Eingaben liest es, die es nicht scannt?" stellt. Die Antwort steht ausschließlich in `MR-055` und im Slice; die Sensors-Tabelle in `harness/README.md` führt in ihrer Bindungs-Spalte `MR-011` und `MR-021`, nicht `MR-055`. `make gate-consistency` prüft Target-**Namen**, nicht Semantik — kein Gate fängt die Drift.
- **verifizierbar:** nein (Urteil; `grep -n "baseline-verify" AGENTS.md harness/README.md` zeigt den Stand)
- **klasse:** `sensor-zusage-nicht-deklariert`

### MEDIUM-3 — `MR-055` schreibt dem Kanon die Pin-Bindung zu und ankert sie auf eine Verzeichnisliste

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 („lies das Feld, nicht den Titel … die **direkteste** Quelle wählen"), `BEO-012`, Reviewer-Frage 9
- **pfad:** `harness/conventions/MR-055-symlink-als-pin-traeger.md:4-9`
- **befund:** Das Feld lautet: *„Der Kanon kennt die Verzeichniskonvention **und die Pin-Bindung** ([`grundlagen-harness-dateien.md` §Verzeichniskonvention])"*. Der zitierte Abschnitt ist ein Code-Fence mit einer Verzeichnisliste (Zeilen 4–24); er nennt `.harness/baseline/<tag>/` nicht einmal — dort steht nur `.harness/  # Skills, Tool-Allowlists, Checklisten-Middlewares`. Eine Pin-Aussage macht der Kanon an anderer Stelle (`modul-02-harness-bootstrap.md:147` Schritt 2, „als präsente, **gepinnte** Referenz"), die hier ungenannt bleibt. Die Bindung von In-Repo-Verweisen an den Pin ist zudem gerade die **Adaption** `MR-021`, nicht Kanon. Damit stützt der Anker die halbe Aussage, für die er steht.
- **verifizierbar:** nein (Urteil; `sed -n '4,24p' .harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md` zeigt den Abschnittsinhalt)
- **klasse:** `zitat-ueber-geltungsbereich`

### MEDIUM-4 — „laden **immer**" ist über die gemessene Menge hinaus verallgemeinert

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 (Schluss reicht nicht weiter als die gemessene Menge), `BEO-024`, Reviewer-Frage 8
- **pfad:** `docs/plan/planning/in-progress/slice-176-planning-rule-pilot.md` §2 („Sie tragen **kein** `paths:` und laden deshalb **immer**"), §4 DoD-Punkt 2, §2 „Was diese Zustellung **nicht** zusagt" (fünf Punkte, keiner nennt den Modus)
- **befund:** Belegt ist **eine** Sitzungsart: eine frische interaktive Hauptsitzung (`/memory`, `/context`). Der Schluss lautet „laden immer". `BEO-024` selbst zählt *„eine Regel, die nur im interaktiven Modus greift"* ausdrücklich zur Klasse, gegen die dieser Slice gebaut ist, und seine Prozedur verlangt, **wovon** das Greifen abhängt zu messen, bevor gebaut wird. Beobachtung aus diesem Lauf: dieser Review ist ein Subagent desselben Werkzeugs, im selben Repo, gestartet nach `0cfe876`; sein Projekt-Kontext führt `CLAUDE.md`, `AGENTS.md` und die Nutzer-Memory-Datei — **keines der vier Regelwerk-Module**. Ich kann nur meinen eigenen Kontext bezeugen, nicht den der Hauptsitzung; genau deshalb ist die unqualifizierte Form „immer" der Befund. Die Tragweite ist konkret: `AGENTS.md` §5 und die DoD dieses Slice verlangen Review und Verifikation **in eigenen Kontexten** — also in der Klasse, in der die Zustellung hier nicht ankommt.
- **verifizierbar:** nein (Werkzeug-Eigenschaft; wiederholbar durch `/context` in einer Subagenten- bzw. Print-Sitzung)
- **klasse:** `zusage-aus-einer-modus-messung`

### MEDIUM-5 — `welle-86` behauptet das Gegenteil des Gelieferten und nennt slice-176 als Beleg

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.7 (Zustandsfeld nennt den Zustand), Reviewer-Frage 7
- **pfad:** `docs/plan/planning/welle-86-closure-uebergang-durchsetzen.md:122-126` und `:160-166`
- **befund:** Das lebende Wellendokument sagt: *„Die vierte Schicht … die **pfad-gebundene** Zustellung ([slice-176])"* und im Nachtrag *„**Der Rules-Kanal fällt weg.** … **Auch die Zustellung läuft deshalb über Hooks**, die das Harness unabhängig vom Werkzeug ausführt; [slice-176] ist entsprechend neu geschnitten"*. Geliefert ist das Gegenteil: `.claude/rules/` **ohne** `paths:`, und §3 des Slice führt „**Kein Hook**" als ausdrückliche Nicht-Leistung. Konkretes Versagen: wer slice-175 aufnimmt, liest in seinem eigenen Wellendokument, der Rules-Kanal sei tot und die Zustellung laufe über Hooks — und plant gegen einen Kanal, den es nicht gibt, während der existierende ungenannt bleibt.
- **verifizierbar:** nein (Urteil)
- **klasse:** `wellendokument-widerspricht-lieferung`

### MEDIUM-6 — Beobachtungs-Register: zwei `Stand`-Zellen sind vom Neuschnitt überholt

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.7, Reviewer-Frage 7; Baseline `modul-06` §Beobachtungs-Register
- **pfad:** `docs/plan/planning/observations.md:11` (BEO-024, Spalte `Stand`), `:26` (BEO-022, Spalte `Stand`)
- **befund:** BEO-024 trägt *„**Antwort:** Hooks statt Regeln; das Harness führt sie aus, unabhängig vom Werkzeug"*; BEO-022 trägt *„**Beobachten**, ob die **pfad-gebundene** Zustellung ([slice-176]) die Klasse verschiebt oder nur verlagert"*. Beide Sätze stammen aus `b35662e`, dem Stand **vor** dem Neuschnitt `8879ff1`; `776c0dd` hat nur die Pfade von `open/` auf `in-progress/` gezogen, nicht die Aussage. Das Register ist Pflicht-Lesestoff jeder Slice-Planung — die nächste Sichtung liest als „Antwort" auf BEO-024 genau das, wogegen dieser Slice entschieden hat. **Mitigation, ausdrücklich:** das Register wird kanonisch bei **Slice-Closure** geschrieben, und slice-176 ist offen; die Korrektur ist also nicht überfällig, der Widerspruch steht aber währenddessen live.
- **verifizierbar:** nein (Urteil)
- **klasse:** `register-stand-ueberholt`

### MEDIUM-7 — Die Beobachtungs-Sichtung übergeht den Eintrag, der slice-176 namentlich führt

- **kategorie:** MEDIUM
- **quelle:** Baseline `modul-05-planning-harness.md:219` („Offene Beobachtungen sichten"), `MR-054`, `BEO-022`
- **pfad:** `docs/plan/planning/in-progress/slice-176-planning-rule-pilot.md` §7, zweiter Vorprüfungs-Block
- **befund:** Der Block sichtet BEO-024, BEO-010, BEO-023 und BEO-008. `BEO-022` fehlt — der einzige Eintrag, dessen `Stand` slice-176 **namentlich** als laufenden Beobachtungsgegenstand führt („Beobachten, ob die … Zustellung ([slice-176]) die Klasse verschiebt oder nur verlagert"). Der Slice ist damit an einer offenen Beobachtung vorbeigelaufen, die ihn adressiert; dass BEO-022 zusätzlich die Prozedur trägt, in derselben Änderung den **Lesepfad** einer neuen Regel zu benennen, hängt sachlich mit MEDIUM-8 zusammen.
- **verifizierbar:** nein (Urteil; die `d-check:cite`-Direktive prüft die Regel-Zeile, nicht die Vollständigkeit der Sichtung)
- **klasse:** `sichtung-uebergeht-eintrag-der-den-slice-nennt`

### MEDIUM-8 — Der neue Zustell-Kanal ist in keiner kanonischen Quelle des Repos genannt

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §1, `BEO-022` (Prozedur: Adressat **und** Lesepfad in derselben Änderung)
- **pfad:** `AGENTS.md` §1 (unverändert), `harness/README.md` (unverändert); `.claude/rules/` neu
- **befund:** `AGENTS.md` §1 weist an: *„Dabei pro Session **nur den benötigten Abschnitt** lesen … nicht das gesamte Regelwerk im Kontext halten."* Seit `a8f2c6b` sind vier Module (805 von 4787 Zeilen) unbedingter Kontext-Anteil; `AGENTS.md` sagt das nirgends, und `.claude/rules/` taucht dort nicht auf — §3.1 nennt als werkzeug-lokale Träger nur `.claude/settings.json` und den Hook. Konkretes Versagen: ein Lauf folgt der §1-Anweisung und liest `modul-05` erneut, den es bereits vollständig im Kontext hat; oder er verlässt sich umgekehrt darauf, dass nichts vorgeladen ist. Der Slice entscheidet bewusst, dass `AGENTS.md` **nichts abgibt** (§2 Punkt 5) — etwas **hinzuzufügen** ist davon nicht gedeckt.
- **verifizierbar:** nein (Urteil)
- **klasse:** `zustellkanal-ohne-lesepfad`

### LOW-1 — Fehlt `readlink`, meldet der Sensor eine falsche Diagnose

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §3.1 (Host-Klasse), Maintainability
- **pfad:** `tools/harness/fetch-baseline-cache.sh:104-107` (Vorprüfung nur `sha256sum`, `find`) gegen `:139`
- **befund:** Gemessen mit einem `readlink`-Stub (Exit 127) im PATH: alle **vier intakten** Aliase werden als tot gemeldet, mit dem Text „Ziel fehlt (Baseline-Bump?)" — ein Bump, den es nicht gab.
  ```
  fetch-baseline-cache: toter Symlink .claude/rules/modul-01-entwicklungszyklus.md — Ziel fehlt (Baseline-Bump?)
  … (vier Zeilen) …    exit=1
  ```
  Fail-closed, also **kein** stilles Grün; die Werkzeug-Vorprüfung existiert im selben Block genau dafür, eine solche Lage als Werkzeug-Fehlen statt als Befund auszuweisen.
- **verifizierbar:** ja — `PATH=<stub>:$PATH make baseline-verify` (oben gefahren)
- **klasse:** `werkzeug-vorpruefung-unvollstaendig`

### LOW-2 — Mess-Etikett nennt Dateien und sagt „Module"

- **kategorie:** LOW
- **quelle:** `BEO-020` („die gezählte Menge benennen, bevor die Zahl fällt")
- **pfad:** `slice-176…md` §2 Punkt 1 („Die übrigen **22 Module** bleiben draußen"), §2 Nachfolge-Entscheid und §3 („wer alle **26** einhängt")
- **befund:** Der vendorte Baum trägt **17** `modul-*.md`, **8** `grundlagen-*.md` und eine `README.md` — 26 **Dateien**, 17 **Module**, davon 4 eingehängt, also 13 übrige Module. Die Zeilenzahlen (805 von 4787) sind nachgemessen und korrekt; falsch ist nur der Nenner-Name. In einem Slice, dessen Anspruch ausdrücklich „gemessen, nicht geschätzt" lautet, und dessen Ausschluss-Argument („wer alle 26 einhängt") auf dieser Zahl steht.
- **verifizierbar:** ja — `ls .harness/baseline/v5.12.0/regelwerk/modul-*.md | wc -l` ⇒ `17`
- **klasse:** `mess-etikett-falscher-nenner`

### LOW-3 — Die Sub-Area-Vorprüfung nennt die einzige berührte Code-Sub-Area nicht als berührt

- **kategorie:** LOW
- **quelle:** Baseline `modul-05:213-214`, `harness/conventions.md` §Modus-Deklaration
- **pfad:** `slice-176…md` §7, erster Vorprüfungs-Block
- **befund:** Der Block führt `.claude/` und `harness/` als berührte Sub-Areas und erwähnt `tools/harness/` nur als Gegenbeispiel („eine eigene Deklaration führt nur `tools/harness/`") — obwohl die einzige Code-Änderung des Slice (`fetch-baseline-cache.sh`) genau dort liegt. Folgenlos für das Ergebnis: `tools/harness/` ist ebenfalls Greenfield (`harness/conventions.md:184`), die Modus-Begründung bleibt gültig.
- **verifizierbar:** ja — `git show --stat a8f2c6b`
- **klasse:** `berührte-sub-area-unvollstaendig`

### LOW-4 — `MR-055` führt kein `Begründung`-Feld

- **kategorie:** LOW
- **quelle:** `.harness/baseline/v5.12.0/templates/harness/conventions/MR-NNN-titel.template.md` („Pflichtfelder sind Datum, Geltungsbereich, **Ersetzt-Baseline-Regel**, Adaption, **Begründung** und Auflösungs-Trigger")
- **pfad:** `harness/conventions/MR-055-symlink-als-pin-traeger.md`
- **befund:** Der Eintrag trägt `Status`, `Ersetzt-Baseline-Regel`, `Datum`, `Geltungsbereich`, `## Adaption`, `## Grenze`, `## Auflösungs-Trigger` — kein Feld und keine Überschrift `Begründung`. Der Grund steht inhaltlich im Adaptions-Abschnitt („Warum die Prüfung dort und nicht in einem Modul"), aber nicht unter dem Namen, unter dem die Vorlage ihn einfordert. `MR-054` teilt die Abweichung; nach dem Nutzer-Maßstab „Regelwerk vor Bestand" ist der Nachbar kein Beleg.
- **verifizierbar:** nein (Form-Urteil gegen die Vorlage)
- **klasse:** `mr-pflichtfeld-fehlt`

### INFO-1 — Symlink-Portabilität ist keine der fünf Nicht-Zusagen

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `.claude/rules/*` (Index-Modus `120000`, belegt via `git ls-files -s`)
- **befund:** Die Ziele sind relativ (`../../.harness/baseline/v5.12.0/regelwerk/…`) und korrekt; `git archive HEAD` materialisiert sie in einem frischen Baum wieder als Symlinks (nachgestellt). Auf einem Checkout ohne Symlink-Unterstützung (`core.symlinks=false`, Windows, exFAT) entstehen daraus **reguläre Dateien mit dem Pfad als Inhalt**: `[ -L "$l" ]` ist dann falsch, der Sensor überspringt sie still, und die „Regel" ist eine Textzeile. Die fünf ausgeschriebenen Nicht-Zusagen decken das nicht.
- **verifizierbar:** nein (Umgebungs-Eigenschaft)
- **klasse:** `symlink-traeger-nicht-portabel`

### INFO-2 — Beleg-Grenze der Zustellungs-Messung steht nur in der Commit-Botschaft

- **kategorie:** INFO
- **quelle:** `AGENTS.md` §5
- **pfad:** `slice-176…md` §4, DoD-Punkte 2 und 3
- **befund:** `/memory` und `/context` sind Werkzeug-Anzeigen einer Sitzung des Auftraggebers und aus dem Repo nicht reproduzierbar. `b0850a9` sagt das („Der Auftraggeber hat eine frische Sitzung gestartet und zwei Anzeigen geliefert"); der **Plan** formuliert beide Punkte als eigene Messung, ohne die Herkunft zu nennen. Zwei Nebenpunkte: die 29,4k sind eine Differenz **zweier Sitzungen**, deren übrige Memory-Dateien nicht als konstant belegt sind; und die Plausibilitätsrechnung stimmt arithmetisch — 59 868 B / 29 400 Tok = 2,036 B/Tok, alle vier Byte- und Zeilenzahlen nachgemessen und korrekt —, aber „für deutsche Prosa … der erwartete Wert" steht ohne Vergleichsmessung da und ist damit Einschätzung, nicht Beleg.
- **verifizierbar:** nein
- **klasse:** `beleg-herkunft-nur-im-commit`

---

## Negativbefunde (geprüft, ohne Befund)

- **§3.1 Docker/make-only:** geprüft — keine Host-Toolchain, kein neuer Interpreter; der Sensor nutzt `readlink` (coreutils, in der Host-Klasse). Ohne Befund.
- **§3.2 Suppression-Verbot:** geprüft — keine `//nolint`, kein neues `d-check:ignore`, keine `.golangci.yml`-Ausnahme. Ohne Befund.
- **§3.3 Move-Zerlegung:** geprüft — `776c0dd` ist reiner `git mv` plus die nach `MR-013` gebündelten Verweise (Roadmap-Flip, `observations.md`, `welle-86`, `slice-171`); Slice-Datei im Move-Commit unverändert, Rename-Detection hält. Ohne Befund.
- **§3.4 Spec-Straten / Referenz-Richtung:** geprüft — kein Spec-Stratum und keine ADR berührt; `MR-055`, `conventions.md` und der Skript-Kommentar tragen keine Slice-, Wellen- oder Commit-Hash-Token. `MR-045` (keine Slice-Verweise in `AGENTS.md`/`harness/README.md`) unberührt. Ohne Befund.
- **§3.5 ADR-Immutabilität:** geprüft — keine ADR berührt. Ohne Befund.
- **§3.6 Gate-Lockerung:** geprüft — die Änderung fügt eine zusätzliche fail-closed Frage hinzu und senkt keine Schwelle; **Verschärfung, kein ADR nötig**. Ohne Befund.
- **§3.7 Kommentar-Klassen:** geprüft — der neue Skript-Kommentar (`:128-134`) trägt Kopplung („bindet denselben Pin"), Abgrenzung („von keinem Modul gescannt, in keiner Manifest-Zeile"), eine ausgewiesene Grenze und **ein** auflösbares Herkunfts-Feld (`MR-055`); keine Slice-Nummer, kein Mess-Label, keine Review-Historie, keine Herkunfts-Prosa. Ohne Befund.
- **§3.9 Workflow-Pins:** geprüft — keine Workflow-Datei berührt. Ohne Befund.
- **Umkehr-Probe des Sensors (`BEO-023`):** gefahren, der Wächter beißt — toter Alias ⇒ `fetch-baseline-cache: toter Symlink … — Ziel fehlt (Baseline-Bump?)`, Skript-Exit 1, `make`-Exit 2; entfernt ⇒ Exit 0. **Zwei** tote Aliase werden **beide** gemeldet (kein `errexit`-Abbruch nach dem ersten). Ein Alias mit **Leerzeichen** im Namen wird korrekt gemeldet. Eine **Schleife** (Symlink auf sich selbst) wird gemeldet. Ohne Befund.
- **Exit-Code-Behauptung in `a8f2c6b` („Exit 2"):** nachgemessen — Skript 1, `make baseline-verify` 2 (GNU-Make-Normalisierung). Die Botschaft misst die make-Ebene, die Angabe stimmt. Ohne Befund.
- **Symlink auf ein Verzeichnis:** gemessen — löst auf, passiert grün. Konsistent mit der Zusage „kein toter Alias"; kein Widerspruch zu einer geschriebenen Zusage. Ohne Befund.
- **Alias außerhalb des Pins:** gemessen — löst auf, passiert grün; genau so in `MR-055` §Grenze **benannt**. Ohne Befund.
- **Modus-Abdeckung:** geprüft — die dritte Frage läuft in `--verify` (`:224`) **und** am Ende des Vendor-Laufs (`:254`), nicht in `--check-latest` (`:225`, eigener Dispatch). Letzteres ist fail-open/informativ; eine fail-closed Frage dort wäre ein Bruch seiner Fehlerpolitik. Auslassung stimmig. Ohne Befund.
- **Bump-Wirksamkeit:** geprüft — der Vendor-Lauf eines **neuen** Tags läuft grün, solange der alte Baum noch steht; rot wird es, sobald `MR-021` Schritt (1) den alten Baum entfernt. Die Kopplung greift also am vorgesehenen Punkt. Ohne Befund.
- **Produkt-Messung „der Scanner folgt Symlinks nicht":** in einem eigenen Mount nachgestellt — von zwei Dateien mit totem Link wird die **echte** unter `.claude/rules/` gemeldet, der **Symlink** nicht (`2 Datei(en) geprüft, 2 Befund(e)`, der Alias fehlt in der Prüfmenge). Korroboriert durch den Dateizähler im inneren Lauf: 597 → 598 nach `a8f2c6b`, obwohl vier Symlinks **und** eine Markdown-Datei dazukamen. Behauptung bestätigt, ohne Befund.
- **`.claude/`-Ausnahme:** geprüft — `scan.ignore` führt nur `.harness/baseline/**` und `.harness/cache/**`; `.claude/` ist **nicht** ausgenommen, die Messung oben hängt also nicht an einem Ventil. Ohne Befund.
- **Zahlen im Plan:** nachgemessen — 86+248+228+243 = **805** Zeilen, **59 868** Bytes, Baum **4787** Zeilen, `AGENTS.md` **527**, Vorlage **236**, 59 868/29 400 = **2,04** B/Tok. Alle korrekt (Etikett: LOW-2). Ohne Befund.
- **`d-check:cite`-Anker der beiden Vorprüfungs-Blöcke:** nachgeschlagen — `modul-05:213-214` und `:219` tragen den zitierten Wortlaut wortgleich; der dritte Block trägt bewusst keine (`MR-054`). Ohne Befund.
- **`MR-021`-Deckung:** geprüft — `MR-021` §Geltungsbereich bindet *„**alle** Markdown-Links auf `.harness/baseline/<tag>/…`"* und überlässt die Menge *„dem Zensus der Bump-Prozedur"*. Beides gibt `MR-055` korrekt wieder, und „ein Symlink ist kein Markdown-Link" trifft zu. Ohne Befund.
- **Kollision mit `MR-011`:** geprüft — `MR-011` pinnt die Baseline auf einen Release-Tag; `MR-055` fügt `--verify` eine dritte, unabhängige Frage hinzu, ohne die beiden bestehenden zu verändern (Reihenfolge: `sha256sum -c` → Manifest-Deckung → Alias-Auflösung, alle drei fail-closed). Keine Kollision. Ohne Befund.
- **`MR-043`-Spannung:** geprüft und **nicht** als Befund gemeldet. `MR-043` argumentiert *„Kein zweiter Kandidat … Jede weitere importierte Datei wäre ein Pflichtanteil am Kontext **jedes** Laufs — genau die Kosten …"*, und slice-176 nimmt diese Kosten über einen anderen Kanal in Kauf; sein `Geltungsbereich` lautet aber ausdrücklich **`CLAUDE.md`**, deckt `.claude/rules/` also nicht. Nach der Prüffrage „lies das Feld, nicht den Titel" trägt der Eintrag den Widerspruch nicht — festgehalten, weil der Slice `MR-043` in seinem `Bezug:`-Feld führt und die Kosten-Aussage dort ungenannt bleibt.
- **`make gates`:** gefahren, **Exit 0** — `600 Datei(en) geprüft, 0 Befund(e)` (doc-check), `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`, Abschlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`.
- **`CHANGELOG.md`:** geprüft — keine nutzersichtbare Produktänderung (Werkzeug-Konfiguration + repo-eigenes Gate-Skript). Ohne Befund.
- **Arbeitsbaum:** `git status --porcelain` leer, HEAD `0cfe876` unverändert; alle Proben liefen in einer `git archive`-Kopie unter `/tmp` und sind entfernt.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 8 |
| LOW | 4 |
| INFO | 2 |

## Verdikt

**Blockiert.** HIGH-1 ist ein stiller Grün-Pfad in einem Gate, das in `make gates` läuft, und liegt innerhalb der Zusage, die `MR-055` an drei Stellen mit dem Wort „jeder" macht — gemessen, nicht vermutet. MEDIUM-1 bis MEDIUM-4 betreffen den Kern des Slice: die unbenannte Vakuum-Grenze des Sensors, die fehlende Deklaration der neuen Frage in allen vier Flächen, einen Anker, der die halbe Aussage nicht trägt, und die zentrale Zustellungs-Zusage „immer" aus einer Ein-Modus-Messung. MEDIUM-5 bis MEDIUM-8 sind Planungs-Zustand, der der Lieferung widerspricht — bei MEDIUM-6 mit der ausdrücklichen Mitigation, dass das Register kanonisch erst bei Closure geschrieben wird.

Der Sensor selbst **beißt** in allen Fällen, für die er gebaut wurde (Umkehr-Probe gefahren, Mehrfach-Befunde, Leerzeichen, Schleife), die Zahlen des Plans sind ausnahmslos nachgemessen und korrekt, und die Produkt-Messung „der Scanner folgt Symlinks nicht" ließ sich unabhängig reproduzieren. Die Substanz trägt; blockierend sind die Reichweiten — der Reichweite des Glob, der Reichweite der Deklarationen und der Reichweite der Zusagen.
