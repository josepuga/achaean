# Ejemplo en C++ #1
Para información adicional sobre la creación de plugins, leer [`plugins/LEEME`](../LEEME.md).

## Descripción
- Muestra el id del plugin por stdout.
- Muestra los parámetros por stdout.
- Realiza un progreso incremental de 5 en 5 y lo imprime en un `progress pipe` (conocido como **named pipe** en Linux).
- Muestra un mensaje de error por stderr
- Mensajes adicionales por stdout. 
- Finaliza Imprimiendo "DONE" en el `progress pipe`.

## Notas específicas para C++
- En principio, usando **stdout** y **stderr** con `std::cout` + `std::endl` parece que se envían los textos correctamente haciendo implícitamente un `flush` al buffer. Otros métodos de mostrar texto, es posible que requieran forzar dicho `flush`.
- Es conveniente crear una compilación estática para mayor portabilidad con: `g++ -static -static-libgcc -static-libstdc++`