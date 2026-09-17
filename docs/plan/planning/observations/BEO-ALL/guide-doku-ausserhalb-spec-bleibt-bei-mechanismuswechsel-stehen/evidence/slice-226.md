`slice-226` änderte den Pack-Sichtbarkeits-Mechanismus in `vcs`/`commits`
(ein Pack mit gültigem Hash-Suffix und passendem Index wird jetzt
unabhängig vom Namens-Präfix gelesen) und zog `spec/spezifikation.md`
korrekt nach. Der unabhängige Review (R1-F-1) fand, dass
`harness/sensors/adr-check.md`, `harness/sensors/trace-check.md` und zwei
Abschnitte von `docs/user/benutzerhandbuch.md` weiterhin den vorherigen
Stand beschrieben (jedes Fremdpräfix unsichtbar, Exit 2 ohne Ausnahme). Im
selben Slice nachgezogen; siehe dessen Closure-Notiz (§7).
