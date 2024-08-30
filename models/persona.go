package models

import (
	"fmt"
	"log"
	"math"
	"time"

	"strconv"
)

type Persona struct {
	IdPersona  uint //`gorm:"primaryKey"`
	Rut        string
	Nombre     string
	Afp        Afp          //`gorm:"foreignKey:IdAfp"`
	Salud      Salud        //`gorm:"foreignKey:IdSalud"`
	SaludComp  []saludCompl //`gorm:"foreignKey:IdSalComp"`
	Apv        Apv          //`gorm:"foreignKey:IdApv"`
	Haberes    []Elemento   //`gorm:"foreignKey:IdElemento"`
	Descuentos []Elemento   //`gorm:"foreignKey:IdElemento"`
	Acumulador []acumulador //`gorm:"foreignKey:IdAcum"`

}

type Afp struct {
	IdAfp  uint
	Codigo string
	Nombre string
	Valor  float64
	Moneda string
}

type Salud struct {
	IdSalud uint //`gorm:"primaryKey"`
	Codigo  string
	Nombre  string
	Valor   float64
	Moneda  string
}

type saludCompl struct {
	IdSalComp uint
	Codigo    string
	Nombre    string
	Valor     float64
	Moneda    string
}

type Apv struct {
	IdApv  uint
	Codigo string
	Nombre string
	Valor  float64
}

type Elemento struct {
	IdElemento uint
	Nombre     string
	Descr      string
	Tipo       string //H haber D descuento
	Calidad    string //IT Haber imponible Tributale NINT No imponile ni tributable NI No imponible NT No tributable
	//Descuento obligatorio OB Descuento Voluntario VO
	Cantidad  float64 //cantidad de unidades trabajas
	Valor     float64 // valor numerico
	Valmax    float64 //valor minimo
	Valmin    float64 //valor maximo tope
	Moneda    string  //CLP UF %
	Formula   string  // formula calculada
	Resultado float64
}
type operador struct {
	IdOper     uint
	IdAcum     uint
	IdElemento uint
	Element    Elemento
	Operacion  string
}

type acumulador struct {
	IdAcum    uint
	Rut       string
	IdPersona uint
	NomAcum   string
	Elementos []operador
	Resultado float64
}

type Moneda struct {
	TMoneda  string
	Fecha    time.Time
	Valorclp float32
}

/*
	type imprime struct {
		nomHaber      string
		cantidadHaber float64
		valHaber      float64
		nomDescu      string
		cantidadDescu float64
		valDescu      float64
	}
*/
func ValorSistema(totalpersonal []Persona, monedas []Moneda) {

	//CONFIGURACION AFP
	//Código Glosa
	//00 no está en AFP
	//03 Cuprum
	//05 Habitat
	//08 Provida
	//29 PlanVital
	//33 Capital
	//34 Modelo
	//35 Uno
	/*
		Afp1 := Afp{
			id:     1,
			codigo: "00",
			nombre: "no está en AFP",
			valor:  0,
			moneda: "",
		}
	*/
	Afp2 := Afp{
		IdAfp:  2,
		Codigo: "03",
		Nombre: "Cuprum",
		Valor:  11.44,
		Moneda: "%",
	}
	Afp3 := Afp{
		IdAfp:  3,
		Codigo: "05",
		Nombre: "Habitat",
		Valor:  11.27,
		Moneda: "%",
	}
	Afp4 := Afp{
		IdAfp:  4,
		Codigo: "08",
		Nombre: "Provida",
		Valor:  11.45,
		Moneda: "%",
	}
	Afp5 := Afp{
		IdAfp:  5,
		Codigo: "29",
		Nombre: "PlanVital",
		Valor:  11.16,
		Moneda: "%",
	}
	/*
		Afp6 := Afp{
			id:     6,
			codigo: "33",
			nombre: "Capital",
			valor:  11.44,
			moneda: "%",
		}
		Afp7 := Afp{
			id:     7,
			codigo: "34",
			nombre: "Modelo",
			valor:  10, 58,
			moneda: "%",
		}
		Afp8 := Afp{
			id:     8,
			codigo: "35",
			nombre: "Uno",
			valor:  10.49,
			moneda: "%",
		}
	*/
	descuento1 := Elemento{
		IdElemento: 1,
		Nombre:     "Fonasa",
		Calidad:    "OB",
		Descr:      "",
		Tipo:       "D",
		Valor:      0,
		Moneda:     "",
		Formula:    "ELEH(SuelBas) * 7% + ELEH(Hextr)",
	}

	descuento2 := Elemento{
		IdElemento: 2,
		Nombre:     "AFP",
		Calidad:    "OB",
		Tipo:       "D",
		Valor:      13,
		Moneda:     "%",
		Formula:    "#ACUMULADOR(TotalHaberesImponible) * 10%",
	}
	descuento3 := Elemento{
		IdElemento: 3,
		Nombre:     "AFC",
		Calidad:    "OB",
		Tipo:       "D",
		Valor:      0.013,
		Moneda:     "%",
		Formula:    "",
	}

	haberes1 := Elemento{
		IdElemento: 1,
		Nombre:     "SuelBas",
		Tipo:       "H",
		Calidad:    "IT",
		Valor:      800000,
		Moneda:     "CLP",
		Formula:    "",
	}

	haberes2 := Elemento{
		IdElemento: 2,
		Nombre:     "Hextr",
		Tipo:       "H",
		Calidad:    "IT",
		Formula:    "(ELEH(SuelBas) /30*28/180*1.5)*#CANTIDAD",
	}
	haberes3 := Elemento{
		IdElemento: 3,
		Nombre:     "HBColacion",
		Tipo:       "H",
		Calidad:    "NI",
		Valor:      80000,
		Moneda:     "CLP",
	}
	haberes4 := Elemento{
		IdElemento: 4,
		Nombre:     "HBMoviliza",
		Tipo:       "H",
		Calidad:    "NI",
		Valor:      60000,
		Moneda:     "CLP",
	}
	//asignando horas extras
	haberes2.Cantidad = 4
	persona1 := Persona{
		IdPersona:  1,
		Rut:        "12312312-3",
		Nombre:     "Juan Pedro Gonzalez Tapia",
		Afp:        Afp2,
		Haberes:    []Elemento{haberes1, haberes2, haberes3, haberes4},
		Descuentos: []Elemento{descuento1, descuento2, descuento3},
	}

	haberes1.Valor = 1600000
	haberes2.Cantidad = 0
	haberes2.Valmax = 0
	persona2 := Persona{
		IdPersona:  2,
		Rut:        "16888333-4",
		Nombre:     "Pedro Tapia Risopatron",
		Afp:        Afp3,
		Haberes:    []Elemento{haberes1, haberes2, haberes3, haberes4},
		Descuentos: []Elemento{descuento1, descuento2, descuento3},
	}
	haberes1.Valor = 850000
	haberes2.Cantidad = 10
	persona3 := Persona{
		IdPersona:  3,
		Rut:        "11222333-4",
		Nombre:     "Jose Leon MArtinez Benitez",
		Afp:        Afp4,
		Haberes:    []Elemento{haberes1, haberes2, haberes3, haberes4},
		Descuentos: []Elemento{descuento1, descuento2, descuento3},
	}
	haberes1.Valor = 800000
	haberes2.Cantidad = 7
	persona4 := Persona{
		IdPersona:  4,
		Rut:        "25454545-2",
		Nombre:     "Juan Pedro Gonzalez Tapia",
		Afp:        Afp5,
		Haberes:    []Elemento{haberes1, haberes2, haberes3, haberes4},
		Descuentos: []Elemento{descuento1, descuento2, descuento3},
	}

	personal := []Persona{persona1, persona2, persona3, persona4}
	//mapPersonas := make(map[string]Persona)
	// LiqSueldo(personal)
	if len(personal) > 1 {
		fmt.Printf("\n Personal calcula \n")
	}
	//fmt.Printf("personal -> %v", personal)
	//Acumuladores
	//db.ListarPersonas()
	//go CalculaNom(personal)
	//go
	go CalculaNom(totalpersonal, monedas)

}

func TotalizaAcumulador(acum []operador) float64 {
	valor := 0.0
	for _, elem := range acum {
		opera := elem.Operacion
		//fmt.Printf("%s----> %v \n", elem.element.nombre, elem.element.resultado)
		switch opera {
		case "Suma":
			valor += elem.Element.Resultado
		case "Resta":
			valor -= elem.Element.Resultado
		default:
			log.Fatal("Error Totaliza acumulador")
		}

	}
	return valor
}

func Fonasa() {}

func LiqSueldo(p []Persona, m []Moneda) {

	valHab := 0.0
	valDesc := 0.0

	for _, per := range p {
		fmt.Print("-----------------------------------------------------------------------------\n")
		fmt.Printf("ID : %d\n", per.IdPersona)
		fmt.Printf("RUT : %s\n", per.Rut)
		fmt.Printf("Nombre : %s\n", per.Nombre)
		fmt.Print("-----------------------------------------------------------------------------\n")
		fmt.Printf("%-20s%-10s%-20s%-10s\n", "Haber", "Monto", "Descuento", "Monto")
		for i := 0; i < len(per.Haberes) || i < len(per.Descuentos); i++ {
			// Obtener el haber actual y el descuento actual
			var haberActual Elemento
			var descuentoActual Elemento

			if i < len(per.Haberes) {
				haberActual = per.Haberes[i]
				if len(haberActual.Formula) > 1 {
					//fmt.Printf("Fórmula a procesar %s\n", haberActual.formula)
					formu, err := reemplazarELEM(haberActual.Formula, per, per.Haberes[i], m)
					if err != nil {
						log.Fatalf("haberes con problemas reemplazarELEM %s", err)
					}
					//fmt.Printf("formula a a calcular Habe %s \n\n", formu)

					valHab, err = ResolverOperacion(formu, haberActual.Valmin, haberActual.Valmax, haberActual.Cantidad)
					if err != nil {
						log.Fatalf("haberes con problemas en ResolverOperacion %s", err)

					}
					haberActual.Resultado = valHab
					//fmt.Printf("valor formula calculada %s=  %f \n", formu, val)

				} else {
					haberActual.Resultado = haberActual.Valor
				}
			}

			if i < len(per.Descuentos) {
				descuentoActual = per.Descuentos[i]
				if len(descuentoActual.Formula) > 1 {
					//fmt.Printf("Fórmula a procesar %s\n", descuentoActual.formula)
					formu, err := reemplazarELEM(descuentoActual.Formula, per, per.Descuentos[i], m)
					if err != nil {
						log.Fatalf("descuentos reemplazarELE %s", err)
					}
					//fmt.Printf("formula a a calcular Desc %s \n\n", formu)
					valDesc, err = ResolverOperacion(formu, descuentoActual.Valmin, descuentoActual.Valmax, descuentoActual.Cantidad)
					if err != nil {
						log.Fatalf("descuentos ResolverOperacion %s", err)

					}
					descuentoActual.Resultado = valDesc
				} else {
					descuentoActual.Resultado = descuentoActual.Valor
				}
			}

			// Mostrar los resultados en una sola fila

			fmt.Printf("%-20s%-10.2f%-20s%-10.2f\n",
				haberActual.Nombre, math.Round(haberActual.Resultado),
				descuentoActual.Nombre, math.Round(descuentoActual.Resultado))
			valHab = 0.0
			valDesc = 0.0
		}

	}
}

func CalculaNom(p []Persona, moneda []Moneda) {

	valHab := 0.0
	valDesc := 0.0
	var PerTemp Persona
	//fmt.Printf("error paso 1 %v \n", p)
	for _, per := range p {

		if PerTemp.IdPersona != per.IdPersona {
			PerTemp = Persona{}
		}
		//fmt.Printf("\n personaT-->%v \n", p)
		//fmt.Printf("persona-->%s \n", per.Nombre)
		//fmt.Printf("persona Haberes-->%v \n", per.Haberes)

		for i := 0; i < len(per.Haberes); i++ {
			// Obtener el haber actual y el descuento actual
			//fmt.Printf("valores haberes--> %v  \n", per.Haberes[i])
			var haberActual Elemento

			if i < len(per.Haberes) {
				haberActual = per.Haberes[i]
				if len(haberActual.Formula) > 1 {
					//fmt.Printf("Fórmula a procesar %s\n", haberActual.formula)
					//fmt.Printf("error paso 2.1 %v \n", per.Haberes[i])
					formu, err := reemplazarELEM(haberActual.Formula, per, per.Haberes[i], moneda)
					if err != nil {

						log.Fatalf("haberes reemplazarELEM %s \n", err)
					}

					//Revisar las asignaciones minimos y maximos
					if len(haberActual.Moneda) > 0 {
						for _, v := range moneda {
							if haberActual.Moneda == v.TMoneda {
								if (haberActual.Valmin) > 0 {
									haberActual.Valmin = haberActual.Valmin * float64(v.Valorclp)
								}

							}
							if (haberActual.Valmax) > 0 {
								haberActual.Valmax = haberActual.Valmax * float64(v.Valorclp)
							}
						}

					}
					//fmt.Printf("formula a a calcular haberes 2 %s \n\n", formu)
					valHab, err = ResolverOperacion(formu, haberActual.Valmin, haberActual.Valmax, haberActual.Cantidad)
					if err != nil {
						log.Fatalf("haberes ResolverOperacion %s", err)

					}
					haberActual.Resultado = valHab
					per.Haberes[i].Resultado = haberActual.Resultado
					//fmt.Printf("valor formula calculada %s=  %f \n", formu, val)

				} else {
					haberActual.Resultado = haberActual.Valor
					per.Haberes[i].Resultado = haberActual.Resultado
				}
			}

		}
		//fmt.Printf("i valor --> %v \n", per)
		//per.nombre = per.nombre
		//PerTemp.Id = per.Id
		//PerTemp.rut = per.rut
		//PerTemp.haberes = per.haberes

		// Total_imponible total_Imponible Total_tributable Total_no_imponible*/
		//fmt.Println("error paso 3")
		totalHaberesReturn, totalHaberesImpReturn, totalHaberesTribReturn, totalHaberesNoTribReturn := TotalHaberes(per)

		//fmt.Printf("total haberes Imp %v ", totalHaberesImpReturn)

		per.Acumulador = append(per.Acumulador, totalHaberesReturn)
		per.Acumulador = append(per.Acumulador, totalHaberesImpReturn)
		per.Acumulador = append(per.Acumulador, totalHaberesTribReturn)
		per.Acumulador = append(per.Acumulador, totalHaberesNoTribReturn)

		for i := 0; i < len(per.Descuentos); i++ {
			// Obtener el haber actual y el descuento actual
			var descuentoActual Elemento

			if i < len(per.Descuentos) {
				descuentoActual = per.Descuentos[i]
				//fmt.Printf("descuentos en form --> %v \n", per.Descuentos[i])
				if len(descuentoActual.Formula) > 1 {
					//fmt.Printf("Fórmula a procesar %s\n", descuentoActual.Formula)
					formu, err := reemplazarELEM(descuentoActual.Formula, per, per.Descuentos[i], moneda)
					//fmt.Printf("Fórmula a procesar calculo %s\n", formu)
					if err != nil {
						log.Fatalf("descuento reemplazarELEM %s", err)
					}
					//fmt.Printf("formula a a calcular descuento 2 %s \n\n", formu)
					valDesc, err = ResolverOperacion(formu, descuentoActual.Valmin, descuentoActual.Valmax, descuentoActual.Cantidad)
					if err != nil {
						log.Fatalf("descuento ResolverOperacion %s", err)

					}
					descuentoActual.Resultado = valDesc
					per.Descuentos[i].Resultado = descuentoActual.Resultado

				} else {
					//agsinar valores si no hay formulas
					//descuentoActual.Moneda
					descuentoActual.Resultado = descuentoActual.Valor
					per.Descuentos[i].Resultado = descuentoActual.Resultado
				}
			}

		}
		//PerTemp.descuentos = per.descuentos

		// Mostrar los resultados en una sola fila
		totalDescuentosReturn := TotalDescuentos(per)
		per.Acumulador = append(per.Acumulador, totalDescuentosReturn)
		valHab = 0.0
		valDesc = 0.0

		//fmt.Printf("Resultado-->%v\n,", PerTemp)

		for _, el := range per.Acumulador {
			fmt.Printf("elemetos en acumulador-->%v resultado %v rut %s\n", el.NomAcum, strconv.FormatFloat(el.Resultado, 'f', 0, 64), el.Rut)
		}

	}
	//fmt.Printf("Resultado-->%v", PerTemp)
}

// /LISTAR PERSONAS DESDE BASE DE DATOS
type Personas []Persona

/*
haberes2.Cantidad = 4
	persona1 := persona{
		IdPersona:  1,
		Rut:        "12312312-3",
		Nombre:     "Juan Pedro Gonzalez Tapia",
		Afp:        Afp2,
		Haberes:    []elemento{haberes1, haberes2, haberes3, haberes4},
		Descuentos: []elemento{descuento1, descuento2, descuento3},
	}

*/
