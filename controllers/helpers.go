package controllers

import (
	"github.com/gorilla/schema"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func ParseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}
func parseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}
func parseValues(values url.Values, dst interface{}) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(dst, values); err != nil {
		return err
	}
	return nil
}

func uploadPicture(r *http.Request, idPicture string) (string, error) {
	file, _, err := r.FormFile(idPicture)
	if err != nil {
		return "", err
	}
	defer file.Close()

	rrand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Genera un número entero aleatorio entre 0 y 100.
	// Verificar si el directorio existe
	//TODO cambiar por el directorio de la carpeta public
	if _, err := os.Stat("/home/runner/NEFTFRONT-2/assets/images/products/"); os.IsNotExist(err) {
		// Si no existe, crear el directorio
		if err := os.MkdirAll("/home/runner/NEFTFRONT-2/assets/images/products/", os.ModePerm); err != nil {
			// Manejar el error si no se puede crear el directorio
			return "", err
		}
	}
	numPicture := rrand.Intn(1000000)
	namePicture := "upload-" + strconv.Itoa(numPicture) + ".png"
	newPicture, err := os.Create("/home/runner/NEFTFRONT-2/assets/images/products/" + namePicture)
	if err != nil {
		return "", err
	}
	defer newPicture.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	// write this byte array to our temporary file
	newPicture.Write(fileBytes)
    return namePicture, nil
}
