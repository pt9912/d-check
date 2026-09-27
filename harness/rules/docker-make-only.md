# Docker/make-only — Begründung, Durchsetzung, Grenzen

Ausgelagerter Körper zu `AGENTS.md` §3.1. Der operative Kern (Falsch/Richtig,
Host-Klasse) steht dort; diese Datei trägt Begründung, Durchsetzungsschichten
und ihre Grenzen — Auslagerungs-Muster aus
[ADR-0096](../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

## Host-Klasse ist eine Zusage, keine vollständige Liste

**Auch keine Host-Skript-Interpreter** (`python`, `perl`, `ruby`, `node`,
`uv`, …) ([`MR-040`](../conventions.md#mr-040)). Datei-Änderungen macht das
Werkzeug ohne Shell; für Messungen gilt die Rangfolge **Produkt vor
`grep`/`awk` vor allem anderen**. Eine Stage für Skripte gibt es nicht, und
diese Regel verweist auch nicht auf eine: ein Fall, der `bash` und die
genannte Host-Klasse übersteigt, ist ein **Entscheid** — keine vierte
Toolchain nebenbei ([`MR-046`](../conventions.md#mr-046)).

**Die Host-Klasse ist die Zusage an den Agenten, nicht die vollständige
Liste dessen, was Gate-Skripte rufen.** Die **Netz**-Targets holen sich
mehr: [`fetch-baseline-cache.sh`](../../tools/harness/fetch-baseline-cache.sh)
braucht `curl` und `unzip`,
[`pin-freshness.sh`](../../tools/harness/pin-freshness.sh) braucht `curl`,
[`nightly-state.sh`](../../tools/harness/nightly-state.sh) ebenso, und
[`image-scan.sh`](../../tools/image-scan.sh) braucht **Docker mit Netz**
(Trivy zieht seine Vuln-DB). **Alle vier** stehen bewusst außerhalb von
`gates`; wer sie fährt, fährt sie mit dieser zusätzlichen Erwartung. **Die
ersten drei sind fail-open, das vierte nicht** — ein gescheiterter CVE-Scan
meldet Exit 2 und ausdrücklich keinen grünen Befundstand
([ADR-0066](../../docs/plan/adr/0066-cve-scan-gegen-das-publizierte-image.md)).

## Begründung

Toolchain-Reproduzierbarkeit + Supply-Chain-Defense.

**Die Regel gilt unabhängig von ihrer Durchsetzung.** Wer sich auf den
Wächter verlässt, verlässt sich auf nichts.

## Durchsetzung, zwei unabhängige Schichten

Die zweite ist eine **Permission-Sperrliste** in
[`.claude/settings.json`](../../.claude/settings.json): sie hängt an keinem
Hook, matcht aber den **ganzen** Befehl ab dem Anfang und sieht deshalb
weder Präfixe noch Sub-Shells noch zusammengesetzte Kommandos. Ihre
git-/docker-Hälfte hat keine zweite Schicht unter sich. Grenzen und
Nicht-Zusagen: [`MR-047`](../conventions.md#mr-047).

Die erste ist ein Tool-Call-Wächter
([`.claude/hooks/pretooluse-command-guard.sh`](../../.claude/hooks/pretooluse-command-guard.sh)):
er prüft die Befehlsposition jedes Segments und Sub-Shell-Strings rekursiv.
Er ist **werkzeug-lokal**, kein Repo-Gate: keine CI ruft ihn, ein Lauf ohne
dieses Werkzeug ist ungebunden. Er ist in `bash` geschrieben und liest die
Hook-Eingabe mit `awk`, läuft also in derselben Klasse, die er durchsetzt
([`MR-042`](../conventions.md#mr-042)); `make guard-probe`
([`harness/README.md`](../README.md) §Sensors) fährt ihn gegen seine Proben.

## Grenze — Stolperdraht, keine Sandbox

Ungeprüft bleiben: ein Shell-Schlüsselwort als Segment-Kopf (`if … then
pip …`), ein Wrapper außerhalb seiner Präfix-Liste (`nohup`, `timeout`), ein
wort-interner Quote-/Backslash-Splice (`p"i"p`) und escapte Quotes in der
Verschachtelung erreichen ein gelistetes Werkzeug, ohne dass er es als Kopf
sieht. Die Regel gilt trotzdem — sie hängt nicht an ihm. Tabelle in
[`MR-042`](../conventions.md#mr-042).
