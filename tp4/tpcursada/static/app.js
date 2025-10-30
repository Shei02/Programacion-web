// ==================== PELICULAS =====================================================================
const API_PELICULAS = "/api/peliculas";
const formPeliculas = document.getElementById("pelicula-form"); //busca en el html el formulario con este ID
const listPeliculas = document.getElementById("pelicula-list"); //busca el contenedor de las peliculas

//Obtener todas las películas GET
async function fetchPeliculas() {
  try {
    const response = await fetch(API_PELICULAS); //peticion GET
    if (!response.ok) throw new Error("Error al obtener películas");

    const peliculas = await response.json(); //convierte cada pelicula en un objeto y las almacena en un array de objs

    // Limpiar la lista antes de volver a renderizar
    listPeliculas.innerHTML = "";

    // Crear dinámicamente cada película en la lista
    peliculas.forEach((peli) => {
      const item = document.createElement("li"); //crea un <li> (elemento de la lista) dinamicamente 
      item.innerHTML = `
        <strong>${peli.titulo}</strong> (${peli.anioestr})<br>
        Director: ${peli.director}<br>
        Duración: ${peli.duracion} min<br>
        Actores ${peli.actores}<br>
        Sinopsis: ${peli.sinopsis}<br>
        Edad mínima: ${peli.edadmin}<br>
        <button class="edit-btn" data-id="${peli.idp}">Editar</button>
        <button class="delete-btn" data-id="${peli.idp}">Eliminar</button>
      `;
      listPeliculas.appendChild(item); //lo agrega al <ul>
    });
  } catch (error) {
    console.error(error);
  }
}
 
//Crear nueva película POST
formPeliculas.addEventListener("submit", async (e) => { //escucha el envio del formulario, submit
  e.preventDefault(); //evita que el formulario recargue la pag 

  // Tomar los valores del formulario
  const nuevaPelicula = {
    titulo: document.getElementById("titulo").value,
    duracion: parseInt(document.getElementById("duracion").value),
    director: document.getElementById("director").value,
    actores: document.getElementById("actores").value,
    edadmin: parseInt(document.getElementById("edadmin").value),
    sinopsis: document.getElementById("sinopsis").value,
    anioestr: parseInt(document.getElementById("anioestr").value),
  };
   
    try {
    await fetch(API_PELICULAS, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevaPelicula),
    });
    formPeliculas.reset();
    fetchPeliculas();
  } catch (err) {
    console.error("Error al crear pelicula:", err);
  }
});

//Modificar pelicula
// --- EDIT INLINE para Películas (mejor que prompt) ---
listPeliculas.addEventListener("click", async (e) => {
  if (!e.target.classList.contains("edit-btn")) return;

  const id = e.target.getAttribute("data-id");
  const li = e.target.closest("li");
  if (!id || !li) return;

  // 1) Obtener la película actual desde el backend
  let peliActual;
  try {
    const resGet = await fetch(`${API_PELICULAS}/${id}`);
    if (!resGet.ok) {
      alert("No se pudo cargar la película para editar.");
      return;
    }
    peliActual = await resGet.json();
  } catch (err) {
    console.error("Error al obtener pelicula:", err);
    alert("Error de red al cargar la película.");
    return;
  }

  // 2) Crear el formulario inline (inputs prellenados)
  const form = document.createElement("form");
  form.className = "inline-edit-form";
  form.innerHTML = `
    <input name="titulo" placeholder="Título" value="${escapeHtml(peliActual.titulo)}" required>
    <input name="duracion" type="number" placeholder="Duración (min)" value="${peliActual.duracion}" required>
    <input name="director" placeholder="Director" value="${escapeHtml(peliActual.director)}">
    <input name="actores" placeholder="Actores" value="${escapeHtml(peliActual.actores)}">
    <input name="edadmin" type="number" placeholder="Edad mínima" value="${peliActual.edadmin}">
    <input name="sinopsis" placeholder="Sinopsis" value="${escapeHtml(peliActual.sinopsis)}">
    <input name="anioestr" type="number" placeholder="Año" value="${peliActual.anioestr}">
    <div class="edit-buttons">
      <button type="submit" class="save-btn">Guardar</button>
      <button type="button" class="cancel-btn">Cancelar</button>
    </div>
  `;

  // Guardar el contenido original para restaurar si cancela
  const originalHTML = li.innerHTML;
  li.innerHTML = "";
  li.appendChild(form);

  // 3) Handler Cancelar
  form.querySelector(".cancel-btn").addEventListener("click", (ev) => {
    ev.preventDefault();
    li.innerHTML = originalHTML;
  });

  // 4) Handler Guardar (merge y PUT)
  form.addEventListener("submit", async (ev) => {
    ev.preventDefault();

    // Construir objeto completo basado en peliActual y los nuevos valores (merge)
    const formData = new FormData(form);
    const payload = {
      titulo: formData.get("titulo") || peliActual.titulo,
      duracion: parseInt(formData.get("duracion")) || parseInt(peliActual.duracion) || 0,
      director: formData.get("director") || peliActual.director,
      actores: formData.get("actores") || peliActual.actores,
      edadmin: parseInt(formData.get("edadmin")) || parseInt(peliActual.edadmin) || 0,
      sinopsis: formData.get("sinopsis") || peliActual.sinopsis,
      anioestr: parseInt(formData.get("anioestr")) || parseInt(peliActual.anioestr) || 0
    };

    try {
      const res = await fetch(`${API_PELICULAS}/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        const text = await res.text().catch(() => null);
        console.error("Error al actualizar:", res.status, text);
        alert("Error al actualizar la película. Revisa la consola.");
        return;
      }

      // Si se actualizó correctamente, refrescar la lista
      await fetchPeliculas();
    } catch (err) {
      console.error("Error de red al actualizar:", err);
      alert("Error de red al actualizar la película.");
    }
  });
});


// --- Helper: escapar texto para insertar en value (evita romper HTML) ---
function escapeHtml(str) {
  if (str == null) return "";
  return String(str)
    .replace(/&/g, "&amp;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}



//Eliminar película (delegación de eventos)
listPeliculas.addEventListener("click", async (e) => { //escucha los clicks dentro del <ul> (lista desordenada)
  if (e.target.classList.contains("delete-btn")) { //verifica que se haya clickeado un boton de delete
    const id = e.target.getAttribute("data-id"); //obtiene el ID de la peli
    if (!id) return;

    const confirmar = confirm("¿Seguro que querés eliminar esta película?");
    if (!confirmar) return;

    try {
      const response = await fetch(`${API_PELICULAS}/${id}`, { //realiza la peticion DELETE
        method: "DELETE",
      });

      if (!response.ok) throw new Error("Error al eliminar película");

      fetchPeliculas(); // refrescar la lista para actualizarla
    } catch (error) {
      console.error(error);
    }
  }
});

//Cargar automáticamente la lista al iniciar la página
fetchPeliculas();

// ==================== USUARIOS ================================================================
const API_USUARIOS = "/api/usuarios";
const formUsuarios = document.getElementById("usuario-form");
const listUsuarios = document.getElementById("usuario-list");

// Obtener usuarios
async function fetchUsuarios() {
  try {
    const res = await fetch(API_USUARIOS);
    const usuarios = await res.json();
    listUsuarios.innerHTML = "";

    usuarios.forEach((u) => {
      const item = document.createElement("li");
      item.innerHTML = `
        <strong>${u.nomusu}</strong> (${u.email})<br>
        Fecha de nacimiento: ${u.fechanac}<br>
        <button class="edit-btn" data-id="${u.idu}">Editar</button>
        <button class="delete-user" data-id="${u.idu}">Eliminar</button>
      `;
      listUsuarios.appendChild(item);
    });
  } catch (err) {
    console.error("Error al obtener usuarios:", err);
  }
}

// Crear nuevo usuario
formUsuarios.addEventListener("submit", async (e) => {
  e.preventDefault();

    const fechaInput = document.getElementById("fechanac").value; // YYYY-MM-DD formato recibido
    const fechaISO = new Date(fechaInput).toISOString(); // YYYY-MM-DDT00:00:00.000Z lo convierte a time.time

  const nuevoUsuario = {
    nomusu: document.getElementById("nomusu").value,
    contrasenia: document.getElementById("contrasenia").value,
    email: document.getElementById("email").value,
    fechanac: fechaISO
  };
    try {
    await fetch(API_USUARIOS, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevoUsuario),
    });
    formUsuarios.reset();
    fetchUsuarios();
  } catch (err) {
    console.error("Error al crear usuario:", err);
  }
});

// Eliminar usuario
listUsuarios.addEventListener("click", async (e) => {
  if (e.target.classList.contains("delete-user")) {
    const id = e.target.getAttribute("data-id");
    if (confirm("¿Eliminar este usuario?")) {
      await fetch(`${API_USUARIOS}/${id}`, { method: "DELETE" });
      fetchUsuarios();
    }
  }
});

// ==================== MIRAS ====================
const formMiras = document.getElementById("mira-form");
const listMiras = document.getElementById("mira-list");
const API_MIRAS = "/api/miras";

// Obtener miras
async function fetchMiras() {
  try {
    const res = await fetch(API_MIRAS);
    const miras = await res.json();
    listMiras.innerHTML = "";

    miras.forEach((m) => {
      const item = document.createElement("li");
      item.innerHTML = `
        Usuario ID: ${m.idu} - Película ID: ${m.idp} <br>
        Gusto: ${m.gustoono} <br>
        Calificación: ${m.calif} <br>
        <button class="edit-btn" data-id="${m.idu}-${m.idp}">Editar</button>
        <button class="delete-mira" data-id="${m.idu}-${m.idp}">Eliminar</button>
      `;
      listMiras.appendChild(item);
    });
  } catch (err) {
    console.error("Error al obtener miras:", err);
  }
}

// Crear nueva mira
formMiras.addEventListener("submit", async (e) => {
  e.preventDefault();

  const nuevaMira = {
    idu: parseInt(document.getElementById("idusuario").value),
    idp: parseInt(document.getElementById("idpelicula").value),
    gustoono: document.getElementById("gustoOno").value,
    calif: parseInt(document.getElementById("calif").value),
  };

  try {
    await fetch(API_MIRAS, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevaMira),
    });
    formMiras.reset();
    fetchMiras();
  } catch (err) {
    console.error("Error al crear mira:", err);
  }
});

// Eliminar mira
listMiras.addEventListener("click", async (e) => {
  if (e.target.classList.contains("delete-mira")) {
    const ids = e.target.getAttribute("data-id").split("-"); // [idu, idp]
    const idu = parseInt(ids[0]);
    const idp = parseInt(ids[1]);

    if (confirm("¿Eliminar esta mira?")) {
      await fetch(`${API_MIRAS}?idu=${idu}&idp=${idp}`, { method: "DELETE" });
      fetchMiras();
    }
  }
});


// ==================== CARGA INICIAL ====================
window.addEventListener("DOMContentLoaded", () => {
  fetchPeliculas();
  fetchUsuarios();
  fetchMiras();
});
