var classBuilder = angular.module('classBuilder', []);

classBuilder.controller('builder', function ($scope, $http, $q) {
    $scope.schedule = [];

    var degree = $.ajax({
        type: "GET",
        url: "/degree/",
        async: false
    }).responseJSON;
    var graph = $.ajax({
        type: "GET",
        url: "/all-classes/",
        async: false
    }).responseJSON;

    console.log("Graph: ", graph, "Degree: ", degree);

    function generateSchedule() {
        var classesTaken = [];
        var classesNeeded = [];

        // get bare bones for the classes a CS student needs to take
        for (var key in degree) {
            if (degree[key].RequirementType === "Class") {
                if (degree[key].RequirementQuantity == -1) {
                    // Get all the core classes because you need them all
                    classesNeeded = classesNeeded.concat(degree[key].Classes);
                } else {
                    var random = Math.floor(Math.random()*degree[key].Classes.length);
                    // randonly pick a science sequence that isnt geo!
                    classesNeeded = classesNeeded.concat(degree[key].Classes[random]);
                }
            } else {
                var credits = 0;
                var added = [];
                while (credits <= degree[key].RequirementQuantity) {
                    // Randomly fill in some electives!
                    var random = Math.floor(Math.random()*degree[key].Classes.length);
                    if (added.indexOf(random) < 0) {
                        classesNeeded.push(degree[key].Classes[random])
                        added.push(random);
                        credits += degree[key].Classes[random].Credits;
                    }
                }
            }
        }
        console.log("Classes needed before filling in prereqs: ", classesNeeded);

        for (var i = 0; i < classesNeeded.length; i++) {
            //console.log("ASD", classesNeeded);
            if (classesNeeded[i].Number != 0 && classesNeeded[i].Number != 999) {
                var name = classesNeeded[i].Department + " " + classesNeeded[i].Number;
                if (graph[name].Prereqs != null) {
                    for (var j = 0; j < graph[name].Prereqs.length; j++) {
                        classesNeeded = classesNeeded.concat(graph[graph[name].Prereqs[j]]);
                    }
                }
            }
        }

        var uniqueNames = [];
        $.each(classesNeeded, function(i, el){
            if($.inArray(el, uniqueNames) === -1) uniqueNames.push(el);
        });
        classesNeeded = uniqueNames;
        //console.log("classes needed: ", classesNeeded);

        $scope.schedules = generateClassGraph([], classesNeeded);
        topoGroupSort($scope.schedules, {"CSCI 141":true}, 14);
    };

    function generateClassGraph(classesTaken, classesToTake) {
        // list of lists
        topo = [];
        classesTaken = ["CSCI 141"];
        for (var i = 0; i < classesToTake.length; i++) {
            var name = classesToTake[i].Department + " " + classesToTake[i].Number;

            //console.log(classesTaken.indexOf(name)); // so sometimes the thing we statically define as being there is there.

            if (graph[name].Prereqs != null && classesTaken.indexOf(name) < 0) {
                for (var k = 0; k < graph[name].Prereqs.length; k++) {
                    topo.push([graph[name].Prereqs[k], name]);
                }
            }
        }
        console.log("to be topo sorted graph of classes to take :", topo);
        return tSort(topo);
    };


    function topoGroupSort(topoNodes, classesTaken, maxCreds) {

        if (topoNodes.length < 1) {
            return;
        }

        var currentQ = [];
        var curCreds = 0;
        var newTopoNodes = [];

        for (var i = 0; i < topoNodes.length; i++) {

            if (graph[topoNodes[i]].Prereqs != null) {

                if (!classesTaken[topoNodes[i]] && graph[topoNodes[i]].Prereqs.every(function(current, index, array) { return classesTaken[current]; })) {
                    console.log("taken this!");
                    if (curCreds + graph[topoNodes[i]].Credits <= maxCreds) {
                        currentQ.push(topoNodes[i]);
                        curCreds += graph[topoNodes[i]].Credits;
                    }

                } else if (!classesTaken[topoNodes[i]] &&
                           curCreds + graph[topoNodes[i]].Credits <= maxCreds) {
                    //graph[topoNodes[i]].Prereqs.every(function(current, index, array) { return classesTaken[current]; })) {

                    console.log("We have not taken and cant: ", topoNodes[i]);

                    for (var j = 0; j < graph[topoNodes[i]].Prereqs.length; j++) {
                        classesTaken[graph[topoNodes[i]].Prereqs[j]] = true;
                    }

                    currentQ.push(topoNodes[i]);
                    curCreds += graph[topoNodes[i]].Credits;


                } else if (graph[topoNodes[i]].Prereqs.every(function(current, index, array) { return classesTaken[current]; })) {

                    if (curCreds + graph[topoNodes[i]].Credits <= maxCreds) {
                        currentQ.push(topoNodes[i]);
                        curCreds += graph[topoNodes[i]].Credits;
                    }

                } else if (!classesTaken[topoNodes[i]] && graph[topoNodes[i]].Prereqs.every(function(current, index, array) { return classesTaken[current]; })) {
                    console.log("taken this!");
                    newTopoNodes.push(topoNodes[i]);

                } else if (!classesTaken[topoNodes[i]] || graph[topoNodes[i]].Prereqs.every(function(current, index, array) { return classesTaken[current]; })) {
                    console.log("taken this!");
                    newTopoNodes.push(topoNodes[i]);

                } else {
                    newTopoNodes.push(topoNodes[i]);
                }
            }
        }
        console.log("PUSHING SAD ASD ASD SAD :",  currentQ);
        $scope.schedule.push(currentQ);
        currentQ = [];
        curCred = 0;
        console.log(topoNodes, newTopoNodes);
        return topoGroupSort(newTopoNodes, classesTaken, maxCreds);
    };

    generateSchedule();
    console.log("ALL THE SCHEDULEEEE: ", $scope.schedules);


    //     } else if (graph[topoNodes[i]].Prereqs != null &&
    //                curCreds + graph[topoNodes[i]].Credits <= maxCreds) {

    //         var val = true;
    //         for (var k = 0; k < graph[topoNodes[i]].Prereqs.length; k++) {
    //             if (!classesTaken[graph[topoNodes[i]].Prereqs[k]]) {
    //                 val = false;
    //             }
    //         }
    //         if (val) {
    //             currentQ.push(topoNodes[i]);
    //             curCreds += graph[topoNodes[i]].Credits;
    //         }

    //     }
    //




    // var temp = [];
    // var curCreds = 0;
    // classesTaken = ["CSCI 141"];

    // for (var i = 0; i < topoNodes.length; i++) {
    //     var curClass = graph[topoNodes[i]];
    //     var curClassName = curClass.Department + " " + curClass.Number;
    //     for (var k = i; k < topoNodes.length; k++) {
    //         if (i != k) {
    //             var walkClass = graph[topoNodes[k]];
    //             var walkClassName = walkClass.Department + " " + walkClass.Number;
    //             //console.log(curClassName, walkClassName);
    //             if (walkClass.Prereqs == null && curCreds + walkClass.Credits <= maxCreds && classesTaken.indexOf(walkClassName) < 0) {
    //                 temp.push(walkClassName);
    //                 curCreds += walkClass.Credits;

    //             } else if (walkClass.Prereqs != null && curCreds + walkClass.Credits <= maxCreds && classesTaken.indexOf(walkClassName) < 0) {
    //                 if (walkClass.Prereqs.every(function(current, index, array) { return classesTaken.indexOf(current) > -1; })) {
    //                     temp.push(walkClassName);
    //                     curCreds += walkClass.Credits;

    //                 }
    //             }
    //         }
    //         $scope.schedule.push(temp);
    //         curCreds = 0;
    //         temp = [];
    //     }
    // }

});


function tSort(edges) {
    return toposort(uniqueNodes(edges), edges);
}

function toposort(nodes, edges) {
    var cursor = nodes.length
    , sorted = new Array(cursor)
    , visited = {}
    , i = cursor
    while (i--) {
        if (!visited[i]) visit(nodes[i], i, [])
    }
    return sorted
    function visit(node, i, predecessors) {
        if(predecessors.indexOf(node) >= 0) {
            throw new Error('Cyclic dependency: '+JSON.stringify(node))
        }
        if (visited[i]) return;
        visited[i] = true
        // outgoing edges
        var outgoing = edges.filter(function(edge){
            return edge[0] === node
        })
        if (i = outgoing.length) {
            var preds = predecessors.concat(node)
            do {
                var child = outgoing[--i][1]
                visit(child, nodes.indexOf(child), preds)
            } while (i)
        }
        sorted[--cursor] = node
    }
}
function uniqueNodes(arr){
    var res = []
    for (var i = 0, len = arr.length; i < len; i++) {
        var edge = arr[i]
        if (res.indexOf(edge[0]) < 0) res.push(edge[0])
        if (res.indexOf(edge[1]) < 0) res.push(edge[1])
    }
    return res
}
