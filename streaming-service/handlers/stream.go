package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"streaming-service/events" // publisher
	"streaming-service/services"
	"streaming-service/stream"
	"streaming-service/utils"

	"github.com/gin-gonic/gin"
)

// StreamSongHandler transmite la canción usando Range
func StreamSongHandler(songService services.SongService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el id desde query (?id=123)
		idStr := c.Query("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id de canción requerido"})
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		// Buscar URL de la canción
		songPath, err := songService.GetSongURL(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "canción no encontrada"})
			return
		}

		// Abrir archivo desde storage
		audioSource, err := stream.AudioFactory("local", songPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		file, err := audioSource.GetReader()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "archivo no encontrado"})
			return
		}
		defer file.Close()

		size := audioSource.GetSize()
		rangeHeader := c.GetHeader("Range")

		// ✅ Caso sin Range: mandar todo el archivo y disparar evento
		if rangeHeader == "" {
			c.Header("Content-Type", "audio/mpeg")
			c.Header("Content-Length", fmt.Sprintf("%d", size))
			c.Status(http.StatusOK)
			io.Copy(c.Writer, file)

			// Emitir evento (ya escuchó toda la canción)
			userIDAny, exists := c.Get("user_id")
			if exists {
				if uid, ok := userIDAny.(float64); ok && uid > 0 {
					userID := uint(uid)
					fmt.Printf("✅ Publicando evento (sin Range): userID=%d, songID=%d\n", userID, id)
					go events.PublishSongPlayed(userID, uint(id))
				}
			} else {
				fmt.Println("⚠️ No se encontró user_id en el contexto, no se publicará el evento (sin Range)")
			}
			return
		}

		// ✅ Caso con Range
		start, end := utils.ParseRange(rangeHeader, size)
		if start > end || start >= size {
			c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": "rango inválido"})
			return
		}

		length := end - start + 1
		c.Header("Content-Type", "audio/mpeg")
		c.Header("Content-Length", fmt.Sprintf("%d", length))
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
		c.Status(http.StatusPartialContent)

		// Posicionarse en el inicio
		file.Seek(start, 0)

		// Enviar en chunks
		buf := make([]byte, 32*1024) // 32KB buffer
		var sent int64 = 0
		eventSent := false

		for sent < length {
			toRead := int64(len(buf))
			if length-sent < toRead {
				toRead = length - sent
			}
			n, err := file.Read(buf[:toRead])
			if n > 0 {
				c.Writer.Write(buf[:n])
				sent += int64(n)

				// 👇 Emitir evento si supera el 30% del archivo completo
				if !eventSent && float64(sent+start)/float64(size) >= 0.3 {
					fmt.Printf("🔎 Superado 30%% del rango: sent=%d, length=%d, size=%d\n", sent, length, size)

					userIDAny, exists := c.Get("user_id")
					if !exists {
						fmt.Println("⚠️ No se encontró user_id en el contexto, no se publicará el evento")
					} else {
						if uid, ok := userIDAny.(float64); ok && uid > 0 {
							userID := uint(uid)
							fmt.Printf("✅ Publicando evento: userID=%d, songID=%d\n", userID, id)
							go events.PublishSongPlayed(userID, uint(id))
							eventSent = true
						}
					}
				}
			}
			if err != nil {
				break
			}
		}
	}
}
