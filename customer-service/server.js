var http = require('http');

var handleRequest = function (request, response) {
  console.log('Customer service received request for URL: ' + request.url);
  response.writeHead(200);
  
  let addressUri = process.env.ADDRESS_URI
  console.log('Contacting address service at URI: ' + addressUri);
  http.get(addressUri, (resp) => {
    let data = ''
    resp.on('data', (chunk) => {
      data += chunk;
    });
    resp.on('end', () => {
      console.log("address data is " + data);
      response.end("Customer Eric at address " + data + ".\n")
    });
  }).on("error", (err) => {
    console.log("Error: " + err.message);
    response.end("ERROR: " + err.message);
  });
}

var www = http.createServer(handleRequest);
www.listen(8080);
