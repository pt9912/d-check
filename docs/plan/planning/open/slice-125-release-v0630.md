# Slice slice-125: Release-Prep über drei Erweiterungen und Release `v0.63.0`

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-82-config-flaechen](../welle-82-config-flaechen.md) (zugeordnet
bei der Eröffnung).

**Bezug:**
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Image-Akzeptanzkriterien), [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus);
die drei Erweiterungen aus slice-122 bis slice-124; die Release-Prozedur in
`docs/user/releasing.md`.

**Berührte Spec-Stellen:** — (Nutzer-Doku und Release-Register; die
Spec-Änderungen liegen in den drei Vorgänger-Slices).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die drei Erweiterungen erreichen den Konsumenten: die Nutzer-Doku ist auf
Stand, das Release `v0.63.0` liegt auf GHCR, und der Digest-Pin ist
zurückgeschrieben. Release-Prep ist hier **eigener** Schritt und nicht Anhang
eines Feature-Slice — die drei Feature-Commits fassen die Nutzer-Doku bewusst
nicht an.

## 2. Vorgehen

1. **Doku-Currency über alle drei Erweiterungen**, nach der §4-Checkliste der
   Release-Prozedur und den bekannten Blind-Spots: Benutzerhandbuch (§5
   Konfigurations-Referenz, §6 Modul-Tabelle, §11-Zeile **chronologisch
   unten**, Kopf-Version), `docs/user/operations.md` (Modul-Liste **und**
   Optionen-Tabelle), README **und** README.de (Modul-Listen), `version.md`
   (Verlaufs-Zeile **und** der Anker wandert auf die neue Version — genau
   dafür trägt ihn nur die aktuelle), `CHANGELOG.md`.
2. **Den Handbuch-Satz aus slice-121 korrigieren:** dort steht, `diagrams` habe
   weder Datei- noch Zeilen-Ventil. Mit slice-124 ist das überholt.
3. **Fenced Config-Beispiele prüfen** — sie sind gate-unsichtbar und können dem
   Validator widersprechen; jedes neue Beispiel gegen den Validator laufen
   lassen.
4. **`make fullbuild`**, dann Tag `v0.63.0`, Push, Release-Pipeline beobachten,
   `docker pull`-Beweis (OCI-Label, Digest, Smoke).
5. **Digest-Backfill** committen (Handbuch, Roadmap, Slice-Verweise auf den
   Digest-Pin).
6. Unabhängiger Review **vor** dem Tag; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Funktions-Änderung.** Was hier passiert, ist Doku und Auslieferung.
- **Keine Profil-Umstellung** (das `diagrams`-Scoping, die Beobachtungs-3×-Form)
  — eigene Entscheide nach dem Release.

## 4. Definition of Done

- [ ] Alle Doku-Flächen der §4-Checkliste nachgezogen; die überholte
      Ventil-Aussage im Handbuch korrigiert; Config-Beispiele gegen den
      Validator geprüft.
- [ ] `make fullbuild` grün (Exit explizit); Review vor dem Tag.
- [ ] Tag `v0.63.0` gepusht, GHCR-Pipeline grün, `docker pull`-Beweis mit
      Digest und OCI-Label; Digest-Backfill committet.
- [ ] `version.md`: Verlaufs-Zeile ergänzt und der Anker auf die neue Version
      gewandert (genau einer).

## 5. Abnahme-Punkte / Risiken

- **Die Doku-Blind-Spots sind gate-unsichtbar:** Modul-Listen, §11-Position,
  Optionen-Tabelle und fenced Config-Beispiele fängt kein Gate. Sie stehen
  deshalb auf der Checkliste und werden einzeln abgehakt. — **Ausgang:** *(bei
  Closure)*
- **Der Anker-Wechsel in `version.md`** ist neu scharf (seit slice-121): wird
  er vergessen, bleibt ein alter Pin auflösbar. — **Ausgang:** *(bei Closure)*
- **Drei Erweiterungen in einem Release** heißt drei Handbuch-Abschnitte und
  eine §11-Zeile, die alle drei nennt — nicht drei Zeilen. — **Ausgang:**
  *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): slice-122, slice-123 und slice-124 in
`done/`.

**Rückführungen:** `in-progress` → `next`, falls der Review vor dem Tag eine
Funktions-Auflage findet (dann erst der betroffene Feature-Slice).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Nutzer-Doku (GF), Release-Register (GF),
  Release-Pipeline (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): **BEO-009**
  (Botschaft behauptet ungeprüfte Probe) ist im Release-Flow einschlägig —
  Botschaften über Pipeline-Läufe erst **nach** deren Exit; BEO-006 beim
  pfad-selektiven Commit; BEO-002 als Spiegel-Pflicht.

Slice-ID: slice-125. Betroffene IDs:
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
Nutzer-Doku, Release-Register, Release-Pipeline. Gates: `make fullbuild`,
`make image-test`; Release-Pipeline.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Auslieferung nach etablierter Prozedur.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
