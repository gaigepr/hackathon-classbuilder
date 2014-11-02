package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmcvetta/neoism"
	"io/ioutil"
	"net/http"
	"strconv"
)

func loadFile(filepath string) (string, error) {
	body, error := ioutil.ReadFile(filepath)
	bodyString := string(body[:])
	if error != nil {
		fmt.Println(error)
		return "", error
	}
	//fmt.Println(bodyString)
	return bodyString, nil
}

func allClassJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println("request for all classes")
	courseList := allCourseList()
	allClasses := make(map[string]classProperties)
	//fmt.Println(len(courseList))
	for i := range courseList {
		allClasses[courseList[i].ClassName] = courseList[i].Properties
		//fmt.Println(courseList[i].ClassName)
	}
	json, err := json.Marshal(allClasses)
	fmt.Println(string(json))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	fmt.Fprintf(w, "%v", string(json))
}

func classNotTakenJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	number, err := strconv.ParseInt(vars["num"], 10, 0)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	notTaken := prereqsNotTaken(int(number), vars["dept"], vars["name"])
	json, err := json.Marshal(struct{ Data []Class }{notTaken})
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	fmt.Fprintf(w, "%v", string(json))
}

func degreeJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	degree := degreeMap()
	json, err := json.Marshal(degree)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	fmt.Fprintf(w, "%v", string(json))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	index, err := loadFile("static/index.html")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	fmt.Fprintf(w, "%s", index)
}

func addClassTaken(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var res []Relation
		cq1 := neoism.CypherQuery{
			Statement: `match (c1:Class), (s:Student) where c1.number = {class} and
                                    c1.department = {department} and s.name = {student} create
                                    unique (c1)<-[r:TAKE]-(s) set r.Label = "taken" `,
			Parameters: neoism.Props{"class": 301, "department": "CSCI", "student": "eric"},
			Result:     &res,
		}
		Db.Cypher(&cq1)
	}
}

func addClassWantToTake(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var res []Relation
		cq := neoism.CypherQuery{
			Statement: `match (c1:Class), (s:Student) where c1.number = {class} and
                                    c1.department = {department} and s.name = {student} create
                                    unique (c1)<-[r:TAKE]-(s) return r`,
			Parameters: neoism.Props{"class": 301, "department": "CSCI", "student": "eric"},
			Result:     &res,
		}
		Db.Cypher(&cq)
		if res[0].Label == "" {
			cq1 := neoism.CypherQuery{
				Statement: `match (c1:Class), (s:Student) where c1.number = {class} and
                                    c1.department = {department} and s.name = {student} create
                                    unique (c1)<-[r:TAKE]-(s) set r.Label = "wants to take" `,
				Parameters: neoism.Props{"class": 301, "department": "CSCI", "student": "eric"},
				Result:     &res,
			}
			Db.Cypher(&cq1)
		}

	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		cq := neoism.CypherQuery{
			Statement:  `create unique (s:Student) set s.name = {name}, s.email = {email}`,
			Parameters: neoism.Props{"name": "eric", "email": "clarke21@students.wwu.edu"},
		}
		Db.Cypher(&cq)
	}
}

func main() {
	var err error
	Db, err = neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		panic(err)
	}
	createEnablesAndPrereqs()

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/", index)
	router.HandleFunc("/degree/", degreeJson)
	router.HandleFunc("/all-classes/", allClassJson)
	router.HandleFunc("/can-take/{name}/{dept}/{num}/", classNotTakenJson)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/", router)

	http.ListenAndServe(":8080", nil)

}
