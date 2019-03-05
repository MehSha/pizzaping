self.addEventListener('push', event => {
  console.log('[Service Worker] Push Received.');
  console.log(`[Service Worker] Push had this data: "${event.data.text()}"`);
  let data = JSON.parse(event.data.text())
  var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance
  let hostname = "http://localhost:8090"
  xmlhttp.open("POST", "/order/"+data.order_id+"/accept");
  xmlhttp.setRequestHeader("Content-Type", "application/json");
  xmlhttp.send(JSON.stringify({"restaurantID": 10}));

  const title = 'order received!';
  const options = {
    body: event.data.text(),
  };

  event.waitUntil(self.registration.showNotification(title, options));
});
