package markov

func (m *Markov) load() {
	states := initialize(m.markovStg.File)
	m.states = states
	m.ready = true
}
