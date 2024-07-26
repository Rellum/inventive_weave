package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	pb "github.com/Rellum/inventive_weave/svc/creators/creatorspb"
	"github.com/Rellum/inventive_weave/svc/creators/types"
	"golang.org/x/xerrors"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

//go:embed static/index.html
var homeTmplFile string
var homeTmpl = template.Must(template.New("home").Parse(homeTmplFile))

func home(w http.ResponseWriter, _ *http.Request) {
	homeTmpl.Execute(w, nil)
}

//go:embed templates/activity.html
var activityTmplFile string
var activityTmpl = template.Must(template.New("activity").Parse(activityTmplFile))

//go:embed static
var static embed.FS

func fsHandler() http.Handler {
	sub, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(sub))
}

func activity(cl pb.CreatorsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			slog.ErrorContext(r.Context(), "mime.ParseMediaType error", slog.Any("error", xerrors.Errorf("%w", err)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !strings.HasPrefix(mediaType, "multipart/") {
			slog.ErrorContext(r.Context(), "not a multipart form request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req types.Data
		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				slog.ErrorContext(r.Context(), "file part missing from request")
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if err != nil {
				slog.ErrorContext(r.Context(), "multipart.Reader.NextPart error", slog.Any("error", xerrors.Errorf("%w", err)))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if p.FormName() != "file" {
				continue
			}

			err = json.NewDecoder(p).Decode(&req)
			if err != nil {
				slog.ErrorContext(r.Context(), "http.Server.ListenAndServe error", slog.Any("error", xerrors.Errorf("%w", err)))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			break
		}

		res, err := cl.MostActiveCreators(r.Context(), pb.ToProto(req))
		if err != nil {
			slog.ErrorContext(r.Context(), "http.Server.ListenAndServe error", slog.Any("error", xerrors.Errorf("%w", err)))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data []struct {
			Rank        int
			CreatorStat *pb.CreatorStats
		}
		for i := 0; i < 3; i++ {
			if len(res.CreatorStats) <= i {
				break
			}
			data = append(data, struct {
				Rank        int
				CreatorStat *pb.CreatorStats
			}{i + 1, res.CreatorStats[i]})
		}
		activityTmpl.Execute(w, data)
	}
}
