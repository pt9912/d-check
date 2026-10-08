**Vorgang:** slice-258
**Fund:** Trotz Prüffrage 19 (seit slice-257) im selben Muster: jede
Behebung der Black-Box-Probe schloss die gerade gemessene Fehlerform — eine
Docker-Meldung, einen Exit-Code, ein gleiches Scheitern — und ließ den
nächsten Ausgang offen (R1 F-1, V-1, R2-1, R3, R4-1). Erst R4 benannte die
eigentliche Fallmenge: die Mount-Quellen, nicht die Ausfallarten. Die
Prüffrage wirkte beim Reviewer, nicht beim Implementer, der sie vor dem
Commit hätte stellen müssen.
