# Goserver **(v0.0.1)**


### Compilation

- Windows:
    ```
    set GOOS=windows; go build -o goserver.exe
    ```
- Linux:
    ```
    set GOOS=linux; go build -o goserver.exe
    ```
- Macos:
    ```
    set GOOS=macos; go build -o goserver.exe
    ```

Tener en cuenta que tambien puede compilarlo respecto a su arquitectura.

## Estructura basica para la devolucion de la api 

```json
{
    "user": [
        {
            "id": 1,
            "name": "John Doe",
            "email": "admin@gmail.com",
            "work" :" programmer"
        },
        {
            "id": 2,
            "name": "Cristian vega",
            "email": "vega@gmail.com",
            "work" :    "medico"
        }
    ],
     "posts": 
    {
        "request": ["POST", "GET"],
        "middleware": {
            "auth": "bearer",
            "logging": true,
            "db": true,
            "security": ["csrf", "xss"]
        },
        "response": [
            {
                "id": "100001_500001",
                "created_time": "2025-07-26T18:03:21+0000",
                "message": "¡Hola a todos! Este es mi primer post 😄",
                "from": {
                  "id": "100001",
                  "name": "Lucía Fernández"
                },
                "permalink_url": "https://facebook.com/100001/posts/500001"
            },
            {
                "id": "100002_500002",
                "created_time": "2025-07-26T18:05:10+0000",
                "message": "¡Feliz viernes! ¿Qué planes tienen para el fin de semana? 🎉",
                "from": {
                  "id": "100002",
                  "name": "Carlos Pérez"
                },
                "permalink_url": "https://facebook.com/100002/posts/500002"
            },
            {
                "id": "100003_500003",
                "created_time": "2025-07-26T18:07:45+0000",
                "message": "¡Saludos desde la playa! 🏖️",
                "from": {
                  "id": "100003",
                  "name": "Ana Torres"
                },
                "permalink_url": "https://facebook.com/100003/posts/500003"
            }
        ],
        "schema": {
            "table_name": "posts",
            "id": "primary_key",
            "created_time": "datetime",
            "message": "text",
            "from": {
                "table_name": "users",
                "id": "primary_key",
                "email": "varchar,255|unique",
                "name": "varchar,255"
            },
            "permalink_url": "url"
        }
    }
}
```

El siguiente ejemplo le creará una ruta "user" con los siguientes metodos HTTP: **GET, POST, PUT, DELETE**. 


## Argumentos al ejecutar el script 
```bash
# default 
go run watcher.go --path=./api/(name).json \
            --mode=static \
            --port=8000


# dev 🪒
go run watcher.go --path=[ruta donde esta el (name).json o yml para crear tu api] \
            --mode=[static | watch] \
            --port=[8000]
```


## Prerequisitos

1. SQLite 