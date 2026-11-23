# CHANGELOG.md wird in der Release-Prep gepflegt, nicht im Feature-Commit

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

`CHANGELOG.md` wird bei nutzersichtbaren Änderungen gepflegt — **in der
Release-Prep, nicht im Feature-Commit.** Die Datei führt **keinen**
`[Unreleased]`-Abschnitt: jeder Eintrag steht unter seiner Versions-Nummer,
und die steht erst fest, wenn das Release geschnitten wird. Ein Slice, der
seine Zeile vorzieht, muss sie beim Bump wieder anfassen. Dieselbe Grenze
gilt den beiden `README*.md` und dem Handbuch-Kopf. Ein fehlender Eintrag
dort im Feature-Commit ist deshalb kein Rückstand.

**Ausgenommen sind Modul-Aufzählungen, die ein Spiegel-Test gegen die
Modul-Registrierung hält.** Sie kommen mit dem Feature-Commit, der das Modul
registriert — sonst ist `make test` dort rot. Welche Stellen das sind, sagt
der Test selbst, nicht eine Liste hier:
[`registry_mirror_test.go`](../../../internal/hexagon/core/app/registry_mirror_test.go)
nennt seine Fundorte (darunter die Modul-Bullets beider `README*.md` und die
Modul-Tabelle des Handbuchs), dazu `TestPrintConfigVerfuegbarDecktRegistry`
für die „Verfügbar“-Zeile von `--print-config`. Der Handbuch-Kopf, die
übrige README-Prosa und die Dogfooding-Zeile bleiben Release-Prep. *(Regel
aus dem Steering Loop,
[`BEO-ALL/regeltext-gegen-spiegel-test`](../../../docs/plan/planning/observations/BEO-ALL/regeltext-gegen-spiegel-test/observation.md),
seit welle-95; Auflösungs-Trigger: die Spiegel-Tests halten keinen Fundort in
einer README oder im Handbuch mehr — dann entfällt dieser Absatz.)*
