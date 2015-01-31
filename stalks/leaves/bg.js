projectlist_initbg = function () {
  // TODO: Use something to replace this if IE 5~8 support is needed
  var bgs = document.getElementsByClassName('bg');
  for (var i = 0; i < bgs.length; i++) {
    var e = bgs[i];
    var prnt = e.parentElement;
    e.style.top = prnt.offsetTop + 'px';
    e.style.width = prnt.clientWidth + 'px';
    e.style.height = prnt.clientHeight + 'px';
  }
};
