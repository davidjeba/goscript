package goscript

import (
	"net/http"
	"path/filepath"
)

type AssetManager struct {
	assetDir string
}

func NewAssetManager(assetDir string) *AssetManager {
	return &AssetManager{assetDir: assetDir}
}

func (am *AssetManager) ServeAssets(prefix string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(am.assetDir))
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = filepath.Join(prefix, r.URL.Path)
		fs.ServeHTTP(w, r)
	}
}

