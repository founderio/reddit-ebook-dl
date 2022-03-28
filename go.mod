module founderio.net/reddit-ebook-dl

go 1.18

require (
	github.com/bmaupin/go-epub v0.11.0
	github.com/go-test/deep v1.0.8
	github.com/gomarkdown/markdown v0.0.0-20220310201231-552c6011c0b8
	github.com/joho/godotenv v1.4.0
	github.com/vartanbeno/go-reddit/v2 v2.0.1
)

require (
	github.com/gabriel-vasile/mimetype v1.3.1 // indirect
	github.com/gofrs/uuid v3.1.0+incompatible // indirect
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50 // indirect
	golang.org/x/net v0.0.0-20210505024714-0287a6fb4125 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	google.golang.org/appengine v1.4.0 // indirect
)

// Replacement because of https://github.com/vartanbeno/go-reddit/pull/21
replace github.com/vartanbeno/go-reddit/v2 => github.com/sethjones/go-reddit/v2 v2.0.1-0.20220211043233-9af4e19ee575

replace github.com/bmaupin/go-epub => github.com/founderio/go-epub v0.11.1-0.20220327234436-bf11d823e362
