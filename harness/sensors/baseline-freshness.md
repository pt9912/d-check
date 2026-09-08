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
   GitHub liefert neueste zuerst. Liegt der Pin außerhalb dieses Fensters,
   findet die Liste ihn nicht — die Currency-Antwort lautet dann *unbestimmt*
   und **nicht** *„aktuell"*. **Gemessen am 2026-09-08: 58 Releases von 100** —
   der Fall ist terminiert, nicht hypothetisch, und tritt ein, sobald der Pin
   um mehr als 100 Releases zurückfällt oder das Repo diese Zahl überschreitet.
5. **Ein unbestimmter Currency-Stand meldet seit slice-215 Exit 3, nicht 0.**
   Vorher endete er mit 0 und war für den Nachtlauf — der ausschließlich den
   Exit-Code liest — von *„Pin ist der neueste Tag"* nicht zu unterscheiden.
   **Bruch-Test, beide Richtungen gefahren:** ein Pin, den es nicht gibt, ergab
   vorher Exit 0 mit derselben stderr-Zeile, jetzt Exit 3. **Die
   fail-open-Linie ist davon unberührt und das ist gemessen:** Ein Netz- oder
   API-Ausfall landet im `skip`-Zweig (Exit 0), nicht hier — geprüft in einer
   isolierten Kopie mit unerreichbarem API-Host.

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
