package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/chai2010/webp"

	"github.com/thanhpk/randstr"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handleAdminCrearFotoExtra() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		uploadPath := os.Getenv("UPLOAD_PATH")

		articuloId := r.FormValue("articulo")
		if articuloId == "" {
			log.Error("el id del artículo no puede estar vacío: %s")
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		file, _, err := r.FormFile("foto")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			log.Error("no se ha subido ninguna imagen: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		photoBytes, err := io.ReadAll(file)
		if err != nil {
			log.Error("no se pudo extraer la imagen del formulario: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		image, err := imagenDesdeMime(photoBytes)
		if err != nil {
			log.Error("error extrayendo el tipo MIME de la imagen: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		var fotoEscribir bytes.Buffer
		err = webp.Encode(&fotoEscribir, image, &webp.Options{Quality: 80})
		if err != nil {
			log.Error("error codificando la imagen: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorDatos)
			return
		}

		var salt = randstr.String(5)
		var imagePath = fmt.Sprintf("%s/extra/%s-%s.webp", uploadPath, articuloId, salt)
		err = os.WriteFile(imagePath, fotoEscribir.Bytes(), 0o644)
		if err != nil {
			log.Error("error escribiendo imagen a %s: %s", imagePath, err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorRender)
		}
	}
}
