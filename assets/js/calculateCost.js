document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('calcularPrecioRelativoBtn').addEventListener('click', function (event) {
        event.preventDefault();

        var dynamicRowsContainer = document.getElementById('dynamicRowsContainer');
        console.log("dynamicRowsContainer:", dynamicRowsContainer);
        var costInputs = dynamicRowsContainer.querySelectorAll('.cost');
        var sellingPriceInputs = dynamicRowsContainer.querySelectorAll('.sale');

        costInputs.forEach(function (costInput, index) {
            var row = costInput.closest('tr');
            var productSelect = row.querySelector('.products');
            var materialSelect = row.querySelector('.materialID');
            console.log("materialSelect:", materialSelect)
            // Obtener el índice seleccionado
            var selectedMaterialIndex = materialSelect.selectedIndex;

            // Obtener el valor del atributo data-price-per-1000g de la opción seleccionada
            var pricePer1000g = materialSelect.options[selectedMaterialIndex].getAttribute('data-price-per-1000g');

            var selectedProductIndices = Array.from(productSelect.selectedOptions).map(option => option.index);
            var totalProductWeight = selectedProductIndices.reduce(function (sum, index) {
                var productWeight = parseFloat(productSelect.options[index].getAttribute('data-weight'));
                return sum + (isNaN(productWeight) ? 0 : productWeight);
            }, 0);

            var totalProductSellingPrice = selectedProductIndices.reduce(function (sum, index) {
                var productSellingPrice = parseFloat(productSelect.options[index].getAttribute('data-selling-price'));
                return sum + (isNaN(productSellingPrice) ? 0 : productSellingPrice);
            }, 0);

            if (!isNaN(totalProductWeight) && !isNaN(pricePer1000g)) {
                var costoRelativo = (totalProductWeight / 1000) * pricePer1000g;
                var precioVentaRelativo = totalProductSellingPrice;

                // Mostrar los resultados en los campos de entrada de la fila actual
                costInput.value = costoRelativo.toFixed(2);
                sellingPriceInputs[index].value = precioVentaRelativo.toFixed(2);
            } else {
                console.error('Los valores del producto, material, peso o precio por 1000g no son números válidos.');
            }
        });
    });

    // Asegúrate de que este script se encuentra después de la generación dinámica de filas para que pueda acceder a esos elementos.






    
    var dynamicRowsContainer = document.getElementById('dynamicRowsContainer');
    var addRowBtn = document.getElementById('addRowBtn');
    var exampleRow = document.querySelector('.white-space-no-wrap'); // Ajusta el selector según tu estructura real
    var rowIndex = 1; // Contador global de filas

    addRowBtn.addEventListener('click', function () {
        // Clonar la fila de ejemplo
        var newRow = exampleRow.cloneNode(true);
        newRow.classList.remove('example-row'); // Elimina la clase 'example-row' si está presente

        // Incrementar los valores que necesitas cambiar
        newRow.querySelectorAll('.pr-0 input[type="checkbox"]').forEach(function(checkbox) {
            checkbox.checked = false; // Para que los nuevos checkboxes no estén marcados por defecto
        });

        // Actualizar los nombres de los campos en la nueva fila
        newRow.querySelectorAll('[name]').forEach(function(field) {
            var name = field.getAttribute('name');
            var matches = name.match(/\[(\d+)\]\[(\w+)\]/);
            if (matches && matches.length === 3) {
                var newIndex = parseInt(matches[1]) + rowIndex; // Usar el contador global de filas
                var newName = name.replace(matches[1], newIndex);
                field.setAttribute('name', newName);
            }
        });

        // Añadir la nueva fila a la tabla
        dynamicRowsContainer.querySelector('tbody').appendChild(newRow);

        // Incrementar el contador de filas
        rowIndex++;
    });

});
