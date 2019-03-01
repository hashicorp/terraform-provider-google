package resource

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/terraform"
)

// testStepConfig runs a config-mode test step
func testStepConfig(
	opts terraform.ContextOpts,
	state *terraform.State,
	step TestStep) (*terraform.State, error) {
	return testStep(opts, state, step)
}

func testStep(
	opts terraform.ContextOpts,
	state *terraform.State,
	step TestStep) (*terraform.State, error) {
	// Pre-taint any resources that have been defined in Taint, as long as this
	// is not a destroy step.
	if !step.Destroy {
		if err := testStepTaint(state, step); err != nil {
			return state, err
		}
	}

	mod, err := testModule(opts, step)
	if err != nil {
		return state, err
	}

	// Build the context
	opts.Module = mod
	opts.State = state
	opts.Destroy = step.Destroy
	ctx, err := terraform.NewContext(&opts)
	if err != nil {
		return state, fmt.Errorf("Error initializing context: %s", err)
	}
	if diags := ctx.Validate(); len(diags) > 0 {
		if diags.HasErrors() {
			return nil, errwrap.Wrapf("config is invalid: {{err}}", diags.Err())
		}

		log.Printf("[WARN] Config warnings:\n%s", diags)
	}

	// Refresh!
	state, err = ctx.Refresh()
	if err != nil {
		return state, fmt.Errorf(
			"Error refreshing: %s", err)
	}

	// If this step is a PlanOnly step, skip over this first Plan and subsequent
	// Apply, and use the follow up Plan that checks for perpetual diffs
	if !step.PlanOnly {
		// Plan!
		if p, err := ctx.Plan(); err != nil {
			return state, fmt.Errorf(
				"Error planning: %s", err)
		} else {
			log.Printf("[WARN] Test: Step plan: %s", p)
		}

		// We need to keep a copy of the state prior to destroying
		// such that destroy steps can verify their behaviour in the check
		// function
		stateBeforeApplication := state.DeepCopy()

		// Apply!
		state, err = ctx.Apply()
		if err != nil {
			return state, fmt.Errorf("Error applying: %s", err)
		}

		// Check! Excitement!
		if step.Check != nil {
			if step.Destroy {
				if err := step.Check(stateBeforeApplication); err != nil {
					return state, fmt.Errorf("Check failed: %s", err)
				}
			} else {
				if err := step.Check(state); err != nil {
					return state, fmt.Errorf("Check failed: %s", err)
				}
			}
		}
	}

	// Now, verify that Plan is now empty and we don't have a perpetual diff issue
	// We do this with TWO plans. One without a refresh.
	var p *terraform.Plan
	if p, err = ctx.Plan(); err != nil {
		return state, fmt.Errorf("Error on follow-up plan: %s", err)
	}
	if p.Diff != nil && !p.Diff.Empty() {
		if step.ExpectNonEmptyPlan {
			log.Printf("[INFO] Got non-empty plan, as expected:\n\n%s", p)
		} else {
			return state, fmt.Errorf(
				"After applying this step, the plan was not empty:\n\n%s", p)
		}
	}

	// And another after a Refresh.
	if !step.Destroy || (step.Destroy && !step.PreventPostDestroyRefresh) {
		state, err = ctx.Refresh()
		if err != nil {
			return state, fmt.Errorf(
				"Error on follow-up refresh: %s", err)
		}
	}
	if p, err = ctx.Plan(); err != nil {
		return state, fmt.Errorf("Error on second follow-up plan: %s", err)
	}
	empty := p.Diff == nil || p.Diff.Empty()

	// Data resources are tricky because they legitimately get instantiated
	// during refresh so that they will be already populated during the
	// plan walk. Because of this, if we have any data resources in the
	// config we'll end up wanting to destroy them again here. This is
	// acceptable and expected, and we'll treat it as "empty" for the
	// sake of this testing.
	if step.Destroy {
		empty = true

		for _, moduleDiff := range p.Diff.Modules {
			for k, instanceDiff := range moduleDiff.Resources {
				if !strings.HasPrefix(k, "data.") {
					empty = false
					break
				}

				if !instanceDiff.Destroy {
					empty = false
				}
			}
		}
	}

	if !empty {
		if step.ExpectNonEmptyPlan {
			log.Printf("[INFO] Got non-empty plan, as expected:\n\n%s", p)
		} else {
			return state, fmt.Errorf(
				"After applying this step and refreshing, "+
					"the plan was not empty:\n\n%s", p)
		}
	}

	// Made it here, but expected a non-empty plan, fail!
	if step.ExpectNonEmptyPlan && (p.Diff == nil || p.Diff.Empty()) {
		return state, fmt.Errorf("Expected a non-empty plan, but got an empty plan!")
	}

	// Made it here? Good job test step!
	return state, nil
}

func testStepTaint(state *terraform.State, step TestStep) error {
	for _, p := range step.Taint {
		m := state.RootModule()
		if m == nil {
			return errors.New("no state")
		}
		rs, ok := m.Resources[p]
		if !ok {
			return fmt.Errorf("resource %q not found in state", p)
		}
		log.Printf("[WARN] Test: Explicitly tainting resource %q", p)
		rs.Taint()
	}
	return nil
}
