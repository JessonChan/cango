function showLink() {
    document.getElementById("link").setAttribute("style", "display:block")
}

function requestJSON(url, data) {
    var httpRequest = new XMLHttpRequest();
    httpRequest.open('POST', url, false);
    httpRequest.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    httpRequest.send(data);
    if (httpRequest.readyState === 4 && httpRequest.status === 200) {
        alert(httpRequest.responseText)
    }
}