# Review — slice-151 / Commit `708bf26` (MR-049)

**Review-Art:** Code/Design (Closure-Profil-Regel) · **Gegenstand:**
`slice-151-urteilsfreie-haelfte-voll.md`, Commit `708bf26` · **Skill:**
`reviewer.md` @ 1.10.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-27

**Selbst gefahren:** `make verify-closure-notes` (Exit 0, 479/0), `make gates`
(Exit 0, 533/0), derselbe Lauf ohne die zwei `exempt-paths` (Exit 2, **121
Befunde: 107 `section-missing` + 14 `section-pattern-missing`**), eine Regel mit
`(?!` (Exit 2, *„invalid or unsupported Perl syntax: `(?!`"*), slice-153 mit
Freitext-Ausgängen (genau **ein** `section-pattern-missing`, Rückbau
sha256-geprüft), 16 konstruierte Proben. Quellen gelesen: `structure.go`,
`sections.go`, `paths.go` (`matchGlob`→segmentweise `path.Match`), `modul-05`
§Offene Risiken, MR-006, MR-049, `slice.template.md`.

Die drei Zahlen der Botschaft (**150 / 107 / 14**) sind exakt reproduziert. Was
nicht trägt, ist ihre **Deutung** und die **Reichweite** der Schlüsse.

## Findings

### F-1 · HIGH · „14 mit Freitext-Ausgang" misst nicht, was es behauptet

**Pfad:** `.d-check.closure.yml:170`; `MR-049-…md:37-39`; Botschaft „Messung 2".

Das Muster erlaubt zwischen `Ausgang` und dem Wort höchstens `:`+` `+`*`+`*`+`
`+`*`. Zwei belegte Schreibweisen fallen heraus:

- **Zeilenumbruch nach dem Marker** — Probe mit dem Produkt gefahren:
  `— **Ausgang:**\n  *eingetreten, …*` ⇒ `section-pattern-missing`. Der
  80-Spalten-Umbruch des Repos erzeugt diese Form von selbst.
- **Fett gesetzter Ausgang** — `**Ausgang:** **eingetreten.**` ⇒ ebenso.

Von den **14** gemeldeten Dateien tragen **7** sehr wohl mindestens einen
kanonischen Ausgang: `slice-125`, `slice-133`, `slice-134`, `slice-138` (nur
Umbruch) sowie `slice-126`, `slice-128`, `slice-129` (Fett + Umbruch).
Kontrolle: `slice-113` und `slice-121` tragen je drei kanonische Ausgänge und
würden nur mit einem erkannt.

**Failure-Szenario:** ein künftiger Slice, dessen §5 alle Risiken korrekt
schließt, wird beim Closure rot — allein wegen des Zeilenumbruchs.
**Empfehlung:** `[ \t\n]*` statt Einzel-Leerzeichen, `\*{0,2}` statt zweier
optionaler Sterne; danach die Bestands-Zahl neu erheben.

### F-2 · HIGH · `weiter offen` ist nicht ungeprüft, sondern verboten

**Pfad:** `MR-049-…md:23-25`; `AGENTS.md:349`; `.d-check.closure.yml:155-157`.

MR-049 formuliert den Verzicht als Nicht-Prüfung. Das begründet, **keine
Zusatzprüfung** zu bauen — nicht, das Wort aus der **akzeptierten Menge** zu
streichen. Belegt: ein §5 mit `**Ausgang:** *weiter offen* — BEO-099 trägt ihn
weiter.` ⇒ `section-pattern-missing`. Die Form existiert im Bestand (10
Vorkommen `**Ausgang:** offen; …`). Verschärfend sagt `AGENTS.md:349` „der
Wortschatz ist geschlossen", während zwei von drei Kanon-Ausgängen akzeptiert
werden.

### F-3 · MEDIUM · Die Ausnahme hört bei `slice-1000` still auf zu greifen

**Pfad:** `.d-check.closure.yml:166-168`; `MR-049-…md:43-47`.

`matchGlob` (`paths.go:98`) matcht segmentweise mit `path.Match`. `slice-1[0-3]*`
trifft `slice-1`, eine Ziffer aus `[0-3]`, dann **beliebigen Rest** — also auch
`slice-1000-…` bis `slice-1399-…`. Gemessen: ausgenommen wurden `slice-099`,
`slice-0999`, **`slice-1000`**, **`slice-1399`**; gemeldet `slice-140`,
`slice-1400`, `slice-153`. Die Botschaft erklärt in derselben Zeile, warum ein
Zahlen-Glob still ausfällt — und baut denselben Fehler mit weiterem Horizont.

### F-4 · MEDIUM · „107 tragen gar keinen §5-Abschnitt" ist falsch

**Alle 150** tragen eine `## 5.`-Überschrift; 107 unter anderem Titel
(`## 5. Trigger` 43 · `## 5. Closure-Trigger` 31 · `## 5. Risiken / offene
Punkte` 19 · weitere). Davon sind **25** echte Risiko-Abschnitte. Und **3 der
14** (`slice-106`, `slice-110`, `slice-111`) tragen ein §5 mit Risiken und
**keinen** `**Ausgang`-Marker — der Kanon-Kernsatz, dreimal verletzt, in der
Sammelzahl verschwunden.

### F-5 · MEDIUM · Die benannte Grenze ist in beide Richtungen weiter

„Ein §5 mit einem kanonischen und einem Freitext-Ausgang läuft grün durch" ist
korrekt; „Gedeckt ist der Abschnitt, in dem **kein** Ausgang kanonisch ist"
nicht. Probe: ein §5 mit einem Risiko (`*behoben*`) plus einem Prosasatz, der
„**Ausgang:** entfallen" zitiert ⇒ **grün**. Falsch-Rot ist gar nicht genannt
(F-1, F-2, F-8, F-9).

### F-6 · MEDIUM · „belegt nicht ausdrückbar" — gemessen ist eine Formulierung

Gemessen: RE2 weist `(?!` ab. Behauptet: die Korrelation sei nicht ausdrückbar.
`forbid-pattern` ist **existenzquantifiziert über jedes Vorkommen** und damit je
Risiko wirksam; das Komplement einer endlichen Wortmenge ist in RE2 darstellbar.
Mit dem Produkt gegengemessen: eine Präfix-Komplement-Regel läuft über die 12
nicht ausgenommenen `done/`-Slices mit **Exit 0**, meldet bei einem
Freitext-Ausgang `section-forbidden` und bleibt bei drei kanonischen
Schreibweisen grün.

### F-7 · MEDIUM · „schärft MR-006" hat keine Stütze

MR-006 regelt die Referenzrichtung; MR-049 nennt ihn im Körper nie, sein Feld
*Ersetzt-Baseline-Regel* sagt „keine". Gestalt von `BEO-012` in der
Erzeugungsrichtung.

### F-8 · LOW · Ein Ausgang in Inline-Code ist unsichtbar

`` — **Ausgang:** `entfallen` — … `` ⇒ `section-pattern-missing`.

### F-9 · LOW · Ein §5 ohne Risiken ist rot

### F-10 · LOW · Der Selektor trifft die Baseline-Vorlage nicht

`section: "## 5. Abnahme-Punkte / Risiken"` ist Klartext-Gleichheit; die
vendorte Vorlage führt `## 6. Risiken und offene Punkte`.

### F-11 · INFO · Der Befund zeigt auf die Überschrift, nicht auf das Risiko

## Negativbefunde

- **Zahlenwerk 150/43/107/14/121** exakt reproduziert; beanstandet sind die
  Etiketten.
- **Ausdrückbarkeits-Messung** — `(?!` wird tatsächlich abgewiesen, als sauberer
  Config-Fehler.
- **Die Regel beißt heute** — slice-153 mit Freitext ⇒ genau ein Befund;
  Rückbau sha256-identisch. Kein Stilles-Grün.
- **Vokabular-Messung** — über alle 144 Vorkommen in `done/`: **28**
  verschiedene Anfangswörter, darunter `behoben` (5), `gemessen` (4), `gehalten`
  (4), `erledigt` (2), `benannt` (2), `aufgelöst` (2).
- **Der gewächterte Bestand ist heute sauber** — alle 12 nicht ausgenommenen
  Slices tragen ausschließlich kanonische Ausgänge.
- **`slice-140` und alles ab 1400 sind nicht ausgenommen.**
- **Kommentar-Form (§3.7)** eingehalten; **Kanon-Zitat** wörtlich korrekt und
  pin-gebunden.

## Kategorie-Summary

HIGH 2 · MEDIUM 5 · LOW 3 · INFO 1

## Urteil

**Schließbar nach Nacharbeit.** Der Slice hat die richtige Frage gestellt und
ehrlich gemessen. Blockierend: das Muster meldet **kanonische Ausgänge als
Verstoß** (F-1) und **lehnt den dritten Kanon-Ausgang ab** (F-2) — in einem
Closure-Gate ein falscher Befund, kein Schönheitsfehler. Dazu eilen drei
tragende Dokumentations-Sätze der Messung voraus (F-4, F-6), und F-3 wiederholt
lautlos die Ausfall-Klasse, deren Vermeidung die Botschaft als eigene Messung
führt.
