# Pokemon API Backend

API para consultar Pokemon, con cache, busquedas rapidas y menos llamadas a la API oficial.

## Endpoints

### 1. Obtener lista de Pokemon
```
GET /api/v1/pokemon?limit=20&offset=0
```


**Parametros**
- `limit`: Cantidad de Pokemon a obtener (maximo 100, por defecto 20)
- `offset`: Desde que posicion empezar (por defecto 0)

### 2. Buscar Pokemon por nombre
```
GET /api/v1/pokemon/search?q=pika&limit=10&offset=0
```


**Parametros**
- `q`: Texto a buscar (requerido)
- `limit`: Cantidad de resultados (maximo 50, por defecto 10)
- `offset`: Desde que posicion empezar (por defecto 0)

### 3. Obtener Pokemon por ID
```
GET /api/v1/pokemon/25
```

### 4. Obtener Pokemon por nombre exacto
```
GET /api/v1/pokemon/name/pikachu
```

### 5. Estado de la aplicación
```
GET /health
```
Verifica que la API este funcionando correctamente

## Instalación

Antes de comenzar, asegúrate de tener instalado
- Go 1.23 o superior

#

1. **Clona o descarga el proyecto**
   ```bash
   git clone https://github.com/gerardstrujills/backend.git
   cd backend
   ```

2. **Instala las dependencias**
   ```bash
   go mod tidy
   ```

3. **Ejecuta la aplicación**
   ```bash
   go run main.go
   ```

4. **API disponible**
   ```bash
   http://localhost:8080
   ```