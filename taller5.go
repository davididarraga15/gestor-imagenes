package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type DatosPagina struct {
	Images   []ImagenBase64
	HostName string
}

type ImagenBase64 struct {
	Encoding template.URL
	Nombre   string
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/index.html")

	if err != nil {
		fmt.Fprint(w, "Página no encontrada")
	} else {
		carpeta := os.Args[2] // Obtener el directorio de los argumentos de la línea de comandos
		directorio, err := os.Open(carpeta)
		check(err)
		defer directorio.Close()
		nombres, err := directorio.Readdirnames(0)
		check(err)

		var archivos []string
		for _, nombre := range nombres {
			if strings.HasSuffix(nombre, ".jpg") || strings.HasSuffix(nombre, ".png") ||
				strings.HasSuffix(nombre, ".jpeg") || strings.HasSuffix(nombre, ".JPG") {
				archivos = append(archivos, nombre)
			}
		}

		fmt.Println("Cantidad de archivos en la carpeta:", len(archivos))

		var imagenAleatoria = archivos[rand.Intn(len(archivos)-1)]

		fmt.Println("Imagen aleatoria escogida del directorio:", imagenAleatoria)

		hostname, err := os.Hostname()
		check(err)

		fmt.Println("Nombre del host:", hostname)

		miMapa := make(map[int]string)
		var listaGenerada []ImagenBase64
		for i := 0; i < 4; i++ {
			indice := rand.Intn(len(archivos))
			existe := miMapa[indice]
			if existe == "" {
				miMapa[indice] = archivos[indice]
				f, err := os.ReadFile(carpeta + miMapa[indice])
				check(err)
				src := "data:image/jpg;base64," + base64.StdEncoding.EncodeToString(f)
				image := ImagenBase64{
					Encoding: template.URL(src),
					Nombre:   miMapa[indice],
				}
				listaGenerada = append(listaGenerada, image)
			} else {
				i--
			}
		}

		data := DatosPagina{
			Images:   listaGenerada,
			HostName: hostname,
		}

		err = tmpl.Execute(w, data)
	}
}

func main() {

	puerto := os.Args[1]
	directorio := os.Args[2]

	// Verifica si el directorio existe
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		fmt.Println("El directorio especificado no existe:", directorio)
		return
	}

	// Verifica si se puede abrir el directorio
	_, err := os.Open(directorio)
	if err != nil {
		fmt.Println("Error al abrir el directorio:", err)
		return
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)

	fmt.Println("Servidor escuchando en el puerto", puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, nil))
}
