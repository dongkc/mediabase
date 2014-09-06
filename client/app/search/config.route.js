(function () {
    'use strict';

    angular
        .module('app.search')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('search', {
                    url: '/search',
                    templateUrl: 'app/layout/main.html',
                    controller: 'Search as vm',
                })            
        });

})();