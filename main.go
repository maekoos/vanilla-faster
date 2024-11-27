//go:generate sqlc generate
//go:generate pnpx tailwindcss -o static/styles.css -c tailwind.config.js -m
//go:generate templ generate .

package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"vanilla-faster/components"
	"vanilla-faster/repository"
	"vanilla-faster/resize"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed static
var staticFs embed.FS

func main() {
	// ctx := context.Background()
	// conn, err := pgx.Connect(context.Background(), "user=postgres dbname=postgres sslmode=disabled")
	urlExample := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	conn, err := sql.Open("pgx", urlExample)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	repo := repository.New(conn)

	http.Handle("GET /{$}", components.SidebarCollectionsMw(repo,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			collections, err := repo.GetCollectionsWithCategories(r.Context())
			if err != nil {
				log.Panicln(err)
			}
			count, err := repo.GetProductCount(r.Context())
			if err != nil {
				log.Panicln(err)
			}

			err = components.HomePage(collections, count).Render(r.Context(), w)
			if err != nil {
				log.Panicln(err)
			}
		})))

	http.Handle("GET /products/{category}", components.SidebarCollectionsMw(repo,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			categorySlug := r.PathValue("category")

			category, err := repo.GetCategory(r.Context(), categorySlug)
			if err != nil {
				// TODO: Not found
				log.Panicln(err)
			}

			productCount, err := repo.GetCategoryProductCount(r.Context(), categorySlug)
			if err != nil {
				log.Panicln(err)
			}

			err = components.CategoryPage(category, productCount).Render(r.Context(), w)
			if err != nil {
				log.Panicln(err)
			}
		})))

	http.Handle("GET /products/{category}/{subcategory}", components.SidebarCollectionsMw(repo,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// categorySlug := r.PathValue("category")
			subcategorySlug := r.PathValue("subcategory")

			products, err := repo.GetProductsForSubcategory(r.Context(), subcategorySlug)
			if err != nil {
				log.Panicln(err)
			}
			productCount, err := repo.GetSubcategoryProductCount(r.Context(), subcategorySlug)
			if err != nil {
				log.Panicln(err)
			}

			err = components.SubcategoryPage(r.URL.Path, products, productCount).Render(r.Context(), w)
			if err != nil {
				log.Panicln(err)
			}
		})))

	http.Handle("GET /products/{category}/{subcategory}/{product}", components.SidebarCollectionsMw(repo,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			categorySlug := r.PathValue("category")
			subcategorySlug := r.PathValue("subcategory")
			productSlug := r.PathValue("product")

			product, err := repo.GetProductDetails(r.Context(), productSlug)
			if err != nil {
				log.Panicln(err)
			}
			related, err := repo.GetProductsForSubcategory(r.Context(), subcategorySlug)
			if err != nil {
				log.Panicln(err)
			}

			for idx, p := range related {
				if p.Slug == productSlug {
					related = append(related[:idx], related[idx+1:]...)
					break
				}
			}

			err = components.ProductPage(fmt.Sprintf("/products/%s/%s", categorySlug, subcategorySlug), product, related).Render(r.Context(), w)
			if err != nil {
				log.Panicln(err)
			}
		})))

	http.Handle("GET /{collection}", components.SidebarCollectionsMw(repo,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			collectionSlug := r.PathValue("collection")

			details, err := repo.GetCollectionDetails(r.Context(), collectionSlug)
			if err != nil {
				// TODO: Not found
				log.Panicln(err)
			}

			if len(details) == 0 {
				http.Error(w, "Collection not found", http.StatusNotFound)
				return
			}

			collection := repository.Collection{}
			categories := make([]*repository.Category, 0)
			for _, c := range details {
				collection = c.Collection
				categories = append(categories, &c.Category)
			}

			err = components.CollectionPage(collection, categories).Render(r.Context(), w)
			if err != nil {
				log.Panicln(err)
			}
		})))

	http.Handle("GET /_image", cacheMw(resize.New(&resize.ResizeOptions{})))

	http.Handle("GET /static/", cacheMw(http.FileServerFS(staticFs)))

	fmt.Println("Running on port 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalln(err)
	}
}

func cacheMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "private, max-age=2592000")
		next.ServeHTTP(w, r)
	})
}
