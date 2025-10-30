const API_USUARIOS = "/api/usuarios";
const formUsuarios = document.getElementById("usuario-form");
const listUsuarios = document.getElementById("usuario-list");

// Obtener usuarios
async function fetchUsuarios() {
  try {
    // detect id in path: /usuarios/3 or /usuario/3
    const path = window.location.pathname || "";
    const maybeId = (path.match(/\/usuarios?\/(\d+)/) || path.match(/\/usuario\/(\d+)/));
    let usuarios = [];
    if (maybeId && maybeId[1]) {
      const id = maybeId[1];
      const res = await fetch(`${API_USUARIOS}/${id}`);
      if (!res.ok) throw new Error('Error al obtener usuario');
      const single = await res.json();
      usuarios = [single];
    } else {
      const res = await fetch(API_USUARIOS);
      if (!res.ok) throw new Error('Error al obtener usuarios');
      usuarios = await res.json();
    }
    listUsuarios.innerHTML = "";

    usuarios.forEach((u) => {
      const item = document.createElement("li");
      item.innerHTML = `
        <small style="color: #bfff00; font-weight:700">ID: ${u.idu}</small><br>
        <strong>${escapeHtml(u.nomusu)}</strong> (${escapeHtml(u.email)})<br>
        Fecha de nacimiento: ${u.fechanac}<br>
        <button class="edit-btn" data-id="${u.idu}">Editar</button>
        <button class="delete-user" data-id="${u.idu}">Eliminar</button>
      `;
      listUsuarios.appendChild(item);
    });
  } catch (err) {
    console.error('Error al obtener usuarios:', err);
    listUsuarios.innerHTML = '<li>El usuario con ese id no existe.</li>';
  }
}

// Crear nuevo usuario
formUsuarios.addEventListener("submit", async (e) => {
  e.preventDefault();
  const fechaInput = document.getElementById("fechanac").value; // YYYY-MM-DD
  const fechaISO = new Date(fechaInput).toISOString();

  const nuevoUsuario = {
    nomusu: document.getElementById("nomusu").value,
    contrasenia: document.getElementById("contrasenia").value,
    email: document.getElementById("email").value,
    fechanac: fechaISO
  };

  try {
    const res = await fetch(API_USUARIOS, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevoUsuario),
    });
    if (!res.ok) throw new Error('Error al crear usuario');
    formUsuarios.reset();
    fetchUsuarios();
  } catch (err) {
    console.error('Error al crear usuario:', err);
    alert('Error al crear usuario. Revisa la consola.');
  }
});

// Eliminar y editar (delegación)
listUsuarios.addEventListener('click', async (e) => {
  if (e.target.classList.contains('delete-user')) {
    const id = e.target.getAttribute('data-id');
    if (!id) return;
    if (!confirm('¿Eliminar este usuario?')) return;
    try {
      const res = await fetch(`${API_USUARIOS}/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error('Error al eliminar usuario');
      fetchUsuarios();
    } catch (err) {
      console.error('Error al eliminar usuario:', err);
      alert('Error al eliminar usuario. Revisa la consola.');
    }
    return;
  }

  if (e.target.classList.contains('edit-btn')) {
    const id = e.target.getAttribute('data-id');
    const li = e.target.closest('li');
    if (!id || !li) return;

    // Cargar usuario actual
    let usuario;
    try {
      const res = await fetch(`${API_USUARIOS}/${id}`);
      if (!res.ok) { alert('No se pudo cargar usuario'); return; }
      usuario = await res.json();
    } catch (err) {
      console.error('Error al obtener usuario:', err);
      alert('Error de red al cargar usuario.');
      return;
    }

    const form = document.createElement('form');
    form.className = 'inline-edit-form';
    form.innerHTML = `
      <input name="nomusu" placeholder="Nombre" value="${escapeHtml(usuario.nomusu)}" required>
      <input name="email" placeholder="Email" value="${escapeHtml(usuario.email)}" required>
      <input name="fechanac" type="date" value="${(usuario.fechanac||'').slice(0,10)}">
      <div class="edit-buttons">
        <button type="submit" class="save-btn">Guardar</button>
        <button type="button" class="cancel-btn">Cancelar</button>
      </div>
    `;

    const original = li.innerHTML;
    li.innerHTML = '';
    li.appendChild(form);

    form.querySelector('.cancel-btn').addEventListener('click', (ev) => { ev.preventDefault(); li.innerHTML = original; });

    form.addEventListener('submit', async (ev) => {
      ev.preventDefault();
      const fd = new FormData(form);
      const payload = {
        nomusu: fd.get('nomusu') || usuario.nomusu,
        email: fd.get('email') || usuario.email,
        fechanac: new Date(fd.get('fechanac')).toISOString() || usuario.fechanac
      };

      try {
        const res = await fetch(`${API_USUARIOS}/${id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
        if (!res.ok) {
          const text = await res.text().catch(() => null);
          console.error('Error al actualizar usuario:', res.status, text);
          alert('Error al actualizar usuario. Revisa la consola.');
          return;
        }
        fetchUsuarios();
      } catch (err) {
        console.error('Error de red al actualizar usuario:', err);
        alert('Error de red al actualizar usuario.');
      }
    });
  }
});

function escapeHtml(str) {
  if (str == null) return '';
  return String(str).replace(/&/g, '&amp;').replace(/"/g,'&quot;').replace(/'/g,'&#39;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

window.addEventListener('DOMContentLoaded', () => { fetchUsuarios(); });
