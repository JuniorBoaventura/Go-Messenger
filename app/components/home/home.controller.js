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
    vm.me = {};
    vm.users    = [];

    // Socket event
    ws.on('connected', function(data) {
      console.log('connected', data);
      $scope.$apply(function() {
        vm.me = data;
        vm.isLogged = true;
      });
    });

    ws.on('ConnectedUsers', function(data) {
      console.log('users', data);

      if (!data.Users) { return; }

      $scope.$apply(function() {
        vm.users = data.Users;
      });
    });

    ws.on('disconnected', function(data) {
      $scope.$apply(function() {
        angular.forEach(vm.users, function(value, key) {
          if (value.id == data.Body) {
            delete vm.users[key];
          }
        });
      });
    });

    ws.on('newUser', function(data) {
      var user = {
        id: data.Body,
        username: data.Name,
      };

      $scope.$apply(function() {
        console.log();
        vm.users.push(user);
      });
    });

    ws.on('message', function(data) {
      $scope.$apply(function() {
        var explode = data.Name.split(':');
        data.id = explode[0];
        data.Name = explode[1];
        vm.messages.push(data);
      });
    });

    vm.sendMessage         = sendMessage;
    vm.sendMessageKeyboard = sendMessageKeyboard;
    vm.setUsername         = setUsername;

    function sendMessage() {
      if (vm.me.Name.length && vm.message.length) {
        ws.sendRequest({
          Name: vm.me.Body + ':' + vm.me.Name,
          Body: vm.message,
          Type: 'message',
        });

        vm.message = '';
      }
    }

    function sendMessageKeyboard(e) {
      if (e.keyCode == 13 && vm.message.length) {
        sendMessage();
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
