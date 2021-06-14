package structs

type AuthorDBModel struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	HasIcelandicQuotes  bool   `json:"has_icelandic_quotes"`
	NrOfIcelandicQuotes int    `json:"nr_of_icelandic_quotes"`
	NrOfEnglishQuotes   int    `json:"nr_of_english_quotes"`
	Count               int    `json:"count"`
}

type AuthorAPIModel struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	HasIcelandicQuotes  bool   `json:"hasIcelandicQuotes"`
	NrOfIcelandicQuotes int    `json:"nrOfIcelandicQuotes"`
	NrOfEnglishQuotes   int    `json:"nrOfEnglishQuotes"`
	Count               int    `json:"count"`
}

func (dbModel *AuthorDBModel) ConvertToAPIModel() AuthorAPIModel {
	return AuthorAPIModel(*dbModel)
}

func (apiModel *AuthorAPIModel) ConvertToDBModel() AuthorDBModel {
	return AuthorDBModel(*apiModel)
}

func ConvertToAuthorsAPIModel(authors []AuthorDBModel) []AuthorAPIModel {
	authorsAPI := []AuthorAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, AuthorAPIModel(author))
	}
	return authorsAPI
}

func ConvertToAuthorsDBModel(authors []AuthorAPIModel) []AuthorDBModel {
	authorsDB := []AuthorDBModel{}
	for _, author := range authors {
		authorsDB = append(authorsDB, AuthorDBModel(author))
	}
	return authorsDB
}

type AodDBModel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

type AuthorsView struct {
	// The author's id
	//Unique: true
	//example: 24952
	Id int `json:"id"`
	// Name of the author
	//example: Muhammad Ali
	Name string `json:"name"`
	// Whether or not this author has some icelandic quotes
	//example: true
	Hasicelandicquotes bool `json:"has_icelandic_quotes"`
	// How many quotes in icelandic this author has
	//example: 6
	Nroficelandicquotes int `json:"nr_of_icelandic_quotes"`
	// How many quotes in icelandic this author has
	//example: 78
	Nrofenglishquotes int `json:"nr_of_english_quotes"`
	//swagger:ignore
	Count int `json:"count"`
	//swagger:ignore
	Date string `json:"date"`
}