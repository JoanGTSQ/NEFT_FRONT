$(document).ready(function () {

    $('#giaFrame').on('load', function() {
        console.log("loading DOM")
        var iframeDoc = document.getElementById('giaFrame').contentWindow;
        console.log("setting events")
        setEvents(iframeDoc)
    });
    const frameGia = document.querySelector('#giaFrame');
    frameGia.addEventListener('change', (event) => {
        console.log("test")
        $('#frameWeb').value = "test"
    });
    const selectElement = document.querySelector('#frameWeb');
    selectElement.addEventListener('change', (event) => {
        $('#giaFrame').attr('src',`${event.target.value}`)
    });

})
var number = 1

function setEvents(element){
    setTimeout(() => {
        console.log("Setting events...")

        $(element).mouseover(function (event) {
            $(event.target).addClass('outline-element');
        }).mouseout(function (event) {
            $(event.target).removeClass('outline-element');
        }).click(function (event) {
            $(event.target).toggleClass('outline-element-clicked')
            var giaFrame = document.getElementById("giaFrame");
            // alert("XPATH located: " +  + "\nObject: " + event.target.type+"\ntagName: " + event.target.tagName)
            document.getElementById('list-gia').insertAdjacentHTML("beforeend",'' +
                '<div class="block-gia" xmlns="http://www.w3.org/1999/html"> ' +
                '<div class="title-block">' +
                '<input form="giaForm" name="Block[Title]" type="text" value="STEP '+number+'" style="width: 30%;"/>' +
                '</div>' +
                '<div class="body-gia">' +
                '<input type=hidden form="giaForm" name="Block[TimeStamp]" value="'+Date.now()+'"/>' +
                '<input type=hidden form="giaForm" name="Block[Xpath]" value="'+getElementXpath(event.target)+'"/>' +
                '<input type=hidden form="giaForm" name="Block[Url]" value="'+giaFrame.getAttribute("src")+'"/>' +
                '<input type=hidden form="giaForm" name="Block[Value]" value="'+event.target.textContent+'"/>' +
                '<input type=text form="giaForm" name="Block[Expected]" value="" placeholder="Expected value"/>' +
                '</div>' +
                '</div>');
            elem = document.getElementById('list-gia')
            elem.scrollTop = elem.scrollHeight
            number++
        })
    }, 500);
}
function getElementXpath(element){
    return "//" + $(element).parents().addBack().map( function (){
        var $this = $(this);
        var tagName = this.nodeName
        if ($this.siblings(tagName).length > 0) {
            tagName += "[" + ($this.prevAll(tagName).length + 1) + "]"
        }
        return tagName;
    }).get().join("/").toLowerCase();
}
