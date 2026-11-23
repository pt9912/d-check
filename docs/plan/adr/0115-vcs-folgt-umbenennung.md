# ADR-0115: Das Modul `vcs` folgt einer Umbenennung der Range — ein Link-Nachzug auf die umbenannte Datei ist keine Drift

**Status:** Proposed

**Datum:** 2026-11-23

**Autor:** pt9912

**Bezug:** [Issue #18](https://github.com/pt9912/d-check/issues/18) (Layout-Migration eines Adopters,
`spezifikation.md` wird `specification.md`),
[`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
(Absatz „Pfad-Nachzug (opt-in)“, Kriterium „Boundary (Pfad-Nachzug)“, Out-of-Scope),
[`DC-FA-VCS-002`](../../../spec/lastenheft.md#dc-fa-vcs-002--leere-commit-range-im-modul-vcs-ist-laut-zu-melden-opt-in)
(der Vorfahren-Walk),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[ADR-0024](0024-vcs-immutable-gate.md) (Modul `vcs`, VCS-Port, eigenes ADR-Gate),
[ADR-0103](0103-adr-gate-laesst-pfad-nachzug-durch.md) (`vcs.ignore-link-targets`, das diese ADR erweitert,
nicht ablöst), [ADR-0109](0109-ci-push-range-faellt-bei-force-push-zurueck.md) (die Range der CI nach einem
Force-Push), [`AGENTS.md`](../../../AGENTS.md) §3.3 (reiner Move-Commit) und §3.6 (Lockerung nur per ADR),
slice-291 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-VCS-001.a`](../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
Schritt 2 (Commit-Menge, Basis kein Vorfahr, welche Bäume gelesen werden) und Schritt 4 (Gleichsetzen der
Link-Ziele) sowie die Schema-Zeile des neuen Schlüssels in
[`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml).

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v6.18.1/regelwerk/modul-04-adrs.md).

---

## Kontext

Eine `Accepted`-ADR verlinkt eine Datei, die umbenannt wird. Das Baseline-Regelwerk nimmt den Nachzug des
Links von der Immutabilität aus (`modul-04-adrs.md` §Nachzug ist keine Überschreibung); dieses Repo hat
dafür `vcs.ignore-link-targets` eingeführt ([ADR-0103](0103-adr-gate-laesst-pfad-nachzug-durch.md)). Der
Schlüssel normiert ein auflösendes Ziel auf Dateiname und Anker. Ein Nachzug auf einen **anderen** Namen
bleibt deshalb `core-drift-vcs`: das Werkzeug weiß nicht, ob der Link nachgezogen oder auf ein anderes
Dokument umgebogen wurde. Eine Layout-Migration mit neuen Dateinamen trifft so jede `Accepted`-ADR, die
die Datei verlinkt; im eigenen Bestand verlinken 59 von 107 `Accepted`-ADRs die Spezifikation (Zählung
2026-11-23, Kommando im Plan des Anforderungs-Slice).

Die Identität, die das Werkzeug nicht kennt, steht in der git-Historie: ein Schritt, in dem ein Pfad
verschwindet und ein anderer mit **demselben Blob** erscheint, ist eine Umbenennung ohne Urteil. Der
Auftraggeber hat am 2026-11-23 die acht Entscheide des Anforderungs-Slice mit den Empfehlungen bestätigt
(Entscheidung 1–8 unten). Das Architect-Verdikt vom selben Tag (unter `docs/reviews/`) hat am
Gegenstand Stellen gefunden, die der Code sonst stillschweigend entschiede; sie stehen als Entscheidung
9–14; der Auftraggeber hat sie am 2026-11-23 bestätigt, ebenso Entscheidung 5 in der seiten-gleichen Fassung.

Am Gegenstand gelesen, auf das die Entscheidung baut:

- Das Modul vergleicht den Core **jeder** geschützten Datei, die in beiden Ständen liegt, nicht nur der
  geänderten, und zwar mit **einer** Auflösungs-Funktion für beide Fassungen
  (`internal/hexagon/core/rules/vcs.go`, Schleife über die geschützten Pfade und `vcsModified`). Eine
  Normierung, die nur auf einer Seite greift, macht eine unveränderte Datei rot.
- Der git-Adapter verlangt für `vcs` heute **keine** Abstammung der Basis; er prüft nur den Leerfall und
  wandert dafür die Vorfahren der Basis (`internal/adapter/driven/git/git.go`, `AllPaths`). `LineChanges`
  (Modul `metrics`) verlangt die Abstammung und bricht sonst mit Exit 2 ab.
- `CommitMessages` liefert die Nicht-Merge-Commits „von HEAD erreichbar, von BASE nicht“, **nach
  Commit-SHA sortiert** — eine deterministische, aber keine topologische Reihenfolge.
- Tree- und Index-Einträge tragen die Blob-Kennung; die Gleichheit zweier Einträge ist ohne Lesen des
  Blobs entscheidbar.
- Die Bibliothek (go-git v5.19.3, im Modul-Cache des Test-Images gelesen, nicht auf dem Host) bietet eine
  Umbenennungs-Erkennung an: Vorbelegung `DetectRenames: true`, `RenameScore: 60`; mit `OnlyExactRenames`
  wählt sie unter mehreren Kandidaten gleichen Blobs den mit dem ähnlichsten Namen und verlangt gleichen
  Modus. Beides weicht von Entscheidung 1 und 3 ab. Der Diff-Walker derselben Bibliothek ist die Stelle,
  die ein unlesbares Unterbaum-Objekt früher still als Ende las (Historie der Spezifikation zu
  [`DC-FA-VCS-001.a`](../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2).

## Entscheidung

Wir geben dem Modul `vcs` einen **eigenen opt-in-Schlüssel `vcs.follow-renames`**, mit dem es die
Umbenennungen der Range aus der git-Historie liest und ein Link-Ziel über sie gleichsetzt. Im Einzelnen
(die Nummern dienen dieser Liste als Bezug; 1–14 vom Auftraggeber am 2026-11-23 bestätigt, Entscheidung 5 in
der seiten-gleichen Fassung nach Entscheidung 9; danach ohne erneute Bestätigung an die Spezifikation
angeglichen — Wortlaut von 1, 2, 3, 9 und 13 — und nach den Reviews an 11 und 13 geschärft; Entscheidung 9 und 10
sind nach Review R2 in der Sache geändert und vom Auftraggeber am 2026-11-23 in dieser Fassung bestätigt, siehe Geschichte):

1. **Erkennung: Inhalts-Gleichheit, keine Ähnlichkeits-Schwelle.** Ein Pfad A, der vor einem Schritt
   existiert und danach nicht, und ein Pfad B, der danach existiert und vorher nicht, sind eine
   Umbenennung A → B genau dann, wenn beide dasselbe Blob tragen (git: `R100`) — gleiche Einträge im Sinne
   von Entscheidung 13 (Blob-Kennung und Art).
2. **Schritt: jeder Commit der Range, mit Kette.** Ein Schritt ist ein Commit mit genau einem Eltern-Commit,
   der von HEAD aus erreichbar ist und von BASE aus nicht (aus derselben Menge wie das Modul `commits`),
   gelesen gegen diesen Eltern-Commit; ein Merge-Commit ist kein Schritt, ein Wurzel-Commit trägt keine
   Umbenennung bei (vor ihm existiert kein Pfad). Im `--staged`-Modus ist der einzige Schritt HEAD → Index.
   Umbenennungen verketten sich (A → B, später B → C ergibt A → C).
3. **Mehrdeutigkeit bleibt Drift.** Tragen in einem Schritt mehrere verschwundene oder mehrere neue Pfade
   einen gleichen Eintrag (Entscheidung 13), oder führt die Kette von A zu zwei verschiedenen Enden, gilt für A
   keine Umbenennung. Ein Merge-Commit trägt selbst keine Umbenennung bei.
4. **Eigener Schalter** `vcs.follow-renames` (bool, Default `false`), unabhängig von
   `vcs.ignore-link-targets`; aus ⇒ byte-identisch
   ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
5. **Gleichgesetzt wird nur das aufgelöste Ziel**, und nur ein **Datei**-Ziel. Anker, Linktext, Titel,
   übrige Zeile und Zahl der Links bleiben roh; was ein Link ist, beantwortet weiter die Erkennung des
   Moduls `links`. Ein Link auf ein Verzeichnis bleibt wie heute; ein Link auf eine Datei darin folgt deren
   eigener Umbenennung. Die Fassung „ein Ziel, das in BASE auf A auflöst, gilt als Ziel auf B“ wird durch
   Entscheidung 9 präzisiert.
6. **Grenzen und fail-closed:** das Modul liest dafür Bäume jedes Commits der Range, über die Scan-Menge
   und `vcs.paths` hinaus; eine Umbenennung irgendwo im Repo zählt. Die Zusage gilt dort für
   Inhalts-Gleichheit, nicht für „derselbe Gegenstand“. Ein unlesbarer Commit oder Baum ist ein Fehler
   (Exit 2). Kein neuer Grund-Code. Die Umbenennung der immutablen Datei selbst bleibt Drift.
7. **Eigenes Profil:** `.d-check.yml` schaltet `follow-renames: true` im `vcs`-Block ein. Das ist eine
   Lockerung der Prüfregel des eigenen ADR-Gates ([`AGENTS.md`](../../../AGENTS.md) §3.6); diese ADR trägt
   sie. Bedingung siehe Entscheidung 14.
8. **Status:** diese ADR bleibt `Proposed`, bis der Mechanismus steht und reviewt ist.
9. **Die Gleichsetzung ist seiten-gleich.** Sie ist **eine** Funktion des aufgelösten Pfads, angewandt auf
   BASE- und HEAD-Fassung gleich, aufgelöst gegen die Vereinigung beider Pfad-Bäume wie heute. Grundlage
   ist die zusammengesetzte Abbildung Quelle → Ende (Entscheidung 11 und 13), nicht die Einzelschritte.
   Quelle ist nur ein Pfad des BASE-Stands; ein Zwischenpfad einer Kette ist weder Quelle noch Ende. Eine
   Umbenennung Quelle → Ende nimmt nur teil, wenn ihr Ende in HEAD existiert, **in BASE nicht** und Ende
   keiner zweiten Quelle ist — ein
   Pfad, der in BASE schon lag, hat dort eine eigene Geschichte (umbenannt, gelöscht und neu belegt,
   getauscht) und wird nicht zum Ende einer anderen Datei. Teil haben dann ihre Quelle (sie gilt als das
   Ende) und ihr Ende. Daraus folgt die Invariante: **gleicher Rohtext ergibt gleichen Core** — ein
   unveränderter Link auf eine umbenannte Datei bleibt gleich. Die Zusage in Befunden: ohne
   `ignore-link-targets` ergibt der Schalter keinen **Befund**, den es ohne ihn nicht gibt; mit beiden
   Schaltern entfällt ein Befund, den `ignore-link-targets` allein meldet, nur dort, wo beide Ziele
   teilnehmen und dasselbe Ende haben (Entscheidung 10). Rot ohne Befund kann ein Lauf mit Schalter
   trotzdem werden: ein unlesbares Objekt, das nur der Schalter liest, ist Exit 2 (Entscheidung 6 und 13).
10. **Zusammenspiel der Schalter.** Ein teilnehmendes Ziel wird zu einer **eigenen** Marke (die Marke
    zweimal), dem repo-relativen Pfad des Endes und dem Anker — mit `follow-renames` allein wie mit beiden
    Schaltern. Jedes andere auflösende Ziel bleibt mit `follow-renames` allein roh und wird mit beiden
    Schaltern wie mit `ignore-link-targets` allein zu Marke, Dateiname und Anker. Weil die beiden Formen
    verschiedene Marken tragen, gleicht ein teilnehmendes Ziel nie einem nicht teilnehmenden, auch bei
    gleichem Dateinamen; zwei teilnehmende gleichen sich nur bei gleichem Ende. Nur `ignore-link-targets`:
    unverändert.
11. **Die Kette folgt der Abstammung.** Eine Umbenennung Y → Z in einem Commit setzt X → Y aus einem
    anderen Commit nur fort, wenn dieser ein Vorfahr des ersten ist. Je Quelle zählen alle
    Abstammungslinien der Range. Eine Linie, auf der der Pfad nicht liegt oder in jedem Schritt bestehen
    bleibt (auch mit geändertem Inhalt), ist unberührt und trägt nichts bei. Eine Linie, die ihn umbenennt,
    trägt ihr Ende bei. Verschwindet der Pfad, auf dem eine Linie steht, in einem Schritt, ohne umbenannt
    zu werden, endet die Linie **ohne Ende** („gelöscht“) — auch wenn ein späterer Schritt den Pfad neu
    anlegt, und auch auf der Seite eines Merges, die ihn gar nicht umbenennt. Löschen und Neuanlegen im
    **selben** Schritt ist kein Verschwinden, sondern eine Inhaltsänderung an Ort und Stelle. Die Quelle
    hat genau dann eine Umbenennung, wenn mindestens eine Linie ein Ende beiträgt, alle beitragenden
    Linien dasselbe Ende tragen und keine Linie ohne Ende endet (Entscheidung 3). Das Ergebnis hängt nicht
    von der Reihenfolge ab, in der Commits oder Einträge aufgezählt werden.
12. **Basis kein Vorfahr von HEAD:** die Range trägt dann **keine** Umbenennung, das Modul urteilt wie ohne
    Schalter. Kein Exit 2, keine Kette ab der Merge-Basis. Die Anforderung nennt das als Grenze.
13. **Vertrag der neuen Port-Methode.** Sie wird nur bei aktivem Schalter gerufen. Sie liefert die
    Abbildung Quelle → Ende nach 1, 2, 3, 9, 11 und 12; Kette, Mehrdeutigkeit und die Bedingung an das
    Ende entscheidet der Adapter, das Gleichsetzen der Kern. Sie liest jeden Commit der Menge, den Baum
    jedes Schritts und seines Eltern-Commits — auch eines Eltern-Commits vor BASE, etwa dort, wo eine
    Merge-Seite vor BASE abzweigt — und jeden Teilbaum, der sich zwischen beiden unterscheidet; ein
    Teilbaum mit gleicher Kennung auf beiden Seiten muss nicht gelesen werden (gleicher Inhalt, keine
    Umbenennung darin). Jedes dieser Objekte, das sie nicht lesen kann, ist ein Fehler (Exit 2) — auch wenn
    BASE und HEAD lesbar sind; das ist die fail-closed-Politik von `AllPaths`, auf die Schritte übertragen.
    Teil haben drei Arten von Einträgen — Datei, ausführbare Datei, Symlink —, in Bäumen wie im Index; ein
    Verzeichnis und ein Gitlink nehmen nicht teil (im Index nimmt `AllPaths` jeden Eintrag, die Abbildung
    nicht). Gleich sind zwei Einträge mit gleicher Blob-Kennung und gleicher Art (Datei oder Symlink), das
    Ausführbar-Bit zählt nicht. Die Umbenennungs-Erkennung der Bibliothek wird **nicht** benutzt. Im
    `--staged`-Modus ist der Schritt der HEAD-Baum gegen den Index, und „BASE“ in Entscheidung 9 ist HEAD;
    ohne HEAD ist die Abbildung leer.
14. **Reihenfolge im eigenen Profil.** Der Schalter in `.d-check.yml` landet im selben Commit wie der
    Übergang dieser ADR nach `Accepted`, nie davor. [`AGENTS.md`](../../../AGENTS.md) §3.5 und
    `harness/README.md` §Traceability rules, die heute
    [ADR-0103](0103-adr-gate-laesst-pfad-nachzug-durch.md) als Träger des Pfad-Nachzugs nennen, nennen
    diese ADR daneben.

Diese ADR erweitert [ADR-0103](0103-adr-gate-laesst-pfad-nachzug-durch.md), sie löst sie nicht ab: deren
Entscheidungen 1–4 und die Wirkung von `ignore-link-targets` bleiben unverändert gültig.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun | keine neue Eingabe, kein neuer Code | das Gate wird bei einer Operation rot, die die Regel erlaubt, die es durchsetzt; Adopter übergehen einen roten Lauf wissentlich — die Gewöhnung, gegen die das Gate antritt |
| B — die Bedeutung von `ignore-link-targets` erweitern (Vorschlag 1 des Issues) | kein neuer Schlüssel | lockert das Gate jedes Repos, das den Schlüssel schon setzt, ohne dass es darum gebeten hat ([`AGENTS.md`](../../../AGENTS.md) §3.6) |
| C — Ähnlichkeits-Schwelle (git-Vorbelegung 50 %, Bibliothek 60) | erkennt auch „umbenennen und im selben Commit ändern“ | ein Parameter, den niemand entschieden hat; lässt im selben Schritt eine Umschrift durch; zwei Werkzeuge, zwei Zahlen |
| D — nur die Endpunkte BASE und HEAD vergleichen | ein Baumpaar, kein Commit-Walk | nach [`AGENTS.md`](../../../AGENTS.md) §3.3 folgt die Inhaltsänderung der umbenannten Datei als eigener Commit, oft im selben Push; an den Endpunkten trügen A und C verschiedene Blobs, und der Nachzug bliebe rot, obwohl jeder Schritt sauber ist |
| E — Verweise über eine Kennung statt einer Adresse (Empfehlung des Baseline-Regelwerks) | kein Nachzug nötig | trägt für künftige ADRs, nicht für den Bestand, dessen Verweise Adressen sind |
| F — Gleichsetzung nur auf der BASE-Seite (Wortlaut von Entscheidung 5) | wörtlich der bestätigte Satz | jede unveränderte ADR mit einem Link auf eine umbenannte Datei wird rot, weil das Modul alle geschützten Dateien beider Stände vergleicht |
| G — Basis kein Vorfahr ⇒ Exit 2 (wie `LineChanges`) | laut, eine Regel | der Rückfall der CI-Range nach einem Force-Push ([ADR-0109](0109-ci-push-range-faellt-bei-force-push-zurueck.md)) kann eine solche Basis liefern; mit Entscheidung 7 würde der Push-Lauf wieder rot, ohne dass eine Eigenschaft verletzt ist |
| H — Basis kein Vorfahr ⇒ Kette ab der Merge-Basis | der Schalter wirkt auch dann | falsches Grün: hat die BASE-Seite einen Pfad A neu belegt (C → A) und die HEAD-Seite das alte A nach B umbenannt, gilt ein Link auf das neue A als nachgezogen auf B |
| I — Umbenennungs-Erkennung der Bibliothek mit `OnlyExactRenames` | weniger eigener Code | wählt bei gleichen Blobs nach Namensähnlichkeit statt Mehrdeutigkeit zu melden, verlangt gleichen Modus; der Diff-Walker darunter las ein unlesbares Objekt früher still als Ende |
| J — `--staged` liest zusätzlich den Commit HEAD | ein Move im vorigen Commit wäre gesehen | willkürlich: ein Move zwei Commits zurück bliebe ungesehen; die CI-Range sieht beide Commits, und der Weg „Verweise im Move-Commit“ steht im Handbuch |
| L — beide Schalter: teilnehmendes Ziel auf den **Dateinamen** des Endes normiert (am 2026-11-23 bestätigte Fassung von Entscheidung 10) | eine Normalform für beide Ziel-Arten | ein Umbiegen von `d1/a.md` auf `d2/c.md`, das die Range nach `d3/a.md` umbenennt, ist grün, mit `ignore-link-targets` allein Drift; eine Umbenennung, die nur das Verzeichnis wechselt, öffnet das Umbiegen auf eine andere gleichnamige Datei |
| M — ein Pfad mit Doppelrolle (Ende einer Umbenennung, in BASE schon vorhanden) nimmt als Ende teil | ein Nachzug auf eine Datei, die einen freigemachten Namen übernimmt, folgt | ein Umbiegen vom alten Inhalt des Pfads auf die umbenannte Datei wird grün; die Regel braucht dann eigene Sätze für Tausch und Neubelegung |
| **K — gewählt: eigener Schalter, Blob-Gleichheit je Commit, Kette entlang der Abstammung, Mehrdeutigkeit und fremde Basis ohne Umbenennung, seiten-gleiche Gleichsetzung, eigene Marke für teilnehmende Ziele, Ende nie in BASE** | entscheidbar ohne Parameter, deterministisch, opt-in; ein reiner Move-Commit nach §3.3 ist genau die Form, die erkannt wird; ohne Schalter byte-identisch; mit Schalter entfällt ein Befund nur, wo beide Ziele dasselbe Ende haben (der fail-closed-Fall der Entscheidung 6 ist kein Befund, aber Exit 2) | liest jeden Commit der Range; Grenzen unter Konsequenzen |

## Konsequenzen

- Positiv: eine Layout-Migration mit neuen Dateinamen zieht die Links in `Accepted`-ADRs nach, ohne dass
  das ADR-Gate rot wird und ohne Folge-ADR. Die Identität „dieselbe Datei“ kommt aus der Historie, nicht
  aus einem Urteil.
- Negativ — umbenennen, dann umschreiben: ein reiner Move-Commit und eine spätere Inhaltsänderung der
  umbenannten Datei gehen durch (Entscheidung 2 will das). Die Zusage ist Pfad-Identität entlang der
  Historie, nicht Inhalts-Identität des Ziels; auch ein unveränderter Link auf eine umgeschriebene Datei
  ist heute keine Drift. Wer eine Umbenennung vortäuscht, um einen Link umzubiegen, sieht das Gate nicht;
  das trägt der Review des Commits.
- Negativ — umbenennen und im selben Commit ändern bleibt Drift; ebenso eine Umbenennung, bei der ein
  Werkzeug den Inhalt ersetzt (Archiv-Stub im selben Akt).
- Negativ — ein Ende, das in BASE schon lag, nimmt nicht teil (Entscheidung 9): ein Tausch zweier Namen
  (A ↔ B über einen Zwischennamen) bleibt Drift, ebenso ein Nachzug auf einen Pfad, den die Range erst
  freimacht und dann mit einer anderen Datei belegt (A gelöscht, danach C → A). Fehlalarm in die sichere
  Richtung.
- Negativ — Dateien gleichen Inhalts (leere Dateien, gleiche Vorlagen) in einem umbenannten Verzeichnis
  sind mehrdeutig, ihre Links folgen nicht.
- Negativ — beide Schalter an ist strenger als `ignore-link-targets` allein, und nur dort, wo ein Ziel
  teilnimmt: ein teilnehmendes Ziel gleicht nur einem teilnehmenden mit demselben Ende (Entscheidung 10).
  Ein Nachzug von einer umbenannten Datei auf eine **andere** gleichnamige Datei ist Drift, auch wenn die
  Umbenennung nur das Verzeichnis wechselt; ein Umbiegen auf eine Datei, die in der Range den alten Namen
  bekommt, ebenso. Die Grenze „gleicher Name, andere Datei“ bleibt nur zwischen zwei Zielen, von denen
  keines teilnimmt — dort ist das Ergebnis dasselbe wie mit `ignore-link-targets` allein.
- Negativ — `--staged` sieht nur den Schritt HEAD → Index. Liegt der Move in einem früheren Commit, bleibt
  der gestagte Nachzug Drift; der Weg ist, die Verweise im Move-Commit mitzunehmen. Die CI-Range sieht
  beide Commits.
- Negativ — Basis kein Vorfahr: der Schalter wirkt nicht, der Nachzug bleibt Drift wie ohne ihn.
- Negativ — trägt die Range einen Merge, beschreiben die Commits seiner Seiten den Weg ab der
  Merge-Basis. Löscht die andere Seite einen Pfad, den diese umbenennt, in einem Commit und legt ihn in
  einem späteren neu an, ist das Mehrdeutigkeit (Entscheidung 11). Überschreibt sie ihn an Ort und Stelle
  — auch durch Löschen und Neuanlegen im selben Commit —, ist das eine Inhaltsänderung, und der Link folgt
  der Umbenennung: dieselbe Grenze wie „umbenennen, dann umschreiben“, die Zusage ist Pfad-Identität,
  nicht Inhalts-Identität.
- Negativ — Neubelegung (Entscheidung 9 und 11): wird ein Pfad einer Kette in einem Commit gelöscht und in
  einem späteren neu belegt, folgt ein Link ihm nicht. Das Ziel nimmt dann nicht teil; mit beiden
  Schaltern wird es wie mit `ignore-link-targets` allein normiert und gleicht keinem teilnehmenden Ziel.
- Negativ — Laufzeit und Lesbarkeit: mit Schalter liest das Modul jeden Commit der Range, die Bäume der Schritte und ihrer Eltern-Commits (auch vor BASE) und die geänderten Teilbäume;
  ohne Schalter liest es nichts davon.
- Folgepflicht: die Anforderung nennt die Grenzen oben; das Kriterium „Boundary (Pfad-Nachzug)“ wird auf
  „ohne `vcs.follow-renames`“ eingeschränkt statt ersetzt, und der Out-of-Scope-Punkt zur
  Umbenennungs-Erkennung über Inhalts-Ähnlichkeit wird auf die immutable Datei selbst und auf Ähnlichkeit
  statt Gleichheit geschärft. Jede Randbedingung hat einen Test, der ohne die Regel aus dem richtigen Grund
  rot läuft.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Tests des git-Adapters gegen echte git-Fixtures | Blob-Gleichheit, Kette über drei Commits, Kette nur entlang der Abstammung, Mehrdeutigkeit (Kopien, zwei Enden, zwei Linien), Merge als HEAD mit einer umbenennenden Seite, Merge mit einer Seite, die den Pfad in einem Commit löscht und in einem späteren neu anlegt ⇒ keine Umbenennung, Löschen und Neuanlegen im selben Commit ⇒ Umbenennung bleibt, Ende in BASE vorhanden (A gelöscht, dann C → A) ⇒ keine Umbenennung, unlesbarer Baum eines Eltern-Commits vor BASE ⇒ Fehler, Basis kein Vorfahr ⇒ leere Abbildung, `--staged`, unlesbares Objekt eines Zwischen-Commits ⇒ Fehler | `make test` |
| Tests der Regel | Nachzug mit Schalter kein Befund, ohne Befund; unveränderter Link auf umbenannte Datei mit erlaubtem Anhang kein Befund; beide Schalter samt gleichnamiger anderer Datei ⇒ Befund; Ende in einem Commit gelöscht und in einem späteren als andere Datei neu angelegt ⇒ Befund; beide Schalter, Umbiegen auf eine Datei, die die Range in den alten Namen umbenennt ⇒ Befund; beide Schalter, Umbenennung nur des Verzeichnisses und Umbiegen auf eine andere gleichnamige Datei ⇒ Befund; Zwischenpfad einer Kette neu belegt und als Ziel nachgezogen ⇒ Befund; unlesbarer Zwischen-Commit ohne Schalter Exit 0, mit ihm Exit 2; bestehende Tests von `ignore-link-targets` unverändert grün; derselbe Lauf zehnmal byte-gleich | `make test` |
| Probe am gebauten Image in einem Wegwerf-Klon | Umbenennung einer verlinkten Spec-Datei samt Nachzug: ohne Schalter je ADR ein Befund, mit ihm keiner; gelöscht und neu angelegt ⇒ Befund | `make adr-check` |
| Vorher-/Nachher-Vergleich ohne Schalter | Ausgabe byte-gleich zur letzten Release-Version | `make blackbox-probe` |

## Re-Evaluierungs-Trigger

- Das Baseline-Regelwerk ändert die Ausnahme für den Referenz-/Pfad-Nachzug.
- Ein Review findet einen Nachzug über eine Umbenennung, der die Aussage einer ADR verschoben hat
  (umbenennen, dann umschreiben als Umweg).
- Eine der benannten Grenzen (Basis kein Vorfahr, `--staged`, Merge, Mehrdeutigkeit) meldet bei einem
  Adopter wiederholt einen Nachzug als Drift, der keiner ist.
- `make adr-check` über eine Range von 100 Commits braucht mit Schalter mehr als doppelt so lange wie ohne.
- Die Bibliothek bietet eine Umbenennungs-Erkennung, die Mehrdeutigkeit meldet statt nach Namen zu wählen
  und auf einem unlesbaren Objekt abbricht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-23 | Proposed: Entscheidung 1–8 wie vom Auftraggeber bestätigt; 9–14 als Empfehlung aus dem Architect-Verdikt (seiten-gleiche Gleichsetzung, Zusammenspiel der Schalter, Kette entlang der Abstammung, fremde Basis ohne Umbenennung, Port-Vertrag, Reihenfolge im eigenen Profil) | slice-291 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, fortgeschrieben: der Auftraggeber hat Entscheidung 9–14 bestätigt, ebenso Entscheidung 5 in der seiten-gleichen Fassung und die Schärfung des Out-of-Scope-Satzes der Anforderung | slice-291 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, an die Spezifikation angeglichen: ein Schritt ist ein Commit mit genau einem Eltern-Commit, ein Wurzel-Commit trägt keine Umbenennung bei (2); „dasselbe Blob“ heißt gleicher Eintrag aus Blob-Kennung und Art (1, 3); Teil haben in Bäumen wie im Index nur Datei, ausführbare Datei und Symlink, ein Gitlink nicht (13); „zugleich Quelle und Ende“ bezieht sich auf die zusammengesetzte Abbildung, nicht auf die Einzelschritte (9) | slice-291 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, nach Review R1 zum Anforderungs-Slice geschärft: Quellen der Abbildung sind nur Pfade des BASE-Stands, ein Zwischenpfad ist weder Quelle noch Ende (9); verschwindet der Pfad einer Linie ohne Umbenennung, endet sie mit „gelöscht“, auch bei späterer Neubelegung (11); die Zusage „nie röter“ heißt „kein Befund mehr“, der fail-closed-Fall ist benannt (9, Alternative K); `Schärft:` nennt Schritt 2 statt 1; die Einleitung der Entscheidung sagt, was nach der Bestätigung angeglichen wurde | slice-291 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, Entscheidung 11 präzisiert: auch eine Merge-Seite, die einen Pfad nicht umbenennt, ihn aber löscht und neu anlegt, ist eine Linie zu „gelöscht“; neben der Umbenennung der anderen Seite ist das Mehrdeutigkeit. Die Merge-Grenze unter Konsequenzen gilt nur noch für das Überschreiben an Ort und Stelle | slice-291 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, nach Review R2 zum Anforderungs-Slice: Entscheidung 10 normiert ein teilnehmendes Ziel auch mit beiden Schaltern auf eine eigene Marke und den vollen Endpfad, damit gleicht es nie einem nicht teilnehmenden (schließt das Umbiegen auf eine Datei, die den alten Namen bekommt, und die Bedingung „Umbenennung wechselt den Dateinamen“ entfällt) — **in der Sache geändert, nur in Richtung Drift, Bestätigung durch den Auftraggeber offen**; Entscheidung 9: ein Ende nimmt nur teil, wenn es in BASE nicht lag und Ende genau einer Quelle ist (ersetzt „zugleich Quelle und Ende“); Entscheidung 11: „gelöscht“ heißt ohne Ende, Löschen und Neuanlegen im selben Schritt ist Inhaltsänderung, eine Linie ohne den Pfad ist unberührt; Entscheidung 13: gelesen werden auch Bäume von Eltern-Commits vor BASE; Alternativen L und M | slice-291 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, fortgeschrieben: der Auftraggeber hat Entscheidung 9 und 10 in der Fassung nach Review R2 bestätigt | slice-291 <!-- d-check:status-provenance --> |
