/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: mchi <mchi@student.42.fr>                  +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2019/04/19 15:21:31 by mchi              #+#    #+#             */
/*   Updated: 2019/04/19 15:21:31 by mchi             ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("need one file.\n")
	}
	fmt.Printf("statements:\n")
	input := ParseFile(os.Args[1])
	factMap := SolveLogics(input)
	fmt.Printf("\nqueries:\n")
	for _, sym := range input.query {
		fmt.Printf("%c is ", sym)
		switch factMap[sym] {
		case True:
			fmt.Printf("True\n")
		case False:
			fmt.Printf("False\n")
		case Unknown:
			fmt.Printf("Unknown\n")
		}
	}
}
