Instrucciones de uso de las plantillas

Archivos creados:
- layout.templ: plantilla base (define "layout"). Incluye Pico.css y un placeholder para `content`.
- entity_list.templ: ejemplo de componente para listar `Peliculas` (define "pelicula_list").
- entity_form.templ: ejemplo de formulario para crear/editar `Pelicula` (define "pelicula_form").

Cómo usar desde Go (ejemplo básico):

1) Parsear todas las plantillas:

   tmpl := template.Must(template.ParseGlob("views/*.templ"))

2) Para renderizar una página compuesta por la layout y un contenido, crear una estructura de datos y ejecutar:

   data := map[string]interface{}{
     "Title": "Listado de películas",
     "Peliculas": peliculasSlice, // []Pelicula
   }
   tmpl.ExecuteTemplate(w, "layout", data)

Dentro de `layout.templ` se invoca `{{template "content" .}}`, y las plantillas que definan `content` o `pelicula_list` pueden ser llamadas desde código o incluidas en un wrapper que defina `content`.

Notas:
- Estas plantillas son ejemplos. Si tus entidades se llaman distinto o tienen campos diferentes, adapta los nombres en `entity_list.templ` y `entity_form.templ`.
- Si quieres que el servidor use templates en vez de servir `index.html` estático, puedo modificar `main.go` para parsear `views/*.templ` y renderizar la página raíz usando plantillas. ¿Querés que haga ese cambio y cree una rama git para mantener los cambios separados?
