package terraform

import (
	"log"

	"github.com/hashicorp/terraform/config/module"
	"github.com/hashicorp/terraform/dag"
)

// GraphNodeDestroyer must be implemented by nodes that destroy resources.
type GraphNodeDestroyer interface {
	dag.Vertex

	// ResourceAddr is the address of the resource that is being
	// destroyed by this node. If this returns nil, then this node
	// is not destroying anything.
	DestroyAddr() *ResourceAddress
}

// GraphNodeCreator must be implemented by nodes that create OR update resources.
type GraphNodeCreator interface {
	// ResourceAddr is the address of the resource being created or updated
	CreateAddr() *ResourceAddress
}

// DestroyEdgeTransformer is a GraphTransformer that creates the proper
// references for destroy resources. Destroy resources are more complex
// in that they must be depend on the destruction of resources that
// in turn depend on the CREATION of the node being destroy.
//
// That is complicated. Visually:
//
//   B_d -> A_d -> A -> B
//
// Notice that A destroy depends on B destroy, while B create depends on
// A create. They're inverted. This must be done for example because often
// dependent resources will block parent resources from deleting. Concrete
// example: VPC with subnets, the VPC can't be deleted while there are
// still subnets.
type DestroyEdgeTransformer struct {
	// These are needed to properly build the graph of dependencies
	// to determine what a destroy node depends on. Any of these can be nil.
	Module *module.Tree
	State  *State
}

func (t *DestroyEdgeTransformer) Transform(g *Graph) error {
	log.Printf("[TRACE] DestroyEdgeTransformer: Beginning destroy edge transformation...")

	// Build a map of what is being destroyed (by address string) to
	// the list of destroyers. In general there will only be one destroyer
	// but to make it more robust we support multiple.
	destroyers := make(map[string][]GraphNodeDestroyer)
	for _, v := range g.Vertices() {
		dn, ok := v.(GraphNodeDestroyer)
		if !ok {
			continue
		}

		addr := dn.DestroyAddr()
		if addr == nil {
			continue
		}

		key := addr.String()
		log.Printf(
			"[TRACE] DestroyEdgeTransformer: %s destroying %q",
			dag.VertexName(dn), key)
		destroyers[key] = append(destroyers[key], dn)
	}

	// If we aren't destroying anything, there will be no edges to make
	// so just exit early and avoid future work.
	if len(destroyers) == 0 {
		return nil
	}

	// Go through and connect creators to destroyers. Going along with
	// our example, this makes: A_d => A
	for _, v := range g.Vertices() {
		cn, ok := v.(GraphNodeCreator)
		if !ok {
			continue
		}

		addr := cn.CreateAddr()
		if addr == nil {
			continue
		}

		key := addr.String()
		ds := destroyers[key]
		if len(ds) == 0 {
			continue
		}

		for _, d := range ds {
			// For illustrating our example
			a_d := d.(dag.Vertex)
			a := v

			log.Printf(
				"[TRACE] DestroyEdgeTransformer: connecting creator/destroyer: %s, %s",
				dag.VertexName(a), dag.VertexName(a_d))

			g.Connect(&DestroyEdge{S: a, T: a_d})
		}
	}

	// This is strange but is the easiest way to get the dependencies
	// of a node that is being destroyed. We use another graph to make sure
	// the resource is in the graph and ask for references. We have to do this
	// because the node that is being destroyed may NOT be in the graph.
	//
	// Example: resource A is force new, then destroy A AND create A are
	// in the graph. BUT if resource A is just pure destroy, then only
	// destroy A is in the graph, and create A is not.
	providerFn := func(a *NodeAbstractProvider) dag.Vertex {
		return &NodeApplyableProvider{NodeAbstractProvider: a}
	}
	steps := []GraphTransformer{
		// Add the local values
		&LocalTransformer{Module: t.Module},

		// Add outputs and metadata
		&OutputTransformer{Module: t.Module},
		&AttachResourceConfigTransformer{Module: t.Module},
		&AttachStateTransformer{State: t.State},

		TransformProviders(nil, providerFn, t.Module),

		// Add all the variables. We can depend on resources through
		// variables due to module parameters, and we need to properly
		// determine that.
		&RootVariableTransformer{Module: t.Module},
		&ModuleVariableTransformer{Module: t.Module},

		&ReferenceTransformer{},
	}

	// Go through all the nodes being destroyed and create a graph.
	// The resulting graph is only of things being CREATED. For example,
	// following our example, the resulting graph would be:
	//
	//   A, B (with no edges)
	//
	var tempG Graph
	var tempDestroyed []dag.Vertex
	for d, _ := range destroyers {
		// d is what is being destroyed. We parse the resource address
		// which it came from it is a panic if this fails.
		addr, err := ParseResourceAddress(d)
		if err != nil {
			panic(err)
		}

		// This part is a little bit weird but is the best way to
		// find the dependencies we need to: build a graph and use the
		// attach config and state transformers then ask for references.
		abstract := &NodeAbstractResource{Addr: addr}
		tempG.Add(abstract)
		tempDestroyed = append(tempDestroyed, abstract)

		// We also add the destroy version here since the destroy can
		// depend on things that the creation doesn't (destroy provisioners).
		destroy := &NodeDestroyResource{NodeAbstractResource: abstract}
		tempG.Add(destroy)
		tempDestroyed = append(tempDestroyed, destroy)
	}

	// Run the graph transforms so we have the information we need to
	// build references.
	for _, s := range steps {
		if err := s.Transform(&tempG); err != nil {
			return err
		}
	}

	log.Printf("[TRACE] DestroyEdgeTransformer: reference graph: %s", tempG.String())

	// Go through all the nodes in the graph and determine what they
	// depend on.
	for _, v := range tempDestroyed {
		// Find all ancestors of this to determine the edges we'll depend on
		vs, err := tempG.Ancestors(v)
		if err != nil {
			return err
		}

		refs := make([]dag.Vertex, 0, vs.Len())
		for _, raw := range vs.List() {
			refs = append(refs, raw.(dag.Vertex))
		}

		refNames := make([]string, len(refs))
		for i, ref := range refs {
			refNames[i] = dag.VertexName(ref)
		}
		log.Printf(
			"[TRACE] DestroyEdgeTransformer: creation node %q references %s",
			dag.VertexName(v), refNames)

		// If we have no references, then we won't need to do anything
		if len(refs) == 0 {
			continue
		}

		// Get the destroy node for this. In the example of our struct,
		// we are currently at B and we're looking for B_d.
		rn, ok := v.(GraphNodeResource)
		if !ok {
			continue
		}

		addr := rn.ResourceAddr()
		if addr == nil {
			continue
		}

		dns := destroyers[addr.String()]

		// We have dependencies, check if any are being destroyed
		// to build the list of things that we must depend on!
		//
		// In the example of the struct, if we have:
		//
		//   B_d => A_d => A => B
		//
		// Then at this point in the algorithm we started with B_d,
		// we built B (to get dependencies), and we found A. We're now looking
		// to see if A_d exists.
		var depDestroyers []dag.Vertex
		for _, v := range refs {
			rn, ok := v.(GraphNodeResource)
			if !ok {
				continue
			}

			addr := rn.ResourceAddr()
			if addr == nil {
				continue
			}

			key := addr.String()
			if ds, ok := destroyers[key]; ok {
				for _, d := range ds {
					depDestroyers = append(depDestroyers, d.(dag.Vertex))
					log.Printf(
						"[TRACE] DestroyEdgeTransformer: destruction of %q depends on %s",
						key, dag.VertexName(d))
				}
			}
		}

		// Go through and make the connections. Use the variable
		// names "a_d" and "b_d" to reference our example.
		for _, a_d := range dns {
			for _, b_d := range depDestroyers {
				if b_d != a_d {
					g.Connect(dag.BasicEdge(b_d, a_d))
				}
			}
		}
	}

	return nil
}
