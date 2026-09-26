package work

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
)

type requestShaperFunc func(context.Context, DraftRequest) (DraftProposal, error)

func (f requestShaperFunc) Shape(ctx context.Context, request DraftRequest) (DraftProposal, error) {
	return f(ctx, request)
}

func validProposal() DraftProposal {
	return DraftProposal{Mission: "Investigate the reported failure.", Constraints: "Unknown: constraints not supplied.", Output: "A reviewable investigation draft."}
}

func TestDraftIntakePreservesSource(t *testing.T) {
	request := DraftRequest{Input: "  Please fix this\r\nKeep café 日本語 and `--flags` exactly.\n", Context: "\tKnown fact only.\n"}
	proposal := validProposal()
	var received DraftRequest
	shaper := requestShaperFunc(func(_ context.Context, input DraftRequest) (DraftProposal, error) {
		received = input
		return proposal, nil
	})
	draft, err := DraftIntake(context.Background(), request, shaper)
	if err != nil {
		t.Fatal(err)
	}
	if received != request || draft.Input != request.Input || draft.Context != request.Context {
		t.Fatalf("source changed: %+v %+v", received, draft)
	}
	encoded, err := json.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]string
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"input": request.Input, "context": request.Context, "mission": proposal.Mission, "constraints": proposal.Constraints, "output": proposal.Output}
	if !reflect.DeepEqual(fields, want) {
		t.Fatalf("unexpected wire fields: %v", fields)
	}
}

func TestDraftIntakeLabelsMissingContext(t *testing.T) {
	for _, supplied := range []string{"", " \n\t"} {
		draft, err := DraftIntake(context.Background(), DraftRequest{Input: "Investigate", Context: supplied}, requestShaperFunc(func(_ context.Context, request DraftRequest) (DraftProposal, error) {
			if request.Context != supplied {
				t.Fatal("provider request context was changed")
			}
			proposal := validProposal()
			proposal.Checkpoint = "Review after investigation."
			return proposal, nil
		}))
		if err != nil || draft.Context != "Unknown: no context supplied." || draft.Checkpoint != "Review after investigation." {
			t.Fatalf("%+v %v", draft, err)
		}
	}
}

func TestDraftIntakeRejectsIncompleteProposal(t *testing.T) {
	for _, field := range []string{"mission", "constraints", "output"} {
		t.Run(field, func(t *testing.T) {
			proposal := validProposal()
			switch field {
			case "mission":
				proposal.Mission = " \n"
			case "constraints":
				proposal.Constraints = ""
			case "output":
				proposal.Output = "\t"
			}
			draft, err := DraftIntake(context.Background(), DraftRequest{Input: "Investigate"}, requestShaperFunc(func(context.Context, DraftRequest) (DraftProposal, error) { return proposal, nil }))
			if err == nil || draft != (domainwork.Draft{}) {
				t.Fatalf("accepted incomplete proposal: %+v %v", draft, err)
			}
		})
	}
}

func TestDraftIntakeValidatesBeforeShaping(t *testing.T) {
	calls := 0
	shaper := requestShaperFunc(func(context.Context, DraftRequest) (DraftProposal, error) { calls++; return validProposal(), nil })
	for _, input := range []string{"", " \n"} {
		if draft, err := DraftIntake(context.Background(), DraftRequest{Input: input}, shaper); err == nil || draft != (domainwork.Draft{}) {
			t.Fatalf("blank input: %+v %v", draft, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if draft, err := DraftIntake(ctx, DraftRequest{Input: "Investigate"}, shaper); !errors.Is(err, context.Canceled) || draft != (domainwork.Draft{}) {
		t.Fatalf("cancelled input: %+v %v", draft, err)
	}
	if calls != 0 {
		t.Fatalf("shaper called %d times for rejected requests", calls)
	}
	if draft, err := DraftIntake(context.Background(), DraftRequest{Input: "Investigate"}, nil); err == nil || draft != (domainwork.Draft{}) {
		t.Fatalf("nil shaper: %+v %v", draft, err)
	}
}

func TestDraftIntakePropagatesFailureAndCancellation(t *testing.T) {
	failure := errors.New("provider failed")
	draft, err := DraftIntake(context.Background(), DraftRequest{Input: "Investigate"}, requestShaperFunc(func(context.Context, DraftRequest) (DraftProposal, error) { return validProposal(), failure }))
	if !errors.Is(err, failure) || draft != (domainwork.Draft{}) {
		t.Fatalf("provider failure: %+v %v", draft, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	draft, err = DraftIntake(ctx, DraftRequest{Input: "Investigate"}, requestShaperFunc(func(passed context.Context, _ DraftRequest) (DraftProposal, error) {
		if passed != ctx {
			t.Fatal("context not propagated")
		}
		cancel()
		return validProposal(), nil
	}))
	if !errors.Is(err, context.Canceled) || draft != (domainwork.Draft{}) {
		t.Fatalf("cancelled output: %+v %v", draft, err)
	}
}
