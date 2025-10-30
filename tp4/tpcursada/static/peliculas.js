const API_PELICULAS = "/api/peliculas";
const listPeliculas = document.getElementById("pelicula-list");

// Formulario de creación (igual comportamiento que en index)
const formPeliculas = document.getElementById("pelicula-form");
if (formPeliculas) {
  formPeliculas.addEventListener("submit", async (e) => {
    e.preventDefault();
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
      const res = await fetch(API_PELICULAS, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(nuevaPelicula),
      });
      if (!res.ok) {
        const txt = await res.text().catch(() => null);
        console.error("Error al crear película:", res.status, txt);
        alert("Error al crear película. Revisa la consola.");
        return;
      }
      formPeliculas.reset();
      fetchPeliculas();
    } catch (err) {
      console.error("Error de red al crear película:", err);
      alert("Error de red al crear película.");
    }
  });
}

// Obtener y renderizar películas con botones Editar/Eliminar
async function fetchPeliculas() {
  try {
    // Detect if the URL contains an ID (e.g. /peliculas/3 or /pelicula/3)
    const path = window.location.pathname || "";
    const maybeId = (path.match(/\/peliculas?\/(\d+)/) || path.match(/\/pelicula\/(\d+)/));
    let peliculas = [];
    if (maybeId && maybeId[1]) {
      const id = maybeId[1];
      const response = await fetch(`${API_PELICULAS}/${id}`);
      if (!response.ok) throw new Error("Error al obtener película");
      const single = await response.json();
      peliculas = [single];
    } else {
      const response = await fetch(API_PELICULAS);
      if (!response.ok) throw new Error("Error al obtener películas");
      peliculas = await response.json();
    }

    listPeliculas.innerHTML = "";

    peliculas.forEach((peli) => {
      const item = document.createElement("li");
      item.innerHTML = `
        <small style="color: #bfff00; font-weight:700">ID: ${peli.idp}</small><br>
        <strong>${escapeHtml(peli.titulo)}</strong> (${peli.anioestr})<br>
        Director: ${escapeHtml(peli.director)}<br>
        Duración: ${peli.duracion} min<br>
        Actores: ${escapeHtml(peli.actores)}<br>
        Sinopsis: ${escapeHtml(peli.sinopsis)}<br>
        Edad mínima: ${peli.edadmin}<br>
        <button class="edit-btn" data-id="${peli.idp}">Editar</button>
        <button class="delete-btn" data-id="${peli.idp}">Eliminar</button>
      `;
      listPeliculas.appendChild(item);
    });
  } catch (error) {
    console.error(error);
    listPeliculas.innerHTML = '<li>La pelicula con ese id no existe.</li>';
  }
}

// Delegación para Editar y Eliminar
listPeliculas.addEventListener("click", async (e) => {
  // Eliminar
  if (e.target.classList.contains("delete-btn")) {
    const id = e.target.getAttribute("data-id");
    if (!id) return;
    if (!confirm("¿Eliminar esta película?")) return;
    try {
      const res = await fetch(`${API_PELICULAS}/${id}`, { method: "DELETE" });
      if (!res.ok) {
        const txt = await res.text().catch(() => null);
        console.error("Error al eliminar:", res.status, txt);
        alert("Error al eliminar la película. Revisa la consola.");
        return;
      }
      fetchPeliculas();
    } catch (err) {
      console.error("Error de red al eliminar:", err);
      alert("Error de red al eliminar la película.");
    }
    return;
  }

  // Editar (inline)
  if (e.target.classList.contains("edit-btn")) {
    const id = e.target.getAttribute("data-id");
    const li = e.target.closest("li");
    if (!id || !li) return;

    // Obtener la película actual
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

    // Crear formulario inline
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

    const originalHTML = li.innerHTML;
    li.innerHTML = "";
    li.appendChild(form);

    // Cancelar
    form.querySelector(".cancel-btn").addEventListener("click", (ev) => {
      ev.preventDefault();
      li.innerHTML = originalHTML;
    });

    // Guardar
    form.addEventListener("submit", async (ev) => {
      ev.preventDefault();
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

        await fetchPeliculas();
      } catch (err) {
        console.error("Error de red al actualizar:", err);
        alert("Error de red al actualizar la película.");
      }
    });
  }
});

// Helper: escapar texto
function escapeHtml(str) {
  if (str == null) return "";
  return String(str)
    .replace(/&/g, "&amp;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}

window.addEventListener("DOMContentLoaded", () => {
  fetchPeliculas();
});
