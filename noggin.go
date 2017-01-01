// A peculiar, experimental neural network model
// Changlog
// AJP 28Dec16 More basic work
// AJP 26Nov16 First real attempt

package noggin

// Spatial defines a dimension in space
type Spatial int16

// Vec defines a spatial vector with all dimensions
type Vec struct {
	x, y, z Spatial
}

// SubVec defines a spatial vector within a tissue (without the final dimension)
type SubVec struct {
	x, y Spatial
}

// Size defines a number of cells in a spatial direction (e.g. row or cols in a grid)
type Size int16

// CellID indexes into an array of electrochemical potentials; conserve bits
type CellID int32

// PotWeight is a weighting of a Potential
type PotWeight int16

// Activation is how an activation level is represented; conserve bits
type Activation int8

// ActWeight is a weighting of an activation
type ActWeight int8

// Cell has a chemical potential, possibly other things
type Cell struct {
	pos SubVec     // Spatial position
	act Activation // Current activation level
}

// Attracter supplies the data required to follow chemical potential
type Attracter interface {
	Attract(want Activation, at SubVec) (w PotWeight, dir SubVec)
}

// Attract is the attracter implementation
func (c Cell) Attract(want Activation, at SubVec) (w PotWeight, dir SubVec) {
	return w, c.pos
}

// DendWeight is a weighting for a dendrite
type DendWeight int16

// MaxDends is the most dendrites a neuron may have
const MaxDends = 100

// DendriteIs says what a dendrite is currently doing; conserve bits
type DendriteIs int8

const (
	growing  DendriteIs = iota // Growing dendrites are looking for a neuron to attach to
	attached                   // Attached dendrites have found a neuron and linked to it
)

// Dendrite represents a connection from a neuron to another cell, it is part of the cell from which it emanates
type Dendrite struct {
	To     CellID     // Neuron/cell it is closest to, or touching
	Weight ActWeight  // How important is this dendrite?; conserve bits
	Doing  DendriteIs // what it is doing
}

// Neuron ,unlike a basic cell, has connections -- dendrites
type Neuron struct {
	Cell                      // a neuron is a cell plus ...
	Axon   Activation         // Current input from axon (used during training)
	nDends int16              // Number of dendrites
	Dends  [MaxDends]Dendrite // Dendrites begining at this neuron -- note: pre-allocated ARRAY
}

// ThingType -- whether neuron or simple cell etc.
type ThingType int8

const (
	simpleCell   ThingType = iota // simpleCell is just a humble cell e.g. photoreceptor
	neuronCell                    // neuronCell is a full-fledged neuron
	dendriteCell                  // dendriteCell is a dendrite (not really a separate cell)
)

// Layer is a flat sheet of cells
type Layer struct {
	Name  string  // For display porpoises
	Z     Spatial // Spatial position of this layer
	ToZ   *Layer  // The layer the dendrites connect to
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
	Name    string   // For display porpoises
	Z       Spatial  // Spatial position of this layer
	ToZ     *Layer   // The layer the dendrites connect to
	nNeur   CellID   // Neuron count
	Neurons []Neuron // All the neurons in this layer
}

// Next returns to index for the next cell to be used
func (lay NeuralLayer) Next() CellID {
	lay.nNeur++
	return lay.nNeur - 1
}

// Init FULLY allocates ALL the structures for a Layer
func (lay NeuralLayer) Init(name string, maxneur int) {
	lay.Neurons = make([]Neuron, maxneur, maxneur)
}

// Nog is a set of layers -- an organoid
type Nog struct {
	Name   string
	Layers map[Spatial]*Layer
}

// Init FULLY allocates ALL the structures for a Noggin
func (ng Nog) Init(name string, mcell, mneur, mdend int) {
	ng.Name = name
}

// AddGrid places a grid of cells or neurons geometrically and then into the layer array of them
func (lay Layer) AddGrid(nrows, ncols Size, spacing Spatial) {
	var newid CellID
	roff := ((Spatial(nrows) - 1) * spacing) >> 1
	coff := ((Spatial(ncols) - 1) * spacing) >> 1
	for r := Size(0); r < nrows; r++ {
		for c := Size(0); c < ncols; c++ {
			x := (Spatial(r) * spacing) - roff
			y := (Spatial(r) * spacing) - coff
			newid = lay.Next()
			lay.Cells[newid].pos = SubVec{x, y}
		}
	}
}

// NeuronGrid makes a layer a grid of s, aligned on 0,0
func (lay Layer) NeuronGrid(ng Nog, nrows, ncols Size, spacing Spatial) {
	//lay.Cells = make([]Neuron, nrows*ncols) // First make the array of the correct type of cell
	lay.AddGrid(nrows, ncols, spacing) // then place them
}

// CellGrid makes a layer a grid of s, aligned on 0,0
func (lay Layer) CellGrid(ng Nog, nrows, ncols Size, spacing Spatial) {
	//lay.Cells = []Pot(make([]Cell, nrows*ncols)) // First make the array of the correct type of cell
	lay.AddGrid(nrows, ncols, spacing) // then place them
}

// AddLayer adds a layer to a noggin at a given depth
func (ng *Nog) AddLayer(lay Layer, z Spatial) {
	ng.Layers[z] = &lay
}
