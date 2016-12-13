(function() {
  'use strict';

  angular
      .module('app')
      .controller('HomeController', HomeController);

  HomeController.$inject = ['WebsocketFactory'];

  function HomeController(ws) {
    function sendMessage() {

    }
  }

})();
