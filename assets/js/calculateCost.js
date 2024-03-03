document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('calcularPrecioRelativoBtn').addEventListener('click', function (event) {
        event.preventDefault();
        console.log('Botón clicado'); 
        var dynamicRowsContainer = document.getElementById('dynamicRowsContainer');
        var costInputs = dynamicRowsContainer.querySelectorAll('.cost');
        var sellingPriceInputs = dynamicRowsContainer.querySelectorAll('.sale');
    
        costInputs.forEach(function (costInput, index) {
            var productSelect = costInput.closest('tr').querySelector('.products');
            var materialSelect = costInput.closest('tr').querySelector('.materialID');
    
            var selectedProductIndices = Array.from(productSelect.selectedOptions).map(option => option.index);
            var selectedMaterialIndex = materialSelect.selectedIndex;
    
            var totalProductWeight = selectedProductIndices.reduce(function (sum, index) {
                var productWeight = parseFloat(productSelect.options[index].getAttribute('data-weight'));
                return sum + (isNaN(productWeight) ? 0 : productWeight);
            }, 0);
    
            var pricePer1000g = parseFloat(materialSelect.options[selectedMaterialIndex].getAttribute('data-price-per-1000g'));
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
                console.log(costoRelativo);
                console.log(precioVentaRelativo);
            } else {
                console.error('Los valores del producto, material, peso o precio por 1000g no son números válidos.');
            }
        });
    });
    var dynamicRowsContainer = document.getElementById('dynamicRowsContainer');
    var addRowBtn = document.getElementById('addRowBtn');
    var exampleRow = document.querySelector('.white-space-no-wrap'); // Ajusta el selector según tu estructura real

        addRowBtn.addEventListener('click', function () {
            // Clonar la fila de ejemplo
            var newRow = exampleRow.cloneNode(true);
            newRow.classList.remove('example-row'); // Elimina la clase 'example-row' si está presente
    
            // Añadir la nueva fila a la tabla
            dynamicRowsContainer.querySelector('tbody').appendChild(newRow);
        });

});
