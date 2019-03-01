package terraform

import (
	"log"

	"github.com/hashicorp/terraform/dag"
)

// GraphNodeAttachResourceState is an interface that can be implemented
// to request that a ResourceState is attached to the node.
type GraphNodeAttachResourceState interface {
	// The address to the resource for the state
	ResourceAddr() *ResourceAddress

	// Sets the state
	AttachResourceState(*ResourceState)
}

// AttachStateTransformer goes through the graph and attaches
// state to nodes that implement the interfaces above.
type AttachStateTransformer struct {
	State *State // State is the root state
}

func (t *AttachStateTransformer) Transform(g *Graph) error {
	// If no state, then nothing to do
	if t.State == nil {
		log.Printf("[DEBUG] Not attaching any state: state is nil")
		return nil
	}

	filter := &StateFilter{State: t.State}
	for _, v := range g.Vertices() {
		// Only care about nodes requesting we're adding state
		an, ok := v.(GraphNodeAttachResourceState)
		if !ok {
			continue
		}
		addr := an.ResourceAddr()

		// Get the module state
		results, err := filter.Filter(addr.String())
		if err != nil {
			return err
		}

		// Attach the first resource state we get
		found := false
		for _, result := range results {
			if rs, ok := result.Value.(*ResourceState); ok {
				log.Printf(
					"[DEBUG] Attaching resource state to %q: %#v",
					dag.VertexName(v), rs)
				an.AttachResourceState(rs)
				found = true
				break
			}
		}

		if !found {
			log.Printf(
				"[DEBUG] Resource state not found for %q: %s",
				dag.VertexName(v), addr)
		}
	}

	return nil
}
