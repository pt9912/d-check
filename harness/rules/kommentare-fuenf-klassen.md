# Kommentare tragen eine der fünf Klassen — Zustandsfelder, Durchsetzung, Bestandsgrenze

Ausgelagerter Körper zu `AGENTS.md` §3.7. Der operative Kern (die fünf
Klassen, das Verbot von Review-Historie/Slice-Nummern/Mess-Labels) steht
dort; diese Datei trägt die Zustandsfeld-Sonderform, die Durchsetzungslage
und die Bestandsgrenze — Auslagerungs-Muster aus
[ADR-0096](../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

## Zustandsfelder tragen eine eigene Form, nicht die fünf Klassen

**Zustandsfelder** sind Zustands-Artefakte wie der Kommentar, nur im Rumpf —
sie tragen **nicht** dessen fünf Klassen, sondern eine **eigene Form**;
übertragen sind die **zwei Tests**: Adressat ist, wer den Zustand liest, um
zu handeln, und die Zeitform ist der Indikativ über das, was ist. Ein Feld,
das einen Zustand trägt — etwa eine `Stand`- oder `Status`-Zelle in
Roadmap, Beobachtungs-Register oder Meilenstein-Tabelle —, nennt **den
Zustand und den Beleg als auflösbaren Anker**, nicht die Chronik, wie es
dazu kam; das Drift-Log der Roadmap trägt **nur Umplanungen**, keine
Schließungen und keine erreichten Meilensteine (die stehen im Closure-Log
bzw. in der Status-Spalte). Ein lebendes Register trägt **keine** Kopfzeile
`Status: Aktiv. Letzte Änderung: <Datum>` — sein Zustand ist sein Inhalt,
sein Änderungsdatum hält `git`; ein Datum, das ein **benannter Trigger**
pflegt, ist davon ausgenommen (der Frische-Marker der Architektur-Sicht).

**Verhältnis zu §3.5:** das `**Status:**`-Feld einer ADR ist ein
Zustandsfeld wie jedes andere — `adr-check` nimmt die Kopf-Status-Zeile
ausdrücklich **aus** dem Kern-Vergleich und lässt den Übergang zu. Es darf
also korrigiert werden, solange der Wert die erlaubte Form behält; §3.5
schützt den **Kern**, nicht dieses Feld. **Benannte Bestands-Ausnahme:** die
historischen `**Status:**`-Felder der `done/`-Slices bleiben, wie sie sind
— sie sind eingefrorene Lauf-Belege, ihr Lifecycle-Zustand ist ohnehin das
Verzeichnis, und das Feld hat dort keine Funktion
([`AGENTS.md` §5](../../AGENTS.md#5-dokumentations-regeln), Regel 11).
Gemeldet wird von
ihnen nur, was dem Verzeichnis **widerspricht**.

Kanon:
[Baseline §Was ein Kommentar trägt](../../.harness/baseline/v6.9.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte).

## Durchsetzung

**Kein Gate prüft das** — weder die fünf Klassen noch die Zustandsfeld-Form;
die Prüfung ist ein Urteil, kein `grep`. Der Reviewer-Skill trägt dazu
**zwei** HIGH-Anker.

## Bestandsgrenze

Vor der Einführung bzw. vor dieser Schärfung geschriebene Kommentare
(Test-Kommentare; ältere Config-Kommentare mit Slice-Nummer) sind
grandfathered — geräumt wird beim nächsten Anfassen der Zeile;
Neuzugänge fallen überall unter den Anker. **Für Zustandsfelder gibt es
keine Bestandsgrenze:** auch Altbestand fällt unter die Form — bis auf die
oben benannte Ausnahme wird nichts grandfathered. *(Auflösungs-Trigger:
permanent.)*
