package driver

// A Collection of drivers
type Collection []*Driver

// Add will add a new driver to the collection, making sure that
// it is not duplicated.
func (c *Collection) Add(d *Driver) {
	if d == nil {
		return
	}

	for i, od := range *c {
		if od.Equals(d) {
			(*c)[i] = d
			return
		}
	}

	*c = append(*c, d)
}

// Remove will remove a driver from the collection
func (c *Collection) Remove(d *Driver) {
	if d == nil {
		return
	}

	for i, od := range *c {
		if od.Equals(d) {
			*c = append((*c)[:i], (*c)[i+1:]...)
			return
		}
	}
}
