# Slice slice-116: ADR-Neuzugangs-Regel + Erstanwendung an den `Proposed`-ADRs

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-80-struktur-ids](../welle-80-struktur-ids.md) (zugeordnet bei
der Eröffnung).

**Bezug:** AGENTS §3.5 (ADR-Immutabilität — `Accepted` bleibt unverändert,
`make adr-check`), die ADR-Index-Konvention aus slice-113, die beiden
`Proposed`-ADRs [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md) und
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md),
Baseline-ADR-Template (`Schärft:` nennt die Kennung, wo das Zielelement eine
trägt; Link auf Überschrift oder Sektion), Reviewer-Skill
`.harness/skills/reviewer.md` (Referenz-Richtungs-Anker), Entscheid D4.

**Berührte Spec-Stellen:** — (ADR-Felder und Harness-Skills; keine Spec-Zeile).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die neue Form wird **angewandt**, nicht nur deklariert: die beiden
`Proposed`-ADRs ziehen `Schärft:`/`Bezug:` auf die jetzt vergebenen
`SPEC-*`/`ARC-*` (kein Immutabilitäts-Bruch — `adr-check` greift nur bei
`Accepted`), das Slice-Kopf-Feld „Berührte Spec-Stellen" im Haus-Stil nennt
Kennungen statt `§N`, der Reviewer-Skill trägt den MEDIUM-Anker „Kennung
statt §-Anker, wo das Zielelement eine trägt" (Zwei-Formen-Welt: `Accepted`-
Bestand bleibt §, Neuzugänge Kennung), und der ADR-Index-Kopf zeigt ein
Beispiel beider Formen. Damit hat die Welle ihren ersten lebenden Konsumenten
auch auf der ADR-Seite — und der Review hat den Anker, um die Regel künftig zu
halten.

## 2. Vorgehen

1. **Messen:** `Schärft:`-Formen über alle ADRs (Stand der Inventur 2026-08-22:
   46 Felder; 28 auf `.a`-Sektionen, 10 „keine Spec-Stelle", 5 Architektur
   §-basiert, 2 ohne Anker, 1 `—`; 11 ohne Feld) — gemessen, die Zahlen in die
   Closure-Notiz.
2. **[ADR-0050](../../adr/0050-fence-unclosed-in-spans.md)/[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) nachziehen:** Felder auf Kennungen, wo das Zielelement
   jetzt eine trägt (sonst bleibt der §-Anker); `## Geschichte`-Zeile je ADR.
   `make adr-check` muss **stumm** bleiben (Proposed) — Probe.
3. **ADR-Index-Kopf** (`docs/plan/adr/README.md`): Beispiel-Zeile für beide
   Formen; **Reviewer-Skill** 1.5.0 → 1.6.0: Anker im MEDIUM-Block, Version
   und Datum.
4. **Haus-Stil der Slices:** das Kopf-Feld „Berührte Spec-Stellen" in diesem
   und den künftigen Slices nennt Kennungen; die welle-80-Slices 114/115
   tragen sie nach der Vergabe bereits in ihrer Closure-Notiz.
5. **Spiegel** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
   ADR-Index-Kopf, Reviewer-Skill, AGENTS §5 (Satz aus slice-113 prüfen),
   `harness/README.md` Guides-Zeile zum Reviewer-Skill (Version?).
6. Unabhängiger Review; Closure — und die **Wellen-Closure** (Ergebnisnotiz,
   Register-Lese-Schritt, `make fullbuild`).

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Retarget eines `Accepted`-ADR** — auch nicht „nur das Schärft-Feld";
  `adr-check` würde es melden, und die Konvention erlaubt die alte Form.
- **Keine Template-Änderung** im vendored Baseline-Baum (derivativ).
- **Kein Produkt-Code.**

## 4. Definition of Done

- [x] [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md)/[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) mit Kennungen in `Schärft:`/`Bezug:` + Geschichte-Zeile;
      `make adr-check` stumm (Probe dokumentiert).
- [x] ADR-Index-Kopf mit Beispiel beider Formen; Reviewer-Skill 1.6.0 mit
      Anker.
- [x] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) ist gleichzeitig Gegenstand laufender Fortschreibungen** (zuletzt
  slice-112); ein Nachzug der Felder darf keine Entscheidung umschreiben — nur
  Adress-Form. — **Ausgang:** entfallen — der Review hat den Diff gegen die
  Immutabilitäts-Grenze geprüft: Status-Zeilen unverändert, Nachträge im
  Geschichte-Anhang, keine Entscheidung umgeschrieben. Der Kern blieb, das
  Feld kam dazu.
- **Reviewer-Anker als MEDIUM** könnte Alt-ADR-Zitate in Reviews als Befund
  lesen — der Anker nennt ausdrücklich die Zwei-Formen-Regel. — **Ausgang:**
  **eingetreten in der Gegenrichtung.** Nicht die Scheinbefund-Seite war offen
  (die schloss der Anker), sondern die andere: eine ADR, die im selben Slice
  `Accepted` wird, wäre unter die Ausnahme gefallen. Der Anker sagt jetzt, dass
  die Form beim **Schreiben** gilt, nicht beim Status-Übergang, und nennt den
  Skopus wörtlich wie die Konventions-Quelle (Review F-4).

## 6. Trigger

**Start** (`open` → `in-progress`): slice-114 **und** slice-115 in `done/`
(Kennungen existieren in beiden Straten).

**Rückführungen:** `in-progress` → `next`, falls ein `Proposed`-ADR vorher
`Accepted` wird (dann entfällt sein Nachzug — Konvention, kein Bruch).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** ADR-Bestand (`docs/plan/adr/`, GF), Harness-Skills
  (`.harness/skills/`, GF), Planungs-Haus-Stil (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-002
  ([`MR-025`](../../../../harness/conventions.md#mr-025)), BEO-006/009 Arbeitsregeln; nichts Einschlägiges darüber hinaus.

Slice-ID: slice-116. Betroffene IDs:
[ADR-0050](../../adr/0050-fence-unclosed-in-spans.md),
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
(Proposed), [`MR-025`](../../../../harness/conventions.md#mr-025). Module:
ADR-Index, Reviewer-Skill, Slice-Haus-Stil. Gates: `make adr-check` (Probe),
`make doc-check`, `make gates`; Wellen-Closure `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Anwendung einer Konvention an eigenen
Artefakten; kein Legacy-Import.

## 9. Closure-Notiz (nach `done/`)

**Geliefert:** die Adressierungs-Form ist angewandt, nicht nur deklariert.
Gemessen zuerst: 57 ADRs, 46 mit `Schärft:`-Feld, 55 `Accepted` und damit
immutabel, **zwei** `Proposed` — genau die zwei tragen die Form jetzt. Eine
bekam die Struktur-Kennung neben ihrer Verfeinerungs-Kennung, die andere
erstmals überhaupt ein `Schärft:`-Feld (es fehlte) mit Algorithmus-Abschnitt,
vier Grund-Code-Festlegungen und dem Konfigurations-Schema. Dazu die drei
Regel-Stellen: der ADR-Index zeigt beide Formen an einem Beispiel, der
Reviewer-Skill (1.6.0) trägt den Anker samt Skopus, und `AGENTS.md` §5 sagt,
was das Slice-Kopf-Feld nennt.

**Review** ([Report](../../../reviews/2026-08-22-slice-116-adr-neuzugangs-regel-review.md)):
0 HIGH, 2 MEDIUM, 3 LOW — der strengste Review der Welle, und alle fünf
Befunde saßen. Eingearbeitet.

**Was ging anders als geplant — dreimal derselbe Kern: die Form, die ich
einführe, gilt auch für mich.**
1. Mein Index-**Beispiel** für die neue Form war selbst falsch ausgezeichnet:
   backslash-escapte Backticks **innerhalb** einer Code-Spanne kennt CommonMark
   nicht. Kein Gate sah es — die Backtick-Parität war gerade, also schwieg
   `spans`. Ein Beispiel ist Prosa mit Anspruch: es zeigt eine Form und muss
   die Form, in der es geschrieben ist, selbst einhalten.
2. Die **Richtung** der Geschichte liest man am Bestand der Datei, nicht an
   der Gewohnheit: dieselbe Zeile ging in eine aufsteigende Liste oben hinein
   und in eine absteigende unten — im selben Commit.
3. Das erstmals angelegte `Schärft:`-Feld war **unvollständig**: es nannte den
   Algorithmus und die Grund-Codes, nicht aber das Konfigurations-Schema,
   obwohl die ADR sieben Schlüssel samt fail-closed-Rändern festlegt. Wer eine
   Adressierungs-Form einführt, muss zuerst die Menge der adressierten Stellen
   vollständig kennen — die Spiegel-Klassen stehen in
   [`MR-025`](../../../../harness/conventions.md#mr-025), das Config-Schema ist
   eine davon.

- **Steering-Loop-Eintrag:** Reviewer-Skill geschärft: die Adressierungs-Form
  eines Neuzugangs ist ein MEDIUM-Anker, der `Accepted`-Bestand ausdrücklich
  ausgenommen — liegt in `.harness/skills/reviewer.md §Repo-spezifische Anker`.
  Kein Auslöser aus dem Register.
- **Beobachtungs-Register (`../observations.md`):** keine neue Beobachtung.
  Die unvollständige Spiegel-Menge ist die verkörperte Klasse BEO-002, hier in
  ihrer Feld-Form; zitiert statt neu formuliert.
- **Folge-Slices:** keiner — mit diesem Slice ist die Welle vollständig; es
  folgt die Wellen-Closure.
- **Risiken aus §6:** beide mit Ausgang (§5) — eines entfallen, eines in der
  Gegenrichtung eingetreten und geschlossen.
- **Drei Paarungen:** Wellen-Slice — die Paarungen prüft die Welle-Closure.
