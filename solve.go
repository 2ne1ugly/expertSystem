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

//SolveLogics : solves and returns truth table.
func SolveLogics(input Input) map[byte]TriBool {
	return make(map[byte]TriBool)
}
