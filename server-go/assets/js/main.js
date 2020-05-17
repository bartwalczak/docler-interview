$(function () {

  console.log("main.js");

  function trimDate(date) {
    return date.substr(0, date.indexOf("T"));
  }

  function renderTasks(data) {
    $("tr.task-item").remove();

    for (let i = 0; i < data.length; i++) {
      let it = data[i];
      let html = `\
      <tr class="task-item" data-id="${it.id}">\
        <td class="title">${it.title}</td>\
        <td class="priority">${it.priority}</td>\
        <td class="due">${trimDate(it.due)}</td>\
      </tr>\
      `
      $(".container .list").append($(html));
    }

  }

  function showError(err) {
    console.log(err);
    $(".message").html(err).slideDown();
  }

  function getAllTasks() {
    $(".message").slideUp().html("");

    fetch("/tasks").then( resp => {
      if (!resp.ok) { throw resp }
      return resp.json()
    })
    .then( data => { renderTasks(data) })
    .catch( err => { showError(err) });
  }

  function getTodayTasks() {
    $(".message").slideUp().html("");

    fetch("/tasks/today").then( resp => {
      if (!resp.ok) { throw resp }
      return resp.json()
    })
    .then( data => { renderTasks(data) })
    .catch( err => { showError(err) });
  }

  getTodayTasks();
});
