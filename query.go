package main

import (
	//"encoding/json"
	"fmt"
	"github.com/jmcvetta/neoism"
	"os"
	"strconv"
)

type req struct {
	Classes             []Class
	RequirementType     string // Class, Credit, or Sequence
	RequirementQuantity int    // -1 for all otherwise a quantity
}

type adjSet struct {
	From string `json:"r.From"`
	To   string `json:"r.To"`
}

type Relation struct {
	Label string `json:"Label"`
}

type course struct {
	ClassName  string          //department + number
	Properties classProperties //see def below
}

type classProperties struct {
	Department string
	Number     int
	Credits    int      // number of credits for this class
	Enables    []string // classes this class enables
	Prereqs    []string // classes this class requires
	Title      string   // the name of this class
	Fulfills   []string // the requirements this class fulfills
}

type Class struct {
	Title      string `json:"Title"`
	Department string `json:"Department"`
	Number     int    `json:"Number"`
	Credits    int    `json:"Credits"`
}

var ClassMap map[string][]string
var Db *neoism.Database

func main1() {
	var err error

	Db, err = neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		panic(err)
	}

	var dept string
	if len(os.Args) < 3 {
		dept = "CSCI"
	} else {
		dept = os.Args[2]
	}
	class64, _ := strconv.ParseInt(os.Args[1], 10, 0)
	class := int(class64)

	createEnablesAndPrereqs()
	//fmt.Println(ClassMap[dept+" "+strconv.FormatInt(class64, 10)+"Enables"])
	prereqs(class, dept) //throw away to use class and dept so i can test.
	//notTaken := prereqsNotTaken(class, dept, "eric")
	//courseList := allCourseList()
	//printClasses(notTaken)
	//courseList = append(courseList, createCourse(141, "CSCI"))
	//courseList = append(courseList, createCourse(145, "CSCI"))
	//courseList = append(courseList, createCourse(241, "CSCI"))

	// testMap := make(map[string]Class)
	// testMap["test1"] = notTaken[0]
	// testMap["test2"] = notTaken[1]

	// jsonMap, _ := json.Marshal(testMap)
	// fmt.Println(string(jsonMap))

	degreeMap()

	//json, _ := json.Marshal(struct{ Data []Class }{notTaken})
	//fmt.Println("printing the json:")
	//fmt.Println(string(json))

	// fmt.Println("the prereq tree fo", class)
	// fmt.Println()
	// res := allPrereqs(class, dept)
	// printPrereqTree(res)
	// fmt.Println()

	// var res1 []Class
	// fmt.Println("prereqs of", class, "are:")
	// res1 = prereqs(class, dept)
	// printClasses(res1)
	// fmt.Println()

	// var res2 []Class
	// fmt.Println(class, "is a prereq of:")
	// res2 = enables(class, dept)
	// printClasses(res2)
}

func degreeMap() map[string]req {
	var cores []Class
	cq := neoism.CypherQuery{
		Statement: `match (c:Class)<-[r:FULFILLED_BY]-(core:Requirement{name: "CORE"})
                            return distinct c.title AS Title, c.department AS Department, c.number AS Number, c.credits AS Credits`,
		Result: &cores,
	}
	Db.Cypher(&cq)
	var electives []Class
	cq1 := neoism.CypherQuery{
		Statement: `match (c:Class)<-[r:FULFILLED_BY]-(core:Requirement{name: "ELECTIVES"})
                             return distinct c.title AS Title, c.department AS Department, c.number AS Number, c.credits AS Credits`,
		Result: &electives,
	}
	Db.Cypher(&cq1)

	// var physics []Class
	// cq2 := neoism.CypherQuery{
	// 	Statement: ` match (a:Class)<-[r]-(Requirement{name: "SCIENCE"}) where
	//                      a.department = "PHYS" match (a)<-[*]-(end:Class)
	//                      return distinct end.title AS Title, end.department AS Department, end.number AS Number end.credits AS Credits UNION match (a:Class)<-[r]-(Requirement{name: "SCIENCE"}) where
	//                      a.department = "PHYS" return a.title AS Title, a.department AS Department, a.number AS Number a.credits AS Credits`
	// 	Result: &physics
	// }
	// Db.Cypher(&cq2)

	// var chemistry []Class
	// cq3 := neoism.CypherQuery{
	// 	Statement: ` match (a:Class)<-[r]-(Requirement{name: "SCIENCE"}) where
	//                      a.department = "CHEM" match (a)<-[*]-(end:Class)
	//                      return distinct end.title AS Title, end.department AS Department, end.number AS Number end.credits AS Credits UNION match (a:Class)<-[r]-(Requirement{name: "SCIENCE"}) where
	//                      a.department = "CHEM" return a.title AS Title, a.department AS Department, a.number AS Number a.credits AS Credits`
	// 	Result: &chemistry
	// }
	// Db.Cypher(&cq3)

	// var biology []Class
	// cq4 := neoism.CypherQuery{
	// 	Statement: ` match (a:Class)<-[r]-(Requirement{name: "SCIENCE"}) where
	//                      a.department = "BIOL" match (a)<-[*]-(end:Class)
	//                      return distinct end.title AS Title, end.department AS Department, end.number AS Number end.credits AS Credits UNION match (a:Class)<-[r]-(Requirement{name: "SCIENCE"}) where
	//                      a.department = "BIOL" return a.title AS Title, a.department AS Department, a.number AS Number a.credits AS Credits`
	// 	Result: &biology
	// }
	// Db.Cypher(&cq4)
	var sci []Class
	cq5 := neoism.CypherQuery{
		Statement: `match (a:Class)<-[r]-(Requirement{name: "SCIENCE"})
                            return a.title AS Title, a.department AS Department, a.number AS Number, a.credits AS Credits `,
		Result: &sci,
	}
	Db.Cypher(&cq5)

	degree := make(map[string]req)

	var core req
	core.Classes = cores
	core.RequirementType = "Class"
	core.RequirementQuantity = -1
	degree["CORE"] = core

	var elective req
	elective.Classes = electives
	elective.RequirementType = "Credit"
	elective.RequirementQuantity = 16
	degree["ELECTIVES"] = elective

	var science req
	science.Classes = sci
	science.RequirementType = "Class"
	science.RequirementQuantity = 1
	degree["SCIENCE"] = science

	return degree
}

func prereqsNotTaken(classNumber int, department string, student string) []Class {
	var prereqs []Class
	cq := neoism.CypherQuery{
		Statement: `match (c:Class)<-[*]-(d:Class)
                            where c.department = {department} and c.number = {class}
                            return distinct d.title AS Title, d.department AS Department, d.number AS Number, d.credits AS Credits`,
		Parameters: neoism.Props{"class": classNumber, "department": department},
		Result:     &prereqs,
	}
	Db.Cypher(&cq)

	var taken []Class
	cq1 := neoism.CypherQuery{
		Statement: `match (c:Class)<-[*]-(d:Class) where c.department = {department}
                            and c.number = {class} match (d)<-[r1:TAKEN]-(s:Student) where s.name = {student}
                            return distinct d.title AS Title, d.department AS Department, d.number AS Number, d.credits AS Credits`,
		Parameters: neoism.Props{"class": classNumber, "department": department, "student": student},
		Result:     &taken,
	}
	Db.Cypher(&cq1)

	// fmt.Println("prereqs:")
	// printClasses(prereqs)
	// fmt.Println()
	// fmt.Println("taken:")
	// printClasses(taken)
	// fmt.Println()

	var notTaken []Class
	for i := range prereqs {
		prereq := prereqs[i]
		count := false
		for j := range taken {
			class := taken[j]
			if prereq.Title == class.Title {
				count = true
			}
		}
		if !count {
			notTaken = append(notTaken, prereq)
		}
	}
	return notTaken
}

//take the difference of the prereqs and the taken class lists

func allCourseList() []course {
	var courseList []course

	var classes []Class
	cq := neoism.CypherQuery{
		Statement: `
                    MATCH (d:Class) RETURN d.title AS Title, d.department AS Department, d.number AS Number, d.credits AS Credits`,
		Result: &classes,
	}
	err := Db.Cypher(&cq)
	if err != nil {
		fmt.Println(err)
	}

	//printClasses(classes)
	//fmt.Println(len(classes))
	for i := range classes {
		//fmt.Println(i)
		class := classes[i]
		//fmt.Println(class.Title, class.Number, class.Department, class.Credits)
		courseList = append(courseList, createCourse(class.Number, class.Department))
	}
	return courseList
}

func createEnablesAndPrereqs() {
	adj := getAdjSet()
	ClassMap = make(map[string][]string)
	for i := range adj {
		relation := adj[i]
		ClassMap[relation.From+"Enables"] = append(ClassMap[relation.From+"Enables"], relation.To)
		ClassMap[relation.To+"Prereqs"] = append(ClassMap[relation.To+"Prereqs"], relation.From)
	}
}
func getAdjSet() []adjSet {
	var adj []adjSet
	cq := neoism.CypherQuery{
		Statement: ` match (a:Class)-[r:PREREQ_OF]->(b:Class) return distinct r.From, r.To `,
		Result:    &adj,
	}
	Db.Cypher(&cq)
	return adj
}

func printClasses(res []Class) {
	for i := range res {
		r := res[i]
		fmt.Println("\t- ", r.Title, r.Department, r.Number, r.Credits)
	}
}

func printPrereqTree(res2 []Class) {
	for i := range res2 {
		r := res2[i]
		resi := allPrereqs(r.Number, r.Department)
		fmt.Println(r.Title, r.Department, r.Number, r.Credits)
		printClasses(resi)
		fmt.Println()
	}

}

func prereqs(classNumber int, department string) []Class {
	var res []Class
	cq := neoism.CypherQuery{
		Statement: `
                           MATCH (c:Class)<-[PREREQ_OF]-(d:Class)
                           WHERE c.number = {class} AND c.department = {department}
                           RETURN DISTINCT d.title AS Title, d.department AS Department, d.number AS Number, d.credits AS Credits `,
		Parameters: neoism.Props{"class": classNumber, "department": department},
		Result:     &res,
	}
	Db.Cypher(&cq)

	return res
}

func enables(classNumber int, department string) []Class {
	var res []Class
	cq := neoism.CypherQuery{
		Statement: `
                           MATCH (c:Class)-[PREREQ_OF]->(d:Class)
                           WHERE c.number = {class} AND c.department = {department}
                           RETURN DISTINCT d.title AS Title, d.department AS Department, d.number AS Number, d.credits AS Credits `,
		Parameters: neoism.Props{"class": classNumber, "department": department},
		Result:     &res,
	}
	Db.Cypher(&cq)

	return res
}

func degreeRequirements(degreeName string, degreeType string) {

}

// for gurs, so passing ACGM, BCGM, etc...
func requirementsOf(requirementType string) {

}

func allPrereqs(classNumber int, department string) []Class {
	var res []Class
	cq := neoism.CypherQuery{
		Statement: `
                           MATCH p = (a:Class)-[*]->(end)
                           WHERE end.number = {class} AND end.department = {department}
                           UNWIND EXTRACT(n IN nodes(p)| n ) AS d
                           WITH DISTINCT d WHERE d.number  <> {class}
                           RETURN d.title AS Title, d.department AS Department, d.number AS Number, d.credits AS Credits
                           ORDER BY d.department AS Department, d.number AS Number`,
		Parameters: neoism.Props{"class": classNumber, "department": department},
		Result:     &res,
	}
	Db.Cypher(&cq)

	return res
}

func createCourse(classNumber int, department string) course {

	var class []Class
	//fmt.Println(classNumber, department)
	cq := neoism.CypherQuery{
		Statement: `Match (d:Class) WHERE d.number  = {classNum}
                            AND d.department  = {classDept}
                            return distinct  d.title AS Title, d.department AS Department, d.number AS Number, d.credits AS Credits`,
		Parameters: neoism.Props{"classNum": classNumber, "classDept": department},
		Result:     &class,
	}
	//fmt.Println(len(class))
	//printClasses(class)
	Db.Cypher(&cq)

	//var enable []Class
	//var prereq []Class
	//enable = enables(classNumber, department)
	//prereq = prereqs(classNumber, department)

	//fmt.Println("prereqs of", classNumber, "are:")
	//printClasses(enable)
	//fmt.Println()

	var enableString []string
	var prereqString []string

	// for i := range enable {
	// 	class := enable[i]
	// 	//fmt.Println(class.Title)
	// 	enableString = append(enableString, class.Title)
	// }
	//fmt.Println("enableString:")
	// for i := range enableString {
	// 	fmt.Println(enableString[i])
	// }

	// for i := range prereq {
	// 	class := prereq[i]
	// 	prereqString = append(prereqString, class.Title)
	// }
	enableString = ClassMap[department+" "+strconv.FormatInt(int64(classNumber), 10)+"Enables"]
	prereqString = ClassMap[department+" "+strconv.FormatInt(int64(classNumber), 10)+"Prereqs"]
	// for i := range class {
	// 	class1 := class[i]
	// 	fmt.Println("class:", i, "title is:", class1.Title, "credits is:", class1.Credits)
	// }
	credits := class[0].Credits
	title := class[0].Title
	//fmt.Println("test", credits, title, class[0].Number, class[0].Department)
	//still need to do fulfills, need to put those relationships into the database.
	var fulfills []string
	//making the classProperties struct.

	//fmt.Println("department and classnumber", department, classNumber)
	cP := classProperties{department, classNumber, credits, enableString, prereqString, title, fulfills}
	//props, _ := json.Marshal(cP)
	//fmt.Println(string(props))

	className := fmt.Sprintf("%s %d", department, classNumber)
	course := course{className, cP}
	//courseJson, err := json.Marshal(course)
	//if err != nil {
	//	fmt.Println(err)
	//}
	return course
}
