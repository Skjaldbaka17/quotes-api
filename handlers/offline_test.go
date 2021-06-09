package handlers

import (
	"testing"

	"github.com/Skjaldbaka17/quotes-api/structs"
)

func TestIncrementCount(t *testing.T) {
	t.Run("Should Increment Authors count from direct fetch by ids", func(t *testing.T) {
		authorIds := Set{1, 2, 3}
		err := DirectFetchAuthorsCountIncrement(authorIds)

		if err != nil {
			t.Fatalf("Expected no error but got %s", err.Error())
		}
		var authors []structs.AuthorsView
		err = Db.Table("authors").Where("id in ?", authorIds).Find(&authors).Error

		if err != nil {
			t.Fatalf("Expected no error but got %s", err.Error())
		}

		if authors[0].Count == 0 || authors[1].Count == 0 || authors[2].Count == 0 {
			t.Fatalf("Expected count of authors given to increase to above 0 but got count 0 for author: %+v", authors[0])
		}
	})

	t.Run("Should Increment Authors count from appearing in a search", func(t *testing.T) {

		quotes := []structs.QuoteView{
			{
				Authorid: 100,
				Quoteid:  100,
			},
		}
		err := AppearInSearchCountIncrement(quotes)

		if err != nil {
			t.Fatalf("Expected no error but got %s", err.Error())
		}

		var authors []structs.AuthorsView
		err = Db.Table("authors").Where("id = ?", quotes[0].Authorid).Find(&authors).Error
		if err != nil {
			t.Fatalf("Expected no error but got %s", err.Error())
		}
		if authors[0].Count == 0 {
			t.Fatalf("Expected count of authors given to increase to above 0 but got count 0 for author: %+v", authors[0])
		}

	})
	t.Run("Should Increment Quotes count from direct fetch by ids", func(t *testing.T) { t.Skip() })
	t.Run("Should Increment Quotes count from appearing in a search", func(t *testing.T) { t.Skip() })
}

type Set []int
