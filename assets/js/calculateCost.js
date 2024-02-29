document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('calcularPrecioRelativoBtn').addEventListener('click', function(event) {
        event.preventDefault();

        var productSelect = document.getElementById('products');
        var materialSelect = document.getElementById('material');
        var costInput = document.getElementById('costInput');
        var sellingPriceInput = document.getElementById('sellingPriceInput');

        var selectedProductIndices = Array.from(productSelect.selectedOptions).map(option => option.index);
        var selectedMaterialIndex = materialSelect.selectedIndex;

        var totalProductWeight = selectedProductIndices.reduce(function(sum, index) {
            var productWeight = parseFloat(productSelect.options[index].getAttribute('data-weight'));
            return sum + (isNaN(productWeight) ? 0 : productWeight);
        }, 0);

        var pricePer1000g = parseFloat(materialSelect.options[selectedMaterialIndex].getAttribute('data-price-per-1000g'));
        var totalProductSellingPrice = selectedProductIndices.reduce(function(sum, index) {
            var productSellingPrice = parseFloat(productSelect.options[index].getAttribute('data-selling-price'));
            return sum + (isNaN(productSellingPrice) ? 0 : productSellingPrice);
        }, 0);

        if (!isNaN(totalProductWeight) && !isNaN(pricePer1000g)) {
            var costoRelativo = (totalProductWeight / 1000) * pricePer1000g;
            var precioVentaRelativo = totalProductSellingPrice;

            // Mostrar los resultados en los campos de entrada
            costInput.value = costoRelativo.toFixed(2);
            sellingPriceInput.value = precioVentaRelativo.toFixed(2);
        } else {
            console.error('Los valores del producto, material, peso o precio por 1000g no son números válidos.');
        }
    });
});
