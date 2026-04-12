package handlers

import (
	"net/url"
	"strings"

	"github.com/mazong/angel_event/internal/models"
)

func normalizeAssetURL(raw string) string {
	if raw == "" {
		return raw
	}

	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "data:") {
		return raw
	}

	parts := strings.SplitN(raw, "?", 2)
	path := parts[0]
	query := ""
	if len(parts) == 2 {
		query = "?" + parts[1]
	}

	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if segment == "" {
			continue
		}

		decoded, err := url.PathUnescape(segment)
		if err != nil {
			decoded = segment
		}
		segments[i] = url.PathEscape(decoded)
	}

	return strings.Join(segments, "/") + query
}

func normalizeSiteContentAssets(items []models.SiteContent) {
	for i := range items {
		if items[i].Type == "image" {
			items[i].Value = normalizeAssetURL(items[i].Value)
		}
	}
}
