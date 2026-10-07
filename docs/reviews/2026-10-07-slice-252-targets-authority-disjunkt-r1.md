# Review R1 — slice-252: `targets` prüft opt-in die Disjunktheit der Autoritäts-Dateien

- **Review-Art:** Code. Geprüft wurde gegen den Slice-Plan (`slice-252`, §1 Abgrenzung,
  §3 Spiegel-Liste, §6 Risiken), gegen `ADR-0101` (Proposed) und `ADR-0100` (Vorgängerin,
  nur ein Geschichte-Anhang), gegen die Baseline-Regel `v6.16.0` · `kurs/de/grundlagen/harness-dateien.md`
  §„Ein Index, mehrere Eigentümer" (Bedingung *Disjunkt*, dazu „Die Target-Zelle trägt den
  nackten Target-Namen"), gegen `MR-025`/`MR-032` und die Hard Rules `AGENTS.md`
  §3.1/§3.2/§3.7/§3.8 sowie §5 Regel 13/15/17. Gegenstand ist die Maintainability; die
  DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-252` · Commit `59d21c07`. Er umfasst `declaredTwiceFindings`
  und `CheckTargets` (Kern), `TargetsConfig.AuthorityDisjoint`, `rawTargets`/`applyTargets`
  (Config-Rand), `AllReasons`/`reasonTexts`, das `--print-config`-Gerüst, Tests,
  Lastenheft 0.97.0, Spezifikation (Schritt 5a, Schema-Zeile, `SPEC-088`),
  `ADR-0101`, den ADR-Index, den Geschichte-Anhang von `ADR-0100` und den Nachtrag in der
  CR-Antwort zur authority-Liste.
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan (`slice-252`); aus dem Lastenheft `DC-FA-TGT-001`
  (0.97.0); aus der Spezifikation `DC-FA-TGT-001.a` (Schritte 1, 3, 5, 5a, 6), `SPEC-005`
  (Schema-Zeile), `SPEC-088`; `ADR-0101`, `ADR-0100`; der eingehende Hinweis von
  `ai-harness-course` vom 2026-10-07; der Baseline-Wortlaut aus dem lokalen Klon
  (`git show v6.16.0:kurs/de/grundlagen/harness-dateien.md`). Vorherige Findings am Modul
  `targets`: R1 zu `slice-250` und `slice-251` (Klassen
  `grenze-gegen-beschreibung-statt-gegenstand-geprueft`, `grenze-ohne-negativtest`,
  `dublette-ohne-pfad-normalisierung`, `release-prep-nachzug`, Kommentar-Nachzug im Kopf
  von `CheckTargets`).
- **Proben:** Zwei Runtime-Images aus `git archive` gebaut — Vorher (`59d21c07^`) und
  Nachher (`59d21c07`) —, je Probe
  `docker run --rm --network none -v <fx>:/repo:ro <img> --enable targets` (Fixtures mit
  `chmod -R a+rX`). Zusatz-Images und Fixtures wieder entfernt; der Arbeitsbaum ist bis
  auf diesen Report unverändert. `make test`/`make lint`/`make gates` wurden nicht
  gefahren (Sache des Verifiers).

Probe-Ergebnisse (Nachher, wo nicht anders genannt):

```text
Fixture                                                        Ergebnis
f1  README.md: | `make gates` |, | `make ci` |                  rc 1: mk/tool.md:5 gates gate-declared-twice
    mk/tool.md: | `make tool-docs` | ... eingehängt in `make gates` | kein Gate |
f1  dieselbe Fixture ohne Schalter: normal/--json/--doctor/     Vorher == Nachher byte-identisch (stdout+stderr+rc)
    --doctor --json
Repo-Selbstlauf (Default und --enable targets --json)           Vorher == Nachher byte-identisch
f3  authority [./z.md, a.md, z.md, b.md, ./sub/c.md]           5 Befunde, Heimat ./z.md; a.md:3 zweimal `make x`
                                                               auf einer Zeile ⇒ ein Befund (SortFindings-Dedupe);
                                                               Ausgabe nach Pfad sortiert, ./sub/c.md zuerst
f5  authority [a.md, b.md] mit echter Doppelung, Schalter an,   rc 0, kein Befund
    ohne targets.makefiles
f8  authority [a.md, nope.md] + Schalter                        rc 2, Datei genannt (aus Richtung 2)
f8  authority-disjoint: "true"                                  rc 2 (strikt)
f8  authority [a.md] + Schalter                                 rc 0 (wirkungslos, wie zugesagt)
f8  authority [a.md, alias.md], alias.md -> a.md (Symlink)      rc 1: alias.md:3 x gate-declared-twice
```

## Findings

### H1 — HIGH — Eine bloße Erwähnung in einer Nachbarspalte zählt als Deklaration; der Werkzeug-Teil, der „eingehängt in `make gates`" schreibt, meldet `gates` als doppelt deklariert

- **quelle:** `DC-FA-TGT-001` / `ADR-0101` Entscheidung 3 gegen Baseline `v6.16.0` ·
  `kurs/de/grundlagen/harness-dateien.md` §Ein Index, mehrere Eigentümer
- **pfad:** `internal/hexagon/core/rules/targets.go:196` (über `extractDocTargets`,
  `targets.go:235` — jedes `` `make X` `` jeder Zelle); `spec/spezifikation.md:2775`
  (Schritt 5a: „jedes Tabellen-Target (Schritt 3)"); `spec/lastenheft.md:3341`
- **befund:** `declaredTwiceFindings` übernimmt die Extraktion von Richtung 2, die jedes
  `make X`-Token irgendwo in einer Tabellenzeile als Target nimmt. Für Richtung 2 ist das
  die nachsichtige Seite (eine Erwähnung entlastet); für die Disjunktheit ist es die
  strenge: Probe f1 — ein Werkzeug-Teil mit der Zeile
  `| \`make tool-docs\` | Doku-Prüfung, eingehängt in \`make gates\` | kein Gate |` liefert
  `harness/mk/tool.md:5 gates gate-declared-twice`, Exit 1, obwohl `gates` nur in einer
  Datei eine Zeile führt. Die Baseline meint mit „Kein Target steht in zwei Teilen" die
  Zeile, die das Target *führt* („sonst führen zwei Zeilen dasselbe Target"; „Die
  Target-Zelle trägt den nackten Target-Namen"), und beschreibt genau diesen Teil als einen,
  der Gates in `make gates` einhängt — der Fehlbefund trifft also den Regelfall, den die
  Regel selbst schildert. Dieselbe Form trägt dieses Repo schon heute: die
  `make doc-complete`-Zeile in `harness/README.md` nennt `make completeness-check` in der
  Bindung-Spalte. Weder Lastenheft noch Schritt 5a noch `ADR-0101` (§Verglichene
  Alternativen) benennen diese Lesart oder wägen „Target-Zelle" gegen „ganze Zeile" ab; der
  `--doctor`-Klartext sagt „deklariert", wo nur erwähnt ist.
- **verifizierbar:** ja — Black-Box-Lauf mit Fixture f1 (oben); ein Unit-Test mit
  Erwähnung in einer Nicht-Target-Zelle fehlt.
- **klasse:** `erwaehnung-als-deklaration-gezaehlt`

### M1 — MEDIUM — Ohne `targets.makefiles` bleibt der Schalter still wirkungslos; die Zusage nennt nur den Ein-Datei-Fall

- **quelle:** `DC-FA-TGT-001`, `AGENTS.md` §5 Regel 13 (Grenzen gegen den Gegenstand)
- **pfad:** `internal/hexagon/core/rules/targets.go:55`; `spec/lastenheft.md:3350`;
  `spec/spezifikation.md:2775`
- **befund:** `CheckTargets` kehrt bei leerem `Makefiles` vor jeder Richtung zurück, also
  auch vor Schritt 5a — Probe f5: zwei Autoritäts-Dateien mit echter Doppelung,
  `authority-disjoint: true`, kein `makefiles`-Eintrag ⇒ Exit 0, kein Befund, keine
  Warnung. Die Disjunktheit braucht die Regelmenge gar nicht; die Modul-Regel
  „leeres `makefiles` ⇒ inert (keine Regelmenge)" begründet sich mit einem Grund, der für
  5a nicht gilt. Lastenheft („mit nur einer Autoritäts-Datei ist der Schalter wirkungslos")
  und Schritt 5a („Nur bei … und mindestens zwei verschiedenen Autoritäts-Dateien") zählen
  die Bedingungen auf, unter denen der Schritt läuft, und lassen genau diese weg — ein
  Leser von 5a schließt, dass die Prüfung läuft. Vierter Slice in Folge an diesem Modul mit
  einer Grenze, die nur im Code steht.
- **verifizierbar:** ja — Black-Box-Lauf mit Fixture f5.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### L1 — LOW — „Dieselbe Datei in zwei Schreibweisen ist eine Datei" gilt nur für den bereinigten Pfad, nicht für einen Symlink-Alias

- **quelle:** `DC-FA-TGT-001`, `AGENTS.md` §5 Regel 13/15
- **pfad:** `spec/lastenheft.md:3350`; Kommentar `internal/hexagon/core/rules/targets.go:180`;
  Commit-Botschaft `59d21c07` („dieselbe Datei in zwei Schreibweisen sind kein Fall")
- **befund:** Der Code vergleicht `path.Clean`; ein Symlink `alias.md -> a.md` als zweite
  Autoritäts-Datei ergibt `alias.md:3 x gate-declared-twice` (Probe f8). Schritt 5a sagt
  korrekt „(bereinigter Pfad)", Lastenheft, Kommentar und Botschaft sagen es ohne
  Einschränkung. Praktische Reichweite gering, die Aussage aber weiter als der Gegenstand.
- **verifizierbar:** ja — Black-Box-Lauf mit Symlink-Alias.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### L2 — LOW — Spiegel im Lastenheft-Rumpf und im Config-Kommentar führen weiter zwei Grund-Codes

- **quelle:** `MR-025`
- **pfad:** `spec/lastenheft.md:3254` („Zwei Richtungen mit je eigenem Grund-Code"),
  `spec/lastenheft.md:3292` („`gate-phantom`/`gate-undocumented` liefern keinen
  `--repair`-Hunk"); `internal/adapter/driven/configyaml/configyaml.go:590`
  (`rawTargets`-Kommentar zählt die Schlüssel ohne `authority-disjoint` auf)
- **befund:** Die Repair-Aussage des Lastenhefts nennt den dritten Grund-Code nicht; ob
  `gate-declared-twice` einen Hunk liefert, sagt nur Schritt 6 der Spezifikation
  modulweit. Der `rawTargets`-Kommentar ist eine Schlüssel-Aufzählung, die jetzt einen
  Schlüssel unterschlägt (dieselbe Klasse wie `slice-251` R1 L1 am Kopfkommentar).
- **verifizierbar:** nein — Lesen.
- **klasse:** `spiegel-nicht-nachgezogen`

### L3 — LOW — Die CR-Antwort sagt eine Release-Version zu, bevor Review, Closure und Release-Prep gelaufen sind

- **quelle:** `AGENTS.md` §5 Regel 15/17
- **pfad:** `docs/plan/cr/2026-10-06-cr-eingehend-ai-harness-init-targets-authority-liste.md:143`
- **befund:** Der Nachtrag nennt „d-check `v0.83.0`" im Feature-Commit; im Vorgänger
  (`slice-251`) kam „Verfügbar ab" erst mit der Closure. Schiebt sich vor diese Closure ein
  anderes Release, steht in einem nach außen gerichteten Dokument eine falsche Version.
- **verifizierbar:** nein — Lesen gegen den späteren Tag.
- **klasse:** `release-prep-nachzug`

### I1 — INFO — Fundstellen tragen die Pfad-Form der Konfiguration

- **quelle:** Maintainability
- **pfad:** `internal/hexagon/core/rules/targets.go:211`
- **befund:** `File` und die Heimat in der Meldung stehen in der konfigurierten Form
  (`./sub/c.md`, „zuerst in ./z.md"); nach `SortFindings` sortiert `./sub/c.md` vor `a.md`.
  Das entspricht Richtung 1 (Bestand) und ist durch Schritt 5a nicht festgelegt;
  undokumentierte Annahme, kein Fehler.
- **verifizierbar:** ja — Probe f3.
- **klasse:** `pfad-form-der-konfiguration`

### I2 — INFO — Release-Prep-Flächen tragen weiter „eine Doppelnennung ist kein Befund"

- **quelle:** `AGENTS.md` §5 Regel 17
- **pfad:** `docs/user/benutzerhandbuch.md:2645` (§6-Modultabelle, auch Grund-Code-Spalte),
  `README.md:125`/`README.de.md`, `CHANGELOG.md` (neuer Eintrag)
- **befund:** Erwartet und laut Plan §1 Release-Prep; hier nur für die Release-Prep-Liste
  festgehalten.
- **verifizierbar:** nein.
- **klasse:** `release-prep-nachzug`

## Negativbefunde

- **Determinismus von `declaredTwiceFindings`:** geprüft, ohne Befund — Iteration über die
  Konfigurations-Liste und die Zeilen, Maps nur zum Nachschlagen; mehrere Vorkommen in
  einer späteren Datei je ein Befund, zwei Tokens auf einer Zeile per `SortFindings` zu
  einem (Probe f3).
- **`path.Clean`-Dubletten:** geprüft, ohne Befund — dieselbe Datei als `./z.md` und `z.md`
  wird einmal gelesen, kein Selbst-Befund (Probe f3, Test-Fall 3).
- **Fehlerfall fehlende Datei:** geprüft, ohne Befund — Exit 2 mit Dateinamen; der Fehler
  entsteht schon in Richtung 2, die alle Autoritäts-Dateien vorher liest (der
  Fehlerzweig in 5a ist dadurch praktisch unerreichbar, aber harmlos).
- **Wechselwirkung mit Richtung 1/2:** geprüft, ohne Befund — 5a ändert keine der beiden
  Mengen; ein exemptes Target in zwei Teilen meldet wie zugesagt.
- **Byte-Identität ohne Schalter / Commit-Probe:** nachgefahren, ohne Befund — Fixture f1
  ohne Schalter in vier Ausgabeformen und der Repo-Selbstlauf (Default sowie
  `--enable targets --json`) Vorher gegen Nachher byte-identisch. Strukturell ist die
  Zusage auch im Code gedeckt (früher Rücksprung bei `!AuthorityDisjoint`). Die
  Botschaft behauptet nur „ohne Schalter"; ihr Schluss geht nicht über die Proben hinaus
  (`AGENTS.md` §5 Regel 15).
- **Config-Rand:** geprüft, ohne Befund — Schlüssel strikt Bool (`"true"` ⇒ Exit 2), Default
  aus, `--print-config` nennt ihn; Vorher-Image lehnt den Schlüssel ab (erwartet).
- **Hexagon-Richtung (`ADR-0005`), Netz, Suppression (§3.1/§3.2):** geprüft, ohne Befund.
- **Kommentare (§3.7):** geprüft, ohne Befund bis auf L1/L2 — die neuen Kommentare tragen
  Zusage, Abgrenzung oder Grenze, keine Review-Historie, keine Slice-Nummern.
- **§3.8 (Eingaben außerhalb der Scan-Menge):** geprüft — die Autoritäts-Dateien sind
  dieselbe Nicht-Scan-Eingabe wie in Richtung 2, mit derselben Fence-/Tabellen-Erkennung;
  neu ist nur der Symlink-Fall (L1).
- **`MR-032`:** geprüft, ohne Befund — Lastenheft 0.96.2 → 0.97.0 mit Historie-Zeile,
  Spezifikation mit Historie-Zeile.
- **`ADR-0100` frozen:** geprüft, ohne Befund — nur ein Geschichte-Anhang, der die
  Doppelnennungs-Aussage auf `gate-undocumented` einschränkt.
- **Spiegel „Doppelnennung ist kein Befund" (`MR-025`):** geprüft — Lastenheft (Rumpf und
  Out-of-Scope), Spezifikation Schritt 5, CR-Antwort, `--print-config`, `--doctor`,
  `AllReasons`/§4 nachgezogen; Rest siehe L2 und I2. Historie-Zeilen und der
  ADR-0100-Titel im Index bleiben als eingefrorener Bestand korrekt stehen.
- **Abgrenzung Plan §1:** geprüft, ohne Befund — keine Datei außerhalb der Plan-Tabelle und
  Spiegel-Liste berührt; Handbuch/README/CHANGELOG nicht angefasst, kein make-Target für
  die Probe.
- **Struktur-ID `SPEC-088`:** geprüft, ohne Befund — nächste freie Nummer, Platzierung bei
  den `targets`-Codes wie im Bestand.
- **Adressierungs-Form (`Schärft:` in `ADR-0101`):** geprüft, ohne Befund — nennt
  `DC-FA-TGT-001.a`, `SPEC-005`, `SPEC-088`.
- **Provenance-Marker in `ADR-0101` §Bezug:** geprüft, ohne Befund — zeigt, wo umgesetzt,
  begründet nichts.

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | H1 |
| MEDIUM | 1 | M1 |
| LOW | 3 | L1, L2, L3 |
| INFO | 2 | I1, I2 |

Wiederkehrende Klasse: `grenze-gegen-beschreibung-statt-gegenstand-geprueft` (M1, L1) —
am Modul `targets` jetzt im vierten Slice in Folge.

## Verdikt

**Nicht freigegeben.** H1 blockiert: der neue Grund-Code meldet in der Form, die die
Baseline als Regelfall eines Werkzeug-Teils beschreibt, ein Target als doppelt deklariert,
das nur erwähnt ist — die Umsetzung trifft damit nicht, was die Baseline unter „disjunkt"
versteht. M1 ist vor dem Merge zu klären (Code oder Zusage). L1–L3 sind nice-to-fix,
I1/I2 zur Kenntnis.
