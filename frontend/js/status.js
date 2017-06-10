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

document.querySelector("button[name=get_status]").onclick = function () {
    let uuid = document.querySelector("input[name=uuid]").value;

    document.querySelector("div#status-box").style = "display:block";
    let body = document.querySelector("div#status-box > .panel-body")

    jQuery.ajax({
        url: config.apiHost+"status/"+ uuid
    }).done(function(data) {
        document.querySelector("div#status-box").className = "panel panel-default";
        body.innerHTML = '<strong>UUID:</strong> ' + data.uuid + `<br />
        <strong>Upload Time:</strong> ` + data.upload_time + `<br />
        <strong>Max Downloads:</strong> ` + data.max_downloads + `<br />
        <strong>Download Count:</strong> ` + data.download_count + `<br />
        <strong>LOG:</strong><p>` + data.message;  
        body.innerHTML += '<a class="btn btn-warning" href="download.html?id='+uuid+'">Download</a>';

    }).fail(function(err) {
        document.querySelector("div#status-box").className = "panel panel-danger";
        body.innerHTML = JSON.stringify(err);
    });

 };