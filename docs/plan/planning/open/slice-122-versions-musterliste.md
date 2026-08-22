# Slice slice-122: `versions` — mehrere Muster-Quellen-Paare statt eines

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-82-config-flaechen](../welle-82-config-flaechen.md) (zugeordnet
bei der Eröffnung).

**Bezug:**
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
(die zu erweiternde Anforderung),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Byte-Identität ohne den neuen Schlüssel),
[`MR-025`](../../../../harness/conventions.md#mr-025) (Semantik-Fläche ⇒
Spiegel-Liste vor dem Editieren); Anlass ist die 3×-Form der Beobachtung
BEO-008 im Register, die mit genau einem Muster nicht baubar ist.

**Berührte Spec-Stellen:**
[`DC-FA-VER-001.a`](../../../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions)
(Algorithmus), §2-Schema (`versions.*`), §4 (`version-stale`) — der Verweis zeigt aufwärts.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Das Modul `versions` prüft heute **ein** Pin-Muster gegen **eine**
Versions-Quelle. Damit ist die im Register benannte mechanische Form nicht
baubar: der Baseline-Tag in URLs und Prosa müsste gegen den §Baseline-Pin
geprüft werden, **zusätzlich** zum bestehenden Image-Pin gegen die
Versions-Datei. Die Erweiterung ist eine **Liste** von Paaren
(Muster + Quelle + eigene Ausnahmen); der Einzel-Schlüssel bleibt als
Kurzform gültig, damit keine bestehende Konfiguration bricht.

## 2. Vorgehen

1. **CR-Commit zuerst (Doc führt):** Lastenheft-Version bump,
   [`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
   um die Listen-Form erweitern — Beschreibung, Akzeptanzkriterien (Happy-Path
   mit zwei Paaren, Negativ je Paar, Default byte-identisch ohne den neuen
   Schlüssel, fail-closed bei leerem/ungültigem Paar), §7-Historie.
2. **ADR der Welle anlegen** mit der ersten Entscheidung: Liste statt Skalar,
   Kurzform bleibt, jedes Paar trägt seine eigenen Ausnahmen; Befund-Form
   unverändert (ein Grund-Code, das `target` unterscheidet).
3. **Spezifikation:** Algorithmus (je Paar derselbe Schritt, Reihenfolge =
   Deklarationsreihenfolge, Dedup über das Befund-Tupel), §2-Schema, §4-Zeile.
4. **Code + Tests:** Config-Rand (Decode, Validierung, Kurzform ⇒ Ein-Paar-Liste),
   Modul-Schleife über die Paare; Tests für alle Akzeptanzkriterien plus
   Mutations-Gegenprobe.
5. **Default-Beweis:** Befundsatz ohne den neuen Schlüssel byte-identisch gegen
   das gepinnte Vorgänger-Image — grün wie rot.
6. **Spiegel** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
   Anforderung, Algorithmus, §2-Schema, §4-Tabelle, Klartexte, `--print-config`-
   Vorlage, Handbuch (Release-Prep, nicht hier), Config-Kommentar.
7. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Die 3×-Form wird nicht konfiguriert.** Dieser Slice macht sie baubar; das
  eigene Profil bleibt unverändert.
- **Kein Handbuch, kein README, kein CHANGELOG** — das ist Release-Prep
  (slice-125).
- **Keine Default-Änderung, kein neuer Grund-Code.**

## 4. Definition of Done

- [ ] CR-Commit (Lastenheft allein) liegt **vor** Spezifikation und Code.
- [ ] Listen-Form implementiert; Kurzform weiter gültig und getestet.
- [ ] Default-Beweis byte-identisch; fail-closed-Ränder je Paar geprüft.
- [ ] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Zwei Schreibweisen für dieselbe Sache** (Kurzform und Liste) sind eine
  Drift-Quelle — die Kurzform wird intern in eine Ein-Paar-Liste übersetzt,
  damit es genau **einen** Auswertungspfad gibt. — **Ausgang:** *(bei Closure)*
- **Dedup über mehrere Paare:** zwei Paare können dieselbe Zeile treffen; das
  Befund-Tupel muss sie unterscheidbar halten. — **Ausgang:** *(bei Closure)*
- **Byte-Identität ist die Zusage**, nicht „ungefähr gleich" — sie wird gegen
  das released Image gemessen. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): Wellen-Eröffnung (eingetreten).

**Rückführungen:** `in-progress` → `next`, falls die Listen-Form eine
Änderung am Befund-Tupel verlangt (dann ist es kein additiver Schnitt mehr).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Kern (`internal/hexagon`, GF), Config-Rand
  (`configyaml`, GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): **BEO-008 bei
  3** ist der Anlass dieses Slice; BEO-002 wirkt als Spiegel-Pflicht,
  BEO-006/009/010 als Arbeitsregeln.

Slice-ID: slice-122. Betroffene IDs:
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
`versions` (Kern `rules/`), Config-Rand, Spec. Gates: `make test` (eng),
`make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — additive Erweiterung einer eigenen,
spezifizierten Anforderung.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
