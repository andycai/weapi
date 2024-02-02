package enum

import "errors"

var ErrUnauthorized = errors.New("unauthorized")
var ErrDraftIsInvalid = errors.New("draft is invalid")
var ErrPageIsNotPublish = errors.New("page is not publish")
var ErrPostIsNotPublish = errors.New("post is not publish")
var ErrInvalidPathAndName = errors.New("invalid path and name")
var ErrUploadsDirNotConfigured = errors.New("uploads dir not configured")
