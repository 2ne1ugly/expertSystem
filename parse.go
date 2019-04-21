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

func parseNode(parent **LogicNode, tokens *[]string) {
	if len(*tokens) == 0 {
		log.Fatalf("expected value but suddenly ended\n")
	} else if len(*tokens) == 1 {
		log.Fatalf("the incomplete Statement ending with: %s\n", (*tokens)[0])
	}
	*parent = &LogicNode{}
	(*parent).token = new(Operator)
	(*parent).token.AssignToken((*tokens)[0])
	*tokens = (*tokens)[1:]
	if (*parent).token.GetToken() == '+' {
		(*parent).right = &LogicNode{}
		(*parent).right.token = new(Symbol)
		(*parent).right.token.AssignToken((*tokens)[0])
		(*parent).right.parent = (*parent)
		*tokens = (*tokens)[1:]
	} else {
		parseTree(&(*parent).right, tokens)
		(*parent).right.parent = (*parent)
	}
}

func parseTree(parent **LogicNode, tokens *[]string) {
	if len(*tokens) == 0 {
		log.Fatalf("no tokens were found\n")
	}
	*parent = &LogicNode{}
	(*parent).token = new(Symbol)
	(*parent).token.AssignToken((*tokens)[0])
	*tokens = (*tokens)[1:]
	for len(*tokens) != 0 {
		temp := *parent
		parseNode(parent, tokens)
		(*parent).left = temp
		if temp != nil {
			temp.parent = (*parent)
		}
	}
}

func copyLogicNode(node *LogicNode) *LogicNode {
	newNode := &LogicNode{}
	newNode.token = node.token.CopyToken()
	if node.left != nil {
		newNode.left = copyLogicNode(node.left)
		newNode.left.parent = newNode
	}
	if node.right != nil {
		newNode.right = copyLogicNode(node.right)
		newNode.right.parent = newNode
	}
	return newNode
}

//Input : input files
type Input struct {
	rules    []Statement
	facts    []Symbol
	query    []Symbol
	refSheet map[byte][]*LogicNode
}

func separateStringsToTokens(str string) []string {
	tokens := make([]string, 0)
	for i := range str {
		if str[i] == '+' || str[i] == '|' || str[i] == '^' || str[i] == '!' || ('A' <= str[i] && str[i] <= 'Z') {
			tokens = append(tokens, string(str[i]))
		} else if str[i] == '(' {
			start := i
			c := 0
			for i != len(str) {
				if str[i] == '(' {
					c++
				} else if str[i] == ')' {
					c--
				}
				i++
				if c == 0 {
					break
				}
			}
			tokens = append(tokens, str[start:i])
		}
	}
	return tokens
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
				result.facts = append(result.facts, Symbol(line[i]))
			}
		} else if line[0] == '?' {
			for i := 1; i < len(line); i++ {
				result.query = append(result.query, Symbol(line[i]))
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
			var prop Statement
			var contr Statement
			tokens := separateStringsToTokens(eq[0])
			parseTree(&prop.expr.root, &tokens)
			tokens = separateStringsToTokens(eq[1])
			parseTree(&prop.concl.root, &tokens)
			result.rules = append(result.rules, prop)
			//create contrapositive by copying previous proposition.
			contr.expr.root = copyLogicNode(prop.concl.root)
			contr.concl.root = copyLogicNode(prop.expr.root)
			result.rules = append(result.rules, contr)
			//additional work for only-if cases.
			if strings.Contains(line, "<=>") {
				var prop2 Statement
				var contr2 Statement
				prop2.expr.root = copyLogicNode(prop.concl.root)
				prop2.concl.root = copyLogicNode(prop.expr.root)
				contr2.expr.root = copyLogicNode(contr.concl.root)
				contr2.concl.root = copyLogicNode(contr.expr.root)
				result.rules = append(result.rules, prop2)
				result.rules = append(result.rules, contr2)
			}
			print(line, "\n")
		}
	}
	return result
}
