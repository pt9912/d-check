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

- [x] CR-Commit (Lastenheft allein) liegt **vor** Spezifikation und Code —
      Lastenheft 0.63.0 allein, danach ADR, Spezifikation, Implementierung;
      die Review-Auflagen als eigener Nachzug (0.63.1).
- [x] Listen-Form implementiert; Kurzform weiter gültig und getestet — die
      Kurzform wird am Config-Rand in die einelementige Liste übersetzt, der
      Kern kennt nur die Liste.
- [x] Default-Beweis byte-identisch; fail-closed-Ränder je Paar geprüft —
      gegen das gepinnte Vorgänger-Image, grün wie rot, **zweimal** gemessen
      (nach der Implementierung und erneut, nachdem der Review den
      Nachrichten-Code geändert hatte).
- [x] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet — der Review war **blockierend** (ein HIGH, vier MEDIUM, drei
      LOW), alle Befunde sind eingearbeitet und im
      [Report](../../../reviews/2026-08-22-slice-122-versions-musterliste-review.md)
      belegt.

## 5. Abnahme-Punkte / Risiken

- **Zwei Schreibweisen für dieselbe Sache** (Kurzform und Liste) sind eine
  Drift-Quelle — die Kurzform wird intern in eine Ein-Paar-Liste übersetzt,
  damit es genau **einen** Auswertungspfad gibt. — **Ausgang:** gehalten, aber
  enger als geplant: die Drift saß nicht im Auswertungspfad, sondern in der
  **Erkennung** der Schreibweise. Sie fragte Werte ab statt Anwesenheit, und
  ein leer gelassener Kurzform-Schlüssel schaltete die Prüfung still auf die
  Liste um (Review F-4). Jetzt Zeiger und Anwesenheit; die Rest-Grenze (ein
  Schlüssel ohne Wert ist im YAML von einem fehlenden nicht unterscheidbar)
  ist benannt, nicht geschlossen.
- **Dedup über mehrere Paare:** zwei Paare können dieselbe Zeile treffen; das
  Befund-Tupel muss sie unterscheidbar halten. — **Ausgang:** **nicht
  gehalten** — das war der HIGH. Das Befund-Tupel kann sie *nicht*
  unterscheidbar halten: die geteilte Adresse trägt die erwartete Version
  nicht, und die Nachrunde verwarf den zweiten Befund samt seiner Erwartung.
  Die Zusage war unerfüllbar, nicht bloß unerfüllt. Geändert ist deshalb die
  Zusage: **eine Adresse, alle Erwartungen** — die Nachricht nennt jede mit
  ihrer Quelle. Der Weg über ein Adress-Feld läge in
  [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--befund) und beträfe
  jedes Modul; er steht als verworfene Alternative mit Re-Evaluierungs-Trigger
  in [ADR-0058](../../adr/0058-konfigurations-flaechen-additiv-weiten.md).
- **Byte-Identität ist die Zusage**, nicht „ungefähr gleich" — sie wird gegen
  das released Image gemessen. — **Ausgang:** gehalten, und die Messung hat
  sich bewährt: sie lief ein zweites Mal, nachdem die Review-Einarbeitung den
  Nachrichten-Code angefasst hatte. Genau dort hätte eine Ein-Paar-Meldung
  ihren Wortlaut verlieren können; der Wortlaut hängt jetzt an der **Paar-Zahl**
  statt an der Schreibweise.

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

Geliefert ist die Fläche: `versions` trägt eine Liste von
Muster-Quellen-Paaren, jedes mit eigener Quelle und eigenen Ausnahmen, und die
Kurzform **ist** die einelementige Liste. Damit ist die 3×-Form der
Beobachtung BEO-008 **baubar** — gebaut wird sie nicht, das ist ein eigener
Entscheid mit eigener Messung (§3, Wellendokument §6).

Die Lehre dieses Slice ist nicht die Erweiterung, sondern was der Review an
ihr fand: **eine Zusage, die man nicht halten kann, ist schlimmer als eine,
die man nicht gibt.** Das Akzeptanzkriterium „zwei Paare, zwei Befunde" war
nicht falsch implementiert, es war unerfüllbar — die geteilte Befund-Adresse
unterscheidet zwei Befunde an derselben Stelle nicht, und die Nachrunde
verwarf den zweiten. Wer eine Zusage über die *Anzahl* von Befunden schreibt,
muss vorher wissen, was die Adresse trägt; ich hatte sie aus der Absicht
abgeleitet statt aus dem Mechanismus.

Zwei weitere Befunde treffen dieselbe Wurzel: der eingebaute Dedup hatte keine
beobachtbare Wirkung (die Mutation überlebte die ganze Suite), und die Zusage
über die Reihenfolge wurde von einem nachgelagerten Sort überschrieben — mein
Test hatte Werte gewählt, deren Sortierung zufällig der Deklaration entsprach
und die beiden Ordnungen gar nicht trennen konnte. Beides sind Tests, die
etwas *bestätigen*, ohne es zu *messen*.

Was gehalten hat, ist die Kern-Zusage der Welle: ohne den neuen Schlüssel ist
der Befundsatz byte-identisch. Sie ist zweimal gegen das gepinnte
Vorgänger-Image gemessen — vor und nach der Review-Einarbeitung —, grün wie
rot, und die Gegenprobe zeigt, dass die Zwei-Paar-Form vorher nicht bloß
unbenutzt, sondern **nicht baubar** war (`field patterns not found`, Exit 2).

Zwei Arbeitsfehler außerhalb des Produkts gehören in dieselbe Notiz: der
Beanspruchungs-Commit trug nur den Rename, weil sein `git add` einen gerade
entfernten Pfad nannte und git deshalb den ganzen Aufruf abbrach — CI meldete
es als `target-missing`. Und der erste Reparaturversuch war ein `--amend` auf
einen bereits veröffentlichten Commit; zurückgenommen, der Nachtrag steht als
eigener Commit, die Kette darauf replayt, Push als Fast-Forward.
