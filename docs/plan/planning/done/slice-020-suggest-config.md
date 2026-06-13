# Slice slice-020: Konfiguration aus Autoritäts-Dokumenten vorschlagen (`--suggest-config`)

**Status:** done.

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

- [x] **Lastenheft-Change-Request** [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)-**Schärfung**:
  Das Out-of-Scope „Automatisches Ermitteln der Muster aus dem
  Repo-Inhalt" gilt unverändert für die **Prüfung**; ein
  **advisory Scaffold-Modus** darf Muster aus *benannten
  Autoritäts-Quellen* ableiten (Ausgabe-only, vom Menschen
  bestätigt). Neue oder erweiterte CLI-Anforderung für
  `--suggest-config` (entweder [`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
  erweitert oder eine eigene neue Anforderung); drei AKs + Out-of-Scope.
- [x] **Spezifikation:** Ableitungs-Algorithmus (ID-Extraktion aus
  Autoritäts-Headings, Regex-Generalisierung, Signal-Probe der
  opt-in-Module), Determinismus-Festlegung.
- [x] **Implementierung** im CLI-/Core-Layer; Ausgabe baut auf dem
  slice-019-Gerüst-Format auf.
- [x] **Regex-Politik:** generalisierter Best-Guess **plus** die
  Quell-IDs als Kommentar („abgeleitet aus: …") — Scaffold, kein
  Orakel; der Mensch verengt.
- [x] **Kalibrierung/Gegentest über die Korpora** ([`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Muster):
  Ableitung gegen die *handgepflegten* `ids.patterns` von d-check (3),
  u-boot (4), b-trace (10) — matchen die abgeleiteten Muster dieselbe
  ID-Menge? b-trace ist der Härtetest der Generalisierung.
- [x] **Read-only-Beleg** ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)):
  liest das Repo, schreibt nie (read-only-Mount genügt).
- [x] **Doku** unter `docs/user/`; `make gates` grün; Closure-Notiz.

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

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Vertrag ([`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten), Lastenheft 0.10.0) +
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)-Schärfung
+ Spezifikation [`DC-FA-CLI-006.a`](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag); `core.SuggestConfig` (ID-Extraktion
aus Headings, Präfix-Alternation mit Round-Trip-Garantie, Modul-Signal-Probe,
YAML-Render), CLI-Flag `--suggest-config`. `make gates` grün.

**Korpora-Gegentest (das Orakel):** `--suggest-config` mit den
Autoritäts-Quellen je Repo gegen die *handgepflegten* Muster:

| Repo | Befund |
|---|---|
| d-check | alle 3 Familien round-trippen: `ADR-\d+`/`MR-\d+` ≈ handgepflegt; `DC` als explizite Präfix-Alternation (handgepflegt generalisiert `FA-[A-Z]+`) |
| u-boot | `LH`-Familie round-trippt (verbose, 32 Präfixe); `slice`/`tranche` + `PH/TC/CO` **verfehlt** |
| b-trace | nur `MR` erkannt — die Requirement-IDs stehen dort nicht als Heading-Token |

- **Was hat funktioniert:** Die Round-Trip-Invariante (abgeleiteter
  `regex` matcht jede Quell-Kennung) macht den Vorschlag *nie falsch*,
  nur evtl. zu weit — und die Quell-Kennungs-Kommentare machen die
  Ableitung prüfbar. Das Orakel war geschenkt: die Korpora *haben*
  schon handgepflegte Muster.
- **Anders/ehrlicher als erhofft:** Die Heading-basierte Extraktion hat
  eine klare, dokumentierte Grenze (großgeschriebene Heading-Token;
  keine `slice-001`/Tabellen-IDs). b-trace war exakt der Härtetest, der
  sie sichtbar machte. Statt das wegzukaschieren, steht die Grenze in
  [`docs/user/operations.md`](../../../../docs/user/operations.md) und
  im Out-of-Scope — „Scaffold, kein Orakel" ist Vertrag, nicht Ausrede.
- **Lerneintrag:** Eine Round-Trip-Garantie + Quell-Nachweis verwandelt
  eine riskante Heuristik (Muster-Raten) in ein vertretbares
  Komfort-Werkzeug — der Mensch verengt, wird aber nie in die Irre
  geführt.
- **Review R1** (`/code-review`, high — [Report](../../../../docs/reviews/2026-06-13-slice-020-suggest-config.md)):
  fünf echte Defekte gefixt, am schwersten ein von den Gates *nicht*
  fangbarer: das Gerüst emittierte `ids`-Muster, ließ aber `ids` aus der
  Modul-Liste — gültiges YAML, semantisch wirkungslos. Genau die Klasse
  „grün ≠ richtig", die den Review rechtfertigt. Dazu Probe-Scope-Fix,
  `target`-Quoting, robustere Heading-Token-Extraktion, leere-Quellen-Fehler.
  Fünf Edge-/Cleanup-Punkte bewusst akzeptiert (im Report dokumentiert).
- **Folge-Slices:** keine; eine Generalisierungs-Verbesserung (Stamm-
  Kollabieren wie `DC-FA-[A-Z]+`) oder kleingeschriebene IDs wären
  eigene, klar abgegrenzte Folge-Slices, falls Bedarf entsteht.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
