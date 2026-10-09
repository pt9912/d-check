# Mit `skip-allows-empty` wird ein Umzug oder ein zu breites Muster still

**Sub-Area:** `*`

Ein Repo, das `skip-allows-empty: true` setzt, erklärt eine Menge, die erst
`skip-pattern` leert, zum Ruhezustand. Trifft das Muster auch Volltexte, oder
liegen die Volltexte nach einem Umzug außerhalb der Kandidatenmenge, meldet
der Lauf nichts mehr. Die Grenze ist in [ADR-0106](../../../../adr/0106-skip-allows-empty-erklaert-den-ruhezustand.md) benannt; beobachtet wird, ob
sie bei einem Adopter eintritt — dann ist das der Re-Evaluierungs-Trigger der
ADR.
