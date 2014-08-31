angular.module( 'mediabase', [
  'templates-app',
  'mediabase.home',
  'mediabase.movies',
  'mediabase.services',
  'ui.router',
  'ui.route'
])

.config( function myAppConfig ( $stateProvider, $urlRouterProvider ) {
  $urlRouterProvider.otherwise( '/home' );
})

.config(['$provide', function($provide) {
  $provide.decorator('$rootScope', ['$delegate', function($delegate) {
    $delegate.constructor.prototype.$onRootScope = function(name, listener) {
      var unsubscribe = $delegate.$on(name, listener);
      this.$on('$destroy', unsubscribe);
    };

    return $delegate;
  }]);
}])

.run( function run () {
})

.controller( 'AppCtrl', function AppCtrl ( $scope, $location, $rootScope ) {
  $scope.$on('$stateChangeSuccess', function(event, toState, toParams, fromState, fromParams){
    if ( angular.isDefined( toState.data.pageTitle ) ) {
      $scope.pageTitle = toState.data.pageTitle + ' | mediabase' ;
    }
  });

  $scope.search = function() {
    var term = $scope.queryTerm
    $rootScope.$emit('app.search', term)
  }
})

;