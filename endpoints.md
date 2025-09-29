# Documentación de API - Backend Microservicios

## 📋 Tabla de Contenidos
- [Microservicio: Subscriptions](#microservicio-subscriptions)
- [Microservicio: Albums](#microservicio-albums)
- [Microservicio: Songs](#microservicio-songs)
- [Microservicio: Artists](#microservicio-artists)
- [Microservicio: Playlists](#microservicio-playlists)
- [Microservicio: History (Go)](#microservicio-history-go)
- [Microservicio: Auth (Go)](#microservicio-auth-go)
- [Microservicio: Streaming (Go)](#microservicio-streaming-go)

---

## Microservicio: Subscriptions

### POST `/subscribe/{artist_id}`
**Descripción:** Suscribirse a un artista

**Parámetros:**
- `artist_id` (path, int): ID del artista al que suscribirse
- Usuario autenticado (header/token)

**Respuesta exitosa (201):**
```json
{
  "message": "Suscripción exitosa",
  "artist_id": 123,
  "user_id": 456,
  "created_at": "2025-09-28T10:30:00"
}
```

**Errores:**
- `400`: ValueError (ej: ya suscrito, artista no existe)

---

### DELETE `/unsubscribe/{artist_id}`
**Descripción:** Cancelar suscripción a un artista

**Parámetros:**
- `artist_id` (path, int): ID del artista
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "message": "Desuscripción exitosa",
  "artist_id": 123,
  "user_id": 456
}
```

**Errores:**
- `400`: ValueError (ej: no estaba suscrito)

---

### GET `/status/{artist_id}`
**Descripción:** Verificar si el usuario está suscrito a un artista

**Parámetros:**
- `artist_id` (path, int): ID del artista
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "artist_id": 123,
  "user_id": 456,
  "is_subscribed": true
}
```

---

### GET `/my-subscriptions`
**Descripción:** Obtener todas las suscripciones del usuario con información completa

**Parámetros:**
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "user_id": 456,
  "total_subscriptions": 3,
  "subscriptions": [
    {
      "artist_id": 123,
      "artist_name": "The Beatles",
      "created_at": "2025-09-28T10:30:00"
    },
    {
      "artist_id": 124,
      "artist_name": "Pink Floyd",
      "created_at": "2025-09-27T15:20:00"
    }
  ]
}
```

---

### GET `/my-artists`
**Descripción:** Obtener solo los IDs de los artistas a los que está suscrito el usuario

**Parámetros:**
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "user_id": 456,
  "total_artists": 3,
  "artist_ids": [123, 124, 125]
}
```

---

### GET `/subscription-count`
**Descripción:** Obtener el número total de suscripciones del usuario

**Parámetros:**
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "user_id": 456,
  "subscription_count": 3
}
```

---

## Microservicio: Albums

### POST `/albums/`
**Descripción:** Crear un nuevo álbum

**Parámetros (Form-data):**
- `title` (string, requerido): Título del álbum
- `release_date` (string, opcional): Fecha de lanzamiento (formato: YYYY-MM-DD)
- `cover_image` (file, opcional): Imagen de portada (PNG, JPG, JPEG)
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Álbum creado exitosamente",
  "data": {
    "id": 789,
    "title": "Abbey Road",
    "release_date": "1969-09-26",
    "cover_url": "https://...",
    "user_id": 456,
    "created_at": "2025-09-28T10:30:00"
  }
}
```

**Errores:**
- `400`: Formato de fecha inválido, archivo no es imagen, extensión no permitida

---

### PUT `/albums/{album_id}`
**Descripción:** Actualizar un álbum existente (solo el propietario)

**Parámetros:**
- `album_id` (path, int): ID del álbum
- `title` (form, string, opcional): Nuevo título
- `release_date` (form, string, opcional): Nueva fecha (YYYY-MM-DD)
- `cover_image` (file, opcional): Nueva imagen de portada
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Álbum actualizado correctamente",
  "data": {
    "id": 789,
    "title": "Abbey Road (Remastered)",
    "release_date": "1969-09-26",
    "cover_url": "https://...",
    "user_id": 456,
    "updated_at": "2025-09-28T11:00:00"
  }
}
```

**Errores:**
- `403`: No es el propietario del álbum
- `404`: Álbum no encontrado
- `400`: Formato de fecha o archivo inválido

---

### GET `/albums/{album_id}`
**Descripción:** Obtener información de un álbum específico

**Parámetros:**
- `album_id` (path, int): ID del álbum

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Álbum recuperado correctamente",
  "data": {
    "id": 789,
    "title": "Abbey Road",
    "release_date": "1969-09-26",
    "cover_url": "https://...",
    "user_id": 456,
    "created_at": "2025-09-28T10:30:00"
  }
}
```

**Errores:**
- `404`: Álbum no encontrado

---

### GET `/albums/{album_id}/songs`
**Descripción:** Listar todas las canciones de un álbum

**Parámetros:**
- `album_id` (path, int): ID del álbum

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Canciones del álbum recuperadas correctamente",
  "data": {
    "songs": [
      {
        "id": 1001,
        "title": "Come Together",
        "track_number": 1,
        "duration": 259,
        "audio_url": "https://...",
        "genre_id": 5,
        "album_id": 789
      },
      {
        "id": 1002,
        "title": "Something",
        "track_number": 2,
        "duration": 183,
        "audio_url": "https://...",
        "genre_id": 5,
        "album_id": 789
      }
    ]
  }
}
```

**Errores:**
- `404`: Álbum no encontrado

---

### DELETE `/albums/{album_id}`
**Descripción:** Eliminar un álbum y todas sus canciones (solo el propietario)

**Parámetros:**
- `album_id` (path, int): ID del álbum
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Álbum eliminado correctamente",
  "data": {}
}
```

**Errores:**
- `403`: No es el propietario del álbum
- `404`: Álbum no encontrado

---

### GET `/albums/artist/{artist_id}`
**Descripción:** Obtener todos los álbumes de un artista con información completa

**Parámetros:**
- `artist_id` (path, int): ID del artista

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Álbumes del artista recuperados correctamente",
  "data": {
    "artist_id": 123,
    "total_albums": 2,
    "albums": [
      {
        "id": 789,
        "title": "Abbey Road",
        "release_date": "1969-09-26",
        "cover_url": "https://...",
        "song_count": 17
      },
      {
        "id": 790,
        "title": "Let It Be",
        "release_date": "1970-05-08",
        "cover_url": "https://...",
        "song_count": 12
      }
    ]
  }
}
```

**Errores:**
- `500`: Error al obtener álbumes

---

### GET `/albums/my-albums`
**Descripción:** Obtener todos los álbumes del usuario autenticado

**Parámetros:**
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Mis álbumes recuperados correctamente",
  "data": {
    "artist_id": 123,
    "total_albums": 2,
    "albums": [
      {
        "id": 789,
        "title": "Abbey Road",
        "release_date": "1969-09-26",
        "cover_url": "https://...",
        "song_count": 17
      }
    ]
  }
}
```

**Errores:**
- `404`: No tienes un perfil de artista creado
- `500`: Error al obtener mis álbumes

---

## Microservicio: Songs

### POST `/songs/`
**Descripción:** Crear una nueva canción en un álbum

**Parámetros (Form-data):**
- `title` (string, requerido): Título de la canción
- `album_id` (int, requerido): ID del álbum
- `audio_file` (file, requerido): Archivo de audio (MP3, WAV, etc., máx 50MB)
- `track_number` (int, opcional): Número de pista
- `genre_id` (int, opcional): ID del género musical
- `artist_ids` (string, opcional): IDs de artistas separados por comas (ej: "1,2,3")
- `override_duration` (int, opcional): Duración en segundos (si no se detecta automáticamente)
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Canción creada exitosamente",
  "data": {
    "id": 1001,
    "title": "Come Together",
    "album_id": 789,
    "track_number": 1,
    "duration": 259,
    "audio_url": "https://...",
    "genre_id": 5,
    "created_at": "2025-09-28T10:30:00"
  }
}
```

**Errores:**
- `403`: El álbum no pertenece al usuario
- `404`: Álbum no encontrado
- `400`: Archivo de audio inválido, artist_ids con formato incorrecto, archivo muy grande (>50MB)

---

### GET `/songs/{song_id}`
**Descripción:** Obtener información de una canción específica

**Parámetros:**
- `song_id` (path, int): ID de la canción

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Canción recuperada correctamente",
  "data": {
    "id": 1001,
    "title": "Come Together",
    "album_id": 789,
    "track_number": 1,
    "duration": 259,
    "audio_url": "https://...",
    "genre_id": 5,
    "artists": [
      {"id": 123, "name": "The Beatles"}
    ]
  }
}
```

**Errores:**
- `404`: Canción no encontrada

---

### PUT `/songs/{song_id}`
**Descripción:** Actualizar información de una canción (solo el propietario)

**Parámetros:**
- `song_id` (path, int): ID de la canción
- `title` (form, string, opcional): Nuevo título
- `track_number` (form, int, opcional): Nuevo número de pista
- `genre_id` (form, int, opcional): Nuevo género
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Canción actualizada correctamente",
  "data": {
    "id": 1001,
    "title": "Come Together (Remastered)",
    "album_id": 789,
    "track_number": 1,
    "duration": 259,
    "audio_url": "https://...",
    "genre_id": 5,
    "updated_at": "2025-09-28T11:00:00"
  }
}
```

**Errores:**
- `403`: No es el propietario de la canción
- `404`: Canción no encontrada

---

### DELETE `/songs/{song_id}`
**Descripción:** Eliminar una canción (solo el propietario)

**Parámetros:**
- `song_id` (path, int): ID de la canción
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Canción eliminada correctamente",
  "data": {}
}
```

**Errores:**
- `403`: No es el propietario de la canción
- `404`: Canción no encontrada

---

## Microservicio: Artists

### POST `/register`
**Descripción:** Registrar un nuevo perfil de artista para el usuario autenticado

**Parámetros (Form-data):**
- `artist_name` (string, requerido): Nombre del artista
- `bio` (string, opcional): Biografía del artista
- `social_links` (string, opcional): Enlaces a redes sociales en formato JSON string
- `profile_pic_file` (file, opcional): Foto de perfil del artista
- Usuario autenticado (header/token)

**Ejemplo de social_links:**
```json
{
  "instagram": "https://instagram.com/artist",
  "twitter": "https://twitter.com/artist",
  "spotify": "https://open.spotify.com/artist/..."
}
```

**Respuesta exitosa (201):**
```json
{
  "success": true,
  "message": "Artista registrado correctamente",
  "data": {
    "id": 123,
    "artist_name": "The Beatles",
    "bio": "Legendary rock band from Liverpool",
    "profile_pic": "https://...",
    "social_links": {
      "instagram": "https://instagram.com/thebeatles",
      "twitter": "https://twitter.com/thebeatles"
    },
    "user_id": 456,
    "created_at": "2025-09-28T10:30:00"
  }
}
```

**Errores:**
- `409`: Ya existe un artista registrado para este usuario
- `400`: Formato inválido en social_links

---

### GET `/me`
**Descripción:** Obtener el perfil de artista del usuario autenticado

**Parámetros:**
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Perfil de artista obtenido",
  "data": {
    "id": 123,
    "artist_name": "The Beatles",
    "bio": "Legendary rock band from Liverpool",
    "profile_pic": "https://...",
    "social_links": {
      "instagram": "https://instagram.com/thebeatles",
      "twitter": "https://twitter.com/thebeatles"
    },
    "user_id": 456,
    "created_at": "2025-09-28T10:30:00"
  }
}
```

**Errores:**
- `404`: Artista no encontrado

---

### PUT `/me`
**Descripción:** Actualizar el perfil de artista del usuario autenticado

**Parámetros (Form-data):**
- `artist_name` (string, opcional): Nuevo nombre del artista
- `bio` (string, opcional): Nueva biografía
- `social_links` (string, opcional): Nuevos enlaces sociales (JSON string)
- `profile_pic_file` (file, opcional): Nueva foto de perfil
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Artista actualizado correctamente",
  "data": {
    "id": 123,
    "artist_name": "The Beatles (Remastered)",
    "bio": "Updated biography",
    "profile_pic": "https://...",
    "social_links": {
      "instagram": "https://instagram.com/thebeatles",
      "twitter": "https://twitter.com/thebeatles",
      "spotify": "https://open.spotify.com/artist/..."
    },
    "user_id": 456,
    "updated_at": "2025-09-28T11:00:00"
  }
}
```

**Errores:**
- `404`: Artista no encontrado
- `400`: Formato inválido en social_links

---

### DELETE `/me`
**Descripción:** Eliminar el perfil de artista del usuario autenticado

**Parámetros:**
- Usuario autenticado (header/token)

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "message": "Artista eliminado correctamente",
  "data": {}
}
```

**Errores:**
- `404`: Artista no encontrado

---

## Notas Generales

### Autenticación
Todos los endpoints requieren autenticación mediante token. El `user_id` se obtiene automáticamente del `request.state.user["user_id"]`.

### Validaciones de Propiedad
- Los endpoints de actualización y eliminación validan que el usuario sea el propietario del recurso
- Las canciones solo pueden ser creadas en álbumes del usuario
- Los álbumes y canciones solo pueden ser modificados por sus propietarios

### Formatos de Respuesta
Todos los endpoints siguen el patrón:
```json
{
  "success": true/false,
  "message": "Mensaje descriptivo",
  "data": { /* contenido */ }
}
```

### Límites
- Imágenes: PNG, JPG, JPEG
- Audio: Máximo 50MB
- Formatos de fecha: ISO 8601 (YYYY-MM-DD)

---

## Microservicio: History (Go)

### GET `/history`
**Descripción:** Obtener el historial de canciones reproducidas por el usuario autenticado

**Parámetros:**
- Usuario autenticado (JWT - `user_id` en contexto)

**Respuesta exitosa (200):**
```json
[
  {
    "id": 1,
    "user_id": 456,
    "song_id": 1001,
    "played_at": "2025-09-28T10:30:00Z",
    "song_title": "Come Together",
    "artist_name": "The Beatles",
    "album_title": "Abbey Road"
  },
  {
    "id": 2,
    "user_id": 456,
    "song_id": 1002,
    "played_at": "2025-09-28T09:15:00Z",
    "song_title": "Something",
    "artist_name": "The Beatles",
    "album_title": "Abbey Road"
  }
]
```

**Errores:**
- `401`: Usuario no autenticado
- `500`: Error obteniendo historial, user_id inválido

---

## Microservicio: Auth (Go)

### POST `/register`
**Descripción:** Registrar un nuevo usuario en el sistema

**Body (JSON):**
```json
{
  "username": "john_lennon",
  "email": "john@beatles.com",
  "password": "imagine123",
  "birth_date": "1940-10-09",
  "first_name": "John",
  "last_name": "Lennon"
}
```

**Respuesta exitosa (201):**
```json
{
  "message": "usuario creado exitosamente",
  "user": {
    "id": 456,
    "username": "john_lennon",
    "email": "john@beatles.com",
    "first_name": "John",
    "last_name": "Lennon",
    "birth_date": "1940-10-09",
    "created_at": "2025-09-28T10:30:00Z"
  }
}
```

**Errores:**
- `400`: Error en validación de datos, fecha inválida (formato esperado YYYY-MM-DD)
- `409`: Usuario o email ya registrados
- `500`: No se pudo encriptar la contraseña

---

### POST `/login`
**Descripción:** Iniciar sesión y obtener tokens de acceso

**Body (JSON):**
```json
{
  "email": "john@beatles.com",
  "password": "imagine123"
}
```

**Respuesta exitosa (200):**
```json
{
  "message": "login exitoso",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "id": 456,
    "username": "john_lennon",
    "email": "john@beatles.com",
    "first_name": "John",
    "last_name": "Lennon"
  }
}
```

**Cookies establecidas:**
- `refresh_token` (HttpOnly, SameSite=Strict): Token de renovación

**Errores:**
- `400`: Error en formato de datos
- `401`: Usuario no encontrado, contraseña incorrecta
- `500`: Error generando tokens

---

### POST `/refresh`
**Descripción:** Renovar el access token usando el refresh token almacenado en cookie

**Parámetros:**
- Cookie `refresh_token` (HttpOnly)

**Respuesta exitosa (200):**
```json
{
  "message": "token renovado",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

**Cookies actualizadas:**
- `refresh_token` (HttpOnly, SameSite=Strict): Nuevo refresh token

**Errores:**
- `401`: Refresh token requerido, refresh token inválido o expirado

---

### POST `/logout`
**Descripción:** Cerrar sesión eliminando el refresh token del usuario

**Parámetros:**
- Usuario autenticado (JWT - `user_id` en contexto)

**Respuesta exitosa (200):**
```json
{
  "message": "logged_out"
}
```

**Errores:**
- `401`: Usuario no autenticado
- `500`: Error cerrando sesión, error procesando ID de usuario

---

### GET `/me`
**Descripción:** Obtener la información del usuario autenticado

**Parámetros:**
- Usuario autenticado (JWT - `user_id` en contexto)

**Respuesta exitosa (200):**
```json
{
  "user": {
    "id": 456,
    "username": "john_lennon",
    "email": "john@beatles.com",
    "first_name": "John",
    "last_name": "Lennon",
    "birth_date": "1940-10-09",
    "created_at": "2025-09-28T10:30:00Z",
    "updated_at": "2025-09-28T10:30:00Z"
  }
}
```

**Errores:**
- `401`: No autorizado
- `404`: Usuario no encontrado
- `500`: Error procesando ID de usuario

---

### PUT `/update`
**Descripción:** Actualizar información del usuario autenticado

**Parámetros:**
- Usuario autenticado (JWT - `user_id` en contexto)

**Body (JSON) - Todos los campos son opcionales:**
```json
{
  "username": "john_lennon_updated",
  "email": "newemail@beatles.com",
  "password": "newpassword123",
  "first_name": "John Winston",
  "last_name": "Lennon"
}
```

**Respuesta exitosa (200):**
```json
{
  "message": "usuario actualizado correctamente",
  "user": {
    "id": 456,
    "username": "john_lennon_updated",
    "email": "newemail@beatles.com",
    "first_name": "John Winston",
    "last_name": "Lennon",
    "birth_date": "1940-10-09",
    "updated_at": "2025-09-28T11:00:00Z"
  }
}
```

**Errores:**
- `400`: Username ya está en uso, email ya está en uso, formato de email inválido, la contraseña debe tener al menos 6 caracteres
- `401`: Usuario no autenticado
- `404`: Usuario no encontrado
- `500`: Error en la autenticación

---

## Microservicio: Streaming (Go)

### GET `/stream?id={song_id}`
**Descripción:** Transmitir audio de una canción con soporte para HTTP Range Requests (streaming parcial)

**Parámetros:**
- `id` (query, int, requerido): ID de la canción a reproducir
- Usuario autenticado (JWT - `user_id` en contexto, opcional pero recomendado)
- Header `Range` (opcional): Para solicitar rangos específicos del archivo (ej: `bytes=0-1023`)

**Comportamiento:**

**1. Sin Range Header (streaming completo):**
- Devuelve el archivo completo de audio
- Registra la reproducción inmediatamente
- Status Code: `200 OK`

**2. Con Range Header (streaming parcial):**
- Devuelve solo el rango solicitado del archivo
- Registra la reproducción cuando se alcanza el 30% del archivo
- Status Code: `206 Partial Content`

**Headers de respuesta:**
```
Content-Type: audio/mpeg
Content-Length: [tamaño en bytes]
Content-Range: bytes [inicio]-[fin]/[total]  (solo con Range)
```

**Ejemplo de solicitud con Range:**
```
GET /stream?id=1001
Range: bytes=0-1048575
```

**Respuesta exitosa (206):**
```
Headers:
Content-Type: audio/mpeg
Content-Length: 1048576
Content-Range: bytes 0-1048575/5242880

Body: [datos binarios de audio]
```

**Respuesta exitosa sin Range (200):**
```
Headers:
Content-Type: audio/mpeg
Content-Length: 5242880

Body: [archivo completo de audio]
```

**Eventos disparados:**
- Al reproducir sin Range: Evento inmediato `SongPlayed(user_id, song_id)`
- Al reproducir con Range: Evento `SongPlayed(user_id, song_id)` cuando se alcanza el 30% del archivo total

**Errores:**
- `400`: ID de canción requerido, ID inválido, error al abrir archivo
- `404`: Canción no encontrada, archivo no encontrado
- `416`: Rango inválido (Range Not Satisfiable)

**Notas:**
- El endpoint soporta reproducción continua y seeking
- El historial solo se registra si el usuario está autenticado
- Se usa buffer de 32KB para la transmisión
- El evento de reproducción se dispara de forma asíncrona (goroutine)

---

## Notas sobre Microservicios en Go

### Framework
Los microservicios en Go utilizan **Gin** como framework web.

### Autenticación
- La autenticación se maneja mediante JWT (JSON Web Tokens)
- El `user_id` se extrae del contexto de Gin: `c.Get("user_id")`
- El refresh token se almacena en cookies HttpOnly para mayor seguridad

### Manejo de tipos
Los handlers de Go incluyen conversión robusta de tipos para `user_id`, soportando `uint`, `float64` e `int`.

### Cookies de seguridad
- `SameSite=Strict`: Protección contra CSRF
- `HttpOnly=true`: No accesible desde JavaScript
- `Secure=false`: Para desarrollo local (cambiar a `true` en producción)