package frontier

type Frontier struct {
	DFSGraph
	visited map[string]bool
}
