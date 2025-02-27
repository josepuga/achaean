# Ejemplo en Bash #1
Para información adicional sobre la creación de plugins, leer [`plugins/LEEME`](../LEEME.md).

## Descripción
- Muestra el id del plugin por stdout.
- Muestra los parámetros por stdout.
- Realiza un progreso incremental de 5 en 5 y lo imprime en un `progress pipe` (conocido como **named pipe** en Linux).
- Muestra un mensaje de error por stderr
- Mensajes adicionales por stdout. 
- Finaliza Imprimiendo "DONE" en el `progress pipe`.