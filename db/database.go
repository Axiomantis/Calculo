package db

import (
	"fmt"

	"database/sql"

	_ "github.com/lib/pq"
)

type Personadb struct {
	IdPersona uint
	Rut       string
	Nombre    string
	Idafp     sql.NullInt32
}

type Personasdb []Personadb

type Afpdb struct {
	Idafp  uint
	Codigo string
	Nombre string
	Valor  sql.NullFloat64
	Moneda sql.NullString
}
type Afpsdb []Afpdb

type ElemEmpldb struct {
	Idpersona  uint
	Idelemento uint
	Cantidad   sql.NullFloat64
	Valor      sql.NullFloat64
	Moneda     sql.NullString
	Formula    sql.NullString
	Resultado  sql.NullFloat64
}
type ElemEmpldbs []ElemEmpldb

type ElementDb struct {
	Idelemento uint
	Nombre     string
	Descr      string
	Tipo       string
	Calidad    string
	Cantidad   sql.NullFloat64
	Valor      sql.NullFloat64
	Valmax     sql.NullFloat64
	Valmin     sql.NullFloat64
	Moneda     sql.NullString
	Formula    sql.NullString //string
	Resultado  sql.NullFloat64
}

type ElementDbs []ElementDb

type MonedaDb struct {
	Moneda   string
	Fecha    sql.NullTime
	Valorclp sql.NullFloat64
}
type MonedaDbs []MonedaDb

var db *sql.DB

const url = "host=localhost user=postgres password=T3st1ng dbname=nomina port=5432 sslmode=disable"

func Connect() {
	conection, err := sql.Open("postgres", url)
	if err != nil {
		fmt.Printf("Error de conect  %s", err)
		//log.Fatal(err)
	} else {
		db = conection
		fmt.Println("Conexión exitosa")

	}

}
func Close() {
	db.Close()

}
func ListarPersonas() (Personasdb, Afpsdb, ElemEmpldbs, ElementDbs) {

	Connect()
	sql := "select idpersona, rut ,	nombre,idafp from personas"
	personasdb := Personasdb{}
	rows, _ := db.Query(sql)

	for rows.Next() {
		personadb := Personadb{}
		err := rows.Scan(&personadb.IdPersona, &personadb.Rut, &personadb.Nombre, &personadb.Idafp)
		if err != nil {
			fmt.Printf("error ListarPersonas-->%s \n", err)
		}
		personasdb = append(personasdb, personadb)
		//fmt.Println(personasdb)
	}
	defer Close()

	afpsdb := ListarAfpsDb()
	elemEmpldbs := ListarElemEmpleado()
	elementDbs := ListarElementos()

	//type personas []models.Persona
	/*
		IdPersona  uint `gorm:"primaryKey"`
		Rut        string
		Nombre     string
		Afp        Afp          //`gorm:"foreignKey:IdAfp"`
		Salud      Salud        //`gorm:"foreignKey:IdSalud"`
		SaludComp  []saludCompl //`gorm:"foreignKey:IdSalComp"`
		Apv        Apv          //`gorm:"foreignKey:IdApv"`
		Haberes    []elemento   //`gorm:"foreignKey:IdElemento"`
		Descuentos []elemento   //`gorm:"foreignKey:IdElemento"`
		Acumulador []acumulador //`gorm:"foreignKey:IdAcum"`
	*/
	/*
		for _, per := range personasdb {
			persona := models.Persona{}
			persona.IdPersona = per.IdPersona
			persona.Rut = per.Rut
			persona.Nombre = per.Nombre
			//persona.Afp=per
			personas = append(personas, persona)

		}
		fmt.Printf("valores -->%v", personas)

	*/
	//fmt.Printf("valores return -->%v \n", elemEmpldbs)
	//fmt.Printf("valores return -->%v \n", personasdb)

	return personasdb, afpsdb, elemEmpldbs, elementDbs
}

func ListarAfpsDb() Afpsdb {
	Connect()
	sql := "select idafp,codigo,nombre,moneda,valor from afps"
	afpsdb := Afpsdb{}
	rows, _ := db.Query(sql)

	for rows.Next() {
		afpdb := Afpdb{}
		err := rows.Scan(&afpdb.Idafp, &afpdb.Codigo, &afpdb.Nombre, &afpdb.Moneda, &afpdb.Valor)
		if err != nil {
			fmt.Printf("error ListarAfpsDb-->%s \n", err)
		}
		afpsdb = append(afpsdb, afpdb)
		//fmt.Println(afpsdb)
	}
	defer Close()
	return afpsdb

}

func ListarMonedasDb() MonedaDbs {
	Connect()
	sql := " select m.moneda, m.fecha,m.valorclp from monedas m where fecha =( select max(m1.fecha) from monedas m1 where m1.moneda=m.moneda and m1.fecha <=now() )"
	monedaDbs := MonedaDbs{}
	rows, _ := db.Query(sql)

	for rows.Next() {
		monedadb := MonedaDb{}
		err := rows.Scan(&monedadb.Moneda, &monedadb.Valorclp)
		if err != nil {
			fmt.Printf("error ListarMonedasDb-->%s \n", err)
		}
		monedaDbs = append(monedaDbs, monedadb)
		//fmt.Println(afpsdb)
	}
	defer Close()
	return monedaDbs

}

func ListarElemEmpleado() ElemEmpldbs {
	Connect()

	sql := `select
		p.idpersona	,p.idelemento	,p.cantidad	,p.valor	, p.moneda	, COALESCE(p.formula,q.formula) as formula,
		case when p.resultado = 0 then
		q.resultado 
		else
		p.resultado
		end
		as resultado
		from  public.elementempl p left outer join public.elementos q
		on p.idelemento=q.idelemento
		union
		select b.idpersona,a.idelemento,a.cantidad,valor,moneda,formula,resultado from public.elementos a
		inner join public.personas b
		on 1=1
		where not exists (select 1 from public.elementempl d where d.idpersona=b.idpersona and d.idelemento=a.idelemento  )
`
	/*
		sql := "Select idpersona,idelemento,cantidad,	valor,	moneda,	formula,resultado from elementempl"
	*/
	elemEmpls := ElemEmpldbs{}
	rows, _ := db.Query(sql)

	for rows.Next() {
		elemEmpl := ElemEmpldb{}
		err := rows.Scan(&elemEmpl.Idpersona, &elemEmpl.Idelemento, &elemEmpl.Cantidad, &elemEmpl.Valor, &elemEmpl.Moneda, &elemEmpl.Formula, &elemEmpl.Resultado)
		if err != nil {
			fmt.Printf("error ListarElemEmpleado-->%s \n", err)
		}

		elemEmpls = append(elemEmpls, elemEmpl)
	}
	defer Close()
	return elemEmpls
}
func ListarElementos() ElementDbs {
	Connect()
	sql := "Select idelemento,nombre,descr,tipo,calidad,cantidad,valor,valmax,valmin,moneda,formula,resultado from elementos"
	elementDbs := ElementDbs{}
	rows, _ := db.Query(sql)

	for rows.Next() {
		elementDb := ElementDb{}

		err := rows.Scan(&elementDb.Idelemento, &elementDb.Nombre, &elementDb.Descr, &elementDb.Tipo, &elementDb.Calidad, &elementDb.Cantidad, &elementDb.Valor, &elementDb.Valmax, &elementDb.Valmin, &elementDb.Moneda, &elementDb.Formula, &elementDb.Resultado)
		if err != nil {
			fmt.Printf("error ListarElementos-->%s \n", err)
		}
		fmt.Printf("%s--->%v \n", elementDb.Nombre, elementDb.Formula)
		elementDbs = append(elementDbs, elementDb)
	}
	defer Close()

	return elementDbs
}
