package views
import (
    "math"
    "fmt"
    "strconv"
    "time"
    "html/template"
    "errors"
)
var funcMap = template.FuncMap{
    "minus": minus,
    "percentage": percentage,
    "csrfField": csrfField,
    "formatDate": formatDate,
    "toMoney": toMoney,
}
func minus(a, b float64) string {
    return strconv.FormatFloat(a-b, 'f', 2, 64)
}
func percentage(a, b float64) string{
    return fmt.Sprintf("%.2f", math.Round((a/b)*100) / 10)
}

func csrfField() (template.HTML, error){
    return "", errors.New("csrf is not implemented")
}
func formatDate(date time.Time) string{
    return date.Format("02/01/2006")
}
func toMoney(amount int64) string{
    // Convertir a float64
    decimalValue := float64(amount) / 100.0

    // Formatear como moneda con dos decimales y el símbolo €
    return fmt.Sprintf("%.2f €", decimalValue)
}