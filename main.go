package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"lazyadmin"
	"log"
	"os"
)

// Alumno Registro para almacenar los alumnos del formulario de recogida
// de datos
type Alumno struct {
	TimeStamp string `json:"timestamp"`
	Curso     string `json:"curso"`
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	Email     string `json:"email"`
	Fnac      string `json:"fecha_nac"`
	Cisco     bool   `json:"cisco"`
}

func main() {

	filePtr := flag.String("f", "", "Nombre del fichero .csv")
	prefixPtr := flag.String("x", "", "Prefijo para los nombres de usurio")
	passwdPtr := flag.String("p", "", "Contraseña para los usuarios")
	grupoPtr := flag.String("g", "", "Grupo al que pertenece los usurios")

	flag.Parse()

	if *filePtr == "" || *passwdPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	csvFile, error := os.Open(*filePtr)
	if error != nil {
		panic(error)
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))

	var alumnos []Alumno

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		alumno := Alumno{
			TimeStamp: line[0],
			Curso:     line[1],
			Nombre:    line[2],
			Apellidos: line[3],
			Email:     line[4],
			Fnac:      line[5],
			Cisco: func(s string) bool {
				if s == "Sí" {
					return true
				}

				return false

			}(line[6]),
		}
		alumnos = append(alumnos, alumno)
	}

	var usuarios []lazyadmin.User
	var usuario lazyadmin.User

	for _, v := range alumnos[1:] {
		usuario.NewUser(*prefixPtr+"_", v.Nombre, v.Apellidos, *passwdPtr, "", "", "", "", "", "", "", "", *grupoPtr)
		usuarios = append(usuarios, usuario)
	}

	for _, v := range usuarios {
		fmt.Println(v.String())
	}
}
