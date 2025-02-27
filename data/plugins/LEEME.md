# Cómo crear un plugin
Instrucciones breves para crear un plugin compatible con Achaean.

## Requisitos del sistema de ficheros
El directorio del plugin debe incluir al menos:
- `plugin.json`: Contendrá la configuración y parámetros del plugin.
- Ejecutable: Ha de tener el atributo +x aunque sea un script.
- `README` o `README.md` o `README.txt` (opcional): Para incluir ayuda adicional, licencia, etc.

## Requisitos del código
Los plugins deben cumplir con:
- (opciona) Leer sus parámetros desde la linea de comandos si los tuviera.
- Usar `stderr` para mensajes de error.
- Usar `stdout` para mensajes estándar.
- Tiene que usar un **named pipe**: `/tmp/plugin_progress-PLUGIN_ID`. Más info sobre los named pipes en https://en.wikipedia.org/wiki/Named_pipe.
  - Si hay un progreso medible, imprimir el porcentaje (`0-100`) en el mismo.
  - Al finalizar el plugin, debe escribir `DONE` en el pipe, para indicarle a Achaean que ya ha terminado.
  - Hay ejemplos en diferente lenguajes sobre cómo usar el **named pipe**.
- Si es un lenguaje compilado, se recomienda un **linkado estático** para minimizar dependencias externas.

## Variables del sistema
- Se creará una **variable de entorno** `PLUGIN_ID` que contendrá el ID del plugin, así se evita tenerlo "hardcoded" en el ejecutable.
- Se creará otra **variable de entorno** `PROGRESS_PIPE` que contiene el nombre del fichero de dicho pipe.

## Formato de `plugin.json`
### Claves obligatorias
```json
{
    "plugin": {
        "id": "tcp-scan",
        "entrypoint": "tcp-scan.py",
        "category": "Networking",
        "name": "TCP Scan",
        "description": "Escanea puertos TCP.",
        "version": "1.0.0",
        "author": "John Doe",
    }
}
```
- id: Identificador único (ejemplo: tcp-scan, backdoor-seeker).
- entrypoint: Nombre del ejecutable (debe tener permisos de ejecución).
- category: Categoría general del plugin.
- name: Nombre breve.
- description: Descripción breve de la funcionalidad.

### Claves opcionales
El plugin si lo requiere, puede tener claves adicionales que representan sus parámetros, las cuales estarán dentro de `parameters`. Ejemplo:
```json

"parameters": {
    "--host": {
        "name": "Hostname",
        "value": "nmap.scanme.org"
    },
    "--ports-range=": {
        "name": "Ports Range",
        "value": [21, 22, 80],
        "limits": [1, 65535]
    },
    "-w": {
        "name": "Workers",
        "value": 100,
        "limits": [1, 200]
    },
    "--show-closed=": {
        "name": "Show Closed ports",
        "value": false
    }
}
```
#### Desglose de las claves

- La clave del parámetro, "--host", "--ports-range=", etc.: Son los parámetros que se incluirán en la linea de comandos. Si terminara en `=` no se le añadiría un espacio al valor. Algunos ejemplos:
  - "--host"  ===> `--host 127.0.0.1`
  - "--ports-range=" ===> `--ports-range=21,22,80,8080`
- "name": Etiqueta que se mostrará en la interface.
- "value": El valor por defecto. Este valor es importante porque le dice a Achaean el tipo de variable que es. Puede ser una cadena, lista de números, número o booleano.
- "limits": (opcional): Rango válido para valores numéricos. Ejemplo: [1, 65535].
