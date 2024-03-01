package controllers

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/schema"
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

func uploadPicture(r *http.Request, idPicture, typeFile string) (string, error) {
	file, _, err := r.FormFile(idPicture)
	if err != nil {
		return "", err
	}
	defer file.Close()

	rrand := rand.New(rand.NewSource(time.Now().UnixNano()))
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
    projectRoot := filepath.Join(basepath, "../")
	// Genera un número entero aleatorio entre 0 y 100.
	// Verificar si el directorio existe
	//TODO cambiar por el directorio de la carpeta public
	if _, err := os.Stat(projectRoot + "/assets/images/products/"); os.IsNotExist(err) {
		// Si no existe, crear el directorio
		if err := os.MkdirAll(projectRoot+"/assets/images/products/", os.ModePerm); err != nil {
			// Manejar el error si no se puede crear el directorio
			return "", err
		}
	}
    if _, err := os.Stat(projectRoot + "/assets/stl/products/"); os.IsNotExist(err) {
        // Si no existe, crear el directorio
        if err := os.MkdirAll(projectRoot+"/assets/stl/products/", os.ModePerm); err != nil {
            // Manejar el error si no se puede crear el directorio
            return "", err
        }
    }
	numFile := rrand.Intn(1000000)
	nameFile := "upload-" + strconv.Itoa(numFile) + ".png"
	var path string
	switch typeFile {
	case "productPicture":
		path = "images/products/"
	case "productSTL":
		path = "stl/products/"
	default:
		return "", errors.New("typeFile not valid")
	}
	newFile, err := os.Create(projectRoot + "/assets/" + path + nameFile)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	// write this byte array to our temporary file
    newFile.Write(fileBytes)
	return nameFile, nil
}
