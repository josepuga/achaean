# Ejemplo en Pascal #1
Para información adicional sobre la creación de plugins, leer [`plugins/LEEME`](../LEEME.md).

## Descripción
- Muestra el id del plugin por stdout.
- Muestra los parámetros por stdout.
- Realiza un progreso incremental de 5 en 5 y lo imprime en un `progress pipe` (conocido como **named pipe** en Linux).
- Muestra un mensaje de error por stderr
- Mensajes adicionales por stdout. 
- Finaliza Imprimiendo "DONE" en el `progress pipe`.

## Notas específicas para Pascal
- Es necesario realizar un `Flush` cada vez que se escribe por alguna salida, al menos en FreePascal, he observado que no siempre se muestran los mensajes al momento. Resulta bastante engorroso tener que usar el `Flush` después de cada `Write`. Se me ocurren 1 solución a esto:
  - Crear una función que use ambos comandos (Write + Flush) o encapsularlos en algún método.
- A la hora de escribir el progreso o "DONE" en el `named pipe`, hay que usar un espacio en blanco adicional.
- Es conveniente crear una compilación estática para mayor portabilidad con: `fpc -Xs -XX -k"static"`