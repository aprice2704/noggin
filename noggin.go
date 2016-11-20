package noggin

// A Neuron is a simple type
type Neuron struct {
	Activation float64
	inputs     float64 // Internal total weight of incoming activations
	Synapses   []*Synapse
}

// Sense inputs
func (n *Neuron) Sense() {
	n.inputs = 0.0
	for _, s := range n.Synapses {
		n.inputs += s.Target.inputs * s.Weight
	}
}

// A Synapse connects two neurons
type Synapse struct {
	Weight float64
	base   *Neuron
	Target *Neuron
}

// SetTarget changes the target neuron of this synapse
func (s Synapse) SetTarget(newt *Neuron) {
	s.Target = newt
}
