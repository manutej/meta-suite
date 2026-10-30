// Package model holds the Go shapes of the core types the kernel reads and writes. Field order here is the
// output order (encoding/json keeps struct order), and the JSON Schemas in schemas/ are the contract: every
// value the kernel writes is validated against its schema first (G1).
package model

// Source is a claim's provenance.
type Source struct {
	Locator         string `json:"locator"`
	RetrievedAt     string `json:"retrieved_at"`
	IndependenceKey string `json:"independence_key"`
}

// StageToken is kernel-issued at `meta stage begin` and stamped onto everything written in that stage.
type StageToken struct {
	TID        string `json:"tid"`
	Run        string `json:"run"`
	Stage      string `json:"stage"`
	Region     string `json:"region"`
	Agent      string `json:"agent"`
	Model      string `json:"model"`
	PromptHash string `json:"prompt_hash"`
}

// Normalized records a unit conversion the kernel applied.
type Normalized struct {
	Dimension string  `json:"dimension"`
	FromUnit  string  `json:"from_unit"`
	FromValue float64 `json:"from_value"`
}

// Claim is one sourced statement.
type Claim struct {
	ID         string      `json:"id"`
	Subject    string      `json:"subject"`
	Relation   string      `json:"relation"`
	Object     string      `json:"object"`
	Value      any         `json:"value,omitempty"`
	ValueType  string      `json:"value_type"`
	Unit       string      `json:"unit,omitempty"`
	AsOf       string      `json:"as_of,omitempty"`
	Source     Source      `json:"source"`
	Region     string      `json:"region"`
	StageToken StageToken  `json:"stage_token"`
	Normalized *Normalized `json:"normalized,omitempty"`
}

// Object is a graph node (an entity).
type Object struct {
	ID      string   `json:"id"`
	Kind    string   `json:"kind"`
	Labels  []string `json:"labels"`
	Regions []string `json:"regions,omitempty"`
	Members []string `json:"members,omitempty"`
}

// Morphism is a typed relation between two objects, backed by claims.
type Morphism struct {
	ID           string   `json:"id"`
	Src          string   `json:"src"`
	Dst          string   `json:"dst"`
	RelationType string   `json:"relation_type"`
	ClaimIDs     []string `json:"claim_ids"`
}

// ShapeMetrics summarise a graph's shape for A1's layout rules (definitions in internal/graph).
type ShapeMetrics struct {
	HubConcentration float64 `json:"hub_concentration"`
	LongestChain     int     `json:"longest_chain"`
	TwoSidedness     float64 `json:"two_sidedness"`
	TemporalDensity  float64 `json:"temporal_density"`
}

// RegionGraph is one region's graph.
type RegionGraph struct {
	Region       string       `json:"region"`
	Objects      []Object     `json:"objects"`
	Morphisms    []Morphism   `json:"morphisms"`
	ShapeMetrics ShapeMetrics `json:"shape_metrics"`
}

// MergeInfo describes how a DomainGraph was merged.
type MergeInfo struct {
	Plain           bool `json:"plain"`
	Interfaces      int  `json:"interfaces"`
	Identifications int  `json:"identifications"`
}

// DomainGraph is the merged graph.
type DomainGraph struct {
	Regions      []string     `json:"regions"`
	Objects      []Object     `json:"objects"`
	Morphisms    []Morphism   `json:"morphisms"`
	ShapeMetrics ShapeMetrics `json:"shape_metrics"`
	Merge        MergeInfo    `json:"merge"`
}

// IdentifiedObject names two objects (one per region) identified under a label.
type IdentifiedObject struct {
	Label   string   `json:"label"`
	Objects []string `json:"objects"`
}

// Interface is a declared (or candidate) identification between two regions.
type Interface struct {
	Regions             []string           `json:"regions"`
	IdentifiedObjects   []IdentifiedObject `json:"identified_objects"`
	IdentifiedRelations []string           `json:"identified_relations"`
	DeclaredBy          string             `json:"declared_by"`
	Reason              string             `json:"reason,omitempty"`
}

// Key is a finding's grouping key (normalized).
type Key struct {
	Subject  string `json:"subject"`
	Relation string `json:"relation"`
	Object   string `json:"object"`
}

// Side is one side of a tension.
type Side struct {
	Label    string   `json:"label"`
	ClaimIDs []string `json:"claim_ids"`
}

// Finding is a verified / single-source (and, later, other kinds) finding.
type Finding struct {
	ID               string   `json:"id"`
	Kind             string   `json:"kind"`
	Sides            []Side   `json:"sides,omitempty"`
	Evidence         []string `json:"evidence"`
	Region           string   `json:"region,omitempty"`
	Regions          []string `json:"regions,omitempty"`
	Confidence       string   `json:"confidence"`
	PlainText        string   `json:"plain_text"`
	Key              *Key     `json:"key,omitempty"`
	Value            any      `json:"value,omitempty"`
	Unit             string   `json:"unit,omitempty"`
	AsOf             string   `json:"as_of,omitempty"`
	IndependenceKeys []string `json:"independence_keys,omitempty"`
}

// GateResult is one recorded gate verdict in run.json.
type GateResult struct {
	Gate       string   `json:"gate"`
	Verdict    string   `json:"verdict"`
	At         string   `json:"at"`
	Summary    string   `json:"summary,omitempty"`
	InputsHash string   `json:"inputs_hash,omitempty"`
	Files      []string `json:"files,omitempty"`
}

// Usage is the run's usage record.
type Usage struct {
	Source       string `json:"source"`
	Calls        int    `json:"calls"`
	InputTokens  *int   `json:"input_tokens,omitempty"`
	OutputTokens *int   `json:"output_tokens,omitempty"`
}

// Retraction is one line of retractions.jsonl: an append-only tombstone for a claim (`meta claim retract`).
// The claim stays in its shard; every command after it skips the claim.
type Retraction struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
	At     string `json:"at"`
	Stage  string `json:"stage,omitempty"`
}

// RunRecord is run.json.
type RunRecord struct {
	RunID         string            `json:"run_id"`
	CreatedAt     string            `json:"created_at"`
	BriefHash     string            `json:"brief_hash"`
	Mode          string            `json:"mode"`
	Surface       string            `json:"surface"`
	Kernel        string            `json:"kernel"`
	KernelVersion string            `json:"kernel_version,omitempty"`
	ModelMap      map[string]string `json:"model_map"`
	PromptHashes  map[string]string `json:"prompt_hashes"`
	ConfigHash    string            `json:"config_hash"`
	ConfigSource  string            `json:"config_source,omitempty"`
	SchemaVersion string            `json:"schema_version"`
	SpecHash      string            `json:"spec_hash,omitempty"`
	PlaybookHash  *string           `json:"playbook_hash"`
	GateResults   []GateResult      `json:"gate_results"`
	Usage         Usage             `json:"usage"`
	WallMs        int64             `json:"wall_ms"`
	Seed          int64             `json:"seed"`
	InternalView  *bool             `json:"internal_view,omitempty"`
}
