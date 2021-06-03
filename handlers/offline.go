package handlers

//Points incremented for appearing in search list
const pointsAppearInSearchList = 1

//Points incremented for direct get
const pointsIdFetch = 10

func searchPoints(authorIds []int) {
	db.Exec("UPDATE users SET name = ?", "jinzhu")
}

func searchPoints(authorIds []string) {

}
