// Package configyaml ist der driven Config-Adapter: striktes Decoding
// und statische Validierung von .d-check.yml (DC-FA-CONF-001;
// Strategie gemäß ADR-0003). Den Datei-Inhalt beschafft das CLI über
// den Filesystem-Adapter — dieser Adapter macht kein I/O
// (spec/architecture.md §2). Die Existenz-Prüfung deklarierter
// Scan-Wurzeln erfolgt im Kern beim Lauf-Start (ebenfalls Exit 2).
package configyaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/pt9912/d-check/internal/hexagon/core"
)

// FileName ist der feste Name der Konfigurationsdatei.
const FileName = ".d-check.yml"

type rawIDPattern struct {
	Regex       string   `yaml:"regex"`
	Target      string   `yaml:"target"`
	LinkPolicy  string   `yaml:"link-policy"`
	ExemptPaths []string `yaml:"exempt-paths"`
}

type rawIDs struct {
	Scope    *rawScope      `yaml:"scope"`
	Patterns []rawIDPattern `yaml:"patterns"`
}

// rawScope ist der modul-lokale Scan-Scope (DC-FA-CONF-002). Roots
// ist Pflicht, wenn scope gesetzt ist — nil (fehlend) ist von der
// expliziten leeren Liste unterscheidbar.
type rawScope struct {
	Roots  []string `yaml:"roots"`
	Ignore []string `yaml:"ignore"`
}

// rawScopeOnly traegt Module, die ausser scope keine eigenen
// Konfigurations-Schluessel haben (links, anchors).
type rawScopeOnly struct {
	Scope *rawScope `yaml:"scope"`
}

type rawMatrix struct {
	Scope   *rawScope `yaml:"scope"`
	Classes []struct {
		Name  string   `yaml:"name"`
		Paths []string `yaml:"paths"`
	} `yaml:"classes"`
	Rules []struct {
		From  string `yaml:"from"`
		To    string `yaml:"to"`
		Allow bool   `yaml:"allow"`
	} `yaml:"rules"`
	Status *struct {
		Forbidden []string `yaml:"forbidden"`
	} `yaml:"status"`
	ExcludeSections []string `yaml:"exclude-sections"`
}

// rawExternal nutzt Pointer, damit ein explizit gesetzter Wert 0 vom
// Nicht-Setzen unterscheidbar bleibt (Constraint 1–300 bzw. 1–16 —
// auch 0 ist ein Konfigurationsfehler, kein stiller Default).
type rawExternal struct {
	Scope          *rawScope `yaml:"scope"`
	TimeoutSeconds *int      `yaml:"timeout-seconds"`
	Parallel       *int      `yaml:"parallel"`
}

// raw bildet das Voll-Schema von .d-check.yml ab
// (spec/spezifikation.md §2). Unbekannte Schlüssel sind durch
// KnownFields(true) Fehler.
type raw struct {
	Scan *struct {
		Roots  []string `yaml:"roots"`
		Ignore []string `yaml:"ignore"`
	} `yaml:"scan"`
	Modules  []string      `yaml:"modules"`
	Links    *rawScopeOnly `yaml:"links"`
	Anchors  *rawScopeOnly `yaml:"anchors"`
	Spans    *rawScopeOnly `yaml:"spans"`
	Hostpaths *struct {
		Scope    *rawScope `yaml:"scope"`
		Prefixes []string  `yaml:"prefixes"`
	} `yaml:"hostpaths"`
	IDs      *rawIDs       `yaml:"ids"`
	Matrix   *rawMatrix    `yaml:"matrix"`
	External *rawExternal  `yaml:"external"`
	Codepaths *struct {
		Scope *rawScope `yaml:"scope"`
		Roots []string  `yaml:"roots"`
	} `yaml:"codepaths"`
}

// Decode parst und validiert den Datei-Inhalt vollständig — Syntax
// und statische Semantik; jeder Fehler bedeutet Exit 2 ohne Prüfung.
// content == nil bedeutet: keine Konfigurationsdatei → Defaults.
func Decode(content []byte) (core.Config, error) {
	var cfg core.Config
	if content == nil {
		return cfg, nil
	}
	r, err := decodeStrict(content)
	if err != nil {
		return cfg, err
	}
	if r == nil {
		return cfg, nil // leeres YAML-Null-Dokument → Defaults
	}

	if r.Scan != nil {
		cfg.Roots = r.Scan.Roots
		cfg.Ignore = r.Scan.Ignore
	}
	cfg.Modules = r.Modules

	if err := applyIDs(r.IDs, &cfg); err != nil {
		return cfg, err
	}
	if err := applyMatrix(r.Matrix, &cfg); err != nil {
		return cfg, err
	}
	if err := applyExternal(r.External, &cfg); err != nil {
		return cfg, err
	}
	if err := applyCodepaths(r, &cfg); err != nil {
		return cfg, err
	}
	if err := applyHostpaths(r, &cfg); err != nil {
		return cfg, err
	}
	if err := applyScopes(r, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// applyCodepaths validiert die codepaths-Präfixe (DC-FA-CODE-001).
func applyCodepaths(r *raw, cfg *core.Config) error {
	if r.Codepaths == nil {
		return nil
	}
	for _, root := range r.Codepaths.Roots {
		if strings.TrimSpace(root) == "" {
			return fmt.Errorf("%s: codepaths.roots enthält ein leeres Präfix", FileName)
		}
		if strings.HasPrefix(root, "/") || root == ".." || strings.Contains(root, "..") {
			return fmt.Errorf("%s: codepaths.roots-Präfix %q muss relativ zur Repo-Wurzel liegen (kein '/', kein '..')", FileName, root)
		}
	}
	cfg.Codepaths = core.CodepathsConfig{Roots: r.Codepaths.Roots}
	return nil
}

// applyHostpaths validiert die hostpaths-Präfixliste (DC-FA-HOST-001:
// nicht-leere Verzeichnisnamen ohne '/').
func applyHostpaths(r *raw, cfg *core.Config) error {
	if r.Hostpaths == nil {
		return nil
	}
	for _, p := range r.Hostpaths.Prefixes {
		if strings.TrimSpace(p) == "" {
			return fmt.Errorf("%s: hostpaths.prefixes enthält einen leeren Namen", FileName)
		}
		if strings.Contains(p, "/") {
			return fmt.Errorf("%s: hostpaths.prefixes-Eintrag %q muss ein Verzeichnisname ohne '/' sein", FileName, p)
		}
	}
	cfg.Hostpaths = core.HostpathsConfig{Prefixes: r.Hostpaths.Prefixes}
	return nil
}

// applyScopes übernimmt die modul-lokalen Scan-Scopes
// (DC-FA-CONF-002): scope ersetzt den globalen Scan für genau dieses
// Modul; roots ist Pflicht (keine stille Vererbung).
func applyScopes(r *raw, cfg *core.Config) error {
	scopes := []struct {
		module string
		scope  *rawScope
	}{
		{"links", scopeOf(r.Links)},
		{"anchors", scopeOf(r.Anchors)},
		{"spans", scopeOf(r.Spans)},
		{"hostpaths", scopeOfHostpaths(r.Hostpaths)},
		{"ids", scopeOfIDs(r.IDs)},
		{"matrix", scopeOfMatrix(r.Matrix)},
		{"external", scopeOfExternal(r.External)},
		{"codepaths", scopeOfCodepaths(r.Codepaths)},
	}
	for _, sc := range scopes {
		if sc.scope == nil {
			continue
		}
		if sc.scope.Roots == nil {
			return fmt.Errorf("%s: %s.scope.roots fehlt — scope ersetzt den globalen Scan und braucht explizite Wurzeln (leere Liste prüft nichts)", FileName, sc.module)
		}
		if cfg.Scopes == nil {
			cfg.Scopes = map[string]*core.ScopeConfig{}
		}
		cfg.Scopes[sc.module] = &core.ScopeConfig{Roots: sc.scope.Roots, Ignore: sc.scope.Ignore}
	}
	return nil
}

// scopeOf-Helfer: nil-sichere Extraktion des scope-Schluessels der
// jeweiligen Modul-Sektion.
func scopeOf(v *rawScopeOnly) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfIDs(v *rawIDs) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfMatrix(v *rawMatrix) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfExternal(v *rawExternal) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfHostpaths(v *struct {
	Scope    *rawScope `yaml:"scope"`
	Prefixes []string  `yaml:"prefixes"`
}) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfCodepaths(v *struct {
	Scope *rawScope `yaml:"scope"`
	Roots []string  `yaml:"roots"`
}) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

// decodeStrict dekodiert mit KnownFields; nil-raw bei leerem Dokument
// (Review R1/B2).
func decodeStrict(content []byte) (*raw, error) {
	var r raw
	dec := yaml.NewDecoder(bytes.NewReader(content))
	dec.KnownFields(true)
	if err := dec.Decode(&r); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		// interne Go-Typnamen aus der yaml-Meldung entfernen
		// (Review R1/B5).
		typeLeak := regexp.MustCompile(` in type \S+`)
		return nil, fmt.Errorf("%s: %s", FileName, typeLeak.ReplaceAllString(err.Error(), ""))
	}
	return &r, nil
}

func applyIDs(ids *rawIDs, cfg *core.Config) error {
	if ids == nil {
		return nil
	}
	for i, p := range ids.Patterns {
		re, err := regexp.Compile(p.Regex)
		if err != nil {
			return fmt.Errorf("%s: ids.patterns[%d].regex: %v", FileName, i, err)
		}
		if re.MatchString("") {
			return fmt.Errorf("%s: ids.patterns[%d].regex matcht den Leerstring", FileName, i)
		}
		if p.Target == "" {
			return fmt.Errorf("%s: ids.patterns[%d].target fehlt", FileName, i)
		}
		if p.LinkPolicy != "" && p.LinkPolicy != "prose" && p.LinkPolicy != core.AlwaysPolicy {
			return fmt.Errorf("%s: ids.patterns[%d].link-policy: ungültig %q (erlaubt: prose, always)", FileName, i, p.LinkPolicy)
		}
		cfg.IDPatterns = append(cfg.IDPatterns, core.IDPattern{
			Regex: re, Target: p.Target,
			LinkPolicy: p.LinkPolicy, ExemptPaths: p.ExemptPaths,
		})
	}
	return nil
}

func applyMatrix(m *rawMatrix, cfg *core.Config) error {
	if m == nil {
		return nil
	}
	classes := map[string]bool{}
	for i, c := range m.Classes {
		if c.Name == "" || classes[c.Name] {
			return fmt.Errorf("%s: matrix.classes[%d].name fehlt oder doppelt", FileName, i)
		}
		classes[c.Name] = true
		cfg.Matrix.Classes = append(cfg.Matrix.Classes, core.MatrixClass{Name: c.Name, Paths: c.Paths})
	}
	for i, rule := range m.Rules {
		if !classes[rule.From] || !classes[rule.To] {
			return fmt.Errorf("%s: matrix.rules[%d] referenziert undeklarierte Klasse", FileName, i)
		}
		cfg.Matrix.Rules = append(cfg.Matrix.Rules, core.MatrixRule{From: rule.From, To: rule.To, Allow: rule.Allow})
	}
	if m.Status != nil {
		cfg.Matrix.StatusForbidden = m.Status.Forbidden
	} else {
		cfg.Matrix.StatusForbidden = []string{"superseded", "deprecated"}
	}
	cfg.Matrix.ExcludeSections = m.ExcludeSections
	return nil
}

func applyExternal(e *rawExternal, cfg *core.Config) error {
	if e == nil {
		return nil
	}
	if t := e.TimeoutSeconds; t != nil {
		if *t < 1 || *t > 300 {
			return fmt.Errorf("%s: external.timeout-seconds außerhalb 1–300", FileName)
		}
		cfg.External.TimeoutSeconds = *t
	}
	if p := e.Parallel; p != nil {
		if *p < 1 || *p > 16 {
			return fmt.Errorf("%s: external.parallel außerhalb 1–16", FileName)
		}
		cfg.External.Parallel = *p
	}
	return nil
}
