package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joeychilson/simplemux"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Repository is the json representation of a GitHub repository.
type Repository struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	URL            string    `json:"url"`
	ForkCount      int       `json:"forkCount"`
	StargazerCount int       `json:"stargazerCount"`
	Language       string    `json:"language"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

func main() {
	ctx := context.Background()

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	))
	client := githubv4.NewClient(httpClient)

	mux := simplemux.New()

	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var query struct {
				RateLimit struct {
					Limit     int
					Remaining int
					ResetAt   githubv4.DateTime
				}
			}

			err := client.Query(ctx, &query, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if query.RateLimit.Remaining < 1 {
				w.Header().Set("Retry-After", query.RateLimit.ResetAt.Time.Format(time.RFC1123))
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(query.RateLimit.Limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(query.RateLimit.Remaining))
			w.Header().Set("X-RateLimit-Reset", query.RateLimit.ResetAt.Time.Format(time.RFC1123))

			next.ServeHTTP(w, r)
		})
	})

	mux.Get("/user/:username", func(w http.ResponseWriter, r *http.Request) {
		username := simplemux.Param(r, "username")
		if username == "" {
			http.Error(w, "username is required", http.StatusBadRequest)
			return
		}

		var query struct {
			User struct {
				Login       string `graphql:"login"`
				PinnedItems struct {
					Edges []struct {
						Node struct {
							Repository struct {
								Name            string
								Description     string
								URL             string
								ForkCount       int
								StargazerCount  int
								PrimaryLanguage struct {
									Name string
								}
								UpdatedAt time.Time
								CreatedAt time.Time
							} `graphql:"... on Repository"`
						} `graphql:"node"`
					}
				} `graphql:"pinnedItems(first: 6)"`
			} `graphql:"user(login: $login)"`
		}

		variables := map[string]interface{}{
			"login": githubv4.String(username),
		}

		if err := client.Query(ctx, &query, variables); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var repos []Repository
		for _, edge := range query.User.PinnedItems.Edges {
			repos = append(repos, Repository{
				Name:           edge.Node.Repository.Name,
				Description:    edge.Node.Repository.Description,
				URL:            edge.Node.Repository.URL,
				ForkCount:      edge.Node.Repository.ForkCount,
				StargazerCount: edge.Node.Repository.StargazerCount,
				Language:       edge.Node.Repository.PrimaryLanguage.Name,
				UpdatedAt:      edge.Node.Repository.UpdatedAt,
				CreatedAt:      edge.Node.Repository.CreatedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repos)
	})

	mux.Get("/org/:orgname", func(w http.ResponseWriter, r *http.Request) {
		orgname := simplemux.Param(r, "orgname")
		if orgname == "" {
			http.Error(w, "orgname is required", http.StatusBadRequest)
			return
		}

		var query struct {
			Organization struct {
				Login       string `graphql:"login"`
				PinnedItems struct {
					Edges []struct {
						Node struct {
							Repository struct {
								Name            string
								Description     string
								URL             string
								ForkCount       int
								StargazerCount  int
								PrimaryLanguage struct {
									Name string
								}
								UpdatedAt time.Time
								CreatedAt time.Time
							} `graphql:"... on Repository"`
						} `graphql:"node"`
					}
				} `graphql:"pinnedItems(first: 6)"`
			} `graphql:"organization(login: $login)"`
		}

		variables := map[string]interface{}{
			"login": githubv4.String(orgname),
		}

		if err := client.Query(ctx, &query, variables); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var repos []Repository
		for _, edge := range query.Organization.PinnedItems.Edges {
			repos = append(repos, Repository{
				Name:           edge.Node.Repository.Name,
				Description:    edge.Node.Repository.Description,
				URL:            edge.Node.Repository.URL,
				ForkCount:      edge.Node.Repository.ForkCount,
				StargazerCount: edge.Node.Repository.StargazerCount,
				Language:       edge.Node.Repository.PrimaryLanguage.Name,
				UpdatedAt:      edge.Node.Repository.UpdatedAt,
				CreatedAt:      edge.Node.Repository.CreatedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repos)
	})

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
