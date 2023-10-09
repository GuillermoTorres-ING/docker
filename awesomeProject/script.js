document.addEventListener("DOMContentLoaded", function () {
    cargarInformacion();

    document.getElementById("insertar-form").addEventListener("submit", function (e) {
        e.preventDefault();
        insertarInformacion();
    });
});

function cargarInformacion() {
    // Realizar una llamada AJAX para obtener la lista de información desde el servidor Go
    fetch("/consultarInformacion")
        .then(response => response.json())
        .then(data => {
            const infoList = document.getElementById("info-list");
            infoList.innerHTML = ""; // Limpiar la lista antes de agregar elementos

            data.forEach(info => {
                const listItem = document.createElement("li");
                listItem.textContent = `ID: ${info.ID}, Usuario: ${info.Usuario}, Contraseña: ${info.Pass}`;
                infoList.appendChild(listItem);
            });
        })
        .catch(error => console.error("Error al cargar información: " + error));
}

function insertarInformacion() {
    // Obtener datos del formulario
    const usuario = document.getElementById("usuario").value;
    const pass = document.getElementById("pass").value;

    // Realizar una llamada AJAX para insertar información en el servidor Go
    fetch("/insertarInformacion", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ usuario, pass }),
    })
        .then(response => response.json())
        .then(data => {
            console.log("Información insertada con éxito.");
            cargarInformacion(); // Volver a cargar la lista después de la inserción
        })
        .catch(error => console.error("Error al insertar información: " + error));
}
