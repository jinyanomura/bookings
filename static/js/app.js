let attention = Prompt();

(function() {
    'use strict';
     window.addEventListener('load', function() {
        // fetch all the forms we want to apply custom bootstrap validations to
        let forms = document.getElementsByClassName('needs-validation');
        // loop over them and prevent submission
        Array.prototype.filter.call(forms, function(form) {
            form.addEventListener('submit', function(event) {
                if (form.checkValidity() === false) {
                  event.preventDefault();
                  event.stopPropagation();
                }
                form.classList.add('was-validated');
            }, false);
        });
     }, false);
})();

function notify(msg, msgType) {
    notie.alert({
      type: msgType,
      text: msg,
    })
}

function notifyModal(title, text, icon, confirmationButtonText) {
    Swal.fire({
       title: title,
       html: text,
       icon, icon,
       confirmationButtonText: confirmationButtonText,
    })
}

function notifyAdmin(target) {
  let e = new Event('click');
  document.getElementById(target).dispatchEvent(e);
}

function Prompt() {
    let toast = function(c) {
      const {
        msg = '',
        icon = 'success',
        position = 'top-end',
      } = c;

      const Toast = Swal.mixin({
        toast: true,
        title: msg,
        position: position,
        icon: icon,
        showConfirmButton: false,
        tiemr: 3000,
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
        msg = '',
        title = '',
        footer = '',
      } = c;

      Swal.fire({
        icon: 'success',
        title: title,
        text: msg,
        footer: footer,
      })
    }

    let error = function(c) {
      const {
        msg = '',
        title = '',
        footer = '',
      } = c;

      Swal.fire({
        icon: 'error',
        title: title,
        text: msg,
        footer: footer,
      })
    }

    async function custom(c) {
      const {
        icon = '',
        msg = '',
        title = '',
        showCancelButton = true,
        showConfirmButton = true,
      } = c;

      const {value: result} = await Swal.fire({
        icon: icon,
        title: title,
        html: msg,
        backdrop: false,
        focusConfirm: false,
        showCancelButton: showCancelButton,
        showConfirmButton: showConfirmButton,
        willOpen: () => {
          if (c.willOpen !== undefined) {
            c.willOpen();
          }
        },
        didOpen: () => {
          if (c.didOpen !== undefined) {
            c.didOpen();
          }
        },
      })

      if (result) {
        if (result.dismiss !== Swal.DismissReason.cancel) {
           if (result.value !== "") {
              if (c.callBack !== undefined) {
                 c.callBack(result);
              }
           } else {
              c.callBack(false);
           }
        } else {
           c.callBack(false);
        }
      }
    }

    return {
          toast: toast,
          success: success,
          error: error,
          custom: custom,
      }
}

function PromptAvail(roomID, csrfToken) {
    document.getElementById("check-availability-button").addEventListener('click', function() {
        let html = `
            <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
            <div class="form-row">
                <div class="col">
                    <div class="row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="start_date" id="start" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="end_date" id="end" placeholder="Departure">
                        </div>
                    </div>
                </div>
            </div>
            </form>
        `;
        attention.custom({
            title: 'Choose your dates',
            msg: html,
            willOpen: () => {
                const el = document.getElementById("reservation-dates-modal");
                const rp = new DateRangePicker(el, {
                    format: 'yyyy-mm-dd',
                    midDate: new Date(),
                    showOnFocus: true,
                    orientation: 'top left'
                })
            },
            didOpen: () => {
                document.getElementById("start").removeAttribute("disabled");
                document.getElementById("end").removeAttribute("disabled");
            },
            callBack: result => {
                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", csrfToken);
                formData.append("room_id", roomID);
    
                fetch('/search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then(res => res.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: "success",
                                msg: "<p>Room is Available!</p>"
                                    + "<p><a href='/book-room?id="
                                    + data.room_id
                                    + "&s="
                                    + data.start_date
                                    + "&e="
                                    + data.end_date
                                    + "'class='btn btn-primary'>Book Now</a></p>",
                                showConfirmButton: false,
                                showCancelButton: false,
                            })
                        } else {
                            attention.error({msg: "Not Available"})
                        }
                    })
            }
        });
    })
}