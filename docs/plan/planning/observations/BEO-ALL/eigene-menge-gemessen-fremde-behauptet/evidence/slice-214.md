**Vorgang:** slice-214
**Fund:** **Die Menge war um fünf Mitglieder zu klein, und der Fehler ist die
Suchmethode.** §3 des Plans schrieb die **Form** eines gepinnten
Fremd-Artefakts vorher aus — nicht im Repo entstanden · an fester Kennung
bezogen · geht in einen Lauf ein. Gesucht wurde dann an **drei Fundorten**
(`Dockerfile`, `a-check.mk`, `tools/semgrep.sh`) statt an der Form. Es fehlten:
drei SHA-gepinnte GitHub-Actions (für die
[`AGENTS.md`](../../../../../../../AGENTS.md) §3.9 eine eigene Hard Rule und §4
drei Freshness-Achsen führt), das digest-gepinnte `trivy`-Image und
`tools/archive-wave/Dockerfile` mit einem **anderen** `golang`-Digest. Sieben
statt zwölf. **Verschärfend:** Aus der zu kleinen Menge folgte der Schluss, das
`semgrep`-Regelset sei *„das einzige gepinnte Fremd-Artefakt außerhalb des
Repos"* — und der stand bereits als **Grenze in einem Gate-Vertrag**, als der
Review ihn widerlegte.
