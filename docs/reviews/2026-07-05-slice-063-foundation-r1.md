# Review-Report — slice-063 (Modul `targets`) — Doc-first-Fundament-Review R1

## Kopf-Metadaten

- **Datum:** 2026-07-05
- **Gegenstand:** Doc-first-Fundament für slice-063 (opt-in Modul `targets`,
  Deklarations-Konsistenz Doku ↔ Build-Targets) — Lastenheft-CR `DC-FA-TGT-001`,
  Spezifikation `DC-FA-TGT-001.a` + §2-Schema + §DC-FA-CLI-010-Anpassung,
  ADR-0031, ADR-Index-Zeile.
- **Reviewer-Rolle:** unabhängig/adversarial; nicht Mit-Autor des Fundaments.
  Dies ist das Folge-Review, das das bereits eingearbeitete Design-Review R1
  (Slice-Plan §8) voraussetzt.
- **Prüfmethode:** rein statisch — `git diff`, `grep`, Lesen. **Keine**
  `make`-/`go`-/`docker`-Läufe (Verifier-getrennt; doc-check/gates hat der Autor
  bereits grün, das bestätige ich nicht).
- **Eingangs-Kontext:** `git diff spec/lastenheft.md spec/spezifikation.md
  docs/plan/adr/README.md`; neue Datei `docs/plan/adr/0031-targets-deklarations-konsistenz-modul.md`;
  Slice-Plan `docs/plan/planning/next/slice-063-targets-modul.md`; Referenzen
  `tools/gate-consistency.sh`, `.d-check.yml`, `DC-FA-PLAN-001` / `DC-FA-PLAN-001.a`,
  `ADR-0028`; Reviewer-Skill v1.2.0; Regelwerk Modul 3 (Lastenheft) + Modul 4
  (ADR).
- **Arbeitsbaum-Scope (bestätigt):** `git status` zeigt genau `spec/lastenheft.md`,
  `spec/spezifikation.md`, `docs/plan/adr/README.md` geändert + ADR-0031/Slice-Plan
  neu; **kein** `.d-check.yml`-, Makefile- oder `internal/`-Diff — sauberer
  Doc-first-Schnitt (Config/Dogfood/Code folgen bei der Impl).

## Findings

### F-1 — MEDIUM — Ziel-Zählung: Lastenheft `DC-FA-CLI-010` (elf) widerspricht Spec §DC-FA-CLI-010.a (zehn)

- **kategorie:** MEDIUM
- **quelle:** `DC-QA-04` / Strata-Konsistenz (Lastenheft sticht Spezifikation,
  untere Schicht darf präzisieren, nie erweitern) / `DC-FA-CLI-010`
- **pfad:** `spec/spezifikation.md:335` (Mechanik) vs. `spec/lastenheft.md:379`,
  `:419`, `:423` (Vertrag) und die beiden Historien `spec/lastenheft.md:1602` /
  `spec/spezifikation.md:1535`
- **befund:** Der Lastenheft-Vertrag `DC-FA-CLI-010` wurde im Fundament
  **vollständig** auf **elf** `##`-Targets inkl. `doc-targets` angehoben:
  Beschreibung „sowie elf" (:379), Happy-Path-AK listet `doc-targets` (:419),
  Boundary-AK (:420), Out-of-Scope „gelisteten elf" (:423), Historie „10→11 …
  trägt ein `doc-targets`-Target". Die spezifizierende Mechanik-Sektion
  §DC-FA-CLI-010.a nennt dagegen unverändert „Zehn `.PHONY`-Targets" (:335) und
  führt `doc-targets` in der Aufzählung (:338–:377) **nicht**. Die
  Spec-Historie legt die Zurückstellung offen („§DC-FA-CLI-010.a (`doc-targets`)
  folgen mit der Modul-Implementierung"), die Lastenheft-Historie deklariert die
  Erweiterung dagegen als **vollzogen** („10→11") — die Offenlegung ist
  einseitig. Die Präzedenz-Slices 057 (`planning`) und 059 (`tracked`) hoben
  **beide** Straten gemeinsam im Fundament (Spec-Historie „8→9"/„9→10 Targets:
  `doc-planning`/`doc-tracked`"); dieses Fundament weicht davon ab. **Failure:**
  ein Implementierer, der §DC-FA-CLI-010.a als maßgebliche `--print-mk`-Mechanik
  nimmt, gibt zehn Targets ohne `doc-targets` aus und verletzt das
  Lastenheft-Happy-Path-AK (:419); ein Leser, der beide Straten abgleicht,
  findet „elf" gegen „Zehn" für dieselbe Fragment-Ausgabe und kann nicht
  entscheiden, was gilt.
- **verifizierbar:** nein durch Gate — **kein** Sensor prüft die
  `--print-mk`-Ziel-Zählung strata-übergreifend (deshalb blieb der Widerspruch
  grün); belegbar per `grep` „elf" (`lastenheft.md`) gegen „Zehn"
  (`spezifikation.md` §DC-FA-CLI-010.a). Dass der Harness dies nicht fängt, ist
  Teil des Befunds.

### F-2 — MEDIUM — Tabellen-Scoping: Spec „getrimmter Anfang" erweitert das von ADR/Lastenheft und Skript gepinnte `^\|`

- **kategorie:** MEDIUM
- **quelle:** `DC-QA-04` (Fidelity ≡ abgelöstes Skript) / Strata-Konsistenz
- **pfad:** `spec/spezifikation.md:1204` (Schritt 3) vs.
  `tools/gate-consistency.sh:26` (`grep -E '^\|'`), `ADR-0031` (Zeilen „`^|`"
  und „`grep -E '^\|'`") und die Lastenheft-Historie („`^\|`")
- **befund:** §DC-FA-TGT-001.a Schritt 3 definiert die Tabellen-Erkennung als
  „Zeilen, deren **getrimmter** Anfang ein Pipe `|` ist" — führende Leerzeichen
  sind damit erlaubt. Das abgelöste Skript nutzt `grep -E '^\|'` (Pipe in
  **Spalte 0**, keine Einrückung), und sowohl ADR-0031 als auch die
  Lastenheft-Historie behaupten ausdrücklich `^\|`- bzw. `grep -E '^\|'`-Parität.
  Die niedrigere Stratum-Sektion **erweitert** damit die von Lastenheft/ADR
  gepinnte Regel (statt sie zu präzisieren) und divergiert vom Skript, mit dem
  sie Parität beansprucht. **Failure:** eine eingerückte Tabellenzeile
  (führende Leerzeichen, dann `|` … `make X` … `|`) wird von einer spec-treuen
  Implementierung als Doku-Behauptung gewertet, vom Skript nicht → auf einem
  Repo mit solcher Zeile meldet das Modul `gate-phantom`, wo das Skript
  schweigt; der DoD-Paritätsnachweis „Modul-Befundsatz ≡ Skript" bricht dort.
  Auf d-checks aktuellem Baum stehen alle Tabellen in Spalte 0 → der
  Paritätslauf liefe **hier** dennoch grün, der Vertrag divergiert aber für
  Consumer-Repos (a-check/belief-agent) und künftige Docs.
- **verifizierbar:** teilweise — der DoD-Paritätstest (Modul vs. Skript) auf
  einem konstruierten Baum mit eingerückter `make X`-Tabellenzeile zeigt die
  Divergenz; auf dem aktuellen Baum nicht.

### F-3 — LOW — Lastenheft-Historie deklariert §4-Grund-Codes „in der Spezifikation" ohne den Zurückstellungs-Vorbehalt der Spec-Historie

- **kategorie:** LOW
- **quelle:** Maintainability / Doku-Drift
- **pfad:** `spec/lastenheft.md:1602` (Historie 0.38.0) vs.
  `spec/spezifikation.md:1535` (Historie) und `spec/spezifikation.md:1449`–`:1473`
  (§4-Tabelle)
- **befund:** Die Lastenheft-Historie 0.38.0 listet „Algorithmus-Sektion
  `DC-FA-TGT-001.a` + Grund-Codes `gate-phantom`/`gate-undocumented` +
  Schema-Keys … **in der Spezifikation**" ohne Vorbehalt. Die §4-Grund-Code-
  Tabelle (die kanonische, „stabil, maschinenlesbar"-Liste) enthält beide Codes
  **nicht** — nur `DC-FA-TGT-001.a` nennt sie als Algorithmus-Ausgabe; die
  Spec-Historie deklariert das offen („Die Grund-Codes … (§4) … folgen mit der
  Modul-Implementierung", `AllReasons`-↔-§4-Lockstep aus slice-060). Die beiden
  Historien beschreiben dieselbe Zurückstellung also **asymmetrisch**.
  **Failure:** ein Leser der Lastenheft-Historie hält `gate-phantom`/
  `gate-undocumented` für in der §4-Liste etabliert und referenziert sie als
  stabile Codes, findet sie dort aber nicht. (Die Zurückstellung von §4 selbst
  ist sachlich korrekt begründet — nur die Lastenheft-Formulierung trägt den
  Vorbehalt nicht mit.)
- **verifizierbar:** ja — `grep` der §4-Tabelle zeigt die Abwesenheit; **kein**
  Gate erzwingt §4-Vollständigkeit gegen die Algorithmus-Sektionen (bewusst, ist
  gerade der Grund der Zurückstellung).

### F-4 — LOW — Kein Akzeptanzkriterium übt `targets.exempt-targets` (die einzige Divergenz zum Skript) aus

- **kategorie:** LOW
- **quelle:** `DC-FA-TGT-001` (neuer öffentlicher Config-Vertrag ohne
  Skript-Präzedenz) / Test-Coverage
- **pfad:** `spec/lastenheft.md:1461`–`:1466` (Akzeptanzkriterien
  `DC-FA-TGT-001`) und `spec/spezifikation.md:1424` (§2-Zeile
  `targets.exempt-targets`, „exakt")
- **befund:** Die sechs Akzeptanzkriterien decken Happy / Boundary(Prosa) /
  Boundary(Modul-aus) / Negative(Phantom) / Negative(undokumentiert) /
  fail-closed ab — aber **keines** übt `targets.exempt-targets` aus. Das ist das
  **einzige** Modul-Verhalten ohne Skript-Präzedenz (`gate-consistency.sh`
  kennt keine Ausnahme, prüft ausnahmslos alle Makefile-Targets gegen §4). Das
  Negative-AK sagt nur „`secret:` (nicht in `targets.exempt-targets`)" — es
  bestätigt die Ausnahme-Seite nicht. Die Exakt-Namen-Match-Semantik (§2:
  „Regelnamen (exakt)") ruht allein auf Prosa. **Failure:** eine Implementierung,
  die `exempt-targets` als Substring- oder Glob-Match (statt exaktem Namen)
  umsetzt, bliebe vom AK-Satz ungefangen — ein exemptes `foo` würde dann auch
  `foobar` grandfathern.
- **verifizierbar:** ja — ein AK „ein in `targets.exempt-targets` gelistetes
  Makefile-Target löst **kein** `gate-undocumented` aus" würde die Semantik
  verriegeln; es fehlt.

### F-5 — INFO — Zeichensatz des dokumentierten `make X`-Tokens ist unspezifiziert; das Skript ist kleinbuchstaben-beschränkt

- **kategorie:** INFO
- **quelle:** `DC-QA-04` (Fidelity-Detail)
- **pfad:** `spec/spezifikation.md:1201`–`:1206` (Schritt 3) vs.
  `tools/gate-consistency.sh:26` (Doku-Token-Regex `make` + `[a-z][a-z0-9_-]*`)
- **befund:** §DC-FA-TGT-001.a legt den Zeichensatz des dokumentierten
  `make X`-Tokens nicht fest. Das Skript beschränkt Doku-Tokens auf
  **Kleinbuchstaben** (Zeichensatz `[a-z][a-z0-9_-]*` nach `make`), während die
  Makefile-Regel-Heuristik (Schritt 2) Großbuchstaben zulässt
  (`[A-Za-z][A-Za-z0-9_-]*`). **Failure:** ein großgeschriebenes Makefile-Target
  (z. B. `Build:`), in einer Tabelle als `make Build` dokumentiert, würde vom
  Skript **nicht** als dokumentiert erkannt (Kleinbuchstaben-Regex) und als
  `gate-undocumented` gemeldet; eine spec-treue case-insensitive Implementierung
  schwiege — Divergenz. Für die durchweg kleingeschriebenen Harness-Targets
  irrelevant, aber ein unbenanntes Fidelity-Detail.
- **verifizierbar:** teilweise — nur auf einem Baum mit großgeschriebenem
  Target sichtbar.

### F-6 — INFO — Richtung 2 ist auf nicht-leeres `doc-tables` gekoppelt, obwohl sie nur `makefiles` + `authority` konsumiert

- **kategorie:** INFO
- **quelle:** `DC-FA-TGT-001` (Design-Kopplung)
- **pfad:** `spec/spezifikation.md:1188`–`:1192` (Schritt 1) und die
  Lastenheft-Inert-Regel („Leeres `targets.makefiles` **oder**
  `targets.doc-tables` ⇒ Modul inert")
- **befund:** Richtung 2 (`gate-undocumented`, Makefile → `authority`)
  konsumiert nur `targets.makefiles` + `targets.authority`; das Modul ist aber
  schon bei leerem `targets.doc-tables` **vollständig** inert (Schritt 1:
  „Leeres `makefiles` **oder** `doc-tables` ⇒ inert"). Ein Consumer, der nur die
  Undokumentiert-Richtung fahren will (Makefile ↔ `authority`), muss dennoch
  `doc-tables` setzen, sonst läuft auch Richtung 2 nicht. Über die Straten
  konsistent dokumentiert (Design-Wahl analog `planning` ohne Roadmap-Pfad),
  aber die Kopplung ist nicht begründet. **Failure:** Config mit gesetztem
  `makefiles` + `authority`, aber leerem `doc-tables` ⇒ kein Befund (inert),
  obwohl Richtung 2 vollständige Eingaben hätte — stiller Nicht-Lauf.
- **verifizierbar:** ja — der genannte Config-Fall liefert einen leeren
  Befundsatz.

## Negativbefunde (geprüft, ohne Befund)

- **ADR-0031-Kernbehauptung (Prämissen-Wende):** geprüft gegen
  `0028-planning-lifecycle-modul.md:111`–`:114` — ADR-0028 stuft
  `gate-consistency` dort tatsächlich als „bewusst **nicht** d-check-fähig"
  („Meta-Spannung Markdown↔Makefile↔YAML") ein. ADR-0031 zitiert das korrekt und
  widerlegt es schlüssig: `planning` liest hermetisch eine deklarierte Datei +
  Repo-Struktur über `ReadFile`; `targets` liest deklarierte Datei + Makefile
  über denselben Port — dieselbe Klasse. Kein Befund.
- **„Kein Supersede von ADR-0028":** geprüft — die **Entscheidung** von ADR-0028
  (Modul `planning`) steht; überholt wird nur eine Konsequenz-Seiten-Einstufung,
  keine Entscheidung. Hard-Rule-konform (neue ADR, expliziter Verweis, kein
  Überschreiben; Accepted-ADR unangetastet). „Kein Supersede" ist korrekt. Kein
  Befund.
- **Scope-Entscheid (nur Kern; Prüfung-3 im Skript-Rest; Skript nicht voll
  entfernt; kein `codepaths.ignore-refs`-Tombstone):** geprüft — konsistent über
  Slice-Plan §2, ADR-0031-Konsequenzen und Lastenheft-/Spec-Out-of-Scope. Kein
  Befund.
- **Makefile-Regel-Heuristik (literale Namen am Zeilenanfang, `:` ohne
  folgendes `=`, kein `.`-Präfix inkl. `.PHONY`, Mehrfach-Target-Zeilen):**
  Schritt 2 gegen `makefile_targets` (`gate-consistency.sh:35`) durchgespielt
  (`X :=`/`X ?=`/`foo:=bar` fallen durch das `([^=]|$)`; `.PHONY` fällt bereits
  durch den `^[a-zA-Z]`-Anker; `a b: dep` liefert beide Namen) — deckungsgleich;
  die Spec macht die `.PHONY`-Ausnahme nur explizit, die das Skript implizit
  über den Anker leistet. Eine marginale Rest-Kante (ein zweiter, ziffern-
  geführter Name wie `foo 9bar:` — Skript liefert `9bar`, die per-Name-Regel
  `[A-Za-z]…` nicht) ist auf realen Targets ohne Wirkung. Kein Befund.
- **§2-Schema ↔ Algorithmus:** `makefiles`→Schritt 2, `doc-tables`→Schritt 3/4,
  `authority`→Schritt 5, `exempt-targets`→Schritt 5; Einzahl/Mehrzahl-Typen
  passen. Kohärent. Kein Befund.
- **fail-closed-Konsistenz:** „fehlende konfigurierte Datei ⇒ Exit 2" ist in
  Lastenheft-Beschreibung, Lastenheft-AK, Spec-Schritt 1 und allen drei
  §2-Schema-Zeilen deckungsgleich; kein stiller-Grün-Pfad spezifiziert. Kein
  Befund.
- **Akzeptanzkriterien-Deckung (Happy/Boundary/Negative je Richtung +
  fail-closed):** die sechs AKs sind vollständig und prüfbar formuliert (je
  Given/When/Then, Exit-Codes, Grund-Codes, Datei/Zeile/Target). Kein Befund
  (die exempt-Lücke ist F-4, nicht Teil der sechs).
- **Versions-/Zähl-Arithmetik:** Modul-Liste in `DC-FA-CLI-002` + Glossar auf
  17 Einträge nachgezählt (`targets` als 17.); Bereichskürzel `TGT` in §3
  ergänzt; ADR-0031 = nächste freie Nummer (nach 0030); v0.38.0 = korrekter
  Minor-Bump von 0.37.1 (neues Produkt-Modul). Kein Befund.
- **Referenz-Richtung (SDP):** ADR-0031 „Schärft" die Spec-Sektion
  (ADR → Spec, korrekt); Spec-Straten verweisen aufwärts (Lastenheft/QA);
  Slice-Tokens erscheinen nur in ADR-Historie/Prosa als Provenance, nicht als
  Entscheidungsgrundlage im Körper. Kein Befund.
- **Anker/Slugs:** die Slugs für `DC-FA-TGT-001` und `DC-FA-TGT-001.a` sind
  zwischen Überschrift und den verweisenden Links (Lastenheft, Spec, ADR,
  ADR-Index) intern deckungsgleich (statisch geprüft; die Gate-Bestätigung ist
  Verifier-Sache). Kein Befund.
- **Doc-first-Reinheit:** der Arbeitsbaum berührt nur die drei Doc-Straten +
  ADR (+ Slice-Plan); `.d-check.yml`-`targets:`-Block, Makefile-Dogfood und
  `internal/`-Code fehlen bewusst (Impl-Phase). Korrekter Schnitt. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 2 |

## Verdikt

**NACHBESSERN.**

Begründung: kein HIGH — das Fundament ist substanziell tragfähig (ADR-Kernthese
belegt und Hard-Rule-konform, Algorithmus deckt die sechs AKs, Makefile-Fidelity
zum Skript gehalten, fail-closed durchgängig, Doc-first sauber geschnitten). Die
zwei MEDIUMs sind aber genau die Klasse, die ein Doc-first-Fundament abfangen
soll — strata-übergreifende Widersprüche, hier am billigsten (Doku) zu heilen:
F-1 (`DC-FA-CLI-010` sagt elf Targets, §DC-FA-CLI-010.a sagt zehn — ein Sensor
fängt es nicht, im Gegensatz zu den 057/059-Präzedenzen, die beide Straten
gemeinsam hoben) und F-2 (Spec „getrimmter Anfang" erweitert das von
ADR/Lastenheft und Skript gepinnte `^\|` und untergräbt die „≡ Skript"-
Paritäts-DoD). Beide sind reine Textangleichungen ohne Code-Berührung. F-3/F-4
(LOW) und F-5/F-6 (INFO) sind nachrangig, F-4 (exempt-AK) lohnt die Aufnahme in
die Impl-Test-DoD.
