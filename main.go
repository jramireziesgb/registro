package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"github.com/jramireziesgb/lazyadmin"
)

// Alumno Registro para almacenar los alumnos del formulario de recogida
// de datos
type Alumno struct {
	TimeStamp string
	Prefijo   string
	Curso     string
	Nombre    string
	Apellidos string
	Email     string
	Fnac      string
	Cisco     bool
}

const AppVersion = "1.0"

func main() {

	prefijos := map[string]string{
		"1º SMR":            "s1",
		"2º SMR":            "s2",
		"1º GA":             "g1",
		"2º GA":             "g2",
		"1º FPB":            "f1",
		"2º FPB":            "f2",
		"1º Bachillerato A": "b1a",
		"1º Bachillerato B": "b1b",
		"1º Bachillerato C": "b1c",
		"2º Bachillerato":   "b2",
	}

	grupos := map[string]string{
		"1º SMR":            "smr1",
		"2º SMR":            "smr2",
		"1º GA":             "ga1",
		"2º GA":             "ga2",
		"1º FPB":            "fpb1",
		"2º FPB":            "fpb2",
		"1º Bachillerato A": "bac1a",
		"1º Bachillerato B": "bac1b",
		"1º Bachillerato C": "bac1c",
		"2º Bachillerato":   "bac2",
	}

	filePtr := flag.String("f", "", "Nombre del fichero .csv")
	prefixPtr := flag.String("x", "", "Prefijo para los nombres de usurio")
	passwdPtr := flag.String("p", "", "Contraseña para los usuarios")
	versionPtr := flag.Bool("v", false, "Muestra la versión del programa")

	flag.Parse()

	if *versionPtr {
		fmt.Println("Versión", AppVersion)
		os.Exit(0)
	}

	if *filePtr == "" {
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

		prefix := ""

		if *prefixPtr == "" {
			prefix = prefijos[line[1]]
		} else {
			prefix = *prefixPtr
		}

		alumno := Alumno{
			TimeStamp: line[0],
			Prefijo:   prefix,
			Curso:     grupos[line[1]],
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
		usuario.NewUser(v.Prefijo+"_", v.Nombre, v.Apellidos, *passwdPtr, "", "", "", "", "", "", "", "", v.Curso)
		usuarios = append(usuarios, usuario)
	}

	for _, v := range usuarios {
		fmt.Println(v.String())
	}
}
