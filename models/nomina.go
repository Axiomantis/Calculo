package models

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var tipoTabla = map[string]reflect.Type{
	"Afp": reflect.TypeOf(Afp{}),
}

// BuscarMoneda utiliza reflexión para buscar una moneda en la tabla especificada
func BuscarValores(tabla interface{}, moneda string) (interface{}, error) {
	v := reflect.ValueOf(tabla)

	// Verifica que la tabla sea un slice
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("tabla no es un slice")
	}

	// Itera sobre los elementos del slice
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()

		// Verifica si el elemento tiene los campos Moneda y Valor
		if monedaValue, ok := getFieldValue(item, "Moneda"); ok {
			if monedaValue == moneda {
				return item, nil
			}
		}
	}

	return nil, fmt.Errorf("moneda no encontrada")
}

// getFieldValue usa reflexión para obtener el valor de un campo de una estructura
func getFieldValue(obj interface{}, fieldName string) (interface{}, bool) {
	v := reflect.ValueOf(obj)
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, false
	}
	return field.Interface(), true
}

func getValueFromString(data interface{}, query string) (interface{}, error) {
	// Validar y parsear el string de consulta
	if !strings.HasPrefix(query, "#TABLA(") || !strings.HasSuffix(query, ")") {
		return nil, fmt.Errorf("formato de consulta no válido")
	}
	//fmt.Printf("query0-->%v \n", query)
	query = strings.TrimPrefix(query, "#TABLA(")
	//fmt.Printf("query1-->%v\n", query)
	query = strings.TrimSuffix(query, ")")
	//fmt.Printf("query2-->%v\n", query)
	parts := strings.Split(query, ".")

	//fmt.Printf("valores en getvaluefromString %s", parts)
	if len(parts) != 2 {
		return nil, fmt.Errorf("formato de consulta no válido")
	}

	tableName, fieldName := parts[0], parts[1]

	// Obtener el tipo de la tabla desde el mapa global
	tipo, ok := tipoTabla[tableName]
	if !ok {
		return nil, fmt.Errorf("tipo de tabla '%s' no encontrado", tableName)
	}

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("el tipo de datos no es un Struct ")
	}

	if v.Type() != tipo {
		return nil, fmt.Errorf("el tipo de datos no coincide con el tipo de la tabla '%s'", tableName)
	}
	// Acceder al campo especificado en el struct
	value, valid := getFieldValue(v.Interface(), fieldName)
	if valid {
		return value, nil
	}
	return nil, fmt.Errorf("campo '%s' no encontrado en la tabla '%s'", fieldName, tableName)

}

func reemplazarELEM(formu string, per Persona, elem Elemento, moneda []Moneda) (string, error) {

	// Definir la expresión regular para encontrar HELE() y su contenido
	err := ""
	flag := true
	//formu = strings.ReplaceAll(formu, "%", "/100")
	//fmt.Printf("valor formula ini %s \n", formu)
	for flag {
		formu = strings.ReplaceAll(formu, "%", "/100")
		formu = strings.ReplaceAll(formu, "#CANTIDAD", strconv.FormatFloat(elem.Cantidad, 'f', 6, 64))
		// ingresar la palabra tabla
		expresionRegularTabla := regexp.MustCompile(`#TABLA\(([^)]+)\)`)
		coincidenciasTabla := expresionRegularTabla.FindAllStringSubmatch(formu, -1)
		if len(coincidenciasTabla) > 0 {
			for _, match := range coincidenciasTabla {
				contenidoTabla := match[1]

				query := match[0]
				valor, err := getValueFromString(per.Afp, query)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					//fmt.Printf("Valor encontrado: %v  contenido t %v\n", valor, contenidoTabla)
					formu = strings.ReplaceAll(formu, match[0], strconv.FormatFloat(valor.(float64), 'f', 6, 64))
				}
				/**/
				for i, a := range per.Acumulador {
					//fmt.Printf("a.nomAcum -->%v\n", a.NomAcum)
					if a.NomAcum == contenidoTabla {
						//fmt.Printf("formula acumuladores--> %s valores %v\n", formu, contenidoTabla)
						formu = strings.ReplaceAll(formu, match[0], strconv.FormatFloat(per.Acumulador[i].Resultado, 'f', 6, 64))
					}
				}
				//fmt.Printf("formu--> %s", formu)
				/**/

			}
		}

		expresionRegularAcum := regexp.MustCompile(`#ACUMULADOR\(([^)]+)\)`)
		coincidenciasAcum := expresionRegularAcum.FindAllStringSubmatch(formu, -1)
		//fmt.Printf("coincidenciasAcum -->%v\n", coincidenciasAcum)
		//fmt.Printf("persona -->%v\n", per.acumulador)
		if len(coincidenciasAcum) > 0 {
			for _, match := range coincidenciasAcum {
				contenidoAcum := match[1]
				for i, a := range per.Acumulador {
					//fmt.Printf("a.nomAcum -->%v\n", a.nomAcum)
					if a.NomAcum == contenidoAcum {
						//fmt.Printf("formula acumuladores--> %s valores %v\n", formu, contenidoAcum)
						formu = strings.ReplaceAll(formu, match[0], strconv.FormatFloat(per.Acumulador[i].Resultado, 'f', 6, 64))
					}

				}
			}
		}
		expresionRegularH := regexp.MustCompile(`ELEH\(([^)]+)\)`)

		// Encontrar todas las coincidencias en la cadena
		coincidenciasH := expresionRegularH.FindAllStringSubmatch(formu, -1)

		//fmt.Printf("coincidencias 1 -->%s  ---> %d\n", coincidenciasH, len(coincidenciasH))

		// Iterar sobre las coincidencias y reemplazar HELE() con el resultado evaluado
		if len(coincidenciasH) > 0 {
			for _, match := range coincidenciasH {
				contenidoELEH := match[1]
				//fmt.Printf("contenidoELEH-->%s\n", contenidoELEH)
				valor := reemplazaValores(contenidoELEH, "H", per.Haberes, per.Descuentos)
				//fmt.Printf()
				if valor == "" {
					err = contenidoELEH
				}
				formu = strings.Replace(formu, match[0], valor, -1)

				//fmt.Printf("formu-->%v \n", formu)
			}
		}
		expresionRegularD := regexp.MustCompile(`ELED\(([^)]+)\)`)

		// Encontrar todas las coincidencias en la cadena
		coincidenciasD := expresionRegularD.FindAllStringSubmatch(formu, -1)

		//fmt.Printf("coincidencias -->%s\n", coincidenciasD)

		// Iterar sobre las coincidencias y reemplazar HELE() con el resultado evaluado
		if len(coincidenciasD) > 0 {
			for _, match := range coincidenciasD {
				contenidoELEH := match[1]
				//fmt.Printf("contenidoELED-->%s\n", contenidoELEH)
				valor := reemplazaValores(contenidoELEH, "D", per.Haberes, per.Descuentos)
				if valor == "" {
					err = contenidoELEH
				}
				formu = strings.Replace(formu, match[0], valor, -1)
				formu = strings.ReplaceAll(formu, "%", "/100")
				formu = strings.ReplaceAll(formu, "#CANTIDAD", strconv.FormatFloat(elem.Cantidad, 'f', 6, 64))

			}
		}
		//fmt.Printf("valor formula antes de evaluar %s\n", formu)
		coincidenciasH = expresionRegularH.FindAllStringSubmatch(formu, -1)
		coincidenciasD = expresionRegularD.FindAllStringSubmatch(formu, -1)
		//fmt.Printf("valores formula encontradasH %s, cantidad %d encontradasD %s flag %v \n", coincidenciasH, len(coincidenciasH), coincidenciasD, flag)

		if (len(coincidenciasH) == 0) && (len(coincidenciasD) == 0) && (len(coincidenciasAcum) == 0) {
			//fmt.Printf("si len(coincidenciasH) %d==0 AND  len(coincidenciasD) %d=0", len(coincidenciasH), len(coincidenciasD))
			flag = false
		}

	}
	//fmt.Printf("valor formula return %s\n", formu)
	if err == "" {
		return formu, nil
	} else {
		return formu, errors.New("el elemento " + fmt.Sprint(err) + " no ha sido creado")
	}
}

func reemplazaValores(elemento string, tipo string, haberes []Elemento, deducciones []Elemento) string {
	valor := ""
	if tipo == "H" {
		for _, match := range haberes {
			//fmt.Printf("%s == %s formula=%s valor=%v\n", match.nombre, elemento, match.formula, match.valor)
			if match.Nombre == elemento {
				if len(match.Formula) > 1 {
					valor = match.Formula
				} else {
					valor = strconv.FormatFloat(match.Valor, 'f', 2, 64)
				}
				break
			}
		}
	} else {
		for _, match := range deducciones {

			if match.Nombre == elemento {
				if len(match.Formula) > 1 {
					valor = match.Formula
				} else {
					valor = strconv.FormatFloat(match.Valor, 'f', 2, 64)
				}
				break
			}

		}
	}
	//fmt.Printf("Return Valor--> %v \n", valor)
	return valor
}

func TotalHaberes(per Persona) (acumulador /*Total imponible*/, acumulador /*total Imponible*/, acumulador /*Total tributable*/, acumulador /*Total no imponible*/) {
	TotHabImponibleRet := acumulador{}
	TotHabTributableRet := acumulador{}
	TotalHabNoImpoRet := acumulador{}

	totalHaberes := acumulador{
		//id int64,
		//rut string,
		IdAcum:  1,
		NomAcum: "TotalHaberes",
		//elementos    []operador,
		Resultado: 0.0,
	}
	totalHaberes.Rut = per.Rut
	totalHaberes.IdPersona = per.IdPersona

	//i:=1
	for _, elem := range per.Haberes {
		//fmt.Printf("Elemento---->%v \n", elem)
		totalHaberes.Elementos = append(totalHaberes.Elementos, operador{
			Element:    elem, // Puedes inicializar elemento según tus necesidades
			Operacion:  "Suma",
			IdElemento: elem.IdElemento,
			IdAcum:     totalHaberes.IdAcum,
		})

		acumcalidad := elem.Calidad
		// IT Haber imponible Tributale NIT No imponile ni tributable NI No imponible NT No tributable
		//fmt.Printf("elem-->,%v, tipo-->%s\n", elem, elem.calidad)
		switch acumcalidad {

		case "IT": // IT Haber imponible Tributale
			TotHabImponibleRet = TotalHaberesImponible(per.Rut, per.IdPersona, elem, TotHabImponibleRet)
			//fmt.Printf("en case TotHabImponibleRet--->%v \n", TotHabImponibleRet)
			TotHabTributableRet = TotalHaberesTributales(per.Rut, per.IdPersona, elem, TotHabTributableRet)
		//case "NIT": // NIT No imponile ni tributable
		//TotalHaberesNoImponible(per.rut, per.Id, elem)
		case "NI": //NI No imponible
			TotalHabNoImpoRet = TotalHaberesNoImponible(per.Rut, per.IdPersona, elem, TotalHabNoImpoRet)
		//case "NT": //NT No tributable
		//	valor -= elem.element.resultado
		default:
			log.Fatal("Error Total haberes %v \n", totalHaberes.Elementos)
		}
	}
	totalHaberes.Resultado = TotalizaAcumulador(totalHaberes.Elementos)
	//TotHabImponibleRet.resultado = TotalizaAcumulador(TotHabImponibleRet.elementos)

	//fmt.Printf("TotHabImponibleRet=== %v \n", TotHabImponibleRet)
	fmt.Printf("%s=== %v \n", totalHaberes.NomAcum, int(totalHaberes.Resultado))
	return totalHaberes, TotHabImponibleRet, TotHabTributableRet, TotalHabNoImpoRet
}

func TotalDescuentos(per Persona) acumulador {

	TotalDescuentos := acumulador{
		//id int64,
		//rut string,
		IdAcum:  1,
		NomAcum: "TotalDescuento",
		//elementos    []operador,
		Resultado: 0.0,
	}

	TotalDescuentos.Rut = per.Rut
	TotalDescuentos.IdPersona = per.IdPersona
	for _, elem := range per.Descuentos {
		//fmt.Printf("valores descuentos--> %v \n", elem)
		TotalDescuentos.Elementos = append(TotalDescuentos.Elementos, operador{
			Element:    elem, // Puedes inicializar elemento según tus necesidades
			Operacion:  "Suma",
			IdElemento: elem.IdElemento,
			IdAcum:     TotalDescuentos.IdAcum,
		})
		acumcalidad := elem.Calidad
		//Descuento obligatorio OB Descuento Voluntario VO
		switch acumcalidad {
		case "OB": // IT Haber imponible Tributale
			TotalDescuentosObligatorios(per.Rut, per.IdPersona, elem)

		//case "NIT": // NIT No imponile ni tributable
		//TotalHaberesNoImponible(per.rut, per.Id, elem)
		case "VO": //NI No imponible
			TotalDescuentosVoluntarios(per.Rut, per.IdPersona, elem)
		//case "NT": //NT No tributable
		//	valor -= elem.element.resultado
		default:
			log.Fatal("Error Total descuentos ")
		}

	}

	TotalDescuentos.Resultado = TotalizaAcumulador(TotalDescuentos.Elementos)

	fmt.Printf("%s=== %v \n", TotalDescuentos.NomAcum, int(TotalDescuentos.Resultado))
	return TotalDescuentos
}
func TotalHaberesImponible(rut string, Id uint, elem Elemento, acum acumulador) acumulador {

	if acum.NomAcum != "TotalHaberesImponible" {
		acum.IdPersona = Id
		acum.Rut = rut
		acum.IdAcum = 3
		acum.NomAcum = "TotalHaberesImponible"
		//elementos    []operador,
		acum.Resultado = 0.0
	}

	acum.Elementos = append(acum.Elementos, operador{
		Element:    elem, // Puedes inicializar elemento según tus necesidades
		Operacion:  "Suma",
		IdElemento: elem.IdElemento,
		IdAcum:     acum.IdAcum,
	})
	acum.Resultado += elem.Resultado

	return acum
}

func TotalHaberesNoImponible(rut string, Id uint, elem Elemento, acum acumulador) acumulador {

	if acum.NomAcum != "TotalHaberesNoImponible" {
		acum.IdPersona = Id
		acum.Rut = rut
		acum.IdAcum = 5
		acum.NomAcum = "TotalHaberesNoImponible"
		//elementos    []operador,
		acum.Resultado = 0.0
	}

	acum.Elementos = append(acum.Elementos, operador{
		Element:    elem, // Puedes inicializar elemento según tus necesidades
		Operacion:  "Suma",
		IdElemento: elem.IdElemento,
		IdAcum:     acum.IdAcum,
	})
	acum.Resultado += elem.Resultado

	return acum

}

func TotalHaberesTributales(rut string, Id uint, elem Elemento, acum acumulador) acumulador {

	if acum.NomAcum != "TotalHaberesTributable" {
		acum.IdPersona = Id
		acum.Rut = rut
		acum.IdAcum = 4
		acum.NomAcum = "TotalHaberesTributable"
		acum.Resultado = 0.0
	}

	acum.Elementos = append(acum.Elementos, operador{
		Element:    elem, // Puedes inicializar elemento según tus necesidades
		Operacion:  "Suma",
		IdElemento: elem.IdElemento,
		IdAcum:     acum.IdAcum,
	})
	acum.Resultado += elem.Resultado

	return acum
}

func TotalDescuentosObligatorios(rut string, Id uint, elem Elemento) *acumulador {
	var DescuentosObligatorios *acumulador
	if DescuentosObligatorios == nil {

		DescuentosObligatorios := &acumulador{
			IdPersona: Id,
			Rut:       rut,
			IdAcum:    4,
			NomAcum:   "DescuentosObligatorios",
			//elementos    []operador,
			Resultado: 0.0,
		}
		DescuentosObligatorios.Elementos = append(DescuentosObligatorios.Elementos, operador{
			Element:    elem, // Puedes inicializar elemento según tus necesidades
			IdElemento: elem.IdElemento,
			IdAcum:     DescuentosObligatorios.IdAcum,
			Operacion:  "Suma",
		})
		DescuentosObligatorios.Resultado += elem.Resultado
	} else {
		DescuentosObligatorios.Elementos = append(DescuentosObligatorios.Elementos, operador{
			Element:   elem, // Puedes inicializar elemento según tus necesidades
			Operacion: "Suma",
		})
		DescuentosObligatorios.Resultado += elem.Resultado
	}
	return DescuentosObligatorios
}

func TotalDescuentosVoluntarios(rut string, Id uint, elem Elemento) *acumulador {
	var DescuentosVoluntarios *acumulador
	if DescuentosVoluntarios == nil {

		DescuentosVoluntarios := &acumulador{
			IdPersona: Id,
			Rut:       rut,
			IdAcum:    6,
			NomAcum:   "DescuentosVoluntarios",
			//elementos    []operador,
			Resultado: 0.0,
		}
		DescuentosVoluntarios.Elementos = append(DescuentosVoluntarios.Elementos, operador{
			Element:    elem, // Puedes inicializar elemento según tus necesidades
			IdElemento: elem.IdElemento,
			IdAcum:     DescuentosVoluntarios.IdAcum,
			Operacion:  "Suma",
		})
		DescuentosVoluntarios.Resultado += elem.Resultado
	} else {
		DescuentosVoluntarios.Elementos = append(DescuentosVoluntarios.Elementos, operador{
			Element:    elem, // Puedes inicializar elemento según tus necesidades
			IdElemento: elem.IdElemento,
			IdAcum:     DescuentosVoluntarios.IdAcum,
			Operacion:  "Suma",
		})
		DescuentosVoluntarios.Resultado += elem.Resultado
	}
	return DescuentosVoluntarios
}
