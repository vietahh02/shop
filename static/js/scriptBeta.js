const image = document.getElementById("img"),
  input = document.getElementById("file");

input.addEventListener("change", () => {
  image.src = URL.createObjectURL(input.files[0]);
  document.getElementById("checkImg").value = input.files[0].name
});