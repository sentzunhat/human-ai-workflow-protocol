package work

import (
	"context"
	"fmt"
	"strings"

	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
)

// DraftRequest contains only caller-supplied text. No repository is required.
type DraftRequest struct {
	Input   string
	Context string
}

// DraftProposal excludes input/context so a shaper cannot replace source text.
// Unknown constraints or output requirements must be labeled, not invented.
type DraftProposal struct {
	Mission     string
	Constraints string
	Output      string
	Checkpoint  string
}

// RequestShaper proposes intake fields. Implementations must honor cancellation.
// No default adapter is selected; context summarizers have a different contract.
type RequestShaper interface {
	Shape(context.Context, DraftRequest) (DraftProposal, error)
}

// DraftIntake returns a reviewable shape without creating a UUID, files, or a
// backlog entry. The use case owns no persistence; injected shapers own their
// execution behavior. Validation is structural, not semantic certification.
func DraftIntake(ctx context.Context, request DraftRequest, shaper RequestShaper) (domainwork.Draft, error) {
	if err := ctx.Err(); err != nil {
		return domainwork.Draft{}, err
	}
	if strings.TrimSpace(request.Input) == "" {
		return domainwork.Draft{}, fmt.Errorf("draft input must not be blank")
	}
	if shaper == nil {
		return domainwork.Draft{}, fmt.Errorf("request shaper is required")
	}
	proposal, err := shaper.Shape(ctx, request)
	if err != nil {
		return domainwork.Draft{}, fmt.Errorf("shape intake: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domainwork.Draft{}, err
	}
	suppliedContext := request.Context
	if strings.TrimSpace(suppliedContext) == "" {
		suppliedContext = "Unknown: no context supplied."
	}
	draft := domainwork.Draft{
		Input: request.Input, Context: suppliedContext,
		Mission: proposal.Mission, Constraints: proposal.Constraints,
		Output: proposal.Output, Checkpoint: proposal.Checkpoint,
	}
	if err := draft.Validate(); err != nil {
		return domainwork.Draft{}, err
	}
	return draft, nil
}
