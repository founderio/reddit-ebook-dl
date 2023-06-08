package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"founderio.net/reddit-ebook-dl/redl"
	"github.com/bmaupin/go-epub"
	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

const (
	epubCSSFile = "epub.css"
)

func main() {

	postID := flag.String("p", "", "Post ID")
	flag.Parse()

	if len(*postID) == 0 {
		log.Println("Missing Post ID")
		flag.Usage()
		os.Exit(1)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}

	appID, ok := os.LookupEnv("REDDIT_API_ID")
	if !ok {
		log.Fatalln("App ID not specified (ensure REDDIT_API_ID env var is set or present in .env file)")
	}
	appSecret, ok := os.LookupEnv("REDDIT_API_SECRET")
	if !ok {
		log.Fatalln("App Secret not specified (ensure REDDIT_API_SECRET env var is set or present in .env file)")
	}

	userAgent := "desktop:net.founderio.redl:0.1 (by /u/founderio)"

	credentials := reddit.Credentials{ID: appID, Secret: appSecret, Username: "", Password: ""}
	client, _ := reddit.NewClient(credentials, reddit.WithUserAgent(userAgent), reddit.WithApplicationOnlyOAuth(true))

	post, _, err := client.Post.Get(context.Background(), *postID)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tags, cleanedTitle := redl.ExtractTags(post.Post.Title)

	log.Printf("Downloading ebook with title \"%s\"", cleanedTitle)

	// TODO: extract flairs as tags, needs this change:
	// https://github.com/vartanbeno/go-reddit/pull/19/files

	var filterUser string
	authorDisplayName := post.Post.Author
	author2DisplayName := ""

	var containsPITag bool
	for _, tag := range tags {
		if tag == "PI" {
			containsPITag = true
			break
		}
	}

	if strings.Contains(post.Post.SubredditName, "Prompt") && !containsPITag {
		log.Println("Detected writing prompt, using longest-first-comment logic")

		tags = append(tags, "Writing Prompt")

		// Identify starting comment + matching user (add as an author)
		longestPostSoFar := 0
		filterUser = ""
		var filterUserName = ""
		for _, comment := range post.Comments {
			length := len(comment.Body)
			if length > longestPostSoFar && comment.Author != "AutoModerator" {
				longestPostSoFar = length
				filterUser = comment.AuthorID
				filterUserName = comment.Author
			}
		}
		if len(filterUser) == 0 {
			log.Fatalln("Could not identify any comments")
		}
		author2DisplayName = filterUserName
	} else {
		log.Println("Detected default self-post, using same-user longest-first-comment logic")
		filterUser = post.Post.AuthorID
	}

	chapters := make([]string, 0)

	if len(post.Post.Body) > 0 {
		chapters = append(chapters, "<h1>"+post.Post.Title+"</h1>\n"+redl.FormatPost(post.Post.Body))
	} else {
		chapters = append(chapters, "<h1>"+post.Post.Title+"</h1>")
	}

	currentComment := getLongestComment(post.Comments, filterUser)
	if currentComment == nil {
		log.Fatalln("Could not identify first comment")
	}
	for currentComment != nil {
		chapters = append(chapters, redl.FormatPost(currentComment.Body))

		// For subsequent comments load ALL comments, to ensure we get all continuations should they be there
		for currentComment.HasMore() {
			_, err := client.Comment.LoadMoreReplies(context.Background(), currentComment)
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
		currentComment = getLongestComment(currentComment.Replies.Comments, filterUser)
	}

	css := `
body {
	font-family: sans-serif;
}
`

	ebook := epub.NewEpub(cleanedTitle)
	epubCSSReferencePath, err := ebook.AddCSS("data:text/css;base64,"+base64.StdEncoding.EncodeToString([]byte(css)), epubCSSFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	ebook.Pkg.AddCreator(authorDisplayName, epub.PropertyRoleAuthor)
	if len(author2DisplayName) > 0 {
		ebook.Pkg.AddCreator(author2DisplayName, epub.PropertyRoleAuthor)
	}
	ebook.Pkg.AddIdentifier("url:"+post.Post.URL, epub.SchemeXSDString, "URL")
	ebook.Pkg.SetLang("en")
	if len(post.Post.Body) > 0 {
		ebook.Pkg.SetDescription("<p>" + cleanedTitle + "</p>\n<p>" + post.Post.Body + "</p>")
	} else {
		ebook.Pkg.SetDescription("<p>" + cleanedTitle + "</p>")
	}
	ebook.Pkg.SetPublisher("reddit.com/" + post.Post.SubredditNamePrefixed)
	ebook.Pkg.SetDate(post.Post.Created.Time)
	ebook.Pkg.SetSubject(tags)

	for i, chapter := range chapters {
		chapterTitle := fmt.Sprintf("file_%d", i)
		_, err := ebook.AddSection(chapter, chapterTitle, chapterTitle+".xhtml", epubCSSReferencePath)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	err = ebook.Write("out.epub")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getLongestComment(comments []*reddit.Comment, username string) *reddit.Comment {
	longestPostSoFar := 0
	var longestPost *reddit.Comment = nil
	for _, comment := range comments {
		length := len(comment.Body)
		if length > longestPostSoFar {
			longestPostSoFar = length
			longestPost = comment
		}
	}
	return longestPost
}
