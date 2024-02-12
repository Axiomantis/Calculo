package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Print("go calcula>")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		operacion := scanner.Text()
		//"2+(3*4+5)" 19

		resultado, err := resolverOperacion(operacion)
		if err != nil {
			fmt.Println("Error al resolver la operación:", err)
			return
		}

		fmt.Println("Resultado de la operación:", resultado)
	}
}
func resolverOperacion(operacion string) (float64, error) {
	// Eliminar espacios en blanco
	operacion = strings.ReplaceAll(operacion, " ", "")

	// Convertir la cadena de operación a una lista de caracteres
	caracteres := strings.Split(operacion, "")

	// Utilizar dos pilas: una para números y otra para operadores
	numeros := []float64{}
	operadores := []string{}

	// Función auxiliar para evaluar operadores
	evaluarOperador := func() {
		if len(operadores) > 0 && operadores[len(operadores)-1] == "(" {
			return
		}

		if len(numeros) >= 2 && len(operadores) >= 1 {
			num2 := numeros[len(numeros)-1]
			num1 := numeros[len(numeros)-2]
			operador := operadores[len(operadores)-1]

			numeros = numeros[:len(numeros)-2]
			operadores = operadores[:len(operadores)-1]

			var resultado float64

			switch operador {
			case "+":
				resultado = num1 + num2
			case "-":
				resultado = num1 - num2
			case "*":
				resultado = num1 * num2
			case "/":
				resultado = num1 / num2
			}

			numeros = append(numeros, resultado)
		}
	}

	// Iterar sobre los caracteres de la operación
	for _, caracter := range caracteres {
		switch caracter {
		case "(":
			operadores = append(operadores, caracter)
		case ")":
			for len(operadores) > 0 && operadores[len(operadores)-1] != "(" {
				evaluarOperador()
			}
			operadores = operadores[:len(operadores)-1] // Pop "("
		case "+", "-", "*", "/":
			for len(operadores) > 0 && prioridadOperador(operadores[len(operadores)-1]) >= prioridadOperador(caracter) {
				evaluarOperador()
			}
			operadores = append(operadores, caracter)
		default:
			// Es un número, convertir a float64 y apilar en la pila de números
			numero, err := strconv.ParseFloat(caracter, 64)
			if err != nil {
				return 0, err
			}
			numeros = append(numeros, numero)
		}
	}

	// Procesar cualquier operador restante
	for len(operadores) > 0 {
		evaluarOperador()
	}

	// El resultado debe estar en la cima de la pila de números
	if len(numeros) != 1 {
		return 0, fmt.Errorf("formato de operación incorrecto")
	}

	return numeros[0], nil
}

// Función para asignar prioridades a los operadores
func prioridadOperador(operador string) int {
	switch operador {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}
