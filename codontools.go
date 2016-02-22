package main

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ruleReg *regexp.Regexp
)

func init() {
	ruleReg = regexp.MustCompile(`([AUCG])\<-\>([AUCG])`)
}

func preparePermutateRule(rules string) ([]string, error) {
	ruleRegList := ruleReg.FindAllStringSubmatch(rules, -1)
	if len(ruleRegList) > 2 || len(ruleRegList) <= 0 {
		return nil, errors.New("Rule not correct")
	}

	ruleList := make([]string, len(ruleRegList)*2)
	for idx, sRule := range ruleRegList {
		ruleList[idx*2] = sRule[1]
		ruleList[idx*2+1] = sRule[2]
	}

	return ruleList, nil
}

// PermutateCodons permutates a list of codons by a rule.
// Rule sample: A<->U;G<->C
func PermutateCodons(codonList []string, rules string) ([]string, error) {
	ruleList, err := preparePermutateRule(rules)
	if err != nil {
		return nil, err
	}

	var lastBase string
	for idx, rule := range ruleList {
		letter := "L"
		if idx%2 == 0 {
			for cIdex, codon := range codonList {
				codonList[cIdex] = strings.Replace(codon, rule, letter, -1)
			}

			lastBase = rule
		} else {
			for cIdex, codon := range codonList {
				codon = strings.Replace(codon, rule, lastBase, -1)
				codonList[cIdex] = strings.Replace(codon, letter, rule, -1)
			}

		}
	}

	return codonList, nil
}
