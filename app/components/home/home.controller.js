(function() {
  'use strict';

  angular
      .module('app')
      .controller('HomeController', HomeController);

  HomeController.$inject = ['WebsocketFactory', '$scope'];

  function HomeController(ws, $scope) {
    var vm = this;

    vm.username    = '';
    vm.usernameInput = '';
    vm.message     = '';

    vm.isLogged = false;
    vm.inUpdate = false;

    vm.messages = [];

    // Socket event
    ws.on('connected', function(data) {
      console.log('connected', data);
      $scope.$apply(function() {
        vm.username = data.Name;
        vm.isLogged = true;
      });
    });

    ws.on('ConnectedUsers', function(data) {
      console.log('users', data);
    });

    ws.on('message', function(data) {
      console.log('onmessage', data);
      $scope.$apply(function() {
        vm.messages.push(data);
      });
    });

    vm.sendMessage = sendMessage;
    vm.setUsername = setUsername;

    function sendMessage(e) {
      if (vm.username.length && vm.message.length) {
        console.log('here');
        ws.sendRequest({
          Name: vm.username,
          Body: vm.message,
          Type: 'message',
        });
      }
    }

    function setUsername(e) {
      if (e.keyCode == 13 && vm.usernameInput.length) {
        console.log('validUsername');
        ws.emit('connect', vm.usernameInput);
      }
    }
  }
})();
