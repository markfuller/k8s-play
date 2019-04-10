//from https://kubernetes.io/docs/tutorials/hello-minikube/
var http = require('http');

var handleRequest = function(request, response) {
  console.log('Address service received request for URL: ' + request.url);
  response.writeHead(200);
  response.end('29 Acacia Road\n');
};
var www = http.createServer(handleRequest);
www.listen(8081);