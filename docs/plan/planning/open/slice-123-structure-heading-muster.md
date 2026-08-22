# Slice slice-123: `structure` — jede Überschrift des Abschnitts matcht ein Muster

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-82-config-flaechen](../welle-82-config-flaechen.md) (zugeordnet
bei der Eröffnung).

**Bezug:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die zu erweiternde Anforderung),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`MR-025`](../../../../harness/conventions.md#mr-025); Anlass ist die
Ersatz-Konstruktion aus welle-80 — eine ausgeschriebene Präfix-Negation, weil
RE2 keinen Lookahead kennt, mit einem belegten **stillen Falsch-Negativ**.

**Berührte Spec-Stellen:**
[`DC-FA-STRUCT-001.a`](../../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
(Bedingungen), §2-Schema (`structure[].*`), §4 (`section-*`) — der Verweis zeigt aufwärts.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die Bedingung „**jede** Überschrift dieses Abschnitts genügt einem Muster" ist
mit den heutigen Schlüsseln nur als Negation ausdrückbar — und die war in
welle-80 nachweislich falsch: sie sprach nicht die Heading-Lexik des Moduls
(führender Weißraum, Tab als Trenner), also entkam eine eingerückte Sektion
still. Ein eigener Schlüssel dreht die Aussage um: **positiv, je Überschrift,
mit der Lexik des Moduls** — das Modul kennt seine Überschriften bereits, es
muss sie nur einzeln prüfen statt den Abschnittstext als Ganzes.

## 2. Vorgehen

1. **CR-Commit zuerst:** Lastenheft
   [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
   um die Bedingung erweitern — Beschreibung, Akzeptanzkriterien (Happy-Path, Negativ je
   verletzender Überschrift mit ihrer Zeile als Befund-Ort, Default
   byte-identisch ohne den Schlüssel, fail-closed bei ungültigem Muster),
   §7-Historie.
2. **ADR der Welle** um die zweite Entscheidung ergänzen: positive
   Je-Überschrift-Prüfung statt Negation im Abschnittstext; Ebenen-Wahl
   (welche Überschriften-Ebenen der Abschnitt umfasst); Befund je Überschrift
   statt einer je Abschnitt, damit die Zeile zeigt, wo es klemmt.
3. **Spezifikation:** Bedingung im Algorithmus, §2-Schema, §4-Zeile
   (bestehender oder neuer Grund-Code — die Entscheidung gehört in die ADR und
   folgt dem Kriterium „andere Reparatur ⇒ eigener Code").
4. **Code + Tests:** die Überschriften des Abschnitts liegen im Modul bereits
   vor (geteilte Lexik) — sie werden einzeln geprüft; Tests für alle
   Akzeptanzkriterien, dazu die Fälle, an denen die alte Negation scheiterte
   (eingerückt, Tab-getrennt, vierte Ebene, Überschrift nur aus Inline-Code).
5. **Das eigene Profil umstellen:** die ausgeschriebene Negation wird durch den
   neuen Schlüssel ersetzt — vorher/nachher gemessen, beide Male grün, und die
   Gegenprobe, an der die Negation still war, wird jetzt rot.
6. **Spiegel** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
   Anforderung, Algorithmus, §2-Schema, §4-Tabelle, Klartexte,
   `--print-config`-Vorlage, Config-Kommentar (Handbuch ist Release-Prep).
7. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine weiteren `structure`-Bedingungen** — nur diese eine.
- **Kein Handbuch, kein CHANGELOG** (slice-125).
- **Keine Default-Änderung.**

## 4. Definition of Done

- [ ] CR-Commit (Lastenheft allein) liegt **vor** Spezifikation und Code.
- [ ] Der Schlüssel prüft **je Überschrift** mit der Modul-Lexik; Befund nennt
      die verletzende Zeile.
- [ ] Das eigene Profil nutzt ihn statt der Negation; die Gegenprobe, an der
      die Negation still war (eingerückte Sektion), ist jetzt rot.
- [ ] Default-Beweis byte-identisch; `make gates` grün; unabhängiger Review;
      Closure-Notiz; Register gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Welche Überschriften gehören zum Abschnitt?** Die Ebenen-Frage entscheidet
  über Falsch-Positive (eine tiefere Ebene, die nie gemeint war). Sie gehört in
  die ADR und in den Vertrag, nicht in den Code. — **Ausgang:** *(bei Closure)*
- **Ein Befund je Überschrift statt je Abschnitt** ändert die Befund-Zahl —
  das ist gewollt (die Zeile zeigt, wo es klemmt), muss aber zugesagt sein. —
  **Ausgang:** *(bei Closure)*
- **Der Umstieg des eigenen Profils ist der eigentliche Beweis:** wenn die neue
  Bedingung die alte nicht deckt, zeigt es sich hier. — **Ausgang:** *(bei
  Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): Wellen-Eröffnung (eingetreten).

**Rückführungen:** `in-progress` → `next`, falls die Ebenen-Frage einen zweiten
Schlüssel verlangt (dann erst Vertrag klären).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Kern (GF), Config-Rand (GF), Spec-Straten (GF),
  eigenes Prüf-Profil (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): **BEO-003**
  (geteilte Lexik driftet an den Rändern) ist einschlägig — der neue Schlüssel
  muss die Lexik des Moduls **benutzen**, nicht nachbauen; BEO-002 als
  Spiegel-Pflicht, BEO-006/009/010 als Arbeitsregeln.

Slice-ID: slice-123. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
`structure` (Kern `rules/`), Config-Rand, Spec, eigenes Profil. Gates:
`make test` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — additive Erweiterung einer eigenen,
spezifizierten Anforderung; sie löst zugleich eine Ersatz-Konstruktion ab.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
