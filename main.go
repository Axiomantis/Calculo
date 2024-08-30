package main

import (
	"Calculo/config"
	"Calculo/db"
	"Calculo/handlers"
	"Calculo/models"
	"bufio"
	"fmt"
	"os"

	"github.com/gorilla/mux"
)

func creaListadoPErsonal(personasdb1 db.Personasdb, afpsdb1 db.Afpsdb, elemEmpldbs1 db.ElemEmpldbs, elementDbs1 db.ElementDbs) []models.Persona {

	var personas []models.Persona

	for _, per := range personasdb1 {
		persona := models.Persona{}
		persona.IdPersona = per.IdPersona
		persona.Rut = per.Rut
		persona.Nombre = per.Nombre
		if per.Idafp.Valid {
			persona.Afp = devuelveAfp(uint(per.Idafp.Int32), afpsdb1)
		}
		//fmt.Printf("creaListadoPersonal--> %v \n", elemEmpldbs1)

		persona.Haberes = devuelveHaberes(per.IdPersona, elemEmpldbs1, elementDbs1)

		persona.Descuentos = devuelveDebes(per.IdPersona, elemEmpldbs1, elementDbs1)
		//monedasAct
		//config.Monedas
		//

		personas = append(personas, persona)
	}

	return personas
}
func devuelveHaberes(idEmpl uint, elemEmpldbs db.ElemEmpldbs, elementDbs []db.ElementDb) []models.Elemento {
	var habElement []models.Elemento
	for _, elemEmp := range elemEmpldbs {
		hab := models.Elemento{}
		if elemEmp.Idpersona == idEmpl {
			//fmt.Printf("valida elemento devuelveHaberes-->%v ,%v \n", elemEmp.Idelemento, elemEmp)
			hab = devuelveElement(elemEmp.Idelemento, elementDbs, "H")
			//fmt.Printf("valida elemento IdElemento-->%v\n", hab.IdElemento)
			if hab.Descr == " " {
				continue
			}
			/*validar si el elemeto de la incidencia posee valor antes de setear*/
			//fmt.Printf("revision devuelveHAberes--> %v \n", hab)
			if elemEmp.Cantidad.Valid {
				hab.Cantidad = elemEmp.Cantidad.Float64
			}

			if elemEmp.Formula.Valid {
				hab.Formula = elemEmp.Formula.String
			}
			hab.IdElemento = elemEmp.Idelemento
			//fmt.Printf("parte 1 %v resultado 1 %v resultado 2 %v \n", hab.Nombre, hab.Resultado, elemEmp.Resultado.Float64)
			hab.Resultado = elemEmp.Resultado.Float64
			if hab.Resultado > 0 {
				hab.Valor = hab.Resultado
			}
			//fmt.Printf("parte s %v resultado 1 %v resultado 2 %v \n", hab.Nombre, hab.Resultado, elemEmp.Resultado.Float64)
			if elemEmp.Moneda.Valid {
				hab.Moneda = elemEmp.Moneda.String
			}

		}
		//fmt.Printf("Rev -->elemento %v %v \n", hab.Nombre, hab.Formula)
		//fmt.Printf("Rev -->elem for %v %v \n", hab.Nombre, hab.Resultado)
		if hab.IdElemento != 0 && len(hab.Descr) > 1 {
			habElement = append(habElement, hab)
		}
	}
	/*
		elemEmpl :=
			elemEmpldbs
		elementDbs
	*/
	return habElement
}
func devuelveDebes(idEmpl uint, elemEmpldbs db.ElemEmpldbs, elementDbs []db.ElementDb) []models.Elemento {
	var debElement []models.Elemento
	for _, elemEmp := range elemEmpldbs {
		deb := models.Elemento{}
		if elemEmp.Idpersona == idEmpl {
			deb = devuelveElement(elemEmp.Idelemento, elementDbs, "D")
			if deb.Descr == "" {
				continue
			}
			deb.Cantidad = elemEmp.Cantidad.Float64
			deb.Formula = elemEmp.Formula.String
			deb.IdElemento = elemEmp.Idelemento
			deb.Resultado = elemEmp.Resultado.Float64
			if elemEmp.Moneda.Valid {
				deb.Moneda = elemEmp.Moneda.String
			}

		}
		if deb.IdElemento != 0 && len(deb.Descr) > 1 {
			debElement = append(debElement, deb)
		}
	}
	/*
		elemEmpl :=
			elemEmpldbs
		elementDbs
	*/
	return debElement
}
func devuelveElement(idelemento uint, elemedb []db.ElementDb, tipo string) models.Elemento {
	elemento := models.Elemento{}
	for _, elem := range elemedb {
		//fmt.Printf(" formula Valid range--> %v %v \n", elem.Descr, elem)
		if elem.Idelemento == idelemento && elem.Tipo == tipo {
			//fmt.Printf(" devuelveElement --> %v %v tipo :%v %v\n", elem.Idelemento, idelemento, elem.Tipo, tipo)
			elemento.Nombre = elem.Nombre
			elemento.Descr = elem.Descr
			elemento.Tipo = elem.Tipo
			elemento.Calidad = elem.Calidad
			if elem.Cantidad.Valid {
				elemento.Cantidad = elem.Cantidad.Float64
			}
			//fmt.Printf(" formula --> %v %v \n", elem.Formula.Valid, elem.Formula.String)
			if elem.Formula.Valid {
				elemento.Formula = elem.Formula.String
			}
			if elem.Moneda.Valid {
				elemento.Moneda = elem.Moneda.String
			}
			if elem.Valmax.Valid {
				elemento.Valmax = elem.Valmax.Float64
			}
			if elem.Valmin.Valid {
				elemento.Valmin = elem.Valmin.Float64
			}
			//fmt.Printf(" valida return --> %v %v %v \n", elem.Valor.Valid, elem.Resultado.Valid, elem.Resultado.Float64)
			if elem.Valor.Valid {
				if elem.Resultado.Valid {
					if elem.Resultado.Float64 > 0 {
						elemento.Valor = elem.Resultado.Float64
					} else {
						elemento.Valor = elem.Valor.Float64
					}
				} else {
					elemento.Valor = elem.Valor.Float64
				}
			}
		} else {
			continue
		}
	}
	//fmt.Printf(" formula return --> %v %v \n", elemento, elemento.Valor)
	return elemento
}
func devuelveAfp(idafp uint, listaAfp db.Afpsdb) models.Afp {
	afp1 := models.Afp{}
	for _, afp := range listaAfp {

		if afp.Idafp == idafp {
			afp1.IdAfp = afp.Idafp
			afp1.Codigo = afp.Codigo
			afp1.Nombre = afp.Nombre
			//afp1.Valor = afp.Valor
			if afp.Valor.Valid {
				afp1.Valor = afp.Valor.Float64
				//} else {
				//afp1.Valor = 0.0 // Por ejemplo, asignar 0.0 como valor predeterminado
			}
			if afp.Moneda.Valid {
				afp1.Moneda = afp.Moneda.String
			}

			return afp1
		}
	}
	return afp1
}

func main() {
	config.Init()
	fmt.Print("go calcula>")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		operacion := scanner.Text()
		//"2+(3*4+5)" 19
		resultado, err := models.ResolverOperacion(operacion, 0, 0, 0)
		if err != nil {
			fmt.Println("Error al resolver la operación:", err, operacion)
			return
		}

		fmt.Println("Resultado de la operación:", resultado)

		var monedasAct []models.Moneda
		for _, mondAct := range config.Monedas {
			monedaAct := models.Moneda{}
			monedaAct.TMoneda = mondAct.TMoneda
			monedaAct.Valorclp = float32(mondAct.ValorCLP)
			monedasAct = append(monedasAct, monedaAct)
		}

		personasdb, afpsdb, elemEmpldbs, elementDbs := db.ListarPersonas()
		personal := creaListadoPErsonal(personasdb, afpsdb, elemEmpldbs, elementDbs)
		//fmt.Printf("personal main %v \n", personal)
		models.ValorSistema(personal, monedasAct)

		mux := mux.NewRouter()
		mux.HandleFunc("/app/user/", handlers.GetEmpleado).Methods("GET")

	}
}
