// Get the template HTML and remove it from the doumenthe template HTML and remove it from the doument
var previewNode = document.querySelector("#template");
previewNode.id = "";
var previewTemplate = previewNode.parentNode.innerHTML;
previewNode.parentNode.removeChild(previewNode);

var myDropzone = new Dropzone("div#my-awesome-dropzone", { // Make the whole body a dropzone
  url: config.apiHost+"upload", // Set the url
  thumbnailWidth: 80,
  thumbnailHeight: 80,
  parallelUploads: 20,
  previewTemplate: previewTemplate,
  autoQueue: false, // Make sure the files aren't queued until manually added
  previewsContainer: "div#previews", // Define the container to display the previews
  maxFiles: 1,
  dictDefaultMessage: "Drop files here",
  method: "POST",
  headers: {
    'Cache-Control': null,
    'X-Requested-With': null
  }
});

myDropzone.on("addedfile", function(file) {
  file.previewElement.querySelector(".start").onclick = function() { myDropzone.enqueueFile(file); };
});

// Update the total progress bar
myDropzone.on("totaluploadprogress", function(progress) {
  document.querySelector("#total-progress .progress-bar").style.width = progress + "%";
});

myDropzone.on("sending", function(file, xhr, formData) {
  // Show the total progress bar when upload starts
  // And disable the start button
  file.previewElement.querySelector(".start").setAttribute("disabled", "disabled");
  formData.append("max_downloads", file.previewElement.querySelector("input[name=max_downloads]").value);
});

// Hide the total progress bar when nothing's uploading anymore
myDropzone.on("queuecomplete", function(progress) {
  document.querySelector("#total-progress").style.opacity = "0";
});

myDropzone.on("error", function(error, xhr) {
  error.previewElement.querySelector("strong.error").innerText = xhr.message || xhr;
});

myDropzone.on("success", function(success, xhr) {
  window.location.href = "status.html?id="+xhr.uuid;
});
