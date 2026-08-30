# Eingehender CR aus `a-check` — `structure`: die legitime Leermenge einer erklärten Teilmenge

**Absender:** a-check (Adopter von d-check).
**Eingegangen:** 2026-08-30, über den Auftraggeber.
**Gegenstand:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
— ein optionaler Schlüssel an der Nullmengen-Härte von
[ADR-0075](../adr/0075-erklaerte-teilmenge-in-structure.md).
**Vorgänger:** [CR 3](2026-08-30-cr-a-check-structure-teilmenge.md) samt
[Antwort](2026-08-30-antwort-a-check-structure-teilmenge.md) — dieser Antrag
kommt aus dessen Anwendung.

Dieses Dokument hält den CR **wie empfangen** fest. Die Bewertung steht nicht
hier, sondern im Slice, der ihn aufnimmt — ein CR-Dokument trägt Bitte und
Beleg, nicht die Antwort darauf.

---

## Anlass — gemessen, und der Messende ist derselbe, der CR 3 gestellt hat

CR 3 hat geliefert, was er beantragte: `exempt-section-pattern` erreicht
Bestände **innerhalb** einer Datei, die `exempt-paths` nicht erreicht. Bei der
Anwendung stellte sich heraus, dass der Antrag seine eigene Konsequenz nicht zu
Ende gedacht hatte.

Der Adopter grandfathert **19** Anforderungen. Sein Lastenheft trägt genau **19**
`### AC-`-Überschriften. Die Regel

<!-- d-check-test:not-config: Antragstext des Absenders, kein .d-check.yml-Input -->
```yaml
- files: "spec/lastenheft.md"
  section-pattern: '^### AC-'
  sections: each
  require-all: [Happy, Boundary, Negative, Out-of-Scope]
  exempt-section-pattern: '^### (AC-FA-RULE-0(0[1-9]|1[01])|AC-FA-EXTRACT-001|…)'
```

meldet darum heute:

```
spec/lastenheft.md:1  …  section-missing
  alle 19 passenden Abschnitte sind von exempt-section-pattern ausgenommen — die Regel liefe leer
```

Das abgelöste Skript meldet an derselben Stelle `0 neue AC(s) geprueft, 19 grandfathered`, Exit 0.
**Das Modul macht mehr rot als der Sensor** — für den Adopter derselbe Bruch wie zu wenig rot, und
der Grund, warum dieser eine Eigenbau-Sensor als einziger nicht abgelöst werden konnte.

**Die Härte ist richtig, ihre Reichweite ist zu weit.** [ADR-0075](../adr/0075-erklaerte-teilmenge-in-structure.md) begründet sie mit: *„Ohne diese
Antwort schaltete ein zu breites Muster die Regel still ab."* Das trifft ein **generisches**
Muster. Es trifft nicht ein **aufzählendes**: das Muster oben nennt 19 Kennungen einzeln und kann
nicht versehentlich zu breit werden — es kann nur veralten, und dann meldet die Regel wieder.

**Zwei Zustände, die heute denselben Grund-Code teilen und verschiedene Dinge bedeuten:**

| Zustand | Bedeutung | heute |
|---|---|---|
| `section-pattern` trifft **nichts** | Konfigurationsdefekt: falsches Muster, falsche Datei, umbenannter Abschnitt | `section-missing` ✔ richtig |
| Muster trifft, **`exempt-section-pattern` nimmt alle** | Bestandszustand: die erklärte Ausnahme deckt gerade alles — *„es gibt noch nichts Neues zu prüfen"* | `section-missing` ✘ |

## Korrektur des Antrags — der erste Entwurf war am Modell nicht prüfbar

Er machte die Sichtbarkeit zur Bedingung („die Regel meldet eine **nicht-fatale** Zeile"). Die gibt
es nicht: `model.Finding` trägt kein Schwere-Feld, und [`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes) ist binär — *ein Befund ist
der Exit-1*, und „differenzierte Exit-Codes pro Befund-Kategorie" führt dieselbe Anforderung
ausdrücklich als **Out-of-Scope**. Der Antrag hätte damit stillschweigend einen Schweregrad
eingeführt. Die Fassung unten entscheidet sich stattdessen für einen der drei gangbaren Wege und
benennt, was sie dabei aufgibt.

**Die drei Wege, und warum es der dritte wird.**

| Weg | Was er kostet |
|---|---|
| Befund verschwindet ersatzlos | genau die Sichtbarkeit, die der Anlass verlangt — der Adopter merkt nicht, wenn seine Regel nichts mehr tut |
| Schweregrad in `model.Finding` | Vertragsänderung an [`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)/[`DC-FA-CLI-004`](../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate), viel größer als dieser Antrag, und ein **eigener** Entscheid des Werkzeugs |
| Sichtbarkeit über `--doctor` | die Zeile steht nicht mehr im Gate-Lauf, sondern im Diagnose-Modus — **schwächer**, aber ohne Vertragsbruch |

**Der dritte Weg ist tragfähig, weil `--doctor` diese Form schon führt.**
[`DC-FA-CLI-007`](../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus) verlangt in seinem Boundary-Kriterium: *„Given ein Repo **ohne Befunde**, when
`d-check --doctor` läuft, then Exit-Code 0 und eine Diagnose, die ‚0 Befunde' ausweist."* Der Modus
rendert also bereits heute etwas, das kein Befund ist, und seine Diagnose erscheint *„auf stdout
unabhängig vom Code"*. Eine Zeile über deklarierte Leermengen fügt sich dort ein, ohne
`model.Finding`, die Exit-Codes oder das Zeilenformat zu berühren.

## Vertrag

Ein optionaler Schlüssel an derselben Bedingung; ohne ihn byte-identisches Verhalten.
Kein neuer Grund-Code, kein neues Befund-Feld, keine Änderung an [`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)/`-004`.

<!-- d-check-test:not-config: Antragstext des Absenders, kein .d-check.yml-Input -->
```yaml
structure:
  - files: "spec/lastenheft.md"
    section-pattern: '^### AC-'
    sections: each
    require-all: [Happy, Boundary, Negative, Out-of-Scope]
    exempt-section-pattern: '^### (AC-FA-RULE-…|…)'
    exempt-may-empty: true
    #   NUR fuer exempt-section-pattern: leert die Ausnahme die Menge, entsteht
    #   KEIN Befund — die Regel wartet auf den ersten nicht ausgenommenen
    #   Abschnitt. Ohne den Schluessel: section-missing, heutiges Verhalten.
    #   Greift NICHT fuer exempt-paths und NICHT, wenn schon section-pattern
    #   nichts trifft — das bleibt ein Defekt.
```

**Die Sichtbarkeit wandert, und das ist der Preis — hier benannt statt verschwiegen.** Im
Gate-Lauf ist der Zustand ab dann stumm. `--doctor` führt ihn: eine Zeile je Regel mit gesetztem
Schlüssel, die die Regel-Identität und die Zahl der ausgenommenen Abschnitte nennt — *„alle 19
Abschnitte ausgenommen; greift beim ersten weiteren"* —, in `--doctor --json` als eigenes Feld
neben `findings`, damit sie nicht als Befund fehlgelesen wird.

**Warum der Preis hier vertretbar ist.** Die Sichtbarkeits-Zusage aus [ADR-0075](../adr/0075-erklaerte-teilmenge-in-structure.md) zielt auf ein
Muster, das **wirkt, aber danebengreift** — dort ist die stille Fehlkonfiguration die Gefahr. Hier
ist der Fall umgekehrt: das Muster trifft **alles**, und es ist eine namentliche Aufzählung von 19
Kennungen. Es kann nicht versehentlich zu breit werden; es kann nur veralten, wenn der Bestand
wächst und jemand die Aufzählung mitpflegt. Genau dafür ist `--doctor` der richtige Ort — ein
Modus, den man befragt, wenn man wissen will, was das Werkzeug gerade *nicht* prüft.

**Warum nicht „die Regel einfach weglassen, bis es etwas zu prüfen gibt".** Das ist die
naheliegende Antwort und die falsche: eine Regel, die man erst einschalten muss, wenn ihr Fall
eintritt, ist kein Gate. Der Sensor existiert genau dafür, beim **ersten** neuen Element zu
greifen — und in diesem Moment denkt niemand an die Konfiguration. Der Adopter hat für dieselbe
Klasse einen eigenen Beleg: das Modul `targets` stand bereit und lief **dreizehn Minor-Versionen**
ins Leere, weil das Target eingebunden, aber nie konfiguriert wurde.

**Abgrenzung gegen `exempt-paths`.** Der Schlüssel gilt bewusst **nicht** dort. Ein Datei-Glob ist
generisch und kann versehentlich einen ganzen Baum verschlucken; eine Abschnitts-Aufzählung
innerhalb *einer* Datei kann das nicht. Wer beides braucht, hat zwei verschiedene Fragen — und die
zweite ist hier nicht beantragt.

**Fence-Treue gilt weiter** — unverändert wie in CR 3.
