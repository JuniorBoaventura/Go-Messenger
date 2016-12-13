(function() {
  'use strict';

  angular
      .module('app')
      .controller('HomeController', HomeController);

  HomeController.$inject = ['WebsocketFactory'];

  function HomeController(ws) {
    var vm = this;

    vm.username    = 'John';
    vm.message     = 'Hello Michel';

    vm.messages = [
      {
        Name: 'Antoine D',
        Body: 'Hello World',
      },
      {
        Name: 'Antoine D',
        Body: 'Hello World',
      },
    ];

    vm.sendMessage = sendMessage;

    function sendMessage() {
      if (vm.username.length && vm.message.length) {
        ws.sendRequest({
          Name: vm.username,
          Body: vm.message,
          Type: 'message',
        });
      }
    }
  }
})();
