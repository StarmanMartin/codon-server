package main

import (
	"errors"
	"log"
	"regexp"
	"sort"
)

var (
	codonReg       *regexp.Regexp
	baseList       [4]string
	complementList map[rune]string
	moveList       [36][]int
)

func init() {
	codonReg = regexp.MustCompile(`^[AUGC]{3}$`)
	baseList = [...]string{"A", "U", "C", "G"}
	complementList = map[rune]string{
		[]rune(baseList[0])[0]: baseList[1],
		[]rune(baseList[1])[0]: baseList[0],
		[]rune(baseList[3])[0]: baseList[2],
		[]rune(baseList[2])[0]: baseList[3],
	}

	moveList[0] = []int{}
	moveList[1] = []int{0, 1, 0, 2}
	moveList[2] = []int{0, 2}
	moveList[3] = []int{0, 2}
	moveList[4] = []int{0, 1, 0, 2}
	moveList[5] = []int{}

	moveList[6] = []int{0, 1}
	moveList[7] = []int{}
	moveList[8] = []int{1, 2, 0, 2}
	moveList[9] = []int{0, 1}
	moveList[10] = []int{1, 2}
	moveList[11] = []int{0, 1}

	moveList[12] = []int{0, 1}
	moveList[13] = []int{1, 3, 0, 1}
	moveList[14] = []int{}
	moveList[15] = []int{0, 2}
	moveList[16] = []int{0, 3, 2, 3}
	moveList[17] = []int{0, 1}

	moveList[18] = []int{}
	moveList[19] = []int{0, 3}
	moveList[20] = []int{1, 3}
	moveList[21] = []int{}
	moveList[22] = []int{1, 2}
	moveList[23] = []int{}

	moveList[24] = []int{2, 3}
	moveList[25] = []int{2, 3, 0, 2}
	moveList[26] = []int{1, 2}
	moveList[27] = []int{2, 3, 0, 3}
	moveList[28] = []int{}
	moveList[29] = []int{2, 3}

	moveList[30] = []int{}
	moveList[31] = []int{1, 2, 2, 3}
	moveList[32] = []int{1, 3}
	moveList[33] = []int{1, 3}
	moveList[34] = []int{0, 1, 1, 3}
	moveList[35] = []int{}

}

type dinucleotideSorter struct {
	dinucleotide []string
	groupe       [][3]int
	sortIndex    int
}

// Swap is part of sort.Interface.
func (s dinucleotideSorter) Swap(i, j int) {
	s.dinucleotide[i], s.dinucleotide[j] = s.dinucleotide[j], s.dinucleotide[i]
	s.groupe[i], s.groupe[j] = s.groupe[j], s.groupe[i]
}

// Len is part of sort.Interface.
func (s dinucleotideSorter) Len() int {
	return len(s.dinucleotide)
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s dinucleotideSorter) Less(i, j int) bool {
	return s.groupe[i][s.sortIndex] < s.groupe[j][s.sortIndex]
}

// CodonGraph is a struct with has all nodes and edges in the
// the graph of a given list of codons
type CodonGraph struct {
	List                       []string
	DinucleotideNodes          []string
	TetranucleotideNodes       [8][]int
	CyclingIndex               int
	SelfComplementary          bool
	StrongNotSelfComplementary bool
	PropertyOne                bool
	PropertyTwo                bool
	Nucleotide                 []string
	MaxPath                    int
}

// OrderNodes orders the list of dinucleotide and nucleotide
// to build the best posible graph
func (c *CodonGraph) OrderNodes() {
	c.orderNucleotide()
	c.orderDinucleotide()
}

func (c *CodonGraph) orderNucleotide() {

	var marked [4]bool
	var position [6]int
	index := 0

	for index < len(c.DinucleotideNodes) {
		for idx, val := range c.TetranucleotideNodes {
			for _, indexVal := range val {
				if indexVal == index {
					marked[idx/2] = true
					continue
				}
			}

		}

		for idxMS, valMS := range marked {
			saveIdx := idxMS
			if idxMS > 0 {
				idxMS += idxMS + 1
			}

			if valMS {
				for idxME, valME := range marked[saveIdx+1:] {
					if valME {
						position[idxME+idxMS]++
					}
				}
			}

			marked[saveIdx] = false
		}

		index++
	}

	indexMax, maxVal, secIndex, secVal := -1, 0, -1, 0

	for idx, val := range position {
		tIdex := idx
		if maxVal < val && val != 0 {
			val, maxVal = maxVal, val
			tIdex, indexMax = indexMax, idx
		}

		if secVal < val && val != 0 {
			secVal = val
			secIndex = tIdex
		}
	}

	if secIndex == -1 {
		secIndex = indexMax
	}

	if indexMax > -1 {
		log.Println(indexMax*6+secIndex, moveList[indexMax*6+secIndex])
		c.swapNucleotide(moveList[indexMax*6+secIndex])
	}
}

func (c *CodonGraph) swapNucleotide(swapList []int) {
	for i := 0; i < len(swapList); i += 2 {
		fIdx, tIdex := swapList[i], swapList[i+1]

		c.Nucleotide[fIdx], c.Nucleotide[tIdex] = c.Nucleotide[tIdex], c.Nucleotide[fIdx]
		fIdx *= 2
		tIdex *= 2
		c.TetranucleotideNodes[fIdx], c.TetranucleotideNodes[tIdex] = c.TetranucleotideNodes[tIdex], c.TetranucleotideNodes[fIdx]
		fIdx++
		tIdex++
		c.TetranucleotideNodes[fIdx], c.TetranucleotideNodes[tIdex] = c.TetranucleotideNodes[tIdex], c.TetranucleotideNodes[fIdx]
	}
}

func (c *CodonGraph) orderDinucleotide() {
	swapNodes := &dinucleotideSorter{
		dinucleotide: c.DinucleotideNodes,
		groupe:       c.getBestOrder(),
		sortIndex:    1,
	}

	sort.Sort(swapNodes)
	swapNodes.sortIndex = 2
	border := (len(swapNodes.dinucleotide) + 1) / 2
	swapNodesRight := &dinucleotideSorter{
		dinucleotide: c.DinucleotideNodes[border:],
		groupe:       swapNodes.groupe[border:],
		sortIndex:    2,
	}
	sort.Sort(swapNodesRight)

	swapNodesLeft := &dinucleotideSorter{
		dinucleotide: c.DinucleotideNodes[:border],
		groupe:       swapNodes.groupe[:border],
		sortIndex:    2,
	}

	sort.Sort(swapNodesLeft)
	swopList := make([]int, len(c.DinucleotideNodes))
	for idx, val := range swapNodes.groupe {
		swopList[val[0]] = idx
	}

	for i := 0; i < 8; i++ {
		for idx := range c.TetranucleotideNodes[i] {
			oldVal := c.TetranucleotideNodes[i][idx]
			c.TetranucleotideNodes[i][idx] = swopList[oldVal]
		}
	}
}

func (c *CodonGraph) getBestOrder() [][3]int {
	rNodes := make([][3]int, len(c.DinucleotideNodes))
	sideValue := 0
	for dIndex := range c.DinucleotideNodes {
		hideValue, amount := 0.0, 0
		factorHide := [...]float64{2, 2, 1, 1, -1, -1, -2, -2}
		for left := 1; left < 8; left += 2 {
			for _, idx := range c.TetranucleotideNodes[left] {
				if idx == dIndex {
					amount++
					hideValue -= factorHide[left]
				}
			}
		}

		for right := 0; right < 8; right += 2 {
			for _, idx := range c.TetranucleotideNodes[right] {
				if idx == dIndex {
					amount++
					hideValue -= factorHide[right]
				}
			}
		}

		rNodes[dIndex][0] = dIndex
		//	rNodes[dIndex][1] = sideValue

		if amount < sideValue {
			// TODO Sort by amount of connections
			rNodes[dIndex][1] = -amount
			sideValue -= amount
		} else {
			rNodes[dIndex][1] = amount
			sideValue += amount
		}

		rNodes[dIndex][2] = int((hideValue / float64(amount)) * 100)

	}

	return rNodes
}

func indexOfInt(list []int, item int) (int, bool) {
	for idx, elm := range list {
		if item == elm {
			return idx, false
		}
	}

	return -1, true
}

func indexOf(list []string, item string) (int, bool) {
	for idx, elm := range list {
		if item == elm {
			return idx, false
		}
	}

	return -1, true
}

// NewCodonGraph returns a new uordered graph struct.
// All dinucleotide are in the DinucleotideNodes list.
// All tetranucleotide ar in the SNode with:
//   [0] -> A in; [1] -> A out
//   [2] -> U in; [3] -> U out
//   [4] -> C in; [5] -> C out
//   [6] -> G in; [7] -> G out
func NewCodonGraph(list []string) (*CodonGraph, error) {
	newGraph := &CodonGraph{}
	newGraph.List = list
	newGraph.DinucleotideNodes = make([]string, 0, 1)

	newGraph.Nucleotide = make([]string, 4)
	copy(newGraph.Nucleotide, baseList[:])

	for i := 0; i < 8; i++ {
		newGraph.TetranucleotideNodes[i] = make([]int, 0)
	}

	for _, codon := range list {
		if !codonReg.MatchString(codon) {
			return nil, errors.New("Wrong Codon")
		}

		index, isNew := indexOf(newGraph.DinucleotideNodes, codon[0:2])

		if isNew {
			index = len(newGraph.DinucleotideNodes)
			newGraph.DinucleotideNodes = append(newGraph.DinucleotideNodes, codon[0:2])
		}

		indexS, _ := indexOf(baseList[:], codon[2:3])
		indexS *= 2
		newGraph.TetranucleotideNodes[indexS] = append(newGraph.TetranucleotideNodes[indexS], index)

		index, isNew = indexOf(newGraph.DinucleotideNodes, codon[1:3])

		if isNew {
			index = len(newGraph.DinucleotideNodes)
			newGraph.DinucleotideNodes = append(newGraph.DinucleotideNodes, codon[1:3])
		}

		indexS, _ = indexOf(baseList[:], codon[0:1])
		indexS *= 2
		newGraph.TetranucleotideNodes[indexS+1] = append(newGraph.TetranucleotideNodes[indexS+1], index)

	}

	return newGraph, nil
}
