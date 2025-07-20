# gitz

`gitz` es una herramienta CLI para facilitar flujos de trabajo con Git,
especialmente en proyectos que requieren subir ramas específicas a múltiples
remotos de manera consistente.

## Características

- Lee configuración desde un archivo `info.yml` localizado en la raíz del repositorio.
- Permite definir ramas a subir y los remotos correspondientes.
- Ejecuta `git push` a múltiples remotos de forma automática.
- Soporte básico para comandos como `init` y `push`.

## Commands

```bash
gitz init
gitz message
gitz push
```

```
```
## Estructura del archivo `info.yml`

```yaml
description: 'Repositorio para ejemplo'
ramas:
  - origin
  - github
branches:
  - main
  - dev
remote-branches:
  - branch: dev
    remotes:
      - origin
