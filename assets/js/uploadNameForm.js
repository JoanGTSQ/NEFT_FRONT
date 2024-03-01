function updateFileName(input) {
    var fileName = input.files[0].name;
    var formId = input.getAttribute('name');
    var fileNameElement = document.querySelector('.selectedFileName[data-form-id="' + formId + '"]');

    if (fileNameElement) {
        fileNameElement.innerText = fileName;
    }
}