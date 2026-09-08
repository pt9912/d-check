**Vorgang:** slice-214
**Fund:** **Fünfte Instanz, und die erste, in der die falsche Grenze bereits
im Gate-Vertrag stand.** Der Slice schrieb `harness/sensors/semgrep.md` eine
neue Grenze 3 — und sie war in ihrer Kernaussage falsch: *„das einzige
gepinnte Fremd-Artefakt außerhalb des Repos"* (es sind elf weitere) und
*„genau einmal geprüft"* (gilt lokal, nicht in CI, wo der Runner keinen Cache
mitbringt). Dazu ein Superlativ ohne Kriterium — *„die stärkste Bindung im
ganzen Pin-Bestand"*, gemessen ein SHA-1 gegen sechs SHA-256-Digests.
**Der Ableiter des Eintrags war angewandt und hat nicht gereicht:** Die Grenze
war gegen das **Skript** geprüft (`[ ! -d … ]` steht wirklich dort) — aber ihre
Aussage über den **Bestand** ruhte auf einer Inventur, die um fünf Mitglieder
zu klein war. **Gegen den Gegenstand zu prüfen genügt nicht, wenn der
Gegenstand größer ist als die Menge, die man gezählt hat.**
