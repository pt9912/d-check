# Konsumenten-CR an `ai-harness-course` — Zeitdokumente-Archivierung (Modul 6 Schritt 4) hat im wellenlosen Betrieb keinen Träger

**Absender:** d-check (Adopter, Baseline-Pin `v5.18.0`)
**Datum:** 2026-09-03
**Gegenstand:** `regelwerk/modul-06-roadmap.md` §*Wann Arbeit eine Welle
braucht* (Träger-Tabelle) + `regelwerk/modul-10-review-harness.md`
§*Reviewer berichtet auch, was er nicht gefunden hat*

Ein Punkt, an eigenem Bestand gemessen. Er bittet um keinen Mechanismus —
nur um zwei ergänzende Sätze an bestehenden Stellen.

---

## Was dasteht

`modul-06` §*Wann Arbeit eine Welle braucht* listet fünf Vorgänge mit ihrem
Träger im wellenlosen Betrieb:

> | Vorgang | Träger im Repo **ohne** Wellen | Wann |
> |---|---|---|
> | **Zähler** | Slice-Closure §7 | vor dem `git mv` nach `done/` |
> | **Lese-Schritt** (was hat 3× erreicht → Ausgang zuweisen) | Slice-Closure §7 | vor dem `git mv`; Anker `seit slice-<NNN>` statt `seit welle-<NN>` |
> | **Sichtungs-Schritt** (offene Beobachtungen unter der Schwelle) | Slice-**Planung**, §8 *Vorgelagert — offene Beobachtungen sichten* | beim Anlegen jedes Slice, unabhängig vom Sub-Area-Modus |
> | **Trigger-Audit** (Carveout · Bootstrap-aware Gate · ADR) | Slice-Closure | bei jeder Closure, zusammen mit dem Lese-Schritt |
> | **Alle drei Paarungen** (a/b/c aus Closure-Schritt 3) | Slice-Closure | **nach** dem `git mv` — sie suchen in `done/` |

`modul-10` §*Reviewer berichtet auch, was er nicht gefunden hat* bindet die
Archivierung eines Review-Reports ausdrücklich an dieses Ereignis:

> **Mit der Closure der Welle, die seinen Slice einsammelt, wandert der Report
> vollständig ins Archiv — ohne Stub** ([Modul 6](../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md), Schritt
> 4). Er hat keine Identität jenseits seines Slice; wer ihn sucht, sucht ihn
> unter dem Slice, den er geprüft hat.

## Was fehlt

Die Träger-Tabelle deckt fünf Vorgänge — **Schritt 4 (Zeitdokumente
archivieren) fehlt als sechste Zeile.** Für wellenlose Arbeit, die die
Baseline selbst als legitimen Dauerbetrieb vorsieht („Wellenlos heißt nicht
wächterlos", `modul-06` Zeile 38) und nicht als Ausnahmefall, tritt das in
`modul-10` genannte Archivierungs-Ereignis nach dem Wortlaut **nie** ein — es
gibt keine Welle, die den Slice einsammelt.

## Warum das jeden Adopter trifft

Jeder Adopter, der wellenlosen Betrieb nutzt, produziert Review-Reports, die
nach der Slice-Closure keinen Konsumenten mehr haben (`modul-10`: „Lauf-Beleg,
kein Wissensspeicher [...] über Läufe hinweg wird er nicht wieder gelesen"),
aber keinen kanonisch benannten Weg, sie loszuwerden. Der Bestand wächst
unbegrenzt und ohne beschriebenes Ende.

## Beleg — gemessen, nicht behauptet

In unserem Repo: **45** bewusst wellenlose Slices (`— **wellenlos**` im
Kopf-Feld, ab `slice-137`) mit **57** zugehörigen Review-Report-Dateien unter
`docs/reviews/`, die nach dem Wortlaut niemals archiviert werden — gegenüber
**85** wellengebundenen Wellen, deren Zeitdokumente vollständig archivierbar
waren und archiviert wurden (unser `tools/archive-wave`, gebaut nach genau
dieser Prozedur).

## Der Kanon kennt die Lösungsform schon — an einer Nachbarstelle

`modul-06`, Wellen-Closure-Prozedur Schritt 4, für den strukturell verwandten
Fall von Alt-Slices ohne Wellen-Zuordnung (Bestand von vor der
Werkzeug-Einführung):

> Sie braucht dafür eine Entscheidung, die die laufende Regel nicht liefert:
> welche Welle die Slices einsammelt, die keiner angehören. Das Repo benennt
> die Zuordnung — die chronologisch nächste geschlossene Welle oder ein
> einzelnes Sammel-Archiv für den Bestand vor der Einführung.

Für wellenlose Slices scheidet die erste Option aus — eine Zuordnung zu einer
zufälligen Nachbar-Welle würde die Wellenlosigkeit rückwirkend verfälschen.
Die zweite, ein eigenständiges **Sammel-Archiv**, passt strukturell auf
denselben Bedarf, hat aber weder eine Zeile in der Träger-Tabelle noch eine
Erwähnung als zulässige Option für **diesen** Fall.

## Was der CR nicht verlangt

Kein Gate, kein vorgeschriebener Mechanismus, keine Pflicht zur Archivierung
— nur zwei ergänzende Sätze:

1. eine sechste Zeile in der Träger-Tabelle für „Zeitdokumente archivieren
   (Review-Reports wellenloser Slices)" — Träger: kein automatischer,
   Repo-Entscheidung, mit Verweis auf die Sammel-Archiv-Option aus Schritt 4;
2. in `modul-10` ein Halbsatz, der diese Sammel-Archiv-Option ausdrücklich
   auch für wellenlose Review-Reports zulässig nennt, statt sie implizit an
   die Wellen-Closure zu binden.

Welche Form (eigenes Werkzeug-Kommando, manueller Vorgang, Zeitpunkt) das
Repo dafür wählt, bleibt — wie beim Alt-Bestand-Fall — seine Entscheidung.
