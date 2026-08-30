# Eingehender CR aus `a-check` — `workflows`: ein SHA, ein Tag-Kommentar

**Absender:** a-check (Adopter von d-check).
**Eingegangen:** 2026-08-30, über den Auftraggeber.
**Gegenstand:** [`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)
— eine dritte Bedingung in der Pin-Familie, hermetisch.
**Vorgänger:** [CR 3](2026-08-30-cr-a-check-structure-teilmenge.md) und
[CR 4](2026-08-30-cr-a-check-leermenge.md), beide angenommen.

Dieses Dokument hält den CR **wie empfangen** fest. Die Bewertung steht nicht
hier, sondern im Slice, der ihn aufnimmt — ein CR-Dokument trägt Bitte und
Beleg, nicht die Antwort darauf.

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

## Warum in die Pin-Familie und nicht daneben

`uses-pin-untagged` erzwingt bereits die **Existenz** des Tag-Kommentars. Damit
ist er im Vertrag keine Dekoration, sondern eine Zusage — und für eine Zusage
ist „vorhanden" die schwächste denkbare Prüfung. Ein Kommentar, den zwei Zeilen
verschieden schreiben, ist an mindestens einer Stelle falsch; welche, sagt die
Regel nicht, und das muss sie auch nicht. Sie sagt: hier stimmt etwas nicht,
seht nach.

## Der Anlass, gemessen

`a-check`s `release.yml` pinnte denselben `docker/login-action`-Digest
**zweimal** — einmal `# v4.2.0`, einmal `# v3.6.0`, 83 Zeilen auseinander,
beide Male derselbe SHA `650006c6…`. Der zweite Kommentar entstand beim
Kopieren der ersten Zeile. Über die GitHub-API aufgelöst: `v4.2.0` →
`650006c6…`, `v3.6.0` → `5e57cd11…`; der Digest ist `v4.2.0`, der zweite
Kommentar war falsch. Ein zweiter Fall derselben Art lag **zwischen zwei
Dateien** (`# v5.0.0` gegen `# v6.0.2` am selben `actions/checkout`-Digest).

**Gemessen, dass das Modul ihn heute nicht meldet:** `d-check --enable
workflows` (Pin `v0.69.0`) lief über `a-check`, während beide widersprüchlichen
Zeilen im Repo standen, und meldete genau **einen** Befund —
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

## Aufwand, soweit von außen einschätzbar

Wir haben die Regel bei uns als Bash-Sensor gebaut; unsere Fassung braucht **19
Zeilen** (Extraktion plus Gruppierung nach SHA) — das sagt etwas über unsere
Umgebung, nicht über eure Codebasis. In eurer liegen die Referenzen bereits
geparst vor; was daraus folgt, wisst ihr besser als wir.

## Was wir nicht beantragen

Keine Aussage darüber, **welcher** der widersprüchlichen Kommentare der richtige
ist — das wäre die Gültigkeitsfrage und damit Netz. Und keine Ausweitung auf
Referenzen **ohne** Tag-Kommentar; die deckt `uses-pin-untagged` bereits.

## Wenn ihr ablehnt

Dann bleibt die Regel bei uns als lokaler Sensor stehen und tut dort ihren
Dienst. Der Antrag ist eine **Einordnungsfrage** — gehört diese Zusage in die
Pin-Familie oder in die Harness des Konsumenten? —, keine Blockade.

---

## Anhang: der Prüf-Durchgang des Absenders

`a-check` schickt CR-Texte erst nach einem Durchgang, der jeden Satz markiert,
der eine **Tatsache über ein System** behauptet, und den Handgriff nennt, der
ihn belegt. Für diesen Text: **neun Tatsachen-Behauptungen, neun belegt.**

| Behauptung | Klasse | Handgriff |
|---|---|---|
| „die Gültigkeit eines SHA ist ausdrücklich ausgenommen" | fremder Vertrag | Out-of-Scope-Absatz gelesen, wörtlich zitiert |
| „`uses-pin-untagged` erzwingt bereits die Existenz des Tag-Kommentars" | fremder Vertrag | Bedingungs-Tabelle, Zeile 2 |
| „`release.yml` pinnte denselben Digest zweimal" | eigener Bestand | `git blame` auf `release.yml:75` und `:158` vor der Korrektur |
| „`v4.2.0` → `650006c6…`, `v3.6.0` → `5e57cd11…`" | fremder Bestand | `gh api repos/docker/login-action/git/ref/tags/<tag>` |
| „ein zweiter Fall lag zwischen zwei Dateien" | eigener Bestand | eigenes Beobachtungs-Register, Beleg-Slice |
| „das Modul meldet ihn heute nicht" | fremder Vertrag | `d-check v0.69.0 --enable workflows` über `a-check` gefahren, während beide Zeilen standen |
| „`d-check` führt drei distinkte SHAs mit je einem Tag-Kommentar" | fremder Bestand | `grep` über `.github/workflows/`, nach SHA gruppiert |
| „`actions/checkout` steht fünfmal mit identischem Kommentar" | fremder Bestand | derselbe `grep`, `uniq -c` |
| „unsere Fassung braucht 19 Zeilen" | eigener Bestand | `awk` über die beiden Funktionen des eigenen Sensors |

**Ein Satz wurde beim Durchgang umgeschrieben.** Der Entwurf nannte die Regel
*„in eurer Codebasis billig"* — eine Tatsache über ein fremdes System, für die
`a-check` keinen Handgriff hat. Sie steht jetzt als Aussage über die **eigene**
Fassung da, mit der Grenze im Satz. Das ist die Klasse, die ihr als *„gemessen
wird die eigene Menge, ausgesagt wird über die fremde"* führt.
