# `make baseline-freshness` — auditiert den Baseline-Pin gegen upstream, in zwei getrennten Teilen

## Vertrag

Upstream-Audit via `--check-latest`, **zwei getrennte Teile:**

- **(A) Currency** — die Release-**Liste** des Kurs-Repos gegen den Pin.
  Bewusst nicht `releases/latest`: das überspringt Prereleases und verbirgt
  einen zurückgezogenen Pin. Exit 3.
- **(B) Content-Drift** am **gepinnten** Tag — die Bytes beider Bäume des
  Release-Assets gegen das committete `SHA256SUMS`. Exit 4.

**Integrität ist nicht Aktualität:** `--verify` (siehe
[`baseline-verify`](baseline-verify.md)) und `--check-latest` beantworten
verschiedene Fragen.

## Grenze — was das Grün nicht abdeckt

1. **Fail-open je Teil** — Netz-, Werkzeug- oder Manifest-Ausfall ergibt
   `SKIP`, mit Zeitgrenzen, damit eine hängende Verbindung nicht zur Job-Decke
   läuft.
2. **Gemeldet wird, nicht gehoben** — die Hebung bleibt ein bewusster Akt.
3. **Nicht [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)** —
   jene Zusage gilt dem **Produkt**; dass `gates` netzlos bleibt, ist eine
   Eigenschaft dieses Repos, keine Produktzusage.
4. **Die Currency-Hälfte sieht nur die jüngsten 100 Releases.** Die
   Release-Liste wird mit `per_page=100` gelesen und **nicht paginiert**;
   GitHub liefert neueste zuerst. Der Pin fällt aus dem Fenster, sobald **mehr
   als 100 Releases über ihm** liegen — die Gesamtzahl des Repos allein genügt
   dafür nicht. **Gemessen am 2026-09-08: 58 Releases insgesamt**, der Pin ist
   der neueste; der Abstand ist also 0.
5. **Ein unbestimmter Currency-Stand meldet seit slice-215 Exit 3, nicht 0.**
   Vorher endete er mit 0 und war für den Nachtlauf — der ausschließlich den
   Exit-Code liest — von *„Pin ist der neueste Tag"* nicht zu unterscheiden.
   **Bruch-Test, beide Richtungen gefahren:** ein Pin, den es nicht gibt, ergab
   vorher Exit 0 mit derselben stderr-Zeile, jetzt Exit 3.
   **Die fail-open-Linie bleibt unberührt, und das ist an der schwierigen
   Ausfall-Form gemessen:** Ein **abgebrochener** Transfer liefert Teildaten
   *und* einen Fehlerstatus (gemessen: 35 010 von 322 579 Bytes, `curl`-Exit
   28). Die erste Fassung wertete nur den Pipeline-Ausgang und hielt die
   abgeschnittene Liste für vollständig — der Pin fehlte darin und der Lauf
   meldete **Exit 3 auf einen Netzausfall**. Seit dem Nachzug wird `curl`s
   Status getrennt geprüft; derselbe Abbruch ergibt jetzt `SKIP`, Exit 0.
6. **Gemessen werden Release-*Objekte*, nicht Tags.** Ein Tag, der geschoben
   ist, bevor (oder ohne dass) sein Release-Objekt existiert, kommt in der
   Liste nicht vor — der Pin gilt dann als *aktuell* und der Lauf meldet
   **Exit 0**. **Das ist die Lage der offenen Beobachtung**
   [`BEO-HARN/check-latest-blind-before-pin`](../../docs/plan/planning/observations/BEO-HARN/check-latest-blind-before-pin/observation.md),
   deren Ursache nicht feststeht; sie ist damit **nicht ausgeschlossen**,
   sondern die naheliegendste verbliebene Erklärung.
7. **Prereleases fallen durch den Filter.** Die Tag-Auswahl verlangt
   `vX.Y.Z` — ein `v7.0.0-rc1` wird verworfen (gemessen). Die Wahl des
   Listen-Endpunkts ist im Skript-Kopf ausdrücklich damit begründet, dass
   `releases/latest` *„Prereleases überspringt"*; der Filter nimmt diesen
   Gewinn wieder zurück. **Für einen Prerelease-Pin ist die Prüfung damit
   blind**, und die `**Stand:**`-Zeile könnte einen solchen Pin tragen, ohne
   dass etwas meldet.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Pin aktuell und Inhalt unverändert — oder `SKIP` |
| 3 | Currency: ein neuerer Release-Tag existiert — **oder** der Pin steht nicht in der Liste (Stand *unbestimmt*, `seit slice-215`) |
| 4 | Content-Drift: der gepinnte Tag trägt upstream andere Bytes |

Beides sind Exit-Codes des **Skripts**; `make` normalisiert einen
fehlgeschlagenen Recipe auf seinen eigenen Exit 2 — **welcher Teil** gemeldet
hat, sagt die Ausgabe, nicht der Exit.

## Bindung

**Netz**, fail-open, kein Gate, bewusst **nicht** in `gates`/`ci`. Bindepunkt ist der Nachtlauf
[`upstream-drift.yml`](../../.github/workflows/upstream-drift.yml), von
`ci.yml` getrennt, damit ein Upstream-Ausfall nie die CI rot färbt.
[`MR-011`](../conventions/done/MR-011-baseline-pin-release-tag.md)-Kette
