var classBrowser = angular.module('classBrowser', []);

classBrowser.controller('classList', function ($scope) {
    console.log("WE are here bb");
    $scope.classes = {
        "CSCI 141": {
            "credits": 4,
            "requires": [],
            "enables": ["CSCI 145"],
            "fulfills": ["QSR", "CSCI MAJOR CORE"],
            "title": "Intro to programming",
            "description": "This is an okay class because python."
        },
        "CSCI 145": {
            "credits": 4,
            "requires": ["CSCI 141"],
            "enables": ["CSCI 241", "CSCI 301", "CSCI 247"],
            "fulfills": ["QSR", "CSCI MAJOR CORE"],
            "title": "Linear data structures",
            "description": "This is a lame class because java."
        }
    };
});



[
    {"className":"CSCI 141",
     "properties":{},
    },
    {"className":"CSCI 145",
     "properties":{},
    }
]