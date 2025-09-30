package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "go-short-url/internal/lib/api/response"
	"go-short-url/internal/lib/logger/sl"
	"go-short-url/internal/lib/random"
	"go-short-url/internal/storage"
	"log/slog"
	"net/http"
)

// если вдруг забуду
//Alias — это короткий код, который заменяет длинный URL

// структура для запроса POST /save
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// структура ответа
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config if needed
const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// Принимает запрос POST /save.
// Проверяет входные данные (JSON + валидация).
// Генерирует alias, если не передан.
// Сохраняет в хранилище.
// Отправляет JSON-ответ с alias.
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		// Декодируем JSON-запрос
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Валидируем структуру запроса
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		// Генерация alias, если не передан
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		// Сохраняем URL через urlSaver
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		// Отправляем успешный JSON-ответ с alias
		responseOK(w, r, alias)

	}
}

// формирует успешный ответ с alias
func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
