document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll('a[data-confirm]').forEach(function (link) {
    link.addEventListener("click", function (e) {
      e.preventDefault();

      if (!confirm(link.dataset.confirm)) {
        e.stopImmediatePropagation();
        return false;
      }
    });
  });

  document.querySelectorAll('a[data-method="post"]').forEach(function (link) {
    console.log(link);
    link.addEventListener("click", function (e) {
      e.preventDefault();
      const form = document.createElement("form");
      form.method = "POST";
      form.action = link.href;
      document.body.appendChild(form);
      form.submit();
    });
  });
});
