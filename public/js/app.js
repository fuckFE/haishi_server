$(() => {
  $.ajaxSetup({
    complete(res, status) {
      if (res.status === 403) {
        location.href = "/login.html"
      }
    }
  })

  init()
  function init(params) {
    $.ajax({
      url: '/api/users',
      success(res) {
        $('#login-user-name').html(res.user)
      }
    })

    $('#menu').delegate('.remove', 'click', e => {
      e.preventDefault()
        if (!window.confirm(`确定删除?`)) {
        return
      }

      const id = $(e.target).data('id')
      $.ajax({
        url: `/api/types/${id}`,
        method: 'DELETE',
        success(res) {
          if (res === 'success') {
            location.reload()
          }
        }
      })

    })

    getTypes(res => {
      const query = parseQuery()
      if (query.isGarbate) {
        getGarbate()
      } else if (query.id) {
        getBooks(query.id)
          $(`.type-list a.item[data-id=${query.id}]`)
          .addClass('active')
      } else if (res.length > 0){
        getBooks(res[0].id)
          $('.type-list a.item').first()
          .addClass('active')
      }
    })

    addType(1)
    addType(2)
    upload()
    updateUpload()
    save()
    managerBook()
    $('#logout').on('click', () => {
      $.ajax({
        url: '/api/users/logout',
        success() {
          location.href = "/login.html"
        }
      })
    })
  }

  function parseQuery() {
    const qstr = location.search
    const query = {};
    const a = (qstr[0] === '?' ? qstr.substr(1) : qstr).split('&');
    for (let i = 0; i < a.length; i++) {
      const b = a[i].split('=');
      query[decodeURIComponent(b[0])] = decodeURIComponent(b[1] || '');
    }
    return query;
  }

  function addType(num) {
    $(`#type${num}-save`).on('click', () => {
      const name = $(`#type${num}-name`).val()
      if (!name || name.length === 0) {
        return alert('主题不能为空')
      }

      $.ajax({
        url: '/api/types',
        method: 'POST',
        data: JSON.stringify({
          category: num,
          name,
        }),
        success: (res) => {
          $(`#type${num}-list a:last-child`)
            .before(`<a href="/?id=${res.id}" data-id=${res.id} class="list-group-item item">${res.name}<span data-id=${res.id} class="glyphicon glyphicon-remove remove" aria-hidden="true"></span></a>`)
          $(`#type${num}-select`)
            .append(`<option value="${res.id}">${res.name}</option>`)
          $(`#addType${num}`).modal('hide')
          $(`#type${num}-name`).val("")
        }
      })
    })
  }

  function getGarbate() {
    $.ajax({
      url: "/api/books/grabate",
      method: "POST",
      success: res => {
        $('#book-body').empty()
        res.forEach(b => {
          $('#book-body')
              .append(`
              <tr>
                <td>${b.id}</td>
                <td><a href="/detail.html?id=${b.id}" target="_blank">${b.filename}</a></td>
                <td>${moment(b.last).format('YYYY-MM-DD HH:mm')}</td>
                <td>
                  <div class="btn-group manager-group" role="group" aria-label="...">
                    <button data-id=${b.id} class="update btn btn-xs btn-info" type="button" data-toggle="modal" data-target="#update">更新</button>
                    <button data-id=${b.id} class="garbate btn-xs btn ${b.isGarbate? 'btn-success' : 'btn-danger' }" data-isGarbate=${b.isGarbate} type="button">${b.isGarbate ? '启用' : '弃用'}</button>
                    <button data-id=${b.id} class="delete btn-xs btn btn-danger" type="button" class="btn btn-danger">删除</button>
                  </div>
                </td>
              </tr>
              `)
        })

        $('#grabate')
          .addClass('active')
      }
    })
  }

  function getTypes(cb) {
    $.ajax({
      url: '/api/types',
      success: (res) => {
        res.forEach(t => {
          $(`#type${t.category}-list a:last-child`)
            .before(`<a href="/?id=${t.id}" data-id=${t.id} class="list-group-item item">${t.name}<span data-id=${t.id} class="glyphicon glyphicon-remove remove" aria-hidden="true"></span></a>`)
          $(`#type${t.category}-select`)
            .append(`<option value="${t.id}">${t.name}</option>`)
        })

        cb(res)
      }
    })
  }

  function upload() {
    $('#fileupload').fileupload({
      url: '/api/upload',
      dataType: 'json',
      done: function (e, data) {
        const result = data.result
        $('#filename').val(result.filename)
        $('#fileId').val(result.id)

        $.each(data.result.files, function (index, file) {
          $('<p/>').text(file.name).appendTo('#files');
        });
      },
      progressall: function (e, data) {
        var progress = parseInt(data.loaded / data.total * 100, 10);
        $('#progress .progress-bar').css(
          'width',
          progress + '%'
        );
      }
    }).prop('disabled', !$.support.fileInput)
      .parent().addClass($.support.fileInput ? undefined : 'disabled');
  }

  function save() {
    $('#save').on('click', () => {
      const type1 = $('#type1-select').val() || 0
      const type2 = $('#type2-select').val() || 0
      const filename = $('#filename').val()
      const fileId = $('#fileId').val()

      if (!filename || !fileId) {
        return alert("请先上传文件")
      }

      const data = JSON.stringify({
        types: [+type1, +type2],
        filename,
        fileId: +fileId,
      })

      $.ajax({
        url: '/api/books',
        method: 'POST',
        data,
        success: (res) => {
          window.location.reload()
        }
      })
    })
  }

  function getBooks(id) {
    $.ajax({
      url: `/api/books?type=${id}&filterGarbate=1`,
      success: (res) => {
        $('#book-body').empty()
        res.forEach(b => {
          $('#book-body')
              .append(`
              <tr>
                <td>${b.id}</td>
                <td><a href="/detail.html?id=${b.id}" target="_blank">${b.filename}</a></td>
                <td>${moment(b.last).format('YYYY-MM-DD HH:mm')}</td>
                <td>
                  <div class="btn-group manager-group" role="group" aria-label="...">
                    <button data-id=${b.id} class="update btn btn-xs btn-info" type="button" data-toggle="modal" data-target="#update">更新</button>
                    <button data-id=${b.id} class="garbate btn-xs btn ${b.isGarbate? 'btn-success' : 'btn-danger' }" data-isGarbate=${b.isGarbate} type="button">${b.isGarbate ? '启用' : '弃用'}</button>
                    <button data-id=${b.id} class="delete btn-xs btn btn-danger" type="button" class="btn btn-danger">删除</button>
                  </div>
                </td>
              </tr>
              `)
        })
      }
    })
  }

  function managerBook() {
    $('#book-body').delegate('.update', 'click', e => {
      const id = $(e.target).data('id')
      $('#update-id').val(id)
    })
     $('#update-book').on('click', e => {
      const fileID = $('#update-fileId').val()
      const id = $('#update-id').val()
      $('#update-fileId').val("")
      $('#update-id').val("")

      if (!fileID || !id) {
        return alert("请上传文件")
      }

      $.ajax({
        url: `/api/books/${id}/payload`,
        method: 'PUT',
        data: { fileID },
        success: (res) => {
          window.location.reload()
        }
      })
    })

    $('#book-body').delegate('.garbate', 'click', e => {
      if (!confirm("确定废弃吗")) {
        return
      }
      const id = $(e.target).data('id')
      $.ajax({
        url: `/api/books/${id}/garbate`,
        method: 'PUT',
        success: () => {
          location.reload()
          // if ($(e.target).text() == '弃用') {
          //   $(e.target).text("启用")
          //   $(e.target).removeClass("btn-danger")
          //   $(e.target).addClass("btn-success")
          // } else {
          //   $(e.target).text("弃用")
          //   $(e.target).removeClass("btn-success")
          //   $(e.target).addClass("btn-danger")
          // }
        }
      })
    })

    $('#book-body').delegate('.delete', 'click', e => {
      if (!confirm("确定更新吗")) {
        return
      }
      const id = $(e.target).data('id')
      $.ajax({
        url: `/api/books/${id}`,
        method: 'DELETE',
        success: () => {
          window.location.reload()
        }
      })
    })
  }

  function updateUpload(params) {
    $('#update-fileupload').fileupload({
      url: '/api/upload',
      dataType: 'json',
      done: function (e, data) {
        const result = data.result
        $('#update-fileId').val(result.id)

        $.each(data.result.files, function (index, file) {
          $('<p/>').text(file.name).appendTo('#files');
        });
      },
      progressall: function (e, data) {
        var progress = parseInt(data.loaded / data.total * 100, 10);
        $('#progress .progress-bar').css(
          'width',
          progress + '%'
        );
      }
    }).prop('disabled', !$.support.fileInput)
      .parent().addClass($.support.fileInput ? undefined : 'disabled');
  }
})