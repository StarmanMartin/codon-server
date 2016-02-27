package main

// FindIfCircular finds a cirkl inside the grph if exists
func (c *CodonGraph) FindIfCircular() (isCyclingCode bool) {
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

	c.MaxPath = c.findPath(connections)
	return true
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

// getFromNucletideByConnectionIdx gets
// target Nucleotide index from connection array
// 0. A<-B, 6.  A->B
// 1. A<-C, 7.  A->C
// 2. A<-D, 8.  A->D
// 3. B<-C, 9.  B->C
// 4. B<-D, 10. B->D
// 5. C<-D, 11. C->D
func getToNucletideByConnectionIdx(dir, idx int) int {
	switch idx + 6*dir {
	case 0, 1, 2:
		return 0
	case 3, 4, 6:
		return 2
	case 5, 7, 9:
		return 4
	case 8, 10, 11:
		return 6
	}

	return -1
}

// getFromNucletideByConnectionIdx gets
// origan Nucleotide index from connection array
// 0. A<-B, 6.  A->B
// 1. A<-C, 7.  A->C
// 2. A<-D, 8.  A->D
// 3. B<-C, 9.  B->C
// 4. B<-D, 10. B->D
// 5. C<-D, 11. C->D
func getFromNucletideByConnectionIdx(dir, idx int) int {
	switch idx + 6*dir {
	case 6, 7, 8:
		return 0
	case 9, 10, 0:
		return 2
	case 11, 1, 3:
		return 4
	case 2, 4, 5:
		return 6
	}

	return -1
}

func (c *CodonGraph) findPath(connections [2][6]bool) int {
	maxLength := 0
	for idx := range connections[0] {
		nucIndex := getFromNucletideByConnectionIdx(0, idx)
		incomming := 0
		if len(c.TetranucleotideNodes[nucIndex]) > 0 {
			incomming = 1
		}

		if temp := c.findFromRightPath(connections, idx, incomming); temp > maxLength {
			maxLength = temp
		}
        
		nucIndex = getFromNucletideByConnectionIdx(1, idx)
		incomming = 0
		if len(c.TetranucleotideNodes[nucIndex]) > 0 {
			incomming = 1
		}

		if temp := c.findFromLeftPath(connections, idx, incomming); temp > maxLength {
			maxLength = temp
		}
	}

	return maxLength
}

func (c *CodonGraph) findFromRightPath(connections [2][6]bool, idx, length int) int {
	if !connections[0][idx] {
		nucIndex := getFromNucletideByConnectionIdx(0, idx)
		if len(c.TetranucleotideNodes[nucIndex+1]) > 0 {
			return length + 1
		}
        
		return length
	}

	length += 2
	startLength := length
	var temp int

	if idx == 0 {
		if temp = c.findFromLeftPath(connections, 1, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromLeftPath(connections, 2, startLength); temp > length {
			length = temp
		}
	} else if idx == 1 {
		if temp = c.findFromLeftPath(connections, 0, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromLeftPath(connections, 2, startLength); temp > length {
			length = temp
		}
	} else if idx == 2 {
		if temp = c.findFromLeftPath(connections, 0, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromLeftPath(connections, 1, startLength); temp > length {
			length = temp
		}
	} else if idx == 3 {
		if temp = c.findFromLeftPath(connections, 4, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromRightPath(connections, 0, startLength); temp > length {
			length = temp
		}
	} else if idx == 4 {
		if temp = c.findFromLeftPath(connections, 3, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromRightPath(connections, 0, startLength); temp > length {
			length = temp
		}
	} else if idx == 5 {
		if temp = c.findFromRightPath(connections, 1, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromRightPath(connections, 3, startLength); temp > length {
			length = temp
		}
	}
    
	nucIndex := getToNucletideByConnectionIdx(0, idx)
	if length == startLength && len(c.TetranucleotideNodes[nucIndex+1]) > 0 {
		return startLength + 1
	}

	return length
}

func (c *CodonGraph) findFromLeftPath(connections [2][6]bool, idx, length int) int {
	if !connections[1][idx] {
		nucIndex := getFromNucletideByConnectionIdx(1, idx)
		if len(c.TetranucleotideNodes[nucIndex+1]) > 0 {
			return length + 1
		}
		return length
	}

	length += 2
	startLength := length
	var temp int

	if idx == 0 {
		if temp = c.findFromLeftPath(connections, 3, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromLeftPath(connections, 4, startLength); temp > length {
			length = temp
		}
	} else if idx == 1 {
		if temp = c.findFromLeftPath(connections, 5, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromRightPath(connections, 3, startLength); temp > length {
			length = temp
		}
	} else if idx == 2 {
		if temp = c.findFromRightPath(connections, 4, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromRightPath(connections, 5, startLength); temp > length {
			length = temp
		}
	} else if idx == 3 {
		if temp = c.findFromLeftPath(connections, 5, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromRightPath(connections, 1, startLength); temp > length {
			length = temp
		}
	} else if idx == 4 {
		if temp = c.findFromRightPath(connections, 5, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromRightPath(connections, 2, startLength); temp > length {
			length = temp
		}
	} else if idx == 5 {
		if temp = c.findFromRightPath(connections, 4, startLength); temp > length {
			length = temp
		}
		if temp = c.findFromRightPath(connections, 2, startLength); temp > length {
			length = temp
		}
	}
    
	nucIndex := getToNucletideByConnectionIdx(1, idx)
	if length == startLength && len(c.TetranucleotideNodes[nucIndex+1]) > 0 {
		return startLength + 1
	}

	return length
}
