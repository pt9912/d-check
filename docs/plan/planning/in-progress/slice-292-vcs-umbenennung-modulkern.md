# Slice slice-292: Kern und git-Adapter — `vcs` liest die Umbenennungen der Range und setzt Link-Ziele über sie gleich

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-96-vcs-folgt-umbenennung`](../welle-96-vcs-folgt-umbenennung.md).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
in der Fassung aus [`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md)
und [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md), [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[`ADR-0024`](../../adr/0024-vcs-immutable-gate.md),
[`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md).

**Berührte Spec-Stellen:** `spezifikation.md`
§[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
Schritt 2, nur zwei Sätze der Enumeration: der Nebensatz zur Index-Enumeration
([Verifikation](../../../reviews/2026-11-24-slice-291-vcs-umbenennung-verify.md) von
[`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md), V-1) und — **nur
wenn** der Testfall „immutable Datei durch Gitlink ersetzt“ in §3 `core-drift-vcs` in beiden Modi bestätigt — der Satz vom „benannten
Preis“ eines Gitlink-Ersatzes
([Architect-Verdikt](../../../reviews/2026-11-24-slice-292-architect-verdikt.md) F-5).
Verfahren und Schema schrieb
[`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md); kein neuer
Grund-Code, §4 der Spezifikation bleibt. Der Verweis zeigt **aufwärts**: Die Spec
nennt diesen Slice nie.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-11-23.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Mit `vcs.follow-renames` liest das Modul die Umbenennungen der Range über den
git-Port (fail-closed), setzt ein Link-Ziel über sie gleich, wie
[`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md) es festschreibt, und
ein reiner Nachzug auf die umbenannte Datei ist keine Drift mehr — ohne den Schlüssel
byte-identisch (mit `--range`; die eine gewollte Abweichung mit `--staged` nennt §3 Zeile
`git.go`); jede Randbedingung der Anforderung hat einen Test, der ohne die Regel aus
dem richtigen Grund rot läuft, und die ADR ist `Accepted`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Anforderung, Verfahren, ADR-Wortlaut ändern** — gehören
  [`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md), **außer** den
  zwei Sätzen in Schritt 2 der Spezifikation, die mit dem Code fallen (Kopf
  `Berührte Spec-Stellen`, §3 Zeile `spec/spezifikation.md`): dem Nebensatz zur
  Index-Enumeration (V-1) und, unter der dort genannten Bedingung, dem Satz zum
  Gitlink-Ersatz (Architect-Verdikt F-5). Zeigt der Code,
  dass ein Entscheid nicht trägt, geht das als Plan-Änderung an den Architect, nicht als
  stille Abweichung im Code; die ADR ändert nur der Architect.
- **`metrics`, `tracked`, `commits`** — lesen git, folgen aber keiner Umbenennung;
  `LineChanges` bleibt ohne Rename-Erkennung (sein Vertrag sagt es).
- **Die Grenze „gleichnamige andere Datei“** von `ignore-link-targets`, **für Ziele, die an
  keiner Umbenennung teilnehmen** — Welle-Out-of-Scope. Zwischen einem teilnehmenden und
  einem nicht teilnehmenden Ziel gilt sie mit beiden Schaltern nicht
  ([`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) Entscheidung 10; Kriterium
  „Boundary (Negativ-Kriterien mit beiden Schlüsseln)“ von
  [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)).
- **CHANGELOG, Handbuch-Kopf, Release** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17). Was ein Spiegel-Test schon im
  Feature-Commit verlangt, kommt mit und steht in §3.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] git-Port und -Adapter liefern die Umbenennungen der Range nach der Anforderung
      (Commit-Menge, Inhalts-Gleichheit, Kette, Mehrdeutigkeit, die Bedingung an das
      Ende — in HEAD, nicht in BASE, Ende keiner zweiten Quelle
      ([`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) Entscheidung 9) —, Basis kein
      Vorfahr ⇒ leere Abbildung (Entscheidung 12), `--staged`); die Index-Enumeration nimmt
      dieselben drei Arten wie der Baum (Entscheidung 13, §3 Zeile `git.go`); ein
      unlesbarer Commit oder Baum ist Exit 2; Tests gegen echte git-Fixtures.
- [ ] Regel und Konfigurations-Rand: `vcs.follow-renames` wird gelesen (falscher Typ ⇒
      Exit 2), das Gleichsetzen trägt jedes Akzeptanzkriterium der Anforderung, ohne
      Schlüssel ist die Ausgabe byte-identisch, außer der einen gewollten Abweichung mit
      `--staged` (§3 Zeile `git.go`); die Negativfälle aus dem Plan von
      [`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md) stehen als
      Testfälle, und für jeden Test ist ein Rot ohne die Regel gezeigt
      (Baseline-Regelwerk `modul-11-verification.md` §Bewusstes Brechen) — ausgenommen
      allein der zweite Lauf eines Negativ-Falls mit beiden Schlüsseln, wo §3 ihn als Netz
      ohne eigenen Rot-Beleg benennt.
- [ ] Spiegel und Profil: `--print-config`, Handbuch-Abschnitt des Moduls und
      [`harness/sensors/adr-check.md`](../../../../harness/sensors/adr-check.md) nennen den
      Schlüssel (vor dem Editieren aufgelistet,
      [`MR-025`](../../../../harness/conventions.md#mr-025)); `.d-check.yml` schaltet ihn
      ein **im selben Commit**, der
      [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) auf `Accepted` hebt und den
      Index nachzieht (Entscheid 14); [`AGENTS.md`](../../../../AGENTS.md) §3.5 und
      [`harness/README.md`](../../../../harness/README.md) §Traceability rules nennen
      [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) neben
      [`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/port/driven/vcs.go` | update | neue Methode am Port, gerufen nur bei aktivem Schalter: die Abbildung Quelle → eindeutiges Ende (Vertrag nach [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) Entscheidung 13: Commit-Menge, Inhalts-Gleichheit nach Blob-Kennung und Art ohne Ausführbar-Bit, Kette entlang der Abstammung, Mehrdeutigkeit, Basis kein Vorfahr ⇒ leer, fail-closed auf jedes nötige Objekt; Kette und Mehrdeutigkeit entscheidet der Adapter, das Gleichsetzen der Kern). **Form (Architect-Verdikt F-6):** eine Zuordnung Quelle → Ende, repo-relativ und `/`-getrennt, nur mit den **teilnehmenden** Umbenennungen, also schon nach den Bedingungen von Entscheidung 9 gefiltert; leer, nicht Fehler, bei „Basis kein Vorfahr“ und im gestagten Modus ohne HEAD. **Invariante, die der Kern prüfen darf:** keine Quelle ist zugleich Ende — sie folgt aus „Ende nicht in BASE“, weil jede Quelle in BASE liegt. Damit ist die Gleichsetzung im Kern ein reines Nachschlagen (Quelle ⇒ ihr Ende, Ende ⇒ es selbst), und die Iterationsreihenfolge der Zuordnung kann die Ausgabe nicht berühren ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)) |
| `internal/adapter/driven/git/git.go` | update | eigener Blob-Vergleich über die Bäume je Commit und Eltern-Commit bzw. HEAD gegen Index; Teilbäume gleicher Kennung bleiben ungelesen; die Rename-Erkennung von go-git wird **nicht** benutzt; kein Diff-Walker, der einen unlesbaren Unterbaum still überspringt. **Dazu (Plan-Änderung nach dem Architect-Befund S-1, 2026-11-23):** `AllPaths` nimmt im gestagten Modus heute jeden Index-Eintrag, auch einen Gitlink (`git.go` Z. 95–98), die Spezifikation (§[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2) sagt „jeder Datei-Pfad (regulär/ausführbar/symlink) … bzw. am staged Index“ — der Code folgt der Spezifikation (Rang 2): dieselben drei Arten wie im Baum, ein Filter für Enumeration und Umbenennungs-Abbildung. **Folge ohne den Schalter (Architect-Verdikt F-4, nach Lesart des Codes, nicht gefahren):** ersetzt im gestagten Modus ein Gitlink eine immutable Datei an ihrem Pfad, liegt der Pfad heute in beiden Mengen, das Lesen über `fileFromIndex` endet mit einem Fehler (Exit 2); nach der Angleichung fehlt er in der Index-Menge, und `vcsDeleted` meldet `core-drift-vcs` (Exit 1) — dieselbe Antwort wie im `--range`-Modus. Das ist **die eine gewollte Abweichung** von „ohne den Schalter nichts Neues“; sie gilt nur mit `--staged`, nur für einen Gitlink, der eine Datei unter `vcs.paths` ersetzt, deren BASE-Fassung `vcs.immutable-when` erfüllt. Jeder andere Gitlink unter `vcs.paths` war und bleibt frei. Das Kriterium „Boundary (ohne `vcs.follow-renames` byte-identisch)“ nennt nur `--range`; die Abweichung gleicht an Rang 2 an und lockert nichts (Exit 2 → Exit 1, der Befund wird sichtbar). Testfall „immutable Datei durch Gitlink ersetzt“ unten |
| `internal/hexagon/core/model/config.go` | update | `VCSConfig` um den Schalter |
| `internal/adapter/driven/configyaml/configyaml.go` | update | Schlüssel lesen und am Rand prüfen (Exit 2), **gelesen wie `ignore-link-targets`** (`rawVCS` in `configyaml.go`) — dieselbe Lesart für einen falschen Typ und für einen Schlüssel ohne Wert (Architect-Verdikt F-6). Der Test des Konfigurations-Rands fährt beide Eingaben für beide Schlüssel nebeneinander und erwartet je dasselbe Ergebnis; rot, sobald einer der beiden anders gelesen wird (etwa als Zeiger oder Zeichenkette) |
| `internal/hexagon/core/rules/vcs.go`, `vcs_link_targets.go` | update | Umbenennungen nur bei aktivem Schalter lesen; das Gleichsetzen im Auflöser der Link-Ziele. **Auflöser (Architect-Verdikt F-6):** `CheckVCS` baut die Vereinigung beider Pfad-Bäume heute nur mit `ignore-link-targets` (`vcs.go` Z. 53–56); mit `follow-renames` allein braucht er sie ebenso, denn Entscheidung 9 löst gegen sie auf. Rot ohne: der Testfall zu „Happy Path (Umbenennung folgen)“ mit `follow-renames` allein — ohne Baum löst kein Ziel auf, und der Nachzug bleibt Befund |
| `internal/adapter/driving/cli/config_template.go` | update | `--print-config`: Zeile des Schlüssels im `vcs`-Block |
| `.d-check.yml` | update | `follow-renames: true` im `vcs`-Block (Entscheid 7), im selben Commit wie der Übergang der ADR nach `Accepted` (Entscheid 14); Kommentar nennt die ADR |
| `docs/user/benutzerhandbuch.md`, `harness/sensors/adr-check.md` | update | nennen `ignore-link-targets` und die Grenze „anderer Name bleibt Drift“, die der Schalter ändert (Workflow-Skelett `implement-slice` Schritt 17); der Handbuch-Kopf bleibt Release-Prep |
| `README.md`, `README.de.md` | — | nennen die Grenze „anderer Name“ nicht; ein Satz zum neuen Schalter ist README-Prosa und gehört in die Release-Prep ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17), außer ein Spiegel-Test verlangt ihn schon hier |
| [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md), `docs/plan/adr/README.md` | update | Status `Accepted` und Index — durch den Architect |
| `AGENTS.md` §3.5, `harness/README.md` §Traceability rules | update | nennen [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) neben [`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) (Entscheid 14) |
| `spec/spezifikation.md` | update | **Plan-Änderung nach der Verifikation von [`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md) (V-1):** Schritt 2 beschreibt die Index-Enumeration im Nebensatz „wo die Enumeration jeden Eintrag nimmt“ mit dem heutigen Code; mit der Angleichung von `AllPaths` wird er falsch und fällt im selben Commit wie der Code. **Dazu (Architect-Verdikt F-5):** der Satz desselben Absatzes vom „benannten Preis, dass eine immutable Datei, die durch einen Gitlink **ersetzt** wird, unbemerkt bleibt“ fällt im selben Commit, **wenn** der Testfall „immutable Datei durch Gitlink ersetzt“ unten `core-drift-vcs` in beiden Modi bestätigt; bestätigt er es nicht, bleibt der Satz, und die Abweichung geht als Plan-Änderung an den Architect, nicht als Code-Änderung, die das Ergebnis erzwingt. Historie-Zeile dazu |
| `docs/user/abdeckung-tests.md` | update | von `make abdeckung` abgeleitet; `make test` hält sie gegen die Ableitung |
| Testdateien neben Adapter, Regel und Konfigurations-Rand; der Test-Doppelgänger des Ports in `internal/hexagon/core/rules/vcs_test.go` | neu / update | Happy/Boundary/Negative je Kriterium von [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in); jeder Implementierer des Ports folgt der neuen Methode |

**Spiegel vor dem Editieren auflisten** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
die Tabelle oben ist der Stand der Planung. Vor dem ersten Edit läuft die Suche nach
`ignore-link-targets` und `IgnoreLinkTargets` über das ganze Repo (ohne `done/`,
`docs/reviews/` und `CHANGELOG.md`, die Historie tragen); jeder Treffer, der in der
Tabelle fehlt, kommt als Plan-Nachtrag vor den Code. Stand der Planung (2026-11-23):

```bash
grep -rln "ignore-link-targets\|IgnoreLinkTargets" --exclude-dir=.git --exclude-dir=done . \
  | grep -vE "^./(docs/plan/planning|docs/reviews)/|CHANGELOG"
```

**Testfälle aus der Anforderung** (Negativliste nach Schritt 19 des Workflow-Skeletts; die
Liste folgt den bestätigten Entscheiden 1–14 und dem Architect-Verdikt;
Kriterien-Namen in Anführungszeichen sind die Akzeptanzkriterien von
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
in Lastenheft 0.114.3).

**Jeder Negativ-Fall läuft mit `follow-renames` allein und mit beiden Schlüsseln**
(tabellengetrieben, dasselbe Fixture zweimal; Architect-Verdikt F-2). Erwartung im zweiten
Lauf nach dem Kriterium „Boundary (Negativ-Kriterien mit beiden Schlüsseln)“: die Zusage
gilt, mit dessen einer Ausnahme — „tragen beide Ziele denselben Dateinamen, gilt die Grenze
‚gleicher Name, andere Datei‘ wie mit `vcs.ignore-link-targets` allein — kein Befund, sofern
keines der beiden Ziele teilnimmt“. Kriterien, die ihre Zusage ohne die Einschränkung „ohne
`vcs.ignore-link-targets`“ geben („Negative (Anker, Linktext)“, „Negative (zwei Quellen auf
ein Ende)“, „Negative (immutable Datei umbenannt)“), erwarten im zweiten Lauf ihren eigenen
Wortlaut. **Rot ohne die Regel:** der zweite Lauf zieht seinen Rot-Beleg aus der Mutation
seines Falls, **wenn** diese Mutation eine Umbenennung auf das Ende des Nachzugs erzeugt —
dann nehmen beide Ziele mit demselben Ende teil, gleichen nach
[`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) Entscheidung 10 auch mit beiden
Schlüsseln, und der Lauf wird grün statt rot. Erzeugt sie keine (etwa beim Umbiegen auf
eine nicht umbenannte Datei), ist der zweite Lauf ein Netz ohne eigenen Rot-Beleg und wird
im Testbericht so geführt; den Rot-Beleg der Schlüssel-Kombination tragen dort die Fälle zu
„Boundary (beide Schlüssel)“, „Negative (beide Schlüssel, Umbiegen auf den alten Namen)“
und „Negative (beide Schlüssel, Umbenennung ohne Namenswechsel)“.

- `git mv A B` und Nachzug in derselben `Accepted`-ADR, ein Commit ⇒ mit Schalter kein
  Befund, ohne Schalter `core-drift-vcs` (heutiges Verhalten, byte-identisch).
- Move und Nachzug in **verschiedenen** Commits (Kriterium „Happy Path (Umbenennung
  folgen)“: „im selben oder in einem anderen Commit“): reiner `git mv A B` in einem Commit,
  der Nachzug A → B in einem späteren; gespiegelt der Nachzug in einem früheren; Range
  über beide, `follow-renames` allein ⇒ kein Befund, Exit 0. **Rot ohne die Regel:** liest
  der Adapter nur den Schritt, der die immutable Datei ändert (Mutation), findet er keine
  Umbenennung, und der Nachzug ist Befund; der Ein-Commit-Fall oben fängt diese Mutation
  nicht, weil dort Move und Nachzug im selben Schritt liegen.
- Derselbe Nachzug, A gelöscht und B mit anderem Inhalt neu angelegt ⇒ Befund.
- Umbiegen von A auf eine bestehende, nicht umbenannte Datei C ⇒ Befund.
- Nachzug mit geändertem Anker oder Linktext ⇒ Befund.
- Umbenennen und im selben Commit ändern ⇒ Befund (keine Inhalts-Gleichheit).
- Reiner Move-Commit A → B, die Inhaltsänderung an B in einem **späteren** Commit,
  Nachzug A → B, `follow-renames` ohne `ignore-link-targets` ⇒ kein Befund (Kriterium
  „Negative (umbenennen und ändern)“, zweiter Satz). **Rot ohne die Regel:** vergleicht der
  Adapter die BASE-Fassung von A mit der HEAD-Fassung von B statt je Schritt (Mutation),
  sind die Blobs ungleich, und der Nachzug ist Befund. Die Kette darunter fährt diesen Fall
  nicht: dort steht die Inhaltsänderung zwischen zwei Umbenennungen, hier am Ende.
- Kette über drei Commits (A → B, Inhaltsänderung an B, B → C) ⇒ folgt (Entscheid 2).
- X → Y auf der einen, Y → Z auf der anderen Seite eines Merges, keine Seite Vorfahr der
  anderen; Y und Z liegen nicht in BASE (die zweite Seite legt Y nach dem Abzweig an und
  benennt es dann um), Z liegt in HEAD; Nachzug X → Z, `follow-renames` ohne
  `ignore-link-targets` ⇒ Befund (Kriterium „Boundary (Kette)“, zweiter Satz: „setzt sich
  die Kette nicht fort“). **Rot ohne die Regel:** eine ordnungsfreie Verkettung über alle
  Kanten ([`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) Entscheidung 11)
  verkettet X → Y → Z mit dem Ende Z in HEAD und lässt den Nachzug durch. Der Fall „parallele Linien, die denselben Pfad verschieden umbenennen“
  weiter unten ist ein anderer (zwei Enden einer Quelle) und fängt diese Mutation nicht.
- Zwei verschwundene oder zwei neue Pfade mit demselben Blob; Kette mit zwei Enden ⇒
  Befund (Mehrdeutigkeit).
- Merge als HEAD mit einer umbenennenden Seite ⇒ die Umbenennung trägt (Entscheid 11);
  parallele Linien, die denselben Pfad verschieden umbenennen ⇒ Befund.
- Basis kein Vorfahr von HEAD ⇒ keine Umbenennung, Urteil wie ohne Schalter
  (Entscheid 12), kein Exit 2.
- Unveränderter Link auf eine umbenannte Datei in einer ADR, die nur einen erlaubten
  Anhang bekommt ⇒ kein Befund (seiten-gleich, Entscheid 9).
- Tausch über einen Zwischennamen (A → T, B → A, T → B in drei Commits) ⇒ kein Pfad
  nimmt teil, weil jedes Ende in BASE existiert (Entscheid 9); Nachzug A → B ⇒ Befund.
  Den Rot-Beleg für die BASE-Bedingung trägt der nächste Fall, nicht dieser (ohne die
  Bedingung hätten A und B hier eine Doppelrolle). Ein Tausch der Inhalte in
  einem Commit ist keine Umbenennung (beide Pfade bestehen weiter) ⇒ Inhaltsänderung.
- A gelöscht, dann C → A, Nachzug C → A ⇒ Befund (Ende existiert in BASE, Entscheid 9) —
  mit `follow-renames` ohne `ignore-link-targets` (Kriterium „Negative (Name freigemacht
  und neu belegt)“) und mit beiden Schlüsseln (Kriterium „Boundary (Negativ-Kriterien mit
  beiden Schlüsseln)“; C und A tragen im Fixture verschiedene Dateinamen, die Ausnahme
  gleichen Dateinamens greift nicht). Den Fall gleichen Dateinamens trägt die nachsichtige
  Hälfte zwei Punkte weiter.
- Rotation A → X, danach C → A ⇒ A nimmt teil (Ende X), C nicht; Nachzug A → X ⇒ kein
  Befund, Nachzug C → A ⇒ Befund — mit `follow-renames` allein und mit beiden.
- Die nachsichtige Hälfte der drei Fälle mit beiden Schaltern (Grenze „Drei Fälle lassen
  etwas durch …“ im Lastenheft): freigemachter Name, Umbiegen von C auf eine andere Datei
  mit dem Dateinamen von C ⇒ kein Befund; zwei Quellen auf ein Ende, Umbiegen von A auf
  eine andere Datei mit dem Dateinamen von A ⇒ kein Befund; mit `follow-renames` allein
  sind beide Fälle Befund.
- Zwei Quellen auf ein Ende (Merge, A → E und B → E), A und B mit verschiedenen
  Dateinamen, Umbiegen von A auf B ⇒ Befund, mit `follow-renames` allein und mit beiden
  (Ende keiner zweiten Quelle, Entscheid 9); ohne die Klausel wäre es grün.
- Beide Schlüssel, A → B mit **anderem** Dateinamen, und ein Link, der von A auf eine
  andere Datei A′ mit dem Dateinamen von A umgebogen wird; A′ liegt in BASE und HEAD und
  nimmt an keiner Umbenennung teil ⇒ Befund; mit `ignore-link-targets` allein derselbe
  Nachzug ⇒ kein Befund (Kriterium „Boundary (beide Schlüssel)“, beide Hälften).
  **Rot ohne die Regel:** normiert die Kombination auch das teilnehmende Ziel A auf
  Dateiname und Anker, statt ihm die eigene Marke des Endes zu geben (Entscheidung 10;
  Mutation), gleichen A und A′, und der Lauf mit beiden Schlüsseln ist grün. Die Fälle zu
  `d1/a.md` und `d1/x.md` darunter sind andere Kriterien (Umbiegen auf den alten Namen,
  Umbenennung ohne Namenswechsel).
- Beide Schalter, Umbiegen von `d1/a.md` auf `d2/c.md`, das die Range nach `d3/a.md`
  umbenennt ⇒ Befund (Entscheid 10).
- Beide Schalter, `d1/x.md` → `d2/x.md`, Umbiegen auf `d3/x.md` ⇒ Befund (Entscheid 10);
  zwei nicht teilnehmende Ziele mit gleichem Namen ⇒ wie mit `ignore-link-targets` allein.
- A → B, danach B im selben Commit gelöscht und neu angelegt ⇒ Inhaltsänderung an Ort
  und Stelle, der Nachzug folgt.
- Ende der Kette in HEAD gelöscht ⇒ keine Umbenennung.
- Ende in einem Commit gelöscht und in einem späteren als andere Datei neu angelegt;
  Zwischenpfad einer Kette ebenso neu belegt ⇒ keine Umbenennung, Nachzug ist Befund
  (Neubelegung, Entscheidung 9 und 11).
- Merge, eine Seite benennt A nach B um, die andere löscht A und legt es später neu an ⇒
  Mehrdeutigkeit, Befund; überschreibt die andere Seite A an Ort und Stelle ⇒ folgt
  (Entscheidung 11).
- Nur das Ausführbar-Bit geändert ⇒ gleich; Symlink gegen Datei mit gleichem Blob ⇒
  nicht gleich (Entscheid 13).
- Link auf ein umbenanntes Verzeichnis ⇒ wie heute; Link auf eine Datei darin ⇒ folgt.
- Die immutable Datei selbst umbenannt ⇒ Befund wie heute.
- `--staged` mit `git mv` und Nachzug im Index ⇒ kein Befund; `git mv` in einem früheren
  Commit, nur der Nachzug gestagt ⇒ Befund (benannte Grenze).
- Unlesbar ist je Fixture genau eines der vier Objekte, die allein der Schalter liest
  (Kriterium „fail-closed (Zwischen-Commit)“): (a) das Commit-Objekt eines Zwischen-Commits,
  (b) der Baum eines Schritts, (c) der Baum des Eltern-Commits vor BASE, an dem eine
  Merge-Seite abzweigt, (d) ein Teilbaum, der sich zwischen einem Schritt und seinem
  Eltern-Commit unterscheidet. **Bedingung des Fixtures:** das Objekt gehört weder zum Baum
  von BASE noch zu dem von HEAD, BASE und HEAD sind lesbar, und zwischen ihnen liegt keine
  Drift — sonst endet der Lauf schon an der Enumeration mit Exit 2 oder mit einem Befund,
  und der Fall belegt nichts über den Schalter. Erwartung: mit `follow-renames` Exit 2 mit
  Hinweis auf stderr; **ohne** den Schlüssel derselbe Lauf Exit 0 ohne Befund — das trägt
  zugleich „kein Commit zwischen BASE und HEAD wird gelesen“ aus dem Kriterium „Boundary
  (ohne `vcs.follow-renames` byte-identisch)“
  ([`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) §Fitness Function nennt dieselben
  Fälle). **Rot ohne die Regel:** mit Schalter, wenn ein unlesbares Objekt die Aufzählung
  still beendet (der verbotene Diff-Walker aus §3 Zeile `git.go`; Mutation) ⇒ Exit 0 statt
  2; ohne Schalter, wenn die Port-Methode unabhängig vom Schalter gerufen wird (Mutation)
  ⇒ Exit 2 statt 0.
- Gitlink im Index unter `vcs.paths` mit `--staged` ⇒ kein Kandidat (wie im Baum); ohne
  die Angleichung wäre er einer (Rot-Beleg am Vorzustand). Das ist eine Aussage über die
  Enumeration des Ports, geprüft an seiner Ausgabe: am Befundsatz des Moduls ändert ein
  Gitlink, der **keine** immutable Datei an ihrem Pfad ersetzt, nichts, er war und bleibt
  frei (Architect-Verdikt F-4). Die eine Stelle, an der die Angleichung den Befundsatz
  ändert, trägt der nächste Fall.
- **Immutable Datei durch Gitlink ersetzt** (Architect-Verdikt F-4, F-5): eine Datei unter
  `vcs.paths`, deren BASE-Fassung `vcs.immutable-when` erfüllt, wird an ihrem Pfad durch
  einen Gitlink ersetzt — einmal als Commit mit `--range`, einmal gestagt mit `--staged`;
  mit und ohne `follow-renames` ⇒ `core-drift-vcs` (Exit 1), in beiden Modi. Erwartung nach
  Lesart des Codes, nicht nach einem Kriterium: keines nennt den Fall, das nächste ist
  „Negative“ (Datei gelöscht oder umbenannt); bestätigt der Test sie nicht, gilt der
  Vorbehalt aus §3 Zeile `spec/spezifikation.md`. **Rot ohne die Regel:** `--staged` am
  Vorzustand (`AllPaths` nimmt jeden Index-Eintrag) ⇒ nach Lesart des Codes Exit 2 statt
  Exit 1, das ist der Rot-Beleg der Angleichung; `--range` mit dem Gitlink als Datei-Pfad
  in der Enumeration des Baums (Mutation) ⇒ der Pfad liegt in beiden Mengen, geht in den
  Core-Vergleich statt in die Pfad-Stabilitäts-Prüfung, und der erwartete Befund „gelöscht
  oder umbenannt“ fehlt.
- Zehn Läufe derselben Eingabe ⇒ byte-gleiche Ausgabe (Messmethode von
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Dieselbe Topologie mit vertauschter Ordnung der Kurz-SHAs (zwei Fixtures, gleiche
  Umbenennungen, andere Commit-Kennungen) ⇒ dieselbe Abbildung; „zehn Läufe“ allein fährt
  keinen Reihenfolge-Wechsel, weil die Aufzählung nach Kurz-SHA sortiert ist
  (Verifikation von slice-291, V-2).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md)
liegt in `done/` (Anforderung, Spezifikation und ADR stehen, die Entscheide sind bestätigt),
das Architect-Verdikt zu diesem Plan liegt unter `docs/reviews/`, und das WIP-Limit ist
frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Port und Adapter brauchen
  allein mehr als eine Review-Sitzung — dann wird der Adapter-Teil ein eigener Slice
  derselben Welle vor dem Regel-Teil.
- `in-progress` → `open` (blockiert): go-git liefert für den Index oder für ältere
  Pack-Formate die Blob-Kennungen nicht so, wie die Inhalts-Gleichheit sie braucht — dann
  geht die Frage an die Anforderung zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Port, Adapter und Regel mit Tests je Kriterium und gezeigtem Rot, Spiegel
grün, ADR `Accepted`, `make gates` grün), Review- und Verifikations-Report liegen unter
`docs/reviews/`, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Auflöser der Link-Ziele trug schon für `ignore-link-targets` acht Review-Runden
  (Reports `docs/reviews/2026-10-09-slice-267-vcs-pfad-nachzug-r1.md` bis `-r8.md`); eine
  zweite Gleichsetzung an derselben Stelle kann die Filter dort verschieben
  ([`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md))
  — **Ausgang:** *entfallen.* Die bestehenden Tests von
  `vcs_pfad_nachzug_test.go` sind unverändert und grün, mehrere Tests fahren beide Schalter
  zusammen (Verifikation, Punkt 2).
- git-Fixtures, die Inhalt gleicher Länge in einem Zug neu schreiben, fallen sporadisch
  aus ([`BEO-ALL/racily-clean-git-fixture`](../observations/BEO-ALL/racily-clean-git-fixture/observation.md))
  — Umbenennungs-Fixtures mit gleichen Blobs liegen nahe an dieser Klasse —
  **Ausgang:** *entfallen.* Die Adapter-Fixtures liegen im Speicher, die CLI-Fixtures laufen über die Helfer von
  `internal/hexagon/core/coretest/`, und kein Fall schreibt dieselbe Datei gleich lang neu
  (Review R2, Verifikation).
- Der Commit-Walk liest mehr Objekte als heute; ein großer Push-Bereich in der CI wird
  langsamer — **Ausgang:** *entfallen.* Gemessen in der Verifikation über 100 Commits mit
  einem Merge, dessen eine Seite vor BASE abzweigt: Median 1,139 s ohne, 1,562 s mit
  `follow-renames`, Faktor 1,37 gegen die Schwelle 2 (fünf Läufe je Konfiguration,
  Kommando im
  [Verifikations-Report](../../../reviews/2026-11-24-slice-292-vcs-umbenennung-verify.md)).
  **Bedingung der
  Messung (Architect-Verdikt F-7):** die Range enthält einen Merge, dessen eine Seite vor
  BASE abzweigt — sonst liest der Lauf keinen Baum eines Eltern-Commits vor BASE
  ([`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) Entscheidung 13) und misst den
  teuren Fall nicht. Die Historie dieses Repos trägt nur wenige Merge-Commits
  (`git log --merges`); die Range wird um einen von ihnen gelegt, BASE zwischen Abzweig
  und Merge.
- Die Angleichung der Index-Enumeration ändert ohne den Schalter ein Ergebnis im gestagten
  Modus (Gitlink ersetzt eine immutable Datei: Exit 2 → `core-drift-vcs`, §3 Zeile
  `git.go`); Welle-Probe oder Verifier könnten das als Bruch von „ohne den Schalter nichts
  Neues“ lesen (Architect-Verdikt F-4) — **Ausgang:** *entfallen.* Der Testfall
  „immutable Datei durch Gitlink ersetzt“ läuft in beiden Modi grün, das Rot am
  Vorzustand ist gemessen (Exit 2 statt `core-drift-vcs`), und der Closure-Trigger der
  [Welle](../welle-96-vcs-folgt-umbenennung.md) §3 benennt die Abweichung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** Der Code folgte der fertigen Spec und brauchte zwei Review-Runden
  statt der vier der Spec; R1 fand keinen Befund im Produktcode, nur zwei fehlende Tests
  (M-1, M-2), die R2 schloss. Das Gegenmittel war die Mutation: Der Implementer fuhr 22,
  der Reviewer weitere 22, der Verifier elf, jede einzeln gesetzt und zurückgenommen; sie
  fand beide Test-Lücken und zeigte, welche Zweige nur äquivalente Optimierung sind. Der
  Anlass aus Issue #18 trägt am gebauten Image: Umbenennung von `spec/spezifikation.md`,
  nur Link-Ziele nachgezogen — 58 `core-drift-vcs` ohne den Schalter, 0 mit ihm, 58 bei der
  Gegenprobe (gelöscht und mit anderem Inhalt neu angelegt).
- **Was ging anders als geplant:** **(1) Eine Datei statt zwei:** der Adapter-Code liegt in
  `git.go`, die Datei wuchs auf über 1000 Zeilen; ein eigenes `renames.go` stand nicht in
  der Tabelle. **(2) Stufen:** Stufe 1 trug schon den ganzen Mechanismus, weil er eine
  Funktion ist; Stufe 2 trug die Randfall-Tests, Spezifikation und Doku. **(3) Der
  Rest-Satz der Spezifikation fiel:** ein Test bestätigte die Lesart des Codes zum Gitlink,
  deshalb fiel der Satz „bleibt unbemerkt“ zusammen mit dem Nebensatz aus V-1 der
  Verifikation von slice-291. **(4) Ein Satz der Welle** („Ohne den Schalter nichts
  Neues“) musste die eine gewollte Abweichung im gestagten Modus nennen, und die
  Blackbox-Probe fährt `vcs` gar nicht (Default-`PROBE_FORMS` ohne `--range`/`--staged`);
  die Gleichheit ohne Schalter trägt der Verifier mit einem Vorher-Image, die Probe
  gehört nicht dazu. **(5) Aufwand:** Der Auftraggeber hat die rund acht Stunden bis zum
  Spec-Stand beanstandet und das Zurückschneiden der Spec verworfen. Gemessen: Der
  Implementer-Lauf für Code, Tests und Doku dauerte rund 46 Minuten, die Review-Runden
  R1 und R2 zusammen rund 22 Minuten; ob ein kleinerer Umfang schneller gewesen wäre,
  ist nicht gemessen.
- **Steering-Loop-Eintrag:** mit diesem Slice ist keine Regel, kein Sensor und keine
  benannte Spec-Lücke verkörpert. Zwei Beobachtungen erreichen 2×
  ([`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md),
  [`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md)).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-292.md` in den zwei
  genannten Einträgen; neu angelegt, 1×:
  [`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](../observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md)
  (Beleg: slice-291, die acht Stunden Spec). Benannt, nicht gezählt: Handbuch und
  Sensor-Datei sagen „Laufzeit wächst mit der Länge der Range“, der `isAncestor`-Walk
  liest bei fremder Basis zusätzlich die Historie vor BASE (R1 I-2); der Gitlink-Satz
  gilt auch ohne den Schlüssel, steht im Handbuch aber nur im Schlüssel-Absatz (I-4);
  `entryUnreadable` führt die drei Datei-Arten ein zweites Mal neben `isFileMode`;
  der Teilbaum-Skip bei gleicher Kennung ist äquivalente Optimierung ohne Test (I-1);
  `TestVCS_FollowRenamesAusByteIdentisch` belegt nur, dass fehlender Schlüssel und
  `false` gleich laufen (I-3).
- **Folge-Slices:** keine neuen. Der Architect-Akt nach Entscheidung 14 ([`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) auf
  `Accepted`, `.d-check.yml`, AGENTS.md §3.5, harness/README.md) läuft als eigener Commit
  derselben Welle; die Probe am Image gehört in `welle-96-results.md`.
- **Risiken aus §6:** alle vier *entfallen* — siehe §6.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) — Port,
git-Adapter, Regel, Konfigurations-Rand, Doku; keine Sub-Area mit eigener
Modus-Deklaration (`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand
2026-11-23), Sub-Area `*`. Treffer:
[`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md)
(verkörpert in welle-74, Restklasse offen — Risiko 1 in §6),
[`BEO-ALL/racily-clean-git-fixture`](../observations/BEO-ALL/racily-clean-git-fixture/observation.md)
(offen — Risiko 2 in §6),
[`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
(4× — die Spiegel-Suche in §3 nimmt Handbuch, READMEs und Sensor-Datei auf),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(7× — die Negativliste in §3),
[`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/observation.md)
(offen, 1× — Welle-Out-of-Scope) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — der Port liest Bäume
außerhalb der Scan-Menge; die Zusage dort schreibt
[`slice-291`](../done/slice-291-vcs-umbenennung-anforderung-und-adr.md) fest). Der Zähler
folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-23: `image-scan` grün (2026-11-23T12:51:15Z),
`upstream-drift` rot (2026-11-23T08:29:02Z, vor den Pin-Hebungen desselben Tages); beim
Start dieses Slice neu lesen.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
