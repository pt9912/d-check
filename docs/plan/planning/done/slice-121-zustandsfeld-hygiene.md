# Slice slice-121: Zustandsfeld-Hygiene — falsche Zustände, ein stiller Anker, eine Vorrangregel

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** ohne Welle — drei unabhängige Hygiene-Punkte ohne gemeinsame
Closure-Bedingung jenseits der DoD (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.5 (ADRs sind nach `Accepted`
immutabel) und §3.7 (Kommentare und Zustandsfelder) — die Klärung ihres
Verhältnisses ist eine **Repo-Aussage über das eigene Gate**, keine
Baseline-Regel; Reviewer-Skill (der
HIGH-Anker *Zustandsfeld trägt Chronik*); Baseline-Regelwerk
`grundlagen-harness-dateien.md` §Was ein Kommentar trägt; die drei Befunde
stammen aus den Reviews der welle-81
([slice-119](../done/slice-119-kopf-zustandsfelder.md),
[slice-120](../done/slice-120-register-und-drift-log.md)) und sind dort als
Folgepunkte benannt.

**Berührte Spec-Stellen:** — (Briefing, Register-Bestand, Release-Register;
keine Spec-Zeile).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Drei Zustandsfelder sagen etwas, das nicht stimmt — jedes auf seine Art:

1. **Elf `done/`-Slices behaupten einen falschen Lifecycle-Zustand.** Von 119
   Slices im Ruheort tragen 90 ein historisches `Status`-Kopffeld; **elf** davon
   beginnen mit `open`, `in-progress` oder `In Arbeit`, obwohl die Datei in
   `done/` liegt — sie sagen etwas Falsches über die Gegenwart. Die übrigen 79
   sagen nichts Falsches; **53** von ihnen tragen neben dem Wert eine Chronik,
   und das ist eine **benannte Bestands-Ausnahme**, keine Nachlässigkeit: sie
   sind eingefrorene Lauf-Belege, ihr Lifecycle-Zustand ist ohnehin das
   Verzeichnis, und das Feld hat dort keine Funktion. Das Briefing erklärt das historische Feld für Alt-Slices — ein Feld,
   das dem Verzeichnis **widerspricht**, deckt es nicht. Die zwölf werden auf
   den wahren Zustand gesetzt; die übrigen 79 bleiben unberührt
   (Auftraggeber-Entscheid 2026-08-22: „nur die Widersprüche heilen").
2. **Das Release-Register trägt einen Anker zu viel.** `version.md` sagt in
   seiner eigenen Regel, der Anker gehöre an **die aktuelle** Version — „sonst
   bleiben veraltete feste Pins auflösbar und die Anker-Kaskade meldet den
   vergessenen Bump nicht". Gemessen tragen **zwei** von 63 Zeilen einen: der
   v0.62.0-Release hat den alten nicht entfernt. Der Vergessens-Detektor ist
   damit für Pins auf die Vorgänger-Version entschärft.
3. **Das Verhältnis von §3.5 und §3.7 ist ungeklärt.** §3.5 macht eine
   `Accepted`-ADR immutabel, §3.7 gilt für Zustandsfelder — und genau eine ADR
   trägt ein Statusfeld mit Chronik. Geklärt wird es nach der **Maschine**:
   `adr-check` nimmt die Kopf-Status-Zeile ausdrücklich aus dem Kern-Vergleich
   und lässt den Übergang zu; das Feld ist also **kein** Sonderfall, §3.5
   schützt den Kern und nicht diese Zeile.

## 2. Vorgehen

1. **Messen statt annehmen:** die zwölf widersprüchlichen Felder namentlich
   auflisten (Datei, behaupteter Zustand) und je Datei den wahren Zustand aus
   der Verzeichnis-Position ableiten.
2. **Heilen:** je Feld den Zustand auf `done` setzen — **ohne** Chronik, ohne
   Datum, ohne Begründung im Feld. Wo das Feld bereits eine Chronik trägt (ein
   Fall sagt `done — abgeschlossen am …`), fällt sie mit.
3. **Anker:** den Anker der Vorgänger-Version aus `version.md` entfernen, so
   dass genau einer bleibt; die Regel im Kopf der Datei bleibt unverändert —
   sie war nie falsch, nur unbefolgt. Gegenprobe: ein Pin auf die
   Vorgänger-Version wird wieder rot.
4. **Vorrangregel:** in §3.7 einen Satz, der die Kollision auflöst — das
   Statusfeld einer `Accepted`-ADR ist Teil des immutablen Kerns, §3.5 sticht;
   für **neue** ADRs gilt §3.7 ab dem ersten Schreiben. Der Reviewer-Anker
   bekommt dieselbe Nicht-Melde-Ausnahme.
5. **Spiegel-Liste vor dem Editieren** (per `grep` abgeleitet, nicht erinnert):
   Briefing §3.5/§3.7, Reviewer-Skill, `AGENTS.md` §5 (die Aussage über
   Alt-Slices), `version.md`-Kopfregel, `harness/README.md`.
6. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Die 79 übrigen `Status`-Felder bleiben.** Sie widersprechen nichts; das
  Briefing erklärt sie. Sie zu entfernen hieße, 78 eingefrorene Lauf-Belege
  anzufassen — ausdrücklich verworfen.
- **Keine ADR wird geheilt.** Die Statuszeile mit Chronik bleibt, wie sie ist;
  dieser Slice klärt nur das Verhältnis der beiden Regeln. Dass sie geheilt
  **werden dürfte** (die Maschine lässt es zu), ist damit gesagt — getan wird
  es nicht.
- **Kein Produkt-Code, kein Release.** Die drei Modul-Erweiterungen sind eine
  eigene Welle.

## 4. Definition of Done

- [x] Kein `done/`-Slice behauptet mehr einen Zustand, der seinem Verzeichnis
      widerspricht (gemessen: vorher elf, nachher null). Die Chronik-Felder
      bleiben — als **benannte** Ausnahme in §3.7, nicht stillschweigend.
- [x] `version.md` trägt genau **einen** Anker; Gegenprobe belegt, dass ein
      Pin auf die Vorgänger-Version wieder meldet.
- [x] §3.7 trägt die Vorrangregel gegenüber §3.5; der Reviewer-Anker nennt
      dieselbe Ausnahme.
- [x] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Ein `done/`-Slice ist ein Lauf-Beleg.** Sein Statusfeld zu ändern heißt,
  ein eingefrorenes Dokument anzufassen. Gerechtfertigt ist das nur, weil das
  Feld eine **falsche Gegenwarts-Aussage** macht — nicht, weil es alt ist. Die
  Grenze steht in §3. — **Ausgang:** **eingetreten und korrigiert.** Der erste
  Anlauf fasste zwei Felder mehr an, als die Klasse hergab — geheilt wurde nach
  Suchstring statt nach Klasse; der Review fand es, die zwei sind zurückgenommen.
- **Der Anker-Fix schärft einen Detektor:** danach melden Pins auf die
  Vorgänger-Version wieder. Vor dem Entfernen prüfen, ob im Baum ein solcher
  Pin lebt — sonst wird der Slice selbst rot. — **Ausgang:** entfallen — kein
  lebender Verweis hing am entfernten Anker (der einzige feste Pin zeigt auf die
  aktuelle Version); die Gegenprobe belegt beide Richtungen: vorher stumpf,
  nachher scharf.
- **Eine Vorrangregel darf kein Schlupfloch werden:** sie gilt für den
  immutablen Kern einer `Accepted`-ADR, nicht als allgemeine Bestandsgrenze
  für Zustandsfelder. — **Ausgang:** **eingetreten, in der schlimmeren Form:**
  die Regel war nicht zu weit, sondern **falsch** — sie behauptete eine
  Immutabilität, die das Gate nicht kennt. Neu gefasst nach dem, was die
  Maschine tut.

## 6. Trigger

**Start** (`open` → `in-progress`): Auftraggeber-Entscheid 2026-08-22 („erst
Hygiene, dann CR-Welle") — eingetreten.

**Rückführungen:** `in-progress` → `next`, falls das Heilen der zwölf Felder
zeigt, dass der wahre Zustand einer Datei nicht aus dem Verzeichnis ableitbar
ist (dann erst klären, dann heilen).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Planungs-Bestand (`docs/plan/planning/done/`, GF),
  Release-Register (GF), Briefing und Reviewer-Skill (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-008 bei 3
  (nicht einschlägig — kein Pin), BEO-006/009/010 als Arbeitsregeln, BEO-002
  als Spiegel-Pflicht (§2 Schritt 5).

Slice-ID: slice-121. Betroffene IDs: — (Form-Regeln des Briefings). Module:
Planungs-Bestand, Release-Register, Briefing, Reviewer-Skill. Gates:
`make doc-check` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Bestands-Hygiene an eigenen Artefakten
nach einer frisch adoptierten Regel; kein Legacy-Import.

## 9. Closure-Notiz (nach `done/`)

**Geliefert:** elf `done/`-Slices behaupten keinen falschen Lifecycle-Zustand
mehr; das Release-Register trägt wieder genau einen Anker, und die Gegenprobe
belegt, dass der Vergessens-Detektor damit scharf ist (vorher stumpf); das
Verhältnis von §3.5 und §3.7 steht — geklärt nach dem, was das Gate tut, nicht
nach dem, was plausibel klingt.

**Review** ([Report](../../../reviews/2026-08-22-slice-121-zustandsfeld-hygiene-review.md)):
merge-blockierend — **2 HIGH**, 5 MEDIUM, 3 LOW, 1 INFO; alle eingearbeitet.
Der strengste Befund-Satz dieser Sitzung.

**Was ging anders als geplant — zweimal habe ich eine Aussage gemacht statt
gemessen:**
1. **Geheilt wurde nach Suchstring, nicht nach Klasse.** Zwei Felder trugen
   zufällig genau den Wortlaut, nach dem ich gesucht hatte; drei
   gleichbedeutende ohne dieses Wort blieben stehen, und 53 der 90 Felder
   tragen weiter Chronik — während die DoD „nachher null" behauptete. Die zwei
   Mitnahmen sind zurückgenommen: der Slice tut, was entschieden war, und die
   Chronik-Felder stehen jetzt als **benannte** Ausnahme im Briefing statt
   stillschweigend.
2. **Die Vorrangregel behauptete eine Immutabilität, die es nicht gibt.** Ich
   schrieb, das Status-Feld einer `Accepted`-ADR gehöre zum immutablen Kern.
   Vier Stellen sagen übereinstimmend das Gegenteil — §3.5 selbst, die
   Sensors-Tabelle, die Prüf-Config (`status-line`/`head-allow`) und der Code,
   der die Kopf-Status-Zeile aus dem Kern-Vergleich **streicht**. Eine
   Kollision, die ich lösen wollte, gab es nie; ich hatte sie erfunden. Wer
   zwei Regeln gegeneinander stellt, muss zuerst die Maschine fragen, welche
   von beiden überhaupt greift.

- **Steering-Loop-Eintrag:** Guide geschärft: §3.7 trägt jetzt das Verhältnis
  zu §3.5 **nach der Maschine** und die benannte Bestands-Ausnahme für
  historische Status-Felder — liegt in [`AGENTS.md`](../../../../AGENTS.md)
  §3.7; der Reviewer-Skill (1.8.0) nennt dieselbe Ausnahme.
- **Beobachtungs-Register (`../observations.md`):** keine neue Beobachtung.
  Beide Fehler sind die verkörperte Klasse BEO-002 — ein Rand, der eine
  Aussage trägt, die der Bestand nicht deckt.
- **Folge-Slices:** keiner. **Benannt, ohne Slice:** 53 historische
  Status-Felder tragen weiter Chronik — bewusst, jetzt als Ausnahme deklariert.
- **Risiken aus §6:** alle drei mit Ausgang (§5) — zwei eingetreten, eines
  entfallen.
- **Drei Paarungen** (wellenloser Slice, hier geprüft): Anker — der
  Steering-Loop-Eintrag nennt seinen Ort; Folge-Slice — keiner, also keine
  Datei in `open/`; Register — keine neue Beobachtung, BEO-002 zitiert.
