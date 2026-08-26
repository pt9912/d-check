# MR-032 — Versions-Bump und Historie-Zeile schon vor `Accepted`

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`grundlagen-source-precedence.md`](../../.harness/baseline/v5.12.0/regelwerk/grundlagen-source-precedence.md)
  §Source Precedence, Absatz *Wann die CR-Pflicht beginnt* — *„Vor `Accepted`
  ist das Lastenheft ein Entwurf — frei änderbar, ohne Change Request, **ohne
  Historie-Zeile**."*
- **Datum:** 2026-08-23
- **Geltungsbereich:** [`spec/lastenheft.md`](../../spec/lastenheft.md) —
  jede Änderung am Dokument, solange sein `**Status:**` unter `Accepted` liegt.
- **Adaption:** Dieses Repo führt **Versions-Bump und Historie-Zeile ab der
  ersten Fassung**, nicht erst ab `Accepted`. Der Kanon erlaubt beides zu
  lassen; wir lassen die Freiheit ungenutzt.

  **Warum das ein Eintrag ist und keine bloße Gewohnheit:** Der Kanon weist
  genau diese Entscheidung dem Konventionsspeicher zu — *„Welche Stelle der
  Version steigt, entscheidet das Repo und gehört in den Adaptions-Block von
  `harness/conventions.md`."* Und der Freshness-Audit liest **diese Liste**;
  ohne Eintrag wäre die Abweichung für ihn unsichtbar. Die zunächst gewählte
  Abgrenzung — *„eine Freiheit nicht nutzen ist keine Abweichung"* —
  widerspricht zudem
  [`MR-031`](../conventions.md#mr-031), der für denselben Fall festhält: *wer
  verschärft, weicht ab, auch wenn er nur mehr verlangt.*

  **Begründung der Strenge:** Die Historie ist der einzige Ort, an dem eine
  Anforderungs-Änderung ihren Beleg hinterlässt. Sie erst ab `Accepted` zu
  führen hieße, den Bestand bis dahin spurlos wachsen zu lassen — bei
  inzwischen 48 Anforderungen und 95 Historie-Zeilen wäre das der Verlust der
  gesamten Entstehungs-Spur.

  **Was die Adaption *nicht* vorwegnimmt:** die `Verweis`-Spalte bleibt
  durchgehend `—`, weil ohne begonnene CR-Pflicht kein externer Vorgang
  existiert, den sie nennen könnte. Und der **Status-Wechsel** selbst ist
  keine Folge dieses Eintrags, sondern eine Auftraggeber-Entscheidung mit
  Vertragswirkung.
- **Ausgelöst durch Baseline-Stand:** v5.11.0
- **Auflösungs-Trigger:** der `**Status:**` des Lastenhefts erreicht
  `Accepted`. Dann verlangt der Kanon Bump und Historie ohnehin, die Adaption
  ist aufgelöst — und die `Verweis`-Spalte beginnt, echte Vorgänge zu tragen.
