// A peculiar, experimental neural network model
// Changlog
// AJP 28Dec16 More basic work
// AJP 26Nov16 First real attempt

package noggin

// XY is for the two large spatial dimensions; conserve bits
type XY int16

// Size defines the number of cells in a direction (e.g. row or cols in a grid)
type Size int16

// Depth is for the layer postion, orth. to XY
type Depth int16

// DendID indexes into a central array of them; conserve bits
//type DendID int32

// DendWeight is a weighting for a dendrite
type DendWeight int16

// CellID indexes into an array of chemical potentials; conserve bits
type CellID int32

// PotWeight is a weighting of a Potential
type PotWeight int16

// AnActivation is how an activation level is represented; conserve bits
type AnActivation int8

// ActWeight is a weighting of an activation
type ActWeight int8

// Cell has a chemical potential, possibly other things
type Cell struct {
	Activation AnActivation // Current activation level
	X, Y       XY           // Spatial position
}

// Attracter supplies the data required to follow chemical potential
type Attracter interface {
	Pot(x XY, y XY) (w PotWeight, dx XY, dy XY)
}

// Pot provides info to follow the potential
func (c Cell) Pot(x XY, y XY) (w PotWeight, dx XY, dy XY) {
	return w, x, y
}

// Neuron , unlike a basic cell, has connections
type Neuron struct {
	Cell                   // a neuron is a cell plus ...
	Axon      AnActivation // Current input from axon (used during training)
	Dendrites []CellID     // Dendrites begining at this neuron
}

// ThingType -- whether neuron or simple cell etc.
type ThingType int8

const (
	// SimpleCell is just a humble cell e.g. photoreceptor
	SimpleCell ThingType = iota
	// NeuronCell is a full-fledged neuron
	NeuronCell
	// DendriteCell is a dendrite (not really a separate cell)
	DendriteCell
)

// DendriteIs says what a dendrite is currently doing; conserve bits
type DendriteIs int8

const (
	// Growing dendrites are looking for a neuron to attach to
	Growing DendriteIs = iota
	// Attached dendrites have found a neuron and linked to it
	Attached
)

// Dendrite represents a connection from a neuron to another cell
type Dendrite struct {
	Doing  DendriteIs // what is it up to
	To     CellID     // Neuron/cell it is closest to, or touching
	Weight ActWeight  // How important is this dendrite?; conserve bits
}

// Layer is a flat sheet of cells
type Layer struct {
	Name  string // For display porpoises
	Z     Depth  // Spatial position of this layer
	ToZ   *Layer // The layer the dendrites connect to
	NCell CellID
	Cells []Cell // All the neurons or cells in this layer
}

// Init FULLY allocates ALL the structures for a Layer
func (lay Layer) Init(name string, mcell int) {
	lay.Cells = make([]Cell, mcell, mcell)
}

// Next returns to index for the next cell to be used
func (lay Layer) Next() CellID {
	lay.NCell++
	return lay.NCell - 1
}

// Nexter has an array of cells which it doles out
type Nexter interface {
	Next() CellID
}

// Activer is a cell that computes an activity - i.e. forward pass
type Activer interface {
	Active()
}

// Trainer is a cell that can be trained somehow - i.e. back prop
type Trainer interface {
	Train()
}

// NeuralLayer is a flat sheet of neurons
type NeuralLayer struct {
	Name    string // For display porpoises
	Z       Depth  // Spatial position of this layer
	ToZ     *Layer // The layer the dendrites connect to
	nNeur   CellID
	nDend   CellID
	Neurons []Neuron
	Dends   []Dendrite
}

// Next returns to index for the next cell to be used
func (lay NeuralLayer) Next() CellID {
	lay.nNeur++
	return lay.nNeur - 1
}

// Init FULLY allocates ALL the structures for a Layer
func (lay NeuralLayer) Init(name string, mneur int, mdend int) {
	lay.Neurons = make([]Neuron, mneur, mneur)
	lay.Dends = make([]Dendrite, mdend, mdend)
}

// Nog is 3D set of layers -- an organoid
type Nog struct {
	Name   string
	Layers map[Depth]*Layer
}

// Init FULLY allocates ALL the structures for a Noggin
func (ng Nog) Init(name string, mcell, mneur, mdend int) {
	ng.Name = name
}

// AddGrid places a grid of cells or neurons geometrically and then into the layer array of them
func (lay Layer) AddGrid(nrows, ncols Size, spacing XY) {
	var newid CellID
	roff := ((XY(nrows) - 1) * spacing) >> 1
	coff := ((XY(ncols) - 1) * spacing) >> 1
	for r := Size(0); r < nrows; r++ {
		for c := Size(0); c < ncols; c++ {
			x := (XY(r) * spacing) - roff
			y := (XY(r) * spacing) - coff
			newid = lay.Next()
			lay.Cells[newid].X, lay.Cells[newid].Y = x, y
		}
	}
}

// NeuronGrid makes a layer a grid of s, aligned on 0,0
func (lay Layer) NeuronGrid(ng Nog, nrows, ncols Size, spacing XY) {
	//lay.Cells = make([]Neuron, nrows*ncols) // First make the array of the correct type of cell
	lay.AddGrid(nrows, ncols, spacing) // then place them
}

// CellGrid makes a layer a grid of s, aligned on 0,0
func (lay Layer) CellGrid(ng Nog, nrows, ncols Size, spacing XY) {
	//lay.Cells = []Pot(make([]Cell, nrows*ncols)) // First make the array of the correct type of cell
	lay.AddGrid(nrows, ncols, spacing) // then place them
}

// AddLayer adds a layer to a noggin at a given depth
func (ng *Nog) AddLayer(lay Layer, z Depth) {
	ng.Layers[z] = &lay
}
