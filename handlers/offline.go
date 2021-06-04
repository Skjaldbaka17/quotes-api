package handlers

//Functions below are for "offline" updating Database
//Points incremented for appearing in search list
const incrementAppearInSearchList = 1

//Points incremented for direct get
const incrementIdFetch = 10

func directFetchAuthorsCountIncrement(authorIds []int) error {
	if len(authorIds) == 0 {
		return nil
	}
	return db.Exec("UPDATE authors SET count = count + ? where id in (?) returning *", incrementIdFetch, authorIds).Error
}

func directFetchQuotesCountIncrement(quoteIds []int) error {
	if len(quoteIds) == 0 {
		return nil
	}
	return db.Exec("UPDATE quotes SET count = count + ? where id in (?) returning *", incrementIdFetch, quoteIds).Error
}

func directFetchTopicCountIncrement(topicId int, topicName string) error {
	return db.Exec("UPDATE topics SET count = count + ? where id = ? or lower(name) = lower(?) returning *", incrementIdFetch, topicId, topicName).Error
}

func authorIdsAppearInSearchCountIncrement(authorIds []int) error {
	if len(authorIds) == 0 {
		return nil
	}
	return db.Exec("UPDATE authors SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, authorIds).Error
}

func authorsAppearInSearchCountIncrement(authors []AuthorsView) error {
	if len(authors) == 0 {
		return nil
	}
	authorIds := []int{}

	for _, author := range authors {
		authorIds = append(authorIds, author.Id)
	}

	return db.Exec("UPDATE authors SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, authorIds).Error
}

func appearInSearchCountIncrement(quotes []QuoteView) error {
	if len(quotes) == 0 {
		return nil
	}
	authorIds := []int{}
	quoteIds := []int{}
	for _, quote := range quotes {
		authorIds = append(authorIds, quote.Authorid)
		quoteIds = append(quoteIds, quote.Quoteid)
	}

	err := db.Exec("UPDATE authors SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, authorIds).Error
	if err != nil {
		return err
	}
	err = db.Exec("UPDATE quotes SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, quoteIds).Error
	return err
}
