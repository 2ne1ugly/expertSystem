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
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("need one file.\n")
	}
	input := ParseFile(os.Args[1])
	truthTable := SolveLogics(input)
	print(truthTable)
}
