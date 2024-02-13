package models

import (
	"fmt"
	"regexp"
)

type persona struct {
	Id         int64
	rut        string
	nombre     string
	haberes    []elemento
	descuentos []elemento
}

type elemento struct {
	Id      int
	nombre  string
	descr   string
	tipo    string  //H haber D descuento
	valor   float64 // valor numerico
	moneda  string  //CLP UF %
	formula string  // formula calculada
}

func ValorSistema() {

	descuento1 := elemento{
		Id:      1,
		nombre:  "Fonasa",
		descr:   "",
		tipo:    "D",
		valor:   0,
		moneda:  "",
		formula: "ELEH(SuelBas) * 7% + ELEH(Hextr)",
	}

	descuento2 := elemento{
		Id:      2,
		nombre:  "AFP",
		tipo:    "D",
		valor:   13,
		moneda:  "%",
		formula: "",
	}
	descuento3 := elemento{
		Id:      3,
		nombre:  "AFC",
		tipo:    "D",
		valor:   0.013,
		moneda:  "%",
		formula: "",
	}

	haberes1 := elemento{
		Id:      1,
		nombre:  "SuelBas",
		tipo:    "H",
		valor:   800000,
		moneda:  "CLP",
		formula: "",
	}

	persona1 := persona{
		Id:         1,
		rut:        "12312312-3",
		nombre:     "Juan Pedro Gonzalez Tapia",
		haberes:    []elemento{haberes1},
		descuentos: []elemento{descuento1, descuento2, descuento3},
	}
	personal := []persona{persona1}
	//mapPersonas := make(map[string]Persona)
	LiqSueldo(personal)
}

func TotalHaberes()    {}
func TotalDescuentos() {}
func Fonasa()          {}

func LiqSueldo(p []persona) {

	for _, per := range p {
		fmt.Print("-----------------------------------------------------------------------------\n")
		fmt.Printf("ID : %d\n", per.Id)
		fmt.Printf("RUT : %s\n", per.rut)
		fmt.Printf("Nombre : %s\n", per.nombre)
		fmt.Print("-----------------------------------------------------------------------------\n")
		fmt.Printf("%-20s%-10s%-20s%-10s\n", "Haber", "Monto", "Descuento", "Monto")
		for i := 0; i < len(per.haberes) || i < len(per.descuentos); i++ {
			// Obtener el haber actual y el descuento actual
			var haberActual elemento
			var descuentoActual elemento

			if i < len(per.haberes) {
				haberActual = per.haberes[i]
				if len(haberActual.formula) > 1 {
					fmt.Printf("Fórmula a procesar %s\n", haberActual.formula)
					reemplazarELEH(haberActual.formula)
				}
			}

			if i < len(per.descuentos) {
				descuentoActual = per.descuentos[i]
				if len(descuentoActual.formula) > 1 {
					fmt.Printf("Fórmula a procesar %s\n", descuentoActual.formula)
					reemplazarELEH(descuentoActual.formula)
				}
			}

			// Mostrar los resultados en una sola fila

			fmt.Printf("%-20s%-10.2f%-20s%-10.2f\n",
				haberActual.nombre, haberActual.valor,
				descuentoActual.nombre, descuentoActual.valor)
		}

	}

}

//func Formular(e []elemento, forml string) {
//	forml.ree

//}

func reemplazarELEH(cadena string) (string, error) {
	// Definir la expresión regular para encontrar HELE() y su contenido
	expresionRegular := regexp.MustCompile(`ELEH\(([^)]+)\)`)

	// Encontrar todas las coincidencias en la cadena
	coincidencias := expresionRegular.FindAllStringSubmatch(cadena, -1)

	fmt.Printf("coincidencias -->%s\n", coincidencias)

	// Iterar sobre las coincidencias y reemplazar HELE() con el resultado evaluado

	for _, match := range coincidencias {
		contenidoELEH := match[1]
		fmt.Printf("contenidoELEH-->%s\n", contenidoELEH)
		/*
			resultado, err := evaluarHELE(contenidoELEH)
			if err != nil {
				return "", err
			}
		*/
		//cadena = strings.ReplaceAll(cadena, "HELE("+contenidoELEH+")", resultado)
	}

	//return cadena, nil

	return "xx", nil
}
