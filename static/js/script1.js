const image = document.querySelector("img"),
  input = document.querySelector("input");

input.addEventListener("change", () => {
  image.src = URL.createObjectURL(input.files[0]);
  document.getElementById("checkImg").value = input.files[0].name
});