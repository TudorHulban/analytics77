package analytics

func (m *MetaActive[T]) setOccupied(ix int) {
	bit := uint32(1 << ix)

	for {
		oldValue := m.occupied.Load()
		newValue := oldValue | bit

		if m.occupied.CompareAndSwap(oldValue, newValue) {
			return
		}
	}
}
