package components

import (
	"fmt"
	"net/url"
)

func ImageUrl(imageUrl string, maxSize int) string {
	return fmt.Sprintf("/_image?url=%s&size=%d", url.QueryEscape(imageUrl), maxSize)
}

type ImageAttrs struct {
	Count int
}

func NewImageAttrs() ImageAttrs {
	return ImageAttrs{
		Count: 0,
	}
}

func (ia *ImageAttrs) Next() {
	ia.Count += 1
}

func (ia *ImageAttrs) Loading() string {
	if ia.Count <= 8*3 {
		return "eager"
	}
	return "lazy"
}

func (ia *ImageAttrs) NextLoading() string {
	ia.Next()
	return ia.Loading()
}

func (ia *ImageAttrs) Decoding() string {
	if ia.Count <= 8*3 {
		return "sync"
	}
	return "async"
}

func (ia *ImageAttrs) NextDecoding() string {
	ia.Next()
	return ia.Decoding()
}
