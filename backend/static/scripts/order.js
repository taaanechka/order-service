const searchForm = document.getElementById("searchOrder");
const searchButton = document.getElementById("search-submit");
const res = document.getElementById("res")

searchButton.addEventListener("click", (e) => {
    e.preventDefault();

    const orderUid = searchForm.orderuid.value;

    const req = new XMLHttpRequest();
    const url = "http://" + HOST + ":" + PORT + "/api/v1/orders/" + orderUid;
    console.log(url);
    req.open("GET", url);
    req.responseType = "json";

    req.onload = function () {
      if (req.status == 404) {
        alert("Заказ не найден")
      } else {
        res.textContent = JSON.stringify(req.response, null, 4);
      }
    };

    req.send();
})
