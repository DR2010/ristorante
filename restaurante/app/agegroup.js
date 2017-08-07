// ---------------------
//     agegroup.js
// ---------------------
var agegroup = angular.module('myApp', []);

agegroup.controller('agegroupController', function($scope, $http) 
	{
		getAgeGroup(); // Load all available age groups 
		function getAgeGroup()
		{  
			$http.get("ajax/getAgeGroup.php").success(function(data)
			{
				$scope.agegroups = data;
			});
		};

		getTeamList(); // Load all available age groups 
		function getTeamList()
		{  
			$http.get("ajax/getTeamList.php").success(function(data)
			{
				$scope.teamlist = data;
			});
		};
		
		
		$scope.getTeamListX = function (agegroupid) {
			$http.get("ajax/getTeamList.php?agegroupid="+agegroupid).success(function(data){
				getTask();
				$scope.teamlist = data;
			});
		};
	});