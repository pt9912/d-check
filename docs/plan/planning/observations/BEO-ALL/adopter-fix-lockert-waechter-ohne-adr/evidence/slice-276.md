**Vorgang:** slice-276
**Fund:** Die Übernahme der Stop-Hook-Fassung von `ai-harness-init` v0.10.0 hätte
das Handoff-Gate mit ihrem Default gelockert (ein Turn-Ende ohne neuen `HEAD` geht
auch mit ungedeckter Änderung frei) — `AGENTS.md` §3.6 verlangt dafür eine ADR. Das
fiel beim Lesen der beiden Hooks **vor** dem Code auf; der Auftraggeber entschied,
die Mechanik zu übernehmen und streng zu lassen. Der Ableiter trug diesmal: die
Lockerung wurde benannt statt still mitgenommen.
