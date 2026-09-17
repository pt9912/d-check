# Eingehender Change Request — `vcs` liest nur Packs mit dem Präfix `pack-`

**Absender:** Adopter `ai-harness-init` · **Eingegangen:** 2026-09-17
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) (Modul `vcs`) und jedes weitere Modul, das die Historie über denselben Objektspeicher liest
**Stand:** **entschieden und umgesetzt am 2026-09-17** — Bitte angenommen,
Träger [`slice-226`](../planning/done/wellenlos/slice-226-vcs-pack-alias-fremdes-praefix.md)
<!-- d-check:status-provenance --> (Belege unten).

**Ablage-Hinweis.** Ein **eingehender** CR ist die dritte Klasse neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf). Die Datei
liegt hier aus demselben Grund wie ihre Vorgänger in diesem Verzeichnis. Der
Absender hat sie abgelegt, aber nicht committet; Commit und Entscheid liegen bei
diesem Repo.

---

## Wortlaut

> **Betreff:** Ein Repo mit Packs unter anderem Präfix lässt `vcs` abbrechen,
> obwohl die Objekte vorhanden sind.
>
> **Einreicher:** ai-harness-init (Adopter), 2026-09-17, gegen d-check
> `@sha256:1470ecdcaa686a5ef4513dee9b0ae522586f54b87d568b06fc6b5b2741b633b3`
> (`v0.76.1`).
>
> ### Anlass
>
> Der Objektspeicher des Adopter-Repos trug neben zwei Packs mit dem Präfix
> `pack-` fünf Packs mit dem Präfix `loose-`, so wie sie
> `git maintenance run --task=loose-objects` anlegt. Unter `v0.76.1` bricht
> ein history-lesender Lauf daran ab:
>
> ```
> make adr-immutable RANGE=8ae647cc~1..8ae647cc
> # d-check … --enable vcs … --range 8ae647cc~1..8ae647cc
> # → Exit 2, „nicht lesbarer Unterbaum ".claude/hooks""
> ```
>
> Unter `v0.76.0` meldete derselbe Aufruf `0 Befund(e)`, Exit 0 — der stille
> Durchgang, den
> [`CO-001`](../carveouts/done/CO-001-vcs-range-stiller-skip.md) beschreibt.
>
> ### Messung
>
> | Lage des Objektspeichers | `v0.76.0` | `v0.76.1` |
> |---|---|---|
> | Packs mit Präfix `pack-` **und** `loose-` | Exit 0, `0 Befund(e)` | Exit 2, „nicht lesbarer Unterbaum" |
> | nach `git repack -a -d` (nur noch Präfix `pack-`) | Exit 0, `0 Befund(e)` | Exit 0, `0 Befund(e)` |
> | frischer Klon (`git clone --no-local`) | Exit 0, `0 Befund(e)` | Exit 0, `0 Befund(e)` |
>
> Die Objekte sind in allen drei Lagen vorhanden: `git cat-file` liest sie, und
> `git fsck --connectivity-only` meldet keinen fehlenden Eintrag. Der Unterschied
> liegt allein im Namen der Pack-Datei.
>
> **Vermutete Ursache:** Der Objektspeicher von go-git (`v5.19.2`) sucht Packs
> nur unter dem Präfix `pack-`. Das hat der Adopter nicht im Quelltext
> nachgeprüft; belegt ist nur das Verhalten oben.
>
> ### Bitte
>
> Lies jede Pack-Datei im Objektspeicher, zu der ein gültiger Index gehört,
> **unabhängig vom Präfix** ihres Namens. Das gilt für jedes Modul, das die
> Historie liest.
>
> Der fail-closed-Abbruch aus `v0.76.1` bleibt für ein wirklich fehlendes
> Objekt richtig; er soll nur nicht greifen, wenn das Objekt in einem anders
> benannten Pack liegt.
>
> ### Gegenprobe für den Test
>
> Ein Repo, dessen Objekte (auch) in einem Pack mit dem Präfix `loose-` liegen,
> meldet unter `vcs` dieselben Befunde wie derselbe Stand nach
> `git repack -a -d`. Das rote Gegenbeispiel: ein wirklich entferntes Objekt
> bricht weiterhin mit Exit 2 ab.
>
> ### Was ausdrücklich **nicht** gebeten ist
>
> - **Den Abbruch zu lockern**, wenn ein Objekt wirklich fehlt.
> - **Den Objektspeicher des Adopters zu verändern.** Ein automatisches
>   Umpacken gehört nicht ins Werkzeug; es liest nur.
>
> ### Wie der Adopter bis dahin arbeitet
>
> Er hat sein Repo einmal umgepackt; danach laufen die history-lesenden Ziele
> wieder. Das hält, bis erneut Packs unter dem Präfix `loose-` entstehen.
> Frische Klone, wie seine CI sie fährt, sind nicht betroffen.

## Antwort

Die Ursachen-Vermutung trifft zu: go-gits `DotGit`-Schicht (`v5.19.2`)
entdeckt und öffnet einen Pack ausschließlich über den rekonstruierten
kanonischen Namen `pack-<hash>.{pack,idx}` — ein valider Pack unter anderem
Präfix (z. B. `loose-<hash>.pack`, wie `git maintenance
run --task=loose-objects` es anlegt) bleibt für sie unsichtbar, unabhängig
davon, ob das Objekt selbst lesbar ist. Empirisch reproduziert: derselbe
Objektspeicher meldet vor und nach `git repack -a -d` identische Befunde,
sobald diese Auflösung ergänzt ist.

Der `vcs`-Port löst seit [`slice-226`](../planning/done/wellenlos/slice-226-vcs-pack-alias-fremdes-praefix.md)
<!-- d-check:status-provenance --> einen Pack zusätzlich unter seinem
kanonischen Namen auf, wenn seine Datei einen gültigen SHA1/SHA256-Hash als
Namens-Suffix trägt **und** eine passende `.idx`-Datei existiert — rein
lesend, ohne den gemounteten Objektspeicher zu verändern
([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Ein Pack ohne ein solches Suffix oder ohne passenden Index bleibt weiterhin
unsichtbar, und die Enumeration bricht fail-closed ab, wenn er die einzige
Quelle eines benötigten Objekts ist — genau die Grenze, die die Bitte
ausdrücklich zieht.

Damit ist erfüllt:

- **Die Bitte** — jede Pack-Datei mit gültigem Index wird unabhängig vom
  Namens-Präfix gelesen, für jedes Modul, das über den `vcs`-Port liest (der
  Fix sitzt im Adapter selbst, nicht in einem einzelnen Modul).
- **Die Gegenprobe** — ein Repo mit `loose-`-präfigierten Packs meldet
  dieselben Befunde wie derselbe Stand nach `git repack -a -d` (zwei
  Go-Tests, `TestAllPathsPackUnterFremdemPraefix`); das rote
  Gegenbeispiel bleibt rot (`TestAllPathsPackMitUnbrauchbaremPraefixBleibtFehlerhaft`
  — ein Pack ohne gültiges Hash-Suffix wird weiterhin nicht gelesen).
- **Die beiden Nicht-Bitten** — der fail-closed-Abbruch für ein wirklich
  fehlendes Objekt ist unverändert, und der Objektspeicher des Adopters wird
  nie beschrieben.

Der Adopter braucht sein Workaround (einmaliges manuelles Umpacken) ab dem
Release, das diesen Fix trägt, nicht mehr.
