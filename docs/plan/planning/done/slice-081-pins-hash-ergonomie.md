# Slice slice-081: `pins`/dpin-Ergonomie — vollen Ist-Hash im `link-stale`-Befund

**Status:** In Arbeit
**Welle:** welle-64-dpin-ergonomie (Trigger: WIP-Slot frei nach welle-63; Nutzer-Entscheid 2026-07-19 — eigener Slice statt Quer-Schnitt in slice-072)
**Bezug:** [`DC-FA-PIN-001`](../../../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in) (Verhalten unverändert — nur die **nicht stabilitätsgarantierte** Befund-`message`); **kein** ADR, **kein** Lastenheft/Spec-Vertragsdelta. Herkunft: die slice-079/080-Erkenntnis, dass `pins`/dpin wie ausgeliefert unbenutzbar ist (der `link-stale`-Befund zeigt nur `shortHash`, kein Weg an den vollen Pin-Hash).
**Autor:** pt9912
**Datum:** 2026-07-19

---

## 1. Ziel

Das Modul `pins` (Content-Pin `<!-- dpin: sha256:… -->`) ist wie ausgeliefert praktisch **nicht adoptierbar**: der `link-stale`-Befund meldet den errechneten Ziel-Span-Hash nur als `shortHash` (12 Stellen + „…"), und es gibt keinen anderen Weg an den vollen 64-Hex-Hash, den der Marker braucht. Der Slice schließt die Ergonomie-Lücke: der Befund führt den **vollen** errechneten `sha256` als Re-Pin-Vorlage — „einmal laufen, gemeldeten Hash in den Marker kopieren", genau wie das Modul `sources` es nativ tut. Plus eine Handbuch-Aufgaben-Sektion, die diesen Workflow beschreibt.

## 2. Entscheidungen

- **Nur die `message`, kein Vertrag.** Die Befund-`message` ist laut Spezifikation „nicht stabilitätsgarantiert" → kein Lastenheft/Spec-Delta, kein ADR. Grund-Code (`link-stale`), Exit-Code und Semantik bleiben unverändert.
- **Voller `errechnet`-Hash, kurzer `erwartet`.** Emittiert wird der **errechnete** (Ist-)Hash voll — den kopiert der Nutzer in den Marker; der `erwartet`-Teil (der bestehende Pin) bleibt gekürzt (Kontext, kein Kopier-Wert). `shortHash` bleibt in Gebrauch.
- **Eigener Slice statt Quer-Schnitt.** Bewusst NICHT in slice-072 (reine §4-Doku, kein Release) gefaltet (Nutzer-Entscheid 2026-07-19).

## 3. Definition of Done

- [ ] **Code:** `pins.go` `link-stale`-`message` emittiert den **vollen** errechneten `sha256` statt `shortHash(got)`; `shortHash(want)` bleibt für den `erwartet`-Teil.
- [ ] **Test:** ein Test belegt, dass die `message` den **vollen** 64-Hex-Ist-Hash trägt (mutations-echt: ohne den Fix trägt sie nur den gekürzten Hash).
- [ ] **Handbuch:** neue Aufgaben-Sektion „Einen Link-Inhalt gegen Drift pinnen (Modul `pins`)" (Ausgangslage → Ziel → Vorgehen: `--enable pins` laufen → vollen errechnet-Hash aus dem `link-stale`-Befund kopieren → in den `<!-- dpin: sha256:… -->`-Marker setzen → Ergebnis); §11-Zeile + Handbuch-Version; `CHANGELOG.md`.
- [ ] **Belege:** `make ci`/`make gates` grün; Review (leichtgewichtig); Release **v0.51.1** + Digest-Backfill; Closure.

## 4. Risiken / offene Punkte

- **Kein Vertrags-Gate** — die `message`-Verbesserung ist nicht gate-erzwungen (message = nicht stabilitätsgarantiert); der Test pinnt sie dennoch.
- **SemVer:** Ausgabe-Verbesserung ohne Verhaltens-/Grund-Code-Änderung → **Patch v0.51.1**. (Dass dpin dadurch erst praktisch benutzbar wird, ist ein Capability-Unlock, rechtfertigt aber keinen Minor, da keine neue Findungs-/Modul-Fähigkeit.)

## 5. Trigger

Nutzer-Entscheid 2026-07-19: der dpin-Ergonomie-Retrofit — in slice-079/080 als „bewusst separat" ausgewiesen — wird als **eigener** Kleinst-Slice gefahren, nicht in slice-072 quergeschnitten.

## 6. Sub-Area-Modus-Begründung

GF: kleiner Produkt-Fix (Befund-`message`) + abgeleitete Nutzer-Doku; kein Vertragsdelta.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend._
