# Music Streaimng Backend - Guía de Configuración y Ejecución

Este proyecto utiliza Docker Compose para orquestar múltiples servicios de backend. Sigue esta guía para configurar y ejecutar correctamente el sistema.

## Requisitos Previos

- Docker instalado en tu sistema
- Docker Compose instalado
- Base de datos con el esquema correcto instalado

## Configuración del Entorno

### 1. Archivo de Configuración (.env)

Crea un archivo `.env` en la **raíz del proyecto** (mismo directorio que el archivo `docker-compose.yml`) con las siguientes variables:

```env
# ===== Base de Datos =====
DB_URL_PY=postgresql+asyncpg://<USUARIO>:<PASSWORD>@<HOST>:<PUERTO>/<NOMBRE_DB>
DB_URL=postgresql://<USUARIO>:<PASSWORD>@<HOST>:<PUERTO>/<NOMBRE_DB>

# ===== JWT =====
# Clave secreta para firmar y validar tokens JWT.
# Usa una cadena larga y aleatoria en producción.
JWT_SECRET=<TU_SECRETO_JWT>
# Algoritmo para firmar JWT (ejemplo: HS256).
JWT_ALGORITHM=HS256

# ===== Tokens (solo auth-service) =====
# Tiempo de vida del access token (ej: 15m, 1h).
ACCESS_TOKEN_TTL=15m
# Tiempo de vida del refresh token (ej: 7d, 168h).
REFRESH_TOKEN_TTL=168h

# ===== Rutas de contenido =====
# Ruta absoluta o relativa donde se almacenará contenido multimedia.
CONTENT_BASE_PATH=/ruta/a/storage

# ===== RabbitMQ =====
# URL de conexión al broker de mensajes.
# Formato: amqp://<USUARIO>:<PASSWORD>@<HOST>:<PUERTO>/
RABBITMQ_URL=amqp://<USUARIO>:<PASSWORD>@<HOST>:5672/

# ===== Puertos de servicios =====
# Cada microservicio puede correr en un puerto diferente.
ARTIST_PORT=8001
CONTENT_PORT=8002
PLAYLIST_PORT=8004
SEARCH_PORT=8005
SUBSCRIPTION_PORT=8006
AUTH_PORT=8080
HISTORY_PORT=8081
STREAMING_PORT=8082
```


### 2. Configuración de Variables Importantes

#### Base de Datos
- **DB_URL_PY**: URL para servicios Python (utiliza el driver `asyncpg`)
- **DB_URL**: URL para servicios Go (formato estándar de PostgreSQL)
- **Nota**: Ambas URLs apuntan a la misma base de datos, pero tienen formatos diferentes debido a los drivers utilizados

#### Seguridad JWT
- **JWT_SECRET**: Clave secreta para la generación y validación de tokens JWT
- **JWT_ALGORITHM**: Algoritmo utilizado para JWT (HS256 por defecto y el único con lo que se ha probado, cualquier otro algoritmo puede dar error)

#### Tokens de Acceso
- **ACCESS_TOKEN_TTL**: Tiempo de vida del token de acceso (15 minutos)
- **REFRESH_TOKEN_TTL**: Tiempo de vida del token de renovación (168 horas = 7 días)

#### Almacenamiento
- **CONTENT_BASE_PATH**: Ruta donde se almacenarán los archivos de contenido

#### Message Queue
- **RABBITMQ_URL**: URL de conexión a RabbitMQ para comunicación entre servicios

## Estructura de Puertos

El sistema utiliza los siguientes puertos para cada servicio:

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| Artist Service | 8001 | Gestión de artistas |
| Content Service | 8002 | Manejo de contenido multimedia |
| Playlist Service | 8004 | Administración de playlists |
| Search Service | 8005 | Servicio de búsqueda |
| Subscription Service | 8006 | Gestión de suscripciones |
| Auth Service | 8080 | Autenticación y autorización |
| History Service | 8081 | Historial de reproducción |
| Streaming Service | 8082 | Streaming de contenido |
> ℹ️ **Nota:** Los puertos listados son solo valores por defecto.  
> Pueden modificarse sin problema alguno, ya que ni los contenedores  
> ni los archivos de configuración están hardcodeados a estos puertos

## Documentación de la API (Postman)

Para facilitar las pruebas de los microservicios, se incluyen las colecciones de Postman exportadas.  

📂 Las encontrarás en: `assets/postman_collections/`  

Cada archivo `.json` corresponde a una colección de endpoints agrupados por servicio.  
Puedes importarlos directamente en Postman siguiendo estos pasos:

1. Abre Postman
2. Ve al menú **Import**
3. Selecciona los archivos `.json` desde `assets/postman_collections/`
4. Las colecciones aparecerán organizadas en tu espacio de trabajo

> ℹ️ **Nota:** Se recomienda usar la versión de exportación **2.1** de Postman para mayor compatibilidad.

## Ejecución del Proyecto

### 1. Preparación
Asegúrate de que el archivo `.env` esté configurado correctamente en la raíz del proyecto.
```markdown
📦 Streaming Backend
├── 📂 artist-service/
├── 📂 auth-service/
├── 📂 content-service/
├── 📂 history-service/
├── 📂 playlist-service/
├── 📂 search-service/
├── 📂 streaming-service/
├── 📂 subscription-service/
├── 📄 docker-compose.yml
└── ⚙️ .env  ← Crear este archivo aquí
```

### 2. Construir y Ejecutar
```bash
# Construir las imágenes y ejecutar los servicios
docker-compose up --build

# Para ejecutar en segundo plano
docker-compose up -d --build
```

### 3. Verificar el Estado
```bash
# Ver el estado de los contenedores
docker-compose ps

# Ver logs de todos los servicios
docker-compose logs

# Ver logs de un servicio específico
docker-compose logs <nombre-del-servicio>
```

### 4. Detener los Servicios
```bash
# Detener los servicios
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v
```

## Resolución de Problemas

### Problemas Comunes

1. **Error de conexión a la base de datos**
   - Verifica que las credenciales en `DB_URL` y `DB_URL_PY` sean correctas
   - Asegúrate de que la base de datos esté accesible desde tu red

2. **Puertos ocupados**
   - Verifica que los puertos definidos en el `.env` no estén siendo utilizados por otros servicios
   - Cambia los puertos en el archivo `.env` si es necesario

3. **Problemas con RabbitMQ**
   - Asegúrate de que el servicio RabbitMQ esté ejecutándose
   - Verifica la configuración de `RABBITMQ_URL`

4. **Archivos no encontrados**
   - Verifica que la ruta `CONTENT_BASE_PATH` exista y sea accesible
   - Asegúrate de que los permisos de directorio sean correctos

### Comandos Útiles

```bash
# Reconstruir un servicio específico
docker-compose up --build <nombre-del-servicio>

# Acceder a un contenedor en ejecución
docker-compose exec <nombre-del-servicio> /bin/bash

# Ver recursos utilizados
docker-compose top
```

## Notas Adicionales

- El archivo `.env` contiene información sensible y **NO debe** ser incluido en el control de versiones
- Asegúrate de configurar correctamente los firewalls y permisos de red para los puertos utilizados
- Para ambiente de producción, considera utilizar secretos más seguros para `JWT_SECRET`
- Los servicios están diseñados para comunicarse entre sí a través de la red interna de Docker Compose