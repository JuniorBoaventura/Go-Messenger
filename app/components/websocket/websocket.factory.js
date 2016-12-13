(function() {
  'use strict';

  angular
    .module('app')
    .factory('WebsocketFactory', WebsocketFactory);

  WebsocketFactory.$inject = ['$rootScope'];

  function WebsocketFactory($rootScope) {
    var onList = {};
    var ws = false;

    function init(username) {
      ws = new WebSocket('ws://localhost:8080/ws');

      ws.onopen    = function() {
        // sendRequest({body: username, type: 'connect'});
        console.log('Connection established!');
      };

      ws.onmessage = function(res) {
        var data = JSON.parse(res.data);

        if (typeof onList[data.Type] == 'function') {
          var callback = onList[data.Type];
          delete data.Type;
          callback(data);
        } else {
          console.warn('Unhandled socket message:', data);
        }
      };
    }

    init();

    function sendRequest(request) {
      if (!ws) {
        console.error('WebSocket is not available');
        return;
      }

      ws.send(JSON.stringify(request));
    };

    function on(name, callback) {
      onList[name] = callback;
    }

    function emit(name, data) {
      sendRequest({
        body: data,
        Type: name
      });
    }

    var service = {
      sendRequest: sendRequest,
      on: on,
      emit: emit,
    };

    return service;
  }

})();
