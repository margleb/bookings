function Prompt() {
    let toast = function(c) {
        // с - переписывается по умолчанию если не определено
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }
    let success = function(c) {

        const {
            msg = "",
            title = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer
        })
    }
    let error = function(c) {

        const {
            msg = "",
            title = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer
        })
    }

    async function custom(c) {
        const {
            msg = "",
            title = ""
        } = c

        const { value: result } = await Swal.fire({
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            // вызывается перед тем как окно будет показано окно
            willOpen: () => {
                if(c.willOpen !== undefined) {
                    c.willOpen()
                }
            },
            // вызывается перед тем как будет подтверждено окно
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
            // вызывается после того, как окно будет показано на странице
            didOpen: () => {
                if(c.didOpen !== undefined) {
                    c.didOpen()
                }
            },
        })

        // если есть введенные результат
        if(result) {
            // если не была нажата кнопка отмены
            if(result.dismiss !== Swal.DismissReason.cancel) {
                // если значение не пустые
                if(result.value !== "") {
                    // если опраделена была функция callback
                    if(c.callback !== undefined) {
                        c.callback(result)
                    }
                } else {
                    c.callback(false)
                }
            } else {
                c.callback(false)
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom
    }
}