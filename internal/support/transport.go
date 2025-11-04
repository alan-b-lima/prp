package support

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/alan-b-lima/prp/internal/xerrors"
	"github.com/alan-b-lima/prp/pkg/errors"
	uuidpkg "github.com/alan-b-lima/prp/pkg/uuid"
)

var (
	reContentTypeApplicationJson = regexp.MustCompile(`^\s*(\*/\*|application/(json|\*))\s*(;.*)?\s*$`)
	reAcceptApplicationJson      = regexp.MustCompile(`(^|.*,)\s*(\*/\*|application/(json|\*))\s*(;.*)?\s*($|,.*)`)
)

func LimitAndOffset(offsetStr, limitStr string, offset, limit *int) error {
	var errs [2]error

	if offsetStr != "" {
		*offset, errs[0] = strconv.Atoi(offsetStr)
	}
	if limitStr != "" {
		*limit, errs[1] = strconv.Atoi(limitStr)
	}

	if errs != [2]error{nil, nil} {
		return xerrors.ErrBadOffsetOrLimit.New(errors.Join(errs[:]...))
	}

	return nil
}

func UUIDFromString(uuidStr string) (uuidpkg.UUID, error) {
	uuid, err := uuidpkg.FromString(uuidStr)
	if err != nil {
		return uuidpkg.UUID{}, xerrors.ErrBadUUID.New(err)
	}

	return uuid, nil
}

func SessionCookie(cookie string, w http.ResponseWriter, r *http.Request) (uuidpkg.UUID, error) {
	s, err := r.Cookie(cookie)
	if err != nil {
		return uuidpkg.UUID{}, xerrors.ErrUnauthenticatedUser
	}

	uuid, err := UUIDFromString(s.Value)
	if err != nil {
		return uuidpkg.UUID{}, xerrors.ErrBadUUID.New(err)
	}

	return uuid, nil
}

func DecodeJSON(req any, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return xerrors.ErrNoContentType
	}

	if !reContentTypeApplicationJson.MatchString(contentType) {
		return xerrors.ErrNoContentType
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err, ok := err.(*json.SyntaxError); ok {
			return xerrors.ErrJsonSyntax.New("JSON syntax error at"+strconv.FormatInt(err.Offset, 10), nil)
		}

		return err
	}

	return nil
}

func EncodeJSON(res any, status int, w http.ResponseWriter, r *http.Request) error {
	accept := r.Header.Get("Accept")
	if !reAcceptApplicationJson.MatchString(accept) {
		return xerrors.ErrNotAcceptableJson
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(res); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := io.Copy(w, &b); err != nil {
		return err
	}

	return nil
}
