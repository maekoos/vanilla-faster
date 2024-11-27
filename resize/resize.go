// Resize images from an external url + cache if possible

package resize

import (
	"bytes"
	"context"
	"fmt"
	"hash/fnv"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nfnt/resize"

	"github.com/coocood/freecache"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	freecache_store "github.com/eko/gocache/store/freecache/v4"
)

type Resize struct {
	Quality int
	cache   *cache.Cache[[]byte]
}

type ResizeOptions struct {
	Quality         int
	CacheSize       int
	CacheExpiration time.Duration
}

func New(opts *ResizeOptions) *Resize {
	if opts.Quality == 0 {
		opts.Quality = 50
	}
	if opts.CacheSize == 0 {
		opts.CacheSize = 200 * 1024 * 1024
	}
	if opts.CacheExpiration == 0 {
		opts.CacheExpiration = 10 * time.Hour
	}

	freecacheStore := freecache_store.NewFreecache(freecache.NewCache(opts.CacheSize), store.WithExpiration(opts.CacheExpiration))
	cacheManager := cache.New[[]byte](freecacheStore)

	return &Resize{Quality: opts.Quality, cache: cacheManager}
}

func (re *Resize) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imageUrl := r.FormValue("url")
	maxSize64, err := strconv.ParseUint(r.FormValue("size"), 10, 32)
	if err != nil {
		http.Error(w, "invalid size", http.StatusBadRequest)
		return
	}
	maxSize := uint(maxSize64)

	cacheKey := makeCacheKey(imageUrl, maxSize)

	// TODO: Smart caching with etag?
	// <https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag>
	// Done in mw instead
	// w.Header().Add("Cache-Control", "public, max-age=2592000")

	{
		value, err := re.cache.Get(r.Context(), cacheKey)
		if err == nil && len(value) != 0 {
			w.Write(value)
			return
		}
	}

	image, err := re.loadJpeg(r.Context(), imageUrl)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("failed to load image: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	resized := resize.Thumbnail(maxSize, maxSize, image, resize.NearestNeighbor)

	buf := new(bytes.Buffer)
	if err = jpeg.Encode(buf, resized, &jpeg.Options{Quality: re.Quality}); err != nil {
		// TODO: Handle this error better?
		log.Println(err)
		return
	}
	w.Write(buf.Bytes())

	// Cache the resized image:
	err = re.cache.Set(context.Background(), cacheKey, buf.Bytes())
	if err != nil {
		log.Panicln(err)
	}
}

func (re *Resize) loadJpeg(ctx context.Context, imageURL string) (image.Image, error) {
	loadedImage, err := re.cache.Get(ctx, makeCacheKey(imageURL, 0))
	if err == nil && len(loadedImage) != 0 {
		return jpeg.Decode(bytes.NewReader(loadedImage))
	}

	log.Printf("Fetching image: %s\n", imageURL)

	rq, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	err = re.cache.Set(context.Background(), makeCacheKey(imageURL, 0), buf.Bytes())
	if err != nil {
		// TODO: Better logging...
		log.Println(err)
	}

	return jpeg.Decode(buf)
}

func makeCacheKey(imageUrl string, maxSize uint) string {
	hash := fnv.New64a()
	hash.Write([]byte(imageUrl))

	return fmt.Sprintf("%X&%d", hash.Sum64(), maxSize)
}
