package main

// FindIfCircular finds a cirkl inside the grph if exists
func (c *CodonGraph) FindIfCircular() {
	var connections [2][6]bool
	subCounter, tempCounter := 0, 0
	for i := 0; i < 4; i++ {
		tempCounter = subCounter
		for _, val := range c.TetranucleotideNodes[i*2] {
			if _, has := indexOfInt(c.TetranucleotideNodes[i*2+1], val); !has {
				c.CyclingIndex = 1
				return
			}
		}

		for baseIn := 0; baseIn < 2; baseIn++ {
			subCounter = tempCounter
			baseIndex := i*2 + baseIn
			for sub := i + 1; sub < 4; sub++ {
				subIn := 1 - baseIn
				subIndex := sub*2 + subIn
				for _, val := range c.TetranucleotideNodes[baseIndex] {
					if _, has := indexOfInt(c.TetranucleotideNodes[subIndex], val); !has {
						connections[baseIn][subCounter] = true
					}
				}

				subCounter++
			}
		}
	}
    
	if isCycling := c.firstStage(connections); isCycling {
		c.CyclingIndex = 2
        return
    }
    
    if isCycling := c.secondStage(connections); isCycling {
		c.CyclingIndex = 3
        return
    }
    
    if isCycling := c.thirdStage(connections); isCycling {
		c.CyclingIndex = 4
        return
    }
}

func (c *CodonGraph) firstStage(connections [2][6]bool) (isCycling bool) {
	for index, val := range connections[0] {
		if val && connections[1][index] {
            isCycling = true
			return
		}
	}

	return
}

func (c *CodonGraph) secondStage(connections [2][6]bool) (isCycling bool) {
	for index, val := range connections {
		switch {
		case val[0] && val[3] && connections[1-index][1]:
			fallthrough
		case val[0] && val[4] && connections[1-index][2]:
			fallthrough
		case val[3] && val[5] && connections[1-index][4]:
			return true
		}
	}

	return
}

func (c *CodonGraph) thirdStage(connections [2][6]bool) (isCycling bool) {
	for index, val := range connections {
		switch {
		case val[0] && val[3] && val[5] && connections[1-index][2]:
			fallthrough
		case val[0] && val[4] && connections[1-index][1] && connections[1-index][5]:
			return true
		}
	}

	return
}