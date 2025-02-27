# Ejemplo en Python #1
Para información adicional sobre la creación de plugins, leer `plugins/LEEME`.

## Descripción
- Muestra el id del plugin por stdout.
- Muestra los parámetros por stdout.
- Realiza un progreso incremental de 5 en 5 y lo imprime en un `progress pipe` (conocido como **named pipe** en Linux).
- Muestra un mensaje de error por stderr
- Mensajes adicionales por stdout. 
- Finaliza Imprimiendo "DONE" en el `progress pipe`.

## Notas específicas para python
- El programa no puede ser llamado con `python your_script.py`, en su lugar, hay que hacerlo ejecutable (`chmod +x`) y poner en el fichero principal en la primera linea el **shebang** `#!/usr/bin/env python` 
- Es necesario forzar un flush con `flush()` del buffer cada vez que se escribe en el `progress pipe`.