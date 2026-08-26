# Welle welle-85-baseline-v5120-migration: Vier Kurs-Wellen, und alle vier sind unsere eigene Bitte

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-85-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug. Ob die Welle mit einem Release
schließt, entscheidet das Delta-Audit — eine Harness-Migration ist
konsumentensichtbar nur dort, wo sie das Produkt berührt.

**Verantwortlich:** pt9912. **Datum:** 2026-08-26.

---

## 1. Welle-Ziel

Der Baseline-Pin steht auf `v5.11.0` (Kurs-Welle **94**); der Kurs steht auf
`v5.12.0` (Kurs-Welle **98**). **Vier Wellen, ein Minor** — und anders als bei
jeder Hebung zuvor wissen wir vorher, was drinsteht: **alle vier sind die
Antwort auf den Konsumenten-CR dieses Repos vom 2026-08-25**, eine Welle je
Punkt. Die Antwort liegt im Repo
([`2026-08-26-antwort-regelwerk-v5110.md`](../cr/2026-08-26-antwort-regelwerk-v5110.md),
[`MR-036`](../../../harness/conventions.md#mr-036)).

**Das Delta ist gemessen, nicht geschätzt.** 52 Bundle-Dateien, **28**
unterscheiden sich; **22 davon ändern ausschließlich den Versions-Stempel**
(zwei Zeilen, maschinell geprüft: kein 2-Zeilen-Delta trägt etwas anderes), und
`regelwerk/README.md` trägt zusätzlich nur seine `**Stand:**`-Zeile. **Fünf
Dateien tragen echten Regel-Inhalt:**

| Datei | Umfang | CR-Punkt |
|---|---|---|
| `regelwerk/modul-05-planning-harness.md` | 11 Zeilen neu | 1 — die urteilsfreie Hälfte der Drei-Ausgänge-Regel ist benannt |
| `regelwerk/modul-13-quality-gates.md` | 8 Zeilen neu | 3 — das Rot muss von *dieser* Regel kommen, Ursache gelesen |
| `regelwerk/modul-03-spec.md` | 5 für 2 | 2 — *Modul-Pfad* heißt **Code**-Modul-Pfad |
| `templates/AGENTS.template.md` | 3 für 2 | 2 — derselbe Satz im Vorlagen-Spiegel |
| `regelwerk/grundlagen-source-precedence.md` | 6 Zeilen neu | 4 — die Reichweitenfrage als Frage, nicht als Katalog |

**Zwei Folgen stehen schon fest**, beide beim Messen gefunden und nicht geraten:

1. **Punkt 2 trifft [`MR-033`](../../../harness/conventions.md#mr-033).** Der
   Kanon sagt jetzt *„darf Pfade zu **Code-Modulen** referenzieren"* — und
   ergänzt: *„Die Erlaubnis ist keine Pflicht: Eine Sicht, die ihre Komponenten
   nur über Rollen und `ARC-*` führt, ist ebenso konform."* Unsere Lesart war
   richtig; ob unser **Verbot** danach noch eine Verschärfung ist oder nur eine
   nicht ausgeübte Erlaubnis, ist die Frage des Freshness-Audits, nicht dieser
   Eröffnung.
2. **Punkt 1 benennt die urteilsfreie Hälfte weiter, als wir sie prüfen.** Der
   neue Text sagt: urteilsfrei ist, *dass* zu jedem Risiko ein Ausgang dasteht
   **und welcher der drei** es ist. Unsere Regel aus
   [slice-143](done/slice-143-structure-abschnitts-skopus.md) prüft davon den
   häufigsten Auslöser — den stehengebliebenen Platzhalter. Der Kanon schließt
   mit *„Welches Werkzeug die urteilsfreie Hälfte prüft, ist Repo-Entscheidung;
   dass sie eine hat, ist es nicht."* Wir haben eine; ob sie die ganze Hälfte
   trägt, ist zu entscheiden.

## 2. Trigger (Welle startet)

`v5.12.0` ist upstream verfügbar — gemessen durch `make baseline-freshness`
(Exit 3, neuer Release), bei zugleich **unverändertem** gepinnten Stand (Bytes
gleich dem vendorten `SHA256SUMS`). Es ist eine reine Fortschreibung, kein
Nachholen stiller Drift.

## 3. Closure-Trigger (Welle schließt)

Die Welle schließt, wenn **alle drei** Aussagen belegt sind — jede beobachtet
mehr als die DoD eines einzelnen Slice, und genau darum ist dies eine Welle und
keine wellenlose Arbeit (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit
eine Welle braucht):

1. **Der Pin ist gehoben, und zwar in allen Spiegel-Klassen.**
   [`BEO-008`](observations.md) führt drei: Pfad-Verweise (gate-gedeckt),
   Release-/Tree-URLs und Prosa-/Ellipsen-Pins (beide ungedeckt). Je Klasse
   eine Zählung vor und nach der Hebung.
2. **Jede aktive Adaption hat eine Antwort auf die Freshness-Frage des Kanons**
   (`modul-02-harness-bootstrap.md` §Freshness-Audit): *Regelt die neue Fassung
   das, wofür diese Adaption angelegt wurde?* — je Eintrag ein Ergebnis, auch
   „unberührt", und keiner ungefragt.
3. **Jede der vier Kurs-Änderungen hat eine Antwort für dieses Repo** —
   Handlung mit Slice-ID oder belegt folgenlos. Die Antwort des Kurses ist
   dabei **Erwartung, nicht Beleg**: geprüft wird gegen den vendorten Text.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-148](done/slice-148-baseline-v5120-vendoring.md) | **Etappe A:** Bundle `v5.12.0` vendored, Pin-Hebung als `MR-`Eintrag, alle drei Spiegel-Klassen gehoben, Alt-Baum entfernt |
| [slice-149](open/slice-149-baseline-v5120-delta-audit.md) | **Etappe B:** Delta-Audit über die vier Kurs-Wellen **und** der Freshness-Audit über alle aktiven Adaptionen; **schneidet Etappe C** |

**Etappe C wird vom Audit geschnitten, nicht hier geraten** — dieselbe Form wie
in den zwei Migrationen davor. Was dabei entsteht, kommt mit Kennung in diese
Tabelle und in eine Drift-Log-Zeile.

## 5. Abhängigkeiten

**Die Reihenfolge ist nicht verhandelbar:** erst pinnen, dann dem Kanon folgen.
Andernfalls führte `AGENTS.md` den Stand von Kurs-Welle 98, während der
vendorte Baum noch 94 sagt — und ein Audit gegen zwei Stände ist kein Audit.

Etappe B setzt Etappe A voraus. Innerhalb von B ist der Freshness-Audit
unabhängig vom Delta-Audit und darf vorgezogen werden.

## 6. Out-of-Scope für diese Welle

- **Kein Release.** Ob die Migration eines braucht, entscheidet das Delta-Audit
  — nur ein Produkt-Delta ist konsumentensichtbar.
- **Keine Ausweitung der Closure-Regel auf die volle urteilsfreie Hälfte.** Der
  Kurs hat das ausdrücklich als unsere Entscheidung und **kein
  Konformitätsthema** bezeichnet. Sie gehört als eigener Slice geschnitten,
  falls das Audit sie befürwortet — nicht als stiller Anhang der Hebung.
- **Keine Nacharbeit an eingefrorenen Zitaten.** Wo `done/`-Dokumente den alten
  Wortlaut zitieren, bleibt er stehen; nur **lebende** Stellen werden gehoben.
