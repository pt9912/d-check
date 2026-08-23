# Welle 82 — Drei Konfigurations-Flächen additiv geweitet, und eine Lehre über Reichweite — Closure-Notiz

**Welle:** welle-82-config-flaechen
**Abschluss:** 2026-08-23
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief · Steering-Loop · Register-Lese-Schritt.

Drei Konfigurations-Flächen, jede um genau die Kerbe geweitet, an der sie
nachweislich zu schmal war — plus die Auslieferung und zwei Doku-Nachträge, die
erst der Betrieb sichtbar gemacht hat:

- [slice-122](slice-122-versions-musterliste.md) — **`versions.patterns`:** eine
  **Liste** von Muster-Quellen-Paaren statt genau eines. Die Kurzform *ist* die
  einelementige Liste (ein Auswertungspfad); das **Paar** ist die Einheit, seine
  Ausnahmen sind paar-lokal, der Zeilen-Marker ist es nicht. Weil eine
  Befund-Adresse zwei Paare nicht unterscheidet, trägt **eine** Meldung alle
  Erwartungen mit ihrer Quelle.
- [slice-123](slice-123-structure-heading-muster.md) —
  **`structure.headings-match`/`headings-level`:** achte Bedingung, *jede*
  Überschrift des Abschnitts **positiv** geprüft, je Überschrift auf ihrer Zeile
  gemeldet (`section-heading-mismatch`), mit **derselben** Heading-Erkennung, mit
  der das Modul den Abschnitt findet.
- [slice-124](slice-124-diagrams-ventile.md) — **beide Ventile für `diagrams`**
  (`exempt-paths`, Zeilen-Marker) samt der §2-Schema-Zeilen, die dem Modul als
  einzigem fehlten. Der Marker ist dort ein **Token**, kein HTML-Kommentar, und
  wirkt auf der Öffnungszeile für den ganzen Block.
- [slice-125](slice-125-release-v0630.md) — **Release `v0.63.0`** samt
  Doku-Currency über alle drei Erweiterungen, Digest-Backfill — und der
  Korrektur einer Aussage, die **älter war als die Welle**: der Zeilen-Marker
  wirkt in **vier** Modulen, nicht in zweien.
- [slice-126](slice-126-handbuch-abschnitts-schnitt.md) — **nach** dem Release
  entstanden, auf Nachfrage des Auftraggebers: die ungesagte Ventil-Lage bei
  `citations` und ein §5-Abschnitt, dessen Überschrift zwei Fähigkeiten nannte,
  während er sechs Module über 183 Zeilen trug.

## Was hat funktioniert?

**Die Zusage der Welle hält:** ohne den jeweils neuen Schlüssel ist der
Befundsatz byte-identisch
([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)) — je
Erweiterung gegen das gepinnte Vorgänger-Image gemessen, grün wie rot. Die
**eine** benannte Ausnahme wurde nicht versteckt, sondern im Wellendokument
ausgesprochen: die Marker-Hälfte von slice-124 hängt am Zeilen-**Inhalt**, nicht
an einem Schlüssel; für sie gilt die Zusage nur für Bäume ohne die Zeichenfolge
in einer gelisteten Fence.

**Je Erweiterung eine konstruierte Gegenprobe**, die ohne sie stumm bliebe:
`versions` mit zwei Paaren gegen ein Vorgänger-Image, das die Liste gar nicht
kennt; die abgelöste Präfix-Negation gegen eine **eingerückte** Sektion (alte
Konstruktion Exit 0, neue Exit 1 auf der Zeile der Überschrift); und für
`diagrams` eine Konfiguration mit `exempt-paths`, die das Vorgänger-Image mit
Exit 2 zurückweist.

**Eine ADR für drei Entscheidungen** hat getragen. Die drei teilen eine
Begründung — eine Fläche additiv weiten statt eine Ersatz-Konstruktion pflegen —
und stehen als eigene Entscheidungen mit eigenen Konsequenzen und eigenen
Re-Evaluierungs-Triggern darin. Drei ADRs hätten dieselbe Begründung dreimal
getragen.

**Der Review vor dem Tag ist der wirksamste Schritt der Prozedur.** Fünf
unabhängige Läufe, **jeder** mit blockierendem Verdikt: 2 HIGH, 25 MEDIUM,
22 LOW, 10 INFO. Beide HIGH betrafen stille Falsch-Negative, die alle Gates
passiert hatten — in slice-122 verschwand eine zweite Erwartung in der
geteilten Nachrunde, in slice-123 war ein Schlüsselname mit einem bestehenden
Selektor kollidiert.

## Was ging anders als geplant?

**Die Welle bekam einen fünften Slice.** [slice-126](slice-126-handbuch-abschnitts-schnitt.md)
entstand **nach** der Closure von [slice-125](slice-125-release-v0630.md) und
nach dem Release — aus einer Nachfrage, nicht aus der Planung. Der
Closure-Trigger dieser Welle sagte weiterhin „alle **vier** Slices" und wäre
damit erfüllbar gewesen, während ein Slice offen lag; der Review hat es
gefunden, die Zahl folgt jetzt §4.

**Die Beobachtungs-Klasse dieser Welle war ihre eigene Arbeitsweise.**
Sechsmal wurde eine Aussage oder Regel aus dem **Anlass** gezogen statt aus dem
**Bestand** — und zwar in wachsender Verschachtelung: erst falsche
Exklusivitäts-Aussagen (slice-123, slice-124), dann ein *Kriterium*, das eine
bereits korrigierte Menge nachträglich begründete (slice-125), dann drei
Instanzen im selben Commit, der die Klasse als
[`BEO-011`](../observations.md) ins Register schrieb (slice-126). Der Zähler
dieses Eintrags war anfangs selbst aus dem Anlass gebildet.

**Eine Klasse aufzuschreiben schützt nicht davor, sie zu begehen.** Das ist die
härteste Erkenntnis der Welle und der Grund, warum die mechanische Form von
[`BEO-011`](../observations.md) ausdrücklich **offen** bleibt statt aus dem
Zähler abgeleitet zu werden. Was in allen sechs Fällen gegriffen hat, war nicht
die Regel, sondern ein zweiter Leser mit Prüfauftrag.

**Der Release-Prep hat einen Blind-Spot sichtbar gemacht, den die Checkliste
nicht kennt.** Modul-Listen, §11-Position, Optionen-Tabelle und fenced
Config-Beispiele stehen dort — *„Prosa-Aussage über eine Modul-Menge"* nicht.
Genau dort lagen beide MEDIUM des Release-Reviews.

## Steering-Loop-Einträge

- **Reviewer-Skill:** der MEDIUM-Anker *„Botschaft verallgemeinert über die
  Messung hinaus"* (Version 1.9.0, aus welle-81) hat in dieser Welle mehrfach
  getroffen. Kein neuer Anker nötig — die Klasse dieser Welle ist als
  [`BEO-011`](../observations.md) registriert, nicht als Skill-Regel, weil sie
  ein Zähl-Auftrag ist und kein Lese-Anker.
- **Release-Prozedur** ([`releasing.md`](../../../user/releasing.md)): die
  Anti-Anlagerungs-Regel gilt jetzt der **Klasse** (jeder gegliederte
  Fließtext), nicht mehr zwei benannten Kapiteln.
- **Nicht** mechanisiert: ein Wortlisten-Lint auf Exklusivitäts-Wörter wäre ein
  Heuristik-Wächter mit hoher Falsch-Positiv-Last in einem Repo, dessen Doku
  bewusst Grenzen benennt — und er deckte nur eine der drei Ausprägungen.

## Beobachtungs-Register (Zeiger)

Gelesen zur Closure ([`observations.md`](../observations.md)):

- **[`BEO-011`](../observations.md) neu angelegt**, Zähler **3**, Schwelle
  erreicht — die Klasse dieser Welle. Drei weitere Instanzen aus
  [slice-126](slice-126-handbuch-abschnitts-schnitt.md) sind dort benannt und
  bringen die dritte Ausprägung erstmals in die Belege.
- **[`BEO-008`](../observations.md) bleibt bei 3, und seine benannte
  mechanische Form ist seit dieser Welle *baubar*** — `versions.patterns` trägt
  den zweiten Abgleich. Gebaut ist er **nicht**: ob das eigene Profil ihn fährt,
  ist ein eigener Entscheid mit eigener Messung und ausdrücklich kein Rest
  dieser Welle.
- **[`BEO-002`](../observations.md) hat wieder getroffen** und bleibt offen: der
  `grep` nach dem **alten** Wortlaut, repo-weit, stand seit
  [slice-125](slice-125-release-v0630.md) als Ableiter in dessen Closure-Notiz
  und wurde in [slice-126](slice-126-handbuch-abschnitts-schnitt.md) wieder
  nicht gefahren — sechs Fundstellen blieben stehen, zwei davon in den
  Spec-Straten. Eine Regel für Menschen kann weiter verfehlt werden.
- **[`BEO-009`](../observations.md)** traf einmal in Richtung (b): eine
  Botschafts-Zahl benannte die gate-gedeckte Menge mit dem Namen der
  gate-blinden. Zähler unverändert bei 3.
- Alle übrigen Einträge unverändert; keine Streichung.

## Folge-Slices

- **Der `diagrams.scope`-Rückbau im eigenen Profil** — die Ventile machen ihn
  möglich; ob er trägt, ist zu **messen**, nicht zu vermuten. Die
  Bestandsaufnahme aus slice-124 („ohne Scope null Befunde") kann *„die Ventile
  helfen"* nicht von *„es hat nie etwas gefeuert"* trennen.
- **Die 3×-Form von [`BEO-008`](../observations.md)** als zweites
  Muster-Quellen-Paar im eigenen Profil.
- **Ob `citations` eine feine Ventil-Achse bekommen soll** — ein Change Request
  am Lastenheft, kein Doku-Slice.

## Verifikation

- `make gates` nach jedem Slice grün; `make fullbuild` zu jeder Closure grün —
  Exit-Codes explizit gelesen, nie hinter einer Pipe
  ([`BEO-007`](../observations.md)).
- **Release ausgeliefert und gegen das gezogene Image belegt:** Tag `v0.63.0`,
  Release-Pipeline `32619522094` success, Digest
  `sha256:7049cefd2d91b367b72f15f789123ab5f51bf09150ad9cd262a9a945cfceb16e`,
  OCI-Label `org.opencontainers.image.version` = `0.63.0`, netzloser Smoke über
  den Digest-Pin gegen dieses Repo Exit 0.
- Fünf unabhängige Reviews, alle blockierend, alle Befunde eingearbeitet und
  jede Einarbeitung nachgemessen.
- Der Anker in [`version.md`](../../../../version.md) ist gewandert — genau
  einer, unabhängig nachgezählt.
- Die fenced YAML-Beispiele, die kein Gate sieht, wurden einzeln gegen den
  Validator gefahren: wie gedruckt Exit 1, in der Mischform Exit 2 mit der
  dokumentierten Ursache, in der einkommentierten Alternative Exit 1.
