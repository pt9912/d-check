# Die Modulliste des eigenen Prüf-Profils hat Spiegel außerhalb der Config, und keiner davon ist gate-gedeckt

**Sub-Area:** `.d-check.yml`, `Makefile`, Gate-Doku

Wer ein Modul in `.d-check.yml` aufnimmt, ändert damit still drei weitere
Stellen mit: `FOCUS_DISABLE` im `Makefile`, die Netzlos-Modulliste im
[`DC-QA-03`](../../../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Test, die Prosa-Beschreibungen des Gates. `gate-consistency`
prüft Target-Namen, nicht Modul-Mengen.
