// A peculiar, experimental neural network model
// Intended to be memory efficient -- not needing model compression -- and hence
// able to continue learning
// Changlog
// AJP 26Nov16 First real attempt

package noggin

// XYPosition is for the two large spatial dimensions; conserve bits
type XYPosition int16

// XYSize defines the number of cells in a direction (e.g. row or cols in a grid)
type XYSize int16

// Depth is for the layer postion, orth. to XY; conserve bits
type Depth int8

// DendriteID indexes into a central array to conserve space
type DendriteID int32

// PotentialID indexes into an array of chemical potentials to conserve space
type PotentialID int32

// AnActivation is how an activation level is represented; conserve bits
type AnActivation int8

// Potential is a chemical potential
type Potential struct {
	Activation AnActivation  // Current activation level
	X, Y       XYPosition    // Spatial position
	nearby     []PotentialID // chemical potentials close to this one
}

// Pot supplies the data required to follow chemical potential
type Pot interface {
	Init(X, Y XYPosition)
	Pot() Potential
}

// Cell is just a potential well
type Cell Potential

// Pot is the chemical potential
func (c Cell) Pot() Potential {
	return Potential(c)
}

// Init the cell
func (c Cell) Init(x, y XYPosition) {
	c.X = x
	c.Y = y
}

// Neuron , unlike a basic cell, has connections
type Neuron struct {
	Cell                   // a neuron is a cell plus ...
	Axon      AnActivation // Current input from axon (used during training)
	Dendrites []DendriteID // Dendrites begining at this neuron
}

// Pot is just that of the cell
func (n Neuron) Pot() Potential {
	return n.Cell.Pot()
}

// Init is just that of the cell
func (n Neuron) Init(x, y XYPosition) {
	n.Cell.Init(x, y)
}

// DendriteIs says what a dendrite is currently doing; conserve bits
type DendriteIs int8

const (
	// Growing dendrites are looking for a neuron to attach to
	Growing DendriteIs = iota
	// Attached dendrites have found a neuron and linked to it
	Attached DendriteIs = iota
)

// Dendrite represents a connection from a neuron to a chemical potential
type Dendrite struct {
	Doing  DendriteIs   // what is it up to
	To     PotentialID  // Neuron/cell it is closest to, or touching
	Weight AnActivation // How important is this dendrite?; conserve bits
}

// Layer is a flat sheet of cells
type Layer struct {
	Name  string // For display porpoises
	Z     Depth  // Spatial position of this layer
	Cells []Pot  // All the neurons or cells in this layer live here
}

// NeuralLayer is a flat sheet of neurons
type NeuralLayer struct {
	Layer                // it is a layer of cells, plus ...
	Dendrites []Dendrite // Any dendrites live here
	ToZ       *Layer     // The layer the dendrites connect to
}

// Brain is a complete NN, made of layers
type Brain struct {
	Name   string
	Layers map[Depth]Layer
}

// AddGrid places a grid of cells or neurons geometrically and then into the layer array of them
func (lay Layer) AddGrid(nrows, ncols XYSize, spacing XYPosition) {
	roff := ((XYPosition(nrows) - 1) * spacing) >> 1
	coff := ((XYPosition(ncols) - 1) * spacing) >> 1
	nid := 0
	for r := XYSize(0); r < nrows; r++ {
		for c := XYSize(0); c < ncols; c++ {
			x := (XYPosition(r) * spacing) - roff
			y := (XYPosition(r) * spacing) - coff
			lay.Cells[nid].Init(x, y)
			nid++
		}
	}
}

// NeuronGrid makes a layer a grid of s, aligned on 0,0
func (lay Layer) NeuronGrid(nrows, ncols XYSize, spacing XYPosition) {
	//lay.Cells = make([]Neuron, nrows*ncols) // First make the array of the correct type of cell
	lay.AddGrid(nrows, ncols, spacing) // then place them
}

// CellGrid makes a layer a grid of s, aligned on 0,0
func (lay Layer) CellGrid(nrows, ncols XYSize, spacing XYPosition) {
	//lay.Cells = []Pot(make([]Cell, nrows*ncols)) // First make the array of the correct type of cell
	lay.AddGrid(nrows, ncols, spacing) // then place them
}

// AddLayer adds a layer to a brain at a given depth
func (nog *Brain) AddLayer(lay Layer, z Depth) {
	nog.Layers[z] = lay
}
