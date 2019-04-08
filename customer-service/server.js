var http = require('http');

var handleRequest = function (request, response) {
  console.log('Received request for URL: ' + request.url);
  response.writeHead(200);

  //TODO this address and port should not be hardcoded
  // http.get('http://localhost:8081', (resp) => {
  http.get('http://my-address-service.local-docker-registry-test:8000', (resp) => {
    let data = ''
    resp.on('data', (chunk) => {
      data += chunk;
    });
    resp.on('end', () => {
      console.log("data is " + data);
      response.end("Customer Eric at address " + data + ".\n")
    });
  }).on("error", (err) => {
    console.log("Error: " + err.message);
    response.end("ERROR: " + err.message);
  });
}

var www = http.createServer(handleRequest);
www.listen(8080);
