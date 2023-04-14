var msg = document.querySelector("input");
var form = document.querySelector("form");
form.addEventListener("submit", handleSubmit);
document.getElementById;
function handleSubmit(event) {
  event.preventDefault();
  console.log(event);
}
msg.addEventListener("input", handleInput);
function handleInput(event) {
  console.log(event);
}
var h1 = document.createElement("h1");
document.body.appendChild(h1);
h1.innerText = "hello";
