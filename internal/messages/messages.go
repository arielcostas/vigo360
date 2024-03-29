package messages

type ErrorMessage string

var ErrorDatos ErrorMessage = "Hubo un error recuperando los datos."
var ErrorNoResultados ErrorMessage = "No se ha encontrado ningún resultado."
var ErrorPaginaNoEncontrada ErrorMessage = "No se ha encontrado ningún resultado."
var ErrorRender ErrorMessage = "Hubo un error mostrando la página solicitada. Inténtelo de nuevo más tarde."
var ErrorFormulario ErrorMessage = "Hubo un error recuperando los datos enviados."
var ErrorValidacion ErrorMessage = "Alguno de los datos del formulario no es válido"
var ErrorSinPermiso ErrorMessage = "No tienes permiso para ver esta página o realizar esta opción."
var ErrorSinAutenticar ErrorMessage = "Debes autenticarte antes de acceder a esta página"

var ErrorIdInvalido ErrorMessage = "El id no es válido. Asegúrate de que solo contiene caracteres alfanuméricos, guiones y guiones bajos."
var ErrorIdDuplicado ErrorMessage = "El id ya está en uso."
var ErrorFatal ErrorMessage = "Ha ocurrido un error fatal. Inténtelo de nuevo más tarde."
