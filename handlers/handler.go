package handlers

import (
	"fmt"
	"net/http"
)

func GetEmpleado(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(rw, "Lista todos los empleados")
}
