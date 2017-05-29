function getParameterByName(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

document.querySelector("input[name=uuid]").value = getParameterByName("id");

document.querySelector("button[name=get_download").onclick = function() {
    let uuid = document.querySelector("input[name=uuid]").value;
    window.location.href = config.apiHost+"/download/"+uuid;
};