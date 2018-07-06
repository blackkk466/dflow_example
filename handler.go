package main

import "fmt"




func getAllNotes() []Note{
	var noteList []Note

	db.Find(&noteList)

	return noteList
}

func getProducts() string {
	var pList []Product

	db.Find(&pList)

	str := ""

	for _, p := range pList {
		str += p.Name
		str += fmt.Sprintf("-%d-%.2f | ", p.Amount,p.Price)
	}

	return str
}
