# Slice slice-082: Release-Prep-Regel — neuer Handbuch-§4-Abschnitt ist eine eigene Aufgabe

**Status:** In Arbeit (welle-66).

**Welle:** welle-66-release-prep-aufgabenregel (Trigger: slice-072-Closure — die
strukturelle Ursache der §4-Erosion blieb nach der redaktionellen Bereinigung
bestehen; Auftraggeber-Aufnahme 2026-07-19).

**Bezug:** Prozess-Härtung der [Releasing-Checkliste](../../../user/releasing.md)
§Release-Prep gegen den in
[slice-072](../done/slice-072-handbuch-aufgabenorientierung.md) §4 benannten
**strukturellen** Erosions-Mechanismus des
[Benutzerhandbuchs](../../../user/benutzerhandbuch.md) §4, gemessen am
[Benutzerhandbuch-Standard](../../../user/benutzerhandbuch-standard.md) §2/§5.
**Kein Change Request** (kein Verhaltens-/Vertragsdelta — das Produkt bleibt
unberührt), **kein ADR** (keine Architekturentscheidung), **kein Release** (nur
Prozess-/Nutzer-Dokumentation). Betrifft ausschließlich `docs/user/releasing.md`.

**Autor:** pt9912. **Datum:** 2026-07-19.

---

## 1. Ziel

slice-072 hat §4 des Benutzerhandbuchs **redaktionell** aufgeräumt (den
~330-Zeilen-§4.12-Monolithen in aufgabenorientierte §4.12–§4.16 aufgetrennt),
aber die **Ursache** nicht beseitigt: §4.12 wuchs über die Slices 074–081, weil
**jeder** Feature-Slice seine neue Fähigkeit im Release-Prep an die bestehende
§4.12-Aufgabe **anhängte**, statt eine eigene Aufgabe zu schreiben. Ohne eine
Regel am Release-Prep-Punkt erodiert §4 nach dem nächsten Feature erneut.

Der Slice ergänzt die Release-Prep-Checkliste um genau diese Disziplin: **ein
neuer §4-Abschnitt für ein neues Feature ist eine eigene Aufgabe
(Ziel/Vorgehen/Ergebnis), keine Anhängung an eine bestehende Aufgabe.**

## 2. Entscheidungen / Regel

- **Ort: `releasing.md` §Release-Prep, Punkt 4** („Prosa-Currency von Hand
  nachziehen"). Das ist der Bindepunkt, an dem die Erosion real passierte — jeder
  Feature-Slice fasst dort das Handbuch an. Die Regel steht neben der bestehenden
  „ggf. neue Feature-Abschnitte (§5/§6)"-Notiz, die den §4-Fall bisher **nicht**
  adressierte.
- **Ehrlich unenforced.** Kein Gate erzwingt Aufgabenorientierung — sie ist eine
  Erkenntnis-, keine Laufzeit-Eigenschaft (so bereits slice-072 §4). Die Regel
  ist die **billigste dauerhafte** Sicherung (Checkliste am richtigen
  Bindepunkt), keine harte Garantie; das wird in der Regel ausdrücklich benannt,
  nicht kaschiert — konsistent mit dem übrigen Punkt 4 („kein Gate erzwingt sie").
- **Verweisen, nicht duplizieren.** Die Regel nennt das Kriterium
  (Benutzerhandbuch-Standard §2/§5) und die Präzedenz (slice-072: §4.12 wuchs auf
  ~330 Zeilen / 8 Themen), statt die Standard-Prosa zu wiederholen.
- **Kein Scope-Kriechen.** Ein heuristisches Gate (z. B. §4-Abschnitt über N
  Zeilen / M Themen) ist ausdrücklich **nicht** Gegenstand — subjektiv, laut,
  falsch-positiv-anfällig; verworfen zugunsten der Checklisten-Disziplin.

## 3. Definition of Done

- [ ] `releasing.md` §Release-Prep Punkt 4 trägt die §4-Aufgabendisziplin-Regel
  (neuer §4-Abschnitt = eigene Aufgabe mit Ziel/Vorgehen/Ergebnis, keine
  Anhängung), mit Kriterium-Verweis (Standard §2/§5) und slice-072-Präzedenz.
- [ ] Die Regel benennt ausdrücklich, dass **kein Gate** sie erzwingt (wie die
  übrige Prosa-Currency-Liste dort).
- [ ] Kein Verhaltens-/Vertragsdelta; kein Release; `make gates` grün
  (`links`/`anchors`/`ids` auf `releasing.md`).
- [ ] Abschluss-Gegenprobe: die Regel steht widerspruchsfrei zur bestehenden
  Punkt-4-Struktur (keine Doppelung mit der §5/§6-Notiz, gleicher „kein
  Gate"-Ton).

## 4. Risiken / offene Punkte

- **Die Regel ist unenforced** — eine Release-Prep-Disziplin, kein Gate. Ein
  künftiger Release kann sie ignorieren; sie macht die Erwartung nur am richtigen
  Punkt **sichtbar**. Das ist die inhärente Grenze (Aufgabenorientierung ist
  nicht maschinell prüfbar), kein Slice-Mangel.
- **Checklisten-Erosion.** Auch die Release-Prep-Liste selbst kann wachsen oder
  veralten; sie ist aber der etablierte, gelebte Ort solcher Disziplin (Punkt 4
  existiert bereits als Prosa-Currency-Sammler).

## 5. Trigger

slice-072-Closure (2026-07-19): der Slice räumte §4 redaktionell auf, aber der in
seinem §4 notierte **offene Designpunkt** (die Release-Prep-Regel) war kein
DoD-Punkt und blieb offen. Auftraggeber-Entscheid 2026-07-19: als eigener
Folge-Slice aufnehmen statt in slice-072 querzuschneiden.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc führt. `docs/user/releasing.md` ist Betriebs-/Prozess-Doku
(Rang 6) ohne Brownfield-Spec; die Regel ist eine Checklisten-Ergänzung im
gelebten Stil des bestehenden Punkt 4.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend — wird bei Abschluss mit Commit-Hash, Review-Verdikt und Lerneintrag
gefüllt._
