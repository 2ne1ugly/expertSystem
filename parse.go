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

//LogicToken : tokens for tree
type LogicToken interface {
	AssignToken(str string)
	GetToken() byte
	CopyToken() LogicToken
}

type symbol struct {
	name byte
	not  bool
}

//LogicNode : nodingForLogic
type LogicNode struct {
	token  LogicToken
	left   *LogicNode
	right  *LogicNode
	parent *LogicNode
	result TriBool
}

type logicTree struct {
	root *LogicNode
}

//StatementType : statement type.
type StatementType int

const (
	proposition    StatementType = 0
	contrapositive StatementType = 1
)

type statement struct {
	expr      logicTree
	concl     logicTree
	stateType StatementType
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

//CopyToken : get token char.
func (op *Operator) CopyToken() LogicToken {
	newOp := *op
	return &newOp
}

//CopyToken : get token char.
func (sym *symbol) CopyToken() LogicToken {
	newSym := *sym
	return &newSym
}

func parseNode(parent **LogicNode, tokens *[]string, ref *map[byte][]*LogicNode) {
	if len(*tokens) == 0 {
		log.Fatalf("expected value but suddenly ended\n")
	} else if len(*tokens) == 1 {
		log.Fatalf("the incomplete statement ending with: %s\n", (*tokens)[0])
	}
	*parent = &LogicNode{}
	(*parent).result = Unknown
	(*parent).token = new(Operator)
	(*parent).token.AssignToken((*tokens)[0])
	*tokens = (*tokens)[1:]
	if (*parent).token.GetToken() == '+' {
		(*parent).right = &LogicNode{}
		(*parent).right.result = Unknown
		(*parent).right.token = new(symbol)
		(*parent).right.token.AssignToken((*tokens)[0])
		(*ref)[(*parent).right.token.GetToken()] = append((*ref)[(*parent).right.token.GetToken()], (*parent).right)
		(*parent).right.parent = (*parent)
		*tokens = (*tokens)[1:]
	} else {
		parseTree(&(*parent).right, tokens, ref)
		(*parent).right.parent = (*parent)
	}
}

func parseTree(parent **LogicNode, tokens *[]string, ref *map[byte][]*LogicNode) {
	if len(*tokens) == 0 {
		log.Fatalf("no tokens were found\n")
	}
	*parent = &LogicNode{}
	(*parent).result = Unknown
	(*parent).token = &symbol{}
	(*parent).token.AssignToken((*tokens)[0])
	(*ref)[(*parent).token.GetToken()] = append((*ref)[(*parent).token.GetToken()], *parent)
	*tokens = (*tokens)[1:]
	for len(*tokens) != 0 {
		temp := *parent
		parseNode(parent, tokens, ref)
		(*parent).left = temp
		if temp != nil {
			temp.parent = (*parent)
		}
	}
}

func copyLogicNode(node *LogicNode, ref *map[byte][]*LogicNode) *LogicNode {
	newNode := &LogicNode{}
	newNode.result = Unknown
	newNode.token = node.token.CopyToken()
	if node.left != nil {
		newNode.left = copyLogicNode(node.left, ref)
		newNode.left.parent = newNode
	}
	if node.right != nil {
		newNode.right = copyLogicNode(node.right, ref)
		newNode.right.parent = newNode
	}
	if node.left == nil && node.right == nil {
		(*ref)[newNode.token.GetToken()] = append((*ref)[newNode.token.GetToken()], newNode)
	}
	return newNode
}

//Input : input files
type Input struct {
	rules    []statement
	facts    []symbol
	query    []symbol
	refSheet map[byte][]*LogicNode
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
	//ref sheets for each variables.
	var result Input
	result.refSheet = make(map[byte][]*LogicNode)
	for _, line := range lines {
		//skip empty lines. Takes Special actions on query and facts.
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
			var eq []string
			if strings.Contains(line, "<=>") {
				eq = strings.Split(line, "<=>")
				if len(eq) != 2 {
					log.Fatalf("incorrect syntex : %s\n", line)
				}
			} else if strings.Contains(line, "=>") {
				eq = strings.Split(line, "=>")
				if len(eq) != 2 {
					log.Fatalf("incorrect syntex : %s\n", line)
				}
			} else {
				log.Fatalf("incorrect syntex : %s\n", line)
			}
			var prop statement
			var contr statement
			spearteByNodes := regexp.MustCompile("[+|^]|[^+|^]*").FindAllString
			//split as expr, conlcusion.
			tokens := spearteByNodes(eq[0], -1)
			parseTree(&prop.expr.root, &tokens, &result.refSheet)
			tokens = spearteByNodes(eq[1], -1)
			parseTree(&prop.concl.root, &tokens, &result.refSheet)
			prop.stateType = proposition
			result.rules = append(result.rules, prop)
			//create contrapositive by copying previous proposition.
			contr.expr.root = copyLogicNode(prop.concl.root, &result.refSheet)
			contr.concl.root = copyLogicNode(prop.expr.root, &result.refSheet)
			contr.stateType = contrapositive
			result.rules = append(result.rules, contr)
			//additional work for only-if cases.
			if strings.Contains(line, "<=>") {
				var prop2 statement
				var contr2 statement
				prop2.expr.root = copyLogicNode(prop.concl.root, &result.refSheet)
				prop2.concl.root = copyLogicNode(prop.expr.root, &result.refSheet)
				prop2.stateType = proposition
				contr2.expr.root = copyLogicNode(contr.concl.root, &result.refSheet)
				contr2.concl.root = copyLogicNode(contr.expr.root, &result.refSheet)
				prop2.stateType = contrapositive
				result.rules = append(result.rules, prop2)
				result.rules = append(result.rules, contr2)
			}
			print(line, "\n")
		}
	}
	return result
}
