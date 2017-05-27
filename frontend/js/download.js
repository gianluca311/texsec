document.querySelector("button[name=get_download").onclick = function() {
    let uuid = document.querySelector("input[name=uuid]").value;
    window.location.href = config.apiHost+"/download/"+uuid;
};