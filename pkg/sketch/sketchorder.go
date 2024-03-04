package sketch

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hansmi/dossier/internal/sketcherror"
	"github.com/hansmi/dossier/pkg/geometry"
)

func nodeNames(nodes []*sketchNode) []string {
	var result []string

	for _, n := range nodes {
		result = append(result, n.name)
	}

	return result
}

// Node search areas are defined as rectangles. Rectangle coordinates can be
// defined relative to other nodes. To find their absolute position the other
// nodes need to be found first. This function walks the node dependencies and
// returns the order in which the nodes need to be evaluated.
func determineNodeOrder(nodes []*sketchNode) ([]int, error) {
	type indexed struct {
		index int
		node  *sketchNode
	}

	byName := map[string]indexed{}

	for idx, n := range nodes {
		if first, ok := byName[n.name]; ok {
			return nil, fmt.Errorf("%w: multiple nodes with name %q, first occurrence at index %d",
				sketcherror.ErrBadConfig, n.name, first.index)
		}

		byName[n.name] = indexed{
			index: idx,
			node:  n,
		}
	}

	var visit func(indexed) error

	visited := make([]bool, len(nodes))
	stack := make([]*sketchNode, 0, len(nodes))
	result := make([]int, 0, len(nodes))

	visit = func(cur indexed) error {
		// Not the most efficient way of checking for an existing item, but
		// most chains are very short and maintaining an additional data
		// structure for faster lookups, e.g. a map, isn't worth it.
		foundOnStack := slices.Contains(stack, cur.node)

		stack = append(stack, cur.node)
		defer func() {
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]
		}()

		if foundOnStack {
			names := nodeNames(stack)
			return fmt.Errorf("%w: recursive node reference: %s",
				sketcherror.ErrBadConfig, strings.Join(names, " \u2192 "))
		}

		if visited[cur.index] {
			return nil
		}

		visited[cur.index] = true

		for _, area := range cur.node.searchAreas {
			for _, i := range area.RequiredNodeFeatures() {
				other, ok := byName[i.NodeName()]
				if !ok {
					return fmt.Errorf("%w: node %q: referenced node %q not found", sketcherror.ErrBadConfig, cur.node.name, i.NodeName())
				}

				if _, err := other.node.featurePosition(geometry.Rect{}, i.Feature()); err != nil {
					return fmt.Errorf("%w: node %q: %w", sketcherror.ErrBadConfig, cur.node.name, err)
				}

				if err := visit(other); err != nil {
					return err
				}
			}
		}

		result = append(result, cur.index)

		return nil
	}

	for idx, n := range nodes {
		if err := visit(indexed{index: idx, node: n}); err != nil {
			return nil, err
		}
	}

	return result, nil
}
