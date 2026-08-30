# Eingehender CR aus `a-check` — `workflows`: ein SHA, ein Tag-Kommentar

**Absender:** a-check (Adopter von d-check).
**Eingegangen:** 2026-08-30, über den Auftraggeber.
**Gegenstand:** [`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)
— eine dritte Bedingung in der Pin-Familie, hermetisch.
**Vorgänger:** [CR 4](2026-08-30-cr-a-check-leermenge.md) samt
[Antwort](2026-08-30-antwort-a-check-leermenge.md).

Dieses Dokument hält den CR **wie empfangen** fest. Die Bewertung steht nicht
hier, sondern im Slice, der ihn aufnimmt — ein CR-Dokument trägt Bitte und
Beleg, nicht die Antwort darauf.

**Zur Herkunft dieser Fassung, damit sie richtig gelesen wird:** der Absender
hat sein Dokument selbst mitgeliefert; es lag kurzzeitig hier und ist ohne
unser Zutun wieder verschwunden. Diese Fassung ist aus dem **übermittelten
Text** rekonstruiert, der Anhang aus dem Beleg-Register seines Slice
`slice-132-cr5-uses-tag-kohaerenz`. Wortlaut und Zahlen sind die des Absenders;
die Gliederung kann von seiner Datei abweichen.

---

## Antrag

Eine dritte Bedingung in der Pin-Familie von
[`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in):
**derselbe 40-stellige SHA trägt innerhalb der Scan-Menge überall denselben
Tag-Kommentar.** Vorschlag für den Grund-Code: `uses-pin-tag-conflict`,
gemeldet auf **jeder** beteiligten Zeile, weil keine von ihnen für sich falsch
ist.

## Das ist nicht euer Out-of-Scope

[`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)
nimmt die **Gültigkeit** eines SHA ausdrücklich aus — *„ob der SHA existiert und
den Commit bezeichnet, den der Tag-Kommentar behauptet, ist eine Netz-Frage und
ausdrücklich außerhalb"*. Dem widerspricht dieser Antrag nicht. Er fragt nicht,
ob ein Kommentar **wahr** ist, sondern ob zwei Kommentare **einander
widersprechen**. Das ist eine Aussage der Scan-Menge gegen sich selbst:
dieselbe Eingabe, die das Modul ohnehin parst, kein Netz, kein git, keine
zusätzliche Datei.

## Warum es in die Pin-Familie gehört und nicht daneben

`uses-pin-untagged` erzwingt bereits die **Existenz** des Tag-Kommentars. Damit
ist er im Vertrag keine Dekoration, sondern eine **Zusage** — und für eine
Zusage ist *„vorhanden"* die schwächste denkbare Prüfung. Ein Kommentar, den
zwei Zeilen verschieden schreiben, ist an mindestens einer Stelle falsch;
welche, sagt die Regel nicht, und das muss sie auch nicht. Sie sagt: hier
stimmt etwas nicht, seht nach.

## Der Anlass, gemessen

`a-check`s `release.yml` pinnte denselben `docker/login-action`-Digest
**zweimal** — einmal `# v4.2.0`, einmal `# v3.6.0`, 83 Zeilen auseinander,
beide Male derselbe SHA `650006c6…`. Der zweite Kommentar entstand beim
Kopieren der ersten Zeile. Über die GitHub-API aufgelöst: `v4.2.0` →
`650006c6…`, `v3.6.0` → `5e57cd11…`; der Digest ist `v4.2.0`, der zweite
Kommentar war falsch. Ein zweiter Fall derselben Art lag **zwischen zwei
Dateien** (`# v5.0.0` gegen `# v6.0.2` am selben `actions/checkout`-Digest).

## Wir haben gemessen, dass das Modul ihn heute nicht meldet

`d-check --enable workflows` (Pin `v0.69.0`) lief über `a-check`, während beide
widersprüchlichen Zeilen im Repo standen, und meldete genau **einen** Befund —
`uses-local-perms-undeclared` in `release.yml:242`. Der Tag-Konflikt war nicht
darunter. Das ist erwartungsgemäß: beide Zeilen tragen einen Tag-Kommentar,
`uses-pin-untagged` greift also bei keiner.

## Bei euch fände die Regel heute nichts

`d-check`s `.github/workflows/` führt **drei** distinkte SHAs mit je einem
Tag-Kommentar; `actions/checkout` steht **fünfmal** mit identischem Kommentar.
**Wiederholung ist kein Befund — nur Widerspruch.** Der Antrag stützt sich also
auf **unseren** Bestand, nicht auf euren; er ist eine **Regressions-Bremse**,
kein Bestandsräumer. Dass beide Repos dieselbe Pin-Konvention fahren und der
Fall trotzdem zweimal bei uns entstand, ist genau das Argument: die Konvention
allein verhindert ihn nicht.

## Aufwand, soweit wir ihn von außen einschätzen können

Wir haben die Regel bei uns als Bash-Sensor gebaut; unsere Fassung braucht
**19 Zeilen** (Extraktion plus Gruppierung nach SHA) — das sagt etwas über
unsere Umgebung, nicht über eure Codebasis. In eurer liegen die Referenzen
bereits geparst vor; was daraus folgt, wisst ihr besser als wir.

## Was wir nicht beantragen

Keine Aussage darüber, **welcher** der widersprüchlichen Kommentare der
richtige ist — das wäre die Gültigkeitsfrage und damit Netz. Und keine
Ausweitung auf Referenzen **ohne** Tag-Kommentar; die deckt `uses-pin-untagged`
bereits.

## Wenn ihr ablehnt

Dann bleibt die Regel bei uns als lokaler Sensor stehen und tut dort ihren
Dienst. Der Antrag ist eine **Einordnungsfrage** — gehört diese Zusage in die
Pin-Familie oder in die Harness des Konsumenten? —, keine Blockade.

## Anhang: der Prüf-Durchgang des Absenders

Der Absender führt zu jeder Tatsachen-Behauptung Herkunft und Handgriff.
Wiedergegeben, weil es die Belege prüfbar macht:

| Behauptung | Bestand | Handgriff |
|---|---|---|
| `release.yml` pinnte denselben Digest zweimal, `v4.2.0` und `v3.6.0` | eigener | `git blame` auf `release.yml:75` und `:158` vor der Korrektur |
| `v4.2.0` → `650006c6…`, `v3.6.0` → `5e57cd11…` | fremder | GitHub-API, beide Tags aufgelöst |
| ein zweiter Fall lag zwischen zwei Dateien (`v5.0.0` gegen `v6.0.2`) | eigener | eigener Register-Eintrag |
| das Modul meldet ihn heute nicht — ein Befund, und der war es nicht | fremder Vertrag | `d-check v0.69.0 --enable workflows` über `a-check`, **während** beide Zeilen im Repo standen; Ausgabe genau `uses-local-perms-undeclared` |
| `d-check` führt drei distinkte SHAs mit je einem Tag-Kommentar | fremder | `grep` über `.github/workflows/`, nach SHA gruppiert |
| `actions/checkout` steht fünfmal mit identischem Kommentar | fremder | derselbe `grep`, `uniq -c` |
| unsere Fassung braucht 19 Zeilen | eigener | `awk` über die beiden Funktionen seines Sensor-Skripts |

**Zwei Sätze sind ohne Handgriff geblieben, und der Absender sagt warum:**
*„für eine Zusage ist ‚vorhanden' die schwächste denkbare Prüfung"* ist ein
Argument, *„der Antrag ist eine Einordnungsfrage, keine Blockade"* eine Haltung
— keine der beiden behauptet eine Tatsache über ein System.
