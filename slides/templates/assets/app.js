$(function () {
  var queryParams = new URLSearchParams(window.location.search);

  let pages = $("page");
  let pl = pages.length;
  let currentPage = 0;

  console.log(queryParams);
  if (queryParams.has("page")) {
    i = parseInt(queryParams.get("page"));
    if (i >= 0 && i < pl) {
      console.log("Setting page to " + i);
      currentPage = i;
    }
  }

  changePage = function () {
    console.log("Changing page to " + currentPage);
    $("page").hide();
    $(pages[currentPage]).show();
    queryParams.set("page", currentPage);
    history.pushState(null, null, "?" + queryParams.toString());
  };

  console.log(currentPage);
  changePage();

  $(document).keydown(function (e) {
    switch (e.keyCode) {
      case 39:
        currentPage++;
        if (currentPage >= pl) {
          currentPage = pl;
          break;
        }
        changePage();
        break;
      case 37:
        currentPage--;
        if (currentPage < 0) {
          currentPage = 0;
          break;
        }
        changePage();
        break;
    }
  });
});
