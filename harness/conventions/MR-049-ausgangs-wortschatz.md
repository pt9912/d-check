# MR-049 — Der Ausgang eines Risikos trägt eines der drei Wörter (schärft MR-006)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon benennt die urteilsfreie Hälfte
  der Drei-Ausgänge-Regel ausdrücklich — *dass* zu jedem Risiko ein Ausgang
  dasteht **und welcher der drei** es ist — und schließt mit: *„Welches Werkzeug
  die urteilsfreie Hälfte prüft, ist Repo-Entscheidung; dass sie eine hat, ist
  es nicht"*
  ([`modul-05-planning-harness.md`](../../.harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md)).
  Diese Adaption ist die Werkzeug-Wahl, nicht eine Abweichung.
- **Datum:** 2026-08-27
- **Geltungsbereich:** [`.d-check.closure.yml`](../../.d-check.closure.yml),
  Abschnitt §5 der Slice-Dateien unter `docs/plan/planning/done/`.
- **Adaption:** Der Ausgangs-Wortschatz ist **geschlossen**: `eingetreten` ·
  `entfallen` · `weiter offen`. Eine `structure`-Regel im Closure-Profil hält
  die ersten beiden; sie läuft in `make verify-closure-notes`.

  **Sie ergänzt die vorhandene Regel, sie ersetzt sie nicht.** Die
  Platzhalter-Regel fängt den **vergessenen** Ausgang, diese den **erfundenen** —
  die Gestalt aus
  [`BEO-015`](../../docs/plan/planning/observations.md).

  **`weiter offen` steht mit Absicht nicht in der Alternation:** dieser Ausgang
  trägt seinen Beleg als `BEO`-Kennung, und die Paarung Register ↔ Beleg ist
  eine eigene Frage mit eigener Mechanik.
- **Grenzen — gemessen, nicht geschätzt.**

  **Je Abschnitt, nicht je Risiko.** Die `structure`-Bedingungen wirken auf den
  **ganzen bereinigten Abschnitts-Text**, und RE2 kennt keinen
  Negativ-Lookahead: `(?!` wird vom Produkt mit *„invalid or unsupported Perl
  syntax"* abgewiesen. Ein §5 mit zwei Risiken, von denen eines kanonisch und
  eines als Freitext ausgeht, läuft deshalb **grün** durch. Gedeckt ist der
  Abschnitt, in dem **kein** Ausgang kanonisch ist. Die Korrelation Risiko ↔
  Ausgang bleibt damit ungedeckt; wer sie will, braucht ein Produkt-Delta.

  **Der Altbestand ist ausgenommen, und die Ausnahme zeigt nach hinten.**
  Ohne sie meldet die Regel **121** Befunde: **107** Slices tragen gar keinen
  §5-Abschnitt (ältere Haus-Form), **14** der übrigen 43 tragen einen
  Freitext-Ausgang. Ein Retrofit über 121 Dateien hat niemand beschlossen, und
  der Kurs hat die Ausweitung ausdrücklich als Repo-Entscheidung und **kein**
  Konformitätsthema bezeichnet.

  Die Ausnahme nennt deshalb den **abgeschlossenen** Altbestand (`slice-0*`,
  `slice-1[0-3]*`) und nicht die offene Zukunft. Ein Zahlen-Glob auf die neuen
  Nummern — `slice-1[4-9][0-9]` — wäre die naheliegende Form und hörte bei
  `slice-200` **still** auf zu greifen: dieselbe Klasse von lautlosem Ausfall,
  die dieses Repo zuletzt zweimal getroffen hat.
- **Begründung:** Der Wortschatz war offen, und das ist gemessen: über den
  Bestand tragen die Ausgänge rund zwanzig verschiedene Anfangswörter, darunter
  `behoben`, `erledigt`, `gemessen`, `gehalten`, `aufgelöst`, `benannt` — und
  blanken Freitext. Jedes einzelne liest sich vernünftig; zusammen ist die
  geschlossene Menge des Kanons keine mehr, und ein vierter Ausgang fällt nicht
  auf, weil er wie der fünfte klingt.
- **Auflösungs-Trigger:** das Produkt kann eine Bedingung **je Listen-Eintrag**
  ausdrücken. Dann ist die Je-Risiko-Korrelation prüfbar, und ein
  Nachfolge-Eintrag löst diesen ab.
