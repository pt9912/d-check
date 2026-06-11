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

	"gopkg.in/yaml.v3"

	"github.com/pt9912/d-check/internal/hexagon/core"
)

// FileName ist der feste Name der Konfigurationsdatei.
const FileName = ".d-check.yml"

type rawIDPattern struct {
	Regex  string `yaml:"regex"`
	Target string `yaml:"target"`
}

type rawIDs struct {
	Patterns []rawIDPattern `yaml:"patterns"`
}

type rawMatrix struct {
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

type rawExternal struct {
	TimeoutSeconds int `yaml:"timeout-seconds"`
	Parallel       int `yaml:"parallel"`
}

// raw bildet das Voll-Schema von .d-check.yml ab
// (spec/spezifikation.md §2). Unbekannte Schlüssel sind durch
// KnownFields(true) Fehler.
type raw struct {
	Scan *struct {
		Roots  []string `yaml:"roots"`
		Ignore []string `yaml:"ignore"`
	} `yaml:"scan"`
	Modules  []string     `yaml:"modules"`
	IDs      *rawIDs      `yaml:"ids"`
	Matrix   *rawMatrix   `yaml:"matrix"`
	External *rawExternal `yaml:"external"`
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
	return cfg, nil
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
		cfg.IDPatterns = append(cfg.IDPatterns, core.IDPattern{Regex: re, Target: p.Target})
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
	if t := e.TimeoutSeconds; t != 0 && (t < 1 || t > 300) {
		return fmt.Errorf("%s: external.timeout-seconds außerhalb 1–300", FileName)
	}
	if p := e.Parallel; p != 0 && (p < 1 || p > 16) {
		return fmt.Errorf("%s: external.parallel außerhalb 1–16", FileName)
	}
	cfg.External = core.ExternalConfig{TimeoutSeconds: e.TimeoutSeconds, Parallel: e.Parallel}
	return nil
}
