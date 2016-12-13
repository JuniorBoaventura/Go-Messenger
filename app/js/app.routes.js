(function() {
  'use strict';

  angular
    .module('app')
    .config(config);

  config.$inject = ['$routeProvider'];

  function config($routeProvider) {
    let templateDir = 'template/';

    $routeProvider.when('/', {
      templateUrl: templateDir + 'home.html',
      controller: 'HomeController',
      controllerAs: 'vm',
    });
  }
})();
