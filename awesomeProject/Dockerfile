FROM ubuntu:latest
LABEL authors="Guillermo Torres"

# Usa la imagen oficial de Golang como imagen base
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /awesomeProyect

# Copia el código fuente de tu aplicación al contenedor
COPY prueba.go .
# Compila la aplicación
RUN go build -o main

# Expón el puerto en el que se ejecutará tu aplicación
EXPOSE 1414

# Ejecuta tu aplicación cuando se inicie el contenedor
CMD ["./main"]

ENTRYPOINT ["top", "-b"]