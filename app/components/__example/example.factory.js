(function() {
  'use strict';

  angular
      .module('app')
      .factory('ExampleFactory', ExampleFactory);

  function ExampleFactory() {
    var service = {
      test: test
    };
    return service;

    function test() {

    }
  }

})();
