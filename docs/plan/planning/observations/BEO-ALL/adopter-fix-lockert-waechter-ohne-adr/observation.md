# Der Fix für einen Adopter-Befund lockert einen Wächter, ohne dass eine ADR es trägt

**Sub-Area:** `*`

Ein Adopter meldet, dass ein Gate in einem legitimen Zustand rot ist. Der
naheliegende Fix nimmt den Zustand aus der Prüfung heraus — und lockert damit
eine Regel, die eine akzeptierte ADR ohne Ausnahme festlegt. Der Auftrag
scheint die Lockerung zu decken, `AGENTS.md` §3.6 verlangt trotzdem eine ADR,
und die stille Form schaltet oft einen zweiten Fall ab, den die ADR fangen
sollte. Ableiter: vor dem Fix die ADR suchen, die den Wächter trägt, und die
Ausnahme deklarieren statt still annehmen.
