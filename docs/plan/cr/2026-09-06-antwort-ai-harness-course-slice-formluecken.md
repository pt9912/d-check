# Antwort auf CR „Der Slice-Plan hat zwei Formlücken"

**Absender:** `ai-harness-course` · **Datum:** 2026-09-06
**Bezug:** [CR vom 2026-09-06](2026-09-06-cr-ai-harness-course-slice-formluecken.md)
**Baseline-Stand des CR:** `v6.3.1`
**Ergebnis:** **beide Bitten angenommen** — Bitte 1 mit geschärfter Diagnose,
Bitte 2 im Befund angenommen, in der Lösung abweichend.

**Ablage nach [`MR-036`](../../../harness/conventions.md#mr-036):** wörtlich,
nicht zusammengefasst — der Wert liegt in der Begründung der Gegenseite, und
die überlebt keine Paraphrase.

---

## Wortlaut

> **Antwort auf CR „Der Slice-Plan hat zwei Formlücken" — ai-harness-course**
>
> Datum: 2026-09-06 · Baseline-Stand des CR: v6.3.1 · Beide Bitten: angenommen.
>
> **Bitte 1 — angenommen, mit geschärfter Diagnose.**
>
> Eure Messung am Bundle stimmt. Die Ursache ist eine andere als „Modul kennt
> den Abschnitt nicht": `kurs/de/03-agenten/modul-09-implementierung.md` führt
> Out-of-Scope an drei Stellen — in der Plan-Ausgabe des Worked Example
> („Out-of-Scope dieses Slice: Refresh-Token (slice-014b)"), als
> Review-Kriterium („Benennt der Plan Out-of-Scope, Risiken, Folge-Slices?")
> und in einem Nebensatz, der es voraussetzt. Alle drei sind Didaktik-Formen,
> und der Spiegel lässt Didaktik weg. Deshalb steht in
> `regelwerk/modul-09-implementierung.md` die 0, die ihr gemessen habt, und
> deshalb erreicht nichts die Vorlage.
>
> Der Befund ist damit nicht „Regel fehlt", sondern „Regel existiert nur als
> Erzählung" — und die trägt weder in den Digest noch in die Ziel-Form. Für die
> Lastenheft-Hälfte derselben Modul-6-Behauptung gilt das übrigens nicht:
> `lastenheft.template.md` führt `**Out-of-Scope:**` je Anforderung und einen
> globalen Abschnitt. Nur die Slice-Hälfte war ohne Träger. Wir encodieren die
> Regel in Modul 9, ziehen den Spiegel mit und geben der Vorlage den Abschnitt.
>
> **Bitte 2 — Befund angenommen, Lösung abweichend.**
>
> Der Widerspruch ist real, und er ist schon zweimal in Prosa überbrückt worden:
> im Template („Die beiden Vorgelagert-Blöcke entfallen nie") und in
> `templates/README.md` („§8 … steht in jedem Slice-Plan"). Zwei Pflaster über
> derselben Stelle sind das Signal, dass die Gliederung falsch ist — da sind wir
> bei euch.
>
> Euren Vorschlag setzen wir trotzdem nicht so um, und der Grund ist messbar:
> §8 ist eine **Adresse**, 37-mal benutzt über Kurs, Regelwerk, Templates und
> Beispiel, davon 4-mal in Zeitdokumenten (Review-Reports,
> Beobachtungs-Belege), die nach unserer Zeitform-Regel nicht nachgezogen
> werden. Eine Umnummerierung ließe die still falsch stehen.
>
> Dazu kommt ein Hinweis, dass die Adresse gar nicht die Lücke ist:
> `modul-06-roadmap.md` adressiert heute schon „Slice-Planung, §8 Vorgelagert —
> offene Beobachtungen sichten". Der Kanon benutzt §8 also längst für beide
> Hälften; falsch ist nur der **Titel**, der allein die bedingte Hälfte nennt —
> dadurch liest sich „darf entfallen" als Aussage über alles darunter.
>
> Wir behalten §8 und geben ihm einen Titel, der beide Hälften trägt: die zwei
> Vorgelagert-Blöcke werden sein **unbedingter Kopf**, die Modus-Begründung sein
> **bedingter Rumpf**. Das löst euren Befund strukturell, macht beide
> Prosa-Pflaster entbehrlich und bewegt keine Adresse. Falls ihr die Trennung in
> zwei Abschnitte aus einem Grund braucht, den wir nicht sehen, sagt es — dann
> reden wir über den Preis.
>
> **Zum Angebot.** Ja, gern: Schickt die Formulierungen für den
> Out-of-Scope-Abschnitt und eure drei Klassen. Wir übernehmen sie, wo sie
> passen, und nennen euren Beleg (fünf geschlossene Slices, dreimal nachweisbar
> getragen) als Anlass im Wellen-Eintrag. Kommt nichts, schreiben wir sie
> selbst — der CR blockiert nicht darauf.
>
> **Landung.** Eine Welle über Modul 9 · Modul 5 · beide Spiegel ·
> `slice.template.md` · `templates/README.md`. Regel-Änderung, also **MINOR**.
> Ihr bekommt Version, CHANGELOG-Zeile und CR-Verweis, sobald der Tag steht;
> euren lokalen Fork könnt ihr dann auflösen.
>
> **Abgrenzung von unserer Seite.** Die Umbenennungen eurer Haus-Form nehmen wir
> nicht auf — ihr habt sie selbst ausgenommen, wir sehen es genauso. Und
> `lab/example` ist in dieser Antwort weder Beleg noch Lieferumfang: Es
> illustriert die Lehre, es normiert sie nicht.

---

## Was daraus für dieses Repo folgt

- **Die Lösung zu Bitte 2 ist besser als unser Vorschlag**, und der Grund ist
  einer, den wir nicht gemessen hatten: die **Adress-Kosten**. §8 wird 37-mal
  adressiert, vier davon in Zeitdokumenten, die nach der Zeitform-Regel nicht
  nachgezogen werden. Unser Vorschlag hätte die still falsch stehen lassen.
  Kein Gegen-Einwand.
- **Die Haus-Form bleibt bis zum Bump** auf die angekündigte MINOR-Version und
  wird dann aufgelöst — nicht vorher. Eine Zwischen-Umstellung auf die heutige
  Acht-Abschnitt-Form hieße, zweimal umzubauen und die heute uniforme
  `done/`-Sektion vorübergehend ungleich zu machen.
- **Bringschuld erfüllt:** die zugesagten Formulierungen für den
  Out-of-Scope-Abschnitt. Der Bestand zeigt **vier** Klassen, nicht drei wie im
  CR behauptet — die Korrektur ging mit der Lieferung raus (siehe unten).

---

## Gelieferte Formulierungen (2026-09-06, weitergegeben)

Die im CR zugesagten Formulierungen sind geliefert. **Mit einer Korrektur:**
Der CR sprach von **drei** Klassen — der Bestand zeigt **vier**. Die vierte war
im CR übersehen und kommt im eigenen Bestand am häufigsten vor.

**Titel-Vorschlag:** `## N. Ausdrücklich NICHT in diesem Slice` — nicht
„Out-of-Scope". Die Formulierung stellt die Frage, die der Abschnitt
beantworten soll: *Was könnte man hier vermuten, das nicht kommt?*
„Out-of-Scope" lädt zur Aufzählung des ohnehin Fernliegenden ein.

**Bedienhinweis:** Je Punkt eine **Begründung**, nicht nur eine Nennung — ein
Ausschluss ohne Grund ist eine Behauptung. Was dort steht, ist die Grenze, an
der ein wachsender Slice sich messen lässt.

| Klasse | Was sie leistet |
|---|---|
| **1. Ein Folge-Slice übernimmt es** — mit Kennung | macht aus „später" eine Adresse |
| **2. Bestand bleibt bewusst stehen** — mit Begründung | verhindert, dass ein Sensor später gegen Altbestand meldet, den niemand entschieden hat |
| **3. Es wäre ein anderer Vorgang** | trennt Arbeit am Gegenstand von Arbeit am Werkzeug |
| **4. Schicht-Abgrenzung** | hält den Slice in seiner Schicht |

**Klasse 4 ist die übersehene** — sie steht in vier von sechs eigenen
Abschnitten und leistet etwas, das die anderen drei nicht tun: Sie sagt nicht,
*wann* etwas kommt, sondern *dass es woanders hingehört*.

**Nicht empfohlen:** eine Mindestzahl (ein echter Ausschluss schlägt vier
erfundene) und ein Sensor darauf (ob ein Ausschluss trägt, ist Urteil; ein
Pflichtfeld erzeugt Pflichterfüllung).

**Wirksamkeit, ehrlich eingeordnet:** Der Abschnitt hat dreimal als **Adresse
für Wachstum** getragen — was der Slice abgab, bekam dort eine Kennung. Er hat
**nicht** verhindert, dass ein Slice wächst: slice-203 wuchs von drei auf 23
Träger, mit vollständigem Abschnitt. Er macht Wachstum benennbar, nicht
unmöglich.
