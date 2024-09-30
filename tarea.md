# Ejercicio
Una empresa tiene 5 servidores (1,2,3,4,5), en el cual almacena todos sus
documentos en un espacio de nombres llamado DOCUMENTOS. Cada servidor contiene
una parte diferente de la documentación de la empresa.

# Estructura del espacio de nombre:
  - Documentos
    - Contabilidad
      - Informes anuales
      - Facturas
      - Revisiones SII
    - Marketing
      - Campañas publicitarias
      - Estudios de mercado
    - Recursos humanos
      - Nominas trabajadores
      - Contratos
    - Produccion
      - Manuales de usuario
        - Equipos electricos
        - Equipos mecanicos
        - Equipos hidraulicos
        - Equipos neumaticos
      - Diagramas de flujo

# Preguntas
1. Construya el grafo (arbol) asociado.
2. Si un empleado busca el nombre "Informe financiero anual", ¿Cómo se aplica
   DFS en este caso?
    Si asumimos que este archivo está dentro de la carpeta "Informes anuales",
    entonces cuando el algoritmo llegue a esa carpeta revisará todos los archivos
    hasta encontrar el que se busca y lo encontrará si existe.

3. Un empleado por error, copió el mismo archivo "Informe financiero anual" en
   revisiones SII, ¿Qué resultados ofrece DFS?, ¿DFS es capaz de distinguir la
   copia?
    Aquí hay muchos casos que considerar. Si la implementación de DFS solo
    utiliza los nombres de los archivos para determinar si es el archivo correcto,
    entonces ambos archivos serán listados como opciones, tal como lo hacen la
    mayoría de sistemas operativas al buscar. Sin embargo, si se configura de modo
    que solo muestre el primer resultado, entonces solo mostrará el primero que
    encuentre. Pero aún así, esto depende de si además se revisa el contenido
    de los archivos u otros factores que modifican el problema aún más.
