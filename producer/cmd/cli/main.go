package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const defaultBaseURL = "http://localhost:8080"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	baseURL := strings.TrimRight(getEnv("API_BASE_URL", defaultBaseURL), "/")

	switch os.Args[1] {
	case "post":
		cmdPost(baseURL, os.Args[2:])
	case "comment":
		cmdComment(baseURL, os.Args[2:])
	case "like":
		cmdLike(baseURL, os.Args[2:])
	case "view":
		cmdView(baseURL, os.Args[2:])
	case "posts":
		cmdPosts(baseURL, os.Args[2:])
	case "comments":
		cmdComments(baseURL, os.Args[2:])
	case "search":
		cmdSearch(baseURL, os.Args[2:])
	case "report-comments":
		cmdReportComments(baseURL, os.Args[2:])
	case "report-days":
		cmdReportDays(baseURL)
	case "report-views":
		cmdReportViews(baseURL, os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func cmdPost(baseURL string, args []string) {
	fs := flag.NewFlagSet("post", flag.ExitOnError)
	title := fs.String("title", "", "post title")
	body := fs.String("body", "", "post body")
	author := fs.String("author", "", "post author")
	_ = fs.Parse(args)

	require(*title != "", "-title is required")
	require(*body != "", "-body is required")
	require(*author != "", "-author is required")

	mustPostJSON(baseURL+"/api/v1/posts", map[string]any{
		"title":  *title,
		"body":   *body,
		"author": *author,
	})
}

func cmdComment(baseURL string, args []string) {
	fs := flag.NewFlagSet("comment", flag.ExitOnError)
	postID := fs.Int("post-id", 0, "post id")
	body := fs.String("body", "", "comment body")
	author := fs.String("author", "", "comment author")
	_ = fs.Parse(args)

	require(*postID > 0, "-post-id is required and must be > 0")
	require(*body != "", "-body is required")
	require(*author != "", "-author is required")

	mustPostJSON(baseURL+"/api/v1/comments", map[string]any{
		"post_id": *postID,
		"body":    *body,
		"author":  *author,
	})
}

func cmdLike(baseURL string, args []string) {
	fs := flag.NewFlagSet("like", flag.ExitOnError)
	postID := fs.Int("post-id", 0, "post id")
	user := fs.String("user", "", "user name")
	_ = fs.Parse(args)

	require(*postID > 0, "-post-id is required and must be > 0")
	require(*user != "", "-user is required")

	mustPostJSON(baseURL+"/api/v1/likes", map[string]any{
		"post_id": *postID,
		"user":    *user,
	})
}

func cmdView(baseURL string, args []string) {
	fs := flag.NewFlagSet("view", flag.ExitOnError)
	postID := fs.Int("post-id", 0, "post id")
	user := fs.String("user", "", "user name")
	_ = fs.Parse(args)

	require(*postID > 0, "-post-id is required and must be > 0")

	mustPostJSON(baseURL+"/api/v1/views", map[string]any{
		"post_id": *postID,
		"user":    *user,
	})
}

func cmdPosts(baseURL string, args []string) {
	fs := flag.NewFlagSet("posts", flag.ExitOnError)
	limit := fs.Int("limit", 20, "posts limit")
	_ = fs.Parse(args)

	mustGet(baseURL + "/api/v1/posts?limit=" + strconv.Itoa(*limit))
}

func cmdComments(baseURL string, args []string) {
	fs := flag.NewFlagSet("comments", flag.ExitOnError)
	postID := fs.Int("post-id", 0, "post id")
	limit := fs.Int("limit", 50, "comments limit")
	_ = fs.Parse(args)

	require(*postID > 0, "-post-id is required and must be > 0")
	mustGet(baseURL + "/api/v1/posts/" + strconv.Itoa(*postID) + "/comments?limit=" + strconv.Itoa(*limit))
}

func cmdSearch(baseURL string, args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	query := fs.String("query", "", "search query")
	_ = fs.Parse(args)

	require(*query != "", "-query is required")
	mustGet(baseURL + "/api/v1/search?query=" + *query)
}

func cmdReportComments(baseURL string, args []string) {
	fs := flag.NewFlagSet("report-comments", flag.ExitOnError)
	limit := fs.Int("limit", 10, "report limit")
	_ = fs.Parse(args)
	mustGet(baseURL + "/api/v1/reports/top-posts-by-comments?limit=" + strconv.Itoa(*limit))
}

func cmdReportDays(baseURL string) {
	mustGet(baseURL + "/api/v1/reports/posts-comments-by-day")
}

func cmdReportViews(baseURL string, args []string) {
	fs := flag.NewFlagSet("report-views", flag.ExitOnError)
	limit := fs.Int("limit", 10, "report limit")
	_ = fs.Parse(args)
	mustGet(baseURL + "/api/v1/reports/top-posts-by-views?limit=" + strconv.Itoa(*limit))
}

func mustPostJSON(url string, payload any) {
	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		fatal(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(bytesPayload))
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()
	printResponse(resp)
}

func mustGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()
	printResponse(resp)
}

func printResponse(resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fatal(err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		fmt.Printf("status: %d\n", resp.StatusCode)
		fmt.Printf("body: %s\n", string(body))
		os.Exit(1)
	}

	var pretty bytes.Buffer
	if json.Indent(&pretty, body, "", "  ") == nil {
		fmt.Println(pretty.String())
		return
	}
	fmt.Println(string(body))
}

func require(ok bool, msg string) {
	if !ok {
		fmt.Println("error:", msg)
		os.Exit(1)
	}
}

func fatal(err error) {
	fmt.Println("error:", err)
	os.Exit(1)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func printUsage() {
	fmt.Println("Kafka lab CLI")
	fmt.Println("Usage:")
	fmt.Println("  go run ./producer/cmd/cli <command> [flags]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  post             -title -body -author")
	fmt.Println("  comment          -post-id -body -author")
	fmt.Println("  like             -post-id -user")
	fmt.Println("  view             -post-id [-user]")
	fmt.Println("  posts            [-limit]")
	fmt.Println("  comments         -post-id [-limit]")
	fmt.Println("  search           -query")
	fmt.Println("  report-comments  [-limit]")
	fmt.Println("  report-days")
	fmt.Println("  report-views     [-limit]")
	fmt.Println("")
	fmt.Println("Env:")
	fmt.Println("  API_BASE_URL (default: http://localhost:8080)")
}
