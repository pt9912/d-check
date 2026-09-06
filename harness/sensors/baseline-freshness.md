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

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Pin aktuell und Inhalt unverändert — oder `SKIP` |
| 3 | Currency: ein neuerer Release-Tag existiert |
| 4 | Content-Drift: der gepinnte Tag trägt upstream andere Bytes |

Beides sind Exit-Codes des **Skripts**; `make` normalisiert einen
fehlgeschlagenen Recipe auf seinen eigenen Exit 2 — **welcher Teil** gemeldet
hat, sagt die Ausgabe, nicht der Exit.

## Bindung

**Netz**, fail-open, kein Gate, bewusst **nicht** in `gates`/`ci`. Bindepunkt ist der Nachtlauf
[`upstream-drift.yml`](../../.github/workflows/upstream-drift.yml), von
`ci.yml` getrennt, damit ein Upstream-Ausfall nie die CI rot färbt.
[`MR-011`](../conventions/done/MR-011-baseline-pin-release-tag.md)-Kette
