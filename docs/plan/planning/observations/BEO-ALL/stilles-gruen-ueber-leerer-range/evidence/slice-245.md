Die stille-Grün-Gegenprobe von slice-245 (shallow-Klon Tiefe 1, lokal über
`file://`) maß am eigenen Adapter: das Modul `vcs` meldet auf eine leere,
auflösbare Range (`HEAD..HEAD`) „0 Befund(e)", Exit 0 — stills Grün über
leerem Prüfbereich. Das Modul `commits` bricht auf derselben Range laut ab
(Exit 2), eine unauflösbare Basis ebenfalls (Exit 2); der Vorlauf-Wächter
history-range-guard fängt beide Fälle vorab (Exit 1 bzw. 2). Der
Verifier reproduzierte alle Messungen unabhängig und ergänzte die positive
Kontrolle: `make trace-check RANGE=HEAD..HEAD` im shallow-Klon bricht am
Wächter mit Exit 2 ab — dieselbe Range, die das Modul ohne Wächter still
grün meldete. Der Produkt-Fix trägt das Anforderungs-Delta im Folge-Slice
slice-247 (in `open/` geschnitten).
