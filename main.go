package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/appengine/blobstore"
	aefile "google.golang.org/appengine/file"
	"google.golang.org/appengine/image"
	"google.golang.org/appengine"
	"encoding/json"
)

type payload struct {
	Url string `json:"url"`
}

type errorPayload struct {
	Message string `json:"message"`
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var p payload
	var ep errorPayload
	var err error
	path := r.URL.Query().Get("path")
	ctx := appengine.NewContext(r)

	w.Header().Set("Content-Type", "application/json")
	p.Url, err = servingURL(ctx, path, 0, true)
	if err != nil {
		ep.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ep)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
	return
}

func servingURL(ctx context.Context, filename string, size int, crop bool) (string, error) {
	dbn, err := aefile.DefaultBucketName(ctx)
	if err != nil {
		return "", err
	}
	fullPath := fmt.Sprintf("/gs/%s/%s", dbn, filename)
	key, err := blobstore.BlobKeyForFile(ctx, fullPath)
	if err != nil {
		return "", err
	}
	opts := &image.ServingURLOptions{
		Secure: true,
		Size:   size,
		Crop:   crop,
	}
	url, err := image.ServingURL(ctx, key, opts)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// Call it
//imageURL, _ := servingURL(ctx, "path/to/image/in/gcs/L2Fwc...", 120, true)
