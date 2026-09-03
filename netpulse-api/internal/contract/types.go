package contract

const (
	ModelVersion       = "0.12.0"
	EngineVersion      = "0.12.0"
	RuleVersion        = "0.12.0-admin-platform"
	MeasurementVersion = "0.6.0-dns-tcp-tls-http"
	ObservedFailures   = "Elevated connectivity failures observed"
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
	OK                bool                 `json:"ok"`
	Diagnosis         *Diagnosis           `json:"diagnosis,omitempty"`
	Services          []Service            `json:"services,omitempty"`
	Service           *Service             `json:"service,omitempty"`
	Intelligence      *ServiceIntelligence `json:"intelligence,omitempty"`
	Incidents         []Incident           `json:"incidents"`
	Incident          *Incident            `json:"incident,omitempty"`
	Page              *Page                `json:"page,omitempty"`
	Health            *Health              `json:"health,omitempty"`
	Map               *MapAggregates       `json:"map,omitempty"`
	User              *User                `json:"user,omitempty"`
	Session           *Session             `json:"session,omitempty"`
	Sessions          []Session            `json:"sessions,omitempty"`
	Events            []SecurityEvent      `json:"events,omitempty"`
	Auth              *AuthState           `json:"auth,omitempty"`
	Dashboard         *Dashboard           `json:"dashboard,omitempty"`
	Diagnoses         []DiagnosisSummary   `json:"diagnoses,omitempty"`
	Reports           []UserReport         `json:"reports,omitempty"`
	Saved             []SavedService       `json:"savedServices,omitempty"`
	Devices           []Device             `json:"devices,omitempty"`
	Alerts            *AlertPreferences    `json:"alerts,omitempty"`
	Billing           *Billing             `json:"billing,omitempty"`
	Privacy           *PrivacySettings     `json:"privacy,omitempty"`
	Share             *ShareLink           `json:"share,omitempty"`
	SessionToken      string               `json:"sessionToken,omitempty"`
	Workspace         *Workspace           `json:"workspace,omitempty"`
	Monitors          []Monitor            `json:"monitors,omitempty"`
	Monitor           *Monitor             `json:"monitor,omitempty"`
	APIKeys           []APIKey             `json:"apiKeys,omitempty"`
	APIKey            *APIKey              `json:"apiKey,omitempty"`
	KeySecret         string               `json:"keySecret,omitempty"`
	Webhooks          []Webhook            `json:"webhooks,omitempty"`
	Webhook           *Webhook             `json:"webhook,omitempty"`
	WebhookSecret     string               `json:"webhookSecret,omitempty"`
	Deliveries        []WebhookDelivery    `json:"deliveries,omitempty"`
	AlertRules        []AlertRule          `json:"alertRules,omitempty"`
	AlertRule         *AlertRule           `json:"alertRule,omitempty"`
	Usage             *Usage               `json:"usage,omitempty"`
	SLA               *SLAReport           `json:"sla,omitempty"`
	DevDashboard      *DeveloperDashboard  `json:"developerDashboard,omitempty"`
	DevIncidents      []DeveloperIncident  `json:"developerIncidents,omitempty"`
	DevIncident       *DeveloperIncident   `json:"developerIncident,omitempty"`
	Checks            []MonitorCheck       `json:"checks,omitempty"`
	Organization      *Organization        `json:"organization,omitempty"`
	Organizations     []Organization       `json:"organizations,omitempty"`
	Members           []Member             `json:"members,omitempty"`
	Member            *Member              `json:"member,omitempty"`
	Teams             []Team               `json:"teams,omitempty"`
	Team              *Team                `json:"team,omitempty"`
	OrgDevices        []OrgDevice          `json:"orgDevices,omitempty"`
	OrgDevice         *OrgDevice           `json:"orgDevice,omitempty"`
	Networks          []OrgNetwork         `json:"networks,omitempty"`
	Network           *OrgNetwork          `json:"network,omitempty"`
	OrgServices       []OrgService         `json:"orgServices,omitempty"`
	OrgService        *OrgService          `json:"orgService,omitempty"`
	OrgIncidents      []OrgIncident        `json:"orgIncidents,omitempty"`
	OrgIncident       *OrgIncident         `json:"orgIncident,omitempty"`
	OrgDashboard      *OrgDashboard        `json:"orgDashboard,omitempty"`
	Analytics         *OrgAnalytics        `json:"analytics,omitempty"`
	OrgReports        []OrgReport          `json:"orgReports,omitempty"`
	OrgReport         *OrgReport           `json:"orgReport,omitempty"`
	AuditEvents       []AuditEvent         `json:"auditEvents,omitempty"`
	Invites           []OrgInvite          `json:"invites,omitempty"`
	OrgDiagnoses      []OrgDiagnosis       `json:"orgDiagnoses,omitempty"`
	Permissions       []string             `json:"permissions,omitempty"`
	Operator          *Operator            `json:"operator,omitempty"`
	AdminSystem       *AdminSystem         `json:"adminSystem,omitempty"`
	AdminUsers        []AdminUser          `json:"adminUsers,omitempty"`
	AdminUser         *AdminUser           `json:"adminUser,omitempty"`
	AdminMeasurements []AdminMeasurement   `json:"adminMeasurements,omitempty"`
	AdminDiagnoses    []AdminDiagnosis     `json:"adminDiagnoses,omitempty"`
	AdminRules        []DiagnosticRule     `json:"adminRules,omitempty"`
	RuleOutcomes      *RuleOutcomes        `json:"ruleOutcomes,omitempty"`
	AbuseEvents       []AbuseEvent         `json:"abuseEvents,omitempty"`
	AdminAudit        []AdminAudit         `json:"adminAudit,omitempty"`
	FeatureFlags      []FeatureFlag        `json:"featureFlags,omitempty"`
	FeatureFlag       *FeatureFlag         `json:"featureFlag,omitempty"`
	RemoteConfig      []RemoteConfigEntry  `json:"remoteConfig,omitempty"`
	AdminIncidents    []AdminIncident      `json:"adminIncidents,omitempty"`
	AdminIncident     *AdminIncident       `json:"adminIncident,omitempty"`
	Error             *APIError            `json:"error,omitempty"`
}

type Service struct {
	Slug     string   `json:"slug"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Summary  string   `json:"summary"`
	Layers   []string `json:"layers"`
}

type Observation struct {
	Value        any     `json:"value"`
	Unit         *string `json:"unit"`
	Measured     bool    `json:"measured"`
	SampleCount  int     `json:"sampleCount"`
	SampleWindow *string `json:"sampleWindow"`
	Summary      string  `json:"summary"`
}

type SliceObservation struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Status      string `json:"status"`
	SampleCount int    `json:"sampleCount"`
	Summary     string `json:"summary"`
}

type ServiceIntelligence struct {
	CurrentState      string             `json:"currentState"`
	Health            *int               `json:"health"`
	LastUpdated       *string            `json:"lastUpdated"`
	Availability      Observation        `json:"availability"`
	Latency           Observation        `json:"latency"`
	Errors            Observation        `json:"errors"`
	RegionalHealth    []SliceObservation `json:"regionalHealth"`
	NetworkHealth     []SliceObservation `json:"networkHealth"`
	RecentIncidentIDs []string           `json:"recentIncidentIds"`
}

type IncidentTimelineEvent struct {
	Stage  string  `json:"stage"`
	Label  string  `json:"label"`
	Status string  `json:"status"`
	At     *string `json:"at"`
	Note   string  `json:"note"`
}

type Incident struct {
	ID                string                  `json:"id"`
	Title             string                  `json:"title"`
	Severity          string                  `json:"severity"`
	Status            string                  `json:"status"`
	Scope             string                  `json:"scope"`
	StartedAt         string                  `json:"startedAt"`
	LastUpdatedAt     string                  `json:"lastUpdatedAt"`
	AffectedServices  []string                `json:"affectedServices"`
	Regions           []string                `json:"regions"`
	Networks          []string                `json:"networks"`
	Evidence          []Evidence              `json:"evidence"`
	Hypotheses        []Hypothesis            `json:"hypotheses"`
	Confidence        Confidence              `json:"confidence"`
	Timeline          []IncidentTimelineEvent `json:"timeline"`
	SampleCount       int                     `json:"sampleCount"`
	SampleRate        *string                 `json:"sampleRate"`
	AffectedUserCount *int                    `json:"affectedUserCount"`
}

type Page struct {
	Number int `json:"number"`
	Size   int `json:"size"`
	Total  int `json:"total"`
}

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	DisplayName    string `json:"displayName"`
	EmailVerified  bool   `json:"emailVerified"`
	CreatedAt      string `json:"createdAt"`
	TelemetryOptIn bool   `json:"telemetryOptIn"`
}

type Session struct {
	ID         string `json:"id"`
	UserID     string `json:"-"`
	TokenHash  string `json:"-"`
	CreatedAt  string `json:"createdAt"`
	LastSeenAt string `json:"lastSeenAt"`
	ExpiresAt  string `json:"expiresAt"`
	Current    bool   `json:"current"`
	Revoked    bool   `json:"revoked"`
	UserAgent  string `json:"userAgent"`
	IP         string `json:"ip"`
	Label      string `json:"label"`
}

type SecurityEvent struct {
	ID      string `json:"id"`
	UserID  string `json:"-"`
	Kind    string `json:"kind"`
	At      string `json:"at"`
	Summary string `json:"summary"`
	IP      string `json:"ip"`
}

type AuthMethods struct {
	Password bool     `json:"password"`
	OAuth    []string `json:"oauth"`
	Passkeys int      `json:"passkeys"`
	MFA      []string `json:"mfa"`
}

type AuthState struct {
	EmailSent   bool        `json:"emailSent"`
	EmailReason string      `json:"emailReason,omitempty"`
	DevToken    string      `json:"devToken,omitempty"`
	Methods     AuthMethods `json:"methods"`
}

type DiagnosisSummary struct {
	ID        string  `json:"id"`
	Target    string  `json:"target"`
	Status    string  `json:"status"`
	Outcome   *string `json:"outcome"`
	CreatedAt string  `json:"createdAt"`
}

type UserReport struct {
	ID        string  `json:"id"`
	Target    string  `json:"target"`
	Status    string  `json:"status"`
	Shared    bool    `json:"shared"`
	CreatedAt string  `json:"createdAt"`
	Outcome   *string `json:"outcome"`
}

type SavedService struct {
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
}

type Device struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	UserAgent string `json:"userAgent"`
	IP        string `json:"ip"`
	LastSeen  string `json:"lastSeenAt"`
	Current   bool   `json:"current"`
	Kind      string `json:"kind"`
}

type AlertPreferences struct {
	EmailEnabled   bool   `json:"emailEnabled"`
	IncidentAlerts bool   `json:"incidentAlerts"`
	DeliveredCount int    `json:"deliveredCount"`
	Summary        string `json:"summary"`
}

type Billing struct {
	HasAccount     bool             `json:"hasAccount"`
	OrganizationID *string          `json:"organizationId"`
	Plan           *string          `json:"plan"`
	Invoices       []BillingInvoice `json:"invoices"`
	Summary        string           `json:"summary"`
}

type BillingInvoice struct {
	ID     string `json:"id"`
	Amount string `json:"amount"`
	Status string `json:"status"`
}

type PrivacySettings struct {
	TelemetryOptIn  bool   `json:"telemetryOptIn"`
	Retention       string `json:"retention"`
	Deletion        string `json:"deletion"`
	Collected       string `json:"collected"`
	Purpose         string `json:"purpose"`
	BrowsingHistory string `json:"browsingHistory"`
}

type ShareLink struct {
	Token   string `json:"token"`
	Path    string `json:"path"`
	Summary string `json:"summary"`
}

type Dashboard struct {
	InternetHealth string             `json:"internetHealth"`
	NetworkInfo    string             `json:"networkInfo"`
	Diagnoses      []DiagnosisSummary `json:"diagnoses"`
	SavedServices  []SavedService     `json:"savedServices"`
	Incidents      []Incident         `json:"incidents"`
	Reports        []UserReport       `json:"reports"`
	Alerts         AlertPreferences   `json:"alerts"`
}

func EmptyAuthMethods() AuthMethods {
	return AuthMethods{Password: true, OAuth: []string{}, Passkeys: 0, MFA: []string{}}
}

func EmptyBilling() Billing {
	return Billing{
		HasAccount:     false,
		OrganizationID: nil,
		Plan:           nil,
		Invoices:       []BillingInvoice{},
		Summary:        "No billing account. Organizations and invoices are not enabled.",
	}
}

func DefaultPrivacy() PrivacySettings {
	return PrivacySettings{
		TelemetryOptIn:  false,
		Collected:       "Account email, password hash, session metadata, diagnoses you start while signed in, and saved service slugs.",
		Purpose:         "Sign-in, session security, and showing your own diagnosis history. Not advertising.",
		Retention:       "Sessions expire after 7 days of inactivity policy or on revoke. Raw probe observations stay short-lived. Exact store TTLs are not claimed beyond that.",
		Deletion:        "You can delete the account. That removes the user, sessions, tokens, saved services, and owned diagnoses from this store.",
		BrowsingHistory: "NetPulse does not collect a browsing history. Diagnosis uses only the target you submit.",
	}
}

type Workspace struct {
	ID        string `json:"id"`
	OwnerID   string `json:"-"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type MonitorThresholds struct {
	AvailabilityBelow *float64 `json:"availabilityBelow"`
	LatencyMsAbove    *int     `json:"latencyMsAbove"`
	ErrorRateAbove    *float64 `json:"errorRateAbove"`
}

type Monitor struct {
	ID          string            `json:"id"`
	WorkspaceID string            `json:"-"`
	OrgID       string            `json:"-"`
	DeviceID    *string           `json:"deviceId"`
	NetworkID   *string           `json:"networkId"`
	ServiceID   *string           `json:"serviceId"`
	Name        string            `json:"name"`
	Target      string            `json:"target"`
	Type        string            `json:"type"`
	Regions     []string          `json:"regions"`
	FrequencyS  int               `json:"frequencySeconds"`
	TimeoutS    int               `json:"timeoutSeconds"`
	Thresholds  MonitorThresholds `json:"thresholds"`
	Status      string            `json:"status"`
	Summary     string            `json:"summary"`
	CheckCount  int               `json:"checkCount"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
}

type MonitorCheck struct {
	ID          string `json:"id"`
	MonitorID   string `json:"monitorId"`
	WorkspaceID string `json:"-"`
	OrgID       string `json:"-"`
	Network     string `json:"network,omitempty"`
	ASN         string `json:"asn,omitempty"`
	Service     string `json:"service,omitempty"`
	Endpoint    string `json:"endpoint,omitempty"`
	Device      string `json:"device,omitempty"`
	Region      string `json:"region"`
	OK          bool   `json:"ok"`
	LatencyMs   *int   `json:"latencyMs"`
	At          string `json:"at"`
	Summary     string `json:"summary"`
}

type APIKey struct {
	ID              string   `json:"id"`
	WorkspaceID     string   `json:"-"`
	OrgID           string   `json:"-"`
	Name            string   `json:"name"`
	Prefix          string   `json:"prefix"`
	Last4           string   `json:"last4"`
	Hash            string   `json:"-"`
	Scopes          []string `json:"scopes"`
	RateLimitPerMin int      `json:"rateLimitPerMin"`
	Revoked         bool     `json:"revoked"`
	CreatedAt       string   `json:"createdAt"`
	LastUsedAt      *string  `json:"lastUsedAt"`
}

type Webhook struct {
	ID          string   `json:"id"`
	WorkspaceID string   `json:"-"`
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Secret      string   `json:"-"`
	SecretHint  string   `json:"secretHint"`
	CreatedAt   string   `json:"createdAt"`
	Disabled    bool     `json:"disabled"`
}

type WebhookDelivery struct {
	ID             string  `json:"id"`
	WebhookID      string  `json:"webhookId"`
	WorkspaceID    string  `json:"-"`
	Event          string  `json:"event"`
	EventID        string  `json:"eventId"`
	Timestamp      string  `json:"timestamp"`
	IdempotencyKey string  `json:"idempotencyKey"`
	Signature      string  `json:"signature"`
	Payload        string  `json:"-"`
	Attempt        int     `json:"attempt"`
	Status         string  `json:"status"`
	NextRetryAt    *string `json:"nextRetryAt"`
	Summary        string  `json:"summary"`
}

type AlertRule struct {
	ID             string  `json:"id"`
	WorkspaceID    string  `json:"-"`
	Kind           string  `json:"kind"`
	MonitorID      *string `json:"monitorId"`
	Threshold      float64 `json:"threshold"`
	Enabled        bool    `json:"enabled"`
	DeliveredCount int     `json:"deliveredCount"`
	Summary        string  `json:"summary"`
	CreatedAt      string  `json:"createdAt"`
}

type DeveloperIncident struct {
	ID          string  `json:"id"`
	WorkspaceID string  `json:"-"`
	MonitorID   string  `json:"monitorId"`
	Title       string  `json:"title"`
	Status      string  `json:"status"`
	StartedAt   string  `json:"startedAt"`
	ResolvedAt  *string `json:"resolvedAt"`
	SampleCount int     `json:"sampleCount"`
	Summary     string  `json:"summary"`
}

type Usage struct {
	Requests     int    `json:"requests"`
	Measurements int    `json:"measurements"`
	Monitors     int    `json:"monitors"`
	Regions      int    `json:"regions"`
	Webhooks     int    `json:"webhooks"`
	Summary      string `json:"summary"`
}

type PercentilePoint struct {
	P50         *float64 `json:"p50"`
	P95         *float64 `json:"p95"`
	P99         *float64 `json:"p99"`
	SampleCount int      `json:"sampleCount"`
	Summary     string   `json:"summary"`
}

type RegionalSlice struct {
	Region      string `json:"region"`
	SampleCount int    `json:"sampleCount"`
	Status      string `json:"status"`
	Summary     string `json:"summary"`
}

type DeveloperDashboard struct {
	Availability Observation         `json:"availability"`
	Latency      PercentilePoint     `json:"latency"`
	Incidents    []DeveloperIncident `json:"incidents"`
	Regional     []RegionalSlice     `json:"regionalPerformance"`
	Summary      string              `json:"summary"`
}

type SLAReport struct {
	Availability Observation         `json:"availability"`
	Downtime     Observation         `json:"downtime"`
	Latency      PercentilePoint     `json:"latency"`
	Incidents    []DeveloperIncident `json:"incidents"`
	Regional     []RegionalSlice     `json:"regionalPerformance"`
	Summary      string              `json:"summary"`
}

func EmptyPercentiles() PercentilePoint {
	return PercentilePoint{
		SampleCount: 0,
		Summary:     "Observed sample count: 0. Percentiles are not estimated from an empty series.",
	}
}

func EmptyDeveloperDashboard() DeveloperDashboard {
	return DeveloperDashboard{
		Availability: UnmeasuredObservation("Availability"),
		Latency:      EmptyPercentiles(),
		Incidents:    []DeveloperIncident{},
		Regional:     []RegionalSlice{},
		Summary:      "No monitor checks are stored. Availability and latency stay unmeasured.",
	}
}

func EmptySLA() SLAReport {
	return SLAReport{
		Availability: UnmeasuredObservation("Availability"),
		Downtime:     UnmeasuredObservation("Downtime"),
		Latency:      EmptyPercentiles(),
		Incidents:    []DeveloperIncident{},
		Regional:     []RegionalSlice{},
		Summary:      "No SLA window can be computed until monitor checks are stored.",
	}
}

func EmptyUsage() Usage {
	return Usage{Summary: "Usage counts are stored request totals, not estimated traffic."}
}

func EmptyAlerts() AlertPreferences {
	return AlertPreferences{
		EmailEnabled:   false,
		IncidentAlerts: false,
		DeliveredCount: 0,
		Summary:        "No alerts have been delivered. Notification delivery is not configured.",
	}
}

type Health struct {
	Status  string            `json:"status"`
	Version string            `json:"version"`
	Storage map[string]string `json:"storage"`
}

type MapCell struct {
	ID          string   `json:"id"`
	Level       string   `json:"level"`
	Label       string   `json:"label"`
	ParentID    *string  `json:"parentId,omitempty"`
	Lon         *float64 `json:"lon,omitempty"`
	Lat         *float64 `json:"lat,omitempty"`
	SampleCount int      `json:"sampleCount"`
	Status      string   `json:"status"`
	Summary     string   `json:"summary"`
	Layer       string   `json:"layer"`
	ChildCount  int      `json:"childCount"`
}

type MapIncidentRef struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	CoarseRegion string `json:"coarseRegion"`
	SampleCount  int    `json:"sampleCount"`
}

type MapAggregates struct {
	Level          string           `json:"level"`
	ParentID       *string          `json:"parentId,omitempty"`
	Cells          []MapCell        `json:"cells"`
	IncidentRefs   []MapIncidentRef `json:"incidentRefs"`
	TotalSamples   int              `json:"totalSamples"`
	Limit          int              `json:"limit"`
	Truncated      bool             `json:"truncated"`
	Precision      string           `json:"precision"`
	Reason         string           `json:"reason"`
	HasCoordinates bool             `json:"hasCoordinates"`
}

func EmptyConfidence() Confidence {
	return Confidence{
		SupportingEvidenceIDs:    []string{},
		AlternativeHypothesisIDs: []string{},
		Caveat:                   ConfidenceCaveat,
	}
}

func UnmeasuredObservation(topic string) Observation {
	return Observation{
		Value:       nil,
		Measured:    false,
		SampleCount: 0,
		Summary:     "Observed sample count: 0. " + topic + " is not estimated for a population.",
	}
}

func EmptyIntelligence() ServiceIntelligence {
	return ServiceIntelligence{
		CurrentState:      "not_measured",
		Health:            nil,
		LastUpdated:       nil,
		Availability:      UnmeasuredObservation("Availability"),
		Latency:           UnmeasuredObservation("Latency"),
		Errors:            UnmeasuredObservation("Error rate"),
		RegionalHealth:    []SliceObservation{},
		NetworkHealth:     []SliceObservation{},
		RecentIncidentIDs: []string{},
	}
}

func NormalizeMap(agg MapAggregates) MapAggregates {
	if agg.Cells == nil {
		agg.Cells = []MapCell{}
	}
	if agg.IncidentRefs == nil {
		agg.IncidentRefs = []MapIncidentRef{}
	}
	if agg.Limit <= 0 {
		agg.Limit = 250
	}
	if agg.Precision == "" {
		agg.Precision = "none"
	}
	if agg.Level == "" {
		agg.Level = "world"
	}
	return agg
}

func NormalizeIncident(item Incident) Incident {
	if item.AffectedServices == nil {
		item.AffectedServices = []string{}
	}
	if item.Regions == nil {
		item.Regions = []string{}
	}
	if item.Networks == nil {
		item.Networks = []string{}
	}
	if item.Evidence == nil {
		item.Evidence = []Evidence{}
	}
	if item.Hypotheses == nil {
		item.Hypotheses = []Hypothesis{}
	}
	if item.Timeline == nil {
		item.Timeline = []IncidentTimelineEvent{}
	}
	item.AffectedUserCount = nil
	if item.Confidence.Caveat == "" {
		item.Confidence = EmptyConfidence()
	}
	return item
}

func CurrentVersions() Versions {
	return Versions{
		DiagnosticEngineVersion: EngineVersion,
		RuleVersion:             RuleVersion,
		MeasurementVersion:      MeasurementVersion,
		ModelVersion:            ModelVersion,
	}
}

type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Role      string `json:"role,omitempty"`
	Summary   string `json:"summary"`
}

type Member struct {
	ID          string   `json:"id"`
	OrgID       string   `json:"-"`
	UserID      string   `json:"userId"`
	Email       string   `json:"email"`
	DisplayName string   `json:"displayName"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"createdAt"`
}

type OrgInvite struct {
	ID        string `json:"id"`
	OrgID     string `json:"-"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	Summary   string `json:"summary"`
}

type Team struct {
	ID        string   `json:"id"`
	OrgID     string   `json:"-"`
	Name      string   `json:"name"`
	MemberIDs []string `json:"memberIds"`
	CreatedAt string   `json:"createdAt"`
	Summary   string   `json:"summary"`
}

type OrgDevice struct {
	ID        string  `json:"id"`
	OrgID     string  `json:"-"`
	Name      string  `json:"name"`
	Label     string  `json:"label"`
	Region    string  `json:"region"`
	NetworkID *string `json:"networkId"`
	Summary   string  `json:"summary"`
	CreatedAt string  `json:"createdAt"`
}

type OrgNetwork struct {
	ID        string `json:"id"`
	OrgID     string `json:"-"`
	Name      string `json:"name"`
	ASN       string `json:"asn"`
	Region    string `json:"region"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"createdAt"`
}

type OrgService struct {
	ID        string `json:"id"`
	OrgID     string `json:"-"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Endpoint  string `json:"endpoint"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"createdAt"`
}

type OrgIncident struct {
	ID          string   `json:"id"`
	OrgID       string   `json:"-"`
	MonitorID   string   `json:"monitorId"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	StartedAt   string   `json:"startedAt"`
	ResolvedAt  *string  `json:"resolvedAt"`
	DeviceIDs   []string `json:"deviceIds"`
	NetworkIDs  []string `json:"networkIds"`
	ServiceIDs  []string `json:"serviceIds"`
	Regions     []string `json:"regions"`
	SampleCount int      `json:"sampleCount"`
	Summary     string   `json:"summary"`
}

type OrgDiagnosis struct {
	ID          string `json:"id"`
	OrgID       string `json:"-"`
	DiagnosisID string `json:"diagnosisId"`
	Target      string `json:"target"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	Summary     string `json:"summary"`
}

type OrgDashboard struct {
	OverallHealth string          `json:"overallHealth"`
	Availability  Observation     `json:"availability"`
	Incidents     []OrgIncident   `json:"incidents"`
	Devices       []OrgDevice     `json:"affectedDevices"`
	Regions       []RegionalSlice `json:"regions"`
	Networks      []OrgNetwork    `json:"networks"`
	Services      []OrgService    `json:"services"`
	Summary       string          `json:"summary"`
}

type OrgAnalytics struct {
	Filters      map[string]string `json:"filters"`
	Availability Observation       `json:"availability"`
	Latency      PercentilePoint   `json:"latency"`
	SampleCount  int               `json:"sampleCount"`
	Incidents    []OrgIncident     `json:"incidents"`
	Summary      string            `json:"summary"`
}

type OrgReport struct {
	ID           string          `json:"id"`
	OrgID        string          `json:"-"`
	Kind         string          `json:"kind"`
	Title        string          `json:"title"`
	Availability Observation     `json:"availability"`
	Latency      PercentilePoint `json:"latency"`
	Incidents    []OrgIncident   `json:"incidents"`
	Regions      []RegionalSlice `json:"regions"`
	Networks     []OrgNetwork    `json:"networks"`
	Findings     []string        `json:"findings"`
	SampleCount  int             `json:"sampleCount"`
	CreatedAt    string          `json:"createdAt"`
	Summary      string          `json:"summary"`
}

type AuditEvent struct {
	ID      string `json:"id"`
	OrgID   string `json:"-"`
	ActorID string `json:"actorId"`
	Kind    string `json:"kind"`
	At      string `json:"at"`
	Summary string `json:"summary"`
}

func EmptyOrgDashboard() OrgDashboard {
	return OrgDashboard{
		OverallHealth: "Not measured. No stored organization checks exist.",
		Availability:  UnmeasuredObservation("Availability"),
		Incidents:     []OrgIncident{},
		Devices:       []OrgDevice{},
		Regions:       []RegionalSlice{},
		Networks:      []OrgNetwork{},
		Services:      []OrgService{},
		Summary:       "Organization health is computed from stored checks only. Empty series stay unmeasured.",
	}
}

func EmptyOrgAnalytics() OrgAnalytics {
	return OrgAnalytics{
		Filters:      map[string]string{},
		Availability: UnmeasuredObservation("Availability"),
		Latency:      EmptyPercentiles(),
		Incidents:    []OrgIncident{},
		Summary:      "No stored samples match the selected filters.",
	}
}

type Operator struct {
	UserID      string   `json:"userId"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type HealthComponent struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Detail   string `json:"detail"`
	Measured bool   `json:"measured"`
}

type AdminSystem struct {
	API                 HealthComponent `json:"api"`
	Worker              HealthComponent `json:"worker"`
	Queue               HealthComponent `json:"queue"`
	Database            HealthComponent `json:"database"`
	Cache               HealthComponent `json:"cache"`
	MeasurementFailures Observation     `json:"measurementFailures"`
	ErrorRate           Observation     `json:"errorRates"`
	Latency             PercentilePoint `json:"latency"`
	Summary             string          `json:"summary"`
}

type AdminUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	DisplayName   string `json:"displayName"`
	EmailVerified bool   `json:"emailVerified"`
	CreatedAt     string `json:"createdAt"`
	Summary       string `json:"summary"`
}

type AdminMeasurement struct {
	DiagnosisID string `json:"diagnosisId"`
	Key         string `json:"key"`
	Label       string `json:"label"`
	Measured    bool   `json:"measured"`
	Summary     string `json:"summary"`
}

type AdminDiagnosis struct {
	ID        string `json:"id"`
	Target    string `json:"target"`
	Status    string `json:"status"`
	UserID    string `json:"userId,omitempty"`
	CreatedAt string `json:"createdAt"`
	Summary   string `json:"summary"`
}

type DiagnosticRule struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Layer      string            `json:"layer"`
	Thresholds map[string]string `json:"thresholds"`
	Summary    string            `json:"summary"`
}

type RuleOutcomes struct {
	Statuses       map[string]int `json:"statuses"`
	FalsePositives int            `json:"falsePositives"`
	FalseNegatives int            `json:"falseNegatives"`
	SampleCount    int            `json:"sampleCount"`
	Summary        string         `json:"summary"`
}

type AbuseEvent struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Actor    string `json:"actor"`
	IP       string `json:"ip"`
	Resource string `json:"resource"`
	At       string `json:"at"`
	Result   string `json:"result"`
	Summary  string `json:"summary"`
}

type AdminAudit struct {
	ID       string `json:"id"`
	ActorID  string `json:"actorId"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
	At       string `json:"at"`
	Result   string `json:"result"`
	Summary  string `json:"summary"`
}

type FeatureFlag struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Environment string   `json:"environment"`
	Enabled     bool     `json:"enabled"`
	Percentage  int      `json:"percentage"`
	UserIDs     []string `json:"userIds"`
	OrgIDs      []string `json:"orgIds"`
	UpdatedAt   string   `json:"updatedAt"`
	Summary     string   `json:"summary"`
	TargetMatch *bool    `json:"targetMatch,omitempty"`
}

type RemoteConfigEntry struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	UpdatedAt string `json:"updatedAt"`
	Summary   string `json:"summary"`
}

type IncidentNote struct {
	ID         string `json:"id"`
	IncidentID string `json:"incidentId"`
	Kind       string `json:"kind"`
	Body       string `json:"body"`
	ActorID    string `json:"actorId"`
	At         string `json:"at"`
}

type IncidentOverride struct {
	Classification string `json:"classification"`
	Reason         string `json:"reason"`
	ActorID        string `json:"actorId"`
	At             string `json:"at"`
}

type AdminIncident struct {
	Incident
	Notes    []IncidentNote    `json:"notes"`
	Override *IncidentOverride `json:"override"`
}

type RuleLabel struct {
	ID          string `json:"id"`
	DiagnosisID string `json:"diagnosisId"`
	Kind        string `json:"kind"`
	ActorID     string `json:"actorId"`
	At          string `json:"at"`
}

func EmptyAdminSystem() AdminSystem {
	return AdminSystem{
		API:                 HealthComponent{Name: "API", Status: "up", Detail: "This process answered the request.", Measured: true},
		Worker:              HealthComponent{Name: "Worker", Status: "unmeasured", Detail: "No worker heartbeat is stored.", Measured: false},
		Queue:               HealthComponent{Name: "Queue", Status: "observed", Detail: "Queue depth: 0.", Measured: true},
		Database:            HealthComponent{Name: "Database", Status: "unmeasured", Detail: "No live database probe is stored.", Measured: false},
		Cache:               HealthComponent{Name: "Cache", Status: "unmeasured", Detail: "No live cache probe is stored.", Measured: false},
		MeasurementFailures: UnmeasuredObservation("Measurement failures"),
		ErrorRate:           UnmeasuredObservation("Error rate"),
		Latency:             EmptyPercentiles(),
		Summary:             "System figures use stored process state and stored measurements only. Empty series stay unmeasured.",
	}
}

func EmptyOrgBilling(orgID string) Billing {
	id := orgID
	return Billing{
		HasAccount:     false,
		OrganizationID: &id,
		Plan:           nil,
		Invoices:       []BillingInvoice{},
		Summary:        "No billing account is stored for this organization. Invoices are not enabled.",
	}
}
