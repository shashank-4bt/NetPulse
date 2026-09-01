package contract

const (
	ModelVersion       = "0.6.0"
	EngineVersion      = "0.6.0"
	RuleVersion        = "0.6.0-worker-vantage"
	MeasurementVersion = "0.6.0-dns-tcp-tls-http"
	ConfidenceCaveat   = "Confidence is not certainty. A level or percentage can be wrong and must not be treated as proof."
	InsufficientCause  = "NetPulse could not safely determine the root cause."
	InsufficientNext   = "Next recommended check: compare a user-path measurement with this worker vantage. Worker probes cannot isolate device, Wi-Fi, or ISP faults on their own."
)

type Target struct {
	Raw         string `json:"raw"`
	Hostname    string `json:"hostname"`
	Kind        string `json:"kind"`
	ServiceSlug string `json:"serviceSlug"`
}

type Measurement struct {
	ID         string  `json:"id"`
	Key        string  `json:"key"`
	Label      string  `json:"label"`
	Value      any     `json:"value"`
	Unit       *string `json:"unit"`
	Measured   bool    `json:"measured"`
	MeasuredAt *string `json:"measuredAt"`
	Layer      *string `json:"layer"`
	Summary    *string `json:"summary"`
}

type Evidence struct {
	ID             string   `json:"id"`
	EvidenceClass  string   `json:"evidenceClass"`
	Title          string   `json:"title"`
	Body           string   `json:"body"`
	MeasurementIDs []string `json:"measurementIds"`
	Layer          *string  `json:"layer"`
	ObservedAt     *string  `json:"observedAt"`
}

type Confidence struct {
	Level                    *string  `json:"level"`
	Percent                  *int     `json:"percent"`
	SupportingEvidenceIDs    []string `json:"supportingEvidenceIds"`
	AlternativeHypothesisIDs []string `json:"alternativeHypothesisIds"`
	Caveat                   string   `json:"caveat"`
}

type Hypothesis struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	EvidenceIDs []string   `json:"evidenceIds"`
	Layer       *string    `json:"layer"`
	Confidence  Confidence `json:"confidence"`
	Kind        string     `json:"kind,omitempty"`
}

type Recommendation struct {
	ID             string   `json:"id"`
	Action         string   `json:"action"`
	Reason         string   `json:"reason"`
	Risk           string   `json:"risk"`
	ExpectedResult string   `json:"expectedResult"`
	Verification   string   `json:"verification"`
	SafetyClass    string   `json:"safetyClass"`
	AutoExecute    bool     `json:"autoExecute"`
	EvidenceIDs    []string `json:"evidenceIds"`
}

type VerificationStep struct {
	ID            string  `json:"id"`
	Label         string  `json:"label"`
	Status        string  `json:"status"`
	Note          string  `json:"note"`
	ComparedRunID *string `json:"comparedRunId"`
}

type EscalationCondition struct {
	ID          string `json:"id"`
	When        string `json:"when"`
	Action      string `json:"action"`
	SafetyClass string `json:"safetyClass"`
}

type GraphNode struct {
	ID             string     `json:"id"`
	Label          string     `json:"label"`
	Status         string     `json:"status"`
	Confidence     Confidence `json:"confidence"`
	Timestamp      *string    `json:"timestamp"`
	EvidenceIDs    []string   `json:"evidenceIds"`
	MeasurementIDs []string   `json:"measurementIds"`
}

type Step struct {
	ID         string  `json:"id"`
	Label      string  `json:"label"`
	State      string  `json:"state"`
	DurationMs *int    `json:"durationMs"`
	Note       *string `json:"note"`
}

type Versions struct {
	DiagnosticEngineVersion string `json:"diagnosticEngineVersion"`
	RuleVersion             string `json:"ruleVersion"`
	MeasurementVersion      string `json:"measurementVersion"`
	ModelVersion            string `json:"modelVersion"`
}

type InsufficientEvidence struct {
	Determined bool   `json:"determined"`
	Message    string `json:"message"`
	NextCheck  string `json:"nextCheck"`
}

type Report struct {
	ReportID              string                `json:"reportId"`
	Target                Target                `json:"target"`
	Timestamp             string                `json:"timestamp"`
	Outcome               string                `json:"outcome"`
	Tests                 []Step                `json:"tests"`
	Measurements          []Measurement         `json:"measurements"`
	Evidence              []Evidence            `json:"evidence"`
	Hypotheses            []Hypothesis          `json:"hypotheses"`
	AlternativeHypotheses []Hypothesis          `json:"alternativeHypotheses"`
	LikelyCause           *string               `json:"likelyCause"`
	Confidence            Confidence            `json:"confidence"`
	Recommendations       []Recommendation      `json:"recommendations"`
	VerificationSteps     []VerificationStep    `json:"verificationSteps"`
	EscalationConditions  []EscalationCondition `json:"escalationConditions"`
	Graph                 []GraphNode           `json:"graph"`
	Versions              Versions              `json:"versions"`
	InsufficientEvidence  InsufficientEvidence  `json:"insufficientEvidence"`
	EngineVersion         string                `json:"engineVersion"`
}

type Diagnosis struct {
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Error   *APIError `json:"error"`
	Report  *Report   `json:"report"`
	Created string    `json:"createdAt"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Envelope struct {
	OK        bool       `json:"ok"`
	Diagnosis *Diagnosis `json:"diagnosis,omitempty"`
	Services  []Service  `json:"services,omitempty"`
	Service   *Service   `json:"service,omitempty"`
	Incidents []Incident `json:"incidents"`
	Health    *Health    `json:"health,omitempty"`
	Error     *APIError  `json:"error,omitempty"`
}

type Service struct {
	Slug     string   `json:"slug"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Summary  string   `json:"summary"`
	Layers   []string `json:"layers"`
}

type Incident struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Scope     string `json:"scope"`
	StartedAt string `json:"startedAt"`
}

type Health struct {
	Status  string            `json:"status"`
	Version string            `json:"version"`
	Storage map[string]string `json:"storage"`
}

func EmptyConfidence() Confidence {
	return Confidence{
		SupportingEvidenceIDs:    []string{},
		AlternativeHypothesisIDs: []string{},
		Caveat:                   ConfidenceCaveat,
	}
}

func CurrentVersions() Versions {
	return Versions{
		DiagnosticEngineVersion: EngineVersion,
		RuleVersion:             RuleVersion,
		MeasurementVersion:      MeasurementVersion,
		ModelVersion:            ModelVersion,
	}
}
