# Slice slice-020: Konfiguration aus Autoritäts-Dokumenten vorschlagen (`--suggest-config`)

**Status:** in-progress.

**Welle:** welle-10-config-ableitung (Trigger: Priorisierung durch den
Auftraggeber; baut auf slice-019 auf).

**Bezug:** [`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
(Ausgabe-Format wird wiederverwendet — Gerüst auf stdout),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(**Schärfung** des Out-of-Scope — siehe §2),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(das vorgeschlagene Gerüst dekodiert über den eigenen Parser),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(liest das Repo, schreibt nie),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(gleiche Eingabe → gleicher Vorschlag),
[`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
(Korpora als Gegentest-Orakel).

**Autor:** pt9912. **Datum:** 2026-06-13.

---

## 1. Ziel

`d-check --suggest-config` macht einen Lese-Durchgang über das Repo und
gibt ein **vorgeschlagenes** `.d-check.yml`-Gerüst auf stdout aus (wie
slice-019: Werkzeug gibt aus, Aufrufer leitet um, nie Schreiben). Statt
eines rein statischen Templates füllt es die **hochkonfidenten**
Schichten aus dem Repo-Inhalt:

- **`ids`-Muster aus Autoritäts-Quellen:** Der Aufrufer benennt die
  Dokumente/Verzeichnisse, in denen Kennungen *definiert* sind
  (Lastenheft, Spezifikation, Architektur, `adr/`, conventions …).
  d-check liest die dort **definierten** IDs (Heading-/Anker-Ebene) und
  leitet je Quelle Muster + `target` ab — die **Umkehrung** der
  bestehenden `ids`-Config. Kein Prosa-Mining: was in keiner
  Autoritäts-Quelle definiert ist, taucht nicht auf.
- **Opt-in-Module nach Signal:** läuft die opt-in-Sensoren probeweise
  und schlägt die vor, die echte Treffer liefern (statt Raterei).
- **Scan-Scope** aus der Verzeichnisstruktur (konservativ).

## 2. Definition of Done

- [ ] **Lastenheft-Change-Request** [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)-**Schärfung**:
  Das Out-of-Scope „Automatisches Ermitteln der Muster aus dem
  Repo-Inhalt" gilt unverändert für die **Prüfung**; ein
  **advisory Scaffold-Modus** darf Muster aus *benannten
  Autoritäts-Quellen* ableiten (Ausgabe-only, vom Menschen
  bestätigt). Neue oder erweiterte CLI-Anforderung für
  `--suggest-config` (entweder [`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
  erweitert oder eine eigene neue Anforderung); drei AKs + Out-of-Scope.
- [ ] **Spezifikation:** Ableitungs-Algorithmus (ID-Extraktion aus
  Autoritäts-Headings, Regex-Generalisierung, Signal-Probe der
  opt-in-Module), Determinismus-Festlegung.
- [ ] **Implementierung** im CLI-/Core-Layer; Ausgabe baut auf dem
  slice-019-Gerüst-Format auf.
- [ ] **Regex-Politik:** generalisierter Best-Guess **plus** die
  Quell-IDs als Kommentar („abgeleitet aus: …") — Scaffold, kein
  Orakel; der Mensch verengt.
- [ ] **Kalibrierung/Gegentest über die Korpora** ([`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Muster):
  Ableitung gegen die *handgepflegten* `ids.patterns` von d-check (3),
  u-boot (4), b-trace (10) — matchen die abgeleiteten Muster dieselbe
  ID-Menge? b-trace ist der Härtetest der Generalisierung.
- [ ] **Read-only-Beleg** ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)):
  liest das Repo, schreibt nie (read-only-Mount genügt).
- [ ] **Doku** unter `docs/user/`; `make gates` grün; Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) | update | [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)-Schärfung + Anforderung für `--suggest-config` |
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | Ableitungs-Algorithmus |
| `internal/…` (CLI + Ableitungs-Logik) | update/neu | ID-Extraktion, Regex-Generalisierung, Modul-Signal-Probe |
| `docs/user/operations.md` | update | Option dokumentieren |

## 4. Trigger

Priorisierung durch den Auftraggeber. Voraussetzung erfüllt: slice-019
liefert das Ausgabe-Format (`--print-config`).

## 5. Closure-Trigger

DoD vollständig inkl. Korpora-Gegentest (abgeleitete vs. handgepflegte
Muster dokumentiert) und grüner Gates.

## 6. Risiken und offene Punkte

- **Regex-Generalisierung** ist der harte Kern: aus einer ID-Menge ein
  treffsicheres Muster ohne Über-/Unter-Matching. b-trace (10 komplexe
  Muster) zeigt die Grenze — der Korpora-Gegentest quantifiziert sie,
  statt sie zu raten.
- **Vertrags-Disziplin:** Die Ableitung darf **nur** im Scaffold-Modus
  passieren, nie in der Prüfung — sonst kippt der Determinismus-/
  Konfigurations-Vertrag der `ids`-Prüfung.
- **Modul-Signal-Probe** bedeutet einen vollständigen Lese-Lauf — Laufzeit
  beachten (aber read-only, netzlos).

## 7. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
