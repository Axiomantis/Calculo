package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

// ResolverOperacion acepta una cadena de operación y devuelve el resultado como float64
func ResolverOperacion(operacion string, min float64, max float64, cantidad float64) (float64, error) {
	// Eliminar espacios en blanco
	operacion = strings.ReplaceAll(operacion, " ", "")
	//println("operacion --> ", operacion)
	// Crear un evaluador de expresiones
	//fmt.Printf("Operacion en calculo--> %s \n", operacion)
	if len(operacion) == 0 {
		return 0, errors.New("error en operacion: la operación no puede estar vacía")
	}

	expr, err := govaluate.NewEvaluableExpression(operacion)
	if err != nil {
		return 0, err
	}

	// Evaluar la expresión
	resultado, err := expr.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	// Convertir el resultado a float64
	switch v := resultado.(type) {
	case float64:
		/*
			if cantidad > 1 {
				v = v * cantidad
			}
		*/
		if v < min {
			v = min
		}
		if v > max && max > 0 {
			v = max
		}
		return v, nil
	default:
		return 0, fmt.Errorf("resultado no es un número")
	}
}

/*
Comentado ResolverOperacion
func ResolverOperacion(operacion string) (float64, error) {
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
	//decimal := false
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

			if caracter != "." {
				numero, err := strconv.ParseFloat(caracter, 64)
				if err != nil {
					return 0, err
				}
				numeros = append(numeros, numero)
			} else {
				numeros = numeros

			}

		}
	}

	// Procesar cualquier operador restante
	for len(operadores) > 0 {
		evaluarOperador()
	}

	// El resultado debe estar en la cima de la pila de números
	if len(numeros) != 1 {
		fmt.Printf("numero --> %f \n", numeros)
		return 0, fmt.Errorf("formato de operación incorrecto")
	}

	return numeros[0], nil
}
Fin Comentado Resolver Operacion
*/
/* comentado PrioridadesOperador
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
fin comentado PrioridadesOperador
*/
