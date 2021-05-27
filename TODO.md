# TODOS

- [ ] Validate RequestBody + Tests
- [ ] add search for topics and searching in topics (maybe just have a single search endpoint with parameters? i.e. want to search for authors / quotes / language inside a specific topic?)
- [ ] Clean up quotes.go
- [ ] Error handling + Tests
- [ ] Make error tests (i.e. made-to-fail-tests)
- [ ] English and Icelandic Authors with same name have same author id
- [ ] Add get random Quotes / Quote of The Day
- [ ] Add icelandic / English Support
- [ ] Make function for building the SQL query (general for all?)
- [ ] Add categories
- [ ] Add categories-search (Motivational etc)
- [ ] Setup AWS server
- [ ] Make Authors Search more efficient (create a similarity-based index ?)
- [ ] Look into and maybe Change rest into GraphQL
- [x] Add Search-"scroll", User is searching and is scrolling through her search and wants next batch of results matching her search i.e. PAGINATION
- [x] setup testing (unit)
- [x] Implement GetQuotesById (multiple quotes route)
- [x] Clean-up test files (Move some lines into their own functions etc.)
- [x] Setup Swagger for api-docs 
      * https://github.com/go-swagger/go-swagger
      * https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
      * https://www.youtube.com/watch?v=07XhTqE-j8k&t=374s
      * https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/handlers/docs.go


Author api: https://quotes.rest




### Helpful resources for full-text search in postgres

* https://www.opsdash.com/blog/postgres-full-text-search-golang.html 
* https://medium.com/@bencagri/implementing-multi-table-full-text-search-with-gorm-632518257d15
* https://www.freecodecamp.org/news/fuzzy-string-matching-with-postgresql/
* https://www.compose.com/articles/mastering-postgresql-tools-full-text-search-and-phrase-search/ 
