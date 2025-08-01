# `gitz add`

Agrega archivos al área de staging de Git, ya sea todos los archivos cambiados o con confirmación interactiva.

## 🧾 Descripción

Este comando identifica archivos modificados o nuevos en el repositorio y permite agregarlos de dos maneras:

- Agregarlos todos automáticamente (`--all`)
- Confirmar uno por uno de manera interactiva (`--confirm`)

Internamente, el comando utiliza `git ls-files` para detectar los archivos no stageados y luego usa `git add` para agregarlos.

## ⚙️ Opciones

| Flag        | Descripción                                 |
|-------------|---------------------------------------------|
| `--confirm` | Pregunta por cada archivo si se debe agregar |
| `--all`     | Agrega todos los archivos sin preguntar      |

> ⚠️ Si no se usa ninguna opción, el comando entra en modo automático y agrega todo excepto archivos ignorados o sin cambios.

## 🧪 Ejemplos de uso

### Agregar todos los archivos modificados

```bash
gitz add --all
