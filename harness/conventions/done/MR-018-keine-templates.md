# MR-018 — d-check verkörpert als Producer-/Self-Hoster keine Templates

- **Status:** Accepted
- **Datum:** 2026-06-25
- **Geltungsbereich:** [`AGENTS.md`](../../../AGENTS.md) §1,
  [§Adoptierte Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen), der
  `.harness/cache/<tag>/templates/`-Cache (Schärfung von
  [`MR-017`](../../conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)),
  die Autoren-Quelle wiederkehrender Artefakte
- **Adaption:** d-check verkörpert **keine** co-located `*.template.md` und
  weicht damit bewusst von der Baseline-Regel **§Ein- vs. wiederkehrende
  Templates** ab (`lab/templates/README.md`: die wiederkehrenden Skelette — ADR,
  Slice, Welle, Carveout, Review-Report — bleiben co-located, jede neue Instanz
  wird daneben kopiert; die Singletons werden beim Bootstrap einmal gefüllt).
  **d-check ist Producer-/Self-Hoster** der Harness-Werkzeuge: es autoriert seine
  wiederkehrenden Artefakte **nativ im Haus-Stil**
  ([`MR-014`](../../conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template))
  aus dem gelebten Bestand seiner ADRs und Slices, nicht aus einem co-located
  Baseline-Skelett. Der `lab-templates.zip`-Cache
  (`.harness/cache/<tag>/templates/`) ist daher **Adoptions-/Drift-Audit-Staging,
  nicht Autorenquelle** — niemand autoriert aus dem ephemeren, gitignorierten
  Cache.
- **Begründung:** Sichtbar im Schwester-Repo-Vergleich (Nutzer-Frage
  2026-06-25): das Consumer/Adopter-Repo **bedrock-eu-guard** verkörpert die fünf
  wiederkehrenden Skelette co-located (folgt §Ein- vs. wiederkehrende); d-check
  als Producer tut es nicht. Der **Kurs stützt die Producer-Lesart für
  `harness.mk` bereits live** (Stand > `v1.4.0`): seine `lab/templates/README.md`
  §Self-Hosting-/Producer-Fall nennt „das Tool-Repo selbst (d-check), das seinen
  Doku-Gate via `make doc-check` direkt dogfooded" und nimmt es von der
  `harness.mk`-Adoption aus. d-checks **gepinnte v1.4.0-Baseline trägt diesen
  Abschnitt noch nicht** (er ist post-v1.4.0); diese MR **überbrückt** bis zum
  nächsten Baseline-Bump und zieht die Producer-Logik zugleich auf die
  wiederkehrenden **Dokument**-Skelette weiter — den Schritt, den der Kurs auch
  live (noch) nicht ausspricht. Ohne die Deklaration bliebe d-checks
  Template-Freiheit eine **stille Setzung**, und ein Agent könnte fälschlich aus
  dem Audit-Cache autoren.
- **Auflösungs-Trigger:** Re-Evaluation beim nächsten Baseline-Bump (Pin auf den
  post-v1.4.0-Kurs, der den Self-Hosting-/Producer-Fall in der Templates-README
  bereits kanonisiert): die MR wird gegen den dann aktuellen Kanon geprüft und
  zur reinen Provenienz — oder aufgelöst, falls der Kurs die Producer-Lesart auch
  für die wiederkehrenden Dokument-Skelette übernimmt.
