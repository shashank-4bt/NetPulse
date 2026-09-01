package incidents

const (
	MinRecoveriesToResolve = 2
)

type ResolutionInput struct {
	RecoverySampleCount int
	IdentifiedCause     bool
}

type ResolutionDecision struct {
	OK     bool
	Reason string
}

func CanMarkResolved(input ResolutionInput) ResolutionDecision {
	if input.RecoverySampleCount < MinRecoveriesToResolve {
		return ResolutionDecision{
			OK:     false,
			Reason: "One recovered measurement is not enough to mark an incident resolved.",
		}
	}
	if !input.IdentifiedCause {
		return ResolutionDecision{
			OK:     false,
			Reason: "Resolution requires an identified cause backed by evidence, not recovery alone.",
		}
	}
	return ResolutionDecision{
		OK:     true,
		Reason: "Independent recoveries and an identified cause are present. Resolution is still a stored judgment, not a population claim.",
	}
}
