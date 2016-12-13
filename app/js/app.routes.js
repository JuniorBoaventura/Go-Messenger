(function() {
  'use strict';
  let templateDir = 'template/';
  angular
    .module('app')
    .config(function($routeProvider) {
      $routeProvider
      .when('/', {
        templateUrl: templateDir + 'home.html',
        controller: 'HomeController'
      });
    });
})();
