# Test-Abdeckung je Anforderung (Go-Suite)

Abgeleitet aus den Go-Tests unter `internal/` und `cmd/`: eine Zeile je
Testfunktion, deren Doc-Kommentar unmittelbar über `func Test…` eine
Anforderungs-Kennung nennt. Gezählt wird, was `go test` als Test ausführt;
eine Kennung an anderer Stelle des Tests zählt nicht. Die Datei schreibt
`make abdeckung`; `make test` hält sie gegen die Testquellen.

**Grenze:** eine Deklaration, kein Beleg — die Zeile sagt, dass der Test die
Anforderung prüfen soll, nicht, dass er es tut.

| Kennung | Test | Datei |
| --- | --- | --- |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestConfigYAMLDiagramsExemptPaths` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestConfigYAMLDiagramsExemptPathsSegmentweise` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md) | `TestConfigYAMLStructureHeadingPattern` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestConfigYAMLStructureMusterfehlerDeterministisch` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestConfigYAMLVersionsFailClosed` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestConfigYAMLVersionsKurzformIstEinPaarListe` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestConfigYAMLVersionsZweiPaare` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-PLAN-001`](../../spec/lastenheft.md) | `TestDecode_ClosureFehler` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-COMMITS-001`](../../spec/lastenheft.md) | `TestDecode_CommitsFehler` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md), [`DC-FA-HOST-001`](../../spec/lastenheft.md) | `TestDecode_HostpathsPrefixesFehler` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestDecode_LinkPolicy` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-MTX-003`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestDecode_MatrixAllowIfSameIDDefaultAus` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-MTX-003`](../../spec/lastenheft.md) | `TestDecode_MatrixAllowIfSameIDFailClosed` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-MTX-003`](../../spec/lastenheft.md) | `TestDecode_MatrixAllowIfSameIDHappy` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-MTX-002`](../../spec/lastenheft.md) | `TestDecode_MatrixDirectionFailClosed` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-MTX-002`](../../spec/lastenheft.md) | `TestDecode_MatrixOrderDirection` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-MTX-001`](../../spec/lastenheft.md) | `TestDecode_MatrixSupersedeLineage` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-MTX-003`](../../spec/lastenheft.md) | `TestDecode_MatrixTokenFailClosed` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-MTX-003`](../../spec/lastenheft.md) | `TestDecode_MatrixTokenUndExemptPaths` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-CONF-002`](../../spec/lastenheft.md) | `TestDecode_ModulScope` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-PLAN-001`](../../spec/lastenheft.md) | `TestDecode_PlanningFehler` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-LINK-001`](../../spec/lastenheft.md) | `TestDecode_ResolveFromFehler` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestDecode_SourcesFailClosed` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestDecode_SourcesHappy` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestDecode_Versions` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestDecode_VersionsFehler` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-CONF-002`](../../spec/lastenheft.md), [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestDecode_VersionsScope` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-PLAN-001`](../../spec/lastenheft.md) | `TestDecode_WavesFehler` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-PLAN-001`](../../spec/lastenheft.md) | `TestDecode_WavesMode` | [`internal/adapter/driven/configyaml/configyaml_test.go`](../../internal/adapter/driven/configyaml/configyaml_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDecodeDiagrams` | [`internal/adapter/driven/configyaml/diagrams_test.go`](../../internal/adapter/driven/configyaml/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDecodeDiagramsErrors` | [`internal/adapter/driven/configyaml/diagrams_test.go`](../../internal/adapter/driven/configyaml/diagrams_test.go) |
| [`DC-FA-CONF-002`](../../spec/lastenheft.md) | `TestDecodeDiagramsScope` | [`internal/adapter/driven/configyaml/diagrams_test.go`](../../internal/adapter/driven/configyaml/diagrams_test.go) |
| [`DC-FA-CLI-012`](../../spec/lastenheft.md) | `TestQA03_ClosureProfil_KeineZweiteNetzTuer` | [`internal/adapter/driven/configyaml/gate_consistency_test.go`](../../internal/adapter/driven/configyaml/gate_consistency_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestDecode_MentionsGlobValidierung` | [`internal/adapter/driven/configyaml/mentions_test.go`](../../internal/adapter/driven/configyaml/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestDecode_MentionsHalberBlockFailClosed` | [`internal/adapter/driven/configyaml/mentions_test.go`](../../internal/adapter/driven/configyaml/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestDecode_MentionsUnbekannterMatch` | [`internal/adapter/driven/configyaml/mentions_test.go`](../../internal/adapter/driven/configyaml/mentions_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md), [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestDecode_MentionsVollblock` | [`internal/adapter/driven/configyaml/mentions_test.go`](../../internal/adapter/driven/configyaml/mentions_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestDecode_CrossConsistencyFehler` | [`internal/adapter/driven/configyaml/trace_cross_test.go`](../../internal/adapter/driven/configyaml/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestDecode_CrossConsistencyHappy` | [`internal/adapter/driven/configyaml/trace_cross_test.go`](../../internal/adapter/driven/configyaml/trace_cross_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestDecode_CrossConsistencyOhneBlock` | [`internal/adapter/driven/configyaml/trace_cross_test.go`](../../internal/adapter/driven/configyaml/trace_cross_test.go) |
| [`DC-FA-VCS-002`](../../spec/lastenheft.md) | `TestAllPathsLeereRange` | [`internal/adapter/driven/git/git_test.go`](../../internal/adapter/driven/git/git_test.go) |
| [`DC-FA-VCS-001`](../../spec/lastenheft.md) | `TestAllPathsPureRenameYieldsDelete` | [`internal/adapter/driven/git/git_test.go`](../../internal/adapter/driven/git/git_test.go) |
| [`DC-FA-COMMITS-001`](../../spec/lastenheft.md) | `TestCommitMessages` | [`internal/adapter/driven/git/git_test.go`](../../internal/adapter/driven/git/git_test.go) |
| [`DC-FA-VCS-001`](../../spec/lastenheft.md) | `TestFileAtEintragOhneBlob` | [`internal/adapter/driven/git/git_test.go`](../../internal/adapter/driven/git/git_test.go) |
| [`DC-FA-VCS-001`](../../spec/lastenheft.md) | `TestFileAtUnlesbaresObjekt` | [`internal/adapter/driven/git/git_test.go`](../../internal/adapter/driven/git/git_test.go) |
| [`DC-FA-TRK-001`](../../spec/lastenheft.md) | `TestTrackedPaths` | [`internal/adapter/driven/git/git_tracked_test.go`](../../internal/adapter/driven/git/git_tracked_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestFetch_BodyLimitUeberschritten` | [`internal/adapter/driven/httpcheck/export_test.go`](../../internal/adapter/driven/httpcheck/export_test.go) |
| [`DC-FA-EXT-001`](../../spec/lastenheft.md) | `TestCheck_Erreichbar` | [`internal/adapter/driven/httpcheck/httpcheck_test.go`](../../internal/adapter/driven/httpcheck/httpcheck_test.go) |
| [`DC-FA-EXT-001`](../../spec/lastenheft.md) | `TestCheck_HeadFallbackAufGet` | [`internal/adapter/driven/httpcheck/httpcheck_test.go`](../../internal/adapter/driven/httpcheck/httpcheck_test.go) |
| [`DC-FA-EXT-001`](../../spec/lastenheft.md) | `TestCheck_RedirectKette` | [`internal/adapter/driven/httpcheck/httpcheck_test.go`](../../internal/adapter/driven/httpcheck/httpcheck_test.go) |
| [`DC-FA-EXT-001`](../../spec/lastenheft.md) | `TestCheck_Status404` | [`internal/adapter/driven/httpcheck/httpcheck_test.go`](../../internal/adapter/driven/httpcheck/httpcheck_test.go) |
| [`DC-FA-EXT-001`](../../spec/lastenheft.md) | `TestCheck_Timeout` | [`internal/adapter/driven/httpcheck/httpcheck_test.go`](../../internal/adapter/driven/httpcheck/httpcheck_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestFetch_ErreichbarMitBody` | [`internal/adapter/driven/httpcheck/httpcheck_test.go`](../../internal/adapter/driven/httpcheck/httpcheck_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestFetch_Status404` | [`internal/adapter/driven/httpcheck/httpcheck_test.go`](../../internal/adapter/driven/httpcheck/httpcheck_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestJSONYAML_MitNoteImSummary` | [`internal/adapter/driven/report/report_test.go`](../../internal/adapter/driven/report/report_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestJSONYAML_OhneNoteKeinNotesSchluessel` | [`internal/adapter/driven/report/report_test.go`](../../internal/adapter/driven/report/report_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md) | `TestText_ErlaeuterungBleibtEinFeld` | [`internal/adapter/driven/report/report_test.go`](../../internal/adapter/driven/report/report_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md), [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestText_NotesAufStderrVorDerZaehlzeile` | [`internal/adapter/driven/report/report_test.go`](../../internal/adapter/driven/report/report_test.go) |
| [`DC-FA-CLI-001`](../../spec/lastenheft.md) | `TestCLI001_Boundary_KeineMarkdownDateien` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-001`](../../spec/lastenheft.md) | `TestCLI001_Happy` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-001`](../../spec/lastenheft.md) | `TestCLI001_Negative_FehlendeWurzel` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-002`](../../spec/lastenheft.md) | `TestCLI002_CLIVorConfig` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-002`](../../spec/lastenheft.md) | `TestCLI002_DisableLinks` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-002`](../../spec/lastenheft.md) | `TestCLI002_UnbekanntesModul` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-003`](../../spec/lastenheft.md) | `TestCLI003_EinBefund` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-003`](../../spec/lastenheft.md) | `TestCLI003_UngueltigeOption` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md) | `TestCLI004_Ausgabeformate` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md) | `TestCLI004_YAML` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCLI004_YAML_Determinismus` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md) | `TestCLI004_YAML_Negative` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-005`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestCLI005_KeinRepoZugriff` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-005`](../../spec/lastenheft.md) | `TestCLI005_PrintConfig` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_AiHarnessInit_VollKanon` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_AiHarness_Boundary` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCLI006_AiHarness_Determinismus` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_AiHarness_Happy` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_AiHarness_KeinQuellenfehler` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_AiHarness_PlanningAktiv` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_KeineKennungen` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_LeereQuellen` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_QuelleFehlt` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI006_SuggestConfig` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCLI007_DoctorJSON_Determinismus` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-007`](../../spec/lastenheft.md) | `TestCLI007_DoctorJSON_Happy` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md), [`DC-FA-CLI-007`](../../spec/lastenheft.md) | `TestCLI007_DoctorYAML` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCLI007_Doctor_Determinismus` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-007`](../../spec/lastenheft.md) | `TestCLI007_Doctor_Happy` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-007`](../../spec/lastenheft.md) | `TestCLI007_Doctor_KeineBefunde` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-008`](../../spec/lastenheft.md) | `TestCLI008_Repair_Boundary_Stufen` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCLI008_Repair_Determinismus` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-008`](../../spec/lastenheft.md) | `TestCLI008_Repair_Inkompatibel` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-008`](../../spec/lastenheft.md) | `TestCLI008_Repair_RoundTrip` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-QA-03`](../../spec/lastenheft.md) | `TestCLI036_Trace_Markdown` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-010`](../../spec/lastenheft.md) | `TestCLI044_PrintMK_TraceTargets` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-011`](../../spec/lastenheft.md) | `TestCLI044_RequireComplete_KeineWaisen` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-011`](../../spec/lastenheft.md) | `TestCLI044_RequireComplete_OhneTrace` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-003`](../../spec/lastenheft.md), [`DC-FA-CLI-011`](../../spec/lastenheft.md), [`DC-QA-03`](../../spec/lastenheft.md) | `TestCLI044_RequireComplete_Waise` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-010`](../../spec/lastenheft.md), [`DC-FA-COMMITS-001`](../../spec/lastenheft.md) | `TestCLI056_PrintMK_CommitsTarget` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-010`](../../spec/lastenheft.md), [`DC-FA-PLAN-001`](../../spec/lastenheft.md) | `TestCLI057_PrintMK_PlanningTarget` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-010`](../../spec/lastenheft.md), [`DC-FA-TGT-001`](../../spec/lastenheft.md) | `TestCLI063_PrintMK_TargetsTarget` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-009`](../../spec/lastenheft.md) | `TestCLI066_Trace_DefaultByteIdentisch` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-009`](../../spec/lastenheft.md) | `TestCLI066_Trace_FremdKonvention` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-011`](../../spec/lastenheft.md) | `TestCLI066_Trace_RequireCompleteErbtConfig` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-009`](../../spec/lastenheft.md) | `TestCLI066_Trace_VollCustomKonvention` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-COV-001`](../../spec/lastenheft.md) | `TestCLI067_Coverage_RangeDecktAb` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-MOD-001`](../../spec/lastenheft.md) | `TestCLI068_Modality_KlassifikationUndGating` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-REQ-001`](../../spec/lastenheft.md) | `TestCLI070_TraceTable_371UndModalitaet` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCLI071_Cross_DefaultByteIdentisch` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCLI071_Cross_FailClosed` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCLI071_Cross_KonsistentGruen` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-009`](../../spec/lastenheft.md), [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCLI071_Cross_RealerDriftAdvisory` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-011`](../../spec/lastenheft.md), [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCLI071_Cross_RequireCompleteGatet` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-REQ-001`](../../spec/lastenheft.md) | `TestCLI074_DirektivenDatenzeileToleriert` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-LINK-001`](../../spec/lastenheft.md) | `TestCLI076_FenceInfozeileVerdecktLinkNicht` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-LINK-001`](../../spec/lastenheft.md) | `TestCLI076_FenceInfozeileVerdecktTabelleNicht` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-REQ-001`](../../spec/lastenheft.md) | `TestCLI076_TrennzeileEinBindestrich` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCLI077_Cross_RuecksichtNichtVerschluckt` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-REQ-001`](../../spec/lastenheft.md) | `TestCLI077_StillerUebersprungGrenze` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-001`](../../spec/lastenheft.md), [`DC-FA-CLI-010`](../../spec/lastenheft.md) | `TestCLI185_PrintMK_UsageTarget` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CLI-006`](../../spec/lastenheft.md) | `TestCLI234_RB_AbwesendByteGleich` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md) | `TestCODE001` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestCONF001_IDTargetFehlt` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestCONF001_IDTargetVerlaesstWurzel` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestCONF001_LeereConfig` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestCONF001_UngueltigeConfig` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-DIST-001`](../../spec/lastenheft.md) | `TestDIST001_GaenzlichLeereWurzel` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-DIST-001`](../../spec/lastenheft.md) | `TestDIST001_OptionenNachPfad` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-EXT-001`](../../spec/lastenheft.md) | `TestEXT001` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestID001` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-LINK-002`](../../spec/lastenheft.md) | `TestLINK002_Symlinks` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-MTX-001`](../../spec/lastenheft.md) | `TestMTX001` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestQA02_Determinismus` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-SCAN-001`](../../spec/lastenheft.md) | `TestSCAN001` | [`internal/adapter/driving/cli/cli_acceptance_test.go`](../../internal/adapter/driving/cli/cli_acceptance_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestConfigPath_UngueltigesSchema` | [`internal/adapter/driving/cli/cli_config_path_test.go`](../../internal/adapter/driving/cli/cli_config_path_test.go) |
| [`DC-FA-CLI-001`](../../spec/lastenheft.md) | `TestCLI001_HilfeNenntHandbuch` | [`internal/adapter/driving/cli/cli_handbuch_link_test.go`](../../internal/adapter/driving/cli/cli_handbuch_link_test.go) |
| [`DC-FA-CLI-010`](../../spec/lastenheft.md) | `TestCLI010_PrintMKNenntHandbuch` | [`internal/adapter/driving/cli/cli_handbuch_link_test.go`](../../internal/adapter/driving/cli/cli_handbuch_link_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md) | `TestHint_TabUndUmbruchWerdenAbgewiesen` | [`internal/adapter/driving/cli/cli_hint_test.go`](../../internal/adapter/driving/cli/cli_hint_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md) | `TestHint_VierteSpalteInDerBefundZeile` | [`internal/adapter/driving/cli/cli_hint_test.go`](../../internal/adapter/driving/cli/cli_hint_test.go) |
| [`DC-FA-REF-001`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestRefs_DefaultAusByteIdentisch` | [`internal/adapter/driving/cli/cli_refs_test.go`](../../internal/adapter/driving/cli/cli_refs_test.go) |
| [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefs_KeepUndTippfehlerEndToEnd` | [`internal/adapter/driving/cli/cli_refs_test.go`](../../internal/adapter/driving/cli/cli_refs_test.go) |
| [`DC-FA-LINK-002`](../../spec/lastenheft.md), [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefs_SymlinkBleibtTrotzVentil` | [`internal/adapter/driving/cli/cli_refs_test.go`](../../internal/adapter/driving/cli/cli_refs_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md), [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefs_UngueltigesGlobExit2` | [`internal/adapter/driving/cli/cli_refs_test.go`](../../internal/adapter/driving/cli/cli_refs_test.go) |
| [`DC-FA-CLI-010`](../../spec/lastenheft.md) | `TestPrintMk_DocTracked` | [`internal/adapter/driving/cli/cli_tracked_test.go`](../../internal/adapter/driving/cli/cli_tracked_test.go) |
| [`DC-FA-TRK-001`](../../spec/lastenheft.md) | `TestTracked_ExemptTargets` | [`internal/adapter/driving/cli/cli_tracked_test.go`](../../internal/adapter/driving/cli/cli_tracked_test.go) |
| [`DC-FA-TRK-001`](../../spec/lastenheft.md) | `TestTracked_FailClosedOhneGit` | [`internal/adapter/driving/cli/cli_tracked_test.go`](../../internal/adapter/driving/cli/cli_tracked_test.go) |
| [`DC-FA-TRK-001`](../../spec/lastenheft.md) | `TestTracked_StagedZielGruen` | [`internal/adapter/driving/cli/cli_tracked_test.go`](../../internal/adapter/driving/cli/cli_tracked_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestTracked_UngueltigesGlobExit2` | [`internal/adapter/driving/cli/cli_tracked_test.go`](../../internal/adapter/driving/cli/cli_tracked_test.go) |
| [`DC-FA-TRK-001`](../../spec/lastenheft.md) | `TestTracked_UntrackedZiel` | [`internal/adapter/driving/cli/cli_tracked_test.go`](../../internal/adapter/driving/cli/cli_tracked_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestVersionsPatterns_KurzformByteIdentisch` | [`internal/adapter/driving/cli/cli_versions_patterns_test.go`](../../internal/adapter/driving/cli/cli_versions_patterns_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsPatterns_MischformExit2` | [`internal/adapter/driving/cli/cli_versions_patterns_test.go`](../../internal/adapter/driving/cli/cli_versions_patterns_test.go) |
| [`DC-FA-CLI-005`](../../spec/lastenheft.md), [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsPatterns_VorlagenBloeckeEinkommentierbar` | [`internal/adapter/driving/cli/cli_versions_patterns_test.go`](../../internal/adapter/driving/cli/cli_versions_patterns_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsPatterns_ZweiReihen` | [`internal/adapter/driving/cli/cli_versions_patterns_test.go`](../../internal/adapter/driving/cli/cli_versions_patterns_test.go) |
| [`DC-FA-COV-001`](../../spec/lastenheft.md) | `TestExpandRange` | [`internal/hexagon/core/app/trace_coverage_test.go`](../../internal/hexagon/core/app/trace_coverage_test.go) |
| [`DC-FA-COV-001`](../../spec/lastenheft.md) | `TestExpandRangeCommaShortform` | [`internal/hexagon/core/app/trace_coverage_test.go`](../../internal/hexagon/core/app/trace_coverage_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestExpandRangeLinkTransparent` | [`internal/hexagon/core/app/trace_coverage_test.go`](../../internal/hexagon/core/app/trace_coverage_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestBuildTraceMatrixWithoutCrossBlock` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyBackwardIDHeaderDoppelt` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyBackwardIDHeaderFehlt` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyBothDirections` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-COV-001`](../../spec/lastenheft.md), [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyCommaShortform` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyConsistent1N` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyDuplicateHeader` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyExcludeReq` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyFailClosed` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyForwardReqPattern` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyForwardReqPatternDefault` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyLeereVorwaertsSichtIstKeinVakuum` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyPraezedenzQuelleVorAbschnitt` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyRangeAware` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyRangeAwareLinkTransparent` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencySections` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencySupersetMode` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyVakuumBeideSichtenLeer` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyVakuumDurchUebergriffigesVentil` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-XREF-001`](../../spec/lastenheft.md) | `TestCrossConsistencyVakuumRueckSichtLeerUnterSuperset` | [`internal/hexagon/core/app/trace_cross_test.go`](../../internal/hexagon/core/app/trace_cross_test.go) |
| [`DC-FA-MOD-001`](../../spec/lastenheft.md) | `TestModalityClassify` | [`internal/hexagon/core/app/trace_modality_test.go`](../../internal/hexagon/core/app/trace_modality_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md) | `TestAnchorsHTMLModul` | [`internal/hexagon/core/rules/anchors_test.go`](../../internal/hexagon/core/rules/anchors_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md) | `TestAnchorsHTMLSelbeDatei` | [`internal/hexagon/core/rules/anchors_test.go`](../../internal/hexagon/core/rules/anchors_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md) | `TestAnchorsModul` | [`internal/hexagon/core/rules/anchors_test.go`](../../internal/hexagon/core/rules/anchors_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md) | `TestAnchorsSchweigtBeiFehlenderDatei` | [`internal/hexagon/core/rules/anchors_test.go`](../../internal/hexagon/core/rules/anchors_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md) | `TestHTMLAnchors` | [`internal/hexagon/core/rules/anchors_test.go`](../../internal/hexagon/core/rules/anchors_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md) | `TestSlugify` | [`internal/hexagon/core/rules/anchors_test.go`](../../internal/hexagon/core/rules/anchors_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsAbsatzweiteSpanneVerschluckt` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsBacktickPfadIstMalform` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsBlockHappy` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsErwaehnungUndDirektiveInEinerZeile` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsFailClosed` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsFenceTrenntDirektiveVomZitat` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsInlineCodeIstErwaehnung` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsInlineHappy` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsInlineMismatch` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsMindestlaenge` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsOhneDirektiveStill` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsOhneInlineCodeGeprueft` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsStrippenErzeugtDirektive` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CITE-001`](../../spec/lastenheft.md) | `TestCitationsZitatFaeule` | [`internal/hexagon/core/rules/citations_test.go`](../../internal/hexagon/core/rules/citations_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md) | `TestCodepathsCheckLines` | [`internal/hexagon/core/rules/codepaths_checklines_test.go`](../../internal/hexagon/core/rules/codepaths_checklines_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCodepathsCheckLinesDefaultAus` | [`internal/hexagon/core/rules/codepaths_checklines_test.go`](../../internal/hexagon/core/rules/codepaths_checklines_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestCodepathsExemptPaths` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md) | `TestCodepathsHTMLAnker` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md) | `TestCodepathsIgnoreMarkerNurDiesesModul` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestCodepathsIgnoreRefs` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md) | `TestCodepathsIgnoreRefsUnterdruecktEscapeUndAnker` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md) | `TestCodepathsModul` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md) | `TestCodepathsNormalisierungUndAnker` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestIgnoreMarkerBleibtRohBeiVersionsUndDiagrams` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestIgnoreMarkerBleibtTokenBeiDiagrams` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md), [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestIgnoreMarkerBleibtTokenBeiVersions` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md), [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestIgnoreMarkerBrauchtDieKommentarForm` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md) | `TestIgnoreMarkerFormGrenzen` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-CODE-001`](../../spec/lastenheft.md), [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestIgnoreMarkerInInlineCodeIstErwaehnung` | [`internal/hexagon/core/rules/codepaths_test.go`](../../internal/hexagon/core/rules/codepaths_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsAllDefined` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsCustomFence` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsDefinedInValidierung` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsDefinedInVerzeichnisFailClosed` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsMarkerAufSchlusszeileWirkungslos` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsMarkerNurExaktesToken` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsMusterPraezedenz` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsNichtAktiv` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsNichtGelisteterFence` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsTokenGrenzeRegexDelegiert` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsUndefinedID` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsVentilBlockNurEigenerFence` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsVentilBlockUeberOeffnungszeile` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsVentilDatei` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-FA-DIAG-001`](../../spec/lastenheft.md) | `TestDiagramsVentilZeile` | [`internal/hexagon/core/rules/diagrams_test.go`](../../internal/hexagon/core/rules/diagrams_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestExternalDeterminismus` | [`internal/hexagon/core/rules/external_test.go`](../../internal/hexagon/core/rules/external_test.go) |
| [`DC-FA-EXT-001`](../../spec/lastenheft.md) | `TestExternalModul` | [`internal/hexagon/core/rules/external_test.go`](../../internal/hexagon/core/rules/external_test.go) |
| [`DC-FA-EXT-001`](../../spec/lastenheft.md), [`DC-QA-03`](../../spec/lastenheft.md) | `TestExternalOptIn` | [`internal/hexagon/core/rules/external_test.go`](../../internal/hexagon/core/rules/external_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCheckFile_OhneRegelnByteIdentisch` | [`internal/hexagon/core/rules/file_test.go`](../../internal/hexagon/core/rules/file_test.go) |
| [`DC-FA-HOST-001`](../../spec/lastenheft.md) | `TestHostpathsExemptTargets` | [`internal/hexagon/core/rules/hostpaths_test.go`](../../internal/hexagon/core/rules/hostpaths_test.go) |
| [`DC-FA-HOST-001`](../../spec/lastenheft.md) | `TestHostpathsModul` | [`internal/hexagon/core/rules/hostpaths_test.go`](../../internal/hexagon/core/rules/hostpaths_test.go) |
| [`DC-FA-HOST-001`](../../spec/lastenheft.md) | `TestHostpathsTilde` | [`internal/hexagon/core/rules/hostpaths_test.go`](../../internal/hexagon/core/rules/hostpaths_test.go) |
| [`DC-FA-HOST-001`](../../spec/lastenheft.md) | `TestHostpathsTildeOhneUnixVerlust` | [`internal/hexagon/core/rules/hostpaths_test.go`](../../internal/hexagon/core/rules/hostpaths_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestIDsExemptPathsProseDefault` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestIDsIgnoreMarkerProseDefault` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestIDsLinkPolicyAlways` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestIDsLinkPolicyProseDefault` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestIDsModul` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestIDsMusterPraezedenz` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestIDsTargetDarfWurzelNichtVerlassen` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-FA-CONF-001`](../../spec/lastenheft.md) | `TestIDsTargetMussExistieren` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-FA-ID-001`](../../spec/lastenheft.md) | `TestIDsVentileNackteVorkommen` | [`internal/hexagon/core/rules/ids_test.go`](../../internal/hexagon/core/rules/ids_test.go) |
| [`DC-FA-IMM-001`](../../spec/lastenheft.md) | `TestImmutableMarkerInInlineCodeInert` | [`internal/hexagon/core/rules/immutable_test.go`](../../internal/hexagon/core/rules/immutable_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestImmutableNichtAktivUndDispatch` | [`internal/hexagon/core/rules/immutable_test.go`](../../internal/hexagon/core/rules/immutable_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestResolveFromDivergenzMeldungSortiert` | [`internal/hexagon/core/rules/links_resolvefrom_test.go`](../../internal/hexagon/core/rules/links_resolvefrom_test.go) |
| [`DC-FA-LINK-001`](../../spec/lastenheft.md) | `TestPreprocessMarkdown_FenceInfozeileMitBacktickIstFliesstext` | [`internal/hexagon/core/rules/markdown_test.go`](../../internal/hexagon/core/rules/markdown_test.go) |
| [`DC-QA-04`](../../spec/lastenheft.md) | `TestPreprocessMarkdown_MehrzeiligerSpan` | [`internal/hexagon/core/rules/markdown_test.go`](../../internal/hexagon/core/rules/markdown_test.go) |
| [`DC-FA-MTX-003`](../../spec/lastenheft.md) | `TestMatrixAllowIfSameID` | [`internal/hexagon/core/rules/matrix_test.go`](../../internal/hexagon/core/rules/matrix_test.go) |
| [`DC-FA-MTX-002`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestMatrixDownwardDefaultAus` | [`internal/hexagon/core/rules/matrix_test.go`](../../internal/hexagon/core/rules/matrix_test.go) |
| [`DC-FA-MTX-002`](../../spec/lastenheft.md) | `TestMatrixDownwardKanten` | [`internal/hexagon/core/rules/matrix_test.go`](../../internal/hexagon/core/rules/matrix_test.go) |
| [`DC-FA-MTX-002`](../../spec/lastenheft.md) | `TestMatrixDownwardRichtung` | [`internal/hexagon/core/rules/matrix_test.go`](../../internal/hexagon/core/rules/matrix_test.go) |
| [`DC-FA-MTX-001`](../../spec/lastenheft.md) | `TestMatrixModul` | [`internal/hexagon/core/rules/matrix_test.go`](../../internal/hexagon/core/rules/matrix_test.go) |
| [`DC-FA-MTX-001`](../../spec/lastenheft.md) | `TestMatrixSupersedeLineage` | [`internal/hexagon/core/rules/matrix_test.go`](../../internal/hexagon/core/rules/matrix_test.go) |
| [`DC-FA-MTX-003`](../../spec/lastenheft.md) | `TestMatrixTokenReferenz` | [`internal/hexagon/core/rules/matrix_test.go`](../../internal/hexagon/core/rules/matrix_test.go) |
| [`DC-FA-MTX-001`](../../spec/lastenheft.md) | `TestStatusOf` | [`internal/hexagon/core/rules/matrix_test.go`](../../internal/hexagon/core/rules/matrix_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestMentionsDeterministisch` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsEffectiveMatch` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsErkennungsform` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsFehlenderSchluesselBrichtAb` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-CLI-004`](../../spec/lastenheft.md), [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsFindet` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsHappyPath` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsIstMengeIstVereinigung` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsLeereIstMengeFailClosed` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsLeereSollMengeFailClosed` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-MENT-001`](../../spec/lastenheft.md) | `TestMentionsNoteNenntDieBezugsmenge` | [`internal/hexagon/core/rules/mentions_test.go`](../../internal/hexagon/core/rules/mentions_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md), [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsAnkerNurImFenceUnaufloesbar` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsAnkerOhneSection` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsExternInert` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsHappy` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsInlineCodeZwischenInert` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsMarkerBindung` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsNichtAktiv` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsReflow` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsRepoEscape` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsSameFileAnker` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-CONF-002`](../../spec/lastenheft.md), [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsScope` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsStale` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-FA-PIN-001`](../../spec/lastenheft.md) | `TestPinsZielWeg` | [`internal/hexagon/core/rules/pins_test.go`](../../internal/hexagon/core/rules/pins_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestClosureDeterministischeReihenfolge` | [`internal/hexagon/core/rules/planning_closure_test.go`](../../internal/hexagon/core/rules/planning_closure_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestClosureInertOhneDir` | [`internal/hexagon/core/rules/planning_closure_test.go`](../../internal/hexagon/core/rules/planning_closure_test.go) |
| [`DC-FA-PLAN-001`](../../spec/lastenheft.md) | `TestPlanningBlockEndetAnGeteilterAbschnittsgrenze` | [`internal/hexagon/core/rules/planning_test.go`](../../internal/hexagon/core/rules/planning_test.go) |
| [`DC-FA-PLAN-001`](../../spec/lastenheft.md) | `TestPlanningUeberschriftImFenceZaehltNicht` | [`internal/hexagon/core/rules/planning_test.go`](../../internal/hexagon/core/rules/planning_test.go) |
| [`DC-FA-PLAN-001`](../../spec/lastenheft.md) | `TestWavesGleicheKennungZaehltEinmal` | [`internal/hexagon/core/rules/planning_waves_test.go`](../../internal/hexagon/core/rules/planning_waves_test.go) |
| [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefsAliasKoexistenz` | [`internal/hexagon/core/rules/refs_test.go`](../../internal/hexagon/core/rules/refs_test.go) |
| [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefsAnkerBleibtScharf` | [`internal/hexagon/core/rules/refs_test.go`](../../internal/hexagon/core/rules/refs_test.go) |
| [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefsGeteiltesVentilCodepaths` | [`internal/hexagon/core/rules/refs_test.go`](../../internal/hexagon/core/rules/refs_test.go) |
| [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefsGeteiltesVentilLinksUndAnchors` | [`internal/hexagon/core/rules/refs_test.go`](../../internal/hexagon/core/rules/refs_test.go) |
| [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefsKeepGewinnt` | [`internal/hexagon/core/rules/refs_test.go`](../../internal/hexagon/core/rules/refs_test.go) |
| [`DC-FA-REF-001`](../../spec/lastenheft.md) | `TestRefsSkopusIsolation` | [`internal/hexagon/core/rules/refs_test.go`](../../internal/hexagon/core/rules/refs_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestDeterminismus` | [`internal/hexagon/core/rules/run_test.go`](../../internal/hexagon/core/rules/run_test.go) |
| [`DC-FA-SCAN-001`](../../spec/lastenheft.md) | `TestDiscoverFiles` | [`internal/hexagon/core/rules/run_test.go`](../../internal/hexagon/core/rules/run_test.go) |
| [`DC-FA-CLI-002`](../../spec/lastenheft.md) | `TestEffectiveModules` | [`internal/hexagon/core/rules/run_test.go`](../../internal/hexagon/core/rules/run_test.go) |
| [`DC-FA-LINK-001`](../../spec/lastenheft.md) | `TestLinksModul` | [`internal/hexagon/core/rules/run_test.go`](../../internal/hexagon/core/rules/run_test.go) |
| [`DC-FA-CONF-002`](../../spec/lastenheft.md) | `TestRun_ModulScope` | [`internal/hexagon/core/rules/run_test.go`](../../internal/hexagon/core/rules/run_test.go) |
| [`DC-FA-CONF-002`](../../spec/lastenheft.md) | `TestRun_ModulScopeAusserhalbGlobal` | [`internal/hexagon/core/rules/run_test.go`](../../internal/hexagon/core/rules/run_test.go) |
| [`DC-FA-LINK-002`](../../spec/lastenheft.md) | `TestSymlinkAblehnung` | [`internal/hexagon/core/rules/run_test.go`](../../internal/hexagon/core/rules/run_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesAbgeschnittenerMarkerHashFailClosed` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesCaseInsensitiv` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesDriftVollerHash` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesFetchFehlerklassen` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesHappyMarkerUndConfigZip` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesMarkerBindungNaechsterLink` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesMarkerOhneSha256Inert` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesModuleOffKeinNetz` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-QA-03`](../../spec/lastenheft.md) | `TestSourcesNilCheckerNoop` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesRepoInternMarkerInert` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesUnreachableNichtDrift` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesZipEntryLimit` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesZipManifestGolden` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesZipReorderInvariant` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SRC-001`](../../spec/lastenheft.md) | `TestSourcesZipUnpackByteLimit` | [`internal/hexagon/core/rules/sources_test.go`](../../internal/hexagon/core/rules/sources_test.go) |
| [`DC-FA-SPAN-001`](../../spec/lastenheft.md) | `TestFenceUnclosed` | [`internal/hexagon/core/rules/spans_test.go`](../../internal/hexagon/core/rules/spans_test.go) |
| [`DC-FA-SPAN-001`](../../spec/lastenheft.md) | `TestFenceUnclosedZiel` | [`internal/hexagon/core/rules/spans_test.go`](../../internal/hexagon/core/rules/spans_test.go) |
| [`DC-FA-SPAN-001`](../../spec/lastenheft.md) | `TestSpansModul` | [`internal/hexagon/core/rules/spans_test.go`](../../internal/hexagon/core/rules/spans_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestStructureMaxLines_AbwesendByteIdentisch` | [`internal/hexagon/core/rules/structure_maxlines_test.go`](../../internal/hexagon/core/rules/structure_maxlines_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCellMaxOhneSchluesselInert` | [`internal/hexagon/core/rules/structure_tablecell_test.go`](../../internal/hexagon/core/rules/structure_tablecell_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md) | `TestStructureHeadingBefundJeUeberschrift` | [`internal/hexagon/core/rules/structure_test.go`](../../internal/hexagon/core/rules/structure_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md) | `TestStructureHeadingEbene` | [`internal/hexagon/core/rules/structure_test.go`](../../internal/hexagon/core/rules/structure_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md) | `TestStructureHeadingEbeneFlacherAlsAbschnitt` | [`internal/hexagon/core/rules/structure_test.go`](../../internal/hexagon/core/rules/structure_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md) | `TestStructureHeadingHappyPath` | [`internal/hexagon/core/rules/structure_test.go`](../../internal/hexagon/core/rules/structure_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md) | `TestStructureHeadingLexikDesModuls` | [`internal/hexagon/core/rules/structure_test.go`](../../internal/hexagon/core/rules/structure_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestStructureHeadingModulAus` | [`internal/hexagon/core/rules/structure_test.go`](../../internal/hexagon/core/rules/structure_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md) | `TestStructureHeadingOhneUeberschriftWirkungslos` | [`internal/hexagon/core/rules/structure_test.go`](../../internal/hexagon/core/rules/structure_test.go) |
| [`DC-FA-STRUCT-001`](../../spec/lastenheft.md) | `TestStructureHeadingSectionsEach` | [`internal/hexagon/core/rules/structure_test.go`](../../internal/hexagon/core/rules/structure_test.go) |
| [`DC-FA-TGT-001`](../../spec/lastenheft.md) | `TestTargetsTabellenzeileImFenceDokumentiertNicht` | [`internal/hexagon/core/rules/targets_test.go`](../../internal/hexagon/core/rules/targets_test.go) |
| [`DC-FA-LINK-002`](../../spec/lastenheft.md) | `TestCheckTracked_SymlinkZieleUeberspringen` | [`internal/hexagon/core/rules/tracked_test.go`](../../internal/hexagon/core/rules/tracked_test.go) |
| [`DC-FA-VCS-001`](../../spec/lastenheft.md) | `TestVCSFailClosed` | [`internal/hexagon/core/rules/vcs_test.go`](../../internal/hexagon/core/rules/vcs_test.go) |
| [`DC-FA-VCS-001`](../../spec/lastenheft.md) | `TestVCSStatusImFenceZaehltNicht` | [`internal/hexagon/core/rules/vcs_test.go`](../../internal/hexagon/core/rules/vcs_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsAusgabeFolgtDerSortierung` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsCurrentFromAnkerFehlt` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md), [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsCurrentFromAnkerNurImFence` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md), [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsCurrentFromAnkerParitaet` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-ANCH-001`](../../spec/lastenheft.md), [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsCurrentFromDuplikatSlug` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsCurrentFromGanzeDatei` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsCurrentFromHTMLAnker` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsCurrentFromKeineVersion` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsCurrentFromSelbstAusgenommen` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsDedupGleichesTupel` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md), [`DC-QA-02`](../../spec/lastenheft.md) | `TestVersionsEinPaarBehaeltWortlaut` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsFehlermeldungNenntDasPaar` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsGleicherWertBeideErwartungen` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsGleichesPaarZweiTrefferEineErwartung` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsHappy` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsNichtAktiv` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsPinOhneGruppe` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsStaleInFence` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsVentile` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsVentilePaarLokal` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsZeilenMarkerGiltAllenPaaren` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsZweiPaareHappy` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsZweiWerteZweiBefunde` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsZweitesPaarMeldet` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-FA-VER-001`](../../spec/lastenheft.md) | `TestVersionsZweitesPaarQuelleFehltFailClosed` | [`internal/hexagon/core/rules/versions_test.go`](../../internal/hexagon/core/rules/versions_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestCheckTagConflictsSortiert` | [`internal/hexagon/core/rules/workflows_test.go`](../../internal/hexagon/core/rules/workflows_test.go) |
| [`DC-QA-02`](../../spec/lastenheft.md) | `TestWorkflowsInert` | [`internal/hexagon/core/rules/workflows_test.go`](../../internal/hexagon/core/rules/workflows_test.go) |
