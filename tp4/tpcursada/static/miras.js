const API_MIRAS = "/api/miras";
const formMiras = document.getElementById("mira-form");
const listMiras = document.getElementById("mira-list");

// Obtener miras
async function fetchMiras() {
  try {
    // If the page URL contains an ID like /miras/2, fetch miras for that user
    const path = window.location.pathname || "";
    const maybeId = (path.match(/\/miras?\/(\d+)/) || path.match(/\/mira\/(\d+)/));
    let resMiras;
    if (maybeId && maybeId[1]) {
      resMiras = await fetch(`${API_MIRAS}/${maybeId[1]}`);
    } else {
      resMiras = await fetch(API_MIRAS);
    }

    const [resUsuarios, resPeliculas] = await Promise.all([fetch('/api/usuarios'), fetch('/api/peliculas')]);

    if (!resMiras.ok) throw new Error('Error al obtener miras');
    if (!resUsuarios.ok) throw new Error('Error al obtener usuarios');
    if (!resPeliculas.ok) throw new Error('Error al obtener peliculas');

    const [miras, usuarios, peliculas] = await Promise.all([resMiras.json(), resUsuarios.json(), resPeliculas.json()]);
    listMiras.innerHTML = '';

    // For each mira, find the corresponding user and pelicula by searching the arrays
    miras.forEach((m) => {
      const item = document.createElement('li');
      const usuarioObj = usuarios.find(u => Number(u.idu) === Number(m.idu));
      const peliculaObj = peliculas.find(p => Number(p.idp) === Number(m.idp));
      const uname = usuarioObj ? usuarioObj.nomusu : 'Desconocido';
      const ptitle = peliculaObj ? peliculaObj.titulo : 'Desconocida';
      item.innerHTML = `
        <small style="color: #bfff00; font-weight:700">Mira: ${escapeHtml(uname)} — ${escapeHtml(ptitle)}</small><br>
        <small style="color: rgba(255,255,255,0.7)">(usu ID ${m.idu} · peli ID ${m.idp})</small><br>
        Gusto: ${escapeHtml(m.gustoono)} <br>
        Calificación: ${m.calif} <br>
        <button class="edit-btn" data-id="${m.idu}-${m.idp}">Editar</button>
        <button class="delete-mira" data-id="${m.idu}-${m.idp}">Eliminar</button>
      `;
      listMiras.appendChild(item);
    });
  } catch (err) {
    console.error('Error al obtener miras:', err);
    listMiras.innerHTML = '<li>Error al cargar las miras.</li>';
  }
}

// Crear nueva mira
formMiras.addEventListener('submit', async (e) => {
  e.preventDefault();
  const nuevaMira = {
    idu: parseInt(document.getElementById('idusuario').value),
    idp: parseInt(document.getElementById('idpelicula').value),
    gustoono: document.getElementById('gustoono').value,
    calif: parseInt(document.getElementById('calif').value),
  };

  try {
    const res = await fetch(API_MIRAS, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(nuevaMira)
    });
    if (!res.ok) throw new Error('Error al crear mira');
    formMiras.reset();
    fetchMiras();
  } catch (err) {
    console.error('Error al crear mira:', err);
    alert('Error al crear mira. Revisa la consola.');
  }
});

// Eliminar mira
listMiras.addEventListener('click', async (e) => {
  if (e.target.classList.contains('delete-mira')) {
    const ids = e.target.getAttribute('data-id').split('-');
    const idu = parseInt(ids[0]);
    const idp = parseInt(ids[1]);
    if (!confirm('¿Eliminar esta mira?')) return;
    try {
      const res = await fetch(`${API_MIRAS}?idu=${idu}&idp=${idp}`, { method: 'DELETE' });
      if (!res.ok) throw new Error('Error al eliminar mira');
      fetchMiras();
    } catch (err) {
      console.error('Error al eliminar mira:', err);
      alert('Error al eliminar mira. Revisa la consola.');
    }
  }

  // Edición básica: abrir prompt para editar calificación y gusto
  if (e.target.classList.contains('edit-btn')) {
    const key = e.target.getAttribute('data-id');
    const parts = key.split('-');
    const idu = parseInt(parts[0]);
    const idp = parseInt(parts[1]);

    // Cargar mira actual
    let mira;
    try {
        const res = await fetch(`${API_MIRAS}?idu=${idu}`);
        if (!res.ok) { alert('No se pudo cargar la(s) mira(s)'); return; }
        const list = await res.json();
        // buscar la mira específica en la lista
        mira = Array.isArray(list) ? list.find(x => parseInt(x.idp) === idp) : null;
        if (!mira) { alert('Mira no encontrada'); return; }
    } catch (err) {
      console.error('Error al obtener mira:', err);
      alert('Error de red al cargar la mira.');
      return;
    }
      // Crear formulario inline para editar gustoono y calif
      const form = document.createElement('form');
      form.className = 'inline-edit-form';
      form.innerHTML = `
        <label>Gusto: <input name="gustoono" value="${escapeHtml(mira.gustoono)}"></label>
        <label>Calificación: <input name="calif" type="number" min="1" max="5" value="${mira.calif}"></label>
        <div class="edit-buttons">
          <button type="submit" class="save-btn">Guardar</button>
          <button type="button" class="cancel-btn">Cancelar</button>
        </div>
      `;

      const li = e.target.closest('li');
      const original = li.innerHTML;
      li.innerHTML = '';
      li.appendChild(form);

      form.querySelector('.cancel-btn').addEventListener('click', (ev) => { ev.preventDefault(); li.innerHTML = original; });

      form.addEventListener('submit', async (ev) => {
        ev.preventDefault();
        const fd = new FormData(form);
        const payload = {
          idp: idp,
          idu: idu,
          gustoono: fd.get('gustoono') || mira.gustoono,
          calif: parseInt(fd.get('calif')) || mira.calif
        };

        try {
          console.log('PUT', `${API_MIRAS}/${idp}/${idu}`, payload);
          const resPut = await fetch(`${API_MIRAS}/${idp}/${idu}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
          });
          if (!resPut.ok) {
            const txt = await resPut.text().catch(() => null);
            console.error('Error al actualizar mira:', resPut.status, txt);
            alert('Error al actualizar mira (ver consola).\n' + (txt || ('HTTP ' + resPut.status)));
            return;
          }
          fetchMiras();
        } catch (err) {
          console.error('Error al actualizar mira:', err);
          alert('Error al actualizar mira (error de red). Revisa la consola.');
        }
      });
  }
});

function escapeHtml(str) {
  if (str == null) return '';
  return String(str).replace(/&/g, '&amp;').replace(/"/g,'&quot;').replace(/'/g,'&#39;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

window.addEventListener('DOMContentLoaded', () => { fetchMiras(); });
