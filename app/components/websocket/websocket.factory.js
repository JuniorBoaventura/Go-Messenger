(function() {
  'use strict';

  angular
    .module('app')
    .factory('WebsocketFactory', WebsocketFactory);

  function WebsocketFactory() {

    var ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen    = function() {
      sendRequest({body: 'junior', type: 'connect'});
      console.log('Connection established!');
    };

    ws.onmessage = function(res) {
      let json = JSON.parse(res.data);
      console.log(json);
    };

    function sendRequest(request) {
      ws.send(JSON.stringify(request));
    };

    var service = {
      sendRequest: sendRequest
    };
    return service;

  }

})();
