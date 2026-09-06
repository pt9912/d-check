# Beobachtung — `BEO-<KUERZEL>/<slug>`

> **Template-Hinweis.** Vorlage für **eine** Beobachtung des
> Beobachtungs-Registers. Lege sie an unter
> `docs/plan/planning/observations/BEO-<KUERZEL>/<slug>/` und ersetze
> Platzhalter. Lösche diesen Block.
>
> Das `<KUERZEL>` wird **nachgeschlagen**, nicht formuliert: Es ist das
> Sub-Area-Kürzel aus der Modus-Deklaration in `harness/conventions.md`.

Regeln dieser Ablage: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — wer schreibt, wer liest, wann gestrichen wird,
welche Form ein Beleg hat und welchen der drei Ausgänge ein Eintrag ab 3×
trägt.

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten.

DREI DATEIEN, DREI LEBENSDAUERN. Das ist der ganze Trick:

  observation.md          unveraenderlich ab Anlage
  state.md                veraenderlich
  evidence/<vorgangs-id>.md   unveraenderlich ab Merge, eine je Auftreten

Es gibt KEIN Zaehler-Feld. Der Zaehler ist die Zahl der Evidence-Dateien.
Wer ihn irgendwo hinschreibt, baut die zweite Quelle wieder ein, die diese
Ablage gerade abgeschafft hat.

Erfinde keine Belege: Ein Verzeichnis entsteht beim ERSTauftreten einer echten
Beobachtung, nicht beim Adoptieren dieser Vorlage.

Eine Evidence-Datei traegt den Namen eines ABGESCHLOSSENEN VORGANGS — Regelfall
Slice, auch Welle oder Review-Report. Ein Vorgang zaehlt EINMAL; das erzwingt
hier das Dateisystem, nicht die Disziplin. Ein Vorkommen ohne abgeschlossenen
Vorgang bekommt keine Datei und bewegt den Zaehler nicht — es gehoert trotzdem
in observation.md unter "Benannt, nicht gezaehlt".

Slug: lowercase ASCII-Kebab-Case, keine fuehrenden, abschliessenden oder
doppelten Bindestriche. Namespace und Slug sind Herkunftsanker und werden nach
dem ersten Merge NICHT umbenannt.
-->

## `observation.md` — die Identität

```markdown
# <kurze, gleichbleibende Bezeichnung>

**Sub-Area:** <Name der Sub-Area, deren Konventions-Härte oder Inventur-Linie
die Beobachtung betrifft — nicht die, in deren Verzeichnis sie auffiel>

<Ein bis drei Sätze: was genau beobachtet wurde.>

## Benannt, nicht gezählt

<Vorkommen ohne abgeschlossenen Vorgang — beim Lesen von Code, im Gespräch,
bei einer Freigabe. Sie bewegen den Zähler nicht. Weglassen, wenn keine.>
```

## `state.md` — der Stand

```markdown
**Stand:** offen
```

Unterhalb der Schwelle ist `offen` der Normalzustand, kein Ausgang. Ab **3×**
trägt `state.md` genau einen von drei Ausgängen — eine geschlossene Menge,
kein Freitext:

| Ausgang | Was dazugehört |
|---|---|
| `verkörpert` | Zielort **und** Herkunfts-Anker (`seit welle-<NN>` / `seit slice-<NNN>`) |
| `geplant` | Kennung des Slice oder der Welle, die die Regel schreibt |
| `gestrichen` | die Begründung, warum die Beobachtung nicht mehr auftreten kann |

`gestrichen` hängt als einziges nicht an der Schwelle: Fällt die Ursache
vorher weg, wandert der Stand sofort dorthin.

**Gestrichen heißt nicht gelöscht.** Das Verzeichnis bleibt liegen — ein still
entferntes ist nicht von einem zu unterscheiden, das es nie gab.

## `evidence/<vorgangs-id>.md` — je Auftreten eine Datei

```markdown
**Vorgang:** <slice-NNN | welle-NN | review-YYYY-MM-DD-<kurz>>
**Fund:** <ein Satz, was in genau diesem Vorgang auftrat>
```

Der Dateiname **ist** die Kennung: `slice-042.md`, nicht `beleg-3.md`.
