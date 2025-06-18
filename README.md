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
}
```

El siguiente ejemplo le creará una ruta "user" con los siguientes metodos HTTP: **GET, POST, PUT, DELETE**. 


## Argumentos al ejecutar el script 
```bash
# default 
go run main --path=./api/(name).json \
            --mode=static \
            --port=8000


# dev 🪒
go run main --path=[ruta donde esta el (name).json o yml para crear tu api] \
            --mode=[static | watch] \
            --port=[8000]
```


## Prerequisitos

1. SQLite 
2. 