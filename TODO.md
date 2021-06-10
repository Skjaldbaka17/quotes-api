# TODOS

- [ ] Clean up routes files
- [ ] Make function for building the SQL query (general for all?)
- [ ] Pagination Everywhere where needed
- [ ] Clean up get/set QOD/AOD

- [ ] Change all searches by a single id into a search of array of ids (i.e. in sql "in ids" rather than "= id")
- [ ] Review /authors for Swagger 
- [ ] Review /meta for Swagger 
- [ ] Review /quotes for Swagger 
- [ ] Review /search for Swagger 
- [ ] Review /topics for Swagger 
- [ ] Review Swagger 
- [ ] Clean up swagger in code
- [ ] Clean up Documentation look (Swagger)

- [ ] only return keys, in the response-json, that are relevant to the request


- [ ] Error handling + Tests
- [ ] Make error tests (i.e. made-to-fail-tests)
- [ ] Add error response to Swagger
- [ ] Go over Swagger + Clean it up and make pretty

- [ ] Add password protection / protected routes capability (at least for SetQuoteOfTheyDay route )
- [ ] Add authentication for access to the api + Creating apiKeys + Documenting usage + admin access vs normal access

- [ ] Add .env variables (i.e. for names of tables etc.)

- [ ] Insert Quote for created author or for a 'real' author (private and public)
- [ ] update inserted quote (priv and pub)
- [ ] Create new Author (private and public)
- [ ] Update created author (priv and pub)
- [ ] Create new Topic (private and public)
- [ ] update created topic (priv and pub)

- [ ] Setup AWS server

- [ ] Draw up DB-Graph (i.e. how tables are connected to view etc)

- [ ] is random truly random (i.e. does the "random" funcitonality truly return randomly or is it biased towards quotes in the "front" of the DB (i.e. in the front where postgres stores them))
- [ ] Make Authors Search more efficient (create a similarity-based index ?)
- [ ] Sort return list alphabetically Icelandic support

- [ ] Look into and maybe Change rest into GraphQL (neeeeeee, frekar fyrir næsta project)
- [ ] Look into payment for some privileges
- [ ] New crawler for new quotes / authors
 --------------

- [x] add search for topics and searching in topics (maybe just have a single search endpoint with parameters? i.e. want to search for authors / quotes / language inside a specific topic?)
- [x] getQuotes route (combine with getQuotesById and add to it to get quotes from a specific author + add pagination)
- [x] Add Counting each time a quote is accessed / sent from Api (also for topics) - i.e. stats
- [x] Add tests for GetAOD and AODHistory and SetAOD
- [x] Add tests for GetQOD and QODHistory and SetQOD
- [x] Add get and set Author of The Day (plus points for available to set authors for multiple days in one request + plus points for AOD history)
- [x] Add get and set Quote of The Day (plus points for available to set quotes for multiple days in one request + plus points for QOD history)
- [x] quoteoftheday << qod
- [x] Get list of authors (with parameters for pagination and alphabet and languages)
- [x] Get random Author
- [x] Validate RequestBody + Tests
- [x] Add get random Quote
- [x] Use TopicsView instead of searchview (and change name to something more general)
- [x] Add Icelandic / English Support
- [x] Add categories
- [x] English and Icelandic Authors with same name have same author id
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
