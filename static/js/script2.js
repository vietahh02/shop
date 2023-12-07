const toggleButton = document.getElementById("toggleButton");
const hoverTarget = document.getElementById("toggleButton");
const info = document.querySelector(".info-cate");

let isInfoVisible = true;

toggleButton.addEventListener("click", function () {
  if (isInfoVisible) {
    info.style.display = "grid";
  } else {
    info.style.display = "none";
  }
  isInfoVisible = !isInfoVisible; 
});

let isHovering = false;

hoverTarget.addEventListener("mouseenter", function () {
  isHovering = true;
  info.style.display = "grid";
});

info.addEventListener("mouseenter", function () {
  isHovering = true;
  info.style.display = "grid";
});

hoverTarget.addEventListener("mouseleave", function () {
  isHovering = false;
  if (isInfoVisible) {
    hideInfo();
  }
});

info.addEventListener("mouseleave", function () {
  isHovering = false;
  if (isInfoVisible) {
    hideInfo();
  }
});

function hideInfo() {
  if (!isHovering) {
    info.style.display = "none";
  }
}



