/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   logicTree.go                                       :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: mchi <mchi@student.42.fr>                  +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2019/04/21 00:14:58 by mchi              #+#    #+#             */
/*   Updated: 2019/04/21 00:14:58 by mchi             ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import "log"

//Operator : enum for operator types
type Operator byte

//Symbol : tokens for symbols
type Symbol byte

//types of operation
const (
	And Operator = '+'
	Or  Operator = '|'
	Xor Operator = '^'
	Not Operator = '!'
)

//LogicToken : tokens for tree
type LogicToken interface {
	AssignToken(str string)
	GetToken() byte
	CopyToken() LogicToken
}

//LogicNode : nodingForLogic
type LogicNode struct {
	token  LogicToken
	left   *LogicNode
	right  *LogicNode
	parent *LogicNode
}

//LogicTree :
type LogicTree struct {
	root *LogicNode
}

//Statement : rules.
type Statement struct {
	expr  LogicTree
	concl LogicTree
}

//AssignToken : assign values from the token
func (sym *Symbol) AssignToken(str string) {
	if len(str) > 1 {
		log.Fatalf("following is not a proper Symbol: %s\n", str)
	}
	if !('A' <= str[0] && str[0] <= 'Z') {
		log.Fatalf("not an alphabet: %c\n", str[0])
	}
	*sym = Symbol(str[0])
}

//AssignToken : assign values from the token
func (op *Operator) AssignToken(str string) {
	if len(str) != 1 {
		log.Fatalf("following is not a proper operator: %s\n", str)
	}
	if !(str[0] == '+' || str[0] == '|' || str[0] == '^' || str[0] == '!') {
		log.Fatalf("not an operator: %c\n", str[0])
	}
	*op = Operator(str[0])
}

//GetToken : get token char.
func (op *Operator) GetToken() byte {
	return byte(*op)
}

//GetToken : get token char.
func (sym *Symbol) GetToken() byte {
	return byte(*sym)
}

//CopyToken : get token char.
func (op *Operator) CopyToken() LogicToken {
	newOp := *op
	return &newOp
}

//CopyToken : get token char.
func (sym *Symbol) CopyToken() LogicToken {
	newSym := *sym
	return &newSym
}
