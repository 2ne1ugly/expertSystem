/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   parse.go                                           :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: mchi <mchi@student.42.fr>                  +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2019/04/19 15:21:07 by mchi              #+#    #+#             */
/*   Updated: 2019/04/19 15:21:07 by mchi             ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

//CheckError : three lines to one.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

//Operator : enum for operator types
type Operator byte

//types of operation
const (
	And Operator = '+'
	Or  Operator = '|'
	Xor Operator = '^'
)

//TriBool : enum for three results
type TriBool int

//types of result
const (
	False   TriBool = 0
	True    TriBool = 1
	Unknown TriBool = 2
)

//LogicToken : tokens for tree
type LogicToken interface {
	AssignToken(str string)
	GetToken() byte
}

type symbol struct {
	name byte
	not  bool
}

type logicNode struct {
	token  LogicToken
	left   *logicNode
	right  *logicNode
	parent *logicNode
}

type logicTree struct {
	root *logicNode
}

type proposition struct {
	expr  logicTree
	concl logicTree
}

func (sym *symbol) AssignToken(str string) {
	if len(str) > 2 {
		log.Fatalf("following is not a proper symbol: %s\n", str)
	}
	var char byte
	if len(str) == 2 {
		if str[0] != '!' {
			log.Fatalf("following has length of 2 but do not start with '!': %s\n", str)
		}
		sym.not = true
		char = str[1]
	} else {
		char = str[0]
	}
	if !(('a' <= char && char <= 'z') ||
		('A' <= char && char <= 'Z')) {
		log.Fatalf("not an alphabet: %c\n", char)
	}
	sym.name = char
}

//AssignToken : assign values from the token
func (op *Operator) AssignToken(str string) {
	if len(str) != 1 {
		log.Fatalf("following is not a proper operator: %s\n", str)
	}
	switch str[0] {
	case '+':
		*op = And
	case '|':
		*op = Or
	case '^':
		*op = Xor
	default:
		log.Fatalf("following is not a proper operator: %s\n", str)
	}
}

//GetToken : get token char.
func (op *Operator) GetToken() byte {
	return byte(*op)
}

//GetToken : get token char.
func (sym *symbol) GetToken() byte {
	return byte(sym.name)
}

func parseNode(parent **logicNode, tokens *[]string) {
	if len(*tokens) == 0 {
		log.Fatalf("expected value but suddenly ended\n")
	} else if len(*tokens) == 1 {
		log.Fatalf("the incomplete statement ending with: %s\n", (*tokens)[0])
	}
	*parent = &logicNode{}
	(*parent).token = new(Operator)
	(*parent).token.AssignToken((*tokens)[0])
	*tokens = (*tokens)[1:]
	if (*parent).token.GetToken() == '+' {
		(*parent).right = &logicNode{}
		(*parent).right.token = new(symbol)
		(*parent).right.token.AssignToken((*tokens)[0])
		(*parent).right.parent = (*parent)
		*tokens = (*tokens)[1:]
	} else {
		parseTree(&(*parent).right, tokens)
		(*parent).right.parent = (*parent)
	}
}

func parseTree(parent **logicNode, tokens *[]string) {
	if len(*tokens) == 0 {
		log.Fatalf("no tokens were found\n")
	}
	*parent = &logicNode{}
	(*parent).token = &symbol{}
	(*parent).token.AssignToken((*tokens)[0])
	*tokens = (*tokens)[1:]
	for len(*tokens) != 0 {
		temp := (*parent).left
		parseNode(parent, tokens)
		(*parent).left = temp
		if temp != nil {
			temp.parent = (*parent)
		}
	}
}

//Input : input files
type Input struct {
	rules []proposition
	facts []symbol
	query []symbol
}

//ParseFile : runs read file and parses.
func ParseFile(path string) Input {
	file, err := ioutil.ReadFile(path)
	CheckError(err)
	deleteComments := regexp.MustCompile("#.*\n")
	deleteSpaces := regexp.MustCompile(" |\t")
	file = deleteComments.ReplaceAll(file, []byte("\n"))
	file = deleteSpaces.ReplaceAll(file, []byte(""))
	data := strings.ReplaceAll(string(file), "\t", "")
	lines := strings.Split(string(data), "\n")
	var result Input
	for _, line := range lines {
		if line == "" || line == "\r" {
			continue
		} else if line[0] == '=' {
			for i := 1; i < len(line); i++ {
				result.facts = append(result.facts, symbol{line[i], false})
			}
		} else if line[0] == '?' {
			for i := 1; i < len(line); i++ {
				result.query = append(result.query, symbol{line[i], false})
			}
		} else {
			eq := strings.Split(line, "=>")
			if len(eq) != 2 {
				log.Fatalf("incorrect syntex : %s\n", line)
			}
			var rule proposition
			spearteByNodes := regexp.MustCompile("[+|^]|[^+|^]*").FindAllString
			tokens := spearteByNodes(eq[0], -1)
			parseTree(&rule.expr.root, &tokens)
			tokens = spearteByNodes(eq[1], -1)
			parseTree(&rule.concl.root, &tokens)
			result.rules = append(result.rules, rule)
			print(line, "\n")
		}
	}
	return result
}
