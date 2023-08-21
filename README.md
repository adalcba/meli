# API Rest para Coordinación de Respuesta a Fraudes

Para coordinar acciones de respuesta ante fraudes, es útil contar con información contextual sobre la ubicación detectada en el momento de la compra, búsqueda y pago. Para lograr esto, se ha tomado la decisión de crear una herramienta que, dada una dirección IP, recupere la información asociada.

## Instrucciones

1. El ejercicio consiste en construir una API REST que permita:

   a. Dada una dirección IP, encontrar el país al que pertenece y mostrar:
      - El nombre y código ISO del país.
   
   b. La moneda local y su tasa de cambio actual en dólares o euros.

2. Lista de Bloqueo de IP: Marcar la IP en una lista de bloqueo y no permitir que consulte la información del punto 1.

## Puntos finales de la API

### 1. Obtener Información del País

- **Punto Final:** `/obtener-info-pais`
- **Método:** GET
- **Parámetros:** `direccion_ip` (cadena)
- **Descripción:** Dada una dirección IP, este punto final devuelve el nombre y código ISO del país correspondiente.

### 2. Obtener Información de la Moneda

- **Punto Final:** `/obtener-info-moneda`
- **Método:** GET
- **Parámetros:** `direccion_ip` (cadena)
- **Descripción:** Dada una dirección IP, este punto final devuelve la moneda local y su tasa de cambio actual en dólares o euros.

### 3. Bloquear IP

- **Punto Final:** `/bloquear-ip`
- **Método:** POST
- **Parámetros:** `direccion_ip` (cadena)
- **Descripción:** Marca la dirección IP especificada en una lista de bloqueo, impidiendo que acceda a la información de los puntos finales anteriores.



