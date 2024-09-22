package views

import (
	"errors"
	"fmt"
	"html/template"
	"jgt.solutions/models"
	"math"
	"strconv"
	"time"
)

var funcMap = template.FuncMap{
	"minus":          minus,
	"percentage":     percentage,
	"csrfField":      csrfField,
	"formatDate":     formatDate,
	"toMoney":        toMoney,
	"plus":           plus,
	"getMainPicture": getMainPicture,
}

func minus(a, b float64) string {
	return strconv.FormatFloat(a-b, 'f', 2, 64)
}
func percentage(a, b float64) string {
	return fmt.Sprintf("%.2f", math.Round((a/b)*100)/10)
}

func csrfField() (template.HTML, error) {
	return "", errors.New("csrf is not implemented")
}
func formatDate(date time.Time) string {
	return date.Format("02/01/2006")
}
func toMoney(amount int) string {
	// Convertir a float64
	decimalValue := float64(amount) / 100.0

	// Formatear como moneda con dos decimales y el símbolo €
	return fmt.Sprintf("%.2f €", decimalValue)
}

func plus(valores ...int) int {
	suma := 0
	for _, valor := range valores {
		suma += valor
	}
	return suma
}

func getMainPicture(product models.Product) string {
	return product.GetMainPicture()
}
