# ADR-0090: `--suggest-config ai-harness` nimmt die Randbedingungs-Reihe `RB` nur bei tatsächlichem Repo-Fund ins Anforderungs-Muster auf, nicht unbedingt wie `FA`/`QA`

**Status:** Accepted

**Datum:** 2026-09-27

**Autor:** claude-sonnet-5

**Bezug:** Change Request des Konsumenten `ai-harness-course` (2026-09-27) —
die Randbedingungs-Reihe selbst steht dort in einer ungetaggten Welle;
slice-234 <!-- d-check:status-provenance -->. Nachgezogen nach unabhängigem
Review (R1-M3): die erste Fassung des Vertrags entschied sich gegen eine ADR
und deckte dabei nur die RTM-Abgrenzung ab, nicht die tatsächlich
nicht-triviale Aktivierungs-Frage unten — dieselbe Lücke, die der Review
benannte.

**Schärft:**
[`DC-FA-CLI-006`](../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(Erweiterung — kein neues Kürzel, dieselbe Präfix-Ableitung wie `FA`/`QA`).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`--suggest-config ai-harness[-init]` erzeugt aus der **bekannten Konvention**
(nicht aus dem Repo abgeleitet) ein kanonisches Anforderungs-`ids`-Muster
`<PREFIX>-(FA-[A-Z]+|QA)-\d+` — `FA`/`QA` stehen dort **unbedingt**, unabhängig
davon, ob das jeweilige Repo tatsächlich beide Kennungsarten führt (belegt
durch `TestCLI037_IDPrefix_Platzhalter`: das Muster erscheint bereits ohne
jedes Lastenheft). Der CR bittet um eine dritte, gleichrangige Reihe `RB`
(Randbedingungen).

Der naheliegende erste Entwurf — `RB` genauso unbedingt wie `FA`/`QA`
hinzufügen — widerspricht der eigenen Zusage von
[`DC-FA-CLI-006`](../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten):
ein Repo
**ohne** `-RB-`-Heading bekäme trotzdem ein geändertes Muster (das neue
`|RB` erscheint immer), und die im selben CR verlangte Byte-Gleichheit für
den `FA`/`QA`-only-Fall wäre unerreichbar. Dieser Widerspruch wurde erst beim
Lesen des bestehenden Codes sichtbar (`harnessIDPatterns` ist eine reine
Funktion von `reqPrefix`, ohne Repo-Kenntnis) — die erste Vertragsfassung
dieses Slice ging noch vom unbedingten Fall aus und musste vor dem Code
korrigiert werden.

## Entscheidung

Das Anforderungs-Muster nimmt `RB` **nur** auf, wenn derselbe
Ableitungs-Durchlauf, der ohnehin schon `reqPrefix` aus
`spec/lastenheft.md` liest (`deriveReqPrefix`, Modus `ai-harness` **ohne**
explizites `--id-prefix`), mindestens eine `-RB-`-Kennung sah. `FA`/`QA`
bleiben unbedingt (unverändertes Bestandsverhalten) — die Asymmetrie ist
bewusst: `FA`/`QA` sind die **Konvention selbst**, die diese Vorlage abbildet;
`RB` ist eine **optionale Erweiterung**, deren Byte-Gleichheits-Zusage nur
über eine Bedingung einlösbar ist.

**Grenze, benannt statt verschwiegen:** Die Bedingung prüft ausschließlich den
`ai-harness`-Ableitungs-Durchlauf. Ein Repo, das `--id-prefix` explizit setzt
oder den Voll-Kanon `ai-harness-init` nutzt, liest `spec/lastenheft.md` dafür
nicht — `RB` bleibt dort außen vor, selbst wenn das Repo `-RB-`-Kennungen
führt. Diese beiden Modi haben keinen Lastenheft-Lesezugriff, der die
Bedingung prüfen könnte, ohne die bestehende Byte-Gleichheits-Zusage für
`--id-prefix` bzw. `ai-harness-init` selbst zu gefährden.

**Zweite Konsequenz, ebenfalls benannt:** Die Präfix-Ableitung selbst zählt
`-RB-`-Kennungen jetzt zur Präfix-Menge (`reqShape` erkennt sie). Ein Repo,
dessen zweites, abweichendes Präfix **nur** über eine `-RB-`-Kennung auftritt,
wechselt dadurch vom stillen Erfolg (Ableitung ignorierte die Kennung vorher
komplett) zum Nutzungsfehler „mehrdeutiges Anforderungs-Präfix" — eine
sichtbare Verhaltensänderung für genau diesen (vom CR selbst verlangten) Fall,
die nicht unter die Byte-Gleichheits-Zusage für den `FA`/`QA`-only-Fall fällt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **`RB` unbedingt wie `FA`/`QA`** (erster Entwurf, verworfen) | einfachste Änderung, symmetrisch zu `FA`/`QA` | widerspricht der eigenen Byte-Gleichheits-Zusage für **jedes** Repo ohne `-RB-`-Heading — der CR verlangt beides zugleich, das ist nicht auflösbar |
| **`RB` gar nicht ins Anforderungs-Muster, nur Präfix-Ableitung erweitern** | keine neue Bedingung im Generator | erfüllt den CR nur zur Hälfte — die Vorlage prüfte `-RB-`-Kennungen dann nie auf Existenz/Link, obwohl der CR genau das verlangt |
| **`RB` bedingt auf tatsächlichen Fund im `ai-harness`-Ableitungs-Durchlauf** (gewählt) | erfüllt beide CR-Anforderungen (Muster erkennt `RB`, `FA`/`QA`-only bleibt byte-gleich); nutzt eine bereits vorhandene Lastenheft-Lesung, kein zweiter Scan | Asymmetrie zu `FA`/`QA` (unbedingt) muss erklärt werden; `--id-prefix`/`ai-harness-init` bleiben ohne `RB`-Erkennung (benannte Grenze) |

## Konsequenzen

- `deriveReqPrefix` liefert jetzt zusätzlich `sawRB bool`; `harnessIDPatterns`
  bekommt den Parameter `includeRB bool`. Beide Signaturen sind intern
  (unexported), keine öffentliche Schnittstellenänderung.
- Drei neue CLI-Akzeptanztests (`TestCLI234_RB_Happy`,
  `TestCLI234_RB_AbwesendByteGleich`, `TestCLI234_RB_Mehrdeutigkeit`) plus ein
  Round-Trip-Test (`TestCLI234_RB_AngewendetMeldetUnverlinkteKennung`), der
  die vorgeschlagene Regel tatsächlich anwendet und eine unverlinkte
  `-RB-`-Kennung als `id-unlinked` bestätigt bekommt — Beleg über die reine
  Muster-Prüfung hinaus.
- **Keine Schwelle oder Aktivierung für ein konkretes Repo** — diese ADR legt
  die Fähigkeit fest, nicht ihre Nutzung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | `TestCLI234_RB_*` (Happy, byte-gleich ohne Fund, Mehrdeutigkeit über RB, Round-Trip-Meldung) | `make test` |
| Manuelle Bestandsprobe | `--suggest-config ai-harness[-init]` gegen dieses Repo vor/nach der Änderung byte-gleich (kein `-RB-` im eigenen Lastenheft) | dokumentiert in der Closure-Notiz |

## Re-Evaluierungs-Trigger

Zwei Bedingungen, jede für sich hinreichend:

1. **Ein Konsument bittet um `RB`-Erkennung auch für `--id-prefix` oder
   `ai-harness-init`.** Dann ist zu prüfen, ob ein zusätzlicher,
   eigenständiger Lastenheft-Scan (unabhängig von der Präfix-Quelle)
   gerechtfertigt ist, oder ob die benannte Grenze bestehen bleibt.
2. **Eine vierte Kennungsreihe wird angefragt.** Dann ist zu prüfen, ob die
   bedingte Aktivierung als generelles Muster (Menge der im Repo tatsächlich
   gesehenen Kennungsarten) verallgemeinert werden sollte, statt jede Reihe
   einzeln zu verdrahten.

Ohne eines von beiden: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-234`, nachgezogen nach unabhängigem Review R1-M3) |
