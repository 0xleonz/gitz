# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v0.2.0]

### Added
- Comando `gitz add` con soporte para `--confirm`, `--dry-run`, `--verbose`.
- Comando `gitz commit` que usa `commitMessage.yml` como plantilla.
- Comando `gitz message` para crear, editar y mostrar mensajes de commit.
- Comando `gitz push` con confirmación de ramas remotas.
- Comando `gitz quick` para `check + add + commit + push` en una sola línea.
- Soporte para flags globales `--verbose`, `--dry-run`.
- Manejo de mensajes cortos con `-s` / `--short` en `commit` y `message`.
- Soporte interactivo para agregar `subject`, `description`, `changes`, `footer`, `issue`, etc.
- Colorización de campos existentes y ejemplos sugeridos cuando no hay datos.

---

## [0.1.0] - 2025-07-xx

### Added
- Estructura base del CLI usando Cobra.
- Lectura de `info.yml` para metainformación del repo.

