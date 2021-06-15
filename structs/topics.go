package structs

type TopicDBModel struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Isicelandic string `json:"is_icelandic"`
}

type TopicAPIModel struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Isicelandic string `json:"is_icelandic"`
}

func (dbModel *TopicDBModel) ConvertToAPIModel() TopicAPIModel {
	return TopicAPIModel(*dbModel)
}

func (apiModel *TopicAPIModel) ConvertToDBModel() TopicDBModel {
	return TopicDBModel(*apiModel)
}

func ConvertToTopicsAPIModel(authors []TopicDBModel) []TopicAPIModel {
	authorsAPI := []TopicAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, TopicAPIModel(author))
	}
	return authorsAPI
}

func ConvertToTopicsDBModel(authors []TopicAPIModel) []TopicDBModel {
	authorsDB := []TopicDBModel{}
	for _, author := range authors {
		authorsDB = append(authorsDB, TopicDBModel(author))
	}
	return authorsDB
}

type TopicViewDBModel struct {
	// The author's id
	//Unique: true
	//example: 24952
	AuthorId int `json:"author_id"`
	// Name of author
	//example: Muhammad Ali
	Name string `json:"name"`
	// The quote's id
	//Unique: true
	//example: 582676
	QuoteId int `json:"quote_id" `
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote string `json:"quote"`
	// Whether or not this quote is in Icelandic or not
	// example: false
	IsIcelandic bool `json:"is_icelandic"`
	//swagger:ignore
	TopicName string `json:"topic_name"`
	//swagger:ignore
	TopicId int `json:"topic_id"`
}

type TopicViewAPIModel struct {
	// The author's id
	//Unique: true
	//example: 24952
	AuthorId int `json:"authorId"`
	// Name of author
	//example: Muhammad Ali
	Name string `json:"name"`
	// The quote's id
	//Unique: true
	//example: 582676
	QuoteId int `json:"quoteId" `
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote string `json:"quote"`
	// Whether or not this quote is in Icelandic or not
	// example: false
	IsIcelandic bool `json:"isIcelandic"`
	//swagger:ignore
	TopicName string `json:"topicName"`
	//swagger:ignore
	TopicId int `json:"topicId"`
}

func (dbModel *TopicViewDBModel) ConvertToAPIModel() TopicViewAPIModel {
	return TopicViewAPIModel(*dbModel)
}

func (apiModel *TopicViewAPIModel) ConvertToDBModel() TopicViewDBModel {
	return TopicViewDBModel(*apiModel)
}

func ConvertToTopicViewsAPIModel(views []TopicViewDBModel) []TopicViewAPIModel {
	viewsAPI := []TopicViewAPIModel{}
	for _, view := range views {
		viewsAPI = append(viewsAPI, TopicViewAPIModel(view))
	}
	return viewsAPI
}

func ConvertToTopicViewsDBModel(views []TopicViewAPIModel) []TopicViewDBModel {
	viewsDB := []TopicViewDBModel{}
	for _, view := range views {
		viewsDB = append(viewsDB, TopicViewDBModel(view))
	}
	return viewsDB
}
