package config

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq" // Importa el controlador de Postgres
)

const url = "host=localhost user=postgres password=T3st1ng dbname=nomina port=5432 sslmode=disable"

type Moneda struct {
	TMoneda  string    `mapstructure:"tmoneda"`
	Fecha    time.Time `mapstructure:"fecha"`
	ValorCLP float64   `mapstructure:"valorclp"`
}

type Afpdb struct {
	Idafp  uint
	Codigo string
	Nombre string
	Valor  sql.NullFloat64
	Moneda sql.NullString
}
type Afpsdb []Afpdb

var (
	Monedas []Moneda
	DB      *sql.DB
	once    sync.Once
)

func initConfig() {
	// Aquí puedes inicializar la configuración usando Viper si es necesario
}

func initDB() {
	// Configura y abre la conexión a la base de datos
	var err error
	DB, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("Error al abrir la conexión a la base de datos: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	loadMonedasFromDB()
}

func loadMonedasFromDB() []Moneda {
	rows, err := DB.Query("SELECT moneda, fecha, valorclp FROM monedas m WHERE fecha = (SELECT max(m1.fecha) FROM monedas m1 WHERE m1.moneda = m.moneda AND m1.fecha <= NOW())")
	if err != nil {
		log.Fatalf("Error al consultar la base de datos: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var moneda Moneda
		var tmoneda string
		var fecha time.Time
		var valorCLP float64
		if err := rows.Scan(&tmoneda, &fecha, &valorCLP); err != nil {
			log.Fatalf("Error al escanear la fila: %v", err)
		}

		moneda.TMoneda = tmoneda
		moneda.Fecha = fecha
		moneda.ValorCLP = valorCLP
		Monedas = append(Monedas, moneda)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error al iterar sobre las filas: %v", err)
	}
	return Monedas
}

/*
	func loadTablaAfpFromDB() []Afp {
		rows, err := DB.Query("SELECT moneda, fecha, valorclp FROM monedas m WHERE fecha = (SELECT max(m1.fecha) FROM monedas m1 WHERE m1.moneda = m.moneda AND m1.fecha <= NOW())")
		if err != nil {
			log.Fatalf("Error al consultar la base de datos: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var moneda Moneda
			var tmoneda string
			var fecha time.Time
			var valorCLP float64
			if err := rows.Scan(&tmoneda, &fecha, &valorCLP); err != nil {
				log.Fatalf("Error al escanear la fila: %v", err)
			}

			moneda.TMoneda = tmoneda
			moneda.Fecha = fecha
			moneda.ValorCLP = valorCLP
			Monedas = append(Monedas, moneda)
		}

		if err := rows.Err(); err != nil {
			log.Fatalf("Error al iterar sobre las filas: %v", err)
		}
		return Afps
	}
*/
func Init() {
	once.Do(func() {
		initConfig()
		initDB()
	})
}
