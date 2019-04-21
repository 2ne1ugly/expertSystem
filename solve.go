/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   solve.go                                           :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: mchi <mchi@student.42.fr>                  +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2019/04/20 02:48:49 by mchi              #+#    #+#             */
/*   Updated: 2019/04/20 02:48:49 by mchi             ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

//TriBool : enum for three results
type TriBool int

//types of result
const (
	False   TriBool = 0
	True    TriBool = 1
	Unknown TriBool = 2
)

/*
	when statement is solved as true,
	apply the conclusion as new rules.
	we discard the whole statement after that.
*/

//SolveLogics : solves and returns fact table.
func SolveLogics(input Input) map[Symbol]TriBool {
	factTable := make(map[Symbol]TriBool)
	//	for b := range input.refSheet {
	//		factTable[b] = Unknown
	//	}
	//	checkQueue := make([]*LogicNode, 0, len(input.facts))
	//	for _, sym := range input.facts {
	//		if factTable[sym] == Unknown {
	//			factTable[sym] = True
	//			checkQueue = append(checkQueue, input.refSheet[sym.name]...)
	//			for _, ref := range input.refSheet[sym.name] {
	//				ref.result = True
	//			}
	//		}
	//	}
	return factTable
}
